---
title: "Rust cargo test Testing Convention"
---

# Rust cargo test Testing Convention

Convention for generating Rust test code using the built-in test framework with `cargo test`.

## framework

- **name**: cargo test (built-in Rust test framework)
- **version**: Rust 1.70+
- **language**: Rust
- **runner_command**: `cargo test`

## discovery

- **test_dir**: `tests/` (integration tests), inline `#[test]` in `src/` (unit tests)
- **file_pattern**: `*.rs` (any `.rs` file in `tests/` is an integration test crate)
- **exclude_pattern**: `target/`, `node_modules/`

## structure

- **suite_pattern**: Module — each file in `tests/` is a separate test crate; `#[cfg(test)] mod` for unit tests
- **case_pattern**: `#[test] fn test_name()` — annotated functions
- **hook_pattern**: No built-in fixtures — use per-test setup functions or helper crates

### Unit Tests (inline)

Place within `src/` files using `#[cfg(test)]` module:

```rust
#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_something() {
        assert_eq!(2 + 2, 4);
    }
}
```

### Integration Tests

Place in `tests/` directory — each file is a separate crate:

```rust
// tests/integration_test.rs
use std::process::Command;

#[test]
fn test_cli_command() {
    let output = Command::new("forge")
        .args(["subcommand", "--flag", "value"])
        .output()
        .expect("failed to run forge");
    assert!(output.status.success());
    assert!(String::from_utf8_lossy(&output.stdout).contains("expected"));
}
```

### Test Organization

```
project/
├── src/
│   └── lib.rs          — #[cfg(test)] mod tests { ... }
├── tests/
│   ├── common/
│   │   └── mod.rs      — shared helpers (not treated as test crate)
│   ├── cli_test.rs     — integration test for CLI
│   └── api_test.rs     — integration test for API
```

### Parameterized Tests

Use the `rstest` crate for parameterized testing:

```rust
use rstest::rstest;

#[rstest]
#[case("hello", "HELLO")]
#[case("", "")]
#[case("123", "123")]
fn test_uppercase(#[case] input: &str, #[case] expected: &str) {
    assert_eq!(input.to_uppercase(), expected);
}
```

Or manual iteration:

```rust
#[test]
fn test_various_inputs() {
    let cases = vec![
        ("hello", "HELLO"),
        ("", ""),
        ("123", "123"),
    ];
    for (input, expected) in cases {
        assert_eq!(input.to_uppercase(), expected);
    }
}
```

### CLI Testing

Use `std::process::Command` to invoke binaries:

```rust
use std::process::Command;

#[test]
fn test_cli_subcommand() {
    let output = Command::new("forge")
        .args(["task", "list"])
        .env("CLAUDE_PROJECT_DIR", "/tmp/test-project")
        .output()
        .expect("failed to execute forge");
    assert!(output.status.success());
    let stdout = String::from_utf8_lossy(&output.stdout);
    assert!(stdout.contains("task"));
}
```

### API Testing

Use `reqwest` for HTTP integration testing:

```rust
#[tokio::test]
async fn test_api_endpoint() {
    let response = reqwest::get("http://localhost:8080/api/resource")
        .await
        .expect("request failed");
    assert_eq!(response.status(), reqwest::StatusCode::OK);
    let body: serde_json::Value = response.json().await.expect("parse failed");
    assert!(body["data"].is_object());
}
```

### Traceability

Each test function should include a traceability comment:

```rust
#[test]
fn test_login_with_valid_credentials() {
    // Traceability: TC-001 -> PRD User Auth section
}
```

## assertions

- **style**: macro-based
- **library**: Rust built-in test macros
- **custom_matchers**: none

### Key Macros

- `assert!(condition)` — boolean assertion
- `assert_eq!(left, right)` — equality check (implements `Debug` + `PartialEq`)
- `assert_ne!(left, right)` — inequality check
- `assert!(condition, "message: {}", value)` — assertion with custom message

### Assertion Details

```rust
assert!(value > 0, "value must be positive, got: {}", value);
assert_eq!(result, expected, "mismatch for input {}", input);
assert_ne!(status_code, 500, "server should not return 500");
```

