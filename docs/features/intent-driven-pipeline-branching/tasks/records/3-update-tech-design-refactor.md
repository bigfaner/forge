---
status: "completed"
started: "2026-05-29 17:11"
completed: "2026-05-29 17:15"
time_spent: "~4m"
---

# Task Record: 3 Update tech-design SKILL.md for refactor intent branch

## Summary
Updated tech-design SKILL.md with refactor/cleanup intent branching: added Intent Detection section, conditional logic in Steps 1/3/5/6/8/9 to skip API handbook, ER diagram, and prd-user-stories.md for refactor/cleanup intents while preserving full new-feature pipeline unchanged.

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/tech-design/SKILL.md

### Key Decisions
无

## Document Metrics
~200 lines added, covering intent detection table, 2 process flows, 6 conditional branches across Steps 1/3/5/6/8/9

## Referenced Documents
- docs/proposals/intent-driven-pipeline-branching/proposal.md
- plugins/forge/skills/tech-design/rules/design-quality-checks.md
- plugins/forge/skills/tech-design/templates/tech-design.md

## Review Status
final

## Acceptance Criteria
- [x] tech-design SKILL.md contains intent detection logic: when proposal.md frontmatter intent is refactor, execute internal-architecture-focused branch
- [x] refactor branch does not generate API handbook file and ER diagram file
- [x] refactor branch does not generate prd-user-stories.md file

## Notes
Implementation follows same pattern as write-prd SKILL.md which already has refactor intent branching. new-feature behavior is completely unchanged — all new content is in conditional blocks.
