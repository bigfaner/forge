---
journey: "start-existing-flags"
step: 2
step-action: "Start with --no-launch on existing worktree"
generated: "2026-06-09"
skip_eval: true
sources:
  - docs/features/worktree-start-idempotent/testing/start-existing-flags/journey.md

anchors:
  cli:
    command: "forge worktree start"
    subcommand: ""
    flags: ["no-launch"]
    aliases: []
last_anchor_sync: ""
---

# Contract: start-existing-flags / Step 2: Start with --no-launch on existing worktree

> **Note**: Contracts generated without eval-journey verification (SKIP_EVAL_GATE=true). Review with extra scrutiny.

<!-- gen-contracts: do not edit manually. Regenerate via /gen-contracts. -->

## Outcome "success"
- Preconditions: "A valid worktree already exists at .forge/worktrees/<slug>"
  fixture_spec:
    entities:
      - entity_type: "Worktree"
        min_count: 1
        field_constraints:
          - field: "state"
            value: "valid (directory exists, .git file present)"
- Input: "forge worktree start <slug> --no-launch"
- Output: "The existing worktree is validated. stdout outputs the worktree path: 'worktree path: <resolvedDir>'. stderr contains 'entering existing worktree'. No Claude session is launched. Exit code is 0."
- State: "No file system mutations occur. The worktree state is unchanged."
- Side-effect: "none"

## Outcome "no-launch-new-worktree"
<!-- source: inferred -->
<!-- reasoning: Journey edge case Step 1b: --no-launch on a non-existent worktree should create a new worktree. Code path at cmd_start.go:239-243 shows --no-launch prints path and exits after creation. -->
- Preconditions: "No worktree exists for the slug"
  fixture_spec:
    entities:
      - entity_type: "Worktree"
        min_count: 0
- Input: "forge worktree start <slug> --no-launch"
- Output: "A new worktree is created. Includes files are copied if configured. stdout contains 'worktree created at <targetDir>'. stderr contains 'created new worktree'. No Claude session is launched. Exit code is 0."
- State: "New git worktree added. New branch created. Includes files copied from project root."
- Side-effect: "none (Claude session explicitly suppressed)"

## Outcome "combined-flags-ignored"
<!-- source: inferred -->
<!-- reasoning: Journey edge case Step 2b: --source-branch + --no-launch combined. Both flags interact: --source-branch is warned and ignored on existing worktree, --no-launch suppresses claude launch. -->
- Preconditions: "A valid worktree already exists for the slug"
  fixture_spec:
    entities:
      - entity_type: "Worktree"
        min_count: 1
        field_constraints:
          - field: "state"
            value: "valid"
- Input: "forge worktree start <slug> --source-branch develop --no-launch"
- Output: "stderr contains warning about ignoring --source-branch. Worktree path is output to stdout. No Claude session is launched. Exit code is 0."
- State: "No branch change. No file mutations. Worktree path printed to stdout."
- Side-effect: "none"

## Journey Invariants

- --no-launch always suppresses Claude session launch regardless of whether worktree is new or existing
- Exit code is 0 for all successful flag combinations
- No file system mutations occur when entering an existing worktree

## Fixture Specification

This Contract requires the following pre-existing data state. See `rules/fixture-spec.md` for schema details.

```yaml
fixture_spec:
  entities:
    - entity_type: "Worktree"           # Existing valid worktree
      min_count: 1
      field_constraints:
        - field: "state"
          value: "valid"
```
