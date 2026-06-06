---
status: "completed"
started: "2026-06-06 13:31"
completed: "2026-06-06 13:35"
time_spent: "~4m"
---

# Task Record: 7 Update gen-journeys and gen-contracts references

## Summary
Updated gen-journeys SKILL.md, journey.md template, and gen-contracts rules/journey-contract-model.md to use surface-key adaptive directory rules instead of flat tests/<journey>/ structure

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/gen-journeys/SKILL.md
- plugins/forge/skills/gen-journeys/templates/journey.md
- plugins/forge/skills/gen-contracts/rules/journey-contract-model.md

### Key Decisions
无

## Document Metrics
3 files modified, 4 AC items validated, coverage: complete

## Referenced Documents
- docs/proposals/surface-key-test-alignment/proposal.md

## Review Status
final

## Acceptance Criteria
- [x] gen-journeys SKILL.md and journey.md template do not use surface-type as directory partition
- [x] gen-contracts SKILL.md tests/<journey>/ references updated to adaptive rules
- [x] journey-contract-model.md contract-to-directory mapping reflects surface-key structure

## Notes
gen-contracts SKILL.md itself has no tests/ path references (output is docs/features/.../contracts/). All tests/<journey>/ references were in journey-contract-model.md which is loaded by gen-contracts. Added surface_keys field to journey frontmatter alongside existing surface_types. Directory Convention, Rules, Migration Guide New Structure, and Tag-Based Promotion table all updated to single-surface flat / multi-surface partitioned by surfaceKey.
