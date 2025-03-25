package indexer

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"prabandh/database"
	"prabandh/models"
)

// FileIndexer is responsible for indexing files into the database.
type FileIndexer struct {
	wg sync.WaitGroup
}

// IndexDirectory recursively indexes all files in the given directory.
func (fi *FileIndexer) IndexDirectory(dirPath string) {
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("Error accessing path %s: %v\n", path, err)
			return nil
		}

		if !info.IsDir() {
			fi.wg.Add(1)
			go fi.indexFile(path, info)
		}
		return nil
	})

	if err != nil {
		fmt.Printf("Error walking the directory %s: %v\n", dirPath, err)
	}

	fi.wg.Wait()
}

// indexFile indexes a single file into the database.
func (fi *FileIndexer) indexFile(filePath string, info os.FileInfo) {
	defer fi.wg.Done()

	db := database.DB

	file := models.FileIndex{
		FilePath:     filePath,
		FileName:     info.Name(),
		Extension:    filepath.Ext(info.Name()),
		CreatedDate:  info.ModTime(), // Assuming creation date is not available
		ModifiedDate: info.ModTime(),
		Size:         info.Size(),
		Hash:         calculateHash(filePath), // Replace with actual hash calculation
	}

	if err := db.Create(&file).Error; err != nil {
		fmt.Printf("Failed to index file %s: %v\n", filePath, err)
	} else {
		fmt.Printf("Indexed file: %s\n", filePath)
	}
}

// calculateHash calculates a dummy hash for the file (replace with actual implementation).
func calculateHash(filePath string) string {
	// Placeholder for actual hash calculation logic (e.g., SHA256)
	return fmt.Sprintf("dummy-hash-%d", time.Now().UnixNano())
}
