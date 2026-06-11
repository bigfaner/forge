---
feature: "test-pipeline-interleaved"
journey: "multi-surface-interleaved"
risk_level: "High"
surface_types: ["cli"]
sources:
  - docs/proposals/test-pipeline-interleaved/proposal.md
generated: "2026-06-08"
---

# Journey: multi-surface-interleaved

**Risk Level**: High

<!-- Risk Classification Criteria:
  High   = Workflow involves state mutation, data loss risk, or irreversible operations
  Multi-surface pipeline DAG generation mutates task dependency state; wrong wiring can cause tests to be silently skipped
-->

## Overview

A user with a multi-surface project (e.g., backend=api + frontend=web) runs the Forge test pipeline and verifies that test tasks are generated with interleaved dependencies: per-surface gen->run pairs execute sequentially rather than all-gen-then-all-run, enabling earlier bug detection.

## Setup

- A Forge project with at least two configured surfaces (e.g., api + web), with `execution_order` defined in `.forge/config.yaml`
- The project has a finalized feature with PRD or proposal documents ready for test pipeline generation
- No pre-existing test tasks for the feature in the task index

## Happy Path

### Step 1: Run pipeline task generation for a multi-surface feature

**User Action**: Execute `forge task index --feature test-pipeline-interleaved` (or equivalent pipeline generation command) to generate test tasks for a feature with multiple surfaces configured.

**Expected Result**: Test tasks are generated for each surface in execution_order: `T-test-gen-scripts-api`, `T-test-run-api`, `T-test-gen-scripts-web`, `T-test-run-web`. The dependency chain follows the interleaved pattern: `gen-scripts-api -> run-api -> gen-scripts-web -> run-web`.

### Step 2: Execute the first surface's gen-scripts task

**User Action**: Run the first surface's gen-scripts task (e.g., `T-test-gen-scripts-api`).

**Expected Result**: The task generates test scripts for the API surface. The task completes successfully with test scripts written to the expected directory.

### Step 3: Execute the first surface's run-tests task

**User Action**: Run the first surface's run-tests task (e.g., `T-test-run-api`), which depends only on `T-test-gen-scripts-api`.

**Expected Result**: API tests execute. If an API bug is discovered, it is reported. The run-tests task includes AC enforcing: all tests must pass, tests must be real (not fake), and production code modification requires explicit confirmation. The task does NOT wait for gen-scripts of other surfaces.

### Step 4: Execute the second surface's gen-scripts task

**User Action**: Run the second surface's gen-scripts task (e.g., `T-test-gen-scripts-web`), which depends on `T-test-run-api` (the prior surface's run-tests, NOT the prior surface's gen-scripts).

**Expected Result**: Web test scripts are generated, potentially informed by feedback from the API test run (e.g., if API behavior was corrected). This is the key interleaving benefit: gen-scripts-web can incorporate learnings from run-api.

### Step 5: Execute the second surface's run-tests task

**User Action**: Run the second surface's run-tests task (e.g., `T-test-run-web`), which depends on `T-test-gen-scripts-web`.

**Expected Result**: Web tests execute with the corrected API behavior baseline. All test-run AC constraints apply (real tests, no fake tests, confirm before modifying production code).

### Step 6: Verify the complete pipeline completes

**User Action**: Inspect the task status for all generated test tasks.

**Expected Result**: All tasks completed. The total pipeline execution discovered issues earlier than the old serial approach because API tests ran and surfaced bugs before web scripts were generated.

## Edge Cases

### Step 1b: Three-surface project with execution_order

**Precondition**: Project has three configured surfaces (e.g., backend=api, frontend=web, cli=cli) with execution_order: [api, web, cli].

**User Action**: Execute pipeline task generation for the three-surface feature.

**Expected Result**: Tasks are generated with chain: `gen-scripts-api -> run-api -> gen-scripts-web -> run-web -> gen-scripts-cli -> run-cli`. Each N-th gen-scripts depends on (N-1)-th run-tests. The interleaving extends naturally to any number of surfaces.

### Step 1c: Pipeline generation with surfaces but no execution_order

**Precondition**: Project has multiple surfaces configured but `execution_order` is not explicitly defined in config.

**User Action**: Execute pipeline task generation.

**Expected Result**: The system uses a default ordering (e.g., alphabetical or surface-type priority). Tasks are still generated with interleaved dependencies. The pipeline does not fail.

### Step 2b: First surface gen-scripts fails

**Precondition**: The API surface's gen-scripts task encounters an error (e.g., no contracts found, template rendering fails).

**User Action**: The task executor runs `T-test-gen-scripts-api` which fails.

**Expected Result**: The failure blocks downstream tasks (`T-test-run-api`, `T-test-gen-scripts-web`, `T-test-run-web`). The Error Handling pause protocol is triggered. The user is notified and can create a fix task.

### Step 3b: First surface run-tests discovers API bug

**Precondition**: API tests execute but one or more tests fail due to a real API bug.

**User Action**: The task executor runs `T-test-run-api` which reports test failures.

**Expected Result**: The task follows the hardened AC: confirms the failure is due to production code (not test script bug), creates a fix task via `forge task add` rather than silently modifying tests to pass. The second surface's gen-scripts is blocked until the fix resolves.

### Step 3c: Run-tests encounters test script bug (not production bug)

**Precondition**: The generated test scripts themselves contain a bug (e.g., wrong assertion, incorrect setup).

**User Action**: The task executor runs `T-test-run-api` which reports failures.

**Expected Result**: Per the prompt template instruction, the agent identifies this as a test script bug (not a production code bug). It fixes the test script directly without modifying production code. The test script fix does not count as "faking" the test -- it corrects a generation error.

### Step 4b: Second surface gen-scripts runs with no API feedback needed

**Precondition**: API tests all pass cleanly with no bugs found.

**User Action**: Run `T-test-gen-scripts-web` after successful `T-test-run-api`.

**Expected Result**: Web test scripts are generated successfully. The dependency on run-api is satisfied but no corrective feedback was needed. The interleaving still provides the guarantee that API behavior was verified before web scripts depend on it.

### Step 5b: Run-tests agent attempts to modify production code without confirmation

**Precondition**: A test fails and the agent suspects a production code issue.

**User Action**: The task executor runs a run-tests task where the agent is tempted to modify production code to make the test pass.

**Expected Result**: Per the prompt template instruction, the agent MUST NOT modify production code without first confirming it is a genuine production code issue (not a test bug). The agent follows the "confirm before modify" protocol. If confirmed as a production issue, the agent creates a fix task via `forge task add`.

## Journey Invariants

- Every generated test task DAG maintains the interleaving invariant: for surface N>0, `T-test-gen-scripts-{surface-N}` depends on `T-test-run-{surface-N-1}`, NOT on `T-test-gen-scripts-{surface-N-1}`
- Every run-tests task includes AC enforcing real tests (no fake/stub-only tests that always pass)
- Every run-tests task enforces the "confirm before modifying production code" rule -- test script bugs may be fixed directly, but production code modification requires explicit confirmation and fix-task creation
- The pipeline does not regress for any number of configured surfaces -- single surface degradation is covered by a separate Journey
- Error handling follows the task-executor's existing Pause Protocol -- new AC instructions supplement, not replace, EXTREMELY-IMPORTANT level directives
