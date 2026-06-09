---
id: "T-test-gen-journeys"
title: "Generate Test Journeys"
priority: "P1"
estimated_time: "20-30min"
dependencies: ["1", "2", "3"]
type: "test.gen-journeys"
surface-key: ""
surface-type: ""
---

Generate test Journey documents for the worktree-start-idempotent feature.
Mode: quick


## Discovery Strategy

Discover the feature's testing directory layout before starting:
```bash
ls docs/features/worktree-start-idempotent/testing/                                 # journeys
ls docs/features/worktree-start-idempotent/testing/<journey>/contracts/              # contracts
```

Invoke the `/gen-journeys` skill to extract Journey narratives from specification documents.

### Input Source by Mode

- **Breakdown mode**: Read PRD user stories from `docs/features/worktree-start-idempotent/prd/prd-user-stories.md` and functional specs from `docs/features/worktree-start-idempotent/prd/prd-spec.md`. These are the primary input sources.
- **Quick mode**: Read the proposal from `docs/proposals/worktree-start-idempotent/proposal.md`. Extract Key Scenarios as Journey candidates. If the proposal lacks `scope` or `success criteria` sections, abort the task with a diagnostic message — Journey generation requires these minimum inputs.

## Process

Follow the `/gen-journeys` skill process flow:

1. **Surface Detection**: Detect the project surface type and persist to `.forge/config.yaml`
2. **Read Sources**: Read PRD user stories (Breakdown) or proposal.md (Quick)
3. **Identify Workflows**: Map each user story or key scenario to a Journey candidate
4. **Classify Risk**: Assign High/Medium/Low risk to each Journey based on workflow characteristics
5. **Generate Files**: Output one `journey.md` per Journey to `docs/features/worktree-start-idempotent/testing/<journey-name>/journey.md`
6. **Validate Output**: Check each Journey for required fields (name, risk level, happy path steps, edge cases, invariants)

## AUTO_COMMIT Directive

When this task runs as an automated pipeline task (not invoked manually by the user), AUTO_COMMIT=true is in effect:

- **If AUTO_COMMIT=true**: Skip the user review-and-approval step. After validation passes, directly commit all generated Journey files:
  ```bash
  git add docs/features/worktree-start-idempotent/testing/
  git commit -m "docs: generate journeys for worktree-start-idempotent"
  ```
- **If AUTO_COMMIT is not set** (manual invocation): Present all Journey files to the user for review. Wait for explicit approval before committing.

## Acceptance Criteria

- [ ] At least 1 Journey file generated under `docs/features/worktree-start-idempotent/testing/`
- [ ] Each Journey has: name, risk level, happy path steps, edge cases, invariants
- [ ] High-risk Journeys have edge case count >= happy path step count
- [ ] All Journey files committed (AUTO_COMMIT=true) or awaiting user review (manual mode)
