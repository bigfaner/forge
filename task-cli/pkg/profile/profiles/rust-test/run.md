# Rust Test Run Strategy

Profile-specific execution and result-parsing rules for the `run-e2e-tests` skill.

## Execution

| Item | Value |
|------|-------|
| Command | `cargo test --test e2e -- --nocapture` |
| Invoked via | `just test-e2e --feature <slug>` |
| Setup | `cargo build` |
| Teardown | No special cleanup needed |

## Timeouts

| Setting | Value | Notes |
|---------|-------|-------|
| Default timeout | None | Rust tests run until completion |
| Custom timeout | `RUST_TEST_TIMEOUT` env var | Per-process timeout in seconds |

## Result Format

Cargo test output (text-based). Parse stdout for test results:

```
test test_tc_001_description ... ok
test test_tc_002_another_test ... FAILED
test test_tc_003_edge_case ... ignored

test result: FAILED. 2 passed; 1 failed; 1 ignored; 0 measured; 0 filtered out
```

With `--format json` (nightly): structured JSON per test event.

## Result Parsing Rules

| Data | Source | Notes |
|------|--------|-------|
| TC ID | Test function name via regex `test_tc_\d+` | Extract from function name |
| Status | `ok` / `FAILED` / `ignored` | After `...` separator |
| Duration | Not available in stable Rust | Use `-- -Z unstable-options --format json` on nightly |

## Test Type Classification

| Type | Indicators |
|------|------------|
| CLI | Uses `std::process::Command` |
| API | Uses `reqwest` or `ureq` |
| TUI | Uses `Command` + stdout comparison |

## Error Handling

| Condition | Action |
|-----------|--------|
| Compilation failure | Report compile errors, abort |
| Binary not found | Run `cargo build`, retry once |
| Test panic | Capture panic message from output |

## Failure Diagnosis

Same threshold-based approach: >30% failures → likely app health issue.
