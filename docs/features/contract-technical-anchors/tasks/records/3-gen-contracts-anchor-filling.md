---
status: "completed"
started: "2026-06-05 17:29"
completed: "2026-06-05 17:35"
time_spent: "~6m"
---

# Task Record: 3 gen-contracts 从 handbook 填充锚点字段

## Summary
Added handbook loading and anchor filling logic to gen-contracts SKILL.md. New Step 3 (Load Handbooks + Anchor Filling) covers handbook location, freshness check, missing handbook graceful degradation, anchor extraction per surface type, and anchor sync timestamp. Updated process flow, step numbering, pipeline position table, core principle description, schema validation checks, and key concepts.

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/gen-contracts/SKILL.md

### Key Decisions
无

## Document Metrics
~70 lines added (Step 3: Load Handbooks + Anchor Filling), 2 validation checks added, 3 key concepts added

## Referenced Documents
- docs/proposals/contract-technical-anchors/proposal.md
- plugins/forge/skills/gen-contracts/templates/contract.md

## Review Status
final

## Acceptance Criteria
- [x] gen-contracts reads api-handbook to auto-fill API Contract endpoint/method fields
- [x] gen-contracts reads cli-handbook to auto-fill CLI/TUI Contract command/subcommand fields
- [x] gen-contracts reads page-map/screen-map to auto-fill Web/Mobile Contract page/screen fields
- [x] Handbook freshness check: warn when handbook generated before tech-design last modified
- [x] Missing handbook: skip anchor filling and prompt user to run tech-design
- [x] last_anchor_sync timestamp auto-updated on fill

## Notes
Backward compatibility ensured via HARD-RULE: missing handbooks do not cause pipeline failure. Anchors come from handbooks as authority source, not reverse-engineered from code.
