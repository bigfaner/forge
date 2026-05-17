---
name: gen-test-scripts
description: Generate executable e2e test scripts from test cases. Profile-aware: reads active test profile from .forge/config.yaml to determine test framework, templates, and conventions.
conventions:
  - testing-isolation.md
---

# Gen Test Scripts

Generate executable e2e test scripts from test cases.

**Core principle**: Every generated `test()` must be independently runnable, repeatable, and have explicit assertions. Test cases are input; scripts are output.

<HARD-GATE>
This skill ONLY writes to `tests/e2e/features/<feature>/` (staging area). It does NOT execute tests (handled by `/run-e2e-tests`).

**FORBIDDEN output paths**: `tests/e2e/<feature>/`, `tests/e2e/<module>/`, any path outside `tests/e2e/features/`. The `features/` prefix is a staging area enforcing: Phase 1: `/gen-test-scripts` → `tests/e2e/features/<feature>/` → Phase 2: `/graduate-tests` → `tests/e2e/<module>/`.
</HARD-GATE>

## Step 0: Resolve Profile

1. Run `forge profile` to get active profile(s). On failure (`PROFILE: (none)`): ask user to choose from `web-playwright`, `go-test`, `maestro`, `java-junit`, `rust-test`, `pytest`. Run `forge profile set <name>` to persist.
2. Run `forge profile get <profile-name> --manifest` and `--generate`.

<HARD-RULE>
Do NOT silently default to any profile. If `forge profile` returns no result and the user cannot decide, abort the skill.
</HARD-RULE>

## Convention Loading

1. **Project-wide**: Read this skill's frontmatter `conventions` field. For each filename, check `docs/conventions/{filename}` — load if exists, skip silently if missing.
2. **Per-type conventions** (per-type mode only): Read `plugins/forge/skills/gen-test-scripts/types/{type}.md` frontmatter `conventions` field. For each filename, check `docs/conventions/{filename}` — load if exists, skip silently if missing.
3. **Legacy mode**: Only project-wide conventions loaded.

<HARD-RULE>
Convention loading is non-blocking. Missing convention files are silently skipped.
</HARD-RULE>

## Type Filter (--type)

Optional `--type <capability>` argument filters script generation to a single test type.

**Valid values**: `ui`, `tui`, `mobile`, `api`, `cli`. On invalid value: error `"invalid type: {value}. Valid types: ui, tui, mobile, api, cli"`. No partial execution.

| Input mode | `--type` specified | `--type` omitted |
|------------|-------------------|-----------------|
| **Per-type** | Load only `{type}-test-cases.md` matching filter | Load all `{type}-test-cases.md` sequentially |
| **Legacy** | Filter test cases by type | Process all types |

Shared infrastructure (Step 3.5) **always runs** regardless of `--type`. Without `--type`, all types generate as before.

## Prerequisites

**Test Case File Discovery** (first match wins):
1. **Per-type mode**: Glob `docs/features/<slug>/testing/*-test-cases.md`
2. **Legacy mode**: `docs/features/<slug>/testing/test-cases.md`

Per-type files take precedence over legacy. Missing? Prompt to run `/gen-test-cases`. `<slug>` from `forge feature`.

### Step Actionability Gate

If eval-test-cases report exists (`docs/features/<slug>/testing/eval/iteration-*.md`) and Step Actionability score < 200: **ABORT** with message to fix test cases first. If no report exists, proceed normally.

## Workflow

```
0. Resolve profile → 1. Read test cases → 1.5. Code Reconnaissance → 3.5. Shared infrastructure → 4. Dispatch to type files
```

<HARD-RULE>
**Task Splitting Guard**: Global Setup Phase (Steps 0→1→1.5→3.5) runs as single pre-task FIRST. Only Step 4 may be parallelized. Each sub-task must include: "Shared infrastructure already configured. Auth: use storageState/cached token, NOT login() in beforeEach."
</HARD-RULE>

### Step 1: Read Test Cases

**Per-type mode**: For each `{type}-test-cases.md` (or only the file matching `--type`), parse test cases directly — file is already single-type. Parse TC ID, title, type, route, feature, pre-conditions, steps, expected result, priority.

