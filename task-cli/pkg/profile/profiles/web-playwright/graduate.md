# Playwright Graduate Strategy

Profile-specific graduation rules for the `graduate-tests` skill.

## Source File Discovery

| Item | Value |
|------|-------|
| File extension | `.spec.ts` |
| Source directory | `tests/e2e/features/<slug>/` (staging) |
| Target directory | `tests/e2e/<module>/` (regression) |

## Import Rewrite

Staging specs import helpers via `'../../helpers.js'` (two levels up). After graduation, rewrite based on target depth:

| Target location | New import path |
|----------------|-----------------|
| `tests/e2e/<module>/file.spec.ts` (1 level) | `'../helpers.js'` |
| `tests/e2e/<module>/sub/file.spec.ts` (2 levels) | `'../../helpers.js'` |
| `tests/e2e/<module>/a/b/file.spec.ts` (3 levels) | `'../../../helpers.js'` |

Formula: count directory levels from spec file to `tests/e2e/`, generate `'../' * levels + 'helpers.js'`.

Other imports (node built-ins, `@playwright/test`) remain unchanged.

## Validation

| Check | Command | Failure action |
|-------|---------|----------------|
| Pre-flight compilation | `just e2e-compile` | Abort before touching anything |
| Post-migration compilation | `just e2e-compile` | Rollback via migration manifest |
| Test discovery | `just e2e-discover` | Rollback via migration manifest |

## Symbol Completeness

Extract imported symbols from spec files (e.g., `screenshot`, `baseUrl`, `curl`, `runCli`). Verify each is exported by `tests/e2e/helpers.ts`. Abort if any missing — prompt to run `/gen-test-scripts` first.

## Shared Infrastructure

These files already exist at `tests/e2e/` and must NOT be copied or modified during graduation:

- `helpers.ts`
- `package.json`
- `tsconfig.json`
- `playwright.config.ts`
- `config.yaml`

## Merge Procedure

When a target file already exists at the graduation destination:

1. Read both source and target spec files
2. Backup target file (only if no backup exists — prevents overwriting original on re-run)
3. Combine imports, deduplicate
4. Match `test.describe` blocks by title — merge their children into a single block
5. Deduplicate `test()` entries by full title string match (identical titles → keep source version; different titles with same TC ID prefix → keep both)
6. Preserve `test.describe` nesting depth — do not flatten
7. Append new describe blocks that don't exist in target
8. Write the merged file

## Graduation Marker

Written only after validation passes (atomic — no marker = not graduated):

```yaml
schema_version: 1
status: completed
timestamp: <UTC ISO timestamp>
source: tests/e2e/features/<slug>/
targets:
  - tests/e2e/<module>/<spec-file>
modules:
  - <module-name>
testCount: <N>
```
