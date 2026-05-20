---
status: "completed"
started: "2026-05-20 17:45"
completed: "2026-05-20 17:55"
time_spent: "~10m"
---

# Task Record: T-quick-doc-drift Detect Spec Drift

## Summary
Drift detection: scanned all 15 project-level spec files (3 business-rules + 12 conventions). Found 1 drifted rule (TEST-isolation-004) — spec described mandatory dedicated binary compilation for all e2e tests, but shared helpers (helpers_test.go, testkit/helpers.go) still use exec.Command("forge", ...). Updated rule from MUST to SHOULD with current state documentation and migration note.

## Changes

### Files Created
无

### Files Modified
- docs/conventions/testing-isolation.md

### Key Decisions
- Updated TEST-isolation-004 from MUST to SHOULD since the codebase is in partial adoption state — some e2e tests compile dedicated binaries but shared helpers have not been migrated

## Test Results
- **Tests Executed**: No
- **Passed**: 0
- **Failed**: 0
- **Coverage**: N/A (task has no tests)

## Acceptance Criteria
- [x] All project-level spec files scanned for drift
- [x] Drifted rules updated to match current codebase
- [x] Changes committed with [auto-specs] tag

## Notes
14 of 15 rules validated as current. Only TEST-isolation-004 was drifted. No orphaned rules found. No implicit new rules discovered.
