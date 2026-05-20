---
feature: "test-knowledge-convention-driven"
type: CLI
generated: "2026-05-20"
---

# CLI Test Cases: test-knowledge-convention-driven

Test cases derived from PRD acceptance criteria. Every case traces to a specific PRD source.

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
  5. Inspect the generated test file for import statements and assertion patterns
  6. Run `just e2e-compile`
- **Expected**:
  - Exit code 0 from `just e2e-compile`
  - Generated test file contains ginkgo-style imports (e.g., `github.com/onsi/ginkgo/v2`)
  - Generated test file uses gomega assertions (e.g., `Expect(...).To(...)`) not testify assertions (e.g., `assert.NoError`)
  - Generated test file includes `//go:build e2e` build tag
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
  5. Inspect the generated test file
  6. Run `just e2e-compile`
- **Expected**:
  - Exit code 0 from `just e2e-compile`
  - Output contains a warning listing the missing Framework section
  - Generated test file uses LLM-detected framework defaults (Go testing + testify for a standard Go project)
- **Priority**: P1

---

## TC-003: Test-Guide Scans Existing Tests and Generates Convention File

- **Source**: Story 2 / AC-1
- **Type**: CLI
- **Target**: cli/test-guide
- **Test ID**: cli/test-guide/scans-tests-and-generates-convention
- **Pre-conditions**: A Go project exists with existing test files using testify assertions (e.g., `import "github.com/stretchr/testify/assert"`). No Convention file exists at `docs/conventions/testing-go.md`.
- **Steps**:
  1. Set up a project fixture with `*_test.go` files containing testify imports and assertion patterns
  2. Run `/forge:test-guide` interactively
  3. Verify the skill scans test files and extracts framework patterns (imports, assertions, tags)
  4. Confirm the extracted patterns when presented
  5. Verify the output Convention file
