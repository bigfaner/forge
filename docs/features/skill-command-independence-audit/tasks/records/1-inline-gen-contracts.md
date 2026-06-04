---
status: "completed"
started: "2026-06-04 00:33"
completed: "2026-06-04 00:35"
time_spent: "~2m"
---

# Task Record: 1 Inline surface detection into gen-contracts + clean Related sections

## Summary
Inlined Surface Detection rules from gen-journeys into gen-contracts, deleted Related Skills and References sections, merged 6 concept definitions into inline Key Concepts section

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/gen-contracts/SKILL.md

### Key Decisions
无

## Document Metrics
inlined ~30 lines of Surface Detection content, removed 18 lines (Related Skills + Reference sections), net +12 lines; HARD block count unchanged at 9

## Referenced Documents
- docs/proposals/skill-command-independence-audit/proposal.md
- plugins/forge/skills/gen-journeys/SKILL.md

## Review Status
final

## Acceptance Criteria
- [x] gen-contracts/SKILL.md contains inlined Surface Detection rules (~20 lines) with INLINE:origin marker
- [x] Related Skills and Integration sections deleted
- [x] References concept definitions (6 concepts) merged into inline Key Concepts section
- [x] gen-contracts no longer contains cross-skill references to gen-journeys internal files

## Notes
Surface Detection HARD-RULE retained gen-contracts-specific wording ('before Convention loading begins') rather than gen-journeys wording ('before Journey generation begins'). INLINE markers added for traceability per Hard Rules.
