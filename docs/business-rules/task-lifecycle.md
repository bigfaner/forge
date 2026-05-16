---
title: "Task Lifecycle Rules"
---

# Task Lifecycle Rules

_Source: feature/forge-cli-v3_

## State Transitions

### BIZ-task-lifecycle-001: Task State Transition Constraints

**Rule**: Tasks follow a state machine with 6 statuses: `pending`, `in_progress`, `completed`, `blocked`, `skipped`, `rejected`. Terminal states are `completed` and `rejected`. ANY transition to `completed` is blocked from the status command (must use `forge task submit`); this applies to all source states, not just `in_progress`. Transitions from non-terminal states (excluding `-> completed`) are generally allowed. `blocked` tasks can auto-restore to `pending` via `autoRestoreSourceTask` when all dependencies complete. Status `skipped` satisfies dependency checks (tasks depending on a skipped task can proceed). Status `rejected` does NOT satisfy dependency checks. Any transition from a terminal state MUST be rejected unless `--force` is used. Invalid transitions produce an AIError with code VALIDATION_ERROR.
**Context**: Prevents invalid task lifecycle progression and ensures data integrity of the task index. The state machine is intentionally permissive for non-terminal transitions to support flexible workflow recovery.
**Source**: feature/forge-cli-v3 BIZ-001

| Current State | Terminal? | Notes |
|---------------|-----------|-------|
| pending | No | Can transition to any non-terminal state; rejected allowed |
| in_progress | No | Cannot go to completed via status command; must use submit |
| completed | Yes | Blocked unless --force |
| blocked | No | Auto-restores to pending when all deps complete |
| skipped | No | Satisfies dependency checks; not terminal |
| rejected | Yes | Blocked unless --force |

### BIZ-task-lifecycle-002: Terminal State Immutability

**Rule**: Submitting a result for a task already in a terminal state (`completed` or `rejected`) via `forge task submit` is currently NOT enforced at the code level -- the submit command does not check whether the task is already terminal. The `forge task status` command enforces terminal-state guards (blocks transitions from `completed`/`rejected` unless `--force`), but `forge task submit` overwrites the status unconditionally. This is a known gap; use `--force` flag on the status command for deliberate recovery, and avoid re-submitting for already-terminal tasks.
**Context**: Guarantees that completed/rejected tasks cannot be accidentally overwritten while allowing deliberate recovery.
**Source**: feature/forge-cli-v3 BIZ-002
