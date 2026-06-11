---
feature: "test-scripts-per-type"
sources:
  - docs/proposals/test-scripts-per-type/proposal.md
generated: "2026-05-15"
---

# Test Cases: test-scripts-per-type

## Summary

| Type | Count |
|------|-------|
| CLI  | 12   |
| **Total** | **12** |

> **Note**: This feature is a CLI-only project (forge CLI tool). No UI or API interfaces are exposed by the product. The go-test profile has [tui, api, cli] capabilities, but only CLI is a product interface.

---

## CLI Test Cases

## TC-001: gen-test-scripts accepts --type cli and generates only CLI scripts
- **Source**: Proposal Success Criterion 1 -- "gen-test-scripts accepts --type <ui|api|cli> and generates only scripts for that type"
- **Type**: CLI
- **Target**: cli/gen-test-scripts
- **Test ID**: cli/gen-test-scripts/type-filter-cli
- **Pre-conditions**: A feature with test-cases.md containing CLI-type test cases exists; go-test profile is active
- **Steps**:
  1. Run `forge gen-test-scripts --type cli` for the feature
  2. Check generated files in the test output directory
- **Expected**: Only CLI test script files are generated (e.g., `*_cli_test.go`). No API or TUI test scripts are generated.
- **Priority**: P0

## TC-002: gen-test-scripts accepts --type api and generates only API scripts
- **Source**: Proposal Success Criterion 1 -- "gen-test-scripts accepts --type <ui|api|cli> and generates only scripts for that type"
- **Type**: CLI
- **Target**: cli/gen-test-scripts
- **Test ID**: cli/gen-test-scripts/type-filter-api
- **Pre-conditions**: A feature with test-cases.md containing API-type test cases exists; a profile with api capability is active
- **Steps**:
  1. Run `forge gen-test-scripts --type api` for the feature
  2. Check generated files in the test output directory
- **Expected**: Only API test script files are generated (e.g., `*_api_test.go`). No CLI or TUI test scripts are generated.
- **Priority**: P0

## TC-003: gen-test-scripts accepts --type tui and generates only TUI scripts
- **Source**: Proposal Success Criterion 1 -- "gen-test-scripts accepts --type <ui|api|cli> and generates only scripts for that type"
- **Type**: CLI
- **Target**: cli/gen-test-scripts
- **Test ID**: cli/gen-test-scripts/type-filter-tui
- **Pre-conditions**: A feature with test-cases.md containing TUI-type test cases exists; a profile with tui capability is active
- **Steps**:
  1. Run `forge gen-test-scripts --type tui` for the feature
  2. Check generated files in the test output directory
- **Expected**: Only TUI test script files are generated (e.g., `*_tui_test.go`). No CLI or API test scripts are generated.
- **Priority**: P0

## TC-004: breakdown-tasks creates per-type gen-scripts tasks
- **Source**: Proposal Success Criterion 2 -- "breakdown-tasks creates separate tasks per detected test type instead of one T-test-2"
- **Type**: CLI
- **Target**: cli/breakdown-tasks
- **Test ID**: cli/breakdown-tasks/per-type-tasks
- **Pre-conditions**: A feature with a finalized tech design; profile has multiple capabilities (e.g., cli + api); test-cases.md contains multiple test types
- **Steps**:
  1. Run `forge breakdown-tasks` for the feature
  2. Inspect generated task files in the tasks directory
- **Expected**: Separate task files are created for each detected test type (e.g., T-test-2-cli, T-test-2-api) instead of a single T-test-2.
- **Priority**: P0

## TC-005: breakdown-tasks creates only tasks for types with test cases
- **Source**: Proposal Key Risk -- "Tasks generated for types with no test cases" and Scope item "Only create per-type task when test-cases.md contains cases of that type"
- **Type**: CLI
- **Target**: cli/breakdown-tasks
- **Test ID**: cli/breakdown-tasks/no-empty-type-tasks
- **Pre-conditions**: A feature where test-cases.md contains only CLI-type test cases; profile has cli + api capabilities
- **Steps**:
  1. Run `forge breakdown-tasks` for the feature
  2. Inspect generated task files
- **Expected**: Only a T-test-2-cli task is created. No T-test-2-api task is created since no API test cases exist.
- **Priority**: P1

## TC-006: quick-tasks creates per-type gen-scripts tasks
- **Source**: Proposal Success Criterion 3 -- "quick-tasks creates separate tasks per detected test type instead of one T-quick-2"
- **Type**: CLI
- **Target**: cli/quick-tasks
- **Test ID**: cli/quick-tasks/per-type-tasks
- **Pre-conditions**: A feature proposal with acceptance criteria covering multiple test types; profile has multiple capabilities
- **Steps**:
  1. Run `forge quick-tasks` for the feature
  2. Inspect generated task files in the tasks directory
- **Expected**: Separate task files are created for each detected test type (e.g., T-quick-2-cli, T-quick-2-api) instead of a single T-quick-2.
- **Priority**: P0

## TC-007: T-test-3 depends on all per-type T-test-2 tasks
- **Source**: Proposal Success Criterion 4 -- "T-test-3 depends on all per-type T-test-2-* tasks completing successfully"
- **Type**: CLI
- **Target**: cli/breakdown-tasks
- **Test ID**: cli/breakdown-tasks/test-3-dependencies
- **Pre-conditions**: A feature with test-cases.md containing CLI and API test cases; breakdown-tasks has been run
- **Steps**:
  1. Run `forge breakdown-tasks` for the feature
  2. Read the T-test-3 task file
  3. Check the dependencies field
