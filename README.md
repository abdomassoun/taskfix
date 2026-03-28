# TaskFix

> Transform messy task descriptions into clean, structured, AI-formatted technical tasks — from your terminal.

## Install

### From .deb Package (Debian/Ubuntu)

```bash
# Download the latest release
wget https://github.com/taskfix/taskfix/releases/latest/download/taskfix_latest_amd64.deb

# Install the package
sudo dpkg -i taskfix_latest_amd64.deb

# If there are dependency issues, run:
sudo apt-get install -f
```

The .deb package will install:
- Binary: `/usr/local/bin/taskfix`
- System config: `/etc/taskfix/config`

### From Source

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

TaskFix requires an [OpenRouter](https://openrouter.ai) API key. You have multiple options:

### Option 1: User Configuration File (Recommended)

Create `~/.tfixrc` with your API key:

```bash
cat > ~/.tfixrc << 'EOF'
{
  "api_key": "sk-or-v1-your-api-key-here",
  "model": "openai/gpt-4o-mini"
}
EOF

# Protect your API key
chmod 600 ~/.tfixrc
```

### Option 2: Environment Variable

```bash
export OPENROUTER_API_KEY=your_key_here
```

Add to your `~/.bashrc` or `~/.zshrc` for persistence.

### Option 3: Command-line Flag

```bash
taskfix "bug description" --api-key your_key_here
```

### Configuration Priority

TaskFix searches for configuration in this order:

1. `--config` flag (explicit path)
2. `~/.tfixrc` (user config)
3. `~/.config/taskfix/config` (XDG config)
4. `/etc/taskfix/config` (system-wide config)

Environment variables and CLI flags override file settings.

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
# Use a custom config file
taskfix "bug description" --config /path/to/config

# Or let TaskFix auto-discover from default locations
taskfix "bug description"
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

Create a config file at `~/.tfixrc`, `~/.config/taskfix/config`, or `/etc/taskfix/config`:

```json
{
  "provider": "openrouter",
  "api_key": "YOUR_KEY",
  "model": "openai/gpt-4o-mini",
  "rules_file": "configs/jira.json"
}
```

**Note:** All fields are optional. API key can be provided via environment variable or CLI flag.

See [configs/README.md](configs/README.md) for detailed configuration examples and best practices.

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
