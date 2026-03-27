package prompt_test

import (
	"strings"
	"testing"

	"github.com/taskfix/taskfix/internal/prompt"
	"github.com/taskfix/taskfix/internal/rules"
)

func TestBuild_ContainsInput(t *testing.T) {
	input := "user cant login"
	rs := rules.Default()
	p := prompt.Build(input, rs)

	if !strings.Contains(p, input) {
		t.Errorf("expected prompt to contain input %q", input)
	}
}

func TestBuild_ContainsRules(t *testing.T) {
	rs := &rules.RuleSet{
		Style: "default",
		Rules: []string{"Fix grammar", "Use bullet points"},
	}
	p := prompt.Build("some task", rs)

	for _, rule := range rs.Rules {
		if !strings.Contains(p, rule) {
			t.Errorf("expected prompt to contain rule %q", rule)
		}
	}
}

func TestBuild_DefaultStyle_ContainsTemplate(t *testing.T) {
	rs := &rules.RuleSet{Style: "default", Rules: []string{"Fix grammar"}}
	p := prompt.Build("task", rs)

	if !strings.Contains(p, "Acceptance Criteria") {
		t.Error("expected default template to include 'Acceptance Criteria'")
	}
	if !strings.Contains(p, "Title:") {
		t.Error("expected default template to include 'Title:'")
	}
}

func TestBuild_JiraStyle_ContainsTemplate(t *testing.T) {
	rs := &rules.RuleSet{Style: "jira", Rules: []string{"Fix grammar"}}
	p := prompt.Build("task", rs)

	if !strings.Contains(p, "Summary") {
		t.Error("expected jira template to include 'Summary'")
	}
}

func TestBuild_GithubStyle_ContainsTemplate(t *testing.T) {
	rs := &rules.RuleSet{Style: "github", Rules: []string{"Fix grammar"}}
	p := prompt.Build("task", rs)

	if !strings.Contains(p, "Steps to Reproduce") {
		t.Error("expected github template to include 'Steps to Reproduce'")
	}
}

func TestBuild_UnknownStyle_FallsBackToDefault(t *testing.T) {
	rs := &rules.RuleSet{Style: "unknown-style", Rules: []string{"Fix grammar"}}
	p := prompt.Build("task", rs)

	// Should not panic and should produce some output format
	if !strings.Contains(p, "Acceptance Criteria") {
		t.Error("unknown style should fall back to default template")
	}
}
