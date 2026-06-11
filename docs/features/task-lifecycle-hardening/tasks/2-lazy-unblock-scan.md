---
id: "2"
title: "Add lazy unblock scan in claimNextTask"
priority: "P0"
estimated_time: "1h"
dependencies: ["1"]
scope: "backend"
breaking: false
type: "implementation"
mainSession: false
---

# 2: Add lazy unblock scan in claimNextTask

## Description

Add a lazy unblock scan at the start of `claimNextTask`: iterate all `blocked` tasks, check if their dependencies are now met (via the fixed `checkDependenciesMet` from task 1), and auto-transition eligible ones to `pending`. This replaces the fragile "submit-time cascade" pattern with DAG-scheduler-standard "claim-time pull" resolution.

## Reference Files
- `docs/proposals/task-lifecycle-hardening/proposal.md` — Source proposal
- `forge-cli/internal/cmd/claim.go` — Target file (claimNextTask at ~line 181)

## Acceptance Criteria

- [ ] Blocked task auto-transitions to pending when `checkDependenciesMet` returns true on claim
- [ ] Auto-unblocked tasks are logged (e.g., `fmt.Printf("Auto-unblocked task %s", id)`)
- [ ] The scan runs before the existing `hasPending` check so newly-unblocked tasks are visible
- [ ] Fix-task in progress targeting the task keeps it blocked (verified via task 1's fix)
- [ ] Existing claim behavior unchanged when no blocked tasks exist

## Hard Rules

- No changes to submit.go
- `checkDependenciesMet` must already include the selfID fix (dependency on task 1)

## Implementation Notes

Insert the lazy scan between `index.TasksMap()` availability and the `hasPending` check. The scan should:

1. Iterate `index.TasksMap()`, filter `Status == "blocked"`
2. For each blocked task, call `checkDependenciesMet(index, task.ID, task)`
3. If met, set `task.Status = "pending"` and `index.SetTask(key, updatedTask)`
4. Log the auto-unblock

The scan is O(n) where n = tasks in the feature, typically <20 — no performance concern.
