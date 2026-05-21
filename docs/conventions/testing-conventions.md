---
title: Testing Convention Files
domains:
  - testing
  - e2e
  - code-generation
---

# Testing Convention Files

Convention files (`docs/conventions/testing-<scope>.md`) define the test generation rules for each language/framework combination. They are user-editable markdown files that provide framework knowledge to LLM-driven test generation, replacing the previous hardcoded Profile system.

This document describes the Convention file structure, section schema, validation rules, and usage patterns.

## File Location and Naming

Convention files live in the user project's `docs/conventions/` directory:

```
docs/conventions/
  testing-go.md          # domains: [testing, go]
  testing-javascript.md  # domains: [testing, javascript]
  testing-python.md      # domains: [testing, python]
  testing-ginkgo.md      # domains: [testing, go, ginkgo]
  testing-vitest.md      # domains: [testing, typescript, javascript, vitest]
```

Naming convention: `testing-<scope>.md` where `<scope>` is a language or framework identifier.

## Frontmatter Schema

Every Convention file must include YAML frontmatter:

```yaml
---
title: "<Framework Name> Testing Convention"
domains: [testing, <scope>]
---
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `title` | string | yes | Human-readable title for the Convention |
| `domains` | string[] | yes | Tags for selective loading. Must include `testing` plus at least one scope identifier |

Common `domains` values:

| Framework | domains |
|-----------|---------|
| Go testing + testify | `[testing, go]` |
| Ginkgo v2 | `[testing, go, ginkgo]` |
| Vitest | `[testing, typescript, javascript, vitest]` |
| pytest | `[testing, python]` |

## Section Schema

Convention files use a fixed section structure with markdown headings (`## Section Name`).

### Required Sections (minimum set)

These four sections must be present. Without them, test generation falls back to LLM defaults with a warning.

#### Framework

Declares the test framework identity and file conventions.

| Field | Type | Description |
|-------|------|-------------|
| Name | string | Framework name and assertion library |
| File pattern | string | Glob pattern for test files |
| Package | string | Test package name (for languages with package scoping) |
| Test runner | string | Command to run tests |
| Build tag | string | Tag or marker syntax for test categorization |

#### Assertion

Declares the assertion library and specific functions to use.

| Field | Type | Description |
|-------|------|-------------|
| Library | string | Assertion library name and import path |
| Key functions | string[] | Primary assertion functions with purpose |
| Rule | string | Usage rule or constraint |

#### Tags

Declares the build tag or marker syntax for test categorization. Must include a code block showing the tag in context.

#### Result Format

Declares how test results are produced and parsed.

| Field | Type | Description |
|-------|------|-------------|
| Output flags | string | Command-line flags for machine-readable output |
| Format type | enum | One of: `json-stream`, `json-report`, `text-verbose` |

**Format type reference:**

| Format type | Description | Typical frameworks |
|-------------|-------------|-------------------|
| `json-stream` | Line-delimited JSON objects (one per test event) | Go testing (`go test -json`) |
| `json-report` | Single JSON object with nested test results | Vitest (`--reporter=json`), Jest (`--json`) |
| `text-verbose` | Human-readable text output | Cargo test, generic CLI tools |

### Optional Sections

