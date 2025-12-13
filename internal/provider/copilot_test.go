package provider

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/trknhr/ai-utils/internal/config"
)

func TestNewCopilotProvider(t *testing.T) {
	cfg := &config.ProviderConfig{
		Command: "copilot",
		Args:    []string{"-s"},
		Model:   "gpt-5",
		Timeout: 60 * time.Second,
	}

	prov := NewCopilotProvider(cfg)
	if prov == nil {
		t.Fatal("NewCopilotProvider returned nil")
	}
	if prov.command != "copilot" {
		t.Errorf("expected command 'copilot', got %q", prov.command)
	}
	if prov.model != "gpt-5" {
		t.Errorf("expected model 'gpt-5', got %q", prov.model)
	}
}

func TestCopilotProvider_Name(t *testing.T) {
	prov := &CopilotProvider{command: "copilot"}
	if prov.Name() != "copilot" {
		t.Fatalf("unexpected name: %q", prov.Name())
	}
}

func TestCopilotProvider_Execute_PromptEmpty(t *testing.T) {
	prov := &CopilotProvider{
		command: "echo",
		timeout: time.Second,
	}
	_, err := prov.Execute(context.Background(), "   ", ExecuteOptions{})
	if err == nil {
		t.Fatal("expected error for empty prompt")
	}
}

func TestCopilotProvider_Execute_WithMockScript_ModelOverride(t *testing.T) {
	tmpDir := t.TempDir()
	mockScript := tmpDir + "/mock-copilot"

	// Echo all args so we can assert flag ordering/override behavior.
	scriptContent := `#!/bin/sh
echo "$@"
`
	if err := os.WriteFile(mockScript, []byte(scriptContent), 0755); err != nil {
		t.Fatalf("failed to create mock script: %v", err)
	}

	prov := &CopilotProvider{
		command: mockScript,
		args:    []string{"-s", "--model", "old-model"},
		model:   "config-model",
		timeout: 10 * time.Second,
	}

	out, err := prov.Execute(context.Background(), "hello", ExecuteOptions{Model: "cli-model"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(out, "--model cli-model") {
		t.Fatalf("expected overridden model in args, got: %q", out)
	}
	if strings.Contains(out, "old-model") {
		t.Fatalf("expected old model removed, got: %q", out)
	}
	if strings.Contains(out, "config-model") {
		t.Fatalf("expected config model ignored when CLI model set, got: %q", out)
	}
	if !strings.Contains(out, "-p hello") {
		t.Fatalf("expected prompt passed with -p, got: %q", out)
	}
}

func TestCopilotProvider_Execute_UsesConfigModel(t *testing.T) {
	tmpDir := t.TempDir()
	mockScript := tmpDir + "/mock-copilot"

	scriptContent := `#!/bin/sh
echo "$@"
`
	if err := os.WriteFile(mockScript, []byte(scriptContent), 0755); err != nil {
		t.Fatalf("failed to create mock script: %v", err)
	}

	prov := &CopilotProvider{
		command: mockScript,
		args:    []string{"-s"},
		model:   "config-model",
		timeout: 10 * time.Second,
	}

	out, err := prov.Execute(context.Background(), "hello", ExecuteOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "--model config-model") {
		t.Fatalf("expected config model in args, got: %q", out)
	}
}

