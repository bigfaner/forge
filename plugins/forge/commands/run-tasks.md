---
name: run-tasks
description: Autonomous task dispatcher that continuously claims tasks and dispatches to subagents.
allowed-tools: Bash Read Agent Skill
---

# /run-tasks

Auto-dispatch tasks. MAIN_SESSION tasks execute in main session; all others dispatch to forge:task-executor subagent.

## Architecture

```mermaid
flowchart TD
    S0["0. Set Active Feature"] --> A["1. Claim Task"]
    A --> B{"MAIN_SESSION?"}
    B -->|"yes"| C["1.5 Follow Task Instructions"]
    C --> LOOP(["Step 3: Continue Loop"])
    B -->|"no"| D["2. Dispatch + Verify"]
    D --> LOOP
    LOOP --> A
```

## Dispatcher Iron Laws

<EXTREMELY-IMPORTANT>
1. Only 3 actions: claim → (main_session? follow task instructions : dispatch+verify) → continue loop
2. NO code reading, NO code writing — EXCEPT for MAIN_SESSION tasks (Step 1.5) where reading the task file and invoking the Skill tool are required
3. NO running tests directly — the CLI submit gate handles quality checks at task submission
4. 30-minute timeout per task
5. 3 consecutive failures → STOP (tracked by failure counter below)
6. NO `run_in_background`, NO `TaskOutput` polling — Agent call is blocking, wait for return
</EXTREMELY-IMPORTANT>

## Execution Loop

**Failure tracking**: maintain `consecutive_failures` (starts at 0). Increment on: fix-task creation, record-missing dispatch, agent timeout. Reset to 0 when a claim→dispatch→verify cycle ends with `STATUS == "completed"` (verified via `forge task status <TASK_ID>`). At 3: print summary (see format below) and STOP.

### Step 0: Set Active Feature

Runs **once** before the claim loop.

1. Determine the feature slug from the current context (proposal directory, manifest, or user input).
2. Run `forge feature set <slug>`. On success (exit code 0), the slug is printed to stdout. Proceed to Step 1.
3. On failure (non-zero exit): the slug is invalid or feature not found. Report the error to the user and STOP — do not enter the claim loop.

### Step 1: Claim Task

```bash
forge task claim
```

**Output**: `ACTION: CLAIMED` (new) | `ACTION: CONTINUE` (resume) | Error (no task, end loop).

**Extract**: `TASK_ID`, `TYPE`, `FILE`, `MAIN_SESSION`, `TASK_CATEGORY`, `SURFACE_KEY`, `SURFACE_TYPE`, `FEATURE`.

### Step 1.5: Main Session Routing

If `MAIN_SESSION == "true"`:

1. Read task file at `FILE`, find `## Main Session Instructions` section.
2. Follow instructions exactly (task document specifies skill, outcome, record logic).
3. If section missing: report error, create fix task to block it (derive fix type from `TASK_CATEGORY` per table in Error Handling), then continue to Step 3:
   ```bash
   forge task add --type <derived-fix-type> --title "Fix: MAIN_SESSION missing instructions" --source-task-id <TASK_ID> --block-source --var SOURCE_FILES="<affected-files>" --var TEST_SCRIPT="<test-path>" --var TEST_RESULTS="<test-output>" --description "MAIN_SESSION task missing Main Session Instructions section"
   ```
4. After execution, verify via `forge task status <TASK_ID>`. If STATUS != "completed", create fix task using `--block-source` (derive fix type from `TASK_CATEGORY` per table in Error Handling):
   ```bash
   forge task add --type <derived-fix-type> --title "Fix: MAIN_SESSION task failed" \
     --source-task-id <TASK_ID> \
     --block-source \
     --var SOURCE_FILES="<affected-files>" \
     --var TEST_SCRIPT="<test-path>" \
     --var TEST_RESULTS="<test-output>" \
     --description "Main session task <TASK_ID> failed — verify output and fix issues"
   ```
5. Skip to Step 3.

Else: proceed to Step 2.

### Step 2: Dispatch + Verify

**2a. Dispatch** — `Agent(subagent_type="forge:task-executor", prompt="Execute task <TASK_ID>")`. **Timeout**: 30 min. NO `run_in_background` — wait for Agent return.

