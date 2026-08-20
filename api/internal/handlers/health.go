package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type DatabasePinger interface {
	Ping() error
}

func HealthHandler(version string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Cache-Control", "no-store, no-cache, must-revalidate")
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"version": version,
		})
	}
}

func ReadyHandler(pinger DatabasePinger, version string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Cache-Control", "no-store, no-cache, must-revalidate")
		if pinger == nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status":  "unavailable",
				"version": version,
				"error":   "database pinger is not configured",
			})
			return
		}
		if err := pinger.Ping(); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status":  "unavailable",
				"version": version,
				"error":   "database ping failed",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status":  "ready",
			"version": version,
		})
	}
}
