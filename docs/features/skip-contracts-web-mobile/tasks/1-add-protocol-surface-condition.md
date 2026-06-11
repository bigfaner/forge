---
id: "1"
title: "Add protocol surface condition to skip gen-contracts pipeline"
priority: "P0"
estimated_time: "1-2h"
complexity: "high"
dependencies: []
surface-key: ""
surface-type: "cli"
breaking: false
type: "coding.enhancement"
mainSession: false
---

# 1: Add protocol surface condition to skip gen-contracts pipeline

## Description

PipelineRegistry 的 gen-contracts 和 eval-contract 节点对所有 feature 无条件生成，但纯 Web/Mobile 交互场景无法产出有效契约文件，导致后续 gen-test-scripts 静默跳过这些 journey。需要新增 `CondHasProtocolSurfaceTask` 条件：检查业务任务的 surface-type 是否包含 tui/cli/api，若全部为 web/mobile 则跳过这两个节点，并将 gen-scripts 的依赖链调整为直接依赖 gen-journeys。

## Reference Files
- `docs/proposals/skip-contracts-web-mobile/proposal.md` — Proposed Solution (Pipeline层), Scope > In Scope, Success Criteria SC-1/2/3/4/7, Key Risks
- `forge-cli/pkg/task/pipeline_validate.go` — PipelineRegistry 定义 (line 238-245)，需修改 T-test-gen-contracts 和 T-eval-contract 的 GenerateCondition
- `forge-cli/pkg/task/build.go` — 现有条件函数 (CondHasTestableTasks 等)，需新增 CondHasProtocolSurfaceTask (ref: GenerateCondFunc)
- `forge-cli/pkg/task/pipeline.go` — GenerateTestTasks 逻辑，需处理依赖链动态调整 (ref: Phase 1: Expand all nodes)

## Acceptance Criteria
- [ ] SC-1: 纯 Web feature（所有业务任务 surface-type: web）不生成 T-test-gen-contracts 和 T-eval-contract
- [ ] SC-2: 纯 Mobile feature 同 SC-1，不生成 gen-contracts/eval-contract
- [ ] SC-3: 混合 surface feature（部分 api + 部分 web）正常生成 T-test-gen-contracts 和 T-eval-contract
- [ ] SC-4: 同一项目多 surface 配置下，前端-only feature 跳过 gen-contracts，后端 feature 保留
- [ ] SC-7: 已有 API/CLI/TUI 流水线行为不变（现有测试全部通过）
- [ ] surface-type 缺失/为空/未知值时保守不跳过，输出 WARN 日志

## Implementation Notes
- 参考 `UISurfaceOnly` 字段实现模式 (build.go:63, pipeline.go:179)
- `CondHasProtocolSurfaceTask` 遍历 businessTasks，检查是否有任何任务的 SurfaceType 为 tui/cli/api。protocol-level 类型集合定义为常量
- 依赖链调整方案：当 gen-contracts 被跳过时，gen-scripts 的 DependsOn resolve 应自然降级到 gen-journeys（通过现有 ResolveUpstream 机制或新增 resolve function）
- INFO 日志格式：`[skip] T-test-gen-contracts: no protocol-level surface tasks found (all tasks: web)`
- 单元测试覆盖：纯 web、纯 mobile、混合 api+web、纯 api、空任务列表、未知 surface-type 场景
