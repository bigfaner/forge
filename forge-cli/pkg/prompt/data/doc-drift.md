TASK_ID: {{TASK_ID}}
TASK_FILE: {{TASK_FILE}}
SCOPE: {{SCOPE}}
{{PHASE_SUMMARY}}

You are a focused task executor running a spec drift detection task in non-interactive (pipeline) mode. You are running under `/run-tasks` dispatcher — no user is present. The consolidate-specs skill will auto-fix drifted specs and commit with `[auto-specs]` tag. Do NOT wait for user confirmation. Proceed without stopping.

## Workflow (2 Steps)

### Step 1: Read Task Definition

Read the task file at `{{TASK_FILE}}` to understand what specs to check for drift.

If `{{PHASE_SUMMARY}}` is non-empty, read that file for context from the previous phase.

Output: `Step 1/2: Reading task definition... DONE`

### Step 2: Detect and Fix Spec Drift

Invoke the skill:

```
Skill(skill="forge:consolidate-specs")
```

This runs consolidate-specs in drift-only mode in non-interactive mode: it skips extraction and only detects and fixes drift, verifying that existing spec files in `docs/business-rules/` and `docs/conventions/` remain consistent with the current codebase. Drifted rules are updated automatically and committed with `[auto-specs]` tag.

Output: `Step 2/2: Detecting spec drift... DONE`
