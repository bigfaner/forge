---
status: "completed"
started: "2026-06-02 23:33"
completed: "2026-06-02 23:37"
time_spent: "~4m"
---

# Task Record: 13 Rewrite test-guide draft-generation.md for surface-first + remove orphan template

## Summary
Rewrote draft-generation.md from 4-section framework-first schema to 7+1 surface-first schema aligned with convention-structure.md and SKILL.md. Removed orphan convention-template.md (old 4-section format, unreferenced by any file).

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/test-guide/rules/draft-generation.md

### Key Decisions
无

## Document Metrics
draft-generation.md: ~220 lines, 7+1 section schema, 6 AC all pass

## Referenced Documents
- docs/proposals/surface-first-testing/proposal.md
- plugins/forge/skills/test-guide/rules/convention-structure.md
- plugins/forge/skills/test-guide/SKILL.md
- docs/conventions/forge-distribution.md

## Review Status
final

## Acceptance Criteria
- [x] draft-generation.md section schema matches convention-structure.md 7+1 sections
- [x] Template paths changed from docs/conventions/testing/<framework>.md to docs/conventions/testing/<surface>/core.md
- [x] File naming from <scope>.md changed to surface-first directory structure
- [x] Built-in Template lookup table updated to templates/surfaces/<surface>.md
- [x] templates/convention-template.md deleted
- [x] draft-generation.md consistent with convention-structure.md and SKILL.md generation flow

## Notes
Old 4-section schema (framework, discovery, structure, assertions) fully replaced. Inferred Defaults table preserved as auxiliary reference. convention-template.md confirmed orphan (0 references across test-guide skill).
