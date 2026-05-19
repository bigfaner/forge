/// E2E test template for Rust.
/// Place in tests/e2e/features/<slug>/e2e_tests.rs
/// Build tag: run with `cargo test --test e2e`

use std::process::Command;

/// Traceability: TC-NNN → {PRD Source}
#[test]
fn test_tc_nnn_description() {
    // Step 1: Setup
    // Step 2: Execute
    let output = Command::new("./target/debug/binary")
        .args(&["--flag", "value"])
        .output()
        .expect("failed to execute binary");

    // Expected: ...
    assert!(output.status.success(), "command failed: {}", String::from_utf8_lossy(&output.stderr));
    assert!(String::from_utf8_lossy(&output.stdout).contains("expected text"));
}
