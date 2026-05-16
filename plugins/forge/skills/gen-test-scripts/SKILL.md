---
name: gen-test-scripts
description: Generate executable e2e test scripts from test cases. Profile-aware: reads active test profile from .forge/config.yaml to determine test framework, templates, and conventions.
---

# Gen Test Scripts

Generate executable e2e test scripts from test cases.

**Core principle**: Every generated `test()` must be independently runnable, repeatable, and have explicit assertions. Test cases are input; scripts are output.

<HARD-GATE>
This skill ONLY writes to `tests/e2e/features/<feature>/` (staging area). It does NOT execute tests (handled by `/run-e2e-tests`).

**FORBIDDEN output paths** — agents MUST NOT write spec files to any of these, regardless of project convention:
- `tests/e2e/<feature>/` (directly under e2e root — this is a post-graduation location)
- `tests/e2e/<module>/` (any module directory — reserved for graduated scripts)
- Any path outside `tests/e2e/features/`

**Why this matters**: The `features/` prefix is a staging area that enforces a two-phase flow:
```
Phase 1: /gen-test-scripts → writes to tests/e2e/features/<feature>/
Phase 2: /graduate-tests   → classifies by functional module → moves to tests/e2e/<module>/
```
Bypassing the staging area skips functional module classification in `/graduate-tests`, scattering tests into wrong directories. Existing directories at `tests/e2e/<module>/` are POST-GRADUATION locations — do NOT copy that convention during generation.
</HARD-GATE>

## Step 0: Resolve Profile

1. **Resolve profile**: Run `forge profile` to get the active test profile(s). This reads `.forge/config.yaml`, falls back to project structure detection.
2. **On failure** (output shows `PROFILE: (none)`): ask the user to choose from known profiles (`web-playwright`, `go-test`, `maestro`, `java-junit`, `rust-test`, `pytest`). Run `forge profile set <name>` to persist their choice.
3. **Load profile manifest**: Run `forge profile get <profile-name> --manifest`.
4. **Load profile strategy**: Run `forge profile get <profile-name> --generate`.

Use the loaded profile manifest and strategy for all subsequent steps.

<HARD-RULE>
Do NOT silently default to any profile. If `forge profile` returns no result and the user cannot decide, abort the skill.
</HARD-RULE>

## Type Filter (--type)

This skill accepts an optional `--type <capability>` argument that filters script generation to a single test type. When `--type` is specified, the skill skips all other type groups entirely — no Fact Table verification, no locator mapping, no spec generation for non-matching types.

**Valid `--type` values**: the profile's declared capability names from its manifest (e.g., `tui` for go-test, `web-ui` for web-playwright, `api`, `cli`).

**Validation**: If `--type` is specified but not found in the profile's capabilities list, the skill MUST error with a clear message listing the valid types for the active profile. For example: `"invalid type: foo. Valid types for profile go-test: tui, api, cli"`.

**Behavior summary**:

| Step | When `--type` matches | When `--type` does not match | Without `--type` |
|------|----------------------|------------------------------|-------------------|
| Step 1: Read & group test cases | Process only test cases of the specified type | Skip non-matching type groups | Process all types (unchanged behavior) |
| Step 1.5: Code Reconnaissance (Fact Table) | Build Fact Table for the specified type only | Skip Fact Table for non-matching types | Build Fact Table for all types |
| Step 2: Resolve Sitemap | Run if type is `web-ui` | Skip | Run if profile has `web-ui` capability |
| Step 3: Map Locators | Run if type is `web-ui` | Skip | Run if profile has `web-ui` capability |
| Step 3.5: Shared Infrastructure | **Always runs** regardless of `--type` | **Always runs** | Always runs |
| Step 4: Generate Spec Files | Generate spec files for the specified type only | Skip non-matching types | Generate for all types |

**Key rules**:
- Shared infrastructure (Step 3.5) **always runs** regardless of `--type` value — helpers, config, and auth-setup are shared across all types
- Without `--type`, behavior is completely unchanged: all types are generated as before
- The type filter is applied early in the pipeline at Step 1 (test case grouping), after profile capabilities are resolved

