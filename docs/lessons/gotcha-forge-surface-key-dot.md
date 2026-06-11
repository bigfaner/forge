---
created: "2026-06-03"
tags: [local-dev-deployment, testing]
---

# Forge Surface Key "." Blocks Surface-Aware Recipe Generation

## Problem

Running `/forge:init-justfile` skipped all surface-aware recipe generation (e.g., `tui-test`, `tui-teardown`). The error reported was:

```
Error: invalid surface-key ".": must match [a-zA-Z0-9_-]+
```

The justfile was generated with only language-level targets, no surface recipes.

## Root Cause

1. `.forge/config.yaml` used the shorthand form `surfaces: tui` (a bare type string).
2. `forge surfaces --json` parsed this as `{"surfaces":[{"key":".","type":"tui"}]}` — defaulting the key to `"."` (current directory) since no explicit key was provided.
3. The `/forge:init-justfile` skill validates surface keys against `[a-zA-Z0-9_-]+` and aborts surface recipe generation for keys that don't match. The dot character fails this check.

## Solution

Use the map/list form in `.forge/config.yaml` with an explicit valid key:

```yaml
surfaces:
  - key: agent-forensic
    type: tui
```

Then re-run `/forge:init-justfile` to generate surface-aware recipes (`agent-forensic-test`, `agent-forensic-teardown` for a multi-surface project, or `tui-test`/`tui-teardown` if it's the only surface).

## Reusable Pattern

When configuring surfaces in `.forge/config.yaml`, always use the explicit map form with a named key. The bare string shorthand (`surfaces: tui`) produces an invalid default key that blocks surface-aware skill features (justfile generation, test orchestration).

## Related Files

- `.forge/config.yaml` — surface configuration
- `justfile` — generated recipes
