package output_test

import (
	"testing"

	"github.com/taskfix/taskfix/internal/output"
)

func TestFormat_TrimWhitespace(t *testing.T) {
	result := output.Format("  hello world  \n")
	if result != "hello world" {
		t.Errorf("expected trimmed string, got %q", result)
	}
}

func TestFormat_StripCodeFence(t *testing.T) {
	raw := "```\nTitle: Fix bug\n\nDescription:\n- Something\n```"
	result := output.Format(raw)
	if result != "Title: Fix bug\n\nDescription:\n- Something" {
		t.Errorf("unexpected result after stripping fences: %q", result)
	}
}

func TestFormat_StripLanguageTaggedFence(t *testing.T) {
	raw := "```markdown\nTitle: Fix bug\n```"
	result := output.Format(raw)
	if result != "Title: Fix bug" {
		t.Errorf("unexpected result: %q", result)
	}
}

func TestFormat_NofencePassthrough(t *testing.T) {
	raw := "Title: Fix login bug\n\nDescription:\n- User cannot log in"
	result := output.Format(raw)
	if result != raw {
		t.Errorf("expected passthrough, got %q", result)
	}
}

func TestFormat_EmptyString(t *testing.T) {
	result := output.Format("")
	if result != "" {
		t.Errorf("expected empty string, got %q", result)
	}
}
