---
name: execute-task
description: Execute single task with focused TDD workflow.
allowed-tools: Bash Read Agent TaskOutput Skill
---

# /execute-task

Execute a single task. MAIN_SESSION tasks execute in main session; all others dispatch to forge:task-executor subagent (which calls `forge prompt get-by-task-id` internally).

## Step 1: Claim Task

```bash
forge task claim
```

**Output parsing**:
- `ACTION: CLAIMED` → New task
- `ACTION: CONTINUE` → Resume existing task
- Error → No available task, stop

**Extract from claim output**:
- `TASK_ID` (e.g., "2.1")
- `FILE` (e.g., full absolute path to task file)
- `MAIN_SESSION` (e.g., "true" or absent)
- `SCOPE` (e.g., "frontend", "backend", or "all" — defaults to "all" if absent)
- `FEATURE` (e.g., "my-feature" — feature slug from claim output)

## Step 1.5: Main Session Routing

If `MAIN_SESSION == "true"`:

1. Read the task file at the FILE path extracted from claim output and find the `## Main Session Instructions` section.
2. Follow the instructions exactly — the task document specifies what skill to invoke, how to check outcome, and how to record the result.
3. The dispatcher does NOT hardcode skill names or record logic — it delegates to the task document.
4. If the task file lacks a `## Main Session Instructions` section, mark the task blocked and report: "MAIN_SESSION task missing Main Session Instructions section — task document is incomplete".
5. After execution, verify the record file exists via `forge task status <TASK_ID>`. If STATUS is not `"completed"`, spawn fix task (same as Step 2 verify logic).
6. Skip to Step 3 (STOP).

Else:
- Proceed to Step 2 (Dispatch + Verify).

## Step 2: Dispatch + Verify

### 2a. Dispatch

```
Agent(
  subagent_type="forge:task-executor",
  prompt="Execute task <TASK_ID>"
)
```

The subagent internally runs `forge prompt get-by-task-id <TASK_ID>` to get the execution strategy.

**Timeout**: 30 minutes

### 2b. Verify Record

After subagent returns, check the task's actual status via CLI:

```bash
forge task status <TASK_ID>
```

- **STATUS == `"completed"`**: proceed to Step 3 (STOP).
- **STATUS != `"completed"`**: task was auto-downgraded (e.g. test failures).
  Spawn fix task using `--block-source` to atomically block the source:
  ```bash
  forge task add --template fix-task --title "Fix: <failure>" \
    --source-task-id <TASK_ID> \
    --block-source \
    --description "<reason>"
  ```
  `forge task add` automatically deduplicates — check output:
  - `ACTION: ADDED` → new fix task created
  - `ACTION: SKIPPED` → active fix task already exists

### 2c. Record-Missing Recovery

When the subagent completes but the record file is missing, delegate recovery to another subagent:

```
Agent(
  subagent_type="forge:task-executor",
  prompt="Fix record for task <TASK_ID>"
)
```

The subagent's Execution Protocol detects the "Fix record for" prefix and calls `forge prompt get-by-task-id <TASK_ID> --fix-record-missed` internally.

## Step 3: STOP

<HARD-RULE>
ONE TASK PER INVOCATION. This is absolute and non-negotiable.

After Step 2, you MUST stop immediately.

<PROHIBITIONS>
- Running `forge task claim` under any circumstances
- Reading the next task file
- Continuing with any additional work
</PROHIBITIONS>

Output your final summary and STOP.
</HARD-RULE>

## Error Handling

| Situation | Action |
|-----------|--------|
| No available task | Stop, report |
| Agent timeout | Mark blocked, stop |
| Record missing | Dispatch `Agent(prompt="Fix record for task <TASK_ID>")` — subagent calls `forge prompt get-by-task-id --fix-record-missed` internally |
| Main session task fails | Follow error handling in task document's `### Error Handling` section; if missing, `forge task add --template fix-task --title "Fix: main session task failed" --block-source --source-task-id <TASK_ID>` |

## Rules

<EXTREMELY-IMPORTANT>
- submit-task is mandatory — No completion without it
- All verifications must pass
- ONE TASK PER INVOCATION — after Step 2, STOP immediately, no exceptions
- FORBIDDEN: run "forge task claim", read index.json, or start any subsequent task
- Do NOT use TASK_FILE parameter when dispatching to forge:task-executor
</EXTREMELY-IMPORTANT>

## Related Commands

| Command | Usage |
|---------|-------|
| `/run-tasks` | Auto-execute all tasks |
| `/submit-task` | Create record + update status |
