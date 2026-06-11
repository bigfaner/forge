---
created: "2026-05-18"
tags: [interface, local-dev-deployment]
---

# forge feature has no `get` subcommand

## Problem

Ran `forge feature get` expecting to retrieve the current active feature. The command returned unexpected output instead of the current feature context.

## Root Cause

1. Assumed `get` was a standard CRUD subcommand pattern for `forge feature`
2. Forge CLI uses a simpler convention: `forge feature` (no arguments) displays the current feature, `forge feature set <slug>` explicitly sets one
3. The help text clearly states "Without arguments: displays the current feature" but this was not consulted before running the command

## Solution

Use `forge feature` (no subcommand, no arguments) to display the current active feature. Use `forge feature set <slug>` to set one.

## Reusable Pattern

Before invoking any CLI command, check its actual interface with `-h` if unfamiliar. Don't assume RESTful CRUD verb patterns (get/set/list) — each CLI has its own conventions.

## Example

```bash
# Wrong
forge feature get

# Correct - display current feature
forge feature

# Correct - set feature
forge feature set my-feature-slug

# Correct - list all features
forge feature list
```

## References

- `forge feature -h` output
