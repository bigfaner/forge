---
status: "completed"
started: "2026-05-14 18:12"
completed: "2026-05-14 18:20"
time_spent: "~8m"
---

# Task Record: 1 Extend BuildIndex with phase detection and stage-gate generation

## Summary
Extended BuildIndex with phase detection and stage-gate generation. Created stage_gates.go with DetectPhases (groups task IDs by phase number using strings.Split+strconv.Atoi), GenerateSummaryMD/GenerateGateMD (programmatic Go string literals, no go:embed), and GenerateStageGates (idempotent file generation with skip-if-exists). Integrated into build.go between orphan detection and test-task generation. Stage-gate generation is independent of --no-test flag (structural, not test tasks). Added StageGatesGenerated to BuildIndexResult. Unit tests cover: phase detection, test-task exclusion, template generation, idempotency, malformed IDs, partial state, threshold check, dependency ordering, sorted output.

## Changes

### Files Created
- forge-cli/pkg/task/stage_gates.go
- forge-cli/pkg/task/stage_gates_test.go

### Files Modified
- forge-cli/pkg/task/build.go
- forge-cli/pkg/task/build_test.go

### Key Decisions
- Generated MD content as Go string literals (not go:embed) per Hard Rule: go:embed cannot reach plugins/ from forge-cli/pkg/task/
- Phase detection uses strings.Split(id, '.') yielding exactly 2 segments parsed via strconv.Atoi per Hard Rule
- Stage-gate generation runs regardless of --no-test flag per Hard Rule: gates are structural, not test tasks
- Dependencies in generated summaries are sorted for deterministic output regardless of input order
- Re-scan directory after generation to index newly created files in the same BuildIndex run

## Test Results
- **Tests Executed**: Yes
- **Passed**: 26
- **Failed**: 0
- **Coverage**: 90.2%

## Acceptance Criteria
- [x] forge task index generates .summary and .gate for every numbered phase with >=2 business tasks
- [x] Single-task phases (<2 business tasks) are skipped
- [x] T-test/T-quick task IDs excluded from phase business task count
- [x] Generated .gate has depends_on: [N.summary] and breaking: true
- [x] Generated .summary has depends_on set to all business task IDs in the same phase
- [x] Re-running forge task index does not overwrite existing files (idempotent)
- [x] Partial state handled: if .summary exists but .gate missing, only .gate is generated
- [x] Malformed task IDs silently skipped, no crash
- [x] Generated tasks appear in index.json with correct type (gate / doc-generation.summary)
- [x] Unit tests cover: phase detection, test-task exclusion, template generation, idempotency, malformed IDs, partial state, threshold check

## Notes
All existing tests continue to pass. Coverage 90.2% for pkg/task. Total 26 new tests (16 unit + 5 integration + 5 threshold/type).
