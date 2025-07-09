package models

import "gorm.io/gorm"

type ApiToken struct {
	gorm.Model
	Token  string `gorm:"uniqueIndex"`
	UserID uint
}
