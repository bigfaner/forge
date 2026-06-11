---
id: "fix.2"
title: "Fix compilation errors from task 2.5 (checkUnmetDeps, status_test)"
priority: "P0"
dependencies: []
status: 
type: "coding.fix"
breaking: true
---

# fix.2: Fix compilation errors from task 2.5 (checkUnmetDeps, status_test)

## Problem

Task 2.5 removed `checkUnmetDeps` from `status.go` but `submit.go` still calls it. Also removed `getTransitionHint` and `getTransitionAction` from `status.go` but `status_test.go` still references them.

## Scope
- `forge-cli/internal/cmd/submit.go:265` — calls `checkUnmetDeps` which no longer exists
- `forge-cli/internal/cmd/status_test.go:279-304` — references `getTransitionHint`, `getTransitionAction`, `strings` which no longer exist

## Acceptance Criteria
- [ ] `go build ./...` passes (0 errors)
- [ ] `go test ./internal/cmd/` passes (0 failures)
- [ ] No test regressions
