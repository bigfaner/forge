# E2E Test Report: task-lifecycle-hardening (full suite)

**Date**: 2026-05-16
**Duration**: ~5s

## Summary

| Type  | Total | Pass | Fail | Skip |
|-------|-------|------|------|------|
| CLI   | 80    | 71    | 9    | 0    |
| TUI   | 0    | 0    | 0    | 0    |
| API   | 0    | 0    | 0    | 0    |
| **All** | **80** | **71** | **9** | **0** |

**Result**: FAILED

---

## Results by Test Case

| Status | Test | Duration | Package |
|--------|------|----------|---------|
| PASS | TestTC_001_ActiveFixTaskWithSourceTaskIDEqSelfBlocksClaim | 0.3s | e2e-tests/features/task-lifecycle-hardening |
| PASS | TestTC_001_DeletedTestFilesDoNotExist | 0s | e2e-tests |
| PASS | TestTC_001_PendingFixTaskBlocksDependentBusinessTask | 0.72s | e2e-tests/features/fix-task-claim-priority |
| PASS | TestTC_001_ProposalListSortedByCreatedDescending | 0.3s | e2e-tests/features/cli-list-reverse-chronological |
| FAIL | TestTC_001_QuickModeSingleProfileTaskCount | 0.08s | e2e-tests |
| PASS | TestTC_001_SetFeatureCreatesDirectoryAndState | 0.08s | e2e-tests |
| PASS | TestTC_001_TaskIndexCreatesPerTypeTasksForMultiType | 0.07s | e2e-tests |
| PASS | TestTC_001_VerifyTuiUiDesignDirectoryDeleted | 0s | e2e-tests |
| PASS | TestTC_002_CompletedFixTaskAllowsDependentBusinessTask | 0.39s | e2e-tests/features/fix-task-claim-priority |
| PASS | TestTC_002_DeletedTestFunctionsDoNotExist | 0s | e2e-tests |
| PASS | TestTC_002_InProgressFixTaskWithSourceTaskIDEqSelfBlocksClaim | 0.17s | e2e-tests/features/task-lifecycle-hardening |
| PASS | TestTC_002_ProposalListMtimeFallback | 0.28s | e2e-tests/features/cli-list-reverse-chronological |
| FAIL | TestTC_002_QuickModeMergedTaskHasGenAndRunType | 0.06s | e2e-tests |
| PASS | TestTC_002_SetFeatureWithEmptySlugReturnsError | 0.05s | e2e-tests |
| PASS | TestTC_002_TaskIndexPerTypeTasksHaveCorrectType | 0.08s | e2e-tests |
| PASS | TestTC_002_VerifyTC020RemovedFromJustfileCanonicalE2e | 0s | e2e-tests |
| PASS | TestTC_003_CompletedFixTaskTargetingSelfDoesNotBlock | 0.37s | e2e-tests/features/task-lifecycle-hardening |
| PASS | TestTC_003_E2ETestSuiteCompilesSuccessfully | 0.27s | e2e-tests |
| PASS | TestTC_003_FixTaskClaimedBeforeBusinessTaskWhenBothEligible | 0.32s | e2e-tests/features/fix-task-claim-priority |
| PASS | TestTC_003_ProposalListEmptyDirectory | 0.29s | e2e-tests/features/cli-list-reverse-chronological |
| PASS | TestTC_003_SetFeaturePrintsSlugToStdout | 0.06s | e2e-tests |
| PASS | TestTC_003_TaskIndexSingleTypeCreatesOneGenTask | 0.07s | e2e-tests |
| PASS | TestTC_004_FeatureListSortedByMtimeDescending | 0.35s | e2e-tests/features/cli-list-reverse-chronological |
| PASS | TestTC_004_FixChainBlocksDependentTaskUntilAllFixTasksComplete | 0.48s | e2e-tests/features/fix-task-claim-priority |
| PASS | TestTC_004_PositionalArgBackwardCompatibility | 0.05s | e2e-tests |
| FAIL | TestTC_004_QuickModePerTypeCreatesIndependentGenAndRun | 0.07s | e2e-tests |
| PASS | TestTC_004_SelfBlockTakesPrecedenceOverMetRegularDependencies | 0.39s | e2e-tests/features/task-lifecycle-hardening |
| PASS | TestTC_004_TaskIndexWithoutTestCasesFallsBackToLegacy | 0.06s | e2e-tests |
| PASS | TestTC_004_ZeroUnconditionalTSkip | 0.01s | e2e-tests |
| PASS | TestTC_005_FeatureListMissingManifestToEnd | 0.19s | e2e-tests/features/cli-list-reverse-chronological |
| PASS | TestTC_005_FixTaskTargetingOtherTaskDoesNotCauseSelfBlock | 0.27s | e2e-tests/features/task-lifecycle-hardening |
| PASS | TestTC_005_QuickModeDependencyChainCorrectAfterMerge | 0.09s | e2e-tests |
| PASS | TestTC_005_QuickModeDependencyChainCorrectAfterMerge/T-quick-2_depends_on_T-quick-1 | 0s | e2e-tests |
| PASS | TestTC_005_QuickModeDependencyChainCorrectAfterMerge/T-quick-3_depends_on_T-quick-2 | 0s | e2e-tests |
| PASS | TestTC_005_QuickModeDependencyChainCorrectAfterMerge/T-quick-4_depends_on_T-quick-3 | 0s | e2e-tests |
| PASS | TestTC_005_QuickModeDependencyChainCorrectAfterMerge/T-quick-5_depends_on_T-quick-4 | 0s | e2e-tests |
| PASS | TestTC_005_SetFeatureWithWhitespaceOnlySlugReturnsError | 0.04s | e2e-tests |
| PASS | TestTC_005_TaskIndexZeroTypeTestCasesFallsBackToLegacy | 0.06s | e2e-tests |
| PASS | TestTC_005_UnrelatedFixTaskDoesNotBlockTaskWithDifferentDependency | 0.19s | e2e-tests/features/fix-task-claim-priority |
| PASS | TestTC_005_ZeroRecursiveGoTestInvocations | 0s | e2e-tests |
| PASS | TestTC_006_FeatureListEmptyDirectory | 0.16s | e2e-tests/features/cli-list-reverse-chronological |
| PASS | TestTC_006_MultipleFixTasksTargetingSelfMustAllComplete | 0.2s | e2e-tests/features/task-lifecycle-hardening |
| PASS | TestTC_006_NoFixTasksPreservesExistingClaimBehavior | 0.2s | e2e-tests/features/fix-task-claim-priority |
| PASS | TestTC_006_NoStaticFileTextGrepTests | 0s | e2e-tests |
| FAIL | TestTC_006_QuickModePerTypeDependencyFanIn | 0.08s | e2e-tests |
| PASS | TestTC_006_SetFeatureIdempotentOnRepeatedCalls | 0.14s | e2e-tests |
| PASS | TestTC_006_TaskIndexRunDependsOnAllPerTypeGenTasks | 0.08s | e2e-tests |
| PASS | TestTC_007_BlockedTaskAutoUnblockedWhenDependenciesMet | 0.17s | e2e-tests/features/task-lifecycle-hardening |
| PASS | TestTC_007_BreakdownModeUnchangedByQuickMerge | 0.06s | e2e-tests |
| PASS | TestTC_007_NoDuplicateTestFilesRootAndFeatures | 0s | e2e-tests |
| PASS | TestTC_007_SetFeatureOverwritesPreviousFeatureInState | 0.12s | e2e-tests |
| PASS | TestTC_007_TaskIndexMultiProfilePerTypeTasks | 0.11s | e2e-tests |
| PASS | TestTC_008_BlockedTaskStaysBlockedWhenDependenciesNotMet | 0.17s | e2e-tests/features/task-lifecycle-hardening |
| PASS | TestTC_008_GetCurrentFeatureReturnsStateJsonFeatureWhenPresent | 0.06s | e2e-tests |
| FAIL | TestTC_008_QuickModeMultiProfileLetterSuffixes | 0.07s | e2e-tests |
| FAIL | TestTC_008_TaskIndexQuickModePerTypeTasks | 0.11s | e2e-tests |
| PASS | TestTC_009_AutoUnblockLoggedToStdout | 0.15s | e2e-tests/features/task-lifecycle-hardening |
| PASS | TestTC_009_GetCurrentFeatureFallsBackWhenStateJsonAbsent | 0.08s | e2e-tests |
| PASS | TestTC_009_PerTypeGenScriptsMdContainsTestType | 0.1s | e2e-tests |
| PASS | TestTC_010_GetCurrentFeatureWithSourceReturnsCorrectSourceType | 0.07s | e2e-tests |
| PASS | TestTC_010_MultipleBlockedTasksUnblockedSimultaneously | 0.27s | e2e-tests/features/task-lifecycle-hardening |
| PASS | TestTC_010_TaskIndexPerTypeIdempotent | 0.17s | e2e-tests |
| PASS | TestTC_011_BlockedTaskWithActiveFixTargetingItStaysBlocked | 0.22s | e2e-tests/features/task-lifecycle-hardening |
| FAIL | TestTC_011_InferTypeMapsMergedIDsCorrectly | 0.29s | e2e-tests |
| PASS | TestTC_011_PerTypeGenScriptsMdHasCorrectTaskIDs | 0.07s | e2e-tests |
| PASS | TestTC_011_StateJsonWithNonexistentDirFallsThrough | 0.08s | e2e-tests |
| PASS | TestTC_012_CorruptStateJsonFallsThroughSilently | 0.08s | e2e-tests |
| PASS | TestTC_012_FixCompletedAutoUnblocksBlockedSourceTask | 0.16s | e2e-tests/features/task-lifecycle-hardening |
| FAIL | TestTC_012_QuickModeSingleProfileProducesFiveTasks | 0.08s | e2e-tests |
| PASS | TestTC_012_TaskIndexSharedInfrastructureNotDuplicated | 0.07s | e2e-tests |
| PASS | TestTC_013_SourceStaysBlockedWhenFixIsStillInProgress | 0.22s | e2e-tests/features/task-lifecycle-hardening |
| PASS | TestTC_013_StateJsonTakesPriorityOverGitWorktree | 0.05s | e2e-tests |
| PASS | TestTC_014_AutoDowngradedTaskAutoUnblockedWhenDepCompletes | 0.14s | e2e-tests/features/task-lifecycle-hardening |
| PASS | TestTC_014_ExistingCallersUnchangedAfterPriorityChainChange | 0.08s | e2e-tests |
| FAIL | TestTC_014_MergedTaskGeneratesCorrectMD | 0.07s | e2e-tests |
| PASS | TestTC_015_DetectTypesFromTestCasesParsesSummaryTable | 0.16s | e2e-tests |
| PASS | TestTC_015_VerboseShowsStateJsonSource | 0.06s | e2e-tests |
| PASS | TestTC_018_VerboseShowsFeaturesDirSource | 0.09s | e2e-tests |
| PASS | TestTC_019_VerboseShowsNoneWhenNoFeatureSet | 0.08s | e2e-tests |
| PASS | TestTC_020_VerboseFlagLocalToFeatureCommandOnly | 0.18s | e2e-tests |

