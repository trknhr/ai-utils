package provider

import (
	"context"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/trknhr/ai-utils/internal/config"
)

func TestNewCodexProvider(t *testing.T) {
	cfg := &config.ProviderConfig{
		Command: "codex",
		Args:    []string{"--model", "gpt-5"},
		Timeout: 60 * time.Second,
	}

	provider := NewCodexProvider(cfg)

	if provider == nil {
		t.Fatal("NewCodexProvider returned nil")
	}

	if provider.command != "codex" {
		t.Errorf("expected command 'codex', got '%s'", provider.command)
	}

	if len(provider.args) != 2 {
		t.Errorf("expected 2 args, got %d", len(provider.args))
	}

	if provider.timeout != 60*time.Second {
		t.Errorf("expected timeout 60s, got %v", provider.timeout)
	}
}

func TestCodexProvider_Name(t *testing.T) {
	provider := &CodexProvider{
		command: "codex",
	}

	if name := provider.Name(); name != "codex" {
		t.Errorf("expected name 'codex', got '%s'", name)
	}
}

func TestCodexProvider_Available(t *testing.T) {
	tests := []struct {
		name     string
		command  string
		expected bool
	}{
		{
			name:     "existing command",
			command:  "echo", // echo is always available
			expected: true,
		},
		{
			name:     "non-existing command",
			command:  "nonexistent-command-12345",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider := &CodexProvider{
				command: tt.command,
			}

			if got := provider.Available(); got != tt.expected {
				t.Errorf("Available() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestCodexProvider_Available_RealCodex(t *testing.T) {
	// Skip if codex is not installed
	_, err := exec.LookPath("codex")
	if err != nil {
		t.Skip("codex CLI not installed, skipping")
	}

	provider := &CodexProvider{
		command: "codex",
	}

	if !provider.Available() {
		t.Error("codex should be available when installed")
	}
}

func TestCodexProvider_Execute_CommandNotFound(t *testing.T) {
	provider := &CodexProvider{
		command: "nonexistent-codex-command",
		args:    []string{},
		timeout: 5 * time.Second,
	}

	ctx := context.Background()
	_, err := provider.Execute(ctx, "test prompt", ExecuteOptions{})

	if err == nil {
		t.Error("expected error when command not found")
	}
}

func TestCodexProvider_Execute_WithTimeout(t *testing.T) {
	provider := &CodexProvider{
		command: "sleep",
		args:    []string{},
		timeout: 100 * time.Millisecond,
	}

	ctx := context.Background()
	_, err := provider.Execute(ctx, "10", ExecuteOptions{
		Timeout: 100 * time.Millisecond,
	})

	if err == nil {
		t.Error("expected timeout error")
	}
}

func TestCodexProvider_Execute_WithMockScript(t *testing.T) {
	// Create a mock script that simulates codex behavior
	tmpDir := t.TempDir()
	mockScript := tmpDir + "/mock-codex"

	// Create mock script that writes to the output file specified by -o
	scriptContent := `#!/bin/sh
# Parse -o flag to get output file
OUTPUT_FILE=""
while [ $# -gt 0 ]; do
    case "$1" in
        -o)
            OUTPUT_FILE="$2"
            shift 2
            ;;
        exec)
            shift
            ;;
        *)
            shift
            ;;
    esac
done

# Write mock output to file
if [ -n "$OUTPUT_FILE" ]; then
    echo "feat(test): mock commit message" > "$OUTPUT_FILE"
fi
`
	if err := os.WriteFile(mockScript, []byte(scriptContent), 0755); err != nil {
		t.Fatalf("failed to create mock script: %v", err)
	}

	provider := &CodexProvider{
		command: mockScript,
		args:    []string{},
		timeout: 10 * time.Second,
	}

	ctx := context.Background()
	result, err := provider.Execute(ctx, "generate commit message", ExecuteOptions{})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "feat(test): mock commit message"
	if result != expected {
		t.Errorf("expected '%s', got '%s'", expected, result)
	}
}

func TestCodexProvider_Execute_WithWorkingDir(t *testing.T) {
	tmpDir := t.TempDir()
	mockScript := tmpDir + "/mock-codex"

	// Create mock script that writes working directory to output
	scriptContent := `#!/bin/sh
OUTPUT_FILE=""
while [ $# -gt 0 ]; do
    case "$1" in
        -o)
            OUTPUT_FILE="$2"
            shift 2
            ;;
        *)
            shift
            ;;
    esac
done

if [ -n "$OUTPUT_FILE" ]; then
    pwd > "$OUTPUT_FILE"
fi
`
	if err := os.WriteFile(mockScript, []byte(scriptContent), 0755); err != nil {
		t.Fatalf("failed to create mock script: %v", err)
	}

	workDir := t.TempDir()
	provider := &CodexProvider{
		command: mockScript,
		args:    []string{},
		timeout: 10 * time.Second,
	}

	ctx := context.Background()
	result, err := provider.Execute(ctx, "test", ExecuteOptions{
		WorkingDir: workDir,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != workDir {
		t.Errorf("expected working dir '%s', got '%s'", workDir, result)
	}
}

func TestCodexProvider_Execute_WithArgs(t *testing.T) {
	tmpDir := t.TempDir()
	mockScript := tmpDir + "/mock-codex"

	// Create mock script that captures all arguments
	scriptContent := `#!/bin/sh
OUTPUT_FILE=""
ALL_ARGS=""
while [ $# -gt 0 ]; do
    case "$1" in
        -o)
            OUTPUT_FILE="$2"
            shift 2
            ;;
        *)
            ALL_ARGS="$ALL_ARGS $1"
            shift
            ;;
    esac
done

if [ -n "$OUTPUT_FILE" ]; then
    echo "$ALL_ARGS" > "$OUTPUT_FILE"
fi
`
	if err := os.WriteFile(mockScript, []byte(scriptContent), 0755); err != nil {
		t.Fatalf("failed to create mock script: %v", err)
	}

	provider := &CodexProvider{
		command: mockScript,
		args:    []string{"--model", "gpt-5"},
		timeout: 10 * time.Second,
	}

	ctx := context.Background()
	result, err := provider.Execute(ctx, "test prompt", ExecuteOptions{})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should contain exec, custom args, and prompt
	if result == "" {
		t.Error("expected non-empty result")
	}
}
