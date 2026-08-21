package ports

import "context"

type OrganisationInfo struct {
	ID   string
	Name string
	Role string
}

type ProjectInfo struct {
	ID             string
	OrganisationID string
	Name           string
	Slug           string
}

type Directory interface {
	OrganisationsForUser(ctx context.Context, userID string) ([]OrganisationInfo, error)
	Project(ctx context.Context, projectID string) (*ProjectInfo, error)
}
