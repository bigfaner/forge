---
name: test-isolation
description: Testing isolation conventions for unit and surface tests -- every test must own its environment
---

<!-- OWNER: run-tests | CONSUMERS: gen-test-scripts (INLINE) -->

# Testing Isolation Conventions

Every test — unit or surface — must own its environment. No test may depend on the real project's filesystem state, git state, or `.forge/` state.

## General Rules

### TEST-isolation-000: Tests Must Be Self-Contained

**Requirement**: Every test MUST create its own world via `t.TempDir()` (unit) or isolated fixture setup (e2e). Tests MUST NOT read, write, or depend on files outside their own sandbox.

**Applies to**: All test files, regardless of layer (unit, integration, e2e).

## Unit Tests

### TEST-isolation-001: Project Root Detection Must Be Isolated

**Requirement**: Unit tests that call `FindProjectRoot()`, `runFeature()`, `executeClaim()`, or any function relying on project root detection MUST set `CLAUDE_PROJECT_DIR` via `t.Setenv()` to isolate from the real project's `.forge/state.json` and workspace markers.

**Why**: `FindRootInfoFrom` walks up the entire directory tree collecting markers (workspace > project > VCS priority). A temp dir with only `go.mod` (project marker) will be overridden by `.forge` (workspace marker) found in the real project root during the upward walk. This causes tests to read real project state instead of the test fixture.

**Pattern**:

```go
func TestMyFeature(t *testing.T) {
    dir := t.TempDir()
    // ... set up test fixtures in dir ...

    t.Setenv("CLAUDE_PROJECT_DIR", dir) // isolate project root

    origWd, _ := os.Getwd()
    t.Cleanup(func() { _ = os.Chdir(origWd) })
    _ = os.Chdir(dir)

    // ... test code ...
}
```

**Applies to**: All tests in `forge-cli/internal/cmd/` and any test that exercises CLI commands reading project state.

**Source**: `gotcha-stale-state-json-feature-mismatch.md` — unit tests read real `.forge/state.json` when project root detection is not isolated.

## E2E CLI Tests

### TEST-isolation-002: CLI Functional Tests Must Use `CLAUDE_PROJECT_DIR`

**Requirement**: CLI functional tests that invoke `forge` CLI commands MUST pass `CLAUDE_PROJECT_DIR=<fixture-dir>` via `cmd.Env` to prevent the CLI from detecting the real project root.

**Why**: The same `FindRootInfoFrom` walk-up logic applies to the CLI binary. Without isolation, `forge feature`, `forge task claim`, etc. may resolve to the real project root and read/write real state.

**Pattern**:

```go
func TestMyE2E(t *testing.T) {
    dir := t.TempDir()
    // ... set up forge project structure in dir ...

    cmd := exec.Command("forge", "task", "claim")
    cmd.Env = append(os.Environ(), "CLAUDE_PROJECT_DIR="+dir)
    out, err := cmd.CombinedOutput()
    // ... assertions ...
}
```

**Applies to**: All tests under `tests/` that spawn `forge` CLI commands.

### TEST-isolation-003: CLI Functional Test Fixtures Must Be Complete

**Requirement**: CLI functional test helpers that create project directories (e.g., `setupProjectDir`, `ensureFeatureDir`) MUST create all files required by the code under test — including `tasks/index.json`, `.forge/config.yaml`, and any other files the production code checks for.

**Why**: Production code validates feature directories by checking for `tasks/index.json`. A helper that creates directory structure without this file will cause false test failures that don't reproduce in production.

**Pattern**:

```go
// BAD: missing index.json causes features-dir scanner to skip
func ensureFeatureDir(t *testing.T, root, slug string) {
    os.MkdirAll(filepath.Join(root, "docs", "features", slug, "tasks"), 0755)
}

// GOOD: includes index.json so scanner recognizes the feature
func ensureFeatureDir(t *testing.T, root, slug string) {
    dir := filepath.Join(root, "docs", "features", slug)
    os.MkdirAll(filepath.Join(dir, "tasks"), 0755)
    os.WriteFile(filepath.Join(dir, "tasks", "index.json"), []byte("{}"), 0644)
}
```

**Source**: `disc-1` fix — `ensureFeatureDir` in `feature_set_command_cli_test.go` was missing `tasks/index.json`.

### TEST-isolation-004: CLI Functional Tests Should Compile a Dedicated Forge Binary

**Requirement**: CLI functional test files and suites that invoke forge CLI commands SHOULD compile a forge binary from the current source tree (via `go build`) and use that binary for all `exec.Command` invocations, rather than relying on the system-installed `forge` binary via `$PATH` resolution. New CLI functional test files MUST follow this pattern. Legacy tests using `exec.Command("forge", ...)` via `helpers_test.go` / `testkit/helpers.go` are acceptable but should be migrated when modified.

**Why**: The system-installed `forge` binary may come from a different branch (e.g., `main`). E2e tests that test branch-specific features or behavior changes will produce false failures (new commands missing, removed commands still present, modified output format unchanged). Building from source guarantees the binary matches the code under test. Additionally, using the production binary risks corrupting its state or confusing test results with real project operations.

**Current state**: Some e2e test files compile a dedicated binary (e.g., `spec_drift_detection_cli_test.go`, `features/task-stage-gates/`, `features/task-type-refinement/`, `features/forge-init-install-just/`). The shared helpers (`helpers_test.go`, `testkit/helpers.go`) still use `exec.Command("forge", ...)` and have not been migrated.

**Pattern**:

```go
var forgeBinary string

func TestMain(m *testing.M) {
    tmpDir, _ := os.MkdirTemp("", "forge-e2e-*")
    defer os.RemoveAll(tmpDir)

    forgeBinary = filepath.Join(tmpDir, "forge-test")
    cmd := exec.Command("go", "build", "-o", forgeBinary, "./cmd/forge")
    if out, err := cmd.CombinedOutput(); err != nil {
        fmt.Fprintf(os.Stderr, "build failed: %s\n%s", err, out)
        os.Exit(1)
    }

    os.Exit(m.Run())
}

// In test functions:
cmd := exec.Command(forgeBinary, "task", "claim") // NOT exec.Command("forge", ...)
```

**Anti-pattern**:

```go
// BAD: resolves to system-installed forge, may be from a different branch
cmd := exec.Command("forge", "task", "claim")
```

**Scope**: All e2e test packages that invoke the forge CLI binary:
- `tests/` (including `tests/<journey>/` sub-packages) — uses `forge_binary.go` init() + `TestMain` alias pattern
- `tests/justfile-canonical-e2e/` — uses `TestMain` direct build pattern

**Source**: `/learn` entry 2026-05-20 — 8/36 e2e tests failed because they ran against system-installed forge from `main` branch instead of the feature branch binary.
