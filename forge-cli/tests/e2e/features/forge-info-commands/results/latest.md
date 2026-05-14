# E2E Test Report: forge-info-commands

**Date**: 2026-05-14
**Duration**: 0.26s

## Summary

| Type  | Total | Pass | Fail | Skip |
|-------|-------|------|------|------|
| CLI   | 73  | 17  | 14  | 42  |
| **All** | **73** | **17** | **14** | **42** |

**Result**: FAIL

---

## Results by Test Case

| TC ID | Test Name | Type | Status | Duration |
|-------|-----------|------|--------|----------|
| TC-001 | TestTC_001_HelpOutputShowsCommandGroups | CLI | PASS | 0.02s |
| TC-002 | TestTC_002_TaskSubcommandHelpShowsAllCommands | CLI | PASS | 0.01s |
| TC-003 | TestTC_003_UnknownCommandReturnsErrorWithSuggestion | CLI | PASS | 0.01s |
| TC-004 | TestTC_004_UnknownTaskSubcommandReturnsErrorWithList | CLI | PASS | 0.01s |
| TC-006 | TestTC_006_GetPromptNonexistentTaskIDReturnsError | CLI | PASS | 0.01s |
| TC-010 | TestTC_010_SubmitTaskMissingResultFlagReturnsError | CLI | PASS | 0.01s |
| TC-017 | TestTC_017_FeatureNoArgsKeepsExistingBehavior | CLI | PASS | 0.01s |
| TC-020 | TestTC_020_E2ERunNonexistentFeatureReturnsError | CLI | PASS | 0.02s |
| TC-021 | TestTC_021_ListTypesOutputsAllWithDescriptions | CLI | PASS | 0.00s |
| TC-023 | TestTC_023_ForensicSearchScansHistoryAndReturnsSessions | CLI | PASS | 0.01s |
| TC-026 | TestTC_026_ForensicExtractNonexistentPathReturnsError | CLI | PASS | 0.01s |
| TC-027 | TestTC_027_ProfileDetectScansAndOutputsProfiles | CLI | PASS | 0.00s |
| TC-029 | TestTC_029_ProfileGetOutputsStrategyFileContent | CLI | PASS | 0.00s |
| TC-030 | TestTC_030_ProfileSetInvalidProfileReturnsErrorWithList | CLI | PASS | 0.00s |
| TC-035 | TestTC_035_TaskStatusNonexistentIDReturnsError | CLI | PASS | 0.02s |
| TC-036 | TestTC_036_ForensicSearchNoResultsReturnsEmptyOutput | CLI | PASS | 0.02s |
| TC-041 | TestTC_041_ProfileGetInvalidProfileReturnsErrorWithList | CLI | PASS | 0.00s |
| TC-004 | TestTC_004_ConfigGetProjectType | CLI | FAIL | 0.01s |
| TC-005 | TestTC_005_ConfigGetCapabilitiesArrayOutput | CLI | FAIL | 0.01s |
| TC-006 | TestTC_006_ConfigGetMissingKey | CLI | FAIL | 0.01s |
| TC-008 | TestTC_008_ProposalListAllProposals | CLI | FAIL | 0.00s |
| TC-009 | TestTC_009_ProposalSlugDetailView | CLI | FAIL | 0.00s |
| TC-010 | TestTC_010_ProposalCreatedDateFromFrontmatter | CLI | FAIL | 0.01s |
| TC-011 | TestTC_011_ProposalPRDColumnChecksPrdSpec | CLI | FAIL | 0.00s |
| TC-012 | TestTC_012_ProposalFeatureColumnReadsManifestStatus | CLI | FAIL | 0.00s |
| TC-013 | TestTC_013_FeatureListAllFeatures | CLI | FAIL | 0.01s |
| TC-015 | TestTC_015_FeatureListScoresFromFrontmatter | CLI | FAIL | 0.01s |
| TC-016 | TestTC_016_FeatureStatusDetailView | CLI | FAIL | 0.00s |
| TC-018 | TestTC_018_LessonListAllLessons | CLI | FAIL | 0.01s |
| TC-019 | TestTC_019_LessonCategoryFromFilenamePrefix | CLI | FAIL | 0.01s |
| TC-020 | TestTC_020_LessonNameDetailView | CLI | FAIL | 0.01s |
| TC-001 | TestTC_001_ConfigInitInteractiveSetup | CLI | SKIP | - |
| TC-002 | TestTC_002_ConfigInitReconfigurePrompt | CLI | SKIP | - |
| TC-003 | TestTC_003_ConfigInitReconfigureAccepted | CLI | SKIP | - |
| TC-005 | TestTC_005_GetPromptByTaskIDReturnsCorrectPrompt | CLI | SKIP | - |
| TC-007 | TestTC_007_ForgeConfigStructFields | CLI | SKIP | - |
| TC-007 | TestTC_007_GetPromptMissingOrInvalidTypeReturnsError | CLI | SKIP | - |
| TC-008 | TestTC_008_SubmitTaskSuccessUpdatesStatusAndCreatesRecord | CLI | SKIP | - |
| TC-009 | TestTC_009_SubmitTaskAlreadyTerminalStateReturnsError | CLI | SKIP | - |
| TC-011 | TestTC_011_ConcurrentSubmitHandlesLockContention | CLI | SKIP | - |
| TC-012 | TestTC_012_CleanupRemovesTerminalStateFiles | CLI | SKIP | - |
| TC-013 | TestTC_013_QualityGateRunsCompileFmtLintTestSequence | CLI | SKIP | - |
| TC-014 | TestTC_014_CleanupNoTerminalTasksOutputsMessage | CLI | SKIP | - |
| TC-014 | TestTC_014_FeatureListProgressFromIndexJSON | CLI | SKIP | - |
| TC-015 | TestTC_015_QualityGateCreatesNewFixTaskOnRepeatedFailure | CLI | SKIP | - |
| TC-016 | TestTC_016_QualityGateStopsCreatingFixTasksAfterMax3 | CLI | SKIP | - |
| TC-017 | TestTC_017_E2ERunWithConfiguredProfileExecutesSuite | CLI | SKIP | - |
| TC-018 | TestTC_018_E2ERunNoProfileConfiguredReturnsError | CLI | SKIP | - |
| TC-019 | TestTC_019_E2ERunUnknownProfileReturnsErrorWithList | CLI | SKIP | - |
| TC-021 | TestTC_021_InitCreatesForgeDir | CLI | SKIP | - |
| TC-022 | TestTC_022_InitGeneratesCLAUDEmd | CLI | SKIP | - |
| TC-022 | TestTC_022_ListTypesEmptyRegistryReturnsEmpty | CLI | SKIP | - |
| TC-023 | TestTC_023_InitAppendsGitignoreWithDedup | CLI | SKIP | - |
| TC-024 | TestTC_024_ForensicExtractOutputsEvidenceSummary | CLI | SKIP | - |
| TC-024 | TestTC_024_InitGitignoreDedupSkipsExisting | CLI | SKIP | - |
| TC-025 | TestTC_025_ForensicSubagentsListsTranscripts | CLI | SKIP | - |
| TC-025 | TestTC_025_InitAppendsJustfileRecipes | CLI | SKIP | - |
| TC-026 | TestTC_026_InitJustfileDedupSkipsExisting | CLI | SKIP | - |
| TC-027 | TestTC_027_InitRunsConfigInitWhenNoConfig | CLI | SKIP | - |
| TC-028 | TestTC_028_InitSkipsExistingFiles | CLI | SKIP | - |
| TC-028 | TestTC_028_ProfileSetUpdatesConfigWithValidProfile | CLI | SKIP | - |
| TC-029 | TestTC_029_InitResultReportFormat | CLI | SKIP | - |
| TC-030 | TestTC_030_ResolveScopeReadsConfigDirectly | CLI | SKIP | - |
| TC-031 | TestTC_031_ResolveScopeMissingConfigReturnsEmpty | CLI | SKIP | - |
| TC-031 | TestTC_031_TaskClaimNoAvailableTasksReturnsError | CLI | SKIP | - |
| TC-032 | TestTC_032_JustfileHasNoProjectTypeRecipe | CLI | SKIP | - |
| TC-032 | TestTC_032_TaskClaimCorruptedIndexReturnsError | CLI | SKIP | - |
| TC-033 | TestTC_033_TaskCheckDepsUnmetDependencyReturnsError | CLI | SKIP | - |
| TC-034 | TestTC_034_TaskValidateIndexInvalidSchemaReturnsError | CLI | SKIP | - |
| TC-037 | TestTC_037_ForensicSearchMissingRecordsDirReturnsError | CLI | SKIP | - |
| TC-038 | TestTC_038_VerifyTaskDoneIncompleteTasksReturnsError | CLI | SKIP | - |
| TC-039 | TestTC_039_TaskSubmitConcurrentWriteConflictReturnsRetryError | CLI | SKIP | - |
| TC-040 | TestTC_040_TaskSubmitMissingIndexReturnsError | CLI | SKIP | - |

