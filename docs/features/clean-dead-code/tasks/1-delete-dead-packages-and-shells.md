---
id: "1"
title: "删除 internal/docsync 目录和空壳文件"
priority: "P1"
estimated_time: "15m"
complexity: "low"
dependencies: []
surface-key: "."
surface-type: "cli"
breaking: false
type: "coding.cleanup"
mainSession: false
---

# 1: 删除 internal/docsync 目录和空壳文件

## Description
删除 3 个无功能代码的文件/目录：(1) `internal/docsync/` 目录仅有 2 个测试文件（~1100 行），无生产代码可供测试；(2) `internal/cmd/errors.go` 仅含包声明和注释重定向到 base 子包；(3) `internal/cmd/worktree/worktree.go` 仅含包文档，零功能代码。

## Reference Files
- `forge-cli/internal/docsync/`: 整个目录删除，确认无外部 import (source: proposal.md#Scope-A)
- `forge-cli/internal/cmd/errors.go`: 仅 5 行包声明+注释，确认同目录其他文件不依赖 (source: proposal.md#Scope-A)
- `forge-cli/internal/cmd/worktree/worktree.go`: 仅 16 行包文档，确认同目录其他文件不依赖 (source: proposal.md#Scope-A)

## Acceptance Criteria
- [ ] `internal/docsync/` 目录已删除，`grep -r "internal/docsync" forge-cli/` 无结果
- [ ] `internal/cmd/errors.go` 已删除
- [ ] `internal/cmd/worktree/worktree.go` 已删除
- [ ] `go build ./...` 零错误

## Implementation Notes
- 删除前 grep 确认无 `//go:embed` 或 build tag 引用这些文件
- `internal/docsync/` 测试文件引用了 `pkg/feature` 和 `pkg/task`，但这些包不导入 docsync，删除不影响任何测试