- **Expected**:
  - Skill output lists detected patterns: `assert` library from testify, Go testing package, build tags if present
  - File `docs/conventions/testing-go.md` is created with minimum sections: Framework (Go testing + testify), Assertion (assert.NoError, assert.Contains), Tags (//go:build e2e), Result Format
- **Priority**: P0

---

## TC-004: Test-Guide Handles Ambiguous Language Signals

- **Source**: Story 2 / AC-2
- **Type**: CLI
- **Target**: cli/test-guide
- **Test ID**: cli/test-guide/handles-ambiguous-language-signals
- **Pre-conditions**: A project exists with both `go.mod` and `package.json` in the root directory. No Convention files exist.
- **Steps**:
  1. Set up a project fixture with both `go.mod` (Go module) and `package.json` (Node.js package)
  2. Run `/forge:test-guide` interactively
  3. Verify the skill detects both language candidates
  4. Select one or both languages when prompted
- **Expected**:
  - Skill output lists detected language candidates: Go and JavaScript/TypeScript
  - Skill asks user to select which language(s) to generate Conventions for
  - After selection, Convention files are generated for the selected language(s)
- **Priority**: P1

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
  - Output contains a hint: "No test Convention files found" (or equivalent)
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
  - No errors or warnings related to the legacy fields
  - All existing e2e tests pass via `just e2e-test`
  - `forge task index` and `forge config init` complete without referencing removed fields
- **Priority**: P0

---

## TC-007: Gen-Test-Scripts Reports Missing E2E-Compile Recipe

- **Source**: Story 4 / AC-2
- **Type**: CLI
- **Target**: cli/gen-test-scripts
- **Test ID**: cli/gen-test-scripts/reports-missing-compile-recipe
- **Pre-conditions**: A forge project exists with a valid Journey and Contract specs. The `justfile` does NOT contain an `e2e-compile` recipe (or no justfile exists).
- **Steps**:
  1. Set up a project without `e2e-compile` recipe in the justfile
  2. Create a valid Journey and Contract spec
  3. Run `forge gen-test-scripts` targeting the Journey
  4. Capture stdout/stderr for error output
- **Expected**:
  - Exit code non-zero
  - Output contains: "Missing justfile e2e-compile recipe. Run `/forge:init-justfile` first, or add a recipe manually."
- **Priority**: P1

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
  4. Inspect which Convention file was loaded (via logs or generated code)
  5. Inspect the generated test file
- **Expected**:
  - Only `testing-go.md` Convention is loaded and applied
  - Generated test code uses Go-specific patterns (Go testing package, not JavaScript/Vitest)
  - No JavaScript-specific patterns appear in the generated CLI test file
- **Priority**: P0

---

## TC-009: Gen-Test-Scripts Merges Overlapping Domain Conventions with Warning

- **Source**: Story 5 / AC-2
- **Type**: CLI
- **Target**: cli/gen-test-scripts
- **Test ID**: cli/gen-test-scripts/merges-overlapping-domain-conventions
- **Pre-conditions**: A project has two Convention files both including `testing` in their domains. The files declare conflicting assertion libraries (one says `assert`, the other says `require`).
- **Steps**:
  1. Create Convention file A: domains [testing, go], Assertion: `assert (not require)`
  2. Create Convention file B: domains [testing, go, cli], Assertion: `require (not assert)`
  3. Create a valid Journey
  4. Run `forge gen-test-scripts` targeting the Journey
  5. Capture output for domain overlap log message
  6. Inspect the generated test file
- **Expected**:
  - Both Convention files are loaded
  - Output contains a note about domain overlap
  - Last-loaded Convention wins for conflicting fields (e.g., `require` assertions if file B is loaded after file A)
  - Generated test file compiles via `just e2e-compile`
- **Priority**: P1

---

## TC-010: Compile Gate Recovery Outputs Actionable Guidance on Exhausted Retries

- **Source**: Story 6 / AC-1
- **Type**: CLI
- **Target**: cli/gen-test-scripts
- **Test ID**: cli/gen-test-scripts/compile-gate-recovery-exhausted-retries
- **Pre-conditions**: A Go project exists with a Convention file that declares an incorrect framework (e.g., declares ginkgo but project has no ginkgo dependency). The project has a valid Journey. `justfile` contains `e2e-compile`. The compile gate will fail consistently.
- **Steps**:
  1. Create a Convention file with incorrect framework declarations (e.g., declares ginkgo imports but `go.mod` has no ginkgo dependency)
  2. Create a valid Journey and Contract spec
  3. Run `forge gen-test-scripts` targeting the Journey
  4. Allow all compile gate retries to exhaust (max 2 retries)
  5. Capture final output
  6. Verify the generated file still exists on disk
- **Expected**:
  - Output contains: (1) the compile error message, (2) the generated file path, (3) recovery actions: check Convention, run `/forge:test-guide`, or manually edit
  - The generated test file is preserved on disk (not deleted)
  - Exit code indicates failure/blocked status
- **Priority**: P0

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
  - Output contains warning about non-loadable Convention file
  - Generation proceeds using LLM defaults (not the Convention file content)
  - `just e2e-compile` passes on generated output
- **Priority**: P1

---

## TC-012: Profile Package Has Zero Import References After Removal

- **Source**: Spec FS-7 / Import Audit Gate
- **Type**: CLI
- **Target**: cli/forge-build
- **Test ID**: cli/forge-build/zero-profile-imports-after-removal
- **Pre-conditions**: The `pkg/profile/` directory has been removed. All consumer rewrites are complete.
- **Steps**:
  1. Run `grep -r "pkg/profile" forge-cli/` in the project root
  2. Run `go build ./...` to verify compilation
  3. Run existing test suite
- **Expected**:
  - `grep -r "pkg/profile" forge-cli/` returns zero results (no output, exit code 1)
  - `go build ./...` exits with code 0
  - All existing tests pass
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
  - Output lists all discovered tasks from the project
  - No errors or warnings about missing Profile
- **Priority**: P0

---

## TC-015: Init-Justfile Generates Recipes Without Profile Dependency

- **Source**: Spec FS-7 / Related Changes
- **Type**: CLI
- **Target**: cli/init-justfile
- **Test ID**: cli/init-justfile/generates-recipes-without-profile
- **Pre-conditions**: A forge project exists with no justfile or a justfile without e2e recipes. Profile package has been removed.
- **Steps**:
  1. Set up a project with no justfile
  2. Run `/forge:init-justfile`
  3. Inspect the generated justfile
  4. Run `just e2e-compile` with a generated test file
- **Expected**:
  - Justfile is generated with `e2e-compile`, `e2e-test`, `e2e-setup` recipes
  - Recipes use Convention + Code Reconnaissance for framework detection (not Profile)
  - `just e2e-compile` executes successfully (exit code 0) against a valid test file
- **Priority**: P0

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
  - Output contains hint: "No test Convention files found in docs/conventions/." and suggests running `/forge:test-guide`
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
  6. Inspect the generated test file
- **Expected**:
  - Output logs a conflict notification between Convention and Reconnaissance
  - Generated test file uses `assert` assertions (Convention wins)
  - `just e2e-compile` passes
- **Priority**: P1

---

## TC-018: Generation Time Does Not Increase Beyond 20% vs Profile Baseline

- **Source**: Spec / Performance Requirements
- **Type**: CLI
- **Target**: cli/gen-test-scripts
- **Test ID**: cli/gen-test-scripts/generation-time-within-budget
- **Pre-conditions**: A forge-cli project exists with 126+ existing test Journeys. Baseline generation time with Profile is recorded.
- **Steps**:
  1. Record the baseline generation time for a representative set of Journeys using Profile-based generation
  2. Run the same set of Journeys using Convention-based generation (gen-test-scripts)
  3. Compare generation times
- **Expected**:
  - Convention-based generation time is within 120% of Profile-based generation time (i.e., no more than 20% increase)
  - All generated tests pass `just e2e-compile`
- **Priority**: P2

---

## TC-019: First-Pass Compile Rate Meets 85% Threshold

- **Source**: Spec / Goals + Performance Requirements
- **Type**: CLI
- **Target**: cli/gen-test-scripts
- **Test ID**: cli/gen-test-scripts/first-pass-compile-rate-threshold
- **Pre-conditions**: The forge-cli project with 126+ existing test Journeys. Convention files may or may not exist.
- **Steps**:
  1. Run gen-test-scripts for all 126+ existing Journeys
  2. Count how many pass `just e2e-compile` on first attempt (no retries)
  3. Calculate first-pass compile rate
- **Expected**:
  - First-pass compile rate >= 85% (at least 107 out of 126 Journeys compile on first attempt)
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
  3. Run consolidate-specs with drift detection enabled
  4. Capture drift report output
- **Expected**:
  - Drift report flags the mismatch between Convention (assert) and actual test usage (require)
  - Report includes the Convention file path and the conflicting test file
- **Priority**: P2

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
| TC-007 | Story 4 / AC-2 | CLI | cli/gen-test-scripts | P1 |
| TC-008 | Story 5 / AC-1 | CLI | cli/gen-test-scripts | P0 |
| TC-009 | Story 5 / AC-2 | CLI | cli/gen-test-scripts | P1 |
| TC-010 | Story 6 / AC-1 | CLI | cli/gen-test-scripts | P0 |
| TC-011 | Spec FS-1 / Validation | CLI | cli/gen-test-scripts | P1 |
| TC-012 | Spec FS-7 / Import Audit | CLI | cli/forge-build | P0 |
| TC-013 | Spec FS-6 / FS-8 | CLI | cli/config-init | P0 |
| TC-014 | Spec FS-7 / Related #5 | CLI | cli/task-index | P0 |
| TC-015 | Spec FS-7 / Related | CLI | cli/init-justfile | P0 |
| TC-016 | Spec FS-2 / Error Handling | CLI | cli/gen-test-scripts | P1 |
| TC-017 | Spec FS-3 / Reliability | CLI | cli/gen-test-scripts | P1 |
| TC-018 | Spec / Performance | CLI | cli/gen-test-scripts | P2 |
| TC-019 | Spec / Goals | CLI | cli/gen-test-scripts | P0 |
| TC-020 | Spec FS-9 / Drift Detection | CLI | cli/consolidate-specs | P2 |
