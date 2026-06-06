---
id: "1"
title: "Switch gen-test-scripts to per-surface-key expansion in pipeline.go"
priority: "P0"
estimated_time: "1h"
complexity: "medium"
dependencies: []
surface-key: ""
surface-type: "cli"
breaking: true
type: "coding.enhancement"
mainSession: false
---

# 1: Switch gen-test-scripts to per-surface-key expansion in pipeline.go

## Description

gen-test-scripts 任务在 pipeline.go 中使用 `per-surface-type` expansion（第 741-746 行），导致任务文件命名为 `gen-test-scripts-api.md` 而非 `gen-test-scripts-backend.md`。将其切换为 `per-surface-key` expansion，与 run-tests（第 749-754 行）保持一致。

## Reference Files
- `docs/proposals/surface-key-test-alignment/proposal.md` — Problem, Evidence, Proposed Solution, Scope, Success Criteria
- `forge-cli/pkg/task/pipeline.go:739-746` — gen-test-scripts registry entry using per-surface-type (ref: Evidence)
- `forge-cli/pkg/task/pipeline.go:747-754` — run-tests registry entry using per-surface-key as reference (ref: Proposed Solution)

## Acceptance Criteria
- [ ] gen-test-scripts registry entry 的 `Expansion` 字段从 `"per-surface-type"` 改为 `"per-surface-key"`
- [ ] Key 和 ID 模板从 `{surface-type}` 改为 `{surface-key}`：`gen-test-scripts-{surface-key}`，`T-test-gen-scripts-{surface-key}`
- [ ] 单 surface 项目通过 `isSingleSurface` 逻辑仍产生无后缀的 `gen-test-scripts.md`（回归验证）
- [ ] `go test ./pkg/task/...` 全部通过

## Implementation Notes

### Test Impact
- Affected test suite(s): `forge-cli/pkg/task/`, `forge-cli/internal/cmd/`
- Expected fixture changes: test fixtures referencing `gen-test-scripts-{type}` naming pattern
- Risk level: medium

run-tests 已使用 `per-surface-key` expansion 并正常工作，可作为参考实现。修改后需确认 `expandPerSurfaceKey` 函数正确处理 gen-test-scripts 的 expansion。
