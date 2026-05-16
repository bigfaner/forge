---
status: "completed"
started: "2026-05-17 00:29"
completed: "2026-05-17 00:29"
time_spent: ""
---

# Task Record: 1 Refactor gen-test-cases SKILL.md into dispatcher

## Summary
Refactored monolithic gen-test-cases SKILL.md (271 lines) into a slim dispatcher (136 lines) that preserves Steps 0-2.5 (profile resolution, PRD reading, AC extraction, interface detection), adds convention loading (Step 2.6), loops through active types loading per-type instruction files (Step 3), keeps route validation (Step 3.5), and generates manifest.md (Step 4). Added conventions frontmatter field. Removed all type-specific Steps 3-4 instructions (delegated to types/{type}.md).

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/gen-test-cases/SKILL.md

### Key Decisions
- Split boundary at Step 2.5/3: Steps 0-2.5 produce type-agnostic data, Steps 3-4 delegated to per-type files
- Added Step 2.6 for convention loading: project-wide from SKILL.md frontmatter, per-type from types/{type}.md frontmatter
- Kept Step 3.5 (Route Validation) in dispatcher since it references all test cases across types
- Added Step 4 (Generate Manifest) to produce testing/manifest.md aggregator after per-type loop
- Compressed interface detection and classification tables into concise prose to meet 150-line target

## Test Results
- **Tests Executed**: No (noTest task)
- **Passed**: 0
- **Failed**: 0
- **Coverage**: N/A (task has no tests)

## Acceptance Criteria
- [x] SKILL.md is under 150 lines
- [x] Steps 0-2.5 are preserved (profile resolution, PRD reading, AC extraction, interface detection)
- [x] After Step 2.5, dispatcher loads conventions from per-type instruction frontmatter conventions field
- [x] Dispatcher loops through each active type, loading types/{type}.md via Read tool
- [x] After per-type loop, dispatcher generates testing/manifest.md with summary table + cross-type traceability
- [x] SKILL.md frontmatter includes conventions: [testing-isolation.md] for project-wide conventions
- [x] Convention loading: read per-type instruction frontmatter, check docs/conventions/{filename} exists, load or skip silently

## Notes
Documentation-only task. No compilation or testing needed. types/ files created in Task 2.
