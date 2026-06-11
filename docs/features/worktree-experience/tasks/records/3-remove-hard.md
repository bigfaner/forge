---
status: "completed"
started: "2026-05-20 17:38"
completed: "2026-05-20 17:45"
time_spent: "~7m"
---

# Task Record: 3 Remove --hard: delete worktree + branch + prune

## Summary
Added --hard and --force flags to 'forge worktree remove'. --hard performs three-step cleanup: worktree remove, local branch deletion, and git worktree prune. --force overrides uncommitted changes check. Branch deletion uses safe -d first, falls back to -D for unmerged branches with a warning. Behavior without --hard is unchanged.

## Changes

### Files Created
无

### Files Modified
- forge-cli/internal/cmd/worktree.go
- forge-cli/internal/cmd/worktree_test.go

### Key Decisions
- Extracted hard cleanup logic into separate runHardCleanup function for clarity
- Branch deletion tries safe -d first, then -D for unmerged (per Hard Rules: --hard without --force warns but allows unmerged deletion, only blocks on uncommitted changes)
- Only local branches are deleted -- never remote (Hard Rule)
- Added --force flag to worktreeRemoveCmd alongside --hard for overriding uncommitted changes protection

## Test Results
- **Tests Executed**: Yes
- **Passed**: 12
- **Failed**: 0
- **Coverage**: 5.0%

## Acceptance Criteria
- [x] forge worktree remove <slug> --hard performs three steps: worktree remove -> branch delete -> prune
- [x] Each step reports its status (success or skip with reason)
- [x] If worktree removal fails, branch deletion and prune are skipped
- [x] Before deleting branch, check for uncommitted changes in the worktree (fail with message if found)
- [x] Before deleting branch, check if branch is merged into source-branch (warn if not, require --force to proceed)
- [x] forge worktree remove <slug> (without --hard) behavior is unchanged

## Notes
无
