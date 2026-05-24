---
name: gen-test-scripts
description: Generate executable test scripts from Contract specifications. Journey-driven: generates test code with @feature tags directly into tests/<journey>/.
---

# Gen Test Scripts

Generate executable test scripts from Contract specifications.

**Core principle**: Tests are generated per Journey, not per interface type. Each Journey step's Contract defines the assertions. Tests go directly to `tests/<journey>/` with `@feature` tags.

## Pipeline Position

```
gen-journeys -> gen-contracts -> gen-test-scripts -> run-tests
                                      ^^^ YOU ARE HERE
```

Input: Contract specifications (from gen-contracts) + Fact Table (from code reconnaissance).
Output: Executable test code with `@feature` tags in `tests/<journey>/`.

## Prerequisites

Check previous stage artifacts. Abort and prompt user if missing:

| Artifact | Missing prompt |
|----------|----------------|
| At least one Contract file in `docs/features/<slug>/testing/<journey>/contracts/` | Run `/gen-contracts` first |
| Eval report for all Contracts (`testing/<journey>/.eval-report.md`) | Run `/eval --type contract` first. **Blocker**: do not proceed if any Contract scored below target. |

## Step 0: Load Convention Files

Load test framework knowledge from Convention files (no Profile/CLI dependency).

### 0.1 Discover Convention Files

Load test framework knowledge from Convention files using a two-level index mechanism.

1. Read `docs/conventions/testing/index.md` — this index file lists all available Conventions with name, description, and applicability conditions.
2. Based on the project's language/framework context, select the matching Convention from the index.
3. Load the selected Convention file from `docs/conventions/testing/<convention>.md`.
4. If `index.md` does not exist, proceed to auto-detection (Step 0.2).

<HARD-RULE>
Do NOT use `domains` frontmatter filtering. Selection is based on index.md descriptions and project context, with LLM autonomous judgment.
</HARD-RULE>

### 0.2 Resolve Target Framework

Determine the target framework from the available signals:

1. **Convention file match**: If a Convention was loaded from index.md, use its framework declaration.
2. **Existing test file scan**: Scan `tests/` for file patterns to confirm the Convention's framework matches the project.
3. **User specification**: If signals are ambiguous, ask the user which framework to use.
4. **No Convention found**: Proceed with LLM defaults + Code Reconnaissance (Step 1). Output hint: "No test Convention files found in `docs/conventions/testing/`. Generation will use LLM defaults. Run `/forge:test-guide` to create one."

<HARD-RULE>
If no Convention files are found and no framework can be detected from existing test files, ask the user which framework to use. Do NOT silently default.
</HARD-RULE>

### 0.3 Validate Convention Content

For the loaded Convention file, check required sections: `framework`, `discovery`, `structure`, `assertions`.

- **Missing required section**: Log warning listing missing sections. Proceed with LLM defaults for that section's area. Example: "Convention file `<path>` is missing sections: discovery, assertions. Using LLM defaults for those sections."
- **Invalid section content** (e.g., empty framework name): Treat as missing. Log warning.

Use the loaded Convention content for all framework-specific rules in subsequent steps.

## Step 0.5: Surface Detection

Determine the project's interface surface type to drive per-surface generation strategy.

### 0.5.1 Read Surface Configuration

Read `.forge/config.yaml` from the project root and extract the `surface` field.

```bash
forge config get surface
```

| Result | Action |
|--------|--------|
| Surface type returned (e.g., `cli`, `tui`, `webui`, `mobile`, `api`) | Use this as the active surface type for generation |
| No config file or field missing | Proceed to auto-detection (Step 0.5.2) |

### 0.5.2 Auto-Detection Fallback

If `.forge/config.yaml` does not contain a `surface` field, infer the surface type from code reconnaissance signals in Step 1. Use the Verification Method defined in each type file (`types/cli.md`, `types/api.md`, etc.) to probe for interface indicators.

Priority: config value > auto-detection. If auto-detection is ambiguous (multiple types detected), ask the user which surface type to prioritize.

