---
id: "2"
title: "Add system type interception in BuildIndex and validate-index"
priority: "P1"
estimated_time: "1h"
dependencies: ["1"]
scope: "backend"
breaking: false
type: "coding.feature"
mainSession: false
---

# 2: Add system type interception in BuildIndex and validate-index

## Description

在 `BuildIndex()`（build.go）和 `validate-index`（validate_index.go）中添加系统类型拦截逻辑：非自动生成任务使用系统类型时拒绝。利用 Task 1 定义的 `IsSystemType()` 函数。

拦截规则：
- 自动生成任务（`isAutoGenTaskID()` 匹配）→ 不受限制
- 非自动生成任务 + 类型命中 SystemTypes → 报错，错误信息包含具体类型和系统类型列表
- 非自动生成任务 + 业务类型 → 通过

## Reference Files
- `docs/proposals/system-type-exclusion/proposal.md` — Source proposal
- `forge-cli/pkg/task/build.go` — BuildIndex 函数、isAutoGenTaskID
- `forge-cli/pkg/task/types.go` — IsSystemType、SystemTypes
- `forge-cli/internal/cmd/validate_index.go` — validateTasks 函数
- `forge-cli/internal/cmd/quality_gate.go` — addFixTask（确认 coding.fix/coding.cleanup 不受影响）

## Acceptance Criteria

- [ ] `BuildIndex()` 拒绝非自动生成任务使用系统类型
- [ ] `BuildIndex()` 允许自动生成任务使用系统类型（现有行为不变）
- [ ] `validate-index` 拒绝非自动生成任务使用系统类型
- [ ] 错误信息包含具体的非法类型和完整系统类型列表
- [ ] 质量门 fix 任务（`coding.fix`、`coding.cleanup`）通过校验
- [ ] `go test ./forge-cli/...` 通过

## Hard Rules

- BuildIndex 和 validate-index 的校验逻辑须一致
- 不修改 `isAutoGenTaskID()` 逻辑
- 不修改 `IsTestableType()` 逻辑

## Implementation Notes

- BuildIndex 中现有类型校验在 ~line 139，新校验应在其后、同级别的位置
- validate_index.go 中 validateTasks 在 ~line 115，新校验应在现有 type 校验之后
- 需要同步 validate_index.go 的 `validateTasks` 签名以传入 isAutoGen 信息（或直接在 validate_index 中也用 `isAutoGenTaskID`）
- Error message format: `"task '%s': type '%s' is a system-reserved type (reserved: %s)"`
