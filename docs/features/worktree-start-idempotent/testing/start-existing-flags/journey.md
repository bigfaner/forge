---
feature: "worktree-start-idempotent"
journey: "start-existing-flags"
risk_level: "Medium"
golden_path: false
surface_types: ["cli"]
surface_keys: ["cli"]
sources:
  - docs/proposals/worktree-start-idempotent/proposal.md
generated: "2026-06-09"
---

# Journey: start-existing-flags

**Risk Level**: Medium

<!-- Risk Classification Criteria:
  Medium = Workflow involves multi-step interaction with various flag combinations, no irreversible side effects (worktree already exists, no creation/deletion)
-->

## Overview

A user runs `forge worktree start` with various flags (`--source-branch`, `--no-launch`, `--interactive`) when the worktree already exists. Each flag combination must behave correctly under the idempotent semantics: existing worktree is reused, inapplicable flags are handled gracefully.

## Setup

- A git repository with forge initialized
- The forge CLI binary is compiled and available in PATH
- A valid worktree already exists at `.forge/worktrees/existing-feature`
- The existing worktree was created from the `main` branch
- The `worktree.includes` config lists at least one file

## Happy Path

### Step 1: Start with --source-branch on existing worktree

**User Action**: Run `forge worktree start existing-feature --source-branch develop`

**Expected Result**: The existing worktree is entered. stderr contains a warning about ignoring `--source-branch`. No branch change occurs. A new Claude session is launched.

### Step 2: Start with --no-launch on existing worktree

**User Action**: Run `forge worktree start existing-feature --no-launch`

**Expected Result**: The existing worktree is validated. stdout outputs the worktree path. No Claude session is launched. Exit code is 0.

### Step 3: Start with --interactive on existing worktree

**User Action**: Run `forge worktree start --interactive` and select `existing-feature` from the list

**Expected Result**: The existing worktree is entered. Behavior is identical to running `forge worktree start existing-feature` explicitly. stderr contains `entering existing worktree`. A new Claude session is launched.

## Edge Cases

### Step 1b: Start with --no-launch on non-existent worktree

**Precondition**: No worktree exists for `new-feature`

**User Action**: Run `forge worktree start new-feature --no-launch`

**Expected Result**: A new worktree is created. Includes files are copied. stdout outputs the worktree path. No Claude session is launched. Exit code is 0.

### Step 2b: Start with --source-branch and --no-launch combined on existing worktree

**Precondition**: A valid worktree exists for `existing-feature`

**User Action**: Run `forge worktree start existing-feature --source-branch develop --no-launch`

**Expected Result**: `--source-branch` is ignored with a warning. Worktree path is output to stdout. No Claude session is launched. Exit code is 0.

### Step 3b: Start with --interactive and no worktrees exist

**Precondition**: No worktrees exist in `.forge/worktrees/`

**User Action**: Run `forge worktree start --interactive`

**Expected Result**: The interactive prompt shows no existing worktrees to select. The user is prompted to enter a new slug. A new worktree is created and a new Claude session is launched.

## Journey Invariants

- `--no-launch` always suppresses Claude session launch regardless of whether worktree is new or existing
- `--source-branch` is always ignored when the worktree already exists, with a warning emitted to stderr
- `--interactive` mode behavior is consistent with explicit slug specification after selection
- Exit code is 0 for all successful flag combinations
- No file system mutations occur when entering an existing worktree (no re-copying of includes)
