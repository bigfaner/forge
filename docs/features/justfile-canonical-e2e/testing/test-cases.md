---
feature: "justfile-canonical-e2e"
sources:
  - docs/proposals/justfile-canonical-e2e/proposal.md
  - docs/features/justfile-canonical-e2e/tasks/1-remove-manifest-command-fields.md
  - docs/features/justfile-canonical-e2e/tasks/2-delegate-actions-to-just.md
  - docs/features/justfile-canonical-e2e/tasks/3-update-actions-tests.md
  - docs/features/justfile-canonical-e2e/tasks/4-version-bump.md
generated: "2026-05-15"
---

# Test Cases: justfile-canonical-e2e

## Summary

| Type | Count |
|------|-------|
| UI   | 0    |
| **Integration** | **0** |
| API  | 0    |
| CLI  | 20   |
| **Total** | **20** |

> **Note**: This feature is a CLI-only refactor. No UI or API interfaces are exposed by forge-cli. Profile capabilities `[tui, api, cli]` describe what profiles can test, not what forge itself exposes. The forge binary is a CLI tool.

---

## CLI Test Cases

### Command Delegation

## TC-001: Run delegates to just test-e2e
- **Source**: Proposal Success Criteria [1], Task 2 AC [1]
- **Type**: CLI
- **Target**: cli/e2e-run
- **Test ID**: cli/e2e-run/run-delegates-to-just-test-e2e
- **Pre-conditions**: A valid test profile is configured in `.forge/config.yaml`
- **Route**: `forge e2e run`
- **Element**: sitemap-missing
- **Steps**:
  1. Configure `go-test` profile in a temp project directory
  2. Stub `runner` to capture the executed command
  3. Call `Run(RunOpts{ProjectRoot: dir, Feature: ""})`
  4. Assert `runner.Run` was called with `("just", "test-e2e")`
- **Expected**: The function calls `exec.Command("just", "test-e2e")` with no additional arguments
- **Priority**: P0

## TC-002: Run passes feature as justfile argument
- **Source**: Proposal Success Criteria [1], Task 2 AC [1], Task 2 Hard Rules [3]
- **Type**: CLI
- **Target**: cli/e2e-run
- **Test ID**: cli/e2e-run/run-passes-feature-as-justfile-argument
- **Pre-conditions**: A valid test profile is configured; `Feature` field is non-empty
- **Route**: `forge e2e run`
- **Element**: sitemap-missing
- **Steps**:
  1. Configure `go-test` profile in a temp project directory
  2. Stub `runner` to capture the executed command
  3. Call `Run(RunOpts{ProjectRoot: dir, Feature: "my-feature"})`
  4. Assert `runner.Run` was called with `("just", "test-e2e", "feature=my-feature")`
- **Expected**: The feature name is appended as `feature=<name>` using just's named argument syntax
- **Priority**: P0

## TC-003: Setup delegates to just e2e-setup
- **Source**: Proposal Success Criteria [2], Task 2 AC [2]
- **Type**: CLI
- **Target**: cli/e2e-setup
- **Test ID**: cli/e2e-setup/setup-delegates-to-just-e2e-setup
- **Pre-conditions**: A valid test profile is configured
- **Route**: `forge e2e setup`
- **Element**: sitemap-missing
- **Steps**:
  1. Configure `go-test` profile in a temp project directory
  2. Stub `runner` to capture the executed command
  3. Call `Setup(RunOpts{ProjectRoot: dir})`
  4. Assert `runner.Run` was called with `("just", "e2e-setup")`
- **Expected**: The function calls `exec.Command("just", "e2e-setup")`
- **Priority**: P0

## TC-004: Compile delegates to just e2e-compile
- **Source**: Proposal Success Criteria [3], Task 2 AC [3]
- **Type**: CLI
- **Target**: cli/e2e-compile
- **Test ID**: cli/e2e-compile/compile-delegates-to-just-e2e-compile
- **Pre-conditions**: A valid test profile is configured
- **Route**: `forge e2e compile`
- **Element**: sitemap-missing
- **Steps**:
  1. Configure `go-test` profile in a temp project directory
  2. Stub `runner` to capture the executed command
  3. Call `Compile(dir)`
  4. Assert `runner.Run` was called with `("just", "e2e-compile")`
