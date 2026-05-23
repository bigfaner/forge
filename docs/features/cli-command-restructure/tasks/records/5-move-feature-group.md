---
status: "completed"
started: "2026-05-23 02:18"
completed: "2026-05-23 03:03"
time_spent: "~45m"
---

# Task Record: 5 Move forge feature group to feature/ subdirectory

## Summary
Moved all forge feature subcommand files from internal/cmd/ to internal/cmd/feature/ subdirectory. Created a Register() function for the feature command group. Moved shared utility functions (CalcSlugColWidth, TruncateSlug, PadRight) to base package for cross-package access. Updated root.go to use the new feature package. Updated all callers of the moved utility functions.

## Changes

### Files Created
- forge-cli/internal/cmd/feature/feature.go
- forge-cli/internal/cmd/feature/feature_complete.go
- forge-cli/internal/cmd/feature/feature_test.go
- forge-cli/internal/cmd/feature/feature_complete_test.go
- forge-cli/internal/cmd/feature/testmain_test.go
- forge-cli/internal/cmd/testing_helpers_test.go

### Files Modified
- forge-cli/internal/cmd/root.go
- forge-cli/internal/cmd/output.go
- forge-cli/internal/cmd/base/output.go
- forge-cli/internal/cmd/proposal.go
- forge-cli/internal/cmd/lesson.go
- forge-cli/internal/cmd/proposal_test.go
- forge-cli/internal/cmd/slug_width_test.go
- forge-cli/internal/cmd/feature_test.go
- forge-cli/internal/cmd/integration_test.go

### Key Decisions
- Moved CalcSlugColWidth, TruncateSlug, PadRight to base package to avoid circular dependency (feature sub-package cannot import internal/cmd)
- Kept feature_test.go tests that use rootCmd (TestRunQuery, TestRunStatus, TestRunCheck) in cmd package since they test task commands
- New feature tests use Cmd.SetArgs directly instead of rootCmd since they are now in the feature sub-package

## Test Results
- **Tests Executed**: Yes
- **Passed**: 729
- **Failed**: 0
- **Coverage**: 84.6%

## Acceptance Criteria
- [x] All feature subcommand files are in internal/cmd/feature/
- [x] New package exports a Register() function
- [x] root.go updated to use the new package
- [x] go build ./... passes
- [x] go test ./... passes
- [x] forge feature subcommands work identically

## Notes
Hard rule satisfied: feature sub-package does NOT import internal/cmd (no circular deps). All 4 utility functions moved to base package are re-exported via cmd/output.go for backward compatibility with proposal.go and lesson.go.
