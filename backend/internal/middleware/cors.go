package middleware

import "github.com/gin-gonic/gin"

// CORS allows the configured origin. An empty or "*" value allows every origin,
// which is convenient for local development. Set CORS_ORIGIN to the deployed
// frontend origin to lock it down in production.
func CORS(allowOrigin string) gin.HandlerFunc {
	if allowOrigin == "" {
		allowOrigin = "*"
	}
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", allowOrigin)
		c.Header("Vary", "Origin")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
