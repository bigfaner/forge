---
status: "completed"
started: "2026-06-08 17:58"
completed: "2026-06-08 18:01"
time_spent: "~3m"
---

# Task Record: 1 gen-journeys 新增 Golden Path Journey 强制要求

## Summary
Added Golden Path Journey mandatory requirement to gen-journeys skill: created rules/golden-path.md with dual constraints (3+ steps + semantic completeness), feature complexity classification heuristics, and semantic completeness proxy; modified SKILL.md to reference the rule; updated journey template with golden_path frontmatter field.

## Changes

### Files Created
- plugins/forge/skills/gen-journeys/rules/golden-path.md

### Files Modified
- plugins/forge/skills/gen-journeys/SKILL.md
- plugins/forge/skills/gen-journeys/templates/journey.md

### Key Decisions
无

## Document Metrics
1 rule file (~80 lines), 2 modified files (~15 lines added to SKILL.md, 1 line added to template)

## Referenced Documents
- docs/proposals/behavioral-test-accuracy/proposal.md

## Review Status
final

## Acceptance Criteria
- [x] gen-journeys SKILL.md references rules/golden-path.md and declares Golden Path requirement
- [x] rules/golden-path.md contains Golden Path dual constraint rules: (a) 3+ steps, (b) steps extracted from PRD/Design primary user story
- [x] rules/golden-path.md contains Feature complexity classification heuristics table (Simple vs Complex)
- [x] rules/golden-path.md declares semantic completeness proxy: domain terminology required, API terminology prohibited
- [x] Journey template supports Golden Path marker (golden_path frontmatter field) for downstream gen-contracts

## Notes
Hard Rules enforced: Golden Path applies to all surface types; entity relationship classification takes priority over workflow description. Proposal.md was the authoritative source for rule content.
