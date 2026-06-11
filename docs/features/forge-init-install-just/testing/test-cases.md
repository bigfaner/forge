---
feature: "forge-init-install-just"
sources:
  - docs/proposals/forge-init-install-just/proposal.md
  - docs/features/forge-init-install-just/tasks/4-tests.md
generated: "2026-05-15"
profile: "go-test"
---

# Test Cases: forge-init-install-just

> **Note**: This feature uses quick mode (no formal PRD). Acceptance criteria are sourced from the proposal success criteria and task 4 acceptance criteria.

> **WARNING**: sitemap.json not found — Element set to `sitemap-missing`. Run `/gen-sitemap` for precise element references. (Non-web-ui profile: sitemap not applicable.)

## Summary

| Type | Count |
|------|-------|
| CLI  | 19    |
| API  | 14    |
| TUI  | 3     |
| **Total** | **36** |

---

## CLI Test Cases

### TC-001: forge init without just triggers installation attempt

- **Source**: Proposal Success Criteria 1 + Task 4 AC "Integration test: forge init without just triggers installation attempt"
- **Type**: CLI
- **Target**: cli/init
- **Test ID**: cli/init/forge-init-without-just-triggers-installation
- **Pre-conditions**: `just` is not in PATH; `forge init` is run on a fresh temp directory
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Set up a temp directory with no pre-existing files
  2. Mock `ensureJustFunc` to return `StatusInstalled` (simulating installation)
  3. Run `forge init --project-root <tmpdir>`
  4. Check stdout for "INSTALLED" status in the summary block
- **Expected**: Output contains "INSTALLED" for "just installation" in the init summary; all other init steps complete normally
- **Priority**: P0

### TC-002: forge init --skip-just skips ensureJust step entirely

- **Source**: Proposal Success Criteria 2 + Task 4 AC "Integration test: forge init with --skip-just skips ensureJust step"
- **Type**: CLI
- **Target**: cli/init
- **Test ID**: cli/init/forge-init-skip-just-skips-step
- **Pre-conditions**: None; `--skip-just` should skip regardless of just availability
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Set up a temp directory
  2. Run `forge init --project-root <tmpdir> --skip-just`
  3. Check stdout for "SKIPPED" for "just installation"
  4. Verify detail mentions "skipped via --skip-just flag"
- **Expected**: Output contains "SKIPPED" for "just installation" with detail "skipped via --skip-just flag"; all other steps still execute normally
- **Priority**: P0

### TC-003: forge init --skip-just still runs all other init steps

- **Source**: Proposal Success Criteria 2 + Task 4 AC "Integration test: forge init with --skip-just skips ensureJust step"
- **Type**: CLI
- **Target**: cli/init
- **Test ID**: cli/init/forge-init-skip-just-runs-other-steps
- **Pre-conditions**: None
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Set up a temp directory
  2. Run `forge init --project-root <tmpdir> --skip-just`
  3. Verify `.forge` directory created
  4. Verify `CLAUDE.md` created
  5. Verify `.gitignore` created
  6. Verify `justfile` created
  7. Verify `config.yaml` created
- **Expected**: All non-just artifacts are created; only the just step is skipped
- **Priority**: P0

### TC-004: forge init with just already installed and meeting minimum version reports SKIPPED

- **Source**: Proposal Success Criteria 3 + Task 4 AC "Integration test: forge init with just already installed reports SKIPPED"
- **Type**: CLI
- **Target**: cli/init
- **Test ID**: cli/init/forge-init-just-installed-reports-skipped
- **Pre-conditions**: Mock `ensureJustFunc` to return `StatusSkipped` with version >= minimum
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Mock `ensureJustFunc` to return `EnsureResult{Status: StatusSkipped, Version: "1.40.0"}`
  2. Run `forge init --project-root <tmpdir>`
  3. Check stdout for "SKIPPED" for "just installation"
  4. Verify detail includes the version string
- **Expected**: Output contains "SKIPPED" for "just installation" with version detail like "just 1.40.0 already available"
- **Priority**: P0

### TC-005: forge init with just below minimum version warns and prompts upgrade

