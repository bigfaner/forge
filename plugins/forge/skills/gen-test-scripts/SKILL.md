---
name: gen-test-scripts
description: Generate executable e2e test scripts from test cases. Language-aware: auto-detects test framework from project files, uses forge testing CLI for strategy and templates.
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

## Step 0: Resolve Language and Strategy

1. **Detect language**: Run `forge testing detect` to auto-detect the project's test language(s) from file signals.
2. **On failure** (no language detected): ask the user to add `languages` to `.forge/config.yaml` (e.g., `languages: [go]`).
3. **Load strategy**: Run `forge testing get generate` to load the generate strategy for the detected language.

Use the loaded strategy for all subsequent steps.

<HARD-RULE>
Do NOT silently default to any language. If `forge testing detect` returns no result and the user cannot configure `languages`, abort the skill.
</HARD-RULE>

## Convention Loading

After language resolution and before entering the workflow steps, load project conventions into context. Conventions ground generated output in project-specific coding standards.

**Resolution algorithm**:

1. **Project-wide conventions**: Read this skill's own frontmatter `conventions` field. For each filename, check `docs/conventions/{filename}` — if it exists, read it into context; if missing, skip silently.
2. **Per-type conventions** (per-type mode only): Read the active type's instruction file from `types/{type}.md`, extract its frontmatter `conventions` field. For each filename, check `docs/conventions/{filename}` — if it exists, read it; if missing, skip silently.
3. **Legacy mode**: Only project-wide conventions are loaded (no per-type instruction files to consult).

<HARD-RULE>
Convention loading is non-blocking. Missing convention files are silently skipped.
</HARD-RULE>

## Type Filter (--type)

This skill accepts an optional `--type <interface>` argument that filters script generation to a single test type. When `--type` is specified, the skill skips all other type groups entirely — no Fact Table verification, no locator mapping, no spec generation for non-matching types.

**Valid `--type` values**: the project's detected interface names (e.g., `tui`, `web-ui`, `api`, `cli`), obtained via `forge testing interfaces`.

**Validation**: If `--type` is specified but not found in the project's interfaces, the skill MUST error with a clear message listing the valid types. For example: `"invalid type: foo. Valid interfaces for this project: tui, api, cli"`.

**`--type` interaction with input mode**:

| Input mode | `--type` specified | `--type` omitted |
|------------|-------------------|-----------------|
| **Per-type** (per-type files found) | Load only `{type}-test-cases.md` matching the filter | Load all `{type}-test-cases.md` files sequentially |
| **Legacy** (single test-cases.md) | Filter test cases by type after reading (unchanged behavior) | Process all types (unchanged behavior) |

**Behavior summary**:

| Step | When `--type` matches | When `--type` does not match | Without `--type` |
|------|----------------------|------------------------------|-------------------|
| Step 1: Read & group test cases | Per-type: read single `{type}-test-cases.md`; Legacy: process only test cases of the specified type | Skip non-matching type groups | Process all types (unchanged behavior) |
| Step 1.5: Code Reconnaissance (Fact Table) | Build Fact Table for the specified type only | Skip Fact Table for non-matching types | Build Fact Table for all types |
| Step 2: Resolve Sitemap | Run if type is `web-ui` | Skip | Run if project has `web-ui` interface |
| Step 3: Map Locators | Run if type is `web-ui` | Skip | Run if project has `web-ui` interface |
| Step 3.5: Shared Infrastructure | **Always runs** regardless of `--type` | **Always runs** | Always runs |
| Step 4: Generate Spec Files | Generate spec files for the specified type only | Skip non-matching types | Generate for all types |

**Key rules**:
- Shared infrastructure (Step 3.5) **always runs** regardless of `--type` value — helpers, config, and auth-setup are shared across all types
- Without `--type`, behavior is completely unchanged: all types are generated as before
- The type filter is applied early in the pipeline at Step 1 (test case grouping), after project interfaces are resolved
- In per-type mode with `--type`, the matching `{type}-test-cases.md` is the only file read — no type grouping step is needed (file is already single-type)

