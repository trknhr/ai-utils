package runner

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func writeResultFile(providerName, content string) (string, error) {
	safeName := strings.ReplaceAll(strings.ToLower(providerName), " ", "-")
	pattern := fmt.Sprintf("aiu-%s-*.md", safeName)
	file, err := os.CreateTemp("", pattern)
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	defer file.Close()

	if _, err := file.WriteString(content); err != nil {
		return "", fmt.Errorf("failed to write temp file: %w", err)
	}

	return file.Name(), nil
}

func openInVSCode(filePaths []string) {
	if len(filePaths) == 0 {
		return
	}

	if _, err := exec.LookPath("code"); err != nil {
		return
	}

	args := append([]string{}, filePaths...)
	if len(filePaths) == 2 {
		args = append([]string{"--diff"}, args...)
	}

	// Normalize paths to avoid surprises when VS Code opens them.
	for i, path := range args {
		if path == "--diff" {
			continue
		}

		absPath, err := filepath.Abs(path)
		if err == nil {
			args[i] = absPath
		}
	}

	_ = exec.Command("code", args...).Start()
}
