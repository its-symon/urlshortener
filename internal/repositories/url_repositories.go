package repositories

import (
	"github.com/its-symon/urlshortener/internal/config"
	"github.com/its-symon/urlshortener/internal/models"
	"gorm.io/gorm"
)

type URLRepository interface {
	Create(url *models.URLMapping) error
	ExistsByShortCode(shortCode string) (bool, error)
	GetByShortCode(shortCode string) (*models.URLMapping, error)
	GetDetailsByShortCode(shortCode string) (*models.URLMapping, error)
	IncrementClickCount(shortCode string) error
	DeleteByShortCode(shortCode string) error
}

type urlRepository struct{}

func NewURLRepository() URLRepository {
	return &urlRepository{}
}

func (r *urlRepository) Create(url *models.URLMapping) error {
	return config.DB.Create(url).Error
}

func (r *urlRepository) ExistsByShortCode(shortCode string) (bool, error) {
	var count int64
	err := config.DB.Model(&models.URLMapping{}).Where("short_code = ?", shortCode).Count(&count).Error
	return count > 0, err
}

func (r *urlRepository) GetByShortCode(shortCode string) (*models.URLMapping, error) {
	var url models.URLMapping
	err := config.DB.Where("short_code = ? AND is_deleted = false", shortCode).First(&url).Error
	if err != nil {
		return nil, err
	}
	return &url, nil
}

func (r *urlRepository) GetDetailsByShortCode(shortCode string) (*models.URLMapping, error) {
	var url models.URLMapping
	err := config.DB.Where("short_code = ? AND is_deleted = false", shortCode).First(&url).Error
	if err != nil {
		return nil, err
	}
	return &url, nil
}

func (r *urlRepository) IncrementClickCount(shortCode string) error {
	return config.DB.Model(&models.URLMapping{}).
		Where("short_code = ?", shortCode).
		UpdateColumn("click_count", gorm.Expr("click_count + ?", 1)).
		Error
}

func (r *urlRepository) DeleteByShortCode(shortCode string) error {
	return config.DB.Unscoped().Where("short_code = ?", shortCode).Delete(&models.URLMapping{}).Error
}
