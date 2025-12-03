package provider

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/trknhr/ai-utils/internal/config"
)

// CodexProvider implements the Provider interface for OpenAI Codex CLI
type CodexProvider struct {
	command string
	args    []string
	timeout time.Duration
}

// NewCodexProvider creates a new Codex provider
func NewCodexProvider(cfg *config.ProviderConfig) *CodexProvider {
	return &CodexProvider{
		command: cfg.Command,
		args:    cfg.Args,
		timeout: cfg.Timeout,
	}
}

// Name returns the provider identifier
func (p *CodexProvider) Name() string {
	return "codex"
}

// Available checks if the Codex CLI is installed
func (p *CodexProvider) Available() bool {
	_, err := exec.LookPath(p.command)
	return err == nil
}

// Execute runs a prompt through Codex CLI
func (p *CodexProvider) Execute(ctx context.Context, prompt string, opts ExecuteOptions) (string, error) {
	if opts.Timeout == 0 {
		opts.Timeout = p.timeout
	}

	ctx, cancel := context.WithTimeout(ctx, opts.Timeout)
	defer cancel()

	// Codex CLI passes prompt as argument
	args := append(p.args, prompt)

	if opts.Verbose {
		fmt.Printf("Executing: %s %s\n", p.command, strings.Join(args, " "))
	}

	cmd := exec.CommandContext(ctx, p.command, args...)
	if opts.WorkingDir != "" {
		cmd.Dir = opts.WorkingDir
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("codex execution failed: %w\nOutput: %s", err, string(output))
	}

	return strings.TrimSpace(string(output)), nil
}
