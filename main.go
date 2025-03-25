package main

import (
	"prabandh/database"
	"prabandh/indexer"
	"prabandh/routers"

	"github.com/gin-gonic/gin"
)

func main() {
	database.Connect()

	fileIndexer := indexer.FileIndexer{}
	fileIndexer.IndexDirectory("/path/to/directory")

	r := gin.Default()

	// Use routers
	routers.RegisterFileRoutes(r)
	routers.RegisterSummaryRoutes(r)

	r.Run(":8080")
}
