---
name: record-task
description: Use after completing a task to create its execution record and update task status.
---

# Record Task

## Overview

任务完成后的收尾操作：创建执行记录 + 更新任务状态。

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
# Step 1: 准备 JSON 数据
echo '{"summary":"...","filesCreated":[...],"filesModified":[...]}' > docs/features/{slug}/tasks/process/record.json

# Step 2: 使用 CLI 命令（必须）
task record <TASK_ID> --data docs/features/{slug}/tasks/process/record.json
```

## ⚠️ Iron Law

```
┌─────────────────────────────────────────────────────────────────┐
│  YOU MUST USE `task record` COMMAND                             │
│                                                                 │
│  DO NOT:                                                        │
│  - Write directly to index.json                                 │
│  - Use Python/JavaScript to modify JSON                         │
│  - Create record files manually                                 │
│                                                                 │
│  The CLI command provides:                                      │
│  - Schema validation                                            │
│  - Consistent output format                                     │
│  - Potential hooks/side-effects                                 │
└─────────────────────────────────────────────────────────────────┘
```

## Related

- `/claim-task` - Claim next available task
- `/set-task-status` - Direct status update only
