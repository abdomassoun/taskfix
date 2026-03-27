package prompt

import (
	"fmt"
	"strings"

	"github.com/taskfix/taskfix/internal/rules"
)

const systemPrompt = `You are TaskFix, an expert technical writer specializing in software task descriptions.
Your job is to transform raw, messy task descriptions into clean, structured, and actionable technical tasks.
You output ONLY the formatted task — no preamble, no explanation, no markdown code blocks.`

// Build constructs the full prompt to send to the AI.
func Build(input string, rs *rules.RuleSet) string {
	var sb strings.Builder

	sb.WriteString(systemPrompt)
	sb.WriteString("\n\n")
	sb.WriteString("Apply the following rules when formatting the task:\n")
	for i, rule := range rs.Rules {
		sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, rule))
	}

	sb.WriteString("\n")
	sb.WriteString(outputTemplate(rs.Style))
	sb.WriteString("\n\n")
	sb.WriteString("Raw task input:\n")
	sb.WriteString(input)
	sb.WriteString("\n\nFormatted task:")

	return sb.String()
}

// outputTemplate returns a style-specific output format instruction.
func outputTemplate(style string) string {
	switch style {
	case "jira":
		return `Output format (Jira style):
**Summary:** <short title>

**Description:**
<h3>Problem</h3>
- <bullet points describing the issue>

**Acceptance Criteria:**
- [ ] <criterion 1>
- [ ] <criterion 2>`

	case "github":
		return `Output format (GitHub Issue style):
## <Title>

### Description
- <bullet describing what is wrong or needed>

### Steps to Reproduce (if applicable)
1. <step>

### Expected Behavior
- <what should happen>

### Acceptance Criteria
- [ ] <criterion>`

	default:
		return `Output format:
Title: <short imperative title>

Description:
- <bullet points describing the problem or task>

Acceptance Criteria:
- <criterion 1>
- <criterion 2>`
	}
}
