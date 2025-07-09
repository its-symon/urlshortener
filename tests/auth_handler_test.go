package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/its-symon/urlshortener/internal/config"
	"github.com/its-symon/urlshortener/internal/handlers"
	"github.com/its-symon/urlshortener/internal/services"
	"github.com/stretchr/testify/assert"
)

func setupAuthRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	config.ConnectTestDB()

	tokenService := services.NewTokenService()
	authHandler := handlers.NewAuthHandler(tokenService)

	router := gin.Default()
	router.POST("/register", authHandler.Register)
	router.POST("/login", authHandler.Login)

	return router
}

func TestRegisterAndLoginHandler(t *testing.T) {
	router := setupAuthRouter()

	payload := map[string]string{
		"email":    "test@example.com",
		"password": "password123",
	}
	body, _ := json.Marshal(payload)

	// Register
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "token")

	// Login
	req = httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "token")
}
