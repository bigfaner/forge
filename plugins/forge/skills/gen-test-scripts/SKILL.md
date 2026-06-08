---
name: gen-test-scripts
description: Generate executable test scripts from Contract specifications. Journey-driven: generates test code with @feature tags. Output directory adapts to surface count: multi-surface projects use tests/<surfaceKey>/<journey>/, single-surface projects use tests/<journey>/. Test type naming follows Surface → Test Type mapping (cli → CLI Functional Test, api → API Functional Test, tui → Terminal Functional Test, web → Web E2E Test, mobile → Mobile E2E Test).
---

# Gen Test Scripts

Generate executable test scripts from Contract specifications.

**Core principle**: Tests are generated per Journey, not per interface type. Each Journey step's Contract defines the assertions. Output directory adapts to surface count: multi-surface projects write to `tests/<surfaceKey>/<journey>/`, single-surface projects write to `tests/<journey>/`.

## Pipeline Position

```
gen-journeys -> gen-contracts -> gen-test-scripts -> run-tests
                                      ^^^ YOU ARE HERE
```

Input: Contract specifications (from gen-contracts) + Fact Table (from code reconnaissance) + Handbook (design document, optional).
Output: Executable test code with `@feature` tags. Output directory adapts to surface count: multi-surface → `tests/<surfaceKey>/<journey>/`, single-surface → `tests/<journey>/`.

## Prerequisites

Check previous stage artifacts. Abort and prompt user if missing:

| Artifact | Missing prompt |
|----------|----------------|
| At least one Contract file in `docs/features/<slug>/testing/<journey>/contracts/` | Run `/gen-contracts` first |
| Eval report for all Contracts (`testing/<journey>/.eval-report.md`) | Run `/eval --type contract` first. **Blocker**: do not proceed if any Contract scored below target. |

### SKIP_EVAL_GATE Mode

When the task context contains `SKIP_EVAL_GATE=true` (injected by Quick mode task templates), the eval report prerequisite is **conditionally waived**:

- **Skip**: eval-contract report check (`testing/<journey>/.eval-report.md`) is bypassed entirely
- **Proceed directly**: move to Step 0 (Load Convention Files — surface-first) and Step 1 (Code Reconnaissance) without eval verification
- **Mark output**: every test file generated under SKIP_EVAL_GATE MUST include a header comment: `// SKIP_EVAL_GATE: generated without eval-contract verification. Review with extra scrutiny.`

**When SKIP_EVAL_GATE is NOT set** (Breakdown mode or manual `/gen-test-scripts` invocation): the eval report Blocker remains mandatory. Behavior is unchanged.

## Step 0: Load Convention Files

Load per-surface test strategy from Convention files (surface-first structure).

Load: `rules/convention-guide.md` — Convention file structure reference, section schema, validation rules, merge semantics, and growth path.

### 0.1 Old Structure Detection

Before loading Convention files, check whether the project uses the legacy (framework-first) structure:

1. Check if `docs/conventions/testing/` contains any `.md` files that are NOT inside a subdirectory (i.e., flat files like `go.md`, `vitest.md`).
2. If legacy files are detected:
   - Output migration prompt: "Legacy Convention structure detected in `docs/conventions/testing/` (framework-first files). Run `/test-guide` to regenerate with the new surface-first structure (`testing/{surface}/core.md`)."
   - Proceed with Step 0.2 using auto-detection (the old files are not loaded).
3. If no legacy files, proceed to Step 0.2.

### 0.2 Discover Convention Files (Surface-First)

1. Determine the active surface type (from Step 0.5).
2. Load the surface Convention from `docs/conventions/testing/{surface}/core.md`.
3. If `core.md` does not exist for the detected surface, proceed to auto-detection (Step 0.3).

<HARD-RULE>
Convention loading is surface-driven, not framework-driven. The `{surface}` segment comes from Step 0.5 surface detection. Do NOT fall back to loading framework-specific flat files.
</HARD-RULE>

### 0.3 Resolve Target Framework

1. **Convention assertion preference table**: If `core.md` was loaded, read its assertion preference table (per-framework rows) to identify the target framework.
2. **Existing test file scan**: Scan `tests/` for file patterns to confirm the framework matches the project.
3. **User specification**: If signals are ambiguous, ask the user which framework to use.
4. **No Convention found**: Proceed with LLM defaults + Code Reconnaissance (Step 1). Output hint: "No test Convention files found for surface `{surface}` in `docs/conventions/testing/{surface}/core.md`. Generation will use LLM defaults. Run `/test-guide` to create one."

