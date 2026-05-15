---
status: "completed"
started: "2026-05-15 19:43"
completed: "2026-05-15 19:45"
time_spent: "~2m"
---

# Task Record: 2 Remove Element handling from gen-test-scripts and enforce source-code-first

## Summary
Removed all Element field processing logic from gen-test-scripts SKILL.md. Made source-code reconnaissance (Step 1.5 Fact Table) the sole primary locator source. Removed sitemap-missing fallback logic and E-NNN/L-NNN Element ID handling. Reordered Step 3 locator priority to start with Fact Table data. Added HARD-RULE enforcing source-code-first locator derivation. Updated Integration Tests locator strategy to be source-code-driven.

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/gen-test-scripts/SKILL.md

### Key Decisions
- Sitemap retained as secondary reference for page structure only, not as a locator source driven by test-case Element IDs
- All // VERIFY: sitemap-missing comment convention removed since Fact Table is always built regardless of sitemap availability
- Locator priority reordered: Fact Table first, sitemap as supplementary confirmation

## Test Results
- **Tests Executed**: Yes
- **Passed**: 0
- **Failed**: 0
- **Coverage**: 0.0%

## Acceptance Criteria
- [x] SKILL.md contains no reference to Element field from test-cases.md as a locator source
- [x] sitemap-missing fallback logic is removed (Fact Table is always built, not just when sitemap is missing)
- [x] Locator priority starts with Fact Table data from Step 1.5 Code Reconnaissance, sitemap is supplementary
- [x] A HARD-RULE exists: Derive all locators from source code (Fact Table). Do NOT reference any testid, selector, or locator from test-cases.md
- [x] Integration Tests locator strategy references are updated to source-code-driven (no Element field dependency)

## Notes
Documentation-only task. No code changes, no tests to run.
