---
status: "completed"
started: "2026-06-02 23:38"
completed: "2026-06-02 23:40"
time_spent: "~2m"
---

# Task Record: 14 Fix test-guide test directory paths in SKILL.md + surface templates

## Summary
Unified test directory paths in SKILL.md quick-reference table and CLI/API/TUI surface templates to tests/journey-based paths. Also fixed residual references in rules/convention-structure.md and rules/draft-generation.md.

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/test-guide/SKILL.md
- plugins/forge/skills/test-guide/templates/surfaces/cli.md
- plugins/forge/skills/test-guide/templates/surfaces/api.md
- plugins/forge/skills/test-guide/templates/surfaces/tui.md
- plugins/forge/skills/test-guide/rules/convention-structure.md
- plugins/forge/skills/test-guide/rules/draft-generation.md

### Key Decisions
无

## Document Metrics
6 files modified, 0 residual stale path references

## Referenced Documents
- docs/proposals/surface-first-testing/proposal.md

## Review Status
final

## Acceptance Criteria
- [x] SKILL.md quick-reference table: all 5 surface rows unified to tests/journey-based paths
- [x] CLI template: removed surface-based fallback, unified to journey-based paths
- [x] API template: removed surface-based fallback, unified to journey-based paths
- [x] TUI template: removed surface-based fallback, unified to journey-based paths
- [x] No residual stale path references in test-guide skill

## Notes
Extra fix: rules/convention-structure.md and rules/draft-generation.md also contained the same stale quick-reference table paths, fixed for consistency.
