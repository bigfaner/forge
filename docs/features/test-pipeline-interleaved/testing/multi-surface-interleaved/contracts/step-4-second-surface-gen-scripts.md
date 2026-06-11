---
journey: "multi-surface-interleaved"
step: 4
step-action: "Execute the second surface's gen-scripts task"
generated: "2026-06-08"
sources:
  - docs/features/test-pipeline-interleaved/testing/multi-surface-interleaved/journey.md
skip_eval: true
---

# Contract: multi-surface-interleaved / Step 4: Execute the second surface's gen-scripts task

> **Note**: Contracts generated without eval-journey verification (SKIP_EVAL_GATE=true). Review with extra scrutiny.

<!-- gen-contracts: do not edit manually. Regenerate via /gen-contracts. -->

## Outcome "success"
- Preconditions: "T-test-gen-scripts-{surface-2} is pending; T-test-run-{surface-1} is completed (interleaved dependency satisfied)"
- Input: "Task executor runs T-test-gen-scripts-{surface-2}"
- Output: "Web test scripts generated, potentially informed by feedback from the API test run. Task completes successfully"
- State: "T-test-gen-scripts-{surface-2} status changed to completed. Web test scripts written to expected directory"
- Side-effect: "Test script files written to disk. This is the key interleaving benefit: gen-scripts can incorporate learnings from prior surface's test run"

## Outcome "success-no-feedback-needed"
- Preconditions: "Prior surface (API) tests all passed cleanly with no bugs found; T-test-run-{surface-1} completed"
- Input: "Task executor runs T-test-gen-scripts-{surface-2} after successful T-test-run-{surface-1}"
- Output: "Web test scripts generated successfully. Dependency on run-{surface-1} is satisfied but no corrective feedback was needed. Interleaving still provides the guarantee that API behavior was verified before web scripts depend on it"
- State: "T-test-gen-scripts-{surface-2} completed. Test scripts generated"
- Side-effect: "Test script files written"

## Outcome "not-found-no-contracts"
- Preconditions: "No contract files exist for the second surface in the testing directory"
- Input: "Task executor runs T-test-gen-scripts-{surface-2} with missing contracts"
- Output: "Error message indicating no contracts found for the surface"
- State: "Task status set to blocked. No test scripts generated"
- Side-effect: "none"

## Journey Invariants

- Every generated test task DAG maintains the interleaving invariant: for surface N>0, T-test-gen-scripts-{surface-N} depends on T-test-run-{surface-N-1}, NOT on T-test-gen-scripts-{surface-N-1}
- Every run-tests task includes AC enforcing real tests (no fake/stub-only tests that always pass)
- Every run-tests task enforces the "confirm before modifying production code" rule
- The pipeline does not regress for any number of configured surfaces
- Error handling follows the task-executor's existing Pause Protocol
