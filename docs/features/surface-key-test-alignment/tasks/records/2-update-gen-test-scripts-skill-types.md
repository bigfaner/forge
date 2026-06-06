---
status: "completed"
started: "2026-06-06 13:17"
completed: "2026-06-06 13:20"
time_spent: "~3m"
---

# Task Record: 2 Update gen-test-scripts SKILL.md and type templates

## Summary
Updated gen-test-scripts SKILL.md and all 6 type template files with adaptive output directory rules: multi-surface -> tests/<surfaceKey>/<journey>/, single-surface -> tests/<journey>/

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/gen-test-scripts/SKILL.md
- plugins/forge/skills/gen-test-scripts/types/_shared.md
- plugins/forge/skills/gen-test-scripts/types/api.md
- plugins/forge/skills/gen-test-scripts/types/cli.md
- plugins/forge/skills/gen-test-scripts/types/mobile.md
- plugins/forge/skills/gen-test-scripts/types/tui.md
- plugins/forge/skills/gen-test-scripts/types/web.md

### Key Decisions
无

## Document Metrics
7 files modified, all output directory references updated to adaptive rule

## Referenced Documents
- docs/proposals/surface-key-test-alignment/proposal.md

## Review Status
final

## Acceptance Criteria
- [x] SKILL.md output directory rule updated: multi-surface -> tests/<surfaceKey>/<journey>/, single-surface -> tests/<journey>/
- [x] All 5 surface type files (api.md, cli.md, mobile.md, tui.md, web.md) output directory description consistent with SKILL.md
- [x] _shared.md contains multi-surface output directory guidance

## Notes
SKILL.md Output Directory section rewritten with multi-surface and single-surface examples. _shared.md gained new 'Output Directory: Adaptive Per Surface Count' section. All 5 type files Output section updated. rules/ files (step-1-contract-loading.md, run-to-learn.md) still reference tests/<journey>/ but are out of scope for this task (covered by task 4).
