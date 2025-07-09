package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/its-symon/urlshortener/internal/config"
	"github.com/its-symon/urlshortener/internal/models"
	"github.com/its-symon/urlshortener/internal/queue"
	"github.com/its-symon/urlshortener/internal/repositories"
	"github.com/its-symon/urlshortener/internal/services"

	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"
)

type URLHandler struct {
	Service *services.URLService
}

func NewURLHandler() *URLHandler {
	repo := repositories.NewURLRepository()
	service := services.NewURLService(repo)
	return &URLHandler{Service: service}
}

func (h *URLHandler) Shorten(c *gin.Context) {
	// 1. Rate limit per IP
	ip := c.ClientIP()
	key := fmt.Sprintf("rate_limit:%s", ip)

	countStr, _ := config.RedisClient.Get(config.RedisCtx, key).Result()
	count, _ := strconv.Atoi(countStr)
	if count >= 10 {
		c.JSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit exceeded. Max 10 requests per minute."})
		return
	}

	// Increment and set expiry (1 min)
	_ = config.RedisClient.Incr(config.RedisCtx, key)
	_ = config.RedisClient.Expire(config.RedisCtx, key, time.Minute*1)

	// 2. Handle URL shortening
	var req models.ShortenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.Service.ShortenURL(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *URLHandler) Redirect(c *gin.Context) {
	shortCode := c.Param("shortCode")

	longURL, err := h.Service.GetOriginalURL(shortCode)
	if err != nil {
		if err.Error() == "URL has expired" {
			c.JSON(http.StatusGone, gin.H{"error": "URL has expired"})
		} else {
			c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
		}
		return
	}

	// Push shortCode to RabbitMQ for async processing
	_ = queue.Channel.Publish(
		"",             // exchange
		"click_events", // routing key (queue name)
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(shortCode),
		},
	)

	// Redirect
	if !strings.HasPrefix(longURL, "http") {
		longURL = "https://" + longURL
	}

	c.Redirect(http.StatusFound, longURL)
}

func (h *URLHandler) GetDetails(c *gin.Context) {
	shortCode := c.Param("shortCode")

	res, err := h.Service.GetURLDetails(shortCode)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *URLHandler) Delete(c *gin.Context) {
	shortCode := c.Param("shortCode")

	err := h.Service.DeleteURL(shortCode)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Short URL deleted successfully"})
}
