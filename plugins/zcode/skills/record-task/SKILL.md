---
name: record-task
description: Use after completing a task to create its execution record and update task status.
---

# Record Task

## Overview

任务完成后的收尾操作：创建执行记录 + 更新任务状态。

## File Locations

| Location | Purpose | Git Status |
|----------|---------|------------|
| `docs/features/{slug}/tasks/process/record.json` | In-progress execution notes | Not committed |
| `docs/features/{slug}/tasks/records/*.md` | Final completed record | Committed to repo |

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
| `coverage`           | float  | 覆盖率                     |
| `acceptanceCriteria` | array  | `{criterion, met}` 对象    |

## Usage

```bash
# Step 1: Write progress to process/record.json
echo '{"summary":"...","filesCreated":[...],"filesModified":[...]}' > docs/features/{slug}/tasks/process/record.json

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

## Related

- `/claim-task` - Claim next available task
- `/set-task-status` - Direct status update only
