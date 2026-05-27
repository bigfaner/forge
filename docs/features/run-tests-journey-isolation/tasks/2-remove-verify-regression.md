---
id: "2"
title: "Remove verify-regression from auto-generated task pipeline"
priority: "P1"
estimated_time: "1-2h"
dependencies: []
surface-key: ""
surface-type: ""
breaking: true
type: "coding.enhancement"
mainSession: false
---

# 2: Remove verify-regression from auto-generated task pipeline

## Description

Remove the `T-test-verify-regression` task from both `GetBreakdownTestTasks()` and `GetQuickTestTasks()` in `autogen.go`. The quality-gate (Stop hook) already runs `just test` for full regression after all tasks complete, making verify-regression redundant in the pipeline.

Changes required:
1. Remove verify-regression task definition from `GetBreakdownTestTasks()` (lines ~219-228)
2. Remove verify-regression task definition from `GetQuickTestTasks()` (lines ~334-343)
3. Update dependency resolution: tasks that depended on verify-regression now depend on the last run-test in the chain
4. Update `resolveBreakdownDeps()`: validation/specs/clean tasks that depended on `T-test-verify-regression` should depend on last run-test instead
5. Update `resolveQuickDeps()`: validation/drift/clean tasks that depended on `T-test-verify-regression` should depend on last run-test instead
6. Update all related tests in `autogen_test.go`

## Reference Files
- `proposal.md#Proposed-Solution` — P2: remove verify-regression, quality-gate covers full regression
- `proposal.md#Key-Risks` — risk of missing pipeline-level regression protection (mitigated by quality-gate)
- `proposal.md#Assumptions-Challenged` — verify-regression overturned by XY Detection

## Acceptance Criteria

- [ ] `GetBreakdownTestTasks()` no longer generates a `T-test-verify-regression` task
- [ ] `GetQuickTestTasks()` no longer generates a `T-test-verify-regression` task
- [ ] Dependency chain is correct: downstream tasks (validation, specs, clean) depend on last run-test, not verify-regression
- [ ] `TypeTestVerifyRegression` constant and its `ValidTypes` entry remain (used elsewhere, not removed)
- [ ] All tests in `forge-cli/pkg/task/autogen_test.go` pass (`TestGetQuickTestTasks_*` and `TestGetBreakdownTestTasks_*`)
- [ ] Version bump in `scripts/version.txt` (patch: bug fix / dead code removal)

## Hard Rules

- Do NOT remove `TypeTestVerifyRegression` constant or its `ValidTypes` registration — it may be used in other contexts
- Only remove the task definition and dependency wiring in `GetBreakdownTestTasks()` and `GetQuickTestTasks()`
- Downstream task dependencies (validation, specs, drift, clean) must be rewired to the last run-test task in the chain, not to a phantom verify-regression

## Implementation Notes

- In `GetBreakdownTestTasks()`, the verify-regression block is lines ~219-228. Remove the block. Then in `resolveBreakdownDeps()`, the variable `verifyIdx` and `tasks[verifyIdx].Dependencies` assignment must be removed. Tasks that previously depended on `T-test-verify-regression` (validation, specs) should now depend on `lastRunID` from `wireRunTestChain()`.
- In `GetQuickTestTasks()`, same pattern: remove the verify-regression block (~lines 334-343). In `resolveQuickDeps()`, remove `verifyIdx` and rewire downstream tasks to `lastRunID`.
- The `findTaskIndexOrPanic` calls for `T-test-verify-regression` in dependency resolution functions must be removed, otherwise they will panic at runtime.
- Check if `forge-cli/pkg/prompt/data/test-verify-regression.md` should be removed or kept for manual use.
