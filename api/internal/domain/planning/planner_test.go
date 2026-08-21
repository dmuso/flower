package planning

import (
	"context"
	"testing"
	"time"

	"flower/api/internal/platform/clock"
	"flower/api/internal/ports"
)

func TestPlannerForProjectUsesWindowAndColdStart(t *testing.T) {
	loc, err := time.LoadLocation("Australia/Melbourne")
	if err != nil {
		t.Fatal(err)
	}
	created := time.Date(2026, 8, 21, 10, 10, 0, 0, loc)
	now := created
	p := NewPlanner(clock.Fixed{T: now}, func(context.Context, string) (ports.ProjectWindow, error) {
		return ports.ProjectWindow{
			Timezone:        "Australia/Melbourne",
			CreatedAt:       created,
			LengthDays:      7,
			StartWeekday:    1,
			InitialVelocity: 10,
		}, nil
	})
	pack, err := p.ForProject(context.Background(), "project-1")
	if err != nil {
		t.Fatal(err)
	}
	if pack.CurrentPoints != 0 || pack.Denominator != 10 || pack.VelocitySource != "initial" {
		t.Fatalf("pack %+v", pack)
	}
	wantEnd := time.Date(2026, 8, 24, 0, 0, 0, 0, loc)
	if !pack.CurrentWindowEndsAt.Equal(wantEnd) {
		t.Fatalf("end %s want %s", pack.CurrentWindowEndsAt, wantEnd)
	}
}
