---
id: "11"
title: "Fix validateRecordData os.Exit to return error"
priority: "P1"
estimated_time: "1.5h"
dependencies: [9]
scope: "backend"
breaking: false
type: "coding.refactor"
mainSession: false
---

# 11: Fix validateRecordData os.Exit to return error

## Description
`validateRecordData()` calls `os.Exit()` internally, making the function untestable. Change signature to return `error` instead. Scope includes refactoring the 30+ test cases in `submit_test.go` that test this function. Phase 4 (anti-pattern fix).

## Reference Files
- `docs/proposals/forge-cli-clean-code/proposal.md` — Source proposal
- `forge-cli/internal/cmd/task/submit.go` — Contains `validateRecordData()`
- `forge-cli/internal/cmd/task/submit_test.go` — 30+ test cases to adapt

## Acceptance Criteria
- [ ] `validateRecordData()` returns `error` instead of calling `os.Exit()`
- [ ] All callers updated to handle the returned error
- [ ] 30+ test cases in `submit_test.go` adapted to the new signature
- [ ] 0 non-top-level `os.Exit` calls in `validateRecordData()`
- [ ] `go build ./...` passes
- [ ] `go test ./...` passes

## Hard Rules
- Only change internal functions — top-level RunE handlers that call `os.Exit(0)` in `quality_gate.go` remain unchanged
- The 2 `os.Exit(0)` calls in `quality_gate.go` RunE handlers are explicitly out of scope

## Implementation Notes
- This is the most impactful testability fix — `os.Exit` prevents proper error assertion in tests
- Test refactoring: `submit_test.go` test cases that expected exit behavior now need error assertion
- The caller (likely a RunE handler) should handle the error and decide whether to exit
