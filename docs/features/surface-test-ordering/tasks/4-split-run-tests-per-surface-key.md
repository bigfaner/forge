---
id: "4"
title: "Split run-tests into per-surface-key serial tasks"
priority: "P0"
estimated_time: "2h"
dependencies: ["1"]
surface-key: ""
surface-type: "cli"
breaking: true
type: "coding.feature"
mainSession: false
---

# 4: Split run-tests into per-surface-key serial tasks

## Description
在 `autogen.go` 中将 `T-test-run` 拆分为 `T-test-run-{surface-key}` 按 execution-order 串行排列。函数签名从 `capabilities []string` 改为接收 surface-key 列表或 surfaces map，所有调用方适配。单 surface 退化为无后缀 `T-test-run`。`T-test-verify-regression` 依赖链尾（execution-order 中最后一个 run-test 子任务）。失败传播通过串行依赖天然实现：上游失败 → 下游 blocked。

## Reference Files
- `proposal.md#Proposed-Solution` — per-surface-key 串行任务模型、3-surface 依赖链示例、单 surface 退化规则
- `proposal.md#Constraints-&-Dependencies` — 函数签名变更影响范围、surface-key 合法性约束
- `proposal.md#Requirements-Analysis` — Key Scenarios：fullstack 默认排序、单 surface 退化、上游失败传播
- `proposal.md#Key-Risks` — 函数签名变更影响所有调用方、gen-scripts type 后缀与 run-tests key 后缀并存

## Acceptance Criteria
- [ ] 配置 `surfaces: { frontend: web, backend: api }` 且无 `execution-order` 时，`T-test-run-backend` 排在 `T-test-run-frontend` 之前
- [ ] `T-test-run-backend` 失败时，`T-test-run-frontend` 状态为 blocked，不执行
- [ ] 单 surface 项目（`surfaces: api`）退化为无后缀 `T-test-run`，任务 ID 和依赖列表与改动前一致
- [ ] Quick 模式：`T-test-gen-journeys` 为 `T-test-run-*` 的直接上游
- [ ] `T-test-verify-regression` 依赖 execution-order 中最后一个 run-test 子任务

## Hard Rules
- 函数签名变更：`capabilities []string` → `surfaceKeys []string` 或新增 surfaces map 参数，逐个函数修改，编译器类型检查确保无遗漏

## Implementation Notes
- 失败传播不需要额外实现——串行依赖链天然实现 blocked 状态传播
- gen-scripts 保持 per-surface 并行（type 后缀 `-api`、`-web`），run-tests 使用 key 后缀（`-backend`、`-frontend`），两套命名共存
- run-tests 的函数签名变更影响 `BuildIndex` 及相关入口函数
