---
id: "1"
title: "Implement worktree command group scaffold + start subcommand"
priority: "P0"
estimated_time: "2h"
dependencies: []
scope: "backend"
breaking: false
type: "feature"
mainSession: false
---

# 1: Implement worktree command group scaffold + start subcommand

## Description

Create the `forge worktree` cobra command group and implement the `start` subcommand. This is the foundational task that establishes the command structure and implements the primary workflow: create a git worktree with branch = slug in a sibling directory, then launch `claude` with `--dangerously-skip-permissions`.

## Reference Files
- `docs/proposals/forge-worktree/proposal.md` — Source proposal
- `forge-cli/internal/cmd/root.go` — Command registration pattern
- `forge-cli/internal/cmd/claude.go` — Claude launch pattern (lookPathFunc, runClaudeFunc, --dangerously-skip-permissions)
- `forge-cli/pkg/git/git.go` — GetWorktreeName(), Run() helper
- `forge-cli/pkg/project/root.go` — FindProjectRoot()

## Acceptance Criteria

- [ ] `forge worktree start <slug>` creates a git worktree at `../<slug>` with branch `<slug>` from HEAD
- [ ] If branch `<slug>` already exists, creates worktree from that branch (resume context)
- [ ] If worktree directory already exists, errors with hint to use `resume`
- [ ] If target directory name conflicts with existing non-worktree directory, errors with clear message
- [ ] Launches `claude --dangerously-skip-permissions` in the worktree directory after creation
- [ ] Detects `claude` binary availability before attempting launch (like `forge claude` does)
- [ ] `forge feature` in the created worktree auto-detects the correct feature via `GetWorktreeName()`

## Hard Rules

- Use `filepath.Join()` for all path construction (Windows compatibility)
- Use direct `git worktree add` via `pkg/git/Run()` — NOT Claude's native `--worktree` mechanism
- Reuse `lookPathFunc` / `runClaudeFunc` pattern from `claude.go` for testability
- Register `worktreeCmd` in `root.go` init() following existing pattern

## Implementation Notes

- The command group (`worktreeCmd`) is a parent cobra command with no `Run` function, only subcommands
- `start` must compute the sibling directory as `filepath.Join(projectRoot, "..", slug)` — use `FindProjectRoot()` to get projectRoot
- Pre-flight checks: verify `claude` in PATH, verify target dir doesn't exist, verify git worktree add succeeds
- Cross-platform risk: `../slug` resolves differently on Windows vs Unix if projectRoot contains symlinks. Use `filepath.Abs()` on the result
