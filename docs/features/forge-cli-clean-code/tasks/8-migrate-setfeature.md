---
id: "8"
title: "Complete SetFeature migration and remove deprecated function"
priority: "P1"
estimated_time: "1h"
dependencies: [4]
scope: "backend"
breaking: true
type: "coding.cleanup"
mainSession: false
---

# 8: Complete SetFeature migration and remove deprecated function

## Description
`SetFeature()` in `pkg/feature/feature.go` is marked Deprecated but still has 7+ call sites. Migrate all callers to the replacement API, then remove `SetFeature()`. Phase 3 (duplicate logic consolidation).

## Reference Files
- `docs/proposals/forge-cli-clean-code/proposal.md` — Source proposal
- `forge-cli/pkg/feature/feature.go` — Contains deprecated `SetFeature()`

## Acceptance Criteria
- [ ] All 7+ call sites migrated from `SetFeature()` to the replacement API
- [ ] `SetFeature()` function deleted
- [ ] 0 Deprecated call sites remain
- [ ] `go build ./...` passes
- [ ] `go test ./...` passes

## Hard Rules
- Identify the replacement API before starting migration
- Migrate one call site at a time, running tests after each

## Implementation Notes
- Deprecated for over one version — the replacement API should be well-established
- Use `grep -r "SetFeature"` to find all call sites
- Risk: some callers may use `SetFeature()` in ways the replacement doesn't support — verify each call site individually
