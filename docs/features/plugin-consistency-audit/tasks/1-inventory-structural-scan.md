---
id: "1"
title: "Component inventory and Layer 1 structural scan"
priority: "P0"
estimated_time: "1h"
dependencies: []
type: "doc"
complexity: "high"
mainSession: false
---

# 1: Component inventory and Layer 1 structural scan

## Description

枚举 forge plugin 全部 41 个组件（21 skills + 18 commands + 1 agent + hooks/guide.md），生成每个组件的完整文件清单，并执行 Layer 1 结构完整性检查（验证所有引用路径存在、识别孤立文件）。产出 REFERENCE 类问题和组件清单，为后续所有审计任务提供基础。

## Reference Files
- `docs/proposals/plugin-consistency-audit/proposal.md#审计方法论`: Layer 1 定义、REFERENCE 失败类型、路径正则提取方法 (source: proposal.md#审计方法论)
- `plugins/forge/skills/`: 21 个 skill 目录，需逐一扫描 SKILL.md 及 templates/rules/data/examples/types (source: proposal.md#审计协议步骤)
- `plugins/forge/commands/`: 18 个 command 文件 (source: proposal.md#审计协议步骤)
- `plugins/forge/agents/task-executor.md`: Agent 定义文件 (source: proposal.md#审计协议步骤)
- `plugins/forge/hooks/guide.md`: Hook guide 文件，含脚本路径引用 (source: proposal.md#In-Scope)

## Affected Files

### Create
| File | Description |
|------|-------------|
| `docs/features/plugin-consistency-audit/reports/01-inventory-structural.md` | 组件清单 + Layer 1 结构完整性检查报告 |

### Modify
| File | Changes |
|------|---------|

### Delete
| File | Reason |
|------|--------|

## Acceptance Criteria
- [ ] 全部 21 个 skill 枚举完成，每个 skill 列出完整文件清单（SKILL.md + templates/ + rules/ + data/ + examples/ + types/）
- [ ] 全部 18 个 command 枚举完成，列出其内部文件引用
- [ ] 1 个 agent (task-executor) 枚举完成，列出引用文件
- [ ] hooks/guide.md 枚举完成，列出引用的脚本路径
- [ ] Layer 1 完成：SKILL.md 中引用的每个路径已与文件系统交叉验证，REFERENCE 类问题已记录
- [ ] 孤立文件（存在于目录中但未被 SKILL.md 引用）已识别并记录
- [ ] 报告包含基准 commit hash

## Implementation Notes
- 报告开头记录基准 commit hash（`git rev-parse HEAD`），确保审计可复现
- 对 SKILL.md 使用正则提取路径引用模式（templates/xxx.md, rules/xxx.md 等），逐个 `ls` 验证存在性
- 比对方法：提取引用路径集合 vs 文件系统实际文件集合，差集即为 REFERENCE 问题或孤立文件
- 按 proposal 定义的报告 schema 输出每条问题：`{component, file_path, layer, category, severity, description, fix_suggestion}`
