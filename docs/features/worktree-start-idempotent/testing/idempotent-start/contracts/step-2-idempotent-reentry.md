---
journey: "idempotent-start"
step: 2
step-action: "Start with existing worktree (idempotent path)"
generated: "2026-06-09"
sources:
  - docs/features/worktree-start-idempotent/testing/idempotent-start/journey.md
skip_eval: true

anchors:
  api:
    endpoint: ""
    method: ""
    content_type: ""
  cli:
    command: "forge worktree start"
    subcommand: ""
    flags: ["source-branch", "no-launch", "interactive"]
    aliases: []
  tui:
    command: ""
    interactive_prompt: ""
    keybindings: []
  web:
    page: ""
    route: ""
  mobile:
    screen: ""
    navigation_path: []
    deeplink: ""
    platform: ""

last_anchor_sync: ""
---

# Contract: idempotent-start / Step 2: Start with existing worktree (idempotent path)

> **Note**: Contracts generated without eval-journey verification (SKIP_EVAL_GATE=true). Review with extra scrutiny.

<!-- gen-contracts: do not edit manually. Regenerate via /gen-contracts. -->

## Outcome "success"
- Preconditions: "Valid worktree already exists at .forge/worktrees/my-feature; .git file present and valid; claude binary available in PATH"
  fixture_spec:
    entities:
      - entity_type: "Worktree"
        min_count: 1
        field_constraints:
          - field: "slug"
            value: "my-feature"
          - field: "valid"
            value: true
- Input: "forge worktree start my-feature"
- Output: "stderr contains 'entering existing worktree: my-feature'; no 'created new worktree' message; a new Claude session is launched in the worktree directory"
- State: "no new worktree created; existing worktree directory unchanged; no includes files re-copied; no git branch changes"
- Side-effect: "claude CLI process spawned with --dangerously-skip-permissions in existing worktree directory"

## Outcome "corrupted-git-file"
- Preconditions: "Directory exists at .forge/worktrees/my-feature but .git file is missing or corrupt"
  fixture_spec:
    entities:
      - entity_type: "Directory"
        min_count: 1
        field_constraints:
          - field: "path"
            value: ".forge/worktrees/my-feature"
          - field: "gitFilePresent"
            value: false
- Input: "forge worktree start my-feature"
- Output: "exit code is non-zero; stderr contains error about worktree directory exists but not a valid git worktree; stderr contains hint to run 'forge worktree remove my-feature'"
- State: "no changes to file system; no Claude session launched"
- Side-effect: "none"

## Outcome "source-branch-ignored"
- Preconditions: "Valid worktree already exists at .forge/worktrees/my-feature; --source-branch flag is explicitly provided"
  fixture_spec:
    entities:
      - entity_type: "Worktree"
        min_count: 1
        field_constraints:
          - field: "slug"
            value: "my-feature"
          - field: "valid"
            value: true
- Input: "forge worktree start my-feature --source-branch develop"
- Output: "stderr contains 'warning: worktree already exists, ignoring --source-branch'; stderr contains 'entering existing worktree: my-feature'; no branch change occurs; a new Claude session is launched"
- State: "existing worktree unchanged; no branch modifications; includes files not re-copied"
- Side-effect: "claude CLI process spawned in existing worktree directory"

## Outcome "diverged-branch-entered"
<!-- source: inferred -->
<!-- reasoning: Journey edge case Step 1b: existing worktree whose branch has diverged from source branch. Code at cmd_start.go:107-143 shows existing path does not check or update branch, only validates .git file -->
- Preconditions: "Valid worktree exists at .forge/worktrees/my-feature; worktree branch has diverged from the original source branch"
  fixture_spec:
    entities:
      - entity_type: "Worktree"
        min_count: 1
        field_constraints:
          - field: "slug"
            value: "my-feature"
          - field: "valid"
            value: true
          - field: "branchDiverged"
            value: true
- Input: "forge worktree start my-feature"
- Output: "stderr contains 'entering existing worktree: my-feature'; no attempt to update or rebase the branch; a new Claude session is launched"
- State: "worktree branch state unchanged; no rebasing or merging"
- Side-effect: "claude CLI process spawned in existing worktree directory"

## Journey Invariants

- Every invocation of `forge worktree start <slug>` either creates a new worktree or enters an existing one -- it never fails due to "worktree already exists" when the worktree is valid
- stderr always contains either 'created new worktree' or 'entering existing worktree' to distinguish the path taken
- includes files are only copied during worktree creation, never during re-entry
- A new Claude session is always launched unless --no-launch is specified
