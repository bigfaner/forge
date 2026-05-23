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

## JSON Data Format (Coding Category)

The format below applies to `coding.*` and `code-quality.simplify` tasks. For other categories, see [Type-Specific Record Formats](#type-specific-record-formats) below.

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

### Shared Fields (All Categories)

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

### Coding-Category Fields

| Field          | Type  | Required | Description                                  |
| -------------- | ----- | -------- | -------------------------------------------- |
| `testsPassed`  | int   | context  | Number of tests passed. See Metrics Collection below |
| `testsFailed`  | int   | context  | Number of tests failed. >0 with completed = auto-downgrade to blocked |
| `coverage`     | float | context  | Coverage percentage. Auto-set to `-1.0` for non-`coding.*` type tasks |

> **context** = required for `completed` tasks with a `coding.*` type; auto-relaxed when type does not start with `coding.`.

### Doc-Category Fields

| Field             | Type     | Required | Description                                      |
| ----------------- | -------- | -------- | ------------------------------------------------ |
| `referencedDocs`  | string[] | optional | List of document paths referenced during work     |
| `reviewStatus`    | string   | optional | Review status (e.g., "draft", "final", "reviewed") |
| `docMetrics`      | string   | optional | Free-text document metrics (word count, sections, etc.) |

Doc-category tasks (`doc`, `doc.eval`, `doc.summary`, `doc.consolidate`, `doc.drift`) should NOT include `testsPassed`, `testsFailed`, or `coverage` fields.

### Test-Category Fields

| Field             | Type     | Required | Description                                        |
| ----------------- | -------- | -------- | -------------------------------------------------- |
| `casesGenerated`  | int      | optional | Number of test cases generated                     |
| `casesEvaluated`  | int      | optional | Number of test cases evaluated                     |
| `scriptsCreated`  | string[] | optional | List of test script files created                  |
| `testResults`     | string   | optional | Free-text summary of test execution results        |

Test-category tasks (`test.gen-cases`, `test.eval-cases`, `test.gen-scripts`, `test.run`, `test.gen-and-run`, `test.graduate`, `test.verify-regression`) use these fields instead of `testsPassed`/`testsFailed`/`coverage`.

### Validation-Category Fields

| Field              | Type     | Required | Description                                    |
| ------------------ | -------- | -------- | ---------------------------------------------- |
| `validationPassed` | bool     | optional | Whether validation passed                      |
| `issuesFound`      | string[] | optional | List of issues found during validation         |

Validation-category tasks (`validation.code`, `validation.ux`) produce pass/fail verdicts rather than test metrics.

### Gate-Category Fields

| Field         | Type     | Required | Description                                    |
| ------------- | -------- | -------- | ---------------------------------------------- |
| `gatePassed`  | bool     | optional | Whether the gate check passed                  |
| `gateChecks`  | string[] | optional | List of individual gate check results          |

Gate-category tasks (`gate`) produce minimal records with gate checks and overall pass status.

## Type-Specific Record Formats

### Category Mapping

The 21 task types map to 5 categories:

| Category | Task Types |
|----------|-----------|
| coding | `coding.feature`, `coding.fix`, `coding.refactor`, `code-quality.simplify` |
| doc | `doc`, `doc.eval`, `doc.summary`, `doc.consolidate`, `doc.drift` |
| test | `test.gen-cases`, `test.eval-cases`, `test.gen-scripts`, `test.run`, `test.gen-and-run`, `test.graduate`, `test.verify-regression` |
| validation | `validation.code`, `validation.ux` |
| gate | `gate` |

### Doc Category Example

For `doc`, `doc.eval`, `doc.summary`, `doc.consolidate`, `doc.drift` tasks:

```json
{
	"taskId": "7",
	"status": "completed",
	"summary": "Updated submit-task SKILL.md with per-type instructions",
	"filesCreated": [],
	"filesModified": ["plugins/forge/skills/submit-task/SKILL.md"],
	"keyDecisions": ["Added type-specific sections after existing Fields table"],
	"referencedDocs": [
		"docs/proposals/typed-task-records/proposal.md",
		"docs/features/typed-task-records/tasks/1-category-and-recorddata.md"
	],
	"reviewStatus": "final",
	"docMetrics": "Added ~120 lines covering 5 categories with JSON examples and field tables",
	"acceptanceCriteria": [{ "criterion": "JSON examples for all 5 categories", "met": true }],
	"notes": "Followed existing structure; added sections rather than reorganizing"
}
```

### Test Category Example

For `test.gen-cases`, `test.eval-cases`, `test.gen-scripts`, `test.run`, `test.gen-and-run`, `test.graduate`, `test.verify-regression` tasks:

```json
{
	"taskId": "12",
	"status": "completed",
	"summary": "Generated test cases for login flow",
	"filesCreated": ["tests/e2e/login/login-flow.spec.ts"],
	"filesModified": [],
	"casesGenerated": 8,
	"casesEvaluated": 8,
	"scriptsCreated": ["tests/e2e/login/login-flow.spec.ts"],
	"testResults": "8 cases generated, 8 passed evaluation criteria",
	"acceptanceCriteria": [{ "criterion": "All login scenarios covered", "met": true }],
	"notes": "Auto-generated via test pipeline"
}
```

### Validation Category Example

For `validation.code`, `validation.ux` tasks:

```json
{
	"taskId": "15",
	"status": "completed",
	"summary": "Validated code quality for auth module",
	"filesCreated": [],
	"filesModified": ["src/auth/middleware.ts"],
	"validationPassed": true,
	"issuesFound": ["Minor: unused import on line 42"],
	"acceptanceCriteria": [{ "criterion": "No critical issues found", "met": true }],
	"notes": "One minor issue auto-fixed during validation"
}
```

### Gate Category Example

For `gate` tasks:

```json
{
	"taskId": "20",
	"status": "completed",
	"summary": "Quality gate check passed for feature branch",
	"gatePassed": true,
	"gateChecks": [
		"All tests passing",
		"No lint errors",
		"Coverage above threshold"
	],
	"notes": "All checks passed"
}
```

### Coding Category Example

For `coding.feature`, `coding.fix`, `coding.refactor`, `code-quality.simplify` tasks — this is the unchanged default format shown in [JSON Data Format](#json-data-format-coding-category) above.

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
Before writing `record.json`, you MUST collect real metrics. Which metrics depend on the task category:

### Coding Category (coding.*, code-quality.simplify)
Collect from the project's test runner. All numeric fields (`coverage`, `testsPassed`, `testsFailed`) must come from actual output, never guessed or defaulted.

- `coverage` = actual percentage from test runner output
- Never write `0.0` unless the runner actually reported 0%
- Capture metrics from the targeted test runs you performed during task development

### Doc Category (doc*)
No test metrics required. Focus on:
- `referencedDocs` — list documents you read/modified
- `reviewStatus` — set to "draft", "final", or "reviewed" as appropriate
- `docMetrics` — describe what was produced (e.g., word count, sections added)

### Test Category (test.*)
Collect from the test pipeline:
- `casesGenerated` / `casesEvaluated` — actual counts from test case generation
- `scriptsCreated` — list of generated script files
- `testResults` — summary of test execution outcomes

### Validation Category (validation.*)
Collect from the validation process:
- `validationPassed` — boolean based on actual validation outcome
- `issuesFound` — list specific issues detected

### Gate Category (gate)
Collect from gate checks:
- `gatePassed` — boolean based on overall gate outcome
- `gateChecks` — list each check and its result

### General Rules
- `coverage` = `-1.0` is auto-set by CLI for non-`coding.*` type tasks
- Never fabricate metrics — if a field does not apply to the task category, omit it
- Do NOT include coding-category fields (testsPassed, testsFailed, coverage) in doc/test/validation/gate records

</HARD-RULE>

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
