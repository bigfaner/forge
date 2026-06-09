---
id: "3"
title: "Simplify 5 surface rule files to replace TODO stubs with Recipe Generation Requirements"
priority: "P1"
estimated_time: "1.5h"
dependencies: [1]
type: "doc"
mainSession: false
# Note: surface-key and surface-type fields are intentionally absent from doc tasks.
# Doc tasks produce non-compilable output (markdown, specs, templates) and do not
# interact with the quality gate or test pipeline, so surface routing is unnecessary.
---

# 3: Simplify 5 surface rule files to replace TODO stubs with Recipe Generation Requirements

## Description

Simplify all 5 surface rule files (api/cli/tui/web/mobile) by replacing the `## Recipe Template (Dual Platform)` section's TODO stub code blocks with a new `## Recipe Generation Requirements` section. The TODO stubs (`echo "TODO: implement ..."`) were designed for template-driven generation where the LLM filled in stubs. With agent-driven generation, the agent reads recipe contracts (names, signatures, exit codes) from the `## Recipe Invocation Contract` section and generates commands directly — stubs are unnecessary noise.

The simplification preserves all contract-level information (orchestration sequence, invocation contract, journey filter strategy) and replaces only the template section. The new section describes structural constraints the agent must follow when generating recipes.

## Reference Files
- `docs/proposals/agent-driven-justfile-generation/proposal.md` — Proposed Solution, Scope > In Scope (ref: ## Proposed Solution, ## Scope > ### In Scope)
- `plugins/forge/skills/init-justfile/rules/surfaces/api.md` — current structure with TODO stubs (ref: ## Recipe Template (Dual Platform))
- `plugins/forge/skills/init-justfile/rules/surfaces/cli.md` — current structure with TODO stubs (ref: ## Recipe Template (Dual Platform))

## Affected Files

### Create
| File | Description |
|------|-------------|

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/init-justfile/rules/surfaces/api.md` | Replace Recipe Template section with Recipe Generation Requirements |
| `plugins/forge/skills/init-justfile/rules/surfaces/cli.md` | Replace Recipe Template section with Recipe Generation Requirements |
| `plugins/forge/skills/init-justfile/rules/surfaces/tui.md` | Replace Recipe Template section with Recipe Generation Requirements |
| `plugins/forge/skills/init-justfile/rules/surfaces/web.md` | Replace Recipe Template section with Recipe Generation Requirements |
| `plugins/forge/skills/init-justfile/rules/surfaces/mobile.md` | Replace Recipe Template section with Recipe Generation Requirements |

### Delete
| File | Reason |
|------|--------|

## Acceptance Criteria
- [ ] 所有 5 个 surface rule 文件（`api.md`/`cli.md`/`tui.md`/`web.md`/`mobile.md`）已更新
- [ ] `## Recipe Template (Dual Platform)` section（含 TODO stub 代码块和 `<Test-Dir-Path>` 块）已替换为 `## Recipe Generation Requirements` section
- [ ] 以下 section 保留不变：Orchestration Sequence、Recipe Invocation Contract、Journey Filter Strategy
- [ ] `## Recipe Generation Requirements` section 包含：agent 生成 recipe 时需遵循的结构约束（recipe 命名规则、`[linux]`/`[windows]` 双平台属性、`# user-customized` 标记、exit code 0/1 语义、test 目录路径规则 single vs multi surface）
- [ ] 双消费者一致性保留：init-justfile（生成 recipe）和 run-tests（消费 recipe）都能从同一份 surface rule 的 Recipe Invocation Contract 获取所需信息

## Implementation Notes

### Recipe Generation Requirements section 内容要素
1. **命名规则**：named surface → `<key>-<verb>`，scalar surface → `<verb>`（无前缀）
2. **双平台**：每个 recipe 必须有 `[linux]` 和 `[windows]` 变体
3. **user-customized 标记**：每个 recipe 的 `[linux]` 变体上方必须有 `# user-customized` 注释
4. **exit code 语义**：必须与 Recipe Invocation Contract 表中定义一致
5. **test 目录路径**：single surface → `tests/<journey>/`，multi surface → `tests/<surfaceKey>/<journey>/`
6. **CLI/TUI 特殊约束**：不生成 dev、probe、aggregate recipe
7. **Aggregate recipe**（web/api/mobile）：`dev->probe->test; rc=$?; teardown; exit $rc` 模式
8. **Gate recipes**：compile/fmt/lint/unit-test 仅 multi-surface 项目生成，scope 限定为对应 surface

### 对 server-lifecycle 的引用
- surface rule 中的 dev/probe/teardown recipe 生成应引用 `rules/server-lifecycle.md` 中的 bash 模式
- 不在 surface rule 中重复 server lifecycle 代码，通过引用保持单一来源

### 文件修改范围
仅修改 `plugins/forge/skills/init-justfile/rules/surfaces/` 下的 5 个文件，每个文件的改动局限于替换 Recipe Template section。
