---
status: "completed"
started: "2026-05-30 22:11"
completed: "2026-05-30 22:15"
time_spent: "~4m"
---

# Task Record: 4 Write constants.md and extend enum-constants.md

## Summary
Created docs/conventions/constants.md with magic value classification rules (paths, colors, timeouts, sentinel values, permissions), extraction rules, centralization table, target state definition, and deviation analysis. Extended docs/conventions/enum-constants.md with 5 new TECH-const rules (TECH-const-001 through TECH-const-005) covering path constants, timeout/duration values, color values, sentinel values, and permission values -- each with deviation analysis referencing specific Evidence cases.

## Changes

### Files Created
- docs/conventions/constants.md

### Files Modified
- docs/conventions/enum-constants.md

### Key Decisions
无

## Document Metrics
constants.md: ~200 lines, 5 classification categories, 20 deviation entries; enum-constants.md: +130 lines, 5 new TECH-const rules

## Referenced Documents
- forge-cli/internal/cmd/quality_gate.go
- forge-cli/internal/cmd/init.go
- forge-cli/internal/cmd/task/list.go
- forge-cli/internal/cmd/task/claim.go
- forge-cli/internal/cmd/init_surfaces.go
- forge-cli/pkg/serverprobe/serverprobe.go
- forge-cli/pkg/index/lock.go
- forge-cli/pkg/feature/constants.go
- forge-cli/pkg/testrunner/test_results.go

## Review Status
final

## Acceptance Criteria
- [x] constants.md exists covering classification rules (paths, colors, timeouts, sentinel values, permissions) and extraction rules (when to extract, centralized management location)
- [x] enum-constants.md extended with non-enum constant management rules for path constants, timeout values, color values
- [x] Both files include target state definitions and deviation analysis referencing Evidence-specific magic value cases
- [x] Constant centralization location specified (constants.go per package)

## Notes
Deviation analysis covers all Evidence cases: quality_gate.go path strings (P1/P2), init.go color #7DCFFF (C1/C2), init_surfaces.go colors (C3/C4), list.go ANSI codes (C5), list.go sentinel 99999 (S1), claim.go sentinel 99999 (S2), quality_gate.go retry params (R1/R3/R4), serverprobe.go timeout (T2), lock.go backoff (T4). Additional audit found: tree.go statusColor function is already centralized and acceptable.