### 0.5.3 Surface Strategy Application

The detected surface type determines the **test ratio strategy** for Step 3 generation:

| Surface Type | Contract : Journey Ratio | Key Generation Constraint |
|-------------|--------------------------|---------------------------|
| CLI | ≥ 80% Contract | Subprocess execution model, binary isolation, environment hermeticity |
| TUI | ≥ 80% Contract | Terminal I/O testing, non-interactive stdin pipe, ANSI sanitization |
| WebUI | Balanced 50/50 | Convention-defined browser framework, session reuse, network interception |
| API | Balanced 50/50 | HTTP client testing, status code coverage, content-type verification |
| Mobile | Best-effort | Maestro YAML skeleton + deep link tests, complex scenarios marked `manual-only` |

**Contract test ratio formula**: `Contract test functions / (Contract test functions + Journey smoke test functions) × 100%`

<HARD-RULE>
The surface type determines generation strategy — test ratio, execution model, and assertion patterns. Type-specific Golden Rules (from `types/<type>.md`) take precedence over generic generation rules. Convention provides framework implementation details, Surface type provides strategy constraints. These two are orthogonal and merged at generation time.
</HARD-RULE>

## Step 1: Code Reconnaissance (Build Fact Table)

Read source code to extract ground-truth values for semantic descriptor resolution AND test framework patterns.

### 1.1 Framework Reconnaissance

Scan existing test files for framework-specific patterns to supplement or validate Convention:

| Source | What to extract | Purpose |
|--------|-----------------|---------|
| Test file names | File pattern (`*_test.go`, `*.test.ts`) | Confirm Convention's `file-pattern` |
| Test file imports | Assertion library, test runner imports | Confirm Convention's `Assertion` section |
| Build tags / markers | Tag syntax (`//go:build e2e`, `@pytest.mark.e2e`) | Confirm Convention's `Tags` section |
| Test function signatures | Naming pattern, parameter types | Infer naming conventions |
| Test helper files | Utility patterns, fixture setup | Infer project-specific test patterns |

If no test files exist and no file signals are recognizable, Reconnaissance produces an empty Fact Table for framework columns. This is expected cold-start behavior — proceed with Convention alone, or LLM defaults if no Convention.

### 1.2 Domain Reconnaissance

Extract ground-truth values from application source code for semantic descriptor resolution:

| Source | What to extract |
|--------|-----------------|
| CLI entry points | Command names, flag names, output format strings |
| API handlers | Request/response schemas, status codes |
| TUI components | Model fields, View output patterns |
| Config files | Ports, base paths, auth mechanisms |
| Auth implementation | Login endpoint, token field, header format |

### 1.3 Build Fact Table

Combine all reconnaissance into a single Fact Table with source citations:

```markdown
## Fact Table
| Key | Value | Source |
|-----|-------|--------|
| CLI_TASK_CLAIM_OUTPUT | claimed task <task_id> | internal/cmd/claim.go:42 |
| CLI_FEATURE_CREATE_OUTPUT | Feature <slug> created successfully | internal/cmd/feature.go:45 |
| TEST_FRAMEWORK | go-testing | tests/e2e/step1_test.go (import analysis) |
| TEST_ASSERTION_LIB | testify/assert | tests/e2e/step1_test.go:3 (import) |
| TEST_BUILD_TAG | //go:build e2e | tests/e2e/step1_test.go:1 |
```

<HARD-RULE>
- Every Fact Table value must cite source file and line number. Unknown sources -> `UNKNOWN`. Do not fabricate.
- Fact Table values drive semantic descriptor to regex conversion. All `// VERIFY:` markers must be resolved using Fact Table values.
- When Reconnaissance finds signals that conflict with Convention -> Convention wins, but log the conflict for user awareness.
</HARD-RULE>

### 1.4 Semantic Descriptor to Regex Conversion

Contract Output dimensions use semantic descriptors (natural language), not regex. This step converts them to precise regex patterns:

