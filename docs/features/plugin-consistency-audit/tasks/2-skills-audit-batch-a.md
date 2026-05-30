---
id: "2"
title: "Skills deep audit - batch A (eval, gen-test-scripts, run-tests, tech-design, write-prd, brainstorm, breakdown-tasks)"
priority: "P0"
estimated_time: "2h"
dependencies: [1]
type: "doc"
complexity: "high"
mainSession: false
---

# 2: Skills deep audit - batch A

## Description

对 7 个 skill（eval, gen-test-scripts, run-tests, tech-design, write-prd, brainstorm, breakdown-tasks）执行 Layer 2（指令一致性）和 Layer 3（时序与流程）深度审计。按 proposal 的"逐组件多轮对话"模式：第一轮提取 SKILL.md 结构化摘要，后续轮次逐一比对关联文件，汇总轮去重分级。

包含 run-tests skill 作为**审计有效性基线验证**——已知其 `rules/env-check.md` 残留 Playwright 硬编码（P1 级 CONFLICT），审计必须能复现此问题。

## Reference Files
- `docs/proposals/plugin-consistency-audit/proposal.md#审计方法论`: Layer 2-3 定义、CONFLICT/REDUNDANT/INCOMPLETE/TIMING 分类、关键词强度映射表 (source: proposal.md#审计方法论)
- `docs/proposals/plugin-consistency-audit/proposal.md#Severity-Level-Definitions`: P0-P3 严重等级定义 (source: proposal.md#Non-Functional-Requirements)
- `docs/features/plugin-consistency-audit/reports/01-inventory-structural.md`: Task 1 产出的组件清单和 REFERENCE 类问题 (source: Task 1)
- `plugins/forge/skills/{eval,gen-test-scripts,run-tests,tech-design,write-prd,brainstorm,breakdown-tasks}/`: 7 个 skill 目录 (source: proposal.md#审计协议步骤)
- `docs/proposals/plugin-consistency-audit/proposal.md#Evidence`: run-tests/env-check.md 第 49 行 Playwright 硬编码的已知问题实例 (source: proposal.md#Evidence)

## Affected Files

### Create
| File | Description |
|------|-------------|
| `docs/features/plugin-consistency-audit/reports/02-skills-batch-a.md` | Skills batch A 审计报告（Layer 2-3 发现） |

### Modify
| File | Changes |
|------|---------|

### Delete
| File | Reason |
|------|--------|

## Acceptance Criteria
- [ ] 7 个 skill 各自的 SKILL.md 已全文读取并提取结构化摘要（步骤、约束、引用路径、字段名）
- [ ] 每个 skill 的所有关联文件（templates/rules/data/examples/types）已逐一与 SKILL.md 摘要比对
- [ ] 使用关键词强度映射表检查"必须/应该/可选/禁止"的一致性，记录 CONFLICT 类问题
- [ ] 多步骤组件的步骤时序已验证（Layer 3），TIMING 类问题已记录
- [ ] **有效性验证**: run-tests skill 的 `rules/env-check.md` Playwright 硬编码已被识别为 P1 级 CONFLICT
- [ ] 每条问题按报告 schema 记录：`{component, file_path, layer, category, severity, description, fix_suggestion, confidence}`

## Hard Rules
- 仅审计以下 7 个 skill 目录：eval, gen-test-scripts, run-tests, tech-design, write-prd, brainstorm, breakdown-tasks。不修改任何文件。

## Implementation Notes
- eval skill 是最大组件（36+ 文件），需分批审计：按子目录分组（rubrics/, experts/, templates/, rules/），每批独立完成 Layer 2-3，最后执行汇总轮进行跨组比对
- 冗余检测启发式：同一语义在 ≥2 个文件中出现且无信息增量才标记 REDUNDANT；rules 中对 SKILL.md 的合理展开（具体示例、边界条件）不算冗余
- 每个 finding 标注置信度（high/medium/low）
- 随机化审计顺序（非字母序），减少系统性偏差
