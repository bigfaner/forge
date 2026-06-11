---
status: "completed"
started: "2026-05-26 10:03"
completed: "2026-05-26 10:15"
time_spent: "~12m"
---

# Task Record: 1 Add execution-order config with surface-key validation

## Summary
Added ExecutionOrder config field, surface-key normalization, validation (key format, execution-order references, same-type conflict detection), and default priority resolution (api > web > cli > tui > mobile). All validations run at config load time (fail fast).

## Changes

### Files Created
- forge-cli/pkg/forgeconfig/execution_order.go
- forge-cli/pkg/forgeconfig/execution_order_test.go

### Files Modified
- forge-cli/pkg/forgeconfig/config.go
- forge-cli/scripts/version.txt

### Key Decisions
- NormalizeSurfaceKey reuses same normalization logic as normalizeSurfaceKeyValue (DRY)
- Default priority uses sort.SliceStable to preserve YAML map order as tiebreaker for non-default types
- Scalar form (key '.') and single-surface configs are exempt from ordering requirements

## Test Results
- **Tests Executed**: Yes
- **Passed**: 32
- **Failed**: 0
- **Coverage**: 88.7%

## Acceptance Criteria
- [x] execution-order references non-existent surface-key errors at config load time
- [x] Same-type conflict (2 api surfaces) errors at config load time, hints execution-order
- [x] Surface-key validation: ADMIN PANEL normalizes to admin-panel; 123bad errors
- [x] Default priority: api > web > cli > mobile without execution-order

## Notes
Bumped version from 5.8.0 to 5.9.0 (minor: new feature). Pre-existing testrunner failure unrelated.
