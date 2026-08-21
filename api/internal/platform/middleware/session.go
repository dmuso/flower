package middleware

import (
	"context"
	"net/http"

	"flower/api/internal/platform/httpx"
	"flower/api/internal/types"

	"github.com/gin-gonic/gin"
)

const (
	CookieName     = "flower_session"
	ContextUserKey = "flower_user"
)

type SessionLookup func(ctx context.Context, tokenHash string) (*types.SessionUser, error)

func Session(lookup SessionLookup, hash func(string) string) gin.HandlerFunc {
	return func(c *gin.Context) {
		raw, err := c.Cookie(CookieName)
		if err != nil || raw == "" {
			c.Next()
			return
		}
		user, err := lookup(c.Request.Context(), hash(raw))
		if err != nil {
			c.Next()
			return
		}
		c.Set(CookieName+"_raw", raw)
		c.Set(ContextUserKey, user)
		c.Next()
	}
}

func CurrentUser(c *gin.Context) *types.SessionUser {
	v, ok := c.Get(ContextUserKey)
	if !ok {
		return nil
	}
	user, ok := v.(*types.SessionUser)
	if !ok || user == nil || user.UserID == "" {
		return nil
	}
	return user
}

func RequireUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		if CurrentUser(c) == nil {
			httpx.Write(c, http.StatusUnauthorized, "unauthorized", "You’re signed out.")
			c.Abort()
			return
		}
		c.Next()
	}
}