<HARD-RULE>
If no Convention files are found and no framework can be detected from existing test files, ask the user which framework to use. Do NOT silently default.
</HARD-RULE>

### 0.4 Validate Convention Content

For the loaded Convention file (`core.md`), check required sections per the surface template: file location, isolation model, assertion focus, timeout strategy, lifecycle, Contract/Journey ratio, anti-patterns.

- **Missing required section**: Log warning listing missing sections. Proceed with LLM defaults for that section's area.
- **Invalid section content** (e.g., empty isolation model): Treat as missing. Log warning.

Use the loaded Convention content for all surface-specific strategy in subsequent steps. Framework implementation details come from the assertion preference table within `core.md`.

<HARD-RULE>
`types/*.md` (loaded in Step 2.5) is the primary authority for generation-time framework strategies. `core.md` is the authority for surface-level strategy (isolation model, assertion focus, etc.). When both cover the same aspect (assertion preferences), `types/*.md` takes precedence.
</HARD-RULE>

## Step 0.5: Surface Detection

Determine the project's interface surface type to drive per-surface generation strategy.

Load `rules/step-0.5-validation.md` for the complete surface detection and strategy application logic, including:
- Reading surface configuration from `.forge/config.yaml`
- Auto-detection fallback from code reconnaissance
- Surface strategy table (CLI/TUI/Web/API/Mobile ratio targets)
- Surface-driven generation strategy for Step 3.0

## Test Type Terminology

Test type names follow the Surface → Test Type mapping:

| Surface | Test Type | Tag |
|---------|-----------|-----|
| `cli` | CLI 功能测试 (CLI Functional Test) | `@cli-functional` |
| `tui` | 终端功能测试 (Terminal Functional Test) | `@tui-functional` |
| `api` | API 功能测试 (API Functional Test) | `@api-functional` |
| `web` | Web 端到端测试 (Web E2E Test) | `@web-e2e` |
| `mobile` | 移动端端到端测试 (Mobile E2E Test) | `@mobile-e2e` |

The "e2e" term is reserved exclusively for Web and Mobile surfaces. CLI, TUI, and API surfaces use "functional" terminology. Generated test code comments and `@feature` tags MUST use these surface-specific test type names, NOT the generic "e2e" label.

## Step 1: Code Reconnaissance (Build Fact Table)

Read source code to extract ground-truth values for semantic descriptor resolution AND test framework patterns.

Load `rules/step-1-contract-loading.md` for the complete reconnaissance logic, including:
- Framework reconnaissance (test file patterns, imports, build tags)
- Domain reconnaissance (CLI entry points, API handlers, config values)
- Fact Table construction with source citations
- Semantic descriptor to regex conversion

## Step 1.5: Cross-Validation (Fact Table vs Contract Anchors)

After building the Fact Table (Step 1), cross-validate it against Contract frontmatter anchor fields. This step detects mismatches between code reality (Fact Table) and design intent (Contract anchors), using the handbook (design document) as the authority source.

Load `rules/step-1.5-cross-validation.md` for the complete cross-validation logic, including:
- Anchor field extraction from Contract frontmatter
- Fact Table vs anchor comparison and classification
- Handbook authority resolution and suggestion generation
- Degradation mode when handbook or anchors are missing

### Cross-Validation Flow

1. **Extract anchors**: For each Contract in the target Journey, read frontmatter anchor fields based on the detected surface type:

   | Surface | Anchor fields to read |
   |---------|----------------------|
   | API | `endpoint`, `method` |
   | CLI | `command`, `subcommand` |
   | TUI | `command` |
   | Web | `page`, `route` |
   | Mobile | `screen`, `deeplink` |

2. **Match against Fact Table**: Compare each anchor value with corresponding Fact Table entries from code reconnaissance.

3. **Classify results**: Each comparison produces one of three classifications:

   | Classification | Criteria | Action |
   |---------------|----------|--------|
   | **High confidence match** | Anchor value matches Fact Table exactly (normalized comparison) | Proceed normally |
   | **Low confidence mismatch** | Anchor value differs from Fact Table, but Fact Table signal is incomplete (e.g., dynamic route registration, partial scan) | Log mismatch with both values, prompt user to confirm |
   | **Cannot verify** | No corresponding Fact Table entry exists, or anchor field is absent from Contract | Log as unverifiable, proceed with anchor value if present, or Fact Table inference if absent |

