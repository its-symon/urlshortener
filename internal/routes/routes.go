package routes

import (
	"github.com/its-symon/urlshortener/internal/handlers"
	"github.com/its-symon/urlshortener/internal/middleware"
	"github.com/its-symon/urlshortener/internal/services"

	"github.com/gin-gonic/gin"

	_ "github.com/its-symon/urlshortener/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func RegisterRoutes(r *gin.Engine) {
	tokenService := services.NewTokenService()
	authHandler := handlers.NewAuthHandler(tokenService)
	urlHandler := handlers.NewURLHandler()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Welcome to the URL Shortener API"})
	})

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy"})
	})

	// Public auth routes
	r.POST("/login", authHandler.Login)
	r.POST("/register", authHandler.Register)

	// Protected route to generate API key (JWT required)
	r.POST("/generate-api-key", middleware.JWTAuthMiddleware(tokenService), authHandler.GenerateApiKey)

	// Protected route using API key
	r.POST("/shorten", middleware.APIKeyAuthMiddleware(), urlHandler.Shorten)

	// Public routes
	r.GET("/:shortCode", urlHandler.Redirect)

	r.GET("/details/:shortCode", urlHandler.GetDetails)

	r.DELETE("delete/:shortCode", urlHandler.Delete)

	// Swagger UI route
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
