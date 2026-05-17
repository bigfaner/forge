---
status: "completed"
started: "2026-05-17 17:44"
completed: "2026-05-17 17:58"
time_spent: "~14m"
---

# Task Record: 3 Implement forge worktree remove subcommand

## Summary
Implemented forge worktree remove <slug> subcommand that removes a git worktree at ../<slug> while preserving the branch. Uses git worktree remove (not manual directory deletion). Errors on missing worktree, uncommitted changes (with hint to commit or stash), and non-git repos. Prints confirmation with branch name after removal.

## Changes

### Files Created
无

### Files Modified
- forge-cli/internal/cmd/worktree.go
- forge-cli/internal/cmd/worktree_test.go
- forge-cli/internal/cmd/root.go

### Key Decisions
- Used git.Run() for git worktree remove to get structured error output, consistent with existing start command
- Resolved worktree branch name before removal by querying listWorktreesFunc, so confirmation message shows actual branch name
- Detected dirty working tree by checking git error output for 'dirty'/'modified'/'local changes' keywords, then surfaced actionable hint

## Test Results
- **Tests Executed**: Yes
- **Passed**: 6
- **Failed**: 0
- **Coverage**: 80.3%

## Acceptance Criteria
- [x] forge worktree remove <slug> removes the git worktree at ../<slug>
- [x] Branch is preserved after removal (not deleted)
- [x] Errors if the specified worktree does not exist
- [x] Errors if the specified worktree has uncommitted changes, with hint to commit or stash
- [x] Prints confirmation with branch name after removal

## Notes
无
