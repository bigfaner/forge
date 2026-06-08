---
journey: "multi-surface-interleaved"
step: 6
step-action: "Verify the complete pipeline completes"
generated: "2026-06-08"
sources:
  - docs/features/test-pipeline-interleaved/testing/multi-surface-interleaved/journey.md
skip_eval: true
---

# Contract: multi-surface-interleaved / Step 6: Verify the complete pipeline completes

> **Note**: Contracts generated without eval-journey verification (SKIP_EVAL_GATE=true). Review with extra scrutiny.

<!-- gen-contracts: do not edit manually. Regenerate via /gen-contracts. -->

## Outcome "success"
- Preconditions: "All four (or six) test tasks have been executed; none remain pending or blocked"
- Input: "User inspects task status for all generated test tasks (e.g., forge task status or reading index.json)"
- Output: "All tasks show completed status. Total pipeline execution discovered issues earlier than the old serial approach because API tests ran and surfaced bugs before web scripts were generated"
- State: "All test tasks in completed status. Pipeline fully verified"
- Side-effect: "none"

## Outcome "partial-failure-pipeline"
- Preconditions: "Some test tasks completed but at least one is blocked or failed"
- Input: "User inspects task status and finds incomplete tasks"
- Output: "Status report showing which tasks completed and which remain blocked with blocked reasons"
- State: "Pipeline not fully complete. Blocked tasks have associated fix tasks or error descriptions"
- Side-effect: "none"

## Outcome "not-found-feature"
- Preconditions: "Feature slug does not exist in the task index"
- Input: "User inspects task status for a nonexistent feature"
- Output: "Error indicating feature not found or no task index exists"
- State: "No change to task state"
- Side-effect: "none"

## Journey Invariants

- Every generated test task DAG maintains the interleaving invariant: for surface N>0, T-test-gen-scripts-{surface-N} depends on T-test-run-{surface-N-1}, NOT on T-test-gen-scripts-{surface-N-1}
- Every run-tests task includes AC enforcing real tests (no fake/stub-only tests that always pass)
- Every run-tests task enforces the "confirm before modifying production code" rule
- The pipeline does not regress for any number of configured surfaces
- Error handling follows the task-executor's existing Pause Protocol
