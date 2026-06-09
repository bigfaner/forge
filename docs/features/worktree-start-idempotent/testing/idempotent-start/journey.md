---
feature: "worktree-start-idempotent"
journey: "idempotent-start"
risk_level: "High"
golden_path: true
surface_types: ["cli"]
surface_keys: ["cli"]
sources:
  - docs/proposals/worktree-start-idempotent/proposal.md
generated: "2026-06-09"
---

# Journey: idempotent-start

**Risk Level**: High

<!-- Risk Classification Criteria:
  High = Workflow involves state mutation (creating worktree, launching Claude session), file system changes
-->

## Overview

A user runs `forge worktree start <slug>` to begin working on a feature. The command must be idempotent: when no worktree exists for the slug, it creates one and launches a new Claude session; when a worktree already exists, it skips creation and launches a new Claude session. This is the primary user workflow and the core value proposition of the feature.

## Setup

- A git repository with forge initialized (`.forge/` directory present)
- The forge CLI binary is compiled and available in PATH
- The repository is on the main/default branch
- No worktree directory exists at `.forge/worktrees/<slug>` for the target slug

## Happy Path

### Step 1: Start with non-existent worktree (create path)

**User Action**: Run `forge worktree start my-feature`

**Expected Result**: A new worktree is created at `.forge/worktrees/my-feature` using the current branch (or configured source branch). The `includes` files are copied into the worktree directory. stderr contains the keyword `created new worktree`. A new Claude session is launched in the worktree directory.

### Step 2: Start with existing worktree (idempotent path)

**User Action**: Run `forge worktree start my-feature` again

**Expected Result**: No new worktree is created. The existing worktree at `.forge/worktrees/my-feature` is reused. stderr contains the keyword `entering existing worktree`. The `includes` file copying is skipped (files were copied during initial creation). A new Claude session is launched in the existing worktree directory.

### Step 3: Verify worktree state after both invocations

**User Action**: Run `git worktree list` and verify `.forge/worktrees/my-feature` appears exactly once

**Expected Result**: The worktree list shows exactly one entry for `my-feature`. The worktree directory structure is intact with the `.git` file pointing to the correct location.

### Step 4: Verify includes files were copied only once

**User Action**: Check that files listed in `worktree.includes` config exist in the worktree and match the originals

**Expected Result**: All includes files from the config are present in the worktree directory. Their contents match the source files from the main repository. No duplicate or overwritten files.

## Edge Cases

### Step 1b: Start with existing worktree that has diverged from source branch

**Precondition**: A worktree for `my-feature` exists but its branch has diverged from the original source branch

**User Action**: Run `forge worktree start my-feature`

**Expected Result**: The existing worktree is entered without recreating or rebasing. stderr contains `entering existing worktree`. No attempt to update the branch or re-sync files. A new Claude session is launched.

### Step 2b: Start with --source-branch on existing worktree

**Precondition**: A worktree for `my-feature` already exists

**User Action**: Run `forge worktree start my-feature --source-branch develop`

**Expected Result**: The `--source-branch` flag is ignored because the worktree already exists and its branch is already determined. stderr contains a warning: `worktree already exists, ignoring --source-branch`. The existing worktree is entered. A new Claude session is launched.

### Step 3b: Start without includes config

**Precondition**: No `worktree.includes` key is set in `.forge/config.yaml`, no worktree exists yet

**User Action**: Run `forge worktree start no-includes-feature`

**Expected Result**: The worktree is created successfully. No file copying step is attempted (no includes configured). stderr contains `created new worktree`. A new Claude session is launched.

### Step 4b: Start with renamed config (includes) after migration from copy-files

**Precondition**: `.forge/config.yaml` contains `worktree.includes` (not `copy-files`), no worktree exists yet

**User Action**: Run `forge worktree start migrated-feature`

**Expected Result**: The worktree is created. Files listed under `worktree.includes` are copied to the worktree. No error or warning about missing `copy-files` key. stderr contains `created new worktree`.

### Step 5b: Start when worktree directory exists but .git file is missing

**Precondition**: Directory `.forge/worktrees/broken-feature` exists but has no `.git` file (manually created or corrupted)

**User Action**: Run `forge worktree start broken-feature`

**Expected Result**: The command exits with a non-zero exit code. stderr contains an error message indicating the worktree is corrupted. stderr suggests running `forge worktree remove broken-feature` to clean up and retry. No Claude session is launched.

### Step 6b: Start with --no-launch on non-existent worktree

**Precondition**: No worktree exists for `quiet-feature`

**User Action**: Run `forge worktree start quiet-feature --no-launch`

**Expected Result**: The worktree is created. Includes files are copied. stdout contains the worktree path. No Claude session is launched. Exit code is 0.

## Journey Invariants

- Every invocation of `forge worktree start <slug>` either creates a new worktree or enters an existing one -- it never fails due to "worktree already exists" when the worktree is valid
- stderr always contains either `created new worktree` or `entering existing worktree` to distinguish the path taken
- `includes` files are only copied during worktree creation, never during re-entry
- A new Claude session is always launched unless `--no-launch` is specified
- The worktree directory structure remains valid (`.git` file present and pointing to correct location) after every invocation
