<p align="center">
  <img src="img/logo.svg" alt="aiu" width="200">
</p>

<h1 align="center">AI Utils (aiu)</h1>

A developer-friendly CLI tool that manages reusable prompt templates and executes them using locally installed AI CLIs (Claude Code, Gemini CLI, Codex CLI). It automatically detects available AI providers, expands shell commands within templates, and eliminates the need for API keys.

## Features

- **Template Management**: Store reusable prompt templates with Markdown + YAML frontmatter
- **Dynamic Command Expansion**: Embed shell commands in templates using `{{$ command }}` syntax
- **Template Arguments**: Pass arguments to templates (e.g., `aiu pr-review feature-branch main`)
- **Built-in Templates**: PR description, PR review, and commit message templates included in binary
- **Multi-Provider Support**: Automatically detects and uses Claude, Gemini, or Codex CLIs
- **No API Keys Required**: Leverages your existing CLI installations

## Installation (macOS / Linux)

Install the latest release via:

```bash
curl -sSfL https://raw.githubusercontent.com/trknhr/ai-utils/main/install.sh | sh
```

This will:
- Detect your OS and CPU architecture
- Download the correct binary
- Install it to `/usr/local/bin/aiu`

## Quick Start

### 1. Install a supported AI CLI

You need at least one of these installed:

