---
journey: "idempotent-start"
step: 4
step-action: "Verify includes files were copied only once"
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
    command: "file comparison"
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

# Contract: idempotent-start / Step 4: Verify includes files were copied only once

> **Note**: Contracts generated without eval-journey verification (SKIP_EVAL_GATE=true). Review with extra scrutiny.

<!-- gen-contracts: do not edit manually. Regenerate via /gen-contracts. -->

## Outcome "files-match"
- Preconditions: "worktree.includes config lists at least one file; Step 1 (create) and Step 2 (re-entry) both completed"
  fixture_spec:
    entities:
      - entity_type: "ForgeConfig"
        min_count: 1
        field_constraints:
          - field: "worktree.includes"
            value: "non-empty list of valid file paths"
      - entity_type: "Worktree"
        min_count: 1
        field_constraints:
          - field: "slug"
            value: "my-feature"
          - field: "valid"
            value: true
- Input: "compare each file in worktree.includes between project root and .forge/worktrees/my-feature"
- Output: "all includes files exist in the worktree directory; file contents match the originals from the project root exactly"
- State: "files were copied only during Step 1 creation; Step 2 re-entry did not modify or re-copy any files"
- Side-effect: "none"

## Outcome "no-includes-config"
<!-- source: inferred -->
<!-- reasoning: Journey edge case Step 3b: no worktree.includes key set. Fact Table shows WorktreeConfig.Includes is []string yaml:includes (config.go:72-75). When nil or empty, no copy attempt is made (helpers.go:152-154) -->
- Preconditions: "No worktree.includes key in .forge/config.yaml; no worktree exists yet"
  fixture_spec:
    entities:
      - entity_type: "ForgeConfig"
        min_count: 1
        field_constraints:
          - field: "worktree.includes"
            value: "empty or not present"
      - entity_type: "Worktree"
        min_count: 1
        field_constraints:
          - field: "slug"
            value: "no-includes-feature"
- Input: "forge worktree start no-includes-feature"
- Output: "stderr contains 'created new worktree: no-includes-feature'; no file copying step is attempted; Claude session launched"
- State: "worktree created with no extra files; only default git worktree contents present"
- Side-effect: "none"

## Outcome "migrated-config"
<!-- source: inferred -->
<!-- reasoning: Journey edge case Step 4b: config uses new 'includes' key after migration from 'copy-files'. Fact Table shows WorktreeConfig struct uses yaml:includes tag (config.go:74) -->
- Preconditions: ".forge/config.yaml contains worktree.includes (not copy-files); no worktree exists yet"
  fixture_spec:
    entities:
      - entity_type: "ForgeConfig"
        min_count: 1
        field_constraints:
          - field: "worktree.includes"
            value: "valid list of file paths"
          - field: "worktree.copy-files"
            value: "not present"
- Input: "forge worktree start migrated-feature"
- Output: "stderr contains 'created new worktree: migrated-feature'; files listed under worktree.includes are copied to the worktree; no error or warning about missing copy-files key"
- State: "worktree created; includes files successfully copied"
- Side-effect: "files copied from project root to worktree directory"

## Journey Invariants

- includes files are only copied during worktree creation, never during re-entry
- The worktree directory structure remains valid (.git file present and pointing to correct location) after every invocation