4. **Authority resolution**: When a mismatch is detected, resolve authority:

   | Handbook exists? | Authority source | Mismatch action |
   |-----------------|-----------------|-----------------|
   | Yes, and matches anchor | Handbook = anchor | Fact Table differs -> **code bug**: handbook says X, code does Y. Generate code bug report. |
   | Yes, and differs from anchor | Handbook | Anchor is stale. Generate suggested fix (diff) to update Contract anchor to match handbook. User confirms before writing. |
   | Yes, and differs from Fact Table | Handbook | Code does not match handbook -> **code bug**. Generate code bug report. |
   | No | Fact Table | Degraded mode: use Fact Table as inference source. Prompt user that handbook is missing and recommend generating one. |
   | No, anchor also missing | Fact Table | Full degradation: use Fact Table inference, no cross-validation possible. Prompt user. |

5. **Suggestion generation**: For each mismatch where handbook exists and differs from anchor, generate a suggested fix:

   - Show a diff of the current Contract frontmatter vs proposed change
   - Include the handbook source citation for the proposed value
   - Present to user for confirmation before writing to Contract
   - If user rejects, keep current anchor value and log the rejection

<HARD-RULE>
**Handbook is the authority source** for cross-validation. When handbook and code implementation disagree, the discrepancy is flagged as a **code bug** (not a test or Contract issue). The user confirmation step is the final gate -- no automatic writes to Contract without explicit user approval.
</HARD-RULE>

<HARD-RULE>
**Low confidence and cannot-verify results are NEVER auto-resolved**. They are reported to the user for manual confirmation. The pipeline does not block on these classifications.
</HARD-RULE>

### Degradation Mode (Backward Compatibility)

When handbook is missing or anchor fields are absent from Contract:

1. **No handbook**: Skip cross-validation for the relevant surface. Use Fact Table values as inference source. Output prompt: "Handbook not found for surface `{surface}`. Cross-validation skipped. Recommend running `/tech-design` to generate handbook for improved anchor accuracy."
2. **No anchor fields in Contract**: The Contract predates technical anchors. Use Fact Table inference as fallback. Output prompt: "Contract `{contract_path}` has no anchor fields. Using Fact Table inference. Consider running `/gen-contracts` to populate anchors from handbook."
3. **Both missing**: Proceed with existing Step 1 Fact Table inference only. No cross-validation. Output both prompts above.

Degradation mode is non-blocking. The pipeline continues normally with reduced confidence.

### Surface Coverage Report

After cross-validation completes (or degrades), output a coverage report:

```
=== Surface Coverage Report ===
Surface: API
  - Contracts with anchors: 3/4
  - Cross-validated (high confidence): 2
  - Mismatches detected: 1 (low confidence: 0, cannot verify: 1)
  - Code bugs flagged: 0
  - Suggested fixes pending user confirmation: 1

Surface: CLI
  - Contracts with anchors: 2/2
  - Cross-validated (high confidence): 2
  - Mismatches detected: 0
  - Code bugs flagged: 0
  - Suggested fixes pending user confirmation: 0

Surfaces not covered:
  - Web: no handbook found
  - Mobile: no contracts in journey
  - TUI: not applicable (surface type = CLI)

Summary: 4/6 anchors verified, 1 mismatch, 0 code bugs, 1 fix pending
```

The report MUST:
- List each surface type present in the Journey's Contracts
- Show anchor coverage ratio (Contracts with anchors / total Contracts)
- Count verification results by classification
- Explicitly list surfaces that were NOT verified and why (no handbook, no contracts, not applicable)
- Provide a summary line with totals

### Lesson Scenario Capture

The cross-validation MUST capture the lesson scenario from the proposal evidence:

**Scenario**: Contract does not specify HTTP method (no `method` anchor) or specifies POST, but the actual route is registered as PUT (per handbook).

