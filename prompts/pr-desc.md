---
name: pr-desc
description: Generate PR description from staged changes
requires:
  - git
---

Generate a clear and concise PR description based on the following git diff.

## Changes
```
{{$ git diff --cached }}
```

## Recent Commit History
```
{{$ git log --oneline -5 }}
```

## Branch Name
```
{{$ git branch --show-current }}
```

---
Output format:
- Title: One-line summary of changes
- Overview: Purpose and context of changes
- Details: List of main changes

IMPORTANT: Output the PR description directly without any introductory phrases. Do not include phrases like "Based on...", "Here's...", or markdown separators (---) at the beginning. Start directly with the ## Title section.
