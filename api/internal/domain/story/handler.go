package story

import (
	"net/http"

	"flower/api/internal/platform/httpx"
	"flower/api/internal/platform/middleware"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc      *Service
	remember func(c *gin.Context, projectID string) error
}

func NewHandler(svc *Service, remember func(c *gin.Context, projectID string) error) *Handler {
	return &Handler{svc: svc, remember: remember}
}

func (h *Handler) List(c *gin.Context) {
	u := middleware.CurrentUser(c)
	if u == nil {
		httpx.Write(c, http.StatusUnauthorized, "unauthorized", "You’re signed out.")
		return
	}
	list, err := h.svc.List(c.Request.Context(), u.UserID, c.Param("id"))
	if err != nil {
		httpx.From(c, err)
		return
	}
	if h.remember != nil {
		if err := h.remember(c, c.Param("id")); err != nil {
			httpx.From(c, err)
			return
		}
	}
	c.JSON(http.StatusOK, list)
}
