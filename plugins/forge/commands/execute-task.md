---
name: execute-task
description: Execute single task with focused TDD workflow.
allowed_tools: ["Bash", "Read", "Write", "Edit", "Grep", "Glob", "Agent", "LSP"]
---

# /execute-task

Execute a single task with streamlined TDD workflow.

## Step 0: MAIN_SESSION Check

If the claimed task has `MAIN_SESSION == "true"`, read the task file's `## Main Session Instructions` section and follow it. After execution, invoke `Skill(skill="record-task")`. Skip Steps 2–5.

## Workflow (5 Steps)

```
Step 1: Read task definition
Step 2: TDD (RED → GREEN → REFACTOR)
Step 3: Full verification
Step 4: Record task (MANDATORY)
Step 5: Git commit
```

## Step 1: Claim & Read

```bash
task claim
```

Parse output for KEY, TASK_ID, FILE, SCOPE, FEATURE. Reading order: project knowledge → task definition.

**Project Knowledge**: Read relevant project knowledge files first (domain constraints):
- Infer relevant domains from task title, scope, and feature slug
- Read matching files from `docs/business-rules/` and `docs/conventions/`
- Example mappings: "auth"/"login"/"permission" → `business-rules/auth.md`; "state"/"validation"/"lifecycle" → `business-rules/<domain>.md`; "API"/"endpoint"/"route" → `conventions/api.md`; "error"/"status code" → `conventions/error-handling.md`; "database"/"schema"/"migration" → `conventions/data-model.md`; "test"/"mock"/"coverage" → `conventions/testing.md`
- If no matching file exists, skip this step

## Step 2: TDD Implementation

```
RED      → Write failing test first
GREEN    → Implement minimal code to pass
REFACTOR → Clean up while keeping tests green
```

## Step 3: Full Verification (Quality Gate)

Execute the quality gate sequence. Apply **Scope Resolution** from the Forge Guide for each command:

```
just compile [scope] → just fmt [scope] → just lint [scope] → just test [scope]
```

Strict sequential order. Stop at first failure:

| Failed step | Action |
|---|---|
| `compile` | Fix compilation errors, then retry from compile |
| `fmt` | Mark task as `blocked` (auto-fix failed = toolchain issue) |
| `lint` | Self-fix (max 1 retry), then mark `blocked` if still failing |
| `test` | Fix failing tests, then retry from compile |

### Fix-Task Escape Hatch

When failures are beyond the current task's scope (pre-existing bugs, environment issues, cross-module failures), create a fix-task instead of struggling:

```bash
task template fix-task  # View template and required variables first
task status <TASK_ID> blocked
task add --template fix-task --title "Fix: <concise description>" \
  --source-task-id <TASK_ID> \
  --var SOURCE_FILES="<affected source paths>" \
  --var TEST_SCRIPT="<failing test file>" \
  --var TEST_RESULTS="<test results path>" \
  --description "<root cause and context>"
```

The fix-task (P0) will be auto-claimed by the next `task claim`. When it completes, `task record` auto-restores the source task to pending via SourceTaskID.

### E2E Reminder (After Step 3)

After the quality gate passes, check whether to show an e2e reminder:

1. SCOPE (from Step 1) is `frontend` or `all`
2. Glob for `tests/e2e/features/<FEATURE>/*.spec.ts` — files exist

If both are true, print:
```
"Reminder: this task may affect e2e specs. Consider running:
 just test-e2e --feature <FEATURE>"
```

This reminder is informational only for manual `/execute-task` invocation. When dispatched by `/run-tasks`, e2e gates are handled by the dispatcher's Step 5b automatically.

## Step 4: Record Task (MANDATORY)

Invoke the skill (it contains file location details):

```
Skill(skill="record-task")
```

The skill provides complete workflow including:
- File locations (process/record.json vs tmp/)
- JSON format and fields
- CLI command usage

## Step 5: Commit

```
Skill(skill="git-commit")
```

## Rules

<EXTREMELY-IMPORTANT>
- record-task is mandatory - No completion without it
- All verifications must pass
- Commit only after record
- ONE TASK PER INVOCATION — after Step 5, STOP immediately, no exceptions
- FORBIDDEN: run "task claim", read index.json, or start any subsequent task
</EXTREMELY-IMPORTANT>

<HARD-GATE>
Task is NOT complete until record-task CLI command succeeds. Commit is blocked until record exists.
</HARD-GATE>

## STOP

<HARD-RULE>
ONE TASK PER INVOCATION. This is absolute and non-negotiable.

After Step 5, you MUST stop immediately.

<PROHIBITIONS>
- Running `task claim` under any circumstances
- Reading the next task file
- Continuing with any additional work
</PROHIBITIONS>

Output your final summary and STOP.
</HARD-RULE>

## Related Commands

| Command | Usage |
|---------|-------|
| `/run-tasks` | Auto-execute all tasks |
| `/record-task` | Create record + update status |
