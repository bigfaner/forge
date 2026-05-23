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
| `docs/features/<slug>/tasks/process/record.json` | In-progress execution notes | Not committed     |
| `docs/features/<slug>/tasks/records/*.md`        | Final completed record      | Committed to repo |

**Workflow:**

```
1. forge task claim    → writes process/state.json (current task)
2. During execution    → write progress to process/record.json
3. forge task submit --data → reads JSON, generates records/*.md, clears process/
```

## Step 1: Determine Record Format

<HARD-RULE>
You MUST read the category-specific format file before writing record.json. The format varies by task type — using the wrong format produces invalid records.
</HARD-RULE>

Your task type is set during `forge task claim`. Map it to a category:

| Category | Task Types | Format File |
|----------|-----------|-------------|
| coding | `coding.feature`, `coding.enhancement`, `coding.cleanup`, `coding.refactor`, `coding.fix`, `code-quality.simplify` | `data/record-format-coding.md` |
| doc | `doc`, `doc.eval`, `doc.summary`, `doc.consolidate`, `doc.drift` | `data/record-format-doc.md` |
| test | `test.gen-cases`, `test.eval-cases`, `test.gen-scripts`, `test.run`, `test.gen-and-run`, `test.graduate`, `test.verify-regression` | `data/record-format-test.md` |
| validation | `validation.code`, `validation.ux` | `data/record-format-validation.md` |
| gate | `gate` | `data/record-format-gate.md` |

**Read the format file now.** It contains the JSON example, category-specific fields, and metrics collection rules for your task type.

**Fallback:** If the task type is empty, unlisted, or you cannot determine the category, treat it as `coding` category and add a `notes` field explaining the fallback. Do NOT skip reading a format file.

## Step 2: Write record.json

Follow the JSON example from your category's format file. All categories share these base fields:

> **IMPORTANT:** The table below shows ONLY shared fields. Your category's format file contains the COMPLETE JSON example with category-specific fields. Do not construct record.json from this table alone.

| Field                  | Type   | Required | Description                                  |
| ---------------------- | ------ | -------- | -------------------------------------------- |
| `taskId`               | string | auto     | Task ID (verified against CLI arg, mismatch = hard error) |
| `status`               | string | auto     | Defaults to `completed`; must be valid enum value |
| `summary`              | string | **hard** | Implementation summary. Empty = hard error (non-overridable) |
| `filesCreated`         | array  | optional | List of newly created files                  |
| `filesModified`        | array  | optional | List of modified files                       |
| `keyDecisions`         | array  | warning  | Key design decisions. Missing = warning (completed status only) |
| `acceptanceCriteria`   | array  | warning  | `{criterion, met}` objects. Missing = warning; any `met:false` = hard error (overridable) |
| `notes`                | string | optional | Optional notes or observations               |
| `typeReclassification` | object | optional | When executor discovers task type doesn't match actual work |

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

## General Metrics Rules

- `coverage` = `-1.0` is auto-set by CLI for non-`coding.*` type tasks
- Never fabricate metrics — if a field does not apply to the task category, omit it
- Do NOT include fields from other categories (e.g., no `testsPassed` in doc records)

## Usage

```bash
# Step 1: Write execution data to process/record.json using the Write tool
# (replace sample values with real metrics collected above)
# Target path: docs/features/<slug>/tasks/process/record.json

# Step 2: Use CLI command (mandatory)
forge task submit <TASK_ID> --data docs/features/<slug>/tasks/process/record.json
```

<EXTREMELY-IMPORTANT>
You MUST use the `forge task submit` CLI command. No exceptions.

**ONLY ALLOWED PATH:** `docs/features/<slug>/tasks/process/record.json`

**DO NOT:**

- Write directly to index.json
- Write to `records/*.md` manually (CLI generates these)
- Use `forge task status <id> completed` to set status (use `forge task submit`)
- Use any other file path (e.g., `.claude/tmp/`)

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
| `forge task status <id> completed` | Status mutation removed — use `forge task submit` to complete tasks |
| Writing record.json to wrong path | CLI only reads from `process/record.json`  |

</EXTREMELY-IMPORTANT>

## Recovery (Only when `forge task submit` fails)

If `forge task submit` fails and cannot be fixed:

1. Identify the root cause from the error output (missing fields, validation failure, state transition error)
2. Fix the `record.json` data and re-run `forge task submit <TASK_ID> --data <path>`
3. If the task is stuck in an invalid state, use `forge task add --template fix-task --title "Fix: submit failed" --source-task-id <TASK_ID> --block-source --description "Submit failed: <error>"` to block it and create a fix task
