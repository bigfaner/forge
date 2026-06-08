---
journey: "multi-surface-interleaved"
step: 5
step-action: "Execute the second surface's run-tests task"
generated: "2026-06-08"
sources:
  - docs/features/test-pipeline-interleaved/testing/multi-surface-interleaved/journey.md
skip_eval: true
---

# Contract: multi-surface-interleaved / Step 5: Execute the second surface's run-tests task

> **Note**: Contracts generated without eval-journey verification (SKIP_EVAL_GATE=true). Review with extra scrutiny.

<!-- gen-contracts: do not edit manually. Regenerate via /gen-contracts. -->

## Outcome "success"
- Preconditions: "T-test-run-{surface-2} is pending and its dependency T-test-gen-scripts-{surface-2} is completed"
- Input: "Task executor runs T-test-run-{surface-2}"
- Output: "Web tests execute with the corrected API behavior baseline. All test-run AC constraints apply (real tests, no fake tests, confirm before modifying production code). Tests pass"
- State: "T-test-run-{surface-2} status changed to completed. Pipeline fully complete for all surfaces"
- Side-effect: "Test execution results recorded"

## Outcome "agent-attempts-unauthorized-production-modification"
- Preconditions: "A test fails and the agent suspects a production code issue"
- Input: "Task executor runs T-test-run-{surface-2} where agent is tempted to modify production code"
- Output: "Agent MUST NOT modify production code without first confirming it is a genuine production code issue. Agent follows the confirm-before-modify protocol. If confirmed as production issue, creates fix task via forge task add"
- State: "Production code unchanged unless explicit confirmation given. Fix task created if confirmed"
- Side-effect: "Fix task created if production bug confirmed. Test scripts may be fixed directly if test script bug"

## Outcome "not-found-missing-scripts"
- Preconditions: "Test script files expected by run-tests task for the second surface are missing"
- Input: "Task executor runs T-test-run-{surface-2} with missing test scripts"
- Output: "Error message indicating test scripts not found. Task blocked"
- State: "Task status set to blocked"
- Side-effect: "none"

## Journey Invariants

- Every generated test task DAG maintains the interleaving invariant: for surface N>0, T-test-gen-scripts-{surface-N} depends on T-test-run-{surface-N-1}, NOT on T-test-gen-scripts-{surface-N-1}
- Every run-tests task includes AC enforcing real tests (no fake/stub-only tests that always pass)
- Every run-tests task enforces the "confirm before modifying production code" rule -- test script bugs may be fixed directly, but production code modification requires explicit confirmation and fix-task creation
- The pipeline does not regress for any number of configured surfaces
- Error handling follows the task-executor's existing Pause Protocol
