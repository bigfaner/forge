# Pytest Run Strategy

Profile-specific execution and result-parsing rules for the `run-e2e-tests` skill.

## Execution

| Item | Value |
|------|-------|
| Command | `python -m pytest tests/e2e/ -v --tb=short` |
| Invoked via | `just test-e2e --feature <slug>` |
| Setup | `python -m py_compile tests/e2e/` |
| Teardown | No special cleanup needed |

## Timeouts

| Setting | Value | Notes |
|---------|-------|-------|
| Per-test timeout | `--timeout=60` (if pytest-timeout installed) | Seconds |
| Global timeout | None by default | Use `--timeout` for per-test |

## Result Format

Pytest verbose output (text-based). With `--json-report` (if pytest-json-report installed):

```json
{
  "tests": [{
    "nodeid": "tests/e2e/test_feature.py::test_tc_001_description",
    "outcome": "passed",
    "duration": 0.123
  }]
}
```

Parse stdout for fallback:
```
tests/e2e/test_feature.py::test_tc_001_description PASSED
tests/e2e/test_feature.py::test_tc_002_another FAILED
=== 1 passed, 1 failed in 0.5s ===
```

## Result Parsing Rules

| Data | Source | Notes |
|------|--------|-------|
| TC ID | Test function name via regex `test_tc_\d+` | From nodeid or function name |
| Status | `PASSED` / `FAILED` / `SKIPPED` / `XFAILED` | From pytest output |
| Duration | `--duration` in seconds | Available with verbose output |
| Error message | Traceback in FAILED output | `--tb=short` format |

## Test Type Classification

| Type | Indicators |
|------|------------|
| API | Uses `requests` or `httpx` imports |
| CLI | Uses `subprocess` imports |

## Error Handling

| Condition | Action |
|-----------|--------|
| Python not found | Report error, abort |
| Module import error | Report missing dependency, abort |
| Test compilation failure | Report syntax error, skip that file |
