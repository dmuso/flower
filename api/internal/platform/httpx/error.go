package httpx

import (
	"errors"
	"net/http"

	"flower/api/internal/types"

	"github.com/gin-gonic/gin"
)

type envelope struct {
	Error body `json:"error"`
}

type body struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func Write(c *gin.Context, status int, code, message string) {
	c.JSON(status, envelope{Error: body{Code: code, Message: message}})
}

func messageOf(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

func From(c *gin.Context, err error) {
	switch {
	case errors.Is(err, types.ErrValidation):
		msg := messageOf(err)
		if msg == "" || msg == types.ErrValidation.Error() {
			msg = "Need an email and a password."
		}
		Write(c, http.StatusBadRequest, "validation_failed", msg)
	case errors.Is(err, types.ErrEmailTaken):
		Write(c, http.StatusConflict, "email_taken", "That email already belongs to a user.")
	case errors.Is(err, types.ErrEmailUnverified):
		Write(c, http.StatusForbidden, "email_unverified", "Verify your email to continue.")
	case errors.Is(err, types.ErrUseMagicLink):
		Write(c, http.StatusUnauthorized, "unauthorized", "Use a magic link to sign in.")
	case errors.Is(err, types.ErrUnauthorized):
		Write(c, http.StatusUnauthorized, "unauthorized", "That email and password don’t match.")
	case errors.Is(err, types.ErrForbidden):
		Write(c, http.StatusForbidden, "forbidden", "Viewers can look, not change.")
	case errors.Is(err, types.ErrNotFound):
		Write(c, http.StatusNotFound, "not_found", "We can’t find that.")
	case errors.Is(err, types.ErrTokenExpired), errors.Is(err, types.ErrTokenConsumed):
		Write(c, http.StatusUnauthorized, "unauthorized", "That link is no longer valid.")
	case errors.Is(err, types.ErrConflict):
		Write(c, http.StatusConflict, "conflict", "That didn’t save.")
	default:
		Write(c, http.StatusInternalServerError, "internal_error", "Something went wrong.")
	}
}
