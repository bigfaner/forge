---
id: "1"
title: "Add auto.validation and worktree prompts to forge init (huh TUI)"
priority: "P1"
estimated_time: "1h"
dependencies: []
scope: "backend"
breaking: false
type: "coding.feature"
mainSession: false
---

# 1: Add auto.validation and worktree prompts to forge init (huh TUI)

## Description

The `forge init` command (huh TUI) is missing two config sections that exist in the current `.forge/config.yaml` schema:

1. `auto.validation` (ModeToggle with quick/full) — exists in `AutoConfig` struct but `askAutoBehavior()` never prompts for it.
2. `worktree` config (source-branch + copy-files) — exists in `forge config init` (stdin) but not in `forge init` (huh TUI).

## Reference Files
- `docs/proposals/forge-init-config-sync/proposal.md` — Source proposal
- `forge-cli/internal/cmd/init.go` — `askAutoBehavior()` at line 236, insert after line 308
- `forge-cli/pkg/forgeconfig/config.go` — `AutoConfig` struct with `Validation` field at line 35, `WorktreeConfig` struct
- `forge-cli/internal/cmd/config.go` — Reference for worktree stdin prompts at lines 129-138

## Acceptance Criteria
- [ ] `askAutoBehavior()` prompts for `auto.validation` quick/full using `huh.Select` (same pattern as e2eTest/consolidateSpecs/cleanCode)
- [ ] Validation prompts appear between cleanCode and gitPush prompts
- [ ] After auto behavior, a worktree section prompts for optional source-branch (text input) and copy-files (multi-select)
- [ ] Both worktree fields are skippable (empty source-branch + no copy-files = no worktree block in config)
- [ ] Existing tests in `init_test.go` still pass (update mocks if needed)

## Hard Rules
- Follow existing huh TUI patterns exactly (`huh.NewSelect`, `huh.NewInput`, etc.)
- Do NOT modify `forge config init` (that's Task 2)

## Implementation Notes
- See `forge-cli/internal/cmd/config.go` lines 129-138 for the worktree prompt structure (source-branch string, copy-files list).
- The `WorktreeConfig` struct in `forgeconfig/config.go` has `SourceBranch` and `CopyFiles` fields.
- Validation ModeToggle follows the same pattern as the other auto fields — look at how `cleanCode` is handled in `askAutoBehavior()`.
