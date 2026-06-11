---
id: "2"
title: "Add tests for remote branch detection and fetch-failure fallback"
priority: "P0"
estimated_time: "1h"
dependencies: ["1"]
scope: "backend"
breaking: false
type: "enhancement"
mainSession: false
---

# 2: Add tests for remote branch detection and fetch-failure fallback

## Description

Add table-driven tests to `worktree_test.go` covering the new remote branch detection paths introduced in task 1. Tests must cover the happy path (remote branch reuse), fetch failure graceful degradation, and source-branch ignored when remote branch exists.

## Reference Files
- `docs/proposals/worktree-remote-branch-reuse/proposal.md` — Source proposal
- `forge-cli/internal/cmd/worktree_test.go` — Existing test file (2110 lines)
- `forge-cli/internal/cmd/worktree.go` — Implementation under test

## Acceptance Criteria

- New test: remote branch `origin/<slug>` exists, local does not → worktree created from remote branch, source-branch ignored
- New test: fetch fails (network error) → falls back to source-branch behavior, no error returned
- New test: both remote and local branches absent → creates from source-branch (existing behavior, regression guard)
- New test: `--source-branch` flag explicitly set but remote branch exists → remote branch wins (consistent with local-branch behavior)
- All existing tests continue to pass

## Hard Rules

- Use table-driven test pattern consistent with existing tests in `worktree_test.go`
- Mock git commands via the existing `gitRunFunc` variable — do NOT shell out to real git
- Follow TDD: write failing tests first, then verify they pass against task 1 implementation

## Implementation Notes

- Existing test `TestWorktreeStart_SourceBranchNotUsedForExistingBranch` covers the local branch case; new tests should mirror its structure but set up remote branch refs instead of local
- The mock for `gitRunFunc` needs to handle `fetch origin` and `rev-parse --verify remotes/origin/<slug>` commands in addition to existing `worktree add` calls
- Check existing mock patterns for how `gitRunFunc` dispatches on command args
