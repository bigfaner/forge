---
status: "completed"
started: "2026-05-16 14:22"
completed: "2026-05-16 14:29"
time_spent: "~7m"
---

# Task Record: 2 Adjust GetCurrentFeature() to read state.json first

## Summary
Modified GetCurrentFeature() to check .forge/state.json as highest-priority source before falling back to git context. Added GetCurrentFeatureWithSource() returning (slug, source, error) where source is one of: state.json, worktree, branch, features-dir. GetCurrentFeature() delegates to GetCurrentFeatureWithSource() internally, preserving the existing API signature unchanged.

## Changes

### Files Created
无

### Files Modified
- forge-cli/pkg/feature/feature.go
- forge-cli/pkg/feature/feature_test.go

### Key Decisions
- GetCurrentFeatureWithSource() calls ReadForgeState() which already silently ignores corrupt/missing files, satisfying the hard rule about silent failure
- When state.json feature directory doesn't exist, we skip to next priority without auto-creating, per hard rule
- Git source distinguishes worktree vs branch by checking GetWorktreeName() separately after GetFeatureFromGit() returns non-empty

## Test Results
- **Tests Executed**: Yes
- **Passed**: 20
- **Failed**: 0
- **Coverage**: 91.8%

## Acceptance Criteria
- [x] GetCurrentFeature() returns the feature from state.json when it exists and the feature directory is valid
- [x] When state.json is absent, behavior falls back to git context identical to current behavior
- [x] GetCurrentFeatureWithSource() returns (slug, source, error) where source is one of: state.json, worktree, branch, features-dir
- [x] GetCurrentFeature() calls GetCurrentFeatureWithSource() internally, preserving existing API
- [x] All existing callers of GetCurrentFeature() continue to work without changes
- [x] Existing tests pass without modification (state.json absent = identical behavior)

## Notes
无
