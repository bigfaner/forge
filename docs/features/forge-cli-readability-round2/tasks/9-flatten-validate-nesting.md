---
id: "9"
title: "平坦化 validate.go 嵌套过深的 validator 方法"
priority: "P1"
estimated_time: "1h"
complexity: "medium"
dependencies: [2, 3, 4, 5]
surface-key: ""
surface-type: "cli"
breaking: false
type: "coding.refactor"
mainSession: false
---

# 9: 平坦化 validate.go 嵌套过深的 validator 方法

## Description
将 `internal/cmd/task/validate.go` 中嵌套过深的 validator 方法（最大嵌套 5 层）平坦化至 ≤ 4 层，使用 early return 和 guard clause。

## Reference Files
- `docs/proposals/forge-cli-readability-round2/proposal.md` — In Scope (InScope-9)、执行阶段排序 (Phase 3)、Success Criteria (SC-1, SC-3)
- `internal/cmd/task/validate.go`: 平坦化嵌套过深的 validator 方法 (ref: In Scope)

## Acceptance Criteria
- [ ] 所有函数嵌套 ≤ 4 层（当前最大 5 层）
- [ ] 所有函数 ≤ 80 行（当前最大 66 行，嵌套平坦化后可能增加行数）
- [ ] `go test ./...` 全绿，零行为变更
- [ ] 文件 ≤ 500 行（当前 573 行，需同步拆分或提取）

## Implementation Notes
这是 Phase 3（中等风险）。`validate.go` 当前 573 行，`validateGateIntegrity` 最大嵌套 5 层。使用 early return 平坦化嵌套。若平坦化后文件仍超 500 行，可提取部分 validator 到同包新文件。
