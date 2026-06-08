---
journey: "multi-surface-interleaved"
step: 2
step-action: "Execute the first surface's gen-scripts task"
generated: "2026-06-08"
sources:
  - docs/features/test-pipeline-interleaved/testing/multi-surface-interleaved/journey.md
skip_eval: true
---

# Contract: multi-surface-interleaved / Step 2: Execute the first surface's gen-scripts task

> **Note**: Contracts generated without eval-journey verification (SKIP_EVAL_GATE=true). Review with extra scrutiny.

<!-- gen-contracts: do not edit manually. Regenerate via /gen-contracts. -->

## Outcome "success"
- Preconditions: "First surface's gen-scripts task (T-test-gen-scripts-{surface-1}) exists and is pending; contracts for the surface are available"
- Input: "Task executor runs T-test-gen-scripts-{surface-1}"
- Output: "Task generates test scripts for the first surface. Task completes successfully with test scripts written to the expected directory"
- State: "T-test-gen-scripts-{surface-1} status changed to completed. Test scripts exist in the appropriate test directory"
- Side-effect: "Test script files written to disk in tests/<journey>/ or equivalent location"

## Outcome "gen-scripts-failure"
- Preconditions: "First surface's gen-scripts task encounters an error (e.g., no contracts found, template rendering fails)"
- Input: "Task executor runs T-test-gen-scripts-{surface-1} which fails"
- Output: "Error reported. Downstream tasks blocked. Error Handling pause protocol triggered"
- State: "T-test-gen-scripts-{surface-1} status set to blocked. T-test-run-{surface-1}, T-test-gen-scripts-{surface-2}, T-test-run-{surface-2} remain blocked as dependents"
- Side-effect: "Fix task may be created via forge task add --block-source"

## Outcome "not-found-no-contracts"
- Preconditions: "No contract files exist for the first surface in the testing directory"
- Input: "Task executor runs T-test-gen-scripts-{surface-1} with missing contracts"
- Output: "Error message indicating no contracts found for the surface"
- State: "Task status set to blocked. No test scripts generated"
- Side-effect: "none"

## Journey Invariants

- Every generated test task DAG maintains the interleaving invariant: for surface N>0, T-test-gen-scripts-{surface-N} depends on T-test-run-{surface-N-1}, NOT on T-test-gen-scripts-{surface-N-1}
- Every run-tests task includes AC enforcing real tests (no fake/stub-only tests that always pass)
- Every run-tests task enforces the "confirm before modifying production code" rule
- The pipeline does not regress for any number of configured surfaces
- Error handling follows the task-executor's existing Pause Protocol
