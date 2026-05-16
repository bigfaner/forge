---
feature: "task-lifecycle-hardening"
sources:
  - docs/proposals/task-lifecycle-hardening/proposal.md
  - docs/features/task-lifecycle-hardening/tasks/1-fix-check-deps-self-block.md
  - docs/features/task-lifecycle-hardening/tasks/2-lazy-unblock-scan.md
  - docs/features/task-lifecycle-hardening/tasks/3-update-tests.md
generated: "2026-05-16"
---

# Test Cases: task-lifecycle-hardening

## Summary

| Type | Count |
|------|-------|
| UI   | 0   |
| TUI  | 0  |
| Mobile | 0 |
| API  | 0  |
| CLI  | 14  |
| **Total** | **14** |

---

## CLI Test Cases

### TC-001: Active fix-task with SourceTaskID == selfID blocks claim

- **Source**: Task 1 / AC-1 ("checkDependenciesMet returns false when an active fix-task has SourceTaskID == selfID")
- **Type**: CLI
- **Target**: cli/claim
- **Test ID**: cli/claim/active-fix-task-with-source-task-id-eq-self-blocks-claim
- **Pre-conditions**: A TaskIndex containing: (1) a pending task "3" with no dependencies, (2) a pending fix-task "fix-1" with SourceTaskID="3" and Type="fix"
- **Steps**:
  1. Call `checkDependenciesMet(index, "3", task3)` where task3 has no dependencies
- **Expected**: Returns `(false, [non-empty unmet])` -- task 3 is blocked by the fix-task targeting itself
- **Priority**: P0

### TC-002: In-progress fix-task with SourceTaskID == selfID blocks claim

- **Source**: Task 1 / AC-1 (extended -- active includes in_progress status)
- **Type**: CLI
- **Target**: cli/claim
- **Test ID**: cli/claim/in-progress-fix-task-with-source-task-id-eq-self-blocks-claim
- **Pre-conditions**: A TaskIndex containing: (1) a pending task "3" with no dependencies, (2) an in_progress fix-task "fix-1" with SourceTaskID="3" and Type="fix"
- **Steps**:
  1. Call `checkDependenciesMet(index, "3", task3)` where task3 has no dependencies
- **Expected**: Returns `(false, [non-empty unmet])` -- in_progress fix-task targeting self also blocks
- **Priority**: P0

### TC-003: Completed fix-task targeting self does not block

- **Source**: Task 1 / AC-2 ("checkDependenciesMet returns true when fix-task targeting self is completed")
- **Type**: CLI
- **Target**: cli/claim
- **Test ID**: cli/claim/completed-fix-task-targeting-self-does-not-block
- **Pre-conditions**: A TaskIndex containing: (1) a pending task "3" with no dependencies, (2) a completed fix-task "fix-1" with SourceTaskID="3" and Type="fix"
- **Steps**:
  1. Call `checkDependenciesMet(index, "3", task3)` where task3 has no dependencies
- **Expected**: Returns `(true, [])` -- completed fix-task does not block
- **Priority**: P0

### TC-004: Self-block takes precedence over met regular dependencies

- **Source**: Task 1 / AC-1 + Proposal flowchart (node A4)
- **Type**: CLI
- **Target**: cli/claim
- **Test ID**: cli/claim/self-block-takes-precedence-over-met-regular-dependencies
- **Pre-conditions**: A TaskIndex containing: (1) a completed task "2", (2) a pending task "3" with Dependencies=["2"], (3) a pending fix-task "fix-1" with SourceTaskID="3" and Type="fix"
- **Steps**:
  1. Call `checkDependenciesMet(index, "3", task3)` where task3 depends on completed task "2"
- **Expected**: Returns `(false, [non-empty unmet])` -- despite regular deps being met, the self-block from fix-1 keeps task 3 blocked
- **Priority**: P0

### TC-005: Fix-task targeting other task does not cause self-block

- **Source**: Task 1 / AC-3 ("Existing behavior unchanged for tasks without active fix-tasks targeting them")
- **Type**: CLI
- **Target**: cli/claim
- **Test ID**: cli/claim/fix-task-targeting-other-task-does-not-cause-self-block
- **Pre-conditions**: A TaskIndex containing: (1) a completed task "2", (2) a pending task "3" with no dependencies, (3) a pending fix-task "fix-1" with SourceTaskID="2" (targeting task 2, not task 3) and Type="fix"
- **Steps**:
  1. Call `checkDependenciesMet(index, "3", task3)` where task3 has no dependencies
