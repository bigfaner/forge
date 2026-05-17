---
created: 2026-05-17
author: "faner"
status: Draft
---

# Proposal: Forge Worktree Management

## Problem

Developers working on multiple Forge features must switch branches in a single working directory, forcing sequential feature development even when features are independent. There is no convenient way to work on multiple features simultaneously in isolated git worktrees.

### Evidence

- `todo.txt` item #76: "forge要支持worktree。doing" — explicitly tracked as in-progress
- `todo.txt` item #31: "task feature 要支持worktree" — listed as needed
- The forge-cli already implements worktree name detection (`GetWorktreeName()`) and maps it to features (`SourceWorktree` priority in feature resolution), but there is no command to create or manage feature worktrees
- The `superpowers plan` (2026-04-06) envisioned worktree-based task isolation, but the management layer was never built

### Urgency

Without worktree management, developers cannot parallelize feature work across multiple Claude Code sessions. Each feature must complete before the next can start, creating a bottleneck when multiple independent features are in flight.

## Proposed Solution

Add a `forge worktree` subcommand group to the forge CLI that provides four operations:

- **`start <slug>`**: Creates a git worktree with branch name = slug in a sibling directory (`../<slug>`), then launches `claude` CLI in that directory. The forge-cli's existing `GetWorktreeName()` automatically detects the feature from the worktree name.
- **`list`**: Shows all forge-managed worktrees — name, branch, path, and feature status.
- **`remove <slug>`**: Removes the git worktree. Keeps the branch for manual merge.
- **`resume <slug>`**: Opens a new `claude` session in an existing worktree.

The user workflow: `forge worktree start auth-login` → Claude Code opens in `../auth-login` → user runs `/run-tasks` or other forge commands → feature is auto-detected from worktree name → `forge worktree remove auth-login` when done.

### Innovation Highlights

This is straightforward adoption of git worktree conventions tailored to the Forge feature workflow. The key insight is that forge-cli already has worktree → feature detection, so the management commands only need to handle creation/cleanup — the feature context is auto-resolved by existing code.

## Requirements Analysis

### Key Scenarios

1. **Happy path**: Developer runs `forge worktree start my-feature` → worktree created at `../my-feature` with branch `my-feature` from HEAD → `claude` launches → developer runs `/run-tasks` → feature auto-detected → tasks execute → developer runs `forge worktree remove my-feature` → worktree removed, branch kept
2. **Multiple features**: Developer runs `forge worktree start feature-a` in terminal 1, `forge worktree start feature-b` in terminal 2 → two independent Claude Code sessions in separate worktrees → both features developed simultaneously
3. **Resume**: Developer closes Claude session in a worktree → later runs `forge worktree resume my-feature` → re-opens Claude in the existing worktree with all state preserved
4. **Branch already exists**: `forge worktree start my-feature` where branch `my-feature` already exists → create worktree from existing branch (resume context)
5. **Worktree already exists**: `forge worktree start my-feature` where worktree already exists → error with hint to use `resume`
6. **Dirty working tree**: `forge worktree start` with uncommitted changes → warn user but proceed (worktree starts clean from HEAD)

### Non-Functional Requirements

- **Start latency**: Worktree creation + Claude launch should complete in < 5 seconds (git worktree creation is near-instant)
- **Cross-platform**: Must work on Windows, macOS, and Linux
- **Zero plugin changes**: This is purely a CLI feature — no plugin skills, agents, or hooks need modification

### Constraints & Dependencies

- Requires git worktree support (git 2.5+)
- Requires Claude Code CLI (`claude`) available in PATH
- Forge-cli's existing `GetWorktreeName()` and feature resolution must work correctly with the new worktree naming convention
- Worktree directory must be outside the main project tree (git worktree constraint)

## Alternatives & Industry Benchmarking

### Industry Solutions

Git worktree is a standard git feature supported by all major git hosting and tooling. The pattern of "one worktree per feature branch" is common in large monorepo workflows. Tools like `git worktree` itself, `mhutchie/git-graph`, and IDEs like VS Code provide worktree management UIs.

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | Zero cost | No parallel feature development; manual worktree setup is error-prone | Rejected: defeats the purpose |
| Document manual workflow | WORKFLOW.md | No code changes | Users must remember git commands; no list/cleanup management | Rejected: poor UX |
| Claude Code `--worktree` + hooks | Claude Code native | Leverages existing infrastructure | Hooks are shell scripts; no feature-specific management commands; less control | Rejected: doesn't provide start/list/remove/resume |
| **Forge CLI `worktree` subcommand** | Forge CLI | Full control, feature-aware, single command for create+launch | New Go code to maintain | **Selected: fills the exact gap with minimal complexity** |

## Feasibility Assessment

### Technical Feasibility

High. All building blocks exist:
- `git worktree add/remove/list` commands are stable and well-documented
- forge-cli already has `GetWorktreeName()`, `FindProjectRoot()`, and feature resolution
- Claude Code CLI can be launched as a subprocess
- Go's `os/exec` package handles subprocess management

### Resource & Timeline

Single developer, estimated 1-2 days for implementation + tests. Scope is well-bounded (4 subcommands, no plugin changes).

### Dependency Readiness

All dependencies are available: git 2.5+, Claude Code CLI, Go toolchain.

## Scope

### In Scope

- `forge worktree start <slug>` subcommand
- `forge worktree list` subcommand
- `forge worktree remove <slug>` subcommand
- `forge worktree resume <slug>` subcommand
- CLI tests for all four commands
- Documentation in `forge-cli/docs/WORKFLOW.md`

### Out of Scope

- Task-level parallelism within a single feature
- Plugin-level changes (hooks, skills, agents)
- Auto-merge on remove
- Worktree conflict resolution
- WorktreeCreate/WorktreeRemove hook integration
- Remote worktree management (push/pull)

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| Claude CLI not in PATH | M | H | Detect `claude` availability at start; print clear error with install instructions |
| `.worktreeinclude` not processed | M | L | Manual worktrees don't copy gitignored files (`.env` etc.). Document this; user handles manually. Not an issue for Forge features which don't typically need local env files |
| Windows path handling issues | M | M | Use `filepath.Join()` consistently; test on Windows from the start |
| Worktree name collision with existing directory | L | M | Check if target directory exists before `git worktree add`; fail with clear message |
| Orphaned worktrees after crash | L | L | `list` command shows worktree state; `remove` handles cleanup regardless |

## Success Criteria

- [ ] `forge worktree start <slug>` creates a git worktree and launches `claude` in it
- [ ] `forge feature` in a worktree auto-detects the correct feature via worktree name
- [ ] `forge worktree list` shows all forge-created worktrees with name, branch, and path
- [ ] `forge worktree remove <slug>` removes the worktree while preserving the branch
- [ ] `forge worktree resume <slug>` re-launches `claude` in an existing worktree
- [ ] All four commands have test coverage in the forge-cli test suite
