---
created: "2026-05-16"
tags: [testing, architecture, error-handling]
---

# Quality Gate Fix-Task Loop: Cap Bypass Due to Missing SourceTaskID

## Problem

The `forge quality-gate` stop hook runs `just test` at session end. When tests fail, it auto-creates a P0 fix task via `addFixTask`. After the fix task is resolved, the session stops again, the hook fires again, the same flaky tests fail, another fix task is created — infinite loop.

## Root Cause

Causal chain (4 levels):

1. **Symptom**: fix-1, fix-2, fix-3... created in sequence for the same transient test failure. The cap (`maxFixTasksPerStep = 3`) never triggers.
2. **Direct cause**: `countActiveFixTasks()` always returns 0 for quality-gate fix tasks, so `addFixTask` never hits the cap.
3. **Root cause**: Two bugs compound:
   - **Bug A**: `addFixTask` puts source info in `Vars["SOURCE_TASK_ID"]` (template rendering) but never sets `opts.SourceTaskID` (struct field). The resulting Task has `SourceTaskID: ""`. But `countActiveFixTasks` requires `t.SourceTaskID != ""` as its first filter — so quality-gate fix tasks are invisible to the counter.
   - **Bug B**: Even if Bug A were fixed, `countActiveFixTasks` only counts non-terminal tasks (`Status != "completed"`). After fix-1 completes, the count resets to 0, allowing fix-2 for the same step.
4. **Trigger condition**: Any flaky test that fails in the full suite but passes in isolation, combined with the Stop hook re-triggering on every session end.

```
Stop hook → quality-gate → just test → flaky test fails
  → addFixTask → countActiveFixTasks returns 0 (SourceTaskID empty)
  → creates fix-1 (cap bypassed)
  → fix-1 resolved → Stop hook fires again
  → same flaky test fails → countActiveFixTasks returns 0 (fix-1 completed + SourceTaskID empty)
  → creates fix-2 (cap bypassed again)
  → infinite loop
```

## Solution

Three fixes needed:

1. **Set `opts.SourceTaskID`**: In `addFixTask`, add `SourceTaskID: "quality-gate"` so quality-gate fix tasks are identifiable by the counter.
2. **Count all fix tasks per step**: Change `countActiveFixTasks` to count ALL fix tasks for the step (including completed), not just active ones. Or track cumulative count per step in a separate mechanism.
3. **Retry before creating fix task**: When `just test` fails, re-run once before creating a fix task. If it passes on retry, emit a warning instead.

## Reusable Pattern

**When implementing rate-limiting or dedup counters, verify the identifier field is actually populated.** A counter that filters on `field != ""` will silently skip every record where the field wasn't set — no error, no warning, just an ineffective cap.

For auto-generated fix/task patterns:
- The struct field and the template variable are two different things. Setting the template var does NOT set the struct field.
- Caps must count across the full lifecycle (completed + active), not just active records.
- Always retry once before auto-creating fix tasks for test failures — transient failures are the norm, not the exception.

## Example

```go
// BAD: SourceTaskID only in Vars (template), not in struct
opts := task.AddTaskOpts{
    Vars: map[string]string{
        "SOURCE_TASK_ID": "N/A (project-wide gate)", // template only
    },
    // SourceTaskID not set → Task.SourceTaskID == ""
}

// countActiveFixTasks checks:
t.SourceTaskID != "" // always false → cap never triggers

// GOOD: Set the struct field
opts := task.AddTaskOpts{
    SourceTaskID: "quality-gate", // struct field → Task.SourceTaskID == "quality-gate"
    Vars: map[string]string{
        "SOURCE_TASK_ID": "N/A (project-wide gate)", // template rendering
    },
}
```

## Related Files

- `forge-cli/internal/cmd/quality_gate.go:304-401` — countActiveFixTasks + addFixTask
- `forge-cli/pkg/task/add.go:27,152,201` — AddTaskOpts.SourceTaskID handling
- `forge-cli/pkg/task/types.go:88-90` — Task.SourceTaskID field
- `docs/forensics/fix-task-loop/report.md` — Full forensic analysis

## References

- Forensic report: `docs/forensics/fix-task-loop/report.md`
