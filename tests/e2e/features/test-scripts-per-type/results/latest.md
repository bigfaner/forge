# E2E Test Report: test-scripts-per-type

**Date**: 2026-05-15
**Duration**: 1.238s

## Summary

| Type  | Total | Pass | Fail | Skip |
|-------|-------|------|------|------|
| UI    | 0     | 0    | 0    | 0    |
| API   | 0     | 0    | 0    | 0    |
| CLI   | 12    | 0    | 12   | 0    |
| **All** | **12** | **0** | **12** | **0** |

**Result**: FAIL

---

## Results by Test Case

| TC ID | Test Name | Type | Status | Duration |
|-------|-----------|------|--------|----------|
| TC-001 | TestTC_001_GenTestScriptsTypeFilterCLI | CLI | FAIL | 0.06s |
| TC-002 | TestTC_002_GenTestScriptsTypeFilterAPI | CLI | FAIL | 0.05s |
| TC-003 | TestTC_003_GenTestScriptsTypeFilterTUI | CLI | FAIL | 0.05s |
| TC-004 | TestTC_004_BreakdownTasksCreatesPerTypeTasks | CLI | FAIL | 0.05s |
| TC-005 | TestTC_005_BreakdownTasksNoEmptyTypeTasks | CLI | FAIL | 0.05s |
| TC-006 | TestTC_006_QuickTasksCreatesPerTypeTasks | CLI | FAIL | 0.05s |
| TC-007 | TestTC_007_Test3DependsOnAllPerTypeTasks | CLI | FAIL | 0.05s |
| TC-008 | TestTC_008_Quick3DependsOnAllPerTypeTasks | CLI | FAIL | 0.05s |
| TC-009 | TestTC_009_FailedGenTaskIndependentRetry | CLI | FAIL | 0.05s |
| TC-010 | TestTC_010_SharedInfrastructureIdempotent | CLI | FAIL | 0.05s |
| TC-011 | TestTC_011_SingleTypeProjectCreatesOnlyOneGenTask | CLI | FAIL | 0.05s |
| TC-012 | TestTC_012_GenTestScriptsWithoutTypeGeneratesAllTypes | CLI | FAIL | 0.05s |

---

## Failed Tests Detail

### 100% Failure Rate -- Infrastructure Problem

**Root Cause**: All 12 tests fail with the same systemic error:

```
Error: unknown command "gen-test-scripts" for "forge"
Error: unknown command "breakdown-tasks" for "forge"
Error: unknown command "quick-tasks" for "forge"
```

**Diagnosis**: The `forge` CLI (v3.11.0) does not implement the commands these tests invoke:

- `forge gen-test-scripts --type <type> --feature <slug>` -- command does not exist
- `forge breakdown-tasks --feature <slug>` -- command does not exist
- `forge quick-tasks --feature <slug>` -- command does not exist

The available `forge e2e` subcommands (`compile`, `discover`, `run`, `setup`, `validate-specs`, `verify`) do not match what these tests expect.

These test scripts validate a **proposed feature** (per-type test script generation) that has not been implemented in the forge CLI yet. The test scripts are correct in structure but the product under test lacks the required functionality.

**Resolution**: The `forge` CLI must implement the `gen-test-scripts`, `breakdown-tasks`, and `quick-tasks` commands with the `--type` flag before these tests can pass.

---

## Screenshots

No screenshots (CLI-only test suite).
