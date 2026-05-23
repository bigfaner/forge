Generate test Journey documents for the {{FEATURE_SLUG}} feature.
Mode: {{MODE}}

## Scope

{{SCOPE}}

## Discovery Strategy

Invoke the `/gen-journeys` skill to extract Journey narratives from specification documents.

### Input Source by Mode

- **Breakdown mode**: Read PRD user stories from `docs/features/{{FEATURE_SLUG}}/prd/prd-user-stories.md` and functional specs from `docs/features/{{FEATURE_SLUG}}/prd/prd-spec.md`. These are the primary input sources.
- **Quick mode**: Read the proposal from `docs/proposals/{{FEATURE_SLUG}}/proposal.md`. Extract Key Scenarios as Journey candidates. If the proposal lacks `scope` or `success criteria` sections, abort the task with a diagnostic message — Journey generation requires these minimum inputs.

## Process

Follow the `/gen-journeys` skill process flow:

1. **Surface Detection**: Detect the project surface type and persist to `.forge/config.yaml`
2. **Read Sources**: Read PRD user stories (Breakdown) or proposal.md (Quick)
3. **Identify Workflows**: Map each user story or key scenario to a Journey candidate
4. **Classify Risk**: Assign High/Medium/Low risk to each Journey based on workflow characteristics
5. **Generate Files**: Output one `journey.md` per Journey to `docs/features/{{FEATURE_SLUG}}/testing/<journey-name>/journey.md`
6. **Validate Output**: Check each Journey for required fields (name, risk level, happy path steps, edge cases, invariants)

## AUTO_COMMIT Directive

When this task runs as an automated pipeline task (not invoked manually by the user), AUTO_COMMIT=true is in effect:

- **If AUTO_COMMIT=true**: Skip the user review-and-approval step. After validation passes, directly commit all generated Journey files:
  ```bash
  git add docs/features/{{FEATURE_SLUG}}/testing/
  git commit -m "docs: generate journeys for {{FEATURE_SLUG}}"
  ```
- **If AUTO_COMMIT is not set** (manual invocation): Present all Journey files to the user for review. Wait for explicit approval before committing.

## Acceptance Criteria

- [ ] At least 1 Journey file generated under `docs/features/{{FEATURE_SLUG}}/testing/`
- [ ] Each Journey has: name, risk level, happy path steps, edge cases, invariants
- [ ] High-risk Journeys have edge case count >= happy path step count
- [ ] All Journey files committed (AUTO_COMMIT=true) or awaiting user review (manual mode)
