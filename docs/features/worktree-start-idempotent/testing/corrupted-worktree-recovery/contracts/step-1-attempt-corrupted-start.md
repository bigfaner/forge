---
journey: "corrupted-worktree-recovery"
step: 1
step-action: "Attempt to start a corrupted worktree"
generated: "2026-06-09"
skip_eval: true
sources:
  - docs/features/worktree-start-idempotent/testing/corrupted-worktree-recovery/journey.md

anchors:
  cli:
    command: "forge worktree start"
    subcommand: ""
    flags: ["slug"]
    aliases: []
last_anchor_sync: ""
---

# Contract: corrupted-worktree-recovery / Step 1: Attempt to start a corrupted worktree

> **Note**: Contracts generated without eval-journey verification (SKIP_EVAL_GATE=true). Review with extra scrutiny.

<!-- gen-contracts: do not edit manually. Regenerate via /gen-contracts. -->

## Outcome "corruption-detected"
- Preconditions: "Directory .forge/worktrees/<slug> exists but the .git file is missing or corrupt"
  fixture_spec:
    entities:
      - entity_type: "Worktree"
        min_count: 1
        field_constraints:
          - field: "state"
            value: "corrupted (directory exists, .git file missing)"
- Input: "forge worktree start <slug> where <slug> corresponds to a corrupted worktree directory"
- Output: "Exit code is non-zero. stderr contains an error message indicating the worktree directory exists but is not a valid git worktree. stderr includes a hint suggesting to run forge worktree remove <slug> and try again."
- State: "No new worktree is created. No git worktree add operation is attempted. No Claude session is launched. The corrupted directory remains unchanged."
- Side-effect: "none"
- Invariants: "Corrupted worktree detection prevents launching a Claude session"

## Outcome "not-found"
<!-- source: surface-required (CLI resource access step) -->
- Preconditions: "No directory exists at .forge/worktrees/<slug> for the given slug"
  fixture_spec:
    entities:
      - entity_type: "Worktree"
        min_count: 0
- Input: "forge worktree start <slug> where no worktree directory exists"
- Output: "A new worktree is created successfully. stderr contains 'created new worktree'. Includes files are copied if configured. A new Claude session is launched."
- State: "New git worktree added at .forge/worktrees/<slug>. New branch created if not existing. Worktree directory structure is valid with .git file present."
- Side-effect: "Claude session launched in worktree directory"

## Outcome "dangling-git-reference"
<!-- source: inferred -->
<!-- reasoning: Fact Table WT_CORRUPTION_DETECTION shows .git file is checked via os.Stat; a .git file pointing to a non-existent git directory passes the Stat check but the worktree is still invalid at the git level -->
- Preconditions: "Directory .forge/worktrees/<slug> exists and has a .git file that references a git directory that no longer exists on disk"
  fixture_spec:
    entities:
      - entity_type: "Worktree"
        min_count: 1
        field_constraints:
          - field: "state"
            value: "dangling (directory exists, .git file points to deleted git dir)"
- Input: "forge worktree start <slug> where the worktree has a dangling .git reference"
- Output: "The command detects the invalid reference. stderr contains a corruption error message. The command suggests running forge worktree remove <slug>. Exit code is non-zero."
- State: "No new worktree is created. No Claude session is launched. The worktree directory with dangling reference remains unchanged."
- Side-effect: "none"
- Invariants: "Error messages always include a suggested recovery command"

## Journey Invariants

- Corrupted worktree detection always prevents launching a Claude session
- Error messages always include a suggested recovery command (forge worktree remove <slug>)
- Exit code is non-zero for all error cases
- The original repository's git state is never corrupted by the detection process

## Fixture Specification

This Contract requires the following pre-existing data state. See `rules/fixture-spec.md` for schema details.

```yaml
fixture_spec:
  entities:
    - entity_type: "Worktree"           # Corrupted or dangling worktree directory
      min_count: 1                      # At least one worktree in corrupted state
      field_constraints:
        - field: "state"
          value: "corrupted or dangling (.git file missing or invalid)"
```
