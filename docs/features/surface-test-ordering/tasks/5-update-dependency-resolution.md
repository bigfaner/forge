---
id: "5"
title: "Update dependency resolution chains"
priority: "P0"
estimated_time: "1.5h"
dependencies: ["2", "4"]
surface-key: ""
surface-type: "cli"
breaking: false
type: "coding.enhancement"
mainSession: false
---

# 5: Update dependency resolution chains

## Description
更新 `autogen.go` 中 `resolveBreakdownDeps` 和 `resolveQuickDeps` 函数，适配 gen-journeys 合并（单任务）和 run-tests 拆分（per-surface-key 串行）后的新任务拓扑。Breakdown 模式：gen-journeys → gen-contracts → gen-scripts(并行) → run-tests(串行) → verify-regression。Quick 模式：gen-journeys → run-tests(串行) → verify-regression。

## Reference Files
- `proposal.md#Proposed-Solution` — 3-surface 依赖链示例（Breakdown 和 Quick 两种模式的完整 DAG）
- `proposal.md#Key-Risks` — 串行执行 happy path 延迟回归风险，回滚策略说明

## Acceptance Criteria
- [ ] Breakdown 模式依赖链：`T-test-gen-journeys → T-test-gen-contracts → T-test-gen-scripts-* → T-test-run-{keys} → T-test-verify-regression`
- [ ] Quick 模式依赖链：`T-test-gen-journeys → T-test-run-{keys} → T-test-verify-regression`
- [ ] `T-test-verify-regression` 仅依赖 execution-order 中最后一个 run-test 子任务

## Implementation Notes
- 串行执行 happy path 无额外开销——per-surface 任务仍由调度器调度，串行仅影响启动时机
- 单 surface 项目依赖链与改动前完全一致
