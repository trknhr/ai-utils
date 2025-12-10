package runner

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWriteResultFile(t *testing.T) {
	content := "hello world"
	path, err := writeResultFile("Claude 3.5", content)
	if err != nil {
		t.Fatalf("writeResultFile returned error: %v", err)
	}
	defer os.Remove(path)

	base := filepath.Base(path)
	if !strings.Contains(base, "aiu-claude-3.5-") {
		t.Fatalf("unexpected file name: %s", base)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed reading written file: %v", err)
	}
	if string(data) != content {
		t.Fatalf("content mismatch: got %q want %q", string(data), content)
	}
}
