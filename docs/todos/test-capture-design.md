# Test Result Capture Design

## Context

`task all-completed` currently streams test output directly to terminal. This design
captures output, writes structured result files, and prints only a brief summary.

## Architecture

```
runAllCompleted()
    → runCmdCapture(e2e)   → TestRunResult{stdout, stderr, exitCode, duration}
    → runCmdCapture(proj)  → TestRunResult{stdout, stderr, exitCode, duration}
    → parseResults(e2e)    → TestStats, []TestFailure
    → parseResults(proj)   → TestStats, []TestFailure
    → writeLatestMd()      → docs/features/{slug}/testing/results/latest.md
    → writeFailureFiles()  → docs/features/{slug}/testing/results/failures/failure-{N}.md
    → printSummary()       → brief output to terminal
```

## Types (`test_results.go`)

```go
type TestRunResult struct {
    Stdout   string
    Stderr   string
    ExitCode int
    Duration time.Duration
}

type TestStats struct {
    Total, Pass, Fail, Skip int
    Framework               string // "go", "npm", "make", "pytest", "unknown"
}

type TestFailure struct {
    N        int
    Name     string  // test name
    Package  string  // Go package or file path
    Framework string
    ErrorMsg string
    Output   string  // relevant lines (max 50 lines)
    Duration float64
}
```

## Parsers

### Go — `go test -json ./...`

Newline-delimited JSON events:
- `action:"fail"` + `Test != ""` → failure; accumulate `action:"output"` lines per test
- `action:"pass"/"skip"` → count stats

### npm (node:test TAP)

- `not ok N - {name}` → failure; YAML block between `---` and `...` → error details
- `# tests/pass/fail` lines → stats

### Generic (make, pytest, unknown)

- Exit code 0 → all pass
- Exit code non-0 → 1 failure with full output

## Output Files

### `latest.md`

Path: `docs/features/{slug}/testing/results/latest.md`

```markdown
# Test Results: {slug}

**Date**: 2026-04-24 10:00
**Feature E2E**: 5/5 passed
**Project Tests**: 3/5 passed
**Overall**: FAIL

## Summary

| Suite | Framework | Total | Pass | Fail | Skip |
|-------|-----------|-------|------|------|------|
| Feature E2E | npm | 5 | 5 | 0 | 0 |
| Project | go | 5 | 3 | 2 | 0 |

## Failures

| # | Test | Suite | Error |
|---|------|-------|-------|
| 1 | TestFoo | go | expected nil, got error |
| 2 | TestBar | go | index out of range |

See `failures/failure-{N}.md` for details.
```

### `failures/failure-{N}.md`

Path: `docs/features/{slug}/testing/results/failures/failure-{N}.md`

```markdown
# Failure {N}: {name}

**Framework**: go
**Package**: task-cli/internal/cmd
**Duration**: 0.50s

## Error

```
expected nil, got error: file not found
```

## Output

```
--- FAIL: TestFoo (0.50s)
    all_completed_test.go:45: expected nil, got error: file not found
```
```

## Terminal Summary

```
=== All tasks completed for feature: my-feat ===
e2e:     5 passed, 0 failed
project: 3 passed, 2 failed
Results: docs/features/my-feat/testing/results/latest.md
Failures (2): docs/features/my-feat/testing/results/failures/
```

## Changes to `all_completed.go`

1. Replace `runCmd`/`runShell` with `runCmdCapture()` — captures stdout+stderr, no streaming
2. Go tests: use `go test -json ./...` for structured output
3. After both runs: parse → write files → print summary
