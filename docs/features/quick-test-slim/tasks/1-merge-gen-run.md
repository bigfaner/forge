---
id: "1"
title: "Merge gen-test-scripts and run-e2e-tests into single quick mode task"
priority: "P1"
estimated_time: "2h"
dependencies: []
scope: "backend"
breaking: true
type: "implementation"
mainSession: false
---

# 1: Merge gen-test-scripts and run-e2e-tests into single quick mode task

## Description

Quick mode test pipeline has 6 test tasks for 3-5 business tasks — test overhead exceeds business work. The gen-test-scripts and run-e2e-tests steps are naturally sequential and should be merged into a single task so the subagent can generate scripts, run them, and fix failures in one session without context rebuild.

Core change: replace separate T-quick-2 (gen-scripts) and T-quick-3 (run) with a single merged T-quick-2 (gen-and-run), then renumber subsequent tasks.

## Reference Files
- `docs/proposals/quick-test-slim/proposal.md` — Source proposal
- `forge-cli/pkg/task/testgen.go` — GetQuickTestTasks(), resolveQuickDeps()
- `forge-cli/pkg/task/infer.go` — InferType() ID-to-type mapping
- `forge-cli/pkg/task/types.go` — Type constants and TaskTypeRegistry
- `forge-cli/pkg/prompt/prompt.go` — typeToTemplate map
- `forge-cli/pkg/prompt/data/test-pipeline-gen-scripts.md` — Existing gen prompt (reference)
- `forge-cli/pkg/prompt/data/test-pipeline-run.md` — Existing run prompt (reference)
- `forge-cli/pkg/task/testgen_test.go` — Existing quick mode tests

## Acceptance Criteria
- [ ] New type constant `TypeTestPipelineGenAndRun` = `"test-pipeline.gen-and-run"` registered in types.go (const block + registry + validTypes)
- [ ] `GetQuickTestTasks()` generates a merged gen-and-run task instead of separate gen-scripts + run tasks
- [ ] Single profile (no per-type): 5 test tasks total (gen-cases, gen-and-run, graduate, verify-regression, drift-detection)
- [ ] Per-type mode: each per-type task (e.g. T-quick-2-tui) is gen-and-run (not gen-only), with fan-in deps to next task
- [ ] Multi-profile: letter suffixes work correctly with merged task
- [ ] `infer.go` maps merged IDs to `TypeTestPipelineGenAndRun` (handles profile suffix and type suffix patterns)
- [ ] `prompt.go` maps `TypeTestPipelineGenAndRun` → `data/test-pipeline-gen-and-run.md`
- [ ] New prompt template `test-pipeline-gen-and-run.md` calls `/gen-test-scripts` then `/run-e2e-tests` sequentially, with in-session fix loop
- [ ] `resolveQuickDeps()` dependency chain updated: verify-regression depends on the merged task (or graduate if present), drift-detection depends on verify-regression
- [ ] All existing quick mode tests in testgen_test.go updated and passing
- [ ] New test cases cover: merged task count (single profile = 5), merged task type, dependency chain after merge, per-type merged tasks

## Hard Rules
- Do NOT modify breakdown mode task generation — this change is quick-mode only
- Do NOT modify the gen-test-scripts or run-e2e-tests skills themselves
- New type constant must follow existing naming convention (TypeTestPipeline*)
- Prompt template must use standard placeholders: {{TASK_ID}}, {{TASK_FILE}}, {{SCOPE}}, {{FEATURE_SLUG}}, {{PROFILE}}, {{TEST_TYPE_ARG}}

## Implementation Notes
- Start from types.go (add constant), then infer.go (add ID pattern), then testgen.go (merge logic + deps), then prompt template, then tests
- The per-type block in GetQuickTestTasks() (line ~127) fans out gen-scripts tasks per type — after merge, these become gen-and-run tasks per type
- resolveQuickDeps() calculates block sizes for dependency resolution — the block size changes when gen+run are merged (reduce by 1 per profile)
- Existing prompt templates (gen-scripts.md, run.md) remain for breakdown mode — only quick mode uses the merged template
- For context efficiency, the merged prompt should be structured as two sequential phases: Phase 1 (gen) then Phase 2 (run + fix)
