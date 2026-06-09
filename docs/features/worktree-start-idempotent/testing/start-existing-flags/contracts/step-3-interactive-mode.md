---
journey: "start-existing-flags"
step: 3
step-action: "Start with --interactive on existing worktree"
generated: "2026-06-09"
skip_eval: true
sources:
  - docs/features/worktree-start-idempotent/testing/start-existing-flags/journey.md

anchors:
  cli:
    command: "forge worktree start"
    subcommand: ""
    flags: ["interactive"]
    aliases: ["i"]
last_anchor_sync: ""
---

# Contract: start-existing-flags / Step 3: Start with --interactive on existing worktree

> **Note**: Contracts generated without eval-journey verification (SKIP_EVAL_GATE=true). Review with extra scrutiny.

<!-- gen-contracts: do not edit manually. Regenerate via /gen-contracts. -->

## Outcome "success"
- Preconditions: "A valid worktree already exists at .forge/worktrees/<slug>. At least one unfinished proposal or feature exists in the project. The terminal is a TTY."
  fixture_spec:
    entities:
      - entity_type: "Worktree"
        min_count: 1
        field_constraints:
          - field: "state"
            value: "valid (directory exists, .git file present)"
      - entity_type: "Proposal"
        min_count: 1
        field_constraints:
          - field: "status"
            value: "not completed"
- Input: "forge worktree start --interactive, then select the slug from the interactive list"
- Output: "The existing worktree is entered. Behavior is identical to running forge worktree start <slug> explicitly. stderr contains 'entering existing worktree'. A new Claude session is launched."
- State: "No new worktree is created. The existing worktree is reused. No file system mutations."
- Side-effect: "Claude session launched in existing worktree directory"

## Outcome "no-existing-worktrees"
<!-- source: inferred -->
<!-- reasoning: Journey edge case Step 3b: --interactive with no worktrees. The interactive mode scans for unfinished proposals/features (not existing worktrees), so the list shows available items to start work on. If no items exist, the output differs per Fact Table WT_INTERACTIVE_NO_ITEMS. -->
- Preconditions: "No worktrees exist in .forge/worktrees/. No unfinished proposals or features exist in the project."
  fixture_spec:
    entities:
      - entity_type: "Worktree"
        min_count: 0
      - entity_type: "Proposal"
        min_count: 0
      - entity_type: "Feature"
        min_count: 0
- Input: "forge worktree start --interactive"
- Output: "stdout contains 'No unfinished proposals or features found.' followed by 'Create one with: forge proposal <slug> or forge feature <slug>'. Exit code is 0. No worktree is created."
- State: "No changes to filesystem or git state."
- Side-effect: "none"

## Outcome "non-tty"
<!-- source: inferred -->
<!-- reasoning: Fact Table WT_INTERACTIVE_TTY_CHECK shows -i requires TTY. Non-TTY invocation produces an error. -->
- Preconditions: "The terminal is not a TTY (piped stdin)"
  fixture_spec:
    entities:
      - entity_type: "Worktree"
        min_count: 1
- Input: "forge worktree start --interactive in a non-TTY environment"
- Output: "stderr contains 'error: interactive mode requires a terminal (TTY)'. Exit code is non-zero."
- State: "No worktree operations are performed."
- Side-effect: "none"

## Journey Invariants

- --interactive mode behavior is consistent with explicit slug specification after selection
- --no-launch always suppresses Claude session launch regardless of interactive selection
- Exit code is 0 for all successful flag combinations

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
    - entity_type: "Proposal"           # Unfinished proposal for interactive selection
      min_count: 1
      field_constraints:
        - field: "status"
          value: "not completed"
```
