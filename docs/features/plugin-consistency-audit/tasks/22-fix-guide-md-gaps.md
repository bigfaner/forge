---
id: "22"
title: "Fix: guide.md surface coverage + CLI reference gaps"
priority: "P2"
estimated_time: "30min"
dependencies: []
type: "doc"
complexity: "medium"
mainSession: false
---

# 22: Fix: guide.md surface coverage + CLI reference gaps

## Description
hooks/guide.md 有 4 个问题：(1) Surface Type orchestration 遗漏 mobile（HOOK-01）；(2) Test Type 示例遗漏 tui/mobile（HOOK-02）；(3) Forge CLI 节缺少多个常用命令（HOOK-03）；(4) 缺少 forge config 机制说明（HOOK-04）。

## Reference Files
- `docs/features/plugin-consistency-audit/reports/05-commands-agent-hooks.md#HOOK-01`: guide.md mobile surface 遗漏 (source: Report 05)
- `docs/features/plugin-consistency-audit/reports/05-commands-agent-hooks.md#HOOK-03`: guide.md CLI 命令不完整 (source: Report 05)
- `docs/features/plugin-consistency-audit/reports/05-commands-agent-hooks.md#HOOK-04`: guide.md 缺少 config 说明 (source: Report 05)
- `plugins/forge/hooks/guide.md`: 需修改的文件 (source: audit finding)

## Affected Files

### Create
| File | Description |
|------|-------------|

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/hooks/guide.md` | (1) 添加 mobile surface；(2) 添加 tui/mobile test type；(3) 补充 CLI 命令；(4) 添加 config 子节 |

### Delete
| File | Reason |
|------|--------|

## Acceptance Criteria
- [ ] Surface Type orchestration 描述包含 mobile（有 probe + teardown，额外有 test-setup）
- [ ] Test Type 示例包含 `tui` → Terminal Functional Test 和 `mobile` → Mobile E2E Test
- [ ] Forge CLI 节补充常用命令：task claim/status/add、feature set/complete、config get、quality-gate、surfaces detect
- [ ] 添加 Configuration 子节说明 `forge config get auto.*` 配置键

## Hard Rules
- 仅修改 `plugins/forge/hooks/guide.md`
