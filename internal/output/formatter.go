package output

import (
	"strings"
)

// Format cleans up the raw AI response for stdout output.
// It trims whitespace and strips any accidental markdown code fences
// the model might have added despite instructions.
func Format(raw string) string {
	s := strings.TrimSpace(raw)

	// Strip ``` fences if the model wrapped the output anyway
	if strings.HasPrefix(s, "```") {
		lines := strings.SplitN(s, "\n", 2)
		if len(lines) == 2 {
			s = lines[1]
		}
	}
	if strings.HasSuffix(s, "```") {
		s = s[:strings.LastIndex(s, "```")]
	}

	return strings.TrimSpace(s)
}
