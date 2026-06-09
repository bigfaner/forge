---
journey: "idempotent-start"
step: 6
step-action: "Start with --no-launch on non-existent worktree"
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
    flags: ["no-launch"]
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

# Contract: idempotent-start / Step 6: Start with --no-launch on non-existent worktree

> **Note**: Contracts generated without eval-journey verification (SKIP_EVAL_GATE=true). Review with extra scrutiny.

<!-- gen-contracts: do not edit manually. Regenerate via /gen-contracts. -->

## Outcome "success-no-launch"
- Preconditions: "No worktree exists for quiet-feature; --no-launch flag set; claude binary check is skipped"
  fixture_spec:
    entities:
      - entity_type: "GitRepository"
        min_count: 1
        field_constraints:
          - field: "forgeInitialized"
            value: true
      - entity_type: "ForgeConfig"
        min_count: 1
        field_constraints:
          - field: "worktree.includes"
            value: "any list of valid file paths or empty"
- Input: "forge worktree start quiet-feature --no-launch"
- Output: "exit code is 0; stderr contains 'created new worktree: quiet-feature'; stdout contains the worktree path; no Claude session is launched"
- State: "new worktree created at .forge/worktrees/quiet-feature; includes files copied; new branch created"
- Side-effect: "git branch and worktree add executed; includes files copied; no claude process spawned"

## Outcome "no-launch-existing"
<!-- source: inferred -->
<!-- reasoning: --no-launch with existing worktree should work too. Fact Table shows cmd_start.go:128-131 outputs worktree path and returns nil for existing worktree with --no-launch -->
- Preconditions: "Valid worktree already exists at .forge/worktrees/quiet-feature; --no-launch flag set"
  fixture_spec:
    entities:
      - entity_type: "Worktree"
        min_count: 1
        field_constraints:
          - field: "slug"
            value: "quiet-feature"
          - field: "valid"
            value: true
- Input: "forge worktree start quiet-feature --no-launch"
- Output: "exit code is 0; stderr contains 'entering existing worktree: quiet-feature'; stdout contains the resolved worktree path; no Claude session is launched"
- State: "existing worktree unchanged; no includes re-copied"
- Side-effect: "none"

## Journey Invariants

- A new Claude session is always launched unless --no-launch is specified
- stderr always contains either 'created new worktree' or 'entering existing worktree' to distinguish the path taken
- includes files are only copied during worktree creation, never during re-entry