- **Source**: Proposal Success Criteria 4
- **Type**: CLI
- **Target**: cli/init
- **Test ID**: cli/init/forge-init-just-outdated-warns-upgrade
- **Pre-conditions**: Mock `DetectJustFunc` to return version < minimum (e.g., "1.30.0")
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Mock detection to return just version "1.30.0" (below minimum 1.40.0)
  2. Provide "n" on stdin for the upgrade prompt
  3. Run `forge init --project-root <tmpdir>`
  4. Check output contains a warning about outdated version
- **Expected**: Output contains "WARNING" about just version; user declining upgrade results in SKIPPED with detail about version mismatch; init continues
- **Priority**: P1

### TC-006: forge init with just installation failure is non-blocking

- **Source**: Task 4 AC — installation failure is non-blocking (implied by WARNING behavior)
- **Type**: CLI
- **Target**: cli/init
- **Test ID**: cli/init/forge-init-just-install-failure-non-blocking
- **Pre-conditions**: Mock `ensureJustFunc` to return `StatusFailed`
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Mock `ensureJustFunc` to return `EnsureResult{Status: StatusFailed, Detail: "no supported package manager found"}`
  2. Run `forge init --project-root <tmpdir>`
  3. Check stdout for "FAILED" for "just installation"
  4. Check stderr for "WARNING" message
  5. Verify all other init steps still complete
- **Expected**: Output shows "FAILED" for just installation; stderr contains WARNING; other steps (CLAUDE.md, .gitignore, etc.) still succeed
- **Priority**: P0

### TC-007: forge init ensureJust step appears before justfile step in summary

- **Source**: Task 4 AC (implied by init step ordering)
- **Type**: CLI
- **Target**: cli/init
- **Test ID**: cli/init/ensurejust-step-before-justfile-step
- **Pre-conditions**: None
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Run `forge init --project-root <tmpdir> --skip-just`
  2. Find index of "just installation" in output
  3. Find index of "justfile" in output
  4. Verify "just installation" appears before "justfile"
- **Expected**: "just installation" step appears before "justfile" step in the init summary
- **Priority**: P1

### TC-008: forge init on fresh machine installs just via package manager

- **Source**: Proposal Key Scenario "Happy path"
- **Type**: CLI
- **Target**: cli/init
- **Test ID**: cli/init/forge-init-happy-path-pkg-manager
- **Pre-conditions**: just not in PATH; package manager available
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Mock `DetectJustFunc` to return not found
  2. Mock `InstallViaPackageManagerFunc` to return `StatusInstalled` with method "brew"
  3. Mock `isTerminalFunc` to return true; provide "y" on stdin
  4. Run the `EnsureJust` function directly
  5. Verify result is `StatusInstalled` with `Method: "brew"`
- **Expected**: `EnsureResult{Status: StatusInstalled, Method: "brew"}`; no fallback to embedded binary
- **Priority**: P0

### TC-009: forge init falls back to embedded binary when package manager fails

- **Source**: Proposal Key Scenario "Fallback path"
- **Type**: CLI
- **Target**: cli/init
- **Test ID**: cli/init/forge-init-fallback-embedded-binary
- **Pre-conditions**: just not in PATH; package manager install fails
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Mock `DetectJustFunc` to return not found
  2. Mock `InstallViaPackageManagerFunc` to return `StatusFailed`
  3. Mock `ExtractEmbeddedBinaryFunc` to return `StatusInstalled` with method "embedded"
  4. Mock `isTerminalFunc` to return true; provide "y" on stdin
  5. Run `EnsureJust` function
  6. Verify result is `StatusInstalled` with `Method: "embedded"`
- **Expected**: `EnsureResult{Status: StatusInstalled, Method: "embedded"}`; embedded fallback is attempted after package manager failure
- **Priority**: P0

### TC-010: forge init on machine with just already installed skips

- **Source**: Proposal Key Scenario "Already installed"
- **Type**: CLI
- **Target**: cli/init
- **Test ID**: cli/init/forge-init-already-installed-skips
- **Pre-conditions**: Mock `DetectJustFunc` to return version >= minimum
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Mock `DetectJustFunc` to return `("path/to/just", "1.50.0", true)`
  2. Run `EnsureJust` function
  3. Verify result is `StatusSkipped` without any prompts
- **Expected**: `EnsureResult{Status: StatusSkipped, Version: "1.50.0"}`; no stdin interaction occurs
- **Priority**: P0

### TC-011: forge CLI binary size increase is within acceptable limit

