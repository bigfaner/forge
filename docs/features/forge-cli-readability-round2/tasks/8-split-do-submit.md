---
id: "8"
title: "拆分 doSubmit 131 行函数"
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

# 8: 拆分 doSubmit 131 行函数

## Description
将 `internal/cmd/task/submit.go` 中 131 行的 `doSubmit` 函数拆分为命名子函数。

## Reference Files
- `docs/proposals/forge-cli-readability-round2/proposal.md` — In Scope (InScope-8)、执行阶段排序 (Phase 3)、Success Criteria (SC-1)
- `internal/cmd/task/submit.go`: 拆分 doSubmit 为命名子函数 (ref: In Scope)

## Acceptance Criteria
- [ ] `doSubmit` 及其所有提取出的子函数均 ≤ 80 行
- [ ] 所有函数嵌套 ≤ 4 层
- [ ] `go test ./...` 全绿，零行为变更
- [ ] 文件 ≤ 500 行

## Implementation Notes
这是 Phase 3（中等风险）。`doSubmit` 当前 131 行，嵌套 4 层（已在限制内）。主要工作是提取验证、数据处理、提交等步骤为独立命名函数。
