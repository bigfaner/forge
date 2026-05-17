---
created: "2026-05-17"
tags: [local-dev-deployment]
---

# Use `forge` not `go run` in Skill/Agent Execution Flows

## Problem

During task execution in `/run-tasks`, all forge CLI invocations used `cd forge-cli && go run ./cmd/forge task <subcommand>` instead of the installed `forge task <subcommand>` binary.

## Root Cause

Causal chain:
1. **Symptom**: Bash commands used `go run ./cmd/forge` unnecessarily
2. **Direct cause**: Development habit — when writing/modifying forge code, `go run` is the natural workflow
3. **Root cause**: Failed to distinguish two contexts: (a) developing forge itself vs (b) using forge to execute tasks. In context (b), forge is a pre-installed CLI at `~/.forge/bin/forge`
4. **Trigger**: Any Bash call to forge CLI in skill/agent prompts defaults to `go run` without checking if the binary is installed

## Solution

Replace all `go run ./cmd/forge` invocations in skill/agent execution flows with `forge` (the installed binary).

## Reusable Pattern

**Rule**: When executing forge commands during skill/agent workflows, use `forge <subcommand>` directly. The binary is installed at `~/.forge/bin/forge` (confirmed in PATH). Only use `go run ./cmd/forge` when actively developing forge itself (iterating on CLI code between runs).

## Example

```bash
# Wrong (skill/agent execution context):
cd Z:/project/ai/forge/forge-cli && go run ./cmd/forge task claim
# compilation overhead, development-only pattern

# Right (skill/agent execution context):
forge task claim
# uses installed binary, fast, production usage
```

## Related Files

- `~/.forge/bin/forge` — installed forge binary
- `forge-cli/` — forge CLI source (development only)
