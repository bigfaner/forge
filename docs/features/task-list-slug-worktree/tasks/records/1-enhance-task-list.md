---
status: "completed"
started: "2026-05-24 10:44"
completed: "2026-05-24 10:49"
time_spent: "~5m"
---

# Task Record: 1 Enhance forge task list with slug parameter and worktree-aware reading

## Summary
Enhanced forge task list with optional slug parameter and --local flag for worktree-aware index reading

## Changes

### Files Created
无

### Files Modified
- forge-cli/internal/cmd/task/list.go
- forge-cli/internal/cmd/task/list_test.go
- forge-cli/scripts/version.txt
- docs/conventions/forge-cli-reference.md

### Key Decisions
- Used cobra.MaximumNArgs(1) instead of cobra.NoArgs to accept optional positional slug
- Added resolveListIndexPath() to handle worktree-aware index.json resolution with fast-path directory check
- Reused existing base.ErrFeatureNotFound() for slug validation errors
- Added --local bool flag via init() to bypass worktree reading

## Test Results
- **Tests Executed**: Yes
- **Passed**: 8
- **Failed**: 0
- **Coverage**: 85.0%

## Acceptance Criteria
- [x] forge task list my-feature shows tasks for the specified feature, reading from worktree if one exists
- [x] forge task list my-feature --local reads from main repo's index.json regardless of worktree existence
- [x] forge task list (no args) behaves exactly as before — existing auto-detection logic unchanged
- [x] Clear error message when slug doesn't match any docs/features/<slug>/ directory
- [x] no tasks found message when feature exists but has no index.json
- [x] Unit tests cover: slug resolution, worktree detection, --local override, backward compatibility
- [x] Version bumped in scripts/version.txt per semver (minor: new user-facing capability)
- [x] docs/conventions/forge-cli-reference.md updated to reflect new usage

## Notes
无
