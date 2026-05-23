---
id: "12"
title: "Extract runE2ERegression to reduce nesting"
priority: "P2"
estimated_time: "1h"
dependencies: [9]
scope: "backend"
breaking: false
type: "coding.refactor"
mainSession: false
---

# 12: Extract runE2ERegression to reduce nesting

## Description
The e2e regression logic in `runQualityGate()` has 4 levels of nesting. Extract into a standalone `runE2ERegression()` function to flatten the control flow. Phase 4 (anti-pattern fix).

## Reference Files
- `docs/proposals/forge-cli-clean-code/proposal.md` — Source proposal
- `forge-cli/internal/cmd/quality_gate.go` — Contains the nested logic

## Acceptance Criteria
- [ ] `runE2ERegression()` extracted as a standalone function
- [ ] Nesting in `runQualityGate()` reduced by at least 2 levels
- [ ] `go build ./...` passes
- [ ] `go test ./...` passes

## Hard Rules
- Follow the project's "prefer flat control flow" convention: use early returns and guard clauses

## Implementation Notes
- The 4-level nesting suggests: if err != nil → if condition → for loop → if condition
- Extraction should preserve the exact error handling and flow control semantics