---

## Failed Tests Detail

### TestTC_001_QuickModeSingleProfileTaskCount

- **Package**: e2e-tests
- **Duration**: 0.08s
- **Error Output**:
```
=== RUN   TestTC_001_QuickModeSingleProfileTaskCount
    quick_test_slim_cli_test.go:226: 
        	Error Trace:	Z:/project/ai/forge/tests/e2e/quick_test_slim_cli_test.go:226
        	Error:      	Not equal: 
        	            	expected: 5
        	            	actual  : 6
        	Test:       	TestTC_001_QuickModeSingleProfileTaskCount
        	Messages:   	quick mode with single profile should generate exactly 5 test pipeline tasks
--- FAIL: TestTC_001_QuickModeSingleProfileTaskCount (0.08s)

```

### TestTC_002_QuickModeMergedTaskHasGenAndRunType

- **Package**: e2e-tests
- **Duration**: 0.06s
- **Error Output**:
```
=== RUN   TestTC_002_QuickModeMergedTaskHasGenAndRunType
    quick_test_slim_cli_test.go:246: 
        	Error Trace:	Z:/project/ai/forge/tests/e2e/quick_test_slim_cli_test.go:246
        	Error:      	Not equal: 
        	            	expected: "test-pipeline.gen-and-run"
        	            	actual  : "test-pipeline.gen-scripts"
        	            	Diff:
        	            	--- Expected
        	            	+++ Actual
        	            	@@ -1 +1 @@
        	            	-test-pipeline.gen-and-run
        	            	+test-pipeline.gen-scripts
        	Test:       	TestTC_002_QuickModeMergedTaskHasGenAndRunType
        	Messages:   	T-quick-2 should have type test-pipeline.gen-and-run
--- FAIL: TestTC_002_QuickModeMergedTaskHasGenAndRunType (0.06s)

```

