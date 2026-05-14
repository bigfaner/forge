# E2E Test Report: forge-cli-v3

**Date**: 2026-05-14
**Duration**: 1.605s
**Profile**: go-test

## Summary

| Type  | Total | Pass | Fail | Skip |
|-------|-------|------|------|------|
| UI    | 0     | 0    | 0    | 0    |
| API   | 0     | 0    | 0    | 0    |
| CLI   | 41    | 17   | 0    | 24   |
| **All** | **41** | **17** | **0** | **24** |

**Result**: PASS (17/17 executed tests passed; 24 skipped due to manual setup requirements)

---

## Results by Test Case

### discovery_cli_test.go (Discovery)

| TC ID | Test Name | Status | Duration | Notes |
|-------|-----------|--------|----------|-------|
| TC-001 | HelpOutputShowsCommandGroups | PASS | 0.15s | |
| TC-002 | TaskSubcommandHelpShowsAllCommands | PASS | 0.03s | |
| TC-003 | UnknownCommandReturnsErrorWithSuggestion | PASS | 0.03s | |
| TC-004 | UnknownTaskSubcommandReturnsErrorWithList | PASS | 0.03s | |

### prompt_cli_test.go (Prompt)

| TC ID | Test Name | Status | Duration | Notes |
|-------|-----------|--------|----------|-------|
| TC-005 | GetPromptByTaskIDReturnsCorrectPrompt | PASS | 0.09s | |
| TC-006 | GetPromptNonexistentTaskIDReturnsError | PASS | 0.09s | |
| TC-007 | GetPromptMissingOrInvalidTypeReturnsError | SKIP | 0.00s | requires manual setup: task with missing/invalid type in index.json |

### submit_cli_test.go (Submit)

| TC ID | Test Name | Status | Duration | Notes |
|-------|-----------|--------|----------|-------|
| TC-008 | SubmitTaskSuccessUpdatesStatusAndCreatesRecord | SKIP | 0.00s | requires manual setup: task in in_progress state with valid test data |
| TC-009 | SubmitTaskAlreadyTerminalStateReturnsError | SKIP | 0.00s | requires manual setup: task in terminal state |
| TC-010 | SubmitTaskMissingResultFlagReturnsError | PASS | 0.09s | |
| TC-011 | ConcurrentSubmitHandlesLockContention | SKIP | 0.00s | requires manual setup: concurrent process simulation |

### lifecycle_cli_test.go (Lifecycle)

| TC ID | Test Name | Status | Duration | Notes |
|-------|-----------|--------|----------|-------|
| TC-012 | CleanupRemovesTerminalStateFiles | SKIP | 0.00s | requires manual setup: feature with terminal-state tasks and state.json |
| TC-013 | QualityGateRunsCompileFmtLintTestSequence | SKIP | 0.00s | requires manual setup: all tasks completed with compilable project |
| TC-014 | CleanupNoTerminalTasksOutputsMessage | SKIP | 0.00s | requires manual setup: feature with no terminal-state tasks |
| TC-015 | QualityGateCreatesNewFixTaskOnRepeatedFailure | SKIP | 0.00s | requires manual setup: failing quality gate with existing fix-task |
| TC-016 | QualityGateStopsCreatingFixTasksAfterMax3 | SKIP | 0.00s | requires manual setup: 3 existing fix-tasks for same step |

### e2e_cli_test.go (E2E Pipeline)

| TC ID | Test Name | Status | Duration | Notes |
|-------|-----------|--------|----------|-------|
| TC-017 | E2ERunWithConfiguredProfileExecutesSuite | SKIP | 0.00s | requires manual setup: .forge/config.yaml with valid profile and feature test data |
| TC-018 | E2ERunNoProfileConfiguredReturnsError | SKIP | 0.00s | requires manual setup: .forge/config.yaml with no profile field |
| TC-019 | E2ERunUnknownProfileReturnsErrorWithList | SKIP | 0.00s | requires manual setup: .forge/config.yaml with unknown profile value |
| TC-020 | E2ERunNonexistentFeatureReturnsError | PASS | 0.10s | |

### task_types_cli_test.go (Task Types)

