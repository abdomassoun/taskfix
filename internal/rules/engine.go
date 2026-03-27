package rules

import (
	"encoding/json"
	"os"
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
		return Default(), nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
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
