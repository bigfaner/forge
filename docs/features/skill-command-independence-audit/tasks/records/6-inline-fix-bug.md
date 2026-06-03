---
status: "completed"
started: "2026-06-04 00:47"
completed: "2026-06-04 00:50"
time_spent: "~3m"
---

# Task Record: 6 Inline templates/rules into fix-bug command + reduce Knowledge Review

## Summary
Inlined learn/templates write formats and consolidate-specs/rules classification rules into fix-bug.md; compressed Knowledge Review from ~165 lines to ~139 lines by converting Notable Knowledge Heuristics to a compact table and streamlining Extraction Flow steps

## Changes

### Files Created
无

### Files Modified
- plugins/forge/commands/fix-bug.md

### Key Decisions
无

## Document Metrics
inline sections: ~48 lines (Write Formats ~28 + Classification ~20); Knowledge Review reduced by ~16%; 0 cross-references remaining

## Referenced Documents
- docs/proposals/skill-command-independence-audit/proposal.md
- plugins/forge/skills/learn/templates/decision-entry.md
- plugins/forge/skills/learn/templates/lesson-entry.md
- plugins/forge/skills/learn/templates/convention-entry.md
- plugins/forge/skills/consolidate-specs/rules/spec-classification.md
- plugins/forge/skills/consolidate-specs/rules/overlap-detection.md
- plugins/forge/skills/consolidate-specs/rules/domain-frontmatter.md

## Review Status
final

## Acceptance Criteria
- [x] fix-bug.md contains inlined template decision points and spec extraction rules (~40 lines) with <!-- INLINE:origin=... --> markers
- [x] Knowledge Review section has been slimmed down
- [x] fix-bug no longer contains cross-references to learn or consolidate-specs internal files

## Notes
Two INLINE markers added per Hard Rules: <!-- INLINE:origin=learn/templates/ --> at line 246, <!-- INLINE:origin=consolidate-specs/rules/ --> at line 276. Heuristics converted from 8 separate blocks (~55 lines) to a single compact table (~10 lines). Extraction Flow steps renumbered from #### headers to numbered list (1-7).
