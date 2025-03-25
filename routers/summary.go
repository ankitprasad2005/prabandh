package routers

import (
	"prabandh/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterSummaryRoutes(r *gin.Engine) {
	summaryGroup := r.Group("/summary")
	{
		summaryGroup.POST("/add", controllers.AddFileSummary)
		summaryGroup.GET("/", controllers.GetFileSummaries)
	}
}
