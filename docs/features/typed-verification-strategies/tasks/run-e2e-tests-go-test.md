---
id: "T-test-3"
title: "Run e2e Tests (go-test)"
priority: "P1"
estimated_time: "30min-1h"
dependencies: ["T-test-2"]
type: "test-pipeline.run"
scope: "all"
profile: "go-test"
---

# Run e2e Tests (go-test)

Profile: **go-test**

# Go Test Run Strategy

Profile-specific execution and result-parsing rules for the `run-e2e-tests` skill.

## Execution

| Item | Value |
|------|-------|
| Command | `go test ./tests/e2e/... -v -tags=e2e -json` |
| Invoked via | `just test-e2e --feature <slug>` |
| Setup | `go build ./...` (verify compilation) |
| Teardown | No special cleanup needed (Go test runner handles process lifecycle) |
| Output format | Go JSON (one JSON object per line, streaming) |

## Build Tag

The `-tags=e2e` flag is **mandatory**. Test files carry `//go:build e2e` and will not compile without it. Omitting this flag results in zero tests discovered.

## Result Format

Go test JSON output (one JSON object per line):

```json
{"Time":"2026-05-12T10:00:00.000Z","Action":"run","Package":"github.com/org/project/tests/e2e","Test":"TestTC_001_Login"}
{"Time":"2026-05-12T10:00:01.234Z","Action":"output","Package":"github.com/org/project/tests/e2e","Test":"TestTC_001_Login","Output":"=== RUN   TestTC_001_Login\n"}
{"Time":"2026-05-12T10:00:02.345Z","Action":"pass","Package":"github.com/org/project/tests/e2e","Test":"TestTC_001_Login","Elapsed":1.111}
```

## Result Parsing Rules

### Event Stream Processing

Read JSON lines sequentially. Track state per test by `Test` name.

### Field Mapping

| Data | Source | Notes |
|------|--------|-------|
| TC ID | `Test` field via regex `TC_\d+` or `TC-\d+` | Extract from test function name |
| Status | `Action` field | `"pass"`, `"fail"`, `"skip"` are terminal states |
| Duration | `Elapsed` field | Seconds (float) |
| Output | `Output` field | Captured stdout/stderr; relevant on `"output"` action |
| Package | `Package` field | Full Go module path |

### Action State Machine

| Action | Meaning |
|--------|---------|
| `run` | Test started |
| `output` | Captured output line |
| `pass` | Test passed (terminal) |
| `fail` | Test failed (terminal) |
| `skip` | Test skipped (terminal) |

## TC ID Extraction

From test function name using pattern `TestTC_NNN_*`:

```
TestTC_001_LoginWithValidCredentials -> TC-001
TestTC_042_ApiErrorHandling          -> TC-042
```

Regex: `TC_(\d+)` or `TC-(\d+)` with separator normalization.

## Test Type Classification

| Type | Indicators |
|------|------------|
| TUI | Golden file comparison, `testdata/*.golden`, snapshot assertion |
| API | `net/http` Client or `httptest`, `http.NewRequest` |
| CLI | `os/exec`, `exec.Command`, `runCLI` helper |

Classification by inspecting test body for characteristic imports/patterns.

## Timeouts

| Setting | Value | Notes |
|---------|-------|-------|
| Default | `-timeout 10m` | 10-minute hard cap for entire test run |
| Per-test | Via `t.Timeout()` or context | Individual test deadline |
| Override | `-timeout` flag | Pass explicitly for long-running suites |

## Error Handling

| Condition | Action |
|-----------|--------|
| Compilation failure | Report build error output, skip all tests |
| Missing binary (CLI tests) | Report which binary is missing, skip affected tests |
| Port conflicts | Report port in use, suggest ephemeral port or kill stale process |
| Test panic | Go runner captures panic as `fail` action with stack trace |
| Module dependency missing | Report `go mod tidy` suggestion, abort |

## Failure Diagnosis

| Failure rate | Response |
|--------------|----------|
| >30% tests fail | Infrastructure problem -- check compilation, server health, dependencies |
| 10-30% | Spot-check 2-3 failing test outputs |
| <10% | Per-test fix tasks |
