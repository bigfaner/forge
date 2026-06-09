---
journey: "idempotent-start"
step: 3
step-action: "Verify worktree state after both invocations"
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
    command: "git worktree list"
    subcommand: ""
    flags: []
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

# Contract: idempotent-start / Step 3: Verify worktree state after both invocations

> **Note**: Contracts generated without eval-journey verification (SKIP_EVAL_GATE=true). Review with extra scrutiny.

<!-- gen-contracts: do not edit manually. Regenerate via /gen-contracts. -->

## Outcome "single-entry"
- Preconditions: "Both Step 1 (create) and Step 2 (re-entry) have been executed successfully"
  fixture_spec:
    entities:
      - entity_type: "Worktree"
        min_count: 1
        field_constraints:
          - field: "slug"
            value: "my-feature"
          - field: "valid"
            value: true
- Input: "git worktree list"
- Output: "worktree list output shows exactly one entry for my-feature; entry has a valid path and branch name"
- State: "worktree state is unchanged; directory structure is intact with .git file pointing to correct location"
- Side-effect: "none"

## Outcome "git-file-valid"
<!-- source: inferred -->
<!-- reasoning: Fact Table shows git worktree validity depends on .git file pointing to correct location (git.go:137-143 IsInsideWorktree checks .git file existence and type). After creation and re-entry, .git must still be a file, not a directory -->
- Preconditions: "Worktree created by Step 1 and re-entered by Step 2"
  fixture_spec:
    entities:
      - entity_type: "Worktree"
        min_count: 1
        field_constraints:
          - field: "slug"
            value: "my-feature"
          - field: "gitFileType"
            value: "file pointing to main repo worktrees directory"
- Input: "check that .forge/worktrees/my-feature/.git is a file containing gitdir reference"
- Output: ".git file exists and contains a valid gitdir reference pointing to the main repository's worktrees directory"
- State: "worktree structure is intact"
- Side-effect: "none"

## Journey Invariants

- The worktree directory structure remains valid (.git file present and pointing to correct location) after every invocation
- stderr always contains either 'created new worktree' or 'entering existing worktree' to distinguish the path taken
