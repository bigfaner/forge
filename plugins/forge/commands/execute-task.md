---
name: execute-task
description: Execute single task via claim, dispatch, and verify orchestrator.
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
- `SURFACE_KEY` (e.g., "admin-panel" or "" — defaults to "" if absent)
- `SURFACE_TYPE` (e.g., "web", "cli", "api" or "" — defaults to "" if absent)
- `FEATURE` (e.g., "my-feature" — feature slug from claim output)

## Step 1.5: Main Session Routing

If `MAIN_SESSION == "true"`:

1. Read the task file at the FILE path extracted from claim output and find the `## Main Session Instructions` section.
2. Follow the instructions exactly — the task document specifies what skill to invoke, how to check outcome, and how to record the result.
3. The dispatcher does NOT hardcode skill names or record logic — it delegates to the task document.
4. If the task file lacks a `## Main Session Instructions` section, create a fix task to block it:
   ```bash
   forge task add --type coding.fix --title "Fix: MAIN_SESSION missing instructions" --source-task-id <TASK_ID> --block-source --var SOURCE_FILES="<affected-files>" --var TEST_SCRIPT="<test-path>" --var TEST_RESULTS="<test-output>" --description "MAIN_SESSION task missing Main Session Instructions section"
   ```
   Then skip to Step 3 (STOP).
5. After execution, verify the record file exists via `forge task status <TASK_ID>`. If STATUS is not `"completed"`, spawn fix task using `--block-source` (same as Step 2b verify logic).
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
- **STATUS == `"blocked"`** (auto-downgraded): create fix task using `--block-source`:
  ```bash
  forge task add --type coding.fix --title "Fix: <failure>" \
    --source-task-id <TASK_ID> \
    --block-source \
    --var SOURCE_FILES="<affected-files>" \
    --var TEST_SCRIPT="<test-path>" \
    --var TEST_RESULTS="<test-output>" \
    --description "Task <TASK_ID> was auto-downgraded to blocked — test failures or record issues"
  ```
  `forge task add` automatically deduplicates — check output:
  - `ACTION: ADDED` → new fix task created
  - `ACTION: SKIPPED` → active fix task already exists
  Proceed to Step 3 (STOP).
- **STATUS == `"in_progress"`** (record missing): proceed to 2c.

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

Output your final summary in this format and STOP:

```
Task <TASK_ID>: <status> | <one-line summary of what was accomplished or why it failed>
```
</HARD-RULE>

## Error Handling

| Situation | Action |
|-----------|--------|
| No available task | Stop, report |
| Agent timeout | Create fix task to block the timed-out task, then stop: `forge task add --type coding.fix --title "Fix: agent timeout" --source-task-id <TASK_ID> --block-source --var SOURCE_FILES="<affected-files>" --var TEST_SCRIPT="<test-path>" --var TEST_RESULTS="<test-output>" --description "Agent timed out after 30 minutes"` |
| Record missing | Dispatch `Agent(subagent_type="forge:task-executor", prompt="Fix record for task <TASK_ID>")` — subagent calls `forge prompt get-by-task-id --fix-record-missed` internally |
| Main session task fails | Follow error handling in task document's `### Error Handling` section; if missing, `forge task add --type coding.fix --title "Fix: main session task failed" --block-source --source-task-id <TASK_ID> --var SOURCE_FILES="<affected-files>" --var TEST_SCRIPT="<test-path>" --var TEST_RESULTS="<test-output>" --description "Main session task failed"` |

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
