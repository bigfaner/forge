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
| `tests/e2e/features/<slug>/` directory | Must exist | Run `/gen-test-scripts` first |
| At least one `.spec.ts` file | Must exist | Run `/gen-test-scripts` first |
| `tests/e2e/helpers.ts` | Must exist (symbol completeness checked after Step 2) | Run `/gen-test-scripts` first |
| `tests/e2e/features/<slug>/results/latest.md` | Must show PASS | Run `/run-e2e-tests` first — only graduate passing tests |
| `tests/e2e/.graduated/<slug>` | Must NOT exist | Already graduated — skip |

<PRINCIPLE>
**共享基础设施优先。** 执行毕业操作前，先验证公共依赖（`helpers.ts`、`config.yaml`、`package.json`、`tsconfig.json`、`playwright.config.ts`）完整且可用。毕业后的 spec 文件 import 路径会从 `'../../helpers.js'` 改写为 `'../helpers.js'`，如果 `helpers.ts` 本身不完整，改写后依然无法通过编译。发现不一致时先回到 `/gen-test-scripts` 修复公共依赖，再执行毕业。
</PRINCIPLE>

```bash
task feature   # get current slug
ls tests/e2e/features/<slug>/
grep -q 'PASS\|passed' tests/e2e/features/<slug>/results/latest.md || echo "Error: e2e tests not passed yet" >&2
cat tests/e2e/.graduated/<slug> 2>/dev/null && echo "already graduated"
```

## When to Use

- After all tasks for a feature are completed and e2e tests pass
- User invokes `/graduate-tests` manually
- `/run-tasks` orchestrator suggests it post-completion

## Workflow

```
1. Check marker → 2. Read scripts → 3. Analyze structure → 4. Decide classification → 5. Migrate → 5.5 Validate → 6. Write marker → 7. Cleanup source
```

### Step 1: Check Graduation Marker

```bash
cat tests/e2e/.graduated/<slug>
```

If marker exists: print "Already graduated on <timestamp>" and stop.

### Step 2: Read Source Scripts

Read all files in `tests/e2e/features/<slug>/`:

```bash
ls tests/e2e/features/<slug>/
```

Read each `.spec.ts` file and `helpers.ts`. Understand:
- What each `describe`/`test` block tests
- Which routes, APIs, or CLI commands are covered
- Whether a single spec mixes multiple functional domains

**Symbol completeness check**: After reading all spec files, extract the set of imported symbols from `helpers.js` (e.g., `screenshot`, `baseUrl`, `curl`, `runCli`). Verify each symbol is exported by `tests/e2e/helpers.ts`. If any are missing, abort and prompt: "helpers.ts is missing exports (X, Y). Run `/gen-test-scripts` first to merge missing symbols."

### Step 3: Analyze Existing Structure

If `tests/e2e/` exists, read its directory structure:

```bash
ls -R tests/e2e/
```

Understand the existing classification convention (by type, by route, by feature domain). New specs must follow the same convention.

**Pre-flight check**: Before migration, verify existing target files are healthy:
```bash
cd tests/e2e && npx tsc --noEmit
```
If compilation fails on pre-existing files, report and abort before touching anything — migrating into a broken codebase compounds errors.

### Step 4: Decide Classification

For each spec file, answer:

| Question | Decision |
|----------|----------|
| What functional module does this spec cover? | One module -> keep as-is; multiple modules -> split by module |
| Which `tests/e2e/<module>/` does it belong to? | Match by functional domain |
| Does a spec file already exist at the target path? | Yes -> merge test blocks (deduplicate by test name), not overwrite |

**Functional module** = the business domain or product area being tested. NOT the test type (UI/API/CLI) and NOT the feature slug.

Classification examples:
```
# Input: tests/e2e/features/user-auth-feature/
#   ui.spec.ts    -> contains login, dashboard, profile tests
#   api.spec.ts   -> all auth-related API tests
#
# Agent decides by functional domain:
#   ui.spec.ts -> split:
#     tests/e2e/auth/login.spec.ts      # authentication module
#     tests/e2e/dashboard/ui.spec.ts    # dashboard module
#   api.spec.ts -> tests/e2e/auth/api.spec.ts  # stays in auth module
```

```
# Input: tests/e2e/features/justfile-integration/
#   cli.spec.ts           -> justfile CLI commands
#   detection-assembly.spec.ts -> project type detection
#   forge-justfile.spec.ts     -> justfile structure validation
#
# All test the same functional domain -> keep together:
#   tests/e2e/justfile/cli.spec.ts
#   tests/e2e/justfile/detection.spec.ts
#   tests/e2e/justfile/structure.spec.ts
```

```
# Input: tests/e2e/features/user-profile/
#   api.spec.ts   -> profile CRUD + settings API
#
# Existing target: tests/e2e/profile/api.spec.ts already has 3 tests
# Merge path — deduplicate by test name, append new tests:
#   tests/e2e/profile/api.spec.ts  (now contains 3 existing + 4 new, deduplicated)
```

### Step 5: Execute Migration

Before migrating, create a slug-scoped backup directory for rollback:

```bash
mkdir -p tests/e2e/.graduated/.backup/<slug>
```

