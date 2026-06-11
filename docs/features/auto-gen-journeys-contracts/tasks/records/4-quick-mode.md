---
status: "completed"
started: "2026-05-24 00:19"
completed: "2026-05-24 00:29"
time_spent: "~10m"
---

# Task Record: 4 autogen.go Quick 模式：替换 gen-and-run 为 staged across types 拓扑

## Summary
Rewrote GetQuickTestTasks to replace gen-and-run with staged across types pipeline (gen-journeys -> gen-contracts -> gen-scripts -> run -> verify-regression), rewrote resolveQuickDeps using findTaskIndexOrPanic for all dependency resolution, updated ResolveFirstTestDep quick branch, and updated infer.go for new task IDs

## Changes

### Files Created
无

### Files Modified
- forge-cli/pkg/task/autogen.go
- forge-cli/pkg/task/autogen_test.go
- forge-cli/pkg/task/infer.go
- forge-cli/pkg/task/infer_test.go
- forge-cli/pkg/task/autoconfig_test.go

### Key Decisions
- Quick mode now uses T-test- prefix (shared with Breakdown) instead of T-quick- prefix for staged pipeline tasks
- Staged across types topology: all gen-journeys parallel -> gen-contracts -> all gen-scripts parallel -> run -> verify-regression
- All dependency resolution uses findTaskIndexOrPanic (no arithmetic indices)
- TypeTestGenAndRun type definition preserved for backward compatibility
- T-quick-doc-drift and T-quick-verify-regression IDs preserved for non-e2e tasks

## Test Results
- **Tests Executed**: Yes
- **Passed**: 16
- **Failed**: 0
- **Coverage**: 90.8%

## Acceptance Criteria
- [x] GetQuickTestTasks no longer generates gen-and-run tasks (TypeTestGenAndRun not in output)
- [x] Per interface type generates T-test-gen-journeys-{type} tasks
- [x] Generates a T-test-gen-contracts task
- [x] Per interface type generates T-test-gen-scripts-{type} tasks
- [x] Generates T-test-run and T-test-verify-regression
- [x] resolveQuickDeps implements staged across types topology (5 stages)
- [x] All dependency lookup uses findTaskIndex/findTaskIndexOrPanic (no arithmetic indices)
- [x] findTaskIndex returns -1 with panic on missing task
- [x] All existing Quick mode tests pass

## Notes
gen-journeys template body contains AUTO_COMMIT=true conditional instruction (from Task 1). gen-contracts template body contains SKIP_EVAL_GATE=true conditional instruction (from Task 1). infer.go updated to recognize T-test-gen-journeys-{type} and T-test-gen-contracts.
