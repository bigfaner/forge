---
id: "1"
title: "Refactor start: branch-first creation + --no-launch flag"
priority: "P1"
estimated_time: "1.5h"
dependencies: []
scope: "backend"
breaking: false
type: "coding.enhancement"
mainSession: false
---

# 1: Refactor start: branch-first creation + --no-launch flag

## Description

Refactor `forge worktree start` in two ways:
1. Change branch creation order: first `git checkout -b <slug>` then `git worktree add`, making the branch creation step observable before worktree creation. Currently the flow is `git worktree add -b <slug>` which bundles branch + worktree creation atomically.
2. Add `--no-launch` flag to skip launching claude after worktree creation. This enables CI scripts and batch operations that only need the environment, not an interactive session.

## Reference Files
- `docs/proposals/worktree-experience/proposal.md` — Source proposal
- `forge-cli/internal/cmd/worktree.go` — Start command implementation (runWorktreeStart)
- `forge-cli/internal/cmd/worktree_test.go` — Existing tests

## Acceptance Criteria
- [ ] Branch creation uses `git checkout -b` before `git worktree add`, decoupling the two steps
- [ ] `forge worktree start <slug> --no-launch` creates the worktree and exits without launching claude
- [ ] Without `--no-launch`, behavior is unchanged (launches claude as before)
- [ ] Branch resolution logic (local existing / remote existing / new from source) still works correctly
- [ ] Copy-files still applied after worktree creation
- [ ] All existing tests pass; new behavior has unit tests

## Hard Rules
- Do NOT change the default behavior when `--no-launch` is not specified
- Branch-first creation must handle the case where the branch already exists (skip checkout)

## Implementation Notes
- The three-layer branch resolution in `runWorktreeStart` (local/remote/new) needs to be refactored. The "new from source" path currently uses `git worktree add -b <slug> <dir> <source>`. This should become: `git branch <slug> <source>` (or `git checkout -b <slug> <source>`) then `git worktree add <dir> <slug>`.
- The `--no-launch` flag wraps the claude launch section. When set, print the worktree path and exit.
- Key risk: branch-first creation changes error behavior — if branch creation succeeds but worktree add fails, cleanup needs to remove the created branch.
