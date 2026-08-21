package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/mail"
	"strings"
	"time"

	"flower/api/internal/platform/auth"
	"flower/api/internal/ports"
	"flower/api/internal/types"

	"golang.org/x/crypto/bcrypt"
)

const bcryptCost = 12

type Service struct {
	repo     *Repository
	clock    ports.Clock
	mailer   ports.Mailer
	origin   string
	dir      ports.Directory
}

func NewService(repo *Repository, clock ports.Clock, mailer ports.Mailer, origin string, dir ports.Directory) *Service {
	return &Service{repo: repo, clock: clock, mailer: mailer, origin: origin, dir: dir}
}

func (s *Service) LookupSession(ctx context.Context, tokenHash string) (*types.SessionUser, error) {
	sess, err := s.repo.SessionByHash(ctx, tokenHash, s.clock.Now())
	if err != nil {
		return nil, err
	}
	out := &types.SessionUser{UserID: sess.UserID, SessionID: sess.ID}
	if sess.LastProjectID != nil {
		out.LastProjectID = *sess.LastProjectID
	}
	return out, nil
}

func (s *Service) SignUp(ctx context.Context, email, password string) error {
	email, err := normalizeEmail(email)
	if err != nil {
		return types.Wrap(types.ErrValidation, "Need an email and a password.")
	}
	if strings.TrimSpace(password) == "" {
		return types.Wrap(types.ErrValidation, "Need an email and a password.")
	}
	if _, err := s.repo.GetByEmail(ctx, email); err == nil {
		return types.ErrEmailTaken
	} else if !errors.Is(err, types.ErrNotFound) {
		return err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return fmt.Errorf("user: hash password: %w", err)
	}
	hashStr := string(hash)
	now := s.clock.Now()
	raw, err := auth.RandomToken()
	if err != nil {
		return err
	}
	return s.repo.WithTx(ctx, func(tx *sql.Tx) error {
		username, err := s.allocateUsername(ctx, tx, email)
		if err != nil {
			return err
		}
		u := &User{
			Username:     username,
			Email:        email,
			PasswordHash: &hashStr,
			DisplayName:  username,
		}
		if err := s.repo.Create(ctx, tx, u); err != nil {
			return err
		}
		if err := s.repo.CreateAuthToken(ctx, tx, KindVerifyEmail, email, auth.HashToken(raw), now.Add(TokenTTL)); err != nil {
			return err
		}
		return s.mailer.Enqueue(ctx, tx, ports.MailMessage{
			Kind:    KindVerifyEmail,
			To:      email,
			Subject: "Verify your Flower email",
			Body:    "Verify your email for Flower.\n\n" + s.origin + "/verify-email?token=" + raw + "\n",
		})
	})
}

func (s *Service) Login(ctx context.Context, email, password string) (string, time.Time, *Me, error) {
	email, err := normalizeEmail(email)
	if err != nil {
		return "", time.Time{}, nil, types.Wrap(types.ErrUnauthorized, "That email and password don’t match.")
	}
	if strings.TrimSpace(password) == "" {
		return "", time.Time{}, nil, types.Wrap(types.ErrUnauthorized, "That email and password don’t match.")
	}
	u, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, types.ErrNotFound) {
			return "", time.Time{}, nil, types.ErrUnauthorized
		}
		return "", time.Time{}, nil, err
	}
	if u.PasswordHash == nil {
		return "", time.Time{}, nil, types.ErrUseMagicLink
	}
	if err := bcrypt.CompareHashAndPassword([]byte(*u.PasswordHash), []byte(password)); err != nil {
		return "", time.Time{}, nil, types.ErrUnauthorized
	}
	raw, expires, err := s.openSession(ctx, u)
	if err != nil {
		return "", time.Time{}, nil, err
	}
	me, err := s.meFromUser(ctx, u)
	return raw, expires, me, err
}

func (s *Service) RequestMagicLink(ctx context.Context, email string) error {
	email, err := normalizeEmail(email)
	if err != nil {
		return types.Wrap(types.ErrValidation, "Need an email.")
	}
	now := s.clock.Now()
	raw, err := auth.RandomToken()
	if err != nil {
		return err
	}
	return s.repo.WithTx(ctx, func(tx *sql.Tx) error {
		if err := s.repo.CreateAuthToken(ctx, tx, KindMagicLink, email, auth.HashToken(raw), now.Add(TokenTTL)); err != nil {
			return err
		}
		return s.mailer.Enqueue(ctx, tx, ports.MailMessage{
			Kind:    KindMagicLink,
			To:      email,
			Subject: "Your Flower sign-in link",
			Body:    "Sign in to Flower.\n\n" + s.origin + "/magic-link?token=" + raw + "\n",
		})
	})
}

