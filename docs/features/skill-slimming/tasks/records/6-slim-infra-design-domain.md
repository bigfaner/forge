---
status: "completed"
started: "2026-05-20 14:28"
completed: "2026-05-20 14:39"
time_spent: "~11m"
---

# Task Record: 6 Slim infra/design domain (init-justfile + ui-design + extract-design-md)

## Summary
Slim infra/design domain skills: init-justfile (387->327), ui-design (314->228), extract-design-md (242->132). Extracted rules/ files for project detection, self-correction, style selection, TUI panel requirements, extraction layers, platform routing, and match strategy.

## Changes

### Files Created
- plugins/forge/skills/init-justfile/rules/project-detection.md
- plugins/forge/skills/init-justfile/rules/self-correction.md
- plugins/forge/skills/ui-design/rules/style-selection.md
- plugins/forge/skills/ui-design/rules/tui-panel-requirements.md
- plugins/forge/skills/extract-design-md/rules/extraction-layers.md
- plugins/forge/skills/extract-design-md/rules/platform-routing.md
- plugins/forge/skills/extract-design-md/rules/match-strategy.md

### Files Modified
- plugins/forge/skills/init-justfile/SKILL.md
- plugins/forge/skills/ui-design/SKILL.md
- plugins/forge/skills/extract-design-md/SKILL.md

### Key Decisions
- init-justfile: extracted project detection signals/classification and self-correction error patterns to rules/
- ui-design: extracted style selection priority chain (web/mobile/TUI) and TUI panel structural requirements to rules/
- extract-design-md: extracted 5-layer extraction strategy, platform-specific routing, and match strategy to rules/
- All SKILL.md files retain complete step numbering and flow skeleton; rules/ contain only detail-level content

## Test Results
- **Tests Executed**: No
- **Passed**: 0
- **Failed**: 0
- **Coverage**: N/A (task has no tests)

## Acceptance Criteria
- [x] Each SKILL.md <= 350 lines
- [x] All step numbers and descriptions preserved
- [x] Referenced auxiliary file paths exist and are readable
- [x] Splitting style consistent with Tier 1

## Notes
无
