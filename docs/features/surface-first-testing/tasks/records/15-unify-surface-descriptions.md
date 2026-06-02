---
status: "completed"
started: "2026-06-02 23:41"
completed: "2026-06-02 23:44"
time_spent: "~3m"
---

# Task Record: 15 Unify gen-journeys test type format + fix gen-contracts stale paths

## Summary
Unified gen-journeys 5 surface rule files test type format to single-line English description, fixed stale path in journey-contract-model.md (testing-<scope>.md -> testing/<surface>/core.md), and removed reference to non-existent model-and-directory-spec.md in gen-journeys SKILL.md

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/gen-journeys/rules/surface-cli.md
- plugins/forge/skills/gen-journeys/rules/surface-api.md
- plugins/forge/skills/gen-journeys/rules/surface-tui.md
- plugins/forge/skills/gen-journeys/rules/surface-web.md
- plugins/forge/skills/gen-journeys/rules/surface-mobile.md
- plugins/forge/skills/gen-contracts/rules/journey-contract-model.md
- plugins/forge/skills/gen-journeys/SKILL.md

### Key Decisions
无

## Document Metrics
7 files modified, 5 format unifications, 2 stale path fixes

## Referenced Documents
- docs/proposals/surface-first-testing/proposal.md
- plugins/forge/skills/test-guide/references/test-type-model.md

## Review Status
final

## Acceptance Criteria
- [x] gen-journeys 5 surface rule files test type unified to **Test type**: {English name}. {one-sentence English description}.
- [x] Removed duplicate Chinese test type descriptions from line 5 (Chinese paragraph on line 3 preserved)
- [x] journey-contract-model.md testing-<scope>.md changed to testing/<surface>/core.md
- [x] gen-journeys SKILL.md HARD-RULE no longer references non-existent model-and-directory-spec.md

## Notes
Test type English names verified against test-type-model.md mapping table: CLI Functional Test, API Functional Test, Terminal Functional Test, Web E2E Test, Mobile E2E Test. SKILL.md Reference section also updated to point to journey-contract-model.md instead of model-and-directory-spec.md.