1. For each Outcome's Output dimension, look up matching Fact Table entries.
2. Convert the Fact Table value to a regex pattern:
   - Placeholder tokens like `<task_id>` become named capture groups: `(?P<task_id>[\w-]+)`
   - Literal text is regex-escaped.
3. If no Fact Table match is found, keep the original descriptor as a `// VERIFY:` marker.

Example pipeline:
```
Semantic descriptor: "success confirmation containing feature-slug"
  -> Fact Table lookup: CLI_FEATURE_CREATE_OUTPUT = "Feature my-feature created successfully"
  -> Generated regex: Feature\s+([\w-]+)\s+created\ successfully
```

## Step 2: Read Contract Specifications

**Input discovery** — find the Contract files for the target Journey:

1. Glob `docs/features/<slug>/testing/<journey>/contracts/step-*.md` for the target Journey.
2. If no Journey is specified, ask the user which Journey to generate tests for.
3. Parse each Contract file to extract: Journey name, Step number, Action, Outcomes (Preconditions/Input/Output/State/Side-effect/Invariants).

<HARD-RULE>
**Batch generation**: Generate tests for one Journey at a time (happy path + edge cases). If the user wants multiple Journeys, process them sequentially.
</HARD-RULE>

### Contract Parsing

Each Contract file has this structure:

```markdown
---
journey: "task-lifecycle"
step: 2
step-action: "forge task claim"
---
# Contract: task-lifecycle / Step 2: forge task claim

## Outcome "success"
- Preconditions: "feature exists; at least one task available"
- Input: no positional args; no flags
- Output: stdout contains "claimed task <task_id>", exit code 0
- State: tasks/<task_id>/status -> "in_progress"; index.json updated
- Side-effect: none

## Outcome "no-tasks-available"
- Preconditions: "feature exists; no tasks available for claiming"
- Input: no positional args; no flags
- Output: stderr contains "no tasks available", exit code 1
- State: unchanged

## Journey Invariants
- feature_slug consistent across all steps
- task_id stable once assigned
```

<HARD-RULE>
**Single Journey per invocation**: Do not attempt to process multiple Journeys in one gen-test-scripts invocation. If Contracts span multiple Journeys, abort and ask the user to specify which Journey to generate.
</HARD-RULE>

## Step 2.5: Load Type Rules

After reading Contract files (Step 2) and before generating test code (Step 3), load type-specific Golden Rules that constrain generation.

<HARD-RULE>
types/ Golden Rules define non-overridable principle constraints. Convention provides framework implementation details. When both cover the same aspect, Golden Rules' principles take precedence, Convention's implementation details supplement areas Golden Rules don't cover.
</HARD-RULE>

<HARD-RULE>
Reconnaissance Hints in type files are discovery aids only. Information discovered via Hints must be converted to Fact Table values, not used directly in generation instructions.
</HARD-RULE>

### 2.5.1 Extract Interface Types from Contracts

Read each Contract file parsed in Step 2 and collect all interface types referenced in Actions and Outcomes:

1. Examine `step-action` fields for interface indicators (CLI commands, HTTP methods, UI interactions, TUI rendering, mobile gestures).
2. Examine Outcome `Output` dimensions for interface-specific assertions (exit codes → CLI, HTTP status codes → API, element selectors → UI/TUI/Mobile).
3. Record the detected type set (e.g., `{CLI, API}`).

### 2.5.2 Load Shared Principles

Always load `types/_shared.md` regardless of detected types. This file defines the five cross-type universal principles (Isolation, Determinism, Timeout Protection, Idempotency, Resource Cleanup) and shared antipattern guards.

### 2.5.3 Load Type-Specific Rules

For each interface type in the detected set, load the corresponding type file:

1. Map interface type to filename: `CLI` → `types/cli.md`, `TUI` → `types/tui.md`, `UI` → `types/ui.md`, `Mobile` → `types/mobile.md`, `API` → `types/api.md`.
2. Read each matched type file via Read tool.
3. Extract Golden Rules (generation constraints) and Reconnaissance Hints (discovery aids).

