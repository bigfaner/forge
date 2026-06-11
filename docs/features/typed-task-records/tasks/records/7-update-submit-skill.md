---
status: "completed"
started: "2026-05-23 13:36"
completed: "2026-05-23 13:38"
time_spent: "~2m"
---

# Task Record: 7 Update submit-task SKILL.md with per-type instructions

## Summary
Updated submit-task SKILL.md with per-type record format instructions covering all 5 categories (coding, doc, test, validation, gate)

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/submit-task/SKILL.md

### Key Decisions
- Split single Fields table into per-category sub-tables (Shared + 5 category-specific)
- Added Type-Specific Record Formats section after Fields with JSON examples for all 5 categories
- Rewrote Metrics Collection section with per-category expectations instead of coding-only rules

## Test Results
- **Tests Executed**: No
- **Passed**: 0
- **Failed**: 0
- **Coverage**: N/A (task has no tests)

## Acceptance Criteria
- [x] SKILL.md has JSON format examples for all 5 categories: coding, doc, test, validation, gate
- [x] Doc example shows referencedDocs, reviewStatus, docMetrics fields (no test fields)
- [x] Test example shows casesGenerated, casesEvaluated, scriptsCreated, testResults fields
- [x] Validation example shows validationPassed, issuesFound fields
- [x] Gate example shows gatePassed, gateChecks fields
- [x] Coding example unchanged from current format
- [x] Field table updated with per-category required/optional annotations
- [x] Metrics Collection section clarifies per-category expectations

## Notes
Followed hard rules: kept existing structure intact, added sections rather than reorganizing, did not change CLI command syntax or workflow steps
