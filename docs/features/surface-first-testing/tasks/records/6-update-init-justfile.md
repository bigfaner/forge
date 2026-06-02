---
status: "completed"
started: "2026-06-02 22:03"
completed: "2026-06-02 22:05"
time_spent: "~2m"
---

# Task Record: 6 更新 init-justfile Test recipe

## Summary
Updated init-justfile SKILL.md Step 0 to load Convention from surface-first directory structure (testing/{surface}/core.md), added legacy fallback for old flat-file structure, updated all Convention path references in output examples and Notes section

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/init-justfile/SKILL.md

### Key Decisions
无

## Document Metrics
7 locations updated: Step 0 convention loading rewrite, HARD-RULE section names, cold-start hint, Step 3a convention source, 2x output examples convention path, Notes test-type-model path

## Referenced Documents
- docs/proposals/surface-first-testing/proposal.md
- docs/conventions/forge-distribution.md

## Review Status
final

## Acceptance Criteria
- [x] SKILL.md Test recipe loads from new directory structure testing/{surface}/
- [x] Test recipe naming follows Surface type convention (e.g. cli-functional build tags)

## Notes
Added legacy fallback for old framework-first flat files to ensure backward compatibility. Build tag example updated from generic <surface>-<type> to concrete cli-functional.
