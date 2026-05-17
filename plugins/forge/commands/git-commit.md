---
name: git-commit
description: Create a git commit following Conventional Commits specification.
allowed-tools: Bash Read
argument-hints:
  - name: scope
    description: Optional commit scope (e.g. api, cli, core). Auto-detected from changes if omitted.
    required: false
---

Create a git commit following Conventional Commits specification.

## Atomic Commits

Group high-related changes; split unrelated changes.

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

1. Lowercase first letter
2. No trailing period
3. Imperative mood
4. Max 72 characters

## Steps

1. Run `git status` and `git diff` to inspect changes.
2. Stage appropriate files with `git add` — listing each file path explicitly.
3. Compose commit message following rules above.
4. Execute commit.

<HARD-RULE>
**NEVER use broad staging commands: `git add -A`, `git add .`, `git add --all`.**

Every `git add` must list explicit file paths. If you do not know which files changed, run `git diff --name-only` and `git diff --cached --name-only` first to identify them.

**Why:** A prior fix task (d7f8a13) that modified only 2 `.go` files resulted in a 169-file commit because the agent used `git add -A`, which staged unrelated untracked files and regenerated `.ts` residue from the working directory. Explicit file paths prevent this class of staging amplification.
</HARD-RULE>

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
