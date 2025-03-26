package routers

import (
	"prabandh/controllers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterSummaryRoutes(r *gin.Engine, db *gorm.DB, ollamaURL string) {
	summaryController := controllers.NewSummaryController(db, ollamaURL)

	summaryGroup := r.Group("/summary")
	{
		summaryGroup.POST("/add", summaryController.AddFileSummary)
		summaryGroup.GET("/", summaryController.GetFileSummaries) // Assuming you implement this method
	}
}
