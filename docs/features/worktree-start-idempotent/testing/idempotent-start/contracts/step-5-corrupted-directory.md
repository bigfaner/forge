---
journey: "idempotent-start"
step: 5
step-action: "Start when worktree directory exists but .git file is missing"
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

# Contract: idempotent-start / Step 5: Start when worktree directory exists but .git file is missing

> **Note**: Contracts generated without eval-journey verification (SKIP_EVAL_GATE=true). Review with extra scrutiny.

<!-- gen-contracts: do not edit manually. Regenerate via /gen-contracts. -->

## Outcome "corrupted-directory"
- Preconditions: "Directory .forge/worktrees/broken-feature exists but has no .git file (manually created or corrupted)"
  fixture_spec:
    entities:
      - entity_type: "Directory"
        min_count: 1
        field_constraints:
          - field: "path"
            value: ".forge/worktrees/broken-feature"
          - field: "gitFilePresent"
            value: false
- Input: "forge worktree start broken-feature"
- Output: "exit code is non-zero; stderr contains error indicating the worktree is corrupted; stderr suggests running 'forge worktree remove broken-feature' to clean up and retry; no Claude session is launched"
- State: "no changes to file system; corrupted directory remains; no new worktree created"
- Side-effect: "none"

## Outcome "symlink-resolution-failure"
<!-- source: inferred -->
<!-- reasoning: Fact Table shows filepath.EvalSymlinks called at cmd_start.go:109-112; if symlink resolution fails, AIError with INVALID_INPUT is returned -->
- Preconditions: "Directory .forge/worktrees/symlink-feature exists but is a broken symlink that cannot be resolved"
  fixture_spec:
    entities:
      - entity_type: "Directory"
        min_count: 1
        field_constraints:
          - field: "path"
            value: ".forge/worktrees/symlink-feature"
          - field: "brokenSymlink"
            value: true
- Input: "forge worktree start symlink-feature"
- Output: "exit code is non-zero; stderr contains error about unable to resolve target path"
- State: "no changes to file system; broken symlink remains"
- Side-effect: "none"

## Journey Invariants

- The worktree directory structure remains valid (.git file present and pointing to correct location) after every invocation
- Error messages always include a suggested recovery command when corruption is detected
