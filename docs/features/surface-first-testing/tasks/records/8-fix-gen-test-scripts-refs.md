---
status: "completed"
started: "2026-06-02 22:43"
completed: "2026-06-02 22:47"
time_spent: "~4m"
---

# Task Record: 8 Fix gen-test-scripts: stale references + rename ui.md → web.md

## Summary
Fixed gen-test-scripts skill: removed all stale docs/reference/test-type-model.md references with inline test type definitions, renamed types/ui.md to types/web.md via git mv, updated all UI surface type references to Web across SKILL.md, types/_shared.md, types/cli.md, types/api.md, types/tui.md, types/mobile.md, and rules/step-0.5-validation.md

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/gen-test-scripts/types/web.md
- plugins/forge/skills/gen-test-scripts/SKILL.md
- plugins/forge/skills/gen-test-scripts/types/_shared.md
- plugins/forge/skills/gen-test-scripts/types/cli.md
- plugins/forge/skills/gen-test-scripts/types/api.md
- plugins/forge/skills/gen-test-scripts/types/tui.md
- plugins/forge/skills/gen-test-scripts/types/mobile.md
- plugins/forge/skills/gen-test-scripts/rules/step-0.5-validation.md

### Key Decisions
无

## Document Metrics
6 acceptance criteria verified, 8 files modified, 1 file renamed/deleted, 0 stale references remaining

## Referenced Documents
- docs/proposals/surface-first-testing/proposal.md
- docs/conventions/forge-distribution.md

## Review Status
final

## Acceptance Criteria
- [x] SKILL.md no longer references docs/reference/test-type-model.md, mapping table inlined
- [x] types/ui.md renamed to types/web.md, type: ui changed to type: web, all UI references changed to Web
- [x] types/_shared.md interface type list uses Web instead of UI, stale path reference removed
- [x] All types/*.md files replaced docs/reference/test-type-model.md references with self-contained inline test type definitions
- [x] rules/step-0.5-validation.md references types/web.md instead of types/ui.md
- [x] All surfaceType references use unified web/api/cli/tui/mobile terminology

## Notes
Used git mv to preserve history for ui.md -> web.md rename. All remaining 'UI' occurrences in types/*.md refer to generic UI concepts (UI components, UI flow), not the surface type name.