func (s *Service) ConsumeMagicLink(ctx context.Context, rawToken string) (string, time.Time, *Me, error) {
	if strings.TrimSpace(rawToken) == "" {
		return "", time.Time{}, nil, types.ErrTokenExpired
	}
	now := s.clock.Now()
	var user *User
	err := s.repo.WithTx(ctx, func(tx *sql.Tx) error {
		tok, err := s.repo.ConsumeAuthToken(ctx, tx, KindMagicLink, auth.HashToken(rawToken), now)
		if err != nil {
			return err
		}
		existing, err := s.repo.GetByEmail(ctx, tok.Email)
		if err != nil && !errors.Is(err, types.ErrNotFound) {
			return err
		}
		if existing != nil {
			if existing.EmailVerifiedAt == nil {
				if err := s.repo.MarkVerified(ctx, tx, existing.ID, now); err != nil {
					return err
				}
				existing.EmailVerifiedAt = &now
			}
			user = existing
			return nil
		}
		username, err := s.allocateUsername(ctx, tx, tok.Email)
		if err != nil {
			return err
		}
		u := &User{
			Username:        username,
			Email:           tok.Email,
			DisplayName:     username,
			EmailVerifiedAt: &now,
		}
		if err := s.repo.Create(ctx, tx, u); err != nil {
			return err
		}
		user = u
		return nil
	})
	if err != nil {
		return "", time.Time{}, nil, err
	}
	raw, expires, err := s.openSession(ctx, user)
	if err != nil {
		return "", time.Time{}, nil, err
	}
	me, err := s.Me(ctx, user.ID)
	return raw, expires, me, err
}

func (s *Service) ConsumeVerifyEmail(ctx context.Context, rawToken string) (string, time.Time, *Me, error) {
	if strings.TrimSpace(rawToken) == "" {
		return "", time.Time{}, nil, types.ErrTokenExpired
	}
	now := s.clock.Now()
	var user *User
	err := s.repo.WithTx(ctx, func(tx *sql.Tx) error {
		tok, err := s.repo.ConsumeAuthToken(ctx, tx, KindVerifyEmail, auth.HashToken(rawToken), now)
		if err != nil {
			return err
		}
		u, err := s.repo.GetByEmail(ctx, tok.Email)
		if err != nil {
			return err
		}
		if err := s.repo.MarkVerified(ctx, tx, u.ID, now); err != nil {
			return err
		}
		u.EmailVerifiedAt = &now
		user = u
		return nil
	})
	if err != nil {
		return "", time.Time{}, nil, err
	}
	raw, expires, err := s.openSession(ctx, user)
	if err != nil {
		return "", time.Time{}, nil, err
	}
	me, err := s.Me(ctx, user.ID)
	return raw, expires, me, err
}

func (s *Service) Logout(ctx context.Context, sessionID string) error {
	return s.repo.DeleteSession(ctx, sessionID)
}

func (s *Service) Me(ctx context.Context, userID string) (*Me, error) {
	u, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return s.meFromUser(ctx, u)
}

func (s *Service) RememberProject(ctx context.Context, userID, sessionID, projectID string) error {
	return s.repo.WithTx(ctx, func(tx *sql.Tx) error {
		if err := s.repo.SetLastProject(ctx, tx, userID, projectID); err != nil {
			return err
		}
		if sessionID == "" {
			return nil
		}
		return s.repo.SetSessionLastProject(ctx, tx, sessionID, projectID)
	})
}

func (s *Service) openSession(ctx context.Context, u *User) (string, time.Time, error) {
	raw, err := auth.RandomToken()
	if err != nil {
		return "", time.Time{}, err
	}
	expires := s.clock.Now().Add(SessionTTL)
	err = s.repo.WithTx(ctx, func(tx *sql.Tx) error {
		_, err := s.repo.CreateSession(ctx, tx, u.ID, auth.HashToken(raw), expires, u.LastProjectID)
		return err
	})
	return raw, expires, err
}

func (s *Service) meFromUser(ctx context.Context, u *User) (*Me, error) {
	me := &Me{
		ID:              u.ID,
		Username:        u.Username,
		Email:           u.Email,
		DisplayName:     u.DisplayName,
		EmailVerifiedAt: u.EmailVerifiedAt,
		Organisations:   []Organisation{},
	}
	if s.dir != nil {
		orgs, err := s.dir.OrganisationsForUser(ctx, u.ID)
		if err != nil {
			return nil, err
		}
		for _, o := range orgs {
			me.Organisations = append(me.Organisations, Organisation{ID: o.ID, Name: o.Name, Role: o.Role})
		}
		if u.LastProjectID != nil && *u.LastProjectID != "" {
			p, err := s.dir.Project(ctx, *u.LastProjectID)
			if err != nil && !errors.Is(err, types.ErrNotFound) {
				return nil, err
			}
			if p != nil {
				me.LastProject = &Project{ID: p.ID, OrganisationID: p.OrganisationID, Name: p.Name, Slug: p.Slug}
			}
		}
	}
	return me, nil
}

func (s *Service) allocateUsername(ctx context.Context, tx *sql.Tx, email string) (string, error) {
	base := inferUsername(email)
	return uniquifyUsername(base, func(name string) (bool, error) {
		return s.repo.UsernameTaken(ctx, tx, name)
	})
}

func normalizeEmail(email string) (string, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	if email == "" {
		return "", types.ErrValidation
	}
	addr, err := mail.ParseAddress(email)
	if err != nil || addr.Address == "" {
		return "", types.ErrValidation
	}
	return addr.Address, nil
}
