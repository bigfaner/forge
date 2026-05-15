---
status: "completed"
started: "2026-05-15 21:44"
completed: "2026-05-15 21:53"
time_spent: "~9m"
---

# Task Record: 1 Add --type filter to gen-test-scripts skill

## Summary
Add --type filter to gen-test-scripts SKILL.md. Added a new 'Type Filter (--type)' section documenting the optional --type argument, its validation rules, and per-step behavior. Updated Steps 1, 1.5, 3, 3.5, and 4 to describe type filtering. Added 2 passing e2e tests (TC_008, TC_009) validating SKILL.md content.

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/gen-test-scripts/SKILL.md
- forge-cli/tests/e2e/gen_test_scripts_cli_test.go

### Key Decisions
- Type filter section placed between Step 0 and Prerequisites for early visibility
- Behavior summary table provides quick reference for each step's filtering behavior
- Type validation error message includes valid types for the active profile (actionable per BIZ-error-reporting-002)

## Test Results
- **Tests Executed**: Yes
- **Passed**: 2
- **Failed**: 0
- **Coverage**: 84.4%

## Acceptance Criteria
- [x] gen-test-scripts accepts --type <capability> argument
- [x] When --type specified, only test cases of that type are processed (Step 1 grouping filters)
- [x] Fact Table verification (Step 1.5) is skipped for non-matching types
- [x] Locator mapping (Step 3) is skipped for non-UI types when --type api or --type cli
- [x] Spec generation (Step 4) only produces files for the specified type
- [x] Shared infrastructure (Step 3.5) always runs regardless of --type value
- [x] Without --type, behavior is unchanged
- [x] Type value matches profile capability names

## Notes
无
