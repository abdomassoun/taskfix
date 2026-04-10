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
	t.Setenv("TASKFIX_RULES_DIR", t.TempDir())

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
	t.Setenv("TASKFIX_RULES_DIR", t.TempDir())

	_, err := rules.Load("missing-rules.json")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestLoad_EmptyRules_BackfillsDefaults(t *testing.T) {
	t.Setenv("TASKFIX_RULES_DIR", t.TempDir())

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

func TestLoad_EmptyPath_UsesSystemDefaultFile(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("TASKFIX_RULES_DIR", dir)

	custom := rules.RuleSet{
		Style: "system-default",
		Rules: []string{"rule from system default"},
	}
	data, _ := json.Marshal(custom)
	if err := os.WriteFile(dir+"/default.json", data, 0o644); err != nil {
		t.Fatal(err)
	}

	rs, err := rules.Load("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rs.Style != "system-default" {
		t.Errorf("expected style 'system-default', got %q", rs.Style)
	}
}

func TestLoad_RelativeFile_ResolvesFromSystemRulesDir(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("TASKFIX_RULES_DIR", dir)

	custom := rules.RuleSet{
		Style: "github",
		Rules: []string{"rule from system file"},
	}
	data, _ := json.Marshal(custom)
	if err := os.WriteFile(dir+"/github.json", data, 0o644); err != nil {
		t.Fatal(err)
	}

	rs, err := rules.Load("github.json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rs.Style != "github" {
		t.Errorf("expected style 'github', got %q", rs.Style)
	}
	if len(rs.Rules) != 1 {
		t.Errorf("expected 1 rule, got %d", len(rs.Rules))
	}
}
