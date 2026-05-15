# Proposal: Migrate Claude Commands to forge-cli

## Problem

Claude-related shortcuts (`just claude`, `just claude-c`, `just claude-w`) live in the project justfile, making them invisible to `forge --help` and disconnected from the forge CLI ecosystem. Users must remember justfile recipe names instead of using the natural `forge claude` entry point. The recipes also provide no access to Claude CLI's native flags (e.g. `--model`, `--max-turns`, `--allowedTools`).

## Solution

Add a `forge claude` subcommand that:
1. **Always injects** `--dangerously-skip-permissions` (the core value prop)
2. **Passes through** all user args directly to the Claude CLI binary (no `--` separator, no explicit flag mapping)
3. Replaces 3 justfile recipes (`claude`, `claude-c`, `claude-w`) with one unified command

Usage examples:
```bash
forge claude                    # → claude --dangerously-skip-permissions
forge claude -c                 # → claude --dangerously-skip-permissions -c
forge claude -w my-feature      # → claude --dangerously-skip-permissions -w my-feature
forge claude --model opus -p "hello"  # → claude --dangerously-skip-permissions --model opus -p "hello"
```

## Alternatives Considered

| Approach | Pros | Cons |
|----------|------|------|
| **Migrate to forge-cli** (chosen) | Discoverable, single command, full flag access, consistent with forge ecosystem | Slightly more code in forge-cli |
| Do nothing (keep justfile) | No implementation work | Undiscoverable, no flag support, fragmented entry points |
| Shell aliases | Quick to set up | Per-shell config, not portable, not discoverable |

## Scope

### In Scope
- Add `forge claude` subcommand to forge-cli with arg passthrough
- Always inject `--dangerously-skip-permissions`
- Validate `claude` binary exists in PATH before execution
- Remove `claude`, `claude-c`, `claude-w` recipes from project justfile
- Update `forge init` to stop generating those justfile recipes

### Out of Scope
- `claude-p` recipe (plugin-dir specific, stays in justfile)
- Any forge-specific flag parsing or transformation
- Configurable permission injection (always-on by design)

## Risks

| Risk | Mitigation |
|------|------------|
| Claude binary not in PATH | Pre-flight check with clear error: "`claude` not found in PATH — install Claude Code first" |
| Cobra flag parsing conflicts with passthrough args | Use `DisableFlagParsing: true` on the command to skip cobra's flag handling entirely |
| Existing projects have justfile recipes that users rely on | `forge init` generates them; migration is additive — users can transition at their own pace |

## Success Criteria

- [ ] `forge claude` launches Claude CLI with `--dangerously-skip-permissions`
- [ ] `forge claude -c` continues the last conversation
- [ ] `forge claude -w <name>` opens a worktree session
- [ ] Any Claude CLI flag passes through: `forge claude --model opus -p "prompt"`
- [ ] Clear error when `claude` binary is not in PATH
- [ ] `claude`, `claude-c`, `claude-w` recipes removed from project justfile
- [ ] `forge init` no longer generates those justfile recipes
