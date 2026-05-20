---
id: "9"
title: "Slim utility domain (clean-code + run-e2e-tests)"
priority: "P1"
estimated_time: "30min"
dependencies: ["8"]
type: "doc"
mainSession: false
---

# 9: Slim utility domain (clean-code + run-e2e-tests)

## Description
对工具域的 2 个 skill 进行精简：clean-code（190 行）、run-e2e-tests（193 行）。这是最后一个任务组。

## Reference Files
- `docs/proposals/skill-slimming/proposal.md` — Splitting Heuristic
- `docs/conventions/skill-self-containment.md` — 自洽原则

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/clean-code/SKILL.md` | 精简冗余、消除歧义 |
| `plugins/forge/skills/run-e2e-tests/SKILL.md` | 精简冗余、消除歧义 |

## Acceptance Criteria
- [ ] 每个 SKILL.md 行数 ≤ 350 行
- [ ] 无歧义术语残留
- [ ] 引用的辅助文件路径均存在可读

## Hard Rules
- 遵守 Splitting Heuristic
- 不改变 skill 的输入/输出契约

## Implementation Notes
- clean-code 有 1 个辅助文件（16 行），run-e2e-tests 有 3 个（147 行）
- 两个文件都较小，重点在消歧和清理
- 作为最后一个任务，完成后验证所有 22 个 skill 的总行数是否达到 25%+ 减少目标
