---
status: "completed"
started: "2026-05-28 02:29"
completed: "2026-05-28 02:35"
time_spent: "~6m"
---

# Task Record: 9 Skill/Command/Agent 概念对齐与 scope 清理

## Summary
Cleaned up all deprecated scope concept remnants: deleted scope-to-surface-key.md migration doc, rewrote mixed.just template to use surface-aware recipe naming (frontend-compile, backend-dev, etc.) instead of scope parameter, verified task-executor.md and commands/*.md have no deprecated scope field references, added // Deprecated: comment to FrontmatterData.Scope field.

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/init-justfile/templates/mixed.just
- plugins/forge/skills/breakdown-tasks/SKILL.md
- plugins/forge/skills/breakdown-tasks/rules/existing-code-split.md
- forge-cli/pkg/task/frontmatter.go

### Key Decisions
- mixed.just: replaced scope parameter with explicit per-surface recipes (frontend-compile, backend-compile, etc.) for all 11 recipe categories
- SKILL.md: replaced rule file reference with inline note since the rule content is already documented in the Surface-Key/Type Assignment section
- FrontmatterData.Scope: added // Deprecated: doc comment; Task.Scope already had it
- CheckLegacyScope() preserved unchanged for backward-compatible migration detection

## Test Results
- **Tests Executed**: Yes
- **Passed**: 28
- **Failed**: 0
- **Coverage**: 86.1%

## Acceptance Criteria
- [x] scope-to-surface-key.md deleted
- [x] mixed.just template has no scope parameter and no frontend/backend scope values, uses surface-aware recipe naming
- [x] task-executor.md and commands/*.md have no deprecated scope field references
- [x] FrontmatterData struct scope field has // Deprecated: comment, CheckLegacyScope() preserved

## Notes
All 28 Go packages pass with race detection. Lint clean (0 issues). Commands files only contain prose scope usage (well-scoped, scope of this task) which is explicitly excluded per AC.
