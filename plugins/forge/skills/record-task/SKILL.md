---
name: record-task
description: Use after completing a task to create its execution record and update task status.
---

# Record Task

## Overview

Post-task completion: create execution record + update task status.

## File Locations

| Location                                         | Purpose                     | Git Status        |
| ------------------------------------------------ | --------------------------- | ----------------- |
| `docs/features/{slug}/tasks/process/record.json` | In-progress execution notes | Not committed     |
| `docs/features/{slug}/tasks/records/*.md`        | Final completed record      | Committed to repo |

**Workflow:**

```
1. task claim           → writes process/state.json (current task)
2. During execution     → write progress to process/record.json
3. task record --data   → reads JSON, generates records/*.md, clears process/
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

| Field                 | Type   | Description                                  |
| --------------------- | ------ | -------------------------------------------- |
| `taskId`              | string | Task ID (verified against CLI arg if provided) |
| `status`              | string | Task status, defaults to `completed`         |
| `summary`             | string | Implementation summary                       |
| `filesCreated`        | array  | List of newly created files                  |
| `filesModified`       | array  | List of modified files                       |
| `keyDecisions`        | array  | Key design decisions                         |
| `testsPassed`         | int    | Number of tests passed                       |
| `testsFailed`         | int    | Number of tests failed                       |
| `coverage`            | float  | Coverage percentage (collected from test runner) |
| `acceptanceCriteria`  | array  | `{criterion, met}` objects                   |
| `notes`               | string | Optional notes or observations               |

## Metrics Collection (MANDATORY before writing record.json)

<HARD-RULE>
Before writing `record.json`, you MUST collect real metrics from the project's test runner. All numeric fields (`coverage`, `testsPassed`, `testsFailed`) must come from actual output, never guessed or defaulted.

Coverage rules:

- `coverage` = actual percentage from test runner output
- `coverage` = `-1.0` is auto-set by CLI for tasks with `noTest: true`. For non-noTest tasks, always report real metrics.
- Never write `0.0` unless the runner actually reported 0%

Example commands (use whatever matches the project's toolchain):

```
just test [scope]
```

</HARD-RULE>

## Usage

```bash
# Step 1: Write progress to process/record.json (replace sample values with real metrics from above)
echo '{"taskId":"3.3.1","status":"completed","summary":"...","filesCreated":["..."],"filesModified":["..."],"keyDecisions":["..."],"testsPassed":12,"testsFailed":0,"coverage":85.6,"acceptanceCriteria":[{"criterion":"...","met":true}]}' > docs/features/{slug}/tasks/process/record.json

# Step 2: Use CLI command (mandatory)
task record <TASK_ID> --data docs/features/{slug}/tasks/process/record.json
```

<EXTREMELY-IMPORTANT>
You MUST use the `task record` CLI command. No exceptions.

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

## What `task record` Does (One Command = 2 Operations)

```
task record <TASK_ID> --data docs/features/{slug}/tasks/process/record.json
```

This single command automatically:

1. ✅ Generates `records/*.md` from JSON
2. ✅ Updates `index.json` status

After running, check the STATUS field in the output:
- `STATUS: completed` → task recorded successfully, proceed to commit
- `STATUS: blocked` → task was auto-downgraded (e.g. test failures), **do NOT commit**

## Validation Rules (enforced by CLI)

### Quality Gate Pre-check

When `status=completed`, `--force` is NOT used, and the task does NOT have `noTest: true` in its index.json entry, `task record` automatically runs the full quality gate before accepting the record:

```
just compile [scope] → just fmt [scope] → just lint [scope] → just test [scope]
```

This gate runs via `validateQualityGate()` in record.go. If any step fails, the record is rejected with an error. Fix the errors and re-run `task record`, or use `--force` to bypass.

### Data Validation

`task record` will reject the following combinations:

| Condition | Error | Fix |
|-----------|-------|-----|
| `status=completed` + `testsFailed > 0` | Auto-downgrade to `blocked` | Fix the test failures, then re-record with `status: "completed"` |
| `status=completed` + `testsPassed=0` + `testsFailed=0` + `coverage >= 0` | No test evidence | Ensure task has `noTest: true` (CLI auto-sets coverage), or run tests and report results |
| `status=completed` + any `acceptanceCriteria.met=false` | Unmet acceptance criteria | Fix the issue, or set `status: "blocked"` |
| `summary` is empty or whitespace | Missing summary | Provide a summary |

Override any validation error with `--force`:
```bash
task record <TASK_ID> --data record.json --force
```

Use `--force` only when you have a specific reason (document it in `notes`).

## Forbidden Operations

<EXTREMELY-IMPORTANT>
These actions will corrupt task state:

| Operation                    | Why Forbidden                              |
| ---------------------------- | ------------------------------------------ |
| `Write("records/*.md")`      | Bypasses CLI validation and hooks          |
| Direct edit to `index.json`  | State becomes inconsistent                 |
| `task status <id> completed` | Only for recovery when `task record` fails |
| Writing to wrong path        | CLI only reads from `process/record.json`  |

</EXTREMELY-IMPORTANT>

## Recovery (Only when `task record` fails)

If `task record` fails and cannot be fixed:

```bash
# Manual status fix (last resort only)
task status <TASK_ID> completed --force
```
