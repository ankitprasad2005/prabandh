package routers

import (
	"prabandh/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterIndexDirRoutes(r *gin.Engine) {
	indexDirGroup := r.Group("/index-dir")
	{
		indexDirGroup.POST("/add", controllers.AddIndexDir)
		indexDirGroup.GET("/", controllers.GetIndexDirs)
	}
}
