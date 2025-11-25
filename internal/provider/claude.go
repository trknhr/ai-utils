package provider

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/trknhr/ai-utils/internal/config"
)

// ClaudeProvider implements the Provider interface for Claude CLI
type ClaudeProvider struct {
	command string
	args    []string
	timeout time.Duration
}

// NewClaudeProvider creates a new Claude provider
func NewClaudeProvider(cfg *config.ProviderConfig) *ClaudeProvider {
	return &ClaudeProvider{
		command: cfg.Command,
		args:    cfg.Args,
		timeout: cfg.Timeout,
	}
}

// Name returns the provider identifier
func (p *ClaudeProvider) Name() string {
	return "claude"
}

// Available checks if the Claude CLI is installed
func (p *ClaudeProvider) Available() bool {
	_, err := exec.LookPath(p.command)
	return err == nil
}

// Execute runs a prompt through Claude CLI
func (p *ClaudeProvider) Execute(ctx context.Context, prompt string, opts ExecuteOptions) (string, error) {
	// Create context with timeout
	if opts.Timeout == 0 {
		opts.Timeout = p.timeout
	}

	ctx, cancel := context.WithTimeout(ctx, opts.Timeout)
	defer cancel()

	// Build command arguments
	// For Claude Code, we pass the prompt directly as an argument
	args := append(p.args, prompt)

	if opts.Verbose {
		fmt.Printf("Executing: %s %s\n", p.command, strings.Join(args, " "))
	}

	// Execute command
	cmd := exec.CommandContext(ctx, p.command, args...)
	if opts.WorkingDir != "" {
		cmd.Dir = opts.WorkingDir
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("claude execution failed: %w\nOutput: %s", err, string(output))
	}

	return strings.TrimSpace(string(output)), nil
}
