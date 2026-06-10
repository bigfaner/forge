---
id: "2"
title: "Remove tech-design dead path (H-2)"
priority: "P1"
estimated_time: "30m"
dependencies: [1]
type: "doc"
mainSession: false
---

# 2: Remove tech-design dead path (H-2)

## Description

tech-design SKILL.md 引用从未被任何 skill 创建的 `docs/features/<slug>/proposal.md` 作为第一个 intent 读取路径，但实际 proposal 位于 `docs/proposals/<slug>/proposal.md`。LLM 浪费上下文查找不存在的文件。

## Reference Files
- `docs/proposals/forge-skill-audit/proposal.md` — H-2: tech-design intent 读取路径误导, Proposed Solution, Success Criteria
- `plugins/forge/skills/tech-design/SKILL.md`: Remove dead path (ref: H-2: tech-design intent 读取路径误导)

## Affected Files

### Create
| File | Description |
|------|-------------|

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/tech-design/SKILL.md` | Remove `docs/features/<slug>/proposal.md` from intent read paths |

### Delete
| File | Reason |
|------|--------|

## Acceptance Criteria
- [ ] tech-design SKILL.md 仅引用 `docs/proposals/<slug>/proposal.md`，无 `docs/features/<slug>/proposal.md` 路径
- [ ] 搜索所有 skill 确认无其他 skill 引用 `docs/features/<slug>/proposal.md` 死路径

## Hard Rules
- Before modifying any plugin file, read `docs/conventions/forge-distribution.md`
- Only modify markdown files; no Go code changes
- 修复前先搜索确认 `docs/features/<slug>/proposal.md` 无并存路径（注意 `docs/features/` 下有 200+ 目录用于其他文件如 PRD）

## Implementation Notes
- `docs/features/` 路径被多 skill 广泛使用（如 tech-design 的 Prerequisites 检查 `docs/features/<slug>/prd/prd-spec.md`），仅移除 proposal 文件的死路径引用，不影响其他用途