### TestTC_004_QuickModePerTypeCreatesIndependentGenAndRun

- **Package**: e2e-tests
- **Duration**: 0.07s
- **Error Output**:
```
=== RUN   TestTC_004_QuickModePerTypeCreatesIndependentGenAndRun
    quick_test_slim_cli_test.go:273: 
        	Error Trace:	Z:/project/ai/forge/tests/e2e/quick_test_slim_cli_test.go:273
        	Error:      	Not equal: 
        	            	expected: "test-pipeline.gen-and-run"
        	            	actual  : "test-pipeline.gen-scripts"
        	            	Diff:
        	            	--- Expected
        	            	+++ Actual
        	            	@@ -1 +1 @@
        	            	-test-pipeline.gen-and-run
        	            	+test-pipeline.gen-scripts
        	Test:       	TestTC_004_QuickModePerTypeCreatesIndependentGenAndRun
        	Messages:   	T-quick-2-api should have type test-pipeline.gen-and-run
    quick_test_slim_cli_test.go:273: 
        	Error Trace:	Z:/project/ai/forge/tests/e2e/quick_test_slim_cli_test.go:273
        	Error:      	Not equal: 
        	            	expected: "test-pipeline.gen-and-run"
        	            	actual  : "test-pipeline.gen-scripts"
        	            	Diff:
        	            	--- Expected
        	            	+++ Actual
        	            	@@ -1 +1 @@
        	            	-test-pipeline.gen-and-run
        	            	+test-pipeline.gen-scripts
        	Test:       	TestTC_004_QuickModePerTypeCreatesIndependentGenAndRun
        	Messages:   	T-quick-2-tui should have type test-pipeline.gen-and-run
--- FAIL: TestTC_004_QuickModePerTypeCreatesIndependentGenAndRun (0.07s)

```

