package user

import "testing"

func TestInferUsernameFromLocalPart(t *testing.T) {
	got := inferUsername("Maya.Harper+dev@example.com")
	if got != "mayaharperdev" {
		t.Fatalf("got %q", got)
	}
}

func TestInferUsernameKeepsUnderscoreAndDigits(t *testing.T) {
	got := inferUsername("luis_2@example.com")
	if got != "luis_2" {
		t.Fatalf("got %q", got)
	}
}

func TestInferUsernameEmptyLocalBecomesUser(t *testing.T) {
	got := inferUsername("@example.com")
	if got != "user" {
		t.Fatalf("got %q", got)
	}
}

func TestInferUsernameTrimsTo100(t *testing.T) {
	local := stringsRepeat("a", 120)
	got := inferUsername(local + "@example.com")
	if len(got) != 100 {
		t.Fatalf("len=%d", len(got))
	}
}

func TestUniquifyUsernameAppendsSuffixWhenTaken(t *testing.T) {
	taken := func(name string) (bool, error) {
		return name == "maya", nil
	}
	got, err := uniquifyUsername("maya", taken)
	if err != nil {
		t.Fatalf("uniquify: %v", err)
	}
	if got == "maya" {
		t.Fatal("expected suffix")
	}
	if len(got) != len("maya-xxxx") {
		t.Fatalf("got %q", got)
	}
	if got[:5] != "maya-" {
		t.Fatalf("got %q", got)
	}
}

func stringsRepeat(s string, n int) string {
	out := make([]byte, 0, n*len(s))
	for i := 0; i < n; i++ {
		out = append(out, s...)
	}
	return string(out)
}
