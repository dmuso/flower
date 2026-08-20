package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func VersionHandler(version string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Cache-Control", "no-store, no-cache, must-revalidate")
		c.Header("Pragma", "no-cache")
		c.Header("Expires", "0")
		c.JSON(http.StatusOK, gin.H{
			"name":    "flower",
			"version": version,
		})
	}
}