- **Source**: Proposal Success Criteria 5 + Non-Functional Requirement "Binary size"
- **Type**: CLI
- **Target**: cli/build
- **Test ID**: cli/build/binary-size-within-limit
- **Pre-conditions**: Forge CLI is built for the current platform
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Build the forge CLI binary with build tags for the current platform
  2. Check the resulting binary file size
  3. Verify size increase is <= 5 MB above baseline (without embedded binary)
- **Expected**: Binary size increase due to embedded just is <= 5 MB per platform
- **Priority**: P2

### TC-012: User declines installation when just is not found

- **Source**: Proposal Key Scenario — user choice preserved
- **Type**: CLI
- **Target**: cli/init
- **Test ID**: cli/init/user-declines-installation
- **Pre-conditions**: just not in PATH
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Mock `DetectJustFunc` to return not found
  2. Mock `isTerminalFunc` to return true; provide "n" on stdin
  3. Run `EnsureJust` function
  4. Verify result is `StatusSkipped` with "user declined installation"
- **Expected**: `EnsureResult{Status: StatusSkipped, Detail: "user declined installation"}`
- **Priority**: P1

### TC-013: Non-interactive stdin when just is not found fails gracefully

- **Source**: Proposal — installation requires user confirmation
- **Type**: CLI
- **Target**: cli/init
- **Test ID**: cli/init/non-interactive-stdin-not-found-fails
- **Pre-conditions**: just not in PATH; stdin is piped (not a terminal)
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Mock `DetectJustFunc` to return not found
  2. Mock `isTerminalFunc` to return false
  3. Run `EnsureJust` function with a bytes.Buffer as stdin
  4. Verify result is `StatusFailed` with "non-interactive stdin" detail
- **Expected**: `EnsureResult{Status: StatusFailed, Detail: "non-interactive stdin (piped); cannot prompt for installation"}`
- **Priority**: P1

### TC-014: forge init --project-root with custom path

- **Source**: Proposal — `--project-root` flag support
- **Type**: CLI
- **Target**: cli/init
- **Test ID**: cli/init/forge-init-custom-project-root
- **Pre-conditions**: Custom temp directory path
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Create a temp directory at a custom path
  2. Run `forge init --project-root <custom-path> --skip-just`
  3. Verify all artifacts are created in the custom directory
- **Expected**: All init artifacts (`.forge`, `CLAUDE.md`, `.gitignore`, `justfile`, `config.yaml`) created in the custom directory
- **Priority**: P1

### TC-015: ensureResultToAction maps installed result with method detail

- **Source**: Task 4 AC — status reporting
- **Type**: CLI
- **Target**: cli/init
- **Test ID**: cli/init/ensure-result-to-action-installed
- **Pre-conditions**: None
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Create `EnsureResult{Status: StatusInstalled, Version: "1.40.0", Method: "brew"}`
  2. Call `ensureResultToAction(result)`
  3. Verify action.status is "INSTALLED"
  4. Verify action.detail contains "brew" and "1.40.0"
- **Expected**: `initAction{status: "INSTALLED", target: "just installation", detail: "installed via brew (1.40.0)"}`
- **Priority**: P1

### TC-016: ensureResultToAction maps skipped result with version detail

- **Source**: Task 4 AC — status reporting
- **Type**: CLI
- **Target**: cli/init
- **Test ID**: cli/init/ensure-result-to-action-skipped
- **Pre-conditions**: None
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Create `EnsureResult{Status: StatusSkipped, Version: "1.40.0"}`
  2. Call `ensureResultToAction(result)`
  3. Verify action.status is "SKIPPED"
  4. Verify action.detail contains "just 1.40.0 already available"
- **Expected**: `initAction{status: "SKIPPED", target: "just installation", detail: "just 1.40.0 already available"}`
- **Priority**: P1

### TC-017: ensureResultToAction maps failed result

- **Source**: Task 4 AC — status reporting
- **Type**: CLI
- **Target**: cli/init
- **Test ID**: cli/init/ensure-result-to-action-failed
- **Pre-conditions**: None
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Create `EnsureResult{Status: StatusFailed, Detail: "no package manager"}`
  2. Call `ensureResultToAction(result)`
  3. Verify action.status is "FAILED"
- **Expected**: `initAction{status: "FAILED", target: "just installation", detail: "no package manager"}`
- **Priority**: P1

### TC-018: User declines upgrade for outdated just version

