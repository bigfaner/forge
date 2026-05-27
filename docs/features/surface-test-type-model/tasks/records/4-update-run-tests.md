---
status: "completed"
started: "2026-05-26 22:34"
completed: "2026-05-26 22:37"
time_spent: "~3m"
---

# Task Record: 4 更新 run-tests skill 文件术语和输出格式

## Summary
更新 run-tests skill 的 SKILL.md 及 5 个 surface orchestration 规则文件中的测试类型术语，使用 surface-specific 测试类型名称替换通用 e2e 标签。

## Changes

### Files Created
无

### Files Modified
- `plugins/forge/skills/run-tests/SKILL.md` — 更新 promote 描述区分 E2E/functional surfaces，新增"测试类型术语"段落引用概念文档
- `plugins/forge/skills/run-tests/rules/surfaces/cli.md` — 标题和编排表使用"CLI 功能测试"术语，新增 Suite 名称段落
- `plugins/forge/skills/run-tests/rules/surfaces/tui.md` — 标题和编排表使用"终端功能测试"术语，新增 Suite 名称段落
- `plugins/forge/skills/run-tests/rules/surfaces/api.md` — 标题和编排表使用"API 功能测试"术语，新增 Suite 名称段落
- `plugins/forge/skills/run-tests/rules/surfaces/web.md` — 标题和编排表使用"Web 端到端测试"术语，Journey filter 标签 @e2e -> @web-e2e，新增 Suite 名称段落
- `plugins/forge/skills/run-tests/rules/surfaces/mobile.md` — 标题和编排表使用"移动端端到端测试"术语，新增 Suite 名称段落

### Key Decisions
- promote 描述保留 "e2e" 一词但限定在 Web/Mobile 端到端测试上下文中（符合 test-type-model.md 约束）
- Suite 名称格式: `<surface>-<scope>/<journey-name>`（如 cli-functional/journey-name, web-e2e/journey-name）

## Document Metrics
6 files modified: 1 SKILL.md + 5 surface rules; 每个规则文件新增 Suite 名称段落和概念文档引用

## Referenced Documents
- `docs/proposals/surface-test-type-model/proposal.md`
- `docs/reference/test-type-model.md`

## Review Status
final

## Acceptance Criteria
- [x] SKILL.md 中 e2e 仅出现在 promote 功能描述中，且标注为 Web/Mobile 端到端测试 — PASS
- [x] rules/surfaces/ 5 个文件各自使用对应 surface 的测试类型名称 — PASS
- [x] 测试执行输出中的 suite 名称使用 surface-specific 测试类型 — PASS
- [x] Web surface 的 Journey filter 标签从 @e2e 更新为 @web-e2e — PASS
- [x] 所有 rules 文件引用概念文档 docs/reference/test-type-model.md — PASS

## Notes
编排序列不变（仅更新术语），CLI/TUI 保持简化序列（test -> teardown）。run-tests rules/ 下的其他文件（failure-diagnosis.md、result-parsing.md、test-isolation.md）中的 "e2e" 引用属于 Forge 自身内部测试基础设施规则，不在本任务范围内。
