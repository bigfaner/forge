---
status: "completed"
started: "2026-05-30 06:14"
completed: "2026-05-30 06:15"
time_spent: "~1m"
---

# Task Record: 20 Fix: align Convention loading in 4 SKILL.md files

## Summary
Aligned Convention loading in 3 SKILL.md Step 0 sections (breakdown-tasks, tech-design, quick-tasks) to use index.md-based discovery instead of domains frontmatter filtering, matching convention-guide.md HARD-RULE. gen-test-scripts was already aligned.

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/breakdown-tasks/SKILL.md
- plugins/forge/skills/tech-design/SKILL.md
- plugins/forge/skills/quick-tasks/SKILL.md

### Key Decisions
无

## Document Metrics
3 files modified, Step 0 sections only, ~3 lines changed per file

## Referenced Documents
- docs/features/plugin-consistency-audit/reports/06-consolidated-report.md
- plugins/forge/skills/gen-test-scripts/rules/convention-guide.md

## Review Status
final

## Acceptance Criteria
- [x] 4 SKILL.md Step 0 Convention loading aligns with convention-guide.md
- [x] If domains filtering is a justified design difference, document the reason

## Notes
All 4 skills now use index.md-based discovery. gen-test-scripts was already correct. The other 3 were updated from 'domains frontmatter filtering' to 'read index.md, select by project context'. No justified design difference for domains filtering exists — index.md-based discovery serves language detection equally well.
