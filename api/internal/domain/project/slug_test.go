package project

import "testing"

func TestSlugify(t *testing.T) {
	if got := slugify("Trail Board"); got != "trail-board" {
		t.Fatalf("got %q", got)
	}
	if got := slugify("  "); got != "project" {
		t.Fatalf("got %q", got)
	}
	if got := slugify("Acme!!"); got != "acme" {
		t.Fatalf("got %q", got)
	}
}
