package planning

import (
	"testing"
	"time"
)

func TestCurrentWindowUsesStartWeekdayOnOrBeforeCreated(t *testing.T) {
	loc, err := time.LoadLocation("Australia/Melbourne")
	if err != nil {
		t.Fatalf("load location: %v", err)
	}
	// Friday 21 Aug 2026 10:10 Melbourne. Start weekday Monday.
	created := time.Date(2026, 8, 21, 10, 10, 0, 0, loc)
	now := created
	win, err := CurrentWindow(now, "Australia/Melbourne", created, 1, 7)
	if err != nil {
		t.Fatalf("CurrentWindow: %v", err)
	}
	wantStart := time.Date(2026, 8, 17, 0, 0, 0, 0, loc)
	wantEnd := time.Date(2026, 8, 24, 0, 0, 0, 0, loc)
	if !win.StartsOn.Equal(wantStart) {
		t.Fatalf("starts: got %s, want %s", win.StartsOn, wantStart)
	}
	if !win.EndsAt.Equal(wantEnd) {
		t.Fatalf("ends: got %s, want %s", win.EndsAt, wantEnd)
	}
}

func TestCurrentWindowAdvancesAfterEnd(t *testing.T) {
	loc, err := time.LoadLocation("Australia/Melbourne")
	if err != nil {
		t.Fatalf("load location: %v", err)
	}
	created := time.Date(2026, 8, 17, 0, 0, 0, 0, loc)
	now := time.Date(2026, 8, 24, 0, 0, 0, 0, loc)
	win, err := CurrentWindow(now, "Australia/Melbourne", created, 1, 7)
	if err != nil {
		t.Fatalf("CurrentWindow: %v", err)
	}
	wantStart := time.Date(2026, 8, 24, 0, 0, 0, 0, loc)
	wantEnd := time.Date(2026, 8, 31, 0, 0, 0, 0, loc)
	if !win.StartsOn.Equal(wantStart) {
		t.Fatalf("starts: got %s, want %s", win.StartsOn, wantStart)
	}
	if !win.EndsAt.Equal(wantEnd) {
		t.Fatalf("ends: got %s, want %s", win.EndsAt, wantEnd)
	}
}

func TestCurrentWindowRejectsInvalidInputs(t *testing.T) {
	now := time.Date(2026, 8, 21, 0, 0, 0, 0, time.UTC)
	if _, err := CurrentWindow(now, "", now, 1, 7); err == nil {
		t.Fatal("expected timezone required")
	}
	if _, err := CurrentWindow(now, "Australia/Melbourne", now, 1, 0); err == nil {
		t.Fatal("expected positive length")
	}
	if _, err := CurrentWindow(now, "Australia/Melbourne", now, 0, 7); err == nil {
		t.Fatal("expected weekday 1-7")
	}
	if _, err := CurrentWindow(now, "Not/AZone", now, 1, 7); err == nil {
		t.Fatal("expected invalid timezone")
	}
}

func TestColdStartPackUsesInitialVelocity(t *testing.T) {
	loc, err := time.LoadLocation("Australia/Melbourne")
	if err != nil {
		t.Fatalf("load location: %v", err)
	}
	win := Window{EndsAt: time.Date(2026, 8, 24, 0, 0, 0, 0, loc)}
	pack := ColdStartPack(win, 10)
	if pack.CurrentPoints != 0 || pack.Denominator != 10 || pack.VelocitySource != "initial" {
		t.Fatalf("unexpected pack: %+v", pack)
	}
	if !pack.CurrentWindowEndsAt.Equal(win.EndsAt) {
		t.Fatalf("ends: got %s", pack.CurrentWindowEndsAt)
	}
}