---

## Failed Tests Detail

### TC-004: TestTC_004_ConfigGetProjectType

```
=== RUN   TestTC_004_ConfigGetProjectType
    forge_info_commands_cli_test.go:39: 
        	Error Trace:	/Users/fanhuifeng/Projects/ai/coding-harness/forge/forge-cli/tests/e2e/features/forge-info-commands/forge_info_commands_cli_test.go:39
        	Error:      	Not equal: 
        	            	expected: 0
        	            	actual  : 1
        	Test:       	TestTC_004_ConfigGetProjectType
        	Messages:   	config get project-type should exit 0
    forge_info_commands_cli_test.go:40: 
        	Error Trace:	/Users/fanhuifeng/Projects/ai/coding-harness/forge/forge-cli/tests/e2e/features/forge-info-commands/forge_info_commands_cli_test.go:40
        	Error:      	Not equal: 
        	            	expected: "backend"
        	            	actual  : "Error: unknown command \"config\" for \"forge\"\nRun 'forge --help' for usage."
        	            	
        	            	Diff:
        	            	--- Expected
        	            	+++ Actual
        	            	@@ -1 +1,2 @@
        	            	-backend
        	            	+Error: unknown command "config" for "forge"
        	            	+Run 'forge --help' for usage.
        	Test:       	TestTC_004_ConfigGetProjectType
        	Messages:   	config get project-type should output plain text 'backend' without formatting
--- FAIL: TestTC_004_ConfigGetProjectType (0.01s)

```

