---
id: "1"
title: "删除死代码：requireSurfaceInference、extractScope、extractBulletItems"
priority: "P0"
estimated_time: "1h"
complexity: "medium"
dependencies: []
surface-key: ""
surface-type: "cli"
breaking: false
type: "coding.cleanup"
mainSession: false
---

# 1: 删除死代码：requireSurfaceInference、extractScope、extractBulletItems

## Description

删除已确认的死代码，为后续重构清理基础。这是 Phase 1（最低风险），所有后续阶段依赖此任务完成后的干净基线。

需要删除的函数：
- `requireSurfaceInference`（`internal/cmd/qualitygate/quality_gate.go`）— 无调用方
- `extractScope`（`pkg/task/extract.go`）— 无调用方
- `extractBulletItems`（`pkg/task/extract.go`）— 仅被 `extractScope` 调用，删除后者后成为死代码

同步删除上述函数对应的测试用例。

## Reference Files
- `docs/proposals/forge-cli-readability-round2/proposal.md` — In Scope (InScope-10)、Success Criteria (SC-4, SC-5, SC-6)
- `internal/cmd/qualitygate/quality_gate.go`: 删除 requireSurfaceInference 函数 (ref: In Scope)
- `pkg/task/extract.go`: 删除 extractScope 和 extractBulletItems 函数 (ref: In Scope)

## Acceptance Criteria
- [ ] `requireSurfaceInference` 函数已从 `internal/cmd/qualitygate/quality_gate.go` 中删除
- [ ] `extractScope` 和 `extractBulletItems` 函数已从 `pkg/task/extract.go` 中删除
- [ ] 上述函数对应的测试用例已同步删除
- [ ] `go test ./...` 全绿，零行为变更

## Implementation Notes
删除前用 `grep` 确认函数确实无其他调用方。`extractBulletItems` 仅被 `extractScope` 调用，删除 `extractScope` 后自然成为死代码，一并删除。允许同步删除对应测试用例（SC-4 放宽）。
