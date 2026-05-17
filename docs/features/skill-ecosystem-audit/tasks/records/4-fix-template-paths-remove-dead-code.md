---
status: "completed"
started: "2026-05-18 00:11"
completed: "2026-05-18 00:17"
time_spent: "~6m"
---

# Task Record: 4 Fix 5 template path references and remove Playwright dead code

## Summary
Fixed 5 hardcoded plugins/forge/ template path references to use relative-from-self paths (../../../eval/rubrics/...) and removed 8 dead Playwright template files from gen-test-scripts that were superseded by the types/ dispatch system.

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/breakdown-tasks/templates/validate-code-task.md
- plugins/forge/skills/breakdown-tasks/templates/validate-ux-task.md
- plugins/forge/skills/quick-tasks/templates/validate-code-task.md
- plugins/forge/skills/quick-tasks/templates/validate-ux-task.md
- plugins/forge/skills/improve-harness/templates/improvements.md

### Key Decisions
- Used ../../../eval/rubrics/ relative paths (3 levels up from templates/ to skills/ root) as specified in Hard Rules
- Deleted only .ts and .json template files, preserving validate-specs.mjs, validate-specs.test.mjs, and __test_fixtures__/ as required

## Test Results
- **Tests Executed**: Yes
- **Passed**: 20
- **Failed**: 0
- **Coverage**: 89.5%

## Acceptance Criteria
- [x] Zero plugins/forge/ hardcoded paths in source files (excluding init-justfile)
- [x] gen-test-scripts/templates/ contains no .ts or .json files
- [x] validate-specs.mjs, validate-specs.test.mjs, __test_fixtures__/ still exist

## Notes
无
