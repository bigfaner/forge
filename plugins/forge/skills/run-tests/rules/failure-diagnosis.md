---
name: failure-diagnosis
description: Diagnostic flow for surface test failures, including app health gate and general failure analysis
---

# Failure Diagnosis

## App Health First Gate

When tests fail, the first step is determining whether the **app itself is healthy**. Surface test error signals cannot distinguish "test wrote wrong selector" from "app crashed and renders nothing" -- both produce "element not found".

<HARD-RULE>
When **>30% of tests fail simultaneously**, do NOT proceed to individual test fix tasks. Run app health diagnostics first. Batch failures almost always indicate an app-level problem, not per-test issues.
</HARD-RULE>

| Failure ratio | Likely cause | First action |
|---|---|---|
| **>30%** tests fail simultaneously | App health problem | Check failure evidence (screenshots, error logs) for systemic issues |
| 10-30% partial failure | Possibly test issues | Spot check 2-3 failures before deciding |
| <10% few failures | Per-test issues | Proceed to per-test fix tasks |

**App health diagnostic flow** (run in order, stop at first positive finding):

1. **Check failure evidence** -- screenshots, error messages, console output for systemic patterns (e.g., all returning same error code, blank screens)
2. **Check app infrastructure** -- verify the application server is running, dependencies are installed, configuration is correct
3. **Verify app responds** -- manually test the app's health endpoint or main interface
4. **Only after app is confirmed healthy** -- proceed to individual test failure analysis

## General Failure Analysis

When tests fail, do not stop at the first visible error message. Follow these rules to avoid misdiagnosis:

<PRINCIPLE>
**Surface-level errors are often secondary effects.** The first error message in test output is frequently a symptom, not the root cause. Look for the underlying issue rather than treating each failure independently.
</PRINCIPLE>

**Diagnostic checklist when batch failures occur:**

1. **Check for cascade patterns** -- if many tests fail with the same symptom (e.g., all 404 on `undefined` in URL paths), the root cause is in shared setup, not individual tests. One setup failure leaves all its module-level variables uninitialized.

2. **Ask the contradiction question** -- "Why does this test fail when other tests with the same pattern pass?" If many tests use the same setup and pass, the failure is specific to this test's setup, not a platform issue.

3. **Verify backend health** -- after code changes, the backend must be rebuilt AND restarted before re-running surface tests. Check:
   - Did the health check (`just probe`) actually pass?
   - Are there backend logs showing the failed requests?
   - Did a rate limit or connection reset occur during setup?

4. **Investigate setup blocks step-by-step** -- when setup is suspected:
   - Add temporary logging between steps to identify which operation fails
   - Run the setup code manually to isolate the issue
   - Use the test runner's debug mode to step through
