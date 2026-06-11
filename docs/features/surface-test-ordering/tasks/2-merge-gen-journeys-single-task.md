---
id: "2"
title: "Merge gen-journeys to single task"
priority: "P0"
estimated_time: "1.5h"
dependencies: []
surface-key: ""
surface-type: "cli"
breaking: false
type: "coding.feature"
mainSession: false
---

# 2: Merge gen-journeys to single task

## Description
移除 `autogen.go` 中 `GetBreakdownTestTasks` 和 `GetQuickTestTasks` 的 per-surface gen-journeys 循环，改为生成单个 `T-test-gen-journeys` 任务，内部遍历所有 surface type 加载对应规则。合并后 TestType 字段留空，`renderBody` 函数适配空 TestType 场景。

## Reference Files
- `proposal.md#Proposed-Solution` — gen-journeys 合并方案：从 per-surface 并行改为单任务，内部遍历所有 surface type
- `proposal.md#Requirements-Analysis` — 单 surface 退化场景（scalar 形式无后缀）和多 surface 验证
- `proposal.md#Key-Risks` — 单任务加载多 surface 规则增加 context 噪音的风险及缓解

## Acceptance Criteria
- [ ] gen-journeys 生成单个 `T-test-gen-journeys` 任务，输出覆盖所有配置 surface 的 Journey 文件
- [ ] 单 surface 项目（`surfaces: api`）退化为无后缀 `T-test-gen-journeys`，行为与改动前一致
- [ ] `renderBody` 正确处理空 TestType 字段

## Hard Rules
- 两个函数（`GetBreakdownTestTasks`、`GetQuickTestTasks`）均需改动

## Implementation Notes
- gen-journeys 以 PRD 为主要输入，surface 规则仅作参考指导（mandatory outcomes、test ratio），加载多份规则文件的噪音影响可忽略
- SKILL.md 中建议增加 `## Multi-Surface Rules Loading` 段落，按 surface type 分节组织规则
