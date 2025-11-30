package config

import (
"os"
"path/filepath"
"strings"
"testing"
)

func TestExpandPath_HomeDir(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Skipf("could not get home dir: %v", err)
	}

	result := expandPath("~/test/path")
	expected := filepath.Join(home, "test/path")

	if result != expected {
		t.Errorf("expected '%s', got '%s'", expected, result)
	}
}

func TestExpandPath_EnvVar(t *testing.T) {
	os.Setenv("TEST_VAR", "testvalue")
	defer os.Unsetenv("TEST_VAR")

	result := expandPath("/path/$TEST_VAR/dir")
	expected := "/path/testvalue/dir"

	if result != expected {
		t.Errorf("expected '%s', got '%s'", expected, result)
	}
}

func TestExpandPath_NoExpansion(t *testing.T) {
	input := "/absolute/path/without/expansion"
	result := expandPath(input)

	if result != input {
		t.Errorf("expected '%s', got '%s'", input, result)
	}
}

func TestGetConfigDir(t *testing.T) {
	dir, err := GetConfigDir()
	if err != nil {
		t.Fatalf("GetConfigDir failed: %v", err)
	}

	if dir == "" {
		t.Error("GetConfigDir returned empty string")
	}

	if !strings.Contains(dir, "ai-utils") {
		t.Errorf("expected path to contain 'ai-utils', got '%s'", dir)
	}
}

func TestGetConfigPath(t *testing.T) {
	path, err := GetConfigPath()
	if err != nil {
		t.Fatalf("GetConfigPath failed: %v", err)
	}

	if path == "" {
		t.Error("GetConfigPath returned empty string")
	}

	if filepath.Base(path) != "config.yaml" {
		t.Errorf("expected path to end with 'config.yaml', got '%s'", path)
	}
}

func TestGetPromptsDir(t *testing.T) {
	dir, err := GetPromptsDir()
	if err != nil {
		t.Fatalf("GetPromptsDir failed: %v", err)
	}

	if dir == "" {
		t.Error("GetPromptsDir returned empty string")
	}

	if filepath.Base(dir) != "prompts" {
		t.Errorf("expected path to end with 'prompts', got '%s'", dir)
	}
}

func TestLoad_DefaultConfig(t *testing.T) {
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if cfg == nil {
		t.Fatal("Load returned nil config")
	}

	if len(cfg.ProviderPriority) == 0 {
		t.Error("expected provider_priority to have defaults")
	}

	if cfg.Providers == nil || len(cfg.Providers) == 0 {
		t.Error("expected providers to have defaults")
	}
}

func TestLoad_PromptDirsDefault(t *testing.T) {
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if len(cfg.PromptDirs) == 0 {
		t.Error("expected prompt_dirs to have default value")
	}
}
