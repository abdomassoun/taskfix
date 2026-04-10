package rules

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// RuleSet defines the formatting rules applied to the AI prompt.
type RuleSet struct {
	Style string   `json:"style"` // e.g. "default", "jira", "github"
	Rules []string `json:"rules"`
}

// Default returns the built-in default rule set used in v1.
func Default() *RuleSet {
	return &RuleSet{
		Style: "default",
		Rules: []string{
			"Fix all grammar and spelling errors",
			"Use clear, professional technical language",
			"Structure the output with a Title, Description, and Acceptance Criteria",
			"Use bullet points for description items and acceptance criteria",
			"Be concise — remove filler words and redundancy",
			"Infer missing context from the description where reasonable",
			"The title should be a short, imperative sentence (max 10 words)",
		},
	}
}

// Load reads a custom rules JSON file. Returns defaults if path is empty.
func Load(path string) (*RuleSet, error) {
	if path == "" {
		if rs, err := loadFromFileIfExists(filepath.Join(systemRulesDir(), "default.json")); err != nil {
			return nil, err
		} else if rs != nil {
			return rs, nil
		}
		return Default(), nil
	}

	candidates := []string{path}
	if !filepath.IsAbs(path) {
		candidates = append(candidates,
			filepath.Join(systemRulesDir(), path),
			filepath.Join(systemRulesDir(), filepath.Base(path)),
		)
	}

	for _, candidate := range candidates {
		rs, err := loadFromFileIfExists(candidate)
		if err != nil {
			return nil, err
		}
		if rs != nil {
			return rs, nil
		}
	}

	return nil, fmt.Errorf("rules file not found: %q", path)
}

func loadFromFileIfExists(path string) (*RuleSet, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}

	var rs RuleSet
	if err := json.Unmarshal(data, &rs); err != nil {
		return nil, err
	}

	// Backfill defaults for missing fields
	if rs.Style == "" {
		rs.Style = "default"
	}
	if len(rs.Rules) == 0 {
		rs.Rules = Default().Rules
	}

	return &rs, nil
}

func systemRulesDir() string {
	if dir := os.Getenv("TASKFIX_RULES_DIR"); dir != "" {
		return dir
	}
	return "/etc/taskfix/config.d"
}
