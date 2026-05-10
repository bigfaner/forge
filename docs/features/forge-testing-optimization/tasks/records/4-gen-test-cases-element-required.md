---
status: "completed"
started: "2026-05-10 13:32"
completed: "2026-05-10 13:35"
time_spent: "~3m"
---

# Task Record: 4 Make gen-test-cases Element field required

## Summary
Made Element field required in gen-test-cases SKILL.md, added sitemap-missing sentinel handling, and updated gen-test-scripts to handle sitemap-missing Element values via Fact Table DOM inference.

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/gen-test-cases/SKILL.md
- plugins/forge/skills/gen-test-scripts/SKILL.md

### Key Decisions
- Element field is now required for every test case - no exceptions
- When sitemap.json is absent, Element is set to 'sitemap-missing' sentinel value rather than blocking generation
- When sitemap exists but lacks element data for a route, Route Validation reports the gap and suggests running /gen-sitemap
- gen-test-scripts handles 'sitemap-missing' by falling back to Fact Table DOM structure from Code Reconnaissance (Step 1.5) to infer locators

## Test Results
- **Tests Executed**: No
- **Passed**: 0
- **Failed**: 0
- **Coverage**: N/A (task has no tests)

## Acceptance Criteria
- [x] gen-test-cases SKILL.md has a HARD-RULE stating Element field is required for every test case
- [x] SKILL.md documents the sitemap-missing sentinel value with WARNING note
- [x] SKILL.md documents Route Validation behavior for missing element data
- [x] test-cases.md template marks Element field as required (not optional)
- [x] gen-test-scripts SKILL.md has a note about handling sitemap-missing Element values using Fact Table DOM structure

## Notes
Markdown documentation changes only - no code compilation or tests applicable. Element field in SKILL.md template changed from 'Optional' to 'Required'. gen-test-scripts updated in Sitemap section, Step 2, and Error Handling table.