- **Source**: Proposal Key Scenario "Outdated version" — user declines
- **Type**: CLI
- **Target**: cli/init
- **Test ID**: cli/init/user-declines-upgrade-outdated
- **Pre-conditions**: Mock detection to return version below minimum
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Mock `DetectJustFunc` to return `("path/to/just", "1.30.0", true)`
  2. Mock `isTerminalFunc` to return true; provide "n" on stdin
  3. Run `EnsureJust` function
  4. Verify result is `StatusSkipped` with detail about declined upgrade
- **Expected**: `EnsureResult{Status: StatusSkipped, Version: "1.30.0", Detail: "user declined upgrade; just 1.30.0 < 1.40.0"}`
- **Priority**: P1

### TC-019: User accepts upgrade for outdated just version

- **Source**: Proposal Key Scenario "Outdated version" — user accepts
- **Type**: CLI
- **Target**: cli/init
- **Test ID**: cli/init/user-accepts-upgrade-outdated
- **Pre-conditions**: Mock detection to return version below minimum
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Mock `DetectJustFunc` to return `("path/to/just", "1.30.0", true)`
  2. Mock `isTerminalFunc` to return true; provide "y" on first stdin read, "y" on second
  3. Mock `InstallViaPackageManagerFunc` to return `StatusInstalled`
  4. Run `EnsureJust` function
  5. Verify result is `StatusInstalled`
- **Expected**: `EnsureResult{Status: StatusInstalled}` — upgrade proceeds after user accepts
- **Priority**: P1

---

## API Test Cases

### TC-020: DetectJust finds just in PATH and parses version

- **Source**: Task 4 AC "Unit tests cover: DetectJust (found/not found)"
- **Type**: API
- **Target**: api/just-detect
- **Test ID**: api/just-detect/detectjust-finds-just-in-path
- **Pre-conditions**: `just` binary is available in PATH (table-driven: with and without just)
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Call `DetectJust()` when just is in PATH
  2. Verify path is non-empty, version is parsed, found is true
  3. Call `DetectJust()` when just is not in PATH (mock `exec.LookPath`)
  4. Verify path is empty, version is empty, found is false
- **Expected**: When found: `(path, version, true)` where version matches expected. When not found: `("", "", false)`
- **Priority**: P0

### TC-021: DetectJust handles binary found but version command fails

- **Source**: Task 4 AC "Unit tests cover: DetectJust (found/not found)"
- **Type**: API
- **Target**: api/just-detect
- **Test ID**: api/just-detect/detectjust-version-command-fails
- **Pre-conditions**: just binary is in PATH but `just --version` returns error
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Mock `exec.LookPath` to return a path
  2. Mock `exec.Command.Output` to return error for `--version`
  3. Call `DetectJust()`
  4. Verify path is returned but version is empty, found is true
- **Expected**: `(path, "", true)` — binary exists but version is unknown
- **Priority**: P1

### TC-022: ParseJustVersion parses valid version output

- **Source**: Task 4 AC "Unit tests cover: ParseJustVersion (valid/invalid formats)"
- **Type**: API
- **Target**: api/just-parse
- **Test ID**: api/just-parse/parsejustversion-valid-formats
- **Pre-conditions**: None
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Table-driven: call `ParseJustVersion` with inputs:
     - `"just 1.40.0\n"` -> `"1.40.0"`
     - `"just 1.37.0\n"` -> `"1.37.0"`
     - `"just 1.50.0-beta.1\n"` -> `"1.50.0-beta.1"`
  2. Verify each returns the expected version string
- **Expected**: Each input correctly extracts the version substring
- **Priority**: P0

### TC-023: ParseJustVersion rejects invalid format

- **Source**: Task 4 AC "Unit tests cover: ParseJustVersion (valid/invalid formats)"
- **Type**: API
- **Target**: api/just-parse
- **Test ID**: api/just-parse/parsejustversion-invalid-formats
- **Pre-conditions**: None
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Table-driven: call `ParseJustVersion` with invalid inputs:
     - `""` (empty string)
     - `"unknown output"`
     - `"1.40.0"` (missing "just" prefix)
  2. Verify each returns an error
- **Expected**: Each invalid input returns a non-nil error
- **Priority**: P1

### TC-024: IsMinimumVersion compares versions correctly

