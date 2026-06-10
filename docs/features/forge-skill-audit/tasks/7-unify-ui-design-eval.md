---
id: "7"
title: "Unify ui-design auto.eval implementation (M-2)"
priority: "P2"
estimated_time: "30m"
dependencies: [5]
type: "doc"
mainSession: false
---

# 7: Unify ui-design auto.eval implementation (M-2)

## Description

ui-design 的 auto.eval 实现使用自然语言描述（依赖 LLM 解释），而其他三个 eval-capable skill（brainstorm、write-prd、tech-design）使用 bash script 模板实现三路分支（确定性执行）。需统一为 bash script 模板以减少变异性。

## Reference Files
- `docs/proposals/forge-skill-audit/proposal.md` — M-2: ui-design auto.eval 实现方式不一致
- `plugins/forge/skills/ui-design/SKILL.md`: Replace NL description with bash script template (ref: M-2: ui-design auto.eval 实现方式不一致)
- `plugins/forge/skills/brainstorm/SKILL.md` or `plugins/forge/skills/tech-design/SKILL.md`: Reference bash script template pattern

## Affected Files

### Create
| File | Description |
|------|-------------|

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/ui-design/SKILL.md` | Replace natural language auto.eval description with bash script template matching other eval-capable skills |

### Delete
| File | Reason |
|------|--------|

## Acceptance Criteria
- [ ] ui-design SKILL.md auto.eval 部分使用 bash script 模板（三路分支：disabled/skip/run），与 brainstorm/write-prd/tech-design 一致

## Hard Rules
- Before modifying any plugin file, read `docs/conventions/forge-distribution.md`
- Only modify markdown files; no Go code changes

## Implementation Notes
- 先读取 brainstorm 或 tech-design 的 bash script 模板作为参考，确保格式完全一致
