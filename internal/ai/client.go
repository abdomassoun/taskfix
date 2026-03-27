package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const openRouterURL = "https://openrouter.ai/api/v1/chat/completions"

// Client handles communication with the AI provider.
type Client struct {
	provider   string
	apiKey     string
	model      string
	http       *http.Client
	endpointURL string // overridable for tests
}

// NewClient constructs a new AI client using the default endpoint URL.
func NewClient(provider, apiKey, model string) *Client {
	return &Client{
		provider:    provider,
		apiKey:      apiKey,
		model:       model,
		http:        &http.Client{Timeout: 60 * time.Second},
		endpointURL: openRouterURL,
	}
}

// NewClientWithURL constructs an AI client with a custom endpoint URL.
// Intended for testing with mock HTTP servers.
func NewClientWithURL(provider, apiKey, model, url string) *Client {
	c := NewClient(provider, apiKey, model)
	c.endpointURL = url
	return c
}

// Complete sends the prompt to the AI and returns the response text.
func (c *Client) Complete(prompt string) (string, error) {
	if c.apiKey == "" {
		return "", fmt.Errorf("no API key configured — set OPENROUTER_API_KEY or use --api-key")
	}

	switch c.provider {
	case "openrouter", "":
		return c.openRouterComplete(prompt)
	default:
		return "", fmt.Errorf("unknown provider %q — only 'openrouter' is supported in v1", c.provider)
	}
}

// ── OpenRouter ────────────────────────────────────────────────────────────────

type openRouterRequest struct {
	Model    string    `json:"model"`
	Messages []message `json:"messages"`
}

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openRouterResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
		Code    int    `json:"code"`
	} `json:"error"`
}

// FetchModels queries the OpenRouter models endpoint and returns a list of model IDs.
func (c *Client) FetchModels() ([]string, error) {
	if c.apiKey == "" {
		return nil, fmt.Errorf("no API key configured — set OPENROUTER_API_KEY or use --api-key")
	}

	url := "https://openrouter.ai/api/v1/models"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Accept", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request: %w", err)
	}
	defer resp.Body.Close()

	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(respData))
	}

	// Try to unmarshal flexible shapes: either an array or an object with a "models" key.
	var raw any
	if err := json.Unmarshal(respData, &raw); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}

	names := make([]string, 0)
	seen := map[string]bool{}

	addName := func(n string) {
		if n == "" {
			return
		}
		if !seen[n] {
			seen[n] = true
			names = append(names, n)
		}
	}

	switch v := raw.(type) {
	case []any:
		for _, el := range v {
			switch e := el.(type) {
			case string:
				addName(e)
			case map[string]any:
				if id, ok := e["id"].(string); ok {
					addName(id)
				} else if name, ok := e["model"].(string); ok {
					addName(name)
				} else if name, ok := e["name"].(string); ok {
					addName(name)
				}
			}
		}
	case map[string]any:
		// common patterns: { "models": [...] } or { "data": [...] } or other keys
		for _, key := range []string{"models", "data", "items"} {
			if arr, ok := v[key]; ok {
				if a, ok := arr.([]any); ok {
					for _, el := range a {
						switch e := el.(type) {
						case string:
							addName(e)
						case map[string]any:
							if id, ok := e["id"].(string); ok {
								addName(id)
							} else if name, ok := e["model"].(string); ok {
								addName(name)
							} else if name, ok := e["name"].(string); ok {
								addName(name)
							}
						}
					}
				}
			}
		}
		// If no known key contained the models, try scanning all top-level arrays.
		if len(names) == 0 {
			for _, val := range v {
				if a, ok := val.([]any); ok {
					for _, el := range a {
						switch e := el.(type) {
						case string:
							addName(e)
						case map[string]any:
							if id, ok := e["id"].(string); ok {
								addName(id)
							} else if name, ok := e["model"].(string); ok {
								addName(name)
							} else if name, ok := e["name"].(string); ok {
								addName(name)
							}
						}
					}
				}
			}
		}
	}

	return names, nil
}

func (c *Client) openRouterComplete(prompt string) (string, error) {
	reqBody := openRouterRequest{
		Model: c.model,
		Messages: []message{
			{Role: "user", Content: prompt},
		},
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("marshaling request: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, c.endpointURL, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("HTTP-Referer", "https://github.com/taskfix/taskfix")
	req.Header.Set("X-Title", "TaskFix CLI")

	resp, err := c.http.Do(req)
	if err != nil {
		return "", fmt.Errorf("HTTP request: %w", err)
	}
	defer resp.Body.Close()

	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("reading response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		// Try to extract a useful error message from the body
		var errResp openRouterResponse
		if json.Unmarshal(respData, &errResp) == nil && errResp.Error != nil {
			return "", fmt.Errorf("API error %d: %s", resp.StatusCode, errResp.Error.Message)
		}
		return "", fmt.Errorf("API error %d: %s", resp.StatusCode, string(respData))
	}

	var result openRouterResponse
	if err := json.Unmarshal(respData, &result); err != nil {
		return "", fmt.Errorf("parsing response: %w", err)
	}

	if result.Error != nil {
		return "", fmt.Errorf("API error: %s", result.Error.Message)
	}

	if len(result.Choices) == 0 {
		return "", fmt.Errorf("AI returned no choices")
	}

	return result.Choices[0].Message.Content, nil
}
