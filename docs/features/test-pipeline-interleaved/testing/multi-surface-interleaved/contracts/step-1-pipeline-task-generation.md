---
journey: "multi-surface-interleaved"
step: 1
step-action: "Run pipeline task generation for a multi-surface feature"
generated: "2026-06-08"
sources:
  - docs/features/test-pipeline-interleaved/testing/multi-surface-interleaved/journey.md
skip_eval: true
---

# Contract: multi-surface-interleaved / Step 1: Run pipeline task generation for a multi-surface feature

> **Note**: Contracts generated without eval-journey verification (SKIP_EVAL_GATE=true). Review with extra scrutiny.

<!-- gen-contracts: do not edit manually. Regenerate via /gen-contracts. -->

## Outcome "success"
- Preconditions: "Project has at least two configured surfaces with execution_order defined; feature has finalized documents; no pre-existing test tasks for the feature"
- Input: "forge task index --feature <slug> for a multi-surface feature"
- Output: "Test tasks generated for each surface in execution_order: T-test-gen-scripts-{surface-1}, T-test-run-{surface-1}, T-test-gen-scripts-{surface-2}, T-test-run-{surface-2}. Dependency chain follows interleaved pattern: gen-scripts-{surface-1} -> run-{surface-1} -> gen-scripts-{surface-2} -> run-{surface-2}"
- State: "Task index file contains all four test tasks with correct dependencies. Each task has surface-key and surface-type fields populated"
- Side-effect: "Task .md files written to docs/features/<slug>/tasks/ for each generated test task"

## Outcome "three-surface-chain"
- Preconditions: "Project has three configured surfaces (e.g., api, web, cli) with execution_order: [api, web, cli]"
- Input: "forge task index --feature <slug> for the three-surface feature"
- Output: "Six test tasks generated with chain: gen-scripts-api -> run-api -> gen-scripts-web -> run-web -> gen-scripts-cli -> run-cli. Each N-th gen-scripts depends on (N-1)-th run-tests"
- State: "Task index contains six test tasks. Interleaving extends naturally to all three surfaces"
- Side-effect: "Task .md files written for all six tasks"

## Outcome "no-execution-order-defaults"
- Preconditions: "Project has multiple surfaces configured but execution_order is not explicitly defined in config"
- Input: "forge task index --feature <slug>"
- Output: "Tasks generated with interleaved dependencies using default ordering (api > web > cli > tui > mobile). Pipeline does not fail"
- State: "Task index contains test tasks ordered by default surface type priority. No error in task generation"
- Side-effect: "none"

## Outcome "not-found-feature-missing"
- Preconditions: "Feature slug does not exist or has no task files"
- Input: "forge task index --feature <nonexistent-slug>"
- Output: "Error message indicating feature not found or no tasks to index"
- State: "No tasks generated. Task index unchanged"
- Side-effect: "none"

## Outcome "already-exists-tasks"
- Preconditions: "Pre-existing test tasks for the feature already exist in the task index"
- Input: "forge task index --feature <slug> re-run with existing test tasks"
- Output: "Existing test tasks preserved (idempotent). Runtime fields (status, blocked-reason) not overwritten"
- State: "Task index retains existing task statuses and runtime state. New/updated tasks correctly merged"
- Side-effect: "none"

## Journey Invariants

- Every generated test task DAG maintains the interleaving invariant: for surface N>0, T-test-gen-scripts-{surface-N} depends on T-test-run-{surface-N-1}, NOT on T-test-gen-scripts-{surface-N-1}
- Every run-tests task includes AC enforcing real tests (no fake/stub-only tests that always pass)
- Every run-tests task enforces the "confirm before modifying production code" rule
- The pipeline does not regress for any number of configured surfaces
- Error handling follows the task-executor's existing Pause Protocol
