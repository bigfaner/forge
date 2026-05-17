---
created: "2026-05-17"
tags: [testing]
---

# Skill Docs Can Reference Nonexistent CLI Flags

## Problem

Executing `forge task index --feature <slug> --languages go` fails with `Error: unknown flag: --languages`. The quick-tasks skill doc (`SKILL.md` Step 5) instructs to use `--languages` when `forge testing detect` returns no result, but the CLI only supports `--test-profiles`.

## Root Cause

Causal chain:
1. **Symptom**: `forge task index --languages go` → "unknown flag: --languages"
2. **Direct cause**: Agent followed the quick-tasks skill doc which documents a `--languages` flag that doesn't exist in the CLI
3. **Root cause**: Skill doc is stale — the CLI was refactored to use `--test-profiles` instead of `--languages`, but the skill doc was never updated
4. **Contributing factor**: `forge testing detect` returned no language because `go.mod` lives in `forge-cli/` subdirectory, not project root — triggering the "pass language explicitly" path
5. **Agent mistake**: Did not verify CLI interface via `forge task index --help` before executing

## Solution

Ran `forge task index --feature stop-hook-completion` without `--languages`. The CLI handled language detection internally and succeeded.

## Reusable Pattern

**Rule**: Before executing any CLI command documented in a skill file, verify the actual flags with `--help`. Skill docs can become stale when the CLI evolves — a quick help check prevents errors.

**When to apply**: Any time a skill doc references specific CLI flags, especially for commands that are actively being developed. The `--help` output is the source of truth, not the skill doc.

**How to apply**: Add `forge <command> --help` as a verification step before first use of any unfamiliar flag combination.

## Example

```bash
# Bad: blindly follow skill doc
forge task index --feature my-feature --languages go  # fails if --languages doesn't exist

# Good: verify first
forge task index --help  # check actual flags
# Then use the correct interface
forge task index --feature my-feature
```

## Related Files

- `plugins/forge/skills/quick-tasks/SKILL.md` — Contains the stale `--languages` reference (Step 5)
- `forge-cli/internal/cmd/index.go` — Actual CLI implementation with `--test-profiles` flag

## References

- `docs/lessons/gotcha-post-completion-commit.md` — Related lesson about docs/code sync
