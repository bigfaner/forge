---
status: "completed"
started: "2026-05-31 14:52"
completed: "2026-05-31 14:53"
time_spent: "~1m"
---

# Task Record: 1 Update brainstorm intent mapping to 6-value enum

## Summary
Updated brainstorm SKILL.md Step 4.5 intent mapping table from 3 values to 6 values (new-feature, enhancement, refactor, cleanup, fix, doc). Removed coding.fix heuristic logic entirely — coding.fix now maps directly to fix. Updated AskUserQuestion to offer all 6 intent options. Added valid values comment to brainstorm/templates/proposal.md intent frontmatter listing all 6 values. Split coding.feature and coding.enhancement into independent mapping paths.

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/brainstorm/SKILL.md
- plugins/forge/skills/brainstorm/templates/proposal.md

### Key Decisions
无

## Document Metrics
2 files modified, 6-value enum fully applied

## Referenced Documents
- docs/proposals/intent-enriched-enum/proposal.md

## Review Status
final

## Acceptance Criteria
- [x] brainstorm/SKILL.md Step 4.5 intent mapping table contains exactly 6 values: new-feature, enhancement, refactor, cleanup, fix, doc
- [x] Fix heuristic logic removed entirely — coding.fix always maps to fix intent without runtime inference
- [x] AskUserQuestion intent selection offers all 6 values as structured options
- [x] brainstorm/templates/proposal.md intent valid values comment lists all 6 values
- [x] coding.feature -> new-feature and coding.enhancement -> enhancement exist as independent mapping paths (not merged)

## Notes
Followed proposal.md Proposed Solution item 1 (mapping), item 3 (remove heuristic), and Architecture Decision (enhancement split from new-feature). doc.consolidate/doc.drift umbrella note added per proposal.
