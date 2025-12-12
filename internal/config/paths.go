package config

import (
	"os"
	"path/filepath"
)

// getWorkspaceAiuDir returns the nearest .aiu directory from the current
// working directory by walking up the directory tree. If none exists, it
// returns the .aiu path under the current working directory.
func getWorkspaceAiuDir() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	dir := cwd
	for {
		candidate := filepath.Join(dir, ".aiu")
		if st, err := os.Stat(candidate); err == nil && st.IsDir() {
			return candidate, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return filepath.Join(cwd, ".aiu"), nil
}

// GetConfigDir returns the global configuration directory path.
// Default: ~/.aiu/
func GetConfigDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	configDir := filepath.Join(home, ".aiu")
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

// GetWorkspaceConfigDir returns the workspace configuration directory path.
// Default: <workspace>/.aiu/
func GetWorkspaceConfigDir() (string, error) {
	return getWorkspaceAiuDir()
}

// GetWorkspaceConfigPath returns the workspace config.yaml path.
func GetWorkspaceConfigPath() (string, error) {
	dir, err := GetWorkspaceConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "config.yaml"), nil
}

// GetPromptsDir returns the global prompts directory path.
// Default: ~/.aiu/prompts
func GetPromptsDir() (string, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(configDir, "prompts"), nil
}

// GetWorkspacePromptsDir returns the workspace prompts directory path.
// Default: <workspace>/.aiu/prompts
func GetWorkspacePromptsDir() (string, error) {
	configDir, err := GetWorkspaceConfigDir()
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
