---
created: 2026-05-17
author: faner
status: Approved
---

# Proposal: Stop Hook Auto-Completion

## Problem

After all tasks complete in `/quick` or full pipeline mode, status updates are written to `manifest.md` and `proposal.md` in the working tree but never committed to git — the user finds modified but uncommitted status files after the pipeline finishes.

### Evidence

**Reproduction**: Run `/quick` on any feature with >= 1 task. After all tasks complete and quality-gate passes, `git status` shows `manifest.md` and `proposal.md` modified but uncommitted. Occurs on 100% of successful runs.

**Root cause** (`docs/lessons/gotcha-post-completion-commit.md`): `/run-tasks` Post-Completion prints a summary and optionally pushes — no commit step. Neither `/quick` nor `/run-tasks` owns the post-completion commit.

**Frequency**: All 10 features completed through `/quick` required manual status commit cleanup.

### Urgency

Every `/quick` run produces uncommitted files. Manual `git add && git commit` is the workaround — friction that defeats autonomous pipeline execution and risks stale status if the user forgets.

## Proposed Solution

Add `forge feature complete --if-done` as a second Stop hook command alongside the existing `forge quality-gate`. After quality-gate passes and all tasks are confirmed done, the command:

1. Updates `manifest.md` status → `completed`
2. Detects pipeline mode by checking whether `<feature-dir>/proposal.md` exists on disk (exists = quick mode; absent = full pipeline). In quick mode only, updates `proposal.md` status → `Completed`
3. Commits both files
4. Pushes to remote if `auto.gitPush` is enabled in config

### User-Facing Behavior

**Success**: Prints `[feature:complete] Status committed: manifest.md, proposal.md` to stderr. If auto-push enabled, adds `[feature:complete] Pushed to remote`.

**Skip**: No output — silent exit when conditions not met.

**Failure**: Prints `[feature:complete] Error: <message>` to stderr, exits non-zero. Non-blocking: agent stop proceeds regardless.

### Innovation Highlights

The design adapts the CI/CD stage pattern to a domain-specific constraint: Claude Code's Stop hooks are the only mechanism to inject logic between "agent finishes responding" and "session ends." Unlike CI/CD where stages are first-class objects, here stages are implicit — achieved by ordering two independent hook commands in a JSON array. The domain-specific adaptation is the `--if-done` flag guard: the command reads `index.json` directly and exits silently when preconditions fail, making it safe to fire on every Stop event. This eliminates inter-hook coordination — each hook independently decides whether to act.

**Data source**: `forge feature complete --if-done` reads `<feature-dir>/index.json` to determine task completion status. It does not depend on `.forge/state.json` (consumed and cleared by quality-gate before this command runs).

## Requirements Analysis

### Key Scenarios

- **Happy path (quick)**: All tasks done → quality-gate passes → `forge feature complete --if-done` commits + pushes → agent stops
- **Happy path (full pipeline)**: Same flow, but no proposal.md to update
- **Quality-gate adds fix tasks**: quality-gate blocks → command sees pending tasks → skips → agent continues → fix tasks done → quality-gate passes → command commits
- **No active feature**: command finds no `index.json` → exits silently in <1s
- **Feature already completed**: command checks index.json → all done → idempotent, skips if already committed
- **Commit fails mid-way**: Both status updates are written to disk first, then staged and committed in a single `git add + git commit`. If file writes fail, no git commands execute. If commit fails, no push is attempted. Error logged to stderr; hook exits non-zero without blocking agent stop
- **index.json missing or corrupt**: Validates structure before reading. On parse failure or missing file, exits silently
- **Multiple active features**: Reads active feature from `.forge/config.yaml` (singleton). At most one feature per session
- **User interrupt during hook**: Ctrl+C terminates the hook process. Working tree may have unstaged updates but no partial commit — single-commit strategy prevents half-committed state

### Non-Functional Requirements

- **Latency**: Hook must complete in <2s when skipping (no file edits, no git operations)
- **Atomicity**: Both files written first, then staged and committed in a single `git add + git commit`. No intermediate commit. If any step fails before commit, no commit is created
- **Security**: `git add` targets only `manifest.md` and `proposal.md` by path — no `git add -A` or `git add .`. Unstaged user changes in other files are never included in the commit
- **Backward-compatibility**: The hook is added to `hooks.json` shipped with the forge plugin. Users upgrading to the new version receive the hook automatically. Existing users with customized `hooks.json` must merge manually (documented in release notes). The command is a no-op when no feature is active, so the upgrade cannot break existing workflows
- **Cross-platform**: Works on Windows (Git Bash) and Unix. File operations use Node.js `fs`; git commands via `child_process.exec` with platform-agnostic arguments

### Constraints & Dependencies

- Stop hooks receive `stop_hook_active` in stdin JSON (Claude Code provides this)
- `forge task submit` sets `allCompleted: true` in `.forge/state.json` after each submission
- quality-gate consumes (clears) `.forge/state.json` before running gates
- `forge feature complete --if-done` must NOT depend on `.forge/state.json` (already consumed by quality-gate)

## Alternatives & Industry Benchmarking

### Industry Solutions

Three established patterns handle post-completion actions:

