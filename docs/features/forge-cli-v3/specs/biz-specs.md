---
feature: "forge-cli-v3"
generated: "2026-05-14"
status: draft
---

# Business Rules: Forge CLI v3

## Task Lifecycle

### BIZ-001: Task State Transition Constraints

**Rule**: Tasks follow a strict state machine: `pending` -> `in_progress` -> `completed`|`blocked`|`rejected`. `blocked` is NOT a terminal state -- it can transition back to `in_progress`. Only `completed` and `rejected` are terminal states. Any transition not explicitly allowed MUST be rejected with exit code 1 and stderr "invalid state transition: <from> -> <to>".
**Context**: Prevents invalid task lifecycle progression and ensures data integrity of the task index.
**Scope**: [CROSS]
**Source**: prd-spec.md > Error Handling > Task State Transitions

| Current State | Terminal? | Allowed Target States |
|---------------|-----------|----------------------|
| pending | No | in_progress, rejected |
| in_progress | No | completed, blocked, rejected |
| completed | Yes | (none) |
| blocked | No | in_progress |
| rejected | Yes | (none) |

### BIZ-002: Terminal State Immutability

**Rule**: Submitting a result for a task already in a terminal state (`completed` or `rejected`) MUST be rejected with exit code 1 and stderr "task already in terminal state: <status>". The index.json MUST NOT be modified.
**Context**: Guarantees that completed/rejected tasks cannot be accidentally overwritten.
**Scope**: [CROSS]
**Source**: prd-spec.md > Error Handling > Command-level failures

### BIZ-003: Concurrent Write Conflict Handling

**Rule**: When two agents simultaneously submit results for the same task, exactly one MUST succeed (exit 0) and the other MUST fail with exit code 1 and stderr "concurrent write conflict, retry". The index.json MUST remain valid JSON parseable by `jq .` after the conflict.
**Context**: File lock mechanism on index.json.lock ensures data integrity under concurrent access.
**Scope**: [LOCAL]
**Source**: prd-spec.md > Error Handling > Command-level failures (forge task submit)

## Quality Gate

### BIZ-004: Quality Gate Sequential Pipeline

**Rule**: `forge quality-gate` executes a sequential pipeline: compile -> fmt -> lint -> test. The first failing step terminates the pipeline with exit code 1. All steps passing yields exit code 0.
**Context**: Provides a single-command CI check that enforces code quality in order of dependency (code must compile before it can be linted).
**Scope**: [CROSS]
**Source**: prd-spec.md > Error Handling > CI/Hook Flow

### BIZ-005: Fix-Task Creation on Quality Gate Failure

**Rule**: When a quality gate step fails, a P0 fix-task is automatically created. If the same step fails again, a new fix-task is created with an incremented sequence number (e.g., "fix-compile-1", "fix-compile-2"). Maximum 3 fix-tasks per step -- beyond that, stderr outputs "max fix-tasks reached for <step>, manual intervention required" with exit code 1, and no new fix-task is created.
**Context**: Prevents unbounded fix-task proliferation while giving agents multiple auto-retry opportunities.
**Scope**: [LOCAL]
**Source**: prd-spec.md > User Stories > Story 4

## Command Discovery

### BIZ-006: Help Entry Count Threshold

**Rule**: `forge --help` MUST display at most 10 top-level entries (5 command groups + 5 top-level commands). This threshold is derived from the `gh` CLI usability benchmark.
**Context**: Ensures AI agents can quickly discover the correct command without information overload.
**Scope**: [LOCAL]
**Source**: prd-spec.md > Goals > Command discoverability

### BIZ-007: Subcommand Description Format

**Rule**: Each subcommand description in `--help` output MUST follow "command name + verb + object" structure (e.g., "submit task execution result") and be <= 80 characters.
**Context**: Self-documenting help output reduces the need for external documentation lookup.
**Scope**: [LOCAL]
**Source**: prd-spec.md > User Stories > Story 1

## Error Reporting

### BIZ-008: Consistent Exit Code Semantics

**Rule**: Exit code 0 = success (or intentional no-op, e.g., "no tasks to clean up"). Exit code 1 = failure with descriptive stderr message. Exit code 2 = reserved for usage errors (Cobra default).
**Context**: AI agents rely on exit codes to determine next action; consistent semantics prevent misinterpretation.
**Scope**: [CROSS]
**Source**: prd-spec.md > Error Handling > Command-level failures (all commands)

### BIZ-009: Actionable Error Messages

**Rule**: Every error message on stderr MUST contain: (1) the specific failure reason, and (2) a hint for recovery when applicable. Example: "unknown profile: <value>" MUST be followed by listing all supported profiles.
**Context**: AI agents need self-correcting feedback loops without human intervention.
**Scope**: [CROSS]
**Source**: prd-spec.md > Error Handling > Command-level failures (pattern across all commands)