**Legacy mode**: Read `test-cases.md`, parse test cases, group by type using profile capabilities. If `--type` specified, keep only matching group.

#### Auth Classification

| Category | Detection | Strategy |
|----------|-----------|----------|
| **login-test** | Target matches auth/login routes | No shared auth. Independent login/logout. Invalidate cached credentials after. |
| **auth-required-test** | Pre-conditions contain "authenticated"/"logged in" or Target implies protected resource | Cached shared auth — credentials acquired once, reused |
| **public-test** | No auth pre-conditions, public resource | No auth needed |
| **custom-auth-test** | "API key", "X-API-Key", "OAuth", "session cookie" in pre-conditions | Detect mechanism from codebase, generate custom auth setup |

**Priority** (first match wins): custom-auth-test → login-test → auth-required-test → public-test

**Auth Plan Output** (consumed by Step 3.5 and Step 4):
```
Auth Plan: auth-required-test: N, login-test: N, public-test: N, custom-auth-test: N, shared-auth-enabled: yes/no
```

<HARD-RULE>
Login tests and authenticated tests must not be mixed in the same `describe` block.
</HARD-RULE>

### Step 1.5: Code Reconnaissance (Build Fact Table)

Read source code to extract ground-truth values. **Never guess** — every value must come from test-cases.md or this Fact Table. If `--type` specified, build entries only for that type.

**Generic reconnaissance reads**:

| Source | What to extract |
|--------|-----------------|
| Router files | Route paths, path parameters, middleware bindings |
| Config files | API port, base path prefix, auth credentials |
| API handlers | Request/response schemas, status codes |
| Auth implementation | Login endpoint, token field name, header format |
| CLI entry points | Command names, flag names, output formats |
| Frontend UI components | `data-testid` values, component structure, dynamic patterns |

Build Fact Table with source citations:

```markdown
## Fact Table
| Key | Value | Source |
|-----|-------|--------|
| API_PORT | 8080 | backend/config.yaml:3 |
```

<HARD-RULE>
- Every Fact Table value must cite source file and line number. Unknown sources → `UNKNOWN`. Do not fabricate.
- Fact Table values override template placeholders. All `// VERIFY:` markers must be resolved using Fact Table values.
</HARD-RULE>

#### Fact Table Completeness Gate

Type-specific required keys and completeness rules are defined in each type file. The generic gate rule: if *all* required keys for a type are UNKNOWN, skip that type with WARNING. Individual UNKNOWN keys are acceptable.

### Step 3.5: Ensure Shared Infrastructure

<PRINCIPLE>
**This step always runs regardless of `--type` value.** Shared infrastructure (helpers, config, auth-setup) is used by all test types. Create-only-if-missing logic ensures safe concurrent execution.
</PRINCIPLE>

Check `tests/e2e/`. If any shared file is missing, create from profile manifest templates. Strip `// VERIFY:` and `// TEMPLATE:` comments from copies. If helper file exists but missing exports needed by Auth Plan, merge missing exports (do NOT overwrite existing).

#### Auth Infrastructure

When `shared-auth-enabled: yes`:
1. Follow profile's `generate.md` for framework-specific auth setup
2. Playwright: generate `auth-setup.ts`, configure `playwright.config.ts` `projects` with `storageState`
3. Other profiles: follow `generate.md`

<HARD-RULE>
Auth-required tests MUST use shared auth (storageState/cached tokens), NOT per-test login in beforeEach. Playwright: calling `login()` in `beforeEach` is FORBIDDEN when `storageState` is configured.
</HARD-RULE>

<HARD-RULE>
**Shared file policy**: Create from templates if missing. Merge missing exports into helpers if exists. Do NOT modify other existing shared files (config, package manifest).
</HARD-RULE>

### Step 4: Dispatch to Type Files

For each active type, load the corresponding type instruction file and follow its instructions.

**Type-to-file mapping** (hard-coded):

