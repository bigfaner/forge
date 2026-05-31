---
id: "12"
title: "Reorganize pkg/ layer"
priority: "P1"
estimated_time: "5h"
complexity: "high"
dependencies: [2, 6, 7, 8, 9, 11]
surface-key: "."
surface-type: "cli"
breaking: true
type: "coding.refactor"
mainSession: false
---

# 12: Reorganize pkg/ layer

## Description
Phase 2c（pkg/ 部分）：以 Task 2 规范和 Task 1 依赖图为指导，重新设计 `pkg/` 层包结构。合并明确目标（`project/` → `task/` 或 `workspace/`、`infocmd/` → `shared/` 或 `internal/cmd/`、~3 个小工具包 → `util/`）。评估 `prompt/` 和 `research/` 的最终归向。不保留兼容层。每个合并后的包须在 `doc.go` 中列出子领域及职责边界。

## Reference Files
- docs/conventions/package-organization.md: 包组织规范和三层模型 (source: proposal.md#Scope item 8)
- docs/features/forge-cli-codebase-standards/pkg-dependency-graph.md: 依赖图事实基线 (source: Task 1)
- forge-cli/pkg/: 当前 17 个包的目标态映射表 (source: proposal.md#Scope item 8)

## Acceptance Criteria
- [ ] `pkg/` 层包数量不超过 14 个（SC-9）
- [ ] 每个合并后的包包含 `doc.go`，列出子领域及职责边界
- [ ] 无循环依赖：合并后 `go build ./...` 通过
- [ ] `go build ./...` 和 `go test ./...` 全部通过（SC-11）
- [ ] `pkg/` 层超大文件的拆分可行性评估已产出并记录（SC-13）

## Hard Rules
- 包合并前必须检查是否引入循环依赖——若 A 依赖 B 且 B 依赖 A 的部分功能，不可合并
- 执行顺序：先处理 leaf 包（`version/`、`types/`），再中间包，最后被依赖最多的包（`task/`）
- 每步移动后立即 `go build` + `go test` 验证
- 若 Task 7 审计发现跨模块依赖，则仅执行 3 个明确合并目标（`project/`、`infocmd/`、~1 个小工具包），上限 16 个包（SC-12f）

## Implementation Notes
- 先检查 Task 7 审计结果决定是否完整执行
- 合并回退标准：循环 import、单一包超过 15 个文件、`go vet` 命名冲突 → `git revert`
- `pkg/util/` 作为唯一允许被多个 `pkg/` 包共享的工具包，不依赖任何其他 forge-cli `pkg/` 包
- `pkg/prompt/` 和 `pkg/research/` 的归向由 Task 1 依赖图分析裁决

### Test Impact
- Affected test suite(s): 全部 `forge-cli/` 测试（import 路径变更）
- Expected fixture changes: 所有引用被合并包的 import 路径需更新
- Risk level: high（最大 blast radius 的任务）

### Breaking Change Assessment
- `pkg/` 层是公共 API 层，但 Task 7 审计确认无外部消费者
- `internal/` 不受此任务影响
