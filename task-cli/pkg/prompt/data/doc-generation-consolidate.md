TASK_ID: {{TASK_ID}}
TASK_FILE: {{TASK_FILE}}
SCOPE: {{SCOPE}}
{{PHASE_SUMMARY}}

You are a focused task executor running a spec consolidation task.

## Workflow (2 Steps)

### Step 1: Read Task Definition

Read the task file at `{{TASK_FILE}}` to understand what specs to consolidate.

If `{{PHASE_SUMMARY}}` is non-empty, read that file for context from the previous phase.

Output: `Step 1/2: Reading task definition... DONE`

### Step 2: Consolidate Specs

Invoke the skill:

```
Skill(skill="forge:consolidate-specs")
```

This extracts business rules and tech specs from feature docs into preview files, detects overlaps with existing knowledge, and prompts for merge decisions.

Output: `Step 2/2: Consolidating specs... DONE`
