---
name: gen-test-scripts
description: Generate executable test scripts from Contract specifications. Journey-driven: generates test code with @feature tags directly into tests/<journey>/.
conventions:
  - testing-isolation.md
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

## Step 0: Resolve Language and Strategy

1. **Detect language**: Run `forge test detect` to auto-detect the project's test language(s) from file signals.
2. **On failure** (no language detected): ask the user to add `languages` to `.forge/config.yaml` (e.g., `languages: [go]`).
3. **Load strategy**: Run `forge test get generate` to load the generate strategy for the detected language.
4. **Resolve framework**: Run `forge test framework` to determine the test framework and its code conventions.

Use the loaded strategy and framework for all subsequent steps.

<HARD-RULE>
Do NOT silently default to any language. If `forge test detect` returns no result and the user cannot configure `languages`, abort the skill.
</HARD-RULE>

## Step 1: Read Contract Specifications

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

## Step 2: Code Reconnaissance (Build Fact Table)

Read source code to extract ground-truth values for semantic descriptor resolution.

**Reconnaissance reads**:

| Source | What to extract |
|--------|-----------------|
| CLI entry points | Command names, flag names, output format strings |
| API handlers | Request/response schemas, status codes |
| TUI components | Model fields, View output patterns |
| Config files | Ports, base paths, auth mechanisms |
| Auth implementation | Login endpoint, token field, header format |

Build Fact Table with source citations:

```markdown
## Fact Table
| Key | Value | Source |
|-----|-------|--------|
| CLI_TASK_CLAIM_OUTPUT | claimed task <task_id> | internal/cmd/claim.go:42 |
| CLI_FEATURE_CREATE_OUTPUT | Feature <slug> created successfully | internal/cmd/feature.go:45 |
```

<HARD-RULE>
- Every Fact Table value must cite source file and line number. Unknown sources -> `UNKNOWN`. Do not fabricate.
- Fact Table values drive semantic descriptor to regex conversion. All `// VERIFY:` markers must be resolved using Fact Table values.
</HARD-RULE>

### Semantic Descriptor to Regex Conversion

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

All generated test files must include `@feature` tags using the framework's native mechanism:

| Framework | Tag format |
|-----------|------------|
| Go testing | `//go:build feature` (build tag at top of file) |
| pytest | `pytestmark = pytest.mark.feature` (module-level mark) |
| mocha | `describe("@feature", () => { ... })` (wrapper describe) |
| JUnit5 | `@Tag("feature")` (class/method annotation) |
| Rust test | `#[cfg(feature = "feature")]` (cfg attribute) |

### Step Test Generation

For each Contract step, generate a test file containing one test function per Outcome:

1. Each test function validates one Outcome's assertions.
2. Use `t.TempDir()` or framework equivalent for isolation.
3. Assert Output matches the regex pattern from Step 2.
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

Follow the active strategy's `generate.md` for all framework-specific patterns:

<EXTREMELY-IMPORTANT>
All framework-specific rules (test runner, assertion library, imports, HTTP client, process execution, anti-patterns) are defined in the active strategy's `generate.md` (loaded in Step 0). Read and follow those rules precisely.

- Use ONLY the framework specified in the strategy's `generate.md`
- Import paths and naming conventions follow the framework's conventions
</EXTREMELY-IMPORTANT>

### Built-in Templates (Default, Overridable)

The 6 built-in language profiles serve as default templates. When a project has zero custom template configuration:

1. `forge test get template <filename>` returns built-in template content.
2. Built-in templates define: test file structure, helper functions, auth setup patterns.
3. Zero-config output equals built-in template output (diff is empty).

Custom template override: When `.forge/config.yaml` declares a custom template directory path, `gen-test-scripts` uses templates from that path instead of built-in ones.

<HARD-RULE>
**Template override**: If `test-template-dir` is set in config, load templates from that directory. Otherwise, use built-in default templates from `forge test get template`.
</HARD-RULE>

## Step 4: Post-Generation Verification

After generating all test files, run verification checks:

### Compilation Check

Run the appropriate compilation command:
- Go: `go build ./tests/<journey>/...` or `go test -c ./tests/<journey>/...`
- Python: `pytest --collect-only tests/<journey>/`
- JavaScript: `tsc --noEmit`

### VERIFY Marker Check

Scan generated files for unresolved `// VERIFY:` markers:
```bash
grep -rn '// VERIFY:' tests/<journey>/
```
Resolve any remaining markers using Fact Table values.

### Antipattern Guard

Verify each generated test function does not match any forbidden pattern:

| # | Forbidden Pattern | Instead |
|---|-------------------|---------|
| 1 | Recursive test invocation | Recursion guard (env var) or `-run` flag |
| 2 | Unconditional `t.Skip` | Implement with fixture or don't generate |
| 3 | Vacuous assertions (conditional assert without else fail) | Every assertion reachable on every code path |
| 4 | Environment-dependent skip without own fixture | `t.TempDir()` + own project structure |
| 5 | Duplicate test function names across packages | Scan for collisions; unique names with journey slug |
| 6 | Static-file text grep (assert on source file content) | Test runtime behavior only |

### Duplicate Name Check

Before writing, scan existing test files in the module for matching function names. If a collision is found, use a unique name that includes the journey slug.

## Error Handling

| Situation | Action |
|-----------|--------|
| Language detection fails | Ask user to configure `languages` in config.yaml |
| Contract files not found | Abort with prompt to run `/gen-contracts` |
| Fact Table lookup fails for a descriptor | Keep `// VERIFY:` marker, do not fabricate regex |
| Compilation fails post-generation | Fix generated code, re-run compile check |
| No test files generated | Abort with clear diagnostic message |
| Custom template path not found | Fall back to built-in templates with WARNING |

## Related Skills

| Skill | Usage |
|-------|-------|
| `/gen-journeys` | Generate Journey narratives from PRD |
| `/gen-contracts` | Generate Contract specifications from Journeys |
| `/run-e2e-tests` | Execute test scripts and report results |
