---
status: "completed"
started: "2026-05-16 20:45"
completed: "2026-05-16 20:55"
time_spent: "~10m"
---

# Task Record: 2 Replace isDocsOnlyFeature with needsTestPipeline and needsDocEval

## Summary
Replace isDocsOnlyFeature() with needsTestPipeline() and needsDocEval() to enable three-tier pipeline decision: docs-only gets T-eval-doc, testable types (feature/enhancement/fix) get test pipeline, cleanup/refactor-only gets neither. Updated quality_gate.go isDocsOnly() to use testable types via exported IsTestableType(). Removed now-unused hasBusinessTasks().

## Changes

### Files Created
无

### Files Modified
- forge-cli/pkg/task/build.go
- forge-cli/pkg/task/build_test.go
- forge-cli/internal/cmd/quality_gate.go
- forge-cli/internal/cmd/quality_gate_test.go

### Key Decisions
- testableTypes uses new constants TypeFeature/TypeEnhancement/TypeFix (not hardcoded strings), per Hard Rules
- needsDocEval returns false for empty maps (no business tasks = no doc eval needed), unlike old isDocsOnlyFeature which returned true for empty
- Exported IsTestableType() for quality_gate.go to call from cmd package without importing internal testableTypes map
- quality_gate.go isDocsOnly() checks ALL tasks including auto-gen (per Implementation Notes), unlike needsTestPipeline which excludes auto-gen tasks

## Test Results
- **Tests Executed**: Yes
- **Passed**: 21
- **Failed**: 0
- **Coverage**: 90.3%

## Acceptance Criteria
- [x] needsTestPipeline(tasks) bool returns true if any non-auto-gen task has type feature, enhancement, or fix
- [x] needsDocEval(tasks) bool returns true if ALL non-auto-gen tasks have type documentation
- [x] isDocsOnlyFeature() removed from build.go
- [x] BuildIndex() updated: uses needsTestPipeline() for test pipeline generation, needsDocEval() for T-eval-doc generation
- [x] quality_gate.go isDocsOnly() updated to use needsTestPipeline() logic (skip gate only when no testable types)
- [x] Three-tier behavior verified: docs-only feature gets T-eval-doc; feature/enhancement/fix feature gets test pipeline; cleanup/refactor-only feature gets neither

## Notes
Updated test helper writeTaskMD default type from TypeImplementation (deprecated) to TypeFeature. Updated TestBuildIndex_CodeFeatureUnchanged and TestBuildIndex_WithTestTasks to use TypeFeature instead of TypeImplementation.
