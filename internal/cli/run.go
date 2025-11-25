package cli

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/trknhr/ai-utils/internal/provider"
	"github.com/trknhr/ai-utils/internal/template"
)

var (
	dryRun  bool
	timeout time.Duration
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run <prompt-name>",
	Short: "Execute a prompt template",
	Long: `Execute a prompt template by name.

The template is loaded from configured prompt directories, any {{$ command }}
placeholders are expanded, and the result is sent to the selected AI provider.

Examples:
  aiu run pr-desc
  aiu run commit-msg --provider claude
  aiu run review --dry-run`,
	Args: cobra.ExactArgs(1),
	RunE: runPrompt,
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().BoolVar(&dryRun, "dry-run", false, "show expanded prompt without executing")
	runCmd.Flags().DurationVar(&timeout, "timeout", 0, "execution timeout (default: from config)")
}

func runPrompt(cmd *cobra.Command, args []string) error {
	promptName := args[0]

	verbose, _ := cmd.Flags().GetBool("verbose")
	providerName, _ := cmd.Flags().GetString("provider")
	workingDir, _ := cmd.Flags().GetString("working-dir")

	// Find template
	if verbose {
		fmt.Printf("Looking for template: %s\n", promptName)
	}

	templatePath, err := template.FindTemplate(promptName, cfg.PromptDirs)
	if err != nil {
		return fmt.Errorf("template not found: %w", err)
	}

	if verbose {
		fmt.Printf("Found template: %s\n", templatePath)
	}

	// Parse template
	tmpl, err := template.Parse(templatePath)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	if verbose {
		fmt.Printf("Template:\n%s\n", tmpl.String())
	}

	// Validate template requirements
	if err := tmpl.Validate(); err != nil {
		return fmt.Errorf("template validation failed: %w", err)
	}

	// Expand template (execute commands)
	executor := template.NewExecutor(workingDir, verbose)

	if dryRun {
		// Just show what would be executed
		fmt.Println(executor.DryRun(tmpl))
		return nil
	}

	expandedPrompt, err := executor.Execute(tmpl)
	if err != nil {
		return fmt.Errorf("failed to expand template: %w", err)
	}

	if verbose {
		fmt.Printf("\n=== Expanded Prompt ===\n%s\n\n", expandedPrompt)
	}

	// Select provider
	var prov provider.Provider
	if providerName != "" {
		// User specified a provider
		prov, err = detector.GetProvider(providerName)
		if err != nil {
			return fmt.Errorf("failed to get provider: %w", err)
		}
		if verbose {
			fmt.Printf("Using specified provider: %s\n", providerName)
		}
	} else {
		// Auto-detect first available
		prov, err = detector.GetFirstAvailable()
		if err != nil {
			return fmt.Errorf("no available providers: %w\n\nPlease install one of: claude, gemini, codex", err)
		}
		if verbose {
			fmt.Printf("Using auto-detected provider: %s\n", prov.Name())
		}
	}

	// Execute with provider
	ctx := context.Background()
	opts := provider.ExecuteOptions{
		Timeout:    timeout,
		WorkingDir: workingDir,
		Verbose:    verbose,
	}

	fmt.Fprintf(os.Stderr, "Executing with %s...\n", prov.Name())

	result, err := prov.Execute(ctx, expandedPrompt, opts)
	if err != nil {
		return fmt.Errorf("execution failed: %w", err)
	}

	// Print result
	fmt.Println(result)

	return nil
}
