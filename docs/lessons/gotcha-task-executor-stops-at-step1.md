# forge:task-executor Stops After "Step 1/5" — Escalation Protocol

## Problem

`forge:task-executor` subagent returns `"Step 1/5: Reading task definition... DONE"` as its only output and stops. No implementation file is created, no record is written. Retrying the same subagent type produces the same result.

Observed pattern:
- Attempt 1: 19 tools, 322s → "Step 1/5: Reading task definition... DONE" → no file
- Attempt 2: 14 tools, 240s → same output → no file
- Attempt 3: 4 tools, 138s → same output → no file (tool count drops each retry — agent is doing less each time)

## Root Cause

Three-level causal chain:

1. **Surface**: `forge:task-executor` outputs the step marker and stops
2. **Direct cause**: The skill's step-by-step format (`Step 1/5 ... Step 2/5 ...`) causes the agent to treat each step as an invocation boundary — it outputs "Step 1 done" and considers its turn complete
3. **Root cause**: The skill prompt structure is ambiguous about whether steps are a checklist to output or actions to perform in sequence. The agent resolves the ambiguity by outputting the marker and stopping, especially when the task requires writing new files (higher-risk action)

Secondary cause: The dispatcher retried the same failing subagent type 3+ times. Each retry had fewer tool uses (19 → 14 → 4), indicating the agent was doing progressively less — a sign of context degradation or the agent "giving up" faster each time.

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

**When `forge:task-executor` returns only "Step 1/5: Reading task definition... DONE", it will not recover on retry.** The dispatcher's escalation path is:

```
forge:task-executor fails once
  → retry once with more explicit prompt
  → still fails → implement directly in main session
```

Do NOT retry the same subagent type 3 times. The tool use count dropping (19 → 14 → 4) is a signal that retries are making things worse, not better.

Also: a general-purpose agent with a vague "implement based on your findings" prompt will also fail. The agent needs a concrete spec, not a research task.
