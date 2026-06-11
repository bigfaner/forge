---
name: execute-task
description: Execute single task via claim, dispatch, and verify orchestrator.
allowed-tools: Bash Read Agent TaskOutput Skill
---

# /execute-task

Execute a single task. MAIN_SESSION tasks execute in main session; all others dispatch to forge:task-executor subagent (which calls `forge prompt get-by-task-id` internally).

> **Note**: This is a manual entry point for executing a single task. The automated pipeline (`/run-tasks`) uses its own dispatch loop and does not invoke `/execute-task`.

## Step 1: Claim Task

```bash
forge task claim
```

**Output**: `ACTION: CLAIMED` (new) | `ACTION: CONTINUE` (resume) | Error (stop).

**Extract**: `TASK_ID`, `TYPE`, `FILE`, `MAIN_SESSION` (`"true"`|`"false"`), `TASK_CATEGORY` (`doc`|`eval`|`coding`|`test`|`validation`|`gate`), `SURFACE_KEY`, `SURFACE_TYPE`, `FEATURE`.

## Step 1.5: Main Session Routing

If `MAIN_SESSION == "true"`:

1. Read task file at `FILE`, find `## Main Session Instructions` section.
2. Follow instructions exactly — the task document specifies what skill to invoke, outcome check, and record logic. The dispatcher does NOT hardcode skill names or record logic.
3. If section missing: create fix task to block it using template in Error Handling. Then skip to Step 3 (STOP).
4. After execution, verify via `forge task status <TASK_ID>`. If STATUS != `"completed"`, create fix task using template in Error Handling.
5. Skip to Step 3 (STOP).

Else: proceed to Step 2 (Dispatch + Verify).

## Step 2: Dispatch + Verify

### 2a. Dispatch

```
Agent(
  subagent_type="forge:task-executor",
  prompt="Execute task <TASK_ID>"
)
```

**Timeout**: 30 minutes

### 2b. Verify Record

After subagent returns, check the task's actual status via CLI:

```bash
forge task status <TASK_ID>
```

- **STATUS == `"completed"`**: proceed to Step 3 (STOP).
- **STATUS == `"blocked"`** (auto-downgraded): create fix task using template in Error Handling. Proceed to Step 3 (STOP).
- **STATUS == `"in_progress"`** (record missing): proceed to 2c.

### 2c. Record-Missing Recovery

When the subagent completes but the record file is missing, delegate recovery to another subagent:

```
Agent(
  subagent_type="forge:task-executor",
  prompt="Fix record for task <TASK_ID>"
)
```

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

**Fix-Type Derivation**: extract `TASK_CATEGORY` from claim output, map to fix type:

| Source Task Category | Fix Task Type |
|----------------------|---------------|
| `doc`, `eval`        | `doc.fix`     |
| `coding`, `test`, `validation`, `gate` | `coding.fix` |

**Fix task template** (used in all situations below):
```bash
forge task add --type <derived-fix-type> --title "Fix: <reason>" --source-task-id <TASK_ID> --block-source --var SOURCE_FILES="<affected-files>" --var TEST_SCRIPT="<test-path>" --var TEST_RESULTS="<test-output>" --description "<summary>"
```

| Situation | Action |
|-----------|--------|
| No available task | Stop, report |
| Agent timeout | Create fix task, stop |
| Record missing | Dispatch `Agent(subagent_type="forge:task-executor", prompt="Fix record for task <TASK_ID>")` |
| Main session task fails | Follow task document's `### Error Handling` section; if missing, create fix task |

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