- **Expected**: Returns `(true, [])` -- fix-task targeting another task does not block task 3
- **Priority**: P0

### TC-006: Multiple fix-tasks targeting self must all complete

- **Source**: Task 1 / AC-1 (extended -- multiple fix-tasks)
- **Type**: CLI
- **Target**: cli/claim
- **Test ID**: cli/claim/multiple-fix-tasks-targeting-self-must-all-complete
- **Pre-conditions**: A TaskIndex containing: (1) a pending task "3" with no dependencies, (2) a completed fix-task "fix-1" with SourceTaskID="3", (3) a pending fix-task "fix-2" with SourceTaskID="3"
- **Steps**:
  1. Call `checkDependenciesMet(index, "3", task3)` where task3 has no dependencies
- **Expected**: Returns `(false, [non-empty unmet])` -- task 3 stays blocked while fix-2 is still pending
- **Priority**: P1

### TC-007: Blocked task auto-unblocked when dependencies met

- **Source**: Task 2 / AC-4 + Proposal SC-1 ("Blocked task auto-transitions to pending when checkDependenciesMet returns true")
- **Type**: CLI
- **Target**: cli/claim
- **Test ID**: cli/claim/blocked-task-auto-unblocked-when-dependencies-met
- **Pre-conditions**: A TaskIndex containing: (1) a completed task "1" with no dependencies, (2) a blocked task "2" with Dependencies=["1"]
- **Steps**:
  1. Call `claimNextTask(index)`
- **Expected**: Returns key="task2", task.Status="in_progress". After the call, task 2 was auto-unblocked from blocked to pending, then claimed as in_progress
- **Priority**: P0

### TC-008: Blocked task stays blocked when dependencies not met

- **Source**: Task 2 / AC-4 (negative case)
- **Type**: CLI
- **Target**: cli/claim
- **Test ID**: cli/claim/blocked-task-stays-blocked-when-dependencies-not-met
- **Pre-conditions**: A TaskIndex containing: (1) a pending task "1" with no dependencies, (2) a blocked task "2" with Dependencies=["1"]
- **Steps**:
  1. Call `claimNextTask(index)`
- **Expected**: Returns key="task1". Task "2" remains in status "blocked" because its dependency "1" is still pending
- **Priority**: P0

### TC-009: Auto-unblock logged to stdout

- **Source**: Task 2 / AC-5 ("Auto-unblocked tasks are logged")
- **Type**: CLI
- **Target**: cli/claim
- **Test ID**: cli/claim/auto-unblock-logged-to-stdout
- **Pre-conditions**: A TaskIndex containing: (1) a completed task "1", (2) a blocked task "2" with Dependencies=["1"]
- **Steps**:
  1. Call `claimNextTask(index)` while capturing stdout
- **Expected**: Stdout contains the string "Auto-unblocked task 2"
- **Priority**: P1

### TC-010: Multiple blocked tasks unblocked simultaneously

- **Source**: Task 2 / AC-4 (extended -- multiple tasks) + AC-6 ("scan runs before hasPending check")
- **Type**: CLI
- **Target**: cli/claim
- **Test ID**: cli/claim/multiple-blocked-tasks-unblocked-simultaneously
- **Pre-conditions**: A TaskIndex containing: (1) a completed task "1", (2) a blocked task "2" with Dependencies=["1"] and Priority="P1", (3) a blocked task "3" with Dependencies=["1"] and Priority="P0"
- **Steps**:
  1. Call `claimNextTask(index)`
- **Expected**: Returns key="task3" (P0 beats P1). Task "2" is auto-unblocked to pending. Task "3" is claimed as in_progress
- **Priority**: P0

### TC-011: Blocked task with active fix targeting it stays blocked

