package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/taskfix/taskfix/internal/ai"
	"github.com/taskfix/taskfix/internal/output"
	"github.com/taskfix/taskfix/internal/prompt"
	"github.com/taskfix/taskfix/internal/rules"
)

var (
	flagConfig    string
	flagRules     string
	flagFile      string
	flagModel     string
	flagAPIKey    string
	flagModelsAll bool
	flagSilent    bool
)

var rootCmd = &cobra.Command{
	Use:   "taskfix [text]",
	Short: "Transform raw task descriptions into structured tasks using AI",
	Long: `TaskFix turns messy task descriptions into clean, structured, AI-formatted
technical tasks — suitable for GitHub Issues, Jira, or any task tracker.

Examples:
  taskfix "user cant login when password wrong"
  echo "login bug" | taskfix
  taskfix -f input.txt
  taskfix "bug" --rules /etc/taskfix/config.d/jira.json
  taskfix "bug" --config taskfix.json
  taskfix "login bug" | gh issue create --title "Bug" --body -`,
	Args: cobra.MaximumNArgs(1),
	RunE: run,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&flagConfig, "config", "c", "", "path to config JSON file")
	rootCmd.Flags().StringVarP(&flagRules, "rules", "r", "", "path to custom rules JSON file")
	rootCmd.Flags().StringVarP(&flagFile, "file", "f", "", "read input from file instead of argument/stdin")
	rootCmd.Flags().StringVarP(&flagModel, "model", "m", "", "AI model override (e.g. anthropic/claude-3-haiku)")
	rootCmd.Flags().BoolVar(&flagModelsAll, "models", false, "list available models; optional pattern may be provided as first positional argument")
	rootCmd.Flags().StringVarP(&flagAPIKey, "api-key", "k", "", "API key (overrides OPENROUTER_API_KEY env var; when used alone, saves to config)")
	rootCmd.Flags().BoolVarP(&flagSilent, "silent", "s", false, "suppress stderr progress output")
}

func run(cmd *cobra.Command, args []string) error {
	// If the user asked to list models, do that and exit early.
	if flagModelsAll {
		cfg, err := loadConfig(flagConfig)
		if err != nil {
			return fmt.Errorf("loading config: %w", err)
		}
		applyFlagOverrides(cfg)

		client := ai.NewClient(cfg.Provider, cfg.APIKey, cfg.Model)
		models, err := client.FetchModels()
		if err != nil {
			return fmt.Errorf("fetching models: %w", err)
		}
		// Accept optional pattern as first positional argument when
		// `--models` is used. Example: `--models free`.
		pattern := ""
		if len(args) > 0 {
			pattern = strings.ToLower(args[0])
		}
		for _, m := range models {
			if pattern == "" {
				fmt.Println(m)
				continue
			}
			if strings.Contains(strings.ToLower(m), pattern) {
				fmt.Println(m)
			}
		}
		return nil
	}

	stdinPiped, err := isStdinPiped()
	if err != nil {
		return err
	}

	// If only --api-key is provided, persist it to config and exit.
	if flagAPIKey != "" && len(args) == 0 && flagFile == "" && !stdinPiped {
		configPath := defaultConfigPath(flagConfig)
		if err := saveAPIKey(configPath, flagAPIKey); err != nil {
			return fmt.Errorf("saving api key: %w", err)
		}
		fmt.Fprintf(os.Stderr, "API key saved to %s\n", configPath)
		return nil
	}

	// ── 1. Resolve input ──────────────────────────────────────────────────────
	input, err := resolveInput(args, stdinPiped)
	if err != nil {
		return err
	}
	input = strings.TrimSpace(input)
	if input == "" {
		return fmt.Errorf(
			"no input provided\n\n" +
				"  Pass text as an argument:  taskfix \"bug description\"\n" +
				"  Pipe via stdin:            echo \"bug\" | taskfix\n" +
				"  Read from file:            taskfix -f input.txt\n\n" +
				"Run 'taskfix --help' for full usage",
		)
	}

	// ── 2. Load config ────────────────────────────────────────────────────────
	cfg, err := loadConfig(flagConfig)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}
	applyFlagOverrides(cfg)

	// ── 3. Load rules ─────────────────────────────────────────────────────────
	rulesPath := flagRules
	if rulesPath == "" {
		rulesPath = cfg.RulesFile
	}
	ruleSet, err := rules.Load(rulesPath)
	if err != nil {
		return fmt.Errorf("loading rules: %w", err)
	}

	// ── 4. Build prompt ───────────────────────────────────────────────────────
	p := prompt.Build(input, ruleSet)

	// ── 5. Call AI ────────────────────────────────────────────────────────────
	logf("→ Sending to AI (%s)...\n", cfg.Model)

	client := ai.NewClient(cfg.Provider, cfg.APIKey, cfg.Model)
	result, err := client.Complete(p)
	if err != nil {
		return fmt.Errorf("AI request failed: %w", err)
	}

	// ── 6. Format & write to stdout ───────────────────────────────────────────
	formatted := output.Format(result)
	fmt.Println(formatted)
	return nil
}

// resolveInput returns the task text from -f flag, stdin pipe, or positional arg.
// Priority: --file > stdin pipe > argument
func resolveInput(args []string, stdinPiped bool) (string, error) {
	// --file flag takes highest priority
	if flagFile != "" {
		data, err := os.ReadFile(flagFile)
		if err != nil {
			return "", fmt.Errorf("reading file %q: %w", flagFile, err)
		}
		return string(data), nil
	}

	// Check if stdin is a pipe/redirect (not an interactive terminal)
	if stdinPiped {
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			return "", fmt.Errorf("reading stdin: %w", err)
		}
		return string(data), nil
	}

	// Fall back to positional argument
	if len(args) > 0 {
		return args[0], nil
	}

	return "", nil
}

func isStdinPiped() (bool, error) {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return false, fmt.Errorf("checking stdin: %w", err)
	}
	return (stat.Mode() & os.ModeCharDevice) == 0, nil
}

func defaultConfigPath(path string) string {
	if path == "" {
		return filepath.Clean(os.ExpandEnv("$HOME/.tfixrc"))
	}

	if strings.HasPrefix(path, "~/") {
		return filepath.Join(os.ExpandEnv("$HOME"), strings.TrimPrefix(path, "~/"))
	}

	return filepath.Clean(os.ExpandEnv(path))
}

func saveAPIKey(path, key string) error {
	cfg := &Config{
		Provider: defaultProvider,
		Model:    defaultModel,
	}

	if data, err := os.ReadFile(path); err == nil {
		if err := json.Unmarshal(data, cfg); err != nil {
			return fmt.Errorf("parsing existing config %q: %w", path, err)
		}
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("reading config %q: %w", path, err)
	}

	cfg.APIKey = key

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("serializing config: %w", err)
	}
	data = append(data, '\n')

	parent := filepath.Dir(path)
	if parent != "." {
		if err := os.MkdirAll(parent, 0o700); err != nil {
			return fmt.Errorf("creating config directory %q: %w", parent, err)
		}
	}

	if err := os.WriteFile(path, data, 0o600); err != nil {
		return fmt.Errorf("writing config %q: %w", path, err)
	}
	return nil
}

// logf writes a styled message to stderr, unless --silent is set.
func logf(format string, a ...any) {
	if flagSilent {
		return
	}
	color.New(color.FgHiBlack).Fprintf(os.Stderr, format, a...)
}