- **Expected**: T-test-3's dependencies field lists all per-type T-test-2 tasks (e.g., T-test-2-cli, T-test-2-api).
- **Priority**: P0

## TC-008: T-quick-3 depends on all per-type T-quick-2 tasks
- **Source**: Proposal Scope item -- "Update T-test-3 / T-quick-3 dependencies to depend on ALL per-type gen tasks"
- **Type**: CLI
- **Target**: cli/quick-tasks
- **Test ID**: cli/quick-tasks/quick-3-dependencies
- **Pre-conditions**: A feature proposal with multiple test types; quick-tasks has been run
- **Steps**:
  1. Run `forge quick-tasks` for the feature
  2. Read the T-quick-3 task file
  3. Check the dependencies field
- **Expected**: T-quick-3's dependencies field lists all per-type T-quick-2 tasks (e.g., T-quick-2-cli, T-quick-2-api).
- **Priority**: P0

## TC-009: Failed gen task can be independently retried
- **Source**: Proposal Success Criterion 5 -- "Failed gen task can be independently retried without affecting other type tasks"
- **Type**: CLI
- **Target**: cli/gen-test-scripts
- **Test ID**: cli/gen-test-scripts/independent-retry
- **Pre-conditions**: A feature with test-cases.md containing CLI and API test cases; CLI scripts generated successfully; API script generation fails
- **Steps**:
  1. Run `forge gen-test-scripts --type cli` -- succeeds
  2. Run `forge gen-test-scripts --type api` -- fails
  3. Verify CLI scripts still exist and are unchanged
  4. Fix the API issue and retry `forge gen-test-scripts --type api`
- **Expected**: CLI scripts remain intact after API failure. Retrying only API generation succeeds without affecting CLI scripts.
- **Priority**: P0

## TC-010: Shared infrastructure generated idempotently
- **Source**: Proposal Success Criterion 6 -- "Shared infrastructure files are generated idempotently (no conflicts from parallel gen)"
- **Type**: CLI
- **Target**: cli/gen-test-scripts
- **Test ID**: cli/gen-test-scripts/idempotent-infra
- **Pre-conditions**: A feature with multiple test types; shared infrastructure files (helpers, config) do not yet exist
- **Steps**:
  1. Run `forge gen-test-scripts --type cli` -- creates shared infrastructure
  2. Run `forge gen-test-scripts --type api` -- should reuse existing infrastructure
  3. Verify shared infrastructure files are unchanged after second run
- **Expected**: Shared infrastructure files (e.g., helpers.go, main_test.go) are created on first run and not overwritten or corrupted on subsequent per-type runs.
- **Priority**: P1

## TC-011: Single type project creates only one gen task
- **Source**: Proposal Key Scenario -- "Single type project: Only API endpoints, no UI. Only T-test-2-api is created."
- **Type**: CLI
- **Target**: cli/breakdown-tasks
- **Test ID**: cli/breakdown-tasks/single-type-project
- **Pre-conditions**: A feature where profile has only CLI capability; test-cases.md contains only CLI test cases
- **Steps**:
  1. Run `forge breakdown-tasks` for the feature
  2. Inspect generated task files
- **Expected**: Only T-test-2-cli is created. No T-test-2-api or T-test-2-tui tasks.
- **Priority**: P1

## TC-012: gen-test-scripts without --type generates all types
- **Source**: Proposal constraint -- backward compatibility; existing behavior must continue for projects that don't opt in
- **Type**: CLI
- **Target**: cli/gen-test-scripts
- **Test ID**: cli/gen-test-scripts/no-type-all-types
- **Pre-conditions**: A feature with test-cases.md containing CLI and API test cases; go-test profile active
- **Steps**:
  1. Run `forge gen-test-scripts` without --type flag
  2. Check generated files
- **Expected**: All test types present in test-cases.md are generated (CLI and API scripts), preserving backward compatibility.
- **Priority**: P1

---

## Traceability

| TC ID | Source | Type | Target | Priority |
|-------|--------|------|--------|----------|
| TC-001 | Proposal SC-1 | CLI | cli/gen-test-scripts | P0 |
| TC-002 | Proposal SC-1 | CLI | cli/gen-test-scripts | P0 |
| TC-003 | Proposal SC-1 | CLI | cli/gen-test-scripts | P0 |
| TC-004 | Proposal SC-2 | CLI | cli/breakdown-tasks | P0 |
| TC-005 | Proposal Risk / Scope | CLI | cli/breakdown-tasks | P1 |
| TC-006 | Proposal SC-3 | CLI | cli/quick-tasks | P0 |
| TC-007 | Proposal SC-4 | CLI | cli/breakdown-tasks | P0 |
| TC-008 | Proposal Scope | CLI | cli/quick-tasks | P0 |
| TC-009 | Proposal SC-5 | CLI | cli/gen-test-scripts | P0 |
| TC-010 | Proposal SC-6 | CLI | cli/gen-test-scripts | P1 |
| TC-011 | Proposal Key Scenario | CLI | cli/breakdown-tasks | P1 |
| TC-012 | Proposal NFR (backward compat) | CLI | cli/gen-test-scripts | P1 |