**Expected behavior**:
1. Fact Table finds `PUT /teams/:teamId/sub-items/:subId/move` from code reconnaissance
2. Handbook defines `PUT /teams/:teamId/sub-items/:subId/move`
3. Contract anchor says `method: POST` (or missing)
4. Cross-validation detects mismatch: Fact Table (PUT) vs Contract anchor (POST)
5. Authority check: handbook says PUT -> suggest fix: change Contract `method` from POST to PUT
6. Generate diff, present to user for confirmation

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

1. Examine `step-action` fields for interface indicators (CLI commands, HTTP methods, Web interactions, TUI rendering, mobile gestures).
2. Examine Outcome `Output` dimensions for interface-specific assertions (exit codes -> CLI, HTTP status codes -> API, element selectors -> Web/TUI/Mobile).
3. Record the detected type set (e.g., `{CLI, API}`).

### 2.5.2 Load Shared Principles

Always load `types/_shared.md` regardless of detected types. This file defines the five cross-type universal principles (Isolation, Determinism, Timeout Protection, Idempotency, Resource Cleanup) and shared antipattern guards.

### 2.5.3 Load Type-Specific Rules

For each interface type in the detected set, load the corresponding type file:

1. Map interface type to filename: `CLI` -> `types/cli.md`, `TUI` -> `types/tui.md`, `Web` -> `types/web.md`, `Mobile` -> `types/mobile.md`, `API` -> `types/api.md`.
2. Read each matched type file via Read tool.
3. Extract Golden Rules (generation constraints) and Reconnaissance Hints (discovery aids).

Do NOT load type files for interface types not detected in the Contracts. No speculative bulk loading.

<HARD-RULE>
`_shared.md` is ALWAYS loaded regardless of detected types. Only type files matching detected interface types are loaded -- no speculative bulk loading.
</HARD-RULE>

### 2.5.4 Assertion Depth and Fixture Quality Rules

Load behavioral test quality rules that constrain assertion depth and fixture data richness:

1. **Load `rules/assertion-depth.md`**: Assertion classification criteria (behavioral vs structural), >=80% behavioral threshold, >=30% deep assertion requirement. This rule is enforced during Step 3 generation — the agent MUST count and classify assertions, supplementing if thresholds are not met.

2. **Load `rules/fixture-from-spec.md`**: Fixture data generation from Contract `fixture_spec` declarations. When a Contract's Preconditions contain `fixture_spec.entities`, test code MUST create >= `min_count` entities with correct relationships and field constraints. When `fixture_spec` is absent, apply backward compatibility handling from `types/_shared.md`.

<HARD-RULE>
**Assertion depth enforcement is mandatory**: Every generated Journey's test suite MUST satisfy both the >=80% behavioral threshold and >=30% deep assertion requirement. Health check / readiness-only Journeys are exempt but MUST document the exemption with a header comment.
</HARD-RULE>

### 2.5.5 Token Budget Warning

If the detected type set contains more than 3 types, emit:

```
WARNING: Detected {N} interface types ({type list}). Loading all type rules may consume significant token budget. Consider splitting the Journey into type-specific sub-Journeys.
```

Proceed with generation -- the warning is advisory, not blocking.

## Step 3: Generate Test Code

For each Contract step, generate test code following the resolved framework's conventions AND the surface-specific strategy from Step 0.5.

### 3.0 Surface-Driven Generation Strategy

Apply the surface type detected in Step 0.5 to constrain the generation plan. Load `rules/step-0.5-validation.md` section "Surface-Driven Generation Strategy" for per-surface ratio targets, execution models, and generation constraints (CLI, TUI, Web, API, Mobile).

**Test isolation**: Every test — unit or surface — must own its environment. No test may depend on the real project's filesystem state, git state, or `.forge/` state.

<!-- INLINE:origin=run-tests/rules/test-isolation.md -->

