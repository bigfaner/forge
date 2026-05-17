---
status: "completed"
started: "2026-05-17 00:41"
completed: "2026-05-17 00:43"
time_spent: "~2m"
---

# Task Record: 6 Update gen-test-scripts input discovery + convention loading

## Summary
Updated gen-test-scripts SKILL.md with per-type input discovery (glob fallback), convention loading from gen-test-cases type instruction frontmatter, updated --type filter for per-type mode, and updated Step Actionability gate for per-type file paths.

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/gen-test-scripts/SKILL.md

### Key Decisions
- Convention Loading placed as standalone section after Step 0 (profile resolution) and before Type Filter, matching the proposal's loading point after profile resolution
- Per-type file discovery uses first-match-wins logic: glob per-type files first, fall back to legacy test-cases.md
- Abort message template uses {type}-test-cases.md placeholder to work for both per-type and legacy modes

## Test Results
- **Tests Executed**: No (noTest task)
- **Passed**: 0
- **Failed**: 0
- **Coverage**: N/A (task has no tests)

## Acceptance Criteria
- [x] Prerequisite check: glob testing/*-test-cases.md first; if found, accept per-type mode; if empty, fall back to testing/test-cases.md for legacy mode
- [x] Step 1 Read Test Cases: when reading per-type files, skip the type grouping step (file is already single-type)
- [x] --type filter: when using per-type files, --type selects which {type}-test-cases.md to read; when using legacy, behavior unchanged
- [x] Step Actionability gate: check eval reports for the specific per-type file being processed
- [x] Convention loading: after profile resolution, read the active type's instruction file frontmatter, extract conventions field, load existing files from docs/conventions/, skip missing silently
- [x] gen-test-scripts frontmatter includes conventions: [testing-isolation.md] for project-wide conventions

## Notes
无
