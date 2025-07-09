package routes

import (
	"github.com/its-symon/urlshortener/internal/handlers"

	"github.com/gin-gonic/gin"

	_ "github.com/its-symon/urlshortener/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func RegisterRoutes(r *gin.Engine) {
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

	r.POST("/shorten", urlHandler.Shorten)

	r.GET("/:shortCode", urlHandler.Redirect)

	r.GET("/details/:shortCode", urlHandler.GetDetails)

	r.DELETE("delete/:shortCode", urlHandler.Delete)

	// Swagger UI route
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