### TC-005: TestTC_005_ConfigGetCapabilitiesArrayOutput

```
=== RUN   TestTC_005_ConfigGetCapabilitiesArrayOutput
    forge_info_commands_cli_test.go:48: 
        	Error Trace:	/Users/fanhuifeng/Projects/ai/coding-harness/forge/forge-cli/tests/e2e/features/forge-info-commands/forge_info_commands_cli_test.go:48
        	Error:      	Not equal: 
        	            	expected: 0
        	            	actual  : 1
        	Test:       	TestTC_005_ConfigGetCapabilitiesArrayOutput
        	Messages:   	config get capabilities should exit 0
    forge_info_commands_cli_test.go:50: 
        	Error Trace:	/Users/fanhuifeng/Projects/ai/coding-harness/forge/forge-cli/tests/e2e/features/forge-info-commands/forge_info_commands_cli_test.go:50
        	Error:      	Should be true
        	Test:       	TestTC_005_ConfigGetCapabilitiesArrayOutput
        	Messages:   	capabilities should output at least 3 lines (one per item), got 2: "Error: unknown command \"config\" for \"forge\"\nRun 'forge --help' for usage.\n"
    forge_info_commands_cli_test.go:55: 
        	Error Trace:	/Users/fanhuifeng/Projects/ai/coding-harness/forge/forge-cli/tests/e2e/features/forge-info-commands/forge_info_commands_cli_test.go:55
        	Error:      	"Error: unknown command \"config\" for \"forge\"" should not contain "\""
        	Test:       	TestTC_005_ConfigGetCapabilitiesArrayOutput
        	Messages:   	each line should not contain quotes: "Error: unknown command \"config\" for \"forge\""
--- FAIL: TestTC_005_ConfigGetCapabilitiesArrayOutput (0.01s)

```

### TC-006: TestTC_006_ConfigGetMissingKey

