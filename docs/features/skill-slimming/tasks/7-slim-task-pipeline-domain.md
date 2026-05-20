---
id: "7"
title: "Slim task pipeline domain (breakdown-tasks + quick-tasks + submit-task)"
priority: "P1"
estimated_time: "30-45min"
dependencies: ["6"]
type: "doc"
mainSession: false
---

# 7: Slim task pipeline domain (breakdown-tasks + quick-tasks + submit-task)

## Description
对任务管线域的 3 个 skill 进行精简：breakdown-tasks（144 行）、quick-tasks（208 行）、submit-task（156 行）。这些文件较小，主要做精简和消歧，不一定需要拆分。

## Reference Files
- `docs/proposals/skill-slimming/proposal.md` — Splitting Heuristic
- `docs/conventions/skill-self-containment.md` — 自洽原则

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/breakdown-tasks/SKILL.md` | 精简冗余、消除歧义 |
| `plugins/forge/skills/quick-tasks/SKILL.md` | 精简冗余、消除歧义 |
| `plugins/forge/skills/submit-task/SKILL.md` | 精简冗余、消除歧义 |

## Acceptance Criteria
- [ ] 每个 SKILL.md 行数 ≤ 350 行
- [ ] 无歧义术语残留（noTest、doc* 等歧义项已替换为精确引用）
- [ ] 引用的辅助文件路径均存在可读

## Hard Rules
- 遵守 Splitting Heuristic
- 不改变 skill 的输入/输出契约

## Implementation Notes
- breakdown-tasks 有 7 个辅助文件（350 行），quick-tasks 有 3 个（93 行），submit-task 无辅助文件
- 这些文件较小（144-208 行），重点在消歧和精简冗余文本，而非拆分
