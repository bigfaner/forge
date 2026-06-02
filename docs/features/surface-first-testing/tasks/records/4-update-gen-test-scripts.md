---
status: "completed"
started: "2026-06-02 21:58"
completed: "2026-06-02 22:01"
time_spent: "~3m"
---

# Task Record: 4 更新 gen-test-scripts Convention 加载路径

## Summary
Updated gen-test-scripts SKILL.md and convention-guide.md to use surface-first Convention loading: testing/{surface}/core.md instead of framework-first flat files. Added legacy structure detection with migration prompt. Preserved per-surface build tag naming (already correct in Test Type Terminology).

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/gen-test-scripts/SKILL.md
- plugins/forge/skills/gen-test-scripts/rules/convention-guide.md

### Key Decisions
无

## Document Metrics
2 files modified, Step 0 rewritten (4 subsections), convention-guide.md fully restructured to surface-first schema (7 sections)

## Referenced Documents
- docs/proposals/surface-first-testing/proposal.md
- docs/conventions/forge-distribution.md

## Review Status
final

## Acceptance Criteria
- [x] SKILL.md Convention loading path changed to testing/{surface}/core.md surface directory traversal
- [x] Legacy structure detection outputs migration prompt instead of silent failure
- [x] Generated test code uses per-surface build tag naming (e.g. cli_functional not e2e)

## Notes
SKILL.md Step 0 fully rewritten: 0.1 Old Structure Detection, 0.2 Surface-First Discovery, 0.3 Framework Resolution, 0.4 Validation with authority Hard Rule (types/*.md > core.md). convention-guide.md updated: new directory structure, 7-section schema, assertion preference table with fixed columns, legacy detection section.
