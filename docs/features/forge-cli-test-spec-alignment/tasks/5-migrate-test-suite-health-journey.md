---
id: "5"
title: "Migrate test-suite-health Journey"
priority: "P1"
estimated_time: "1h"
dependencies: ["1"]
scope: "backend"
breaking: false
type: "coding.refactor"
mainSession: false
---

# 5: Migrate test-suite-health Journey

## Description

Migrate 3 test files into the `test-suite-health/` Journey directory. This Journey covers meta-testing: ensuring the test suite itself remains healthy and well-structured.

**Source → Target mapping:**
- `tests/e2e/e2e_test_quality_cleanup_cli_test.go` → `tests/test-suite-health/quality_cleanup_test.go`
- `tests/e2e/gen_journeys_skill_cli_test.go` → `tests/test-suite-health/gen_journeys_skill_test.go`
- `tests/e2e/simplify_e2e_tests_cli_test.go` → `tests/test-suite-health/simplify_test.go`
- New: `tests/test-suite-health/contracts/step-1-test-quality.md`
- New: `tests/test-suite-health/main_test.go`

## Reference Files
- `docs/proposals/forge-cli-test-spec-alignment/proposal.md` — Journey mapping table
- `tests/testkit/` — Shared infrastructure
- `tests/e2e/e2e_test_quality_cleanup_cli_test.go` — Source
- `tests/e2e/gen_journeys_skill_cli_test.go` — Source
- `tests/e2e/simplify_e2e_tests_cli_test.go` — Source

## Acceptance Criteria
- [ ] 3 test files migrated with correct package name and imports
- [ ] `tests/test-suite-health/contracts/` contains 1 Contract spec file
- [ ] `tests/test-suite-health/main_test.go` initializes binary via testkit
- [ ] All tests pass: `go test ./tests/test-suite-health/... -tags=e2e -count=1`

## Hard Rules
- Package name: `testsuitehealth`
- All test functions use `//go:build e2e` build tag

## Implementation Notes
- `gen_journeys_skill_cli_test.go` (16KB) is a substantial file — verify all imports are updated
- The "quality cleanup" tests are meta-tests that validate other tests' structure — ensure they still work after the directory reorganization
