---
name: run-tasks
description: Autonomous task dispatcher that continuously claims tasks and dispatches to subagents.
allowed_tools: ["Bash", "Read", "Agent", "TaskOutput"]
---

# /run-tasks

Auto-dispatch tasks to subagents. Main session only handles dispatching.

## Architecture

```
MAIN SESSION (Dispatcher)
   ┌─────────────┐    ┌─────────────┐    ┌─────────────┐
   │ 1. Claim    │───▶│ 2. Dispatch │───▶│ 3. Verify   │
   │    Task     │    │   + Timeout │    │   Record    │
   └─────────────┘    └─────────────┘    └─────────────┘
          ▲                                    │
          └────────────────────────────────────┘
                      LOOP
```

## Dispatcher Iron Laws

<EXTREMELY-IMPORTANT>
1. Only 3 actions: claim → dispatch → verify
2. NO code reading, NO code writing
3. NO running tests directly
4. 30-minute timeout per task
5. 3 consecutive failures → STOP
</EXTREMELY-IMPORTANT>

## Execution Loop

### Step 1: Claim Task

```bash
task claim
```

**Output parsing**:
- `ACTION: CLAIMED` → New task
- `ACTION: CONTINUE` → Resume existing task
- Error → No available task, end loop

### Step 2: Dispatch with Timeout

```
Agent(
  subagent_type="task-executor",
  prompt="TASK_KEY: {{KEY}}
TASK_ID: {{ID}}
TASK_FILE: {{FILE}}

IMPORTANT: Do NOT claim or start any other tasks after completing this one. Stop after recording the task result."
)
```

**Timeout**: 30 minutes

### Step 3: Verify Record

Check if record file exists after agent completes.

### Step 4: Continue Loop

Return to Step 1.

## Error Handling

| Situation | Action |
|-----------|--------|
| No available task | End loop, print summary |
| Agent timeout | Mark blocked, continue next |
| Record missing | Dispatch error-fixer (include: "Use /record-task skill to create record") |
| 3 consecutive failures | STOP dispatcher |

### Error-Fixer Dispatch

When dispatching error-fixer for missing record, include explicit instruction:

```
Agent(
  subagent_type="error-fixer",
  prompt="TASK_ID: {{ID}}
  ERROR_MESSAGES: Missing task record
  INSTRUCTION: Use /record-task skill to create the record (task record CLI is mandatory)"
)
```

## Related Commands

| Command | Usage |
|---------|-------|
| `/execute-task` | Manual single task |
| `/claim-task` | Claim task only |
| `/record-task` | Create record + update status |
