---
id: "3"
title: "Skills deep audit - batch B (quick-tasks, consolidate-specs, clean-code, deep-research, forensic, ui-design, learn)"
priority: "P1"
estimated_time: "2h"
dependencies: [1]
type: "doc"
complexity: "high"
mainSession: false
---

# 3: Skills deep audit - batch B

## Description

对 7 个 skill（quick-tasks, consolidate-specs, clean-code, deep-research, forensic, ui-design, learn）执行 Layer 2（指令一致性）和 Layer 3（时序与流程）深度审计。按 proposal 的"逐组件多轮对话"模式执行。

## Reference Files
- `docs/proposals/plugin-consistency-audit/proposal.md#审计方法论`: Layer 2-3 定义、CONFLICT/REDUNDANT/INCOMPLETE/TIMING 分类、关键词强度映射表 (source: proposal.md#审计方法论)
- `docs/proposals/plugin-consistency-audit/proposal.md#Severity-Level-Definitions`: P0-P3 严重等级定义 (source: proposal.md#Non-Functional-Requirements)
- `docs/features/plugin-consistency-audit/reports/01-inventory-structural.md`: Task 1 产出的组件清单 (source: Task 1)
- `plugins/forge/skills/{quick-tasks,consolidate-specs,clean-code,deep-research,forensic,ui-design,learn}/`: 7 个 skill 目录 (source: proposal.md#审计协议步骤)

## Affected Files

### Create
| File | Description |
|------|-------------|
| `docs/features/plugin-consistency-audit/reports/03-skills-batch-b.md` | Skills batch B 审计报告（Layer 2-3 发现） |

### Modify
| File | Changes |
|------|---------|

### Delete
| File | Reason |
|------|--------|

## Acceptance Criteria
- [ ] 7 个 skill 各自的 SKILL.md 已全文读取并提取结构化摘要（步骤、约束、引用路径、字段名）
- [ ] 每个 skill 的所有关联文件（templates/rules/data）已逐一与 SKILL.md 摘要比对
- [ ] 使用关键词强度映射表检查"必须/应该/可选/禁止"的一致性，记录 CONFLICT 类问题
- [ ] 多步骤组件的步骤时序已验证（Layer 3），TIMING 类问题已记录
- [ ] 每条问题按报告 schema 记录：`{component, file_path, layer, category, severity, description, fix_suggestion, confidence}`

## Hard Rules
- 仅审计以下 7 个 skill 目录：quick-tasks, consolidate-specs, clean-code, deep-research, forensic, ui-design, learn。不修改任何文件。

## Implementation Notes
- 冗余检测启发式：同一语义在 ≥2 个文件中出现且无信息增量才标记 REDUNDANT；rules 中对 SKILL.md 的合理展开不算冗余
- 每个 finding 标注置信度（high/medium/low）
- INCOMPLETE 类型：rules/templates 中存在但 SKILL.md 未提及的约束
