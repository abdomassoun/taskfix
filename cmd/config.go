package cmd

import (
	"encoding/json"
	"os"
)

// Config holds all runtime configuration for TaskFix.
type Config struct {
	Provider  string `json:"provider"`   // "openrouter" (default)
	APIKey    string `json:"api_key"`    // falls back to OPENROUTER_API_KEY env var
	Model     string `json:"model"`      // e.g. "openai/gpt-4o-mini"
	RulesFile string `json:"rules_file"` // optional path to custom rules JSON
}

const (
	defaultProvider = "openrouter"
	defaultModel    = "openai/gpt-4o-mini"
)

// loadConfig loads a JSON config file if provided, otherwise returns defaults.
func loadConfig(path string) (*Config, error) {
	cfg := &Config{
		Provider: defaultProvider,
		Model:    defaultModel,
	}

	// Resolve API key from environment if not set
	cfg.APIKey = os.Getenv("OPENROUTER_API_KEY")

	if path == "" {
		// Try default config locations (user-local, XDG, system)
		for _, candidate := range []string{
			"taskfix.json",
			os.ExpandEnv("$HOME/.taskfixrc"),
			os.ExpandEnv("$HOME/.config/taskfix/config.json"),
			"/etc/taskfix/config.json",
		} {
			if _, err := os.Stat(candidate); err == nil {
				path = candidate
				break
			}
		}
	}

	if path == "" {
		return cfg, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

// applyFlagOverrides applies CLI flag values on top of the loaded config.
// Flags always win over config file values.
func applyFlagOverrides(cfg *Config) {
	if flagAPIKey != "" {
		cfg.APIKey = flagAPIKey
	}
	if flagModel != "" {
		cfg.Model = flagModel
	}
}
