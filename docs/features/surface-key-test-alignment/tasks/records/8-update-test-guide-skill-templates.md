---
status: "completed"
started: "2026-06-06 13:36"
completed: "2026-06-06 13:38"
time_spent: "~2m"
---

# Task Record: 8 Update test-guide SKILL.md and surface templates

## Summary
Updated test-guide SKILL.md quick-reference table and multi-surface notes with surface-key adaptive directory rules; updated 5 surface templates (api/cli/web/tui/mobile) directory path to tests/<surfaceKey>/<journey/> for multi-surface or tests/<journey/> for single-surface; added test directory structure section to test-type-model.md with scenario table and principles; synced convention-structure.md index template.

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/test-guide/SKILL.md
- plugins/forge/skills/test-guide/templates/surfaces/api.md
- plugins/forge/skills/test-guide/templates/surfaces/cli.md
- plugins/forge/skills/test-guide/templates/surfaces/web.md
- plugins/forge/skills/test-guide/templates/surfaces/tui.md
- plugins/forge/skills/test-guide/templates/surfaces/mobile.md
- plugins/forge/skills/test-guide/references/test-type-model.md
- plugins/forge/skills/test-guide/rules/convention-structure.md

### Key Decisions
无

## Document Metrics
7 files modified, 3 document types updated (SKILL.md, templates, references+rules)

## Referenced Documents
- docs/proposals/surface-key-test-alignment/proposal.md

## Review Status
final

## Acceptance Criteria
- [x] test-guide SKILL.md 目录约定部分包含多 surface 规则：tests/<surfaceKey>/<journey/>
- [x] 5 个 surface 模板文件测试目录路径与新规则一致
- [x] test-type-model.md 目录结构模型反映 surface-key 分区

## Notes
Also updated convention-structure.md rules file to keep index template in sync with the new directory convention. All changes are text-only updates to path descriptions — no structural or code changes.
