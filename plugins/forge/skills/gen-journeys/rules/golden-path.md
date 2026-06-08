---
name: golden-path
description: "Golden Path Journey mandatory rule — every feature must have at least one Golden Path Journey covering the primary user story's core domain action sequence."
---

# Golden Path Journey Rule

## Mandatory Requirement

Every feature MUST have **at least one Golden Path Journey**. This Journey must satisfy both constraints below simultaneously.

## Dual Constraints

A Journey qualifies as a Golden Path if and only if it satisfies **both** constraints:

### Constraint A: Multi-Step Span

The Journey must span **3+ step operations**. Each step must represent a distinct user action, not a sub-operation of a single action.

Valid examples:
- "Create item -> Edit item -> Verify updated state" (3 steps)
- "Create map -> Add milestone -> Move milestone -> Verify layout" (4 steps)

Invalid examples:
- "Submit form" -> "See confirmation" (only 2 steps)
- "Create milestone" -> "Create another milestone" -> "Create third milestone" (not a meaningful action sequence — no workflow progression)

### Constraint B: Semantic Completeness

The step sequence MUST be extracted from the PRD/Design document's **primary user story** or core workflow description. The steps must cover the core domain action sequence — not arbitrary operations assembled to meet the step count.

**Extraction rule**: Golden Path steps must describe what the user achieves (domain actions), not what the system does mechanically (API calls).

**Semantic completeness proxy**: Every Golden Path step description must reference domain terminology from the PRD/Design document (e.g., "create a milestone on the map", "transition task to in-progress"), NOT API/technical terminology (e.g., "POST /milestones", "send PATCH request").

**Anti-pattern — padding with unrelated operations**: The following are explicitly prohibited:
- Adding CRUD operations for entities not involved in the primary user story
- Repeating the same operation type (e.g., multiple "create" steps for different entities) without workflow justification
- Inserting verification-only steps (e.g., "list all items to confirm") that do not advance the user's goal

<HARD-RULE>
Golden Path applies to ALL surface types. There is no surface-specific differentiation of Golden Path requirements.
</HARD-RULE>

## Feature Complexity Classification Heuristics

Before generating Golden Path Journeys, classify the feature's complexity using the following heuristics:

| Criterion | Simple Feature | Complex Feature |
|-----------|---------------|-----------------|
| Entity type count | 1 entity type | >=2 entity types with parent-child or association relationships |
| Workflow description in PRD/Design | Single CRUD operation or linear flow | Multi-step workflow with state transitions or cross-entity interactions |
| Golden Path expectation | Complete CRUD cycle (create -> read -> update -> delete), 3-5 steps | Coverage of primary user story's core domain action sequence, 5+ steps |
| Fixture Specification expectation | Single entity + minimum count | Entity relationships + child entity minimum count + state constraints |

**Classification priority**: Entity relationships > workflow description. If parent-child entity relationships exist, the feature MUST be classified as Complex regardless of step count.

### Simple Feature Golden Path

For simple features (single entity, no parent-child relationships):
- A complete CRUD cycle (create -> read -> update -> delete) is a valid Golden Path
- 3-5 steps is the expected range
- No requirement for 5+ steps — the CRUD cycle naturally spans multiple steps

### Complex Feature Golden Path

For complex features (>=2 entity types with relationships):
- The Golden Path MUST cover entity interactions, not just isolated operations on each entity
- 5+ steps is expected, covering the primary user story's action sequence
- The sequence must demonstrate how entities relate and interact (e.g., parent creation -> child creation -> state transition affecting both)

## Enforcement in Step 2 (Identify User Workflows)

When identifying user workflows in Step 2 of the gen-journeys process:

1. After listing all Journey candidates from PRD/Proposal extraction, classify the feature complexity using the heuristics above
2. Identify which Journey candidate best represents the primary user story's core workflow — this becomes the Golden Path Journey
3. If no existing candidate covers the primary user story as a multi-step workflow, **create one** by merging related user stories that form a cohesive workflow
4. Validate the Golden Path candidate against both constraints (multi-step span + semantic completeness)
5. If the Golden Path candidate fails validation, expand it by incorporating additional steps from the PRD/Design's workflow description until both constraints are satisfied

## Validation in Step 5 (Validate Output)

In addition to the existing validation checks, verify:

| Check | Rule |
|-------|------|
| Golden Path exists | At least one Journey has `golden_path: true` in frontmatter |
| Golden Path step count | Golden Path Journey has >=3 happy path steps |
| Golden Path semantic completeness | Step descriptions reference domain terminology from PRD/Design, not API/technical terminology |
| Feature complexity classified | The complexity classification (simple/complex) is consistent with the Golden Path's step count and scope |
| Complex feature depth | For complex features: Golden Path has >=5 steps and covers cross-entity interactions |

If any Golden Path validation check fails, fix the Journey file before proceeding.