- **Expected**: The function calls `exec.Command("just", "e2e-compile")`
- **Priority**: P0

## TC-005: Discover delegates to just e2e-discover
- **Source**: Proposal Success Criteria [4], Task 2 AC [4]
- **Type**: CLI
- **Target**: cli/e2e-discover
- **Test ID**: cli/e2e-discover/discover-delegates-to-just-e2e-discover
- **Pre-conditions**: A valid test profile is configured
- **Route**: `forge e2e discover`
- **Element**: sitemap-missing
- **Steps**:
  1. Configure `go-test` profile in a temp project directory
  2. Stub `runner` to capture the executed command
  3. Call `Discover(dir)`
  4. Assert `runner.Run` was called with `("just", "e2e-discover")`
- **Expected**: The function calls `exec.Command("just", "e2e-discover")`
- **Priority**: P0

### Verify Unchanged

## TC-006: Verify does not delegate to just
- **Source**: Proposal Success Criteria [5], Task 2 AC [5], Task 2 Hard Rules [1]
- **Type**: CLI
- **Target**: cli/e2e-verify
- **Test ID**: cli/e2e-verify/verify-does-not-delegate-to-just
- **Pre-conditions**: A valid test profile is configured; e2e test directory exists with no VERIFY markers
- **Route**: `forge e2e verify`
- **Element**: sitemap-missing
- **Steps**:
  1. Configure `go-test` profile in a temp project directory
  2. Create `tests/e2e/` directory with a test file containing no VERIFY markers
  3. Call `Verify(RunOpts{ProjectRoot: dir, Feature: ""})`
  4. Assert no subprocess is spawned (runner is not invoked)
- **Expected**: Verify scans files locally using `filepath.WalkDir`; no `just` subprocess is called
- **Priority**: P0

## TC-007: Verify finds VERIFY markers in test files
- **Source**: Task 2 AC [5] (Verify unchanged behavior)
- **Type**: CLI
- **Target**: cli/e2e-verify
- **Test ID**: cli/e2e-verify/verify-finds-verify-markers-in-test-files
- **Pre-conditions**: A valid test profile is configured; e2e test directory exists with a file containing `VERIFY`
- **Route**: `forge e2e verify`
- **Element**: sitemap-missing
- **Steps**:
  1. Configure `go-test` profile in a temp project directory
  2. Create `tests/e2e/` directory with a file containing `// VERIFY: placeholder`
  3. Call `Verify(RunOpts{ProjectRoot: dir, Feature: ""})`
  4. Assert the error contains "VERIFY markers found"
- **Expected**: Error listing files with unresolved VERIFY markers
- **Priority**: P1

### Error Handling

## TC-008: Just not on PATH returns actionable error for Run
- **Source**: Proposal Success Criteria [7], Proposal Error Scenarios row 1, Task 2 AC [7]
- **Type**: CLI
- **Target**: cli/e2e-run
- **Test ID**: cli/e2e-run/just-not-on-path-returns-actionable-error-for-run
- **Pre-conditions**: A valid test profile is configured; `just` binary is not on PATH
- **Route**: `forge e2e run`
- **Element**: sitemap-missing
- **Steps**:
  1. Configure `go-test` profile in a temp project directory
  2. Stub `runner` to return an error containing "executable file not found" with "just" in the message
  3. Call `Run(RunOpts{ProjectRoot: dir})`
  4. Assert the error message contains `'just' is required but not found on PATH`
- **Expected**: Clear, actionable error message directing user to install `just`
- **Priority**: P0

## TC-009: Just not on PATH returns actionable error for Setup
- **Source**: Proposal Error Scenarios row 1, Task 2 AC [7]
- **Type**: CLI
- **Target**: cli/e2e-setup
- **Test ID**: cli/e2e-setup/just-not-on-path-returns-actionable-error-for-setup
- **Pre-conditions**: A valid test profile is configured; `just` binary is not on PATH
- **Route**: `forge e2e setup`
- **Element**: sitemap-missing
- **Steps**:
  1. Configure `go-test` profile in a temp project directory
  2. Stub `runner` to return an error containing "executable file not found" with "just" in the message
  3. Call `Setup(RunOpts{ProjectRoot: dir})`
  4. Assert the error message contains `'just' is required but not found on PATH`