- **Source**: Task 4 AC "Unit tests cover: IsMinimumVersion (equal/above/below/edge cases)"
- **Type**: API
- **Target**: api/just-version
- **Test ID**: api/just-version/isminimumversion-comparison
- **Pre-conditions**: None
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Table-driven: call `IsMinimumVersion` with:
     - `("1.40.0", "1.40.0")` -> true (equal)
     - `("1.50.0", "1.40.0")` -> true (above)
     - `("1.30.0", "1.40.0")` -> false (below)
     - `("2.0.0", "1.40.0")` -> true (major above)
     - `("0.9.0", "1.40.0")` -> false (major below)
     - `("1.40.1", "1.40.0")` -> true (patch above)
     - `("1.40.0-pre", "1.40.0")` -> false (pre-release less than release)
- **Expected**: Each comparison returns the correct boolean result
- **Priority**: P0

### TC-025: IsMinimumVersion handles edge cases

- **Source**: Task 4 AC "Unit tests cover: IsMinimumVersion (edge cases)"
- **Type**: API
- **Target**: api/just-version
- **Test ID**: api/just-version/isminimumversion-edge-cases
- **Pre-conditions**: None
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Table-driven: call `IsMinimumVersion` with edge cases:
     - `("0.0.0", "0.0.0")` -> true (zero versions)
     - `("999.999.999", "1.40.0")` -> true (very high version)
     - `("invalid", "1.40.0")` -> false (non-parseable version returns zero semver)
- **Expected**: Edge case comparisons return correct boolean results
- **Priority**: P2

### TC-026: Package manager dispatch per OS (macOS brew)

- **Source**: Task 4 AC "Unit tests cover: package manager dispatch logic per OS"
- **Type**: API
- **Target**: api/just-install
- **Test ID**: api/just-install/pkg-manager-dispatch-macos-brew
- **Pre-conditions**: Mock `runtime.GOOS` and `exec.LookPath` for macOS with brew available
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Mock OS to "darwin", LookPath to find "brew"
  2. Call `detectPackageManager()`
  3. Verify result is "brew"
- **Expected**: Returns "brew" when brew is available on macOS
- **Priority**: P1

### TC-027: Package manager dispatch per OS (macOS cargo fallback)

- **Source**: Task 4 AC "Unit tests cover: package manager dispatch logic per OS"
- **Type**: API
- **Target**: api/just-install
- **Test ID**: api/just-install/pkg-manager-dispatch-macos-cargo-fallback
- **Pre-conditions**: Mock `runtime.GOOS` and `exec.LookPath` for macOS without brew but with cargo
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Mock OS to "darwin", LookPath to not find "brew" but find "cargo"
  2. Call `detectPackageManager()`
  3. Verify result is "cargo"
- **Expected**: Returns "cargo" when brew is unavailable but cargo is on macOS
- **Priority**: P1

### TC-028: Package manager dispatch per OS (Windows scoop)

- **Source**: Task 4 AC "Unit tests cover: package manager dispatch logic per OS"
- **Type**: API
- **Target**: api/just-install
- **Test ID**: api/just-install/pkg-manager-dispatch-windows-scoop
- **Pre-conditions**: Mock `runtime.GOOS` and `exec.LookPath` for Windows with scoop
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Mock OS to "windows", LookPath to find "scoop"
  2. Call `detectPackageManager()`
  3. Verify result is "scoop"
- **Expected**: Returns "scoop" when scoop is available on Windows
- **Priority**: P1

### TC-029: Package manager dispatch per OS (Windows choco fallback)

- **Source**: Task 4 AC "Unit tests cover: package manager dispatch logic per OS"
- **Type**: API
- **Target**: api/just-install
- **Test ID**: api/just-install/pkg-manager-dispatch-windows-choco
- **Pre-conditions**: Mock `runtime.GOOS` and `exec.LookPath` for Windows without scoop but with choco
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Mock OS to "windows", LookPath to not find "scoop" but find "choco"
  2. Call `detectPackageManager()`
  3. Verify result is "choco"
- **Expected**: Returns "choco" when scoop is unavailable but choco is on Windows
- **Priority**: P1

### TC-030: Package manager dispatch returns empty when no package manager found

