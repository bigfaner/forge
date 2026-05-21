# Contract: justfile-integration / Step 3: Execution

## Outcome "compile-passes"
- Preconditions: "just command available in PATH, toolchains may or may not be available"
- Input: "runJust compile"
- Output: "No scope error in output; exit code 0 if toolchain available, non-zero allowed for missing toolchain"
- State: "compile recipe execution verified"
- Side-effect: none

## Outcome "compile-error-propagation"
- Preconditions: "forge project justfile exists"
- Input: "Read justfile content"
- Output: "Compile recipe contains set -euo pipefail for error propagation"
- State: "error propagation mechanism verified"
- Side-effect: none

## Outcome "compile-stderr-output"
- Preconditions: "forge project justfile exists"
- Input: "Read compile recipe section"
- Output: "Compile recipe uses set -euo pipefail"
- State: "stderr output mechanism verified"
- Side-effect: none

## Outcome "consecutive-commands"
- Preconditions: "just command available, toolchains may be missing"
- Input: "runJust install then runJust compile"
- Output: "Both exit 0 if toolchains available; test skips if toolchains missing"
- State: "command chaining verified"
- Side-effect: none

## Outcome "build-invalid-scope"
- Preconditions: "forge project justfile exists with build scope parameter"
- Input: "runJust build"
- Output: "Build recipe exists with scope parameter; exit 1 on failure"
- State: "build scope behavior verified"
- Side-effect: none

## Outcome "project-type-deterministic"
- Preconditions: "forge CLI binary built and available"
- Input: "forge probe twice"
- Output: "Both calls return exit code 0, same trimmed output"
- State: "deterministic output verified"
- Side-effect: none

## Outcome "idempotent-recipes"
- Preconditions: "just command available"
- Input: "runJust install twice, runJust install-forge twice"
- Output: "Second run exits 0 if first run succeeded"
- State: "idempotency verified"
- Side-effect: none

## Outcome "scope-dispatch-backend"
- Preconditions: "forge project justfile exists"
- Input: "Read test and build recipes"
- Output: "Test recipe has go test, build recipe has go build for backend scope"
- State: "backend scope dispatch verified"
- Side-effect: none

## Outcome "scope-dispatch-frontend"
- Preconditions: "forge project justfile exists"
- Input: "Read build recipe"
- Output: "Build has go build; if frontend branch exists, it has npm run build"
- State: "frontend scope dispatch verified (conditional on project type)"
- Side-effect: none

## Journey Invariants
- runJust helper consistent across all execution tests
- testkit paths resolve correctly from justfile-integration package location
