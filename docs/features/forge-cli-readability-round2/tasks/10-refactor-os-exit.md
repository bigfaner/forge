---
id: "10"
title: "重构 quality_gate.go 的 os.Exit 反模式"
priority: "P1"
estimated_time: "2h"
complexity: "high"
dependencies: [6, 7, 8, 9]
surface-key: ""
surface-type: "cli"
breaking: false
type: "coding.refactor"
mainSession: false
---

# 10: 重构 quality_gate.go 的 os.Exit 反模式

## Description
将 `internal/cmd/qualitygate/quality_gate.go` 中 4 处 `os.Exit(0)` 替换为 error return，由调用方处理退出码语义。这是 Phase 4（最高风险），因为 `RunQualityGate` 当前零测试覆盖。

替换策略：
- 已处理失败路径（`os.Exit(0)`）→ `return nil`（调用方继续以退出码 0 退出）
- 真正异常路径 → `return fmt.Errorf(...)`（调用方决定退出码）
- `os.Exit` 仅保留在 `cmd/forge/` 入口和 `base.Exit` 统一出口

## Reference Files
- `docs/proposals/forge-cli-readability-round2/proposal.md` — In Scope (InScope-6)、os.Exit 语义替换方案、Key Risks (RunQualityGate 零测试覆盖)、Success Criteria (SC-5, SC-7)
- `internal/cmd/qualitygate/quality_gate.go`: 替换 4 处 os.Exit(0) 为 error return (ref: os.Exit 语义替换方案)

## Acceptance Criteria
- [ ] `quality_gate.go` 无直接 `os.Exit` 调用
- [ ] `os.Exit` 仅存在于 `cmd/forge/` 入口和 `base.Exit` 统一出口
- [ ] `go test ./...` 全绿
- [ ] CLI 退出码语义不变（通过基线输出捕获验证）

## Implementation Notes
这是 Phase 4（最高风险）。`RunQualityGate` 当前零测试覆盖（`os.Exit` 导致不可测试），需通过基线输出捕获（SC-5）验证行为一致性。替换时严格区分"已处理失败"（`return nil`）和"真正异常"（`return error`），确保 CLI 退出码语义不变。

复杂度判定为 high：虽然 AC 仅 4 项，但 `RunQualityGate` 零测试覆盖且含 4 处 `os.Exit`，回归风险高。