- **Source**: Task 4 AC "Unit tests cover: package manager dispatch logic per OS"
- **Type**: API
- **Target**: api/just-install
- **Test ID**: api/just-install/pkg-manager-dispatch-none-found
- **Pre-conditions**: Mock no package managers available
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Mock `exec.LookPath` to not find any package manager
  2. Call `detectPackageManager()`
  3. Verify result is empty string
- **Expected**: Returns `""` when no supported package manager is found
- **Priority**: P1

### TC-031: Embedded binary extraction to ~/.forge/bin/ succeeds

- **Source**: Task 4 AC "Unit tests cover: embedded binary extraction to ~/.forge/bin/ (with temp dirs)"
- **Type**: API
- **Target**: api/just-extract
- **Test ID**: api/just-extract/embedded-binary-extraction-success
- **Pre-conditions**: Mock `userHomeDir` to return a temp directory; mock `EmbeddedBinaryFunc` to return non-empty byte slice
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Mock `userHomeDir` to return `t.TempDir()`
  2. Mock `EmbeddedBinaryFunc` to return `[]byte("fake binary data")`
  3. Call `ExtractEmbeddedBinaryFunc()`
  4. Verify result is `StatusInstalled` with method "embedded"
  5. Verify file exists at `<tmpdir>/.forge/bin/just` (or `just.exe` on Windows)
  6. Verify file contents match the mock binary data
  7. Verify file permissions are 0o755
- **Expected**: Binary extracted to `~/.forge/bin/just` with correct contents and permissions; result `StatusInstalled, Method: "embedded"`
- **Priority**: P0

### TC-032: Embedded binary extraction fails with empty binary data

- **Source**: Task 4 AC "Unit tests cover: embedded binary extraction to ~/.forge/bin/ (with temp dirs)"
- **Type**: API
- **Target**: api/just-extract
- **Test ID**: api/just-extract/embedded-binary-empty-data-fails
- **Pre-conditions**: Mock `EmbeddedBinaryFunc` to return empty byte slice
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Mock `userHomeDir` to return `t.TempDir()`
  2. Mock `EmbeddedBinaryFunc` to return `[]byte{}`
  3. Call `ExtractEmbeddedBinaryFunc()`
  4. Verify result is `StatusFailed`
  5. Verify detail mentions "no embedded just binary"
- **Expected**: `EnsureResult{Status: StatusFailed, Detail: "no embedded just binary available for this platform"}`
- **Priority**: P1

### TC-033: Embedded binary extraction handles permission denied on write

- **Source**: Task 4 AC "Edge cases: permission denied on extraction"
- **Type**: API
- **Target**: api/just-extract
- **Test ID**: api/just-extract/embedded-binary-permission-denied
- **Pre-conditions**: Mock `userHomeDir` to return a read-only directory
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Mock `userHomeDir` to return a read-only temp directory
  2. Mock `EmbeddedBinaryFunc` to return non-empty byte slice
  3. Call `ExtractEmbeddedBinaryFunc()`
  4. Verify result is `StatusFailed`
  5. Verify detail mentions "cannot write" or "cannot create"
- **Expected**: `EnsureResult{Status: StatusFailed}` with detail about write failure
- **Priority**: P2

---

## TUI Test Cases

### TC-034: EnsureJust prompts for confirmation when just is not found

- **Source**: Proposal Key Scenario "Happy path" — user confirmation step
- **Type**: TUI
- **Target**: tui/ensurejust
- **Test ID**: tui/ensurejust/prompt-when-just-not-found
- **Pre-conditions**: Mock detection to return not found; mock terminal to return true
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Mock `DetectJustFunc` to return not found
  2. Mock `isTerminalFunc` to return true
  3. Provide "y" via stdin buffer
  4. Capture output written to the writer
  5. Verify output contains "just is not installed. Install just? [Y/n]:"
- **Expected**: Output contains the installation confirmation prompt
- **Priority**: P0

### TC-035: EnsureJust prompts for upgrade when version is outdated

- **Source**: Proposal Key Scenario "Outdated version" — upgrade prompt
- **Type**: TUI
- **Target**: tui/ensurejust
- **Test ID**: tui/ensurejust/prompt-upgrade-outdated-version
- **Pre-conditions**: Mock detection to return version below minimum; mock terminal to return true
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Mock `DetectJustFunc` to return `("path", "1.30.0", true)`
  2. Mock `isTerminalFunc` to return true
  3. Provide "y" via stdin buffer (for both upgrade and install prompts)
  4. Capture output
  5. Verify output contains "WARNING: just 1.30.0 found" and "Upgrade? [y/N]:"