### TestTC_006_QuickModePerTypeDependencyFanIn

- **Package**: e2e-tests
- **Duration**: 0.08s
- **Error Output**:
```
=== RUN   TestTC_006_QuickModePerTypeDependencyFanIn
    quick_test_slim_cli_test.go:366: 
        	Error Trace:	Z:/project/ai/forge/tests/e2e/quick_test_slim_cli_test.go:366
        	Error:      	[]string{"T-quick-5"} does not contain "T-quick-3"
        	Test:       	TestTC_006_QuickModePerTypeDependencyFanIn
        	Messages:   	T-quick-4 should depend on T-quick-3
--- FAIL: TestTC_006_QuickModePerTypeDependencyFanIn (0.08s)

```

### TestTC_008_QuickModeMultiProfileLetterSuffixes

- **Package**: e2e-tests
- **Duration**: 0.07s
- **Error Output**:
```
=== RUN   TestTC_008_QuickModeMultiProfileLetterSuffixes
    quick_test_slim_cli_test.go:459: 
        	Error Trace:	Z:/project/ai/forge/tests/e2e/quick_test_slim_cli_test.go:459
        	Error:      	Not equal: 
        	            	expected: "test-pipeline.gen-and-run"
        	            	actual  : "test-pipeline.gen-scripts"
        	            	Diff:
        	            	--- Expected
        	            	+++ Actual
        	            	@@ -1 +1 @@
        	            	-test-pipeline.gen-and-run
        	            	+test-pipeline.gen-scripts
        	Test:       	TestTC_008_QuickModeMultiProfileLetterSuffixes
        	Messages:   	T-quick-2a should have type test-pipeline.gen-and-run
    quick_test_slim_cli_test.go:459: 
        	Error Trace:	Z:/project/ai/forge/tests/e2e/quick_test_slim_cli_test.go:459
        	Error:      	Not equal: 
        	            	expected: "test-pipeline.gen-and-run"
        	            	actual  : "test-pipeline.gen-scripts"
        	            	Diff:
        	            	--- Expected
        	            	+++ Actual
        	            	@@ -1 +1 @@
        	            	-test-pipeline.gen-and-run
        	            	+test-pipeline.gen-scripts
        	Test:       	TestTC_008_QuickModeMultiProfileLetterSuffixes
        	Messages:   	T-quick-2b should have type test-pipeline.gen-and-run
    quick_test_slim_cli_test.go:467: 
        	Error Trace:	Z:/project/ai/forge/tests/e2e/quick_test_slim_cli_test.go:467
        	Error:      	Not equal: 
        	            	expected: "test-pipeline.graduate"
        	            	actual  : "test-pipeline.run"
        	            	Diff:
        	            	--- Expected
        	            	+++ Actual
        	            	@@ -1 +1 @@
        	            	-test-pipeline.graduate
        	            	+test-pipeline.run
        	Test:       	TestTC_008_QuickModeMultiProfileLetterSuffixes
        	Messages:   	T-quick-3a should have type test-pipeline.graduate
    quick_test_slim_cli_test.go:467: 
        	Error Trace:	Z:/project/ai/forge/tests/e2e/quick_test_slim_cli_test.go:467
        	Error:      	Not equal: 
        	            	expected: "test-pipeline.graduate"
        	            	actual  : "test-pipeline.run"
        	            	Diff:
        	            	--- Expected
        	            	+++ Actual
        	            	@@ -1 +1 @@
        	            	-test-pipeline.graduate
        	            	+test-pipeline.run
        	Test:       	TestTC_008_QuickModeMultiProfileLetterSuffixes
        	Messages:   	T-quick-3b should have type test-pipeline.graduate
    quick_test_slim_cli_test.go:474: 
        	Error Trace:	Z:/project/ai/forge/tests/e2e/quick_test_slim_cli_test.go:474
        	Error:      	Should be true
        	Test:       	TestTC_008_QuickModeMultiProfileLetterSuffixes
        	Messages:   	T-quick-4 should exist as shared task
--- FAIL: TestTC_008_QuickModeMultiProfileLetterSuffixes (0.07s)

```

