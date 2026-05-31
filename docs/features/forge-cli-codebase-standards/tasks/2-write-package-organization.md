---
id: "2"
title: "Write package-organization.md with PR review checklist"
priority: "P1"
estimated_time: "3h"
dependencies: [1]
type: "doc"
mainSession: false
---

# 2: Write package-organization.md with PR review checklist

## Description
新增 `docs/conventions/package-organization.md`，定义 `internal/cmd/` 和 `pkg/` 的包职责划分、依赖方向规则（`cmd → internal → pkg`）、`pkg/` 内三层层级模型（leaf/基础设施/领域）、文件组织原则。包含目标态定义（规范性）和基于 Task 1 依赖图的偏差分析。附加 PR review checklist（包结构变更需 review 确认符合依赖方向规则和包职责定义）。

## Reference Files
- docs/features/forge-cli-codebase-standards/pkg-dependency-graph.md: 依赖图事实基线，偏差分析须引用此文件中的具体数据 (source: proposal.md#Proposed-Solution Phase 1)
- forge-cli/internal/cmd/: 15 个顶层命令文件 + 7 个子包，需描述目标态结构 (source: proposal.md#Scope item 7)
- forge-cli/pkg/: 17 个包的目标态映射 (source: proposal.md#Scope item 8)

## Affected Files

### Create
| File | Description |
|------|-------------|
| docs/conventions/package-organization.md | 包组织规范，含依赖方向、三层模型、偏差分析、PR review checklist |

### Modify
| File | Changes |
|------|---------|

### Delete
| File | Reason |
|------|--------|

## Acceptance Criteria
- [ ] `docs/conventions/package-organization.md` 存在，包含目标态定义（非描述性）和偏差分析表
- [ ] 依赖方向规则明确：`cmd → internal → pkg`（严格单向），`pkg/` 内三层模型定义清晰（leaf/基础设施/领域）
- [ ] 偏差分析表引用 Task 1 依赖图的具体数据（如 `pkg/infocmd` 被 4 个领域包导入的事实）
- [ ] PR review checklist 包含：包结构变更需 review 确认符合依赖方向规则和包职责定义
- [ ] 开发者工作流描述完整：新增命令时应在 `internal/cmd/<command-group>/` 下创建文件

## Implementation Notes
- 偏差分析格式：每行一个偏差项，包含 [当前状态] → [目标状态] → [差距描述]
- PR review checklist 应可复制粘贴到 GitHub PR template
