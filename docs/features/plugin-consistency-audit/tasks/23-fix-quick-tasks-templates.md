---
id: "23"
title: "Fix: quick-tasks template hardcoded defaults"
priority: "P2"
estimated_time: "15min"
dependencies: []
type: "doc"
complexity: "low"
mainSession: false
---

# 23: Fix: quick-tasks template hardcoded defaults

## Description
`plugins/forge/skills/quick-tasks/templates/task.md` 的 frontmatter 硬编码 `complexity: "medium"` 和 `type: "coding.feature"`，但 SKILL.md Step 2 定义了多值启发式规则（low/medium/high + 8 种 type）。模板的硬编码值覆盖了 SKILL.md 的逻辑，除非 agent 手动编辑生成的内容。

## Reference Files
- `docs/features/plugin-consistency-audit/reports/03-skills-batch-b.md#QT-02`: P3 CONFLICT, complexity 硬编码 (source: Report 03)
- `docs/features/plugin-consistency-audit/reports/03-skills-batch-b.md#QT-03`: P3 CONFLICT, type 硬编码 (source: Report 03)
- `plugins/forge/skills/quick-tasks/templates/task.md`: 需修改的模板文件 (source: audit finding)

## Affected Files

### Create
| File | Description |
|------|-------------|

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/quick-tasks/templates/task.md` | 替换 `complexity: "medium"` → `complexity: "{{COMPLEXITY}}"`，`type: "coding.feature"` → `type: "{{TYPE}}"` |

### Delete
| File | Reason |
|------|--------|

## Acceptance Criteria
- [ ] `complexity` 字段使用 `{{COMPLEXITY}}` 占位符
- [ ] `type` 字段使用 `{{TYPE}}` 占位符
- [ ] 模板注释说明默认值（complexity 默认 medium，type 默认 coding.feature）

## Hard Rules
- 仅修改 `plugins/forge/skills/quick-tasks/templates/task.md`
