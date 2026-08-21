package story

import (
	"context"

	"flower/api/internal/ports"
	"flower/api/internal/types"
)

type Service struct {
	stories *Repository
	access  ports.Access
	planner ports.Planner
}

func NewService(stories *Repository, access ports.Access, planner ports.Planner) *Service {
	return &Service{stories: stories, access: access, planner: planner}
}

func (s *Service) List(ctx context.Context, userID, projectID string) (*List, error) {
	if _, _, err := s.access.OpenProject(ctx, userID, projectID); err != nil {
		return nil, types.ErrNotFound
	}
	pack, err := s.planner.ForProject(ctx, projectID)
	if err != nil {
		return nil, err
	}
	stories, err := s.stories.List(ctx, projectID)
	if err != nil {
		return nil, err
	}
	return &List{
		Stories: stories,
		Pack: Pack{
			CurrentPoints:       pack.CurrentPoints,
			Denominator:         pack.Denominator,
			VelocitySource:      pack.VelocitySource,
			CurrentWindowEndsAt: pack.CurrentWindowEndsAt.UTC(),
		},
	}, nil
}
