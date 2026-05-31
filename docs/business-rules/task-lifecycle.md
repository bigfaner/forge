---
title: "Task Lifecycle Rules"
domains: [state-machine, transition, terminal, blocked, completed, pending, suspended, skipped, reopen, type-validation, system-types, topological, ordering, claim-priority, sizing, complexity, audit, ac-count]
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

**Rule**: Non-auto-generated tasks (`.md` files created by Skills or users) MUST NOT use system-reserved types. System types are defined in `SystemTypes` map (`pkg/task/types.go`) with exactly 12 base types: `gate`, `test.gen-journeys`, `test.gen-contracts`, `test.gen-scripts`, `test.run`, `eval.journey`, `eval.contract`, `validation.code`, `validation.ux`, `doc.review`, `doc.summary`, `code-quality.simplify`. Surface-specific variants (e.g., `test.gen-scripts.cli`, `test.run.api`) are dynamically recognized by `IsSystemType()` stripping the last `.<surface>` segment and checking the base type. Auto-generated tasks (identified by `IsAutoGenTaskID()` matching `T-test-*`, `T-quick-*`, `T-specs-*`, `T-clean-*`, `T-validate-*`, `T-eval-*`, `T-review-doc`, `*.gate`, `*.summary` ID patterns) are exempt. Enforcement occurs in both `BuildIndex()` and `forge task validate`. Error message includes the specific invalid type and full system type list.

**Dual-identity exception**: `doc.consolidate` and `doc.drift` are NOT in SystemTypes — they can be both auto-generated (by `forge task index`) and manually created by Skills for legacy projects.

测试类型命名遵循 Surface → Test Type 映射，权威定义参见 `docs/reference/test-type-model.md`。

**Context**: Prevents Skills from accidentally assigning pipeline-managed types to business tasks, which would cause scheduling anomalies (wrong stage-gate routing, test pipeline misdetection). The blacklist approach avoids maintenance burden since system types form a stable closed set while business types grow.
**Source**: feature/system-type-exclusion

## Topological Ordering

### BIZ-task-lifecycle-004: Topological Task Ordering

**Rule**: `forge task list` and `forge task claim` use topological ordering (Kahn's algorithm) to sort tasks by dependency depth. `forge task list` defaults to topological order (`--sort topo`), with `--sort id` as fallback for natural ID ordering. `forge task claim` uses `computeTopoDepths()` to prefer claiming tasks with lower topological depth (closer to dependency roots) among eligible tasks, breaking ties by priority (P0 > P1 > P2), then by natural ID order. `forge task list --tree` renders an interactive TUI dependency tree. Wildcard dependencies (e.g., `1.x`) are expanded via `ResolveWildcardDep` before building the adjacency list. Cycle detection is built-in: tasks participating in dependency cycles are reported and excluded from ordering.
**Context**: Topological ordering ensures that dependency-first execution is the default, preventing agents from claiming tasks whose dependencies haven't completed. This replaces the previous natural-ID-only ordering.
**Source**: feature/task-list-topological-order [auto-specs]

## Task Sizing

### BIZ-task-lifecycle-005: Task Sizing Constraints

**Rule**: `forge task validate` enforces AC (Acceptance Criteria) count constraints: each task must have >= 1 and <= 6 AC items (detected by parsing `## Acceptance Criteria` section for `- [ ]` lines). Tasks with 0 AC or > 6 AC cause validation failure (exit 1) with the specific task file name and AC count in the error message. The `breakdown-tasks` and `quick-tasks` skills include a mandatory Task Sizing Audit step after generating all task files but before running `forge task index`. The audit checks for: (1) multi-verb titles linking independent actions, (2) AC cross-domain coverage, (3) operational ceiling (> 8 files with same pattern). Violations trigger automatic task splitting.
**Context**: Prevents oversized tasks that increase the risk of agent timeout, macOS sleep interruption, and scope creep. Task 11 (reorganize internal/cmd/) demonstrated that rules written in documentation are insufficient — LLM ignored multi-verb and file-count rules, resulting in 26 minutes of work followed by 9.2 hours of sleep-suspended execution.
**Source**: feature/task-sizing-gate [auto-specs]

### BIZ-task-lifecycle-006: Task Complexity Classification

**Rule**: Each task carries a `complexity` field in its YAML frontmatter with values `low`, `medium`, or `high`. Default heuristic: `low` = AC <= 3 AND no Hard Rules AND Reference Files <= 1; `high` = AC >= 5 OR has Hard Rules; `medium` = everything else. LLM may override when static metrics conflict with cognitive judgment (e.g., AC <= 3 but involves multi-file architectural change), recording the reason in the task's Implementation Notes.
**Context**: Complexity classification enables task routing and executor behavior adaptation. The heuristic provides a fast default, while LLM override handles edge cases that static metrics miss.
**Source**: feature/task-sizing-gate [auto-specs]
