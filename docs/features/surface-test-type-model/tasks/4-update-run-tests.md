---
id: "4"
title: "更新 run-tests skill 文件术语和输出格式"
priority: "P1"
estimated_time: "1.5h"
dependencies: ["1"]
type: "doc"
mainSession: false
---

# 4: 更新 run-tests skill 文件术语和输出格式

## Description

更新 run-tests skill 的 SKILL.md 及 5 个 surface orchestration 规则文件中的测试类型术语。run-tests 负责测试执行编排，需要确保输出标签（suite 名称、进度提示）使用 surface-specific 的测试类型名称。

## Reference Files
- `docs/proposals/surface-test-type-model/proposal.md#User-Facing-Experience` — 定义了测试执行时的输出标签变更（如 "Running CLI functional tests..."）
- `docs/proposals/surface-test-type-model/proposal.md#Test-Type-Mapping` — 每种 surface 的测试类型名称和执行模型
- `docs/proposals/surface-test-type-model/proposal.md#Success-Criteria` — SC9（测试执行输出使用 surface-specific 名称）

## Affected Files

### Create
| File | Description |
|------|-------------|

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/run-tests/SKILL.md` | 更新 "e2e" 引用（promote 描述和提示），引用概念文档 |
| `plugins/forge/skills/run-tests/rules/surfaces/cli.md` | 使用 "CLI 功能测试" 术语，更新 suite 名称 |
| `plugins/forge/skills/run-tests/rules/surfaces/tui.md` | 使用 "终端功能测试" 术语，更新 suite 名称 |
| `plugins/forge/skills/run-tests/rules/surfaces/api.md` | 使用 "API 功能测试" 术语，更新 suite 名称 |
| `plugins/forge/skills/run-tests/rules/surfaces/web.md` | 使用 "Web 端到端测试" 术语，更新 suite 名称，Journey filter 标签更新 |
| `plugins/forge/skills/run-tests/rules/surfaces/mobile.md` | 使用 "移动端端到端测试" 术语，更新 suite 名称 |

### Delete
| File | Reason |
|------|--------|

## Acceptance Criteria

- [ ] run-tests/SKILL.md 中 "e2e" 仅出现在 promote 功能描述中，且标注为 "Web/Mobile 端到端测试"
- [ ] rules/surfaces/ 5 个文件各自使用对应 surface 的测试类型名称
- [ ] 测试执行输出中的 suite 名称使用 surface-specific 测试类型（如 `cli-functional/journey-name` 而非 `e2e/journey-name`）
- [ ] Web surface 的 Journey filter 标签从 `@e2e` 更新为 `@web-e2e`（或等效的 surface-specific 标签）
- [ ] 所有 rules 文件引用概念文档 `docs/reference/test-type-model.md`

## Hard Rules

- orchestration 序列不变（dev → probe → test → teardown 等），仅更新术语
- CLI/TUI surface 保持简化序列（test → teardown）

## Implementation Notes

SKILL.md 第 265 和 267 行有 `forge test promote` 相关的 "e2e" 引用。promote 功能将 passing tests 推到 regression suite，这个 "regression suite" 的命名也需要与新的测试类型模型一致。rules/surfaces/web.md 的 Journey filter 当前包含 `@e2e`，需要更新为 `@web-e2e` 或类似标签。
