package ai_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/taskfix/taskfix/internal/ai"
)

// newMockServer spins up a fake OpenRouter-compatible HTTP server.
func newMockServer(t *testing.T, status int, body any) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify required headers
		if r.Header.Get("Authorization") == "" {
			t.Error("expected Authorization header")
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("unexpected Content-Type: %s", r.Header.Get("Content-Type"))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(body)
	}))
}

func TestComplete_Success(t *testing.T) {
	mockResp := map[string]any{
		"choices": []map[string]any{
			{"message": map[string]string{"content": "Title: Fix login bug\n\nDescription:\n- User cannot log in"}},
		},
	}
	srv := newMockServer(t, http.StatusOK, mockResp)
	defer srv.Close()

	client := ai.NewClientWithURL("openrouter", "test-key", "gpt-4o-mini", srv.URL+"/v1/chat/completions")
	result, err := client.Complete("some prompt")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "Title: Fix login bug\n\nDescription:\n- User cannot log in" {
		t.Errorf("unexpected result: %q", result)
	}
}

func TestComplete_APIError(t *testing.T) {
	mockResp := map[string]any{
		"error": map[string]any{
			"message": "invalid API key",
			"code":    401,
		},
	}
	srv := newMockServer(t, http.StatusUnauthorized, mockResp)
	defer srv.Close()

	client := ai.NewClientWithURL("openrouter", "bad-key", "gpt-4o-mini", srv.URL+"/v1/chat/completions")
	_, err := client.Complete("some prompt")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestComplete_NoAPIKey(t *testing.T) {
	client := ai.NewClient("openrouter", "", "gpt-4o-mini")
	_, err := client.Complete("some prompt")
	if err == nil {
		t.Fatal("expected error for missing API key, got nil")
	}
}

func TestComplete_EmptyChoices(t *testing.T) {
	mockResp := map[string]any{"choices": []any{}}
	srv := newMockServer(t, http.StatusOK, mockResp)
	defer srv.Close()

	client := ai.NewClientWithURL("openrouter", "test-key", "gpt-4o-mini", srv.URL+"/v1/chat/completions")
	_, err := client.Complete("some prompt")
	if err == nil {
		t.Fatal("expected error for empty choices, got nil")
	}
}

func TestComplete_UnknownProvider(t *testing.T) {
	client := ai.NewClient("unknown-provider", "key", "model")
	_, err := client.Complete("some prompt")
	if err == nil {
		t.Fatal("expected error for unknown provider, got nil")
	}
}
