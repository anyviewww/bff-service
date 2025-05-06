package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware represents middleware for authentication.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		// Check the token format (Bearer <token>)
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
			c.Abort()
			return
		}

		// If the token is valid, continue executing the request
		c.Set("user_token", token)
		c.Next()
	}
}
