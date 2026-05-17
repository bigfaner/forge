---
created: 2026-05-17
author: faner
status: Draft
---

# Proposal: Stop Hook Auto-Completion

## Problem

After all tasks complete in `/quick` or full pipeline mode, `manifest.md` and `proposal.md` status updates are never committed — the user finds uncommitted status files after the pipeline finishes.

### Evidence

Documented in `docs/lessons/gotcha-post-completion-commit.md`. The gap exists between `/quick` and `/run-tasks` responsibilities: neither skill reliably owns the post-completion status transition + commit when quality-gate adds fix tasks.

### Urgency

Every `/quick` execution that completes all tasks produces uncommitted files. The workaround is manual git commit — friction that defeats the purpose of autonomous pipeline execution.

## Proposed Solution

Add `forge feature complete-if-done` as a second Stop hook command alongside the existing `forge quality-gate`. After quality-gate passes and all tasks are confirmed done, the command:

1. Updates `manifest.md` status → `completed`
2. Updates `proposal.md` status → `Completed` (if exists, quick mode only)
3. Commits both files
4. Pushes to remote if `auto.gitPush` is enabled in config

### Innovation Highlights

Leverages Claude Code's Stop hook mechanism for lifecycle management. The key insight is that two sequential Stop hooks create a natural pipeline: gate-then-commit. The second hook (`complete-if-done`) is safe to run on every Stop event because it simply checks index.json and skips when conditions aren't met. This is a straightforward application of the existing hook infrastructure — no new patterns needed.

## Requirements Analysis

### Key Scenarios

- **Happy path (quick)**: All tasks done → quality-gate passes → complete-if-done commits + pushes → agent stops
- **Happy path (full pipeline)**: Same flow, but no proposal.md to update
- **Quality-gate adds fix tasks**: quality-gate blocks → complete-if-done sees pending tasks → skips → agent continues → fix tasks done → quality-gate passes → complete-if-done commits
- **No active feature**: complete-if-done finds no state → exits silently
- **Feature already completed**: complete-if-done checks index.json → all done → idempotent, skips if already committed

### Non-Functional Requirements

- **Latency**: Hook must complete in <2s when skipping (no file edits, no git operations)
- **Reliability**: Status transition + commit must be atomic — either both files update or neither
- **Cross-platform**: Must work on Windows (bash + cmd.exe) and Unix

### Constraints & Dependencies

- Stop hooks receive `stop_hook_active` in stdin JSON (Claude Code provides this)
- `forge task submit` sets `allCompleted: true` in `.forge/state.json` after each submission
- quality-gate consumes (clears) `.forge/state.json` before running gates
- `complete-if-done` must NOT depend on `.forge/state.json` (already consumed by quality-gate)

## Alternatives & Industry Benchmarking

### Industry Solutions

CI/CD pipelines (GitHub Actions, GitLab CI) handle post-completion actions through separate job stages that run after the main job succeeds. This is the same pattern: gate → deploy.

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | Zero effort | Uncommitted files every run, manual cleanup | Rejected: defeats automation goal |
| Combined into quality-gate | Single command | One config point | Couples gate logic with lifecycle management | Rejected: violates single responsibility |
| **Two separate Stop hooks** | CI/CD stage pattern | Clean separation, idempotent, testable | Two hook commands to maintain | **Selected: best SoC + safety** |

## Feasibility Assessment

### Technical Feasibility

The `feature` CLI command group already exists with `list`, `status`, `set` subcommands. Adding `complete-if-done` follows the existing pattern. The hooks.json format already supports multiple commands in the `hooks` array.

### Resource & Timeline

~3-5 tasks. Straightforward CLI implementation + hook config + skill doc updates.

### Dependency Readiness

All dependencies exist: CLI infrastructure, hooks.json format, config system (`auto.gitPush`), git operations via existing helpers.

## Scope

### In Scope

- `forge feature complete-if-done` CLI command
- Update `plugins/forge/hooks/hooks.json` to add second Stop hook
- Remove post-completion status transition from `plugins/forge/commands/quick.md`
- Move auto-push from `plugins/forge/commands/run-tasks.md` to the new hook
- Support both quick and full pipeline status flows

### Out of Scope

- Changes to quality-gate behavior
- Changes to cleanup hook (SessionEnd/SubagentStop)
- New config options (reuses existing `auto.gitPush`)
- E2E test execution (handled by quality-gate)

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| Hook runs before quality-gate finishes | L | L | Hooks execute sequentially in array order |
| Double-commit on edge case | L | L | Idempotent — git add on already-committed files is a no-op |
| complete-if-done runs during fix-task loop | M | M | Checks index.json for all tasks completed — pending fix tasks prevent action |
| Git push fails (network, auth) | M | M | Commit still succeeds locally. Push failure is logged but not blocking |

## Success Criteria

- [ ] After `/quick` completes all tasks + quality-gate passes, manifest.md and proposal.md show status `completed` in a git commit
- [ ] After full pipeline completes, manifest.md shows status `completed` in a git commit
- [ ] When quality-gate adds fix tasks, no premature status commit occurs
- [ ] When no feature is active, `complete-if-done` exits silently in <1s
- [ ] Auto-push works when `auto.gitPush: true` is set in config

## Next Steps

- Proceed directly to `/quick-tasks` (no PRD/design needed — well-scoped infrastructure change)
