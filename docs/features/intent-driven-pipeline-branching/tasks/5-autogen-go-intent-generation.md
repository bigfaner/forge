---
id: "5"
title: "autogen.go — intent-driven task generation and wiring"
priority: "P0"
estimated_time: "3-4h"
complexity: "high"
dependencies: [4]
surface-key: "."
surface-type: "cli"
breaking: true
type: "coding.feature"
mainSession: false
---

# 5: autogen.go — intent-driven task generation and wiring

## Description

在 `autogen.go` 中实现 intent 驱动的测试任务生成和依赖接线。核心改动：
1. `GetBreakdownTestTasks()` 和 `GetQuickTestTasks()` 接收 intent 参数——当 intent 为 `refactor`/`cleanup` 时跳过测试任务（gen-journeys/gen-contracts/gen-scripts/run-tests），但仍生成 validate-code/clean-code 等验证任务
2. `resolveBreakdownDeps()` 和 `resolveQuickDeps()` 感知 intent——refactor/cleanup 时跳过 run-tests 节点，把下游任务直接接到 business tasks 尾部

接线矩阵：
- **new-feature（Breakdown）**：business tasks → gen-journeys → ... → run-tests → validate-code → clean-code → consolidate-specs（不变）
- **refactor（Breakdown）**：business tasks → validate-code → clean-code → consolidate-specs（validate-code 依赖最后一个 business task）
- **refactor（Quick）**：business tasks → clean-code → doc-drift（clean-code 依赖最后一个 business task）
- **cleanup（Quick）**：business tasks → clean-code → doc-drift（同 refactor Quick）
- **零 business task 保护**：refactor/cleanup 下无 business task 时不生成下游任务

## Reference Files
- `forge-cli/pkg/task/autogen.go`: `GetBreakdownTestTasks()`（L182-303）、`GetQuickTestTasks()`（L316-405）、`resolveBreakdownDeps()`（L591-643）、`resolveQuickDeps()`（L649-683）、`ResolveFirstTestDep()`（L820-870） (source: proposal.md#In-Scope, items 7-8)
- `forge-cli/pkg/task/build.go`: `GenerateTestTasks()` 函数（L480-489）调用 GetBreakdownTestTasks/GetQuickTestTasks，需传入 intent 参数 (source: proposal.md#Feasibility-Assessment, autogen 层)
- `docs/proposals/intent-driven-pipeline-branching/proposal.md#Feasibility-Assessment`: autogen.go 需覆盖 5 种有效接线场景 + 零 business task 边界分支

## Acceptance Criteria
- [ ] `GetBreakdownTestTasks()` 和 `GetQuickTestTasks()` 接收 intent 参数——当 intent 为 `refactor`/`cleanup` 时跳过 `auto.Test.*` 管控的测试任务生成（gen-journeys/gen-contracts/gen-scripts/run-tests），但 `auto.Validation.*`、`auto.ConsolidateSpecs.*`、`auto.CleanCode.*` 管控的任务（validate-code/clean-code/consolidate-specs）仍正常生成
- [ ] `resolveBreakdownDeps()` 和 `resolveQuickDeps()` 感知 intent——refactor/cleanup 时 validate-code/clean-code 直接依赖最后一个 business task（不查找 lastRunID），new-feature 保持现有逻辑不变
- [ ] 零 business task 保护：intent 为 refactor/cleanup 但 business task 列表为空时，不生成 validate-code/clean-code 等下游任务，避免悬空的 `depends_on` 引用
- [ ] `intent: new-feature` 的 pipeline 行为与当前完全一致（golden file 兼容）——测试任务和依赖链结构不变
- [ ] 单元测试覆盖全部 5 种接线场景（new-feature Breakdown/Quick、refactor Breakdown/Quick、cleanup Quick）及零 business task 边界情况，测试全部通过

## Implementation Notes

### Test Impact
- Affected test suite(s): `forge-cli/pkg/task/`
- Expected fixture changes: 无 fixture 变更，纯单元测试
- Risk level: high

- `GenerateTestTasks()` 函数（build.go:480）是 `GetBreakdownTestTasks()`/`GetQuickTestTasks()` 的入口，需同步接收 intent 参数
- cleanup 的 Breakdown 路径不会到达 autogen.go——build.go 在 intent=`cleanup` 时已将 mode 强制为 Quick。autogen.go 无需处理此组合
- refactor（Quick）和 cleanup（Quick）的接线逻辑相同：最后一个 business task → clean-code → doc-drift
- `ResolveFirstTestDep()` 也需感知 intent——refactor/cleanup 时无 gen-journeys 任务，不需要 resolve first test dep
