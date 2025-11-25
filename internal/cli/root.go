package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/trknhr/ai-utils/internal/config"
	"github.com/trknhr/ai-utils/internal/provider"
)

var (
	cfg      *config.Config
	detector *provider.Detector
)

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "aiu",
	Short: "AI Utils - Execute prompts with local AI CLIs",
	Long: `AI Utils (aiu) is a command-line tool that manages prompt templates
and executes them using locally installed AI CLIs (claude, gemini, codex).

It automatically detects available AI CLIs and uses them according to
your configured priority, eliminating the need for API keys.`,
	SilenceUsage:  true,
	SilenceErrors: true,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().StringP("provider", "P", "", "specify provider (claude, gemini, codex)")
	rootCmd.PersistentFlags().StringP("working-dir", "w", ".", "working directory for command execution")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	var err error
	cfg, err = config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Failed to load config: %v\n", err)
		fmt.Fprintf(os.Stderr, "Using default configuration.\n")
		// Create default config
		cfg = &config.Config{
			ProviderPriority: []string{"claude", "gemini", "codex"},
			Providers: map[string]config.ProviderConfig{
				"claude": {
					Command: "claude",
					Args:    []string{"-p", "--output-format", "text"},
					Timeout: 120,
				},
			},
		}
	}

	detector = provider.NewDetector(cfg)
}
