---
id: "fix-2"
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
--- FAIL: TestTSG_001_GeneratesSummaryAndGateForQualifyingPhases (0.02s)
--- FAIL: TestTSG_002_CorrectDependencyWiringForGateTasks (0.02s)
--- FAIL: TestTSG_003_SkipsSingleTaskPhases (0.01s)
--- FAIL: TestTSG_004_ExcludesTestOnlyPhases (0.01s)
--- FAIL: TestTSG_005_FiltersTestTasksFromBusinessCount (0.02s)
--- FAIL: TestTSG_006_IdempotentRerunPreservesExistingFiles (0.01s)
--- FAIL: TestTSG_007_GeneratesOnlyMissingGateWhenSummaryExists (0.02s)
--- FAIL: TestTSG_008_GeneratesOnlyMissingSummaryWhenGateExists (0.02s)
--- FAIL: TestTSG_009_SilentlySkipsMalformedTaskIDs (0.02s)
--- FAIL: TestTSG_010_PreservesPreexistingHandCraftedGateFiles (0.01s)
--- FAIL: TestTSG_011_GeneratedTasksInIndexJsonWithCorrectType (0.01s)
--- FAIL: TestTSG_012_PrintsSummaryLinePerQualifyingPhase (0.02s)
--- FAIL: TestTSG_013_PrintsNoQualificationMessage (0.01s)
--- FAIL: TestTSG_015_QuickModeGeneratesStageGatesIdentically (0.02s)
--- FAIL: TestTSG_016_DoesNotBreakExistingIndexBehavior (0.02s)
--- FAIL: TestTSG_018_ConcurrentExecutionIdenticalOutput (0.01s)
--- FAIL: TestTSG_019_RejectsPathTraversalInTaskIDs (0.02s)
--- FAIL: TestTSG_020_GenerationCompletesWithinTimeBudget (0.09s)
--- FAIL: TestTC_TypeRefine_007_BuildIndexSkipsPipelineForCleanupOnly (0.02s)
--- FAIL: TestTC_TypeRefine_008_BuildIndexSkipsPipelineForRefactorOnly (0.02s)
--- FAIL: TestTC_TypeRefine_009_BuildIndexGeneratesReviewDocForDocumentationOnly (0.02s)
--- FAIL: TestTC_TypeRefine_010_BuildIndexNoPipelineNoEvalForCleanupRefactor (0.02s)
--- FAIL: TestTC_TypeRefine_011_QualityGateSkipsForCleanupOnly (0.03s)
--- FAIL: TestTC_001_QuickModeSingleProfileTaskCount (0.03s)
--- FAIL: TestTC_005_QuickModeDependencyChainCorrect (0.03s)
    --- FAIL: TestTC_005_QuickModeDependencyChainCorrect/T-test-run-backend_depends_on_T-test-gen-journeys (0.00s)
    --- FAIL: TestTC_005_QuickModeDependencyChainCorrect/T-test-run-frontend_depends_on_T-test-run-backend_(serial_chain) (0.00s)
--- FAIL: TestTC_006_QuickModeRunTestSerialChainFanIn (0.04s)
--- FAIL: TestTC_012_QuickModeSingleProfileProducesCorrectTaskCount (0.03s)
--- FAIL: TestPerType_TC_006_TaskIndexRunDependsOnAllPerTypeGenTasks (0.03s)
```

## Reference Files

- Source: smoke_test.go, step2_load_strategy_rule_test.go, step3_execute_dev_test.go, step4_execute_probe_test.go, step5_execute_test_test.go, step6_teardown_test.go, step7_alternative_surfaces_test.go, forge_info_commands_test.go, mixed_cli_test.go, prompt_test.go
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
