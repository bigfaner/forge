---
journey: "single-surface-degradation"
step: 1
step-action: "Run pipeline task generation for a single-surface feature"
generated: "2026-06-08"
sources:
  - docs/features/test-pipeline-interleaved/testing/single-surface-degradation/journey.md
skip_eval: true
---

# Contract: single-surface-degradation / Step 1: Run pipeline task generation for a single-surface feature

> **Note**: Contracts generated without eval-journey verification (SKIP_EVAL_GATE=true). Review with extra scrutiny.

<!-- gen-contracts: do not edit manually. Regenerate via /gen-contracts. -->

## Outcome "success"
- Preconditions: "Project has exactly one configured surface (e.g., cli). Feature has finalized documents. No pre-existing test tasks for the feature"
- Input: "forge task index --feature <slug> for a single-surface project"
- Output: "Exactly two test tasks generated: T-test-gen-scripts-cli and T-test-run-cli. Dependency chain is gen-scripts-cli -> run-cli. No interleaving occurs because there is only one surface"
- State: "Task index contains exactly two test tasks with correct gen->run dependency. Each task has surface-key and surface-type fields populated"
- Side-effect: "Task .md files written for both generated tasks"

## Outcome "no-execution-order-single-surface"
- Preconditions: "Project has one surface but execution_order is not explicitly set"
- Input: "forge task index --feature <slug>"
- Output: "System detects a single surface and generates the standard gen -> run chain. Absence of execution_order does not cause an error because single-surface mode is the natural degenerate case"
- State: "Task index contains two tasks with correct dependency. No error"
- Side-effect: "none"

## Outcome "not-found-feature-missing"
- Preconditions: "Feature slug does not exist or has no task files"
- Input: "forge task index --feature <nonexistent-slug>"
- Output: "Error message indicating feature not found or no tasks to index"
- State: "No tasks generated. Task index unchanged"
- Side-effect: "none"

## Journey Invariants

- Single-surface projects always produce exactly two test tasks: one gen-scripts and one run-tests
- The dependency chain for single-surface projects is always gen-scripts -> run-tests with no intermediate dependencies
- The hardened test-run AC (real tests, no fakes, confirm before production code changes) applies regardless of surface count
- The interleaving dependency logic is never activated for single-surface projects -- no surface N>0 exists to trigger the cross-surface dependency
