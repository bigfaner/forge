---
status: "completed"
started: "2026-06-01 21:30"
completed: "2026-06-01 21:32"
time_spent: "~2m"
---

# Task Record: 4 eval/validate-ux-pipeline: Add explicit surface guard for web sitemap reference

## Summary
Added explicit web surface guard to validate-ux-pipeline.md PRD-to-Operation Translation table: sitemap.json now annotated as web-surface-only data source with instruction to skip sitemap lookup for non-web projects

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/eval/rules/validate-ux-pipeline.md

### Key Decisions
无

## Document Metrics
1 table cell updated, 0 naming inconsistencies, 2 AC verified

## Referenced Documents
- docs/features/sitemap-surface-guard/tasks/4-eval-surface-guard.md

## Review Status
final

## Acceptance Criteria
- [x] Web 行的 sitemap.json 引用增加显式守卫：仅在项目有 web surface 时使用 sitemap 作为辅助信息
- [x] 非 web 类型（CLI、TUI）行的辅助信息不含 sitemap 相关内容，确认无遗漏

## Notes
Guard added inline in Auxiliary Information cell: '**web surface only**; skip sitemap lookup entirely for non-web projects (CLI, TUI)'
