---
status: "completed"
started: "2026-05-16 14:02"
completed: "2026-05-16 14:22"
time_spent: "~20m"
---

# Task Record: 1 Merge gen-test-scripts and run-e2e-tests into single quick mode task

## Summary
Merged gen-test-scripts and run-e2e-tests into a single gen-and-run task in quick mode. Added TypeTestPipelineGenAndRun constant, updated InferType mapping, restructured GetQuickTestTasks to produce 3 per-profile tasks (gen-cases, gen-and-run, graduate) instead of 4, renumbered shared tasks (T-quick-4=verify-regression, T-quick-5=drift), created merged prompt template, and updated all tests.

## Changes

### Files Created
- forge-cli/pkg/prompt/data/test-pipeline-gen-and-run.md

### Files Modified
- forge-cli/pkg/task/types.go
- forge-cli/pkg/task/infer.go
- forge-cli/pkg/task/testgen.go
- forge-cli/pkg/task/testgen_test.go
- forge-cli/pkg/task/infer_test.go
- forge-cli/pkg/task/types_test.go
- forge-cli/pkg/task/build_test.go
- forge-cli/pkg/prompt/prompt.go
- forge-cli/pkg/prompt/prompt_test.go
- forge-cli/tests/e2e/spec_drift_detection_cli_test.go
- forge-cli/docs/OVERVIEW.md
- plugins/forge/hooks/guide.md
- forge-cli/scripts/version.txt

### Key Decisions
- Used TypeTestPipelineGenAndRun = 'test-pipeline.gen-and-run' following existing TypeTestPipeline* naming convention
- Merged prompt template uses two-phase approach: Phase 1 (gen) then Phase 2 (run + fix loop)
- Renumbered shared tasks: T-quick-4=verify-regression, T-quick-5=drift-detection (was T-quick-5/T-quick-6)
- Per-type mode: each T-quick-2<L>-<type> is gen-and-run (not gen-only), with fan-in deps to graduate
- Breakdown mode completely untouched per Hard Rules

## Test Results
- **Tests Executed**: Yes
- **Passed**: 20
- **Failed**: 0
- **Coverage**: 90.3%

## Acceptance Criteria
- [x] New type constant TypeTestPipelineGenAndRun registered in types.go (const block + registry + validTypes)
- [x] GetQuickTestTasks() generates a merged gen-and-run task instead of separate gen-scripts + run tasks
- [x] Single profile: 5 test tasks total (gen-cases, gen-and-run, graduate, verify-regression, drift-detection)
- [x] Per-type mode: each per-type task is gen-and-run with fan-in deps to next task
- [x] Multi-profile: letter suffixes work correctly with merged task
- [x] infer.go maps merged IDs to TypeTestPipelineGenAndRun
- [x] prompt.go maps TypeTestPipelineGenAndRun to data/test-pipeline-gen-and-run.md
- [x] New prompt template calls /gen-test-scripts then /run-e2e-tests sequentially
- [x] resolveQuickDeps() dependency chain updated correctly
- [x] All existing quick mode tests updated and passing
- [x] New test cases cover: merged task count, merged task type, dependency chain, per-type merged tasks

## Notes
Breakdown mode task generation untouched. Existing prompt templates (gen-scripts.md, run.md) preserved for breakdown mode. E2e spec-drift-detection tests updated for renumbered task IDs.
