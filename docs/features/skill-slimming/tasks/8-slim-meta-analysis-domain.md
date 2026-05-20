---
id: "8"
title: "Slim meta/analysis domain (brainstorm + learn + forensic + improve-harness)"
priority: "P1"
estimated_time: "30-45min"
dependencies: ["7"]
type: "doc"
mainSession: false
---

# 8: Slim meta/analysis domain (brainstorm + learn + forensic + improve-harness)

## Description
对元分析域的 4 个 skill 进行精简：brainstorm（139 行）、learn（259 行）、forensic（198 行）、improve-harness（163 行）。

## Reference Files
- `docs/proposals/skill-slimming/proposal.md` — Splitting Heuristic
- `docs/conventions/skill-self-containment.md` — 自洽原则

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/brainstorm/SKILL.md` | 精简冗余、消除歧义 |
| `plugins/forge/skills/learn/SKILL.md` | 精简冗余、消除歧义 |
| `plugins/forge/skills/forensic/SKILL.md` | 精简冗余、消除歧义 |
| `plugins/forge/skills/improve-harness/SKILL.md` | 精简冗余、消除歧义 |

## Acceptance Criteria
- [ ] 每个 SKILL.md 行数 ≤ 350 行
- [ ] 无歧义术语残留
- [ ] 引用的辅助文件路径均存在可读

## Hard Rules
- 遵守 Splitting Heuristic
- 不改变 skill 的输入/输出契约

## Implementation Notes
- brainstorm 有 1 个辅助文件（101 行），learn 有 3 个（183 行），forensic 有 1 个（86 行），improve-harness 有 1 个（59 行）
- 所有文件都 ≤ 259 行，重点在消歧和清理冗余
