package organisation

import (
	"context"
	"strings"

	"flower/api/internal/ports"
	"flower/api/internal/types"
)

type Service struct {
	repo     *Repository
	identity ports.Identity
}

func NewService(repo *Repository, identity ports.Identity) *Service {
	return &Service{repo: repo, identity: identity}
}

func (s *Service) Create(ctx context.Context, userID, name string) (*Organisation, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, types.Wrap(types.ErrValidation, "Name the organisation.")
	}
	verified, err := s.identity.EmailVerifiedAt(ctx, userID)
	if err != nil {
		return nil, err
	}
	if verified == nil {
		return nil, types.ErrEmailUnverified
	}
	return s.repo.Create(ctx, userID, name)
}
