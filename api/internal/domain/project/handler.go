package project

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
		httpx.Write(c, http.StatusBadRequest, "validation_failed", "Name the project.")
		return
	}
	p, err := h.svc.Create(c.Request.Context(), u.UserID, c.Param("id"), body.Name)
	if err != nil {
		httpx.From(c, err)
		return
	}
	if h.remember != nil {
		if err := h.remember(c, p.ID); err != nil {
			httpx.From(c, err)
			return
		}
	}
	c.JSON(http.StatusCreated, p)
}

func (h *Handler) List(c *gin.Context) {
	u := middleware.CurrentUser(c)
	if u == nil {
		httpx.Write(c, http.StatusUnauthorized, "unauthorized", "You’re signed out.")
		return
	}
	projects, err := h.svc.List(c.Request.Context(), u.UserID, c.Param("id"))
	if err != nil {
		httpx.From(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"projects": projects})
}

func (h *Handler) Get(c *gin.Context) {
	u := middleware.CurrentUser(c)
	if u == nil {
		httpx.Write(c, http.StatusUnauthorized, "unauthorized", "You’re signed out.")
		return
	}
	p, err := h.svc.Get(c.Request.Context(), u.UserID, c.Param("id"))
	if err != nil {
		httpx.From(c, err)
		return
	}
	if h.remember != nil {
		if err := h.remember(c, p.ID); err != nil {
			httpx.From(c, err)
			return
		}
	}
	c.JSON(http.StatusOK, p)
}
