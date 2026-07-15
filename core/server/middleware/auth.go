package middleware

import (
	"net/http"

	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mmycin/goforge/internal/config"
)

// ValidateAppKey checks if the provided key matches the configured App Key
func ValidateAppKey(key string) bool {
	return key != "" && key == config.App.Key
}

// AppKey checks for X-App-Key header and validates it against config
func AppKey() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Allow public access to documentation
		if strings.HasPrefix(c.Request.URL.Path, "/api/docs") || strings.HasPrefix(c.Request.URL.Path, "/openapi") || strings.HasPrefix(c.Request.URL.Path, "/docs") {
			c.Next()
			return
		}

		key := c.GetHeader("X-App-Key")

		if !ValidateAppKey(key) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Invalid or missing App Key",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
