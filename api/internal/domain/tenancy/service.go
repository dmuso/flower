package tenancy

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"flower/api/internal/types"
)

type Service struct {
	db *sql.DB
}

func NewService(database *sql.DB) *Service {
	return &Service{db: database}
}

func (s *Service) OrganisationRole(ctx context.Context, userID, organisationID string) (string, error) {
	var role string
	err := s.db.QueryRowContext(ctx, `
		SELECT role FROM organisation_memberships
		WHERE organisation_id = $1 AND user_id = $2
	`, organisationID, userID).Scan(&role)
	if errors.Is(err, sql.ErrNoRows) {
		var exists int
		if qerr := s.db.QueryRowContext(ctx, `SELECT 1 FROM organisations WHERE id = $1`, organisationID).Scan(&exists); errors.Is(qerr, sql.ErrNoRows) {
			return "", types.ErrNotFound
		} else if qerr != nil {
			return "", fmt.Errorf("tenancy: organisation exists: %w", qerr)
		}
		return "", types.ErrNotFound
	}
	if err != nil {
		return "", fmt.Errorf("tenancy: organisation role: %w", err)
	}
	return role, nil
}

func (s *Service) OpenProject(ctx context.Context, userID, projectID string) (string, string, error) {
	var orgID string
	err := s.db.QueryRowContext(ctx, `SELECT organisation_id FROM projects WHERE id = $1`, projectID).Scan(&orgID)
	if errors.Is(err, sql.ErrNoRows) {
		return "", "", types.ErrNotFound
	}
	if err != nil {
		return "", "", fmt.Errorf("tenancy: project org: %w", err)
	}
	orgRole, err := s.OrganisationRole(ctx, userID, orgID)
	if err == nil && orgRole == "owner" {
		return orgID, "owner", nil
	}
	if err != nil && !errors.Is(err, types.ErrNotFound) {
		return "", "", err
	}
	var role string
	err = s.db.QueryRowContext(ctx, `
		SELECT role FROM project_memberships
		WHERE project_id = $1 AND user_id = $2
	`, projectID, userID).Scan(&role)
	if errors.Is(err, sql.ErrNoRows) {
		return "", "", types.ErrNotFound
	}
	if err != nil {
		return "", "", fmt.Errorf("tenancy: project role: %w", err)
	}
	return orgID, role, nil
}
