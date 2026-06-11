---
id: "5"
title: "Add worktree status subcommand"
priority: "P2"
estimated_time: "1.5h"
dependencies: []
scope: "backend"
breaking: false
type: "coding.feature"
mainSession: false
---

# 5: Add worktree status subcommand

## Description

Add a new `forge worktree status [<slug>]` subcommand that displays the status of a worktree: branch name, latest commit, and uncommitted file list. When no slug is provided, shows status for all forge-managed worktrees. This is a read-only command that never modifies filesystem state.

## Reference Files
- `docs/proposals/worktree-experience/proposal.md` — Source proposal
- `forge-cli/internal/cmd/worktree.go` — Existing worktree commands
- `forge-cli/pkg/git/worktree_list.go` — Worktree listing (WorktreeEntry type)
- `forge-cli/pkg/git/git.go` — Git utilities

## Acceptance Criteria
- [ ] `forge worktree status <slug>` shows: branch name, latest commit (hash + message), uncommitted files list
- [ ] `forge worktree status` (no slug) shows status for all forge-managed worktrees
- [ ] Output uses structured format (consistent with project's output.go pattern)
- [ ] Non-existent slug prints clear error message
- [ ] Command is strictly read-only — never modifies any file

## Hard Rules
- Status command must not modify any filesystem state — no side effects
- Completion response time should be < 200ms

## Implementation Notes
- Use `git.ListWorktrees()` to resolve slug to path.
- Branch name from `WorktreeEntry.Branch` or `git.GetCurrentBranch()`.
- Latest commit: `git log -1 --oneline` in the worktree directory.
- Uncommitted files: `git status --porcelain` in the worktree directory.
- Use `PrintBlockStart/PrintField/PrintBlockEnd` from output.go for consistent formatting.
