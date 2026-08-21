package types

import "errors"

var (
	ErrNotFound        = errors.New("not_found")
	ErrUnauthorized    = errors.New("unauthorized")
	ErrForbidden       = errors.New("forbidden")
	ErrValidation      = errors.New("validation_failed")
	ErrEmailTaken      = errors.New("email_taken")
	ErrEmailUnverified = errors.New("email_unverified")
	ErrConflict        = errors.New("conflict")
	ErrUseMagicLink    = errors.New("use_magic_link")
	ErrTokenExpired    = errors.New("token_expired")
	ErrTokenConsumed   = errors.New("token_consumed")
)

type Error struct {
	Kind error
	Msg  string
}

func (e *Error) Error() string {
	if e == nil {
		return ""
	}
	if e.Msg != "" {
		return e.Msg
	}
	if e.Kind != nil {
		return e.Kind.Error()
	}
	return "error"
}

func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Kind
}

func Wrap(kind error, msg string) error {
	return &Error{Kind: kind, Msg: msg}
}
