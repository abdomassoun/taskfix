package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfigPath(t *testing.T) {
	t.Setenv("HOME", "/tmp/taskfix-home")

	if got := defaultConfigPath(""); got != "/tmp/taskfix-home/.tfixrc" {
		t.Fatalf("default path mismatch: got %q", got)
	}

	if got := defaultConfigPath("~/custom/config.json"); got != "/tmp/taskfix-home/custom/config.json" {
		t.Fatalf("tilde expansion mismatch: got %q", got)
	}
}

func TestSaveAPIKey_PreservesExistingFields(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")

	initial := Config{
		Provider:  "openrouter",
		APIKey:    "old-key",
		Model:     "openai/gpt-4o-mini",
		RulesFile: "/tmp/rules.json",
	}
	data, _ := json.Marshal(initial)
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatalf("failed to write initial config: %v", err)
	}

	if err := saveAPIKey(path, "new-key"); err != nil {
		t.Fatalf("saveAPIKey returned error: %v", err)
	}

	updatedData, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read updated config: %v", err)
	}

	var updated Config
	if err := json.Unmarshal(updatedData, &updated); err != nil {
		t.Fatalf("failed to parse updated config: %v", err)
	}

	if updated.APIKey != "new-key" {
		t.Fatalf("expected api key to be updated, got %q", updated.APIKey)
	}
	if updated.RulesFile != initial.RulesFile {
		t.Fatalf("expected rules file to be preserved, got %q", updated.RulesFile)
	}
	if updated.Model != initial.Model {
		t.Fatalf("expected model to be preserved, got %q", updated.Model)
	}
}

func TestSaveAPIKey_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")
	if err := os.WriteFile(path, []byte("{invalid"), 0o600); err != nil {
		t.Fatalf("failed to write invalid config: %v", err)
	}

	if err := saveAPIKey(path, "new-key"); err == nil {
		t.Fatal("expected error for invalid config, got nil")
	}
}
