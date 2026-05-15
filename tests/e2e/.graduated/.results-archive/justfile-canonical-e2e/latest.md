# E2E Test Report: justfile-canonical-e2e

**Date**: 2026-05-15
**Duration**: 1.465s
**Profile**: go-test

## Summary

| Type  | Total | Pass | Fail | Skip |
|-------|-------|------|------|------|
| UI    | 0     | 0    | 0    | 0    |
| API   | 0     | 0    | 0    | 0    |
| CLI   | 20    | 20   | 0    | 0    |
| **All** | **20** | **20** | **0** | **0** |

**Result**: ALL PASSED

---

## Results by Test Case

| TC ID  | Test Name | Type | Status | Duration |
|--------|-----------|------|--------|----------|
| TC-001 | RunDelegatesToJustTestE2e | CLI | PASS | 0.070s |
| TC-002 | RunPassesFeatureAsJustfileArgument | CLI | PASS | 0.070s |
| TC-003 | SetupDelegatesToJustE2eSetup | CLI | PASS | 0.070s |
| TC-004 | CompileDelegatesToJustE2eCompile | CLI | PASS | 0.070s |
| TC-005 | DiscoverDelegatesToJustE2eDiscover | CLI | PASS | 0.070s |
| TC-006 | VerifyDoesNotDelegateToJust | CLI | PASS | 0.040s |
| TC-007 | VerifyFindsVerifyMarkersInTestFiles | CLI | PASS | 0.040s |
| TC-008 | JustNotOnPathReturnsActionableErrorForRun | CLI | PASS | 0.030s |
| TC-009 | JustNotOnPathReturnsActionableErrorForSetup | CLI | PASS | 0.030s |
| TC-010 | JustNotOnPathReturnsActionableErrorForCompile | CLI | PASS | 0.030s |
| TC-011 | JustNotOnPathReturnsActionableErrorForDiscover | CLI | PASS | 0.030s |
| TC-012 | NonZeroJustExitReturnsErrorForRun | CLI | PASS | 0.070s |
| TC-013 | ZeroJustExitReturnsNilErrorForRun | CLI | PASS | 0.070s |
| TC-014 | NonZeroJustExitReturnsErrorForCompile | CLI | PASS | 0.070s |
| TC-015 | NonZeroJustExitReturnsErrorForDiscover | CLI | PASS | 0.070s |
| TC-016 | NoProfileReturnsErrNoProfileForRun | CLI | PASS | 0.030s |
| TC-017 | NoProfileReturnsErrNoProfileForSetup | CLI | PASS | 0.030s |
| TC-018 | NoProfileReturnsErrNoProfileForCompile | CLI | PASS | 0.030s |
| TC-019 | NoProfileReturnsErrNoProfileForDiscover | CLI | PASS | 0.030s |
| TC-020 | AllManifestsContainZeroRunAndGraduateFields | CLI | PASS | 0.000s |

### Subtest Results (TC-020)

| Subtest | Status |
|---------|--------|
| TC-020/go-test | PASS |
| TC-020/java-junit | PASS |
| TC-020/maestro | PASS |
| TC-020/pytest | PASS |
| TC-020/rust-test | PASS |
| TC-020/web-playwright | PASS |

---

## Failed Tests Detail

No failures.

---

## Screenshots

No screenshots captured (CLI tests only).
