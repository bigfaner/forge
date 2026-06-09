---
journey: "corrupted-worktree-recovery"
step: 2
step-action: "Remove the corrupted worktree"
generated: "2026-06-09"
skip_eval: true
sources:
  - docs/features/worktree-start-idempotent/testing/corrupted-worktree-recovery/journey.md

anchors:
  cli:
    command: "forge worktree remove"
    subcommand: ""
    flags: ["slug"]
    aliases: []
last_anchor_sync: ""
---

# Contract: corrupted-worktree-recovery / Step 2: Remove the corrupted worktree

> **Note**: Contracts generated without eval-journey verification (SKIP_EVAL_GATE=true). Review with extra scrutiny.

<!-- gen-contracts: do not edit manually. Regenerate via /gen-contracts. -->

## Outcome "success"
- Preconditions: "A corrupted worktree directory exists at .forge/worktrees/<slug>. The git worktree entry may or may not be present in git worktree list."
  fixture_spec:
    entities:
      - entity_type: "Worktree"
        min_count: 1
        field_constraints:
          - field: "state"
            value: "corrupted (directory exists, .git file missing or invalid)"
- Input: "forge worktree remove <slug> where <slug> corresponds to the corrupted worktree"
- Output: "Exit code is 0. stdout contains a message confirming removal. The corrupted directory is cleaned up. If git recognizes the worktree, the worktree entry is removed from git's worktree list."
- State: "The .forge/worktrees/<slug> directory is removed. The git worktree list no longer contains an entry for the slug. Any stale administrative files are pruned."
- Side-effect: "git worktree prune is executed to clean up stale administrative files"
- Invariants: "Recovery always restores the worktree to a valid state"

## Outcome "orphan-directory"
<!-- source: inferred -->
<!-- reasoning: Fact Table WT_REMOVE_NOT_FOUND checks os.Stat on targetDir, but the git worktree remove subcommand may fail if git does not recognize the directory as a worktree (orphan). The remove command at cmd_remove.go:97 uses git.Run directly which may error on non-git-recognized directories. -->
- Preconditions: "Directory .forge/worktrees/<slug> exists but git worktree list does not show it (orphan directory, not a real git worktree)"
  fixture_spec:
    entities:
      - entity_type: "Worktree"
        min_count: 1
        field_constraints:
          - field: "state"
            value: "orphan (directory exists but not tracked by git worktree)"
- Input: "forge worktree remove <slug> where the directory is an orphan not tracked by git"
- Output: "The orphan directory is removed. No error from git about missing worktree entry. Exit code is 0."
- State: "The orphan directory at .forge/worktrees/<slug> is deleted. Git worktree list is unchanged (no entry to remove)."
- Side-effect: "none"
- Invariants: "The original repository's git state is never corrupted by the recovery process"

## Outcome "not-found"
<!-- source: surface-required (CLI resource access step) -->
- Preconditions: "No directory exists at .forge/worktrees/<slug> for the given slug"
  fixture_spec:
    entities:
      - entity_type: "Worktree"
        min_count: 0
- Input: "forge worktree remove <slug> where no worktree directory exists"
- Output: "Exit code is non-zero. stderr contains an error message indicating worktree not found with the target path."
- State: "No directories are removed. No git operations are performed."
- Side-effect: "none"

## Journey Invariants

- Recovery (remove + start) always restores the worktree to a valid state
- Exit code is non-zero for error cases, zero for successful operations
- The original repository's git state is never corrupted by the recovery process
- Error messages always include a suggested recovery command

## Fixture Specification

This Contract requires the following pre-existing data state. See `rules/fixture-spec.md` for schema details.

```yaml
fixture_spec:
  entities:
    - entity_type: "Worktree"           # Corrupted worktree to be removed
      min_count: 1                      # Exactly one worktree in corrupted state
      field_constraints:
        - field: "state"
          value: "corrupted, orphan, or invalid"
```
