---
name: convention-guide
description: Convention file structure, validation rules, merge semantics, and growth path for test framework knowledge files
---

# Convention File Structure Reference

Convention files define test generation rules for each language/framework combination. They are user-editable markdown files that provide framework knowledge to LLM-driven test generation.

## File Location and Discovery

Convention files live in `docs/conventions/testing/` with a surface-first directory structure:

```
docs/conventions/testing/
  index.md         — Speed-reference: Surface → Type → Location → Assertion Focus
  cli/
    index.md       — CLI documentation index
    core.md        — CLI test strategy (language-agnostic)
  api/
    index.md
    core.md        — API test strategy (language-agnostic)
  web/
    index.md
    core.md        — Web E2E test strategy (language-agnostic)
  tui/
    index.md
    core.md        — TUI test strategy (language-agnostic)
  mobile/
    index.md
    core.md        — Mobile test strategy (language-agnostic)
```

### Loading Mechanism (Surface-First)

1. Determine the active surface type (from `.forge/config.yaml` or auto-detection).
2. Load the surface Convention from `docs/conventions/testing/{surface}/core.md`.
3. If `core.md` does not exist for the surface, fall back to auto-detection from existing test files.
4. Each `core.md` contains an assertion preference table with per-framework rows — use this to resolve the target framework.

### Legacy Structure Detection

If `docs/conventions/testing/` contains flat `.md` files (e.g., `go.md`, `vitest.md`) instead of surface subdirectories, this indicates the legacy (framework-first) structure. Output migration prompt and do NOT load these files. Refer the user to run `/test-guide` to regenerate with the new structure.

<HARD-RULE>
Convention loading is surface-driven. Do NOT fall back to loading framework-specific flat files from the legacy structure.
</HARD-RULE>

## Section Schema (core.md)

Each `core.md` Convention file uses a fixed 7-section structure. These sections are required — without them, test generation falls back to LLM defaults with a warning.

### `file-location`

Declares where test files for this surface should be placed.

| Field | Type | Description |
|-------|------|-------------|
| test_dir | string | Root directory for test files |
| file_pattern | string | File naming pattern |
| build_tag | string | Per-surface build tag (e.g., `cli_functional`, `api_functional`) |

### `isolation-model`

Declares the isolation strategy for tests on this surface.

| Field | Type | Description |
|-------|------|-------------|
| model | string | Isolation approach (e.g., subprocess, HTTP server, browser context) |
| mechanism | string | How isolation is achieved (e.g., `t.TempDir()`, test containers) |

### `assertion-focus`

Declares what aspects tests on this surface should primarily assert.

| Field | Type | Description |
|-------|------|-------------|
| primary | string[] | Primary assertion targets (e.g., exit code, stdout, stderr) |
| secondary | string[] | Secondary assertion targets |

### `timeout-strategy`

Declares timeout recommendations for tests on this surface.

| Field | Type | Description |
|-------|------|-------------|
| default | string | Default timeout for individual test functions |
| smoke | string | Timeout for Journey smoke tests |

### `lifecycle`

Declares setup/teardown patterns for tests on this surface.

| Field | Type | Description |
|-------|------|-------------|
| setup | string | Setup pattern |
| teardown | string | Teardown pattern |
| hooks | string | Hook syntax for the framework |

### `contract-journey-ratio`

Declares the ratio target for Contract vs Journey smoke tests.

| Field | Type | Description |
|-------|------|-------------|
| ratio | string | Target ratio (e.g., ">= 80% Contract", "Balanced 50/50") |

### `anti-patterns`

Declares surface-specific forbidden patterns with replacements.

### Assertion Preference Table

Each `core.md` contains a per-framework assertion preference table with columns:

| Column | Description |
|--------|-------------|
| Assertion Library | Framework assertion library name |
| Mock Mechanism | Mock/stub approach for the framework |
| Fixture Pattern | Setup/teardown fixture pattern |

<HARD-RULE>
The assertion preference table columns are fixed to: Assertion Library, Mock Mechanism, Fixture Pattern. Adding new columns requires a proposal review to prevent core.md from regressing into per-surface framework files.
</HARD-RULE>

## Validation Rules

When skills load Convention files, these validation rules apply:

### Missing required sections

| Condition | Behavior |
|-----------|----------|
| All 7 required sections missing | Proceed with LLM defaults, output warning listing missing sections |
| Some required sections missing | Proceed with LLM defaults for missing sections only, log warning |
| Section heading present but content empty | Treat as missing section |
| Required field within section missing | Treat field as absent, use LLM default for that field |

### File access errors

| Condition | Behavior |
|-----------|----------|
| `core.md` not found for surface | Fall back to auto-detection from existing test files |
| No Convention files in `docs/conventions/testing/{surface}/` | Proceed with LLM defaults + Reconnaissance |
| File unreadable (permissions, encoding) | Skip file, output warning |

### Convention vs Reconnaissance conflict

| Condition | Behavior |
|-----------|----------|
| Convention declares X but Reconnaissance detects Y | Convention wins (user-edited knowledge overrides auto-detection) |

## Growth Path

Convention files support progressive enrichment. Start minimal, expand as needed.

### Level 1: Minimal (auto-generated by `/forge:test-guide`)

The 7 required sections only (file-location, isolation-model, assertion-focus, timeout-strategy, lifecycle, contract-journey-ratio, anti-patterns). Sufficient for new projects.

### Level 2: Extended (user enriches assertion preference table)

Add framework-specific rows to the assertion preference table after generating tests and observing patterns.

### How to add sections

1. **Direct edit**: Open `docs/conventions/testing/{surface}/core.md` and update the assertion preference table or section content.
2. **Re-run `/forge:test-guide`**: After creating test files, re-running the skill detects patterns and proposes updated Convention content (presents diff for confirmation).

## Relationship to Forge Components

| Component | Interaction |
|-----------|-------------|
| `/forge:test-guide` | Generates Level 1 Convention files (per-surface `core.md` with 7 required sections) |
| `/forge:gen-test-scripts` | Reads per-surface `core.md` for test generation strategy |
| `/forge:run-tests` | Reads `core.md` timeout and lifecycle sections |
| `/forge:init-justfile` | Reads `core.md` file-location section for test recipe generation |
| `/forge:consolidate-specs` | Treats Convention files as standard conventions for drift detection |
