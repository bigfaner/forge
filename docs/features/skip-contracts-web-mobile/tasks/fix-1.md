---
id: "fix-1"
title: "fix test: just test failure in quality gate"
priority: "P0"
estimated_time: "30min"
dependencies: []
status: pending
breaking: true
type: "coding.fix"
surface-key: "."
surface-type: "cli"
---

# fix test: just test failure in quality gate

## Root Cause

Quality gate step `just test` failed during quality-gate hook.

Error output saved to: `tests/results/raw-output.txt`

Concise error:
```
--- FAIL: TestTC_DET_001_DetectionChecksPackageJsonAsFrontendSignal (0.00s)
--- FAIL: TestTC_DET_002_DetectionChecksGoModAsBackendSignal (0.00s)
--- FAIL: TestTC_DET_003_DetectionChecksCargoTomlAsBackendSignal (0.00s)
--- FAIL: TestTC_DET_004_DetectionChecksPyprojectTomlAsBackendSignal (0.00s)
--- FAIL: TestTC_DET_005_ClassificationProducesMixedWhenBothSignalsDetected (0.00s)
--- FAIL: TestTC_DET_008_ClassificationProducesErrorWhenNoMarkersDetected (0.00s)
--- FAIL: TestTC_DET_009_SelectsBackendTemplateForPureBackendProjects (0.00s)
--- FAIL: TestTC_DET_010_SelectsFrontendTemplateForPureFrontendProjects (0.00s)
--- FAIL: TestTC_DET_011_SelectsMixedTemplateForMixedProjects (0.00s)
--- FAIL: TestTC_DET_014_AllThreeTemplatesHaveBoundaryMarkers (0.00s)
--- FAIL: TestTC_DET_018_AllThreeProjectTypeRecipeVariantsExist (0.00s)
--- FAIL: TestTC_DET_019_DetectionCorrectlyMapsSignalsToFrontendBackendCategories (0.00s)
--- FAIL: TestTC_004_FrontendProjectDetectionGeneratesScopeFreeJustfile (0.00s)
--- FAIL: TestTC_005_BackendProjectDetectionGeneratesScopeFreeJustfile (0.00s)
--- FAIL: TestTC_006_MixedProjectDetectionGeneratesScopeAwareJustfile (0.00s)
--- FAIL: TestTC_018_NoMarkerFilesDetectedCausesInitJustfileToError (0.00s)
--- FAIL: TestTC_MIX_001_ProjectTypeOutputsMixed (0.00s)
--- FAIL: TestTC_MIX_002_CompileRecipeHasPerSurfaceRecipes (0.00s)
--- FAIL: TestTC_MIX_003_BuildRecipeHasPerSurfaceRecipes (0.00s)
--- FAIL: TestTC_MIX_004_RunRecipeHasPerSurfaceRecipes (0.00s)
--- FAIL: TestTC_MIX_005_DevRecipeHasPerSurfaceRecipes (0.00s)
--- FAIL: TestTC_MIX_006_UnitTestRecipeHasPerSurfaceRecipes (0.00s)
--- FAIL: TestTC_MIX_007_LintRecipeHasPerSurfaceRecipes (0.00s)
--- FAIL: TestTC_MIX_008_FmtRecipeHasPerSurfaceRecipes (0.00s)
--- FAIL: TestTC_MIX_009_CheckRecipeHasPerSurfaceRecipes (0.00s)
--- FAIL: TestTC_MIX_010_CleanRecipeHasPerSurfaceRecipes (0.00s)
--- FAIL: TestTC_MIX_011_InstallRecipeHasPerSurfaceRecipes (0.00s)
--- FAIL: TestTC_MIX_012_PerSurfaceRecipesCoverFrontendAndBackend (0.00s)
--- FAIL: TestTC_MIX_013_EmptyBranchExecutesBothFrontendAndBackendCommands (0.00s)
--- FAIL: TestTC_MIX_014_AllBashRecipesUseSetEuoPipefail (0.00s)
--- FAIL: TestTC_MIX_015_FrontendCommandsUseNpmBackendUsesPlaceholders (0.00s)
--- FAIL: TestTC_MIX_016_ProjectTypeHasNoScopeParameter (0.00s)
--- FAIL: TestTC_MIX_017_TestRecipeHasNoScopeParameter (0.00s)
--- FAIL: TestTC_MIX_018_CiHasNoScopeParameter (0.00s)
--- FAIL: TestTC_MIX_019_TestSetupHasNoScopeParameter (0.00s)
--- FAIL: TestTC_MIX_020_ProbeHasNoScopeParameter (0.00s)
--- FAIL: TestTC_MIX_021_MixedTemplateHasForgeBoundaryMarkers (0.00s)
--- FAIL: TestTC_MIX_022_AllRecipesArePresentInMixedTemplate (0.00s)
--- FAIL: TestTC_MIX_023_CiRecipeChainsStandardCommands (0.00s)
--- FAIL: TestTC_003_TestSetupRecipePresent (0.00s)
--- FAIL: TestTC_004_TestRecipeHasJourneyParameter (0.00s)
--- FAIL: TestTC_009_JustTestSetupExits1WhenPackageJsonMissing (0.00s)
--- FAIL: TestTC_011_TestRecipeAcceptsJourneyParameter (0.00s)
--- FAIL: TestTC_012_TestRecipeFiltersByJourneyWhenProvided (0.00s)
--- FAIL: TestTC_018_InitJustfileGeneratesTestSetupTarget (0.00s)
--- FAIL: TestTC_019_InitJustfileGeneratesTestTarget (0.00s)
--- FAIL: TestTSG_001_GeneratesSummaryAndGateForQualifyingPhases (0.02s)
--- FAIL: TestTSG_002_CorrectDependencyWiringForGateTasks (0.02s)
--- FAIL: TestTSG_003_SkipsSingleTaskPhases (0.01s)
--- FAIL: TestTSG_004_ExcludesTestOnlyPhases (0.01s)
--- FAIL: TestTSG_005_FiltersTestTasksFromBusinessCount (0.02s)
--- FAIL: TestTSG_006_IdempotentRerunPreservesExistingFiles (0.01s)
--- FAIL: TestTSG_007_GeneratesOnlyMissingGateWhenSummaryExists (0.02s)
--- FAIL: TestTSG_008_GeneratesOnlyMissingSummaryWhenGateExists (0.01s)
--- FAIL: TestTSG_009_SilentlySkipsMalformedTaskIDs (0.01s)
--- FAIL: TestTSG_010_PreservesPreexistingHandCraftedGateFiles (0.01s)
--- FAIL: TestTSG_011_GeneratedTasksInIndexJsonWithCorrectType (0.01s)
--- FAIL: TestTSG_012_PrintsSummaryLinePerQualifyingPhase (0.02s)
--- FAIL: TestTSG_013_PrintsNoQualificationMessage (0.01s)
--- FAIL: TestTSG_015_QuickModeGeneratesStageGatesIdentically (0.02s)
--- FAIL: TestTSG_016_DoesNotBreakExistingIndexBehavior (0.02s)
--- FAIL: TestTSG_018_ConcurrentExecutionIdenticalOutput (0.02s)
--- FAIL: TestTSG_019_RejectsPathTraversalInTaskIDs (0.01s)
--- FAIL: TestTSG_020_GenerationCompletesWithinTimeBudget (0.07s)
--- FAIL: TestTC_TypeRefine_007_BuildIndexSkipsPipelineForCleanupOnly (0.02s)
--- FAIL: TestTC_TypeRefine_008_BuildIndexSkipsPipelineForRefactorOnly (0.04s)
--- FAIL: TestTC_TypeRefine_009_BuildIndexGeneratesReviewDocForDocumentationOnly (0.02s)
--- FAIL: TestTC_TypeRefine_010_BuildIndexNoPipelineNoEvalForCleanupRefactor (0.02s)
--- FAIL: TestTC_TypeRefine_011_QualityGateSkipsForCleanupOnly (0.02s)
```

