---
id: "2"
title: "Slim tech-design (190→≤150 lines)"
priority: "P1"
estimated_time: "1-2h"
dependencies: ["1"]
type: "doc"
mainSession: false
---

# 2: Slim tech-design (190→≤150 lines)

## Description
按 Splitting Heuristic 规则，精简 `tech-design/SKILL.md`（当前 190 行 + 649 行辅助文件）。重点在精简冗余文本和消除歧义。参照 Task 1 的处理风格保持一致性。

## Reference Files
- `docs/proposals/skill-slimming/proposal.md` — Splitting Heuristic
- `docs/conventions/skill-self-containment.md` — 自洽原则

## Affected Files

### Create
| File | Description |
|------|-------------|
| `plugins/forge/skills/tech-design/rules/*.md` | 规则细节（如有需从 SKILL.md 移出的内容） |

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/tech-design/SKILL.md` | 精简冗余、消除歧义 |

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
- tech-design 有 10 个辅助文件（649 行），已有丰富的 templates/ 和 rules/ 目录——检查是否可复用现有目录结构，避免重复创建
