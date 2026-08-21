package planning

import (
	"context"

	"flower/api/internal/ports"
)

type SettingsLoader func(ctx context.Context, projectID string) (ports.ProjectWindow, error)

type Planner struct {
	clock ports.Clock
	load  SettingsLoader
}

func NewPlanner(clock ports.Clock, load SettingsLoader) *Planner {
	return &Planner{clock: clock, load: load}
}

func (p *Planner) ForProject(ctx context.Context, projectID string) (ports.Pack, error) {
	settings, err := p.load(ctx, projectID)
	if err != nil {
		return ports.Pack{}, err
	}
	win, err := CurrentWindow(p.clock.Now(), settings.Timezone, settings.CreatedAt, settings.StartWeekday, settings.LengthDays)
	if err != nil {
		return ports.Pack{}, err
	}
	pack := ColdStartPack(win, settings.InitialVelocity)
	return ports.Pack{
		CurrentPoints:       pack.CurrentPoints,
		Denominator:         pack.Denominator,
		VelocitySource:      pack.VelocitySource,
		CurrentWindowEndsAt: pack.CurrentWindowEndsAt,
	}, nil
}
