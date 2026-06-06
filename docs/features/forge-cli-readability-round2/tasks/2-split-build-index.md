---
id: "2"
title: "拆分 BuildIndex 390 行上帝函数"
priority: "P1"
estimated_time: "2h"
complexity: "medium"
dependencies: [1]
surface-key: ""
surface-type: "cli"
breaking: false
type: "coding.refactor"
mainSession: false
---

# 2: 拆分 BuildIndex 390 行上帝函数

## Description
将 `pkg/task/build.go` 中 390 行的 `BuildIndex` 函数拆分为 ~9 个命名步骤函数，每个步骤函数 ≤ 80 行。保持同包拆分，不改变导出 API 签名。

## Reference Files
- `docs/proposals/forge-cli-readability-round2/proposal.md` — In Scope (InScope-1)、执行阶段排序 (Phase 2)、Success Criteria (SC-1, SC-2)
- `pkg/task/build.go`: 将 BuildIndex 拆分为 ~9 个命名步骤函数 (ref: In Scope)

## Acceptance Criteria
- [ ] `BuildIndex` 及其所有提取出的子函数均 ≤ 80 行
- [ ] `build.go` 文件 ≤ 500 行
- [ ] `go test ./...` 全绿，零行为变更
- [ ] 导出 API 签名不变（同包拆分，不改包名或导入路径）

## Implementation Notes
这是 Phase 2（低风险，机械操作）的文件拆分。`BuildIndex` 当前包含 9 个步骤注释，每个步骤提取为独立命名函数，`BuildIndex` 本体变为顺序调用链。所有新函数保持在 `pkg/task` 包内。