- **Claude Code**: Follow [ClaudeCode installation](https://github.com/anthropics/claude-code)
- **Gemini CLI**: Follow [Gemini CLI installation](https://github.com/google-gemini/gemini-cli)
- **Codex CLI**: Follow [Codex installation](https://github.com/openai/codex)

### 2. Run a prompt

```bash
# Generate PR description
aiu pr-desc

# Review a PR branch (compare feature/xxx with origin/main)
aiu pr-review feature/xxx

# Review with custom base branch
aiu pr-review feature/xxx develop

# Generate commit message from staged changes
git add .
aiu commit-msg
```

### 3. Enable auto commit messages (optional)

```bash
# Enable AI-powered commit messages globally
aiu enable-auto-commit

# Now just use git commit normally - AI generates the message!
git add .
git commit

# Skip AI for a specific commit
AI_IGNORE=1 git commit

# Disable auto commit
aiu disable-auto-commit
```

## How It Works

```
┌─────────────┐    1. Read Template    ┌──────────────────┐
│ Prompt File │◄──────────────────────│   aiu CLI        │
│  (*.md)     │                        │                  │
└─────────────┘                        │  - Parse YAML    │
                                       │  - Expand {{$}}  │
┌─────────────┐    2. Execute Commands │  - Detect AI CLI │
│ Shell       │◄──────────────────────│                  │
│ (git, etc)  │                        └────────┬─────────┘
└─────────────┘                                 │
                                                │ 3. Send Prompt
                 ┌──────────────────────────────▼─────┐
                 │  AI Provider (Claude/Gemini/Codex) │
                 └────────────────────────────────────┘
```

1. **Template Loading**: Reads Markdown files with YAML frontmatter
2. **Command Expansion**: Executes `{{$ command }}` placeholders and injects output
3. **Provider Detection**: Automatically selects available AI CLI (Claude → Gemini → Codex)
4. **Execution**: Sends expanded prompt to the provider and displays response

## Template System

### Template Format

Templates are Markdown files with optional YAML frontmatter:

```markdown
---
name: pr-desc
description: Generate PR description from staged changes
requires:
  - git
---

Generate a clear and concise PR description based on the following git diff.

## Changes
\```
{{$ git diff --cached }}
\```

## Recent Commit History
\```
{{$ git log --oneline -5 }}
\```

Output format:
- Title: One-line summary
- Overview: Purpose and context
- Details: List of main changes
```

### Command Syntax

| Syntax | Description | Example |
|--------|-------------|---------|
| `{{$ command }}` | Execute shell command and insert output | `{{$ git diff }}` |
| `$ARG1`, `$ARG2`... | Template arguments (passed from CLI) | `{{$ echo $ARG1 }}` |
| `$ARGS` | All arguments joined by space | `{{$ echo $ARGS }}` |

### Template Location

Built-in templates are embedded in the binary. Custom templates can be stored in:
- `<workspace>/.aiu/prompts/` (overrides global; searched first)
- `~/.aiu/prompts/` (global)
- Additional directories via `config.yaml` (`prompt_dirs`)

## Configuration

Configuration files (merged in this order; workspace overrides global):
- `~/.aiu/config.yaml`
- `<workspace>/.aiu/config.yaml`

```yaml
# Prompt directories (searched in order)
prompt_dirs:
  - <workspace>/.aiu/prompts
  - ~/.aiu/prompts
  - ~/.local/share/ai-utils/prompts

# Provider priority (tries in order)
provider_priority:
  - claude
  - gemini
  - codex

# Provider-specific settings
providers:
  claude:
    command: claude
    args: ["-p", "--output-format", "text"]
    timeout: 120s
  gemini:
    command: gemini
    args: []
    timeout: 120s
  codex:
    command: codex
    args: []
    timeout: 120s

# Defaults
defaults:
  timeout: 60s
  verbose: false
```

## Usage Examples

### Generate PR Description

```bash
# From current branch vs origin/main
aiu pr-desc
```

### PR Code Review

```bash
# Review current branch vs origin/main
aiu pr-review

# Review specific branch vs origin/main
aiu pr-review feature/new-api

# Review with custom base branch
aiu pr-review feature/new-api develop
```

### Generate Commit Message

```bash
# Stage changes first
git add .

# Generate commit message from staged changes
aiu commit-msg

# Or use auto-commit hook
aiu enable-auto-commit
git commit  # AI generates message automatically
```

### Dry Run (Preview)

```bash
# Show expanded prompt without sending to AI
aiu pr-review --dry-run
```

### Specify Provider

```bash
# Use specific provider
aiu pr-desc --provider gemini
```

## Development

### Build from Source

```bash
# Clone repository
git clone https://github.com/trknhr/ai-utils.git
cd ai-utils

# Build
make build

# Install locally
make install

# Run tests
make test
```

### Project Structure

```
ai-utils/
├── cmd/aiu/              # CLI entry point
├── internal/
│   ├── cli/              # Command implementations
│   ├── config/           # Configuration management
│   ├── provider/         # AI provider implementations
│   └── template/         # Template parser & executor
│       └── prompts/      # Built-in templates (embedded)
├── build/                # Build output (gitignored)
├── go.mod
├── Makefile
└── README.md
```

### Requirements

- Go 1.21+
- At least one supported AI CLI installed

## Built-in Templates

| Template | Description | Usage |
|----------|-------------|-------|
| `pr-desc` | Generate PR description from git diff | `aiu pr-desc` |
| `pr-review` | Comprehensive code review checklist | `aiu pr-review [target] [base]` |
| `commit-msg` | Generate commit message from staged changes | `aiu commit-msg` |

### pr-review Output

The PR review template generates a detailed review covering:
- Code Quality (readability, DRY, complexity)
- Security (input validation, credentials, vulnerabilities)
- Test Coverage (unit tests, edge cases)
- Performance (N+1, memory, optimization)
- Maintainability (docs, error handling, logging)
- Architecture & Design (separation of concerns, dependencies)

## Roadmap

- [x] **Phase 1**: Basic template system + Claude provider
- [x] **Phase 2**: Built-in templates, template arguments, auto-commit hook
- [ ] **Phase 3**: Gemini & Codex providers, `aiu list` command
- [ ] **Phase 4**: Interactive mode, colored output, better error messages
- [ ] **Phase 5**: Plugin system, custom providers

## Feedback & Contributions

PRs and issues welcome → [github.com/trknhr/ai-utils](https://github.com/trknhr/ai-utils)

## License

Apache 2.0 License — © 2025 Teruo Kunihiro
