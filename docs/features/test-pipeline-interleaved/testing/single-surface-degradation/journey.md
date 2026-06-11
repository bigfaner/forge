---
feature: "test-pipeline-interleaved"
journey: "single-surface-degradation"
risk_level: "Low"
surface_types: ["cli"]
sources:
  - docs/proposals/test-pipeline-interleaved/proposal.md
generated: "2026-06-08"
---

# Journey: single-surface-degradation

**Risk Level**: Low

<!-- Risk Classification Criteria:
  Low = Workflow is read-only or purely observational
  This Journey verifies that single-surface projects retain their existing gen->run behavior unchanged.
-->

## Overview

A user with a single-surface project (e.g., only CLI) runs the Forge test pipeline and verifies that the dependency chain remains the simple gen -> run pattern, with no interleaving overhead or behavioral change.

## Setup

- A Forge project with exactly one configured surface (e.g., cli)
- The project has a finalized feature with documents ready for test pipeline generation
- No pre-existing test tasks for the feature

## Happy Path

### Step 1: Run pipeline task generation for a single-surface feature

**User Action**: Execute `forge task index --feature test-pipeline-interleaved` (or equivalent) for a project with only one surface configured.

**Expected Result**: Exactly two test tasks are generated: `T-test-gen-scripts-cli` and `T-test-run-cli`. The dependency chain is `gen-scripts-cli -> run-cli`. No interleaving occurs because there is only one surface.

### Step 2: Execute gen-scripts and run-tests in sequence

**User Action**: Run the generated tasks in dependency order: first `T-test-gen-scripts-cli`, then `T-test-run-cli`.

**Expected Result**: gen-scripts completes and produces test scripts. run-tests executes those scripts. The hardened AC from the test-run template still applies (real tests, no fakes, confirm before modifying production code). The behavior is identical to the pre-interleaving implementation.

## Edge Cases

### Step 1b: Single surface with no execution_order configured

**Precondition**: The project has one surface but `execution_order` is not explicitly set.

**User Action**: Execute pipeline task generation.

**Expected Result**: The system detects a single surface and generates the standard gen -> run chain. The absence of execution_order does not cause an error because single-surface mode is the natural degenerate case.

## Journey Invariants

- Single-surface projects always produce exactly two test tasks: one gen-scripts and one run-tests
- The dependency chain for single-surface projects is always `gen-scripts -> run-tests` with no intermediate dependencies
- The hardened test-run AC (real tests, no fakes, confirm before production code changes) applies regardless of surface count
- The interleaving dependency logic is never activated for single-surface projects -- no surface N>0 exists to trigger the cross-surface dependency
