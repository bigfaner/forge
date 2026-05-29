---
id: "3"
title: "清理 go.mod 残留依赖并验证构建"
priority: "P1"
estimated_time: "15m"
complexity: "low"
dependencies: [1, 2]
surface-key: "."
surface-type: "cli"
breaking: false
type: "coding.cleanup"
mainSession: false
---

# 3: 清理 go.mod 残留依赖并验证构建

## Description
`go.mod` 中 `mitchellh/hashstructure/v2` 和 `dustin/go-humanize` 两个间接依赖零 import 引用，`go mod tidy` 会自动移除。本任务在任务 1、2 完成后执行，作为最终构建和测试验证关卡。

## Reference Files
- `forge-cli/go.mod`: 行 25 (dustin/go-humanize) 和行 32 (mitchellh/hashstructure/v2) 需被 go mod tidy 移除 (source: proposal.md#Scope-C)

## Acceptance Criteria
- [ ] `go.mod` 中 `mitchellh/hashstructure/v2` 不再出现
- [ ] `go.mod` 中 `dustin/go-humanize` 不再出现
- [ ] `go build ./...` 零错误
- [ ] `go test ./...` 全部通过

## Implementation Notes
- 在 forge-cli/ 目录下执行 `go mod tidy`
- tidy 后运行完整测试套件验证无依赖被误移除
- `go mod graph` 确认这两个依赖不属于任何传递链