```
=== RUN   TestTC_006_ConfigGetMissingKey
    forge_info_commands_cli_test.go:65: 
        	Error Trace:	/Users/fanhuifeng/Projects/ai/coding-harness/forge/forge-cli/tests/e2e/features/forge-info-commands/forge_info_commands_cli_test.go:65
        	Error:      	Not equal: 
        	            	expected: ""
        	            	actual  : "Error: unknown command \"config\" for \"forge\"\nRun 'forge --help' for usage."
        	            	
        	            	Diff:
        	            	--- Expected
        	            	+++ Actual
        	            	@@ -1 +1,2 @@
        	            	-
        	            	+Error: unknown command "config" for "forge"
        	            	+Run 'forge --help' for usage.
        	Test:       	TestTC_006_ConfigGetMissingKey
        	Messages:   	config get with missing key should produce no stdout output
--- FAIL: TestTC_006_ConfigGetMissingKey (0.01s)

```

### TC-008: TestTC_008_ProposalListAllProposals

```
=== RUN   TestTC_008_ProposalListAllProposals
    forge_info_commands_cli_test.go:82: 
        	Error Trace:	/Users/fanhuifeng/Projects/ai/coding-harness/forge/forge-cli/tests/e2e/features/forge-info-commands/forge_info_commands_cli_test.go:82
        	Error:      	Not equal: 
        	            	expected: 0
        	            	actual  : 1
        	Test:       	TestTC_008_ProposalListAllProposals
        	Messages:   	proposal list should exit 0
    forge_info_commands_cli_test.go:88: 
        	Error Trace:	/Users/fanhuifeng/Projects/ai/coding-harness/forge/forge-cli/tests/e2e/features/forge-info-commands/forge_info_commands_cli_test.go:88
        	Error:      	Should be true
        	Test:       	TestTC_008_ProposalListAllProposals
        	Messages:   	proposal list output should contain SLUG column header: Error: unknown command "proposal" for "forge"
        	            	Run 'forge --help' for usage.
    forge_info_commands_cli_test.go:90: 
        	Error Trace:	/Users/fanhuifeng/Projects/ai/coding-harness/forge/forge-cli/tests/e2e/features/forge-info-commands/forge_info_commands_cli_test.go:90
        	Error:      	Should be true
        	Test:       	TestTC_008_ProposalListAllProposals
        	Messages:   	proposal list output should contain CREATED column header: Error: unknown command "proposal" for "forge"
        	            	Run 'forge --help' for usage.
    forge_info_commands_cli_test.go:92: 
        	Error Trace:	/Users/fanhuifeng/Projects/ai/coding-harness/forge/forge-cli/tests/e2e/features/forge-info-commands/forge_info_commands_cli_test.go:92
        	Error:      	Should be true
        	Test:       	TestTC_008_ProposalListAllProposals
        	Messages:   	proposal list output should contain STATUS column header: Error: unknown command "proposal" for "forge"
        	            	Run 'forge --help' for usage.
    forge_info_commands_cli_test.go:94: 
        	Error Trace:	/Users/fanhuifeng/Projects/ai/coding-harness/forge/forge-cli/tests/e2e/features/forge-info-commands/forge_info_commands_cli_test.go:94
        	Error:      	Should be true
        	Test:       	TestTC_008_ProposalListAllProposals
        	Messages:   	proposal list output should contain PRD column header: Error: unknown command "proposal" for "forge"
        	            	Run 'forge --help' for usage.
    forge_info_commands_cli_test.go:96: 
        	Error Trace:	/Users/fanhuifeng/Projects/ai/coding-harness/forge/forge-cli/tests/e2e/features/forge-info-commands/forge_info_commands_cli_test.go:96
        	Error:      	Should be true
        	Test:       	TestTC_008_ProposalListAllProposals
        	Messages:   	proposal list output should contain FEATURE column header: Error: unknown command "proposal" for "forge"
        	            	Run 'forge --help' for usage.
--- FAIL: TestTC_008_ProposalListAllProposals (0.00s)

```

### TC-009: TestTC_009_ProposalSlugDetailView

