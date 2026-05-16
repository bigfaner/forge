# E2E Test Report: task-type-refinement

**Date**: 2026-05-16
**Duration**: 2.803s
**Profile**: go-test

## Summary

| Type  | Total | Pass | Fail | Skip |
|-------|-------|------|------|------|
| CLI   | 20    | 6    | 14   | 0    |
| **All** | **20** | **6** | **14** | **0** |

**Result**: FAIL (70% failure rate)

---

## Results by Test Case

| TC ID | Test Name | Status | Duration |
|-------|-----------|--------|----------|
| TC-001 | ListTypesDisplaysFourNewBusinessTypes | FAIL | 0.08s |
| TC-002 | ListTypesShowsDeprecatedImplementation | FAIL | 0.05s |
| TC-003 | ValidateIndexAcceptsNewTypeValues | FAIL | 0.06s |
| TC-004 | BuildIndexGeneratesPipelineForFeature | FAIL | 0.08s |
| TC-005 | BuildIndexGeneratesPipelineForEnhancement | FAIL | 0.09s |
| TC-006 | BuildIndexGeneratesPipelineForFix | FAIL | 0.07s |
| TC-007 | BuildIndexSkipsPipelineForCleanupOnly | PASS | 0.07s |
| TC-008 | BuildIndexSkipsPipelineForRefactorOnly | PASS | 0.05s |
| TC-009 | BuildIndexGeneratesEvalDocForDocumentationOnly | PASS | 0.07s |
| TC-010 | BuildIndexNoPipelineNoEvalForCleanupRefactor | FAIL | 0.09s |
| TC-011 | QualityGateSkipsForCleanupOnly | PASS | 0.10s |
| TC-012 | PromptReturnsFeatureTemplate | FAIL | 0.11s |
| TC-013 | PromptReturnsCleanupTemplate | FAIL | 0.09s |
| TC-014 | PromptReturnsRefactorTemplate | FAIL | 0.10s |
| TC-015 | QualityGateCreatesFixTypeOnCompileFailure | PASS | 0.07s |
| TC-016 | QualityGateCreatesCleanupTypeOnFmtFailure | FAIL | 0.07s |
| TC-017 | QualityGateCreatesCleanupTypeOnLintFailure | FAIL | 0.04s |
| TC-018 | RecordHasReclassificationWhenTypeShifts | FAIL | 0.14s |
| TC-019 | RecordOmitsReclassificationWhenNoShift | FAIL | 0.13s |
| TC-020 | MigrateMapsImplementationToFeature | FAIL | 0.11s |

---

## Failed Tests Detail

### TC-001: ListTypesDisplaysFourNewBusinessTypes
**Error**: `forge task list-types` output does not contain "enhancement", "cleanup", or "refactor".
**Root cause**: The new business type constants have not been registered in the forge CLI yet. The output currently lists only the legacy types: implementation, documentation, doc-evaluation, fix, gate, and test-pipeline.*.

### TC-002: ListTypesShowsDeprecatedImplementation
**Error**: Output does not contain "deprecated".
**Root cause**: Same as TC-001 -- the deprecated marker for `implementation` type has not been added to the CLI output.

### TC-003: ValidateIndexAcceptsNewTypeValues
**Error**: validate-index rejects new type values (enhancement, cleanup, refactor).
**Root cause**: The validation logic has not been updated to accept the new type constants.

### TC-004: BuildIndexGeneratesPipelineForFeature
**Error**: `forge task index` does not generate test pipeline for feature-typed tasks.
**Root cause**: The `needsTestPipeline` logic has not been updated to recognize the "feature" type.

### TC-005: BuildIndexGeneratesPipelineForEnhancement
**Error**: Same pattern as TC-004 for "enhancement" type.
**Root cause**: Same as TC-004.

### TC-006: BuildIndexGeneratesPipelineForFix
**Error**: Same pattern as TC-004 for "fix" type.
**Root cause**: Same as TC-004.

### TC-010: BuildIndexNoPipelineNoEvalForCleanupRefactor
**Error**: Combined cleanup+refactor task pair does not skip pipeline generation as expected.
**Root cause**: The pipeline skip logic for cleanup/refactor-only features is not fully implemented.

### TC-012: PromptReturnsFeatureTemplate
**Error**: Feature-type prompt template not found/returned.
**Root cause**: Type-specific prompt templates have not been wired into the prompt generation logic.

### TC-013: PromptReturnsCleanupTemplate
**Error**: Cleanup-type prompt template not found.
**Root cause**: Same as TC-012.

### TC-014: PromptReturnsRefactorTemplate
**Error**: Refactor-type prompt template not found.
**Root cause**: Same as TC-012.

### TC-016: QualityGateCreatesCleanupTypeOnFmtFailure
**Error**: "cleanup" type not registered in list-types output.
**Root cause**: Same root cause as TC-001.

### TC-017: QualityGateCreatesCleanupTypeOnLintFailure
**Error**: "cleanup" type not registered in list-types output.
**Root cause**: Same root cause as TC-001.

### TC-018: RecordHasReclassificationWhenTypeShifts
**Error**: Record file not created at expected path after `forge task submit`.
**Root cause**: The `forge task submit --data` command with type reclassification is not implemented.

### TC-019: RecordOmitsReclassificationWhenNoShift
**Error**: Record file not created at expected path.
**Root cause**: Same as TC-018.

### TC-020: MigrateMapsImplementationToFeature
**Error**: `forge task migrate` does not map "implementation" to "feature". Type remains "implementation" after migration.
**Root cause**: The migration logic has not been implemented.

---

## Failure Diagnosis

**Failure rate: 70% (>30% threshold) -- Infrastructure/app-level problem.**

All 14 failures stem from a single root cause: the `task-type-refinement` feature implementation is incomplete. The forge CLI does not yet recognize the new business types (`enhancement`, `cleanup`, `refactor`) or the updated pipeline/prompt/migration logic. The 6 passing tests (TC-007 through TC-009, TC-011, TC-015) test behavior that works correctly with the current codebase (cleanup/refactor pipeline skipping, quality gate skip, and fix type registration).

**Recommendation**: Complete the feature implementation tasks (new type constants, pipeline logic, prompt templates, type reclassification, migration) before re-running these tests.

---

## Screenshots

No screenshots (CLI tests only).
