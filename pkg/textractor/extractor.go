package textractor

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type TextExtractor struct{}

func NewTextExtractor() *TextExtractor {
	return &TextExtractor{}
}

func (te *TextExtractor) ExtractText(filePath string) (string, error) {
	// Only allow plain text files
	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".txt", ".md", ".csv", ".log":
		// These are safe to read directly
	default:
		return "", fmt.Errorf("unsupported file type: %s - only plain text files (.txt, .md, .csv, .log) supported", ext)
	}

	// Simple file read (equivalent to 'cat')
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	return string(content), nil
}