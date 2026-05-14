---
id: "3"
title: "forge init command"
priority: "P1"
estimated_time: "1.5h"
dependencies: ["1"]
scope: "backend"
breaking: false
type: "implementation"
mainSession: false
---

# 3: forge init command

## Description

Implement `forge init` — a one-stop project initialization command that creates `.forge/` directory, generates `CLAUDE.md` from embedded template, appends runtime entries to `.gitignore`, appends `claude`/`claude-c` recipes to `justfile`, and runs interactive config if `.forge/config.yaml` doesn't exist.

## Reference Files
- `docs/proposals/forge-info-commands/proposal.md` — Source proposal

## Affected Files

### Create
| File | Description |
|------|-------------|
| `forge-cli/internal/cmd/init.go` | `forge init` command implementation |
| `forge-cli/internal/cmd/init_test.go` | Tests for init command |
| `forge-cli/internal/embedded/claudemd.go` | `go:embed` directive for CLAUDE.md template |
| `forge-cli/internal/embedded/claudemd_test.go` | Tests for embedded template |

### Modify
| File | Changes |
|------|---------|
| `forge-cli/internal/cmd/root.go` | Register `initCmd` |

## Acceptance Criteria

- [ ] `forge init` creates `.forge/` directory (skip if exists)
- [ ] `forge init` generates `CLAUDE.md` from embedded template (skip if exists)
- [ ] `forge init` appends forge runtime entries to `.gitignore` with dedup check
- [ ] `forge init` appends `claude` / `claude-c` recipes to `justfile` with dedup check
- [ ] `forge init` runs interactive config (`forge config init`) when `.forge/config.yaml` doesn't exist
- [ ] Each step reports CREATED / APPENDED / SKIPPED status
- [ ] Execution result report matches proposal format
- [ ] CLAUDE.md template is embedded via `go:embed`
- [ ] Test coverage ≥ 80% for new code

## Hard Rules

- Never overwrite existing files — skip with SKIPPED status
- `.gitignore` append must check each line individually for duplicates
- `justfile` append must check recipe names exist before adding

## Implementation Notes

- CLAUDE.md template content: generic behavioral guidelines (Think Before Coding / Simplicity First / Surgical Changes / Goal-Driven Execution) — content defined in proposal
- Use `os.MkdirAll` for `.forge/` (idempotent)
- For `.gitignore` dedup: read existing lines into a set, only append lines not in the set
- For justfile dedup: read file, check if recipe name `claude:` already exists as a line prefix
- The init command should call `forge config init` logic directly (import function, not subprocess) when config doesn't exist
- Print a summary block at the end showing all actions taken

**gitignore entries to append:**
```
# Forge runtime
docs/features/*/tasks/process/
.forge/state.json
tests/results/.last-run.json
tests/e2e/results/.last-run.json
tests/e2e/results/*/error-context.md
```

**justfile recipes to append:**
```just
claude:
    claude --dangerously-skip-permissions

claude-c:
    claude --dangerously-skip-permissions -c
```
