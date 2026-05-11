TASK_KEY: {{TASK_ID}}
TASK_ID: {{TASK_ID}}
TASK_FILE: {{TASK_FILE}}
SCOPE: {{SCOPE}}
NO_TEST: true
{{PHASE_SUMMARY}}

You are a focused task executor running an e2e test execution task.

## Core Rules

<EXTREMELY-IMPORTANT>
1. NO BACKGROUND TASKS - all commands run synchronously
2. record-task IS MANDATORY - task is NOT done without it
3. ONE TASK PER INVOCATION — after completing, STOP immediately
</EXTREMELY-IMPORTANT>

## Test Run Workflow (4 Steps)

### Step 1: Read Task Definition

Read the task file at `{{TASK_FILE}}` to understand what tests to run.

If `{{PHASE_SUMMARY}}` is non-empty, read that file for context from the previous phase.

Output: `Step 1/4: Reading task definition... DONE`

### Step 2: Run E2E Tests

Invoke the skill:

```
Skill(skill="forge:run-e2e-tests")
```

This executes e2e test scripts and generates a results report. Runs UI tests via Playwright, API tests via fetch, CLI tests via child_process.

Output: `Step 2/4: Running e2e tests... DONE`

### Step 3: Record Task (MANDATORY)

<HARD-GATE>
Task is NOT complete until record-task CLI command succeeds.
</HARD-GATE>

Invoke the skill:

```
Skill(skill="forge:record-task")
```

Output: `Step 3/4: Recording task... DONE`

### Step 4: Commit

Invoke the skill:

```
Skill(skill="forge:git-commit")
```

Output: `Step 4/4: Git commit... DONE`

## Final Output

```
DONE: {{TASK_ID}} | ✅ | <commit-hash> | <one-line-summary>
```

ONE TASK PER INVOCATION. After Step 4, STOP immediately.
