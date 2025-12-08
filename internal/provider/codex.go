package provider

import (
	"context"
	"fmt"
	"os"
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

	// Create a temp file to capture the output
	tmpFile, err := os.CreateTemp("", "codex-output-*.txt")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	tmpPath := tmpFile.Name()
	tmpFile.Close()
	defer os.Remove(tmpPath)

	// Codex CLI uses "exec" subcommand for non-interactive mode
	// Use -o to output only the last message to a file
	args := []string{"exec", "-o", tmpPath}
	args = append(args, p.args...)
	args = append(args, prompt)

	if opts.Verbose {
		fmt.Printf("Executing: %s %s\n", p.command, strings.Join(args, " "))
	}

	cmd := exec.CommandContext(ctx, p.command, args...)
	if opts.WorkingDir != "" {
		cmd.Dir = opts.WorkingDir
	}

	// Run command, ignore stdout/stderr (session info goes there)
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("codex execution failed: %w", err)
	}

	// Read the output from the temp file
	output, err := os.ReadFile(tmpPath)
	if err != nil {
		return "", fmt.Errorf("failed to read codex output: %w", err)
	}

	return strings.TrimSpace(string(output)), nil
}
