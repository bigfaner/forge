---
status: "completed"
started: "2026-05-16 20:39"
completed: "2026-05-16 20:49"
time_spent: "~10m"
---

# Task Record: 1 Fix P0: Set step-scoped SourceTaskID and cumulative counting

## Summary
Fix two P0 bugs in quality_gate.go that cause infinite fix-task loops: (A) addFixTask now sets opts.SourceTaskID to step-scoped sentinel 'quality-gate:<step>' so the counter can identify quality-gate fix tasks; (B) renamed countActiveFixTasks to countFixTasks and changed it to count ALL fix tasks per step regardless of status (cumulative counting). Vars['SOURCE_TASK_ID'] remains 'N/A (project-wide gate)' for template rendering.

## Changes

### Files Created
无

### Files Modified
- forge-cli/internal/cmd/quality_gate.go
- forge-cli/internal/cmd/quality_gate_test.go

### Key Decisions
- Sentinel 'quality-gate:<step>' is constructed programmatically from the same 'step' variable used for title prefix, eliminating typo risk
- Vars['SOURCE_TASK_ID'] intentionally diverges from opts.SourceTaskID to preserve backward-compatible template rendering
- countFixTasks counts all statuses (completed + skipped + active + blocked) to provide cumulative lifetime cap per step

## Test Results
- **Tests Executed**: Yes
- **Passed**: 15
- **Failed**: 0
- **Coverage**: 80.8%

## Acceptance Criteria
- [x] addFixTask sets opts.SourceTaskID to 'quality-gate:<step>'
- [x] Vars['SOURCE_TASK_ID'] remains 'N/A (project-wide gate)' for template rendering
- [x] countActiveFixTasks renamed to countFixTasks
- [x] countFixTasks counts ALL fix tasks per step regardless of status
- [x] After 3 cumulative fix tasks for a step, no more are created
- [x] Pending fix for step A does NOT block fix creation for step B
- [x] Existing test at quality_gate_test.go:598 updated to assert step-scoped sentinel
- [x] New tests added for step-scoped sentinel per step, cumulative counting, cross-step independence
- [x] All existing quality-gate tests pass

## Notes
无
