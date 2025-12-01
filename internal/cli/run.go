package cli

import (
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run <prompt-name> [args...]",
	Short: "Execute a prompt template",
	Long: `Execute a prompt template by name.

The template is loaded from configured prompt directories, any {{$ command }}
placeholders are expanded, and the result is sent to the selected AI provider.

Additional arguments are passed to the template as environment variables:
  $1, $2, $3... for positional arguments
  $ARGS for all arguments joined by space

Examples:
  aiu run pr-desc
  aiu run pr-review develop          # $1=develop, compare with origin/develop
  aiu run commit-msg --provider claude
  aiu run review --dry-run`,
	Args: cobra.MinimumNArgs(1),
	RunE: runPrompt,
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().Bool("dry-run", false, "show expanded prompt without executing")
	runCmd.Flags().Duration("timeout", 0, "execution timeout (default: from config)")
}

func runPrompt(cmd *cobra.Command, args []string) error {
	promptName := args[0]
	templateArgs := args[1:]
	return executeTemplate(cmd, promptName, templateArgs)
}
