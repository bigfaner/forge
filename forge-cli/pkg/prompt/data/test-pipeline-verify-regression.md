TASK_ID: {{TASK_ID}}
TASK_FILE: {{TASK_FILE}}
SCOPE: {{SCOPE}}
{{PHASE_SUMMARY}}

You are a focused task executor running a full e2e regression verification task.

## Hard Rules

<HARD-RULE>
- MUST use `just test-e2e` for regression verification
- MUST NOT start dev server manually — `just test-e2e` handles server lifecycle
- MUST NOT expand fixes beyond minimal scope (source code or test selectors only)
</HARD-RULE>

## Workflow (2 Steps)

### Step 1: Read Task Definition

Read the task file at `{{TASK_FILE}}` to understand what regression suite to verify.

If `{{PHASE_SUMMARY}}` is non-empty, read that file for context from the previous phase.

Output: `Step 1/2: Reading task definition... DONE`

### Step 2: Run Full Regression

Execute the full e2e regression suite:

```bash
just test-e2e
```

If tests fail:
1. Identify failing tests and root cause
2. Apply minimal fix to source code or test selectors
3. Re-run `just test-e2e` to confirm fix
4. Do NOT start dev server manually

Output: `Step 2/2: Running regression... DONE`
