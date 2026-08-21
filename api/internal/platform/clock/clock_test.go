package clock

import (
	"testing"
	"time"
)

func TestFixedNow(t *testing.T) {
	want := time.Date(2026, 8, 21, 0, 10, 0, 0, time.UTC)
	if got := (Fixed{T: want}).Now(); !got.Equal(want) {
		t.Fatalf("got %s", got)
	}
}

func TestSystemNow(t *testing.T) {
	before := time.Now().UTC().Add(-time.Second)
	got := System{}.Now()
	after := time.Now().UTC().Add(time.Second)
	if got.Before(before) || got.After(after) {
		t.Fatalf("system now %s", got)
	}
}
