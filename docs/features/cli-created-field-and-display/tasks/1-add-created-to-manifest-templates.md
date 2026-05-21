---
id: "1"
title: "Manifest 模板添加 created frontmatter"
priority: "P2"
estimated_time: "30m"
dependencies: []
type: "doc"
mainSession: false
---

# 1: Manifest 模板添加 created frontmatter

## Description

为两个 manifest 模板添加 `created` frontmatter 字段，使新生成的 feature manifest 包含创建日期。

当前状态：标准 manifest 模板和 quick manifest 模板均无 `created` 字段。Proposal 和 Lesson 已有该字段。

## Reference Files
- `docs/proposals/cli-created-field-and-display/proposal.md` — Source proposal

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/write-prd/templates/manifest.md` | frontmatter 添加 `created` 字段 |
| `plugins/forge/skills/quick-tasks/templates/manifest-quick.md` | frontmatter 添加 `created` 字段 |

## Acceptance Criteria
- [ ] `plugins/forge/skills/write-prd/templates/manifest.md` frontmatter 包含 `created: "{{DATE}}"` (YYYY-MM-DD 格式占位符)
- [ ] `plugins/forge/skills/quick-tasks/templates/manifest-quick.md` frontmatter 包含 `created: "{{DATE}}"` (YYYY-MM-DD 格式占位符)
- [ ] 使用该模板的 skill（write-prd、quick-tasks）在生成 manifest 时填充实际日期

## Hard Rules
- 日期格式必须为 `YYYY-MM-DD`，与 proposal/lesson 的 `created` 字段保持一致

## Implementation Notes
- 检查使用这两个模板的 skill 代码，确保生成 manifest 时将 `{{DATE}}` 替换为当天日期
- 不需要回填已有 feature manifest 的 `created` 字段（Out of Scope）
