---
name: run-tasks
description: Autonomous task dispatcher that continuously claims tasks and dispatches to subagents.
allowed_tools: ["Bash", "Read", "Agent", "TaskOutput", "Skill"]
---

# /run-tasks

Auto-dispatch tasks. MAIN_SESSION tasks execute in main session; eval-cases type executes in main session; all others go to forge:task-executor subagent.

## Architecture

```mermaid
flowchart TD
    A["1. Claim Task"] --> B{"MAIN_SESSION?"}
    B -->|"yes"| C["1.5 Follow Task Instructions"]
    C --> LOOP(["Step 6: Continue Loop"])
    B -->|"no"| D["2. task prompt → Dispatch"]
    D --> D1{"prompt exit != 0?"}
    D1 -->|"yes"| D2["block task, continue"]
    D2 --> LOOP
    D1 -->|"no"| D3{"TYPE == eval-cases?"}
    D3 -->|"yes"| D4["Execute in main session"]
    D3 -->|"no"| D5["Agent(forge:task-executor)"]
    D4 --> E["3. Verify Record"]
    D5 --> E
    E --> F["4. Context Check"]
    F --> G{"Breaking task?"}
    G -->|"yes"| H["5. Breaking Gate"]
    G -->|"no"| LOOP
    H --> LOOP
    LOOP --> A
```

## Dispatcher Iron Laws

<EXTREMELY-IMPORTANT>
1. Only 5 actions: claim → (main_session? follow task instructions : dispatch) → verify → context check → breaking gate
2. NO code reading, NO code writing — EXCEPT for MAIN_SESSION tasks (Step 1.5) where the Skill tool is invoked in the main session
3. NO running tests directly — EXCEPT in Step 5 (Breaking Task Gate) where `just test` and `just test-e2e` are executed as quality gates
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

**Extract from claim output**:
- `TASK_ID` (e.g., "2.1")
- `KEY` (e.g., "2.1-implementation")
- `FILE` (e.g., full absolute path to task file)
- `BREAKING` (e.g., "true" or absent)
- `MAIN_SESSION` (e.g., "true" or absent)
- `SCOPE` (e.g., "frontend", "backend", or "all" — defaults to "all" if absent; may be omitted entirely by claim output when not set)
- `NO_TEST` (e.g., "true" or "false")
- `FEATURE` (e.g., "my-feature" — feature slug from claim output)
- `TYPE` (e.g., "test-pipeline.eval-cases" — task type from claim output; may be absent)

### Step 1.5: Main Session Routing

If `MAIN_SESSION == "true"`:

1. Read the task file at `{{FILE}}` and find the `## Main Session Instructions` section.
2. Follow the instructions exactly — the task document specifies what skill to invoke, how to check outcome, and how to record the result.
3. The dispatcher does NOT hardcode skill names or record logic — it delegates to the task document.
4. If the task file lacks a `## Main Session Instructions` section, mark the task blocked and report: "MAIN_SESSION task missing Main Session Instructions section — task document is incomplete".
5. After execution, verify the record file exists (same as Step 3 for subagent tasks).
6. Skip to Step 6 (Continue Loop).

Else:
- Proceed to Step 2 (Dispatch with Timeout).

### Step 2: Dispatch with Timeout

**Synthesize prompt via CLI**:

```bash
SYNTHESIZED_PROMPT=$(task prompt <TASK_ID> 2>/tmp/prompt_error.txt)
PROMPT_EXIT=$?
PROMPT_ERROR=$(cat /tmp/prompt_error.txt)
```

**If `PROMPT_EXIT != 0`**:
```bash
task status <KEY> blocked --reason "$PROMPT_ERROR"
```
Then continue loop (skip Steps 3–5 for this iteration).

**If `PROMPT_EXIT == 0`**:

Check `TYPE` extracted from Step 1 claim output:

**If `TYPE == "test-pipeline.eval-cases"`**:
- Execute `SYNTHESIZED_PROMPT` directly in the main session (do NOT dispatch to subagent).
- Proceed to Step 3.

**All other types**:
```
Agent(
  subagent_type="forge:task-executor",
  prompt=SYNTHESIZED_PROMPT
)
```
Proceed to Step 3.

**Timeout**: 30 minutes

### Step 3: Verify Record

Check if record file exists after agent completes. Then check the task's actual status via CLI:

```bash
task query <TASK_ID>
```

- If STATUS is not `"completed"`: task was auto-downgraded (e.g. test failures).
  Spawn fix task using `--block-source` to atomically block the source:
  ```bash
  task add --template fix-task --title "Fix: <failure>" \
    --source-task-id <TASK_ID> \
    --block-source \
    --description "<reason>"
  ```
  `task add` automatically deduplicates — check output:
  - `ACTION: ADDED` → new fix task created, continue loop
  - `ACTION: SKIPPED` → active fix task already exists, continue loop
- Only proceed if STATUS is `"completed"`

### Step 4: Context Check

After verifying the record, check if the completed task was a phase summary task (ID ends with `.summary`):
- If yes: This phase's summary is now available for subsequent phases
- No additional action needed — the summary will be injected on next phase boundary

### Step 5: Breaking Task Gate

Determine which gates to run based on claim output from Step 1:

| BREAKING=true? | SCOPE frontend\|all + specs exist? | Run 5a? | Run 5b? |
|----------------|-------------------------------------|---------|---------|
| Yes | No | Yes | No |
| No | Yes | No | Yes |
| Yes | Yes | Yes | Yes |
| No | No | Skip Step 5 entirely | Skip Step 5 entirely |

