---
created: 2026-05-22
author: "forge:quick"
status: Draft
---

# Proposal: Enforce `forge task add` in Auto-Scheduling Loop for Accurate Dispatch Order

## Problem

The run-tasks auto-scheduling loop uses ambiguous "spawn fix task" instructions that leave the AI agent to interpret how to create fix tasks. This can result in:
- Direct file writes bypassing `index.json` synchronization
- Incorrect `--source-task-id` or missing `--block-source` flags
- Fix tasks added without proper dependency chaining, breaking dispatch order

Meanwhile, the `task-executor` subagent lacks the ability to self-heal when encountering complex compilation errors — it can only mark tasks as `blocked` and return, forcing the dispatcher to reason about what happened.

### Evidence

- `run-tasks.md` lines 63, 74, 91: all say "spawn fix task" without specifying the command
- `execute-task.md` (lines 70-76) already has the correct pattern: `forge task add --template fix-task --source-task-id <ID> --block-source`
- The `task.AddTask()` function already has dedup (via `HasActiveFixTasks()`), auto-ID generation, and dependency injection — the plumbing exists, the instructions don't use it
- task-executor.md currently FORBIDS all task creation, even when AI detects a complex error during execution

### Urgency

Every run-tasks cycle without explicit `forge task add` risks inconsistent index state. As the system scales to more concurrent features, the probability of dispatch ordering bugs increases. The cost of fixing now is low (two files); the cost of debugging later is high.

## Proposed Solution

Make task creation in auto-scheduling loops explicit and consistent:

1. **run-tasks.md**: Replace all "spawn fix task" with the exact `forge task add` command (matching the already-correct pattern in execute-task.md)
2. **task-executor.md**: Grant limited `forge task add` capability for AI-detected complex errors — when the executor encounters a blocking error (e.g., large compilation failure), it can pause the current task, add a specialized coding.* fix task, and return to the dispatcher

### Innovation Highlights

This is a straightforward adoption of an existing mechanism. `forge task add` with `--block-source` and the dedup logic in `task.AddTask()` already handle all the ordering concerns. The proposal simply closes the gap between instruction and correct tool usage.

## Requirements Analysis

### Key Scenarios

- **Happy path**: task-executor executes a task successfully, no fix tasks needed
- **Complex error during execution**: task-executor encounters a large compilation failure, creates a fix task via `forge task add`, blocks the current task, returns to dispatcher
- **Dispatched task blocked**: run-tasks verifies a dispatched task is blocked, creates fix task via `forge task add`, continues loop
- **Main session task blocked**: run-tasks detects main session task failure, creates fix task via `forge task add`, continues loop
- **Dedup safety**: if task-executor already created a fix task, run-tasks calling `forge task add` again hits `HasActiveFixTasks()` and skips gracefully

### Non-Functional Requirements

- Fix tasks must have correct dependency chaining (`--block-source` + `--source-task-id`)
- Must not introduce duplicate fix tasks (built-in dedup)
- task-executor must stop immediately after pausing (no continued execution)

### Constraints & Dependencies

- Must respect the `ONE TASK PER INVOCATION` constraint in task-executor — pausing IS ending
- Must not break existing "mark blocked on prompt failure" behavior
- `forge task add` CLI must be available in the execution environment (it's built into forge-cli)

## Alternatives & Industry Benchmarking

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | No work | Ambiguous instructions persist, dispatch order risk | Rejected: cost of inaction exceeds effort |
| Direct file writes | Current "spawn fix task" | Simple for single case | Bypasses index, no dependency tracking, no dedup | Rejected: fragile |
| Dedicated error-handler agent | Custom subagent | Clean separation | Over-engineered for this scope | Rejected: too complex |
| **Chosen: explicit forge task add** | Pattern from execute-task.md | Already proven, built-in dedup & ordering, minimal change | Requires AI agent to use correct flags | **Selected: lowest risk, highest leverage** |

## Feasibility Assessment

### Technical Feasibility

Fully feasible. `forge task add` is already implemented, tested (1958 lines of add_test.go), and used in execute-task.md. Changes are limited to markdown instruction text.

### Resource & Timeline

Estimated: 2-3 tasks (1-2 doc updates + 1 validation). Straightforward, no new code needed.

### Dependency Readiness

All dependencies are in-tree: forge-cli provides `forge task add`, `forge task status`, and the underlying `task.AddTask()` with dedup and dependency injection.

## Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| "task-executor should never create tasks" | XY Detection | The real need is correct task ordering — the executor creating fix tasks under AI judgment and then stopping achieves the same ordering guarantee as the dispatcher doing it, with less round-trip latency |
| "spawn fix task is clear enough" | Occam's Razor | The ambiguity forces the AI to guess. The simplest unambiguous instruction is the exact command |

## Scope

### In Scope
- `run-tasks.md`: replace "spawn fix task" with explicit `forge task add` commands, including proper flags (`--template fix-task`, `--source-task-id`, `--block-source`)
- `task-executor.md`: add "complex error → pause → add fix task → STOP" capability with explicit forge task add command
- Both files: maintain backward compatibility with all existing workflows

### Out of Scope
- `execute-task.md` (already has correct pattern)
- `quality_gate.go` (Go code uses internal API, not CLI)
- Other loops or agent files
- Changes to `forge task add` CLI or backend logic
- New tool/subagent creation

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| task-executor creates fix tasks unnecessarily | L | M | AI judgment with clear guidelines: only for complex blocking errors, not for routine failures |
| Dedup fails in edge case | L | H | Already covered by 1958 lines of tests for AddTask(); no logic changes proposed |
| AI ignores explicit command and still writes files | L | M | Instructions use `forge task add` with exact code block; AI follows instructions from markdown commands |

## Success Criteria

- [x] `run-tasks.md` no longer contains the ambiguous phrase "spawn fix task" — all fix task creation uses explicit `forge task add` commands
- [x] `task-executor.md` documents the complex-error pause flow with exact `forge task add` command and flag set
- [x] Existing "mark blocked on prompt failure" behavior is preserved in task-executor.md
- [x] Run-tasks dispatcher handles the case where a fix task was already created by task-executor (no duplicate fix tasks)

## Next Steps

- Proceed to task generation via `/quick-tasks`, then auto-execute via `/run-tasks`