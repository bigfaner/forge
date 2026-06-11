---
id: "11"
title: "Fix quick-tasks placeholder mapping (MEDIUM-C1)"
priority: "P1"
estimated_time: "15m"
dependencies: [10]
type: "doc"
mainSession: false
---

# 11: Fix quick-tasks placeholder mapping

## Description

quick-tasks/SKILL.md 的 Task Template Placeholders 表缺少 `{{COMPLEXITY}}` 和 `{{TYPE}}` 两行映射。模板 task.md 已使用这两个占位符（含注释块列出有效值），但 SKILL.md 的映射表未同步更新，LLM 可能忽略这两个字段。

## Reference Files
- `plugins/forge/skills/quick-tasks/SKILL.md`: Add COMPLEXITY and TYPE to placeholder mapping table (ref: Proposed Solution)
- `plugins/forge/skills/quick-tasks/templates/task.md`: Verify placeholders in template

## Affected Files

### Create
| File | Description |
|------|-------------|

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/quick-tasks/SKILL.md` | Add `{{COMPLEXITY}}` and `{{TYPE}}` rows to Task Template Placeholders table |

### Delete
| File | Reason |
|------|--------|

## Acceptance Criteria
- [ ] quick-tasks/SKILL.md Task Template Placeholders 表包含 `{{COMPLEXITY}}` 和 `{{TYPE}}` 映射行（含 Value Source 说明）

## Hard Rules
- Before modifying any plugin file, read `docs/conventions/forge-distribution.md`

## Implementation Notes
- 参考 breakdown-tasks/SKILL.md 中是否有类似的占位符映射表作为格式参考