## Prerequisites

Check previous stage artifacts. Abort and prompt user if missing:

| Artifact | Missing prompt |
|----------|----------------|
| `docs/features/<slug>/testing/test-cases.md` | Run `/gen-test-cases` first |
| `docs/sitemap/sitemap.json` (only when profile has `web-ui` capability) | Run `/gen-sitemap` first |

`<slug>` is the current feature name, obtained via `forge feature` command. `docs/sitemap/sitemap.json` is a project-level file (one per application), not isolated per feature.

### Step Actionability Gate

Check whether an eval-test-cases report exists for the current feature. If it does, parse the Step Actionability score and enforce the blocking threshold.

**Check procedure**:

1. List evaluation reports: `ls docs/features/<slug>/testing/eval/iteration-*.md 2>/dev/null`
2. If no reports exist, proceed normally (backward compatible — eval is optional)
3. If reports exist, read the latest iteration report (highest N)
4. Find the Step Actionability dimension score in the report's dimension table
5. If Step Actionability < 200: **ABORT** with the message below

**Abort message**:
```
ABORT: Step Actionability score is below the blocking threshold (score < 200).

The test cases in test-cases.md have insufficient step-level detail for reliable
script generation. Low Step Actionability means test steps are vague, missing
concrete actions, or lack sufficient detail for automated translation.

Action required: Run /eval-test-cases, review the Step Actionability dimension
feedback, and fix test-cases.md to improve step-level specificity before
re-running /gen-test-scripts.
```

<HARD-RULE>
If eval-test-cases report exists and Step Actionability < 200, gen-test-scripts MUST abort. This gate prevents generating scripts from poorly-specified test cases, which historically leads to 3-5 rounds of agent-driven fixes. The threshold of 200 (out of 250) is the minimum for unambiguous step-to-script translation.
</HARD-RULE>

### Sitemap (web-ui capability only)

When the active profile has `web-ui` capability, the sitemap provides supplementary page structure for UI tests: `docs/sitemap/sitemap.json` (generated by `/gen-sitemap`, full example: `plugins/forge/references/shared/sitemap.json`).

**Note**: The sitemap is a secondary reference for page structure and route matching. All locators are derived from source code via the Fact Table (Step 1.5). The sitemap does NOT provide locators driven by test-case Element IDs.

**Key page structure fields**:

| Field | Usage |
|-------|-------|
| `layout.elements[].role` + `name` | Layout-level shared elements (navbar, etc.), available for all pages wrapped by `layout.wraps` |
| `elements[].role` + `name` | Semantic role+name pair — supplementary confirmation of Fact Table data |
| `elements[].level` | Heading level for page structure understanding |
| `elements[].label` | Accessible label — supplementary confirmation |
| `elements[].placeholder` | Placeholder text — supplementary confirmation |
| `states[].trigger` | Click trigger element to enter dynamic state |
| `states[].elements` | Elements within dynamic states |

Test cases match sitemap pages via the `Route` field. The `Element` field from test-cases.md is NOT used as a locator source — all locator derivation comes from the Fact Table built in Step 1.5 Code Reconnaissance.

## When to Use

**Trigger:**

- User asks to "generate test scripts" or "create e2e scripts"
- User provides `/gen-test-scripts` command
- After `/gen-test-cases` has produced `testing/test-cases.md`

**Skip:**
- `tests/e2e/features/<slug>/` already exists and contains generated scripts (re-run only if explicitly requested)

## Workflow

<HARD-RULE>
**Task Splitting Guard**: When gen-test-scripts is split into parallel sub-tasks (e.g., by output file or test case range), the following MUST hold:

1. **Global Setup Phase** (Steps 0→1→1.5→3.5) runs as a single pre-task FIRST. This produces:
   - Auth Plan (from Step 1 classification)
   - Fact Table (from Step 1.5 reconnaissance)
   - Shared infrastructure (from Step 3.5: auth-setup.ts, playwright.config.ts, helpers.ts, config files)
