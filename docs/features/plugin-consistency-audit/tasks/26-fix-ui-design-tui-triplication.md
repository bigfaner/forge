---
id: "26"
title: "Fix: ui-design TUI requirements triplication"
priority: "P2"
estimated_time: "30min"
dependencies: []
type: "doc"
complexity: "medium"
mainSession: false
---

# 26: Fix: ui-design TUI requirements triplication

## Description
ui-design 的 5 个 TUI panel 强制结构要求定义在 3 个文件中：(1) `rules/tui-panel-requirements.md`（权威定义）；(2) `templates/platforms/tui.md`（逐字重复）；(3) `templates/ui-design.md`（内联重复）。任何变更需同步 3 处，维护风险高。应将 tui.md 和 ui-design.md 中的重复内容替换为对 rules 文件的引用。

## Reference Files
- `docs/features/plugin-consistency-audit/reports/03-skills-batch-b.md#UD-01`: P2 REDUNDANT, TUI requirements 三重定义 (source: Report 03)
- `plugins/forge/skills/ui-design/rules/tui-panel-requirements.md`: 权威定义（保留） (source: audit finding)
- `plugins/forge/skills/ui-design/templates/platforms/tui.md`: 需替换重复内容为引用 (source: audit finding)
- `plugins/forge/skills/ui-design/templates/ui-design.md`: 需替换内联 TUI section 为引用 (source: audit finding)

## Affected Files

### Create
| File | Description |
|------|-------------|

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/ui-design/templates/platforms/tui.md` | 将内联的 5 个 structural requirements 替换为 "Per `rules/tui-panel-requirements.md`" 引用 |
| `plugins/forge/skills/ui-design/templates/ui-design.md` | 将 TUI Component 模板节中的 5 个 requirements 替换为引用 |

### Delete
| File | Reason |
|------|--------|

## Acceptance Criteria
- [ ] `templates/platforms/tui.md` 不再逐字重复 5 个 structural requirements，改为引用 `rules/tui-panel-requirements.md`
- [ ] `templates/ui-design.md` 的 TUI Component 节同样改为引用
- [ ] `rules/tui-panel-requirements.md` 保持不变（权威来源）

## Hard Rules
- 仅修改 `templates/platforms/tui.md` 和 `templates/ui-design.md`
- 不修改 `rules/tui-panel-requirements.md`
