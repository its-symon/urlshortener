package middleware

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func APIKeyAuthMiddleware() gin.HandlerFunc {
	requiredAPIKey := os.Getenv("API_KEY")
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")

		if apiKey == "" || apiKey != requiredAPIKey {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or missing API key"})
			return
		}

		c.Next()
	}
}
