package organisation

import (
	"context"
	"errors"
	"testing"
	"time"

	"flower/api/internal/types"
)

type stubIdentity struct {
	at *time.Time
}

func (s stubIdentity) EmailVerifiedAt(context.Context, string) (*time.Time, error) {
	return s.at, nil
}

func TestCreateRejectsBlankName(t *testing.T) {
	now := time.Date(2026, 8, 21, 0, 10, 0, 0, time.UTC)
	svc := NewService(nil, stubIdentity{at: &now})
	_, err := svc.Create(context.Background(), "user-1", "   ")
	if !errors.Is(err, types.ErrValidation) {
		t.Fatalf("got %v", err)
	}
}

func TestCreateRejectsUnverifiedEmail(t *testing.T) {
	svc := NewService(nil, stubIdentity{})
	_, err := svc.Create(context.Background(), "user-1", "Acme")
	if !errors.Is(err, types.ErrEmailUnverified) {
		t.Fatalf("got %v", err)
	}
}
