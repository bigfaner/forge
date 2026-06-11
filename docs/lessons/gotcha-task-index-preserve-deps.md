---
created: "2026-05-28"
tags: [architecture]
---

# forge task index 不更新已有任务的依赖，导致 claim 跳过业务任务

## Problem

拆分 Task 2 为 2a/2b/2c 后，更新了 Task 4 的 .md 文件依赖从 `[2, 3]` 改为 `[2a, 2b, 2c, 3]`，重新运行 `forge task index`。但 Task 4 在 index.json 中的依赖仍为 `['2', '3']`（指向被 skip 的 Task 2），导致 Task 4 永远无法 claim。

更严重的是，`forge task claim` 在此状态下跳过了 pending P1 业务任务 2c，转而 claim 了 P2 auto-gen 的 quick-drift-detection 任务。`forge task list` 中所有业务任务标注 `[cycle]`。

## Root Cause

因果链（3 层）：

1. **表面现象**：`forge task claim` 跳过 P1 业务任务（2c），选择 P2 auto-gen 任务（quick-drift-detection），违反 P0 > P1 > P2 优先级规则
2. **直接原因**：
   - `forge task index` 在 re-run 时保留已有任务的 index.json 条目，不重新读取 .md 文件的 frontmatter。修改 .md 文件中的 `dependencies` 后重新 index，index.json 中的 dependencies 不变
   - Task 4 在 index.json 中仍依赖被 skip 的 Task 2 → claim 算法检测到依赖链断裂 → 将相关任务标注 `[cycle]` → 降低/排除这些任务的 claim 优先级
3. **根因**：
   - **index.json 与 .md 文件的双源真相问题**：index.json 是 claim 算法的唯一数据源，但 .md 文件是人类/agent 的编辑入口。两者不同步时，以 index.json 为准，但编辑行为发生在 .md 文件上
   - **`forge task index` 的 "preserve" 语义过于宽泛**：保留已有条目的全部字段（包括 dependencies），即使 .md 文件已更新。正确行为应该是：保留 status/record 等运行时状态，但从 .md frontmatter 重新读取 dependencies/priority 等声明式字段

## Solution

### 短期修复（手动同步）

修改 index.json 中 Task 4 的 dependencies：

```json
"dependencies": ["2a", "2b", "2c", "3"]
```

### 长期修复（forge task index）

`forge task index` 应区分两类字段：
- **运行时状态**（保留）：status, record, startedTime
- **声明式元数据**（从 .md 重新读取）：dependencies, priority, title, estimatedTime

## Reusable Pattern

- **index.json 是 claim 的唯一数据源**：修改 .md 文件不会影响已 index 的任务。如果需要更新 dependencies/priority，必须直接修改 index.json 或删除旧条目后重新 index
- **拆分任务后必须同步 index.json**：当拆分/合并任务改变了依赖关系时，手动检查 index.json 中相关任务的 dependencies 是否与 .md 文件一致
- **auto-gen 任务无依赖，claim 时可能抢先**：quick-drift-detection（P2, no deps）比 2c（P1, deps=['1'] met）更容易被 claim。当依赖链断裂导致业务任务被标注 `[cycle]` 时，auto-gen 任务会抢占 claim 顺序

## Related Files

- `forge-cli/internal/cmd/task.go` — `forge task index` 和 `forge task claim` 实现
- `docs/features/slim-task-prompt-templates/tasks/index.json` — Task 4 的 dependencies 需要手动修复
- [[gotcha-split-rules-operational-blindness]] — 本次任务拆分的上游原因
