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

func (te *TextExtractor) SupportedExtensions() []string {
	return []string{
		".txt", ".md", ".csv", ".log",
		".go", ".py", ".js", ".ts",
		".html", ".css", ".json",
		".yaml", ".yml", ".sh",
	}
}

func (te *TextExtractor) CanExtract(filePath string) bool {
	ext := strings.ToLower(filepath.Ext(filePath))
	for _, supported := range te.SupportedExtensions() {
		if ext == supported {
			return true
		}
	}
	return false
}

func (te *TextExtractor) ExtractText(filePath string) (string, error) {
	if !te.CanExtract(filePath) {
		return "", fmt.Errorf("unsupported file type: %s", filepath.Ext(filePath))
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	return string(content), nil
}
