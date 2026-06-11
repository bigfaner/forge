---
status: "completed"
started: "2026-05-26 22:39"
completed: "2026-05-26 22:42"
time_spent: "~3m"
---

# Task Record: 5 更新 init-justfile skill 文件术语和 recipe 别名

## Summary
Updated init-justfile SKILL.md and 5 surface rule files with surface-specific test type terminology, backward-compatible aliases, and concept doc references

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/init-justfile/SKILL.md
- plugins/forge/skills/init-justfile/rules/surfaces/cli.md
- plugins/forge/skills/init-justfile/rules/surfaces/tui.md
- plugins/forge/skills/init-justfile/rules/surfaces/api.md
- plugins/forge/skills/init-justfile/rules/surfaces/web.md
- plugins/forge/skills/init-justfile/rules/surfaces/mobile.md

### Key Decisions
无

## Document Metrics
6 files updated: 5 surface rules + 1 SKILL.md; 5 aliases added; 5 concept doc references added

## Referenced Documents
- docs/proposals/surface-test-type-model/proposal.md
- docs/reference/test-type-model.md

## Review Status
final

## Acceptance Criteria
- [x] Every surface rule file contains backward-compatible alias (alias test-e2e := <surface>-test)
- [x] Alias lines have DEPRECATED comment with v3.2.0 removal target
- [x] Recipe descriptions use surface-specific test type names (CLI functional, API functional, Web E2E, Mobile E2E)
- [x] just --list recipe names and descriptions clearly distinguish test types
- [x] Aggregate recipe descriptions updated for surface test types
- [x] All rules files reference concept doc docs/reference/test-type-model.md

## Notes
Current Forge version is 3.0.0-rc.26; alias removal target set to v3.2.0 (current+2 per proposal NFR). Existing recipe execution logic, commands, and environment variables remain unchanged — only names, descriptions, aliases, and concept references were updated.
