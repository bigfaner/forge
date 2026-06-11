---
feature: "test-knowledge-convention-driven"
---

# User Stories: test-knowledge-convention-driven

## Story 1: Non-Default Framework Support

**As a** Forge CLI user using a non-default test framework (e.g., ginkgo, vitest)
**I want to** generate e2e tests that use my project's actual framework instead of the Profile-detected default
**So that** generated test code compiles and runs correctly without manual post-generation edits

**Acceptance Criteria:**
- Given a Go project using ginkgo (not go-testing), and a Convention file declaring ginkgo in `docs/conventions/testing-go.md`
- When the user runs gen-test-scripts for a Journey
- Then the generated test code uses ginkgo imports, assertions, and style conventions as declared in the Convention file
- And `just e2e-compile` passes on first attempt

**Error scenarios:**
- Given a Convention file with an empty Framework section
- When gen-test-scripts loads the Convention
- Then the skill outputs a warning listing the missing section, and proceeds with LLM defaults for framework detection
- And the generated code compiles via `just e2e-compile`

---

## Story 2: Convention Bootstrap

**As a** Forge CLI user starting with a new project
**I want to** create a test Convention file with minimal effort through guided detection
**So that** subsequent test generations produce consistent, framework-accurate code

**Acceptance Criteria:**
- Given a project with existing test files (e.g., `*_test.go` using testify)
- When the user runs `/forge:test-guide`
- Then the skill scans existing test files, extracts framework patterns (imports, assertions, tags), and presents them for confirmation
- And after confirmation, writes `docs/conventions/testing-<scope>.md` with the minimum set (Framework + Assertion + Tags + Result Format)

**Error scenarios:**
- Given a project with ambiguous file signals (both go.mod and package.json present)
- When the user runs `/forge:test-guide`
- Then the skill lists all detected language candidates and asks the user to select which one(s) to generate Conventions for

---

## Story 3: Cold Start Test Generation

**As a** Forge CLI user on a brand new project with no existing tests
**I want to** generate e2e tests that compile without first creating a Convention file
**So that** I can get started quickly and formalize conventions later

**Acceptance Criteria:**
- Given a new project with no Convention files and no existing test files
- When the user runs gen-test-scripts
- Then the skill proceeds with LLM defaults and Code Reconnaissance (file signal detection)
- And outputs a hint that no Convention was found
- And the generated test file passes `just e2e-compile`

**Population-level goal** (not individually verifiable per run): across forge-cli's 126+ existing test Journeys, >= 85% of generated files compile on first attempt without Convention files. This target is justified because forge-cli uses Go testing + testify — a mainstream combination with abundant LLM training data.

---

## Story 4: Backward Compatibility After Profile Removal

**As a** Forge CLI user on an existing project using default frameworks (Go + go-testing + testify)
**I want to** upgrade forge and have my tests continue to generate and run identically
**So that** my workflow is not disrupted by the Profile removal

**Acceptance Criteria:**
- Given an existing forge project with `languages`, `interfaces`, `test-framework` in `.forge/config.yaml`
- When the user upgrades forge to the Convention-based version
- Then those fields are silently ignored (no errors, no warnings, no migration prompts)
- And all existing e2e tests continue to pass via `just e2e-test`
- And `forge task index`, `forge config init`, and other CLI commands work without the removed fields

**Error scenarios:**
- Given the upgraded forge with no justfile `e2e-compile` recipe
- When the user runs gen-test-scripts
- Then gen-test-scripts outputs: "Missing justfile e2e-compile recipe. Run `/forge:init-justfile` first, or add a recipe manually."

---

## Story 5: Multi-Framework Convention Management

**As a** Forge CLI user on a mixed-language project (e.g., Go backend + TypeScript frontend)
**I want to** maintain separate Convention files for each framework
**So that** test generation uses the correct framework knowledge for each Journey's interface type

**Acceptance Criteria:**
- Given a project with `docs/conventions/testing-go.md` (domains: [testing, go]) and `docs/conventions/testing-javascript.md` (domains: [testing, javascript, web-ui])
- When gen-test-scripts generates tests for a Journey involving a CLI interface (Go)
- Then only the Go Convention file is loaded and applied
- And when generating for a Journey involving a web-ui interface (TypeScript), only the JavaScript Convention is loaded

**Error scenarios:**
- Given two Convention files with overlapping domains (e.g., both include `testing` in domains)
- When gen-test-scripts loads Conventions for a Journey
- Then both files are loaded and merged; last-loaded wins for conflicting fields
- And the skill logs a note about the domain overlap for user awareness

---

## Story 6: Compile Gate Recovery

**As a** Forge CLI user whose generated test fails the compile gate after all retries
**I want to** receive actionable guidance on how to fix the problem
**So that** I am not stuck at a dead-end with no path forward

**Acceptance Criteria:**
- Given gen-test-scripts has exhausted all compile gate retries (2 retries failed)
- When the compile gate reports failure
- Then the skill outputs: (1) the compile error message, (2) the generated file path for manual inspection, (3) suggested recovery actions: check Convention declarations, run `/forge:test-guide` to regenerate Convention, or manually edit the generated file
- And the generated file is preserved (not deleted) for user inspection