```
=== RUN   TestTC_009_ProposalSlugDetailView
    forge_info_commands_cli_test.go:104: 
        	Error Trace:	/Users/fanhuifeng/Projects/ai/coding-harness/forge/forge-cli/tests/e2e/features/forge-info-commands/forge_info_commands_cli_test.go:104
        	Error:      	Not equal: 
        	            	expected: 0
        	            	actual  : 1
        	Test:       	TestTC_009_ProposalSlugDetailView
        	Messages:   	proposal detail should exit 0
    forge_info_commands_cli_test.go:110: 
        	Error Trace:	/Users/fanhuifeng/Projects/ai/coding-harness/forge/forge-cli/tests/e2e/features/forge-info-commands/forge_info_commands_cli_test.go:110
        	Error:      	Should be true
        	Test:       	TestTC_009_ProposalSlugDetailView
        	Messages:   	proposal detail should show SLUG field: Error: unknown command "proposal" for "forge"
        	            	Run 'forge --help' for usage.
    forge_info_commands_cli_test.go:112: 
        	Error Trace:	/Users/fanhuifeng/Projects/ai/coding-harness/forge/forge-cli/tests/e2e/features/forge-info-commands/forge_info_commands_cli_test.go:112
        	Error:      	Should be true
        	Test:       	TestTC_009_ProposalSlugDetailView
        	Messages:   	proposal detail should show CREATED field: Error: unknown command "proposal" for "forge"
        	            	Run 'forge --help' for usage.
    forge_info_commands_cli_test.go:114: 
        	Error Trace:	/Users/fanhuifeng/Projects/ai/coding-harness/forge/forge-cli/tests/e2e/features/forge-info-commands/forge_info_commands_cli_test.go:114
        	Error:      	Should be true
        	Test:       	TestTC_009_ProposalSlugDetailView
        	Messages:   	proposal detail should show STATUS field: Error: unknown command "proposal" for "forge"
        	            	Run 'forge --help' for usage.
    forge_info_commands_cli_test.go:116: 
        	Error Trace:	/Users/fanhuifeng/Projects/ai/coding-harness/forge/forge-cli/tests/e2e/features/forge-info-commands/forge_info_commands_cli_test.go:116
        	Error:      	Should be true
        	Test:       	TestTC_009_ProposalSlugDetailView
        	Messages:   	proposal detail should show FILE path: Error: unknown command "proposal" for "forge"
        	            	Run 'forge --help' for usage.
--- FAIL: TestTC_009_ProposalSlugDetailView (0.00s)

```

### TC-010: TestTC_010_ProposalCreatedDateFromFrontmatter

```
=== RUN   TestTC_010_ProposalCreatedDateFromFrontmatter
    forge_info_commands_cli_test.go:124: 
        	Error Trace:	/Users/fanhuifeng/Projects/ai/coding-harness/forge/forge-cli/tests/e2e/features/forge-info-commands/forge_info_commands_cli_test.go:124
        	Error:      	Not equal: 
        	            	expected: 0
        	            	actual  : 1
        	Test:       	TestTC_010_ProposalCreatedDateFromFrontmatter
        	Messages:   	proposal list should exit 0
    forge_info_commands_cli_test.go:126: 
        	Error Trace:	/Users/fanhuifeng/Projects/ai/coding-harness/forge/forge-cli/tests/e2e/features/forge-info-commands/forge_info_commands_cli_test.go:126
        	Error:      	Should be true
        	Test:       	TestTC_010_ProposalCreatedDateFromFrontmatter
        	Messages:   	proposal list should show created date '2026-05-14' from frontmatter: Error: unknown command "proposal" for "forge"
        	            	Run 'forge --help' for usage.
--- FAIL: TestTC_010_ProposalCreatedDateFromFrontmatter (0.01s)

```

### TC-011: TestTC_011_ProposalPRDColumnChecksPrdSpec

```
=== RUN   TestTC_011_ProposalPRDColumnChecksPrdSpec
    forge_info_commands_cli_test.go:134: 
        	Error Trace:	/Users/fanhuifeng/Projects/ai/coding-harness/forge/forge-cli/tests/e2e/features/forge-info-commands/forge_info_commands_cli_test.go:134
        	Error:      	Not equal: 
        	            	expected: 0
        	            	actual  : 1
        	Test:       	TestTC_011_ProposalPRDColumnChecksPrdSpec
        	Messages:   	proposal list should exit 0
    forge_info_commands_cli_test.go:148: 
        	Error Trace:	/Users/fanhuifeng/Projects/ai/coding-harness/forge/forge-cli/tests/e2e/features/forge-info-commands/forge_info_commands_cli_test.go:148
        	Error:      	Should be true
        	Test:       	TestTC_011_ProposalPRDColumnChecksPrdSpec
        	Messages:   	should find forge-info-commands row in proposal list: Error: unknown command "proposal" for "forge"
        	            	Run 'forge --help' for usage.
--- FAIL: TestTC_011_ProposalPRDColumnChecksPrdSpec (0.00s)

```

