package provider

import (
"testing"

"github.com/trknhr/ai-utils/internal/config"
)

func TestNewDetector(t *testing.T) {
	cfg := &config.Config{
		ProviderPriority: []string{"claude", "gemini"},
		Providers: map[string]config.ProviderConfig{
			"claude": {
				Command: "claude",
				Args:    []string{"-p"},
			},
		},
	}

	detector := NewDetector(cfg)

	if detector == nil {
		t.Fatal("NewDetector returned nil")
	}

	if len(detector.priority) != 2 {
		t.Errorf("expected priority length 2, got %d", len(detector.priority))
	}
}

func TestDetector_GetProvider_NotFound(t *testing.T) {
	cfg := &config.Config{
		ProviderPriority: []string{},
		Providers:        map[string]config.ProviderConfig{},
	}

	detector := NewDetector(cfg)

	_, err := detector.GetProvider("nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent provider")
	}
}

func TestDetector_DetectAvailable(t *testing.T) {
	cfg := &config.Config{
		ProviderPriority: []string{"claude"},
		Providers: map[string]config.ProviderConfig{
			"claude": {
				Command: "claude",
			},
		},
	}

	detector := NewDetector(cfg)
	available := detector.DetectAvailable()

	// Result depends on whether claude is installed
	_ = available
}

func TestDetector_ListProviders(t *testing.T) {
	cfg := &config.Config{
		ProviderPriority: []string{"claude"},
		Providers: map[string]config.ProviderConfig{
			"claude": {
				Command: "claude",
			},
		},
	}

	detector := NewDetector(cfg)
	providers := detector.ListProviders()

	if len(providers) == 0 {
		t.Error("expected at least one provider in list")
	}

	for _, p := range providers {
		if p.Name == "" {
			t.Error("provider has empty name")
		}
	}
}

func TestDetector_GetFirstAvailable_NoneAvailable(t *testing.T) {
	cfg := &config.Config{
		ProviderPriority: []string{"nonexistent"},
		Providers:        map[string]config.ProviderConfig{},
	}

	detector := NewDetector(cfg)
	_, err := detector.GetFirstAvailable()

	if err == nil {
		t.Error("expected error when no providers available")
	}
}
