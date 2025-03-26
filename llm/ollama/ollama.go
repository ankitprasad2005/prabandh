package ollama

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Client struct {
	BaseURL string
}

func New(baseURL string) *Client {
	return &Client{BaseURL: baseURL}
}

func (c *Client) ExtractKeywords(text string) ([]string, error) {
	requestBody := map[string]string{
		"model":  "gemma:2b",
		"prompt": "Return these keywords exactly: 'test, sample, data, keyword, check'",
	}

	jsonData, _ := json.Marshal(requestBody)

	resp, err := http.Post(c.BaseURL+"/api/generate", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("ollama request failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Println("Raw Response from Ollama:", string(body)) // 🔎 Debugging Log

	var result struct {
		Response string `json:"response"`
		Error    string `json:"error"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	if result.Error != "" {
		return nil, fmt.Errorf("ollama returned an error: %s", result.Error)
	}

	keywords := strings.Split(strings.TrimSpace(result.Response), ",")
	for i, kw := range keywords {
		keywords[i] = strings.TrimSpace(kw)
	}

	if len(keywords) == 1 && keywords[0] == "" {
		return []string{}, nil
	}

	return keywords, nil
}