## Reference Files

- Source: smoke_test.go, step2_load_strategy_rule_test.go, step3_execute_dev_test.go, step4_execute_probe_test.go, step5_execute_test_test.go, step6_teardown_test.go, step7_alternative_surfaces_test.go, forge_info_commands_test.go, forge_detection_test.go
- Test script: just test
- Test results: tests/results/raw-output.txt

## Fix Boundaries

When fixing test failures, observe these boundaries:

**Forbidden:**
- Starting dev server (`npx expo start`, `npm run dev`, etc.)
- Running `npm install` more than 3 times — mark task as blocked if dependency installation fails 3 times
- Running full test suite — regression is verified by the dispatcher after fix completes
- Manually opening browser to verify rendering

**Correct workflow:**
1. Read failing test + corresponding component source
2. Compare test's expected testID/selectors vs actual DOM structure
3. Modify component (add testID) or test (adjust selectors/assertions)
4. Run targeted tests on affected packages — unit tests must pass
5. Record completion

## Verification

After fixing, verify the fix works:
1. Run targeted tests on changed packages: `go test -race ./affected/package/...`
2. Replace the path with the actual packages you modified

> **Note:** Full project-wide tests run at CLI submit (`forge task submit`) — agent runs targeted tests only.

Full regression is verified by the dispatcher, not by this fix task.

When this task is recorded as completed via `task record`, the source task N/A (project-wide gate) is automatically restored to pending if all its dependencies are completed.
