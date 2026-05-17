---
id: "2"
title: "Implement forge worktree list subcommand"
priority: "P0"
estimated_time: "1h"
dependencies: ["1"]
scope: "backend"
breaking: false
type: "feature"
mainSession: false
---

# 2: Implement forge worktree list subcommand

## Description

Implement `forge worktree list` that shows all git worktrees with name, branch, and path. Worktrees whose name matches a feature slug in `docs/features/` are visually marked as forge-managed.

## Reference Files
- `docs/proposals/forge-worktree/proposal.md` — Source proposal
- `forge-cli/pkg/git/git.go` — Run() helper for git commands
- `forge-cli/internal/cmd/root.go` — Command registration pattern

## Acceptance Criteria

- [ ] `forge worktree list` displays all git worktrees (from `git worktree list`)
- [ ] Each entry shows: worktree name, branch, path
- [ ] Entries whose worktree name matches a directory in `docs/features/` are marked as forge-managed
- [ ] Main worktree (current project) is included but distinguished from feature worktrees
- [ ] No worktrees → prints "No worktrees found" message

## Hard Rules

- Parse `git worktree list --porcelain` for reliable machine-readable output (not the human-readable default)
- Extract worktree name from the `.git/worktrees/<name>` path or from the directory basename

## Implementation Notes

- `git worktree list --porcelain` outputs blocks separated by blank lines, each block has lines like `worktree /path`, `HEAD abc123`, `branch refs/heads/feature-x`
- To determine forge-managed: check if `docs/features/<worktree_basename>/` exists relative to the main project root
- Use tabwriter or simple aligned columns for output formatting
