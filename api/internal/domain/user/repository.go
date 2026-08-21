package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"flower/api/internal/platform/db"
	"flower/api/internal/types"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(database *sql.DB) *Repository {
	return &Repository{db: database}
}

func (r *Repository) WithTx(ctx context.Context, fn func(*sql.Tx) error) error {
	return db.WithTx(ctx, r.db, fn)
}

func (r *Repository) EmailVerifiedAt(ctx context.Context, userID string) (*time.Time, error) {
	var verified sql.NullTime
	err := r.db.QueryRowContext(ctx, `SELECT email_verified_at FROM users WHERE id = $1`, userID).Scan(&verified)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, types.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("user: email verified: %w", err)
	}
	if !verified.Valid {
		return nil, nil
	}
	t := verified.Time.UTC()
	return &t, nil
}

func (r *Repository) GetByID(ctx context.Context, id string) (*User, error) {
	return scanUser(r.db.QueryRowContext(ctx, userSelect+` WHERE id = $1`, id))
}

func (r *Repository) GetByEmail(ctx context.Context, email string) (*User, error) {
	return scanUser(r.db.QueryRowContext(ctx, userSelect+` WHERE lower(email) = lower($1)`, email))
}

func (r *Repository) UsernameTaken(ctx context.Context, tx *sql.Tx, username string) (bool, error) {
	var n int
	err := tx.QueryRowContext(ctx, `SELECT COUNT(1) FROM users WHERE username = $1`, username).Scan(&n)
	if err != nil {
		return false, fmt.Errorf("user: username taken: %w", err)
	}
	return n > 0, nil
}

func (r *Repository) Create(ctx context.Context, tx *sql.Tx, u *User) error {
	var hash any
	if u.PasswordHash != nil {
		hash = *u.PasswordHash
	}
	err := tx.QueryRowContext(ctx, `
		INSERT INTO users (username, email, password_hash, display_name, email_verified_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`, u.Username, u.Email, hash, u.DisplayName, u.EmailVerifiedAt).Scan(&u.ID, &u.CreatedAt)
	if err != nil {
		if isUniqueViolation(err) {
			return types.ErrEmailTaken
		}
		return fmt.Errorf("user: create: %w", err)
	}
	return nil
}

func (r *Repository) MarkVerified(ctx context.Context, tx *sql.Tx, userID string, at time.Time) error {
	res, err := tx.ExecContext(ctx, `UPDATE users SET email_verified_at = $2, updated_at = $2 WHERE id = $1`, userID, at)
	if err != nil {
		return fmt.Errorf("user: mark verified: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("user: mark verified rows: %w", err)
	}
	if n == 0 {
		return types.ErrNotFound
	}
	return nil
}

func (r *Repository) SetLastProject(ctx context.Context, tx *sql.Tx, userID, projectID string) error {
	_, err := tx.ExecContext(ctx, `UPDATE users SET last_project_id = $2, updated_at = CURRENT_TIMESTAMP WHERE id = $1`, userID, projectID)
	if err != nil {
		return fmt.Errorf("user: set last project: %w", err)
	}
	return nil
}

func (r *Repository) CreateAuthToken(ctx context.Context, tx *sql.Tx, kind, email, tokenHash string, expiresAt time.Time) error {
	_, err := tx.ExecContext(ctx, `
		INSERT INTO auth_tokens (kind, email, token_hash, expires_at)
		VALUES ($1, $2, $3, $4)
	`, kind, email, tokenHash, expiresAt)
	if err != nil {
		return fmt.Errorf("user: create auth token: %w", err)
	}
	return nil
}

func (r *Repository) ConsumeAuthToken(ctx context.Context, tx *sql.Tx, kind, tokenHash string, now time.Time) (*AuthToken, error) {
	tok := &AuthToken{}
	var consumed sql.NullTime
	err := tx.QueryRowContext(ctx, `
		SELECT id, kind, email, token_hash, expires_at, consumed_at
		FROM auth_tokens
		WHERE token_hash = $1 AND kind = $2
		FOR UPDATE
	`, tokenHash, kind).Scan(&tok.ID, &tok.Kind, &tok.Email, &tok.TokenHash, &tok.ExpiresAt, &consumed)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, types.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("user: load auth token: %w", err)
	}
	if consumed.Valid {
		return nil, types.ErrTokenConsumed
	}
	if !now.Before(tok.ExpiresAt) {
		return nil, types.ErrTokenExpired
	}
	if _, err := tx.ExecContext(ctx, `UPDATE auth_tokens SET consumed_at = $2 WHERE id = $1`, tok.ID, now); err != nil {
		return nil, fmt.Errorf("user: consume auth token: %w", err)
	}
	tok.ConsumedAt = &now
	return tok, nil
}

func (r *Repository) CreateSession(ctx context.Context, tx *sql.Tx, userID, tokenHash string, expiresAt time.Time, lastProjectID *string) (*Session, error) {
	s := &Session{UserID: userID, TokenHash: tokenHash, ExpiresAt: expiresAt, LastProjectID: lastProjectID}
	err := tx.QueryRowContext(ctx, `
		INSERT INTO sessions (user_id, token_hash, expires_at, last_project_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`, userID, tokenHash, expiresAt, lastProjectID).Scan(&s.ID)
	if err != nil {
		return nil, fmt.Errorf("user: create session: %w", err)
	}
	return s, nil
}

func (r *Repository) SessionByHash(ctx context.Context, tokenHash string, now time.Time) (*Session, error) {
	s := &Session{}
	var last sql.NullString
	err := r.db.QueryRowContext(ctx, `
		SELECT id, user_id, token_hash, expires_at, last_project_id
		FROM sessions
		WHERE token_hash = $1 AND expires_at > $2
	`, tokenHash, now).Scan(&s.ID, &s.UserID, &s.TokenHash, &s.ExpiresAt, &last)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, types.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("user: session by hash: %w", err)
	}
	if last.Valid {
		s.LastProjectID = &last.String
	}
	return s, nil
}

func (r *Repository) DeleteSession(ctx context.Context, sessionID string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM sessions WHERE id = $1`, sessionID)
	if err != nil {
		return fmt.Errorf("user: delete session: %w", err)
	}
	return nil
}

func (r *Repository) SetSessionLastProject(ctx context.Context, tx *sql.Tx, sessionID, projectID string) error {
	_, err := tx.ExecContext(ctx, `UPDATE sessions SET last_project_id = $2 WHERE id = $1`, sessionID, projectID)
	if err != nil {
		return fmt.Errorf("user: session last project: %w", err)
	}
	return nil
}

const userSelect = `
	SELECT id, username, email, password_hash, display_name, email_verified_at, last_project_id, created_at
	FROM users
`

func scanUser(row *sql.Row) (*User, error) {
	u := &User{}
	var hash sql.NullString
	var verified sql.NullTime
	var last sql.NullString
	err := row.Scan(&u.ID, &u.Username, &u.Email, &hash, &u.DisplayName, &verified, &last, &u.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, types.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("user: scan: %w", err)
	}
	if hash.Valid {
		u.PasswordHash = &hash.String
	}
	if verified.Valid {
		t := verified.Time.UTC()
		u.EmailVerifiedAt = &t
	}
	if last.Valid {
		u.LastProjectID = &last.String
	}
	return u, nil
}

func isUniqueViolation(err error) bool {
	if err == nil {
		return false
	}
	s := err.Error()
	return strings.Contains(s, "duplicate key") || strings.Contains(s, "unique constraint") || strings.Contains(s, "idx_users_email") || strings.Contains(s, "idx_users_username")
}
