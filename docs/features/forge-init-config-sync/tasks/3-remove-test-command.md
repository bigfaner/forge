---
id: "3"
title: "Remove test-command from Config and refactor consumers"
priority: "P1"
estimated_time: "2h"
dependencies: []
scope: "backend"
breaking: true
type: "coding.cleanup"
mainSession: false
---

# 3: Remove test-command from Config and refactor consumers

## Description

`test-command` violates Forge's "just as unified test abstraction layer" design principle. Two consumers exist — `RunProjectTests` (uses it as bypass override) and `journey_isolation.go` (runs raw command). Both should go through `just` instead. Remove the field from all config types and refactor consumers.

## Reference Files
- `docs/proposals/forge-init-config-sync/proposal.md` — Source proposal
- `forge-cli/pkg/forgeconfig/config.go` — `TestCommand` field at line 148, `GetConfigValue` case at line 321
- `forge-cli/pkg/task/types.go` — `TestCommand` in `TaskIndex` at line 176, `taskIndexJSON` at line 191, Marshal/Unmarshal at lines 207/227
- `forge-cli/pkg/testrunner/testrunner.go` — `RunProjectTests` signature at line 59, `testCommand` parameter
- `forge-cli/internal/cmd/quality_gate.go` — `AllCompletedResult.TestCommand` at line 58, `runUnitTestStep` at line 282
- `forge-cli/internal/cmd/journey_isolation.go` — `JourneyExecutionConfig.TestCommand` at line 74, `readTestCommand` at line 91, `executeJourneyInIsolation` at line 153

## Acceptance Criteria
- [ ] `TestCommand` field removed from `forgeconfig.Config` struct
- [ ] `test-command` case removed from `forgeconfig.GetConfigValue`
- [ ] `TestCommand` field removed from `task.TaskIndex` and `taskIndexJSON`
- [ ] `MarshalJSON`/`UnmarshalJSON` in types.go no longer reference `TestCommand`
- [ ] `RunProjectTests` signature simplified: `func RunProjectTests(projectRoot string) (string, bool)` — always uses fallback chain (just → make → go → npm → pytest)
- [ ] `AllCompletedResult.TestCommand` removed, `runUnitTestStep` signature simplified
- [ ] `journey_isolation.go`: `readTestCommand` removed, `executeJourneyInIsolation` runs `just e2e-test` from project root with journey filter instead of raw test command
- [ ] All callers of modified functions updated
- [ ] All existing tests pass (update test fixtures and mocks)

## Hard Rules
- This is a breaking change to public struct fields and function signatures — all callers must be updated
- `executeJourneyInIsolation` must use `just e2e-test` from the **project root** (not the isolated temp dir), passing the journey filter
- `RunProjectTests` must always use the fallback chain — no special-case bypass

## Implementation Notes
- For `journey_isolation.go`: the current code runs `cfg.TestCommand` in the isolated temp dir. The new code should run `just e2e-test` from project root. Check how the journey filter is passed (likely via env var or argument).
- For `RunProjectTests`: when `testCommand` is empty, it already falls through to the `just test` path. Removing the parameter just eliminates the override shortcut.
- Check `all_completed.go` or any other file that constructs `AllCompletedResult` or calls `RunProjectTests`.
