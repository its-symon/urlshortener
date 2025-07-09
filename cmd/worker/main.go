package main

import (
	"github.com/its-symon/urlshortener/internal/config"
	"github.com/its-symon/urlshortener/internal/queue"
	"github.com/its-symon/urlshortener/internal/workers"
)

func main() {
	config.LoadConfig()
	config.ConnectDatabase()
	queue.InitRabbitMQ()
	defer queue.Conn.Close()
	defer queue.Channel.Close()

	workers.StartClickWorker()
}
