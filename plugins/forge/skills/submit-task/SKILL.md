---
name: submit-task
description: Use after completing a task to create its execution record and update task status.
argument-hint: "[task-id]"
---

# Submit Task

## Common Fields

All categories share these fields:

| Field                | Type                          | Required    | Description                                                |
| -------------------- | ----------------------------- | ----------- | ---------------------------------------------------------- |
| `taskId`             | string                        | auto        | Omit or match the CLI argument — CLI validates consistency |
| `status`             | `"completed"` \| `"blocked"`  | required    | Task outcome                                               |
| `summary`            | string                        | required    | What was accomplished                                      |
| `filesCreated`       | string[]                      | optional    | Files created during execution                             |
| `filesModified`      | string[]                      | optional    | Files modified during execution                            |
| `acceptanceCriteria` | `Array<{criterion, met}>`     | recommended | Criteria with pass/fail                                    |
| `notes`              | string                        | optional    | Additional observations                                    |

## Common Rules

- Format file rules take precedence over Common Fields when they conflict.
- `acceptanceCriteria` with any `met: false` entry is rejected for `completed` status — use `blocked` instead.
- `testsPassed`, `testsFailed`, `coverage` are coding-only. Do NOT include them in doc, test, validation, or gate records.

## File Locations

| Location                                         | Purpose                     |
| ------------------------------------------------ | --------------------------- |
| `docs/features/<slug>/tasks/process/record.json` | In-progress (not committed) |
| `docs/features/<slug>/tasks/records/*.md`        | Final record (committed)    |

## Step 1: Determine Record Format

<HARD-RULE>
You MUST read the category-specific format file before writing record.json.
</HARD-RULE>

The `TASK_CATEGORY` field is in your task prompt. Use that value to locate the format file.

The format files are located in the `data/` directory alongside this SKILL.md. Read the file at:

```
data/record-format-{TASK_CATEGORY}.md
```

This is a path **relative to this SKILL.md's directory** (`skills/submit-task/`). Claude resolves it from the skill file location, not from the project root.

**Fallback:** If `TASK_CATEGORY` is not in your context, run `forge task query <TASK_ID>` to retrieve it. Do NOT skip reading a format file.

> `code-quality.simplify` maps to the `coding` category — use `record-format-coding.md`.

## Step 2: Write record.json

Use the JSON example from your category's format file as the template. Combine common fields with category-specific fields.

## Step 3: Submit via CLI

```bash
forge task submit <TASK_ID> --data docs/features/<slug>/tasks/process/record.json
```

<EXTREMELY-IMPORTANT>
You MUST use the `forge task submit` CLI command. No exceptions.

ONLY ALLOWED PATH: `docs/features/<slug>/tasks/process/record.json`
</EXTREMELY-IMPORTANT>

## Forbidden Operations

| Operation                          | Why Forbidden                      |
| ---------------------------------- | ---------------------------------- |
| `Write("records/*.md")`            | Bypasses CLI validation and hooks  |
| Direct edit to `index.json`        | State becomes inconsistent         |
| `forge task status <id> completed` | Use `forge task submit` to complete|
| Wrong record.json path             | CLI only reads from `process/`     |

## What `forge task submit` Does

Generates the execution record and updates task status. After running, check STATUS:
- `STATUS: completed` → recorded successfully, proceed to commit
- `STATUS: blocked` → auto-downgraded (e.g. test failures), **do NOT commit**

## Type Reclassification

When the assigned type doesn't match the actual work, process according to the **actual type** and include a `typeReclassification` block:

```json
{
  "typeReclassification": {
    "originalType": "coding.fix",
    "actualType": "coding.cleanup",
    "reason": "Root cause was test state leak between runs, not a code bug"
  }
}
```

The original type in `index.json` is **not** changed.

## Recovery (when `forge task submit` fails)

1. Fix `record.json` and re-run `forge task submit <TASK_ID> --data <path>`
2. If stuck: `forge task add --type coding.fix --title "Fix: submit failed" --source-task-id <TASK_ID> --block-source --var SOURCE_FILES="<affected-files>" --var TEST_SCRIPT="<test-path>" --var TEST_RESULTS="<test-output>" --description "Submit failed: <error>"`
