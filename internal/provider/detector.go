package provider

import (
	"fmt"

	"github.com/trknhr/ai-utils/internal/config"
)

// Detector manages provider detection and selection
type Detector struct {
	providers map[string]Provider
	priority  []string
}

// NewDetector creates a new Detector with the given configuration
func NewDetector(cfg *config.Config) *Detector {
	d := &Detector{
		providers: make(map[string]Provider),
		priority:  cfg.ProviderPriority,
	}

	// Register providers
	if claudeCfg, ok := cfg.Providers["claude"]; ok {
		d.providers["claude"] = NewClaudeProvider(&claudeCfg)
	}

	// TODO: Add gemini and codex providers in Phase 2

	return d
}

// DetectAvailable returns a list of available providers
func (d *Detector) DetectAvailable() []string {
	var available []string
	for _, name := range d.priority {
		if provider, ok := d.providers[name]; ok && provider.Available() {
			available = append(available, name)
		}
	}
	return available
}

// GetProvider returns a provider by name
func (d *Detector) GetProvider(name string) (Provider, error) {
	provider, ok := d.providers[name]
	if !ok {
		return nil, fmt.Errorf("provider %s not found", name)
	}

	if !provider.Available() {
		return nil, fmt.Errorf("provider %s is not available", name)
	}

	return provider, nil
}

// GetFirstAvailable returns the first available provider according to priority
func (d *Detector) GetFirstAvailable() (Provider, error) {
	for _, name := range d.priority {
		provider, ok := d.providers[name]
		if ok && provider.Available() {
			return provider, nil
		}
	}

	return nil, fmt.Errorf("no available providers found")
}

// ListProviders returns information about all registered providers
func (d *Detector) ListProviders() []ProviderInfo {
	var infos []ProviderInfo
	for name, provider := range d.providers {
		infos = append(infos, ProviderInfo{
			Name:      name,
			Command:   name, // Simplified for now
			Available: provider.Available(),
		})
	}
	return infos
}
