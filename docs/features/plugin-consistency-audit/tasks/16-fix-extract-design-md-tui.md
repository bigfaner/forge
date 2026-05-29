---
id: "16"
title: "Fix: extract-design-md missing TUI match strategy in SKILL.md"
priority: "P1"
estimated_time: "15min"
dependencies: []
type: "doc"
complexity: "low"
mainSession: false
---

# 16: Fix: extract-design-md missing TUI match strategy in SKILL.md

## Description
extract-design-md SKILL.md Step 3 的 match strategy 仅提到 web built-in styles（5 个），未提到 TUI themes（modern-dark-tui / minimal-ascii-tui）。但 `rules/platform-routing.md` 为 TUI 平台提供了独立的 match strategy。SKILL.md Step 3 应明确区分 web/mobile 和 TUI 的 match 流程。

## Reference Files
- `docs/features/plugin-consistency-audit/reports/04-skills-batch-c.md#C-05`: P1 CONFLICT, TUI match strategy 仅在 rules 中 (source: Report 04)
- `plugins/forge/skills/extract-design-md/SKILL.md`: 需修改的 Step 3 (source: audit finding)
- `plugins/forge/skills/extract-design-md/rules/platform-routing.md`: TUI match strategy 定义 (source: audit finding)

## Affected Files

### Create
| File | Description |
|------|-------------|

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/extract-design-md/SKILL.md` | 在 Step 3 添加 TUI match strategy 描述 |

### Delete
| File | Reason |
|------|--------|

## Acceptance Criteria
- [ ] Step 3 明确区分 web/mobile 使用 5 个 web built-in styles
- [ ] Step 3 明确 TUI 使用 2 个 TUI themes（modern-dark-tui / minimal-ascii-tui）
- [ ] TUI match strategy 引用 `rules/platform-routing.md` 的完整定义

## Hard Rules
- 仅修改 `plugins/forge/skills/extract-design-md/SKILL.md`
