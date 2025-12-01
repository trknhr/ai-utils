package template

import (
	"embed"
	"io/fs"
	"path/filepath"
	"strings"
)

//go:embed prompts/*.md
var builtinPrompts embed.FS

// builtinTemplates caches parsed builtin templates
var builtinTemplates map[string]*Template

func init() {
	builtinTemplates = make(map[string]*Template)

	entries, err := fs.ReadDir(builtinPrompts, "prompts")
	if err != nil {
		return
	}

	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".md" {
			continue
		}

		content, err := builtinPrompts.ReadFile("prompts/" + entry.Name())
		if err != nil {
			continue
		}

		name := strings.TrimSuffix(entry.Name(), ".md")
		tmpl, err := ParseContent(string(content), "builtin://"+name)
		if err != nil {
			continue
		}

		builtinTemplates[name] = tmpl
	}
}

// GetBuiltinTemplate returns a builtin template by name
func GetBuiltinTemplate(name string) (*Template, bool) {
	tmpl, ok := builtinTemplates[name]
	return tmpl, ok
}

// ListBuiltinTemplates returns all builtin templates
func ListBuiltinTemplates() []*Template {
	templates := make([]*Template, 0, len(builtinTemplates))
	for _, tmpl := range builtinTemplates {
		templates = append(templates, tmpl)
	}
	return templates
}

// ListBuiltinTemplateNames returns all builtin template names
func ListBuiltinTemplateNames() []string {
	names := make([]string, 0, len(builtinTemplates))
	for name := range builtinTemplates {
		names = append(names, name)
	}
	return names
}