### TestTC_008_TaskIndexQuickModePerTypeTasks

- **Package**: e2e-tests
- **Duration**: 0.11s
- **Error Output**:
```
=== RUN   TestTC_008_TaskIndexQuickModePerTypeTasks
    test_scripts_per_type_cli_test.go:438: 
        	Error Trace:	Z:/project/ai/forge/tests/e2e/test_scripts_per_type_cli_test.go:438
        	Error:      	Should be true
        	Test:       	TestTC_008_TaskIndexQuickModePerTypeTasks
        	Messages:   	index should contain quick-gen-and-run-go-test-ui for quick mode
    test_scripts_per_type_cli_test.go:438: 
        	Error Trace:	Z:/project/ai/forge/tests/e2e/test_scripts_per_type_cli_test.go:438
        	Error:      	Should be true
        	Test:       	TestTC_008_TaskIndexQuickModePerTypeTasks
        	Messages:   	index should contain quick-gen-and-run-go-test-api for quick mode
    test_scripts_per_type_cli_test.go:438: 
        	Error Trace:	Z:/project/ai/forge/tests/e2e/test_scripts_per_type_cli_test.go:438
        	Error:      	Should be true
        	Test:       	TestTC_008_TaskIndexQuickModePerTypeTasks
        	Messages:   	index should contain quick-gen-and-run-go-test-cli for quick mode
    test_scripts_per_type_cli_test.go:447: 
        	Error Trace:	Z:/project/ai/forge/tests/e2e/test_scripts_per_type_cli_test.go:447
        	Error:      	"---\nid: \"T-quick-4\"\ntitle: \"Graduate Quick Test Scripts (go-test)\"\npriority: \"P1\"\nestimated_time: \"15min\"\ndependencies: [\"T-quick-5\"]\ntype: \"test-pipeline.graduate\"\nscope: \"all\"\nprofile: \"go-test\"\n---\n\n# Graduate Quick Test Scripts (go-test)\n\nProfile: **go-test**\n\n# Go Test Graduate Strategy\n\nProfile-specific graduation rules for the `graduate-tests` skill.\n\n## Source File Discovery\n\n| Item | Value |\n|------|-------|\n| File extension | `_test.go` |\n| Source directory | `tests/e2e/features/<slug>/` (staging) |\n| Target directory | `tests/e2e/` (regression) |\n\n## Import Rewrite\n\n**None required.** Go uses module paths (defined in `go.mod`) rather than relative file paths. All imports resolve identically regardless of file location within the module.\n\n## Validation\n\n| Check | Command | Failure action |\n|-------|---------|----------------|\n| Pre-flight compilation | `just e2e-compile` | Abort before touching anything |\n| Post-migration compilation | `just e2e-compile` | Rollback via migration manifest |\n| Test discovery | `just e2e-discover` | Rollback via migration manifest |\n\n## Merge Procedure\n\nGo file-level merge when a target file already exists at the graduation destination:\n\n1. Read both source and target test files\n2. Backup target file (only if no backup exists -- prevents overwriting original on re-run)\n3. Combine imports, deduplicate by import path\n4. Match test functions by name -- `func TestTC_NNN_*`\n5. Deduplicate: identical function names keep source version; different names with same TC ID prefix keep both\n6. Append new test functions that don't exist in target\n7. Write the merged file\n\nMerge strategy is `package` -- all test functions reside in the same `package e2e` declaration.\n\n## Shared Infrastructure\n\nThese files already exist at `tests/e2e/` and must NOT be copied or modified during graduation:\n\n- `main_test.go` (TestMain setup/teardown)\n- `helpers_test.go` (shared test helpers)\n- `testdata/` (golden files, fixtures)\n\n## Compilation Check\n\nAfter migration, verify all packages compile correctly with the new test files in place:\n\n```bash\njust e2e-compile\n```\n\n## Test Discovery\n\nVerify all expected tests are discoverable:\n\n```bash\njust e2e-discover\n```\n\nOutput is a plain list of test function names, one per line. Compare against expected TC IDs.\n\n## Graduation Marker\n\nWritten only after validation passes (atomic -- no marker = not graduated):\n\n```yaml\nschema_version: 1\nstatus: completed\ntimestamp: <UTC ISO timestamp>\nsource: tests/e2e/features/<slug>/\ntargets:\n  - tests/e2e/<test-file>\nmodules:\n  - <module-name>\ntestCount: <N>\n```\n" does not contain "T-quick-2-ui"
        	Test:       	TestTC_008_TaskIndexQuickModePerTypeTasks
        	Messages:   	quick graduate task should depend on T-quick-2-ui
    test_scripts_per_type_cli_test.go:448: 
        	Error Trace:	Z:/project/ai/forge/tests/e2e/test_scripts_per_type_cli_test.go:448
        	Error:      	"---\nid: \"T-quick-4\"\ntitle: \"Graduate Quick Test Scripts (go-test)\"\npriority: \"P1\"\nestimated_time: \"15min\"\ndependencies: [\"T-quick-5\"]\ntype: \"test-pipeline.graduate\"\nscope: \"all\"\nprofile: \"go-test\"\n---\n\n# Graduate Quick Test Scripts (go-test)\n\nProfile: **go-test**\n\n# Go Test Graduate Strategy\n\nProfile-specific graduation rules for the `graduate-tests` skill.\n\n## Source File Discovery\n\n| Item | Value |\n|------|-------|\n| File extension | `_test.go` |\n| Source directory | `tests/e2e/features/<slug>/` (staging) |\n| Target directory | `tests/e2e/` (regression) |\n\n## Import Rewrite\n\n**None required.** Go uses module paths (defined in `go.mod`) rather than relative file paths. All imports resolve identically regardless of file location within the module.\n\n## Validation\n\n| Check | Command | Failure action |\n|-------|---------|----------------|\n| Pre-flight compilation | `just e2e-compile` | Abort before touching anything |\n| Post-migration compilation | `just e2e-compile` | Rollback via migration manifest |\n| Test discovery | `just e2e-discover` | Rollback via migration manifest |\n\n## Merge Procedure\n\nGo file-level merge when a target file already exists at the graduation destination:\n\n1. Read both source and target test files\n2. Backup target file (only if no backup exists -- prevents overwriting original on re-run)\n3. Combine imports, deduplicate by import path\n4. Match test functions by name -- `func TestTC_NNN_*`\n5. Deduplicate: identical function names keep source version; different names with same TC ID prefix keep both\n6. Append new test functions that don't exist in target\n7. Write the merged file\n\nMerge strategy is `package` -- all test functions reside in the same `package e2e` declaration.\n\n## Shared Infrastructure\n\nThese files already exist at `tests/e2e/` and must NOT be copied or modified during graduation:\n\n- `main_test.go` (TestMain setup/teardown)\n- `helpers_test.go` (shared test helpers)\n- `testdata/` (golden files, fixtures)\n\n## Compilation Check\n\nAfter migration, verify all packages compile correctly with the new test files in place:\n\n```bash\njust e2e-compile\n```\n\n## Test Discovery\n\nVerify all expected tests are discoverable:\n\n```bash\njust e2e-discover\n```\n\nOutput is a plain list of test function names, one per line. Compare against expected TC IDs.\n\n## Graduation Marker\n\nWritten only after validation passes (atomic -- no marker = not graduated):\n\n```yaml\nschema_version: 1\nstatus: completed\ntimestamp: <UTC ISO timestamp>\nsource: tests/e2e/features/<slug>/\ntargets:\n  - tests/e2e/<test-file>\nmodules:\n  - <module-name>\ntestCount: <N>\n```\n" does not contain "T-quick-2-api"
        	Test:       	TestTC_008_TaskIndexQuickModePerTypeTasks
        	Messages:   	quick graduate task should depend on T-quick-2-api
    test_scripts_per_type_cli_test.go:449: 
        	Error Trace:	Z:/project/ai/forge/tests/e2e/test_scripts_per_type_cli_test.go:449
        	Error:      	"---\nid: \"T-quick-4\"\ntitle: \"Graduate Quick Test Scripts (go-test)\"\npriority: \"P1\"\nestimated_time: \"15min\"\ndependencies: [\"T-quick-5\"]\ntype: \"test-pipeline.graduate\"\nscope: \"all\"\nprofile: \"go-test\"\n---\n\n# Graduate Quick Test Scripts (go-test)\n\nProfile: **go-test**\n\n# Go Test Graduate Strategy\n\nProfile-specific graduation rules for the `graduate-tests` skill.\n\n## Source File Discovery\n\n| Item | Value |\n|------|-------|\n| File extension | `_test.go` |\n| Source directory | `tests/e2e/features/<slug>/` (staging) |\n| Target directory | `tests/e2e/` (regression) |\n\n## Import Rewrite\n\n**None required.** Go uses module paths (defined in `go.mod`) rather than relative file paths. All imports resolve identically regardless of file location within the module.\n\n## Validation\n\n| Check | Command | Failure action |\n|-------|---------|----------------|\n| Pre-flight compilation | `just e2e-compile` | Abort before touching anything |\n| Post-migration compilation | `just e2e-compile` | Rollback via migration manifest |\n| Test discovery | `just e2e-discover` | Rollback via migration manifest |\n\n## Merge Procedure\n\nGo file-level merge when a target file already exists at the graduation destination:\n\n1. Read both source and target test files\n2. Backup target file (only if no backup exists -- prevents overwriting original on re-run)\n3. Combine imports, deduplicate by import path\n4. Match test functions by name -- `func TestTC_NNN_*`\n5. Deduplicate: identical function names keep source version; different names with same TC ID prefix keep both\n6. Append new test functions that don't exist in target\n7. Write the merged file\n\nMerge strategy is `package` -- all test functions reside in the same `package e2e` declaration.\n\n## Shared Infrastructure\n\nThese files already exist at `tests/e2e/` and must NOT be copied or modified during graduation:\n\n- `main_test.go` (TestMain setup/teardown)\n- `helpers_test.go` (shared test helpers)\n- `testdata/` (golden files, fixtures)\n\n## Compilation Check\n\nAfter migration, verify all packages compile correctly with the new test files in place:\n\n```bash\njust e2e-compile\n```\n\n## Test Discovery\n\nVerify all expected tests are discoverable:\n\n```bash\njust e2e-discover\n```\n\nOutput is a plain list of test function names, one per line. Compare against expected TC IDs.\n\n## Graduation Marker\n\nWritten only after validation passes (atomic -- no marker = not graduated):\n\n```yaml\nschema_version: 1\nstatus: completed\ntimestamp: <UTC ISO timestamp>\nsource: tests/e2e/features/<slug>/\ntargets:\n  - tests/e2e/<test-file>\nmodules:\n  - <module-name>\ntestCount: <N>\n```\n" does not contain "T-quick-2-cli"
        	Test:       	TestTC_008_TaskIndexQuickModePerTypeTasks
        	Messages:   	quick graduate task should depend on T-quick-2-cli
--- FAIL: TestTC_008_TaskIndexQuickModePerTypeTasks (0.11s)

```

