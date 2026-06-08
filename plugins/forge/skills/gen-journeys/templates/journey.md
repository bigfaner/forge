---
feature: "{{FEATURE_SLUG}}"
journey: "{{JOURNEY_NAME}}"
risk_level: "{{RISK_LEVEL}}"
golden_path: false
surface_types: ["{{SURFACE_TYPE_1}}", "{{SURFACE_TYPE_2}}"]
surface_keys: ["{{SURFACE_KEY_1}}", "{{SURFACE_KEY_2}}"]
sources:
  - docs/features/{{FEATURE_SLUG}}/prd/prd-user-stories.md
  - docs/features/{{FEATURE_SLUG}}/prd/prd-spec.md
generated: "{{DATE}}"
---

# Journey: {{JOURNEY_NAME}}

**Risk Level**: {{RISK_LEVEL}}

<!-- Risk Classification Criteria:
  High   = Workflow involves state mutation, data loss risk, or irreversible operations
  Medium = Workflow involves multi-step interaction without irreversible side effects
  Low    = Workflow is read-only or purely observational
-->

## Overview

{{JOURNEY_OVERVIEW}}

<!-- One-sentence description of the user workflow and its goal -->

## Setup

<!-- Preconditions that must be established before the Journey starts.
     These are environment states, not user actions. -->

- {{SETUP_PRECONDITION_1}}
- {{SETUP_PRECONDITION_2}}

## Happy Path

<!-- The primary success scenario: steps the user takes to accomplish the goal.
     Each step describes a user action and the expected outcome.
     High-risk Journeys MUST have edge case count >= happy path step count. -->

### Step 1: {{STEP_1_ACTION}}

**User Action**: {{STEP_1_USER_ACTION}}

**Expected Result**: {{STEP_1_EXPECTED_RESULT}}

### Step 2: {{STEP_2_ACTION}}

**User Action**: {{STEP_2_USER_ACTION}}

**Expected Result**: {{STEP_2_EXPECTED_RESULT}}

<!-- Repeat for additional happy path steps -->

## Edge Cases

<!-- Alternative scenarios where things go wrong or take an unexpected path.
     Each edge case references a happy path step (variant) and describes the
     divergent precondition and expected outcome.
     High-risk Journeys: number of edge cases MUST be >= number of happy path steps. -->

### Step 1b: {{EDGE_1_ACTION}}

**Precondition**: {{EDGE_1_PRECONDITION}}

<!-- The precondition that differs from the happy path, causing this outcome -->

**User Action**: {{EDGE_1_USER_ACTION}}

**Expected Result**: {{EDGE_1_EXPECTED_RESULT}}

### Step 2b: {{EDGE_2_ACTION}}

**Precondition**: {{EDGE_2_PRECONDITION}}

**User Action**: {{EDGE_2_USER_ACTION}}

**Expected Result**: {{EDGE_2_EXPECTED_RESULT}}

<!-- Repeat for additional edge cases -->

## Journey Invariants

<!-- Cross-step properties that must hold throughout the entire Journey.
     At least one invariant is required per Journey.
     These are verified across all steps, not within a single step. -->

- {{INVARIANT_1}}
- {{INVARIANT_2}}
