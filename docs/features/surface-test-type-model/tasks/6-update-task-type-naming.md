---
id: "6"
title: "更新 forge-cli task type 命名携带 surface 信息"
priority: "P0"
estimated_time: "2h"
dependencies: ["1"]
surface-key: ""
surface-type: ""
breaking: true
type: "coding.enhancement"
mainSession: false
---

# 6: 更新 forge-cli task type 命名携带 surface 信息

## Description

修改 forge-cli 的 task type 命名，从当前的统一类型名（`test.gen-scripts`、`test.run`）改为携带 surface 信息的三段式命名（`test.gen-scripts.cli`、`test.run.backend`）。同时更新任务模板中的 "e2e" 术语为 surface-specific 测试类型名称，并更新 Go 测试中的相关断言。

## Reference Files
- `docs/proposals/surface-test-type-model/proposal.md#Technical-Direction` — task type 三段式命名规则（{action}.{skill}.{surface}），parser 已支持
- `docs/proposals/surface-test-type-model/proposal.md#Test-Type-Mapping` — 每种 surface 对应的测试类型名称，用于 task title 生成
- `docs/proposals/surface-test-type-model/proposal.md#Success-Criteria` — SC6（index.json 中 task type 携带 surface 信息）和 SC8（向后兼容）
- `docs/reference/test-type-model.md` — 权威测试类型定义（由 Task 1 创建）

## Acceptance Criteria

- [ ] gen-scripts 任务的 type 字段从 `test.gen-scripts` 变为 `test.gen-scripts.<surfaceType>`（如 `test.gen-scripts.cli`）
- [ ] run-test 任务的 type 字段从 `test.run` 变为 `test.run.<surfaceKey>`（单 surface 项目使用 surface type）
- [ ] verify-regression 任务的 type 字段从 `test.verify-regression` 变为 `test.verify-regression.<surfaceKey>`
- [ ] 任务标题从 "Run e2e Tests" 变为 surface-specific 名称（如 "Run CLI Functional Tests"）
- [ ] 任务模板（test-run.md、test-verify-regression.md）使用 surface-specific 测试类型描述
- [ ] Go 测试中的相关断言（如 autogen_test.go 中的 StrategyKind "just test-e2e"）更新为新命名
- [ ] `go test ./pkg/task/...` 通过

## Hard Rules

- **不要修改 task-lifecycle parser**。proposal 明确说明 parser 已支持带点的类型名。只需修改 autogen.go 中的 type 值生成逻辑
- 保持向后兼容：生成的 fix task ID 模式（如 `fix-e2e-1-1`）暂不改动，作为后续迭代
- TaskIndex.E2ERound 字段暂不重命名，本 task 仅修改 type 值和 display strings

## Implementation Notes

关键改动点：
1. **types.go**: Task type 常量可改为函数式生成（如 `GenTestType(action, surface string) string`），或新增 surface-specific 常量
2. **autogen.go**: GetBreakdownTestTasks/GetQuickTestTasks 中 per-surface 任务的 Type 字段使用 surface-specific 值
3. **autogen.go**: Task title 格式从 `fmt.Sprintf("Run e2e Tests (%s)", key)` 改为 `fmt.Sprintf("Run %s", TestTypeTitle(surfaceType))`
4. **data/test-run.md**: 标题从 "Execute staged e2e test scripts" 改为使用 surface-specific 测试类型名
5. **data/test-verify-regression.md**: "Run full e2e regression suite" 和 `just test-e2e` 引用更新
6. **autogen_test.go**: 更新 type 和 title 断言值

TestTypeTitle 辅助函数建议放在 types.go 中，返回 "CLI Functional Tests"、"API Functional Tests"、"Web E2E Tests" 等。该函数的映射表应与 `docs/reference/test-type-model.md` 中的映射保持一致。