### TestTC_011_InferTypeMapsMergedIDsCorrectly

- **Package**: e2e-tests
- **Duration**: 0.29s
- **Error Output**:
```
=== RUN   TestTC_011_InferTypeMapsMergedIDsCorrectly
    quick_test_slim_cli_test.go:512: 
        	Error Trace:	Z:/project/ai/forge/tests/e2e/quick_test_slim_cli_test.go:512
        	Error:      	Not equal: 
        	            	expected: "test-pipeline.gen-and-run"
        	            	actual  : "test-pipeline.gen-scripts"
        	            	Diff:
        	            	--- Expected
        	            	+++ Actual
        	            	@@ -1 +1 @@
        	            	-test-pipeline.gen-and-run
        	            	+test-pipeline.gen-scripts
        	Test:       	TestTC_011_InferTypeMapsMergedIDsCorrectly
        	Messages:   	T-quick-2-api should have type test-pipeline.gen-and-run
    quick_test_slim_cli_test.go:512: 
        	Error Trace:	Z:/project/ai/forge/tests/e2e/quick_test_slim_cli_test.go:512
        	Error:      	Not equal: 
        	            	expected: "test-pipeline.gen-and-run"
        	            	actual  : "test-pipeline.gen-scripts"
        	            	Diff:
        	            	--- Expected
        	            	+++ Actual
        	            	@@ -1 +1 @@
        	            	-test-pipeline.gen-and-run
        	            	+test-pipeline.gen-scripts
        	Test:       	TestTC_011_InferTypeMapsMergedIDsCorrectly
        	Messages:   	T-quick-2-tui should have type test-pipeline.gen-and-run
    quick_test_slim_cli_test.go:529: 
        	Error Trace:	Z:/project/ai/forge/tests/e2e/quick_test_slim_cli_test.go:529
        	Error:      	Not equal: 
        	            	expected: "test-pipeline.gen-and-run"
        	            	actual  : "test-pipeline.gen-scripts"
        	            	Diff:
        	            	--- Expected
        	            	+++ Actual
        	            	@@ -1 +1 @@
        	            	-test-pipeline.gen-and-run
        	            	+test-pipeline.gen-scripts
        	Test:       	TestTC_011_InferTypeMapsMergedIDsCorrectly
        	Messages:   	T-quick-2a should have type test-pipeline.gen-and-run
    quick_test_slim_cli_test.go:529: 
        	Error Trace:	Z:/project/ai/forge/tests/e2e/quick_test_slim_cli_test.go:529
        	Error:      	Not equal: 
        	            	expected: "test-pipeline.gen-and-run"
        	            	actual  : "test-pipeline.gen-scripts"
        	            	Diff:
        	            	--- Expected
        	            	+++ Actual
        	            	@@ -1 +1 @@
        	            	-test-pipeline.gen-and-run
        	            	+test-pipeline.gen-scripts
        	Test:       	TestTC_011_InferTypeMapsMergedIDsCorrectly
        	Messages:   	T-quick-2b should have type test-pipeline.gen-and-run
--- FAIL: TestTC_011_InferTypeMapsMergedIDsCorrectly (0.29s)

```

