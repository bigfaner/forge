---
name: submit-task
description: Use after completing a task to create its execution record and update task status.
argument-hint: "[task-id]"
---

# Submit Task

## Overview

Post-task completion: create execution record + update task status.

## File Locations

| Location                                         | Purpose                     | Git Status        |
| ------------------------------------------------ | --------------------------- | ----------------- |
| `docs/features/{slug}/tasks/process/record.json` | In-progress execution notes | Not committed     |
| `docs/features/{slug}/tasks/records/*.md`        | Final completed record      | Committed to repo |

**Workflow:**

```
1. forge task claim    → writes process/state.json (current task)
2. During execution    → write progress to process/record.json
3. forge task submit --data → reads JSON, generates records/*.md, clears process/
```

## JSON Data Format

```json
{
	"taskId": "3.3.1",
	"status": "completed",
	"summary": "What was implemented",
	"filesCreated": ["src/components/Button.tsx"],
	"filesModified": ["src/utils/helpers.ts"],
	"keyDecisions": ["Decision 1"],
	"testsPassed": 12,
	"testsFailed": 0,
	"coverage": 85.6,
	"acceptanceCriteria": [{ "criterion": "Acceptance criterion 1", "met": true }],
	"notes": "Optional observations"
}
```

## Fields

| Field                 | Type   | Required | Description                                  |
| --------------------- | ------ | -------- | -------------------------------------------- |
| `taskId`              | string | auto     | Task ID (verified against CLI arg, mismatch = hard error) |
| `status`              | string | auto     | Defaults to `completed`; must be valid enum value |
| `summary`             | string | **hard** | Implementation summary. Empty = hard error (non-overridable) |
| `filesCreated`        | array  | optional | List of newly created files                  |
| `filesModified`       | array  | optional | List of modified files                       |
| `keyDecisions`        | array  | warning  | Key design decisions. Missing = warning (completed status only) |
| `testsPassed`         | int    | context  | Number of tests passed. See Metrics Collection below |
| `testsFailed`         | int    | context  | Number of tests failed. >0 with completed = auto-downgrade to blocked |
| `coverage`            | float  | context  | Coverage percentage. Auto-set to `-1.0` for non-`coding.*` type tasks |
| `acceptanceCriteria`  | array  | warning  | `{criterion, met}` objects. Missing = warning; any `met:false` = hard error (overridable) |
| `notes`               | string | optional | Optional notes or observations               |
| `typeReclassification` | object | optional | When executor discovers task type doesn't match actual work |

> **context** = required for `completed` tasks with a `coding.*` type; auto-relaxed when type does not start with `coding.`.

## Type Reclassification

When executing a task, you may discover that the assigned type doesn't match the actual work. For example, a `fix` task might turn out to be a flaky test cleanup, or a `feature` task might only involve refactoring existing code.

In such cases, process the task according to its **actual type** and include a `typeReclassification` block in the JSON data:

```json
{
  "taskId": "fix-1",
  "status": "completed",
  "summary": "Fixed flaky test by cleaning up test isolation",
  "typeReclassification": {
    "originalType": "fix",
    "actualType": "cleanup",
    "reason": "Root cause was test state leak between runs, not a code bug"
  }
}
```

The reclassification is recorded in the task's execution log for traceability. The original type in `index.json` is **not** changed — only the record documents the discrepancy.

## Metrics Collection (MANDATORY before writing record.json)

<HARD-RULE>
Before writing `record.json`, you MUST collect real metrics from the project's test runner. All numeric fields (`coverage`, `testsPassed`, `testsFailed`) must come from actual output, never guessed or defaulted.

Coverage rules:

- `coverage` = actual percentage from test runner output
- `coverage` = `-1.0` is auto-set by CLI for non-`coding.*` type tasks (e.g., `doc`, `test.*`, `validation.*`). For testable tasks (any `coding.*` type), always report real metrics.
- Never write `0.0` unless the runner actually reported 0%

Capture metrics from the targeted test runs you performed during task development (framework-native commands on changed code). Report the actual pass/fail counts and coverage from those runs.

</HARD-RULE>

## Usage

```bash
# Step 1: Write progress to process/record.json (replace sample values with real metrics from above)
echo '{"taskId":"3.3.1","status":"completed","summary":"...","filesCreated":["..."],"filesModified":["..."],"keyDecisions":["..."],"testsPassed":12,"testsFailed":0,"coverage":85.6,"acceptanceCriteria":[{"criterion":"...","met":true}]}' > docs/features/{slug}/tasks/process/record.json

# Step 2: Use CLI command (mandatory)
forge task submit <TASK_ID> --data docs/features/{slug}/tasks/process/record.json
```

<EXTREMELY-IMPORTANT>
You MUST use the `forge task submit` CLI command. No exceptions.

**ONLY ALLOWED PATH:** `docs/features/{slug}/tasks/process/record.json`

**DO NOT:**

- Write directly to index.json
- Use Python/JavaScript to modify JSON
- Create record files manually
- Use Bash echo/cat to write JSON
- Think "both approaches achieve the same result"
- Use any other file path (e.g., .claude/tmp/)

The CLI command provides schema validation, consistent output format, and potential hooks/side-effects.
Bypassing the command defeats the purpose of the skill.
</EXTREMELY-IMPORTANT>

## What `forge task submit` Does

`forge task submit` generates the execution record and updates task status. After running, check the STATUS field:
- `STATUS: completed` → task recorded successfully, proceed to commit
- `STATUS: blocked` → task was auto-downgraded (e.g. test failures), **do NOT commit**

## Forbidden Operations

<EXTREMELY-IMPORTANT>
These actions will corrupt task state:

| Operation                    | Why Forbidden                              |
| ---------------------------- | ------------------------------------------ |
| `Write("records/*.md")`      | Bypasses CLI validation and hooks          |
| Direct edit to `index.json`  | State becomes inconsistent                 |
| `forge task status <id> completed` | Only for recovery when `forge task submit` fails |
| Writing to wrong path        | CLI only reads from `process/record.json`  |

</EXTREMELY-IMPORTANT>

## Recovery (Only when `forge task submit` fails)

If `forge task submit` fails and cannot be fixed:

```bash
# Manual status fix (last resort only)
forge task status <TASK_ID> completed --force
```
