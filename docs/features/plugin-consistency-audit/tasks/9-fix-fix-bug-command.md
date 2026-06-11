---
id: "9"
title: "Fix: fix-bug command - add AskUserQuestion, replace Playwright, add mobile/tui"
priority: "P1"
estimated_time: "30min"
dependencies: []
type: "doc"
complexity: "medium"
mainSession: false
---

# 9: Fix: fix-bug command - add AskUserQuestion, replace Playwright, add mobile/tui

## Description
fix-bug command 有 3 个相关问题：(1) `allowed-tools` 缺少 `AskUserQuestion`，但 Step 5 使用该工具 (CMD-01)；(2) Bug surface 表格硬编码 `Playwright` (CMD-02)；(3) 表格缺少 `mobile` 和 `tui` 行 (CMD-14)。

## Reference Files
- `docs/features/plugin-consistency-audit/reports/05-commands-agent-hooks.md#CMD-01`: P1 级 CONFLICT，allowed-tools 缺少 AskUserQuestion (source: Report 05)
- `docs/features/plugin-consistency-audit/reports/05-commands-agent-hooks.md#CMD-02`: P1 级 CONFLICT，Playwright 硬编码 (source: Report 05)
- `docs/features/plugin-consistency-audit/reports/05-commands-agent-hooks.md#CMD-14`: P2 级 INCOMPLETE，缺少 mobile/tui 行 (source: Report 05)
- `plugins/forge/commands/fix-bug.md`: 需修复的文件 (source: audit finding)

## Affected Files

### Create
| File | Description |
|------|-------------|

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/commands/fix-bug.md` | (1) 添加 AskUserQuestion 到 allowed-tools；(2) 替换 Playwright 硬编码；(3) 添加 mobile/tui 行 |

### Delete
| File | Reason |
|------|--------|

## Acceptance Criteria
- [ ] `allowed-tools` 包含 `AskUserQuestion`
- [ ] Bug surface 表格中 `Playwright` 替换为通用描述（如 "test profile runner" 或 Convention-derived 命令）
- [ ] Bug surface 表格包含 `mobile` 行（Maestro YAML / mobile test profile）
- [ ] Bug surface 表格包含 `tui` 行（子进程 + stdin pipe / tui test profile）

## Hard Rules
- 仅修改 `plugins/forge/commands/fix-bug.md`

## Implementation Notes
- 参考 `docs/reference/test-type-model.md` 中 mobile 和 tui 的 test type 定义
- 参考 run-tests 的 surface rules 中 mobile/tui 的 test runner 描述
