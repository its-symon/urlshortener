package workers

import (
	"log"

	"github.com/its-symon/urlshortener/internal/config"
	"github.com/its-symon/urlshortener/internal/models"
	"github.com/its-symon/urlshortener/internal/queue"
	"gorm.io/gorm"
)

func StartClickWorker() {
	msgs, err := queue.Channel.Consume(
		"click_events",
		"",
		true, // auto-acknowledge
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to register consumer: %v", err)
	}

	log.Println("Click Worker started...")
	for msg := range msgs {
		shortCode := string(msg.Body)
		err := config.DB.Model(&models.URLMapping{}).
			Where("short_code = ?", shortCode).
			UpdateColumn("click_count", gorm.Expr("click_count + ?", 1)).Error
		if err != nil {
			log.Printf("Failed to update click count for %s: %v", shortCode, err)
		}
	}
}
