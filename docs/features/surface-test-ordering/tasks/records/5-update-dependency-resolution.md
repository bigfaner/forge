---
status: "completed"
started: "2026-05-26 13:22"
completed: "2026-05-26 13:40"
time_spent: "~18m"
---

# Task Record: 5 Update dependency resolution chains

## Summary
Updated resolveQuickDeps for Quick mode to use simplified dependency chain: gen-journeys -> run-tests(serial) -> verify-regression. Removed gen-contracts and gen-scripts from Quick mode task generation. Added wireQuickRunTestChain helper for Quick-specific run-test wiring where first run-test depends directly on gen-journeys.

## Changes

### Files Created
无

### Files Modified
- forge-cli/pkg/task/autogen.go
- forge-cli/pkg/task/autogen_test.go

### Key Decisions
- Quick mode skips gen-contracts and gen-scripts entirely -- first run-test depends directly on gen-journeys
- Created separate wireQuickRunTestChain function to keep Breakdown mode's wireRunTestChain unchanged
- Kept resolveQuickDeps signature compatible with existing callers (unused surfaceTypes param renamed to _)

## Test Results
- **Tests Executed**: Yes
- **Passed**: 6
- **Failed**: 0
- **Coverage**: 87.4%

## Acceptance Criteria
- [x] Breakdown mode dependency chain: T-test-gen-journeys -> T-test-gen-contracts -> T-test-gen-scripts-* -> T-test-run-{keys} -> T-test-verify-regression
- [x] Quick mode dependency chain: T-test-gen-journeys -> T-test-run-{keys} -> T-test-verify-regression
- [x] T-test-verify-regression depends only on last run-test in execution order

## Notes
Breakdown mode unchanged -- all existing Breakdown tests continue to pass. Single-surface projects degenerate correctly in both modes.
