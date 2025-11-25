package template

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Template represents a parsed prompt template
type Template struct {
	// Frontmatter metadata
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	Requires    []string `yaml:"requires"`

	// Template content (after frontmatter)
	Content string

	// Original file path
	FilePath string
}

// Parse parses a template file
func Parse(filePath string) (*Template, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read template file: %w", err)
	}

	return ParseContent(string(content), filePath)
}

// ParseContent parses template content
func ParseContent(content string, filePath string) (*Template, error) {
	tmpl := &Template{
		FilePath: filePath,
	}

	// Check if content has frontmatter (starts with ---)
	if !strings.HasPrefix(content, "---\n") {
		// No frontmatter, use entire content as template
		tmpl.Content = content
		tmpl.Name = filepath.Base(filePath)
		return tmpl, nil
	}

	// Split frontmatter and content
	parts := strings.SplitN(content[4:], "\n---\n", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid frontmatter format")
	}

	// Parse frontmatter
	if err := yaml.Unmarshal([]byte(parts[0]), tmpl); err != nil {
		return nil, fmt.Errorf("failed to parse frontmatter: %w", err)
	}

	// Set content (trim leading newline)
	tmpl.Content = strings.TrimPrefix(parts[1], "\n")

	// If name is not set in frontmatter, use filename
	if tmpl.Name == "" {
		tmpl.Name = strings.TrimSuffix(filepath.Base(filePath), filepath.Ext(filePath))
	}

	return tmpl, nil
}

// FindTemplate finds a template file by name in the given directories
func FindTemplate(name string, dirs []string) (string, error) {
	// Try exact name first
	for _, dir := range dirs {
		// Try with .md extension
		path := filepath.Join(dir, name+".md")
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}

		// Try without extension (in case user provided it)
		path = filepath.Join(dir, name)
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	return "", fmt.Errorf("template %q not found in prompt directories", name)
}

// ListTemplates lists all available templates in the given directories
func ListTemplates(dirs []string) ([]*Template, error) {
	var templates []*Template
	seen := make(map[string]bool)

	for _, dir := range dirs {
		entries, err := os.ReadDir(dir)
		if err != nil {
			// Skip directories that don't exist
			if os.IsNotExist(err) {
				continue
			}
			return nil, fmt.Errorf("failed to read directory %s: %w", dir, err)
		}

		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}

			// Only process .md files
			if filepath.Ext(entry.Name()) != ".md" {
				continue
			}

			name := strings.TrimSuffix(entry.Name(), ".md")

			// Skip if already seen (first directory wins)
			if seen[name] {
				continue
			}
			seen[name] = true

			// Parse template to get metadata
			path := filepath.Join(dir, entry.Name())
			tmpl, err := Parse(path)
			if err != nil {
				// Skip templates that fail to parse
				continue
			}

			templates = append(templates, tmpl)
		}
	}

	return templates, nil
}

// Validate checks if the template's requirements are met
func (t *Template) Validate() error {
	for _, req := range t.Requires {
		// Check if required command exists
		if _, err := os.Stat(req); err != nil {
			// Not a file, check if it's a command in PATH
			// For now, we'll skip validation - can be added later
		}
	}

	return nil
}

// String returns a string representation of the template
func (t *Template) String() string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("Name: %s\n", t.Name))
	if t.Description != "" {
		buf.WriteString(fmt.Sprintf("Description: %s\n", t.Description))
	}
	if len(t.Requires) > 0 {
		buf.WriteString(fmt.Sprintf("Requires: %s\n", strings.Join(t.Requires, ", ")))
	}
	return buf.String()
}
