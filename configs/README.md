# TaskFix Configuration Examples

## User Configuration (~/.tfixrc)

Create a file at `~/.tfixrc` with your personal configuration:

```json
{
  "provider": "openrouter",
  "api_key": "sk-or-v1-your-api-key-here",
  "model": "openai/gpt-4o-mini",
  "rules_file": ""
}
```

## XDG Configuration (~/.config/taskfix/config)

Alternatively, follow XDG Base Directory specification:

```bash
mkdir -p ~/.config/taskfix
cat > ~/.config/taskfix/config << 'EOF'
{
  "provider": "openrouter",
  "api_key": "sk-or-v1-your-api-key-here",
  "model": "anthropic/claude-3-haiku",
  "rules_file": "/home/user/.config/taskfix/rules.json"
}
EOF
```

## System-wide Configuration (/etc/taskfix/config)

For system administrators, edit `/etc/taskfix/config` for all users:

```json
{
  "provider": "openrouter",
  "model": "openai/gpt-4o-mini",
  "api_key": "",
  "rules_file": "/etc/taskfix/rules.json"
}
```

**Note:** Individual users can override system settings with their own `~/.tfixrc` file.

## Configuration Priority

TaskFix uses the following priority order (highest to lowest):

1. Command-line flags (`--api-key`, `--model`, etc.)
2. Environment variables (`OPENROUTER_API_KEY`)
3. `~/.tfixrc` (user configuration)
4. `~/.config/taskfix/config` (XDG configuration)
5. `/etc/taskfix/config` (system-wide configuration)

## Configuration Fields

- **provider**: AI provider to use (default: "openrouter")
- **api_key**: Your API key (can also use `OPENROUTER_API_KEY` environment variable)
- **model**: AI model to use (e.g., "openai/gpt-4o-mini", "anthropic/claude-3-haiku")
- **rules_file**: Optional path to custom rules JSON file for task formatting

## Security Best Practices

### For individual users:
```bash
# Create config with restricted permissions
touch ~/.tfixrc
chmod 600 ~/.tfixrc
echo '{"api_key": "your-key", "model": "openai/gpt-4o-mini"}' > ~/.tfixrc
```

### For system administrators:
```bash
# System config should not contain API keys
# Users should provide their own keys via ~/.tfixrc or environment variables
sudo vim /etc/taskfix/config
```

## Using Environment Variables

Instead of storing your API key in a file, you can use environment variables:

```bash
# In your ~/.bashrc or ~/.zshrc
export OPENROUTER_API_KEY="sk-or-v1-your-api-key-here"

# Or for a single session
OPENROUTER_API_KEY="your-key" taskfix "bug description"
```

## Examples

```bash
# Use default config from ~/.tfixrc
taskfix "user can't login"

# Override with command-line flags
taskfix "bug" --model "anthropic/claude-3-opus" --api-key "your-key"

# Use custom config file
taskfix "feature request" --config /path/to/custom/config

# Use custom rules
taskfix "bug" --rules /path/to/custom/rules.json
```