## Prerequisites

Check previous stage artifacts. Abort and prompt user if missing:

**Test Case File Discovery** (first match wins):
1. **Per-type mode**: Glob `docs/features/<slug>/testing/*-test-cases.md`
2. **Legacy mode**: `docs/features/<slug>/testing/test-cases.md`

Per-type files take precedence over legacy. Missing? Prompt to run `/gen-test-cases`. `<slug>` from `forge feature`.

### Step Actionability Gate

If eval-test-cases report exists (`docs/features/<slug>/testing/eval/iteration-*.md`) and Step Actionability score < 200: **ABORT** with message to fix test cases first. If no report exists, proceed normally.

## Workflow

```
0. Resolve language + strategy → 1. Read test cases → 1.5. Code Reconnaissance → 3.5. Shared infrastructure → 4. Dispatch to type files
```

<HARD-RULE>
**Task Splitting Guard**: Global Setup Phase (Steps 0→1→1.5→3.5) runs as single pre-task FIRST. Only Step 4 may be parallelized. Each sub-task must include: "Shared infrastructure already configured. Auth: use storageState/cached token, NOT login() in beforeEach."
</HARD-RULE>

### Step 1: Read Test Cases

**Per-type mode**: For each `{type}-test-cases.md` (or only the file matching `--type`), parse test cases directly — file is already single-type. Parse TC ID, title, type, route, feature, pre-conditions, steps, expected result, priority.

**Legacy mode**: Read `docs/features/<slug>/testing/test-cases.md`. Parse each test case — extract TC ID, title, type, route, feature, pre-conditions, steps, expected result, priority. Group by type using the project's interfaces (e.g., `web-ui` → UI, `api` → API, `cli` → CLI).

**Type filter**: If `--type` was specified, keep only the group matching the specified type and discard all other groups. Only test cases of the specified type are processed in all subsequent steps. If no test cases match the specified type, emit a WARNING and proceed (shared infrastructure in Step 3.5 still runs).

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

**Credential caching**: Follow the rules in the active strategy's `generate.md` for framework-specific credential caching (e.g., storageState for browser-based UI, cached tokens for API). Login tests must invalidate cached credentials after completion.

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

Check `tests/e2e/`. If any shared file is missing, create from strategy templates. Strip `// VERIFY:` and `// TEMPLATE:` comments from copies. If helper file exists but missing exports needed by Auth Plan, merge missing exports (do NOT overwrite existing).

#### Auth Infrastructure

When `shared-auth-enabled: yes`:
1. Follow the strategy's `generate.md` for framework-specific auth setup
2. Playwright: generate `auth-setup.ts`, configure `playwright.config.ts` `projects` with `storageState`
3. Other strategies: follow `generate.md`

<HARD-RULE>
Auth-required tests MUST use shared auth (storageState/cached tokens), NOT per-test login in beforeEach. Playwright: calling `login()` in `beforeEach` is FORBIDDEN when `storageState` is configured.
</HARD-RULE>

<HARD-RULE>
**Shared file policy**: Create from templates if missing. Merge missing exports into helpers if exists. Do NOT modify other existing shared files (config, package manifest).
</HARD-RULE>

### Step 4: Dispatch to Type Files

For each active type, load the corresponding type instruction file and follow its instructions.

For each non-empty, verified type group, generate a test file from the corresponding template in the language's `templates/` directory:

- UI tests: language template for UI tests
- API tests: language template for API tests
- CLI tests: language template for CLI tests

**Follow the active strategy's `generate.md`** for all framework-specific patterns:
- Test runner imports, annotations, and lifecycle hooks
- Assertion library usage and conventions
- HTTP client for API tests
- Process execution for CLI tests
- Auth credential injection mechanism
- Import path conventions (staging: relative path to shared helpers)
- Test class/file naming conventions

