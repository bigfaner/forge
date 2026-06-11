---
id: "3"
title: "Fix breakdown-tasks template placeholders (H-3)"
priority: "P1"
estimated_time: "30m"
dependencies: [1]
type: "doc"
mainSession: false
---

# 3: Fix breakdown-tasks template placeholders (H-3)

## Description

breakdown-tasks/templates/task.md 硬编码 `complexity: "medium"` 和 `type: "coding.feature"`，而非使用占位符。SKILL.md 有完整的 complexity 和 type 判定规则，但模板不使用占位符，导致 LLM 在疲劳上下文下可能直接采用模板值。

## Reference Files
- `docs/proposals/forge-skill-audit/proposal.md` — H-3: breakdown-tasks task.md 模板硬编码 complexity 和 type, Proposed Solution, Success Criteria
- `plugins/forge/skills/breakdown-tasks/templates/task.md`: Replace hardcoded values with placeholders (ref: H-3: breakdown-tasks task.md 模板硬编码 complexity 和 type)
- `plugins/forge/skills/quick-tasks/templates/task.md`: Reference for correct placeholder usage (ref: H-3: breakdown-tasks task.md 模板硬编码 complexity 和 type)

## Affected Files

### Create
| File | Description |
|------|-------------|

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/breakdown-tasks/templates/task.md` | Replace `complexity: "medium"` → `complexity: "{{COMPLEXITY}}"` and `type: "coding.feature"` → `type: "{{TYPE}}"`; add comment block |

### Delete
| File | Reason |
|------|--------|

## Acceptance Criteria
- [ ] breakdown-tasks/templates/task.md 使用 `complexity: "{{COMPLEXITY}}"` 和 `type: "{{TYPE}}"` 占位符
- [ ] 模板包含与 quick-tasks/templates/task.md 一致的注释块，列出 COMPLEXITY 和 TYPE 的可选值及默认值

## Hard Rules
- Before modifying any plugin file, read `docs/conventions/forge-distribution.md`
- Only modify markdown files; no Go code changes

## Implementation Notes
- 修复后运行回归验证：`grep -r 'complexity: "medium"' plugins/forge/skills/breakdown-tasks/templates/` 确认无残留
- 注释块格式：`# Template placeholders: #   COMPLEXITY — low | medium | high (default: medium) #   TYPE — coding.feature | coding.enhancement | ... (default: coding.feature)`
