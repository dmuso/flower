package ports

import "context"

type Access interface {
	OrganisationRole(ctx context.Context, userID, organisationID string) (string, error)
	OpenProject(ctx context.Context, userID, projectID string) (organisationID string, role string, err error)
}
