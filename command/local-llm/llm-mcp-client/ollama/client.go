package ollama

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// Client interface for Ollama API interactions
type Client interface {
	Generate(prompt string) (string, error)
}

// HTTPClient implements the Client interface using HTTP requests
type HTTPClient struct {
	URL        string
	Model      string
	HTTPClient *http.Client
}

// Request represents an Ollama API request
type Request struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

// Response represents an Ollama API response
type Response struct {
	Model     string    `json:"model"`
	CreatedAt time.Time `json:"created_at"`
	Response  string    `json:"response"`
	Done      bool      `json:"done"`
}

// NewHTTPClient creates a new Ollama HTTP client
func NewHTTPClient(url, model string) *HTTPClient {
	return &HTTPClient{
		URL:   url,
		Model: model,
		HTTPClient: &http.Client{
			Timeout: 300 * time.Second, // 5 minute timeout for LLM generation
		},
	}
}

// Generate sends a prompt to Ollama and returns the generated response
func (c *HTTPClient) Generate(prompt string) (string, error) {
	reqBody := Request{
		Model:  c.Model,
		Prompt: prompt,
		Stream: false,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := c.HTTPClient.Post(
		c.URL+"/api/generate",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return "", fmt.Errorf("failed to send request to Ollama: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Ollama API error: status %d", resp.StatusCode)
	}

	var ollamaResp Response
	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
		return "", fmt.Errorf("failed to decode Ollama response: %w", err)
	}

	return strings.TrimSpace(ollamaResp.Response), nil
}