| Type | Instruction File |
|------|-----------------|
| ui | `plugins/forge/skills/gen-test-scripts/types/ui.md` |
| tui | `plugins/forge/skills/gen-test-scripts/types/tui.md` |
| mobile | `plugins/forge/skills/gen-test-scripts/types/mobile.md` |
| api | `plugins/forge/skills/gen-test-scripts/types/api.md` |
| cli | `plugins/forge/skills/gen-test-scripts/types/cli.md` |

**Dispatch procedure** for each active type:
1. Verify type file exists. If missing: **HALT** with `"Type file missing: {path}. Cannot generate {type} test scripts. Available types: ui, tui, mobile, api, cli."`
2. Read type file, execute its instructions (reconnaissance, Fact Table keys, verification, generation patterns, antipattern guards)
3. Generate spec files following type file patterns and profile's `generate.md`

If `--type` specified, only dispatch to the matching type file.

#### beforeAll Safety

All setup blocks must: (1) wrap operations in try/catch, (2) explicit null/undefined check after assignment, (3) prefer reusing existing resources, (4) use sequential execution for dependent tests, (5) use retry for backend-dependent operations.

<HARD-RULE>
Uninitialized module-level variables cause misleading downstream failures. Try/catch + explicit check ensures the real failure is reported.
</HARD-RULE>

#### Antipattern Guard (Mandatory Pre-Emission Check)

Verify each generated test function does not match any forbidden pattern. Fix before writing.

| # | Pattern | Harmful Because | Instead |
|---|---------|-----------------|---------|
| 1 | Recursive test invocation (spawning test runner within test) | Process explosion; 6GB+ RAM on Windows | Recursion guard (env var) or `-run` flag filtering |
| 2 | Unconditional `t.Skip` with no env detection | Dead code inflating coverage signal | Implement with fixture setup or don't generate |
| 3 | Vacuous assertions (`if { assert }` without else skip/fail) | Assertion never reached; false confidence | Every assertion reachable on every code path |
| 4 | Skip based on environment state, no self-contained fixture | Non-deterministic pass/skip | `t.TempDir()` + own project structure |
| 5 | Duplicate test function names across packages | Double CI time; divergent copies | Scan for collisions; use unique names with feature slug |
| 6 | Static-file text grep (assert on source file content) | Coupled to prose, not behavior | Test runtime behavior only (invoke + assert outputs) |

#### Generation HARD-RULEs

**Traceability**: Each test function includes traceability comment to PRD. **Skip empty groups**. **Empty result guard**: if zero test files generated, abort with diagnostic. **VERIFY markers**: resolve from Fact Table; if no value, keep `// VERIFY:` as-is. **Post-generation**: run `just e2e-compile` (fix errors), `just e2e-verify --feature <slug>` (check unresolved VERIFY), verify helper exports cover generated imports. **No fixed delays/sleeps. No hardcoded URLs. No debug output.** Selector strategy follows `generate.md`.

<EXTREMELY-IMPORTANT>
All framework-specific rules are in the active profile's `generate.md`. Use ONLY that framework. `@eN` refs must not appear in any generated test file.
</EXTREMELY-IMPORTANT>

## Output

All generated test files go to `tests/e2e/features/<feature>/` (staging area). After `/run-e2e-tests` verification, use `/graduate-tests` to migrate.

## Error Handling

| Situation | Action |
|-----------|--------|
| Profile resolution fails | Ask user to select a profile |
| No test case files found | Abort, prompt `/gen-test-cases` |
| Step Actionability < 200 | Abort, prompt fix test cases |
| Unknown `--type` value | Error listing valid types. No generation. |
| Type file missing | Halt naming missing file and valid types |
| Compilation fails | Fix generated code, re-compile |
| Zero test files generated | Abort with diagnostic |
| Source not found for Fact Table | Mark `UNKNOWN`, don't fabricate |
| All required keys UNKNOWN for a type | Skip type with WARNING |
| Helper missing needed exports | Merge from template |

## Related Skills

| Skill | Usage |
|-------|-------|
| `/gen-test-cases` | Generate test cases from PRD |
| `/gen-sitemap` | Generate and maintain sitemap.json |
| `/run-e2e-tests` | Execute test scripts and report results |
