---
title: "Testing Isolation Conventions"
---

# Testing Isolation Conventions

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
