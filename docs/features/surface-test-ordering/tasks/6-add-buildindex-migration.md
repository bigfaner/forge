---
id: "6"
title: "Add BuildIndex migration for existing fix-tasks"
priority: "P1"
estimated_time: "1h"
dependencies: ["4"]
surface-key: ""
surface-type: "cli"
breaking: false
type: "coding.enhancement"
mainSession: false
---

# 6: Add BuildIndex migration for existing fix-tasks

## Description
在 `BuildIndex` 阶段（非 task 生成函数），检测 `index.json` 中已有 `SourceTaskID: "T-test-run"` 的 fix-tasks，自动重映射到对应的 `T-test-run-{surface-key}`。多 surface 项目将旧 `T-test-run` 的 status/blocked-reason 复制到 `T-test-run-{execution-order 首个 surface-key}`。单 surface 项目保持 `T-test-run` 不变，无迁移成本。

## Reference Files
- `proposal.md#Feasibility-Assessment` — 迁移逻辑位置：BuildIndex 阶段拥有 index.json 读写权限和旧任务状态上下文
- `proposal.md#Key-Risks` — index.json 已有 T-test-run 条目变为孤儿的风险，迁移策略

## Acceptance Criteria
- [ ] 多 surface 项目 index.json 中已有 `SourceTaskID: "T-test-run"` 的 fix-tasks 在 BuildIndex 阶段自动重映射到 `T-test-run-{execution-order 首个 surface-key}`
- [ ] 单 surface 项目（`surfaces: api`）无迁移行为，`T-test-run` 保持不变
- [ ] 迁移后旧 `T-test-run` 条目从 index.json 中移除（不再为孤儿）

## Hard Rules
- 迁移逻辑仅在 `BuildIndex` 阶段执行，不在 `GetBreakdownTestTasks` 或 `GetQuickTestTasks` 中

## Implementation Notes
- `BuildIndex` 拥有 index.json 读写权限和旧任务状态上下文，适合放置迁移逻辑
- 多 surface 项目需按 execution-order 首个 surface 继承旧任务状态
