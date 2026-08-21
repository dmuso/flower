package ports

import (
	"context"
	"time"
)

type Identity interface {
	EmailVerifiedAt(ctx context.Context, userID string) (*time.Time, error)
}
