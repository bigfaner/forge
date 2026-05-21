---
step: 1
title: E2E Command Delegation
journey: e2e-pipeline
---

# Step 1: E2E Command Delegation

## Given
- A forge project directory with .forge/config.yaml (with e2e profile)
- Optionally 'just' installed on PATH

## When
- `forge e2e run` is executed
- `forge e2e run --feature <name>` is executed
- `forge e2e setup` is executed
- `forge e2e compile` is executed
- `forge e2e discover` is executed

## Then
- run delegates to `just test-e2e`
- run --feature passes feature as argument to just
- setup delegates to `just e2e-setup`
- compile delegates to `just e2e-compile`
- discover delegates to `just e2e-discover`
- verify does NOT delegate to just (scans files locally)
- Missing 'just' returns actionable error with install instructions

## Contract Dimensions
- **Actor**: CLI user running e2e pipeline commands
- **Input**: .forge/config.yaml with profile, optional feature flag
- **Output**: CLI stdout/stderr, exit codes propagated from just
- **Error Cases**: just not on PATH -> actionable error; no profile configured -> error
- **Invariants**: verify never delegates to just; all other commands require just
