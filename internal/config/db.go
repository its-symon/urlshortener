package config

import (
	"fmt"
	"log"

	"github.com/its-symon/urlshortener/internal/models"
	"gorm.io/driver/postgres"
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