Based on Step 1 auth classification, uncomment matching CONDITIONAL blocks in templates, remove non-matching blocks, then fill in test data. Replace example content with actual test cases from `test-cases.md`, keeping the same structure. Do not rewrite template structure from scratch.

**Type-to-file mapping** (hard-coded):

| Type | Instruction File |
|------|-----------------|
| ui | `types/ui.md` |
| tui | `types/tui.md` |
| mobile | `types/mobile.md` |
| api | `types/api.md` |
| cli | `types/cli.md` |

**Dispatch procedure** for each active type:
1. Verify type file exists. If missing: **HALT** with `"Type file missing: {path}. Cannot generate {type} test scripts. Available types: ui, tui, mobile, api, cli."`
2. Read type file, execute its instructions (reconnaissance, Fact Table keys, verification, generation patterns, antipattern guards)
3. Generate spec files following type file patterns and strategy's `generate.md`

If `--type` specified, only dispatch to the matching type file.

#### Integration Test Scripts

Integration test cases verify that a component embedded on an existing page is visible and renders correctly. Follow the strategy's `generate.md` for framework-specific patterns.

**Identification**: A test case is an integration test when:
- `Source` field contains "Placement + Integration Spec", **or**
- `Test ID` field contains "integration-"

Integration tests use the same framework and file structure as standard tests. They are grouped with UI test cases and emitted into the same test file.

**Script generation strategy** — for each integration test case:

1. **Locate the target page** by route (using sitemap resolution from Step 2)
2. **Locate the embedded component** using the locator strategy from `generate.md`
3. **Assert visibility**: use framework-specific visibility assertion
4. **Assert position** (if feasible): verify the component appears at expected position
5. **Verify data rendering**: assert expected text content is present within the component locator

#### beforeAll Safety

All setup blocks must: (1) wrap operations in try/catch, (2) explicit null/undefined check after assignment, (3) prefer reusing existing resources, (4) use sequential execution for dependent tests, (5) use retry for backend-dependent operations.

<HARD-RULE>
Uninitialized module-level variables cause misleading downstream failures. Try/catch + explicit check ensures the real failure is reported.
</HARD-RULE>

#### Antipattern Guard (Mandatory Pre-Emission Check)

Verify each generated test function does not match any forbidden pattern. Fix before writing.

These patterns are derived from common e2e test quality issues observed in automated test generation: recursive invocation causing process explosion, dead tests inflating coverage signal, vacuous truth providing false confidence, environment-dependent skips hiding regressions, duplicate tests wasting CI time, and static-file checks coupled to prose rather than behavior.

