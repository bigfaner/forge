---
created: 2026-05-15
author: "faner"
status: Draft
---

# Proposal: `forge feature set` — Explicit Feature Selection

## Problem

`GetCurrentFeature()` relies entirely on implicit resolution (worktree name, branch name, task state), with no way for the user to explicitly declare which feature they're working on. The `feature` field in `.forge/state.json` is written by task claim/submit but never read for resolution.

### Evidence

- `.forge/state.json` contains a `feature` field written by `EnsureForgeState()` and `WriteForgeState()`, but `GetCurrentFeature()` never reads it
- Worktree names don't always match feature slugs — worktree names follow git conventions, feature slugs follow forge conventions
- Multiple features in a repo requires git context (branch/worktree) to disambiguate; no manual override exists

### Urgency

Without explicit feature selection, users working on branches with non-matching names or in worktrees with unrelated names must resort to workarounds. The data (`ForgeState.Feature`) already exists but is unused — the gap is minimal to close.

## Proposed Solution

Add `forge feature set <slug>` subcommand that writes the feature to `.forge/state.json` and ensures the feature directory exists. Adjust `GetCurrentFeature()` to check `.forge/state.json` as the highest-priority source.

### Innovation Highlights

Straightforward application of "explicit overrides implicit." The data structure already exists; this is wiring it into the resolution chain.

## Requirements Analysis

### Key Scenarios

- **Happy path**: User runs `forge feature set my-feature` → state.json written, directory ensured → all subsequent `forge task` commands resolve to `my-feature`
- **Worktree mismatch**: Worktree named `fix-auth-bug` but feature slug is `oauth-rewrite` → `forge feature set oauth-rewrite` overrides the implicit worktree-derived slug
- **Feature completion reset**: Quality-gate runs `ClearForgeState()` → deletes state.json → resolution falls back to git context (worktree name)
- **Verbose query**: `forge feature -v` shows `my-feature (from: state.json)` or `my-feature (from: worktree)`

### Non-Functional Requirements

- **Backward compatibility**: When `.forge/state.json` doesn't exist, behavior is identical to current
- **Performance**: No additional filesystem scans; state.json is a single file read

### Constraints & Dependencies

- `.forge/state.json` lifecycle is managed by hooks (cleanup, quality-gate); `set` integrates into existing lifecycle
- `ForgeState` struct unchanged — no schema migration needed

## Alternatives & Industry Benchmarking

### Industry Solutions

Most CLI tools with "current context" use an explicit config file (e.g., `kubectl config use-context`, `gcloud config set project`).

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | No code change | No explicit control; worktree mismatch unresolvable | Rejected: user pain exists |
| Environment variable | Unix convention | No file I/O | Lost between sessions; doesn't work with subagents | Rejected: less persistent than file |
| Separate config file | kubectl pattern | Clean separation | Yet another file to manage | Rejected: state.json already exists |
| **Reuse `.forge/state.json`** | Current codebase | Minimal change; data exists; per-worktree scope | Cleared on quality-gate (acceptable: feature complete = reset) | **Selected: lowest cost, highest leverage** |

## Feasibility Assessment

### Technical Feasibility

Fully supported by current Go codebase. Changes limited to:
- `feature.go`: New `SetFeatureState()` function, priority chain adjustment in `GetCurrentFeature()`
- `feature.go` (cmd): New `set` subcommand handler
- No new dependencies

### Resource & Timeline

Single-developer, estimated 2-3 tasks.

### Dependency Readiness

No external dependencies. All required APIs (`ReadForgeState`, `EnsureForgeState`, `EnsureFeatureDir`) exist.

## Scope

### In Scope

- `forge feature set <slug>` subcommand: ensure directory + write state.json
- `GetCurrentFeature()` priority chain: state.json > git > task state > single feature
- `forge feature -v`: show feature name + resolution source
- Reuse existing `ForgeState` struct

### Out of Scope

- `forge feature unset` command (set a different feature instead, or let cleanup handle it)
- `ClearForgeState()` behavior change (keep deleting entire file)
- Git integration (auto-create branch / switch worktree)
- `ForgeState` struct changes (no new fields)

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| `set` + `task claim` both write state.json | L | L | Same feature slug; `claim` overwrites with same value — no conflict |
| Priority change breaks existing scripts | L | M | state.json absent = identical behavior; fully backward compatible |
| `-v` flag conflicts with subcommands | L | L | Only applied to bare `forge feature`; `set`/`list`/`status` unaffected |

## Success Criteria

- [ ] `forge feature set <slug>` writes `.forge/state.json` with correct feature slug and `allCompleted=false`
- [ ] `forge feature set <slug>` creates feature directory if it doesn't exist
- [ ] `forge feature set <slug>` returns error when slug is empty
- [ ] `GetCurrentFeature()` returns state.json feature when present, overriding git context
- [ ] `forge feature -v` shows feature name and resolution source
- [ ] When state.json is deleted (quality-gate), resolution falls back to git context
- [ ] Existing tests pass without modification (backward compatibility)

## Next Steps

- Proceed to `/quick-tasks` to generate implementation tasks
