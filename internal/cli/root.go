package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/trknhr/ai-utils/internal/config"
	"github.com/trknhr/ai-utils/internal/provider"
	"github.com/trknhr/ai-utils/internal/runner"
	"github.com/trknhr/ai-utils/internal/template"
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
	rootCmd.PersistentFlags().String("lang", "", "output language (overrides config output_lang)")

	// Register template commands at init time (so they show in --help)
	registerTemplateCommands()
}

// registerTemplateCommands dynamically registers builtin templates as subcommands
func registerTemplateCommands() {
	for _, tmpl := range template.ListBuiltinTemplates() {
		// Capture the template in closure
		t := tmpl
		cmd := &cobra.Command{
			Use:   t.Name + " [args...]",
			Short: t.Description,
			Long:  fmt.Sprintf("Execute the %s template.\n\n%s", t.Name, t.Description),
			RunE: func(cmd *cobra.Command, args []string) error {
				return executeTemplate(cmd, t.Name, args)
			},
		}
		cmd.Flags().Bool("dry-run", false, "show expanded prompt without executing")
		cmd.Flags().Duration("timeout", 0, "execution timeout (default: from config)")
		cmd.Flags().BoolP("parallel", "p", false, "run available providers in parallel and choose the output")
		rootCmd.AddCommand(cmd)
	}
}

// executeTemplate is the shared logic for running a template
func executeTemplate(cmd *cobra.Command, promptName string, args []string) error {
	verbose, _ := cmd.Flags().GetBool("verbose")
	providerName, _ := cmd.Flags().GetString("provider")
	workingDir, _ := cmd.Flags().GetString("working-dir")
	dryRun, _ := cmd.Flags().GetBool("dry-run")
	timeout, _ := cmd.Flags().GetDuration("timeout")
	parallel, _ := cmd.Flags().GetBool("parallel")
	langFlag, _ := cmd.Flags().GetString("lang")

	// Find template
	if verbose {
		fmt.Printf("Looking for template: %s\n", promptName)
	}

	var tmpl *template.Template
	var err error

	// First, try to find in user directories
	templatePath, findErr := template.FindTemplate(promptName, cfg.PromptDirs)
	if findErr == nil {
		if verbose {
			fmt.Printf("Found template: %s\n", templatePath)
		}
		tmpl, err = template.Parse(templatePath)
		if err != nil {
			return fmt.Errorf("failed to parse template: %w", err)
		}
	} else {
		// Fallback to builtin templates
		var ok bool
		tmpl, ok = template.GetBuiltinTemplate(promptName)
		if !ok {
			return fmt.Errorf("template not found: %s", promptName)
		}
		if verbose {
			fmt.Printf("Using builtin template: %s\n", promptName)
		}
	}

	if verbose {
		fmt.Printf("Template:\n%s\n", tmpl.String())
	}

	// Validate template requirements
	if err := tmpl.Validate(); err != nil {
		return fmt.Errorf("template validation failed: %w", err)
	}

	// Expand template (execute commands)
	executor := template.NewExecutor(workingDir, verbose, args)

	expandedPrompt, err := executor.Execute(tmpl)
	if err != nil {
		return fmt.Errorf("failed to expand template: %w", err)
	}

	expandedPrompt = appendLanguageHint(expandedPrompt, langFlag, cfg.OutputLang)

	if dryRun {
		fmt.Println(expandedPrompt)
		return nil
	}

	if verbose {
		fmt.Printf("\n=== Expanded Prompt ===\n%s\n\n", expandedPrompt)
	}

	// Select providers
	providers, err := resolveProviders(parallel, providerName, verbose)
	if err != nil {
		return err
	}

	// Execute with runner
	result, err := runner.Run(cmd.Context(), runner.ExecuteOptions{
		Parallel:   parallel,
		Prompt:     expandedPrompt,
		Providers:  providers,
		Timeout:    timeout,
		WorkingDir: workingDir,
		Verbose:    verbose,
	})
	if err != nil {
		return fmt.Errorf("execution failed: %w", err)
	}

	fmt.Println(result)

	return nil
}

func resolveProviders(parallel bool, providerName string, verbose bool) ([]provider.Provider, error) {
	if providerName != "" {
		prov, err := detector.GetProvider(providerName)
		if err != nil {
			return nil, fmt.Errorf("failed to get provider: %w", err)
		}
		if verbose {
			fmt.Fprintf(os.Stderr, "Using specified provider: %s\n", providerName)
		}
		return []provider.Provider{prov}, nil
	}

	if parallel {
		availableNames := detector.DetectAvailable()
		if len(availableNames) == 0 {
			return nil, fmt.Errorf("no available providers\n\nPlease install one of: claude, gemini, codex")
		}

		providers := make([]provider.Provider, 0, len(availableNames))
		for _, name := range availableNames {
			prov, err := detector.GetProvider(name)
			if err != nil {
				continue
			}
			providers = append(providers, prov)
		}

		if len(providers) == 0 {
			return nil, fmt.Errorf("no available providers\n\nPlease install one of: claude, gemini, codex")
		}

		if verbose {
			names := make([]string, 0, len(providers))
			for _, prov := range providers {
				names = append(names, prov.Name())
			}
			fmt.Fprintf(os.Stderr, "Running in parallel with: %s\n", strings.Join(names, ", "))
		}

		return providers, nil
	}

	prov, err := detector.GetFirstAvailable()
	if err != nil {
		return nil, fmt.Errorf("no available providers: %w\n\nPlease install one of: claude, gemini, codex", err)
	}
	if verbose {
		fmt.Fprintf(os.Stderr, "Using auto-detected provider: %s\n", prov.Name())
	}
	return []provider.Provider{prov}, nil
}

func appendLanguageHint(prompt string, flagLang string, configLang string) string {
	lang := strings.TrimSpace(flagLang)
	if lang == "" {
		lang = strings.TrimSpace(configLang)
	}
	if lang == "" {
		return prompt
	}

	var sb strings.Builder
	sb.WriteString(strings.TrimSpace(prompt))
	sb.WriteString("\n\nPlease respond in ")
	sb.WriteString(lang)
	sb.WriteString(".")
	return sb.String()
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
					Args: []string{
						"-p",
						"--output-format", "text",
						"--settings", `{"attribution":{"commit":"","pr":""},"includeCoAuthoredBy":false,"gitAttribution":false}`,
					},
					Timeout: 120,
				},
				"gemini": {
					Command: "gemini",
					Args:    []string{},
					Timeout: 120,
				},
				"codex": {
					Command: "codex",
					Args:    []string{},
					Timeout: 120,
				},
			},
		}
	}

	detector = provider.NewDetector(cfg)
}
