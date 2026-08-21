package handlers

import (
	"fmt"

	"flower/api/internal/domain/organisation"
	"flower/api/internal/domain/project"
	"flower/api/internal/domain/story"
	"flower/api/internal/domain/user"
	"flower/api/internal/platform/middleware"

	"github.com/gin-gonic/gin"
)

type Dependencies struct {
	Version       string
	Pinger        DatabasePinger
	Users         *user.Handler
	Organisations *organisation.Handler
	Projects      *project.Handler
	Stories       *story.Handler
}

func SetupRoutes(router *gin.Engine, deps *Dependencies) error {
	if router == nil {
		return fmt.Errorf("router is required")
	}
	if deps == nil {
		return fmt.Errorf("dependencies are required")
	}
	if deps.Version == "" {
		return fmt.Errorf("version is required")
	}
	if deps.Users == nil {
		return fmt.Errorf("users handler is required")
	}
	if deps.Organisations == nil {
		return fmt.Errorf("organisations handler is required")
	}
	if deps.Projects == nil {
		return fmt.Errorf("projects handler is required")
	}
	if deps.Stories == nil {
		return fmt.Errorf("stories handler is required")
	}

	router.GET("/health", HealthHandler(deps.Version))
	router.GET("/ready", ReadyHandler(deps.Pinger, deps.Version))
	router.GET("/api/version", VersionHandler(deps.Version))

	v1 := router.Group("/api/v1")
	v1.POST("/auth/signup", deps.Users.SignUp)
	v1.POST("/auth/login", deps.Users.Login)
	v1.POST("/auth/magic-link", deps.Users.MagicLink)
	v1.POST("/auth/magic-link/consume", deps.Users.ConsumeMagicLink)
	v1.POST("/auth/verify-email/consume", deps.Users.ConsumeVerifyEmail)
	v1.POST("/auth/verify-email", middleware.RequireUser(), deps.Users.ResendVerifyEmail)
	v1.POST("/auth/logout", deps.Users.Logout)
	v1.GET("/me", middleware.RequireUser(), deps.Users.Me)
	v1.POST("/organisations", middleware.RequireUser(), deps.Organisations.Create)
	v1.POST("/organisations/:id/projects", middleware.RequireUser(), deps.Projects.Create)
	v1.GET("/organisations/:id/projects", middleware.RequireUser(), deps.Projects.List)
	v1.GET("/projects/:id", middleware.RequireUser(), deps.Projects.Get)
	v1.GET("/projects/:id/stories", middleware.RequireUser(), deps.Stories.List)
	return nil
}
