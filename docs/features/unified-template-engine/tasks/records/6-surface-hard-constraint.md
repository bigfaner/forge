---
status: "completed"
started: "2026-05-28 01:48"
completed: "2026-05-28 01:59"
time_spent: "~11m"
---

# Task Record: 6 Surface 硬性约束与 forge init 集成

## Summary
Implemented surface hard constraint in quality_gate.go: addSingleFixTask() now returns error when surface inference fails via new requireSurfaceInference() function. forge task add path retains soft behavior. forge init already includes surface config step (verified).

## Changes

### Files Created
无

### Files Modified
- forge-cli/internal/cmd/quality_gate.go
- forge-cli/internal/cmd/quality_gate_test.go
- forge-cli/scripts/version.txt

### Key Decisions
- requireSurfaceInference() wraps inferSurface() with hard-failure policy, returning error with 'forge surfaces detect' guidance
- helperSetup() now configures default scalar surface ('.': 'cli') so existing addFixTask tests pass under hard constraint
- Tests with no extractable source files (empty output, no file references) now expect errors per hard constraint
- Revert strategy: single-line change from requireSurfaceInference back to inferSurface restores soft behavior

## Test Results
- **Tests Executed**: Yes
- **Passed**: 72
- **Failed**: 0
- **Coverage**: 66.4%

## Acceptance Criteria
- [x] quality_gate.go addSingleFixTask() surface inference failure returns error, no longer creates empty-surface task
- [x] forge task add command path retains soft behavior (empty string on inference failure)
- [x] forge init includes surface config step (forge surfaces detect)
- [x] Hard failure logic encapsulated in independent requireSurfaceInference() function with 'forge surfaces detect' guidance

## Notes
AC-2 (soft behavior) and AC-3 (init surface step) were already implemented in existing code. This task primarily added AC-1 (hard failure) and AC-4 (requireSurfaceInference function). Version bumped to 5.13.0 (minor: new behavior).