1. **GitHub Actions `needs:` + job stages**: Downstream jobs declare `needs:` to depend on upstream success. The deploy job runs only after test passes. Directly analogous: `forge feature complete --if-done` depends on quality-gate passing.
2. **ArgoCD PostSync Hooks**: `PostSync` hooks run after successful sync, separate from sync logic. Idempotent and fails non-blocking. Maps to: separate hook acting after quality-gate succeeds.
3. **Jenkins `post` blocks**: `post { success { ... } }` runs cleanup/deploy after main stage, with separate `failure` and `always` blocks. Separates lifecycle management from stage logic.

Common principle: **separate the gate from the action**, make the action idempotent.

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | Zero effort | Uncommitted files every run | Rejected: defeats automation |
| Combined into quality-gate | Single command | One config point | Couples quality validation with lifecycle management — a status formatting bug could break quality checks. Violates SRP | Rejected: error isolation failure |
| Inline commit in `/quick` or `/run-tasks` | Skill code | No new hook | Status commit runs before quality-gate (Stop hook runs after skill returns). Premature if quality-gate adds fix tasks | Rejected: ordering violation |
| **Two separate Stop hooks** | CI/CD stage pattern (GitHub Actions `needs:`, ArgoCD PostSync) | Clean SoC, idempotent, independently testable | Two hooks to maintain; assumes sequential array-order execution | **Selected: best SoC + safety** |

## Feasibility Assessment

### Technical Feasibility

The `feature` CLI command group already exists with `list`, `status`, `set` subcommands. Adding `complete` with `--if-done` flag follows the existing pattern. The hooks.json format already supports multiple commands in the `hooks` array.

**Pre-condition**: This design requires Claude Code to execute Stop hooks sequentially in array order. This must be verified before implementation by testing: two Stop hooks in hooks.json, first returns block, confirm second executes after first completes. If sequential execution is not guaranteed, the fallback is a single combined hook that performs both quality-gate and completion logic.

### Resource & Timeline

~3-5 tasks. Straightforward CLI implementation + hook config + skill doc updates.

### Dependency Readiness

All dependencies exist: CLI infrastructure, hooks.json format, config system (`auto.gitPush`), git operations via existing helpers.

## Scope

### In Scope

- `forge feature complete --if-done` CLI command
- Update `plugins/forge/hooks/hooks.json` to add second Stop hook
- Remove post-completion status transition from `plugins/forge/commands/quick.md`
- Move auto-push from `plugins/forge/commands/run-tasks.md` to the new hook
- The command auto-detects pipeline mode by checking for `proposal.md` existence in the feature directory — updates proposal.md status only when the file exists (quick mode), and skips it in full pipeline mode

### Out of Scope

- Changes to quality-gate behavior
- Changes to cleanup hook (SessionEnd/SubagentStop)
- New config options (reuses existing `auto.gitPush`)
- E2E test execution (handled by quality-gate)

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| Sequential hook execution assumption wrong | M | H | Test before implementation: two Stop hooks, first returns block, confirm second skips. Fallback: single combined hook |
| Command runs during fix-task loop | M | M | Checks index.json for pending tasks — early exit with no file mutations. Fallback: if index.json is stale (tasks added but not yet reflected), the commit would set status to completed prematurely — guarded by also checking that all task files have `status: completed` in their frontmatter, not relying solely on index.json summary |
| Regression from removing status transition in `/quick` | M | M | Remove transition only after command verified in both modes. E2E test before removal |
| Double-commit on edge case | L | L | Idempotent — duplicate commit produces empty commit which git rejects |
| Git push fails (network, auth) | M | M | Commit succeeds locally. Push failure logged. Non-blocking |
| Git merge conflict on auto-push | L | M | Detects non-zero `git push` exit. No force-push. User resolves manually |

## Success Criteria

- [ ] After `/quick` completes all tasks + quality-gate passes, `manifest.md` and `proposal.md` show status `completed` in a single git commit (verify: `git log -1 --name-only` lists both files)
- [ ] After full pipeline completes, `manifest.md` shows status `completed` in a git commit
- [ ] When quality-gate adds fix tasks, command creates zero commits between quality-gate rounds (verify: intentional quality-gate failure, check `git log` for no intermediate commits)
- [ ] When no feature is active, exits with code 0 in <1s, no stdout
- [ ] When `auto.gitPush: true`, `git push` succeeds and status files appear on remote branch
- [ ] `quick.md` contains no status transition logic after refactoring (verify: `grep -c "status.*completed" quick.md` returns 0)
- [ ] `run-tasks.md` contains no auto-push logic after refactoring (verify: `grep -c "gitPush\|git push" run-tasks.md` returns 0)
- [ ] Atomic: after simulated failure (corrupt index.json), `git log` shows either a complete two-file commit or no commit — no single-file partial commit
- [ ] Cross-platform: runs on Windows (Git Bash) and Linux without platform-specific code paths
- [ ] `hooks.json` Stop hooks array contains `forge quality-gate` before `forge feature complete --if-done` in that order (verify: `jq '.Stop.hooks' hooks.json` shows both entries with quality-gate first)
- [ ] Mode detection: running `forge feature complete --if-done` in a feature directory with `proposal.md` present updates both files; running in a feature directory without `proposal.md` updates only `manifest.md` — same command, no flags or configuration difference

## Next Steps

- Proceed directly to `/quick-tasks` (no PRD/design needed — well-scoped infrastructure change)
