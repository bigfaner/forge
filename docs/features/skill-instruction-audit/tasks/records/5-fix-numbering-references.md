---
status: "completed"
started: "2026-05-28 23:01"
completed: "2026-05-28 23:04"
time_spent: "~3m"
---

# Task Record: 5 Fix numbering/reference errors in skills

## Summary
Fixed numbering/reference errors across 6 skill SKILL.md files: tech-design flow gap, run-tests misleading cross-reference, write-prd decimal numbering, breakdown-tasks Step 4b converted to informational note, gen-contracts Section X.Y references removed.

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/tech-design/SKILL.md
- plugins/forge/skills/run-tests/SKILL.md
- plugins/forge/skills/write-prd/SKILL.md
- plugins/forge/skills/breakdown-tasks/SKILL.md
- plugins/forge/skills/gen-contracts/SKILL.md

### Key Decisions
无

## Document Metrics
5 files modified, 6 AC items addressed, 0 decimal steps remaining

## Referenced Documents
- docs/features/skill-instruction-audit/tasks/5-fix-numbering-references.md

## Review Status
final

## Acceptance Criteria
- [x] tech-design Process Flow: 0->1->...->8->9->10->11 no gaps
- [x] run-tests Step 5 does not reference 'Convention loaded in Step 0'
- [x] write-prd has no decimal step numbers
- [x] quick-tasks has no decimal step numbers
- [x] breakdown-tasks Step 4b is informational note, not numbered step
- [x] gen-contracts has no 'Section X.Y' references

## Notes
quick-tasks had no decimal step numbers to fix (AC already met). Only 5 of 6 files needed modification.
