# E2E Test Report: test-scripts-per-type

**Date**: 2026-05-15
**Duration**: 1.80s

## Summary

| Type  | Total | Pass | Fail | Skip |
|-------|-------|------|------|------|
| UI    | 0   | 0   | 0   | 0   |
| API   | 0  | 0  | 0  | 0  |
| CLI   | 12  | 12  | 0  | 0  |
| **All** | **12** | **12** | **0** | **0** |

**Result**: PASS

---

## Results by Test Case

| TC ID | Test Name | Type | Status | Duration |
|-------|-----------|------|--------|----------|
| TC-001 | TaskIndexCreatesPerTypeTasksForMultiType | CLI | PASS | 0.26s |
| TC-002 | TaskIndexPerTypeTasksHaveCorrectType | CLI | PASS | 0.07s |
| TC-003 | TaskIndexSingleTypeCreatesOneGenTask | CLI | PASS | 0.07s |
| TC-004 | TaskIndexWithoutTestCasesFallsBackToLegacy | CLI | PASS | 0.06s |
| TC-005 | TaskIndexZeroTypeTestCasesFallsBackToLegacy | CLI | PASS | 0.07s |
| TC-006 | TaskIndexRunDependsOnAllPerTypeGenTasks | CLI | PASS | 0.06s |
| TC-007 | TaskIndexMultiProfilePerTypeTasks | CLI | PASS | 0.08s |
| TC-008 | TaskIndexQuickModePerTypeTasks | CLI | PASS | 0.07s |
| TC-009 | PerTypeGenScriptsMdContainsTestType | CLI | PASS | 0.06s |
| TC-010 | TaskIndexPerTypeIdempotent | CLI | PASS | 0.10s |
| TC-011 | PerTypeGenScriptsMdHasCorrectTaskIDs | CLI | PASS | 0.06s |
| TC-012 | TaskIndexSharedInfrastructureNotDuplicated | CLI | PASS | 0.11s |

---

## Failed Tests Detail

No failures.

---

## Screenshots

No screenshots (CLI tests only).

## Notes

- Fixed compilation error: unused variable `runTask` in TC-006 (line 346), changed to `_`.
- All 12 tests use `os/exec` / `exec.Command` patterns, classified as CLI type.
- Raw JSON output: `tests/e2e/features/test-scripts-per-type/results/test-output.json`