### TC-012: TestTC_012_ProposalFeatureColumnReadsManifestStatus

```
=== RUN   TestTC_012_ProposalFeatureColumnReadsManifestStatus
    forge_info_commands_cli_test.go:155: 
        	Error Trace:	/Users/fanhuifeng/Projects/ai/coding-harness/forge/forge-cli/tests/e2e/features/forge-info-commands/forge_info_commands_cli_test.go:155
        	Error:      	Not equal: 
        	            	expected: 0
        	            	actual  : 1
        	Test:       	TestTC_012_ProposalFeatureColumnReadsManifestStatus
        	Messages:   	proposal list should exit 0
    forge_info_commands_cli_test.go:168: 
        	Error Trace:	/Users/fanhuifeng/Projects/ai/coding-harness/forge/forge-cli/tests/e2e/features/forge-info-commands/forge_info_commands_cli_test.go:168
        	Error:      	Should be true
        	Test:       	TestTC_012_ProposalFeatureColumnReadsManifestStatus
        	Messages:   	should find forge-info-commands row in proposal list: Error: unknown command "proposal" for "forge"
        	            	Run 'forge --help' for usage.
--- FAIL: TestTC_012_ProposalFeatureColumnReadsManifestStatus (0.00s)

```

### TC-013: TestTC_013_FeatureListAllFeatures

```
=== RUN   TestTC_013_FeatureListAllFeatures
    forge_info_commands_cli_test.go:185: 
        	Error Trace:	/Users/fanhuifeng/Projects/ai/coding-harness/forge/forge-cli/tests/e2e/features/forge-info-commands/forge_info_commands_cli_test.go:185
        	Error:      	Should be true
        	Test:       	TestTC_013_FeatureListAllFeatures
        	Messages:   	feature list should contain SLUG column header: ---
        	            	FEATURE: list
        	            	---
    forge_info_commands_cli_test.go:187: 
        	Error Trace:	/Users/fanhuifeng/Projects/ai/coding-harness/forge/forge-cli/tests/e2e/features/forge-info-commands/forge_info_commands_cli_test.go:187
        	Error:      	Should be true
        	Test:       	TestTC_013_FeatureListAllFeatures
        	Messages:   	feature list should contain STATUS column header: ---
        	            	FEATURE: list
        	            	---
    forge_info_commands_cli_test.go:189: 
        	Error Trace:	/Users/fanhuifeng/Projects/ai/coding-harness/forge/forge-cli/tests/e2e/features/forge-info-commands/forge_info_commands_cli_test.go:189
        	Error:      	Should be true
        	Test:       	TestTC_013_FeatureListAllFeatures
        	Messages:   	feature list should contain PROGRESS column header: ---
        	            	FEATURE: list
        	            	---
--- FAIL: TestTC_013_FeatureListAllFeatures (0.01s)

```

### TC-015: TestTC_015_FeatureListScoresFromFrontmatter

```
=== RUN   TestTC_015_FeatureListScoresFromFrontmatter
    forge_info_commands_cli_test.go:205: 
        	Error Trace:	/Users/fanhuifeng/Projects/ai/coding-harness/forge/forge-cli/tests/e2e/features/forge-info-commands/forge_info_commands_cli_test.go:205
        	Error:      	Should be true
        	Test:       	TestTC_015_FeatureListScoresFromFrontmatter
        	Messages:   	feature list should show em-dash for missing scores: ---
        	            	FEATURE: list
        	            	---
--- FAIL: TestTC_015_FeatureListScoresFromFrontmatter (0.01s)

```

### TC-016: TestTC_016_FeatureStatusDetailView

