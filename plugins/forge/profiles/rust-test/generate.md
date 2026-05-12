# Rust Test Generate Strategy

Profile-specific test generation rules for the `gen-test-scripts` skill.

## Test Runner & Imports

| Item | Value |
|------|-------|
| Test runner | Built-in `#[test]` attribute |
| Assertions | `assert!`, `assert_eq!`, `assert_ne!` |
| Test attribute | `#[test]` on each test function |
| Async support | `#[tokio::test]` with `tokio` crate |
| Serial tests | `serial_test` crate `#[serial]` attribute |

## Spec Template Mapping

| Test type | Template file | Output filename |
|-----------|--------------|-----------------|
| Mixed | `templates/test_file.rs` | `e2e_tests.rs` |

Rust tests are typically in a single file per integration test binary (`tests/` directory).

## CLI Testing

Use `std::process::Command`:

```rust
let output = Command::new("./binary")
    .args(&["--flag", "value"])
    .output()
    .expect("failed to execute");

assert!(output.status.success());
assert!(String::from_utf8_lossy(&output.stdout).contains("expected text"));
```

## API Testing

Use `reqwest` crate (or `ureq` for sync):

```rust
let client = reqwest::Client::new();
let resp = client.get("http://localhost:8080/api/health")
    .send().await
    .expect("request failed");
assert_eq!(resp.status(), StatusCode::OK);
```

## TUI Testing

Test terminal output via subprocess execution and stdout comparison:

```rust
let output = Command::new("./tui-app").output().expect("failed");
let stdout = String::from_utf8_lossy(&output.stdout);
assert!(stdout.contains("expected rendering"));
```

## Auth

- CLI: environment variables (`env::var("API_KEY")`) or config file
- API: `reqwest` header injection via `.header("Authorization", format!("Bearer {}", token))`

## Import Conventions

Rust integration tests in `tests/` directory import from the crate under test via `use crate_name::*` or direct dependency usage.

## Anti-Patterns (Forbidden)

- No `std::thread::sleep` — use `tokio::time::sleep` with timeout, or poll loops
- No hardcoded URLs — use environment variables or config
- No unwrapping without error context — use `.expect("descriptive message")`

## Compilation Check

```bash
just e2e-compile
```

## Traceability

Each test function includes a doc comment with TC ID and PRD source:

```rust
/// Traceability: TC-NNN → {PRD Source}
#[test]
fn test_tc_nnn_description() { ... }
```
