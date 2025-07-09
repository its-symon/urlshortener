package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/its-symon/urlshortener/internal/config"
	"github.com/its-symon/urlshortener/internal/models"
	"github.com/its-symon/urlshortener/internal/utils"
)

func (h *AuthHandler) GenerateApiKey(c *gin.Context) {
	emailAny, exists := c.Get("email")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	email := emailAny.(string)

	// Find user by email
	var user models.User
	if err := config.DB.Where("email = ?", email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	// Generate random API key
	apiKey, err := utils.GenerateRandomToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate API key"})
		return
	}

	// Save API key linked to user
	tokenRecord := models.ApiToken{
		Token:  apiKey,
		UserID: user.ID,
	}

	if err := config.DB.Create(&tokenRecord).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save API key"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"api_key": apiKey})
}
