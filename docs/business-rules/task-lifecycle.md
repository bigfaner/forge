---
title: "Task Lifecycle Rules"
---

# Task Lifecycle Rules

_Source: feature/forge-cli-v3_

## State Transitions

### BIZ-task-lifecycle-001: Task State Transition Constraints

**Rule**: Tasks follow a state machine with 6 statuses: `pending`, `in_progress`, `completed`, `blocked`, `skipped`, `rejected`. Terminal states are `completed` and `rejected`. The `in_progress -> completed` transition is blocked from the status command (must use `forge task submit`). Transitions from non-terminal states are generally allowed (the code returns true for any non-blocked transition). `blocked` tasks can auto-restore to `pending` via `autoRestoreSourceTask` when all dependencies complete. Status `skipped` satisfies dependency checks (tasks depending on a skipped task can proceed). Any transition from a terminal state MUST be rejected unless `--force` is used. Invalid transitions produce an AIError with code VALIDATION_ERROR.
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

**Rule**: Submitting a result for a task already in a terminal state (`completed` or `rejected`) MUST be rejected with an AIError (code VALIDATION_ERROR, message "Invalid transition: <from> -> <to>"). The index.json MUST NOT be modified. Override requires `--force` flag.
**Context**: Guarantees that completed/rejected tasks cannot be accidentally overwritten while allowing deliberate recovery.
**Source**: feature/forge-cli-v3 BIZ-002
