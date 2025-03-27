package indexer

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"

	"prabandh/database"
	"prabandh/llm/ollama"
	"prabandh/models"
	"prabandh/pkg/textractor"
)

type FileIndexer struct {
	wg            sync.WaitGroup
	textExtractor *textractor.TextExtractor
	ollamaClient  *ollama.Client
	verbose       bool
}

func NewFileIndexer(ollamaURL, model string, verbose bool) *FileIndexer {
	return &FileIndexer{
		textExtractor: textractor.NewTextExtractor(),
		ollamaClient:  ollama.New(ollamaURL, model),
		verbose:       verbose,
	}
}

func (fi *FileIndexer) IndexDirectory(dirPath string) {
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			if fi.verbose {
				fmt.Printf("Error accessing path %s: %v\n", path, err)
			}
			return nil
		}

		if !info.IsDir() {
			fi.wg.Add(1)
			go fi.indexFile(path, info)
		}
		return nil
	})

	if err != nil && fi.verbose {
		fmt.Printf("Error walking directory %s: %v\n", dirPath, err)
	}

	fi.wg.Wait()
}

func (fi *FileIndexer) indexFile(filePath string, info os.FileInfo) {
	defer fi.wg.Done()

	// 1. Collect file metadata
	creationTime := getCreationTime(info)
	hash, err := calculateHash(filePath)
	if err != nil && fi.verbose {
		fmt.Printf("Error calculating hash for %s: %v\n", filePath, err)
		hash = "error-hash"
	}

	file := models.FileIndex{
		FilePath:     filePath,
		FileName:     info.Name(),
		Extension:    filepath.Ext(info.Name()),
		CreatedDate:  creationTime,
		ModifiedDate: info.ModTime(),
		Size:         info.Size(),
		Hash:         hash,
	}

	// 2. Save file metadata first
	if err := database.DB.Create(&file).Error; err != nil {
		if fi.verbose {
			fmt.Printf("Failed to save metadata for %s: %v\n", filePath, err)
		}
		return
	}

	// 3. Skip if file type not supported
	if !fi.textExtractor.CanExtract(filePath) {
		return
	}

	// 4. Extract text content
	content, err := fi.textExtractor.ExtractText(filePath)
	if err != nil {
		if fi.verbose && !strings.Contains(err.Error(), "unsupported file type") {
			fmt.Printf("Extraction error for %s: %v\n", filePath, err)
		}
		return
	}

	// 5. Generate keywords from content + metadata
	metadata := fmt.Sprintf(
		"File: %s\nPath: %s\nSize: %d bytes\nCreated: %s\nModified: %s\nContent:\n%s",
		file.FileName,
		file.FilePath,
		file.Size,
		file.CreatedDate.Format(time.RFC3339),
		file.ModifiedDate.Format(time.RFC3339),
		content,
	)

	keywords, err := fi.ollamaClient.ExtractKeywords(metadata)
	if err != nil {
		if fi.verbose {
			fmt.Printf("Keyword generation failed for %s: %v\n", filePath, err)
		}
		return
	}

	// 6. Save each keyword as a separate row
	var summaries []models.FileSummary
	for _, keyword := range keywords {
		keyword = strings.TrimSpace(keyword)
		if keyword != "" {
			summaries = append(summaries, models.FileSummary{
				FileIndexID:    file.ID,
				SummaryKeyword: keyword,
			})
		}
	}

	if len(summaries) > 0 {
		if err := database.DB.CreateInBatches(&summaries, 100).Error; err != nil {
			if fi.verbose {
				fmt.Printf("Failed to save keywords for %s: %v\n", filePath, err)
			}
		} else if fi.verbose {
			fmt.Printf("Indexed %s with %d keywords\n", filePath, len(summaries))
		}
	}
}

func calculateHash(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

func getCreationTime(info os.FileInfo) time.Time {
	if stat, ok := info.Sys().(*syscall.Stat_t); ok {
		return time.Unix(stat.Ctim.Sec, stat.Ctim.Nsec)
	}
	return info.ModTime()
}
