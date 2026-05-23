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

**Failure tracking**: maintain `consecutive_failures` (starts at 0). Increment on: fix-task creation, record-missing dispatch, agent timeout. Reset to 0 on successful claim→dispatch→verify cycle. At 3: print summary and STOP.

### Step 0: Set Active Feature

Runs **once** before the claim loop.

1. Determine the feature slug from the current context (proposal directory, manifest, or user input).
2. Run `forge feature set <slug>`. On success (exit code 0), the slug is printed to stdout. Proceed to Step 1.

### Step 1: Claim Task

```bash
forge task claim
```

**Output**: `ACTION: CLAIMED` (new) | `ACTION: CONTINUE` (resume) | Error (no task, end loop).

**Extract**: `TASK_ID`, `FILE`, `MAIN_SESSION`, `SCOPE` (defaults "all"), `FEATURE`.

### Step 1.5: Main Session Routing

If `MAIN_SESSION == "true"`:

1. Read task file at `FILE`, find `## Main Session Instructions` section.
2. Follow instructions exactly (task document specifies skill, outcome, record logic).
3. If section missing: report error, create fix task to block it, then continue to Step 3:
   ```bash
   forge task add --template fix-task --title "Fix: MAIN_SESSION missing instructions" --source-task-id <TASK_ID> --block-source --description "MAIN_SESSION task missing Main Session Instructions section"
   ```
4. After execution, verify via `forge task status <TASK_ID>`. If STATUS != "completed", create fix task using `--block-source`:
   ```bash
   forge task add --template fix-task --title "Fix: MAIN_SESSION task failed" \
     --source-task-id <TASK_ID> \
     --block-source \
     --description "Main session task <TASK_ID> failed — verify output and fix issues"
   ```
   `forge task add` automatically deduplicates — check output:
   - `ACTION: ADDED` → new fix task created
   - `ACTION: SKIPPED` → active fix task already exists
5. Skip to Step 3.

Else: proceed to Step 2.

### Step 2: Dispatch + Verify

**2a. Dispatch** — `Agent(subagent_type="forge:task-executor", prompt="Execute task <TASK_ID>")`. Subagent calls `forge prompt get-by-task-id` internally. **Timeout**: 30 min. NO `run_in_background` — wait for Agent return.

**2b. Verify Record** — Run `forge task status <TASK_ID>`:
- **STATUS == "completed"**: proceed to Step 3 (Continue Loop).
- **STATUS == "blocked"** (auto-downgraded): create fix task using `--block-source`:
  ```bash
  forge task add --template fix-task --title "Fix: <failure>" \
    --source-task-id <TASK_ID> \
    --block-source \
    --description "Dispatched task <TASK_ID> was auto-downgraded to blocked — test failures or record issues"
  ```
  `forge task add` automatically deduplicates — check output:
  - `ACTION: ADDED` → new fix task created
  - `ACTION: SKIPPED` → active fix task already exists
  Continue loop.
- **STATUS == "in_progress"** (no record created): proceed to 2c.

**2c. Record-Missing Recovery** — `Agent(subagent_type="forge:task-executor", prompt="Fix record for task <TASK_ID>")`. Subagent detects "Fix record for" prefix and calls `forge prompt get-by-task-id <TASK_ID> --fix-record-missed` internally. After 2c, re-verify via 2b logic.

### Step 3: Continue Loop

Return to Step 1.

## Error Handling

| Situation | Action |
|-----------|--------|
| No available task | End loop, print summary |
| Agent timeout | Mark blocked, continue |
| Record missing | Dispatch fix-record subagent (2c) |
| 3 consecutive failures | STOP |
| Main session fails | Follow task doc's error section; if missing, `forge task add --template fix-task --title "Fix: main session task failed" --source-task-id <TASK_ID> --block-source --description "Main session task failed"` then continue |

## Post-Completion

After loop ends, print: "All tasks completed. T-test-run, T-test-graduate, and T-test-verify-regression handle e2e verification, graduation, and regression automatically."

If index lacks T-test-run/T-test-graduate, suggest: "Run `/run-tests` then `forge test promote <journey>`."

Do NOT run e2e tests from the dispatcher.

Do NOT commit post-loop artifacts. The `forge feature complete --if-done` Stop hook detects uncommitted artifacts and blocks the agent to commit them via `/git-commit`. This ensures artifacts are committed only after quality-gate passes.
