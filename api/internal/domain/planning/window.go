package planning

import (
	"fmt"
	"time"
)

const DefaultTimezone = "Australia/Melbourne"

type Window struct {
	StartsOn time.Time
	EndsAt   time.Time
}

func isoWeekday(t time.Time) int {
	d := int(t.Weekday())
	if d == 0 {
		return 7
	}
	return d
}

func CurrentWindow(now time.Time, timezone string, createdAt time.Time, startWeekday, lengthDays int) (Window, error) {
	if timezone == "" {
		return Window{}, fmt.Errorf("planning: timezone is required")
	}
	if lengthDays <= 0 {
		return Window{}, fmt.Errorf("planning: iteration_length_days must be positive")
	}
	if startWeekday < 1 || startWeekday > 7 {
		return Window{}, fmt.Errorf("planning: iteration_start_weekday must be 1-7")
	}
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return Window{}, fmt.Errorf("planning: load timezone %q: %w", timezone, err)
	}

	createdLocal := createdAt.In(loc)
	start := time.Date(createdLocal.Year(), createdLocal.Month(), createdLocal.Day(), 0, 0, 0, 0, loc)
	delta := isoWeekday(start) - startWeekday
	if delta < 0 {
		delta += 7
	}
	start = start.AddDate(0, 0, -delta)

	nowLocal := now.In(loc)
	for {
		end := start.AddDate(0, 0, lengthDays)
		if !nowLocal.Before(end) {
			start = end
			continue
		}
		return Window{StartsOn: start, EndsAt: end}, nil
	}
}

type Pack struct {
	CurrentPoints       int
	Denominator         int
	VelocitySource      string
	CurrentWindowEndsAt time.Time
}

func ColdStartPack(window Window, initialVelocity int) Pack {
	return Pack{
		CurrentPoints:       0,
		Denominator:         initialVelocity,
		VelocitySource:      "initial",
		CurrentWindowEndsAt: window.EndsAt,
	}
}
