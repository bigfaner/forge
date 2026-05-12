# Pytest Graduate Strategy

Profile-specific graduation rules for the `graduate-tests` skill.

## Source File Discovery

| Item | Value |
|------|-------|
| File extension | `.py` |
| Naming pattern | `test_*.py` or `*_test.py` |
| Source directory | `tests/e2e/features/<slug>/` (staging) |
| Target directory | `tests/e2e/<module>/` (regression) |

## Import Rewrite

No import rewriting needed. Python imports use module paths that are position-independent when `__init__.py` files are present.

## Validation

| Check | Command | Failure action |
|-------|---------|----------------|
| Pre-flight compilation | `just e2e-compile` | Abort before touching anything |
| Post-migration compilation | `just e2e-compile` | Rollback via migration manifest |
| Test discovery | `just e2e-discover` | Rollback via migration manifest |

## Merge Procedure

When a target file already exists:

1. Read both source and target Python files
2. Backup target file
3. Merge imports at the top, deduplicate
4. Merge test functions — deduplicate by function name
5. Keep fixtures and helper functions intact
6. Write the merged file

## Shared Infrastructure

- `tests/e2e/conftest.py` — shared pytest fixtures (stays in place)
- `tests/e2e/helpers.py` — shared test utilities (stays in place)
- `tests/e2e/__init__.py` — package marker (stays in place)

## Graduation Marker

Standard marker format (see `plugins/forge/skills/graduate-tests/templates/graduation-marker.yaml`).