- **Source**: Task 2 / AC-7 + Proposal SC-2 ("Fix-task in progress targeting the task keeps it blocked")
- **Type**: CLI
- **Target**: cli/claim
- **Test ID**: cli/claim/blocked-task-with-active-fix-targeting-it-stays-blocked
- **Pre-conditions**: A TaskIndex containing: (1) a completed task "1", (2) a blocked task "2" with Dependencies=["1"], (3) a pending fix-task "fix-1" with SourceTaskID="2" and Type="fix"
- **Steps**:
  1. Call `claimNextTask(index)`
- **Expected**: Returns key="fix-1" (the fix-task is claimed). Task "2" remains blocked because fix-1 (targeting it) is still active (though now in_progress, the lazy scan ran before the claim)
- **Priority**: P0

### TC-012: Fix completed auto-unblocks blocked source task (block-source lifecycle)

- **Source**: Proposal SC-3 ("--block-source scenario: fix done -> claim -> source auto-unblocked") + Task 2 / AC-4
- **Type**: CLI
- **Target**: cli/claim
- **Test ID**: cli/claim/fix-completed-auto-unblocks-blocked-source-task
- **Pre-conditions**: A TaskIndex containing: (1) a blocked source task "1" with no dependencies, (2) a completed fix-task "fix-1" with SourceTaskID="1" and Type="fix"
- **Steps**:
  1. Call `claimNextTask(index)`
- **Expected**: Returns key="source", task.Status="in_progress". The source task is auto-unblocked from blocked to pending (fix-task is completed, no self-block), then claimed
- **Priority**: P0

### TC-013: Source stays blocked when fix is still in-progress

- **Source**: Proposal SC-3 (negative case -- fix still active)
- **Type**: CLI
- **Target**: cli/claim
- **Test ID**: cli/claim/source-stays-blocked-when-fix-is-still-in-progress
- **Pre-conditions**: A TaskIndex containing: (1) a blocked source task "1" with no dependencies, (2) an in_progress fix-task "fix-1" with SourceTaskID="1" and Type="fix"
- **Steps**:
  1. Call `claimNextTask(index)`
- **Expected**: Returns error (no eligible tasks). Source task "1" remains in status "blocked"
- **Priority**: P0

### TC-014: Auto-downgraded task auto-unblocked when dep completes

- **Source**: Proposal SC-4 ("Auto-downgrade scenario: task blocked -> dep completed -> claim auto-unblocks")
- **Type**: CLI
- **Target**: cli/claim
- **Test ID**: cli/claim/auto-downgraded-task-auto-unblocked-when-dep-completes
- **Pre-conditions**: A TaskIndex containing: (1) a completed task "1" with no dependencies, (2) a blocked task "2" with Dependencies=["1"] and BlockedReason="auto-downgrade: testsFailed=2"
- **Steps**:
  1. Call `claimNextTask(index)`
- **Expected**: Returns key="task2", task.Status="in_progress". The auto-downgraded task is auto-unblocked because its dependency is met and no fix-task targets it
- **Priority**: P0

---

## Traceability

| TC ID | Source | Type | Target | Priority |
|-------|--------|------|--------|----------|
| TC-001 | Task 1 / AC-1 | CLI | cli/claim | P0 |
| TC-002 | Task 1 / AC-1 (in_progress) | CLI | cli/claim | P0 |
| TC-003 | Task 1 / AC-2 | CLI | cli/claim | P0 |
| TC-004 | Task 1 / AC-1 + Proposal A4 | CLI | cli/claim | P0 |
| TC-005 | Task 1 / AC-3 | CLI | cli/claim | P0 |
| TC-006 | Task 1 / AC-1 (multiple) | CLI | cli/claim | P1 |
| TC-007 | Task 2 / AC-4 + Proposal SC-1 | CLI | cli/claim | P0 |
| TC-008 | Task 2 / AC-4 (negative) | CLI | cli/claim | P0 |
| TC-009 | Task 2 / AC-5 | CLI | cli/claim | P1 |
| TC-010 | Task 2 / AC-4 + AC-6 | CLI | cli/claim | P0 |
| TC-011 | Task 2 / AC-7 + Proposal SC-2 | CLI | cli/claim | P0 |
| TC-012 | Proposal SC-3 + Task 2 / AC-4 | CLI | cli/claim | P0 |
| TC-013 | Proposal SC-3 (negative) | CLI | cli/claim | P0 |
| TC-014 | Proposal SC-4 | CLI | cli/claim | P0 |
