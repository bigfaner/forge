---
id: "3"
title: "Implement forge worktree remove subcommand"
priority: "P0"
estimated_time: "1h"
dependencies: ["1"]
scope: "backend"
breaking: false
type: "feature"
mainSession: false
---

# 3: Implement forge worktree remove subcommand

## Description

Implement `forge worktree remove <slug>` that removes the git worktree while preserving the branch for manual merge.

## Reference Files
- `docs/proposals/forge-worktree/proposal.md` — Source proposal
- `forge-cli/pkg/git/git.go` — Run() helper for git commands

## Acceptance Criteria

- [ ] `forge worktree remove <slug>` removes the git worktree at `../<slug>`
- [ ] Branch is preserved after removal (not deleted)
- [ ] Errors if the specified worktree does not exist
- [ ] Errors if the specified worktree has uncommitted changes, with hint to commit or stash
- [ ] Prints confirmation with branch name after removal

## Hard Rules

- Use `git worktree remove` command, not manual directory deletion
- Do NOT auto-delete the branch — always keep it

## Implementation Notes

- `git worktree remove <path>` fails if there are uncommitted changes. The `--force` flag bypasses this, but we should NOT use it by default — instead surface the error message
- After removal, print the branch name so the user can `git merge <branch>` later
- Resolve worktree path the same way as `start`: `filepath.Join(projectRoot, "..", slug)`
