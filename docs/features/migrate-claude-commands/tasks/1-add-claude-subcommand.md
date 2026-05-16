---
id: "1"
title: "Add forge claude subcommand with arg passthrough"
priority: "P0"
estimated_time: "1-2h"
dependencies: []
scope: "backend"
breaking: false
type: "implementation"
mainSession: false
---

# 1: Add forge claude subcommand with arg passthrough

## Description

Claude-related shortcuts live in the project justfile, making them invisible to `forge --help`. Add a `forge claude` subcommand that always injects `--dangerously-skip-permissions` and passes through all user args directly to the Claude CLI binary.

## Reference Files
- `docs/proposals/migrate-claude-commands/proposal.md` — Source proposal
- `forge-cli/internal/cmd/root.go` — Command registration
- `forge-cli/internal/cmd/init.go` — Existing claude-related code (justfile recipes)

## Acceptance Criteria

- [ ] `forge claude` launches Claude CLI with `--dangerously-skip-permissions`
- [ ] `forge claude -c` continues the last conversation
- [ ] `forge claude -w <name>` opens a worktree session
- [ ] Any Claude CLI flag passes through: `forge claude --model opus -p "prompt"`
- [ ] Clear error when `claude` binary is not in PATH
- [ ] Unit tests for: PATH validation, arg passthrough, flag injection

## Hard Rules

- Use `DisableFlagParsing: true` on the cobra command to avoid flag parsing conflicts
- No `--` separator required; all user args pass through transparently
- `--dangerously-skip-permissions` is always prepended, not configurable

## Implementation Notes

- Pre-flight check: verify `claude` binary exists in PATH before execution
- Use `exec.LookPath("claude")` for PATH validation
- Use `os/exec.CommandContext` with args: `["--dangerously-skip-permissions", ...userArgs]`
- Register in `root.go` alongside other subcommands
