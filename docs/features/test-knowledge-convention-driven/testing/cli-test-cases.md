---
feature: "test-knowledge-convention-driven"
type: CLI
generated: "2026-05-20"
---

# CLI Test Cases: test-knowledge-convention-driven

Test cases derived from PRD acceptance criteria. Every case traces to a specific PRD source. 36 test cases covering all PRD stories, functional specs, and scope items.

---

## TC-001: Gen-Test-Scripts Uses Convention-Declared Framework for Non-Default Projects

- **Source**: Story 1 / AC-1
- **Type**: CLI
- **Target**: cli/gen-test-scripts
- **Test ID**: cli/gen-test-scripts/uses-convention-framework-for-ginkgo
- **Pre-conditions**: A Go project exists using ginkgo (not go-testing). A Convention file exists at `docs/conventions/testing-go.md` declaring ginkgo as the framework. The project has a valid Journey with Contract specs. `justfile` contains an `e2e-compile` recipe.
- **Steps**:
  1. Set up a Go project fixture with ginkgo imports in existing test files (`*_test.go` with `import . "github.com/onsi/ginkgo/v2"`)
  2. Create Convention file `docs/conventions/testing-go.md` with `Framework: ginkgo`, `Assertion: gomega`, `Tags: //go:build e2e`, `Result Format: json-stream`
  3. Create a valid Journey and Contract spec in the feature directory
  4. Run `forge gen-test-scripts` targeting the Journey
  5. Run `grep "ginkgo/v2" <generated-test-path>` to verify ginkgo imports are present
  6. Run `grep -c "Expect.*To\|Expect.*Should" <generated-test-path>` to verify gomega assertion style (>= 1 match)
  7. Run `grep "//go:build e2e" <generated-test-path>` to verify build tag
  8. Run `just e2e-compile`
- **Expected**:
  - Exit code 0 from `just e2e-compile`
  - Step 5: `grep` exits with code 0 (ginkgo import found)
  - Step 6: grep count >= 1 (gomega assertions present, not testify `assert.NoError`)
  - Step 7: `grep` exits with code 0 (`//go:build e2e` build tag found)
- **Priority**: P0

---

## TC-002: Gen-Test-Scripts Warns and Falls Back on Empty Convention Framework Section

- **Source**: Story 1 / AC-2
- **Type**: CLI
- **Target**: cli/gen-test-scripts
- **Test ID**: cli/gen-test-scripts/warns-on-empty-convention-framework
- **Pre-conditions**: A Go project exists. A Convention file exists at `docs/conventions/testing-go.md` with an empty or missing `Framework` section. The project has a valid Journey with Contract specs. `justfile` contains an `e2e-compile` recipe.
- **Steps**:
  1. Create Convention file `docs/conventions/testing-go.md` with `Assertion`, `Tags`, `Result Format` sections populated, but `Framework` section empty or absent
  2. Create a valid Journey and Contract spec
  3. Run `forge gen-test-scripts` targeting the Journey
  4. Capture stdout/stderr for warning output
  5. Run `grep -c "testing\|testify\|assert" <generated-test-path>` to verify generated test file uses LLM-detected defaults
  6. Run `just e2e-compile`
- **Expected**:
  - Exit code 0 from `just e2e-compile`
  - Output matches regex: `warning.*Framework.*missing|missing.*Framework section` (skill logs warning listing the missing Framework section)
  - Step 5: `grep` count >= 1 (generated test file uses Go testing + testify defaults)
- **Priority**: P1

---

## TC-003: Test-Guide Convention File Created from Existing Test Patterns

