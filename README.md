# TaskFix

> Transform messy task descriptions into clean, structured, AI-formatted technical tasks — from your terminal.

## Install

```bash
git clone https://github.com/taskfix/taskfix
cd taskfix
go mod tidy
go build -o taskfix .
```

Move the binary somewhere on your PATH:

```bash
mv taskfix /usr/local/bin/taskfix
```

## Setup

TaskFix requires an [OpenRouter](https://openrouter.ai) API key. Set it as an environment variable:

```bash
export OPENROUTER_API_KEY=your_key_here
```

Or pass it at runtime:

```bash
taskfix "bug description" --api-key your_key_here
```

## Usage

### Argument input

```bash
taskfix "user cant login when password wrong"
```

### Stdin (pipe-friendly)

```bash
echo "login bug" | taskfix
cat input.txt | taskfix
```

### Custom rules

```bash
taskfix "bug description" --rules configs/jira.json
taskfix "bug description" --rules configs/github.json
taskfix "bug description" --rules my-rules.json
```

### Config file

```bash
taskfix "bug description" --config taskfix.json
```

### Pipe to GitHub CLI

```bash
taskfix "login bug" | gh issue create --title "Bug Fix" --body -
```

### Save to file

```bash
taskfix "task description" > task.md
```

## Config file format

```json
{
  "provider": "openrouter",
  "api_key": "YOUR_KEY",
  "model": "openai/gpt-4o-mini",
  "rules_file": "configs/jira.json"
}
```

Default config locations (auto-discovered if `--config` is not set):
- `./taskfix.json`
- `~/.config/taskfix/config.json`

## Custom rules format

```json
{
  "style": "default",
  "rules": [
    "Fix grammar",
    "Use bullet points",
    "Add acceptance criteria",
    "Be concise"
  ]
}
```

Built-in styles: `default`, `jira`, `github`

## Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--config` | `-c` | Path to config JSON file |
| `--rules` | `-r` | Path to custom rules JSON file |
| `--model` | `-m` | AI model override (e.g. `anthropic/claude-3-haiku`) |
| `--api-key` | `-k` | API key (overrides env var) |
| `--silent` | `-s` | Suppress stderr progress messages |

## Example output

**Input:**
```
user cant login when password wrong
```

**Output:**
```
Title: Fix login failure on incorrect password

Description:
- Users are unable to log in when entering an incorrect password
- The system does not display a clear error message on failed attempts
- Login process may become unstable after repeated failures

Acceptance Criteria:
- System displays a clear, user-friendly error message on login failure
- Login form remains functional and stable after incorrect attempts
- Correct credentials continue to work as expected
```

## Project structure

```
taskfix/
├── cmd/
│   ├── root.go          # CLI entry point, input resolution, orchestration
│   └── config.go        # Config loading and flag overrides
├── internal/
│   ├── ai/
│   │   └── client.go    # OpenRouter API client
│   ├── rules/
│   │   └── engine.go    # Rule loading (default + custom JSON)
│   ├── prompt/
│   │   └── builder.go   # Prompt construction with style templates
│   └── output/
│       └── formatter.go # Output cleanup for stdout
├── configs/
│   ├── default.json     # Default rules preset
│   ├── jira.json        # Jira-style preset
│   └── github.json      # GitHub Issue preset
├── main.go
├── go.mod
└── README.md
```
