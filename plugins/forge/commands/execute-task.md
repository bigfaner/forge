---
name: execute-task
description: Execute single task with workflow-driven execution.
allowed_tools: ["Bash", "Read", "Write", "Edit", "Grep", "Glob", "Agent", "LSP"]
---

# /execute-task

Execute a single task with workflow-driven execution.

## Step 0: MAIN_SESSION Check

If the claimed task has `MAIN_SESSION == "true"`, read the task file's `## Main Session Instructions` section and follow it. After execution, invoke `Skill(skill="record-task")`. Skip Steps 2–4.

## Workflow (4 Steps)

```
Step 1: Claim & Read task definition
Step 2: Execute Workflow
Step 3: Record task (MANDATORY)
Step 4: Git commit
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

## Step 2: Execute Workflow

<EXTREMELY-IMPORTANT>
You MUST determine the execution workflow for this task by following this exact procedure:

1. Read the task file specified in Step 1.
2. Search for a `## Execution Workflow` heading in the task file.
3. Based on what you find:

   **CASE A — `## Execution Workflow` heading exists with non-empty content:**
   The content under the heading (excluding the heading line itself, up to the next
   `##` heading or end of file) is your execution instructions. Follow these steps
   EXACTLY. Do not deviate, add, or skip steps.

   **CASE B — No `## Execution Workflow` heading found:**
   Read the default workflow template at:
   `plugins/forge/skills/breakdown-tasks/templates/task.md`
   Find its `## Execution Workflow` section and follow those steps.

   **CASE C — `## Execution Workflow` heading exists but content is empty:**
   Log: "WARNING: ## Execution Workflow heading present but empty. Falling back to default template."
   Then proceed as Case B.

4. Output after execution:
   - Success: `Step 2/4: [workflow description]... DONE`
   - Failure: `Step 2/4: [workflow description]... FAILED: [reason]`
</EXTREMELY-IMPORTANT>

### Fix-Task Escape Hatch

When failures are beyond the current task's scope (pre-existing bugs, environment issues, cross-module failures), create a fix-task instead of struggling:

```bash
task template fix-task  # View template and required variables first
task add --template fix-task --title "Fix: <concise description>" \
  --source-task-id <TASK_ID> \
  --block-source \
  --var SOURCE_FILES="<affected source paths>" \
  --var TEST_SCRIPT="<failing test file>" \
  --var TEST_RESULTS="<test results path>" \
  --description "<root cause and context>"
```

**`--block-source`**: atomically sets source task to blocked before resolution, preserving the fix-chain model.
**`--source-task-id` auto-resolves**: if `<TASK_ID>` is a **completed** fix-task, the CLI automatically resolves to the root blocked task. Always pass the current failing task's ID — no manual chain tracing needed.
`task add` automatically deduplicates — check output: `ACTION: ADDED` (new fix task) or `ACTION: SKIPPED` (active fix already exists).

The fix-task (P0) will be auto-claimed by the next `task claim`. When it completes, `task record` auto-restores the source task to pending via SourceTaskID.

### E2E Reminder (After Step 2)

After the workflow completes successfully, check whether to show an e2e reminder:

1. SCOPE (from Step 1) is `frontend` or `all`
2. Glob for `tests/e2e/features/<FEATURE>/*.spec.ts` — files exist

If both are true, print:
```
"Reminder: this task may affect e2e specs. Consider running:
 just test-e2e --feature <FEATURE>"
```

This reminder is informational only for manual `/execute-task` invocation. When dispatched by `/run-tasks`, e2e gates are handled by the dispatcher's Step 5b automatically.

## Step 3: Record Task (MANDATORY)

Invoke the skill (it contains file location details):

```
Skill(skill="record-task")
```

The skill provides complete workflow including:
- File locations (process/record.json vs tmp/)
- JSON format and fields
- CLI command usage

## Step 4: Commit

```
Skill(skill="git-commit")
```

## Rules

<EXTREMELY-IMPORTANT>
- record-task is mandatory - No completion without it
- All verifications must pass
- Commit only after record
- ONE TASK PER INVOCATION — after Step 4, STOP immediately, no exceptions
- FORBIDDEN: run "task claim", read index.json, or start any subsequent task
</EXTREMELY-IMPORTANT>

<HARD-GATE>
Task is NOT complete until record-task CLI command succeeds. Commit is blocked until record exists.
</HARD-GATE>

## STOP

<HARD-RULE>
ONE TASK PER INVOCATION. This is absolute and non-negotiable.

After Step 4, you MUST stop immediately.

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
