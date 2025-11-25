package provider

import (
	"context"
	"time"
)

// Provider is the interface that all AI CLI providers must implement
type Provider interface {
	// Name returns the provider identifier
	Name() string

	// Available checks if the CLI is installed and accessible
	Available() bool

	// Execute runs the prompt and returns the response
	Execute(ctx context.Context, prompt string, opts ExecuteOptions) (string, error)
}

// ExecuteOptions contains options for executing a prompt
type ExecuteOptions struct {
	Timeout    time.Duration
	WorkingDir string
	Verbose    bool
}

// ProviderInfo contains information about a provider
type ProviderInfo struct {
	Name      string
	Command   string
	Version   string
	Available bool
}
