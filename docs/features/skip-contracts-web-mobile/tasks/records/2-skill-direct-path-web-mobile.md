---
status: "completed"
started: "2026-06-09 18:24"
completed: "2026-06-09 18:29"
time_spent: "~5m"
---

# Task Record: 2 Add direct test generation path for web/mobile journeys in gen-test-scripts

## Summary
Added direct test generation path for web/mobile journeys in gen-test-scripts: SKILL.md Step 2 routing (Contract Path vs Direct Path), coverage self-check (Step 5), and Direct Path Generation Rules in types/web.md and types/mobile.md

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/gen-test-scripts/SKILL.md
- plugins/forge/skills/gen-test-scripts/types/web.md
- plugins/forge/skills/gen-test-scripts/types/mobile.md

### Key Decisions
无

## Document Metrics
SKILL.md: +~180 lines (Step 2.0 routing, Step 2.2 direct path, Step 5 coverage self-check); types/web.md: +~80 lines (direct path rules); types/mobile.md: +~80 lines (direct path rules)

## Referenced Documents
- docs/proposals/skip-contracts-web-mobile/proposal.md
- docs/conventions/forge-distribution.md

## Review Status
final

## Acceptance Criteria
- [x] SC-5: Direct path generated scripts include user action calls (click/type/navigate) and at least one non-trivial visual assertion
- [x] SC-6: Coverage self-check by surface type with gap list and hard FAIL on mismatches
- [x] SC-8: types/web.md and types/mobile.md direct path rules produce tests with meaningful assertions

## Notes
Three files modified. SKILL.md Step 2 restructured into Step 2.0 (route decision), Step 2.1 (contract path), Step 2.2 (direct path). Step 2.5.1 updated to support surface_types-based type extraction for direct path. New Step 5 adds coverage self-check with Surface->Test Type mapping. Both type files received comprehensive Direct Path Generation Rules sections with step-to-action mapping templates, precondition-to-fixture mapping, expected-result-to-assertion mapping, and minimum assertion requirements.