2. **Only Step 4** (spec file generation) may be parallelized into sub-tasks, each handling a subset of test cases
3. Each sub-task's description MUST include: "Shared infrastructure already configured in pre-task. Auth: use storageState (Playwright) / cached token (API), NOT login() in beforeEach."

**Why**: Steps 0-3.5 produce shared artifacts that ALL spec files depend on. Splitting them causes each subagent to skip shared infrastructure and fall back to per-test login (see `docs/lessons/gotcha-split-task-missing-shared-setup.md`).
</HARD-RULE>

```
0. Resolve profile → 1. Read test cases → 1.5. Code Reconnaissance → 2. Resolve sitemap (web-ui) → 3. Map locators (web-ui) → 3.5. Ensure shared infrastructure (auth + helpers + config) → 4. Generate spec files
```

### Step 1: Read Test Cases

Read `docs/features/<slug>/testing/test-cases.md`. Parse each test case — extract TC ID, title, type, route, feature, pre-conditions, steps, expected result, priority. Group by type using the profile's capabilities (e.g., `web-ui` → UI, `api` → API, `cli` → CLI).

**Type filter**: If `--type` was specified, keep only the group matching the specified type and discard all other groups. Only test cases of the specified type are processed in all subsequent steps. If no test cases match the specified type, emit a WARNING and proceed (shared infrastructure in Step 3.5 still runs).

#### Auth Classification

For each test case, classify by authentication requirements:

| Category | Detection Rules | Generation Strategy |
|----------|----------------|---------------------|
| **login-test** | `Target` matches auth/login routes (e.g., `ui/login`, `api/auth`, `api/token`) | No shared auth. Generate independent login/logout logic. Must invalidate cached credentials after to prevent stale state. |
| **auth-required-test** | `Pre-conditions` contain "authenticated", "logged in", or `Target` implies protected resource | Use **cached shared auth** — credentials acquired once, reused across all tests |
| **public-test** | No auth-related pre-conditions, target is a public resource | No auth needed |
| **custom-auth-test** | `Pre-conditions` mention "API key", "X-API-Key", "OAuth", "session cookie", or other non-Bearer auth patterns | Detect auth mechanism from codebase during Step 1.5. Generate custom auth setup |

**Classification priority** (evaluated top-down, first match wins):
1. `custom-auth-test` — takes priority over `auth-required-test` when both match (e.g., "authenticated via API key" → custom-auth, not auth-required)
2. `login-test`
3. `auth-required-test`
4. `public-test`

Count each category to decide whether to enable shared auth (enabled when `auth-required-test` exists).

**Auth Plan Output**: Record the classification result. This is the contract consumed by Step 3.5 (auth infrastructure) and Step 4 (spec generation):
```
Auth Plan:
  auth-required-test: N (routes: ...)
  login-test: N (routes: ...)
  public-test: N
  custom-auth-test: N (mechanism: ...)
  shared-auth-enabled: yes/no
```

**Credential caching**: Follow the rules in the active profile's `generate.md` for framework-specific credential caching (e.g., storageState for browser-based UI, cached tokens for API). Login tests must invalidate cached credentials after completion.

<HARD-RULE>
Login tests and authenticated tests must not be mixed in the same `describe` block.
</HARD-RULE>

### Step 1.5: Code Reconnaissance (Build Fact Table)

Read actual source code files to extract ground-truth values. **Never guess or assume values** — every value in a generated script must come from the test-cases.md input or the Fact Table built here.

**Type filter**: If `--type` is specified, skip Fact Table verification for non-matching types. Only build Fact Table entries required by the specified type's test cases.

**Check test-cases.md** for Route Validation results (warnings from gen-test-cases Step 3.5). Use corrected routes where available.

**Required reads** (adapt to project structure):

