package controllers

import (
	"net/http"
	"prabandh/database"
	"prabandh/models"

	"github.com/gin-gonic/gin"
)

func AddFileSummary(c *gin.Context) {
	var summary models.FileSummary
	if err := c.ShouldBindJSON(&summary); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Create(&summary).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Summary added successfully", "summary": summary})
}

func GetFileSummaries(c *gin.Context) {
	fileIndexID := c.Query("file_index_id")
	var summaries []models.FileSummary

	if err := database.DB.Where("file_index_id = ?", fileIndexID).Find(&summaries).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, summaries)
}