If running both: execute 5a first. Only proceed to 5b if 5a passes.

#### 5a. Unit/Integration Gate (BREAKING: true)

```bash
# Pre-flight: verify justfile and test recipe exist
if [ ! -f justfile ] && [ ! -f Justfile ]; then
    echo "Error: justfile not found — run /init-justfile first" >&2
    exit 1
fi
just --list 2>/dev/null | grep -q "^    test " || {
    echo "Error: 'test' recipe not found in justfile" >&2
    exit 1
}
```

```bash
just test [scope]
```

Apply the **Scope Resolution** protocol from the Forge Guide — use the `SCOPE` extracted from the claim output in Step 1.

**If tests fail**:
- Run `task template fix-task` to view the template, then add fix task:
  ```bash
  task add --template fix-task --title "Fix: <failure>" \
    --source-task-id <TASK_ID> \
    --block-source \
    --var SOURCE_FILES="<affected paths>" \
    --var TEST_SCRIPT="<failing test>" \
    --var TEST_RESULTS="<results path>" \
    --description "<root cause>"
  ```
  **`--block-source`**: atomically sets source task to blocked before resolution, preserving the fix-chain model.
  **`--source-task-id` auto-resolves**: if `<TASK_ID>` is a **completed** fix-task, the CLI automatically resolves to the root blocked task. Always pass the current failing task's ID — no manual chain tracing needed.
- Continue loop — fix task (P0) will be claimed on next iteration
- Do NOT proceed to next task until fix task completes

**If tests pass**: if the routing table indicates 5b should also run (SCOPE frontend|all + specs exist), proceed to 5b. Otherwise continue to next iteration (Step 1).

#### 5b. Feature E2E Gate (SCOPE=frontend|all, specs exist)

<EXTREMELY-IMPORTANT>
The dispatcher evaluates SCOPE and FEATURE from Step 1 claim output BEFORE executing any bash commands below. If SCOPE is `backend` or FEATURE is empty, skip this entire section.
</EXTREMELY-IMPORTANT>

Pre-conditions (all must be true):
- SCOPE is `frontend` or `all` (defaults to "all" if absent from claim output)
- FEATURE is non-empty (always true after successful claim)
- Feature has e2e spec files: `tests/e2e/features/$FEATURE/` contains `.spec.ts` files
- `test-e2e` recipe exists in justfile

```bash
# Pre-flight: verify test-e2e recipe exists — if missing, skip to next iteration
SKIP=""
just --list 2>/dev/null | grep -q "test-e2e" || { echo "Skip: test-e2e recipe not found"; SKIP=true; }

# Check if specs exist for this feature
if [ -z "$(ls "tests/e2e/features/$FEATURE/"*.spec.ts 2>/dev/null)" ]; then
    echo "Skip: no .spec.ts files in tests/e2e/features/$FEATURE/"
    SKIP=true
fi

# If pre-flights passed, run e2e
if [ -z "$SKIP" ]; then
    just e2e-setup
    just test-e2e --feature "$FEATURE"
fi
```

**If e2e fails**:
- Add fix task using the fix-task template:
  ```bash
  task add --template fix-task --title "Fix: <concise description>" \
    --source-task-id <TASK_ID> \
    --block-source \
    --var SOURCE_FILES="<affected source paths>" \
    --var TEST_SCRIPT="tests/e2e/features/$FEATURE/<failing-spec>.spec.ts" \
    --var TEST_RESULTS="tests/e2e/features/$FEATURE/results/latest.md" \
    --description "<root cause and context>"
  ```

**If e2e passes or pre-flight skipped**: continue to next iteration (Step 1)

### Step 6: Continue Loop

Return to Step 1.

## Error Handling

| Situation | Action |
|-----------|--------|
| No available task | End loop, print summary |
| Agent timeout | Mark blocked, continue next |
| `task prompt` exit != 0 | Write stderr to blockedReason, `task status <KEY> blocked`, continue loop |
| Record missing | Run `task prompt <TASK_ID> --fix-record-missed` → `Agent(forge:task-executor, prompt=<stdout>)` |
| 3 consecutive failures | STOP dispatcher |
| Breaking task tests fail (5a) | `task add --template fix-task --block-source`, continue loop |
| Feature e2e tests fail (5b) | `task add --template fix-task --block-source`, continue loop |
| Main session task fails | Follow error handling in task document's `### Error Handling` section; if missing, `task add --template fix-task --block-source`, continue loop |

### Record-Missing Recovery

When a task agent completes but the record file is missing, recover using `task prompt` with the fix flag:

```bash
RECOVERY_PROMPT=$(task prompt <TASK_ID> --fix-record-missed)
```

Then dispatch:
```
Agent(
  subagent_type="forge:task-executor",
  prompt=RECOVERY_PROMPT
)
```

## Post-Completion

After all tasks are completed (loop ends with "No available task"):

```
Print summary to user:
"All tasks completed. T-test-3, T-test-4, and T-test-4.5 in the task chain handle
e2e verification, graduation, and regression automatically."
```

If the feature's task index does not include T-test-3/T-test-4, suggest:
```
"Run `/run-e2e-tests` to verify against PRD acceptance criteria,
then `/graduate-tests` to migrate scripts to the regression suite."
```

Do NOT run e2e tests outside of the Breaking Task Gate (Step 5) — the dispatcher only executes tests as quality gates.

## Related Commands

| Command | Usage |
|---------|-------|
| `/execute-task` | Manual single task |
| `/record-task` | Create record + update status |
| `/run-e2e-tests` | E2E verification against PRD |
