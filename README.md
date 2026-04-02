![TaskFix Banner](TFIX-README-header.png)

# TaskFix

> Transform messy task descriptions into clean, structured, AI-formatted technical tasks — from your terminal.

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/github.com/abdomassoun/taskfix)](https://goreportcard.com/report/github.com/abdomassoun/taskfix)

---

## 📖 Table of Contents

- [What is TaskFix?](#what-is-taskfix)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Usage](#usage)
- [Configuration](#configuration)
- [Examples](#examples)
- [For Developers](#for-developers)
- [Contributing](#contributing)
- [License](#license)

---

## What is TaskFix?

TaskFix is a command-line tool that helps developers and teams transform rough task descriptions into well-formatted, professional technical tasks. Whether you're creating GitHub issues, Jira tickets, or internal documentation, TaskFix uses AI to clean up your descriptions and ensure they follow best practices.

**Key Features:**
- 🤖 AI-powered formatting using OpenRouter
- 📝 Support for multiple output formats (GitHub, Jira, custom)
- ⚡ Fast and pipe-friendly CLI interface
- 🎨 Customizable rules and templates
- 🔒 Secure configuration management
- 🚀 No dependencies beyond the binary

---

## Installation

### Quick Install (Debian/Ubuntu)

Download and install the latest `.deb` package:

```bash
wget https://github.com/abdomassoun/taskfix/releases/latest/download/taskfix_latest_amd64.deb
sudo dpkg -i taskfix_latest_amd64.deb
```

If you encounter dependency issues:
```bash
sudo apt-get install -f
```

### Other Platforms

For other platforms, download pre-built binaries from the [releases page](https://github.com/abdomassoun/taskfix/releases) or [build from source](#building-from-source).

---

## Quick Start

### 1. Get an API Key

TaskFix uses [OpenRouter](https://openrouter.ai) for AI processing:

1. Visit [openrouter.ai](https://openrouter.ai) and sign up
2. Navigate to [API Keys](https://openrouter.ai/keys)
3. Create a new key (starts with `sk-or-v1-`)

### 2. Configure TaskFix

Create a config file with your API key:

```bash
cat > ~/.tfixrc << 'EOF'
{
  "api_key": "sk-or-v1-your-api-key-here",
  "model": "openai/gpt-4o-mini"
}
EOF

chmod 600 ~/.tfixrc
```

### 3. Use TaskFix

```bash
taskfix "user cant login when password wrong"
```

**Output:**
```markdown
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

---

## Usage

### Basic Usage

```bash
# Process a task description
taskfix "bug description here"

# Read from stdin
echo "login bug" | taskfix
cat input.txt | taskfix

# Save to file
taskfix "task description" > task.md
```

### Using Different Formats

TaskFix includes built-in templates for common platforms:

```bash
# GitHub Issues format
taskfix "bug description" --rules configs/github.json

# Jira tickets format
taskfix "bug description" --rules configs/jira.json

# Default format
taskfix "bug description" --rules configs/default.json
```

### Integration with Other Tools

```bash
# Create GitHub issue directly
taskfix "login bug" | gh issue create --title "Bug Fix" --body -

# Copy to clipboard (Linux)
taskfix "task description" | xclip -selection clipboard

# Copy to clipboard (macOS)
taskfix "task description" | pbcopy
```

### Command-Line Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--config` | `-c` | Path to config file |
| `--rules` | `-r` | Path to custom rules file |
| `--model` | `-m` | AI model override |
| `--api-key` | `-k` | API key (overrides config) |
| `--silent` | `-s` | Suppress progress messages |

---

## Configuration

### Configuration Files

TaskFix looks for configuration in these locations (in order of priority):

1. `--config` flag (explicit path)
2. `~/.tfixrc` (user config - **recommended**)
3. `~/.config/taskfix/config` (XDG config)
4. `/etc/taskfix/config` (system-wide config)

### Configuration Format

```json
{
  "provider": "openrouter",
  "api_key": "sk-or-v1-your-api-key-here",
  "model": "openai/gpt-4o-mini",
  "rules_file": "path/to/custom-rules.json"
}
```

**All fields are optional.** API key can be provided via:
- Config file
- Environment variable: `OPENROUTER_API_KEY`
- Command-line flag: `--api-key`

### Custom Rules

Create custom formatting rules in JSON:

```json
{
  "style": "my-custom-style",
  "rules": [
    "Fix all grammar and spelling errors",
    "Use clear, professional technical language",
    "Include estimated time to complete",
    "Add relevant labels or tags",
    "Be concise and actionable"
  ]
}
```

Use your custom rules:
```bash
taskfix "task description" --rules my-rules.json
```

### Security Best Practices

- **Never commit API keys** to version control
- **Protect config files** with appropriate permissions:
  ```bash
  chmod 600 ~/.tfixrc
  ```
- **Use environment variables** for CI/CD environments
- **Rotate API keys** regularly via OpenRouter dashboard

For detailed configuration options, see [INSTALL.md](INSTALL.md) and [configs/README.md](configs/README.md).

---

## Examples

### Example 1: Bug Report

**Input:**
```bash
taskfix "button doesnt work when clicked twice fast"
```

**Output:**
```markdown
Title: Fix button double-click handling issue

Description:
- Button becomes unresponsive when clicked twice in rapid succession
- Expected behavior is to handle multiple clicks gracefully
- May cause user frustration and missed actions

Acceptance Criteria:
- Button responds correctly to rapid successive clicks
- No UI freezing or unresponsive states
- Proper debouncing or state management implemented
```

### Example 2: Feature Request

**Input:**
```bash
taskfix "add dark mode to settings page" --rules configs/github.json
```

**Output:**
```markdown
## Add dark mode support to settings page

**Description:**
The settings page currently only supports light mode. Users have requested a dark mode option for better accessibility and reduced eye strain in low-light environments.

**Expected Behavior:**
- Settings page should include a dark mode toggle
- Theme preference should persist across sessions
- All UI elements should be readable in dark mode

**Acceptance Criteria:**
- [ ] Dark mode toggle added to settings page
- [ ] Theme preference saved to user preferences
- [ ] All text and UI elements properly styled for dark mode
- [ ] Smooth transition between light and dark modes
```

---

## For Developers

### Building from Source

**Prerequisites:**
- Go 1.21 or higher
- Git

**Build instructions:**

```bash
# Clone the repository
git clone https://github.com/abdomassoun/taskfix.git
cd taskfix

# Install dependencies
go mod download

# Build the binary
make build

# Or build manually
go build -o taskfix .

# Install system-wide (optional)
sudo cp taskfix /usr/local/bin/
```

### Project Structure

```
taskfix/
├── cmd/
│   ├── root.go          # CLI entry point and orchestration
│   └── config.go        # Configuration loading and management
├── internal/
│   ├── ai/
│   │   └── client.go    # OpenRouter API client
│   ├── rules/
│   │   └── engine.go    # Rule loading and processing
│   ├── prompt/
│   │   └── builder.go   # Prompt construction
│   └── output/
│       └── formatter.go # Output formatting
├── configs/
│   ├── default.json     # Default formatting rules
│   ├── jira.json        # Jira-style preset
│   └── github.json      # GitHub Issue preset
├── main.go              # Application entry point
├── go.mod               # Go module definition
└── Makefile             # Build automation
```

### Development Commands

```bash
# Build the project
make build

# Run tests
make test

# Run linter
make lint

# Clean build artifacts
make clean

# Build .deb package
make deb
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...
```

### Code Style

This project follows standard Go conventions:
- Use `gofmt` for formatting
- Follow [Effective Go](https://golang.org/doc/effective_go.html) guidelines
- Write tests for new features
- Keep functions focused and modular

---

## Contributing

We welcome contributions from the community! Whether it's:

- 🐛 Bug reports
- 💡 Feature requests
- 📖 Documentation improvements
- 🔧 Code contributions

Please read our [Contributing Guidelines](CONTRIBUTING.md) to get started.

### Quick Contribution Guide

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/my-new-feature`
3. Make your changes and commit: `git commit -am 'Add some feature'`
4. Push to the branch: `git push origin feature/my-new-feature`
5. Submit a pull request

For detailed guidelines, see [CONTRIBUTING.md](CONTRIBUTING.md).

---

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## Support

- 📚 [Documentation](https://github.com/abdomassoun/taskfix/wiki)
- 🐛 [Issue Tracker](https://github.com/abdomassoun/taskfix/issues)
- 💬 [Discussions](https://github.com/abdomassoun/taskfix/discussions)

---

**Made with ❤️ by the TaskFix community**
