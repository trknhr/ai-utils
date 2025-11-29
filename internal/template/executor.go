package template

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

// commandPattern matches {{$ command }} syntax
var commandPattern = regexp.MustCompile(`\{\{\$\s+([^}]+?)\s*\}\}`)

// Executor executes template expansion
type Executor struct {
	workingDir string
	verbose    bool
	args       []string
}

// NewExecutor creates a new template executor
func NewExecutor(workingDir string, verbose bool, args []string) *Executor {
	return &Executor{
		workingDir: workingDir,
		verbose:    verbose,
		args:       args,
	}
}

// Execute expands a template by executing shell commands
func (e *Executor) Execute(tmpl *Template) (string, error) {
	content := tmpl.Content

	// Find all command placeholders
	matches := commandPattern.FindAllStringSubmatch(content, -1)
	if len(matches) == 0 {
		// No commands to execute, return as-is
		return content, nil
	}

	// Execute each command and replace
	for _, match := range matches {
		placeholder := match[0] // Full match: {{$ command }}
		command := match[1]     // Captured command

		if e.verbose {
			fmt.Printf("Executing: %s\n", command)
		}

		output, err := e.executeCommand(command)
		if err != nil {
			return "", fmt.Errorf("failed to execute command %q: %w", command, err)
		}

		// Replace placeholder with command output
		content = strings.ReplaceAll(content, placeholder, output)
	}

	return content, nil
}

// executeCommand runs a shell command and returns its output
func (e *Executor) executeCommand(command string) (string, error) {
	// Use sh -c to execute the command (supports pipes, redirects, etc.)
	cmd := exec.Command("sh", "-c", command)

	if e.workingDir != "" {
		cmd.Dir = e.workingDir
	}

	// Set template arguments as environment variables
	cmd.Env = e.buildEnv()

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("command failed: %w\nOutput: %s", err, string(output))
	}

	return strings.TrimSpace(string(output)), nil
}

// buildEnv creates environment variables for template arguments
func (e *Executor) buildEnv() []string {
	// Start with current environment
	env := os.Environ()

	// Add positional arguments as $1, $2, $3...
	for i, arg := range e.args {
		env = append(env, fmt.Sprintf("ARG%d=%s", i+1, arg))
	}

	// Add all arguments as $ARGS
	env = append(env, fmt.Sprintf("ARGS=%s", strings.Join(e.args, " ")))

	return env
}

// DryRun returns the template content without executing commands
// It shows what commands would be executed
func (e *Executor) DryRun(tmpl *Template) string {
	content := tmpl.Content

	matches := commandPattern.FindAllStringSubmatch(content, -1)
	if len(matches) == 0 {
		return content
	}

	var commands []string
	for _, match := range matches {
		commands = append(commands, match[1])
	}

	var result strings.Builder
	result.WriteString("=== Template Content ===\n")
	result.WriteString(content)
	result.WriteString("\n\n=== Commands to Execute ===\n")
	for i, cmd := range commands {
		result.WriteString(fmt.Sprintf("%d. %s\n", i+1, cmd))
	}

	return result.String()
}