| # | Forbidden Pattern | What It Is | Why It's Harmful | What To Do Instead |
|---|---|---|---|---|
| 1 | **Recursive test invocation** | Calling `exec.Command("go", "test"` (or the equivalent test-runner invocation for any language) from within a test that belongs to the same package/module being tested | Causes process explosion — each recursion level spawns more processes. On Windows, orphaned children persist indefinitely, consuming 6GB+ RAM (see `gotcha-recursive-go-test-process-explosion.md`) | If a meta-test must verify "all tests pass", use a recursion guard: set an environment variable before spawning the subprocess and check it at the top of the calling test. Or exclude the meta-test via `-run` flag filtering |
| 2 | **Unconditional `t.Skip` (dead tests)** | `t.Skip("requires X")` with no environment-detection logic and no plan to create the precondition | Dead code that inflates coverage signal without verifying anything. The test count increases but no behavior is tested | Either implement the test with proper fixture setup (create the precondition), or do not generate the test at all. A test that always skips is worse than no test — it hides the gap |
| 3 | **Vacuous assertions** | Patterns like `if condition { assert.X(...) }` where the assertion may never execute on some code paths | Vacuous truth: the test passes because the assertion is never reached, not because the behavior is correct. This is the most harmful outcome — it provides false confidence (see `gotcha-e2e-test-quality-antipatterns.md`) | Every assertion must be reachable on every code path. If a condition must be checked, guard with `require` (fail-fast) not `assert` inside an `if`. Or: test must explicitly skip with a rationale when the condition is unmet |
| 4 | **Conditional skip without self-contained fixture** | Tests that skip based on environment state ("no pending tasks", "feature branch not found") instead of creating their own isolated world | Tests silently pass or skip depending on environment state. Results are non-deterministic — the same test can pass in one environment and skip in another, hiding real regressions | Every test that calls an external CLI or service must use `t.TempDir()` (or framework equivalent) and set up its own project structure. Fixtures must create the exact state the test needs — no reliance on pre-existing environment state |
| 5 | **Duplicate test functions across packages** | Generating a `func TestTC_*` name that already exists in another package under the same module (e.g., both `tests/e2e/` and `tests/e2e/features/slug/`) | Both copies run under `go test ./...`, doubling CI time. When they diverge (one fixed, one not), one always fails. This happens when graduation copies but does not remove the source (see `gotcha-e2e-test-quality-antipatterns.md`) | Before generating, scan existing test files in the module for matching `func TestTC_*` names. If a collision is found, either skip generating the duplicate or use a unique function name that includes the feature slug |
| 6 | **Static-file text grep tests** | Reading static source files (`.md`, `.go`, `.json`, `.yaml`) and asserting on text content via `assert.Contains(out, "some string")` | Tests documentation text, not runtime behavior. A typo fix in a markdown file breaks the test without any functional regression. Zero verification value — the test is coupled to prose, not behavior | Only test runtime behavior: invoke the CLI/API/UI and assert on outputs, exit codes, HTTP responses, or rendered content. Never read source files as test input |

**Validation procedure** — run after generating each test file:

1. **Recursion check**: Grep the generated file for test-runner invocations (`exec.Command("go", "test"`, `subprocess.run(["pytest"`, etc.). If found, verify a recursion guard is present. If not, reject.
2. **Skip check**: Find all `t.Skip` (or framework equivalent). Each must have an environment-detection rationale (e.g., `os.Getenv("CI") == ""`). If any skip has no rationale, reject.
3. **Vacuous assertion check**: Find all `if ... { assert` patterns. If the assertion is inside a conditional block without a corresponding `else { t.Skip/t.Fatal }`, reject.
4. **Fixture check**: For each test that invokes an external command or service, verify it creates its own isolated directory (`t.TempDir()` or equivalent). If it relies on the working directory or pre-existing state, reject.
5. **Duplicate check**: Compare generated function names against existing test files in the module. If any name already exists elsewhere, rename or skip.
6. **Static grep check**: Find all file-reading patterns (`os.ReadFile`, `ioutil.ReadFile`, `fs.readFileSync`, `open()`) in the generated code. If the read content is used in assertions (not just as test fixture input), reject.

#### Anti-Patterns (Forbidden in Generated Code)

Follow the anti-pattern rules in the active strategy's `generate.md`. Common anti-patterns across all languages:

<HARD-RULE>
**No fixed delays/sleeps.** These mask real timing issues, waste time when the app is fast, and fail when the app is slow. Every wait must be event-driven or use polling with timeout.
| # | Pattern | Harmful Because | Instead |
|---|---------|-----------------|--------|
| 1 | Recursive test invocation (spawning test runner within test) | Process explosion; 6GB+ RAM on Windows | Recursion guard (env var) or `-run` flag filtering |
| 2 | Unconditional `t.Skip` with no env detection | Dead code inflating coverage signal | Implement with fixture setup or don't generate |
| 3 | Vacuous assertions (`if { assert }` without else skip/fail) | Assertion never reached; false confidence | Every assertion reachable on every code path |
| 4 | Skip based on environment state, no self-contained fixture | Non-deterministic pass/skip | `t.TempDir()` + own project structure |
| 5 | Duplicate test function names across packages | Double CI time; divergent copies | Scan for collisions; use unique names with feature slug |
| 6 | Static-file text grep (assert on source file content) | Coupled to prose, not behavior | Test runtime behavior only (invoke + assert outputs) |

