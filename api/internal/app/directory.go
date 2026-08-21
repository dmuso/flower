package app

import (
	"context"

	"flower/api/internal/domain/organisation"
	"flower/api/internal/domain/project"
	"flower/api/internal/ports"
)

type directory struct {
	orgs     *organisation.Repository
	projects *project.Repository
}

func (d *directory) OrganisationsForUser(ctx context.Context, userID string) ([]ports.OrganisationInfo, error) {
	return d.orgs.OrganisationsForUser(ctx, userID)
}

func (d *directory) Project(ctx context.Context, projectID string) (*ports.ProjectInfo, error) {
	return d.projects.Project(ctx, projectID)
}
