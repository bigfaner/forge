---
journey: "corrupted-worktree-recovery"
step: 3
step-action: "Retry start after cleanup"
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

# Contract: corrupted-worktree-recovery / Step 3: Retry start after cleanup

> **Note**: Contracts generated without eval-journey verification (SKIP_EVAL_GATE=true). Review with extra scrutiny.

<!-- gen-contracts: do not edit manually. Regenerate via /gen-contracts. -->

## Outcome "success"
- Preconditions: "The corrupted worktree has been removed. No directory exists at .forge/worktrees/<slug>. The git repository is in a clean state."
  fixture_spec:
    entities:
      - entity_type: "Worktree"
        min_count: 0
      - entity_type: "Repository"
        min_count: 1
        field_constraints:
          - field: "state"
            value: "clean, forge initialized, on main/default branch"
- Input: "forge worktree start <slug> where <slug> was previously corrupted and has been cleaned up"
- Output: "A new worktree is created successfully. stderr contains 'created new worktree'. Includes files are copied if configured. A new Claude session is launched in the worktree directory."
- State: "New git worktree added at .forge/worktrees/<slug>. New branch created. Worktree directory structure is valid with .git file pointing to correct location. Includes files copied from project root."
- Side-effect: "Claude session launched in the new worktree directory"

## Outcome "worktrees-dir-not-directory"
<!-- source: inferred -->
<!-- reasoning: The start command creates .forge/worktrees/ via os.MkdirAll (cmd_start.go:102). If .forge/worktrees exists as a file (not a directory), MkdirAll will fail and produce a filesystem error. -->
- Preconditions: ".forge/worktrees exists as a regular file (not a directory) in the project"
  fixture_spec:
    entities:
      - entity_type: "Repository"
        min_count: 1
        field_constraints:
          - field: "worktrees_path_type"
            value: "file (not directory)"
- Input: "forge worktree start <slug> when .forge/worktrees is a file"
- Output: "Exit code is non-zero. stderr contains an error message about failing to create worktrees directory."
- State: "No worktree is created. The filesystem state is unchanged."
- Side-effect: "none"

## Outcome "source-branch-not-found"
<!-- source: inferred -->
<!-- reasoning: Fact Table WT_BRANCH_RESOLUTION shows Layer 3 validates source branch via git rev-parse --verify. If an invalid source-branch is specified, the command fails with source branch not found error (cmd_start.go:207-210). -->
- Preconditions: "No worktree exists for the slug. A --source-branch flag specifies a branch that does not exist locally or remotely."
  fixture_spec:
    entities:
      - entity_type: "Worktree"
        min_count: 0
      - entity_type: "Branch"
        min_count: 0
- Input: "forge worktree start <slug> --source-branch nonexistent-branch"
- Output: "Exit code is non-zero. stderr contains error message indicating the source branch was not found. stderr includes a hint to verify the branch exists."
- State: "No worktree is created. No branch is created. The repository state is unchanged."
- Side-effect: "none"

## Journey Invariants

- Recovery (remove + start) always restores the worktree to a valid state
- Exit code is zero for successful operations
- A new Claude session is always launched on successful creation
- The worktree directory structure remains valid (.git file present and pointing to correct location) after creation

## Fixture Specification

This Contract requires the following pre-existing data state. See `rules/fixture-spec.md` for schema details.

```yaml
fixture_spec:
  entities:
    - entity_type: "Worktree"
      min_count: 0
    - entity_type: "Repository"
      min_count: 1
      field_constraints:
        - field: "state"
          value: "clean, forge initialized"
```
