package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/mmycin/goforge/core/config"
)

// CORS returns a middleware for handling CORS
func CORS() gin.HandlerFunc {
	allowedOrigins := config.HTTP.CORSAllowedOrigins
	if allowedOrigins == "" {
		allowedOrigins = "*"
	}
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", allowedOrigins)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
