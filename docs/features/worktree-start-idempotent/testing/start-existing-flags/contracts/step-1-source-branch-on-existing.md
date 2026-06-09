---
journey: "start-existing-flags"
step: 1
step-action: "Start with --source-branch on existing worktree"
generated: "2026-06-09"
skip_eval: true
sources:
  - docs/features/worktree-start-idempotent/testing/start-existing-flags/journey.md

anchors:
  cli:
    command: "forge worktree start"
    subcommand: ""
    flags: ["source-branch"]
    aliases: []
last_anchor_sync: ""
---

# Contract: start-existing-flags / Step 1: Start with --source-branch on existing worktree

> **Note**: Contracts generated without eval-journey verification (SKIP_EVAL_GATE=true). Review with extra scrutiny.

<!-- gen-contracts: do not edit manually. Regenerate via /gen-contracts. -->

## Outcome "success-ignored-flag"
- Preconditions: "A valid worktree already exists at .forge/worktrees/<slug>. The worktree was created from the main branch. The worktree.includes config lists at least one file."
  fixture_spec:
    entities:
      - entity_type: "Worktree"
        min_count: 1
        field_constraints:
          - field: "state"
            value: "valid (directory exists, .git file present)"
          - field: "source_branch"
            value: "main"
      - entity_type: "Config"
        min_count: 1
        field_constraints:
          - field: "worktree.includes"
            value: "at least one file path"
- Input: "forge worktree start <slug> --source-branch develop"
- Output: "The existing worktree is entered. stderr contains a warning about ignoring --source-branch: 'warning: worktree already exists, ignoring --source-branch'. stderr also contains 'entering existing worktree'. A new Claude session is launched."
- State: "No branch change occurs. The worktree remains on its existing branch. No file system mutations occur (no re-copying of includes)."
- Side-effect: "Claude session launched in existing worktree directory"

## Outcome "diverged-branch"
<!-- source: inferred -->
<!-- reasoning: Journey edge case Step 1b describes a worktree whose branch has diverged. The code path at cmd_start.go:107-143 enters existing worktree unconditionally when directory exists and .git is valid, regardless of branch divergence. -->
- Preconditions: "A valid worktree for <slug> exists but its branch has diverged from the original source branch"
  fixture_spec:
    entities:
      - entity_type: "Worktree"
        min_count: 1
        field_constraints:
          - field: "state"
            value: "valid (directory exists, .git file present)"
          - field: "branch_status"
            value: "diverged from source branch"
- Input: "forge worktree start <slug> --source-branch develop"
- Output: "The existing worktree is entered without recreating or rebasing. stderr contains 'warning: worktree already exists, ignoring --source-branch'. stderr contains 'entering existing worktree'. No attempt to update the branch. A new Claude session is launched."
- State: "The worktree branch remains unchanged. No git operations that modify branch state are performed."
- Side-effect: "Claude session launched in existing worktree directory"

## Journey Invariants

- --source-branch is always ignored when the worktree already exists, with a warning emitted to stderr
- No file system mutations occur when entering an existing worktree (no re-copying of includes)
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
```
