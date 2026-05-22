---
id: "fix.5"
title: "Fix ErrFileNotFound type mismatch in worktree.go"
priority: "P0"
dependencies: []
status: 
type: "coding.fix"
breaking: true
---

# fix.5: Fix ErrFileNotFound type mismatch in worktree.go

## Problem

`forge-cli/internal/cmd/worktree.go:762` passes `ErrFileNotFound` (a function `func(path string) *AIError`) as an `ErrorCode` value to `NewAIError`. `ErrFileNotFound` is a factory function, not an `ErrorCode` constant.

## Acceptance Criteria
- [ ] `go build ./...` passes (0 errors)
