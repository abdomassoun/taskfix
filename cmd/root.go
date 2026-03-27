package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/taskfix/taskfix/internal/ai"
	"github.com/taskfix/taskfix/internal/output"
	"github.com/taskfix/taskfix/internal/prompt"
	"github.com/taskfix/taskfix/internal/rules"
)

var (
	flagConfig string
	flagRules  string
	flagFile   string
	flagModel  string
	flagAPIKey string
	flagModelsAll     bool
	flagSilent bool
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
  taskfix "bug" --rules configs/jira.json
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
	rootCmd.Flags().StringVarP(&flagAPIKey, "api-key", "k", "", "API key (overrides OPENROUTER_API_KEY env var)")
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

	// ── 1. Resolve input ──────────────────────────────────────────────────────
	input, err := resolveInput(args)
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
func resolveInput(args []string) (string, error) {
	// --file flag takes highest priority
	if flagFile != "" {
		data, err := os.ReadFile(flagFile)
		if err != nil {
			return "", fmt.Errorf("reading file %q: %w", flagFile, err)
		}
		return string(data), nil
	}

	// Check if stdin is a pipe/redirect (not an interactive terminal)
	stat, err := os.Stdin.Stat()
	if err != nil {
		return "", fmt.Errorf("checking stdin: %w", err)
	}
	if (stat.Mode() & os.ModeCharDevice) == 0 {
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

// logf writes a styled message to stderr, unless --silent is set.
func logf(format string, a ...any) {
	if flagSilent {
		return
	}
	color.New(color.FgHiBlack).Fprintf(os.Stderr, format, a...)
}