Users can add these sections to improve generation quality. See [Growth Path](#growth-path) for when to add each.

| Section | Purpose |
|---------|---------|
| Import Patterns | Standard import blocks for e2e tests |
| Code Style | Test naming, table-driven patterns, traceability conventions |
| Anti-patterns | Framework-specific forbidden patterns with replacements |
| Helpers | Common helper functions in the project's test infrastructure |

## Complete Examples

### Go testing + testify

```markdown
---
title: "Go Testing Convention"
domains: [testing, go]
---

# Go Testing Convention

Convention for generating Go test code using Go testing package with testify/assert.

## Framework

- **Name**: Go testing package + testify/assert
- **File pattern**: *_test.go
- **Package**: e2e
- **Test runner**: go test
- **Build tag**: //go:build e2e

## Assertion

- **Library**: assert from github.com/stretchr/testify/assert
- **Key functions**:
  - assert.NoError(err) -- verify no error returned
  - assert.Contains(haystack, needle) -- verify string/array contains element
  - assert.Equal(expected, actual) -- verify exact equality
  - assert.True(condition) -- verify boolean is true
  - assert.False(condition) -- verify boolean is false
- **Rule**: Always use assert (NOT require). assert allows subsequent checks within the same test.

## Tags

Build tag placed at the top of every test file, before the package declaration.

\`\`\`go
//go:build e2e

package e2e
\`\`\`

## Result Format

- **Output flags**: -json
- **Format type**: json-stream

## Import Patterns

Standard imports for e2e test files:

\`\`\`go
import (
    "os/exec"
    "strings"
    "testing"

    "github.com/stretchr/testify/assert"
)
\`\`\`

## Code Style

- **Test naming**: TestTC_NNN_Description (e.g., TestTC_001_CreateProject)
- **Traceability**: Each test function includes a comment referencing the TC ID
- **Subtests**: Use t.Run for parameterized or grouped test cases

## Anti-patterns

- **Forbidden**: time.Sleep for synchronization -- use assert.Eventually or retry loops
- **Forbidden**: Hardcoded ports -- use dynamically assigned ports
- **Forbidden**: require.* from testify -- use assert.* to allow multiple assertions per test

## Helpers

- runCLI(args ...string) *exec.Cmd -- Execute forge CLI with arguments
- withRetry(fn func() error, maxAttempts int) error -- Retry function with backoff
```

### Python pytest

```markdown
---
title: "Python pytest Convention"
domains: [testing, python]
---

# Python pytest Convention

Convention for generating Python test code using pytest.

## Framework

- **Name**: pytest
- **File pattern**: test_*.py
- **Package**: tests/e2e
- **Test runner**: pytest
- **Build tag**: @pytest.mark.e2e (decorator)

## Assertion

- **Library**: Built-in assert statement (Python native)
- **Key functions**:
  - assert value -- native Python assertion
  - assert a == b -- equality check
  - assert a in b -- containment check
  - assert isinstance(obj, cls) -- type check
- **Rule**: Use plain assert statements. pytest provides detailed introspection on failure.

## Tags

Markers applied as decorators to test functions or classes.

\`\`\`python
import pytest

@pytest.mark.e2e
def test_tc_001_feature():
    ...
\`\`\`

## Result Format

- **Output flags**: --json-report --json-report-file=results.json
- **Format type**: json-report

## Import Patterns

Standard imports for e2e test files:

\`\`\`python
import pytest
import subprocess
\`\`\`

## Code Style

- **Test naming**: test_tc_nnn_description (e.g., test_tc_001_create_project)
- **Traceability**: Each test function name includes the TC ID
- **Fixtures**: Use pytest fixtures for setup/teardown

## Anti-patterns

- **Forbidden**: time.sleep() for synchronization -- use pytest-repeat or polling
- **Forbidden**: Bare print() for debugging -- use capsys fixture
- **Forbidden**: unittest.TestCase style -- use plain functions and fixtures
```

### JavaScript Vitest

```markdown
---
title: "Vitest Testing Convention"
domains: [testing, typescript, javascript, vitest]
---

# Vitest Testing Convention

Convention for generating TypeScript test code using Vitest.

## Framework

- **Name**: Vitest
- **File pattern**: *.test.ts
- **Package**: tests/e2e
- **Test runner**: vitest run
- **Build tag**: { tags: ['@e2e'] } in describe/it options

## Assertion

- **Library**: Built-in expect from vitest (Jest-compatible)
- **Key functions**:
  - expect(value).toBe(expected) -- strict equality
  - expect(value).toContain(item) -- containment check
  - expect(value).toEqual(expected) -- deep equality
  - expect(fn).toThrow() -- error check
  - expect(value).toBeTruthy() -- truthy check
- **Rule**: Use Vitest built-in expect. Do not import chai or jest-matchers separately.

## Tags

Tags applied via describe/it options object.

\`\`\`typescript
import { describe, it, expect } from 'vitest'

describe('Feature', { tags: ['@e2e'] }, () => {
  it('should work', () => {
    // ...
  })
})
\`\`\`

## Result Format

- **Output flags**: --reporter=json
- **Format type**: json-report

## Import Patterns

Standard imports for e2e test files:

\`\`\`typescript
import { describe, it, expect, beforeAll, afterAll } from 'vitest'
import { execSync } from 'child_process'
\`\`\`

## Code Style

- **Test naming**: describe('Feature: ...', () => { it('should ...', ...) })
- **Traceability**: Describe block includes TC ID reference
- **Async**: Use async/await pattern for all async operations

## Anti-patterns

- **Forbidden**: await new Promise(r => setTimeout(r, N)) -- use vi.waitFor or polling
- **Forbidden**: Synchronous exec in async tests -- use execAsync or execSync consistently
- **Forbidden**: done() callback style -- use async/await
```

## Validation Rules

When skills load Convention files, these validation rules apply:

### Missing frontmatter

| Condition | Behavior |
|-----------|----------|
| File has no YAML frontmatter | Skip file, output warning |
| `domains` field missing | Skip file, output warning |
| `domains` field empty | Skip file, output warning |

### Missing required sections

| Condition | Behavior |
|-----------|----------|
| All 4 required sections missing | Proceed with LLM defaults, output warning |
| Some required sections missing | Proceed with LLM defaults for missing sections only |

### Invalid section content

| Condition | Behavior |
|-----------|----------|
| Section heading present but content empty | Treat as missing section |
| Required field within section missing (e.g., Framework has no Name) | Treat field as absent, use LLM default for that field |

### File access errors

| Condition | Behavior |
|-----------|----------|
| No Convention files found | Proceed with LLM defaults + Reconnaissance |
| File unreadable (permissions, encoding) | Skip file, output warning |

### Convention vs Reconnaissance conflict

| Condition | Behavior |
|-----------|----------|
| Convention declares X but Reconnaissance detects Y | Convention wins (user-edited knowledge overrides auto-detection) |

## Merge Semantics

When multiple Convention files have overlapping `domains`, the LLM merges them at the **section level**.

### Rules

1. **Last-loaded wins for conflicting sections**: If two files both declare a `Framework` section, the later file's `Framework` section overwrites the earlier one entirely.
2. **Unique sections are preserved**: If file A has `Helpers` and file B does not, file A's `Helpers` section is kept.
3. **Within a section, last-loaded wins at the field level**: If both files declare `Framework.Name`, the later file's value is used.
4. **The skill logs overlap notes**: When domain overlap is detected, the skill outputs a note listing which sections were overwritten.

### Example

Given two files loaded in order:

**File 1: `testing-go.md`** (domains: [testing, go]) -- has Framework, Assertion, Tags, Helpers sections.

**File 2: `testing-ginkgo.md`** (domains: [testing, go, ginkgo]) -- has Framework, Assertion, Tags sections.

**Merged result** (for a Journey matching [testing, go, ginkgo]):
- Framework, Assertion, Tags: from File 2 (overwritten)
- Helpers: from File 1 (preserved, unique section)

Skill output: "Convention files testing-go.md and testing-ginkgo.md have overlapping domains [testing, go]. Sections Framework, Assertion, Tags from testing-ginkgo.md overwrite those in testing-go.md."

## Growth Path

Convention files support progressive enrichment. Start minimal, expand as needed.

### Level 1: Minimal (auto-generated by `/forge:test-guide`)

The 4 required sections only. Sufficient for new projects and standard framework configurations.

```
Framework + Assertion + Tags + Result Format
```

### Level 2: Extended (user adds optional sections)

Add Import Patterns and Code Style after generating tests and observing patterns. Sufficient for projects with established test conventions.

```
+ Import Patterns
+ Code Style
```

### Level 3: Comprehensive (full knowledge capture)

Add Anti-patterns and Helpers for the complete testing knowledge. Sufficient for mature projects with custom test infrastructure and strict testing standards.

```
+ Anti-patterns
+ Helpers
```

### When to expand

| Trigger | Section to add |
|---------|---------------|
| LLM generates wrong imports consistently | Import Patterns |
| Test naming differs from project convention | Code Style |
| LLM generates patterns that compile but violate team rules | Anti-patterns |
| LLM does not use existing helper functions | Helpers |

### How to add sections

1. **Direct edit**: Open `docs/conventions/testing-<scope>.md` and add the section heading with content.
2. **Re-run `/forge:test-guide`**: After creating test files, re-running the skill detects patterns and proposes updated Convention content (presents diff for confirmation).

## Relationship to Forge Components

| Component | Interaction |
|-----------|-------------|
| `/forge:test-guide` | Generates Level 1 Convention files (required sections only) |
| `/forge:gen-test-scripts` | Reads Convention files to generate test code |
| `/forge:run-e2e-tests` | Reads Result Format section to parse test output |
| `/forge:init-justfile` | Reads Framework section to generate justfile recipes |
| `/forge:consolidate-specs` | Treats Convention files as standard conventions for drift detection |

## Reference

For the detailed technical specification behind this structure, see [design/convention-file-structure.md](../features/test-knowledge-convention-driven/design/convention-file-structure.md).
