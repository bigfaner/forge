---
id: "1"
title: "Implement forge feature complete --if-done CLI command"
priority: "P0"
estimated_time: "2h"
dependencies: []
scope: "backend"
breaking: false
type: "feature"
mainSession: false
---

# 1: Implement forge feature complete --if-done CLI command

## Description

Implement the `forge feature complete` subcommand with `--if-done` flag in the forge CLI. This command is the core of the stop-hook-completion feature: it runs as a second Stop hook after `forge quality-gate` and handles the post-completion status transition + git commit + optional push.

The command reads the active feature's `index.json` to determine if all tasks are completed/skipped. If so, it updates `manifest.md` status to `completed`, optionally updates `proposal.md` status to `Completed` (quick mode detection via file existence), commits both files, and pushes if `auto.gitPush` is enabled.

## Reference Files

- `docs/proposals/stop-hook-completion/proposal.md` — Source proposal
- `forge-cli/internal/cmd/feature.go` — Existing feature command group (add `complete` subcommand here)
- `forge-cli/internal/cmd/quality_gate.go` — Reference for hook protocol (stdin JSON, PrintHookJSON)
- `forge-cli/pkg/feature/forge_state.go` — ForgeState reading, GetCurrentFeature
- `forge-cli/pkg/feature/constants.go` — Status constants, path constants
- `forge-cli/pkg/task/` — TaskIndex reading
- `plugins/forge/hooks/hooks.json` — Hook registration (add second Stop hook)
- `forge-cli/internal/cmd/submit.go` — Reference for index.json reading pattern

## Acceptance Criteria

- [ ] `forge feature complete --if-done` exits 0 silently when no feature is active (no index.json found)
- [ ] `forge feature complete --if-done` exits 0 silently when index.json has pending/in_progress tasks
- [ ] When all tasks in index.json are completed/skipped, command updates manifest.md status to `completed`
- [ ] When proposal.md exists in feature directory (quick mode), command also updates its status to `Completed`
- [ ] When proposal.md is absent (full pipeline), only manifest.md is updated
- [ ] Git commit targets only manifest.md and proposal.md by explicit path — never `git add -A` or `git add .`
- [ ] When `auto.gitPush` is `true` in config, command pushes to remote after commit
- [ ] When `auto.gitPush` is `false` or absent, no push occurs
- [ ] Push failure is logged to stderr but does not block (exit 0)
- [ ] hooks.json Stop array contains `forge quality-gate` first, then `forge feature complete --if-done` second
- [ ] Go tests pass with >= 80% coverage on new code
- [ ] `forge task validate-index` still passes after changes

## Hard Rules

- Do NOT depend on `.forge/state.json` — it is consumed and cleared by quality-gate before this command runs
- `git add` must target only specific file paths, never broad staging
- Exit code is always 0 (hook protocol: non-blocking). Communicate via stderr for logging
- Follow TDD: write tests first, then implement
- Bump version in `scripts/version.txt` (minor: new command)

## Implementation Notes

1. **Follow existing patterns**: The quality-gate command (`quality_gate.go`) shows how to read stdin JSON from hooks, use `PrintHookJSON`, and follow the hook protocol. The `feature.go` file shows the cobra command registration pattern.

2. **Hook stdin**: The Stop hook receives JSON on stdin with `stop_hook_active` and `last_assistant_message`. The `--if-done` command does NOT need to parse stdin JSON — it only checks index.json.

3. **Pipeline mode detection**: Check `os.Stat(<featureDir>/proposal.md)`. If exists → quick mode → update both files. If absent → full pipeline → update manifest.md only.

4. **Config reading**: Use the existing `forge config get auto.gitPush` pattern from run-tasks.md, or read directly from `.forge/config.yaml`.

5. **Git operations**: Use `os/exec` for git commands. Stage both files explicitly: `git add manifest.md proposal.md`. Commit with a descriptive message. Push only if configured.

6. **Atomicity**: Write both status files to disk first, then stage + commit in one operation. If file writes fail, no git commands execute.

7. **Version bump**: New CLI command = minor version bump in `forge-cli/scripts/version.txt`.

8. **Pre-condition verification**: Before implementing, verify that Claude Code executes Stop hooks sequentially by testing with two hooks in hooks.json where the first returns block.
