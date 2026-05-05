---
name: graduate-tests
description: Migrate feature test scripts to the regression suite (tests/e2e/). Agent-driven: reads scripts, analyzes content, decides classification, splits/merges as needed, rewrites imports, creates graduation marker.
---

# /graduate-tests

Migrate feature test scripts from the staging area `tests/e2e/features/<feature>/` to the project-wide regression suite at `tests/e2e/<target>/`.

**Core principle**: Graduation is a decision, not a file copy. The agent reads and understands each spec before deciding where it belongs.

<HARD-GATE>
- Do NOT overwrite existing files in `tests/e2e/` without merging
- Do NOT graduate if marker already exists (idempotent)
- Do NOT modify the source scripts in `tests/e2e/features/<slug>/`
</HARD-GATE>

## Prerequisites

Check before running. Abort and prompt user if missing:

| Artifact | Condition | Action if not met |
|----------|-----------|-------------------|
| `tests/e2e/features/<slug>/` directory | Must exist (or scripts at `tests/e2e/<slug>/` via staging bypass — see Step 1.5) | Run `/gen-test-scripts` first |
| At least one `.spec.ts` file | Must exist | Run `/gen-test-scripts` first |
| `tests/e2e/helpers.ts` | Must exist (symbol completeness checked after Step 2) | Run `/gen-test-scripts` first |
| `tests/e2e/features/<slug>/results/latest.md` | Must show PASS | Run `/run-e2e-tests` first — only graduate passing tests |
| `tests/e2e/.graduated/<slug>` | Must NOT exist | Already graduated — skip |

<PRINCIPLE>
**Shared infrastructure first.** Before executing graduation, verify that shared dependencies (`helpers.ts`, `config.yaml`, `package.json`, `tsconfig.json`, `playwright.config.ts`) are complete and functional. After graduation, spec file import paths are rewritten from `'../../helpers.js'` to `'../helpers.js'` — if `helpers.ts` itself is incomplete, the rewritten imports will still fail to compile. When inconsistencies are found, go back to `/gen-test-scripts` to fix shared dependencies before graduating.
</PRINCIPLE>

## When to Use

- After all tasks for a feature are completed and e2e tests pass
- User invokes `/graduate-tests` manually
- `/run-tasks` orchestrator suggests it post-completion

## Workflow

```
1. Check marker → 2. Read scripts → 3. Analyze structure → 4. Decide classification → 5. Migrate → 5.5 Validate → 6. Write marker → 7. Cleanup source
```

### Step 1: Check Graduation Marker

If `tests/e2e/.graduated/<slug>` exists: print "Already graduated on <timestamp>" and stop.

### Step 1.5: Detect Staging Bypass

Check whether spec files exist at the staging path `tests/e2e/features/<slug>/`:
- **Found at staging path** → proceed normally with Step 2
- **NOT found at staging path** → check `tests/e2e/<slug>/` (direct e2e subdirectory without `features/` prefix)

If spec files are found at `tests/e2e/<slug>/` but NOT at `tests/e2e/features/<slug>/`, this indicates a staging bypass — the generation step wrote directly to a post-graduation location. **Do NOT skip classification.** Instead:
1. Treat `tests/e2e/<slug>/` as the source directory for this graduation run
2. Proceed with Step 2 using this source path
3. Step 4 (functional module classification) is MANDATORY — classify by business domain and move to the correct `tests/e2e/<module>/`
4. Log a warning: "Staging bypass detected: source at tests/e2e/<slug>/ instead of tests/e2e/features/<slug>/. Performing classification anyway."

<HARD-RULE>
Graduation MUST ALWAYS perform Step 4 (functional module classification), regardless of where the source scripts are found. The agent MUST NOT treat scripts already at `tests/e2e/<any-dir>/` as "already graduated" — only the marker at `tests/e2e/.graduated/<slug>` determines graduation status.
</HARD-RULE>

### Step 2: Read Source Scripts

Read all `.spec.ts` files from the source directory determined in Step 1.5 (either `tests/e2e/features/<slug>/` or `tests/e2e/<slug>/` for bypass) and `helpers.ts`. Understand what each `describe`/`test` block tests, which routes/APIs/CLI commands are covered, and whether a single spec mixes multiple functional domains.

**Symbol completeness check**: Extract the set of imported symbols from `helpers.js` (e.g., `screenshot`, `baseUrl`, `curl`, `runCli`). Verify each symbol is exported by `tests/e2e/helpers.ts`. If any are missing, abort and prompt: "helpers.ts is missing exports (X, Y). Run `/gen-test-scripts` first to merge missing symbols."

### Step 3: Analyze Existing Structure

Read `tests/e2e/` directory structure to understand the existing classification convention (by type, by route, by feature domain). New specs must follow the same convention.

**Pre-flight check**: Run `cd tests/e2e && npx tsc --noEmit`. If compilation fails on pre-existing files, abort before touching anything — migrating into a broken codebase compounds errors.

### Step 4: Decide Classification

For each spec file, answer:

| Question | Decision |
|----------|----------|
| What functional module does this spec cover? | One module -> keep as-is; multiple modules -> split by module |
| Which `tests/e2e/<module>/` does it belong to? | Match by functional domain |
| Does a spec file already exist at the target path? | Yes -> merge test blocks (deduplicate by test name), not overwrite |

**Functional module** = the business domain or product area being tested. NOT the test type (UI/API/CLI) and NOT the feature slug.

