package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/its-symon/urlshortener/internal/config"
	"github.com/its-symon/urlshortener/internal/handlers"
	"github.com/its-symon/urlshortener/internal/middleware"
	"github.com/its-symon/urlshortener/internal/services"
	"github.com/stretchr/testify/assert"
)

func setupURLRouter() *gin.Engine {

	config.AppConfig = &config.Config{
		Port: "8080",
	}

	gin.SetMode(gin.TestMode)
	config.ConnectTestDB()
	config.ConnectTestRedis()

	tokenService := services.NewTokenService()
	authHandler := handlers.NewAuthHandler(tokenService)
	urlHandler := handlers.NewURLHandler()

	router := gin.Default()
	router.POST("/register", authHandler.Register)
	router.POST("/login", authHandler.Login)

	router.POST("/generate-api-key", middleware.JWTAuthMiddleware(tokenService), authHandler.GenerateApiKey)
	router.POST("/shorten", middleware.APIKeyAuthMiddleware(), urlHandler.Shorten)
	router.GET("/:shortCode", urlHandler.Redirect)

	return router
}

func TestRegister(t *testing.T) {
	router := setupURLRouter()

	payload := map[string]string{
		"email":    "register@example.com",
		"password": "strongpass",
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "token")
}

func TestLogin(t *testing.T) {
	router := setupURLRouter()

	payload := map[string]string{
		"email":    "login@example.com",
		"password": "mypassword",
	}
	body, _ := json.Marshal(payload)

	// Register
	registerReq := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
	registerReq.Header.Set("Content-Type", "application/json")
	registerResp := httptest.NewRecorder()
	router.ServeHTTP(registerResp, registerReq)
	assert.Equal(t, http.StatusOK, registerResp.Code)

	// Login
	loginReq := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	loginReq.Header.Set("Content-Type", "application/json")
	loginResp := httptest.NewRecorder()
	router.ServeHTTP(loginResp, loginReq)

	assert.Equal(t, http.StatusOK, loginResp.Code)
	assert.Contains(t, loginResp.Body.String(), "token")
}

func TestGenerateAPIKey(t *testing.T) {
	router := setupURLRouter()

	// Register user
	payload := map[string]string{
		"email":    "user@example.com",
		"password": "mypassword",
	}
	body, _ := json.Marshal(payload)

	registerReq := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
	registerReq.Header.Set("Content-Type", "application/json")
	registerResp := httptest.NewRecorder()
	router.ServeHTTP(registerResp, registerReq)
	assert.Equal(t, http.StatusOK, registerResp.Code)

	// Login to get JWT token
	loginReq := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	loginReq.Header.Set("Content-Type", "application/json")
	loginResp := httptest.NewRecorder()
	router.ServeHTTP(loginResp, loginReq)
	assert.Equal(t, http.StatusOK, loginResp.Code)

	var loginData map[string]string
	_ = json.Unmarshal(loginResp.Body.Bytes(), &loginData)
	token := loginData["token"]

	// Generate API key
	apiKeyReq := httptest.NewRequest(http.MethodPost, "/generate-api-key", nil)
	apiKeyReq.Header.Set("Authorization", "Bearer "+token)
	apiKeyResp := httptest.NewRecorder()
	router.ServeHTTP(apiKeyResp, apiKeyReq)

	fmt.Println("Generated API Key:", apiKeyResp.Body.String())
	assert.Equal(t, http.StatusOK, apiKeyResp.Code)
	assert.Contains(t, apiKeyResp.Body.String(), "api_key")
}

// Test shortening a URL
func TestShortenURL(t *testing.T) {
	router := setupURLRouter()

	// Register user
	payload := map[string]string{
		"email":    "user@example.com",
		"password": "mypassword",
	}
	body, _ := json.Marshal(payload)
	registerReq := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
	registerReq.Header.Set("Content-Type", "application/json")
	registerResp := httptest.NewRecorder()
	router.ServeHTTP(registerResp, registerReq)
	assert.Equal(t, http.StatusOK, registerResp.Code)

	// Login to get JWT token
	loginReq := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	loginReq.Header.Set("Content-Type", "application/json")
	loginResp := httptest.NewRecorder()
	router.ServeHTTP(loginResp, loginReq)
	assert.Equal(t, http.StatusOK, loginResp.Code)

	var loginData map[string]string
	_ = json.Unmarshal(loginResp.Body.Bytes(), &loginData)
	token := loginData["token"]

	// Generate API key
	apiKeyReq := httptest.NewRequest(http.MethodPost, "/generate-api-key", nil)
	apiKeyReq.Header.Set("Authorization", "Bearer "+token)
	apiKeyResp := httptest.NewRecorder()
	router.ServeHTTP(apiKeyResp, apiKeyReq)

	// Extract API key from response
	var apiKeyData map[string]string
	err := json.Unmarshal(apiKeyResp.Body.Bytes(), &apiKeyData)
	assert.NoError(t, err)
	apiKey := apiKeyData["api_key"]
	assert.NotEmpty(t, apiKey)

	fmt.Println("Generated API Key:", apiKey)

	// Shorten URL
	shortenPayload := map[string]string{
		"long_url": "https://www.youtube.com/watch?v=HTcL9WkB_wg&list=RDOxXKDGO-MYQ&index=18",
	}
	shortenBody, _ := json.Marshal(shortenPayload)
	shortenReq := httptest.NewRequest(http.MethodPost, "/shorten", bytes.NewBuffer(shortenBody))
	shortenReq.Header.Set("Content-Type", "application/json")
	shortenReq.Header.Set("X-API-Key", apiKey)
	shortenResp := httptest.NewRecorder()
	router.ServeHTTP(shortenResp, shortenReq)

	assert.Equal(t, http.StatusOK, shortenResp.Code)

	// Verify the response contains the shortened URL
	var shortenData map[string]interface{}
	err = json.Unmarshal(shortenResp.Body.Bytes(), &shortenData)
	assert.NoError(t, err)
	assert.Contains(t, shortenData, "short_url")
	assert.NotEmpty(t, shortenData["short_url"])

	fmt.Println("Shortened URL:", shortenData["short_url"])

}
