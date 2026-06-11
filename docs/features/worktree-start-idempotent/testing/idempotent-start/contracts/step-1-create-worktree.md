---
journey: "idempotent-start"
step: 1
step-action: "Start with non-existent worktree (create path)"
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

# Contract: idempotent-start / Step 1: Start with non-existent worktree (create path)

> **Note**: Contracts generated without eval-journey verification (SKIP_EVAL_GATE=true). Review with extra scrutiny.

<!-- gen-contracts: do not edit manually. Regenerate via /gen-contracts. -->

## Outcome "success"
- Preconditions: "No worktree directory exists at .forge/worktrees/my-feature; forge initialized in a git repository; claude binary available in PATH"
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
- Input: "forge worktree start my-feature"
- Output: "stderr contains 'created new worktree: my-feature'; a new Claude session is launched in the worktree directory with --dangerously-skip-permissions flag"
- State: "new worktree directory created at .forge/worktrees/my-feature with .git file pointing to main repo; new git branch 'my-feature' created; includes files copied if configured"
- Side-effect: "claude CLI process spawned with --dangerously-skip-permissions; git branch and worktree add executed; includes files copied from project root to worktree directory"

## Outcome "slug-not-found"
<!-- source: inferred -->
<!-- reasoning: CLI surface requires not-found for resource access. forge worktree start validates the slug argument is provided; when missing, ErrSlugRequired is returned (cmd_start.go:70-71, base/errors.go:380-388) -->
- Preconditions: "No slug argument provided and --interactive flag not set"
  fixture_spec:
    entities:
      - entity_type: "GitRepository"
        min_count: 1
        field_constraints:
          - field: "forgeInitialized"
            value: true
- Input: "forge worktree start (no arguments)"
- Output: "exit code is non-zero; stderr contains error about slug being required"
- State: "no worktree created; no branch changes; no Claude session launched"
- Side-effect: "none"

## Outcome "not-git-repository"
<!-- source: inferred -->
<!-- reasoning: Fact Table shows git.IsGitRepository check at cmd_start.go:90-92; base.ErrNotGitRepository returns AIError with INVALID_INPUT code -->
- Preconditions: "Current directory is not a git repository"
  fixture_spec:
    entities:
      - entity_type: "Directory"
        min_count: 1
        field_constraints:
          - field: "hasGitDir"
            value: false
- Input: "forge worktree start my-feature"
- Output: "exit code is non-zero; stderr contains error indicating not a git repository"
- State: "no worktree created; no branch changes"
- Side-effect: "none"

## Outcome "claude-not-found"
<!-- source: inferred -->
<!-- reasoning: Fact Table shows claude binary lookup via exec.LookPath at cmd_start.go:78-80; base.NewAIError with NOT_FOUND code returned when binary missing -->
- Preconditions: "No worktree exists; claude binary not in PATH; --no-launch flag not set"
  fixture_spec:
    entities:
      - entity_type: "GitRepository"
        min_count: 1
        field_constraints:
          - field: "forgeInitialized"
            value: true
      - entity_type: "SystemPath"
        min_count: 1
        field_constraints:
          - field: "claudeBinary"
            value: "not present"
- Input: "forge worktree start my-feature"
- Output: "exit code is non-zero; stderr contains error about claude binary not found in PATH"
- State: "no worktree created; no branch changes"
- Side-effect: "none"

## Outcome "source-branch-not-found"
<!-- source: inferred -->
<!-- reasoning: Fact Table shows source branch validation via git rev-parse at cmd_start.go:206-210; ErrSourceBranchNotFound returned when branch does not exist -->
- Preconditions: "No worktree exists; --source-branch specified with non-existent branch name; no local or remote branch matching the slug"
  fixture_spec:
    entities:
      - entity_type: "GitRepository"
        min_count: 1
        field_constraints:
          - field: "forgeInitialized"
            value: true
- Input: "forge worktree start my-feature --source-branch nonexistent-branch"
- Output: "exit code is non-zero; stderr contains error about source branch not found; stderr contains hint to verify branch exists"
- State: "no worktree created; no branch changes"
- Side-effect: "none"

## Journey Invariants

- Every invocation of `forge worktree start <slug>` either creates a new worktree or enters an existing one -- it never fails due to "worktree already exists" when the worktree is valid
- stderr always contains either 'created new worktree' or 'entering existing worktree' to distinguish the path taken
- The worktree directory structure remains valid (.git file present and pointing to correct location) after every invocation