- **Expected**: Same actionable error message for all just-delegating functions
- **Priority**: P0

## TC-010: Just not on PATH returns actionable error for Compile
- **Source**: Proposal Error Scenarios row 1, Task 2 AC [7]
- **Type**: CLI
- **Target**: cli/e2e-compile
- **Test ID**: cli/e2e-compile/just-not-on-path-returns-actionable-error-for-compile
- **Pre-conditions**: A valid test profile is configured; `just` binary is not on PATH
- **Route**: `forge e2e compile`
- **Element**: sitemap-missing
- **Steps**:
  1. Configure `go-test` profile in a temp project directory
  2. Stub `runner` to return an error containing "executable file not found" with "just" in the message
  3. Call `Compile(dir)`
  4. Assert the error message contains `'just' is required but not found on PATH`
- **Expected**: Same actionable error message across all delegating functions
- **Priority**: P0

## TC-011: Just not on PATH returns actionable error for Discover
- **Source**: Proposal Error Scenarios row 1, Task 2 AC [7]
- **Type**: CLI
- **Target**: cli/e2e-discover
- **Test ID**: cli/e2e-discover/just-not-on-path-returns-actionable-error-for-discover
- **Pre-conditions**: A valid test profile is configured; `just` binary is not on PATH
- **Route**: `forge e2e discover`
- **Element**: sitemap-missing
- **Steps**:
  1. Configure `go-test` profile in a temp project directory
  2. Stub `runner` to return an error containing "executable file not found" with "just" in the message
  3. Call `Discover(dir)`
  4. Assert the error message contains `'just' is required but not found on PATH`
- **Expected**: Same actionable error message across all delegating functions
- **Priority**: P0

### Exit Code Propagation

## TC-012: Non-zero just exit returns error for Run
- **Source**: Proposal Success Criteria [6], Proposal Error Scenarios row 3, Task 2 AC [6]
- **Type**: CLI
- **Target**: cli/e2e-run
- **Test ID**: cli/e2e-run/non-zero-just-exit-returns-error-for-run
- **Pre-conditions**: A valid test profile is configured; `just` exits with non-zero code
- **Route**: `forge e2e run`
- **Element**: sitemap-missing
- **Steps**:
  1. Configure `go-test` profile in a temp project directory
  2. Stub `runner` to return `fmt.Errorf("exit status 1")` with stderr output
  3. Call `Run(RunOpts{ProjectRoot: dir})`
  4. Assert the returned error is non-nil and contains the first line of stderr
- **Expected**: Non-nil error wrapping the just failure with formatted stderr context
- **Priority**: P0

## TC-013: Zero just exit returns nil error for Run
- **Source**: Proposal Success Criteria [6], Task 2 AC [6]
- **Type**: CLI
- **Target**: cli/e2e-run
- **Test ID**: cli/e2e-run/zero-just-exit-returns-nil-error-for-run
- **Pre-conditions**: A valid test profile is configured; `just` exits with zero code
- **Route**: `forge e2e run`
- **Element**: sitemap-missing
- **Steps**:
  1. Configure `go-test` profile in a temp project directory
  2. Stub `runner` to return `([]byte("ok\n"), nil)`
  3. Call `Run(RunOpts{ProjectRoot: dir})`
  4. Assert the returned error is nil
- **Expected**: nil error when just succeeds
- **Priority**: P0

## TC-014: Non-zero just exit returns error for Compile
- **Source**: Proposal Success Criteria [6], Task 2 AC [6]
- **Type**: CLI
- **Target**: cli/e2e-compile
- **Test ID**: cli/e2e-compile/non-zero-just-exit-returns-error-for-compile
- **Pre-conditions**: A valid test profile is configured; `just` exits with non-zero code
- **Route**: `forge e2e compile`
- **Element**: sitemap-missing
- **Steps**:
  1. Configure `go-test` profile in a temp project directory
  2. Stub `runner` to return `fmt.Errorf("exit status 1")` with stderr output containing a compilation error
  3. Call `Compile(dir)`
  4. Assert error is non-nil and contains "just e2e-compile failed:"
