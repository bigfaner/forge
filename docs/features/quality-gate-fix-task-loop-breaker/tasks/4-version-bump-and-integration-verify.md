---
id: "4"
title: "Version bump and integration verification"
priority: "P2"
estimated_time: "30m"
dependencies: ["1", "2", "3"]
scope: "backend"
breaking: false
type: "implementation"
mainSession: false
---

# 4: Version bump and integration verification

## Description

Patch version bump and final integration check — verify all fixes work together, no regressions in non-quality-gate fix task flows.

## Reference Files
- `docs/proposals/quality-gate-fix-task-loop-breaker/proposal.md` — Source proposal
- `scripts/version.txt` — Version file
- `forge-cli/internal/cmd/quality_gate.go` — All modified files
- `forge-cli/internal/cmd/quality_gate_test.go` — All modified tests

## Acceptance Criteria
- [ ] Version bumped in `scripts/version.txt` (patch)
- [ ] All quality-gate tests pass (`go test ./forge-cli/internal/cmd/ -run QualityGate -v`)
- [ ] No regression in non-quality-gate fix task flows (dispatcher-created fix tasks still work)
- [ ] Full test suite passes (`just test`)

## Hard Rules
- Only patch version bump (e.g., 3.0.0 → 3.0.1 or next patch number based on current version).

## Implementation Notes
- Run `go test ./forge-cli/... -v` for full CLI test coverage.
- Verify dispatcher-created fix tasks (from `execute-task` quality gate) are unaffected — they use different SourceTaskID values.
