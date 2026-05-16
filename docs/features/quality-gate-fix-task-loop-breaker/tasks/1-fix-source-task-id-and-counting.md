---
id: "1"
title: "Fix P0: Set step-scoped SourceTaskID and cumulative counting"
priority: "P0"
estimated_time: "60m"
dependencies: []
scope: "backend"
breaking: true
type: "implementation"
mainSession: false
---

# 1: Fix P0: Set step-scoped SourceTaskID and cumulative counting

## Description

Two bugs in `quality_gate.go` render the fix-task cap ineffective, causing infinite fix-task loops:

- **Bug A**: `addFixTask` sets `Vars["SOURCE_TASK_ID"]` (template variable) but never `opts.SourceTaskID` (struct field). The counter's `t.SourceTaskID != ""` filter always excludes quality-gate fix tasks, so the cap permanently reads 0.
- **Bug B**: `countActiveFixTasks` excludes completed/skipped tasks. After fix-1 completes, count resets to 0, allowing fix-2 for the same step.

Both fixes must be applied together — the counter's `SourceTaskID != ""` filter only passes after Fix A populates the field.

## Reference Files
- `docs/proposals/quality-gate-fix-task-loop-breaker/proposal.md` — Source proposal
- `forge-cli/internal/cmd/quality_gate.go` — Lines 26 (maxFixTasksPerStep), 307-319 (countActiveFixTasks), 324-401 (addFixTask)
- `forge-cli/internal/cmd/quality_gate_test.go` — Existing tests (line 598 asserts SourceTaskID == "")
- `forge-cli/pkg/task/add.go` — HasActiveFixTasks, ResolveSourceTask
- `docs/forensics/hook-feedback-loop/report.md` — Incident evidence
- `docs/forensics/fix-task-loop/report.md` — Incident evidence

## Acceptance Criteria
- [ ] `addFixTask` sets `opts.SourceTaskID` to `"quality-gate:<step>"` (e.g., `"quality-gate:compile"`, `"quality-gate:unit-test"`)
- [ ] `Vars["SOURCE_TASK_ID"]` remains `"N/A (project-wide gate)"` for template rendering (intentionally diverges from struct field)
- [ ] `countActiveFixTasks` renamed to `countFixTasks`
- [ ] `countFixTasks` counts ALL fix tasks per step regardless of status (completed + active + blocked + skipped)
- [ ] After 3 cumulative fix tasks for a step, no more are created
- [ ] Pending fix for step A does NOT block fix creation for step B (cross-step independence)
- [ ] Existing test at `quality_gate_test.go:598` updated to assert step-scoped sentinel
- [ ] New tests added for: step-scoped sentinel per step, cumulative counting, cross-step independence
- [ ] All existing quality-gate tests pass

## Hard Rules
- Sentinel `"quality-gate:<step>"` must not collide with real task IDs. `FindTask` returns nil for it, so source resolution/blocking is correctly skipped.
- The `step` variable used in sentinel construction must be the same one used for title prefix — no typo risk.

## Implementation Notes
- The sentinel is constructed programmatically: `"quality-gate:" + step` where `step` is the same variable already used in title prefix.
- `HasActiveFixTasks(index, "quality-gate:compile")` only blocks compile fix tasks, not lint or unit-test ones — step-scoping avoids cross-step dedup conflicts.
- `FindTask` returns nil for the sentinel, which means source resolution and blocking are correctly skipped for quality-gate fix tasks.
