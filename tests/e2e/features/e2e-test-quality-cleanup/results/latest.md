# E2E Test Report: e2e-test-quality-cleanup

**Date**: 2026-05-16
**Duration**: 3.19s
**Profile**: go-test

## Summary

| Type  | Total | Pass | Fail | Skip |
|-------|-------|------|------|------|
| CLI   | 66    | 61   | 5    | 0    |
| TUI   | 0     | 0    | 0    | 0    |
| API   | 0     | 0    | 0    | 0    |
| **All** | **66** | **61** | **5** | **0** |

**Result**: FAIL (5 failures)

---

## Results by Test Case

### e2e-test-quality-cleanup (feature package)

| TC ID | Test Name | Status | Duration |
|-------|-----------|--------|----------|
| TC-001 | DeletedTestFilesDoNotExist | PASS | 0.00s |
| TC-002 | DeletedTestFunctionsDoNotExist | PASS | 0.00s |
| TC-003 | E2ETestSuiteCompilesSuccessfully | PASS | 0.26s |
| TC-004 | ZeroUnconditionalTSkip | PASS | 0.01s |
| TC-005 | ZeroRecursiveGoTestInvocations | PASS | 0.00s |
| TC-006 | NoStaticFileTextGrepTests | PASS | 0.00s |
| TC-007 | NoDuplicateTestFilesRootAndFeatures | PASS | 0.00s |

### feature_set_command (root package)

| TC ID | Test Name | Status | Duration |
|-------|-----------|--------|----------|
| TC-001 | SetFeatureCreatesDirectoryAndState | PASS | 0.06s |
| TC-002 | SetFeatureWithEmptySlugReturnsError | PASS | 0.04s |
| TC-003 | SetFeaturePrintsSlugToStdout | PASS | 0.06s |
| TC-004 | PositionalArgBackwardCompatibility | PASS | 0.05s |
| TC-005 | SetFeatureWithWhitespaceOnlySlugReturnsError | PASS | 0.04s |
| TC-006 | SetFeatureIdempotentOnRepeatedCalls | PASS | 0.10s |
| TC-007 | SetFeatureOverwritesPreviousFeatureInState | PASS | 0.11s |
| TC-008 | GetCurrentFeatureReturnsStateJsonFeatureWhenPresent | PASS | 0.04s |
| TC-009 | GetCurrentFeatureFallsBackWhenStateJsonAbsent | FAIL | 0.06s |
| TC-010 | GetCurrentFeatureWithSourceReturnsCorrectSourceType | PASS | 0.04s |
| TC-011 | StateJsonWithNonexistentDirFallsThrough | FAIL | 0.07s |
| TC-012 | CorruptStateJsonFallsThroughSilently | FAIL | 0.08s |
| TC-013 | StateJsonTakesPriorityOverGitWorktree | PASS | 0.04s |
| TC-014 | ExistingCallersUnchangedAfterPriorityChainChange | FAIL | 0.07s |
| TC-015 | VerboseShowsStateJsonSource | PASS | 0.05s |
| TC-018 | VerboseShowsFeaturesDirSource | FAIL | 0.06s |
| TC-019 | VerboseShowsNoneWhenNoFeatureSet | PASS | 0.06s |
| TC-020 | VerboseFlagLocalToFeatureCommandOnly | PASS | 0.13s |

### simplify_e2e_tests (root package)

| TC ID | Test Name | Status | Duration |
|-------|-----------|--------|----------|
| TC-001 | VerifyTuiUiDesignDirectoryDeleted | PASS | 0.00s |
| TC-002 | VerifyTC020RemovedFromJustfileCanonicalE2e | PASS | 0.00s |

### quick_mode (root package)

| TC ID | Test Name | Status | Duration |
|-------|-----------|--------|----------|
| TC-001 | QuickModeSingleProfileTaskCount | PASS | 0.06s |
| TC-002 | QuickModeMergedTaskHasGenAndRunType | PASS | 0.05s |
| TC-004 | QuickModePerTypeCreatesIndependentGenAndRun | PASS | 0.06s |
| TC-005 | QuickModeDependencyChainCorrectAfterMerge | PASS | 0.05s |
| TC-006 | QuickModePerTypeDependencyFanIn | PASS | 0.05s |
| TC-007 | BreakdownModeUnchangedByQuickMerge | PASS | 0.05s |
| TC-008 | QuickModeMultiProfileLetterSuffixes | PASS | 0.06s |
| TC-011 | InferTypeMapsMergedIDsCorrectly | PASS | 0.11s |
| TC-012 | QuickModeSingleProfileProducesFiveTasks | PASS | 0.06s |
| TC-014 | MergedTaskGeneratesCorrectMD | PASS | 0.05s |
| TC-015 | DetectTypesFromTestCasesParsesSummaryTable | PASS | 0.11s |

### per_type_gen_scripts (root package)

