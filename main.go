package main

import (
	"os"
	"prabandh/database"
	"prabandh/indexer"
	"prabandh/routers"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	err := godotenv.Load(".env")
	if err != nil {
		panic("Error loading .env file")
	}

	database.Connect()

	fileIndexer := indexer.FileIndexer{}
	directoryPath := os.Getenv("DATA_PATH")
	if directoryPath == "" {
		panic("DATA_PATH is not set in the environment")
	}
	fileIndexer.IndexDirectory(directoryPath)

	r := gin.Default()

	// Use routers
	routers.RegisterFileRoutes(r)
	routers.RegisterSummaryRoutes(r)
	routers.RegisterIndexDirRoutes(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port
	}
	r.Run(":" + port)
}
