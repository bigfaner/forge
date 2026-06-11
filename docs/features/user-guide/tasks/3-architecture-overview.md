---
id: "3"
title: "编写架构概览文档 architecture-overview.md"
priority: "P1"
estimated_time: "1.5h"
dependencies: []
type: "doc"
mainSession: false
---

# 3: 编写架构概览文档 architecture-overview.md

## Description
创建 `docs/user-guide/architecture-overview.md`，以用户视角介绍 Forge 的插件机制、组件角色（skill/command/agent/hook）、数据流向、状态管理和目录约定。区别于面向开发者的 `docs/ARCHITECTURE.md`（554 行，包含 Go 包结构等内部实现细节），本文档聚焦于"用户需要了解什么才能更好地使用 Forge"。

## Reference Files
- `docs/ARCHITECTURE.md`: 架构信息来源，提取用户视角内容，排除 Go 包结构等开发者细节 (source: proposal.md#Constraints-&-Dependencies)
- `README.md`: 组件角色（Skills、Commands、Agents）概览和目录结构 (source: proposal.md#Proposed-Solution)
- `docs/conventions/forge-distribution.md`: 分发模型和路径解析机制，帮助用户理解文件组织 (source: proposal.md#Constraints-&-Dependencies)
- `.forge/config.yaml`: 配置结构作为目录约定的一部分 (source: proposal.md#Proposed-Solution)

## Affected Files

### Create
| File | Description |
|------|-------------|
| `docs/user-guide/architecture-overview.md` | 用户视角架构概览：插件机制、组件角色、数据流、状态管理、目录约定 |

### Modify
| File | Changes |
|------|---------|
| (无) | |

### Delete
| File | Reason |
|------|--------|
| (无) | |

## Acceptance Criteria
- [ ] 包含插件机制说明（Claude Code 插件加载方式和 Forge 的定位）
- [ ] 包含四大组件角色表格（skill、command、agent、hook），每个有名称、用途、触发方式
- [ ] 包含数据流向图解（从用户输入 → Forge 处理 → 文件系统变更的可视化说明）
- [ ] 包含目录约定说明（`.forge/` 目录结构、`docs/features/` 结构、`manifest.md` 作用）
- [ ] 不包含 Go 包结构、CLI 内部命令注册、ResolveScope 等开发者内部实现细节

## Hard Rules
- 不包含 Go 包结构（如 `forge-cli/internal/cmd/`）、CLI 内部命令注册机制、或任何 `resolveScope` 等内部实现细节。本文档面向终端用户，只描述"是什么"和"怎么用"，不描述"怎么实现的"。

## Implementation Notes
- 文档使用中文
- 数据流向可以使用 Mermaid 图或 ASCII 图，确保在 Markdown 中可读
- 组件角色描述从 README.md 和 ARCHITECTURE.md 中提取，但过滤掉内部实现细节
- 目录约定需与 `.forge/` 目录实际结构一致
- 文档顶部标注"最后更新"日期和对应版本（v3.0.0）
