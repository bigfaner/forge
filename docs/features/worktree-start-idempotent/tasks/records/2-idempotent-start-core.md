---
status: "completed"
started: "2026-06-09 14:29"
completed: "2026-06-09 14:37"
time_spent: "~8m"
---

# Task Record: 2 Make `forge worktree start` idempotent — core behavior

## Summary
Made `forge worktree start` idempotent: when worktree already exists, skip creation and launch fresh Claude session instead of erroring out. Validates existing worktree via symlink resolve + .git check, skips includes copy, ignores --source-branch with warning.

## Changes

### Files Created
无

### Files Modified
- forge-cli/internal/cmd/worktree/cmd_start.go
- forge-cli/internal/cmd/worktree/worktree_test.go

### Key Decisions
- Reused cmd_resume.go validation pattern (filepath.EvalSymlinks + os.Stat(.git)) for existing worktree verification
- Existing worktree path returns early before config loading, includes validation, and branch resolution
- Added 'entering existing worktree' and 'created new worktree' keywords to stderr for script-parseable output

## Test Results
- **Tests Executed**: Yes
- **Passed**: 6
- **Failed**: 0
- **Coverage**: 87.4%

## Acceptance Criteria
- [x] start on existing worktree skips creation and launches fresh Claude session
- [x] stderr contains 'entering existing worktree' when existing, 'created new worktree' when new
- [x] includes copy skipped when worktree already exists
- [x] corrupt worktree (missing .git) errors with 'forge worktree remove' suggestion
- [x] worktree not exist behavior unchanged (regression tests pass)

## Notes
All existing worktree tests pass (full suite 87.4% coverage). Replaced TestWorktreeStart_ErrorWhenTargetDirExists with TestWorktreeStart_EntersExistingWorktree plus 5 new tests for idempotent paths.
