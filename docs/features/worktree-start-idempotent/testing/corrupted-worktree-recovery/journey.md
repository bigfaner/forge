---
feature: "worktree-start-idempotent"
journey: "corrupted-worktree-recovery"
risk_level: "Medium"
golden_path: false
surface_types: ["cli"]
surface_keys: ["cli"]
sources:
  - docs/proposals/worktree-start-idempotent/proposal.md
generated: "2026-06-09"
---

# Journey: corrupted-worktree-recovery

**Risk Level**: Medium

<!-- Risk Classification Criteria:
  Medium = Workflow involves error detection and user guidance, no direct data loss but requires user to take corrective action
-->

## Overview

A user attempts to start a worktree that has a directory present but is in a corrupted state (missing `.git` file, invalid git worktree reference). The command must detect the corruption, report a clear error, and guide the user toward recovery.

## Setup

- A git repository with forge initialized
- The forge CLI binary is compiled and available in PATH
- A directory exists at `.forge/worktrees/corrupted-feature` but the `.git` file is missing or points to a non-existent location

## Happy Path

### Step 1: Attempt to start a corrupted worktree

**User Action**: Run `forge worktree start corrupted-feature`

**Expected Result**: The command exits with a non-zero exit code. stderr contains an error message indicating the worktree is corrupted or invalid. The error message suggests running `forge worktree remove corrupted-feature` to clean up. No Claude session is launched.

### Step 2: Remove the corrupted worktree

**User Action**: Run `forge worktree remove corrupted-feature`

**Expected Result**: The corrupted directory is cleaned up. The worktree entry is removed from git's worktree list (if present). Exit code is 0.

### Step 3: Retry start after cleanup

**User Action**: Run `forge worktree start corrupted-feature`

**Expected Result**: A new worktree is created successfully. stderr contains `created new worktree`. Includes files are copied. A new Claude session is launched.

## Edge Cases

### Step 1b: Worktree directory exists with .git file pointing to deleted git directory

**Precondition**: `.forge/worktrees/dangling-feature` has a `.git` file that references a git directory that no longer exists

**User Action**: Run `forge worktree start dangling-feature`

**Expected Result**: The command detects the invalid reference. stderr contains a corruption error. The command suggests `forge worktree remove dangling-feature`. Exit code is non-zero.

### Step 2b: Attempt to start when .forge/worktrees is not a directory

**Precondition**: `.forge/worktrees` exists as a file (not a directory)

**User Action**: Run `forge worktree start any-feature`

**Expected Result**: The command reports a filesystem error indicating the worktrees path is not a directory. Exit code is non-zero.

### Step 3b: Remove corrupted worktree that git does not recognize

**Precondition**: Directory `.forge/worktrees/orphan-feature` exists but `git worktree list` does not show it (orphan directory, not a real git worktree)

**User Action**: Run `forge worktree remove orphan-feature`

**Expected Result**: The orphan directory is removed. No error from git about missing worktree entry. Exit code is 0.

## Journey Invariants

- Corrupted worktree detection always prevents launching a Claude session
- Error messages always include a suggested recovery command (`forge worktree remove <slug>`)
- Recovery (remove + start) always restores the worktree to a valid state
- Exit code is non-zero for all error cases, zero for successful operations
- The original repository's git state is never corrupted by the recovery process
