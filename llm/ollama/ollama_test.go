package ollama

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractKeywords(t *testing.T) {
	mockResponse := `{"response": "test, sample, data, keyword, check"}`
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockResponse))
	}))
	defer server.Close()

	client := New(server.URL)

	keywords, err := client.ExtractKeywords("AI in healthcare is revolutionizing data analysis.")
	t.Logf("Received Keywords: %v", keywords)

	assert.NoError(t, err, "Expected no error")
	assert.Equal(t, []string{"test", "sample", "data", "keyword", "check"}, keywords)
}

func TestExtractKeywords_Failure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Internal server error"}`))
	}))
	defer server.Close()

	client := New(server.URL)

	_, err := client.ExtractKeywords("Some sample text")
	assert.Error(t, err, "Expected an error when the server fails")
}

func TestExtractKeywords_EmptyResponse(t *testing.T) {
	mockResponse := `{"response": ""}`
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockResponse))
	}))
	defer server.Close()

	client := New(server.URL)

	keywords, err := client.ExtractKeywords("Text with no keywords")
	t.Logf("Received Keywords (empty case): %v", keywords)

	assert.NoError(t, err, "Expected no error for empty response")
	assert.Empty(t, keywords, "Expected empty keyword list")
}

func TestExtractKeywords_MalformedResponse(t *testing.T) {
	mockResponse := `{"response": test sample data keyword check}`
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockResponse))
	}))
	defer server.Close()

	client := New(server.URL)

	_, err := client.ExtractKeywords("Sample text")
	assert.Error(t, err, "Expected an error for malformed JSON response")
}
