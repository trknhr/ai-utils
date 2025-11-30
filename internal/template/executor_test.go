package template

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExecutor_Execute_NoCommands(t *testing.T) {
	tmpl := &Template{
		Name:    "test",
		Content: "This is plain content without any commands.",
	}

	executor := NewExecutor("", false, nil)
	result, err := executor.Execute(tmpl)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if result != tmpl.Content {
		t.Errorf("expected '%s', got '%s'", tmpl.Content, result)
	}
}

func TestExecutor_Execute_WithCommand(t *testing.T) {
	tmpl := &Template{
		Name:    "test",
		Content: "Output: {{$ echo hello }}",
	}

	executor := NewExecutor("", false, nil)
	result, err := executor.Execute(tmpl)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	expected := "Output: hello"
	if result != expected {
		t.Errorf("expected '%s', got '%s'", expected, result)
	}
}

func TestExecutor_Execute_MultipleCommands(t *testing.T) {
	tmpl := &Template{
		Name:    "test",
		Content: "First: {{$ echo one }}\nSecond: {{$ echo two }}",
	}

	executor := NewExecutor("", false, nil)
	result, err := executor.Execute(tmpl)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	expected := "First: one\nSecond: two"
	if result != expected {
		t.Errorf("expected '%s', got '%s'", expected, result)
	}
}

func TestExecutor_Execute_WithWorkingDir(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "executor-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a test file
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("content"), 0644); err != nil {
		t.Fatalf("failed to write file: %v", err)
	}

	tmpl := &Template{
		Name:    "test",
		Content: "Files: {{$ ls -1 }}",
	}

	executor := NewExecutor(tmpDir, false, nil)
	result, err := executor.Execute(tmpl)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if result != "Files: test.txt" {
		t.Errorf("expected 'Files: test.txt', got '%s'", result)
	}
}

func TestExecutor_Execute_WithArgs(t *testing.T) {
	tmpl := &Template{
		Name:    "test",
		Content: "Arg1: {{$ echo $ARG1 }}\nArg2: {{$ echo $ARG2 }}\nAll: {{$ echo $ARGS }}",
	}

	executor := NewExecutor("", false, []string{"first", "second"})
	result, err := executor.Execute(tmpl)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	expected := "Arg1: first\nArg2: second\nAll: first second"
	if result != expected {
		t.Errorf("expected '%s', got '%s'", expected, result)
	}
}

func TestExecutor_Execute_FailingCommand(t *testing.T) {
	tmpl := &Template{
		Name:    "test",
		Content: "Output: {{$ exit 1 }}",
	}

	executor := NewExecutor("", false, nil)
	_, err := executor.Execute(tmpl)
	if err == nil {
		t.Error("expected error for failing command, got nil")
	}
}

func TestExecutor_BuildEnv(t *testing.T) {
	executor := NewExecutor("", false, []string{"arg1", "arg2", "arg3"})
	env := executor.buildEnv()

	// Check that ARG1, ARG2, ARG3 are set
	envMap := make(map[string]string)
	for _, e := range env {
		for i := 0; i < len(e); i++ {
			if e[i] == '=' {
				envMap[e[:i]] = e[i+1:]
				break
			}
		}
	}

	if envMap["ARG1"] != "arg1" {
		t.Errorf("expected ARG1='arg1', got '%s'", envMap["ARG1"])
	}
	if envMap["ARG2"] != "arg2" {
		t.Errorf("expected ARG2='arg2', got '%s'", envMap["ARG2"])
	}
	if envMap["ARG3"] != "arg3" {
		t.Errorf("expected ARG3='arg3', got '%s'", envMap["ARG3"])
	}
	if envMap["ARGS"] != "arg1 arg2 arg3" {
		t.Errorf("expected ARGS='arg1 arg2 arg3', got '%s'", envMap["ARGS"])
	}
}

func TestExecutor_DryRun(t *testing.T) {
	tmpl := &Template{
		Name:    "test",
		Content: "Output: {{$ echo hello }}",
	}

	executor := NewExecutor("", false, nil)
	result := executor.DryRun(tmpl)

	// DryRun should contain the original content
	if result == "" {
		t.Error("DryRun returned empty string")
	}

	// Should contain the command
	if !containsSubstring(result, "echo hello") {
		t.Error("DryRun missing command")
	}
}

func TestCommandPattern(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{"{{$ echo hello }}", []string{"echo hello"}},
		{"{{$echo hello}}", nil}, // No space after $
		{"{{$ git diff }}", []string{"git diff"}},
		{"{{$ ls -la | grep test }}", []string{"ls -la | grep test"}},
		{"prefix {{$ cmd1 }} middle {{$ cmd2 }} suffix", []string{"cmd1", "cmd2"}},
	}

	for _, tt := range tests {
		matches := commandPattern.FindAllStringSubmatch(tt.input, -1)
		var cmds []string
		for _, m := range matches {
			if len(m) > 1 {
				cmds = append(cmds, m[1])
			}
		}

		if len(cmds) != len(tt.expected) {
			t.Errorf("input '%s': expected %d matches, got %d", tt.input, len(tt.expected), len(cmds))
			continue
		}

		for i, cmd := range cmds {
			if cmd != tt.expected[i] {
				t.Errorf("input '%s': expected '%s', got '%s'", tt.input, tt.expected[i], cmd)
			}
		}
	}
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
