---
status: "completed"
started: "2026-05-23 10:02"
completed: "2026-05-23 10:07"
time_spent: "~5m"
---

# Task Record: T-quick-doc-drift Detect Spec Drift

## Summary
Detected and fixed spec drift in 3 files: (1) forge-cli-reference.md had orphaned commands (forge probe, entire forge e2e group) and missing forge task list command; (2) error-reporting.md BIZ-001 had outdated exit code 2 semantics; (3) task-lifecycle.md BIZ-002 referenced wrong error type string. All fixes committed with [auto-specs] tag.

## Changes

### Files Created
无

### Files Modified
- docs/business-rules/error-reporting.md
- docs/business-rules/task-lifecycle.md
- docs/conventions/forge-cli-reference.md

### Key Decisions
- Moved orphaned forge e2e commands to removed-commands section rather than silently deleting, for traceability
- Updated exit code semantics to match AIError.ExitCode() implementation rather than outdated Cobra-only description

## Test Results
- **Tests Executed**: No
- **Passed**: 0
- **Failed**: 0
- **Coverage**: N/A (task has no tests)

## Acceptance Criteria
- [x] All project-level spec files validated against current codebase
- [x] Drifted rules updated to match current code
- [x] Orphaned entries removed or moved to removed-commands section
- [x] Changes committed with [auto-specs] tag

## Notes
Drift-only mode (no PRD/design files for list-tasks feature). Scanned 17 spec files total (3 business-rules + 14 conventions). Found 3 drifted files, 0 orphaned rules requiring deletion, 0 implicit new rules to add.
