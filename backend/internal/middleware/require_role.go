package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		r, _ := c.Get("role")
		for _, x := range roles {
			if r == x {
				c.Next()
				return
			}
		}
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
	}
}
