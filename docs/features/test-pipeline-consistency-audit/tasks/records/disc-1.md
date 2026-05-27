---
status: "completed"
started: "2026-05-27 15:55"
completed: "2026-05-27 16:06"
time_spent: "~11m"
---

# Task Record: disc-1 Fix: pre-existing TestDetectModeFromPath failure

## Summary
Fixed two pre-existing test failures in TestDetectModeFromPath: Windows backslash path handling and symlink resolution for non-existent subdirectories

## Changes

### Files Created
无

### Files Modified
- forge-cli/internal/cmd/config.go

### Key Decisions
- Added resolveExistingPrefix helper to incrementally resolve symlinks only for existing path components, avoiding EvalSymlinks failure on non-existent tails
- Normalized backslashes to OS path separator at the start of detectModeFromPath, before any filepath operations

## Test Results
- **Tests Executed**: Yes
- **Passed**: 9
- **Failed**: 0
- **Coverage**: 9.3%

## Acceptance Criteria
无

## Notes
Root cause: (1) filepath.ToSlash does not convert backslashes on Unix, so Windows-style paths with backslashes never matched /docs/features/ pattern. (2) filepath.EvalSymlinks fails when the full path does not exist on disk, causing symlink test to fall back to the un-resolved symlink path which lacks /docs/features/.