| TC ID | Test Name | Status | Duration | Notes |
|-------|-----------|--------|----------|-------|
| TC-021 | ListTypesOutputsAllWithDescriptions | PASS | 0.03s | |
| TC-022 | ListTypesEmptyRegistryReturnsEmpty | SKIP | 0.00s | requires manual setup: empty task type registry |

### forensic_cli_test.go (Forensic)

| TC ID | Test Name | Status | Duration | Notes |
|-------|-----------|--------|----------|-------|
| TC-023 | ForensicSearchScansHistoryAndReturnsSessions | PASS | 0.08s | |
| TC-024 | ForensicExtractOutputsEvidenceSummary | SKIP | 0.00s | requires manual setup: valid session JSONL file path |
| TC-025 | ForensicSubagentsListsTranscripts | SKIP | 0.00s | requires manual setup: session directory with subagent transcripts |
| TC-026 | ForensicExtractNonexistentPathReturnsError | PASS | 0.04s | |

### error_handling_cli_test.go (Error Handling)

| TC ID | Test Name | Status | Duration | Notes |
|-------|-----------|--------|----------|-------|
| TC-031 | TaskClaimNoAvailableTasksReturnsError | SKIP | 0.00s | requires manual setup: feature with no available tasks to claim |
| TC-032 | TaskClaimCorruptedIndexReturnsError | SKIP | 0.00s | requires manual setup: corrupted or missing index.json |
| TC-033 | TaskCheckDepsUnmetDependencyReturnsError | SKIP | 0.00s | requires manual setup: task with unmet dependency in index.json |
| TC-034 | TaskValidateIndexInvalidSchemaReturnsError | SKIP | 0.00s | requires manual setup: index.json with schema validation errors |
| TC-035 | TaskStatusNonexistentIDReturnsError | PASS | 0.09s | |
| TC-036 | ForensicSearchNoResultsReturnsEmptyOutput | PASS | 0.04s | |
| TC-037 | ForensicSearchMissingRecordsDirReturnsError | SKIP | 0.00s | requires manual setup: missing ~/.claude/history.jsonl |
| TC-038 | VerifyTaskDoneIncompleteTasksReturnsError | SKIP | 0.00s | requires manual setup: feature with incomplete tasks and active state.json |
| TC-039 | TaskSubmitConcurrentWriteConflictReturnsRetryError | SKIP | 0.00s | requires manual setup: concurrent lock contention scenario |
| TC-040 | TaskSubmitMissingIndexReturnsError | SKIP | 0.00s | requires manual setup: feature directory without index.json |
| TC-041 | ProfileGetInvalidProfileReturnsErrorWithList | PASS | 0.03s | |

### profile_cli_test.go (Profile)

| TC ID | Test Name | Status | Duration | Notes |
|-------|-----------|--------|----------|-------|
| TC-027 | ProfileDetectScansAndOutputsProfiles | PASS | 0.03s | |
| TC-028 | ProfileSetUpdatesConfigWithValidProfile | SKIP | 0.00s | requires manual setup: writable .forge/config.yaml, verify config restoration |
| TC-029 | ProfileGetOutputsStrategyFileContent | PASS | 0.03s | |
| TC-030 | ProfileSetInvalidProfileReturnsErrorWithList | PASS | 0.03s | |

---

## Failed Tests Detail

No test failures.

---

## Skipped Tests Analysis

24 tests were skipped due to manual setup requirements. These tests require specific filesystem state (e.g., corrupted index.json, in_progress task states, writable config files) that must be provisioned before execution. These are not failures -- they are intentionally guarded with setup preconditions.

Categories of skipped tests:
- **Lifecycle tests** (TC-012 to TC-016): Require feature directories with specific task states and state.json files
- **Submit tests** (TC-008, TC-009, TC-011): Require tasks in specific states (in_progress, terminal) with lock contention scenarios
- **E2E pipeline tests** (TC-017 to TC-019): Require specific .forge/config.yaml configurations
- **Error handling tests** (TC-031 to TC-034, TC-037 to TC-040): Require specific error-state filesystem setups
- **Other setup-dependent tests** (TC-007, TC-022, TC-024, TC-025, TC-028): Various manual preconditions

---

## Screenshots

No screenshots (CLI tests only -- no UI tests in this suite).
