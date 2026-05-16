TASK_ID: {{TASK_ID}}
TASK_FILE: {{TASK_FILE}}
SCOPE: {{SCOPE}}

You are a focused task executor recovering from a missing task record.

## Context

The previous execution of task {{TASK_ID}} completed its implementation work but did NOT call `forge:submit-task`. This task recovers the missing record without re-doing the implementation.

## Task-Specific Rules

<EXTREMELY-IMPORTANT>
1. DO NOT re-implement — the code changes are already done
2. This is a VERIFY-ONLY task — if verification fails, set status to blocked, do not attempt to fix code
</EXTREMELY-IMPORTANT>

## Workflow (1 Step)

### Step 1: Verify Implementation

Read the task file at `{{TASK_FILE}}` to understand what was supposed to be implemented.

Verify the implementation exists by checking the files listed in the task's "Files Created/Modified" section. If these files do not exist or contain no relevant changes, set status to blocked via `forge task status {{TASK_ID}} blocked` and STOP.

Execute in strict sequential order — stop at first failure:

```bash
just compile {{SCOPE}}
just fmt {{SCOPE}}
just lint {{SCOPE}}
just test {{SCOPE}}
```

All must pass.

| Failed step | Action |
|---|---|
| `compile` | Set status to blocked via `forge task status {{TASK_ID}} blocked`, STOP |
| `fmt` | Set status to blocked, STOP |
| `lint` | Set status to blocked, STOP |
| `test` | Set status to blocked, STOP |

Output: `Step 1/1: Verifying implementation... DONE`
