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
   ┌─────────────┐    ┌─────────────┐    ┌─────────────┐    ┌─────────────┐
   │ 1. Claim    │───▶│ 2. Dispatch │───▶│ 3. Verify   │───▶│ 4. Context  │
   │    Task     │    │   + Timeout │    │   Record    │    │   Check     │
   └─────────────┘    └─────────────┘    └─────────────┘    └─────────────┘
          ▲                                    │                    │
          └────────────────────────────────────┘                    │
                      LOOP                                        │
                                                                   ▼
                                                        ┌─────────────┐
                                                        │ 5. Breaking │
                                                        │    Gate     │
                                                        └─────────────┘
```

## Dispatcher Iron Laws

<EXTREMELY-IMPORTANT>
1. Only 5 actions: claim → dispatch → verify → context check → breaking gate
2. NO code reading, NO code writing
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
- `SCOPE` (e.g., "frontend", "backend", or "all" — defaults to "all" if absent)
- `FEATURE` (e.g., "my-feature" — feature slug from claim output)

### Step 2: Dispatch with Timeout

Determine if phase summary context should be injected:

**Phase boundary detection**:
1. From `TASK_ID`, extract the phase number (integer before first dot, e.g., "2.1" → phase 2)
2. If this is the first task of a new phase (phase number changed from previous task):
   - Check if previous phase's summary record exists: `docs/features/<slug>/tasks/records/<prev-phase>-summary.md`
   - If exists, include in dispatch prompt
   - Skip injection for gate tasks (ID ends with `.gate`) and phase summary tasks (ID ends with `.summary`)

```
Agent(
  subagent_type="task-executor",
  prompt="TASK_KEY: {{KEY}}
TASK_ID: {{TASK_ID}}
TASK_FILE: {{FILE}}
SCOPE: {{SCOPE}}
{{PHASE_SUMMARY_SECTION}}

IMPORTANT: Do NOT claim or start any other tasks after completing this one. Stop after recording the task result."
)
```

Where `{{PHASE_SUMMARY_SECTION}}` is:
```
PHASE_SUMMARY: Read the following file for context from previous phases:
docs/features/<slug>/tasks/records/<prev-phase>-summary.md
```

Phase summaries follow a fixed 5-section structure (Tasks Completed, Key Decisions, Types & Interfaces Changed, Conventions Established, Deviations from Design). The task-executor is trained to parse these sections.

Or empty string if no phase summary exists.

**Timeout**: 30 minutes

### Step 3: Verify Record

Check if record file exists after agent completes.

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
- Run `task template fix-task` to view the template, then add fix task and mark source blocked:
  ```bash
  task status <TASK_ID> blocked
  task add --template fix-task --title "Fix: <failure>" \
    --source-task-id <TASK_ID> \
    --var SOURCE_FILES="<affected paths>" \
    --var TEST_SCRIPT="<failing test>" \
    --var TEST_RESULTS="<results path>" \
    --description "<root cause>"
  ```
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
- Mark source task blocked, then add fix task using the fix-task template:
  ```bash
  task status <TASK_ID> blocked
  task add --template fix-task --title "Fix: <concise description>" \
    --source-task-id <TASK_ID> \
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
| Record missing | Dispatch error-fixer (include: "Use /record-task skill to create record") |
| 3 consecutive failures | STOP dispatcher |
| Breaking task tests fail (5a) | `task status <ID> blocked` + `task add --template fix-task`, continue loop |
| Feature e2e tests fail (5b) | `task status <ID> blocked` + `task add --template fix-task`, continue loop |

### Error-Fixer Dispatch

When dispatching error-fixer for missing record, include explicit instruction:

```
Agent(
  subagent_type="error-fixer",
  prompt="TASK_ID: {{TASK_ID}}
TASK_FILE: {{FILE}}
ERROR_MESSAGES: Missing task record
INSTRUCTION: Use /record-task skill to create the record (task record CLI is mandatory)"
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
