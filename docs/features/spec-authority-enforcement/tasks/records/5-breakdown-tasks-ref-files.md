---
status: "completed"
started: "2026-05-24 09:53"
completed: "2026-05-24 09:54"
time_spent: "~1m"
---

# Task Record: 5 Improve breakdown-tasks SKILL.md Reference Files generation for non-UI tasks

## Summary
Added Reference Files Population guidance for non-UI tasks in breakdown-tasks SKILL.md Step 4a

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/breakdown-tasks/SKILL.md

### Key Decisions
无

## Document Metrics
1 file modified, ~25 lines added to Step 4a

## Referenced Documents
- docs/proposals/spec-authority-enforcement/proposal.md
- docs/conventions/forge-distribution.md
- plugins/forge/skills/breakdown-tasks/rules/ui-placement.md

## Review Status
final

## Acceptance Criteria
- [x] SKILL.md Step 4a includes explicit instructions for populating Reference Files with section-level precision
- [x] Instructions require format: path/to/tech-design.md#Section-Title — brief description
- [x] Instructions specify 4-step extraction heuristic (extract paths, search sections, find arch decisions, merge 2-5)
- [x] Instructions include checklist: every generated task must have >=1 design-level Reference File
- [x] Guidance described as heuristic strategy, not deterministic algorithm
- [x] UI tasks continue to use existing rules/ui-placement.md requirements — no conflict

## Notes
Inserted as #### subsection within Step 4a Business Tasks, after Hard Rules and before Scope Assignment. No structural changes to SKILL.md.
