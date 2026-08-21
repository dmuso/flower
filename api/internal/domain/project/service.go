package project

import (
	"context"
	"errors"
	"strings"

	"flower/api/internal/ports"
	"flower/api/internal/types"
)

type Service struct {
	repo   *Repository
	access ports.Access
}

func NewService(repo *Repository, access ports.Access) *Service {
	return &Service{repo: repo, access: access}
}

func (s *Service) Create(ctx context.Context, userID, organisationID, name string) (*Project, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, types.Wrap(types.ErrValidation, "Name the project.")
	}
	role, err := s.access.OrganisationRole(ctx, userID, organisationID)
	if err != nil {
		if errors.Is(err, types.ErrNotFound) {
			return nil, types.ErrNotFound
		}
		return nil, err
	}
	if role != "owner" {
		return nil, types.ErrNotFound
	}
	return s.repo.Create(ctx, userID, organisationID, name)
}

func (s *Service) List(ctx context.Context, userID, organisationID string) ([]Project, error) {
	if _, err := s.access.OrganisationRole(ctx, userID, organisationID); err != nil {
		return nil, types.ErrNotFound
	}
	return s.repo.ListForOrganisation(ctx, organisationID)
}

func (s *Service) Get(ctx context.Context, userID, projectID string) (*Project, error) {
	if _, _, err := s.access.OpenProject(ctx, userID, projectID); err != nil {
		return nil, types.ErrNotFound
	}
	return s.repo.Get(ctx, projectID)
}
