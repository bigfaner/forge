---
id: "11"
title: "Reorganize internal/cmd/ package structure and split large files"
priority: "P1"
estimated_time: "4h"
complexity: "high"
dependencies: [2, 6]
surface-key: "."
surface-type: "cli"
breaking: true
type: "coding.refactor"
mainSession: false
---

# 11: Reorganize internal/cmd/ package structure and split large files

## Description
Phase 2c（internal/cmd/ 部分）：将 `internal/cmd/` 下 15 个顶层命令文件子包化，统一命令注册模式。同步拆分超大文件（`quality_gate.go` 1067 行、`init.go` 591 行、`init_surfaces.go` 550 行、`task/validate_index.go` 521 行、`task/tree.go` 504 行）。

## Reference Files
- docs/conventions/package-organization.md: 包组织规范，定义目标态结构 (source: proposal.md#Scope item 7)
- forge-cli/internal/cmd/: 15 个顶层 .go 文件，需子包化 (source: proposal.md#Scope item 7)
- forge-cli/internal/cmd/root.go: 命令注册入口，需更新以引用子包化后的命令 (source: proposal.md#Scope item 7)
- forge-cli/internal/cmd/quality_gate.go: 1067 行，需拆分为质量检查核心逻辑 + 报告生成逻辑 (source: proposal.md#Scope item 10a)

## Acceptance Criteria
- [ ] `internal/cmd/` 下仅保留基础设施文件：`root.go`、`output.go`、`surfaces.go`、`surfaces_detect.go`（SC-5，分类标准：不含 Cobra `Run`/`RunE`）
- [ ] 所有命令实现文件已移入对应子包（如 `internal/cmd/quality/`、`internal/cmd/init/`）
- [ ] 超大文件已拆分：`quality_gate.go`、`init.go`、`init_surfaces.go`、`validate_index.go`、`tree.go` 各拆分后均 <500 行（SC-10）
- [ ] `go build ./...` 和 `go test ./...` 全部通过（SC-11）

## Hard Rules
- 每个子包移动后立即 `go build ./...` + `go test ./...` 验证
- 不改变任何命令的 CLI 行为（命令名、参数、输出格式不变）

## Implementation Notes
- 执行顺序：先拆分超大文件（在当前位置），再执行子包化移动
- 每个子包移动为独立提交，可单独回退
- 子包内的文件可共享该子包的 `output.go` 等工具文件

### Test Impact
- Affected test suite(s): `forge-cli/internal/cmd/`, `forge-cli/internal/cmd/task/`
- Expected fixture changes: import 路径更新
- Risk level: high（大规模文件移动 + import 更新）

### Breaking Change Assessment
- 修改的是 `internal/` 层，无外部消费者（`internal` 语义保证）
- 但 monorepo 内其他模块可能有引用（Task 7 审计确认）