| Source | What to extract | Discovery guidance |
|--------|----------------|---------------------|
| Router files | Route paths, path parameters, middleware bindings | Search for route registration patterns, configuration files, and entry points using Grep. Look for URL path strings, HTTP method bindings, and path parameter definitions. |
| Config files | API port, base path prefix, auth credentials | Search for config/settings files (`.env`, `config.*`, `settings.*`). Look for port numbers, base URLs, and credential variable names. |
| API handlers | Request/response schemas, status codes, validation rules | Search for request handler functions and response definitions. Look for status code usage, input validation, and response body shaping. |
| Auth implementation | Login endpoint path, token field name, header format | Search for authentication/authorization modules. Look for login endpoints, token generation/parsing, and header middleware. |
| CLI entry points | Command names, flag names, output formats | Search for command registration and argument parsing. Look for CLI framework usage (e.g., cobra, commander, click) and output formatting. |
| Frontend UI components | `data-testid` values, component structure, dynamic testid patterns | Search frontend source for `data-testid` attributes. Grep for `data-testid` in component files (`.tsx`, `.vue`, `.svelte`). Note static vs dynamic patterns (e.g., `map-card-${id}` vs `map-card`). |

**Build the Fact Table**: After reading, record verified facts with source citations:

```markdown
## Fact Table
| Key | Value | Source |
|-----|-------|--------|
| API_PORT | 8080 | backend/config.yaml:3 |
| AUTH_ENDPOINT | POST /v1/auth/login | internal/handler/router.go:42 |
| TESTID_MAP_CARD | map-card-${bizKey} | frontend/src/pages/MilestonesPage.tsx:110 |
```

<HARD-RULE>
- Every value in the Fact Table must cite the source file and line number.
- If a source file cannot be found, note it as `UNKNOWN`. Do not fabricate values.
- The agent must use Fact Table values when generating spec files in Step 4. When a Fact Table value contradicts a template placeholder, the Fact Table wins.
- All `// VERIFY:` markers in templates must be resolved using Fact Table values.
</HARD-RULE>

#### Fact Table Completeness Gate

After building the Fact Table, validate completeness per test type before proceeding. For each test type that has test cases in test-cases.md:

| Test type | Required Fact Table keys (all must be non-UNKNOWN) |
|-----------|-----------------------------------------------------|
| UI | `FRONTEND_BASE` or `TESTID_*` entries |
| API | `API_PORT` or route path entries (`AUTH_ENDPOINT`, etc.) |
| CLI | At least one CLI command name entry |

**If all required keys for a test type are UNKNOWN**: skip that test type, emit a WARNING explaining why, and suggest the user verify the relevant source files exist before re-running. Do NOT generate tests for that type — UNKNOWN values would force the agent to guess, reproducing the anti-pattern this gate prevents.

**If some keys are non-UNKNOWN**: proceed for that type. Individual UNKNOWN keys are acceptable (e.g., `TESTID_MAP_CARD` unknown but `TESTID_BTN_CONFIRM` known) — only skip when *all* required keys for the entire type are UNKNOWN.

### Step 2: Resolve Sitemap (web-ui capability only)

**Only execute when the active profile has `web-ui` capability AND UI-type test cases exist.**

Read `docs/sitemap/sitemap.json`. For each route referenced in test cases:

1. Match `Route` field to `sitemap.json` `pages[].route` (dynamic routes like `/tasks/:id` match specific paths like `/tasks/123`)
2. Collect page structure data for matched pages (base elements + related states) as supplementary reference
3. **If sitemap has `layout` field**: for each route wrapped by `layout.wraps`, merge `layout.elements` as available page structure elements

The sitemap serves as a secondary reference for understanding page structure. Locators are NOT derived from sitemap element data keyed by test-case Element IDs — all locator derivation comes from the Fact Table (Step 1.5).

If a route referenced in test cases does not exist in sitemap, **emit a WARNING** listing the missing routes and suggest re-running `/gen-sitemap`. Proceed using the Fact Table from Step 1.5 for the missing routes — do not abort.

### Step 3: Map Locators (web-ui capability only)

