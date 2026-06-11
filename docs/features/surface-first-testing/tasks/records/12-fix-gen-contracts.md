---
status: "completed"
started: "2026-06-02 22:55"
completed: "2026-06-02 22:57"
time_spent: "~2m"
---

# Task Record: 12 Fix gen-contracts: surface detection, Convention loading, surfaceType naming

## Summary
Fixed gen-contracts: surface detection now uses forge surfaces CLI, Convention loading switched to surface-first path (testing/{surface}/core.md), risk-density.md WebUI→Web, journey-contract-model.md UI→Web, removed cross-skill internal file references

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/gen-contracts/SKILL.md
- plugins/forge/skills/gen-contracts/rules/risk-density.md
- plugins/forge/skills/gen-contracts/rules/journey-contract-model.md

### Key Decisions
无

## Document Metrics
3 files modified, 5 AC items met, all surfaceType naming unified to web/api/cli/tui/mobile

## Referenced Documents
- docs/proposals/surface-first-testing/proposal.md
- docs/conventions/forge-distribution.md
- plugins/forge/skills/gen-journeys/SKILL.md
- plugins/forge/skills/gen-test-scripts/SKILL.md

## Review Status
final

## Acceptance Criteria
- [x] gen-contracts SKILL.md surface detection uses forge surfaces CLI
- [x] gen-contracts SKILL.md Convention loading path changed to docs/conventions/testing/{surface}/core.md with legacy fallback
- [x] risk-density.md WebUI changed to Web
- [x] journey-contract-model.md UI changed to Web in type list
- [x] All surfaceType unified to web/api/cli/tui/mobile

## Notes
Also removed cross-skill internal file references: SKILL.md Step 1 no longer references gen-journeys rules/surface-<type>.md, risk-density.md Surface-Required Outcome Derivation table now inline instead of referencing gen-journeys skill files
