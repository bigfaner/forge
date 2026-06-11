---
status: "completed"
started: "2026-05-15 01:06"
completed: "2026-05-15 01:26"
time_spent: "~20m"
---

# Task Record: 2 Implement ensureJust core logic with detection and installation

## Summary
Implement ensureJust core logic: detect just binary in PATH, parse version, compare against minimum (1.40.0), install via system package manager (brew/cargo/scoop/choco) with embedded binary fallback to ~/.forge/bin/. Includes user confirmation prompt with terminal detection (abort on piped stdin).

## Changes

### Files Created
- forge-cli/pkg/just/ensure.go
- forge-cli/pkg/just/ensure_test.go

### Files Modified
无

### Key Decisions
- Used function variables (detectJustFunc, isTerminalFunc, detectPackageManager, embeddedBinaryFunc, userHomeDir) for testability rather than interfaces, matching the existing pattern in just.go
- ParseJustVersion uses regex (^just + non-space) to handle varied output formats including pre-release suffixes
- IsMinimumVersion implements manual semver comparison with pre-release considered less than release
- extractEmbeddedBinary stores the binary path in EnsureResult.Version field for verification in tests
- writeStr helper discards write errors to satisfy errcheck lint since output errors are not actionable in CLI context
- minimumVersion set to 1.40.0 per task spec (proposal.md says 1.50.0 -- task file takes precedence)

## Test Results
- **Tests Executed**: Yes
- **Passed**: 32
- **Failed**: 0
- **Coverage**: 85.1%

## Acceptance Criteria
- [x] DetectJust() (path string, version string, found bool) locates just in PATH and parses its version
- [x] ParseJustVersion(output string) (semver, error) parses just --version output
- [x] IsMinimumVersion(version string, minimum string) bool compares versions
- [x] EnsureJust(io.Reader, io.Writer) EnsureResult orchestrates the full flow with user interaction
- [x] Package manager installation works for brew/cargo/scoop/choco
- [x] Embedded binary fallback extracts to ~/.forge/bin/just
- [x] EnsureResult struct carries Status/Version/Method/Detail
- [x] User is prompted for confirmation before installation (unless CI/pipe: check os.Stdin)
- [x] Outdated version triggers upgrade prompt with clear warning
- [x] All unit tests pass with table-driven test patterns

## Notes
Did NOT modify internal/cmd/init.go per Hard Rules. Package dependency direction maintained: pkg/just imports internal/embedded/just for binary access. Terminal detection uses os.File Stat ModeCharDevice check.
