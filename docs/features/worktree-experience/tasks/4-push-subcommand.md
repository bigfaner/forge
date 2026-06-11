---
id: "4"
title: "Add worktree push subcommand"
priority: "P2"
estimated_time: "1h"
dependencies: []
scope: "backend"
breaking: false
type: "coding.feature"
mainSession: false
---

# 4: Add worktree push subcommand

## Description

Add a new `forge worktree push` subcommand that pushes the current worktree's branch to the remote. This saves users from manually navigating to the worktree directory and running git push. The command detects the current worktree context and pushes accordingly.

## Reference Files
- `docs/proposals/worktree-experience/proposal.md` — Source proposal
- `forge-cli/internal/cmd/worktree.go` — Existing worktree commands
- `forge-cli/internal/cmd/feature_complete.go` — Existing `gitPush()` function (line 246)
- `forge-cli/pkg/git/git.go` — Git utilities

## Acceptance Criteria
- [ ] `forge worktree push` pushes the current worktree branch to origin with `-u` (set upstream)
- [ ] When run inside a worktree, automatically detects the branch name
- [ ] When run outside a worktree, prints an error message and exits
- [ ] Prints the push output (remote URL or branch info) for confirmation
- [ ] Handles push failure gracefully (network error, auth failure, rejected push)

## Hard Rules
- Must detect worktree context before pushing — refuse to push from main worktree's main branch
- Use `git push -u origin HEAD` pattern consistent with existing `gitPush()` in feature_complete.go

## Implementation Notes
- Can reuse or extract `gitPush()` from `feature_complete.go`.
- Detection: use `git.GetWorktreeName(projectRoot)` or check `.git` file to determine if inside a worktree.
- The command has no positional args (operates on current directory context).
