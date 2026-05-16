---
created: "2026-05-16"
tags: [architecture, dependencies]
---

# Blocked Tasks Never Auto-Unblock on Dependency Completion

## Problem

`forge task claim` reports "No pending tasks available" even though a task exists whose dependencies are all completed. The task is stuck in `blocked` status permanently.

Example: T-quick-6 (drift detection) depends on T-quick-5. T-quick-5 completed, but T-quick-6 remains `blocked` and is invisible to the claim system.

## Root Cause

Causal chain (3 levels deep):

1. **Symptom** → `forge task claim` says "no pending tasks" while T-quick-6 is stuck in `blocked`
2. **Direct cause** → `claimNextTask()` (claim.go:189-197) only scans tasks with `status == "pending"`. Tasks with `status == "blocked"` are completely skipped, even if dependencies are met.
3. **Root cause** → `autoRestoreSourceTask()` (submit.go:222-236) only handles fix-task → source-task restoration. There is NO general dependency cascade: when a regular task completes, no code scans blocked tasks to see if their dependencies are now satisfied and restores them to `pending`.
4. **Trigger** → Any task that is `blocked` (either set during creation or during execution) whose dependency later completes will never become claimable.

## Solution

**Immediate workaround**: Manually set the blocked task to pending:
```bash
forge task status T-quick-6 pending
```

**Proper fix**: Add a cascade unblock step to `saveIndexAndSignalCompletion()` (submit.go). After marking a task as completed, scan ALL tasks in the index: for any task with `status == "blocked"` whose dependencies are now fully met, transition it to `status == "pending"`. This mirrors `autoRestoreSourceTask()` but applies broadly instead of only for fix-task sources.

## Reusable Pattern

**When a task status blocks downstream work, the completion handler MUST cascade**. Any state machine with a "blocked" state needs an explicit unblock trigger — it cannot rely on downstream consumers to re-check. Without cascade:
- Blocked tasks accumulate and never surface
- The dispatcher thinks all work is done when it isn't
- The only recovery is manual intervention

**Design rule**: If a system has a `blocked` state, the event that satisfies the blocking condition MUST actively transition blocked items. Passive checking (only checking at claim time) fails if the item is never considered for claiming.

## Example

```go
// In saveIndexAndSignalCompletion, after saving the index:
for key, t := range idx.TasksMap() {
    if t.Status != "blocked" {
        continue
    }
    unmet := checkUnmetDeps(idx, &t)
    if len(unmet) == 0 {
        t.Status = "pending"
        idx.SetTask(key, t)
    }
}
```

## Related Files

- `forge-cli/internal/cmd/claim.go` — `claimNextTask()` only checks `pending` tasks
- `forge-cli/internal/cmd/submit.go` — `autoRestoreSourceTask()` only handles fix-task sources
- `forge-cli/internal/cmd/status.go` — `checkUnmetDeps()` helper already exists

## References

- Feature: `extract-design-md-platform-adapters` (T-quick-6 stuck in blocked)
