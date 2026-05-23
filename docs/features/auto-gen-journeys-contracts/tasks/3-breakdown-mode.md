---
id: "3"
title: "autogen.go Breakdown 模式：插入 gen-journeys/gen-contracts 并重写依赖解析"
priority: "P0"
estimated_time: "2h"
dependencies: ["1", "2"]
scope: "backend"
breaking: true
type: "coding.feature"
mainSession: false
---

# 3: autogen.go Breakdown 模式：插入 gen-journeys/gen-contracts 并重写依赖解析

## Description

修改 `GetBreakdownTestTasks()` 在 Breakdown 模式的任务列表头部插入 gen-journeys 和 gen-contracts 任务，并将 `resolveBreakdownDeps()` 从硬编码索引重写为基于 `findTaskIndexByPrefix` 的 ID 查找。

## Reference Files
- `docs/proposals/auto-gen-journeys-contracts/proposal.md` — Source proposal
- `forge-cli/pkg/task/autogen.go` — GetBreakdownTestTasks (L84-171), resolveBreakdownDeps (L419-467)
- `forge-cli/pkg/task/types.go` — TypeTestGenJourneys, TypeTestGenContracts (Task 1 新增)

## Acceptance Criteria

- [ ] `GetBreakdownTestTasks()` 为每个 interface type 生成一个 `T-test-gen-journeys-{type}` 任务（TestType=interface type, StrategyKind=interface, Type=TypeTestGenJourneys）
- [ ] `GetBreakdownTestTasks()` 生成一个 `T-test-gen-contracts` 任务（Type=TypeTestGenContracts）
- [ ] gen-journeys 任务排在 eval-journey 之前，gen-contracts 排在 eval-journey 之后、eval-contract 之前
- [ ] 新任务使用 embed 模板渲染 body（通过 autogenTypeToFile 映射）
- [ ] `resolveBreakdownDeps()` 不再使用硬编码索引（evalJourneyIdx=0, evalContractIdx=1, genStart=2）
- [ ] 所有依赖关系通过 `findTaskIndex` 或 `findTaskIndexByPrefix` 查找
- [ ] 依赖链：gen-journeys → eval-journey → gen-contracts → eval-contract → gen-scripts → run → verify-regression
- [ ] findTaskIndex 返回 -1 时 panic 并输出明确错误信息（包含未找到的任务 ID 前缀）
- [ ] 所有现有 Breakdown 模式单测通过

## Hard Rules

- 完全消除 resolveBreakdownDeps 中的算术索引，不允许任何 `tasks[n]` 形式的位置访问
- panic 时的错误信息必须包含：未找到的任务 ID 前缀和当前 tasks 列表中所有 ID（便于调试）
- breaking change：修改了 GetBreakdownTestTasks 的返回值结构（任务列表长度和顺序变化）

## Implementation Notes

- 当前 GetBreakdownTestTasks 按 eval-journey, eval-contract, gen-scripts-per-type, run, verify-regression 顺序生成。新顺序为：gen-journeys-per-type, eval-journey, gen-contracts, eval-contract, gen-scripts-per-type, run, verify-regression
- gen-journeys 任务的 BodyContext 需包含 Mode="breakdown"，以便模板渲染正确的输入源指令
- resolveBreakdownDeps 的非 E2E 分支（validation, specs, clean-code）已使用 findTaskIndex，保持不变
- 测试覆盖：新拓扑 + 旧回归（verify-regression 仍依赖 run，run 仍依赖 gen-scripts）