**Only execute when the active profile has `web-ui` capability AND the type filter (if specified) includes a UI-related type (`web-ui` or `tui`).** If `--type` is `api` or `cli`, skip locator mapping entirely.

Locator strategy is owned by this step — `generate.md` provides only the framework-specific syntax (how to write a Playwright/Go/pytest locator), not the decision of *which* locator to use.

**Locator priority** (all locator types, including integration tests):

1. **Fact Table data** (from Step 1.5 Code Reconnaissance): `TESTID_*` entries, semantic roles, labels extracted from source code → translate via `generate.md`. For dynamic testids (e.g., `map-card-${bizKey}`), use prefix selector (`page.locator('[data-testid^="map-card-"]')`)
2. **Sitemap page structure** (from Step 2): role+name, label, placeholder — used as supplementary confirmation of Fact Table data only, NOT as primary locator source → translate via `generate.md`
3. **Semantic inference**: heading/section title, aria-label from source code → translate via `generate.md`

<HARD-RULE>
Never guess `data-testid` values. Every testid locator must come from a Fact Table `TESTID_*` entry. If the Fact Table has no testid entry for the needed element, use semantic locators (role, label, text) instead.
</HARD-RULE>

<HARD-RULE>
**Source-code-first locator derivation**: Derive all locators from source code (Fact Table). Do NOT reference any testid, selector, or locator from test-cases.md — test-cases provides scenario context only. The Fact Table built in Step 1.5 is the sole primary locator source. The sitemap (Step 2) is supplementary confirmation, not a locator source driven by test-case Element IDs.
</HARD-RULE>

For test steps within dynamic states: first click the trigger element's locator, then map locators for in-state elements. Build an in-memory mapping table for use in Step 4.

**Integration test component locators**: Use the same priority chain above. Locate the embedded component via Fact Table data (source code) first, then confirm with sitemap page structure, then semantic inference.

### Step 3.5: Ensure Shared Infrastructure

<PRINCIPLE>
**Shared infrastructure first.** Before generating any test files, ensure that shared dependencies are complete and functional. Test files depend on these shared files via import — if they are missing or incomplete, all tests will fail at the import stage. Downstream skills (`/run-e2e-tests`, `/graduate-tests`) follow the same principle.

**This step always runs regardless of `--type` value.** Shared infrastructure (helpers, config, auth-setup) is used by all test types and MUST be present before any spec files are generated. The idempotent create-only-if-missing logic ensures safe concurrent execution across per-type tasks.
</PRINCIPLE>

Check if shared infrastructure exists at `tests/e2e/`:

```bash
ls tests/e2e/ 2>/dev/null
```

**If any shared file is missing**, create it from the corresponding template (paths from profile manifest `templates.*`):

- Helper file: copy from `{profile-templates-dir}/{manifest.templates.helpers}` (only if missing). When copying, strip all `// VERIFY:` and `// TEMPLATE:` comments — these are generation-time markers that should not appear in runtime files.
- Config files: copy from manifest template references (only if missing)
- Follow `generate.md` for additional shared infrastructure requirements specific to the framework

**If helper file already exists**: check whether it exports all symbols required by the Auth Plan (e.g., `ensureAuthState`, `getApiToken`). If symbols are missing, add the missing exports AND all their private dependencies from the template. Do NOT overwrite existing exports — merge.

#### Auth Infrastructure

Based on the Auth Plan from Step 1, configure auth infrastructure for auth-required tests.

**When `shared-auth-enabled: yes`** (auth-required-test count > 0):

1. Follow the active profile's `generate.md` for framework-specific auth setup
2. For playwright profile:
   - Generate `tests/e2e/auth-setup.ts` from template (if not already present)
   - Configure `tests/e2e/playwright.config.ts`: uncomment the `projects` section with `setup` + `authenticated` projects using `storageState`
3. For other profiles: follow `generate.md` auth setup instructions
4. Verify `helpers.ts` exports all auth-related symbols needed by spec files

**When `shared-auth-enabled: no`**: Skip auth infrastructure setup.

