package ports

import (
	"context"
	"time"
)

type Pack struct {
	CurrentPoints       int
	Denominator         int
	VelocitySource      string
	CurrentWindowEndsAt time.Time
}

type ProjectWindow struct {
	Timezone        string
	CreatedAt       time.Time
	LengthDays      int
	StartWeekday    int
	InitialVelocity int
}

type Planner interface {
	ForProject(ctx context.Context, projectID string) (Pack, error)
}
