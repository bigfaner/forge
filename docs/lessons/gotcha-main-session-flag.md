---
created: "2026-05-09"
tags: [architecture, task-dispatcher, main-session]
---

# MAIN_SESSION Flag Must Be Respected by Task Dispatcher

## Problem

The `/run-tasks` dispatcher received a task with `MAIN_SESSION: true` in the claim output but still dispatched it to a `forge:task-executor` subagent. This causes the task to fail because the task-executor cannot spawn the subagents that skills like `/eval-test-cases` require.

## Root Cause

Causal chain (3 levels deep):

1. **Symptom**: T-test-1b (eval-test-cases) was dispatched to a task-executor agent despite `MAIN_SESSION: true`
2. **Direct cause**: The dispatcher ignored the `MAIN_SESSION` flag from the claim output and treated all tasks uniformly as "dispatch to subagent"
3. **Root cause**: The `/run-tasks` description said "Main session only handles dispatching" — this framing created a false assumption that ALL tasks must be dispatched, overriding the `MAIN_SESSION` signal
4. **Trigger condition**: Tasks like `/eval-test-cases` that involve interactive skill invocation (doc-scorer/doc-reviser orchestration) are designed for main-session execution where the skill can orchestrate sub-agents directly

## Solution

In `/run-tasks` dispatcher loop, after Step 1 (claim), check `MAIN_SESSION`:

```
if MAIN_SESSION == true:
    Read task file → follow Main Session Instructions → verify record → continue loop
else:
    Dispatch to forge:task-executor subagent (Step 2)
```

The `MAIN_SESSION` flag overrides the default dispatch behavior. Not all tasks are subagent-compatible — some skills require main-session orchestration capabilities.

## Why Subagents Cannot Spawn Subagents (Architectural Constraint)

- **Subagents lack the Agent tool** — Only the main session has access to the Agent tool. Subagents run with a restricted toolset and cannot spawn sub-subagents.
- **Even explicit frontmatter declaration doesn't help** — Adding `allowed_tools: ["Agent"]` in subagent definition frontmatter does not grant the ability.
- **Design rationale** — Subagents run in isolated contexts with limited tool access and return a single final message. This prevents infinite recursion and unbounded resource consumption.

## Changes Applied

1. **`run-tasks.md`**: Fixed contradictory description, strengthened Step 1.5 with missing-instruction handling and record verification. Error fallback uses `rejected` status.
2. **`execute-task.md`**: Added Step 0 MAIN_SESSION routing for manual invocation.
3. **`task-executor.md`**: Added Step 0 MAIN_SESSION guard (defense in depth).
4. **`breakdown-tasks/SKILL.md`**: Added HARD-RULE documenting when `mainSession: true` is needed.
5. **`quick-tasks/SKILL.md`**: Added `mainSession` note to index.json rules.
6. **task-cli**: Added `rejected` status — terminal state for tasks that ran but did not pass acceptance criteria.
7. **`guide.md`**: Updated task status list to include `rejected`.

## Related Files

- `plugins/forge/commands/run-tasks.md` — Dispatcher protocol
- `plugins/forge/commands/execute-task.md` — Manual task executor
- `plugins/forge/agents/task-executor.md` — Subagent with MAIN_SESSION guard
- `plugins/forge/skills/eval-test-cases/SKILL.md` — Eval skill definition
- `plugins/forge/skills/breakdown-tasks/SKILL.md` — Task breakdown with mainSession rule
- `plugins/forge/skills/breakdown-tasks/templates/eval-test-cases.md` — Task template (T-test-1b)

## References

- [GitHub Issue #4182 — Sub-Agent Task Tool Not Exposed When Launching Nested Agents](https://github.com/anthropics/claude-code/issues/4182)
- [GitHub Issue #50306 — Allow subagents to declare nested-Agent capability](https://github.com/anthropics/claude-code/issues/50306)
