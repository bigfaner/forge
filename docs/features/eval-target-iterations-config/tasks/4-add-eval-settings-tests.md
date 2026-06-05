---
id: "4"
title: "Add EvalSettings unit tests"
priority: "P1"
estimated_time: "1h"
complexity: "medium"
dependencies: [1, 2]
surface-key: ""
surface-type: "cli"
breaking: false
type: "coding.feature"
mainSession: false
---

# 4: Add EvalSettings unit tests

## Description

Add comprehensive unit tests for the new `EvalSettings` config block: get/set operations on eval paths, nil pointer fallback behavior, and init default value generation. Verify no regression in existing config tests.

## Reference Files
- `docs/proposals/eval-target-iterations-config/proposal.md` — Success Criteria (test coverage expectations)
- `forge-cli/pkg/forgeconfig/config_test.go` — Existing config tests (pattern reference)
- `forge-cli/internal/cmd/config_test.go` — Existing CLI config tests (pattern reference)

## Acceptance Criteria
- [ ] Tests for `GetConfigValue` on eval paths: returns correct value when configured, returns error when nil (`*int` not set)
- [ ] Tests for `SetConfigValue` on eval paths: sets target and iterations correctly, writes valid YAML
- [ ] Tests for `forge config init` eval block: generated config contains all 7 types with rubric-default values
- [ ] Existing config_test.go and config_test.go (cmd) tests pass with no regression

## Implementation Notes
- Follow existing test patterns in config_test.go (table-driven tests, temp directory fixtures).
- Key test scenarios: nil `*int` → errKeyNotFound, set `*int` → correct value, partial config (only target set, iterations nil), full config (both set).
- Init tests may need to verify the generated YAML output contains the expected eval block structure.
