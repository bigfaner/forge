---
id: "fix.4"
title: "Fix undefined PreserveRuntimeFields in preserve_test.go"
priority: "P0"
dependencies: []
status: 
type: "coding.fix"
breaking: true
---

# fix.4: Fix undefined PreserveRuntimeFields in preserve_test.go

## Problem

`forge-cli/pkg/index/preserve_test.go:37` references `PreserveRuntimeFields` which is not defined/exported. Likely the function was renamed or not exported during task 2.7, but the test file still references the old name.

## Acceptance Criteria
- [ ] `go build ./...` passes (0 errors)
- [ ] `go test ./pkg/index/...` passes (0 failures)