- **Expected**: Output contains the upgrade warning and prompt; user can accept or decline
- **Priority**: P1

### TC-036: Non-interactive stdin with outdated just fails with descriptive message

- **Source**: Proposal — non-interactive stdin handling for outdated version
- **Type**: TUI
- **Target**: tui/ensurejust
- **Test ID**: tui/ensurejust/non-interactive-stdin-outdated-fails
- **Pre-conditions**: Mock detection to return version below minimum; mock terminal to return false
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Mock `DetectJustFunc` to return `("path", "1.30.0", true)`
  2. Mock `isTerminalFunc` to return false
  3. Run `EnsureJust` with a bytes.Buffer as stdin
  4. Verify result is `StatusFailed` with "non-interactive stdin" and "outdated" detail
- **Expected**: `EnsureResult{Status: StatusFailed, Detail: "outdated just and non-interactive stdin (piped); cannot prompt for upgrade"}`
- **Priority**: P1

---

## Traceability

| TC ID | Source | Type | Target | Priority |
|-------|--------|------|--------|----------|
| TC-001 | Proposal SC 1 + Task 4 AC | CLI | cli/init | P0 |
| TC-002 | Proposal SC 2 + Task 4 AC | CLI | cli/init | P0 |
| TC-003 | Proposal SC 2 + Task 4 AC | CLI | cli/init | P0 |
| TC-004 | Proposal SC 3 + Task 4 AC | CLI | cli/init | P0 |
| TC-005 | Proposal SC 4 | CLI | cli/init | P1 |
| TC-006 | Task 4 AC | CLI | cli/init | P0 |
| TC-007 | Task 4 AC (implied) | CLI | cli/init | P1 |
| TC-008 | Proposal "Happy path" | CLI | cli/init | P0 |
| TC-009 | Proposal "Fallback path" | CLI | cli/init | P0 |
| TC-010 | Proposal "Already installed" | CLI | cli/init | P0 |
| TC-011 | Proposal SC 5 | CLI | cli/build | P2 |
| TC-012 | Proposal "User choice" | CLI | cli/init | P1 |
| TC-013 | Proposal (non-interactive) | CLI | cli/init | P1 |
| TC-014 | Proposal (flags) | CLI | cli/init | P1 |
| TC-015 | Task 4 AC (status mapping) | CLI | cli/init | P1 |
| TC-016 | Task 4 AC (status mapping) | CLI | cli/init | P1 |
| TC-017 | Task 4 AC (status mapping) | CLI | cli/init | P1 |
| TC-018 | Proposal "Outdated version" | CLI | cli/init | P1 |
| TC-019 | Proposal "Outdated version" | CLI | cli/init | P1 |
| TC-020 | Task 4 AC (DetectJust) | API | api/just-detect | P0 |
| TC-021 | Task 4 AC (DetectJust) | API | api/just-detect | P1 |
| TC-022 | Task 4 AC (ParseJustVersion) | API | api/just-parse | P0 |
| TC-023 | Task 4 AC (ParseJustVersion) | API | api/just-parse | P1 |
| TC-024 | Task 4 AC (IsMinimumVersion) | API | api/just-version | P0 |
| TC-025 | Task 4 AC (IsMinimumVersion) | API | api/just-version | P2 |
| TC-026 | Task 4 AC (pkg manager) | API | api/just-install | P1 |
| TC-027 | Task 4 AC (pkg manager) | API | api/just-install | P1 |
| TC-028 | Task 4 AC (pkg manager) | API | api/just-install | P1 |
| TC-029 | Task 4 AC (pkg manager) | API | api/just-install | P1 |
| TC-030 | Task 4 AC (pkg manager) | API | api/just-install | P1 |
| TC-031 | Task 4 AC (embedded binary) | API | api/just-extract | P0 |
| TC-032 | Task 4 AC (embedded binary) | API | api/just-extract | P1 |
| TC-033 | Task 4 AC (edge cases) | API | api/just-extract | P2 |
| TC-034 | Proposal "Happy path" | TUI | tui/ensurejust | P0 |
| TC-035 | Proposal "Outdated version" | TUI | tui/ensurejust | P1 |
| TC-036 | Proposal (non-interactive) | TUI | tui/ensurejust | P1 |
