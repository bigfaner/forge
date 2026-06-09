---
status: "completed"
started: "2026-06-09 16:23"
completed: "2026-06-09 16:28"
time_spent: "~5m"
---

# Task Record: fix-1 Fix: forge worktree remove fails on corrupted worktrees

## Summary
Added corrupted worktree fallback in cmd_remove.go: when git worktree remove fails with .git validation errors, fall back to os.RemoveAll + git worktree prune for manual cleanup

## Changes

### Files Created
无

### Files Modified
- forge-cli/internal/cmd/worktree/cmd_remove.go

### Key Decisions
- Detected corrupted worktree by matching git error messages (.git, validation failed, not a valid, could not identify) and falling back to os.RemoveAll + git worktree prune instead of returning error

## Test Results
- **Tests Executed**: Yes
- **Passed**: 10
- **Failed**: 0
- **Coverage**: 0.0%

## Acceptance Criteria
- [x] forge worktree remove handles corrupted worktrees (missing .git) via fallback to os.RemoveAll + git worktree prune
- [x] Existing uncommitted-changes error handling preserved
- [x] Compilation and lint pass
- [x] TestStep2_Success_RemovesCorruptedWorktree passes
- [x] TestJourney_CorruptedWorktreeRecovery passes

## Notes
All 10 corrupted-worktree-recovery tests PASS including the 2 previously failing tests (TestStep2_Success_RemovesCorruptedWorktree and TestJourney_CorruptedWorktreeRecovery).