| TC ID | Test Name | Status | Duration |
|-------|-----------|--------|----------|
| TC-001 | TaskIndexCreatesPerTypeTasksForMultiType | PASS | 0.05s |
| TC-002 | TaskIndexPerTypeTasksHaveCorrectType | PASS | 0.05s |
| TC-003 | TaskIndexSingleTypeCreatesOneGenTask | PASS | 0.07s |
| TC-004 | TaskIndexWithoutTestCasesFallsBackToLegacy | PASS | 0.05s |
| TC-005 | TaskIndexZeroTypeTestCasesFallsBackToLegacy | PASS | 0.05s |
| TC-006 | TaskIndexRunDependsOnAllPerTypeGenTasks | PASS | 0.05s |
| TC-007 | TaskIndexMultiProfilePerTypeTasks | PASS | 0.06s |
| TC-008 | TaskIndexQuickModePerTypeTasks | PASS | 0.06s |
| TC-009 | PerTypeGenScriptsMdContainsTestType | PASS | 0.06s |
| TC-010 | TaskIndexPerTypeIdempotent | PASS | 0.09s |
| TC-011 | PerTypeGenScriptsMdHasCorrectTaskIDs | PASS | 0.06s |
| TC-012 | TaskIndexSharedInfrastructureNotDuplicated | PASS | 0.05s |

### cli-list-reverse-chronological (feature package)

| TC ID | Test Name | Status | Duration |
|-------|-----------|--------|----------|
| TC-001 | ProposalListSortedByCreatedDescending | PASS | 0.07s |
| TC-002 | ProposalListMtimeFallback | PASS | 0.09s |
| TC-003 | ProposalListEmptyDirectory | PASS | 0.06s |
| TC-004 | FeatureListSortedByMtimeDescending | PASS | 0.08s |
| TC-005 | FeatureListMissingManifestToEnd | PASS | 0.07s |
| TC-006 | FeatureListEmptyDirectory | PASS | 0.05s |

### fix-task-claim-priority (feature package)

| TC ID | Test Name | Status | Duration |
|-------|-----------|--------|----------|
| TC-001 | PendingFixTaskBlocksDependentBusinessTask | PASS | 0.12s |
| TC-002 | CompletedFixTaskAllowsDependentBusinessTask | PASS | 0.11s |
| TC-003 | FixTaskClaimedBeforeBusinessTaskWhenBothEligible | PASS | 0.11s |
| TC-004 | FixChainBlocksDependentTaskUntilAllFixTasksComplete | PASS | 0.17s |
| TC-005 | UnrelatedFixTaskDoesNotBlockTaskWithDifferentDependency | PASS | 0.09s |
| TC-006 | NoFixTasksPreservesExistingClaimBehavior | PASS | 0.09s |

---

## Failed Tests Detail

### TC-009: GetCurrentFeatureFallsBackWhenStateJsonAbsent (feature_set_command)

**File**: `feature_set_command_cli_test.go:325`
**Error**: Expected output to contain `FEATURE: lone-feature`, got `FEATURE: (none)`.
**Diagnosis**: The `forge feature` fallback chain does not resolve features from the `.forge/` config or features directory when `state.json` is absent. Feature detection only reads from state.json and returns `(none)` when it is missing.

### TC-011: StateJsonWithNonexistentDirFallsThrough (feature_set_command)

**File**: `feature_set_command_cli_test.go:358`
**Error**: Expected output to contain `FEATURE: fallback-feature`, got `FEATURE: (none)`.
**Diagnosis**: Same root cause as TC-009. The feature resolution chain does not fall through to alternative sources when state.json points to a nonexistent directory.

### TC-012: CorruptStateJsonFallsThroughSilently (feature_set_command)

**File**: `feature_set_command_cli_test.go:382`
**Error**: Expected output to contain `FEATURE: recoverable-feature`, got `FEATURE: (none)`.
**Diagnosis**: Corrupt state.json causes the feature command to return `(none)` instead of falling through to other detection methods.

### TC-014: ExistingCallersUnchangedAfterPriorityChainChange (feature_set_command)

**File**: `feature_set_command_cli_test.go:410`
**Error**: Expected output to contain `FEATURE: solo-feature`, got `FEATURE: (none)`.
**Diagnosis**: Same fallback chain issue. The feature detection only uses state.json and does not try the `.forge/features/` directory as a fallback source.

### TC-018: VerboseShowsFeaturesDirSource (feature_set_command)

**File**: `feature_set_command_cli_test.go:436`
**Error**: Expected output to contain `FEATURE: detected-feature (from: features-dir)`, got `FEATURE: (none)`.
**Diagnosis**: The features-dir detection source is not implemented. The `forge feature` command does not discover features from the `.forge/features/` directory structure.

---

## Failure Diagnosis

**Failure rate**: 5/66 = 7.6% (under 10% threshold -- per-test issues)

**Root cause**: All 5 failures share the same root cause. The `forge feature` command's feature detection priority chain only implements state.json detection. The tests expect a multi-source priority chain (state.json -> git worktree -> features dir) but only state.json is currently functional. The fallback sources (git worktree, features-dir) are not implemented in the CLI.

**Per-test fix needed**: Implement the `forge feature` fallback chain to check git worktree name and `.forge/features/` directory when state.json is absent, corrupt, or points to a nonexistent directory.

---

## Screenshots

No screenshots available (CLI-only test suite).
