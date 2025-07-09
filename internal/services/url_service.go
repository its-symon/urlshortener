package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/its-symon/urlshortener/internal/config"
	"github.com/its-symon/urlshortener/internal/models"
	"github.com/its-symon/urlshortener/internal/repositories"
	"github.com/its-symon/urlshortener/internal/utils"
)

type URLService struct {
	Repo repositories.URLRepository
}

func NewURLService(repo repositories.URLRepository) *URLService {
	return &URLService{Repo: repo}
}

func (s *URLService) ShortenURL(req models.ShortenRequest) (*models.ShortenResponse, error) {
	// 1. Use custom alias if provided
	shortCode := ""
	if req.CustomAlias != nil {
		shortCode = *req.CustomAlias
		exists, _ := s.Repo.ExistsByShortCode(shortCode)
		if exists {
			return nil, errors.New("custom alias already exists")
		}
	} else {
		// 2. Generate unique short code
		for {
			shortCode = utils.GenerateShortCode(6)
			exists, _ := s.Repo.ExistsByShortCode(shortCode)
			if !exists {
				break
			}
		}
	}

	// 3. Create DB record
	url := models.URLMapping{
		LongURL:   req.LongURL,
		ShortCode: shortCode,
		CreatedAt: time.Now(),
		ExpiresAt: req.ExpiresAt,
	}

	if req.CustomAlias != nil {
		url.CustomAlias = req.CustomAlias
	}

	err := s.Repo.Create(&url)
	if err != nil {
		return nil, err
	}

	shortURL := fmt.Sprintf("http://localhost:%s/%s", config.AppConfig.Port, shortCode)

	res := &models.ShortenResponse{
		ShortCode: shortCode,
		ShortURL:  shortURL,
		LongURL:   req.LongURL,
		CreatedAt: url.CreatedAt,
		ExpiresAt: url.ExpiresAt,
	}

	return res, nil
}

func (s *URLService) GetOriginalURL(shortCode string) (string, error) {
	// First check Redis
	longURL, err := config.RedisClient.Get(config.RedisCtx, shortCode).Result()
	if err == nil {
		// Found in cache
		_ = s.Repo.IncrementClickCount(shortCode) // count still updated
		fmt.Println("Cache hit for short code:", shortCode)
		return longURL, nil
	}

	// If not found in Redis, hit the DB
	url, err := s.Repo.GetByShortCode(shortCode)
	if err != nil {
		return "", err
	}

	if url.ExpiresAt != nil && url.ExpiresAt.Before(time.Now()) {
		return "", errors.New("URL has expired")
	}

	// Save to Redis for next time
	_ = config.RedisClient.Set(config.RedisCtx, shortCode, url.LongURL, time.Hour*24).Err()

	_ = s.Repo.IncrementClickCount(shortCode)

	return url.LongURL, nil
}

func (s *URLService) GetURLDetails(shortCode string) (*models.ShortenResponse, error) {
	url, err := s.Repo.GetDetailsByShortCode(shortCode)
	if err != nil {
		return nil, err
	}

	shortURL := fmt.Sprintf("http://localhost:%s/%s", config.AppConfig.Port, shortCode)

	return &models.ShortenResponse{
		ShortCode:  shortCode,
		ShortURL:   shortURL,
		LongURL:    url.LongURL,
		CreatedAt:  url.CreatedAt,
		ExpiresAt:  url.ExpiresAt,
		ClickCount: url.ClickCount,
	}, nil

}

func (s *URLService) DeleteURL(shortCode string) error {
	return s.Repo.DeleteByShortCode(shortCode)
}