| Rule ID | Scope | Requirement | Pattern |
|---------|-------|-------------|---------|
| TEST-isolation-000 | All tests | Every test MUST create its own world via `t.TempDir()` (unit) or isolated fixture setup (e2e). Tests MUST NOT read, write, or depend on files outside their own sandbox. | Use `t.TempDir()` / `t.Setenv()` for all stateful tests |
| TEST-isolation-001 | Unit tests (project root) | Tests calling `FindProjectRoot()`, `runFeature()`, `executeClaim()`, or any function relying on project root detection MUST set `CLAUDE_PROJECT_DIR` via `t.Setenv()` to isolate from real `.forge/state.json` and workspace markers. | `t.Setenv("CLAUDE_PROJECT_DIR", dir)` + `os.Chdir(dir)` |
| TEST-isolation-002 | CLI functional tests | Tests invoking `forge` CLI commands MUST pass `CLAUDE_PROJECT_DIR=<fixture-dir>` via `cmd.Env` to prevent CLI from detecting real project root. | `cmd.Env = append(os.Environ(), "CLAUDE_PROJECT_DIR="+dir)` |
| TEST-isolation-003 | CLI functional test fixtures | Helpers creating project directories MUST create all files required by code under test — including `tasks/index.json`, `.forge/config.yaml`, and any other files the production code checks for. | Include `index.json` in `ensureFeatureDir` |
| TEST-isolation-004 | CLI functional tests (binary) | Test files invoking forge CLI commands SHOULD compile a dedicated forge binary from current source tree via `go build` and use it for all `exec.Command` invocations, rather than relying on system-installed `forge` via `$PATH`. | `TestMain` builds binary; tests use `exec.Command(forgeBinary, ...)` |

<!-- END INLINE:origin=run-tests/rules/test-isolation.md -->

### Output Directory

Output directory adapts to the project's surface count. Determine the correct path by checking whether the project has multiple surfaces (via `forge surfaces` or `.forge/config.yaml`):

- **Multi-surface** (2+ surfaces): `tests/<surfaceKey>/<journey>/`
- **Single-surface** (1 surface): `tests/<journey>/`

Contract specs are read from `docs/features/<slug>/testing/<journey>/contracts/`.

**Multi-surface example** (surfaces: `backend=api`, `frontend=web`):

```
docs/features/<slug>/testing/
  task-lifecycle/                  <- Journey directory
    journey.md                     <- Journey narrative
    contracts/                     <- Contract specs (input, from gen-contracts)
      step-1-feature-create.md
      step-2-task-claim.md
      step-3-task-submit.md

tests/
  backend/                         <- surfaceKey directory
    task-lifecycle/                <- Generated test scripts for backend surface
      step1_feature_create_test.go
      step2_task_claim_test.go
  frontend/                        <- surfaceKey directory
    task-lifecycle/                <- Generated test scripts for frontend surface
      step1_feature_create_test.go
      step2_task_claim_test.go
```

**Single-surface example** (surfaces: `tui` or `surfaces: [{key: app, type: tui}]`):

```
docs/features/<slug>/testing/
  task-lifecycle/                  <- Journey directory
    journey.md                     <- Journey narrative
    contracts/                     <- Contract specs (input, from gen-contracts)
      step-1-feature-create.md
      step-2-task-claim.md
      step-3-task-submit.md

tests/
  task-lifecycle/                  <- Generated test scripts (no surfaceKey layer)
    step1_feature_create_test.go
    step2_task_claim_test.go
    step3_task_submit_test.go
    task_lifecycle_smoke_test.go
```

<HARD-RULE>
**No staging area**: Tests go directly to the output directory (adaptive per surface count), NOT to `tests/e2e/features/` staging. The old staging model is replaced by tag-based lifecycle management.
</HARD-RULE>

<HARD-RULE>
**Surface-key directory**: When the project has multiple surfaces, the `<surfaceKey>` segment MUST use the surface's key (e.g., `backend`, `frontend`), NOT the surface's type (e.g., `api`, `web`). The key is the user-defined identifier; the type is the technical classification. Keys are unique; types can repeat.
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
| Go | `gofmt -e <file>` | `go vet ./tests/...` (adaptive path per surface count) |
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

- **VERIFY Marker Check**: `grep -rn '// VERIFY:' tests/` (adaptive per surface count) -- resolve remaining markers using Fact Table.
- **Antipattern Guard & Duplicate Name Check**: See `rules/quality-gates.md`.
- **Assertion Depth Check**: See `rules/assertion-depth.md` — verify >=80% behavioral assertion threshold and >=30% deep assertion requirement are met per Journey. If thresholds are not met and no exemption header exists, regenerate with supplemented assertions.
- **Fixture Spec Compliance**: See `rules/fixture-from-spec.md` — verify that when Contract contains `fixture_spec.entities`, the generated test creates >= `min_count` entities with correct relationships.

## Error Handling

See `rules/quality-gates.md` for the complete error handling table covering Convention, Contract, Fact Table, compile gate, and template resolution failures.
