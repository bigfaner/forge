---
id: "9"
title: "Clean test-bridge alias functions"
priority: "P1"
estimated_time: "2h"
complexity: "medium"
dependencies: [5, 6]
surface-key: "."
surface-type: "cli"
breaking: false
type: "coding.cleanup"
mainSession: false
---

# 9: Clean test-bridge alias functions

## Description
Phase 2b：清理 test-bridge 别名函数。区分两类：(1) 纯粹重导出别名（如 `checkExistingTaskState` 等仅转发调用的包装，可直接删除，测试改为直接调用 `pkg/task` 的函数）；(2) 内部函数导出（如 `getTaskPhase` 在 `validate_index.go` 有 5 处生产调用，需评估测试迁移策略后再决定处理方式）。

## Reference Files
- forge-cli/internal/cmd/task/claim.go:169: `var checkExistingTaskState = task.CheckExistingTaskState` (source: proposal.md#Evidence)
- forge-cli/internal/cmd/task/claim.go:241: `var getTaskPhase = task.GetTaskPhase` — 有 5 处生产调用在 validate_index.go (source: proposal.md#Evidence)
- forge-cli/internal/cmd/task/claim.go:283: `var compareVersionIDs = task.CompareVersionIDs` (source: proposal.md#Evidence)
- forge-cli/internal/cmd/task/validate_index.go:299,333,369,381,410: `getTaskPhase` 生产调用点 (source: proposal.md#Evidence)

## Acceptance Criteria
- [ ] 纯粹重导出别名（`checkExistingTaskState`、`compareVersionIDs`）已删除，测试代码已迁移为直接调用 `pkg/task` 的函数
- [ ] `getTaskPhase` 的处理策略已确定：(a) 保留为生产代码可用的别名，或 (b) 将 validate_index.go 改为直接调用 `task.GetTaskPhase`，视迁移成本而定
- [ ] 所有 test-bridge 别名函数的处理结果记录在案
- [ ] `go build ./...` 和 `go test ./...` 全部通过（SC-11）

## Hard Rules
- `getTaskPhase` 不得简单删除（有 5 处生产调用），必须确保 `validate_index.go` 编译通过

## Implementation Notes
- 使用 `grep -rn '函数名' --include='*.go' | grep -v _test.go | grep -v 'var 函数名'` 验证每个别名在生产代码中的实际调用次数
- 测试文件中的别名调用改为直接导入 `pkg/task` 包的函数

### Test Impact
- Affected test suite(s): `forge-cli/internal/cmd/task/`
- Expected fixture changes: 测试文件中别名调用需改为直接调用
- Risk level: medium（测试代码迁移，需逐文件验证）
