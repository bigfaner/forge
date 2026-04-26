---
name: record-task
description: Use after completing a task to create its execution record and update task status.
---

# Record Task

## Overview

任务完成后的收尾操作：创建执行记录 + 更新任务状态。

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
	"summary": "实现了什么",
	"filesCreated": ["src/components/Button.tsx"],
	"filesModified": ["src/utils/helpers.ts"],
	"keyDecisions": ["决策 1"],
	"testsPassed": 12,
	"testsFailed": 0,
	"coverage": 85.6,
	"acceptanceCriteria": [{ "criterion": "验收标准 1", "met": true }]
}
```

## Fields

| 字段                 | 类型   | 说明                       |
| -------------------- | ------ | -------------------------- |
| `status`             | string | 任务状态，默认 `completed` |
| `summary`            | string | 实现摘要                   |
| `filesCreated`       | array  | 新建文件列表               |
| `filesModified`      | array  | 修改文件列表               |
| `keyDecisions`       | array  | 关键设计决策               |
| `testsPassed`        | int    | 通过测试数                 |
| `testsFailed`        | int    | 失败测试数                 |
| `coverage`           | float  | 覆盖率（须从测试工具采集） |
| `acceptanceCriteria` | array  | `{criterion, met}` 对象    |

## Metrics Collection (MANDATORY before writing record.json)

<HARD-RULE>
Before writing `record.json`, you MUST collect real metrics from the project's test runner. All numeric fields (`coverage`, `testsPassed`, `testsFailed`) must come from actual output, never guessed or defaulted.

Coverage rules:

- `coverage` = actual percentage from test runner output
- `coverage` = `-1.0` when task has no tests
- Never write `0.0` unless the runner actually reported 0%

Example commands (use whatever matches the project's toolchain):

```
Go:        go test -cover ./changed/package/...
TypeScript: npm test -- --coverage --watchAll=false
Python:    pytest --cov=<module> --cov-report=term-missing
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

## What `task record` Does (One Command = 3 Operations)

```
task record <TASK_ID> --data docs/features/{slug}/tasks/process/record.json
```

This single command automatically:

1. ✅ Generates `records/*.md` from JSON
2. ✅ Updates `index.json` status to `completed`

**You don't need to do anything else after calling this command.**

## Validation Rules (enforced by CLI)

`task record` will reject the following combinations:

| Condition | Error | Fix |
|-----------|-------|-----|
| `status=completed` + `testsPassed=0` + `testsFailed=0` + `coverage >= 0` | No test evidence | Run tests and report results, or set `coverage: -1.0` for no-test tasks |
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
task status <TASK_ID> completed
```
