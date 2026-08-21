package project

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"flower/api/internal/platform/auth"
	"flower/api/internal/platform/db"
	"flower/api/internal/ports"
	"flower/api/internal/types"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(database *sql.DB) *Repository {
	return &Repository{db: database}
}

func (r *Repository) Create(ctx context.Context, userID, organisationID, name string) (*Project, error) {
	p := &Project{
		OrganisationID:        organisationID,
		Name:                  name,
		PointScale:            DefaultPointScale,
		Timezone:              DefaultTimezone,
		VelocityStrategy:      DefaultVelocityStrategy,
		InitialVelocity:       DefaultInitialVelocity,
		IterationStartWeekday: DefaultIterationStartWeekday,
		IterationLengthDays:   DefaultIterationLengthDays,
	}
	base := slugify(name)
	err := db.WithTx(ctx, r.db, func(tx *sql.Tx) error {
		slug, err := r.allocateSlug(ctx, tx, organisationID, base)
		if err != nil {
			return err
		}
		p.Slug = slug
		// iteration_length_weeks is a leftover NOT NULL column from 000001. Set 1 so inserts succeed. Do not write iterations.
		if err := tx.QueryRowContext(ctx, `
			INSERT INTO projects (
				organisation_id, name, slug, point_scale, iteration_length_weeks,
				timezone, velocity_strategy, initial_velocity, iteration_start_weekday, iteration_length_days
			) VALUES ($1, $2, $3, $4, 1, $5, $6, $7, $8, $9)
			RETURNING id, created_at
		`, organisationID, name, slug, p.PointScale, p.Timezone, p.VelocityStrategy, p.InitialVelocity, p.IterationStartWeekday, p.IterationLengthDays).Scan(&p.ID, &p.CreatedAt); err != nil {
			return fmt.Errorf("project: insert: %w", err)
		}
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO project_memberships (project_id, user_id, role)
			VALUES ($1, $2, 'owner')
		`, p.ID, userID); err != nil {
			return fmt.Errorf("project: membership: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (r *Repository) allocateSlug(ctx context.Context, tx *sql.Tx, organisationID, base string) (string, error) {
	taken := func(slug string) (bool, error) {
		var n int
		if err := tx.QueryRowContext(ctx, `SELECT COUNT(1) FROM projects WHERE organisation_id = $1 AND slug = $2`, organisationID, slug).Scan(&n); err != nil {
			return false, fmt.Errorf("project: slug taken: %w", err)
		}
		return n > 0, nil
	}
	ok, err := taken(base)
	if err != nil {
		return "", err
	}
	if !ok {
		return base, nil
	}
	root := base
	if len(root) > 95 {
		root = root[:95]
	}
	for i := 0; i < 32; i++ {
		suffix, err := auth.RandomUnambiguous(4)
		if err != nil {
			return "", err
		}
		candidate := strings.Trim(root, "-") + "-" + suffix
		inUse, err := taken(candidate)
		if err != nil {
			return "", err
		}
		if !inUse {
			return candidate, nil
		}
	}
	return "", fmt.Errorf("project: could not allocate a unique slug")
}

func (r *Repository) Get(ctx context.Context, id string) (*Project, error) {
	p := &Project{}
	var desc sql.NullString
	err := r.db.QueryRowContext(ctx, `
		SELECT id, organisation_id, name, slug, description, point_scale, timezone,
		       velocity_strategy, initial_velocity, iteration_start_weekday, iteration_length_days, created_at
		FROM projects WHERE id = $1
	`, id).Scan(&p.ID, &p.OrganisationID, &p.Name, &p.Slug, &desc, &p.PointScale, &p.Timezone, &p.VelocityStrategy, &p.InitialVelocity, &p.IterationStartWeekday, &p.IterationLengthDays, &p.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, types.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("project: get: %w", err)
	}
	if desc.Valid {
		p.Description = &desc.String
	}
	return p, nil
}

func (r *Repository) ListForOrganisation(ctx context.Context, organisationID string) ([]Project, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, organisation_id, name, slug, description, point_scale, timezone,
		       velocity_strategy, initial_velocity, iteration_start_weekday, iteration_length_days, created_at
		FROM projects WHERE organisation_id = $1 ORDER BY created_at ASC
	`, organisationID)
	if err != nil {
		return nil, fmt.Errorf("project: list: %w", err)
	}
	defer rows.Close()
	var out []Project
	for rows.Next() {
		p := Project{}
		var desc sql.NullString
		if err := rows.Scan(&p.ID, &p.OrganisationID, &p.Name, &p.Slug, &desc, &p.PointScale, &p.Timezone, &p.VelocityStrategy, &p.InitialVelocity, &p.IterationStartWeekday, &p.IterationLengthDays, &p.CreatedAt); err != nil {
			return nil, fmt.Errorf("project: list scan: %w", err)
		}
		if desc.Valid {
			p.Description = &desc.String
		}
		out = append(out, p)
	}
	if out == nil {
		out = []Project{}
	}
	return out, rows.Err()
}

func (r *Repository) Project(ctx context.Context, projectID string) (*ports.ProjectInfo, error) {
	p, err := r.Get(ctx, projectID)
	if err != nil {
		return nil, err
	}
	return &ports.ProjectInfo{ID: p.ID, OrganisationID: p.OrganisationID, Name: p.Name, Slug: p.Slug}, nil
}
