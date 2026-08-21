package story

import (
	"context"
	"errors"
	"testing"
	"time"

	"flower/api/internal/platform/clock"
	"flower/api/internal/types"
)

type denyAccess struct{}

func (denyAccess) OrganisationRole(context.Context, string, string) (string, error) {
	return "", types.ErrNotFound
}

func (denyAccess) OpenProject(context.Context, string, string) (string, string, error) {
	return "", "", types.ErrNotFound
}

func TestListUnknownProjectIsNotFound(t *testing.T) {
	svc := NewService(nil, denyAccess{}, clock.Fixed{T: time.Date(2026, 8, 21, 0, 10, 0, 0, time.UTC)})
	_, err := svc.List(context.Background(), "user-1", "project-1")
	if !errors.Is(err, types.ErrNotFound) {
		t.Fatalf("got %v", err)
	}
}
