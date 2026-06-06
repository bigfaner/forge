---
id: "6"
title: "拆分 runExtract 304 行函数"
priority: "P1"
estimated_time: "2h"
complexity: "medium"
dependencies: [2, 3, 4, 5]
surface-key: ""
surface-type: "cli"
breaking: false
type: "coding.refactor"
mainSession: false
---

# 6: 拆分 runExtract 304 行函数

## Description
将 `internal/cmd/forensic/extract.go` 中 304 行的 `runExtract` 函数拆分为解析/聚合/输出三个阶段，使用 early return 平坦化 7+ 层嵌套。

## Reference Files
- `docs/proposals/forge-cli-readability-round2/proposal.md` — In Scope (InScope-2)、执行阶段排序 (Phase 3)、Success Criteria (SC-1, SC-3)
- `internal/cmd/forensic/extract.go`: 将 runExtract 拆分为解析/聚合/输出阶段 (ref: In Scope)

## Acceptance Criteria
- [ ] `runExtract` 及其所有提取出的子函数均 ≤ 80 行
- [ ] 所有函数嵌套 ≤ 4 层（使用 early return / guard clause 平坦化）
- [ ] `go test ./...` 全绿，零行为变更
- [ ] 文件 ≤ 500 行

## Implementation Notes
这是 Phase 3（中等风险）的函数提取与嵌套平坦化。当前最大嵌套 7+ 层，需使用 early return 和 guard clause 将嵌套降至 ≤ 4 层。按解析、聚合、输出三个阶段提取为独立函数。
