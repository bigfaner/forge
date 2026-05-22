---
id: "1"
title: "Enforce explicit forge task add in run-tasks dispatcher"
priority: "P1"
estimated_time: "1h"
dependencies: []
type: "doc"
mainSession: false
---

# 1: Enforce explicit forge task add in run-tasks dispatcher

## Description

The run-tasks dispatcher (`plugins/forge/commands/run-tasks.md`) currently uses ambiguous "spawn fix task" instructions. Replace every instance with the explicit `forge task add` command using correct flags (`--template fix-task`, `--source-task-id`, `--block-source`), matching the already-correct pattern in `execute-task.md`. The `forge task add` built-in dedup handles the case where task-executor already created a fix task.

## Reference Files
- `docs/proposals/enforce-forge-task-add-in-loop/proposal.md` — Source proposal

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/commands/run-tasks.md` | Replace all "spawn fix task" with explicit `forge task add` commands |

## Acceptance Criteria
- [ ] run-tasks.md no longer contains the phrase "spawn fix task" — every fix task creation uses explicit `forge task add --template fix-task --source-task-id <ID> --block-source --description "<reason>"`
- [ ] Step 1.5 (main session failure) uses explicit `forge task add` with proper flags
- [ ] Step 2b (dispatched task blocked) uses explicit `forge task add` with proper flags
- [ ] Step 3 (main session error table) uses explicit `forge task add` with proper flags
- [ ] Each instance includes correct `--source-task-id <TASK_ID>` and `--block-source` flags
- [ ] Each instance includes an informative `--description` explaining the failure reason
- [ ] Failure tracking comment (consecutive_failures) remains accurate

## Implementation Notes
- Pattern to follow from `execute-task.md`:
  ```bash
  forge task add --template fix-task \
    --source-task-id <TASK_ID> \
    --block-source \
    --description "<reason>"
  ```
- `forge task add` has built-in dedup via `HasActiveFixTasks()` — calling it when a fix task already exists is harmless
- The `consecutive_failures` counter is still incremented on fix-task creation (line 37)
- No logic changes to the dispatcher workflow — just making instructions unambiguous