---
name: gen-test-scripts
description: Generate executable test scripts from Contract specifications. Journey-driven: generates test code with @feature tags directly into tests/<journey>/. Test type naming follows Surface → Test Type mapping (see docs/reference/test-type-model.md).
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

Load: `rules/convention-guide.md` — Convention file structure reference, section schema, validation rules, merge semantics, and growth path.

### 0.1 Discover Convention Files

1. Read `docs/conventions/testing/index.md` -- lists all available Conventions with name, description, and applicability conditions.
2. Based on the project's language/framework context, select the matching Convention from the index.
3. Load the selected Convention file from `docs/conventions/testing/<convention>.md`.
4. If `index.md` does not exist, proceed to auto-detection (Step 0.2).

<HARD-RULE>
Do NOT use `domains` frontmatter filtering. Selection is based on index.md descriptions and project context, with LLM autonomous judgment.
</HARD-RULE>

### 0.2 Resolve Target Framework

1. **Convention file match**: If a Convention was loaded from index.md, use its framework declaration.
2. **Existing test file scan**: Scan `tests/` for file patterns to confirm the Convention's framework matches the project.
3. **User specification**: If signals are ambiguous, ask the user which framework to use.
4. **No Convention found**: Proceed with LLM defaults + Code Reconnaissance (Step 1). Output hint: "No test Convention files found in `docs/conventions/testing/`. Generation will use LLM defaults. Run `/forge:test-guide` to create one."

<HARD-RULE>
If no Convention files are found and no framework can be detected from existing test files, ask the user which framework to use. Do NOT silently default.
</HARD-RULE>

### 0.3 Validate Convention Content

For the loaded Convention file, check required sections: `framework`, `discovery`, `structure`, `assertions`.

- **Missing required section**: Log warning listing missing sections. Proceed with LLM defaults for that section's area.
- **Invalid section content** (e.g., empty framework name): Treat as missing. Log warning.

Use the loaded Convention content for all framework-specific rules in subsequent steps.

## Step 0.5: Surface Detection

Determine the project's interface surface type to drive per-surface generation strategy.

Load `rules/step-0.5-validation.md` for the complete surface detection and strategy application logic, including:
- Reading surface configuration from `.forge/config.yaml`
- Auto-detection fallback from code reconnaissance
- Surface strategy table (CLI/TUI/WebUI/API/Mobile ratio targets)
- Surface-driven generation strategy for Step 3.0

## Test Type Terminology

Test type names follow the Surface → Test Type mapping defined in `docs/reference/test-type-model.md`:

| Surface | Test Type | Tag |
|---------|-----------|-----|
| `cli` | CLI 功能测试 (CLI Functional Test) | `@cli-functional` |
| `tui` | 终端功能测试 (Terminal Functional Test) | `@tui-functional` |
| `api` | API 功能测试 (API Functional Test) | `@api-functional` |
| `web` | Web 端到端测试 (Web E2E Test) | `@web-e2e` |
| `mobile` | 移动端端到端测试 (Mobile E2E Test) | `@mobile-e2e` |

Generated test code comments and `@feature` tags MUST use these surface-specific test type names, NOT the generic "e2e" label. The "e2e" term is reserved exclusively for Web and Mobile surfaces.

## Step 1: Code Reconnaissance (Build Fact Table)

Read source code to extract ground-truth values for semantic descriptor resolution AND test framework patterns.

Load `rules/step-1-contract-loading.md` for the complete reconnaissance logic, including:
- Framework reconnaissance (test file patterns, imports, build tags)
- Domain reconnaissance (CLI entry points, API handlers, config values)
- Fact Table construction with source citations
- Semantic descriptor to regex conversion

## Step 2: Read Contract Specifications

**Optional -- Run-to-Learn (R2L)**: If R2L is enabled in `.forge/config.yaml`, execute the R2L loop between Step 1 and Step 3. Load: `rules/run-to-learn.md` for the complete R2L mechanism (skeleton test generation, runtime fact enrichment, graceful degradation).

**Input discovery** -- find the Contract files for the target Journey:

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

1. Examine `step-action` fields for interface indicators (CLI commands, HTTP methods, UI interactions, TUI rendering, mobile gestures).
2. Examine Outcome `Output` dimensions for interface-specific assertions (exit codes -> CLI, HTTP status codes -> API, element selectors -> UI/TUI/Mobile).
3. Record the detected type set (e.g., `{CLI, API}`).

### 2.5.2 Load Shared Principles

Always load `types/_shared.md` regardless of detected types. This file defines the five cross-type universal principles (Isolation, Determinism, Timeout Protection, Idempotency, Resource Cleanup) and shared antipattern guards.

### 2.5.3 Load Type-Specific Rules

For each interface type in the detected set, load the corresponding type file:

