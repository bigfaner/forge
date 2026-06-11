---
id: "10"
title: "Fix gen-contracts path references (HIGH-A1, HIGH-A2, MEDIUM-A3)"
priority: "P0"
estimated_time: "30m"
dependencies: []
type: "doc"
mainSession: false
---

# 10: Fix gen-contracts path references

## Description

gen-contracts/SKILL.md 存在 3 处路径问题：(1) Handbook Location 表中 5 个 handbook 路径缺少 `design/` 子目录前缀；(2) freshness check 引用的 tech-design 路径缺少 `design/` 前缀；(3) INLINE 内容引用 gen-journeys 的 surface-*.md 文件但 gen-contracts 目录下不存在这些文件。

## Reference Files
- `docs/proposals/forge-skill-audit/proposal.md` — Proposed Solution, Scope
- `plugins/forge/skills/gen-contracts/SKILL.md`: Fix handbook paths, freshness check path, and INLINE surface file references (ref: Proposed Solution)
- `plugins/forge/skills/tech-design/SKILL.md`: Verify actual output path for handbooks and tech-design

## Affected Files

### Create
| File | Description |
|------|-------------|

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/gen-contracts/SKILL.md` | Add `design/` prefix to handbook paths; fix freshness check tech-design path; clarify INLINE surface references |

### Delete
| File | Reason |
|------|--------|

## Acceptance Criteria
- [ ] Handbook Location 表中所有路径包含 `design/` 子目录前缀（如 `docs/features/<slug>/design/api-handbook.md`）
- [ ] freshness check 引用 `docs/features/<slug>/design/tech-design.md`
- [ ] INLINE 区块中 surface-*.md 引用添加说明标注这些文件属于 gen-journeys（非本 skill 目录）

## Hard Rules
- Before modifying any plugin file, read `docs/conventions/forge-distribution.md`
- 修改前先用 grep 验证 tech-design skill 实际输出路径（确认 handbook 和 tech-design 确实生成到 `design/` 子目录）

## Implementation Notes
- HIGH-A1: 5 个 handbook 路径（api-handbook, cli-handbook, page-map, screen-map 等）全部缺少 `design/`
- HIGH-A2: 第 158 行 freshness check 路径 `docs/features/<slug>/tech-design.md` 应为 `docs/features/<slug>/design/tech-design.md`
- MEDIUM-A3: 第 63 行 INLINE 区块包含 `Load rules/surface-<type>.md` 指令，但 gen-contracts 无此文件；需标注来源为 gen-journeys
