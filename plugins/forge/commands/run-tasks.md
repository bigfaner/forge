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

**Extract from claim output**:
- `TASK_ID` (e.g., "2.1")
- `TASK_KEY` (e.g., "2.1-implementation")
- `TASK_FILE` (e.g., "2.1-implementation.md")
- `BREAKING` (e.g., "true" or absent)

### Step 2: Dispatch with Timeout

Determine if phase summary context should be injected:

**Phase boundary detection**:
1. From `TASK_ID`, extract the phase number (integer before first dot, e.g., "2.1" → phase 2)
2. If this is the first task of a new phase (phase number changed from previous task):
   - Check if previous phase's summary record exists: `docs/features/<slug>/tasks/records/<prev-phase>.summary-phase-summary.md`
   - If exists, include in dispatch prompt
   - Skip injection for gate tasks (ID ends with `.gate`) and phase summary tasks (ID ends with `.summary`)

```
Agent(
  subagent_type="task-executor",
  prompt="TASK_KEY: {{KEY}}
TASK_ID: {{ID}}
TASK_FILE: {{FILE}}
{{PHASE_SUMMARY_SECTION}}

IMPORTANT: Do NOT claim or start any other tasks after completing this one. Stop after recording the task result."
)
```

Where `{{PHASE_SUMMARY_SECTION}}` is:
```
PHASE_SUMMARY: Read the following file for context from previous phases:
docs/features/<slug>/tasks/records/<prev-phase>.summary-phase-summary.md
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

If the claimed task had `BREAKING: true` in the claim output:

```bash
# Run project-level full test suite
# Use the same detection logic as `task all-completed`:
# go test ./... OR npm test OR the testCommand from index.json
```

**If tests fail**:
- Option A: Dispatch error-fixer with failure context (existing behavior)
- Option B: Add fix task via `task add --title "Fix: <failure>" --priority P0 --breaking --description "..."` and continue loop
- Do NOT proceed to next task until error-fixer resolves or fix task completes

**If tests pass**:
- Continue to next iteration

### Step 6: Continue Loop

Return to Step 1.

## Error Handling

| Situation | Action |
|-----------|--------|
| No available task | End loop, print summary |
| Agent timeout | Mark blocked, continue next |
| Record missing | Dispatch error-fixer (include: "Use /record-task skill to create record") |
| 3 consecutive failures | STOP dispatcher |
| Breaking task tests fail | Dispatch error-fixer with failure details |

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

## Post-Completion: E2E Verification

After all tasks are completed (loop ends with "No available task"):

```
Suggest to user:
"All tasks completed. Run `/run-e2e-tests` to verify against PRD acceptance criteria."
```

Do NOT run e2e tests automatically — the dispatcher must not execute tests. Only suggest.

## Related Commands

| Command | Usage |
|---------|-------|
| `/execute-task` | Manual single task |
| `/record-task` | Create record + update status |
| `/run-e2e-tests` | E2E verification against PRD |
