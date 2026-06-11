---
id: "2"
title: "Surface rules 增加 compile/fmt/lint/unit-test gate recipe 模板"
priority: "P1"
estimated_time: "1h"
dependencies: []
type: "doc"
mainSession: false
---

# 2: Surface rules 增加 compile/fmt/lint/unit-test gate recipe 模板

## Description
当前 surface rule 文件（api.md、web.md、cli.md、tui.md、mobile.md）只生成 lifecycle recipe（dev/probe/test/teardown），不生成 compile/fmt/lint/unit-test 的 prefixed recipe。需要为每个 surface rule 增加 `<key>-compile`、`<key>-fmt`、`<key>-lint`、`<key>-unit-test` 的 stub recipe 定义和 Recipe Invocation Contract 条目，使非 mixed 多 surface 项目也能生成 prefixed recipe。

## Reference Files
- `docs/proposals/per-task-surface-scoped-gate/proposal.md` — Proposed Solution, Constraints & Dependencies, Success Criteria
- `plugins/forge/skills/init-justfile/rules/surfaces/api.md`: 增加 gate recipe stub 和 Recipe Invocation Contract (ref: Proposed Solution)
- `plugins/forge/skills/init-justfile/rules/surfaces/web.md`: 增加 gate recipe stub 和 Recipe Invocation Contract (ref: Proposed Solution)
- `plugins/forge/skills/init-justfile/rules/surfaces/cli.md`: 增加 gate recipe stub 和 Recipe Invocation Contract (ref: Proposed Solution)
- `plugins/forge/skills/init-justfile/rules/surfaces/tui.md`: 增加 gate recipe stub 和 Recipe Invocation Contract (ref: Proposed Solution)

## Affected Files

### Create
| File | Description |
|------|-------------|

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/init-justfile/rules/surfaces/api.md` | 增加 compile/fmt/lint/unit-test gate recipe stub 定义和 Recipe Invocation Contract 条目 |
| `plugins/forge/skills/init-justfile/rules/surfaces/web.md` | 增加 compile/fmt/lint/unit-test gate recipe stub 定义和 Recipe Invocation Contract 条目 |
| `plugins/forge/skills/init-justfile/rules/surfaces/cli.md` | 增加 compile/fmt/lint/unit-test gate recipe stub 定义和 Recipe Invocation Contract 条目 |
| `plugins/forge/skills/init-justfile/rules/surfaces/tui.md` | 增加 compile/fmt/lint/unit-test gate recipe stub 定义和 Recipe Invocation Contract 条目 |
| `plugins/forge/skills/init-justfile/rules/surfaces/mobile.md` | 增加 compile/fmt/lint/unit-test gate recipe stub 定义和 Recipe Invocation Contract 条目 |

### Delete
| File | Reason |
|------|--------|

## Acceptance Criteria
- [ ] 5 个 surface rule 文件（api.md、web.md、cli.md、tui.md、mobile.md）均包含 `<key>-compile`/`<key>-fmt`/`<key>-lint`/`<key>-unit-test` 的 stub recipe 定义
- [ ] 5 个 surface rule 文件均包含 compile/fmt/lint/unit-test 的 Recipe Invocation Contract 条目

## Hard Rules
- 修改 `plugins/forge/` 下的文件前必须先加载 `docs/conventions/forge-distribution.md`

## Implementation Notes
- Stub recipe 使用 surface **key**（非 type）作为前缀，与 SKILL.md line 259-263 一致
- Stub recipe 模板须包含约束注释（如 `# This recipe compiles ONLY the <key> surface code`），确保 LLM 理解 recipe 的 surface 隔离语义
- `mobile.md` 遵循与 `api.md`/`web.md` 相同的模板模式
- 单语言模板（go.just、node.just 等）不需要修改 — 单 surface 项目无需 prefixed recipe
- `mixed.just` 模板已硬编码 `backend-compile`/`frontend-lint` 等 recipe，不需要修改
