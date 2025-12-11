package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config represents the application configuration
type Config struct {
	PromptDirs       []string                  `mapstructure:"prompt_dirs"`
	ProviderPriority []string                  `mapstructure:"provider_priority"`
	Providers        map[string]ProviderConfig `mapstructure:"providers"`
	Defaults         DefaultConfig             `mapstructure:"defaults"`
	OutputLang       string                    `mapstructure:"output_lang"`
}

// ProviderConfig represents a provider's configuration
type ProviderConfig struct {
	Command string        `mapstructure:"command"`
	Args    []string      `mapstructure:"args"`
	Timeout time.Duration `mapstructure:"timeout"`
}

// DefaultConfig represents default settings
type DefaultConfig struct {
	Timeout time.Duration `mapstructure:"timeout"`
	Verbose bool          `mapstructure:"verbose"`
}

// Load loads the configuration from file
func Load() (*Config, error) {
	v := viper.New()

	// Set config file path
	configPath, err := GetConfigPath()
	if err != nil {
		return nil, fmt.Errorf("failed to get config path: %w", err)
	}

	v.SetConfigFile(configPath)
	v.SetConfigType("yaml")

	// Set defaults
	setDefaults(v)

	// Read config file if it exists
	if err := v.ReadInConfig(); err != nil {
		// If config file doesn't exist, use defaults
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			// Also check for os.ErrNotExist (file not found)
			if !os.IsNotExist(err) {
				return nil, fmt.Errorf("failed to read config: %w", err)
			}
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// If prompt_dirs is empty, use default
	if len(cfg.PromptDirs) == 0 {
		promptsDir, err := GetPromptsDir()
		if err != nil {
			return nil, err
		}
		cfg.PromptDirs = []string{promptsDir}
	} else {
		// Expand ~ and environment variables in prompt_dirs
		for i, dir := range cfg.PromptDirs {
			cfg.PromptDirs[i] = expandPath(dir)
		}
	}

	// If output_lang is empty, use default
	if cfg.OutputLang == "" {
		cfg.OutputLang = v.GetString("output_lang")
	}

	return &cfg, nil
}

// expandPath expands ~ and environment variables in a path
func expandPath(path string) string {
	// Expand environment variables
	path = os.ExpandEnv(path)

	// Expand ~ to home directory
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err == nil {
			path = filepath.Join(home, path[2:])
		}
	}

	return path
}

// setDefaults sets default configuration values
func setDefaults(v *viper.Viper) {
	promptsDir, _ := GetPromptsDir()

	v.SetDefault("prompt_dirs", []string{promptsDir})
	v.SetDefault("provider_priority", []string{"claude", "gemini", "codex"})
	v.SetDefault("output_lang", "en")

	// Claude defaults
	v.SetDefault("providers.claude.command", "claude")
	v.SetDefault("providers.claude.args", []string{
		"-p",
		"--output-format", "text",
		"--settings", `{"attribution":{"commit":"","pr":""},"includeCoAuthoredBy":false,"gitAttribution":false}`,
	})
	v.SetDefault("providers.claude.timeout", 120*time.Second)

	// Gemini defaults
	v.SetDefault("providers.gemini.command", "gemini")
	v.SetDefault("providers.gemini.args", []string{})
	v.SetDefault("providers.gemini.timeout", 120*time.Second)

	// Codex defaults
	v.SetDefault("providers.codex.command", "codex")
	v.SetDefault("providers.codex.args", []string{})
	v.SetDefault("providers.codex.timeout", 120*time.Second)

	// Defaults
	v.SetDefault("defaults.timeout", 60*time.Second)
	v.SetDefault("defaults.verbose", false)
}

// CreateDefaultConfig creates a default config file
func CreateDefaultConfig() error {
	if err := EnsureConfigDir(); err != nil {
		return err
	}

	configPath, err := GetConfigPath()
	if err != nil {
		return err
	}

	v := viper.New()
	setDefaults(v)

	return v.WriteConfigAs(configPath)
}
