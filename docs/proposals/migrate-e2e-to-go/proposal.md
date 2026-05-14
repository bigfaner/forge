---
created: 2026-05-14
author: "fanhuifeng"
status: Draft
---

# Proposal: Migrate E2E Tests from Playwright/Node.js to Go

## Problem

The project maintains two test stacks: Playwright (Node.js/TypeScript) for e2e tests in `tests/e2e/` and Go for unit/integration tests in `forge-cli/`. This dual-stack requires maintaining separate dependencies (package.json vs go.mod), CI configurations, and developer toolchains for what is functionally the same type of CLI e2e testing.

### Evidence

- `tests/e2e/` contains 12 Playwright `.spec.ts` files with ~152 test cases
- `forge-cli/tests/e2e/` already has Go e2e tests using `testkit` helpers with proven patterns
- All Playwright tests are CLI execution + file content assertion — no browser UI interaction
- The project's primary language is Go (forge-cli is a Go binary); Node.js is used solely for testing

### Urgency

The Node.js test stack is the sole remaining Node.js dependency in the project. Eliminating it removes a whole runtime from CI, reduces Docker image size, and simplifies the developer setup to just Go.

## Proposed Solution

1:1 direct translation of all Playwright test cases to Go, following the existing `forge-cli/tests/e2e/` conventions (build tag `//go:build e2e`, `testkit` helpers, `stretchr/testify` assertions, `TestTC_NNN_Description` naming). After verification, remove the entire Node.js test infrastructure.

### Innovation Highlights

Straightforward migration — the target patterns are already proven in the codebase. The innovation is eliminating the Node.js runtime dependency entirely rather than maintaining two stacks.

## Requirements Analysis

### Key Scenarios

- Each Playwright test case maps to one Go test function with identical TC number and assertion logic
- `runCli()` calls become `testkit.RunCLIExitCode()` / `testkit.RunCLIWithResult()`
- `fileContains()` / `fileNotContains()` helpers are added to `testkit`
- `beforeAll`/`afterAll` fixture setup maps to `t.TempDir()` + `require.NoError`
- Feature-scoped tests go to `forge-cli/tests/e2e/features/<slug>/`
- Cross-cutting tests go to `forge-cli/tests/e2e/<area>_cli_test.go`

### Non-Functional Requirements

- All converted tests must pass: `go test ./tests/e2e/... -v -tags=e2e`
- No change to test coverage scope — every Playwright test case has a Go equivalent

### Constraints & Dependencies

- Depends on existing `testkit` package (`forge-cli/tests/e2e/testkit/`)
- Depends on `stretchr/testify` (already in go.mod)
- `forge` binary must be installed (`forge init-forge`) for `testkit.RunCLI*` to work

## Alternatives & Industry Benchmarking

### Industry Solutions

Polyglot test stacks are common in microservices but unnecessary for single-language CLI tools. Industry standard is to keep tests in the same language as the code under test.

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | Zero migration cost | Dual stack maintenance, Node.js dependency for Go project | Rejected: ongoing cost |
| Gradual migration (add Go, keep both) | Common migration pattern | Low risk, rollback easy | Extended dual-stack period, CI runs both | Rejected: unnecessary complexity |
| **Full 1:1 translation + remove Node.js** | Standard practice | Clean single stack, no ongoing dual maintenance | One-time migration effort | **Selected: unified tech stack** |

## Feasibility Assessment

### Technical Feasibility

All Playwright tests are CLI + file assertion. No browser automation. Go `os/exec` + `testkit` covers all patterns.

### Resource & Timeline

~152 test cases across 12 files. Each follows one of ~3 patterns (CLI exit code, file content, fixture setup). Estimated 1-2 hours for experienced developer with AI assistance.

### Dependency Readiness

`testkit` package exists and is proven. `stretchr/testify` is already a dependency. No external blockers.

## Scope

### In Scope

- Add `fileContains` / `fileNotContains` helpers to `testkit` package
- Convert 12 Playwright test files to Go (1:1, preserving TC numbers):
  1. `gen-test-scripts/cli.spec.ts` (15 tests)
  2. `justfile-execution/justfile-execution.spec.ts` (9 tests)
  3. `task-cli/typed-task-dispatch.spec.ts` (20 tests)
  4. `scope-resolution/scope-resolution.spec.ts` (8 tests)
  5. `justfile-e2e-integration/forge-justfile.spec.ts` (15 tests)
  6. `justfile-e2e-integration/detection-assembly.spec.ts` (19 tests)
  7. `justfile-e2e-integration/mixed-template.spec.ts` (23 tests)
  8. `justfile-e2e-integration/cli.spec.ts` (20 tests)
  9. `init-justfile/init-justfile.spec.ts` (7 tests)
  10. `plugin-content/skill-content.spec.ts` (1 test)
  11. `features/forge-testing-optimization/cli.spec.ts` (15 tests — duplicate of #1, merge into #1)
- Remove Node.js test infrastructure:
  - `tests/e2e/package.json`
  - `tests/e2e/playwright.config.ts`
  - `tests/e2e/tsconfig.json`
  - `tests/e2e/helpers.ts`
  - All `.spec.ts` files
  - `node_modules/` (if present)

### Out of Scope

- Refactoring test structure (table-driven tests, etc.) — can be done incrementally after migration
- Changing `go-test` profile templates or `gen-test-scripts` skill
- Updating justfile recipes (they already support Go e2e tests)
- Template files in `plugins/forge/skills/gen-test-scripts/templates/` — these are for generating new tests, not existing ones

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| Test behavior divergence during translation | M | H | Run both suites side-by-side before removing Playwright |
| Missing helper functions in testkit | L | L | Add as needed during conversion |
| Feature-scoped test path mismatch | L | M | Follow existing `features/<slug>/` convention from go-test profile |

## Success Criteria

- [ ] All ~152 test cases have corresponding Go test functions with matching TC numbers
- [ ] `go test ./tests/e2e/... -v -tags=e2e` passes with 0 failures
- [ ] No `.spec.ts` files remain in `tests/e2e/`
- [ ] No `package.json` or `node_modules` in `tests/e2e/`
- [ ] `forge-cli/tests/e2e/testkit/` exports new file content helpers

## Next Steps

- Proceed to `/quick-tasks` to generate task files for execution
