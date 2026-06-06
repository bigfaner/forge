---
id: "7"
title: "拆分 runList 217 行函数"
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

# 7: 拆分 runList 217 行函数

## Description
将 `internal/cmd/task/list.go` 中 217 行的 `runList` 函数拆分为命名子函数，降低单函数复杂度。

## Reference Files
- `docs/proposals/forge-cli-readability-round2/proposal.md` — In Scope (InScope-7)、执行阶段排序 (Phase 3)、Success Criteria (SC-1, SC-3)
- `internal/cmd/task/list.go`: 拆分 runList 为命名子函数 (ref: In Scope)

## Acceptance Criteria
- [ ] `runList` 及其所有提取出的子函数均 ≤ 80 行
- [ ] 所有函数嵌套 ≤ 4 层
- [ ] `go test ./...` 全绿，零行为变更
- [ ] 文件 ≤ 500 行

## Implementation Notes
这是 Phase 3（中等风险）。`runList` 当前 217 行，最大嵌套 5 层。提取格式化、过滤、排序等逻辑为独立命名函数，同时用 early return 平坦化嵌套。
