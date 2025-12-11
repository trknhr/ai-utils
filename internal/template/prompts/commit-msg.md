---
name: commit-msg
description: "Generate commit message from staged changes"
requires:
  - git
---

You are an expert at writing clear and concise git commit messages following best practices.

## Staged Files
```
{{$ git diff --cached --name-only }}
```

## Staged Changes
```
{{$ git diff --cached }}
```

## Recent Commit History (for context)
```
{{$ git log --oneline -5 }}
```


## Guidelines

Generate a commit message following the Conventional Commits format:

```
<type>(<scope>): <subject>

<body>
```

### Type (required)
- `feat`: A new feature
- `fix`: A bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, semicolons, etc.)
- `refactor`: Code refactoring (no feature or bug fix)
- `perf`: Performance improvements
- `test`: Adding or updating tests
- `chore`: Maintenance tasks (build, dependencies, etc.)
- `ci`: CI/CD changes

### Scope (optional)
The area of the codebase affected (e.g., `auth`, `api`, `ui`, `config`)

### Subject (required)
- Use imperative mood ("add" not "added" or "adds")
- Don't capitalize the first letter
- No period at the end
- Max 50 characters

### Body (optional, for complex changes)
- Explain what and why, not how
- Wrap at 72 characters
- Use bullet points if needed

## Output Format

Provide the commit message directly. If the changes are simple, just provide the subject line. For complex changes, include the body.

Example outputs:

Simple:
```
feat(auth): add JWT token refresh endpoint
```

With body:
```
feat(auth): add JWT token refresh endpoint

- Add /auth/refresh endpoint for token renewal
- Implement sliding window expiration
- Add refresh token rotation for security
```

IMPORTANT: Output only the commit message without any explanation or markdown code blocks. Just the raw commit message text that can be directly used with `git commit -m`.
