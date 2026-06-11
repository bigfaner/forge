---
created: "2026-05-24"
status: Draft
intent: "Enhance forge task list with slug parameter and worktree-aware reading"
---

# Proposal: forge task list Cross-Worktree Enhancement

## Problem

`forge task list` currently has no way to specify which feature to query — it relies entirely on auto-detection (`.forge/state.json` → worktree name → branch name → features dir scan). When a feature's tasks are being worked on in a worktree, the main repo's `index.json` copy may be stale. Users must `cd` into the worktree to see latest task status.

### Evidence

- Running `forge task list` from main repo with no active feature shows error or wrong feature
- Worktree has its own copy of `docs/features/<slug>/tasks/index.json` with up-to-date statuses
- No existing command can query a specific feature's tasks from an arbitrary working directory

### Urgency

Working on v3.0.0 with multiple parallel features — frequently need to check progress of features running in worktrees without context switching.

## Proposed Solution

Enhance `forge task list` with an optional positional slug parameter and worktree-aware index reading:

```bash
# Current behavior (unchanged)
forge task list

# New: specify feature directly
forge task list my-feature

# New: force read from main repo (ignore worktree)
forge task list my-feature --local
```

### Behavior

1. **No slug**: existing auto-detection logic (unchanged)
2. **With slug**: bypass auto-detection, directly locate `docs/features/<slug>/tasks/index.json`
3. **Reading priority** (when slug specified):
   - Default: check if a worktree exists for the slug → read from worktree's `index.json`
   - `--local` flag: always read from main repo's `index.json`, ignore worktree
4. **Worktree detection**: use existing `git worktree list` parsing or check `.forge/worktrees/<slug>` path

## Requirements Analysis

### Key Scenarios

- **Happy path**: `forge task list my-feature` from main repo → shows latest tasks from worktree
- **No worktree**: `forge task list my-feature` → falls back to main repo's index.json
- **Force local**: `forge task list my-feature --local` → always reads from main repo
- **Feature not found**: slug doesn't match any `docs/features/<slug>/` → clear error
- **No index.json**: feature exists but no tasks generated → "no tasks found" message

### Non-Functional Requirements

- Zero performance regression when slug is not provided (existing path)
- Worktree detection should be fast (<100ms, single git command or path check)

### Constraints & Dependencies

- Depends on existing `git` package for worktree detection
- Must not break existing `forge task list` behavior (backward compatible)

## Alternatives & Industry Benchmarking

| Approach | Pros | Cons | Verdict |
|----------|------|------|---------|
| **Do nothing** — cd into worktree | No code changes | Poor DX, context switching | **Rejected**: defeats purpose of CLI |
| `forge worktree status <slug>` — extend worktree cmd | Reuses existing worktree awareness | Wrong command family; task listing belongs under `forge task` | **Rejected**: breaks convention |
| **Chosen: enhance task list** | Minimal change, intuitive API, backward compatible | None significant | **Selected** |

## Feasibility Assessment

### Technical Feasibility

Fully supported by existing Go CLI architecture. Key functions already exist: `ListWorktrees()` in `pkg/git/`, `FindProjectRoot()`, `LoadIndex()`.

### Resource & Timeline

Single coding task — modify `list.go` to accept optional arg and add `--local` flag.

## Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| Worktree detection requires new code | Codebase-first check | Confirmed: `git.ListWorktrees()` and `git.GetWorktreeName()` already exist |
| Need separate command for cross-feature queries | Occam's Razor | Overturned: optional positional arg on existing command is simpler |

## Scope

### In Scope

- Modify `list.go`: accept optional positional slug arg
- Add `--local` flag to force main repo reading
- Worktree detection: check if worktree exists for given slug
- Read index.json from worktree path when available
- Update help text and examples

### Out of Scope

- Modifying task statuses from outside worktree
- Changes to other task commands (claim, submit, query, etc.)
- New filtering flags (`--status`, `--type`, etc.)
- Cross-feature aggregation

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| Worktree path resolution differs by OS | L | L | Use existing `git.ListWorktrees()` which handles OS differences |
| Stale worktree index.json (worktree abandoned but not removed) | M | L | Accept as-is; `--local` flag provides escape hatch |
| Breaking existing `forge task list` behavior | L | H | Arg is optional; existing path unchanged; add tests |

## Success Criteria

- [ ] `forge task list my-feature` shows tasks for specified feature from worktree (if active)
- [ ] `forge task list my-feature --local` reads from main repo regardless of worktree
- [ ] `forge task list` (no args) behaves exactly as before
- [ ] Clear error when slug doesn't match any feature
- [ ] Unit tests cover: slug resolution, worktree detection, --local override, backward compatibility

## Next Steps

- Proceed to `/quick-tasks` to generate implementation tasks
