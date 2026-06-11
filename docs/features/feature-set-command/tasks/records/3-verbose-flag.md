---
status: "completed"
started: "2026-05-16 14:31"
completed: "2026-05-16 14:39"
time_spent: "~8m"
---

# Task Record: 3 Add verbose flag to forge feature command

## Summary
Add -v (verbose) flag to bare forge feature command that displays feature name with resolution source (state.json, worktree, branch, features-dir). Flag registered as local flag only, does not leak to subcommands. Uses GetCurrentFeatureWithSource() from Task 2.

## Changes

### Files Created
无

### Files Modified
- forge-cli/internal/cmd/feature.go
- forge-cli/internal/cmd/feature_test.go

### Key Decisions
- Registered -v as local flag via Flags().BoolVarP, not PersistentFlags, to prevent leakage to set/list/status subcommands
- Reset verbose global var at test start to prevent cross-test state leakage from cobra's flag binding behavior

## Test Results
- **Tests Executed**: Yes
- **Passed**: 5
- **Failed**: 0
- **Coverage**: 80.8%

## Acceptance Criteria
- [x] forge feature -v shows FEATURE: my-feature (from: state.json) when resolved from state.json
- [x] forge feature -v shows FEATURE: my-feature (from: worktree) when resolved from git worktree
- [x] forge feature -v shows FEATURE: my-feature (from: branch) when resolved from git branch
- [x] forge feature -v shows FEATURE: my-feature (from: features-dir) when resolved from directory scanning
- [x] forge feature -v shows FEATURE: (none) when no feature is set
- [x] forge feature (without -v) behavior unchanged — shows only the slug
- [x] -v flag only applies to bare forge feature — does not affect set, list, or status subcommands

## Notes
Worktree and branch source tests covered by existing GetCurrentFeatureWithSource tests in pkg/feature. CMD-layer tests verify state.json and features-dir sources plus the (none) case and non-leakage to subcommands.
