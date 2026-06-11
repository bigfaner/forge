---
id: "3"
title: "拆分 config.go 为三文件（配置读写 / reflect 路径遍历 / AutoConfig）"
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

# 3: 拆分 config.go 为三文件

## Description
将 `pkg/forgeconfig/config.go`（1365 行）按职责边界拆分为三个文件：
- `config.go` — 配置读写
- `config_reflect.go` — reflect 路径遍历（`setByPath` 等函数）
- `config_auto.go` — AutoConfig 默认值

保持同包拆分，不改变导出 API 签名。

## Reference Files
- `docs/proposals/forge-cli-readability-round2/proposal.md` — In Scope (InScope-3)、Success Criteria (SC-2)
- `pkg/forgeconfig/config.go`: 按 3 种职责拆分为 3 个文件 (ref: In Scope)

## Acceptance Criteria
- [ ] `config.go` 仅包含配置读写相关函数
- [ ] `config_reflect.go` 仅包含 reflect 路径遍历函数
- [ ] `config_auto.go` 仅包含 AutoConfig 默认值函数
- [ ] 每个文件 ≤ 500 行
- [ ] `go test ./...` 全绿，零行为变更

## Implementation Notes
这是 Phase 2（低风险，机械操作）。`config.go` 当前混合了 3 种职责：配置读写、reflect 路径遍历（最大嵌套 7 层）、AutoConfig 默认值。按职责边界拆分到同包不同文件，不改包名或导入路径。
