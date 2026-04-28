# forge:task-executor Stops After "Step 1/5" — Escalation Protocol

## Problem

`forge:task-executor` subagent returns `"Step 1/5: Reading task definition... DONE"` as its only output and stops. No implementation file is created, no record is written.

Observed pattern (manual retries):
- Attempt 1: 19 tools, 322s → "Step 1/5: Reading task definition... DONE" → no file
- Attempt 2: 14 tools, 240s → same output → no file
- Attempt 3: 4 tools, 138s → same output → no file (tool count drops each retry)

## Root Cause (Uncertain)

The exact root cause is unknown. What we know:

1. **The step format itself is not the cause.** `Step N/5: <name> DONE` is an output format rule ("when done, output this marker"), not an execution boundary. The agent is supposed to do the work AND output the marker.
2. **The agent is choosing not to proceed past Step 1.** Possible reasons: task file content is ambiguous, the implementation looks high-risk (new files vs. edits), context window pressure, or permission friction.
3. **Tool count dropping (19→14→4) suggests the agent adapts its behavior across retries** — doing progressively less each invocation, consistent with LLM context degradation.

Note: The `run-tasks` dispatcher does NOT retry the same task with `task-executor` — it escalates to `error-fixer` when a record is missing. The retry pattern above was from manual re-dispatch by the operator.

## Solution

**After 1 failed `forge:task-executor` attempt, do not retry the same type.** Escalate immediately:

### Option A: Direct implementation in main session

Read the task file and context files directly, implement the code, run tests, write record.json, call `task record`. This is always available and doesn't depend on subagent behavior.

### Option B: General-purpose agent with concrete spec

Replace the skill-based prompt with a fully explicit prompt that:
- Names every file to create with exact paths
- Specifies the exact code structure (structs, function signatures)
- Lists every bash command to run in order
- Does NOT say "based on what you find" — provide the spec directly

Bad prompt: "Read the task file and implement what's needed"
Good prompt: "Create file X with struct Y and function Z. Run command W. Write record.json with these exact fields. Run `task record`."

## Key Takeaway

**When `forge:task-executor` returns only "Step 1/5: Reading task definition... DONE", retrying the same subagent type will not fix it.** The escalation path for manual dispatch:

```
forge:task-executor fails once
  → do not retry same type
  → implement directly (Option A) or use general-purpose agent with concrete spec (Option B)
```

A general-purpose agent with a vague "implement based on your findings" prompt will also fail. The agent needs a concrete spec, not a research task.
