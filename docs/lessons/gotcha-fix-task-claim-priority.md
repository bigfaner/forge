---
created: "2026-05-15"
tags: [dependencies, architecture]
---

# Fix Tasks Lose Claim Priority to Numeric Business Tasks

## Problem

After creating a fix task (`forge task add --template fix-task`) to address compile errors from task 3, the dispatcher expected the fix task to be claimed next. Instead, `forge task claim` returned business task 4. Both tasks were P0 with met dependencies.

## Root Cause

Causal chain:

1. **Symptom**: `forge task claim` returned task 4 instead of fix-1
2. **Direct cause**: `compareVersionIDs()` in `claim.go:271-289` sorts numeric IDs before alphabetic IDs
3. **Root cause**: Fix tasks get alphabetic IDs (`"fix-1"`) while business tasks get numeric IDs (`"4"`). The claim logic sorts by priority first (P0 = both), then by version ID as tiebreaker. The version comparison treats numeric segments as sorting before alphabetic segments (`return aIsNum`), so `"4"` (numeric) always sorts before `"fix-1"` (alphabetic)
4. **Trigger condition**: A fix task and an unblocked business task coexist at the same priority level

Key code in `compareVersionIDs`:
```go
if aIsNum != bIsNum {
    return aIsNum // numeric sorts before alphabetic
}
```

This means `"4"` < `"fix-1"` in the ordering, so business task 4 is claimed first.

## Solution

Two options (not yet implemented):

1. **Fix task ID format**: Change fix task IDs to sort before business tasks (e.g., `"0-fix-1"` instead of `"fix-1"`)
2. **Dispatcher awareness**: The dispatcher (run-tasks) should track pending fix tasks and explicitly claim them by ID after adding, rather than relying on the generic claim ordering

## Reusable Pattern

**When `forge task add --template fix-task` creates a fix task, do NOT assume it will be claimed next.** The CLI's claim ordering treats alphabetic IDs as lower priority than numeric IDs in the tiebreaker. If the fix task must run before other P0 business tasks, either:

- Verify the fix task was actually claimed (check `ACTION: CLAIMED` output)
- Or claim it explicitly: `forge task claim` → if wrong task returned, the dispatcher needs logic to prioritize fix tasks

**Warning sign**: Any `forge task add` output with an alphabetic ID (fix-N, T-test-*, T-quick-*) competing against numeric business task IDs at the same priority level.

## Related Files

- `forge-cli/internal/cmd/claim.go` — `claimNextTask()`, `compareVersionIDs()`
- `forge-cli/pkg/task/add.go` — fix task ID generation
- `plugins/forge/commands/run-tasks.md` — dispatcher protocol (assumes fix tasks are claimed first)
