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

## Step 0: Load Convention Files

Load test framework knowledge from Convention files (no Profile/CLI dependency).

### 0.1 Discover Convention Files

1. Glob `docs/conventions/testing-*.md` in the project root.
2. Read each file's YAML frontmatter `domains` field.
3. Keep files whose `domains` contain `testing`.
4. Skip files with no `domains` frontmatter — output warning: "Convention file `<path>` has no domains frontmatter. Skipping."
5. Skip files that cannot be read (permissions, encoding) — output warning: "Cannot read Convention file `<path>`: `<error>`. Skipping."

### 0.2 Resolve Target Framework

Determine the target framework from the available signals:

1. **Convention file match**: If Convention files exist, their `domains` indicate the framework (e.g., `[testing, go]` = Go testing, `[testing, javascript]` = JS testing).
2. **Existing test file scan**: Scan `tests/` for file patterns to confirm the Convention's framework matches the project.
3. **User specification**: If multiple Convention files match or signals are ambiguous, ask the user which framework to use.
4. **No Convention found**: Proceed with LLM defaults + Code Reconnaissance (Step 1). Output hint: "No test Convention files found in `docs/conventions/`. Generation will use LLM defaults. Run `/forge:test-guide` to create one."

<HARD-RULE>
If no Convention files are found and no framework can be detected from existing test files, ask the user which framework to use. Do NOT silently default.
</HARD-RULE>

### 0.3 Validate Convention Content

For the loaded Convention file(s), check required sections: `Framework`, `Assertion`, `Tags`, `Result Format`.

- **Missing required section**: Log warning listing missing sections. Proceed with LLM defaults for that section's area. Example: "Convention file `<path>` is missing sections: Assertion, Tags. Using LLM defaults for those sections."
- **Invalid section content** (e.g., empty Framework name): Treat as missing. Log warning.
- **Multiple Convention files with overlapping domains**: Merge at section level — last-loaded file's section wins for conflicting sections. Log a note about the overlap for user awareness.
- **Convention vs Reconnaissance conflict** (detected in Step 1): Convention wins. Log the conflict for user awareness.

Use the loaded Convention content for all framework-specific rules in subsequent steps.

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

1. Glob `tests/<journey>/_contracts/step-*.md` for the target Journey.
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

For each Contract step, generate test code following the resolved framework's conventions.

### Output Directory

Tests go directly into `tests/<journey>/`:

```
tests/
  task-lifecycle/                  <- Journey directory
    _contracts/                    <- Contract specs (input, from gen-contracts)
      step-1-feature-create.md
      step-2-task-claim.md
      step-3-task-submit.md
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
- Convention sections: Framework, Assertion, Tags, Result Format, Import Patterns, Code Style, Anti-patterns, Helpers
</EXTREMELY-IMPORTANT>

### Templates (Convention-Driven)

Convention files contain all framework-specific patterns (imports, assertion syntax, helpers, anti-patterns). Use the Convention's Code Style and Helpers sections as the template for generated code.

If the project has a custom template directory configured (`.forge/config.yaml` `test-template-dir`), load templates from that path. Otherwise, use the Convention file content as the authoritative template source.

<HARD-RULE>
**Template override**: If `test-template-dir` is set in config, load templates from that directory. Otherwise, use Convention file patterns as the template source.
</HARD-RULE>

## Step 4: Compile Gate

After generating all test files, run a compile gate to verify generated code correctness.

### 4.1 Prerequisite Check

Verify `just e2e-compile` recipe exists:

```bash
just --list | grep e2e-compile
```

If the recipe is missing:
- Block generation
- Output: "Missing justfile `e2e-compile` recipe. Run `/forge:init-justfile` first, or add a recipe manually."
- Do not proceed with compile check

### 4.2 Compile and Retry

Run the compile gate:

```bash
just e2e-compile
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
   - (a) Check Convention file for incorrect framework/assertion declarations
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
| `/run-e2e-tests` | Execute test scripts and report results |
| `/forge:test-guide` | Generate a Convention file for test framework configuration |
