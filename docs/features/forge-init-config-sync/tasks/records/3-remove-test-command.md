---
status: "completed"
started: "2026-05-20 17:10"
completed: "2026-05-20 17:37"
time_spent: "~27m"
---

# Task Record: 3 Remove test-command from Config and refactor consumers

## Summary
Removed test-command from all config types (forgeconfig.Config, task.TaskIndex) and refactored all consumers: RunProjectTests simplified to always use fallback chain (just->make->go->npm->pytest), quality_gate.go updated to remove TestCommand propagation, journey_isolation.go refactored to use just e2e-test from project root instead of raw test command, test.go/test_promote.go updated with new signatures. Updated schema, example YAML, docs (OVERVIEW.md, WORKFLOW.md), and all test files.

## Changes

### Files Created
无

### Files Modified
- forge-cli/pkg/forgeconfig/config.go
- forge-cli/pkg/task/types.go
- forge-cli/pkg/testrunner/testrunner.go
- forge-cli/internal/cmd/quality_gate.go
- forge-cli/internal/cmd/journey_isolation.go
- forge-cli/internal/cmd/test.go
- forge-cli/internal/cmd/test_promote.go
- forge-cli/internal/cmd/testdata/forge-config.schema.json
- forge-cli/internal/cmd/testdata/forge-config.example.yaml
- forge-cli/docs/OVERVIEW.md
- forge-cli/docs/WORKFLOW.md
- forge-cli/pkg/testrunner/testrunner_test.go
- forge-cli/pkg/task/types_test.go
- forge-cli/pkg/forgeconfig/config_test.go
- forge-cli/internal/cmd/quality_gate_test.go
- forge-cli/internal/cmd/journey_isolation_test.go
- forge-cli/internal/cmd/config_test.go
- forge-cli/internal/cmd/config_schema_test.go
- forge-cli/internal/docsync/docsync_test.go

### Key Decisions
- Removed testCommand bypass entirely from RunProjectTests - always uses just->make->go->npm->pytest fallback chain
- executeJourneyInIsolation now runs 'just e2e-test <journey>' from project root with FORGE_JOURNEY env var
- resolveJourneyExecutionConfig simplified to just store ProjectRoot (no longer reads config)
- Config schema and example YAML updated to remove test-command property
- Docs (OVERVIEW.md, WORKFLOW.md) updated to reflect 5-step detection order (was 6 with testCommand)

## Test Results
- **Tests Executed**: Yes
- **Passed**: 42
- **Failed**: 0
- **Coverage**: 88.1%

## Acceptance Criteria
- [x] TestCommand field removed from forgeconfig.Config struct
- [x] test-command case removed from forgeconfig.GetConfigValue
- [x] TestCommand field removed from task.TaskIndex and taskIndexJSON
- [x] MarshalJSON/UnmarshalJSON in types.go no longer reference TestCommand
- [x] RunProjectTests signature simplified to func(projectRoot string) (string, bool)
- [x] AllCompletedResult.TestCommand removed, runUnitTestStep signature simplified
- [x] journey_isolation.go: readTestCommand removed, executeJourneyInIsolation runs just e2e-test
- [x] All callers of modified functions updated
- [x] All existing tests pass

## Notes
Breaking change to public struct fields and function signatures. All callers updated across forgeconfig, task, testrunner, cmd packages.
