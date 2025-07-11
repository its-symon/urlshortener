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

func TestRegisterSuccess(t *testing.T) {
	router := setupAuthRouter()

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

func TestLoginSuccess(t *testing.T) {
	router := setupAuthRouter()

	// First register
	payload := map[string]string{
		"email":    "login@example.com",
		"password": "mypassword",
	}
	body, _ := json.Marshal(payload)

	registerReq := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
	registerReq.Header.Set("Content-Type", "application/json")
	registerResp := httptest.NewRecorder()
	router.ServeHTTP(registerResp, registerReq)
	assert.Equal(t, http.StatusOK, registerResp.Code)

	// Then login
	loginReq := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	loginReq.Header.Set("Content-Type", "application/json")
	loginResp := httptest.NewRecorder()
	router.ServeHTTP(loginResp, loginReq)

	assert.Equal(t, http.StatusOK, loginResp.Code)
	assert.Contains(t, loginResp.Body.String(), "token")
}

func TestLoginWrongPassword(t *testing.T) {
	router := setupAuthRouter()

	// Register user
	payload := map[string]string{
		"email":    "wrongpass@example.com",
		"password": "correctpass",
	}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Login with wrong password
	wrongPayload := map[string]string{
		"email":    "wrongpass@example.com",
		"password": "wrongpass",
	}
	wrongBody, _ := json.Marshal(wrongPayload)
	loginReq := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(wrongBody))
	loginReq.Header.Set("Content-Type", "application/json")
	loginResp := httptest.NewRecorder()
	router.ServeHTTP(loginResp, loginReq)

	assert.Equal(t, http.StatusUnauthorized, loginResp.Code)
	assert.Contains(t, loginResp.Body.String(), "Invalid credentials")
}

func TestRegisterDuplicateEmail(t *testing.T) {
	router := setupAuthRouter()

	payload := map[string]string{
		"email":    "duplicate@example.com",
		"password": "somepass",
	}
	body, _ := json.Marshal(payload)

	// First registration
	req1 := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
	req1.Header.Set("Content-Type", "application/json")
	w1 := httptest.NewRecorder()
	router.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusOK, w1.Code)

	// Second registration with same email
	req2 := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	assert.Equal(t, 409, w2.Code)                                    // match status
	assert.Contains(t, w2.Body.String(), "Email already registered") // match message

}