- **Source**: Story 2 / AC-1
- **Type**: CLI
- **Target**: cli/test-guide
- **Test ID**: cli/test-guide/convention-file-created-from-existing-tests
- **Pre-conditions**: A Go project exists with existing test files using testify assertions (e.g., `import "github.com/stretchr/testify/assert"`). No Convention file exists at `docs/conventions/testing-go.md`.
- **Steps**:
  1. Set up a project fixture with `*_test.go` files containing testify imports (`github.com/stretchr/testify/assert`) and `//go:build e2e` build tags
  2. Verify no `docs/conventions/testing-go.md` file exists (`test -f docs/conventions/testing-go.md` returns exit code 1)
  3. Invoke test-guide skill and confirm pattern extraction (or simulate by creating the expected Convention file content via the skill's documented output contract)
  4. Verify file `docs/conventions/testing-go.md` exists on disk (`test -f docs/conventions/testing-go.md` returns exit code 0)
  5. Verify the Convention file contains required sections by running `grep -c "^## Framework" docs/conventions/testing-go.md`, `grep -c "^## Assertion" docs/conventions/testing-go.md`, `grep -c "^## Tags" docs/conventions/testing-go.md`, `grep -c "^## Result Format" docs/conventions/testing-go.md`
  6. Verify detected framework is present: `grep "testify" docs/conventions/testing-go.md` returns exit code 0
- **Expected**:
  - Step 4: `docs/conventions/testing-go.md` file exists (exit code 0)
  - Step 5: All four `grep` commands return exit code 0 (each required section header found)
  - Step 6: Convention file contains `testify` in the Framework or Assertion section (exit code 0)
- **Priority**: P0
- **Note**: test-guide is a multi-turn conversational skill (FS-5), not a standalone forge CLI command. This TC verifies the observable file-system artifact produced by the skill. Execution requires a Claude Code agent session; it cannot be automated via forge binary alone.

---

## TC-004: Test-Guide Creates Convention Files for Multiple Languages in Mixed Project

- **Source**: Story 2 / AC-2
- **Type**: CLI
- **Target**: cli/test-guide
- **Test ID**: cli/test-guide/convention-files-for-multiple-languages
- **Pre-conditions**: A project exists with both `go.mod` and `package.json` in the root directory. No Convention files exist in `docs/conventions/`.
- **Steps**:
  1. Set up a project fixture with `go.mod` (Go module) and `package.json` (Node.js package) in the root directory
  2. Verify no Convention files exist: `ls docs/conventions/testing-*.md` returns exit code non-zero (no files)
  3. Invoke test-guide skill and confirm selection of both languages (or simulate by creating Convention files for both languages via the skill's documented output contract)
  4. Verify at least two Convention files exist: `ls docs/conventions/testing-*.md | wc -l` returns >= 2
  5. Verify Go Convention contains go-specific framework: `grep -l "go" docs/conventions/testing-*.md | head -1` and verify JavaScript Convention contains JS-specific content: `grep -l -i "javascript\|typescript\|vitest\|jest" docs/conventions/testing-*.md | head -1`
- **Expected**:
  - Step 4: At least 2 Convention files exist in `docs/conventions/`
  - Step 5: One Convention file references Go framework patterns, another references JavaScript/TypeScript framework patterns
- **Priority**: P1
- **Note**: test-guide is a multi-turn conversational skill (FS-5), not a standalone forge CLI command. This TC verifies observable file-system artifacts. Execution requires a Claude Code agent session.

---

## TC-005: Gen-Test-Scripts Proceeds Without Convention on Cold Start

- **Source**: Story 3 / AC-1
- **Type**: CLI
- **Target**: cli/gen-test-scripts
- **Test ID**: cli/gen-test-scripts/cold-start-no-convention
- **Pre-conditions**: A new Go project exists with no Convention files in `docs/conventions/` and no existing test files. The project has a valid Journey with Contract specs. `justfile` contains an `e2e-compile` recipe.
- **Steps**:
  1. Set up a clean Go project with no `docs/conventions/` directory and no `*_test.go` files
  2. Create a valid Journey and Contract spec
  3. Run `forge gen-test-scripts` targeting the Journey
  4. Capture stdout/stderr for hint output
  5. Run `just e2e-compile`
- **Expected**:
  - Exit code 0 from `just e2e-compile`
  - Output matches regex: `No test Convention files found|hint.*Convention|no.*Convention.*found` (hint about missing Convention files)
  - Generated test file compiles successfully using LLM defaults + Code Reconnaissance
- **Priority**: P0

---

## TC-006: Upgraded Forge Silently Ignores Legacy Config Fields

- **Source**: Story 4 / AC-1
- **Type**: CLI
- **Target**: cli/forge-commands
- **Test ID**: cli/forge-commands/backward-compat-ignores-legacy-config
- **Pre-conditions**: An existing forge project has `.forge/config.yaml` containing legacy fields: `languages`, `interfaces`, `test-framework`, `project-type`. Forge has been upgraded to the Convention-based version.
- **Steps**:
  1. Create a `.forge/config.yaml` with legacy fields: `languages: [go]`, `interfaces: [cli]`, `test-framework: go-testing`, `project-type: backend`
  2. Run `forge task index`
  3. Run `forge config init`
  4. Run `just e2e-test` (or equivalent test command)
  5. Capture stderr for any errors or warnings
- **Expected**:
  - Exit code 0 for all commands
  - stderr does NOT match regex: `languages|interfaces|test-framework|project-type|legacy.*field|deprecated.*field` (no errors or warnings referencing legacy fields)
  - All existing e2e tests pass via `just e2e-test`
  - `forge task index` and `forge config init` complete without referencing removed fields
- **Priority**: P0
- **Note**: The pre-condition "Forge has been upgraded" requires the Convention-based forge binary to be installed. In an isolated test fixture, this is satisfied by building from the current branch (`go build -o /tmp/forge-test ./cmd/forge` or equivalent) and ensuring the test PATH resolves to this binary. The TC does not require a version-migration scenario — only that the running binary lacks Profile support.

---

## TC-008: Gen-Test-Scripts Loads Convention by Interface Type Selectively

- **Source**: Story 5 / AC-1
- **Type**: CLI
- **Target**: cli/gen-test-scripts
- **Test ID**: cli/gen-test-scripts/selective-convention-loading-by-interface
- **Pre-conditions**: A project has two Convention files: `docs/conventions/testing-go.md` (domains: [testing, go]) and `docs/conventions/testing-javascript.md` (domains: [testing, javascript, web-ui]). A CLI Journey (Go) exists.
- **Steps**:
  1. Create both Convention files with distinct framework declarations
  2. Create a CLI Journey (Go interface)
  3. Run `forge gen-test-scripts` targeting the CLI Journey
  4. Run `grep -c "testing-go" <forge-log-file>` to verify `testing-go.md` Convention was loaded (exit code 0, count >= 1); alternatively run `grep "testing-go" <generated-test-path>` to verify Go Convention was applied
  5. Run `grep -c "vitest\|jest\|javascript" <generated-test-path>` to verify no JavaScript patterns exist (exit code 1, count 0)
  6. Run `grep "testing\|testify\|go\." <generated-test-path>` to verify Go-specific patterns are present (exit code 0)
- **Expected**:
  - Step 4: `testing-go.md` Convention is loaded (grep finds reference in logs or generated code)
  - Step 5: No JavaScript/Vitest patterns in generated test code (`grep -c` returns exit code 1)
  - Step 6: Go-specific patterns present (`grep` exits with code 0)
  - `just e2e-compile` passes (exit code 0)
- **Priority**: P0

---

## TC-009: Gen-Test-Scripts Loads Overlapping Domain Conventions and Logs Overlap Warning

- **Source**: Story 5 / AC-2
- **Type**: CLI
- **Target**: cli/gen-test-scripts
- **Test ID**: cli/gen-test-scripts/merges-overlapping-domain-conventions
- **Pre-conditions**: A project has two Convention files both including `testing` in their domains. The files declare conflicting assertion libraries (one says `assert`, the other says `require`). Files are named to control load order: `testing-00-assert.md` (alphabetically first) and `testing-01-require.md` (alphabetically second, assumed last-loaded on platforms with sorted directory listing).
- **Steps**:
  1. Create Convention file `docs/conventions/testing-00-assert.md`: domains [testing, go], Assertion: `assert (not require)`
  2. Create Convention file `docs/conventions/testing-01-require.md`: domains [testing, go, cli], Assertion: `require (not assert)`
  3. Create a valid Journey
  4. Run `forge gen-test-scripts` targeting the Journey
  5. Capture output for domain overlap log message
  6. Run `grep -c "assert\.\|require\." <generated-test-path>` to verify an assertion library is used
- **Expected**:
  - Both Convention files are loaded
  - Output matches regex: `overlap|domain.*overlap` (log note about overlapping domains)
  - Generated test file uses EITHER `assert` OR `require` assertions (merge resolves conflict to one library): `grep -c "assert\.\|require\." <generated-test-path>` returns >= 1
  - Note: "last-loaded wins" ordering depends on filesystem listing; this TC accepts either resolution but requires the overlap warning to be present
  - Generated test file compiles via `just e2e-compile` (exit code 0)
- **Priority**: P1

---

## TC-011: Convention File Missing Domains Frontmatter Treated as Non-Loadable

- **Source**: Spec FS-1 / Validation Rules
- **Type**: CLI
- **Target**: cli/gen-test-scripts
- **Test ID**: cli/gen-test-scripts/warns-on-missing-domains-frontmatter
- **Pre-conditions**: A Convention file exists at `docs/conventions/testing-go.md` without `domains` frontmatter. The project has a valid Journey. `justfile` contains `e2e-compile`.
- **Steps**:
  1. Create Convention file `docs/conventions/testing-go.md` with valid Framework/Assertion/Tags sections but no `domains` frontmatter
  2. Create a valid Journey and Contract spec
  3. Run `forge gen-test-scripts` targeting the Journey
  4. Capture output for warning
  5. Verify generation proceeds with LLM defaults
- **Expected**:
  - Output matches regex: `warning.*non-loadable|warning.*domains.*missing|Convention.*skipped` (warning about non-loadable Convention file)
  - Generation proceeds using LLM defaults (not the Convention file content)
  - `just e2e-compile` passes on generated output
- **Priority**: P1

---

## TC-012: Forge Commands Function Correctly After Profile Removal

- **Source**: Spec FS-7 / Import Audit Gate
- **Type**: CLI
- **Target**: cli/forge-commands
- **Test ID**: cli/forge-commands/commands-work-after-profile-removal
- **Pre-conditions**: The `pkg/profile/` directory has been removed. Forge binary is rebuilt and installed. A forge project exists with valid task structure, Journeys, and a `justfile` with `e2e-compile`.
- **Steps**:
  1. Run `forge task index` in the project
  2. Run `forge config init` in the project
  3. Run `forge gen-test-scripts` targeting a Journey
  4. Run `just e2e-compile`
- **Expected**:
  - All forge commands exit with code 0
  - `forge task index` lists discovered tasks (output contains at least one line matching `^[^/]+/[^/]+$`)
  - `forge config init` generates config without legacy fields: `grep -c "languages\|interfaces\|test-framework\|project-type" .forge/config.yaml` returns exit code 1
  - `just e2e-compile` passes on generated test file (exit code 0)
- **Priority**: P0

---

## TC-013: Config Init Works Without Legacy Fields

- **Source**: Spec FS-6 / FS-8
- **Type**: CLI
- **Target**: cli/config-init
- **Test ID**: cli/config-init/works-without-legacy-fields
- **Pre-conditions**: A new project directory exists with no `.forge/` directory. Forge CLI is installed.
- **Steps**:
  1. Run `forge config init` in the new project directory
  2. Inspect the generated `.forge/config.yaml`
  3. Verify no legacy fields (`languages`, `interfaces`, `test-framework`, `project-type`) are present
- **Expected**:
  - Exit code 0
  - Generated `config.yaml` contains only: `auto.*`, `worktree`, and/or `test-command` fields
  - No `languages`, `interfaces`, `test-framework`, or `project-type` fields in the output
- **Priority**: P0

---

## TC-014: Forge Task Index Works Without Profile Dependency

- **Source**: Spec FS-7 / Related Changes #5
- **Type**: CLI
- **Target**: cli/task-index
- **Test ID**: cli/task-index/works-without-profile-dependency
- **Pre-conditions**: A forge project exists with valid task structure. Profile package has been removed.
- **Steps**:
  1. Set up a project with features containing tasks
  2. Run `forge task index`
  3. Verify output lists all tasks correctly
- **Expected**:
  - Exit code 0
  - Output lists all discovered tasks from the project; each line matches pattern: `<feature-name>/<task-id>` or `<feature>/<task-title>` (regex: `^[^/]+/[^/]+$` with at least 1 match)
  - No errors or warnings about missing Profile
- **Priority**: P0

---

## TC-015: Init-Justfile Generates Recipes Without Profile Dependency

- **Source**: Spec FS-7 / Related Changes #4 (pkg/just/ — Convention-driven recipe generation)
- **Type**: CLI
- **Target**: cli/init-justfile
- **Test ID**: cli/init-justfile/generates-recipes-without-profile
- **Pre-conditions**: A forge project exists with no justfile or a justfile without e2e recipes. Profile package has been removed.
- **Steps**:
  1. Set up a project with no justfile
  2. Run `/forge:init-justfile`
  3. Run `grep -c "e2e-compile\|e2e-test\|e2e-setup" justfile` to verify recipes exist
  4. Run `just e2e-compile` with a generated test file
- **Expected**:
  - Justfile is generated with `e2e-compile`, `e2e-test`, `e2e-setup` recipes
  - Recipes use Convention + Code Reconnaissance for framework detection (not Profile)
  - `just e2e-compile` executes successfully (exit code 0) against a valid test file
- **Priority**: P0
- **Note**: init-justfile is invoked via Claude Code slash command (`/forge:init-justfile`), not a standalone `forge` CLI binary command. Execution requires a Claude Code agent session. The TC verifies observable file-system artifacts (justfile content) and runtime behavior (`just e2e-compile`).

---

## TC-016: No Convention Files Hint During Gen-Test-Scripts

- **Source**: Spec FS-2 / Error Handling
- **Type**: CLI
- **Target**: cli/gen-test-scripts
- **Test ID**: cli/gen-test-scripts/hint-when-no-convention-files
- **Pre-conditions**: A project exists with no files in `docs/conventions/`. The project has existing test files and a valid Journey. `justfile` contains `e2e-compile`.
- **Steps**:
  1. Remove all files from `docs/conventions/` (or ensure the directory is empty)
  2. Create a valid Journey and Contract spec
  3. Run `forge gen-test-scripts` targeting the Journey
  4. Capture output for hint message
- **Expected**:
  - Output matches regex: `No test Convention files found in docs/conventions/|hint.*test-guide` (hint message with test-guide suggestion)
  - Generation proceeds using LLM defaults + Code Reconnaissance
  - `just e2e-compile` passes
- **Priority**: P1

---

## TC-017: Convention Wins Over Conflicting Reconnaissance Signals

- **Source**: Spec FS-3 / Reliability Expectations
- **Type**: CLI
- **Target**: cli/gen-test-scripts
- **Test ID**: cli/gen-test-scripts/convention-overrides-reconnaissance-conflict
- **Pre-conditions**: A project has a Convention file declaring `assert` assertions. Existing test files in the project use `require` assertions. `justfile` contains `e2e-compile`.
- **Steps**:
  1. Create Convention file with `Assertion: assert (not require)`
  2. Create test files using `require.NoError` in the project
  3. Create a valid Journey
  4. Run `forge gen-test-scripts` targeting the Journey
  5. Capture logs for conflict notification
  6. Run `grep -c "assert\.\|assertNoError\|assertEqual" <generated-test-path>` to verify `assert` assertion style is used (exit code 0, count >= 1)
  7. Run `grep -c "require\.\|requireNoError" <generated-test-path>` to verify `require` assertions are NOT used (exit code 1 or count 0)
  8. Run `just e2e-compile`
- **Expected**:
  - Output matches regex: `conflict|Convention.*Reconnaissance|override.*detected` (conflict notification logged)
  - Step 6: `assert` assertions present in generated test file (Convention wins over Reconnaissance)
  - Step 7: `require` assertions absent from generated test file
  - Step 8: `just e2e-compile` passes (exit code 0)
- **Priority**: P1

---

## TC-018: Generation Time Does Not Exceed Absolute Budget Per Journey

- **Source**: Spec / Performance Requirements
- **Type**: CLI
- **Target**: cli/gen-test-scripts
- **Test ID**: cli/gen-test-scripts/generation-time-within-budget
- **Pre-conditions**: A forge-cli project exists with at least 10 Journeys. Convention files may or may not exist. `justfile` contains `e2e-compile`. No Profile package exists (removed per FS-7).
- **Steps**:
  1. Measure generation time for 10 representative Journeys using `time forge gen-test-scripts` for each
  2. Record the wall-clock time (seconds) for each generation
  3. Compute the average generation time across all 10 Journeys
- **Expected**:
  - Average generation time per Journey does not exceed 60 seconds (absolute budget; derived from Profile-era baseline of ~50s/Journey * 1.2 tolerance)
  - All generated tests pass `just e2e-compile` (exit code 0)
  - No single Journey exceeds 120 seconds
- **Priority**: P2

---

## TC-019: Gen-Test-Scripts Produces Compilable Output on First Attempt for Standard Go Project

- **Source**: Spec / Goals + Performance Requirements
- **Type**: CLI
- **Target**: cli/gen-test-scripts
- **Test ID**: cli/gen-test-scripts/first-pass-compile-for-standard-project
- **Pre-conditions**: A standard Go project exists using go-testing + testify (the default framework combination). A Convention file exists at `docs/conventions/testing-go.md` declaring the standard framework. A valid Journey with Contract specs exists. `justfile` contains `e2e-compile`.
- **Steps**:
  1. Create Convention file `docs/conventions/testing-go.md` with Framework: Go testing + testify, Assertion: assert, Tags: `//go:build e2e`, Result Format: json-stream
  2. Create a valid Journey and Contract spec for a CLI interface
  3. Run `forge gen-test-scripts` targeting the Journey
  4. Immediately run `just e2e-compile` (no manual edits, no retries)
  5. Capture exit code and output from `just e2e-compile`
- **Expected**:
  - `forge gen-test-scripts` exits with code 0
  - `just e2e-compile` exits with code 0 on first attempt (no compile gate retries needed)
  - Generated test file contains `testing` and `testify` imports (verified by `grep`)
- **Priority**: P0

---

## TC-020: Consolidate-Specs Detects Convention Drift

- **Source**: Spec FS-9 / Drift Detection
- **Type**: CLI
- **Target**: cli/consolidate-specs
- **Test ID**: cli/consolidate-specs/detects-convention-drift
- **Pre-conditions**: A project has a Convention file declaring `assert (not require)`. Existing test files use `require` assertions. `consolidate-specs` is configured to run drift audits.
- **Steps**:
  1. Create Convention file with `Assertion: assert (not require)`
  2. Create test files using `require.NoError` in the project
  3. Run `forge consolidate-specs` (drift detection is performed by default per FS-9)
  4. Capture drift report output
- **Expected**:
  - Output matches regex: `Convention.*assert.*conflict.*test.*require|drift.*Convention:.*assert.*actual:.*require` (drift report flags the assert/require mismatch)
  - Output includes the Convention file path: regex `docs/conventions/.*\.md` matches
  - Output includes the conflicting test file path: regex `_test\.go` matches
- **Priority**: P2

---

## TC-021: Removed Command forge test detect Returns Error

- **Source**: PRD Scope / Remove CLI commands
- **Type**: CLI
- **Target**: cli/forge-test-detect
- **Test ID**: cli/forge-test-detect/removed-command-errors
- **Pre-conditions**: Forge CLI is installed with the Convention-based version (Profile removed).
- **Steps**:
  1. Run `forge test detect`
  2. Capture exit code and stderr
- **Expected**:
  - Exit code 2 (BIZ-error-reporting-001: usage error for unknown/removed command)
  - stderr matches regex: `unknown command|command not found|removed`
- **Priority**: P0

---

## TC-022: Removed Command forge test get Returns Error

- **Source**: PRD Scope / Remove CLI commands
- **Type**: CLI
- **Target**: cli/forge-test-get
- **Test ID**: cli/forge-test-get/removed-command-errors
- **Pre-conditions**: Forge CLI is installed with the Convention-based version (Profile removed).
- **Steps**:
  1. Run `forge test get`
  2. Capture exit code and stderr
- **Expected**:
  - Exit code 2 (usage error for unknown/removed command)
  - stderr matches regex: `unknown command|command not found|removed`
- **Priority**: P0

---

## TC-023: Removed Command forge test interfaces Returns Error

- **Source**: PRD Scope / Remove CLI commands
- **Type**: CLI
- **Target**: cli/forge-test-interfaces
- **Test ID**: cli/forge-test-interfaces/removed-command-errors
- **Pre-conditions**: Forge CLI is installed with the Convention-based version (Profile removed).
- **Steps**:
  1. Run `forge test interfaces`
  2. Capture exit code and stderr
- **Expected**:
  - Exit code 2 (usage error for unknown/removed command)
  - stderr matches regex: `unknown command|command not found|removed`
- **Priority**: P0

---

## TC-024: Removed Command forge test framework Returns Error

- **Source**: PRD Scope / Remove CLI commands
- **Type**: CLI
- **Target**: cli/forge-test-framework
- **Test ID**: cli/forge-test-framework/removed-command-errors
- **Pre-conditions**: Forge CLI is installed with the Convention-based version (Profile removed).
- **Steps**:
  1. Run `forge test framework`
  2. Capture exit code and stderr
- **Expected**:
  - Exit code 2 (usage error for unknown/removed command)
  - stderr matches regex: `unknown command|command not found|removed`
- **Priority**: P0

---

## TC-025: Forge Task Add Works Without Profile Dependency

- **Source**: Spec FS-7 / Related Changes #5 (internal/cmd/ rewrite)
- **Type**: CLI
- **Target**: cli/task-add
- **Test ID**: cli/task-add/works-without-profile-dependency
- **Pre-conditions**: A forge project exists with valid task structure. Profile package has been removed.
- **Steps**:
  1. Set up a project with features directory containing at least one task
  2. Run `forge task add --title "test-task" --description "test description"`
  3. Capture exit code and output
  4. Verify the task was added by running `forge task index` and checking output contains "test-task"
- **Expected**:
  - `forge task add` exits with code 0
  - Output confirms task creation
  - `forge task index` output includes the newly added task
  - No errors or warnings about missing Profile
- **Priority**: P0

---

## TC-026: Forge Init Creates Project Without Legacy Fields

- **Source**: Spec FS-7 / Related Changes #5 (internal/cmd/ — Simplify forge init)
- **Type**: CLI
- **Target**: cli/init
- **Test ID**: cli/init/creates-project-without-legacy-fields
- **Pre-conditions**: An empty directory exists. Forge CLI is installed.
- **Steps**:
  1. Run `forge init` in the empty directory
  2. Capture exit code and output
  3. Verify `.forge/` directory was created: `test -d .forge` returns exit code 0
  4. If `config.yaml` was generated, verify it contains no legacy fields: `grep -c "languages\|interfaces\|test-framework\|project-type" .forge/config.yaml` returns exit code 1 (no matches)
- **Expected**:
  - `forge init` exits with code 0
  - `.forge/` directory exists
  - Generated `config.yaml` (if present) contains no legacy fields (`languages`, `interfaces`, `test-framework`, `project-type`)
- **Priority**: P0

---

## TC-027: Run-E2E-Tests Parses Results Using Convention Result Format

- **Source**: PRD Scope / Rewrite run-e2e-tests skill
- **Type**: CLI
- **Target**: cli/run-e2e-tests
- **Test ID**: cli/run-e2e-tests/parses-results-with-convention-format
- **Pre-conditions**: A forge project exists with a Convention file declaring `Result Format: json-stream`. A valid Journey with Contract specs exists. `justfile` contains `e2e-test` recipe.
- **Steps**:
  1. Create Convention file with `Result Format: json-stream`
  2. Generate test scripts using `forge gen-test-scripts` for a Journey
  3. Run `just e2e-test` (or equivalent e2e test execution) and capture output
  4. Run `/forge:run-e2e-tests` to parse the test results
  5. Verify the parsed results report matches the actual test outcomes
- **Expected**:
  - run-e2e-tests exits with code 0 (if tests pass) or code 1 (if tests fail, but parsing succeeds)
  - Results report includes per-Journey pass/fail status
  - Results report correctly interprets json-stream format output: report contains fields matching regex `"Test":|"Action":|"Elapsed":` (json-stream fields present), and does NOT contain regex `^--- PASS:|^--- FAIL:|^ok ` (go-test format fields absent)
- **Priority**: P0
- **Note**: Step 4 invokes the run-e2e-tests Claude Code skill (`/forge:run-e2e-tests`), not a standalone forge CLI binary command. Execution requires a Claude Code agent session. Steps 1-3 are forge CLI binary commands; step 4 is the skill invocation.

---

## TC-028: Convention File Unreadable Due to Permissions Is Skipped with Warning

- **Source**: Spec FS-2 / Error Handling
- **Type**: CLI
- **Target**: cli/gen-test-scripts
- **Test ID**: cli/gen-test-scripts/warns-on-unreadable-convention-file
- **Pre-conditions**: A Convention file exists at `docs/conventions/testing-go.md` but has file permissions set to 000 (no read access). The project has a valid Journey. `justfile` contains `e2e-compile`.
- **Steps**:
  1. Create a valid Convention file at `docs/conventions/testing-go.md`
  2. Set file permissions to unreadable: `chmod 000 docs/conventions/testing-go.md`
  3. Run `forge gen-test-scripts` targeting the Journey
  4. Capture output for warning
  5. Run `just e2e-compile`
  6. Restore permissions: `chmod 644 docs/conventions/testing-go.md`
- **Expected**:
  - Output matches regex: `warning.*permission|warning.*unreadable|skip.*Convention` (warning about unreadable Convention file with file path)
  - Generation proceeds using LLM defaults (skips the unreadable file)
  - `just e2e-compile` passes (exit code 0)
- **Priority**: P1

---

## TC-029: End-to-End Flow — Convention Creation Through Test Execution

- **Source**: PRD / Main flow (Business Flow Description)
- **Type**: CLI
- **Target**: cli/gen-test-scripts (integration)
- **Test ID**: cli/integration/e2e-convention-through-execution
- **Pre-conditions**: A standard Go project exists using go-testing + testify. No Convention files exist. A valid Journey with Contract specs exists. `justfile` contains `e2e-compile` and `e2e-test` recipes.
- **Steps**:
  1. Verify no Convention files exist: `ls docs/conventions/testing-*.md 2>/dev/null` returns exit code non-zero
  2. Run `forge gen-test-scripts` targeting the Journey (cold start — no Convention)
  3. Capture output; verify hint message about missing Convention is present
  4. Run `just e2e-compile`; verify exit code 0 (generated code compiles)
  5. Create a Convention file `docs/conventions/testing-go.md` with standard Go testing + testify declarations
  6. Run `forge gen-test-scripts` targeting the same Journey again (with Convention)
  7. Run `just e2e-compile`; verify exit code 0
  8. Verify the Convention-based generation produces consistent output: `grep -c "testify" <generated-test-file>` returns >= 1
- **Expected**:
  - Step 2: `forge gen-test-scripts` exits with code 0, output contains hint about missing Convention
  - Step 4: `just e2e-compile` exits with code 0 (cold start generation compiles)
  - Step 7: `just e2e-compile` exits with code 0 (Convention-based generation compiles)
  - Step 8: Generated test file uses Convention-declared testify assertions
- **Priority**: P0

---

## TC-030: Test-Guide Presents Framework Candidates When No Test Files Exist

- **Source**: Story 2 / AC-1 (cold start path), PRD FS-5
- **Type**: CLI
- **Target**: cli/test-guide
- **Test ID**: cli/test-guide/presents-candidates-when-no-test-files
- **Pre-conditions**: A project exists with `go.mod` (recognizable language signal) but no `*_test.go` files. No Convention files exist in `docs/conventions/`.
- **Steps**:
  1. Set up a project with `go.mod` but no test files (`find . -name "*_test.go"` returns no results)
  2. Verify no Convention files: `ls docs/conventions/testing-*.md 2>/dev/null` returns exit code non-zero
  3. Invoke test-guide skill
  4. Verify the skill presents candidate frameworks for user selection (e.g., "go-testing", "ginkgo")
  5. Select a framework and verify Convention file is created: `test -f docs/conventions/testing-go.md` returns exit code 0
- **Expected**:
  - Step 4: Skill output includes candidate framework names (not auto-selecting a single framework)
  - Step 5: `docs/conventions/testing-go.md` is created with user-selected framework in the Framework section
- **Priority**: P1
- **Note**: test-guide is a multi-turn conversational skill (FS-5), not a standalone forge CLI command. Execution requires a Claude Code agent session.

---

## TC-031: Compile Gate Intermediate Retry Feeds Error Back to LLM

- **Source**: Spec FS-4 / Retry Semantics, PRD Main Flow step 6
- **Type**: CLI
- **Target**: cli/gen-test-scripts
- **Test ID**: cli/gen-test-scripts/compile-gate-intermediate-retry-feedback
- **Pre-conditions**: A Go project exists with a Convention file declaring a slightly incorrect framework detail (e.g., declares wrong import path for testify). The project has a valid Journey. `justfile` contains `e2e-compile`. The first compile attempt will fail but the error provides enough information for LLM to correct on retry.
- **Steps**:
  1. Create Convention file with a minor framework import error (e.g., `github.com/stretchr/testify/require` declared but project uses `assert` style)
  2. Create a valid Journey and Contract spec
  3. Run `forge gen-test-scripts` targeting the Journey
  4. Capture full stdout/stderr output including retry attempts
  5. Verify the generated file compiles after retry (not necessarily first attempt)
- **Expected**:
  - Output contains evidence of compile error feedback being used: regex matches `retry|attempt.*2|regenerat|re-generat`
  - Final `just e2e-compile` exit code 0 (retry succeeded after error feedback)
  - Generated test file exists on disk
- **Priority**: P1

---

## TC-032: Convention File with All Required Sections Missing Falls Back Fully

- **Source**: Spec FS-1 / Validation Rules (boundary case)
- **Type**: CLI
- **Target**: cli/gen-test-scripts
- **Test ID**: cli/gen-test-scripts/fallback-when-all-sections-missing
- **Pre-conditions**: A Convention file exists at `docs/conventions/testing-go.md` with valid `domains` frontmatter but ALL four required sections (`Framework`, `Assertion`, `Tags`, `Result Format`) are missing. The project has a valid Journey. `justfile` contains `e2e-compile`.
- **Steps**:
  1. Create Convention file with `domains: [testing, go]` frontmatter but no `## Framework`, `## Assertion`, `## Tags`, or `## Result Format` section headers
  2. Create a valid Journey and Contract spec
  3. Run `forge gen-test-scripts` targeting the Journey
  4. Capture output for warnings
  5. Run `just e2e-compile`
- **Expected**:
  - Output matches regex: `warning.*missing.*section|missing.*Framework|missing.*Assertion|missing.*Tags|missing.*Result Format` (warnings listing all missing sections)
  - Generation proceeds fully with LLM defaults (Convention provides no usable content)
  - `just e2e-compile` passes (exit code 0)
- **Priority**: P1

---

## TC-034: Convention File with Invalid Encoding Is Skipped with Warning

- **Source**: Spec FS-2 / Error Handling (encoding)
- **Type**: CLI
- **Target**: cli/gen-test-scripts
- **Test ID**: cli/gen-test-scripts/warns-on-invalid-encoding-convention-file
- **Pre-conditions**: A Convention file exists at `docs/conventions/testing-go.md` containing invalid UTF-8 byte sequences (e.g., raw `\xff\xfe` BOM or binary garbage). The project has a valid Journey. `justfile` contains `e2e-compile`.
- **Steps**:
  1. Create a Convention file at `docs/conventions/testing-go.md` with binary content: `printf '\xff\xfe invalid binary content' > docs/conventions/testing-go.md`
  2. Create a valid Journey and Contract spec
  3. Run `forge gen-test-scripts` targeting the Journey
  4. Capture output for warning
  5. Run `just e2e-compile`
- **Expected**:
  - Output matches regex: `warning.*encoding|warning.*unreadable|skip.*Convention` (warning about unparseable Convention file with file path)
  - Generation proceeds using LLM defaults (skips the unparseable file)
  - `just e2e-compile` passes (exit code 0)
- **Priority**: P1

---

## TC-035: Gen-Test-Scripts Warns on Invalid Section Content (Empty Framework Name)

- **Source**: Spec FS-1 / Validation Rules (Invalid section content)
- **Type**: CLI
- **Target**: cli/gen-test-scripts
- **Test ID**: cli/gen-test-scripts/warns-on-invalid-section-content
- **Pre-conditions**: A Convention file exists at `docs/conventions/testing-go.md` with a `## Framework` section whose name field is empty (e.g., `## Framework\n- name:\n- File pattern: *_test.go`). The project has a valid Journey. `justfile` contains `e2e-compile`.
- **Steps**:
  1. Create Convention file `docs/conventions/testing-go.md` with `domains: [testing, go]`, `## Framework` section containing an empty `name:` field, and valid `## Assertion`, `## Tags`, `## Result Format` sections
  2. Create a valid Journey and Contract spec
  3. Run `forge gen-test-scripts` targeting the Journey
  4. Capture stdout/stderr for warning output
  5. Run `just e2e-compile`
- **Expected**:
  - Output matches regex: `warning.*Framework.*invalid|warning.*Framework.*empty|treat.*Framework.*as missing` (skill logs warning about invalid Framework section content, treats as missing)
  - Generation proceeds with LLM defaults for framework selection (ignores the empty Framework name)
  - `just e2e-compile` passes (exit code 0)
- **Priority**: P1

---

## TC-036: Integration — Test-Guide Creates Convention Consumed by Gen-Test-Scripts

- **Source**: PRD / Main flow + Bootstrap flow (FS-5 → FS-2 → FS-4)
- **Type**: CLI
- **Target**: cli/integration
- **Test ID**: cli/integration/test-guide-to-gen-test-scripts-pipeline
- **Pre-conditions**: A standard Go project exists using go-testing + testify with existing test files (`*_test.go` containing `assert.NoError` patterns). No Convention files exist. A valid Journey with Contract specs exists. `justfile` contains `e2e-compile`. Forge binary is built from the Convention-based source.
- **Steps**:
  1. Verify no Convention files exist: `ls docs/conventions/testing-*.md 2>/dev/null` returns exit code non-zero
  2. Invoke test-guide skill and confirm framework detection (or simulate by creating the Convention file with testify content per the skill's documented output contract)
  3. Verify Convention file created: `test -f docs/conventions/testing-go.md` returns exit code 0
  4. Verify Convention contains required sections: `grep -c "^## Framework" docs/conventions/testing-go.md` returns >= 1, `grep -c "^## Assertion" docs/conventions/testing-go.md` returns >= 1
  5. Run `forge gen-test-scripts` targeting the Journey
  6. Run `grep -c "testify\|assert\." <generated-test-path>` to verify Convention-declared framework is applied (count >= 1)
  7. Run `just e2e-compile`
- **Expected**:
  - Step 3: Convention file exists (exit code 0)
  - Step 4: Required sections present (both grep counts >= 1)
  - Step 6: Generated test file uses Convention-declared testify/assert patterns (grep count >= 1)
  - Step 7: `just e2e-compile` exits with code 0 (Convention from test-guide produces compilable output)
- **Priority**: P0
- **Note**: This integration TC chains the test-guide Convention creation (TC-003) with gen-test-scripts Convention consumption (TC-001) in a single pipeline. Steps 2-3 require a Claude Code agent session for test-guide invocation; steps 5-7 are forge CLI binary commands.

---

## Traceability

| TC ID | Source | Type | Target | Priority |
|-------|--------|------|--------|----------|
| TC-001 | Story 1 / AC-1 | CLI | cli/gen-test-scripts | P0 |
| TC-002 | Story 1 / AC-2 | CLI | cli/gen-test-scripts | P1 |
| TC-003 | Story 2 / AC-1 | CLI | cli/test-guide | P0 |
| TC-004 | Story 2 / AC-2 | CLI | cli/test-guide | P1 |
| TC-005 | Story 3 / AC-1 | CLI | cli/gen-test-scripts | P0 |
| TC-006 | Story 4 / AC-1 | CLI | cli/forge-commands | P0 |
| TC-008 | Story 5 / AC-1 | CLI | cli/gen-test-scripts | P0 |
| TC-009 | Story 5 / AC-2 | CLI | cli/gen-test-scripts | P1 |
| TC-011 | Spec FS-1 / Validation | CLI | cli/gen-test-scripts | P1 |
| TC-012 | Spec FS-7 / Import Audit | CLI | cli/forge-commands | P0 |
| TC-013 | Spec FS-6 / FS-8 | CLI | cli/config-init | P0 |
| TC-014 | Spec FS-7 / Related #5 | CLI | cli/task-index | P0 |
| TC-015 | Spec FS-7 / Related #4 | CLI | cli/init-justfile | P0 |
| TC-016 | Spec FS-2 / Error Handling | CLI | cli/gen-test-scripts | P1 |
| TC-017 | Spec FS-3 / Reliability | CLI | cli/gen-test-scripts | P1 |
| TC-018 | Spec / Performance | CLI | cli/gen-test-scripts | P2 |
| TC-019 | Spec / Goals | CLI | cli/gen-test-scripts | P0 |
| TC-020 | Spec FS-9 / Drift Detection | CLI | cli/consolidate-specs | P2 |
| TC-021 | PRD Scope / Remove commands | CLI | cli/forge-test-detect | P0 |
| TC-022 | PRD Scope / Remove commands | CLI | cli/forge-test-get | P0 |
| TC-023 | PRD Scope / Remove commands | CLI | cli/forge-test-interfaces | P0 |
| TC-024 | PRD Scope / Remove commands | CLI | cli/forge-test-framework | P0 |
| TC-025 | Spec FS-7 / Related #5 | CLI | cli/task-add | P0 |
| TC-026 | Spec FS-7 / Related #5 | CLI | cli/init | P0 |
| TC-027 | PRD Scope / Rewrite skill | CLI | cli/run-e2e-tests | P0 |
| TC-028 | Spec FS-2 / Error Handling | CLI | cli/gen-test-scripts | P1 |
| TC-029 | PRD / Main flow | CLI | cli/integration | P0 |
| TC-030 | Story 2 / AC-1 (cold start), FS-5 | CLI | cli/test-guide | P1 |
| TC-031 | Spec FS-4 / Retry | CLI | cli/gen-test-scripts | P1 |
| TC-032 | Spec FS-1 / Validation (boundary) | CLI | cli/gen-test-scripts | P1 |
| TC-034 | Spec FS-2 / Error Handling (encoding) | CLI | cli/gen-test-scripts | P1 |
| TC-035 | Spec FS-1 / Validation Rules (invalid content) | CLI | cli/gen-test-scripts | P1 |
| TC-036 | PRD / Main flow + Bootstrap flow | CLI | cli/integration | P0 |