Do NOT load type files for interface types not detected in the Contracts. No speculative bulk loading.

<HARD-RULE>
`_shared.md` is ALWAYS loaded regardless of detected types. Only type files matching detected interface types are loaded — no speculative bulk loading.
</HARD-RULE>

### 2.5.4 Token Budget Warning

If the detected type set contains more than 3 types, emit:

```
WARNING: Detected {N} interface types ({type list}). Loading all type rules may consume significant token budget. Consider splitting the Journey into type-specific sub-Journeys.
```

Proceed with generation — the warning is advisory, not blocking.

## Step 3: Generate Test Code

For each Contract step, generate test code following the resolved framework's conventions AND the surface-specific strategy from Step 0.5.

### 3.0 Surface-Driven Generation Strategy

Before generating, apply the surface type detected in Step 0.5 to constrain the generation plan:

#### CLI Surface (Contract ≥ 80%)

- **Primary focus**: Contract test functions — one test per Outcome per step
- **Execution model**: Subprocess execution (binary isolation, environment hermeticity per `types/cli.md`)
- **Journey smoke tests**: Generate exactly 1 smoke test per Journey (happy path only)
- **Ratio enforcement**: For N Contract steps with M total Outcomes, generate M Contract test functions + 1 Journey smoke test. This ensures the ratio stays well above 80%.
- **Binary check**: Verify the binary can be built before generating tests. Auto-detect binary name and build command from Fact Table.

#### TUI Surface (Contract ≥ 80%)

- **Primary focus**: Contract test functions — one test per Outcome per step
- **Execution model**: Non-interactive stdin pipe with terminal output capture (per `types/tui.md`)
- **Journey smoke tests**: Generate exactly 1 smoke test per Journey (happy path only)
- **Ratio enforcement**: Same formula as CLI — M Contract functions + 1 Journey smoke test

#### WebUI Surface (Balanced 50/50)

- **Balanced approach**: Generate Contract tests for each Outcome AND enrich the Journey smoke test with multi-step verification
- **Execution model**: Convention-defined browser framework (per `types/ui.md`)
- **Journey smoke tests**: Generate 1 comprehensive smoke test that verifies the happy path AND at least 1 failure path through the Journey
- **Ratio target**: Approximately equal Contract test functions and Journey smoke test functions

#### API Surface (Balanced 50/50)

- **Balanced approach**: Generate Contract tests for each Outcome AND enrich the Journey smoke test
- **Execution model**: HTTP client testing (per `types/api.md`)
- **Journey smoke tests**: Generate 1 comprehensive smoke test that verifies the happy path AND at least 1 error path through the Journey
- **Ratio target**: Approximately equal Contract test functions and Journey smoke test functions

#### Mobile Surface (Best-Effort)

- **Maestro YAML skeleton**: Generate Maestro YAML flows instead of code-based test functions
- **Skeleton structure**: Each generated Maestro YAML file MUST contain:
  1. `appId` declaration (from Fact Table `MOBILE_APP_ID`)
  2. `onFlowStart: [launchApp]` lifecycle hook
  3. `onFlowEnd: [killApp]` lifecycle hook
  4. Navigation flow for the Contract step's happy path
  5. Deep link test variant when the Contract involves navigation to a specific screen
- **Deep link tests**: For each Journey step that navigates to a specific screen, generate an additional Maestro YAML that opens the app via URL scheme and asserts the target screen is visible
- **Complex scenario handling**: If a test case involves gestures not expressible in Maestro (pinch, rotate, multi-finger swipe), or requires physical device capabilities (sensors, camera), mark the test with `manual-only` annotation and skip generation. Add a comment explaining which capability requires manual testing.
- **Convention reference**: Use Maestro YAML syntax conventions. If no Mobile Convention file exists, use the Maestro reference from `types/mobile.md` as the authoritative syntax guide.

