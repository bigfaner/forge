---
status: "completed"
started: "2026-05-26 22:14"
completed: "2026-05-26 22:30"
time_spent: "~16m"
---

# Task Record: 6 更新 forge-cli task type 命名携带 surface 信息

## Summary
Updated forge-cli task type naming to carry surface information using three-segment format ({action}.{skill}.{surface}). Added TestTypeTitle and GenSurfaceTestType helpers. Updated gen-scripts, run-test, and verify-regression task types to include surface suffix. Changed titles from generic 'Run e2e Tests' to surface-specific names like 'Run CLI Functional Tests'. Updated templates to remove 'e2e' terminology.

## Changes

### Files Created
无

### Files Modified
- forge-cli/pkg/task/types.go
- forge-cli/pkg/task/autogen.go
- forge-cli/pkg/task/types_test.go
- forge-cli/pkg/task/autogen_test.go
- forge-cli/pkg/task/data/test-run.md
- forge-cli/pkg/task/data/test-verify-regression.md

### Key Decisions
- Used GenSurfaceTestType helper to append surface segment to base types instead of creating per-surface constants
- Single-surface projects use surface type as suffix (e.g. test.run.cli); multi-surface use surface key (e.g. test.run.backend)
- verify-regression uses surface suffix only for single-surface projects; multi-surface keeps base type since it is a shared task
- autogenTemplatePath now falls back to base type template for surface-specific types (strips last segment)

## Test Results
- **Tests Executed**: Yes
- **Passed**: 142
- **Failed**: 0
- **Coverage**: 87.7%

## Acceptance Criteria
- [x] gen-scripts task type changes from test.gen-scripts to test.gen-scripts.<surfaceType>
- [x] run-test task type changes from test.run to test.run.<surfaceKey>
- [x] verify-regression task type changes from test.verify-regression to test.verify-regression.<surfaceKey>
- [x] Task titles change from 'Run e2e Tests' to surface-specific names
- [x] Task templates (test-run.md, test-verify-regression.md) use surface-specific descriptions
- [x] Go test assertions updated for new naming
- [x] go test ./pkg/task/... passes

## Notes
无
