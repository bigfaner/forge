---
id: "13"
title: "Fix extract-design-md cross-references and templates (MEDIUM-A4, MINOR-C3)"
priority: "P2"
estimated_time: "15m"
dependencies: [10]
type: "doc"
mainSession: false
---

# 13: Fix extract-design-md cross-references and templates

## Description

extract-design-md SKILL.md 跨 skill 引用 `ui-design/templates/styles/<name>.md`，在 Forge 分发模型中 LLM 无法可靠解析另一个 skill 目录的文件。同时 extract-design-md 的 3 个 design 模板使用混合格式占位符（如 `{{App Name or Domain}}`）不符合 `{{UPPER_CASE}}` 标准。

## Reference Files
- `plugins/forge/skills/extract-design-md/SKILL.md`: Fix cross-skill reference to ui-design styles
- `plugins/forge/skills/extract-design-md/templates/design-web.md`: Standardize placeholder format
- `plugins/forge/skills/extract-design-md/templates/design-tui.md`: Standardize placeholder format
- `plugins/forge/skills/extract-design-md/templates/design-mobile.md`: Standardize placeholder format

## Affected Files

### Create
| File | Description |
|------|-------------|

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/extract-design-md/SKILL.md` | Clarify cross-skill reference to ui-design templates |
| `plugins/forge/skills/extract-design-md/templates/design-web.md` | Standardize placeholder format to UPPER_CASE |
| `plugins/forge/skills/extract-design-md/templates/design-tui.md` | Standardize placeholder format to UPPER_CASE |
| `plugins/forge/skills/extract-design-md/templates/design-mobile.md` | Standardize placeholder format to UPPER_CASE |

### Delete
| File | Reason |
|------|--------|

## Acceptance Criteria
- [ ] extract-design-md SKILL.md 中跨 skill 引用添加说明（如 "This file is in ui-design skill directory; in Forge distribution, use the full plugin path"）
- [ ] 3 个 design 模板中的占位符统一为 `{{UPPER_CASE}}` 格式

## Hard Rules
- Before modifying any plugin file, read `docs/conventions/forge-distribution.md`

## Implementation Notes
- 模板中 `{{App Name or Domain}}` → `{{APP_NAME_OR_DOMAIN}}`，`{{YYYY-MM-DD}}` → `{{DATE}}`，`{{value}}` → `{{VALUE}}` 等
- 跨 skill 引用问题可能需要在 SKILL.md 中添加路径解析说明
