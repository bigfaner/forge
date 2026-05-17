# Rust Test Graduate Strategy

Profile-specific graduation rules for the `graduate-tests` skill.

## Source File Discovery

| Item | Value |
|------|-------|
| File extension | `.rs` |
| Source directory | `tests/e2e/features/<slug>/` (staging) |
| Target directory | `tests/e2e/<module>/` (regression) |

## Import Rewrite

No import rewriting needed. Rust uses module paths that are position-independent.

## Validation

| Check | Command | Failure action |
|-------|---------|----------------|
| Pre-flight compilation | `just e2e-compile` | Abort before touching anything |
| Post-migration compilation | `just e2e-compile` | Rollback via migration manifest |
| Test discovery | `just e2e-discover` | Rollback via migration manifest |

## Merge Procedure

When a target file already exists:

1. Read both source and target Rust files
2. Backup target file
3. Merge `use` statements at the top, deduplicate
4. Merge test functions — deduplicate by function name
5. Keep helper functions and modules intact
6. Write the merged file

## Shared Infrastructure

Rust test helpers in `tests/e2e/` are typically `mod helpers;` declarations. The helpers module stays in place during graduation.

## Graduation Marker

Standard marker format (see `plugins/forge/skills/graduate-tests/templates/graduation-marker.yaml`).
