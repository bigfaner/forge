# Record Task 必须捕获实际测试覆盖率

## Problem

Task 3.1 记录显示 `Coverage: 0.0%`，实际 `go test -cover` 为 100%。字段存在但数据失真。

## Root Cause

`record-task` skill 定义了 `coverage` 字段，但 workflow 缺少采集步骤。省略该字段时 CLI 默认填 0.0%，无报错。

**有字段定义，无采集动作。**

## Solution

写 `record.json` 前，运行覆盖率命令并将结果填入 `coverage` 字段：

| 语言 | 命令 |
|------|------|
| Go | `go test -cover ./changed/package/...` |
| TypeScript | `npm test -- --coverage --watchAll=false` |

若项目无测试，`coverage` 应省略或填 `null`，不可填 0——0 表示"测了但没覆盖"，不是"没测"。

## Takeaway

Workflow 中的字段 = 格式定义 + 采集步骤。缺采集，字段会被默认值静默填充，数据失真且不易察觉。
