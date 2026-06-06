---
id: "4"
title: "拆分 pipeline.go 提取校验逻辑"
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

# 4: 拆分 pipeline.go 提取校验逻辑

## Description
将 `pkg/task/pipeline.go`（1103 行）中的校验逻辑提取到 `pipeline_validate.go`，重组 var 块减少阅读打断。保持同包拆分。

## Reference Files
- `docs/proposals/forge-cli-readability-round2/proposal.md` — In Scope (InScope-4)、Success Criteria (SC-2, SC-3)
- `pkg/task/pipeline.go`: 提取校验逻辑到新文件，重组 var 块 (ref: In Scope)

## Acceptance Criteria
- [ ] `pipeline.go` 和 `pipeline_validate.go` 各 ≤ 500 行
- [ ] var 块和类型定义集中放置，不穿插在函数间打断阅读流
- [ ] `go test ./...` 全绿，零行为变更
- [ ] 所有函数 ≤ 80 行

## Implementation Notes
这是 Phase 2（低风险，机械操作）。`pipeline.go` 当前有 445 行非函数代码（var 块、类型定义）打断阅读流。提取校验逻辑到独立文件，同时将共享的 var 块和类型定义集中在文件头部或独立区域。
