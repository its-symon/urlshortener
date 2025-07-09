package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/its-symon/urlshortener/internal/config"
	"github.com/its-symon/urlshortener/internal/models"
	"github.com/its-symon/urlshortener/internal/routes"
)

func main() {
	// Load config
	config.LoadConfig()

	// Connect to DB
	config.ConnectDatabase()

	// Connect to Redis
	config.InitRedis()

	// Auto-migrate URLMapping model
	config.DB.AutoMigrate(&models.URLMapping{})

	fmt.Println("PORT:", config.AppConfig.Port)

	r := gin.Default()
	routes.RegisterRoutes(r)
	r.Run(":" + config.AppConfig.Port)
}