```
=== RUN   TestTC_016_FeatureStatusDetailView
    forge_info_commands_cli_test.go:213: 
        	Error Trace:	/Users/fanhuifeng/Projects/ai/coding-harness/forge/forge-cli/tests/e2e/features/forge-info-commands/forge_info_commands_cli_test.go:213
        	Error:      	Not equal: 
        	            	expected: 0
        	            	actual  : 1
        	Test:       	TestTC_016_FeatureStatusDetailView
        	Messages:   	feature status should exit 0
    forge_info_commands_cli_test.go:220: 
        	Error Trace:	/Users/fanhuifeng/Projects/ai/coding-harness/forge/forge-cli/tests/e2e/features/forge-info-commands/forge_info_commands_cli_test.go:220
        	Error:      	Should be true
        	Test:       	TestTC_016_FeatureStatusDetailView
        	Messages:   	feature status should show STATUS field: Error: accepts at most 1 arg(s), received 2
        	            	Usage:
        	            	  forge feature [slug] [flags]
        	            	
        	            	Flags:
        	            	  -h, --help   help for feature
        	            	
    forge_info_commands_cli_test.go:222: 
        	Error Trace:	/Users/fanhuifeng/Projects/ai/coding-harness/forge/forge-cli/tests/e2e/features/forge-info-commands/forge_info_commands_cli_test.go:222
        	Error:      	Should be true
        	Test:       	TestTC_016_FeatureStatusDetailView
        	Messages:   	feature status should show TASKS section: Error: accepts at most 1 arg(s), received 2
        	            	Usage:
        	            	  forge feature [slug] [flags]
        	            	
        	            	Flags:
        	            	  -h, --help   help for feature
        	            	
--- FAIL: TestTC_016_FeatureStatusDetailView (0.00s)

```

### TC-018: TestTC_018_LessonListAllLessons

```
=== RUN   TestTC_018_LessonListAllLessons
    forge_info_commands_cli_test.go:243: 
        	Error Trace:	/Users/fanhuifeng/Projects/ai/coding-harness/forge/forge-cli/tests/e2e/features/forge-info-commands/forge_info_commands_cli_test.go:243
        	Error:      	Not equal: 
        	            	expected: 0
        	            	actual  : 1
        	Test:       	TestTC_018_LessonListAllLessons
        	Messages:   	lesson list should exit 0
    forge_info_commands_cli_test.go:249: 
        	Error Trace:	/Users/fanhuifeng/Projects/ai/coding-harness/forge/forge-cli/tests/e2e/features/forge-info-commands/forge_info_commands_cli_test.go:249
        	Error:      	Should be true
        	Test:       	TestTC_018_LessonListAllLessons
        	Messages:   	lesson list should contain NAME column header: Error: unknown command "lesson" for "forge"
        	            	Run 'forge --help' for usage.
    forge_info_commands_cli_test.go:251: 
        	Error Trace:	/Users/fanhuifeng/Projects/ai/coding-harness/forge/forge-cli/tests/e2e/features/forge-info-commands/forge_info_commands_cli_test.go:251
        	Error:      	Should be true
        	Test:       	TestTC_018_LessonListAllLessons
        	Messages:   	lesson list should contain CREATED column header: Error: unknown command "lesson" for "forge"
        	            	Run 'forge --help' for usage.
    forge_info_commands_cli_test.go:253: 
        	Error Trace:	/Users/fanhuifeng/Projects/ai/coding-harness/forge/forge-cli/tests/e2e/features/forge-info-commands/forge_info_commands_cli_test.go:253
        	Error:      	Should be true
        	Test:       	TestTC_018_LessonListAllLessons
        	Messages:   	lesson list should contain CATEGORY column header: Error: unknown command "lesson" for "forge"
        	            	Run 'forge --help' for usage.
    forge_info_commands_cli_test.go:255: 
        	Error Trace:	/Users/fanhuifeng/Projects/ai/coding-harness/forge/forge-cli/tests/e2e/features/forge-info-commands/forge_info_commands_cli_test.go:255
        	Error:      	Should be true
        	Test:       	TestTC_018_LessonListAllLessons
        	Messages:   	lesson list should contain TAGS column header: Error: unknown command "lesson" for "forge"
        	            	Run 'forge --help' for usage.
--- FAIL: TestTC_018_LessonListAllLessons (0.01s)

```

### TC-019: TestTC_019_LessonCategoryFromFilenamePrefix

