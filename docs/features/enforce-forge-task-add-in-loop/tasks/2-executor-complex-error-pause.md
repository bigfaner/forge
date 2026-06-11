---
id: "2"
title: "Add complex error pause capability to task-executor"
priority: "P1"
estimated_time: "1h"
dependencies: ["1"]
type: "doc"
mainSession: false
---

# 2: Add complex error pause capability to task-executor

## Description

The task-executor subagent (`plugins/forge/agents/task-executor.md`) needs the ability to pause the current task and create a dedicated `coding.fix` task when it encounters a complex, recurring error during execution.

### Error Classification

Not all errors warrant a fix task. The executor uses AI judgment:

| Error Type | Examples | Action |
|------------|----------|--------|
| Simple/transient | Network timeout, missing dependency, single command failure, formatting lint | Inline fix (retry or auto-fix), continue execution |
| Complex/recurring | Same error persists after ~3 attempts, large compilation failure, cross-file refactoring needed | Create `coding.fix` task, pause current task |

### Decision Flow

```
Execute strategy step → error
  → Can AI fix inline? (low effort, obvious cause)
    → Yes: fix it, continue
    → No: is this the same/similar error persisting after ~3 attempts?
      → No: try another approach
      → Yes: this is a complex error → create coding.fix task via forge task add
```

All fix tasks use `--template fix-task` (type `coding.fix`), regardless of error category.

### Pause Protocol

1. Run: `forge task add --template fix-task --source-task-id <TASK_ID> --block-source --description "<error classification and summary>"`
2. Output: `PAUSE: <TASK_ID> | added fix-task <FIX_ID> | <reason>`
3. STOP immediately — return to dispatcher

## Reference Files
- `docs/proposals/enforce-forge-task-add-in-loop/proposal.md` — Source proposal

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/agents/task-executor.md` | Add complex error pause flow with error classification and forge task add |

## Acceptance Criteria
- [ ] task-executor.md documents the error classification (inline fix vs. fix task)
- [ ] task-executor.md documents the decision flow: try ~3 approaches → create fix task if recurring
- [ ] task-executor.md includes the exact `forge task add` command with `--template fix-task`, `--source-task-id`, `--block-source`, `--description`
- [ ] task-executor.md specifies the `PAUSE` output format
- [ ] task-executor.md specifies that the executor must STOP immediately after pausing
- [ ] Existing hard constraints are updated: `forge task add` is allowed for the pause flow; `forge task claim`, `read index.json`, `start any subsequent task` remain forbidden
- [ ] Existing "mark blocked on prompt failure" behavior is preserved
- [ ] The simple-error flow (one-off failure resolved inline) is documented as NOT warranting a fix task

## Implementation Notes
- The `forge task add` command has built-in dedup — if the dispatcher later tries to add a fix for the same source, `HasActiveFixTasks()` skips gracefully
- `--block-source` automatically sets the source task to `blocked`, preventing re-claim until the fix resolves
- `submit.go` auto-restore mechanism unblocks the source task when the fix task completes
- The executor must NOT continue execution after pausing — hard constraint: ONE TASK PER INVOCATION
- The executor's Hard Constraints #5 currently says "FORBIDDEN: run forge task claim, read index.json, or start any subsequent task" — update to allow `forge task add` for pause flow while keeping others forbidden