---
id: "1"
title: "Slim consolidate-specs (348→≤200 lines)"
priority: "P1"
estimated_time: "1-2h"
dependencies: []
type: "doc"
mainSession: false
---

# 1: Slim consolidate-specs (348→≤200 lines)

## Description
按 proposal 中的 Splitting Heuristic 规则，精简 `consolidate-specs/SKILL.md`（当前 348 行 + 282 行辅助文件）。重点是精简冗余文本和消除歧义。这是 Tier 1 的第一个任务，处理风格将作为后续 task 的参考。

## Reference Files
- `docs/proposals/skill-slimming/proposal.md` — Source proposal (Splitting Heuristic 节)
- `docs/conventions/skill-self-containment.md` — 自洽原则约束

## Affected Files

### Create
| File | Description |
|------|-------------|
| `plugins/forge/skills/consolidate-specs/rules/*.md` | 从 SKILL.md 移出的规则细节和术语定义（如需新建） |

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/consolidate-specs/SKILL.md` | 保留流程骨架，移除规则细节和模板，精简冗余文本，消除歧义 |

## Acceptance Criteria
- [ ] SKILL.md 行数 ≤ 350 行
- [ ] 所有步骤编号及其描述保留在 SKILL.md 中
- [ ] 所有条件分支逻辑保留在 SKILL.md 中
- [ ] 输入/输出契约定义保留在 SKILL.md 中
- [ ] SKILL.md 中引用的所有 rules/templates 路径存在且文件可读
- [ ] 无流程步骤遗漏

## Hard Rules
- 遵守 Splitting Heuristic：步骤编号+条件分支+I/O 契约留在 SKILL.md，>5 行规则定义移至 rules/，>10 行模板移至 templates/
- 边界规则：混合流程+规则的内容，流程留 SKILL.md，规则移 rules/ 并添加引用
- 不改变 skill 的输入/输出契约

## Implementation Notes
- 参照 proposal 中 Worked Example: consolidate-specs 的拆分结构
- commit message 需注明拆分了哪些内容到哪里
