package auth

import (
	"strings"
	"testing"
)

func TestHashTokenIsHexAndStable(t *testing.T) {
	a := HashToken("abc")
	b := HashToken("abc")
	if a != b {
		t.Fatal("hash must be deterministic")
	}
	if len(a) != 64 {
		t.Fatalf("len %d", len(a))
	}
	if HashToken("abd") == a {
		t.Fatal("different input must differ")
	}
}

func TestRandomToken(t *testing.T) {
	one, err := RandomToken()
	if err != nil {
		t.Fatal(err)
	}
	two, err := RandomToken()
	if err != nil {
		t.Fatal(err)
	}
	if one == two {
		t.Fatal("tokens must be unique")
	}
	if len(one) != 64 {
		t.Fatalf("len %d", len(one))
	}
}

func TestRandomUnambiguous(t *testing.T) {
	got, err := RandomUnambiguous(4)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 4 {
		t.Fatalf("len %d", len(got))
	}
	for _, r := range got {
		if !strings.ContainsRune(Unambiguous, r) {
			t.Fatalf("ambiguous %q", r)
		}
	}
	if _, err := RandomUnambiguous(0); err == nil {
		t.Fatal("expected error")
	}
}
