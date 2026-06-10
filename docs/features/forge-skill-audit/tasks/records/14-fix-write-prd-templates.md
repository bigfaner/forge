---
status: "completed"
started: "2026-06-10 21:12"
completed: "2026-06-10 21:13"
time_spent: "~1m"
---

# Task Record: 14 Fix write-prd template consistency (MINOR-C4, MINOR-C5)

## Summary
Added Template Placeholder Mapping table to write-prd SKILL.md covering all 8 placeholders (FEATURE_NAME, SLUG, DATE, DB_SCHEMA, PRD_SUMMARY, USER_STORIES_SUMMARY, UI_FUNCTIONS_SUMMARY, PLATFORM) with value sources and examples. Standardized prd-ui-functions.md enum-style placeholder from {{web | mobile | mini-program | tablet | tui}} to {{PLATFORM}} with HTML comment listing valid values.

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/write-prd/SKILL.md
- plugins/forge/skills/write-prd/templates/prd-ui-functions.md

### Key Decisions
无

## Document Metrics
8 placeholders mapped in SKILL.md, 1 placeholder standardized in prd-ui-functions.md

## Referenced Documents
- docs/conventions/forge-distribution.md

## Review Status
final

## Acceptance Criteria
- [x] SKILL.md 中关键占位符（{{FEATURE_NAME}}、{{DB_SCHEMA}}、{{PRD_SUMMARY}} 等）有赋值逻辑说明或映射
- [x] prd-ui-functions.md 中 {{web | mobile | ...}} 格式替换为标准 {{PLATFORM}} 占位符并附注释

## Notes
Template Placeholder Mapping section added between Output Documents and Step 1 in SKILL.md. Enum-style placeholder replaced with {{PLATFORM}} and valid values listed in HTML comment.
