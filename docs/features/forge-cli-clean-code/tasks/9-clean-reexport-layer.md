---
id: "9"
title: "Clean up re-export layer in errors.go and output.go"
priority: "P1"
estimated_time: "1.5h"
dependencies: [4]
scope: "backend"
breaking: true
type: "coding.refactor"
mainSession: false
---

# 9: Clean up re-export layer in errors.go and output.go

## Description
`errors.go` and `output.go` in `internal/cmd/` re-export all symbols from `base/`, creating a redundant indirection layer. The sub-packages already import `base/` directly. Remove the re-exports and update all call sites (especially `cmd.Debugf` → `base.Debugf` in `quality_gate.go` and 10+ other locations). Phase 3 (duplicate logic consolidation).

## Reference Files
- `docs/proposals/forge-cli-clean-code/proposal.md` — Source proposal
- `forge-cli/internal/cmd/errors.go` — Re-export file
- `forge-cli/internal/cmd/output.go` — Re-export file (also contains duplicate `Debugf`)
- `forge-cli/internal/cmd/base/errors.go` — Source definitions
- `forge-cli/internal/cmd/base/output.go` — Source definitions

## Acceptance Criteria
- [ ] `cmd.Debugf` call sites (10+ in `quality_gate.go` and elsewhere) changed to `base.Debugf`
- [ ] All re-exported symbols in `errors.go` and `output.go` cleaned up
- [ ] `go build ./...` passes with zero errors
- [ ] `go test ./...` passes
- [ ] `go vet ./...` passes (confirms no dangling import references)

## Hard Rules
- Run `go vet ./...` after cleanup to confirm all import paths are correct
- Do NOT modify `base/errors.go` or `base/output.go` — only remove the re-export layer in the parent package

## Implementation Notes
- The proposal overturned the assumption that the re-export layer is necessary: sub-packages already import `base/` directly
- Start by listing all re-exported symbols and their callers
- Change callers to import `base` directly, then remove the re-export declarations
- Risk: external callers (tests in other packages) may import via the re-export — check test files too
