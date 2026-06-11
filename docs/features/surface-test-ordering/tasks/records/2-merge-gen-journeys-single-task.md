---
status: "completed"
started: "2026-05-26 10:17"
completed: "2026-05-26 10:34"
time_spent: "~17m"
---

# Task Record: 2 Merge gen-journeys to single task

## Summary
Merged per-surface gen-journeys tasks into a single T-test-gen-journeys task in both GetBreakdownTestTasks and GetQuickTestTasks. Updated dependency resolution (resolveBreakdownDeps, resolveQuickDeps) to reference the single task. Updated all affected tests.

## Changes

### Files Created
无

### Files Modified
- forge-cli/pkg/task/autogen.go
- forge-cli/pkg/task/autogen_test.go
- forge-cli/pkg/task/autoconfig_test.go
- forge-cli/pkg/task/build_test.go

### Key Decisions
- Single gen-journeys task uses empty SurfaceType field, allowing renderBody to omit {{TEST_TYPE}} line automatically
- Key changed from per-surface 'gen-journeys-{type}' to single 'gen-journeys'
- ID changed from per-surface 'T-test-gen-journeys-{type}' to single 'T-test-gen-journeys'

## Test Results
- **Tests Executed**: Yes
- **Passed**: 88
- **Failed**: 0
- **Coverage**: 87.3%

## Acceptance Criteria
- [x] gen-journeys generates single T-test-gen-journeys task covering all configured surfaces
- [x] Single surface project degenerates to unsuffixed T-test-gen-journeys, behavior consistent with before
- [x] renderBody correctly handles empty TestType field

## Notes
Both functions (GetBreakdownTestTasks, GetQuickTestTasks) modified as required by Hard Rules. renderBody already handled empty SurfaceType — no changes needed there.
