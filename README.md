# 🧠 AI Utils (`aiu`)

A developer-friendly CLI tool that manages reusable prompt templates and executes them using locally installed AI CLIs (Claude Code, Gemini CLI, Codex CLI). It automatically detects available AI providers, expands shell commands within templates, and eliminates the need for API keys.

## ✨ Features

- **Template Management**: Store reusable prompt templates with Markdown + YAML frontmatter
- **Dynamic Command Expansion**: Embed shell commands in templates using `{{$ command }}` syntax
- **Multi-Provider Support**: Automatically detects and uses Claude, Gemini, or Codex CLIs
- **Cost Optimization**: Uses subscription-based CLIs instead of pay-per-use APIs
- **No API Keys Required**: Leverages your existing CLI installations
- **Git Integration**: Built-in support for `git diff`, `git log`, and more
- **Vendor Agnostic**: Provider fallback system for high availability

## 📦 Installation (macOS / Linux)

Install the latest release via:

```bash
curl -sSfL https://raw.githubusercontent.com/trknhr/ai-utils/main/install.sh | sh
```

This will:
- Detect your OS and CPU architecture
- Download the correct binary
- Install it to `/usr/local/bin/aiu`

## 🚀 Quick Start

### 1. Install a supported AI CLI

You need at least one of these installed:

- **Claude Code**: `npm install -g @anthropic-ai/claude-code`
- **Gemini CLI**: Follow [Gemini CLI installation](https://github.com/google-gemini/gemini-cli)
- **Codex CLI**: Follow [Codex installation](https://github.com/openai/codex)

### 2. Initialize configuration

```bash
# Create config directory and sample prompts
make dev-setup

# Or manually
mkdir -p ~/.config/ai-utils/prompts
```

### 3. Run a prompt

```bash
# Generate PR description from staged changes
aiu run pr-desc

# Or use the shorthand
aiu pr-desc
```

## 🧠 How It Works

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

## 📝 Template System

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
| Future: `{{.Var}}` | Variable substitution | `{{.FilePath}}` |
| Future: `{{env "VAR"}}` | Environment variable | `{{env "USER"}}` |

### Template Location

Templates are stored in:
- `~/.config/ai-utils/prompts/` (default)
- Additional directories via `config.yaml`

List available templates:

```bash
aiu list  # Coming in Phase 2
```

## ⚙️ Configuration

Configuration file: `~/.config/ai-utils/config.yaml`

```yaml
# Prompt directories (searched in order)
prompt_dirs:
  - ~/.config/ai-utils/prompts
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

## 💡 Usage Examples

### Generate PR Description

```bash
# From staged changes
git add .
aiu run pr-desc
```

### Code Review (Coming Soon)

```bash
# Review staged changes
aiu run review
```

### Commit Message (Coming Soon)

```bash
# Generate commit message
aiu run commit-msg
```

### Specify Provider

```bash
# Use specific provider
aiu run pr-desc --provider gemini

# Verbose output
aiu run pr-desc -v

# Dry run (show expanded prompt without executing)
aiu run pr-desc --dry-run
```

## 🛠 Development

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
├── cmd/aiu/           # CLI entry point
├── internal/
│   ├── cli/           # Command implementations
│   ├── config/        # Configuration management
│   ├── provider/      # AI provider implementations
│   └── template/      # Template parser & executor
├── prompts/           # Sample templates
├── go.mod
├── Makefile
└── README.md
```

### Requirements

- Go 1.24+
- At least one supported AI CLI installed

## 🧩 Built-in Templates

- **pr-desc**: Generate PR description from `git diff --cached`
- More templates coming in future releases

## 🗺️ Roadmap

- [x] **Phase 1 (MVP)**: Basic template system + Claude provider
- [ ] **Phase 2**: Gemini & Codex providers, `aiu list` command
- [ ] **Phase 3**: Variable expansion, file includes, environment variables
- [ ] **Phase 4**: Interactive mode, colored output, better error messages
- [ ] **Phase 5**: Plugin system, custom providers

## 📬 Feedback & Contributions

PRs and issues welcome → [github.com/trknhr/ai-utils](https://github.com/trknhr/ai-utils)

## 📄 License

Apache 2.0 License — © 2025 Teruo Kunihiro
