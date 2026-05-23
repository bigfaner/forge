---
status: "completed"
started: "2026-05-23 10:08"
completed: "2026-05-23 10:12"
time_spent: "~4m"
---

# Task Record: 5 Update existing tests and add body content verification tests

## Summary
Added 5 new test functions (18 sub-tests total) for body content verification: TestRenderBody_FeatureSlug, TestRenderBody_ScopeAndInterfaces, TestRenderBody_AcceptanceCriteria, TestRenderBody_EmptyFields, and TestBodyContentPerStrategy (table-driven, all 13 types). All tests pass with 90.1% coverage for autogen.go.

## Changes

### Files Created
无

### Files Modified
- forge-cli/pkg/task/autogen_test.go

### Key Decisions
- Used table-driven tests for TestBodyContentPerStrategy per Hard Rules
- Used real embedded templates (no mock embed.FS) per Hard Rules
- TestBodyContentPerStrategy covers all 13 types grouped by strategy (A: 7 types, B: 2 types, C: 4 types)

## Test Results
- **Tests Executed**: Yes
- **Passed**: 303
- **Failed**: 0
- **Coverage**: 90.1%

## Acceptance Criteria
- [x] All existing TestGenerateTestTaskMD* tests updated to pass BodyContext{} (backward compat)
- [x] New test: TestRenderBody_FeatureSlug — verifies {{FEATURE_SLUG}} substitution
- [x] New test: TestRenderBody_ScopeAndInterfaces — verifies {{SCOPE}} and {{INTERFACES}} substitution
- [x] New test: TestRenderBody_AcceptanceCriteria — verifies Strategy B criteria filling
- [x] New test: TestRenderBody_EmptyFields — verifies graceful handling (omit sections, use fallbacks)
- [x] New test: TestBodyContentPerStrategy — verifies each of the 13 types gets correct body content
- [x] Strategy A body tests: feature slug + scope + interfaces injected
- [x] Strategy B body tests: acceptance criteria pre-filled as validation checklist
- [x] Strategy C body tests: discovery strategy steps present (git diff, directory scan)
- [x] go test -race -cover ./pkg/task/... passes with 80%+ coverage

## Notes
Existing tests already used BodyContext{} signature from prior tasks. Coverage 90.1% exceeds 80% target.
