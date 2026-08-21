package organisation

import (
	"context"
	"database/sql"
	"fmt"

	"flower/api/internal/platform/db"
	"flower/api/internal/ports"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(database *sql.DB) *Repository {
	return &Repository{db: database}
}

func (r *Repository) Create(ctx context.Context, userID, name string) (*Organisation, error) {
	org := &Organisation{Name: name}
	err := db.WithTx(ctx, r.db, func(tx *sql.Tx) error {
		if err := tx.QueryRowContext(ctx, `
			INSERT INTO organisations (name) VALUES ($1) RETURNING id, created_at
		`, name).Scan(&org.ID, &org.CreatedAt); err != nil {
			return fmt.Errorf("organisation: insert: %w", err)
		}
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO organisation_memberships (organisation_id, user_id, role)
			VALUES ($1, $2, 'owner')
		`, org.ID, userID); err != nil {
			return fmt.Errorf("organisation: membership: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return org, nil
}

func (r *Repository) OrganisationsForUser(ctx context.Context, userID string) ([]ports.OrganisationInfo, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT o.id, o.name, m.role
		FROM organisation_memberships m
		JOIN organisations o ON o.id = m.organisation_id
		WHERE m.user_id = $1
		ORDER BY o.created_at ASC
	`, userID)
	if err != nil {
		return nil, fmt.Errorf("organisation: list for user: %w", err)
	}
	defer rows.Close()
	var out []ports.OrganisationInfo
	for rows.Next() {
		var info ports.OrganisationInfo
		if err := rows.Scan(&info.ID, &info.Name, &info.Role); err != nil {
			return nil, fmt.Errorf("organisation: list scan: %w", err)
		}
		out = append(out, info)
	}
	if out == nil {
		out = []ports.OrganisationInfo{}
	}
	return out, rows.Err()
}
