---
id: "2"
title: "编写初始化文档 initialization.md"
priority: "P1"
estimated_time: "1.5h"
dependencies: []
type: "doc"
mainSession: false
---

# 2: 编写初始化文档 initialization.md

## Description
创建 `docs/user-guide/initialization.md`，详细说明 `forge init` 的完整流程、`.forge/config.yaml` 全字段含义、Surface 检测机制和首个项目设置。当前这些关键初始化知识分散在多个文件中，没有独立的初始化指南。

## Reference Files
- `.forge/config.yaml`: 当前项目的配置文件，包含所有配置项的实际值 (source: proposal.md#Proposed-Solution)
- `plugins/forge/commands/init-forge.md`: forge init 命令的完整流程 (source: proposal.md#Proposed-Solution)
- `docs/ARCHITECTURE.md`: Surface Detection 部分，理解 Surface 检测机制 (source: proposal.md#Proposed-Solution)
- `docs/reference/test-type-model.md`: Surface Type 与 Test Type 的映射关系 (source: proposal.md#Proposed-Solution)
- `README.md`: 命令参考中 init 相关命令 (source: proposal.md#Problem)

## Affected Files

### Create
| File | Description |
|------|-------------|
| `docs/user-guide/initialization.md` | 初始化指南：forge init 流程、config.yaml 全字段参考、Surface 检测 |

### Modify
| File | Changes |
|------|---------|
| (无) | |

### Delete
| File | Reason |
|------|--------|
| (无) | |

## Acceptance Criteria
- [ ] 包含 `forge init` 的完整流程说明（从命令执行到项目就绪）
- [ ] 包含 config.yaml 全字段表格，至少 8 个配置项，每个字段有名称、类型、默认值、说明
- [ ] 包含 Surface 检测机制说明（`forge surfaces detect` 的使用和结果解读）
- [ ] 包含首个项目设置的端到端示例（从 init 到可以开始使用 Forge）
- [ ] 所有代码示例可直接复制执行，无需额外修改

## Implementation Notes
- 文档使用中文
- config.yaml 字段表格需从 `.forge/config.yaml` 的实际结构提取，不要遗漏配置项
- Surface 检测说明需覆盖 `forge surfaces detect` 和 `forge surfaces detect --apply` 两种用法
- 文档顶部标注"最后更新"日期和对应版本（v3.0.0）
