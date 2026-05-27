---
status: "completed"
started: "2026-05-26 22:30"
completed: "2026-05-26 22:33"
time_spent: "~3m"
---

# Task Record: 3 更新 gen-test-scripts 和 gen-journeys skill 文件术语

## Summary
Updated gen-test-scripts and gen-journeys skill files with surface-specific test type terminology, replacing generic 'e2e' labels with CLI/TUI/API functional test and Web/Mobile E2E test names per the test-type-model.md mapping

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/gen-test-scripts/SKILL.md
- plugins/forge/skills/gen-test-scripts/types/cli.md
- plugins/forge/skills/gen-test-scripts/types/tui.md
- plugins/forge/skills/gen-test-scripts/types/api.md
- plugins/forge/skills/gen-test-scripts/types/ui.md
- plugins/forge/skills/gen-test-scripts/types/mobile.md
- plugins/forge/skills/gen-test-scripts/types/_shared.md
- plugins/forge/skills/gen-journeys/rules/surface-cli.md
- plugins/forge/skills/gen-journeys/rules/surface-tui.md
- plugins/forge/skills/gen-journeys/rules/surface-api.md
- plugins/forge/skills/gen-journeys/rules/surface-web.md
- plugins/forge/skills/gen-journeys/rules/surface-mobile.md

### Key Decisions
无

## Document Metrics
12 files updated, 0 logic changes, terminology-only

## Referenced Documents
- docs/proposals/surface-test-type-model/proposal.md
- docs/reference/test-type-model.md
- docs/conventions/forge-distribution.md

## Review Status
final

## Acceptance Criteria
- [x] gen-test-scripts/SKILL.md no longer uses 'e2e' as blanket term, references test-type-model.md
- [x] gen-test-scripts/types/ 5 files use surface-specific test type names and semantic definitions
- [x] gen-journeys/rules/surface-*.md 5 files use corresponding test type names
- [x] All rules files reference docs/reference/test-type-model.md
- [x] Generated test code tags use surface-specific names (e.g. @cli-functional not @e2e)

## Notes
Only terminology and tag annotations were changed. No generation strategy logic was modified. The 'tests/e2e/features/' path restriction on SKILL.md line 228 was intentionally preserved as it refers to a legacy directory path, not test type naming.
