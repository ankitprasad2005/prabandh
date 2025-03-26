package database

import (
	"fmt"
	"log"
	"os"

	"prabandh/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	var err error
	dsn := os.Getenv("DATABASE_URL")
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Migrate models
	err = DB.AutoMigrate(&models.FileIndex{}, &models.FileSummary{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Create index for keyword search
	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_summary_keyword_search ON file_summaries USING gin(to_tsvector('english', summary_keyword))").Error; err != nil {
		log.Printf("Warning: Could not create full-text search index: %v", err)
	}

	fmt.Println("Database connection established and migrated")
}
