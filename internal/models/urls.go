package models

import (
	"time"
)

type URLMapping struct {
	ID          uint    `gorm:"primaryKey"`
	LongURL     string  `gorm:"not null"`
	ShortCode   string  `gorm:"uniqueIndex;not null"`
	CustomAlias *string `gorm:"uniqueIndex"`
	ClickCount  int     `json:"click_count"`
	CreatedAt   time.Time
	ExpiresAt   *time.Time
	IsDeleted   bool `gorm:"default:false"`
}

type ShortenRequest struct {
	LongURL     string     `json:"long_url" binding:"required"`
	CustomAlias *string    `json:"custom_alias,omitempty"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
}

type ShortenResponse struct {
	ShortCode  string     `json:"short_code"`
	ShortURL   string     `json:"short_url"`
	LongURL    string     `json:"long_url"`
	ExpiresAt  *time.Time `json:"expires_at,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	ClickCount int        `json:"click_count"`
}
