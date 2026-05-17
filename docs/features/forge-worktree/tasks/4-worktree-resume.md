---
id: "4"
title: "Implement forge worktree resume subcommand"
priority: "P1"
estimated_time: "1h"
dependencies: ["1"]
scope: "backend"
breaking: false
type: "feature"
mainSession: false
---

# 4: Implement forge worktree resume subcommand

## Description

Implement `forge worktree resume <slug>` that re-launches `claude` with `--dangerously-skip-permissions` in an existing worktree directory.

## Reference Files
- `docs/proposals/forge-worktree/proposal.md` — Source proposal
- `forge-cli/internal/cmd/claude.go` — Claude launch pattern

## Acceptance Criteria

- [ ] `forge worktree resume <slug>` launches `claude --dangerously-skip-permissions` in the worktree directory
- [ ] Errors if the specified worktree does not exist
- [ ] Detects `claude` binary availability before attempting launch

## Hard Rules

- Reuse `lookPathFunc` / `runClaudeFunc` pattern from `claude.go` for testability
- Use `filepath.Join()` for path construction

## Implementation Notes

- This is essentially the second half of `start` — same claude launch logic, but skips worktree creation
- Verify the worktree exists by checking if the directory exists and is a git worktree (`.git` file exists)
- Could extract a shared `launchClaude(dir string)` helper used by both `start` and `resume`
