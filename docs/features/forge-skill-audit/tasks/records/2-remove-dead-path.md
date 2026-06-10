---
status: "completed"
started: "2026-06-10 19:18"
completed: "2026-06-10 19:20"
time_spent: "~2m"
---

# Task Record: 2 Remove tech-design dead path (H-2)

## Summary
Removed dead path docs/features/<slug>/proposal.md from tech-design SKILL.md Intent Detection section; verified no other skill references this dead path

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/tech-design/SKILL.md

### Key Decisions
无

## Document Metrics
1 file modified, 1 line changed (dead path removed), 0 other dead paths found across all skills

## Referenced Documents
- docs/proposals/forge-skill-audit/proposal.md
- plugins/forge/skills/tech-design/SKILL.md
- docs/conventions/forge-distribution.md

## Review Status
final

## Acceptance Criteria
- [x] tech-design SKILL.md only references docs/proposals/<slug>/proposal.md, no docs/features/<slug>/proposal.md path
- [x] Search all skills confirms no other skill references docs/features/<slug>/proposal.md dead path

## Notes
docs/features/ path is still used correctly by tech-design for other files (e.g. prd/prd-spec.md prerequisite check). Only the proposal.md reference was a dead path. eval/SKILL.md and git-checkout.md reference docs/features/<slug>/ generically for non-proposal files, which is correct.
