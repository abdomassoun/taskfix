package rules_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/taskfix/taskfix/internal/rules"
)

func TestDefault(t *testing.T) {
	rs := rules.Default()
	if rs.Style != "default" {
		t.Errorf("expected style 'default', got %q", rs.Style)
	}
	if len(rs.Rules) == 0 {
		t.Error("default rules should not be empty")
	}
}

func TestLoad_EmptyPath_ReturnsDefault(t *testing.T) {
	rs, err := rules.Load("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rs.Style != "default" {
		t.Errorf("expected default style, got %q", rs.Style)
	}
	if len(rs.Rules) == 0 {
		t.Error("expected non-empty rules")
	}
}

func TestLoad_CustomFile(t *testing.T) {
	custom := rules.RuleSet{
		Style: "jira",
		Rules: []string{"Fix grammar", "Add acceptance criteria"},
	}
	data, _ := json.Marshal(custom)

	f, err := os.CreateTemp("", "rules-*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	f.Write(data)
	f.Close()

	rs, err := rules.Load(f.Name())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rs.Style != "jira" {
		t.Errorf("expected style 'jira', got %q", rs.Style)
	}
	if len(rs.Rules) != 2 {
		t.Errorf("expected 2 rules, got %d", len(rs.Rules))
	}
}

func TestLoad_MissingFile_ReturnsError(t *testing.T) {
	_, err := rules.Load("/nonexistent/path/rules.json")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestLoad_EmptyRules_BackfillsDefaults(t *testing.T) {
	partial := map[string]string{"style": "github"}
	data, _ := json.Marshal(partial)

	f, err := os.CreateTemp("", "rules-*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	f.Write(data)
	f.Close()

	rs, err := rules.Load(f.Name())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rs.Rules) == 0 {
		t.Error("expected default rules to be backfilled when rules array is empty")
	}
}
