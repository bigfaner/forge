---
id: "4"
title: "eval/validate-ux-pipeline: Add explicit surface guard for web sitemap reference"
priority: "P1"
estimated_time: "0.5h"
dependencies: [1]
type: "doc"
mainSession: false
complexity: "low"
---

# 4: eval/validate-ux-pipeline: Add explicit surface guard for web sitemap reference

## Description

validate-ux-pipeline.md 的 PRD-to-Operation Translation 表格中，Web 行使用 sitemap.json 作为辅助信息来源，但无显式 surface 守卫阻止非 web 项目访问 sitemap。需为 Web 行的 sitemap 引用增加显式守卫。

## Reference Files
- `plugins/forge/skills/eval/rules/validate-ux-pipeline.md`: PRD-to-Operation Translation 表格 Web 行 Auxiliary Information 列引用 `sitemap.json`，需增加显式 web surface 守卫 (source: proposal.md#Scope-In-Scope)

## Affected Files

### Create
| File | Description |
|------|-------------|

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/eval/rules/validate-ux-pipeline.md` | PRD-to-Operation Translation 表格 Web 行增加 web surface 守卫说明 |

### Delete
| File | Reason |
|------|--------|

## Acceptance Criteria
- [ ] Web 行的 sitemap.json 引用增加显式守卫：仅在项目有 web surface 时使用 sitemap 作为辅助信息
- [ ] 非 web 类型（CLI、TUI）行的辅助信息不含 sitemap 相关内容，确认无遗漏

## Implementation Notes

表格按类型分行已有一定隔离效果，但需显式标注 sitemap 是 web surface 专属数据源，防止 agent 在非 web 场景误用。
