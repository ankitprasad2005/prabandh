package controllers

import (
	"net/http"
	"prabandh/database"
	"prabandh/models"

	"github.com/gin-gonic/gin"
)

func AddFile(c *gin.Context) {
	var file models.FileIndex
	if err := c.ShouldBindJSON(&file); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Create(&file).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File added successfully", "file": file})
}

func SearchFiles(c *gin.Context) {
	query := c.Query("query")
	var files []models.FileIndex

	if err := database.DB.Where("file_path LIKE ?", "%"+query+"%").Find(&files).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, files)
}
