---
status: "completed"
started: "2026-06-02 23:54"
completed: "2026-06-02 23:57"
time_spent: "~3m"
---

# Task Record: 18 Reduce cross-skill redundancy in surface detection + orchestration

## Summary
Simplified cross-skill redundancy in three locations: gen-contracts SKILL.md Surface Detection (25 lines -> 7 lines with HARD-RULE + reference to gen-journeys), run-tests SKILL.md Step 4 orchestration (removed duplicated failure handling details, now references surface rule files), journey-contract-model.md Semantic Descriptors (17 lines -> 5 lines with reference to dimension-rules.md).

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/gen-contracts/SKILL.md
- plugins/forge/skills/run-tests/SKILL.md
- plugins/forge/skills/gen-contracts/rules/journey-contract-model.md

### Key Decisions
无

## Document Metrics
3 files modified, ~50 lines removed, ~15 lines added, all HARD-RULEs preserved

## Referenced Documents
- docs/proposals/surface-first-testing/proposal.md
- plugins/forge/skills/gen-journeys/SKILL.md
- plugins/forge/skills/gen-contracts/rules/dimension-rules.md
- plugins/forge/skills/run-tests/rules/surfaces/web.md
- plugins/forge/skills/run-tests/rules/surfaces/cli.md
- docs/conventions/forge-distribution.md

## Review Status
final

## Acceptance Criteria
- [x] gen-contracts SKILL.md Surface Detection section <= 10 lines (HARD-RULE + reference)
- [x] run-tests SKILL.md Step 4 orchestration does not duplicate surface rule failure handling details
- [x] journey-contract-model.md Semantic Descriptors section simplified to reference dimension-rules.md
- [x] All simplified locations still retain key constraints (HARD-RULE), no information loss

## Notes
Reference format follows task Hard Rules: 'See `<relative-path>` for details'. gen-contracts references gen-journeys SKILL.md for full Surface Detection flow. run-tests references loaded surface rule files for failure handling. journey-contract-model references dimension-rules.md for full Semantic Descriptors rules.
