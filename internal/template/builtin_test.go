package template

import (
	"testing"
)

func TestGetBuiltinTemplate(t *testing.T) {
	// Test that builtin templates are loaded
	templates := ListBuiltinTemplates()
	if len(templates) == 0 {
		t.Error("expected at least one builtin template")
	}

	// Check for expected templates
	expectedTemplates := []string{"pr-desc", "pr-review", "commit-msg"}
	for _, name := range expectedTemplates {
		tmpl, ok := GetBuiltinTemplate(name)
		if !ok {
			t.Errorf("expected builtin template '%s' to exist", name)
			continue
		}

		if tmpl.Name != name {
			t.Errorf("expected template name '%s', got '%s'", name, tmpl.Name)
		}

		if tmpl.Content == "" {
			t.Errorf("template '%s' has empty content", name)
		}
	}
}

func TestGetBuiltinTemplate_NotFound(t *testing.T) {
	_, ok := GetBuiltinTemplate("nonexistent-template")
	if ok {
		t.Error("expected false for nonexistent template")
	}
}

func TestListBuiltinTemplates(t *testing.T) {
	templates := ListBuiltinTemplates()

	// Should have at least the templates we created
	if len(templates) < 3 {
		t.Errorf("expected at least 3 builtin templates, got %d", len(templates))
	}

	// Check that all templates have required fields
	for _, tmpl := range templates {
		if tmpl.Name == "" {
			t.Error("builtin template has empty name")
		}
		if tmpl.Content == "" {
			t.Errorf("builtin template '%s' has empty content", tmpl.Name)
		}
	}
}

func TestBuiltinTemplates_HaveGitRequirement(t *testing.T) {
	gitTemplates := []string{"pr-desc", "pr-review", "commit-msg"}

	for _, name := range gitTemplates {
		tmpl, ok := GetBuiltinTemplate(name)
		if !ok {
			t.Errorf("template '%s' not found", name)
			continue
		}

		hasGit := false
		for _, req := range tmpl.Requires {
			if req == "git" {
				hasGit = true
				break
			}
		}

		if !hasGit {
			t.Errorf("template '%s' should require 'git'", name)
		}
	}
}
