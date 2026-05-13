/// Shared test helpers for Rust E2E tests.
/// Place in tests/e2e/helpers.rs

use std::process::{Command, Output};
use std::time::Duration;

/// Execute a CLI command and return output
pub fn run_cli(args: &[&str]) -> Output {
    Command::new(args[0])
        .args(&args[1..])
        .output()
        .expect("failed to execute command")
}

/// Retry a function until it succeeds or max retries exceeded
pub fn with_retry<F, R>(mut f: F, max_retries: usize, delay: Duration) -> R
where
    F: FnMut() -> Result<R, String>,
{
    let mut last_err = String::new();
    for _ in 0..max_retries {
        match f() {
            Ok(r) => return r,
            Err(e) => {
                last_err = e;
                std::thread::sleep(delay);
            }
        }
    }
    panic!("retry exhausted after {} attempts: {}", max_retries, last_err);
}
