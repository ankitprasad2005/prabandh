package models

import (
	"gorm.io/gorm"
)

type FileSummary struct {
	gorm.Model
	FileIndexID    uint      `gorm:"not null"` // Foreign key linking to FileIndex
	SummaryKeyword string    `gorm:"not null"`
	FileIndex      FileIndex `gorm:"foreignKey:FileIndexID;constraint:OnDelete:CASCADE"`
}