- **Expected**: Non-nil error with formatted compilation failure context
- **Priority**: P0

## TC-015: Non-zero just exit returns error for Discover
- **Source**: Proposal Success Criteria [6], Task 2 AC [6]
- **Type**: CLI
- **Target**: cli/e2e-discover
- **Test ID**: cli/e2e-discover/non-zero-just-exit-returns-error-for-discover
- **Pre-conditions**: A valid test profile is configured; `just` exits with non-zero code
- **Route**: `forge e2e discover`
- **Element**: sitemap-missing
- **Steps**:
  1. Configure `go-test` profile in a temp project directory
  2. Stub `runner` to return `fmt.Errorf("exit status 1")` with stderr output
  3. Call `Discover(dir)`
  4. Assert error is non-nil and contains "just e2e-discover failed:"
- **Expected**: Non-nil error with formatted discovery failure context
- **Priority**: P0

### Profile Resolution Errors

## TC-016: No profile returns ErrNoProfile for Run
- **Source**: Proposal Error Scenarios row 5, Task 2 AC [8]
- **Type**: CLI
- **Target**: cli/e2e-run
- **Test ID**: cli/e2e-run/no-profile-returns-err-no-profile-for-run
- **Pre-conditions**: No `.forge/config.yaml` exists; no profile configured
- **Route**: `forge e2e run`
- **Element**: sitemap-missing
- **Steps**:
  1. Create a temp directory with no `.forge/config.yaml`
  2. Call `Run(RunOpts{ProjectRoot: dir})`
  3. Assert the error wraps `ErrNoProfile`
- **Expected**: Error is `ErrNoProfile`; no subprocess is spawned
- **Priority**: P0

## TC-017: No profile returns ErrNoProfile for Setup
- **Source**: Proposal Error Scenarios row 5
- **Type**: CLI
- **Target**: cli/e2e-setup
- **Test ID**: cli/e2e-setup/no-profile-returns-err-no-profile-for-setup
- **Pre-conditions**: No `.forge/config.yaml` exists; no profile configured
- **Route**: `forge e2e setup`
- **Element**: sitemap-missing
- **Steps**:
  1. Create a temp directory with no `.forge/config.yaml`
  2. Call `Setup(RunOpts{ProjectRoot: dir})`
  3. Assert the error wraps `ErrNoProfile`
- **Expected**: Error is `ErrNoProfile`; no subprocess is spawned
- **Priority**: P0

## TC-018: No profile returns ErrNoProfile for Compile
- **Source**: Proposal Error Scenarios row 5
- **Type**: CLI
- **Target**: cli/e2e-compile
- **Test ID**: cli/e2e-compile/no-profile-returns-err-no-profile-for-compile
- **Pre-conditions**: No `.forge/config.yaml` exists; no profile configured
- **Route**: `forge e2e compile`
- **Element**: sitemap-missing
- **Steps**:
  1. Create a temp directory with no `.forge/config.yaml`
  2. Call `Compile(dir)`
  3. Assert the error wraps `ErrNoProfile`
- **Expected**: Error is `ErrNoProfile`; no subprocess is spawned
- **Priority**: P0

## TC-019: No profile returns ErrNoProfile for Discover
- **Source**: Proposal Error Scenarios row 5
- **Type**: CLI
- **Target**: cli/e2e-discover
- **Test ID**: cli/e2e-discover/no-profile-returns-err-no-profile-for-discover
- **Pre-conditions**: No `.forge/config.yaml` exists; no profile configured
- **Route**: `forge e2e discover`
- **Element**: sitemap-missing
- **Steps**:
  1. Create a temp directory with no `.forge/config.yaml`
  2. Call `Discover(dir)`
  3. Assert the error wraps `ErrNoProfile`
- **Expected**: Error is `ErrNoProfile`; no subprocess is spawned
- **Priority**: P0

### Manifest Cleanup

