---
created: "2026-05-20"
tags: [local-dev-deployment, testing]
---

# gotcha-quick-tasks-stale-detect-command

## Problem

The `quick-tasks` skill Step 0 instructs to run `forge test detect` to auto-detect the project's test language(s). This command does not exist in forge CLI — running it produces `Error: unknown command "detect" for "forge test"`.

## Root Cause

1. The `quick-tasks` skill template references `forge test detect` as the language detection mechanism.
2. This command was likely planned but never implemented, or was removed in a refactor.
3. The skill template was not updated to reflect the actual CLI surface, creating a dead reference that blocks the quick-tasks pipeline at Step 0.

The actual language detection works via `.forge/config.yaml` `languages` field — if configured, `forge task index` uses it directly without needing a `detect` subcommand.

## Solution

Skip `forge test detect` in Step 0. Instead:
1. Read `.forge/config.yaml` and check for the `languages` field.
2. If `languages` is present (e.g., `languages: [go]`), use it as the resolved language.
3. If `languages` is missing, pass `--languages` flag to `forge task index` explicitly.

## Reusable Pattern

When a skill template references a CLI command that fails with "unknown command":
1. Check the actual CLI surface (`forge <command> -h`).
2. Look for alternative config-based mechanisms (config.yaml, environment variables).
3. Update the skill template if the command was removed or never existed.
4. Never block the pipeline on a missing subcommand — fall back to config or explicit input.

## Related Files

- `skills/quick-tasks/SKILL.md` — references `forge test detect` in Step 0
- `.forge/config.yaml` — actual language configuration (`languages: [go]`)
