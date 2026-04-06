---
name: git-commit
description: Use when creating git commits. Ensures commit messages follow Conventional Commits format.
---

# Git Commit

## Overview

Follow Conventional Commits specification with project-defined type and scope rules.

## Atomic Commits

**Core principle**: group high-related changes; split unrelated changes.

| Practice           | Description                                     |
| ------------------ | ----------------------------------------------- |
| **Group together** | feature + its tests + its docs in one commit    |
| **Split apart**    | unrelated features, fixes, independent refactor |

## Format

```
<type>(<scope>): <subject>

[optional body]

[optional footer(s)]
```

## Allowed Types

| Type       | When to Use                |
| ---------- | -------------------------- |
| `feat`     | New feature                |
| `fix`      | Bug fix                    |
| `docs`     | Documentation only         |
| `test`     | Adding/modifying tests     |
| `refactor` | Code refactoring           |
| `chore`    | Maintenance, tooling, deps |

## Scope Examples

| Scope  | Module              |
| ------ | ------------------- |
| `api`  | API layer           |
| `app`  | Application layer   |
| `cli`  | CLI entry point     |
| `core` | Core logic          |
| `docs` | Documentation       |
| `test` | Test infrastructure |

## Subject Rules

1. **Lowercase first letter** - `add` not `Add`
2. **No trailing period**
3. **Imperative mood** - `add` not `added`
4. **Max 72 characters**

## Examples

```bash
# Good
feat(api): add streaming support
fix(parser): handle empty input
docs(readme): update install steps

# Bad
Update(api): add streaming    # Wrong type
feat(api): Added support.     # Past tense, period
```

## Task Completion Template

```bash
git add <changed-files>
git commit -m "$(cat <<'EOF'
<type>(<scope>): <subject>

Task: <TASK_ID>

Co-Authored-By: Agent
EOF
)"
```

## Quick Checklist

- [ ] Type is one of: feat / fix / docs / test / refactor / chore
- [ ] Scope matches affected module
- [ ] Subject starts with lowercase
- [ ] Subject has no trailing period
- [ ] Subject is imperative mood