<HARD-RULE>
**Mobile best-effort**: Do not aim for comprehensive coverage. Generate skeleton flows + deep link tests for core Journeys only. Complex scenarios MUST be marked `manual-only` rather than generating incomplete or fragile tests. Mobile test generation MUST NOT fail the pipeline — any generation issue should result in a skeleton with `manual-only` markers.
</HARD-RULE>

### Output Directory

Tests go directly into `tests/<journey>/`. Contract specs are read from `docs/features/<slug>/testing/<journey>/contracts/`:

```
docs/features/<slug>/testing/
  task-lifecycle/                  <- Journey directory
    journey.md                     <- Journey narrative
    contracts/                     <- Contract specs (input, from gen-contracts)
      step-1-feature-create.md
      step-2-task-claim.md
      step-3-task-submit.md

tests/
  task-lifecycle/                  <- Generated test scripts
    step1_feature_create_test.go   <- Generated step test
    step2_task_claim_test.go       <- Generated step test (multiple Outcomes)
    step3_task_submit_test.go      <- Generated step test
    task_lifecycle_smoke_test.go   <- Journey smoke test (happy path E2E)
```

<HARD-RULE>
**No staging area**: Tests go directly to `tests/<journey>/`, NOT to `tests/e2e/features/` staging. The old staging model is replaced by tag-based lifecycle management.
</HARD-RULE>

### @feature Tags

All generated test files MUST include `@feature` tags. The tag format and mechanism are defined by the Convention file's **Tags** section. Read and follow the Convention's tag syntax precisely.

If the Convention file does not have a Tags section, ask the user which tag format to use. Common formats for reference only — do NOT auto-select without Convention or user input.

### Step Test Generation

For each Contract step, generate a test file containing one test function per Outcome:

1. Each test function validates one Outcome's assertions.
2. Use `t.TempDir()` or framework equivalent for isolation.
3. Assert Output matches the regex pattern from Step 1.
4. Assert State changes as specified in the Contract.
5. Include traceability comment linking back to the Contract.

### Journey Smoke Test

Generate exactly one smoke test per Journey that:

1. Runs the complete happy path end-to-end (all steps in sequence).
2. Passes state between steps (feature_slug, task_id, etc.).
3. Asserts each step's output matches the "success" Outcome's Output dimension.
4. Verifies Journey Invariants hold across all steps.

<HARD-RULE>
Every Journey MUST have at least 1 smoke test. The smoke test MUST only test the happy path (success Outcomes).
</HARD-RULE>

### Test Data Safety

<HARD-RULE>
**No hardcoded secrets**: Generated test code MUST NOT contain real secret/token values. Sensitive fields (token, api_key, password, secret, credential) MUST use environment variable placeholders:
- Go: `os.Getenv("E2E_API_TOKEN")`
- Python: `os.environ.get("E2E_API_TOKEN")`
- JavaScript: `process.env.E2E_API_TOKEN`
</HARD-RULE>

### Framework-Specific Rules

<EXTREMELY-IMPORTANT>
All framework-specific rules (test runner, assertion library, imports, HTTP client, process execution, anti-patterns) are defined in the Convention file loaded in Step 0. Read and follow those rules precisely.

- Use ONLY the framework specified in the Convention file
- Import paths and naming conventions follow the Convention's conventions
- Convention sections: framework, discovery, structure, assertions, Import Patterns, Code Style, Anti-patterns, Helpers
</EXTREMELY-IMPORTANT>

### Templates (Convention-Driven)

Convention files contain all framework-specific patterns (imports, assertion syntax, helpers, anti-patterns). Use the Convention's Code Style and Helpers sections as the template for generated code.

If the project has a custom template directory configured (`.forge/config.yaml` `test-template-dir`), load templates from that path. Otherwise, use the Convention file content as the authoritative template source.

<HARD-RULE>
**Template override**: If `test-template-dir` is set in config, load templates from that directory. Otherwise, use Convention file patterns as the template source.
</HARD-RULE>

## Step 4: Compile Gate

After generating all test files, run a compile gate to verify generated code correctness.

