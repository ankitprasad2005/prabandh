package textractor

import (
	"os"
	"path/filepath"
	"testing"
)

func TestTextExtractor(t *testing.T) {
	// Setup test files
	setupTestFiles := func() {
		os.Mkdir("testdata", 0755)
		os.WriteFile("testdata/sample.txt", []byte("This is a test text file"), 0644)
		os.WriteFile("testdata/empty.txt", []byte(""), 0644)
	}
	setupTestFiles()
	defer os.RemoveAll("testdata")

	extractor := NewTextExtractor()

	t.Run("Valid text file", func(t *testing.T) {
		text, err := extractor.ExtractText(filepath.Join("testdata", "sample.txt"))
		if err != nil {
			t.Fatalf("ExtractText failed: %v", err)
		}

		if text != "This is a test text file" {
			t.Errorf("Expected specific text, got: %q", text)
		}
	})

	t.Run("Empty file", func(t *testing.T) {
		_, err := extractor.ExtractText(filepath.Join("testdata", "empty.txt"))
		if err != nil {
			t.Errorf("Should handle empty files, got error: %v", err)
		}
	})

	t.Run("Unsupported file", func(t *testing.T) {
		_, err := extractor.ExtractText(filepath.Join("testdata", "sample.pdf"))
		if err == nil {
			t.Error("Expected error for unsupported file type")
		}
	})
}
