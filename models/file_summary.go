package models

import (
	"gorm.io/gorm"
)

type FileSummary struct {
	gorm.Model
	FileIndexID    uint   `gorm:"not null;index"` // Foreign key linking to FileIndex
	SummaryKeyword string `gorm:"not null;index"` // Each keyword gets its own row
}

// Index on FileIndexID and SummaryKeyword for faster searches
