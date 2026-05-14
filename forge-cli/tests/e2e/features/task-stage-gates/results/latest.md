# E2E Test Report: task-stage-gates

**Date**: 2026-05-14
**Duration**: 1.926s

## Summary

| Type  | Total | Pass | Fail | Skip |
|-------|-------|------|------|------|
| UI    | 0     | 0    | 0    | 0    |
| API   | 0     | 0    | 0    | 0    |
| CLI   | 20    | 19   | 0    | 1    |
| **All** | **20** | **19** | **0** | **1** |

**Result**: PASS

---

## Results by Test Case

| TC ID | Test Name | Type | Status | Duration |
|-------|-----------|------|--------|----------|
| TC-001 | GeneratesSummaryAndGateForQualifyingPhases | CLI | PASS | 0.84s |
| TC-002 | CorrectDependencyWiringForGateTasks | CLI | PASS | 0.04s |
| TC-003 | SkipsSingleTaskPhases | CLI | PASS | 0.03s |
| TC-004 | ExcludesTestOnlyPhases | CLI | PASS | 0.03s |
| TC-005 | FiltersTestTasksFromBusinessCount | CLI | PASS | 0.02s |
| TC-006 | IdempotentRerunPreservesExistingFiles | CLI | PASS | 0.05s |
| TC-007 | GeneratesOnlyMissingGateWhenSummaryExists | CLI | PASS | 0.03s |
| TC-008 | GeneratesOnlyMissingSummaryWhenGateExists | CLI | PASS | 0.02s |
| TC-009 | SilentlySkipsMalformedTaskIDs | CLI | PASS | 0.02s |
| TC-010 | PreservesPreexistingHandCraftedGateFiles | CLI | PASS | 0.01s |
| TC-011 | GeneratedTasksInIndexJsonWithCorrectType | CLI | PASS | 0.02s |
| TC-012 | PrintsSummaryLinePerQualifyingPhase | CLI | PASS | 0.02s |
| TC-013 | PrintsNoQualificationMessage | CLI | PASS | 0.01s |
| TC-014 | ExitsWithErrorOnTemplateRenderFailure | CLI | SKIP | 0.00s |
| TC-015 | QuickModeGeneratesStageGatesIdentically | CLI | PASS | 0.02s |
| TC-016 | DoesNotBreakExistingIndexBehavior | CLI | PASS | 0.01s |
| TC-017 | NoTestFlagDoesNotAffectStageGates | CLI | PASS | 0.02s |
| TC-018 | ConcurrentExecutionIdenticalOutput | CLI | PASS | 0.02s |
| TC-019 | RejectsPathTraversalInTaskIDs | CLI | PASS | 0.01s |
| TC-020 | GenerationCompletesWithinTimeBudget | CLI | PASS | 0.11s |

---

## Failed Tests Detail

No failed tests.

---

## Skipped Tests Detail

| TC ID | Reason |
|-------|--------|
| TC-014 | Requires binary modification to corrupt embedded template - unit test scenario |

---

## Screenshots

No screenshots (CLI tests only).
