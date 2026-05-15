---
status: "completed"
started: "2026-05-15 19:40"
completed: "2026-05-15 19:42"
time_spent: "~2m"
---

# Task Record: 1 Remove Element field from gen-test-cases skill and template

## Summary
Removed Element field from gen-test-cases SKILL.md and template. Deleted sitemap presence check, Route Validation enhancement for element gaps, and Element-required assertions. Replaced with HARD-RULE prohibiting testid/CSS selector/XPath/implementation-specific locators in test-cases.md output.

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/gen-test-cases/SKILL.md

### Key Decisions
- Replaced the entire Element-required HARD-RULE block with a single prohibition rule: test-cases.md must NOT contain any testid, CSS selector, XPath, or implementation-specific locator
- Removed sitemap presence check and Route Validation element-gap logic entirely since they only served the Element field
- Template file (templates/test-cases.md) required no changes as it uses structural placeholders rather than explicit Element fields

## Test Results
- **Tests Executed**: Yes
- **Passed**: 0
- **Failed**: 0
- **Coverage**: 0.0%

## Acceptance Criteria
- [x] gen-test-cases SKILL.md contains no reference to 'Element' field (as a test-case output field)
- [x] SKILL.md includes a HARD-RULE prohibiting provisional testid/selector/implementation details
- [x] templates/test-cases.md has no Element column/field
- [x] Integration Test Case Generation section no longer references Element field (pattern preserved, Element line removed)

## Notes
Documentation-only task. No code changes, no tests applicable.
