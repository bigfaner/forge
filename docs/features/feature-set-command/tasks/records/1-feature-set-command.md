---
status: "completed"
started: "2026-05-16 14:09"
completed: "2026-05-16 14:21"
time_spent: "~12m"
---

# Task Record: 1 Add forge feature set subcommand

## Summary
Add forge feature set <slug> subcommand that writes feature slug to .forge/state.json via EnsureForgeState() and ensures feature directory structure via EnsureFeatureDir(). Added exactArgsNonEmpty() Cobra validator for empty-slug rejection at the arg level. Registered featureSetCmd in init() alongside existing subcommands. Existing positional arg behavior (forge feature <slug>) unchanged.

## Changes

### Files Created
无

### Files Modified
- forge-cli/internal/cmd/feature.go
- forge-cli/internal/cmd/feature_test.go
- forge-cli/scripts/version.txt

### Key Decisions
- Used exactArgsNonEmpty() custom Cobra Args validator to reject empty slugs at the arg-parsing level (returns error instead of os.Exit), making the empty-slug path testable
- Kept slug == "" defense-in-depth check in runFeatureSet per Hard Rules, even though Cobra validation catches it first
- Bumped version from 3.13.0 to 3.14.0 (minor: new subcommand)

## Test Results
- **Tests Executed**: Yes
- **Passed**: 3
- **Failed**: 0
- **Coverage**: 80.7%

## Acceptance Criteria
- [x] forge feature set <slug> creates the feature directory structure under docs/features/<slug>/
- [x] forge feature set <slug> writes .forge/state.json with feature=<slug> and allCompleted=false
- [x] forge feature set <slug> returns an error when slug is empty
- [x] forge feature set <slug> prints the feature slug to stdout on success
- [x] Existing forge feature <slug> positional arg behavior unchanged (backward compatible)

## Notes
无