### 4.0 Syntax and Import Validation

Before compile, run lightweight syntax and import checks:

#### (a) Syntax Correctness

Verify each generated test file has valid syntax:

| Framework | Validation Method |
|-----------|-------------------|
| Go | `gofmt -e <file>` (exit code 0 = valid syntax) |
| JavaScript/TypeScript | `node --check <file>` or `tsc --noEmit` |
| Python | `python -m py_compile <file>` |
| Maestro YAML | `maestro validate <file>` (if CLI available) or YAML lint |

#### (b) Import Path Resolution

Verify generated imports can be resolved:

| Framework | Validation Method |
|-----------|-------------------|
| Go | `go vet ./tests/<journey>/...` (reports unresolved imports) |
| JavaScript/TypeScript | `tsc --noEmit` or `node -e "require('<import>')"` |
| Python | `python -c "import <module>"` for each generated import |
| Maestro YAML | N/A (no import resolution needed) |

#### Validation Failure Handling

For each file that fails validation:

1. **Auto-retry (1 attempt)**: Feed the syntax/import error back to LLM with the generated file content. Regenerate the failing file. Re-validate.
2. **Retry also fails**: Mark the file as `gen-failed` by adding a header comment:
   ```
   // GEN-FAILED: <error summary>
   // This file was generated but failed validation. Manual review required.
   ```
   For Maestro YAML:
   ```yaml
   # GEN-FAILED: <error summary>
   # This file was generated but failed validation. Manual review required.
   ```
3. **Skip and continue**: The `gen-failed` file is skipped in subsequent compile/test steps. Other generated files proceed normally.

<HARD-RULE>
**gen-failed files are NOT deleted**. They remain in `tests/<journey>/` for user inspection. The pipeline does NOT block on `gen-failed` files — other tests continue to compile and execute. At most 1 auto-retry per file; do not retry more than once.
</HARD-RULE>

### 4.1 Prerequisite Check

Verify `just compile` recipe exists:

```bash
just --list | grep compile
```

If the recipe is missing:
- Block generation
- Output: "Missing justfile `compile` recipe. Run `/forge:init-justfile` first, or add a recipe manually."
- Do not proceed with compile check

### 4.2 Compile and Retry

Run the compile gate:

```bash
just compile
```

| Result | Action |
|--------|--------|
| Pass | Proceed to post-generation checks (Step 4.3) |
| Fail (attempt 1) | Feed compile error + generated file content back to LLM. Regenerate the failing file. Retry compile. |
| Fail (attempt 2) | Feed compile error again. Regenerate with explicit error analysis. Retry compile. |
| Fail (attempt 3) | Block task. Output error details + recovery guidance. Do NOT delete generated files. |

### 4.3 Recovery on Exhaustion

If all compile attempts fail:

1. Output the compile error to the user with the generated file path
2. Suggest recovery actions:
   - (a) Check Convention file for incorrect framework/assertions section declarations
   - (b) Run `/forge:test-guide` to regenerate Convention from project analysis
   - (c) Manually edit the generated test file to fix compilation
3. Do not auto-delete the generated file — leave it for user inspection

### 4.4 Post-Compile Checks

After compile passes, run these additional checks:

#### VERIFY Marker Check

Scan generated files for unresolved `// VERIFY:` markers:
```bash
grep -rn '// VERIFY:' tests/<journey>/
```
Resolve any remaining markers using Fact Table values.

#### Antipattern Guard & Duplicate Name Check

See `rules/quality-gates.md` for the forbidden pattern list and duplicate name collision rules.

## Error Handling

See `rules/quality-gates.md` for the complete error handling table covering Convention, Contract, Fact Table, compile gate, and template resolution failures.

## Related Skills

| Skill | Usage |
|-------|-------|
| `/gen-journeys` | Generate Journey narratives from PRD |
| `/gen-contracts` | Generate Contract specifications from Journeys |
| `/run-tests` | Execute test scripts and report results |
| `/forge:test-guide` | Generate a Convention file for test framework configuration |
