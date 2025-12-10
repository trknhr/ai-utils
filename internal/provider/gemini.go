package provider

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/trknhr/ai-utils/internal/config"
)

// GeminiProvider implements the Provider interface for Gemini CLI
type GeminiProvider struct {
	command string
	args    []string
	timeout time.Duration
}

// NewGeminiProvider creates a new Gemini provider
func NewGeminiProvider(cfg *config.ProviderConfig) *GeminiProvider {
	return &GeminiProvider{
		command: cfg.Command,
		args:    cfg.Args,
		timeout: cfg.Timeout,
	}
}

// Name returns the provider identifier
func (p *GeminiProvider) Name() string {
	return "gemini"
}

// Available checks if the Gemini CLI is installed
func (p *GeminiProvider) Available() bool {
	_, err := exec.LookPath(p.command)
	return err == nil
}

// Execute runs a prompt through Gemini CLI
func (p *GeminiProvider) Execute(ctx context.Context, prompt string, opts ExecuteOptions) (string, error) {
	if opts.Timeout == 0 {
		opts.Timeout = p.timeout
	}

	ctx, cancel := context.WithTimeout(ctx, opts.Timeout)
	defer cancel()

	// Gemini CLI uses -p flag for prompt
	args := append(p.args, "-p", prompt)

	if opts.Verbose {
		fmt.Printf("Executing: %s %s\n", p.command, strings.Join(args, " "))
	}

	cmd := exec.CommandContext(ctx, p.command, args...)
	if opts.WorkingDir != "" {
		cmd.Dir = opts.WorkingDir
	}

	output, err := cmd.CombinedOutput()
	clean := filterGeminiNoise(string(output))
	if err != nil {
		return "", fmt.Errorf("gemini execution failed: %w\nOutput: %s", err, clean)
	}

	return strings.TrimSpace(clean), nil
}

// filterGeminiNoise removes IDEClient warnings the user cannot act on here.
func filterGeminiNoise(out string) string {
	lines := strings.Split(out, "\n")
	filtered := make([]string, 0, len(lines))
	for _, line := range lines {
		if strings.Contains(line, "[ERROR] [IDEClient]") {
			continue
		}
		filtered = append(filtered, line)
	}
	return strings.Join(filtered, "\n")
}
