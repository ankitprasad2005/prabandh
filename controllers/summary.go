package controllers

import (
	"net/http"
	"strings"
	"time"

	"prabandh/llm/ollama"
	"prabandh/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SummaryController struct {
	db           *gorm.DB
	ollamaClient *ollama.Client
}

func NewSummaryController(db *gorm.DB, ollamaURL string) *SummaryController {
	return &SummaryController{
		db:           db,
		ollamaClient: ollama.New(ollamaURL, "gemma:2b"),
	}
}

func (sc *SummaryController) AddFileSummary(c *gin.Context) {
	var input struct {
		FileIndexID uint   `json:"file_index_id" binding:"required"`
		Content     string `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	keywords, err := sc.generateKeywords(input.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Keyword generation failed",
			"details": err.Error(),
		})
		return
	}

	summary := models.FileSummary{
		FileIndexID:    input.FileIndexID,
		SummaryKeyword: strings.Join(keywords, ","),
	}

	if err := sc.db.Create(&summary).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Summary created successfully",
		"summary": summary,
	})
}

func (sc *SummaryController) GetFileSummaries(c *gin.Context) {
	fileName := c.Query("file_name")
	if fileName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file_name query parameter is required"})
		return
	}

	var fileIndex models.FileIndex
	if err := sc.db.Where("file_name = ?", fileName).First(&fileIndex).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	var summaries []models.FileSummary
	if err := sc.db.Where("file_index_id = ?", fileIndex.ID).Find(&summaries).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve summaries"})
		return
	}

	keywords := make([]string, len(summaries))
	for i, summary := range summaries {
		keywords[i] = summary.SummaryKeyword
	}

	c.JSON(http.StatusOK, gin.H{
		"file_name": fileName,
		"keywords":  keywords,
	})
}

func (sc *SummaryController) generateKeywords(text string) ([]string, error) {
	const maxRetries = 3
	var keywords []string
	var err error

	prompt := `Extract 5-7 most relevant keywords from this text.
Return ONLY lowercase, comma-separated terms.
Example: "ai,healthcare,data analysis"

Text: ` + text

	for i := 0; i < maxRetries; i++ {
		keywords, err = sc.ollamaClient.ExtractKeywords(prompt)
		if err == nil {
			break
		}
		time.Sleep(time.Second * time.Duration(i+1))
	}

	if err != nil {
		return nil, err
	}

	cleaned := make([]string, 0, len(keywords))
	for _, kw := range keywords {
		kw = strings.ToLower(strings.TrimSpace(kw))
		kw = strings.Trim(kw, `.,;:"'!?"`)
		if len(kw) > 2 && len(kw) < 50 {
			cleaned = append(cleaned, kw)
		}
	}

	return cleaned, nil
}
