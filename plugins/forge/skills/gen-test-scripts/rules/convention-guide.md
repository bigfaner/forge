---
name: convention-guide
description: Convention file structure, validation rules, merge semantics, and growth path for test framework knowledge files
---

# Convention File Structure Reference

Convention files define test generation rules for each language/framework combination. They are user-editable markdown files that provide framework knowledge to LLM-driven test generation.

## File Location and Discovery

Convention files live in `docs/conventions/testing/` with a two-level index:

```
docs/conventions/testing/
  index.md         — Lists all available Conventions with name, description, applicability
  go.md            — Go testing + testify Convention
  ginkgo.md        — Ginkgo v2 Convention
  vitest.md        — Vitest Convention
  pytest.md        — pytest Convention
  junit.md         — JUnit 5 Convention
  rust.md          — Rust Convention
```

### Loading Mechanism

1. Read `docs/conventions/testing/index.md` to discover available Conventions
2. Based on the project's language/framework context, select the matching Convention
3. Load the selected file from `docs/conventions/testing/<convention>.md`
4. If `index.md` does not exist, fall back to auto-detection from existing test files

<HARD-RULE>
Do NOT use `domains` frontmatter filtering. Selection is based on index.md descriptions and project context, with LLM autonomous judgment.
</HARD-RULE>

## Section Schema

Convention files use a fixed 4-section structure. These four sections are required — without them, test generation falls back to LLM defaults with a warning.

### `framework`

Declares the test framework identity and file conventions.

| Field | Type | Description |
|-------|------|-------------|
| name | string | Framework name and assertion library |
| file-pattern | string | Glob pattern for test files |
| test-dir | string | Default test directory |
| runner-command | string | Command to run tests |
| build-tag | string | Tag or marker syntax for test categorization |

### `discovery`

Declares how tests are discovered and organized.

| Field | Type | Description |
|-------|------|-------------|
| test_dir | string | Root directory for test files |
| file_pattern | string | File naming pattern |
| exclude_pattern | string | Files/directories to exclude |

### `structure`

Declares test file structure patterns and result format.

| Field | Type | Description |
|-------|------|-------------|
| suite_pattern | string | How test suites/groups are declared |
| case_pattern | string | How individual test cases are declared |
| hook_pattern | string | Setup/teardown hook syntax |
| output-flags | string | Command-line flags for machine-readable output |
| format-type | enum | One of: `json-stream`, `json-report`, `text-verbose` |

**Format type reference:**

| Format type | Description | Typical frameworks |
|-------------|-------------|-------------------|
| `json-stream` | Line-delimited JSON objects (one per test event) | Go testing (`go test -json`) |
| `json-report` | Single JSON object with nested test results | Vitest (`--reporter=json`), Jest (`--json`) |
| `text-verbose` | Human-readable text output | Cargo test, generic CLI tools |

### `assertions`

Declares the assertion library and usage rules.

| Field | Type | Description |
|-------|------|-------------|
| style | enum | One of: `assert`, `expect`, `should` |
| library | string | Assertion library name and import path |
| custom_matchers | string[] | Project-specific matchers (optional) |

### Optional Sections

Users can add these sections to improve generation quality. See [Growth Path](#growth-path) for when to add each.

| Section | Purpose |
|---------|---------|
| Import Patterns | Standard import blocks for e2e tests |
| Code Style | Test naming, table-driven patterns, traceability conventions |
| Anti-patterns | Framework-specific forbidden patterns with replacements |
| Helpers | Common helper functions in the project's test infrastructure |

## Validation Rules

When skills load Convention files, these validation rules apply:

### Missing required sections

| Condition | Behavior |
|-----------|----------|
| All 4 required sections missing | Proceed with LLM defaults, output warning listing missing sections |
| Some required sections missing | Proceed with LLM defaults for missing sections only, log warning |
| Section heading present but content empty | Treat as missing section |
| Required field within section missing | Treat field as absent, use LLM default for that field |

### File access errors

| Condition | Behavior |
|-----------|----------|
| `index.md` not found | Fall back to auto-detection from existing test files |
| No Convention files in `docs/conventions/testing/` | Proceed with LLM defaults + Reconnaissance |
| File unreadable (permissions, encoding) | Skip file, output warning |

### Convention vs Reconnaissance conflict

| Condition | Behavior |
|-----------|----------|
| Convention declares X but Reconnaissance detects Y | Convention wins (user-edited knowledge overrides auto-detection) |

## Merge Semantics

When multiple Convention files are loaded, merge at the **section level**:

1. **Later-loaded wins for conflicting sections**: If two files both declare a `framework` section, the later file's section overwrites the earlier one entirely.
2. **Unique sections are preserved**: If file A has `Helpers` and file B does not, file A's `Helpers` section is kept.
3. **Within a section, later-loaded wins at the field level**.
4. **Log overlap notes**: When overlap is detected, output a note listing which sections were overwritten.

## Growth Path

Convention files support progressive enrichment. Start minimal, expand as needed.

### Level 1: Minimal (auto-generated by `/forge:test-guide`)

The 4 required sections only (`framework`, `discovery`, `structure`, `assertions`). Sufficient for new projects.

### Level 2: Extended (user adds optional sections)

Add Import Patterns and Code Style after generating tests and observing patterns.

### Level 3: Comprehensive (full knowledge capture)

Add Anti-patterns and Helpers. Sufficient for mature projects with custom test infrastructure.

### When to expand

| Trigger | Section to add |
|---------|---------------|
| LLM generates wrong imports consistently | Import Patterns |
| Test naming differs from project convention | Code Style |
| LLM generates patterns that compile but violate team rules | Anti-patterns |
| LLM does not use existing helper functions | Helpers |

### How to add sections

1. **Direct edit**: Open `docs/conventions/testing/<convention>.md` and add the section heading with content.
2. **Re-run `/forge:test-guide`**: After creating test files, re-running the skill detects patterns and proposes updated Convention content (presents diff for confirmation).

## Relationship to Forge Components

| Component | Interaction |
|-----------|-------------|
| `/forge:test-guide` | Generates Level 1 Convention files (required sections only) |
| `/forge:gen-test-scripts` | Reads Convention files to generate test code |
| `/forge:run-tests` | Reads `structure` section for result format parsing |
| `/forge:init-justfile` | Reads `framework` section to generate justfile recipes |
| `/forge:consolidate-specs` | Treats Convention files as standard conventions for drift detection |