### TestTC_012_QuickModeSingleProfileProducesFiveTasks

- **Package**: e2e-tests
- **Duration**: 0.08s
- **Error Output**:
```
=== RUN   TestTC_012_QuickModeSingleProfileProducesFiveTasks
    quick_test_slim_cli_test.go:556: 
        	Error Trace:	Z:/project/ai/forge/tests/e2e/quick_test_slim_cli_test.go:556
        	Error:      	Not equal: 
        	            	expected: 5
        	            	actual  : 6
        	Test:       	TestTC_012_QuickModeSingleProfileProducesFiveTasks
        	Messages:   	single profile quick mode should produce exactly 5 test tasks
--- FAIL: TestTC_012_QuickModeSingleProfileProducesFiveTasks (0.08s)

```

### TestTC_014_MergedTaskGeneratesCorrectMD

- **Package**: e2e-tests
- **Duration**: 0.07s
- **Error Output**:
```
=== RUN   TestTC_014_MergedTaskGeneratesCorrectMD
    quick_test_slim_cli_test.go:576: 
        	Error Trace:	Z:/project/ai/forge/tests/e2e/quick_test_slim_cli_test.go:576
        	Error:      	Received unexpected error:
        	            	open C:\Users\panda\AppData\Local\Temp\TestTC_014_MergedTaskGeneratesCorrectMD4166453834\001\docs\features\test-qts-014\tasks\quick-gen-and-run-go-test.md: The system cannot find the file specified.
        	Test:       	TestTC_014_MergedTaskGeneratesCorrectMD
        	Messages:   	quick-gen-and-run-go-test.md should exist
--- FAIL: TestTC_014_MergedTaskGeneratesCorrectMD (0.07s)

```



---

## Screenshots

No screenshots (CLI tests only).
