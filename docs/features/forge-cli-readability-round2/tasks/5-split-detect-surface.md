---
id: "5"
title: "拆分 detect_surface.go 提取信号表"
priority: "P1"
estimated_time: "1h"
complexity: "medium"
dependencies: [1]
surface-key: ""
surface-type: "cli"
breaking: false
type: "coding.refactor"
mainSession: false
---

# 5: 拆分 detect_surface.go 提取信号表

## Description
将 `pkg/forgeconfig/detect_surface.go`（963 行）中的信号映射表（前 ~150 行）提取到 `detect_surface_signals.go`，统一各生态的推断模式。

## Reference Files
- `docs/proposals/forge-cli-readability-round2/proposal.md` — In Scope (InScope-5)、Success Criteria (SC-2)
- `pkg/forgeconfig/detect_surface.go`: 提取信号表和推断函数到新文件 (ref: In Scope)

## Acceptance Criteria
- [ ] `detect_surface.go` 和 `detect_surface_signals.go` 各 ≤ 500 行
- [ ] 信号映射表集中在 `detect_surface_signals.go`
- [ ] `go test ./...` 全绿，零行为变更
- [ ] 所有函数 ≤ 80 行

## Implementation Notes
这是 Phase 2（低风险，机械操作）。`detect_surface.go` 前 150 行全是信号映射表，推断函数按生态重复模式。提取信号表到独立文件可显著提升主文件可读性。
