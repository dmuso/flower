package organisation

import (
	"net/http"

	"flower/api/internal/platform/httpx"
	"flower/api/internal/platform/middleware"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

type createBody struct {
	Name string `json:"name"`
}

func (h *Handler) Create(c *gin.Context) {
	u := middleware.CurrentUser(c)
	if u == nil {
		httpx.Write(c, http.StatusUnauthorized, "unauthorized", "You’re signed out.")
		return
	}
	var body createBody
	if err := c.ShouldBindJSON(&body); err != nil {
		httpx.Write(c, http.StatusBadRequest, "validation_failed", "Name the organisation.")
		return
	}
	org, err := h.svc.Create(c.Request.Context(), u.UserID, body.Name)
	if err != nil {
		httpx.From(c, err)
		return
	}
	c.JSON(http.StatusCreated, org)
}
