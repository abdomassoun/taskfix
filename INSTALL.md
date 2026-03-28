# TaskFix Installation and Configuration Guide

## Installation

### Method 1: Install from .deb Package (Recommended for Debian/Ubuntu)

```bash
# Download the latest release
wget https://github.com/taskfix/taskfix/releases/latest/download/taskfix_latest_amd64.deb

# Install the package
sudo dpkg -i taskfix_latest_amd64.deb

# If there are dependency issues, run:
sudo apt-get install -f
```

The package installs:
- **Binary**: `/usr/local/bin/taskfix`
- **System config**: `/etc/taskfix/config` (template with empty API key)

### Method 2: Build from Source

```bash
git clone https://github.com/taskfix/taskfix
cd taskfix
make build

# Optional: Install system-wide
sudo cp taskfix /usr/local/bin/
sudo mkdir -p /etc/taskfix
sudo cp taskfix-deb/etc/taskfix/config /etc/taskfix/
```

## Configuration

TaskFix follows Linux best practices for configuration management. You have three configuration locations to choose from:

### 1. User Configuration: `~/.tfixrc` (Recommended)

This is the recommended location for individual user settings:

```bash
# Create your personal config
cat > ~/.tfixrc << 'EOF'
{
  "provider": "openrouter",
  "api_key": "sk-or-v1-your-api-key-here",
  "model": "openai/gpt-4o-mini",
  "rules_file": ""
}
EOF

# Protect your API key
chmod 600 ~/.tfixrc
```

### 2. XDG Config: `~/.config/taskfix/config`

Follow the XDG Base Directory specification:

```bash
# Create config directory
mkdir -p ~/.config/taskfix

# Create config file
cat > ~/.config/taskfix/config << 'EOF'
{
  "provider": "openrouter",
  "api_key": "sk-or-v1-your-api-key-here",
  "model": "anthropic/claude-3-haiku",
  "rules_file": ""
}
EOF

# Protect your API key
chmod 600 ~/.config/taskfix/config
```

### 3. System-wide: `/etc/taskfix/config`

For system administrators managing multiple users:

```bash
# Edit system config (requires sudo)
sudo vim /etc/taskfix/config

# Example system config (don't put API keys here!)
{
  "provider": "openrouter",
  "model": "openai/gpt-4o-mini",
  "api_key": "",
  "rules_file": "/etc/taskfix/rules.json"
}
```

**Security Note**: Never put API keys in `/etc/taskfix/config`. Users should provide their own keys via personal config or environment variables.

## Configuration Priority

TaskFix uses the following priority (highest to lowest):

1. **CLI flags**: `--api-key`, `--model`, `--config`
2. **Environment variables**: `OPENROUTER_API_KEY`
3. **User config**: `~/.tfixrc`
4. **XDG config**: `~/.config/taskfix/config`
5. **System config**: `/etc/taskfix/config`

## Quick Start Examples

### Example 1: Using ~/.tfixrc (Simplest)

```bash
# 1. Create config
echo '{"api_key": "sk-or-v1-xxx"}' > ~/.tfixrc
chmod 600 ~/.tfixrc

# 2. Use taskfix
taskfix "user can't login when password is wrong"
```

### Example 2: Using Environment Variable

```bash
# 1. Set environment variable (add to ~/.bashrc for persistence)
export OPENROUTER_API_KEY="sk-or-v1-your-key"

# 2. Use taskfix
taskfix "login bug"
```

### Example 3: Using Custom Config File

```bash
# 1. Create custom config anywhere
cat > ~/my-taskfix-config << 'EOF'
{
  "api_key": "sk-or-v1-xxx",
  "model": "anthropic/claude-3-opus",
  "rules_file": "~/my-rules.json"
}
EOF

# 2. Use with --config flag
taskfix "bug description" --config ~/my-taskfix-config
```

## Configuration Fields

All fields are optional:

| Field | Description | Default |
|-------|-------------|---------|
| `provider` | AI provider | `"openrouter"` |
| `api_key` | Your API key | From `OPENROUTER_API_KEY` env |
| `model` | AI model to use | `"openai/gpt-4o-mini"` |
| `rules_file` | Path to custom rules | Built-in rules |

## Getting an API Key

1. Visit [OpenRouter](https://openrouter.ai)
2. Sign up or log in
3. Go to [API Keys](https://openrouter.ai/keys)
4. Create a new key
5. Copy the key (starts with `sk-or-v1-`)

## Verifying Installation

```bash
# Check version
taskfix version

# List available models
taskfix --models

# Test with a simple task (requires API key configured)
taskfix "test bug description"
```

## Security Best Practices

1. **Never commit API keys** to version control
2. **Use file permissions** to protect config files:
   ```bash
   chmod 600 ~/.tfixrc
   ```
3. **Use environment variables** for CI/CD:
   ```bash
   OPENROUTER_API_KEY="xxx" taskfix "bug"
   ```
4. **System configs** should not contain API keys
5. **Rotate keys** regularly via OpenRouter dashboard

## Troubleshooting

### "no API key provided"

Make sure you've set your API key via:
- Config file (`~/.tfixrc`)
- Environment variable (`OPENROUTER_API_KEY`)
- CLI flag (`--api-key`)

### "config file not found"

TaskFix auto-discovers config files. Create one of:
- `~/.tfixrc`
- `~/.config/taskfix/config`
- `/etc/taskfix/config`

Or specify explicitly:
```bash
taskfix --config /path/to/config "task description"
```

### Permission denied

Make sure the binary is executable:
```bash
chmod +x /usr/local/bin/taskfix
```

### Package installation fails

```bash
# Fix dependencies
sudo apt-get install -f

# Or install manually
sudo dpkg -i taskfix_*_amd64.deb
```

## Uninstallation

### If installed from .deb:
```bash
sudo dpkg -r taskfix
```

### If installed manually:
```bash
sudo rm /usr/local/bin/taskfix
sudo rm -rf /etc/taskfix
rm ~/.tfixrc
rm -rf ~/.config/taskfix
```

## More Information

- [Configuration Examples](README.md) - Detailed config examples
- [Custom Rules](../README.md#custom-rules-format) - Creating custom formatting rules
- [GitHub Repository](https://github.com/taskfix/taskfix)
- [OpenRouter Documentation](https://openrouter.ai/docs)
