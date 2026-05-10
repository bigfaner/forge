---
status: "completed"
started: "2026-05-10 15:26"
completed: "2026-05-10 15:29"
time_spent: "~3m"
---

# Task Record: fix-4 Fix: TestCheckAllCompleted_NoProject env leakage

## Summary
Fix TestCheckAllCompleted_NoProject env leakage by converting to subprocess pattern

## Changes

### Files Created
无

### Files Modified
- task-cli/internal/cmd/all_completed_test.go

### Key Decisions
- Used subprocess pattern (same as TestRunValidate_NoProjectRoot and TestRunCheck_NoProjectRoot) to isolate the test from ancestor project markers that cause FindProjectRoot to walk up and find Z:\project\ai\forge

## Test Results
- **Tests Executed**: Yes
- **Passed**: 1
- **Failed**: 0
- **Coverage**: 85.7%

## Acceptance Criteria
- [x] TestCheckAllCompleted_NoProject passes reliably without picking up ancestor project markers

## Notes
Root cause: t.Setenv("CLAUDE_PROJECT_DIR", "") sets env to empty string but FindProjectRoot falls back to directory walk from cwd, finding ancestor project markers in Z:\project\ai\forge. Fix: run test as subprocess with clean env and isolated tmpDir as working directory.
