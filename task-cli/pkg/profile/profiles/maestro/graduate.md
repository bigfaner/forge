# Maestro Graduate Strategy

Profile-specific graduation rules for the `graduate-tests` skill.

## Source File Discovery

| Item | Value |
|------|-------|
| File extension | `.yaml` |
| Source directory | `tests/e2e/features/<slug>/` (staging) |
| Target directory | `tests/e2e/<module>/` (regression) |

## Import Rewrite

No import rewriting needed (YAML has no import mechanism).

## Validation

| Check | Command | Failure action |
|-------|---------|----------------|
| YAML syntax | `just e2e-compile` | Report syntax error |
| Test discovery | `just e2e-discover` | Rollback via migration manifest |

No compilation check (YAML is interpreted).

## Merge Procedure

When a target file already exists:

1. Read both YAML files
2. Backup target file
3. Parse YAML, concatenate flow lists
4. Deduplicate flows by `name` field
5. Write merged YAML

## Shared Infrastructure

- `tests/e2e/config.yaml` — Maestro environment config (stays in place)

## Graduation Marker

Standard marker format (see `plugins/forge/skills/graduate-tests/templates/graduation-marker.yaml`).