Classification patterns:
```
# Split: ui.spec.ts (login+dashboard) → tests/e2e/auth/login.spec.ts + tests/e2e/dashboard/ui.spec.ts
# Keep: justfile-integration/*.spec.ts → tests/e2e/justfile/cli.spec.ts, tests/e2e/justfile/detection.spec.ts
# Merge: new profile/api.spec.ts into existing tests/e2e/profile/api.spec.ts (deduplicate by test title)
```

### Step 5: Execute Migration

Create backup directory: `mkdir -p tests/e2e/.graduated/.backup/<slug>`

**Backup path convention**: `<sanitized-path>` = target path relative to `tests/e2e/` with `/` and `\` replaced by `__`. Example: `justfile/cli.spec.ts` → `justfile__cli.spec.ts`.

For each target file:

1. **Create directory** if it doesn't exist
2. **Write spec file** (or merge into existing)
3. **Record in migration manifest**: append entries to `tests/e2e/.graduated/.backup/<slug>/manifest.json` after each file operation (write-ahead log). See template: `plugins/forge/skills/graduate-tests/templates/manifest.json`. On re-run after partial failure, read existing manifest and continue — do not reset.

**Merge procedure** (when target file already exists). Full example: `plugins/forge/skills/graduate-tests/templates/merge-example.md`:
1. Read both source and target spec files
2. **Backup** the target file (only if no backup exists — prevents overwriting original on re-run): `test -f <backup-path> || cp <target-path> <backup-path>`
3. **Merge rules**:
   - Combine imports, deduplicate
   - Match `test.describe` blocks by title — merge their children into a single block
   - Deduplicate `test()` entries by full title string match (identical titles → keep source version; different titles with same TC ID prefix → keep both)
   - Preserve `test.describe` nesting depth — do not flatten
   - Append new describe blocks that don't exist in target
4. Write the merged file

<HARD-RULE>
Specs in the staging area (`tests/e2e/features/<slug>/`) import helpers via `'../../helpers.js'` (two levels up). After migration to the regression suite, the import must be rewritten based on target depth:

| Target location | Import path |
|----------------|-------------|
| `tests/e2e/<module>/file.spec.ts` (1 level deep) | `'../helpers.js'` |
| `tests/e2e/<module>/sub/file.spec.ts` (2 levels deep) | `'../../helpers.js'` |
| `tests/e2e/<module>/a/b/file.spec.ts` (3 levels deep) | `'../../../helpers.js'` |

Formula: count directory levels from spec file to `tests/e2e/`, then generate `'../' * levels + 'helpers.js'`. Every migrated spec file MUST have its helpers import path updated. Other imports (node built-ins, @playwright/test) remain unchanged.
</HARD-RULE>

Shared infrastructure (`helpers.ts`, `package.json`, `tsconfig.json`) already exists at `tests/e2e/` — no merging or copying needed.

### Step 5.5: Validate Migration

After migrating all spec files:

1. Verify TypeScript compilation: `cd tests/e2e && npx tsc --noEmit`
2. Verify Playwright discovers all tests: `cd tests/e2e && npx playwright test --list`

If validation fails and is unfixable, rollback using the migration manifest:
- **Newly created** target files: delete them entirely
- **Merged** target files: revert by restoring from `tests/e2e/.graduated/.backup/<slug>/`
- **Update manifest entries**: change `"status": "done"` to `"status": "rolled-back"` for affected entries (prevents stale entries on re-run)
- Do NOT write the marker. Source directory remains intact for retry.

### Step 6: Create Graduation Marker

Write marker only after Step 5.5 validation passes (atomic — no marker = not graduated). Template: `plugins/forge/skills/graduate-tests/templates/graduation-marker.yaml`:

```bash
mkdir -p tests/e2e/.graduated
cat > tests/e2e/.graduated/<slug> <<EOF
schema_version: 1
status: completed
timestamp: <UTC ISO timestamp>
source: tests/e2e/features/<slug>/
targets:
  - tests/e2e/<module>/<spec-file>
modules:
  - <module-name>
testCount: <N>
EOF
```

**Atomicity**: The marker is written ONLY after validation passes. (Legacy markers may have `source:` paths without the `features/` prefix — match by slug filename alone for idempotency.)

### Step 7: Source Cleanup

After the marker is written, archive test results then remove the source directory:

```bash
# Archive results before cleanup
if [ -d "tests/e2e/features/<slug>/results/" ]; then
    mkdir -p tests/e2e/.graduated/.results-archive
    cp -r tests/e2e/features/<slug>/results/ tests/e2e/.graduated/.results-archive/<slug>/
fi
rm -rf tests/e2e/features/<slug>/
rm -rf tests/e2e/.graduated/.backup/<slug>/
```

<HARD-RULE>
Source cleanup MUST NOT happen before the marker is written. If Step 5.5 validation fails, the source directory MUST remain intact for retry. Only remove the source if ALL spec files were successfully migrated, validated, and the marker is written.
</HARD-RULE>

## Output

Report to user:

```
Graduated <slug>:
  ui.spec.ts → tests/e2e/ui/login/login.spec.ts
  api.spec.ts → tests/e2e/api/auth/auth.spec.ts

Marker: tests/e2e/.graduated/<slug>
```

## Related Skills

| Skill | Usage |
|-------|-------|
| `/gen-test-scripts` | Generate source scripts before graduation |
| `/run-e2e-tests` | Execute scripts before graduating |
| `/run-tasks` | Suggests graduation after all tasks complete |