<HARD-RULE>
Auth-required tests MUST use the shared auth mechanism (e.g., storageState for Playwright, cached tokens for API), NOT per-test login in beforeEach. The auth infrastructure generated here provides shared credentials — spec files in Step 4 MUST reference it.

Specifically for Playwright: when the `projects` section in playwright.config.ts is configured with `storageState`, calling `login()` or `loginViaUI()` in `beforeEach` is FORBIDDEN for auth-required tests.
</HARD-RULE>

<HARD-RULE>
**Shared file policy**: Shared files at `tests/e2e/` are shared across all features.
- If they do not exist: create them from templates.
- If helper file exists: merge missing exports from template into it (do NOT overwrite existing exports).
- Other shared files (config, package manifest, etc.): if they exist, DO NOT modify.
- Only test files are written per-feature.
</HARD-RULE>

### Step 4: Generate Spec Files

**Verify project interfaces before generating**: For each type group from Step 1, confirm the project actually exposes that interface. Build/test/lint commands are developer tooling, not a CLI product interface.

**Type filter**: If `--type` is specified, only generate spec files for the specified type. Skip generating spec files for all other types. Without `--type`, generate spec files for all non-empty, verified type groups as before.

| Type | Verification method | Evidence of product interface |
|------|---------------------|-------------------------------|
| UI | Fact Table has `TESTID_*` or `FRONTEND_*` entries (from Step 1.5) | Frontend source code was found and read |
| API | `grep -rn "router\|handler\|endpoint\|HandleFunc\|app.get\|app.post" --include='*.go' --include='*.ts' --include='*.js' .` | HTTP handler registration patterns found |
| CLI | `grep '"bin"' package.json` or `ls cmd/` or `grep -rn "cobra.Command" --include='*.go' .` | CLI entry point or command framework detected |

For each non-empty, verified type group, generate a test file from the corresponding template in the profile's `templates/` directory (paths from `manifest.templates.*`):

- UI tests: `{profile-templates-dir}/{manifest.templates.test-file-ui}`
- API tests: `{profile-templates-dir}/{manifest.templates.test-file-api}`
- CLI tests: `{profile-templates-dir}/{manifest.templates.test-file-cli}`

**Follow the active profile's `generate.md`** for all framework-specific patterns:
- Test runner imports, annotations, and lifecycle hooks
- Assertion library usage and conventions
- HTTP client for API tests
- Process execution for CLI tests
- Auth credential injection mechanism
- Import path conventions (staging: relative path to shared helpers)
- Test class/file naming conventions

Based on Step 1 auth classification, uncomment matching CONDITIONAL blocks in templates, remove non-matching blocks, then fill in test data. Replace example content with actual test cases from `test-cases.md`, keeping the same structure. Do not rewrite template structure from scratch.

#### Integration Test Scripts

Integration test cases verify that a component embedded on an existing page is visible and renders correctly. Follow the profile's `generate.md` for framework-specific patterns.

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

All setup blocks (`beforeAll`, `@BeforeAll`, `setUp`, etc.) must follow defensive patterns to prevent cascade failures where module-level variables stay uninitialized:

1. **Wrap each operation in try/catch** with error logging identifying which step failed
2. **Explicit null/undefined/None check** — after assignment, validate the value is set
3. **Prefer reusing existing resources** over creating new ones — fewer operations means fewer failure modes
4. **Use sequential execution** when tests have dependencies, instead of shared mutable state
5. **Use retry mechanisms** for backend-dependent operations (token acquisition, resource creation) to handle transient failures

<HARD-RULE>
When setup blocks throw and module-level variables stay uninitialized, all downstream tests fail with misleading errors (e.g., `undefined` in URL paths → 404). The try/catch + explicit check pattern ensures the *real* failure is reported, not a secondary symptom.
</HARD-RULE>

#### Antipattern Guard (Mandatory Pre-Emission Check)

Before emitting each generated test function, verify it does not match **any** of the 6 forbidden patterns below. If a pattern matches, fix the generated code before writing it. This guard is mandatory — it is not optional quality advice.

