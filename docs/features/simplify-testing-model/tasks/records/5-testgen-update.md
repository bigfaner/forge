---
status: "completed"
started: "2026-05-17 15:01"
completed: "2026-05-17 15:13"
time_spent: "~12m"
---

# Task Record: 5 Update testgen.go to use Language instead of ProfileName

## Summary
Updated testgen.go to use profile.Language instead of string-based ProfileName. Changed GetBreakdownTestTasks and GetQuickTestTasks signatures from (profiles []string, capabilities []string) to (languages []profile.Language, interfaces []string). Replaced TestTaskDef.ProfileName with TestTaskDef.Language of type profile.Language. Updated build.go adapter to convert []string to []profile.Language. Updated all test files (testgen_test.go, autoconfig_test.go) to use the new signatures. Output format (frontmatter profile field, task IDs, keys) remains unchanged.

## Changes

### Files Created
无

### Files Modified
- forge-cli/pkg/task/testgen.go
- forge-cli/pkg/task/testgen_test.go
- forge-cli/pkg/task/autoconfig_test.go
- forge-cli/pkg/task/build.go

### Key Decisions
- Used profile.Language type for TestTaskDef.Language field rather than a local type alias, maintaining single source of truth in the profile package
- Kept frontmatter field name as 'profile' (not 'language') to satisfy the hard rule that task file template structure must not change
- Added string() conversion in build.go adapter layer (generateTestTasks) to bridge existing []string callers to new []profile.Language API

## Test Results
- **Tests Executed**: Yes
- **Passed**: 207
- **Failed**: 0
- **Coverage**: 89.4%

## Acceptance Criteria
- [x] GetBreakdownTestTasks signature uses []Language and []string (interfaces) parameters
- [x] Per-language-per-interface task expansion produces correct test tasks
- [x] TestTaskDef.ProfileName field replaced with TestTaskDef.Language of type Language
- [x] Task IDs and slugs use language keys (not profile names)
- [x] Existing test generation test cases updated and passing
- [x] go test ./... passes

## Notes
The hard rule 'task generation output must be compatible with forge task index command' is preserved: task IDs, slugs, frontmatter format, and key naming are unchanged. Only the internal Go types changed.
