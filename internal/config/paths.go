package config

import (
	"os"
	"path/filepath"
)

// GetConfigDir returns the configuration directory path
// Default: ~/.config/ai-utils/
func GetConfigDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	configDir := filepath.Join(home, ".config", "ai-utils")
	return configDir, nil
}

// GetConfigPath returns the full path to config.yaml
func GetConfigPath() (string, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(configDir, "config.yaml"), nil
}

// GetPromptsDir returns the prompts directory path
// Default: ~/.config/ai-utils/prompts
func GetPromptsDir() (string, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(configDir, "prompts"), nil
}

// EnsureConfigDir creates the config directory if it doesn't exist
func EnsureConfigDir() error {
	configDir, err := GetConfigDir()
	if err != nil {
		return err
	}

	return os.MkdirAll(configDir, 0755)
}

// EnsurePromptsDir creates the prompts directory if it doesn't exist
func EnsurePromptsDir() error {
	promptsDir, err := GetPromptsDir()
	if err != nil {
		return err
	}

	return os.MkdirAll(promptsDir, 0755)
}
