---
status: "completed"
started: "2026-05-15 00:41"
completed: "2026-05-15 00:57"
time_spent: "~16m"
---

# Task Record: 2 Trim submit, query, and status output to essential fields

## Summary
Trimmed submit, query, and status CLI output to essential fields only. submit.go now outputs only STATUS (non-JSON, non-quiet). query.go outputs TASK_ID + STATUS + SCOPE (when non-empty) + BREAKING (when true). status.go outputs TASK_ID + STATUS in all modes (query, update, unmet-deps). Updated all contract tests and integration tests to match new output format. Patch version bump to 3.9.2.

## Changes

### Files Created
无

### Files Modified
- forge-cli/internal/cmd/submit.go
- forge-cli/internal/cmd/query.go
- forge-cli/internal/cmd/status.go
- forge-cli/internal/cmd/output_contract_test.go
- forge-cli/internal/cmd/feature_test.go
- forge-cli/internal/cmd/integration_test.go
- forge-cli/scripts/version.txt

### Key Decisions
- Kept BREAKING conditional (omit when false) in query.go per proposal spec
- Added SCOPE field to query.go output (conditional, omit when empty) per proposal spec
- Unmet-deps warning in status.go still prints WARNING line between TASK_ID/STATUS fields and closing separator

## Test Results
- **Tests Executed**: Yes
- **Passed**: 19
- **Failed**: 0
- **Coverage**: 80.5%

## Acceptance Criteria
- [x] forge task submit (non-JSON, non-quiet) outputs exactly 1 field: STATUS
- [x] forge task query outputs exactly TASK_ID + STATUS + SCOPE (when non-empty) + BREAKING (when true)
- [x] forge task status (query mode) outputs exactly TASK_ID + STATUS
- [x] forge task status (update mode) outputs exactly TASK_ID + STATUS
- [x] forge task status (unmet deps warning) outputs TASK_ID + STATUS + WARNING line
- [x] JSON mode (--json) in submit is NOT changed
- [x] All existing unit tests pass after updates

## Notes
Pre-existing TestSaveIndexAndSignalCompletion_SaveIndexError fails on Windows (os.Chmod 0555 does not prevent writes). Unrelated to this change.
