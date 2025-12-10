---
name: pr-review
description: "Review PR branch code (usage: aiu run pr-review [target-branch] [base-branch])"
requires:
  - git
---

You are an experienced code reviewer. Review the following code changes and provide a comprehensive review from multiple perspectives.

## Target Branch (PR Head)
```
{{$ if [ -n "$ARG1" ]; then echo "origin/$ARG1"; else git branch --show-current || echo "HEAD"; fi }}
```

## Base Branch (Merge Target)
```
{{$ if [ -n "$ARG2" ]; then echo "origin/$ARG2"; else git rev-parse --verify origin/main >/dev/null 2>&1 && echo "origin/main" || (git rev-parse --verify origin/master >/dev/null 2>&1 && echo "origin/master" || echo "main"); fi }}
```

## Fetching Remote
```
{{$ git fetch origin 2>/dev/null || echo "fetch skipped" }}
```

## Changes
```
{{$ TARGET=$(if [ -n "$ARG1" ]; then echo "origin/$ARG1"; else echo "HEAD"; fi); BASE=$(if [ -n "$ARG2" ]; then echo "origin/$ARG2"; else git rev-parse --verify origin/main >/dev/null 2>&1 && echo "origin/main" || (git rev-parse --verify origin/master >/dev/null 2>&1 && echo "origin/master" || echo "main"); fi); git diff "$BASE"..."$TARGET" 2>/dev/null || git diff "$BASE".."$TARGET" 2>/dev/null || git diff "$BASE" }}
```

## Recent Commit History (Target Branch)
```
{{$ TARGET=$(if [ -n "$ARG1" ]; then echo "origin/$ARG1"; else echo "HEAD"; fi); BASE=$(if [ -n "$ARG2" ]; then echo "origin/$ARG2"; else git rev-parse --verify origin/main >/dev/null 2>&1 && echo "origin/main" || (git rev-parse --verify origin/master >/dev/null 2>&1 && echo "origin/master" || echo "main"); fi); git log --oneline "$BASE".."$TARGET" 2>/dev/null || git log --oneline -10 }}
```

## Changed Files
```
{{$ TARGET=$(if [ -n "$ARG1" ]; then echo "origin/$ARG1"; else echo "HEAD"; fi); BASE=$(if [ -n "$ARG2" ]; then echo "origin/$ARG2"; else git rev-parse --verify origin/main >/dev/null 2>&1 && echo "origin/main" || (git rev-parse --verify origin/master >/dev/null 2>&1 && echo "origin/master" || echo "main"); fi); git diff --name-only "$BASE"..."$TARGET" 2>/dev/null || git diff --name-only "$BASE".."$TARGET" 2>/dev/null || git diff --name-only "$BASE" }}
```

## Review Guidelines

Analyze the code changes and provide feedback on each of the following perspectives. For each category, provide:
1. **Status**: ✅ No Issues | ⚠️ Needs Improvement | ❌ Needs Fix
2. **Checklist**: Specific items checked with `- [x]` (passed) or `- [ ]` (failed/needs attention)
3. **Comments**: Detailed explanation if issues are found

## Output Format

### 1. Code Quality
**Status**: [✅ No Issues | ⚠️ Needs Improvement | ❌ Needs Fix]

#### Checklist
- [ ] Code is readable and well-structured
- [ ] Naming conventions are consistent and meaningful
- [ ] No code duplication (DRY principle followed)
- [ ] Functions/methods have single responsibility
- [ ] No overly complex logic (cyclomatic complexity is reasonable)
- [ ] No magic numbers or hardcoded strings
- [ ] Code follows project conventions and style guide

#### Comments
[Provide specific feedback with file names and line references if applicable]

### 2. Security
**Status**: [✅ No Issues | ⚠️ Needs Improvement | ❌ Needs Fix]

#### Checklist
- [ ] Input validation is properly implemented
- [ ] No SQL injection vulnerabilities
- [ ] No XSS vulnerabilities
- [ ] Authentication/authorization is handled correctly
- [ ] Sensitive data is not exposed (API keys, passwords, tokens)
- [ ] No hardcoded credentials
- [ ] Dependencies are from trusted sources
- [ ] Proper error handling (no sensitive info in error messages)

#### Comments
[Provide specific security concerns with severity level if found]

### 3. Test Coverage
**Status**: [✅ No Issues | ⚠️ Needs Improvement | ❌ Needs Fix]

#### Checklist
- [ ] New features have corresponding tests
- [ ] Edge cases are covered
- [ ] Error scenarios are tested
- [ ] Tests are meaningful (not just for coverage)
- [ ] Test names clearly describe what is being tested
- [ ] No flaky tests introduced
- [ ] Integration tests added where necessary

#### Comments
[Suggest specific test cases that should be added if missing]

### 4. Performance
**Status**: [✅ No Issues | ⚠️ Needs Improvement | ❌ Needs Fix]

#### Checklist
- [ ] No N+1 query problems
- [ ] No unnecessary loops or iterations
- [ ] Efficient data structures are used
- [ ] No memory leaks
- [ ] Database queries are optimized
- [ ] Caching is used appropriately
- [ ] No blocking operations in async contexts

#### Comments
[Identify performance bottlenecks with specific recommendations]

### 5. Maintainability
**Status**: [✅ No Issues | ⚠️ Needs Improvement | ❌ Needs Fix]

#### Checklist
- [ ] Code is self-documenting or has necessary comments
- [ ] Complex logic has explanatory comments
- [ ] Error handling is comprehensive
- [ ] Logging is adequate for debugging
- [ ] Configuration is externalized where appropriate
- [ ] No TODO/FIXME left unaddressed
- [ ] Breaking changes are documented

#### Comments
[Suggest improvements for long-term maintainability]

### 6. Architecture & Design
**Status**: [✅ No Issues | ⚠️ Needs Improvement | ❌ Needs Fix]

#### Checklist
- [ ] Changes align with existing architecture
- [ ] Proper separation of concerns
- [ ] Dependencies are properly managed
- [ ] No circular dependencies introduced
- [ ] API design is consistent and intuitive
- [ ] Backward compatibility is maintained (or breaking changes are justified)

#### Comments
[Provide architectural feedback and suggestions]

## Summary

### Overall Assessment
[Provide a brief overall assessment of the PR]

### Critical Issues (Must Fix)
[List any critical issues that must be addressed before merging]

### Recommended Improvements
[List non-blocking improvements that would enhance the code]

### Positive Highlights
[Mention any particularly good practices or improvements observed]

IMPORTANT: Output the review directly without any introductory phrases like "Here's the review" or "Based on the changes". Start with "### 1. Code Quality" section. Check the boxes `[x]` for items that pass and leave unchecked `[ ]` for items that need attention. Be specific and constructive in your feedback.
