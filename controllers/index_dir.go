package controllers

import (
	"net/http"
	"prabandh/database"
	"prabandh/models"

	"github.com/gin-gonic/gin"
)

func AddIndexDir(c *gin.Context) {
	var indexDir models.IndexDir
	if err := c.ShouldBindJSON(&indexDir); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Create(&indexDir).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Index directory added successfully", "index_dir": indexDir})
}

func GetIndexDirs(c *gin.Context) {
	var indexDirs []models.IndexDir
	if err := database.DB.Find(&indexDirs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, indexDirs)
}
