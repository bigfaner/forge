---
name: execute-task
description: Execute single task with focused TDD workflow.
allowed_tools: ["Bash", "Read", "Agent", "TaskOutput", "Skill"]
---

# /execute-task

Execute a single task. MAIN_SESSION tasks execute in main session; eval-cases type executes in main session; all others dispatch to forge:task-executor subagent.

## Step 0: MAIN_SESSION Check

If the claimed task has `MAIN_SESSION == "true"`, read the task file's `## Main Session Instructions` section and follow it. After execution, invoke `Skill(skill="record-task")`. Skip Steps 2–4.

## Step 1: Claim & Read

```bash
task claim
```

**Output parsing**:
- `ACTION: CLAIMED` → New task
- `ACTION: CONTINUE` → Resume existing task
- Error → No available task, stop

**Extract from claim output**:
- `TASK_ID` (e.g., "2.1")
- `KEY` (e.g., "2.1-implementation")
- `FILE` (e.g., full absolute path to task file)
- `BREAKING` (e.g., "true" or absent)
- `MAIN_SESSION` (e.g., "true" or absent)
- `SCOPE` (e.g., "frontend", "backend", or "all" — defaults to "all" if absent)
- `NO_TEST` (e.g., "true" or "false")
- `FEATURE` (e.g., "my-feature" — feature slug from claim output)
- `TYPE` (e.g., "test-pipeline.eval-cases" — task type from claim output; may be absent)

## Step 2: Dispatch

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
Stop — do not proceed to Steps 3–4.

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

## Step 3: Verify Record

Check the task's actual status via CLI:

```bash
task query <TASK_ID>
```

- If STATUS is `"completed"`: proceed normally.
- If STATUS is not `"completed"`: task was auto-downgraded (e.g. test failures).
  Spawn fix task using `--block-source` to atomically block the source:
  ```bash
  task add --template fix-task --title "Fix: <failure>" \
    --source-task-id <TASK_ID> \
    --block-source \
    --description "<reason>"
  ```
  `task add` automatically deduplicates — check output:
  - `ACTION: ADDED` → new fix task created
  - `ACTION: SKIPPED` → active fix task already exists

### Record-Missing Recovery

When the agent completes but the record file is missing, recover using `task prompt` with the fix flag:

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

## Step 4: Breaking Task Gate

Determine which gates to run based on claim output from Step 1:

| BREAKING=true? | SCOPE frontend\|all + specs exist? | Run 4a? | Run 4b? |
|----------------|-------------------------------------|---------|---------|
| Yes | No | Yes | No |
| No | Yes | No | Yes |
| Yes | Yes | Yes | Yes |
| No | No | Skip Step 4 entirely | Skip Step 4 entirely |

If running both: execute 4a first. Only proceed to 4b if 4a passes.

### 4a. Unit/Integration Gate (BREAKING: true)

```bash
just test [scope]
```

Apply the **Scope Resolution** protocol from the Forge Guide — use the `SCOPE` extracted from the claim output in Step 1.

**If tests fail**:
```bash
task add --template fix-task --title "Fix: <failure>" \
  --source-task-id <TASK_ID> \
  --block-source \
  --var SOURCE_FILES="<affected paths>" \
  --var TEST_SCRIPT="<failing test>" \
  --var TEST_RESULTS="<results path>" \
  --description "<root cause>"
```

### 4b. Feature E2E Gate (SCOPE=frontend|all, specs exist)

Pre-conditions (all must be true):
- SCOPE is `frontend` or `all`
- FEATURE is non-empty
- Feature has e2e spec files: `tests/e2e/features/$FEATURE/` contains `.spec.ts` files
- `test-e2e` recipe exists in justfile

```bash
just e2e-setup
just test-e2e --feature "$FEATURE"
```

**If e2e fails**:
```bash
task add --template fix-task --title "Fix: <concise description>" \
  --source-task-id <TASK_ID> \
  --block-source \
  --var SOURCE_FILES="<affected source paths>" \
  --var TEST_SCRIPT="tests/e2e/features/$FEATURE/<failing-spec>.spec.ts" \
  --var TEST_RESULTS="tests/e2e/features/$FEATURE/results/latest.md" \
  --description "<root cause and context>"
```

## Error Handling

| Situation | Action |
|-----------|--------|
| No available task | Stop, report |
| Agent timeout | Mark blocked, stop |
| `task prompt` exit != 0 | Write stderr to blockedReason, `task status <KEY> blocked`, stop |
| Record missing | Run `task prompt <TASK_ID> --fix-record-missed` → `Agent(forge:task-executor, prompt=<stdout>)` |
| Breaking task tests fail (4a) | `task add --template fix-task --block-source` |
| Feature e2e tests fail (4b) | `task add --template fix-task --block-source` |
| Main session task fails | Follow error handling in task document's `### Error Handling` section; if missing, `task add --template fix-task --block-source` |

## Rules

<EXTREMELY-IMPORTANT>
- record-task is mandatory — No completion without it
- All verifications must pass
- ONE TASK PER INVOCATION — after Step 4, STOP immediately, no exceptions
- FORBIDDEN: run "task claim", read index.json, or start any subsequent task
- Do NOT use TASK_FILE or NO_TEST parameters when dispatching to forge:task-executor
</EXTREMELY-IMPORTANT>

## STOP

<HARD-RULE>
ONE TASK PER INVOCATION. This is absolute and non-negotiable.

After Step 4, you MUST stop immediately.

<PROHIBITIONS>
- Running `task claim` under any circumstances
- Reading the next task file
- Continuing with any additional work
</PROHIBITIONS>

Output your final summary and STOP.
</HARD-RULE>

## Related Commands

| Command | Usage |
|---------|-------|
| `/run-tasks` | Auto-execute all tasks |
| `/record-task` | Create record + update status |
