---
status: "completed"
started: "2026-05-23 09:38"
completed: "2026-05-23 09:46"
time_spent: "~8m"
---

# Task Record: 2 Add BodyContext struct and renderBody placeholder substitution

## Summary
Added BodyContext struct and renderBody() placeholder substitution to autogen.go. BodyContext carries planning-time data (FeatureSlug, Mode, Scope, SuccessCriteria, AcceptanceCriteria, ProjectType, Interfaces) from BuildIndex() to template rendering. renderBody() handles 6 placeholder tokens with empty-field handling per spec. Updated GenerateTestTaskMD() signature from (AutoGenTaskDef, string) to (AutoGenTaskDef, BodyContext). Updated callers in build.go to pass empty BodyContext{} for backward compatibility. All existing tests pass without modification to their assertions.

## Changes

### Files Created
无

### Files Modified
- forge-cli/pkg/task/autogen.go
- forge-cli/pkg/task/autogen_test.go
- forge-cli/pkg/task/build.go

### Key Decisions
- renderBody() uses strings.ReplaceAll per Hard Rules (no regex)
- Empty SCOPE with no ## Scope heading: removeSection is a no-op, then ReplaceAll clears residual {{SCOPE}} placeholder
- removeSection simplified to single headingTitle parameter since target substring was unused
- Existing tests updated only at call site (string -> BodyContext{}), no assertion changes needed

## Test Results
- **Tests Executed**: Yes
- **Passed**: 292
- **Failed**: 0
- **Coverage**: 88.2%

## Acceptance Criteria
- [x] BodyContext struct added to autogen.go with fields: FeatureSlug, Mode, Scope, SuccessCriteria, AcceptanceCriteria, ProjectType, Interfaces
- [x] renderBody(templateContent string, def AutoGenTaskDef, ctx BodyContext) string function added
- [x] Placeholder substitution handles all 6 tokens: {{FEATURE_SLUG}}, {{MODE}}, {{SCOPE}}, {{INTERFACES}}, {{TEST_TYPE}}, {{ACCEPTANCE_CRITERIA}}
- [x] Empty fields handled per spec: {{MODE}} omit line, {{SCOPE}} omit section, {{INTERFACES}} default, {{ACCEPTANCE_CRITERIA}} default
- [x] GenerateTestTaskMD() signature updated to (def AutoGenTaskDef, ctx BodyContext) ([]byte, error)
- [x] Existing callers in build.go pass empty BodyContext{} for backward compatibility
- [x] Existing tests pass without modification (backward compatible)

## Notes
9 new test functions added: TestRenderBody_SubstitutesAllPlaceholders, TestRenderBody_EmptyMode_OmitsLine, TestRenderBody_EmptyScope_OmitsSection, TestRenderBody_EmptyInterfaces_Default, TestRenderBody_EmptyTestType_OmitsLine, TestRenderBody_EmptyAcceptanceCriteria_Default, TestRenderBody_EmptyBodyContext_KnownPlaceholdersResolved, TestGenerateTestTaskMD_WithBodyContext, TestGenerateTestTaskMD_BackwardCompat_EmptyBodyContext. Coverage target 80% exceeded at 88.2%.
