package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/its-symon/urlshortener/internal/config"
	"github.com/its-symon/urlshortener/internal/models"
)

func APIKeyAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("x-api-key")
		if apiKey == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing API key"})
			return
		}

		var tokenRecord models.ApiToken
		if err := config.DB.Where("token = ?", apiKey).First(&tokenRecord).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid API key"})
			return
		}

		// Optionally, load user or token details to context if needed
		c.Next()
	}
}
