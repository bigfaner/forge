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

The task-executor subagent (`plugins/forge/agents/task-executor.md`) executes strategy steps for a coding.* task. When it encounters a recurring error that resists multiple attempts within the same execution, it should be able to pause the current task, create a specialized fix task, and return to the dispatcher.

The decision flow:
1. Execute strategy steps normally
2. If a step fails → AI tries to resolve it (~3 approaches max)
3. If the same/similar error keeps recurring → AI judges this needs a dedicated fix task
4. Run `forge task add --template fix-task --source-task-id <TASK_ID> --block-source --description "<error summary>"`
5. Output `PAUSE: <TASK_ID> | added fix-task <FIX_ID> | <reason>`
6. STOP — return to dispatcher

Existing behaviors preserved:
- If `forge prompt get-by-task-id` fails → mark blocked, STOP (unchanged)
- If task completes normally → submit, commit, output DONE (unchanged)
- One-shot command failures resolved on retry → continue normally (no fix task)

## Reference Files
- `docs/proposals/enforce-forge-task-add-in-loop/proposal.md` — Source proposal

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/agents/task-executor.md` | Add complex error pause flow with forge task add |

## Acceptance Criteria
- [ ] task-executor.md documents the AI decision criteria for creating a fix task (recurring error after ~3 attempts within the same execution)
- [ ] task-executor.md includes the exact `forge task add` command with correct flags
- [ ] task-executor.md specifies the `PAUSE` output format
- [ ] task-executor.md specifies that the executor must STOP after pausing
- [ ] Existing hard constraints remain intact: ONE TASK PER INVOCATION, submit-task mandatory, NO background tasks
- [ ] Existing "mark blocked on prompt failure" behavior is preserved
- [ ] The simple-error flow (one-off failure resolved on retry) is documented as NOT warranting a fix task
- [ ] The forbidden actions list is updated to allow `forge task add` only for the pause flow

## Implementation Notes
- The `forge task add` command already has dedup — if the dispatcher later tries to add another fix for the same source, `HasActiveFixTasks()` will detect the one created by the executor and skip
- The `--block-source` flag automatically sets the source task to `blocked`, preventing re-claim until the fix task is resolved
- The `submit.go` auto-restore mechanism will automatically unblock the source task when the fix task completes
- The executor must NOT continue execution after pausing — it's ONE TASK PER INVOCATION
- The executor's Hard Constraints #5 currently says "FORBIDDEN: run forge task claim, read index.json, or start any subsequent task" — this should be clarified to allow `forge task add` for the pause flow