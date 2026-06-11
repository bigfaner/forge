---
id: "1"
title: "Implement per-surface gen→run interleaved dependency wiring"
priority: "P0"
estimated_time: "2h"
complexity: "high"
dependencies: []
surface-key: ""
surface-type: "cli"
breaking: true
type: "coding.enhancement"
mainSession: false
---

# 1: Implement per-surface gen→run interleaved dependency wiring

## Description

当前 pipeline 中 `TypeTestGenScripts` 和 `TypeTestRun` 的 per-surface-key 扩展采用串行链式依赖（gen-A → gen-B → run-A → run-B），导致后端 API bug 必须等所有 surface 的脚本生成完才能发现，且前端测试脚本基于未验证的 API 行为生成。

需要修改 `GenerateTestTasks` 的依赖接线逻辑，将 gen-scripts 和 test-run 按 surface 配对交错执行（gen-A → run-A → gen-B → run-B），使前序 surface 的测试反馈能传递给后续 surface 的脚本生成。

## Reference Files
- `docs/proposals/test-pipeline-interleaved/proposal.md` — Problem, Proposed Solution, Constraints & Dependencies, Key Risks, Success Criteria
- `forge-cli/pkg/task/pipeline.go`: change dependency wiring in GenerateTestTasks and expandPerSurfaceKey (ref: Proposed Solution)
- `forge-cli/pkg/task/pipeline_validate.go`: update PipelineRegistry entries for gen-scripts and test-run nodes (ref: Constraints & Dependencies)
- `forge-cli/pkg/task/pipeline_test.go`: update test expectations for new dependency chain (ref: Key Risks)

## Acceptance Criteria

- [ ] 多 surface 项目中，`T-test-run-{surface-N}` 依赖 `T-test-gen-scripts-{surface-N}`，而非 `T-test-gen-scripts-{surface-N+1}`
- [ ] 多 surface 项目中，`T-test-gen-scripts-{surface-N}` (N>0) 依赖 `T-test-run-{surface-N-1}`，而非 `T-test-gen-scripts-{surface-N-1}`
- [ ] 单 surface 项目中依赖链不变（gen → run，即 T-test-gen-scripts 依赖 upstream，T-test-run 依赖 T-test-gen-scripts）
- [ ] Pipeline registry 验证通过（init-time ValidatePipelineRegistry 不 panic）
- [ ] 所有现有单元测试通过（允许因 registry node 数量变化更新 `ExpectedNodeCount` 测试）

## Hard Rules

- 仅修改以下文件：`forge-cli/pkg/task/pipeline.go`、`forge-cli/pkg/task/pipeline_validate.go`、`forge-cli/pkg/task/pipeline_test.go`
- 不改 task-executor agent 定义、quality-gate 机制、gen-scripts 模板

## Implementation Notes

### 核心改动思路

当前架构中，registry 按 node 顺序依次展开，per-surface-key 节点的后续 task 依赖前一个 task（serial chain）。交错需要打破 gen-scripts 和 test-run 两个独立 node 的边界。

方案：修改 `GenerateTestTasks` 中的展开逻辑，当遇到 `TypeTestRun` 节点时，不使用默认的 serial chain wiring，而是让每个 test-run-{surface-N} 依赖对应的 gen-scripts-{surface-N}。同时，在 gen-scripts 的后续 task（i>0）中，改为依赖前一个 surface 的 test-run 而非前一个 gen-scripts。

具体实现：
1. 在 `GenContext` 中添加 `GenScriptsMap map[string]string`，记录 surface-key → gen-scripts-task-ID 的映射
2. gen-scripts 展开时，第一个 task 依赖 upstream（不变），后续 task（i>0）改为依赖前一个 surface 的 test-run task（需要查找 `ctx.RunTestMap`）
3. test-run 展开时，每个 task 依赖对应 surface 的 gen-scripts task（通过 `ctx.GenScriptsMap` 查找），而非 serial chain
4. 添加 `RunTestMap map[string]string`，记录 surface-key → test-run-task-ID 的映射

### Test Impact
- Affected test suite(s): `forge-cli/pkg/task/pipeline_test.go`
- Expected fixture changes: `ExpectedNodeCount` 可能不变（registry node 数量不变，只改接线）
- Risk level: medium — 依赖接线是核心逻辑，需确保 cycle detection 仍通过

### Key Risks

- prompt 模板中的新指令与 task-executor Pause Protocol 冲突（L/M）— 新指令作为 TASK-CONSTRAINTS 级别补充，不覆盖 EXTREMELY-IMPORTANT 层级
