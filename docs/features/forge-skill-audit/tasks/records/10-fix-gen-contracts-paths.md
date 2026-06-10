---
status: "completed"
started: "2026-06-10 21:02"
completed: "2026-06-10 21:04"
time_spent: "~2m"
---

# Task Record: 10 Fix gen-contracts path references (HIGH-A1, HIGH-A2, MEDIUM-A3)

## Summary
Fixed 3 path reference issues in gen-contracts/SKILL.md: (1) added design/ prefix to all 5 handbook paths in Handbook Location table, (2) fixed freshness check tech-design path to include design/, (3) clarified that surface-*.md files referenced in INLINE block belong to gen-journeys

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/gen-contracts/SKILL.md

### Key Decisions
无

## Document Metrics
3 fixes across 1 file, 5 handbook paths corrected, 1 freshness path corrected, 1 INLINE clarification added

## Referenced Documents
- docs/proposals/forge-skill-audit/proposal.md
- plugins/forge/skills/tech-design/SKILL.md

## Review Status
final

## Acceptance Criteria
- [x] Handbook Location table paths include design/ subdirectory prefix
- [x] Freshness check references docs/features/<slug>/design/tech-design.md
- [x] INLINE surface-*.md references annotated as belonging to gen-journeys

## Notes
Verified against tech-design/SKILL.md lines 270 and 296-300 confirming all handbook/tech-design output goes to design/ subdirectory
