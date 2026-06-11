---
id: "5"
title: "Commands + Agent + Hooks deep audit"
priority: "P1"
estimated_time: "1.5h"
dependencies: [1]
type: "doc"
complexity: "high"
mainSession: false
---

# 5: Commands + Agent + Hooks deep audit

## Description

对 18 个 command（plugins/forge/commands/）、1 个 agent（task-executor）、和 hooks/guide.md 执行全部三层审计。Commands 审计聚焦内部流程步骤的一致性；Agent 审计检查指令之间的矛盾；hooks/guide.md 审计验证引用的脚本路径存在性、参数描述与脚本声明的一致性。

## Reference Files
- `docs/proposals/plugin-consistency-audit/proposal.md#审计方法论`: Layer 1-3 定义、分类标准 (source: proposal.md#审计方法论)
- `docs/proposals/plugin-consistency-audit/proposal.md#In-Scope`: hooks/guide.md 审计范围——验证脚本路径存在、参数描述一致、内部步骤无矛盾 (source: proposal.md#In-Scope)
- `docs/features/plugin-consistency-audit/reports/01-inventory-structural.md`: Task 1 产出的组件清单 (source: Task 1)
- `plugins/forge/commands/`: 18 个 command 文件 (source: proposal.md#审计协议步骤)
- `plugins/forge/agents/task-executor.md`: Agent 定义文件 (source: proposal.md#审计协议步骤)

## Affected Files

### Create
| File | Description |
|------|-------------|
| `docs/features/plugin-consistency-audit/reports/05-commands-agent-hooks.md` | Commands + Agent + Hooks 审计报告（Layer 1-3 发现） |

### Modify
| File | Changes |
|------|---------|

### Delete
| File | Reason |
|------|--------|

## Acceptance Criteria
- [ ] 全部 18 个 command 文件已读取，内部流程步骤的时序和引用一致性已验证
- [ ] task-executor agent 文件已读取，指令之间的矛盾、冗余、时序问题已检查
- [ ] hooks/guide.md 已读取，引用的脚本路径存在性已验证（Layer 1 REFERENCE）
- [ ] hooks/guide.md 中脚本参数描述与实际脚本声明的一致性已验证（Layer 2 CONFLICT）
- [ ] 每条问题按报告 schema 记录：`{component, file_path, layer, category, severity, description, fix_suggestion, confidence}`

## Implementation Notes
- hooks/guide.md 是"单一组件自洽"原则的明确例外——它是跨 hook 脚本的索引文档。审计范围限于：(1) 脚本路径存在性；(2) 参数描述一致性；(3) 内部步骤无矛盾。不深入验证脚本文件的逻辑正确性
- Commands 通常较轻量（单文件），审计重点在内部流程逻辑（步骤时序、条件分支覆盖）
- 每个 finding 标注置信度（high/medium/low）
