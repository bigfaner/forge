---
id: "14"
title: "Fix write-prd template consistency (MINOR-C4, MINOR-C5)"
priority: "P2"
estimated_time: "15m"
dependencies: [10]
type: "doc"
mainSession: false
---

# 14: Fix write-prd template consistency

## Description

write-prd SKILL.md 缺少占位符映射表，多个模板占位符（如 `{{FEATURE_NAME}}`、`{{DB_SCHEMA}}`）没有显式映射。prd-ui-functions.md 使用非标准格式占位符如 `{{web | mobile | mini-program | tablet | tui}}`。

## Reference Files
- `plugins/forge/skills/write-prd/SKILL.md`: Add or complete placeholder mapping
- `plugins/forge/skills/write-prd/templates/prd-ui-functions.md`: Standardize placeholder format

## Affected Files

### Create
| File | Description |
|------|-------------|

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/write-prd/SKILL.md` | Add placeholder mapping annotations where missing |
| `plugins/forge/skills/write-prd/templates/prd-ui-functions.md` | Standardize `{{web | mobile | ...}}` to `{{PLATFORM}}` with comment listing valid values |

### Delete
| File | Reason |
|------|--------|

## Acceptance Criteria
- [ ] write-prd/SKILL.md 中关键占位符（`{{FEATURE_NAME}}`、`{{DB_SCHEMA}}`、`{{PRD_SUMMARY}}` 等）有赋值逻辑说明或映射
- [ ] prd-ui-functions.md 中 `{{web | mobile | ...}}` 格式替换为标准 `{{PLATFORM}}` 占位符并附注释

## Hard Rules
- Before modifying any plugin file, read `docs/conventions/forge-distribution.md`

## Implementation Notes
- 混合格式占位符 `{{web | mobile | mini-program | tablet | tui}}` 是枚举选择型，应改为 `{{PLATFORM}}` 并在模板注释中列出可选值
