---
status: "completed"
started: "2026-05-17 13:30"
completed: "2026-05-17 13:36"
time_spent: "~6m"
---

# Task Record: 6 Refactor main SKILL.md to dispatcher pattern

## Summary
Refactored gen-test-scripts SKILL.md from 530-line monolith to 231-line dispatcher pattern. Removed all type-specific content (Steps 2-3 sitemap/locators, type-specific Fact Table keys, type-specific generation patterns) that was extracted into types/{type}.md files in tasks 1-5. Added Step 4 dispatch loop with hard-coded type-to-file mapping table (ui/tui/mobile/api/cli). Updated convention loading to read from own type files instead of gen-test-cases/type files. Preserved Step 3.5 shared infrastructure (always runs), Steps 0-1 (profile resolution, test case reading, auth classification), Step 1.5 generic Fact Table framework, and all type-agnostic HARD-RULEs.

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/gen-test-scripts/SKILL.md

### Key Decisions
- Hard-coded type-to-file mapping table in Step 4 per proposal design decision (not dynamic)
- Convention loading reads from gen-test-scripts/types/{type}.md (not gen-test-cases/types/) to eliminate cross-skill coupling
- Antipattern guard table condensed to tabular format to fit within 250-line budget while preserving all 6 patterns
- Step 1.5 Fact Table completeness gate delegates type-specific required keys to type files, keeping only generic UNKNOWN/gate logic in dispatcher

## Test Results
- **Tests Executed**: No
- **Passed**: 20
- **Failed**: 0
- **Coverage**: N/A (task has no tests)

## Acceptance Criteria
- [x] SKILL.md line count is at most 250 lines
- [x] grep -c 'gen-test-cases/types' SKILL.md returns 0
- [x] grep -c 'Sitemap' SKILL.md returns 0
- [x] grep -c 'Locator' SKILL.md returns 0
- [x] Convention loading references gen-test-scripts/types/{type}.md
- [x] Step 4 contains dispatch loop with type-to-file mapping table
- [x] Error scenarios documented: unknown --type, missing type file
- [x] Step 3.5 shared infrastructure preserved, always runs
- [x] Steps 0-1 preserved
- [x] Step 1.5 retains generic Fact Table framework, delegates type-specific requirements
- [x] Frontmatter conventions field still lists testing-isolation.md
- [x] --type filter documentation section preserved

## Notes
This is a documentation-only task (no code changes). Coverage -1.0 because there are no code tests for markdown instruction files. All 20 existing test packages pass unchanged.
