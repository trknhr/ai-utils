package provider

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/trknhr/ai-utils/internal/config"
)

// CopilotProvider implements the Provider interface for GitHub Copilot CLI.
type CopilotProvider struct {
	command string
	args    []string
	timeout time.Duration
	model   string
}

// NewCopilotProvider creates a new Copilot provider.
func NewCopilotProvider(cfg *config.ProviderConfig) *CopilotProvider {
	return &CopilotProvider{
		command: cfg.Command,
		args:    cfg.Args,
		timeout: cfg.Timeout,
		model:   strings.TrimSpace(cfg.Model),
	}
}

func (p *CopilotProvider) Name() string { return ProviderCopilot }

func (p *CopilotProvider) Available() bool {
	_, err := exec.LookPath(p.command)
	return err == nil
}

func (p *CopilotProvider) Execute(ctx context.Context, prompt string, opts ExecuteOptions) (string, error) {
	if opts.Timeout == 0 {
		opts.Timeout = p.timeout
	}

	ctx, cancel := context.WithTimeout(ctx, opts.Timeout)
	defer cancel()

	if strings.TrimSpace(prompt) == "" {
		return "", errors.New("prompt is empty")
	}

	effectiveModel := strings.TrimSpace(opts.Model)
	if effectiveModel == "" {
		effectiveModel = p.model
	}

	baseArgs := append([]string{}, p.args...)
	baseArgs = stripFlagWithValue(baseArgs, flagModel)
	if effectiveModel != "" {
		baseArgs = append(baseArgs, flagModel, effectiveModel)
	}

	args := append(baseArgs, "-p", prompt)

	if opts.Verbose {
		fmt.Printf("Executing: %s %s\n", p.command, strings.Join(args, " "))
	}

	cmd := exec.CommandContext(ctx, p.command, args...)
	if opts.WorkingDir != "" {
		cmd.Dir = opts.WorkingDir
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("copilot execution failed: %w\nOutput: %s", err, string(output))
	}

	return strings.TrimSpace(string(output)), nil
}
