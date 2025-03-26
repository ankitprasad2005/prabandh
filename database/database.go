package database

import (
	"fmt"
	"log"
	"os"

	"prabandh/models"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	// Load environment variables
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// var err error
	dsn := os.Getenv("DB_URL")
	if dsn == "" {
		log.Fatal("DB_URL is not set in the environment")
	}
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	err = DB.AutoMigrate(&models.FileIndex{}, &models.FileSummary{}, &models.IndexDir{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Create index for keyword search
	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_summary_keyword_search ON file_summaries USING gin(to_tsvector('english', summary_keyword))").Error; err != nil {
		log.Printf("Warning: Could not create full-text search index: %v", err)
	}

	fmt.Println("Database connection established and migrated")
}
