---
feature: "migrate-claude-commands"
status: tasks
mode: quick
---

# Feature (Quick): migrate-claude-commands

<!-- Status flow: tasks -> in-progress -> completed -->

Migrate Claude-related shortcuts from the project justfile to a `forge claude` subcommand that always injects `--dangerously-skip-permissions` and passes through all user args directly to the Claude CLI binary.

## Documents

| Document | Path |
|----------|------|
| Proposal | ../../proposals/migrate-claude-commands/proposal.md |

## Tasks

| ID | Title | Status | File |
|----|-------|--------|------|
| 1 | Add forge claude subcommand with arg passthrough | completed | tasks/1-add-claude-subcommand.md |
| 2 | Remove claude/claude-c/claude-w recipes from project justfile | completed | tasks/2-remove-justfile-recipes.md |
| 3 | Update forge init to stop generating claude justfile recipes | completed | tasks/3-update-forge-init.md |
| T-eval-doc | Evaluate Documentation Quality | in_progress | tasks/eval-doc.md |
