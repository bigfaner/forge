---
status: "completed"
started: "2026-05-17 19:56"
completed: "2026-05-17 20:16"
time_spent: "~20m"
---

# Task Record: 2 Add --source-branch flag to worktree start command

## Summary
Add --source-branch / -b CLI flag to forge worktree start command. Implement source resolution with priority: flag > config > HEAD. Pre-validate branch exists via git rev-parse before worktree creation. Source branch applies to new-branch path only (existing branches ignore it). Updated CLI help text and Long description.

## Changes

### Files Created
无

### Files Modified
- forge-cli/internal/cmd/worktree.go
- forge-cli/internal/cmd/worktree_test.go

### Key Decisions
- Used cmd.Flags().Changed("source-branch") instead of checking string value to avoid Cobra flag state leakage between test runs
- Added resetSourceBranchFlag test helper to prevent Cobra Command flag persistence across sequential tests
- Pre-validate source branch via git rev-parse --verify before calling git worktree add, per Hard Rules

## Test Results
- **Tests Executed**: Yes
- **Passed**: 11
- **Failed**: 0
- **Coverage**: 80.3%

## Acceptance Criteria
- [x] forge worktree start <slug> --source-branch develop creates worktree from develop branch
- [x] forge worktree start <slug> -b v3.0.0 creates worktree from v3.0.0 branch
- [x] Without flag or config, behavior is identical to current version (creates from HEAD)
- [x] worktree.source-branch in config.yaml sets default source branch
- [x] Flag overrides config default
- [x] Clear error message when specified branch does not exist locally or remotely
- [x] CLI help text updated to show --source-branch flag and usage
- [x] Long description of worktreeStartCmd updated to reflect source-branch support

## Notes
Hard Rules followed: (1) pre-validate branch via git rev-parse --verify before worktree creation, (2) source branch applies to new-branch path only -- existing branches use existing-branch path unchanged.
