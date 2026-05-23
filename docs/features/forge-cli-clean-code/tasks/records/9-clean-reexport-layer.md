---
status: "completed"
started: "2026-05-24 02:38"
completed: "2026-05-24 03:09"
time_spent: "~31m"
---

# Task Record: 9 Clean up re-export layer in errors.go and output.go

## Summary
Removed re-export layer in errors.go and output.go. All cmd package callers now import base directly. Debugf kept as local inline (base version has variadic expansion bug).

## Changes

### Files Created
无

### Files Modified
- forge-cli/internal/cmd/errors.go
- forge-cli/internal/cmd/output.go
- forge-cli/internal/cmd/quality_gate.go
- forge-cli/internal/cmd/research.go
- forge-cli/internal/cmd/proposal.go
- forge-cli/internal/cmd/lesson.go
- forge-cli/internal/cmd/version.go
- forge-cli/internal/cmd/errors_test.go
- forge-cli/internal/cmd/output_test.go
- forge-cli/internal/cmd/output_contract_test.go
- forge-cli/internal/cmd/slug_width_test.go
- forge-cli/internal/cmd/characterization_test.go
- forge-cli/internal/cmd/integration_test.go
- forge-cli/internal/cmd/proposal_test.go

### Key Decisions
- Kept Debugf as inline function in output.go because base.Debugf has a bug (passes args instead of args... to Fprintf)
- Hard Rule prohibits modifying base/output.go, so the variadic bug remains in base but the correct version is preserved in cmd

## Test Results
- **Tests Executed**: Yes
- **Passed**: 271
- **Failed**: 0
- **Coverage**: 67.4%

## Acceptance Criteria
- [x] cmd.Debugf call sites changed to use correct Debugf
- [x] All re-exported symbols in errors.go cleaned up
- [x] All re-exported symbols in output.go cleaned up (except Debugf inline)
- [x] go build ./... passes with zero errors
- [x] go test ./... passes
- [x] go vet ./... passes

## Notes
Debugf kept as inline copy in output.go because base.Debugf does not correctly expand variadic args. The re-export layer for all other symbols has been fully removed.
