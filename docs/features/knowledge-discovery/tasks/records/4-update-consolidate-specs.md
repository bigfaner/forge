---
status: "completed"
started: "2026-05-17 13:44"
completed: "2026-05-17 13:46"
time_spent: "~2m"
---

# Task Record: 4 Update consolidate-specs to manage domains frontmatter

## Summary
Updated consolidate-specs SKILL.md to generate and maintain domains frontmatter on convention/business-rule files. Added Domain Frontmatter section with derivation rules (spec ID keywords + source keywords, 3-7 per file), domain overlap detection (>50% triggers warning), new file frontmatter template, and drift re-derivation logic in Step 10.

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/consolidate-specs/SKILL.md

### Key Decisions
- Domains derived programmatically from spec ID tokens and source keywords rather than freeform agent invention
- Domain overlap uses intersection/min ratio with >50% threshold flagged during user confirmation
- Existing title frontmatter behavior preserved unchanged; domains is purely additive

## Test Results
- **Tests Executed**: Yes
- **Passed**: 0
- **Failed**: 0
- **Coverage**: 0.0%

## Acceptance Criteria
- [x] SKILL.md instructs agent to write domains frontmatter when creating new files in docs/conventions/ or docs/business-rules/
- [x] Domains are derived from spec content (ID keywords, source keywords) -- not freeform
- [x] During drift detection (Steps 9-11), domains are re-derived when file content changes substantially
- [x] Domain overlap >50% between files triggers a warning during the user confirmation step
- [x] Each file gets 3-7 specific keywords
- [x] The existing title frontmatter behavior is unchanged

## Notes
Documentation-only task. No test changes. Added new Domain Frontmatter section before Step 5, updated Steps 6, 7, 10, and Rules section.