```
=== RUN   TestTC_019_LessonCategoryFromFilenamePrefix
    forge_info_commands_cli_test.go:263: 
        	Error Trace:	/Users/fanhuifeng/Projects/ai/coding-harness/forge/forge-cli/tests/e2e/features/forge-info-commands/forge_info_commands_cli_test.go:263
        	Error:      	Not equal: 
        	            	expected: 0
        	            	actual  : 1
        	Test:       	TestTC_019_LessonCategoryFromFilenamePrefix
        	Messages:   	lesson list should exit 0
    forge_info_commands_cli_test.go:274: 
        	Error Trace:	/Users/fanhuifeng/Projects/ai/coding-harness/forge/forge-cli/tests/e2e/features/forge-info-commands/forge_info_commands_cli_test.go:274
        	Error:      	Should be true
        	Test:       	TestTC_019_LessonCategoryFromFilenamePrefix
        	Messages:   	lesson list should show 'gotcha' category for gotcha-* lesson files: Error: unknown command "lesson" for "forge"
        	            	Run 'forge --help' for usage.
--- FAIL: TestTC_019_LessonCategoryFromFilenamePrefix (0.01s)

```

### TC-020: TestTC_020_LessonNameDetailView

```
=== RUN   TestTC_020_LessonNameDetailView
    forge_info_commands_cli_test.go:305: 
        	Error Trace:	/Users/fanhuifeng/Projects/ai/coding-harness/forge/forge-cli/tests/e2e/features/forge-info-commands/forge_info_commands_cli_test.go:305
        	Error:      	Not equal: 
        	            	expected: 0
        	            	actual  : 1
        	Test:       	TestTC_020_LessonNameDetailView
        	Messages:   	lesson detail should exit 0
    forge_info_commands_cli_test.go:310: 
        	Error Trace:	/Users/fanhuifeng/Projects/ai/coding-harness/forge/forge-cli/tests/e2e/features/forge-info-commands/forge_info_commands_cli_test.go:310
        	Error:      	Should be true
        	Test:       	TestTC_020_LessonNameDetailView
        	Messages:   	lesson detail should show NAME field: Error: unknown command "lesson" for "forge"
        	            	Run 'forge --help' for usage.
    forge_info_commands_cli_test.go:312: 
        	Error Trace:	/Users/fanhuifeng/Projects/ai/coding-harness/forge/forge-cli/tests/e2e/features/forge-info-commands/forge_info_commands_cli_test.go:312
        	Error:      	Should be true
        	Test:       	TestTC_020_LessonNameDetailView
        	Messages:   	lesson detail should show FILE path: Error: unknown command "lesson" for "forge"
        	            	Run 'forge --help' for usage.
--- FAIL: TestTC_020_LessonNameDetailView (0.01s)

```

---

## Screenshots

No screenshots (CLI tests only).

---

## Failure Diagnosis

**Failure rate**: 14/31 non-skipped tests failed (45%).

### Root Cause Analysis

>30% failure rate detected. Investigating as potential infrastructure/app health issue.

All 14 failures are caused by **unimplemented CLI commands** -- not infrastructure problems:

1. **`config` subcommand (3 failures)**: `forge config get` does not exist yet. Error: `unknown command "config" for "forge"`.
   - TC-004: ConfigGetProjectType
   - TC-005: ConfigGetCapabilitiesArrayOutput
   - TC-006: ConfigGetMissingKey

2. **`proposal` subcommand (5 failures)**: `forge proposal` does not exist yet. Error: `unknown command "proposal" for "forge"`.
   - TC-008: ProposalListAllProposals
   - TC-009: ProposalSlugDetailView
   - TC-010: ProposalCreatedDateFromFrontmatter
   - TC-011: ProposalPRDColumnChecksPrdSpec
   - TC-012: ProposalFeatureColumnReadsManifestStatus

3. **`lesson` subcommand (3 failures)**: `forge lesson` does not exist yet. Error: `unknown command "lesson" for "forge"`.
   - TC-018: LessonListAllLessons
   - TC-019: LessonCategoryFromFilenamePrefix
   - TC-020: LessonNameDetailView

4. **`feature list` output format mismatch (2 failures)**: `forge feature list` outputs `FEATURE: list` (slug resolution format) instead of a table with column headers.
   - TC-013: FeatureListAllFeatures
   - TC-015: FeatureListScoresFromFrontmatter

5. **`feature status` subcommand (1 failure)**: `forge feature status <slug>` is not implemented -- `feature` accepts at most 1 arg.
   - TC-016: FeatureStatusDetailView

### Conclusion

No infrastructure issues. All failures are expected -- the test scripts exercise commands defined by the forge-info-commands feature proposal that are not yet implemented in the CLI binary.
