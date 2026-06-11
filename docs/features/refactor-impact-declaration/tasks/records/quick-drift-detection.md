---
status: "completed"
started: "2026-05-22 22:24"
completed: "2026-05-22 22:28"
time_spent: "~4m"
---

# Task Record: T-quick-doc-drift Detect Spec Drift

## Summary
Drift detection scan of 14 project-level spec files (3 business-rules + 14 conventions). Found 3 drifted rules: BIZ-task-lifecycle-001 (6->7 statuses, terminal state changes, --force removed), BIZ-task-lifecycle-002 (submit now enforces terminal state), and forge-cli-reference (missing reopen/transition commands). Fixed all drift, committed with [auto-specs] tag. Regenerated vocabulary index.

## Changes

### Files Created
无

### Files Modified
- docs/business-rules/task-lifecycle.md
- docs/conventions/forge-cli-reference.md
- docs/.vocabulary.md

### Key Decisions
- Updated state machine description from 6 to 7 statuses (added suspended)
- Changed terminal state description to include skipped alongside completed/rejected
- Replaced --force recovery mechanism with forge task reopen/transition
- Updated error code from VALIDATION_ERROR to INVALID_TRANSITION
- Added forge task reopen and forge task transition to CLI reference

## Test Results
- **Tests Executed**: No
- **Passed**: 0
- **Failed**: 0
- **Coverage**: N/A (task has no tests)

## Acceptance Criteria
- [x] All project-level spec files scanned for drift
- [x] Drifted rules updated to match current codebase
- [x] Changes committed with [auto-specs] tag
- [x] Vocabulary index regenerated

## Notes
14 spec files scanned: 3 business-rules + 14 conventions. 11 rules current, 3 drifted, 0 orphaned.
