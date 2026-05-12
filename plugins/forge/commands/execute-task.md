---
name: execute-task
description: Execute single task with focused TDD workflow.
allowed_tools: ["Bash", "Read", "Agent", "TaskOutput", "Skill"]
---

# /execute-task

Execute a single task. MAIN_SESSION tasks execute in main session; all others dispatch to forge:task-executor subagent (which calls `task prompt` internally).

## Step 1: Claim Task

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
- `FEATURE` (e.g., "my-feature" — feature slug from claim output)

## Step 1.5: Main Session Routing

If `MAIN_SESSION == "true"`:

1. Read the task file at the FILE path extracted from claim output and find the `## Main Session Instructions` section.
2. Follow the instructions exactly — the task document specifies what skill to invoke, how to check outcome, and how to record the result.
3. The dispatcher does NOT hardcode skill names or record logic — it delegates to the task document.
4. If the task file lacks a `## Main Session Instructions` section, mark the task blocked and report: "MAIN_SESSION task missing Main Session Instructions section — task document is incomplete".
5. After execution, verify the record file exists via `task query <TASK_ID>`. If STATUS is not `"completed"`, spawn fix task (same as Step 2 verify logic).
6. Skip to Step 4 (STOP).

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

The subagent internally runs `task prompt <TASK_ID>` to get the execution strategy.

**Timeout**: 30 minutes

### 2b. Verify Record

After subagent returns, check the task's actual status via CLI:

```bash
task query <TASK_ID>
```

- **STATUS == `"completed"`**: proceed to Step 3 (Breaking Gate).
- **STATUS != `"completed"`**: task was auto-downgraded (e.g. test failures).
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

### 2c. Record-Missing Recovery

When the subagent completes but the record file is missing, delegate recovery to another subagent:

```
Agent(
  subagent_type="forge:task-executor",
  prompt="Fix record for task <TASK_ID>"
)
```

The subagent's Execution Protocol detects the "Fix record for" prefix and calls `task prompt <TASK_ID> --fix-record-missed` internally.

## Step 3: Breaking Task Gate

Determine which gates to run based on claim output from Step 1:

| BREAKING=true? | SCOPE frontend\|all + specs exist? | Run 3a? | Run 3b? |
|----------------|-------------------------------------|---------|---------|
| Yes | No | Yes | No |
| No | Yes | No | Yes |
| Yes | Yes | Yes | Yes |
| No | No | Skip Step 3 entirely | Skip Step 3 entirely |

If running both: execute 3a first. Only proceed to 3b if 3a passes.

### 3a. Unit/Integration Gate (BREAKING: true)

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

**If tests pass**: if the routing table indicates 3b should also run (SCOPE frontend|all + specs exist), proceed to 3b. Otherwise proceed to Step 4 (STOP).

### 3b. Feature E2E Gate (SCOPE=frontend|all, specs exist)

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

## Step 4: STOP

<HARD-RULE>
ONE TASK PER INVOCATION. This is absolute and non-negotiable.

After Step 3, you MUST stop immediately.

<PROHIBITIONS>
- Running `task claim` under any circumstances
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
| Record missing | Dispatch `Agent(prompt="Fix record for task <TASK_ID>")` — subagent calls `task prompt --fix-record-missed` internally |
| Breaking task tests fail (3a) | `task add --template fix-task --block-source` |
| Feature e2e tests fail (3b) | `task add --template fix-task --block-source` |
| Main session task fails | Follow error handling in task document's `### Error Handling` section; if missing, `task add --template fix-task --block-source` |

## Rules

<EXTREMELY-IMPORTANT>
- record-task is mandatory — No completion without it
- All verifications must pass
- ONE TASK PER INVOCATION — after Step 3, STOP immediately, no exceptions
- FORBIDDEN: run "task claim", read index.json, or start any subsequent task
- Do NOT use TASK_FILE parameter when dispatching to forge:task-executor
</EXTREMELY-IMPORTANT>

## Related Commands

| Command | Usage |
|---------|-------|
| `/run-tasks` | Auto-execute all tasks |
| `/record-task` | Create record + update status |
