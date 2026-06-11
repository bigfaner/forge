---
status: "completed"
started: "2026-05-23 09:54"
completed: "2026-05-23 10:01"
time_spent: "~7m"
---

# Task Record: 1 Add forge task list subcommand

## Summary
Add forge task list subcommand that displays all tasks for the current feature in a table format with natural sort ordering

## Changes

### Files Created
- forge-cli/internal/cmd/task/list.go
- forge-cli/internal/cmd/task/list_test.go

### Files Modified
- forge-cli/internal/cmd/task/register.go
- forge-cli/scripts/version.txt

### Key Decisions
- Natural sort groups by numeric prefix: business IDs (1, 1.gate, 2) before T-prefixed IDs (T-1, T-2)
- Reused base.TruncateSlug and base.PadRight for formatting instead of local helpers
- No extra flags per Hard Rules -- initial version has no --status filter

## Test Results
- **Tests Executed**: Yes
- **Passed**: 13
- **Failed**: 0
- **Coverage**: 92.3%

## Acceptance Criteria
- [x] forge task list resolves feature, loads index.json, displays tasks in table
- [x] Table columns: ID, Type, Title (truncated), Status
- [x] Header shows total count and feature slug
- [x] Tasks sorted by ID in natural order: numeric first, then test/gate IDs
- [x] Features with no index.json or empty task list print clear message
- [x] Unit tests cover: normal list output, empty feature, sorted output order
- [x] Version bumped in scripts/version.txt (minor: new command)

## Notes
13 test cases: 3 metadata tests, 3 output tests, 1 sorting integration test, 1 title truncation test, 5 natural sort unit tests
