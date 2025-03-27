package ollama

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type Client struct {
	BaseURL string
	Model   string
	Timeout time.Duration
}

func New(baseURL, model string) *Client {
	return &Client{
		BaseURL: strings.TrimSuffix(baseURL, "/"),
		Model:   model,
		Timeout: 300 * time.Second,
	}
}

func (c *Client) ExtractKeywords(text string) ([]string, error) {
	if len(text) > 10000 {
		text = text[:10000]
	}

	prompt := `You are tasked with being a search optimizer. Given the text content of a file and its metadata (such as creation date, path, file name, etc.), generate only 5-10 relevant keywords that can help users search for this file efficiently. Return the keywords as a list where each keyword is prefixed with a '-'. Example output: - academics - module1 - sem4 - os
Text: ` + text

	requestBody := map[string]interface{}{
		"model":  c.Model,
		"prompt": prompt,
		"stream": false,
		"options": map[string]interface{}{
			"temperature": 0.3,
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}

	client := &http.Client{Timeout: c.Timeout}
	resp, err := client.Post(
		c.BaseURL+"/api/generate",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return nil, fmt.Errorf("API request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
	}

	var response struct {
		Response string `json:"response"`
		Error    string `json:"error"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("decode error: %w", err)
	}

	if response.Error != "" {
		return nil, fmt.Errorf("model error: %s", response.Error)
	}

	// Extract keywords with regex for lines starting with "- "
	re := regexp.MustCompile(`(?m)^-\s*(\w[\w\s]*)$`)
	matches := re.FindAllStringSubmatch(response.Response, -1)

	var keywords []string
	for _, match := range matches {
		kw := strings.TrimSpace(match[1])
		kw = strings.ToLower(kw)
		kw = strings.Trim(kw, `.,;:"'!?`)
		if len(kw) > 2 { // Minimum keyword length
			keywords = append(keywords, kw)
		}
	}

	if len(keywords) == 0 {
		return nil, fmt.Errorf("no valid keywords generated")
	}

	return keywords, nil
}
