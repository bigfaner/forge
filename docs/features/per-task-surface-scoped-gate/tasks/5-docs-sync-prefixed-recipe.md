---
id: "5"
title: "文档同步：quality gate prefixed recipe 规范更新"
priority: "P2"
estimated_time: "30min"
dependencies: ["1"]
type: "doc"
mainSession: false
---

# 5: 文档同步：quality gate prefixed recipe 规范更新

## Description
Task 1 将 per-task gate 的 scope 机制从参数模式切换为 prefixed recipe 模式，相关文档仍描述旧的参数模式。需要同步更新 init-justfile SKILL.md、business-rules/quality-gate.md 和 conventions/dispatcher-quality.md，使文档反映新的 prefixed recipe 语义。

## Reference Files
- `docs/proposals/per-task-surface-scoped-gate/proposal.md` — Proposed Solution, Constraints & Dependencies
- `plugins/forge/skills/init-justfile/SKILL.md`: Standard Target Contract 需增加 prefixed recipe 说明 (ref: Proposed Solution)
- `docs/business-rules/quality-gate.md`: 描述 `just unit-test [scope]` 参数模式，需更新 (ref: Proposed Solution)
- `docs/conventions/dispatcher-quality.md`: 描述 `just compile` after each task，需增加 multi-surface 说明 (ref: Proposed Solution)

## Affected Files

### Create
| File | Description |
|------|-------------|

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/init-justfile/SKILL.md` | Standard Target Contract 段落增加 multi-surface prefixed recipe 模式说明 |
| `docs/business-rules/quality-gate.md` | 更新 gate 描述从参数模式到 prefixed recipe 语义 |
| `docs/conventions/dispatcher-quality.md` | 增加 multi-surface 场景下 prefixed recipe 使用说明 |

### Delete
| File | Reason |
|------|--------|

## Acceptance Criteria
- [ ] `init-justfile/SKILL.md` Standard Target Contract 包含 multi-surface 项目的 `<key>-compile`/`<key>-fmt`/`<key>-lint`/`<key>-unit-test` prefixed recipe 模式说明
- [ ] `docs/business-rules/quality-gate.md` 中 `just unit-test [scope]` 的参数模式描述更新为 prefixed recipe 语义（`just <key>-unit-test`）
- [ ] `docs/conventions/dispatcher-quality.md` 说明 multi-surface 场景下 dispatcher 使用 surface-key 对应的 prefixed recipe

## Hard Rules
- 修改 `plugins/forge/` 下的文件前必须先加载 `docs/conventions/forge-distribution.md`

## Implementation Notes
- `init-justfile/SKILL.md` line 45-50 定义 Standard Target Contract，需在现有 `compile`/`unit-test`/`lint`/`fmt` 条目后增加 multi-surface prefixed recipe 的命名规则和生成条件
- `docs/business-rules/quality-gate.md` line 17 描述了 `ResolveScope` 探测机制，需替换为 `resolvePrefixedRecipe()` 的前缀匹配语义
- `docs/conventions/dispatcher-quality.md` line 12/36 描述 dispatcher 每个任务后运行 `just compile`，需增加 `surface-key` 存在时使用 prefixed recipe 的说明
- 保持 feature-level gate（全量验证）的描述不变，仅更新 per-task gate 的描述
