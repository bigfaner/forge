---
id: "9"
title: "Migrate forge-cli/tests/e2e/ Journeys: test-generation + spec-drift"
priority: "P1"
estimated_time: "1h"
dependencies: []
scope: "backend"
breaking: false
type: "coding.refactor"
mainSession: false
---

# 9: Migrate forge-cli/tests/e2e/ Journeys: test-generation + spec-drift

## Description

Migrate 2 test files into 2 Journey directories.

**test-generation Journey:**
- `gen_test_scripts_cli_test.go` → `forge-cli/tests/test-generation/gen_test_scripts_test.go`
- Contracts: `forge-cli/tests/test-generation/contracts/step-1-validate-specs.md`, `step-2-gen-scripts.md`

**spec-drift Journey:**
- `spec_drift_detection_cli_test.go` → `forge-cli/tests/spec-drift/spec_drift_detection_test.go`
- Contracts: `forge-cli/tests/spec-drift/contracts/step-1-drift-type.md`, `step-2-detect.md`

## Reference Files
- `forge-cli/tests/e2e/testkit/` — Already existing testkit
- All source files listed above

## Acceptance Criteria
- [ ] 2 Journey directories with 2 migrated test files
- [ ] Each Journey has `contracts/` with spec files and `main_test.go` via testkit
- [ ] Tests pass: `go test ./forge-cli/tests/test-generation/... ./forge-cli/tests/spec-drift/... -tags=e2e -count=1`

## Hard Rules
- Package names: `testgeneration`, `specdrift`
- All test functions use `//go:build e2e` build tag

## Implementation Notes
- gen_test_scripts_cli_test.go validates spec rules and script generation — uses forge binary
- spec_drift_detection_cli_test.go has 20 tests covering drift type, strategy template, and pipeline integration
