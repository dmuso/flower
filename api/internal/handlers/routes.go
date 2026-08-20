package handlers

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type Dependencies struct {
	Version string
	Pinger  DatabasePinger
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

	router.GET("/health", HealthHandler(deps.Version))
	router.GET("/ready", ReadyHandler(deps.Pinger, deps.Version))
	router.GET("/api/version", VersionHandler(deps.Version))
	return nil
}
