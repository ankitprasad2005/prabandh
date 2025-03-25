package models

import (
	"time"

	"gorm.io/gorm"
)

type FileIndex struct {
	gorm.Model
	FilePath     string    `gorm:"unique;not null"`
	FileName     string    `gorm:"not null"`
	Extension    string    `gorm:"not null"`
	CreatedDate  time.Time `gorm:"not null"`
	ModifiedDate time.Time `gorm:"not null"`
	Size         int64     `gorm:"not null"`
	Hash         string    `gorm:"not null"`
}
