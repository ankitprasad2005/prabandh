package routers

import (
	"prabandh/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterFileRoutes(r *gin.Engine) {
	fileGroup := r.Group("/file")
	{
		fileGroup.POST("/add", controllers.AddFile)
		fileGroup.GET("/search", controllers.SearchFiles)
	}
}
