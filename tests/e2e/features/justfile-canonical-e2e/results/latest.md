# E2E Test Report: justfile-canonical-e2e

**Date**: 2026-05-15
**Duration**: 1.317s

## Summary

| Type  | Total | Pass | Fail | Skip |
|-------|-------|------|------|------|
| UI    | 0     | 0    | 0    | 0    |
| API   | 0     | 0    | 0    | 0    |
| CLI   | 20    | 15   | 5    | 0    |
| **All** | **20** | **15** | **5** | **0** |

**Result**: FAIL (25% failure rate)

---

## Results by Test Case

| TC ID   | Test Name                                                    | Type | Status | Duration  |
|---------|--------------------------------------------------------------|------|--------|-----------|
| TC-001  | RunDelegatesToJustTestE2e                                    | CLI  | FAIL   | 0.05s     |
| TC-002  | RunPassesFeatureAsJustfileArgument                           | CLI  | FAIL   | 0.05s     |
| TC-003  | SetupDelegatesToJustE2eSetup                                 | CLI  | FAIL   | 0.05s     |
| TC-004  | CompileDelegatesToJustE2eCompile                             | CLI  | FAIL   | 0.05s     |
| TC-005  | DiscoverDelegatesToJustE2eDiscover                           | CLI  | FAIL   | 0.05s     |
| TC-006  | VerifyDetectsUnresolvedMarkers                               | CLI  | PASS   | 0.00s     |
| TC-007  | VerifyPassesWhenNoUnresolvedMarkers                          | CLI  | PASS   | 0.00s     |
| TC-008  | GraduateMergesFeatureTestIntoRegression                      | CLI  | PASS   | 0.04s     |
| TC-009  | RunDetectsMissingJustfileGracefully                          | CLI  | PASS   | 0.04s     |
| TC-010  | SetupDetectsMissingJustfileGracefully                        | CLI  | PASS   | 0.04s     |
| TC-011  | CompileDetectsMissingJustfileGracefully                      | CLI  | PASS   | 0.04s     |
| TC-012  | DiscoverDetectsMissingJustfileGracefully                     | CLI  | PASS   | 0.04s     |
| TC-013  | VerifyDetectsMissingJustfileGracefully                       | CLI  | PASS   | 0.04s     |
| TC-014  | GraduateDetectsMissingJustfileGracefully                     | CLI  | PASS   | 0.04s     |
| TC-015  | RunErrorsWhenFeatureDirMissing                               | CLI  | PASS   | 0.03s     |
| TC-016  | SetupErrorsWhenFeatureDirMissing                             | CLI  | PASS   | 0.03s     |
| TC-017  | AllManifestsContainRequiredFields                            | CLI  | PASS   | 0.00s     |
| TC-018  | AllManifestsHaveValidFileExtension                           | CLI  | PASS   | 0.00s     |
| TC-019  | AllManifestsContainCapabilitiesSubset                        | CLI  | PASS   | 0.00s     |
| TC-020  | AllManifestsContainZeroRunAndGraduateFields                  | CLI  | PASS   | 0.00s     |

---

## Failed Tests Detail

### TC-001: RunDelegatesToJustTestE2e
**Error**: `failed to build forge binary: exit status 1: no Go files in Z:\project\ai\coding-harness\forge-3\forge-cli`
**Root cause**: The `forgeBinary()` helper runs `go build -o <bin> ./` from `forge-cli/` directory, but there is no Go package at the root of `forge-cli/`. The main entry point is under `forge-cli/cmd/`.

### TC-002: RunPassesFeatureAsJustfileArgument
**Error**: `failed to build forge binary: exit status 1: no Go files in Z:\project\ai\coding-harness\forge-3\forge-cli`
**Root cause**: Same as TC-001 -- forge binary cannot be built.

### TC-003: SetupDelegatesToJustE2eSetup
**Error**: `failed to build forge binary: exit status 1: no Go files in Z:\project\ai\coding-harness\forge-3\forge-cli`
**Root cause**: Same as TC-001 -- forge binary cannot be built.

### TC-004: CompileDelegatesToJustE2eCompile
**Error**: `failed to build forge binary: exit status 1: no Go files in Z:\project\ai\coding-harness\forge-3\forge-cli`
**Root cause**: Same as TC-001 -- forge binary cannot be built.

### TC-005: DiscoverDelegatesToJustE2eDiscover
**Error**: `failed to build forge binary: exit status 1: no Go files in Z:\project\ai\coding-harness\forge-3\forge-cli`
**Root cause**: Same as TC-001 -- forge binary cannot be built.

### Diagnosis

**Failure pattern**: All 5 failing tests (TC-001 through TC-005) share the same root cause: the `forgeBinary()` helper in `helpers_test.go` attempts `go build ./` from the `forge-cli/` directory, but the Go module root at `forge-cli/` has no `.go` files -- the entry point is under `forge-cli/cmd/`. This is a shared infrastructure issue, not per-test logic bugs.

**Impact**: 25% failure rate (5/20). Falls in the 10-30% range. Spot-checking confirms all failures are identical -- one shared helper function issue.

**Recommended fix**: Update `forgeBinary()` in `helpers_test.go` to build from the correct subdirectory (e.g., `go build -o <bin> ./cmd/...` or `go build -o <bin> .` from `forge-cli/cmd/forge/`).

---

## Screenshots

No screenshots (CLI tests only).