#### Generation HARD-RULEs

**Traceability**: Each test function includes traceability comment to PRD. **Skip empty groups**. **Empty result guard**: if zero test files generated, abort with diagnostic. **VERIFY markers**: resolve from Fact Table; if no value, keep `// VERIFY:` as-is. **Post-generation**: run `just e2e-compile` (fix errors), `just e2e-verify --feature <slug>` (check unresolved VERIFY), verify helper exports cover generated imports. **No fixed delays/sleeps. No hardcoded URLs. No debug output.** Selector strategy follows `generate.md`.

<EXTREMELY-IMPORTANT>
All framework-specific rules (test runner, assertion library, imports, HTTP client, process execution, anti-patterns) are defined in the active strategy's `generate.md` (loaded in Step 0). Read and follow those rules precisely.

- Use ONLY the framework specified in the strategy's `generate.md`
- **`@eN` refs must not appear in any generated test file.**
</EXTREMELY-IMPORTANT>

<HARD-RULE>
**Traceability**: Each test function must include a traceability comment ensuring complete traceability chain from test script to PRD.

**Skip empty groups**: If no test cases of a given type exist, or the project lacks that interface (Step 4 verification), skip generating that test file.

**Empty result guard**: If zero test files were generated (all groups empty or all interfaces undetected), abort with a clear message: "No test files generated — either the test case files have no testable cases, or all interface types were undetected. Re-run /gen-test-cases or verify project structure."

**VERIFY marker resolution**: Resolve all `// VERIFY:` comments using Fact Table values. If no Fact Table value exists, keep the `// VERIFY:` comment as-is.

**Post-generation compilation check**: Run `just e2e-compile` to verify generated files compile. Fix any compilation errors before completing.

**Post-generation VERIFY check**: Verify no unresolved `// VERIFY:` markers remain:
1. Run `just e2e-verify --feature <slug>` if the recipe exists in the Justfile
2. Otherwise: `grep -rn '// VERIFY:' tests/e2e/features/<slug>/` (adapt file extension to the detected language)

**Post-generation helper merge**: After generating all spec files, verify the language's helpers file exports cover all imports used by generated specs. If any import is missing from helpers, merge it from the template (do NOT overwrite existing exports).
</HARD-RULE>

## Output

All generated test files go to `tests/e2e/features/<feature>/` (staging area). After verification via `/run-e2e-tests`, use `/graduate-tests` to migrate them to the regression suite. Output file names are determined by the language's template mappings.

## Error Handling

| Situation | Action |
|-----------|--------|
| `.forge/config.yaml` missing and auto-detection fails | Ask user to configure `languages` in config.yaml |
| No per-type files AND `test-cases.md` missing | Abort with prompt to run `/gen-test-cases` |
| `sitemap.json` missing (web-ui tests) | Abort with prompt to run `/gen-sitemap` |
| eval-test-cases Step Actionability < 200 | Abort with prompt to fix `test-cases.md` first |
| Compilation fails post-generation | Fix generated code, re-run compile check |
| No test files generated (all empty groups) | Abort with clear diagnostic message |
| Source files not found for Fact Table | Mark as `UNKNOWN`, do not fabricate values |
| All required Fact Table keys UNKNOWN for a test type | Skip that test type with WARNING; suggest verifying source files |
| Helper file exists but missing needed exports | Merge missing exports from template into existing file |
| Language resolution fails | Ask user to configure `languages` in config.yaml |
| Unknown `--type` value | Error listing valid types. No generation. |
| Type file missing | Halt naming missing file and valid types |

## Related Skills

| Skill | Usage |
|-------|-------|
| `/gen-test-cases` | Generate test cases from PRD |
| `/gen-sitemap` | Generate and maintain sitemap.json |
| `/run-e2e-tests` | Execute test scripts and report results |