**Rule**: Use built-in `assert!`, `assert_eq!`, `assert_ne!`. Do not use external assertion crates unless specifically required (e.g., `spectral`).

### Should-Panic Tests

Use `#[should_panic]` attribute for expected failures:

```rust
#[test]
#[should_panic(expected = "index out of bounds")]
fn test_out_of_bounds() {
    let v = vec![1, 2, 3];
    v[99];
}
```

## Tags

- **Format**: `#[ignore]` attribute for slow tests, feature flags for categorization
- **Running ignored tests**: `cargo test -- --ignored`
- **Filtering**: `cargo test test_name_pattern`

```rust
#[test]
#[ignore] // slow test — run with `cargo test -- --ignored`
fn test_large_dataset() {
    // ...
}
```

## Result Format

- **Output format**: `--format=plain` (default) or `--format=json` (unstable, nightly only)
- **Format type**: text (stdout)

### Text Output Example

```
running 3 tests
test test_something ... ok
test test_cli_command ... ok
test test_api_endpoint ... ok

test result: ok. 3 passed; 0 failed; 0 ignored; 0 measured; 0 filtered out
```

### Exit Codes

- `0` — all tests passed
- `101` — one or more tests failed

## Import Patterns

Standard imports for Rust integration tests:

```rust
use std::process::Command;
use std::path::PathBuf;
use std::env;
```

- HTTP tests add: `use reqwest;` with `tokio::test` macro
- File tests add: `use std::fs;`, `use std::path::Path;`
- Async tests use `#[tokio::test]` instead of `#[test]`

### Cargo.toml dev-dependencies

```toml
[dev-dependencies]
tokio = { version = "1", features = ["rt", "macros"] }
reqwest = { version = "0.12", features = ["json"] }
rstest = "0.19"
serde_json = "1"
tempfile = "3"
assert_cmd = "2"
```

## Anti-patterns (Forbidden)

| Forbidden | Replacement |
|-----------|-------------|
| `std::thread::sleep()` for synchronization | Retry loop with timeout |
| `unwrap()` in test body | `expect("descriptive message")` or explicit assertion |
| Hardcoded ports | `TcpListener::bind("127.0.0.1:0")` for dynamic port |
| Real secrets/tokens in code | `env::var("E2E_API_TOKEN")` |
| `println!` for debug output | Remove or use `eprintln!` |
| `#[ignore]` without comment | Implement properly or don't generate |
| Panicking in helper functions | Return `Result` and use `?` operator |

## Helpers

### run_cli helper

```rust
use std::process::Command;

struct CliResult {
    stdout: String,
    stderr: String,
    success: bool,
}

fn run_cli(args: &[&str]) -> CliResult {
    let output = Command::new("forge")
        .args(args)
        .output()
        .expect("failed to execute forge");
    CliResult {
        stdout: String::from_utf8_lossy(&output.stdout).to_string(),
        stderr: String::from_utf8_lossy(&output.stderr).to_string(),
        success: output.status.success(),
    }
}
```

### retry helper

```rust
use std::time::{Duration, Instant};

fn retry<F, R>(mut f: F, max_attempts: u32, interval: Duration) -> R
where
    F: FnMut() -> Result<R, String>,
{
    let start = Instant::now();
    for i in 0..max_attempts {
        match f() {
            Ok(result) => return result,
            Err(e) if i == max_attempts - 1 => panic!("retry exhausted after {:?}: {}", start.elapsed(), e),
            Err(_) => std::thread::sleep(interval),
        }
    }
    unreachable!()
}
```

### temp_project helper

```rust
use std::fs;
use std::path::PathBuf;
use tempfile::TempDir;

fn setup_test_project() -> (TempDir, PathBuf) {
    let dir = tempfile::tempdir().expect("failed to create temp dir");
    let forge_dir = dir.path().join(".forge");
    fs::create_dir_all(&forge_dir).expect("failed to create .forge dir");
    fs::write(forge_dir.join("config.yaml"), "{}").expect("failed to write config");
    (dir, dir.path().to_path_buf())
}
```
