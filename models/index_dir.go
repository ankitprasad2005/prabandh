package models

import (
	"gorm.io/gorm"
)

type IndexDir struct {
	gorm.Model
	DirectoryLocation string `gorm:"not null;unique"`
	IsWhitelisted     bool   `gorm:"default:true"` // Default to whitelisted
}
