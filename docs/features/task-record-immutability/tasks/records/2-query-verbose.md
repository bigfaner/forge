---
status: "completed"
started: "2026-05-17 01:33"
completed: "2026-05-17 01:43"
time_spent: "~10m"
---

# Task Record: 2 Query verbose mode with RELATED_FIXES

## Summary
Extended forge task query with --verbose/-v flag. Default output unchanged (TASK_ID, STATUS, SCOPE if set, BREAKING if true). Verbose mode adds KEY, TITLE, PRIORITY, TYPE, SCOPE, DEPENDENCIES (multi-line), TASK_FILE, RECORD_FILE, and RELATED_FIXES reverse lookup. Refactored runQuery into printDefaultQuery and printVerboseQuery functions.

## Changes

### Files Created
- forge-cli/internal/cmd/query_test.go

### Files Modified
- forge-cli/internal/cmd/query.go

### Key Decisions
- Used BoolVarP for --verbose/-v flag registration via init() function
- Extracted printDefaultQuery and printVerboseQuery to keep runQuery clean
- RELATED_FIXES sorted by task ID for deterministic output
- DEPENDENCIES uses multi-line format with indented items when present
- Used existing PrintField/PrintFieldIfNotEmpty/PrintBlockStart/PrintBlockEnd helpers per hard rules

## Test Results
- **Tests Executed**: Yes
- **Passed**: 10
- **Failed**: 0
- **Coverage**: 81.0%

## Acceptance Criteria
- [x] forge task query <id> output unchanged (TASK_ID, STATUS, SCOPE if set, BREAKING if true)
- [x] forge task query <id> --verbose displays: KEY, TASK_ID, TITLE, STATUS, PRIORITY, TYPE, SCOPE (if set), DEPENDENCIES (multi-line if multiple), TASK_FILE, RECORD_FILE
- [x] forge task query <id> --verbose shows RELATED_FIXES when fix tasks exist: <id> [<status>] <title> per line
- [x] forge task query <id> --verbose omits RELATED_FIXES when no fixes exist
- [x] forge task query <id> -v works as shorthand for --verbose
- [x] TASK_FILE and RECORD_FILE paths are constructed from feature slug + task File/Record fields

## Notes
无
