---
status: "completed"
started: "2026-06-02 21:49"
completed: "2026-06-02 21:53"
time_spent: "~4m"
---

# Task Record: 1 新增 surface 策略模板和 test-type-model 参考

## Summary
Created 5 per-surface strategy templates (cli/api/web/tui/mobile) with 7 required sections + assertion preference table each, and migrated complete test-type-model reference document into plugin layer

## Changes

### Files Created
- plugins/forge/skills/test-guide/templates/surfaces/cli.md
- plugins/forge/skills/test-guide/templates/surfaces/api.md
- plugins/forge/skills/test-guide/templates/surfaces/web.md
- plugins/forge/skills/test-guide/templates/surfaces/tui.md
- plugins/forge/skills/test-guide/templates/surfaces/mobile.md
- plugins/forge/skills/test-guide/references/test-type-model.md

### Files Modified
无

### Key Decisions
无

## Document Metrics
6 files, ~600 lines total; all 5 templates have exactly 8 sections (7 required + assertion table); test-type-model 100% migrated from source

## Referenced Documents
- docs/proposals/surface-first-testing/proposal.md
- docs/reference/test-type-model.md
- plugins/forge/skills/test-guide/templates/convention-template.md
- plugins/forge/skills/gen-test-scripts/types/cli.md
- plugins/forge/skills/gen-test-scripts/types/api.md
- plugins/forge/skills/gen-test-scripts/types/tui.md
- plugins/forge/skills/gen-test-scripts/types/ui.md
- plugins/forge/skills/gen-test-scripts/types/mobile.md
- plugins/forge/skills/gen-test-scripts/types/_shared.md

## Review Status
final

## Acceptance Criteria
- [x] cli.md contains 7 required sections + assertion preference table
- [x] api.md contains 7 required sections + assertion preference table
- [x] web.md contains 7 required sections + assertion preference table
- [x] tui.md contains 7 required sections + assertion preference table
- [x] mobile.md contains 7 required sections + assertion preference table
- [x] references/test-type-model.md contains complete classification, mapping table, e2e constraints and semantic definitions

## Notes
Hard Rules verified: each template has exactly 7 fixed sections + assertion preference table; assertion table columns fixed to 3 (断言库/mock机制/fixture模式); assertion content sourced from gen-test-scripts types/*.md as primary authority; no content bloat beyond the required sections
