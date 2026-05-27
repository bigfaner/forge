---
title: "Task Lifecycle Rules"
domains: [state-machine, transition, terminal, blocked, completed, pending, suspended, skipped, reopen, type-validation, system-types]
---

# Task Lifecycle Rules

_Source: feature/forge-cli-v3_

## State Transitions

### BIZ-task-lifecycle-001: Task State Transition Constraints

**Rule**: Tasks follow a state machine with 7 statuses: `pending`, `in_progress`, `completed`, `blocked`, `suspended`, `skipped`, `rejected`. Terminal states are `completed`, `rejected`, and `skipped`. Transition to `completed` is only allowed via `forge task submit` (role-based enforcement via `ValidateTransition`). Terminal states are immutable: `completed` can never be transitioned; `rejected` and `skipped` can only go to `pending` via `forge task reopen`. Manual operator overrides (e.g., unblocking, skipping, rejecting) use `forge task transition` (requires `--reason` flag). `blocked` tasks can auto-restore to `pending` via `autoRestoreSourceTask` when all dependencies complete. `suspended` is a manual-hold state (entered/exited only via `forge task transition`); it does not satisfy dependency checks. Status `skipped` satisfies dependency checks (tasks depending on a skipped task can proceed). Status `rejected` does NOT satisfy dependency checks. Invalid transitions produce an AIError with code `INVALID_TRANSITION`.
**Context**: Prevents invalid task lifecycle progression and ensures data integrity of the task index. The state machine uses a role-based transition table (`ValidateTransition` with roles: submit, claim, reopen, auto, manual) for fine-grained access control. Terminal states are strictly protected — no `--force` override exists; `forge task reopen` is the only recovery path for rejected/skipped tasks.
**Source**: feature/forge-cli-v3 BIZ-001

| Current State | Terminal? | Notes |
|---------------|-----------|-------|
| pending | No | Can transition to any non-terminal state |
| in_progress | No | Cannot go to completed via status command; must use submit |
| completed | Yes | Irreversible; no recovery path |
| blocked | No | Auto-restores to pending when all deps complete; manual unblock via `forge task transition` |
| suspended | No | Manual hold only; entered/exited via `forge task transition`; does not satisfy deps |
| skipped | Yes | Satisfies dependency checks; can reopen to pending via `forge task reopen` |
| rejected | Yes | Does not satisfy dependency checks; can reopen to pending via `forge task reopen` |

### BIZ-task-lifecycle-002: Terminal State Immutability

**Rule**: Terminal states (`completed`, `rejected`, `skipped`) are enforced by the `ValidateTransition` state machine in `pkg/task/statemachine.go`. `forge task submit` validates transitions via `ValidateTransition(current, "completed", RoleSubmit)` before proceeding — attempting to submit a task already in a terminal state will fail with a `TransitionError`. The only recovery path for `rejected`/`skipped` tasks is `forge task reopen` (transitions to `pending`). `completed` tasks are truly irreversible — no command can transition them. `forge task status` no longer supports `--force`; all manual overrides go through `forge task transition` (which also respects terminal state protection).
**Context**: Guarantees that completed/rejected/skipped tasks cannot be accidentally overwritten. The role-based transition table is the single authority for state validation.
**Source**: feature/forge-cli-v3 BIZ-002

## Type Validation

### BIZ-task-lifecycle-003: System Type Exclusion

**Rule**: Non-auto-generated tasks (`.md` files created by Skills or users) MUST NOT use system-reserved types. System types are defined in `SystemTypes` map (`pkg/task/types.go`) with exactly 12 base types: `gate`, `test.gen-journeys`, `test.gen-contracts`, `test.gen-scripts`, `test.run`, `eval.journey`, `eval.contract`, `validation.code`, `validation.ux`, `doc.review`, `doc.summary`, `code-quality.simplify`. Surface-specific variants (e.g., `test.gen-scripts.cli`, `test.run.api`) are dynamically recognized by `IsSystemType()` stripping the last `.<surface>` segment and checking the base type. Auto-generated tasks (identified by `IsAutoGenTaskID()` matching `T-test-*`, `T-quick-*`, `T-specs-*`, `T-clean-*`, `T-validate-*`, `T-eval-*`, `T-review-doc`, `*.gate`, `*.summary` ID patterns) are exempt. Enforcement occurs in both `BuildIndex()` and `validate-index`. Error message includes the specific invalid type and full system type list.

**Dual-identity exception**: `doc.consolidate` and `doc.drift` are NOT in SystemTypes — they can be both auto-generated (by `forge task index`) and manually created by Skills for legacy projects.

测试类型命名遵循 Surface → Test Type 映射，权威定义参见 `docs/reference/test-type-model.md`。

**Context**: Prevents Skills from accidentally assigning pipeline-managed types to business tasks, which would cause scheduling anomalies (wrong stage-gate routing, test pipeline misdetection). The blacklist approach avoids maintenance burden since system types form a stable closed set while business types grow.
**Source**: feature/system-type-exclusion
