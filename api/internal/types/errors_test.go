package types

import (
	"errors"
	"testing"
)

func TestWrapUnwrap(t *testing.T) {
	err := Wrap(ErrNotFound, "We can’t find that.")
	if !errors.Is(err, ErrNotFound) {
		t.Fatal("unwrap")
	}
	if err.Error() != "We can’t find that." {
		t.Fatalf("msg %s", err.Error())
	}
	var empty *Error
	if empty.Error() != "" {
		t.Fatal("nil error")
	}
	if empty.Unwrap() != nil {
		t.Fatal("nil unwrap")
	}
	plain := &Error{Kind: ErrForbidden}
	if plain.Error() != ErrForbidden.Error() {
		t.Fatalf("kind msg %s", plain.Error())
	}
	if (&Error{}).Error() != "error" {
		t.Fatal("empty")
	}
}
