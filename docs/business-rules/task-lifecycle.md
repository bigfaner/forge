---
title: "Task Lifecycle Rules"
---

# Task Lifecycle Rules

_Source: feature/forge-cli-v3_

## State Transitions

### BIZ-task-lifecycle-001: Task State Transition Constraints

**Rule**: Tasks follow a strict state machine: `pending` -> `in_progress` -> `completed`|`blocked`|`rejected`. `blocked` is NOT a terminal state -- it can transition back to `in_progress`. Only `completed` and `rejected` are terminal states. Any transition not explicitly allowed MUST be rejected with exit code 1 and stderr "invalid state transition: <from> -> <to>".
**Context**: Prevents invalid task lifecycle progression and ensures data integrity of the task index.
**Source**: feature/forge-cli-v3 BIZ-001

| Current State | Terminal? | Allowed Target States |
|---------------|-----------|----------------------|
| pending | No | in_progress, rejected |
| in_progress | No | completed, blocked, rejected |
| completed | Yes | (none) |
| blocked | No | in_progress |
| rejected | Yes | (none) |

### BIZ-task-lifecycle-002: Terminal State Immutability

**Rule**: Submitting a result for a task already in a terminal state (`completed` or `rejected`) MUST be rejected with exit code 1 and stderr "task already in terminal state: <status>". The index.json MUST NOT be modified.
**Context**: Guarantees that completed/rejected tasks cannot be accidentally overwritten.
**Source**: feature/forge-cli-v3 BIZ-002
