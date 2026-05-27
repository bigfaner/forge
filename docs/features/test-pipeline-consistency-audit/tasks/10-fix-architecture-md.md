---
id: "10"
title: "修复 ARCHITECTURE.md 与实际代码的对齐"
priority: "P1"
estimated_time: "1h"
dependencies: [3, 4]
type: "doc"
mainSession: false
---

# 10: 修复 ARCHITECTURE.md 与实际代码的对齐

## Description
修复 `ARCHITECTURE.md` 中的多处不一致：Quick 模式流程图移除 gen-contracts/gen-scripts（实际代码 `GetQuickTestTasks()` 跳过两者）；任务 ID 从 `T-test-1~5` 更新为描述性 ID；移除 `T-test-promote` 幽灵条目；"profile type" → "Convention" / "surface type"；"profile 路由" → "Convention 路由"；第 305 行并行执行描述按实际代码修正。

## Reference Files
- `proposal.md#Layer-2-Skill-文档层术语统一` — 第 12 项定义了 ARCHITECTURE.md 所有修复点
- `proposal.md#Success-Criteria` — 验证条件：Quick 模式流程图、任务 ID、术语替换
- `proposal.md#Risks` — ARCHITECTURE.md 更新引入新错误的风险，以 autogen.go 为 ground truth

## Acceptance Criteria
- [ ] Quick 模式流程图不含 gen-contracts/gen-scripts
- [ ] 任务 ID 为描述性名称（非 `T-test-1~5`）
- [ ] 不含 `T-test-promote` 条目
- [ ] "profile type" 已替换为 "Convention" / "surface type"
- [ ] "profile 路由" 已替换为 "Convention 路由"
- [ ] 第 305 行并行执行描述与 `autogen.go` 实际代码一致

## Hard Rules
- 以 `autogen.go` 代码为 ground truth，交叉验证每个流程图描述

## Implementation Notes
- 先读 `autogen.go` 确认 `GetQuickTestTasks()` 和 `GetFullTestTasks()` 的实际流程，再修改 ARCHITECTURE.md
- 流程图使用 Mermaid 语法

## Affected Files

### Create
| File | Description |
|------|-------------|

### Modify
| File | Changes |
|------|---------|
| `docs/ARCHITECTURE.md` | Quick 模式流程图、任务 ID、术语替换 |

### Delete
| File | Reason |
|------|--------|
