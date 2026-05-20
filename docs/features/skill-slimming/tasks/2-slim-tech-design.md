---
id: "2"
title: "Slim tech-design (472→≤350 lines)"
priority: "P1"
estimated_time: "1-2h"
dependencies: ["1"]
type: "doc"
mainSession: false
---

# 2: Slim tech-design (472→≤350 lines)

## Description
按 Splitting Heuristic 规则，将 `tech-design/SKILL.md`（当前 472 行）拆分为 SKILL.md + rules/ + templates/。参照 Task 1 的拆分结构和粒度保持一致性。

## Reference Files
- `docs/proposals/skill-slimming/proposal.md` — Splitting Heuristic
- `docs/conventions/skill-self-containment.md` — 自洽原则

## Affected Files

### Create
| File | Description |
|------|-------------|
| `plugins/forge/skills/tech-design/rules/*.md` | 规则细节（评估标准、约束条件等） |
| `plugins/forge/skills/tech-design/templates/*.md` | 输出模板（如有内联模板） |

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/tech-design/SKILL.md` | 保留流程骨架，移除规则细节和模板 |

## Acceptance Criteria
- [ ] SKILL.md 行数 ≤ 350 行
- [ ] 所有步骤编号及其描述保留
- [ ] 条件分支逻辑和 I/O 契约保留
- [ ] 引用的 rules/templates 路径均存在可读
- [ ] 拆分风格与 Task 1 保持一致

## Hard Rules
- 遵守 Splitting Heuristic
- 不改变 skill 的输入/输出契约

## Implementation Notes
- tech-design 有 8 个辅助文件（templates/、rules/、examples/），部分已存在——检查是否可复用现有目录结构，避免重复创建
