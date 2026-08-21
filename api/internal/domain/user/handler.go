package user

import (
	"net/http"
	"time"

	"flower/api/internal/platform/httpx"
	"flower/api/internal/platform/middleware"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc    *Service
	secure bool
}

func NewHandler(svc *Service, cookieSecure bool) *Handler {
	return &Handler{svc: svc, secure: cookieSecure}
}

type emailPasswordBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type emailBody struct {
	Email string `json:"email"`
}

type tokenBody struct {
	Token string `json:"token"`
}

func (h *Handler) SignUp(c *gin.Context) {
	var body emailPasswordBody
	if err := c.ShouldBindJSON(&body); err != nil {
		httpx.Write(c, http.StatusBadRequest, "validation_failed", "Need an email and a password.")
		return
	}
	if err := h.svc.SignUp(c.Request.Context(), body.Email, body.Password); err != nil {
		httpx.From(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"email": body.Email})
}

func (h *Handler) Login(c *gin.Context) {
	var body emailPasswordBody
	if err := c.ShouldBindJSON(&body); err != nil {
		httpx.Write(c, http.StatusUnauthorized, "unauthorized", "That email and password don’t match.")
		return
	}
	raw, expires, me, err := h.svc.Login(c.Request.Context(), body.Email, body.Password)
	if err != nil {
		httpx.From(c, err)
		return
	}
	h.setCookie(c, raw, expires)
	c.JSON(http.StatusOK, me)
}

func (h *Handler) MagicLink(c *gin.Context) {
	var body emailBody
	if err := c.ShouldBindJSON(&body); err != nil {
		httpx.Write(c, http.StatusBadRequest, "validation_failed", "Need an email.")
		return
	}
	if err := h.svc.RequestMagicLink(c.Request.Context(), body.Email); err != nil {
		httpx.From(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"email": body.Email})
}

func (h *Handler) ConsumeMagicLink(c *gin.Context) {
	var body tokenBody
	if err := c.ShouldBindJSON(&body); err != nil {
		httpx.From(c, err)
		return
	}
	raw, expires, me, err := h.svc.ConsumeMagicLink(c.Request.Context(), body.Token)
	if err != nil {
		httpx.From(c, err)
		return
	}
	h.setCookie(c, raw, expires)
	c.JSON(http.StatusOK, me)
}

func (h *Handler) ConsumeVerifyEmail(c *gin.Context) {
	var body tokenBody
	if err := c.ShouldBindJSON(&body); err != nil {
		httpx.From(c, err)
		return
	}
	raw, expires, me, err := h.svc.ConsumeVerifyEmail(c.Request.Context(), body.Token)
	if err != nil {
		httpx.From(c, err)
		return
	}
	h.setCookie(c, raw, expires)
	c.JSON(http.StatusOK, me)
}

func (h *Handler) ResendVerifyEmail(c *gin.Context) {
	u := middleware.CurrentUser(c)
	if u == nil {
		httpx.Write(c, http.StatusUnauthorized, "unauthorized", "You’re signed out.")
		return
	}
	if err := h.svc.ResendVerifyEmail(c.Request.Context(), u.UserID); err != nil {
		httpx.From(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func (h *Handler) Logout(c *gin.Context) {
	if u := middleware.CurrentUser(c); u != nil {
		if err := h.svc.Logout(c.Request.Context(), u.SessionID); err != nil {
			httpx.From(c, err)
			return
		}
	}
	h.clearCookie(c)
	c.Status(http.StatusNoContent)
}

func (h *Handler) Me(c *gin.Context) {
	u := middleware.CurrentUser(c)
	if u == nil {
		httpx.Write(c, http.StatusUnauthorized, "unauthorized", "You’re signed out.")
		return
	}
	me, err := h.svc.Me(c.Request.Context(), u.UserID)
	if err != nil {
		httpx.From(c, err)
		return
	}
	c.JSON(http.StatusOK, me)
}

func (h *Handler) setCookie(c *gin.Context, raw string, expires time.Time) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(middleware.CookieName, raw, int(time.Until(expires).Seconds()), "/", "", h.secure, true)
}

func (h *Handler) clearCookie(c *gin.Context) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(middleware.CookieName, "", -1, "/", "", h.secure, true)
}
