---
step: 2
title: E2E Error Handling and Verify
journey: e2e-pipeline
---

# Step 2: E2E Error Handling and Verify

## Given
- A forge project directory
- Just may or may not be on PATH
- E2e test files may contain VERIFY markers

## When
- E2E commands are run without 'just' on PATH
- `forge e2e verify --feature <name>` scans test files
- E2E commands are run without a configured profile

## Then
- Missing 'just' returns "'just' is required but not found on PATH" for run/setup/compile/discover
- Non-zero just exit propagates as error
- Zero just exit returns nil error
- verify scans files locally without invoking just
- verify detects VERIFY markers and reports them
- No profile returns "no e2e profile configured" error

## Contract Dimensions
- **Actor**: CLI user in various error scenarios
- **Input**: Environment state (PATH, config), test files with optional VERIFY markers
- **Output**: Error messages, exit codes
- **Error Cases**: just missing, just fails, no profile, feature not found, VERIFY markers found
- **Invariants**: exit codes propagated from just; verify is purely local file scanning
