package config

import (
	"fmt"
	"log"

	"github.com/its-symon/urlshortener/internal/models"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		AppConfig.DBHost,
		AppConfig.DBUser,
		AppConfig.DBPassword,
		AppConfig.DBName,
		AppConfig.DBPort,
	)

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	DB = database
	fmt.Println("Connected to DB")

	// Auto-migrate models (add your models here)
	err = DB.AutoMigrate(
		&models.User{},
		&models.ApiToken{},
	)
	if err != nil {
		log.Fatal("Failed DB migration:", err)
	}
}

func ConnectTestDB() {
	database, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to test database:", err)
	}
	DB = database

	DB.AutoMigrate(&models.User{}, &models.ApiToken{}, &models.URLMapping{})
}

func ConnectTestRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	_, err := RedisClient.Ping(RedisCtx).Result()
	if err != nil {
		log.Fatal("Failed to connect to Redis in tests:", err)
	}
}