**2b. Verify Record** — Run `forge task status <TASK_ID>`:
- **STATUS == "completed"**: proceed to Step 3 (Continue Loop).
- **STATUS == "blocked"** (auto-downgraded): create fix task using `--block-source` (derive fix type from `TASK_CATEGORY` per table in Error Handling):
  ```bash
  forge task add --type <derived-fix-type> --title "Fix: <failure>" \
    --source-task-id <TASK_ID> \
    --block-source \
    --var SOURCE_FILES="<affected-files>" \
    --var TEST_SCRIPT="<test-path>" \
    --var TEST_RESULTS="<test-output>" \
    --description "Dispatched task <TASK_ID> was auto-downgraded to blocked — test failures or record issues"
  ```
  Continue loop.
- **STATUS == "in_progress"** (no record created): proceed to 2c.

**2c. Record-Missing Recovery** — `Agent(subagent_type="forge:task-executor", prompt="Fix record for task <TASK_ID>")`. After 2c, re-verify via 2b logic.

### Step 3: Continue Loop

Return to Step 1.

## Error Handling

**Fix-Type Derivation**: When creating a fix task, extract `TASK_CATEGORY` from the claim output of the source task and derive the fix type:

| Source Task Category | Fix Task Type |
|----------------------|---------------|
| `doc`, `eval`        | `doc.fix`     |
| `coding`, `test`, `validation`, `gate` | `coding.fix` |

| Situation | Action |
|-----------|--------|
| No available task | End loop, print summary (see format below) |
| Agent timeout | Create fix task to block the timed-out task, increment `consecutive_failures`, continue loop: `forge task add --type <derived-fix-type> --title "Fix: agent timeout" --source-task-id <TASK_ID> --block-source --var SOURCE_FILES="<affected-files>" --var TEST_SCRIPT="<test-path>" --var TEST_RESULTS="<test-output>" --description "Agent timed out after 30 minutes"` |
| Record missing | Dispatch fix-record subagent (2c) |
| 3 consecutive failures | STOP |
| Main session fails | Follow task doc's error section; if missing, `forge task add --type <derived-fix-type> --title "Fix: main session task failed" --source-task-id <TASK_ID> --block-source --var SOURCE_FILES="<affected-files>" --var TEST_SCRIPT="<test-path>" --var TEST_RESULTS="<test-output>" --description "Main session task failed"` then continue |

### Summary Format

When the loop ends (no available task or 3 consecutive failures), print:

```
## Dispatch Summary

- Total claimed: <N>
- Completed: <N>
- Blocked: <N>
- Failed (fix-task created): <N>
- Consecutive failures at stop: <N>

<If any blocked/failed tasks, list them:>
- Task <ID>: <status> — <short reason>
```

## Post-Completion

After loop ends, print a conditional completion message.

**T-test-run** is a conventionally named auto-generated task (type `test`, title containing "test-run" or "run-tests") that `/breakdown-tasks` may include in the task index for surface-level test execution. When referenced below, it means: the task index contains at least one task whose type is `test` and whose purpose is running the full test suite for the feature.

- **Full pipeline mode** (tasks generated via `/breakdown-tasks`): "All tasks completed. T-test-run handles surface-level test execution automatically." If index lacks T-test-run, suggest: "Run `/run-tests` to execute surface-level tests."
- **Quick mode** (tasks generated via `/quick-tasks`): "All tasks completed. Test tasks generated by quick-tasks handle verification." If no test tasks exist in the index, suggest: "Run `/run-tests` to execute any available tests."

Do NOT run surface-level tests from the dispatcher.

Do NOT commit post-loop artifacts. The `forge feature complete --if-done` Stop hook detects uncommitted artifacts and blocks the agent to commit them via `/git-commit`. This ensures artifacts are committed only after quality-gate passes.

### Git Status Summary

After printing the completion message, display a concise git summary so the user has immediate visibility into the current state. Run the following commands, wrapping each in error handling — if any command fails, skip it silently and continue:

1. **Branch info**: `git branch --show-current` and `git rev-list --left-right --count main...HEAD` (shows ahead/behind relative to main). Print as: `On branch <name>, <n> ahead / <m> behind main`.

2. **Working tree changes**: `git status --short`. Print the raw output if non-empty, or "Working tree clean" if empty.

If all git commands fail (e.g., not in a git repository), print nothing — the completion message above is sufficient.