**Backup path convention**: `<sanitized-path>` = the target path relative to `tests/e2e/` with `/` and `\` replaced by `__`. Example: `tests/e2e/justfile/cli.spec.ts` → `justfile__cli.spec.ts`.

For each target file:

1. **Create directory** if it doesn't exist
2. **Write spec file** (or merge into existing)
3. **Record in migration manifest**: append `{targetPath, wasExistingBeforeMerge: boolean}` for each migrated file. Persist the manifest to `tests/e2e/.graduated/.backup/<slug>/manifest.json` after each file operation (write-ahead log pattern). On re-run after partial failure, read the existing manifest and continue from where it left off — do not reset. If a target was already recorded in a prior iteration, update its entry (do not duplicate). This manifest is used in Step 5.5 rollback.

**Merge procedure** (when target file already exists):
1. Read both source and target spec files
2. **Backup** the target file (only if no backup exists for this file — prevents overwriting original on re-run after partial failure): `test -f tests/e2e/.graduated/.backup/<slug>/<sanitized-path> || cp <target-path> tests/e2e/.graduated/.backup/<slug>/<sanitized-path>` (slug-scoped to avoid collision with concurrent graduations)
3. Extract all `import` statements from both files, deduplicate, and combine into a single import block
4. Walk the AST-like nesting tree: for each `test.describe` block, collect its direct `test()` children and nested `test.describe` children recursively
5. **Merge describe blocks** (matched by describe title): `test.describe` blocks with the same title are *merged*, not deduplicated. Combine their children into a single block. This is distinct from test dedup — describe blocks are containers, not leaf nodes.
6. **Deduplicate individual `test()` entries** (by full title string match): if two tests have identical titles, keep the source version. If titles differ but share a TC ID prefix (e.g., `TC-001: Login` vs `TC-001: Different test`), keep both — TC IDs alone are not globally unique across features. Only tests with identical titles are considered duplicates. Re-graduation of the same feature is blocked by the idempotency check (Step 1). **Preserve `test.describe` nesting**: do not flatten nested describe blocks into the parent. If source has `test.describe('A', () => { test.describe('B', () => { test('TC-002') }) })`, the merged file must retain that nesting, not extract `TC-002` to the top level
7. Combine: append new describe blocks that don't exist in target
7. Write the merged file preserving the target's existing structure where possible

<HARD-RULE>
Specs in the staging area (`tests/e2e/features/<slug>/`) import helpers via `'../../helpers.js'` (two levels up). After migration to the regression suite (`tests/e2e/<target>/`), the import must be rewritten to `'../helpers.js'` (one level up). Every migrated spec file MUST have its helpers import path updated from `'../../helpers.js'` to `'../helpers.js'`. Other imports (node built-ins, @playwright/test) remain unchanged.

Note: This rule assumes targets are at `tests/e2e/<target>/` (one level deep). If the agent places specs in a nested directory (e.g., `tests/e2e/<target>/sub/`), compute the relative path to `tests/e2e/helpers.ts` accordingly — two levels would still need `'../../helpers.js'`.
</HARD-RULE>

Shared infrastructure (`helpers.ts`, `package.json`, `tsconfig.json`) already exists at `tests/e2e/` — no merging or copying needed.

### Step 5.5: Validate Migration

After migrating all spec files:

1. Verify TypeScript compilation:
```bash
cd tests/e2e && npx tsc --noEmit
```

2. Verify Playwright discovers all tests:
```bash
cd tests/e2e && npx playwright test --list
```

If validation fails:
1. Read the compilation/discovery error
2. Fix the migrated spec files (usually import path issues)
3. Re-run validation
4. If unfixable: rollback using the migration manifest from Step 5:
   - **Newly created** target files (did not exist before this migration): delete them entirely
   - **Merged** target files (existed before this migration, with prior graduation content): revert to pre-merge state by restoring from `tests/e2e/.graduated/.backup/<slug>/`.
   Report the error to the user and do NOT write the marker. Source directory remains intact for retry.

### Step 6: Create Graduation Marker

Write marker only after Step 5.5 validation passes (atomic — no marker = not graduated):

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

**Atomicity**: The marker is written ONLY after validation passes. If migration is interrupted, no marker exists and re-running will re-attempt.

**Note on legacy markers**: Markers created before the staging architecture was introduced may have `source:` paths without the `features/` prefix (e.g., `tests/e2e/<slug>/` instead of `tests/e2e/features/<slug>/`). This reflects the actual source at the time of graduation. When checking idempotency (Step 1), match the marker filename (`tests/e2e/.graduated/<slug>`) by slug alone — do not rely on the `source:` path format for idempotency checks. New markers MUST use `tests/e2e/features/<slug>/`.

### Step 7: Source Cleanup

After the marker is written, remove the source directory and clean up backups:

```bash
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
  ui.spec.ts → tests/e2e/ui/dashboard/dashboard.spec.ts
  api.spec.ts → tests/e2e/api/auth/auth.spec.ts
  cli.spec.ts → tests/e2e/cli/cli.spec.ts

Marker: tests/e2e/.graduated/<slug>
```

## Related Skills

| Skill | Usage |
|-------|-------|
| `/gen-test-scripts` | Generate source scripts before graduation |
| `/run-e2e-tests` | Execute scripts before graduating |
| `/run-tasks` | Suggests graduation after all tasks complete |