Authoritative sources: `docs/lessons/gotcha-e2e-test-quality-antipatterns.md` and `docs/lessons/gotcha-recursive-go-test-process-explosion.md`.

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

Follow the anti-pattern rules in the active profile's `generate.md` strategy. Common anti-patterns across all profiles:

<HARD-RULE>
**No fixed delays/sleeps.** These mask real timing issues, waste time when the app is fast, and fail when the app is slow. Every wait must be event-driven or use polling with timeout.

**No hardcoded URLs.** Use config/environment variables for all endpoint URLs.

**No debug output in generated code.** Only assertions — no print/console.log for debugging.

**Selector strategy**: Follow `generate.md` for selector priority. Never bind to internal implementation details (CSS classes, DOM structure) when semantic selectors are available.
</HARD-RULE>

<EXTREMELY-IMPORTANT>
All framework-specific rules (test runner, assertion library, imports, HTTP client, process execution, anti-patterns) are defined in the active profile's `generate.md` strategy file (loaded in Step 0). Read and follow those rules precisely.

- Use ONLY the framework specified in the profile's `generate.md`
- **`@eN` refs must not appear in any generated test file.**
</EXTREMELY-IMPORTANT>

<HARD-RULE>
**Traceability**: Each test function must include a traceability comment ensuring complete traceability chain from test script to PRD.

**Skip empty groups**: If no test cases of a given type exist, or the project lacks that interface (Step 4 verification), skip generating that test file.

**Empty result guard**: If zero test files were generated (all groups empty or all interfaces undetected), abort with a clear message: "No test files generated — either test-cases.md has no testable cases, or all interface types were undetected. Re-run /gen-test-cases or verify project structure."

**VERIFY marker resolution**: Resolve all `// VERIFY:` comments using Fact Table values. If no Fact Table value exists, keep the `// VERIFY:` comment as-is.

**Post-generation compilation check**: Run `just e2e-compile` to verify generated files compile. Fix any compilation errors before completing.

**Post-generation VERIFY check**: Verify no unresolved `// VERIFY:` markers remain:
1. Run `just e2e-verify --feature <slug>` if the recipe exists in the Justfile
2. Otherwise: `grep -rn '// VERIFY:' tests/e2e/features/<slug>/` (adapt file extension to profile)

**Post-generation helper merge**: After generating all spec files, verify `helpers.ts` exports cover all imports used by generated specs. If any import is missing from helpers, merge it from the template (do NOT overwrite existing exports).
</HARD-RULE>

## Output

All generated test files go to `tests/e2e/features/<feature>/` (staging area). After verification via `/run-e2e-tests`, use `/graduate-tests` to migrate them to the regression suite. Output file names are determined by the profile's manifest template mappings.

## Error Handling

| Situation | Action |
|-----------|--------|
| `.forge/config.yaml` missing and auto-detection fails | Ask user to select a profile |
| `docs/features/<slug>/testing/test-cases.md` missing | Abort with prompt to run `/gen-test-cases` |
| `sitemap.json` missing (web-ui tests) | Abort with prompt to run `/gen-sitemap` |
| eval-test-cases Step Actionability < 200 | Abort with prompt to fix `test-cases.md` first |
| Compilation fails post-generation | Fix generated code, re-run compile check |
| No test files generated (all empty groups) | Abort with clear diagnostic message |
| Source files not found for Fact Table | Mark as `UNKNOWN`, do not fabricate values |
| All required Fact Table keys UNKNOWN for a test type | Skip that test type with WARNING; suggest verifying source files |
| Helper file exists but missing needed exports | Merge missing exports from template into existing file |

## Related Skills

| Skill             | Usage                                   |
| ----------------- | --------------------------------------- |
| `/gen-test-cases` | Generate test cases from PRD            |
| `/gen-sitemap`    | Generate and maintain sitemap.json      |
| `/run-e2e-tests`  | Execute test scripts and report results |
