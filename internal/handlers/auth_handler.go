package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/its-symon/urlshortener/internal/services"
)

type AuthHandler struct {
	TokenService *services.TokenService
}

func NewAuthHandler(tokenService *services.TokenService) *AuthHandler {
	return &AuthHandler{TokenService: tokenService}
}

type TokenRequest struct {
	Username string `json:"username" binding:"required"`
}

func (h *AuthHandler) GenerateToken(c *gin.Context) {
	var req TokenRequest

	// Default to "guest" if no username provided
	username := req.Username
	if username == "" {
		username = "guest"
	}

	// Generate JWT token
	token, err := h.TokenService.GenerateToken(username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Return token in response
	c.JSON(http.StatusOK, gin.H{"token": token})
}
