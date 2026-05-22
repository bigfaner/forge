---
id: "fix.6"
title: "Fix undefined validateJourneyName in test_promote_test.go"
priority: "P0"
dependencies: []
status: 
type: "coding.fix"
breaking: true
---

# fix.6: Fix undefined validateJourneyName in test_promote_test.go

## Problem

6 lines in `forge-cli/internal/cmd/test_promote_test.go` reference `validateJourneyName` which is undefined. Task 2.10 added `validateJourneyName` but tests cannot resolve it.

## Acceptance Criteria
- [ ] `go build ./...` passes (0 errors)
- [ ] `go test ./internal/cmd/...` passes (0 failures)
