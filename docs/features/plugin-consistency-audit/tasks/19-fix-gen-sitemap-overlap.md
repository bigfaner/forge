---
id: "19"
title: "Fix: gen-sitemap Step 2b/4 overlap handling"
priority: "P1"
estimated_time: "15min"
dependencies: []
type: "doc"
complexity: "low"
mainSession: false
---

# 19: Fix: gen-sitemap Step 2b/4 overlap handling

## Description
gen-sitemap SKILL.md 的 Step 2b 使用 agent-browser 访问 3-5 个 wrapped routes（用于 layout identification），而 Step 4 "Explore Pages" 对每条路由逐个探索。如果 Step 2b 已探索了 5 条路由，Step 4 是否应跳过这 5 条？SKILL.md 未说明已探索页面在 Step 4 中的处理方式。

## Reference Files
- `docs/features/plugin-consistency-audit/reports/04-skills-batch-c.md#T-03`: P1 TIMING, Step 2b/4 重叠 (source: Report 04)
- `plugins/forge/skills/gen-sitemap/SKILL.md`: 需修改的 Step 4 (source: audit finding)

## Affected Files

### Create
| File | Description |
|------|-------------|

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/gen-sitemap/SKILL.md` | 在 Step 4 中明确说明 Step 2b 已探索页面的处理方式 |

### Delete
| File | Reason |
|------|--------|

## Acceptance Criteria
- [ ] Step 4 明确说明 Step 2b 已探索的页面如何处理（复用结果或跳过 layout comparison）
- [ ] Step 4 与 Step 2b 之间无重复探索

## Hard Rules
- 仅修改 `plugins/forge/skills/gen-sitemap/SKILL.md`
