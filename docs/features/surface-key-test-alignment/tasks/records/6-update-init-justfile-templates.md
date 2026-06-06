---
status: "completed"
started: "2026-06-06 13:49"
completed: "2026-06-06 13:53"
time_spent: "~4m"
---

# Task Record: 6 Update init-justfile justfile templates

## Summary
Added multi-surface test directory path convention comments to all 6 justfile templates, documenting the single-surface (tests/<journey>/) and multi-surface (tests/<surfaceKey>/<journey>/) path rules with language-specific command examples.

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/init-justfile/templates/generic.just
- plugins/forge/skills/init-justfile/templates/go.just
- plugins/forge/skills/init-justfile/templates/node.just
- plugins/forge/skills/init-justfile/templates/python.just
- plugins/forge/skills/init-justfile/templates/rust.just
- plugins/forge/skills/init-justfile/templates/mixed.just

### Key Decisions
无

## Document Metrics
6 templates updated with path convention comments; coverage: 100% of affected files

## Referenced Documents
- docs/proposals/surface-key-test-alignment/proposal.md

## Review Status
final

## Acceptance Criteria
- [x] All 6 templates document multi-surface path: tests/<surfaceKey>/<journey>/
- [x] Single surface default preserved: tests/<journey>/
- [x] Recipe parameter passing consistent with new directory structure

## Notes
Templates use just syntax which does not support conditional logic. The test recipe body remains the single-surface default; LLM customization during init-justfile replaces paths based on actual surface count. This approach is consistent with SKILL.md Step 3b which documents the adaptive path rule.