## TC-020: All 6 manifest.yaml files contain zero run and graduate fields
- **Source**: Proposal Success Criteria [4], Task 1 AC [1-2]
- **Type**: CLI
- **Target**: cli/manifest
- **Test ID**: cli/manifest/all-manifests-contain-zero-run-and-graduate-fields
- **Pre-conditions**: All 6 profile directories exist under `pkg/profile/profiles/`
- **Route**: `forge profile get <profile> --manifest`
- **Element**: sitemap-missing
- **Steps**:
  1. Read each of the 6 `manifest.yaml` files
  2. Assert none contain `run.command`, `run.compile`, `run.result-format` fields
  3. Assert none contain `graduate.target-directory`, `graduate.merge-strategy`, `graduate.import-rewrite`, `graduate.compile-check`, `graduate.list-tests` fields
  4. Assert remaining fields (name, display, language, file-extension, test-directory, capabilities, templates) are present and unchanged
- **Expected**: All manifests have metadata-only fields; zero command fields
- **Priority**: P0

---

## Traceability

| TC ID | Source | Type | Target | Priority |
|-------|--------|------|--------|----------|
| TC-001 | Proposal SC [1], Task 2 AC [1] | CLI | cli/e2e-run | P0 |
| TC-002 | Proposal SC [1], Task 2 AC [1], Task 2 HR [3] | CLI | cli/e2e-run | P0 |
| TC-003 | Proposal SC [2], Task 2 AC [2] | CLI | cli/e2e-setup | P0 |
| TC-004 | Proposal SC [3], Task 2 AC [3] | CLI | cli/e2e-compile | P0 |
| TC-005 | Proposal SC [4], Task 2 AC [4] | CLI | cli/e2e-discover | P0 |
| TC-006 | Proposal SC [5], Task 2 AC [5], Task 2 HR [1] | CLI | cli/e2e-verify | P0 |
| TC-007 | Task 2 AC [5] (Verify unchanged) | CLI | cli/e2e-verify | P1 |
| TC-008 | Proposal SC [7], Proposal Error [1], Task 2 AC [7] | CLI | cli/e2e-run | P0 |
| TC-009 | Proposal Error [1], Task 2 AC [7] | CLI | cli/e2e-setup | P0 |
| TC-010 | Proposal Error [1], Task 2 AC [7] | CLI | cli/e2e-compile | P0 |
| TC-011 | Proposal Error [1], Task 2 AC [7] | CLI | cli/e2e-discover | P0 |
| TC-012 | Proposal SC [6], Proposal Error [3], Task 2 AC [6] | CLI | cli/e2e-run | P0 |
| TC-013 | Proposal SC [6], Task 2 AC [6] | CLI | cli/e2e-run | P0 |
| TC-014 | Proposal SC [6], Task 2 AC [6] | CLI | cli/e2e-compile | P0 |
| TC-015 | Proposal SC [6], Task 2 AC [6] | CLI | cli/e2e-discover | P0 |
| TC-016 | Proposal Error [5], Task 2 AC [8] | CLI | cli/e2e-run | P0 |
| TC-017 | Proposal Error [5] | CLI | cli/e2e-setup | P0 |
| TC-018 | Proposal Error [5] | CLI | cli/e2e-compile | P0 |
| TC-019 | Proposal Error [5] | CLI | cli/e2e-discover | P0 |
| TC-020 | Proposal SC [4], Task 1 AC [1-2] | CLI | cli/manifest | P0 |

---

## Route Validation

This is a CLI project using cobra command registration. Route validation maps to CLI subcommand discovery.

| CLI Route | Status | TC IDs | Matched Command |
|-----------|--------|--------|-----------------|
| `forge e2e run` | Matched | TC-001, TC-002, TC-008, TC-012, TC-013, TC-016 | `internal/cmd/e2e_run.go:16` Use: "run" |
| `forge e2e setup` | Matched | TC-003, TC-009, TC-017 | `internal/cmd/e2e_setup.go:16` Use: "setup" |
| `forge e2e compile` | Matched | TC-004, TC-010, TC-014, TC-018 | `internal/cmd/e2e_compile.go:14` Use: "compile" |
| `forge e2e discover` | Matched | TC-005, TC-011, TC-015, TC-019 | `internal/cmd/e2e_discover.go:14` Use: "discover" |
| `forge e2e verify` | Matched | TC-006, TC-007 | `internal/cmd/e2e_verify.go:16` Use: "verify" |
| `forge profile get <profile> --manifest` | Matched | TC-020 | Profile CLI command (existing, unchanged) |
