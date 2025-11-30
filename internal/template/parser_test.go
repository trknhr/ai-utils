package template

import (
"os"
"path/filepath"
"strings"
"testing"
)

func TestParseContent_WithFrontmatter(t *testing.T) {
	content := `---
name: test-template
description: "A test template"
requires:
  - git
---

This is the template content.
`
	tmpl, err := ParseContent(content, "test.md")
	if err != nil {
		t.Fatalf("ParseContent failed: %v", err)
	}

	if tmpl.Name != "test-template" {
		t.Errorf("expected Name 'test-template', got '%s'", tmpl.Name)
	}

	if tmpl.Description != "A test template" {
		t.Errorf("expected Description 'A test template', got '%s'", tmpl.Description)
	}

	if len(tmpl.Requires) != 1 || tmpl.Requires[0] != "git" {
		t.Errorf("expected Requires ['git'], got %v", tmpl.Requires)
	}

	expectedContent := "This is the template content.\n"
	if tmpl.Content != expectedContent {
		t.Errorf("expected Content '%s', got '%s'", expectedContent, tmpl.Content)
	}
}

func TestParseContent_WithoutFrontmatter(t *testing.T) {
	content := "Just plain content without frontmatter."

	tmpl, err := ParseContent(content, "plain.md")
	if err != nil {
		t.Fatalf("ParseContent failed: %v", err)
	}

	if tmpl.Name != "plain.md" {
		t.Errorf("expected Name 'plain.md', got '%s'", tmpl.Name)
	}

	if tmpl.Content != content {
		t.Errorf("expected Content '%s', got '%s'", content, tmpl.Content)
	}
}

func TestParseContent_InvalidFrontmatter(t *testing.T) {
	content := `---
name: test
invalid yaml: [
---

Content here.
`
	_, err := ParseContent(content, "invalid.md")
	if err == nil {
		t.Error("expected error for invalid YAML, got nil")
	}
}

func TestParseContent_MissingEndDelimiter(t *testing.T) {
	content := `---
name: test
description: missing end delimiter

This should fail.
`
	_, err := ParseContent(content, "missing.md")
	if err == nil {
		t.Error("expected error for missing end delimiter, got nil")
	}
}

func TestParseContent_EmptyName(t *testing.T) {
	content := `---
description: "No name specified"
---

Content without name.
`
	tmpl, err := ParseContent(content, "noname.md")
	if err != nil {
		t.Fatalf("ParseContent failed: %v", err)
	}

	if tmpl.Name != "noname" {
		t.Errorf("expected Name 'noname', got '%s'", tmpl.Name)
	}
}

func TestFindTemplate(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "template-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	templatePath := filepath.Join(tmpDir, "test-template.md")
	if err := os.WriteFile(templatePath, []byte("test content"), 0644); err != nil {
		t.Fatalf("failed to write template file: %v", err)
	}

	found, err := FindTemplate("test-template", []string{tmpDir})
	if err != nil {
		t.Fatalf("FindTemplate failed: %v", err)
	}
	if found != templatePath {
		t.Errorf("expected path '%s', got '%s'", templatePath, found)
	}
}

func TestFindTemplate_NotFound(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "template-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	_, err = FindTemplate("nonexistent", []string{tmpDir})
	if err == nil {
		t.Error("expected error for nonexistent template, got nil")
	}
}

func TestListTemplates(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "template-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	templates := []string{"template1.md", "template2.md"}
	for _, name := range templates {
		content := "---\nname: " + strings.TrimSuffix(name, ".md") + "\n---\n\nContent"
		if err := os.WriteFile(filepath.Join(tmpDir, name), []byte(content), 0644); err != nil {
			t.Fatalf("failed to write template: %v", err)
		}
	}

	if err := os.WriteFile(filepath.Join(tmpDir, "readme.txt"), []byte("ignored"), 0644); err != nil {
		t.Fatalf("failed to write file: %v", err)
	}

	result, err := ListTemplates([]string{tmpDir})
	if err != nil {
		t.Fatalf("ListTemplates failed: %v", err)
	}

	if len(result) != 2 {
		t.Errorf("expected 2 templates, got %d", len(result))
	}
}

func TestListTemplates_NonexistentDir(t *testing.T) {
	result, err := ListTemplates([]string{"/nonexistent/path"})
	if err != nil {
		t.Fatalf("ListTemplates failed: %v", err)
	}

	if len(result) != 0 {
		t.Errorf("expected 0 templates, got %d", len(result))
	}
}

func TestTemplate_String(t *testing.T) {
	tmpl := &Template{
		Name:        "test",
		Description: "Test description",
		Requires:    []string{"git", "go"},
	}

	str := tmpl.String()
	if str == "" {
		t.Error("String() returned empty string")
	}

	if !strings.Contains(str, "test") {
		t.Error("String() missing name")
	}
	if !strings.Contains(str, "Test description") {
		t.Error("String() missing description")
	}
}

func TestTemplate_Validate(t *testing.T) {
	tmpl := &Template{
		Name:     "test",
		Requires: []string{"git"},
	}

	if err := tmpl.Validate(); err != nil {
		t.Errorf("Validate failed: %v", err)
	}
}