1. Map interface type to filename: `CLI` -> `types/cli.md`, `TUI` -> `types/tui.md`, `UI` -> `types/ui.md`, `Mobile` -> `types/mobile.md`, `API` -> `types/api.md`.
2. Read each matched type file via Read tool.
3. Extract Golden Rules (generation constraints) and Reconnaissance Hints (discovery aids).

Do NOT load type files for interface types not detected in the Contracts. No speculative bulk loading.

<HARD-RULE>
`_shared.md` is ALWAYS loaded regardless of detected types. Only type files matching detected interface types are loaded -- no speculative bulk loading.
</HARD-RULE>

### 2.5.4 Token Budget Warning

If the detected type set contains more than 3 types, emit:

```
WARNING: Detected {N} interface types ({type list}). Loading all type rules may consume significant token budget. Consider splitting the Journey into type-specific sub-Journeys.
```

Proceed with generation -- the warning is advisory, not blocking.

## Step 3: Generate Test Code

For each Contract step, generate test code following the resolved framework's conventions AND the surface-specific strategy from Step 0.5.

### 3.0 Surface-Driven Generation Strategy

Apply the surface type detected in Step 0.5 to constrain the generation plan. Load `rules/step-0.5-validation.md` section "Surface-Driven Generation Strategy" for per-surface ratio targets, execution models, and generation constraints (CLI, TUI, WebUI, API, Mobile).

**Test isolation**: Generated tests must follow isolation conventions per `rules/test-isolation.md` (located in the run-tests skill directory, resolve relative to the skills parent directory) — every test must own its environment, no dependency on real project state.

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
    task_lifecycle_smoke_test.go   <- Journey smoke test (happy path, full Journey sequence)
```

<HARD-RULE>
**No staging area**: Tests go directly to `tests/<journey>/`, NOT to `tests/e2e/features/` staging. The old staging model is replaced by tag-based lifecycle management.
</HARD-RULE>

### @feature Tags

All generated test files MUST include `@feature` tags. The tag format and mechanism are defined by the Convention file's **Tags** section. Read and follow the Convention's tag syntax precisely.

If the Convention file does not have a Tags section, ask the user which tag format to use. Common formats for reference only -- do NOT auto-select without Convention or user input.

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

Use the Convention file content as the authoritative template source for all framework-specific patterns.

## Step 4: Compile Gate

After generating all test files, run a compile gate to verify generated code correctness.

### 4.0 Syntax and Import Validation

Before compile, run lightweight syntax and import checks:

| Framework | Syntax Validation | Import Validation |
|-----------|-------------------|-------------------|
| Go | `gofmt -e <file>` | `go vet ./tests/<journey>/...` |
| JavaScript/TypeScript | `node --check <file>` or `tsc --noEmit` | `tsc --noEmit` or `node -e "require('<import>')"` |
| Python | `python -m py_compile <file>` | `python -c "import <module>"` per import |
| Maestro YAML | `maestro validate <file>` or YAML lint | N/A |

#### Validation Failure Handling

1. **Auto-retry (1 attempt)**: Feed error + file content back to LLM. Regenerate. Re-validate.
2. **Retry also fails**: Mark file as `gen-failed` with header comment: `// GEN-FAILED: <error summary>`.
3. **Skip and continue**: `gen-failed` file is skipped in subsequent steps. Others proceed normally.

<HARD-RULE>
**gen-failed files are NOT deleted**. They remain for user inspection. Pipeline does NOT block on `gen-failed` files. At most 1 auto-retry per file.
</HARD-RULE>

### 4.1 Prerequisite Check

Verify `just compile` recipe exists:

```bash
just --list | grep compile
```

If missing, block generation with: "Missing justfile `compile` recipe. Run `/forge:init-justfile` first, or add a recipe manually."

### 4.2 Compile and Retry

Run `just compile`. On failure: retry up to 3 attempts (feed compile error + file content back to LLM, regenerate). After 3 failures: block task, output error details, preserve generated files.

### 4.3 Recovery on Exhaustion

If all compile attempts fail:
1. Output compile error with generated file path.
2. Suggest: (a) check Convention for incorrect declarations, (b) run `/forge:test-guide` to regenerate Convention, (c) manually edit.
3. Do not auto-delete generated files.

### 4.4 Post-Compile Checks

- **VERIFY Marker Check**: `grep -rn '// VERIFY:' tests/<journey>/` -- resolve remaining markers using Fact Table.
- **Antipattern Guard & Duplicate Name Check**: See `rules/quality-gates.md`.

## Error Handling

See `rules/quality-gates.md` for the complete error handling table covering Convention, Contract, Fact Table, compile gate, and template resolution failures.

## Related Skills

| Skill | Usage |
|-------|-------|
| `/gen-journeys` | Generate Journey narratives from PRD |
| `/gen-contracts` | Generate Contract specifications from Journeys |
| `/run-tests` | Execute test scripts and report results |
| `/forge:test-guide` | Generate a Convention file for test framework configuration |
