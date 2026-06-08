---
journey: "single-surface-degradation"
step: 2
step-action: "Execute gen-scripts and run-tests in sequence"
generated: "2026-06-08"
sources:
  - docs/features/test-pipeline-interleaved/testing/single-surface-degradation/journey.md
skip_eval: true
---

# Contract: single-surface-degradation / Step 2: Execute gen-scripts and run-tests in sequence

> **Note**: Contracts generated without eval-journey verification (SKIP_EVAL_GATE=true). Review with extra scrutiny.

<!-- gen-contracts: do not edit manually. Regenerate via /gen-contracts. -->

## Outcome "success"
- Preconditions: "T-test-gen-scripts-cli is pending. T-test-run-cli depends on gen-scripts-cli"
- Input: "Run tasks in dependency order: first T-test-gen-scripts-cli, then T-test-run-cli"
- Output: "gen-scripts completes and produces test scripts. run-tests executes those scripts. Hardened AC from test-run template applies (real tests, no fakes, confirm before modifying production code). Behavior is identical to the pre-interleaving implementation"
- State: "Both tasks completed. gen-scripts-cli -> run-cli simple chain fully executed"
- Side-effect: "Test scripts written then executed. Results recorded"

## Outcome "not-found-scripts-missing"
- Preconditions: "gen-scripts task completed but test scripts are missing from expected location"
- Input: "run-tests executor attempts to run T-test-run-cli but scripts not found"
- Output: "Error indicating test scripts not found. Task blocked"
- State: "run-tests task blocked"
- Side-effect: "none"

## Journey Invariants

- Single-surface projects always produce exactly two test tasks: one gen-scripts and one run-tests
- The dependency chain for single-surface projects is always gen-scripts -> run-tests with no intermediate dependencies
- The hardened test-run AC (real tests, no fakes, confirm before production code changes) applies regardless of surface count
- The interleaving dependency logic is never activated for single-surface projects -- no surface N>0 exists to trigger the cross-surface dependency
