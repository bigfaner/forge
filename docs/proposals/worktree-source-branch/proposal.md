---
created: 2026-05-17
author: "faner"
status: Approved
---

# Proposal: Worktree Source Branch and File Copy

## Problem

`forge worktree start` always creates the worktree from HEAD, with no way to specify a different source branch and no mechanism to copy gitignored local files (e.g., `.env`) into the new worktree.

### Evidence

- The `runWorktreeStart` function in `forge-cli/internal/cmd/worktree.go` hard-codes `git worktree add -b <slug> <path>` from HEAD — no branch parameter exists
- `.env` and other gitignored local config files are missing in new worktrees, requiring manual copy each time
- The original `forge-worktree` proposal (status: Completed) noted `.worktreeinclude` processing as a known limitation (line 122) but did not address it

### Urgency

Developers frequently need to base worktrees on non-main branches (e.g., `develop`, `v3.0.0`) and must manually switch branches and copy `.env` files for every new worktree. This friction grows linearly with the number of parallel features.

## Proposed Solution

Add a `worktree` config section to `.forge/config.yaml` with `source-branch` and `copy-files` fields, plus a `--source-branch` / `-b` CLI flag on `forge worktree start`. Priority: flag > config > HEAD. Copy files from project root to worktree after creation. Abort if any configured copy-file is missing.

### Innovation Highlights

Straightforward enhancement to the existing worktree command. The copy-files mechanism addresses the gitignored-file gap by copying from the project root (where local env files exist) rather than relying on git tracking.

## Requirements Analysis

### Key Scenarios

1. **Default (no config, no flag)**: `forge worktree start my-feature` — creates from HEAD, same as today. No copy-files.
2. **Config-only**: `worktree.source-branch: develop` in config — creates from `develop`. No flag needed.
3. **Flag override**: `forge worktree start my-feature --source-branch v3.0.0` — overrides config to `v3.0.0`.
4. **Copy-files**: `worktree.copy-files: [.env]` in config — copies `.env` from project root to worktree after creation.
5. **Missing copy-file**: Config specifies `.env` but it doesn't exist in project root — error and abort before creating worktree.
6. **Combined**: Config has `source-branch: develop` and `copy-files: [.env, .env.local]` — creates from `develop`, copies both files.
7. **Flag with copy-files**: `--source-branch main` overrides config default but copy-files still apply.

### Non-Functional Requirements

- **Backward compatible**: No config = no behavior change. Existing `worktree start` calls work identically.
- **Cross-platform**: Path handling via `filepath.Join()`, tested on Windows.

### Constraints & Dependencies

- Requires existing `forge worktree start` command (already implemented)
- Config changes must be backward-compatible with existing `.forge/config.yaml` files that lack the `worktree` section

## Alternatives & Industry Benchmarking

### Industry Solutions

Git worktree itself supports specifying a branch via `git worktree add -b <branch> <path> <start-point>`. Most worktree wrapper tools (like `gw` or IDE integrations) allow specifying the base branch. Copying gitignored files is typically handled by project-specific scripts or make targets.

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | Zero cost | Manual branch switching and file copying every time | Rejected: friction scales with usage |
| CLI flag only | `--source-branch` | Flexible | No project-level default; copy-files still manual | Rejected: half solution |
| CLI flag + config.yaml | Forge CLI | Project defaults + per-invocation override; copy-files automated | New config section to maintain | **Selected: solves both pain points with minimal surface area** |

## Feasibility Assessment

### Technical Feasibility

High. All building blocks exist:
- `git worktree add -b <branch> <path> <start-point>` natively supports specifying a start point
- `ForgeConfig` struct is already modular — adding a `Worktree` section is straightforward
- File copy uses standard `os` / `io` packages

### Resource & Timeline

Single developer, estimated 0.5-1 day. Scope is well-bounded (one command enhancement + config extension).

### Dependency Readiness

All dependencies available: existing worktree command, config loading infrastructure, git worktree support.

## Scope

### In Scope

- `WorktreeConfig` struct in `ForgeConfig` with `source-branch` and `copy-files` fields
- `--source-branch` / `-b` flag on `forge worktree start`
- Source resolution: flag > config > HEAD
- Copy-files from project root to worktree after `git worktree add`
- Error and abort if any copy-file is missing from project root
- Update config JSON schema and example YAML
- Update CLI help text
- Tests for all new logic

### Out of Scope

- Lifecycle hooks (post-create hooks for arbitrary commands)
- Glob patterns in copy-files (literal paths only)
- Changes to `worktree resume`, `worktree remove`, or `worktree list`
- Auto-detection of default branch from remote

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| Source branch doesn't exist locally or remotely | M | M | Pre-validate branch exists before worktree creation; clear error message |
| Copy-file path traversal (e.g., `../../etc/passwd`) | L | H | Validate paths are relative and within project root; reject absolute or `..` paths |
| Config migration confusion | L | L | Missing `worktree` section = no behavior change; fully backward compatible |
| Copy file conflicts with worktree checkout | L | L | Copy after `git worktree add`; if file exists in checkout, overwrite with project root version |

## Success Criteria

- [ ] `forge worktree start <slug>` without config/flag behaves identically to current version
- [ ] `forge worktree start <slug> --source-branch develop` creates worktree from `develop`
- [ ] `worktree.source-branch` in config.yaml sets default source branch
- [ ] Flag overrides config default
- [ ] `worktree.copy-files` copies listed files from project root to worktree
- [ ] Worktree creation aborts with clear error if any copy-file is missing
- [ ] Copy-files rejects absolute paths and `..` traversals
- [ ] All new logic has test coverage in the forge-cli test suite

## Next Steps

- Proceed to task generation via `/quick-tasks`
