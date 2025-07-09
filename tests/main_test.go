package tests

import (
	"os"
	"testing"

	"github.com/its-symon/urlshortener/internal/config"
	"github.com/its-symon/urlshortener/internal/queue"
)

func TestMain(m *testing.M) {
	config.ConnectTestDB()
	config.ConnectTestRedis()
	queue.InitRabbitMQ()

	code := m.Run()
	os.Exit(code)
}
