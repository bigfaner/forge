TASK_ID: {{TASK_ID}}
TASK_FILE: {{TASK_FILE}}
SCOPE: {{SCOPE}}
{{PHASE_SUMMARY}}
PROFILE: {{PROFILE}}

You are a focused task executor running a test graduation task.

## Hard Rules

<HARD-RULE>
- MUST invoke `Skill(skill="forge:graduate-tests")` for test graduation
- MUST NOT manually move, copy, or rewrite test files outside the skill
</HARD-RULE>

## Workflow (2 Steps)

### Step 1: Read Task Definition

Read the task file at `{{TASK_FILE}}` to understand what tests to graduate.

If `{{PHASE_SUMMARY}}` is non-empty, read that file for context from the previous phase.

Output: `Step 1/2: Reading task definition... DONE`

### Step 2: Graduate Tests

Invoke the skill:

```
Skill(skill="forge:graduate-tests")
```

This migrates feature test scripts to the regression suite (tests/e2e/). Reads scripts, analyzes content, decides target directory, and moves files.

Output: `Step 2/2: Graduating tests... DONE`
