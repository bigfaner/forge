---
created: 2026-05-20
design: tech-design.md
prd: prd/prd-spec.md
status: Reference
---

# Convention File Structure Reference

This document defines the fixed structure, section schema, and validation rules for test Convention files (`docs/conventions/testing-<scope>.md`). Convention files are the user-editable knowledge layer that drives test code generation in forge.

## File Location and Naming

Convention files live in the user project's `docs/conventions/` directory:

```
docs/conventions/
  testing-go.md          # domains: [testing, go]
  testing-javascript.md  # domains: [testing, javascript]
  testing-python.md      # domains: [testing, python]
```

Naming convention: `testing-<scope>.md` where `<scope>` is a language or framework identifier (e.g., `go`, `javascript`, `python`, `ginkgo`, `vitest`).

## Frontmatter Schema

Every Convention file must include YAML frontmatter with these fields:

```yaml
---
title: "<Framework Name> Testing Convention"
domains: [testing, <scope>]
---
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `title` | string | yes | Human-readable title for the Convention |
| `domains` | string[] | yes | Tags used by skills to load relevant Conventions. Must include `testing` plus at least one scope identifier |

**Example frontmatter values:**

| Framework | domains |
|-----------|---------|
| Go testing + testify | `[testing, go]` |
| Ginkgo v2 | `[testing, go, ginkgo]` |
| Vitest | `[testing, typescript, javascript, vitest]` |
| pytest | `[testing, python]` |

## Section Schema

Convention files use a fixed section structure. Sections are markdown headings (`## Section Name`) with bullet-list content.

### Required Sections (always present)

These four sections form the minimum set. Without them, test generation falls back to LLM defaults with a warning.

#### Framework

Declares the test framework identity and file conventions.

| Field | Type | Description |
|-------|------|-------------|
| Name | string | Framework name and assertion library |
| File pattern | string | Glob pattern for test files |
| Package | string | Test package name (for languages with package scoping) |
| Test runner | string | Command to run tests |
| Build tag | string | Tag or marker syntax for test categorization (if applicable) |

**Example (Go testing):**

```markdown
## Framework

- **Name**: Go testing package + testify/assert
- **File pattern**: *_test.go
- **Package**: e2e
- **Test runner**: go test
- **Build tag**: //go:build e2e
```

#### Assertion

Declares the assertion library and the specific functions to use.

| Field | Type | Description |
|-------|------|-------------|
| Library | string | Assertion library name and import path |
| Key functions | string[] | Primary assertion functions with purpose |
| Rule | string | Usage rule or constraint |

**Example (Go testify):**

```markdown
## Assertion

- **Library**: assert from github.com/stretchr/testify/assert
- **Key functions**:
  - assert.NoError -- verify no error returned
  - assert.Contains -- verify string/array contains element
  - assert.Equal -- verify exact equality
  - assert.True -- verify boolean is true
- **Rule**: Always use assert, never require. assert allows subsequent checks within the same test to execute.
```

#### Tags

Declares the build tag or marker syntax used to categorize test files.

| Field | Type | Description |
|-------|------|-------------|
| Build tag | string | Tag syntax with placement rules |

The section must include a code block showing the tag in context.

**Example (Go):**

```markdown
## Tags

Build tag placed at the top of every test file, before the package declaration.

```go
//go:build e2e

package e2e
```
```

#### Result Format

Declares how test results are produced and parsed.

| Field | Type | Description |
|-------|------|-------------|
| Output flags | string | Command-line flags that produce machine-readable output |
| Format type | enum | One of: `json-stream`, `json-report`, `text-verbose` |
| Execution command | string | Full command to execute tests |

**Format type reference:**

| Format type | Description | Typical frameworks |
|-------------|-------------|-------------------|
| `json-stream` | Line-delimited JSON objects (one per test event) | Go testing (`go test -json`) |
| `json-report` | Single JSON object with nested test results | Vitest (`--reporter=json`), Jest (`--json`) |
| `text-verbose` | Human-readable text output | Cargo test, generic CLI tools |

**Example (Go):**

```markdown
## Result Format

- **Output flags**: -json
- **Format type**: json-stream
- **Execution command**: go test -json -tags=e2e ./tests/e2e/...
```

### Optional Sections

These sections are included when the extracted patterns or framework knowledge provide enough detail. Users can add them by editing the Convention file directly or re-running `/forge:test-guide`.

#### Import Patterns

Standard import blocks for e2e tests. Helps the LLM generate correct import declarations.

**Example (Go):**

```markdown
## Import Patterns

Standard imports for e2e test files:

```go
import (
    "os/exec"
    "strings"
    "testing"

    "github.com/stretchr/testify/assert"
)
```
```

#### Code Style

Test function naming conventions, table-driven test patterns, traceability conventions.

**Example (Go):**

```markdown
## Code Style

- **Test naming**: TestTC_NNN_Description (e.g., TestTC_001_CreateProject)
- **Traceability**: Each test function includes a comment with the TC ID
- **Table-driven**: Use subtests (t.Run) for parameterized cases
```

#### Anti-patterns

Framework-specific forbidden patterns with suggested replacements.

**Example (Go):**

```markdown
## Anti-patterns

- **Forbidden**: `time.Sleep` for waiting -- use assert.Eventually or retry loops
- **Forbidden**: Hardcoded ports -- use dynamic port assignment
- **Forbidden**: `require.*` assertions -- use `assert.*` instead
```

#### Helpers

Common helper functions available in the project's test infrastructure.

**Example (Go):**

```markdown
## Helpers

- `runCLI(args ...string) *exec.Cmd` -- Execute forge CLI command
- `withRetry(fn func() error, maxAttempts int) error` -- Retry with backoff
```

## Complete Examples

### Go testing + testify

```markdown
---
title: "Go Testing Convention"
domains: [testing, go]
---

<!-- auto-generated by forge:test-guide -->

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

```go
//go:build e2e

package e2e
```

## Result Format

- **Output flags**: -json
- **Format type**: json-stream
- **Execution command**: go test -json -tags=e2e ./tests/e2e/...

## Import Patterns

Standard imports for e2e test files:

```go
import (
    "os/exec"
    "strings"
    "testing"

    "github.com/stretchr/testify/assert"
)
```

## Code Style

- **Test naming**: TestTC_NNN_Description (e.g., TestTC_001_CreateProject)
- **Traceability**: Each test function includes a comment referencing the TC ID
- **Subtests**: Use t.Run for parameterized or grouped test cases

## Anti-patterns

- **Forbidden**: `time.Sleep` for synchronization -- use assert.Eventually or retry loops
- **Forbidden**: Hardcoded ports -- use dynamically assigned ports
- **Forbidden**: `require.*` from testify -- use `assert.*` to allow multiple assertions per test

## Helpers

- `runCLI(args ...string) *exec.Cmd` -- Execute forge CLI with arguments
- `withRetry(fn func() error, maxAttempts int) error` -- Retry function with backoff
```

### Python pytest

```markdown
---
title: "Python pytest Convention"
domains: [testing, python]
---

<!-- auto-generated by forge:test-guide -->

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

```python
import pytest

@pytest.mark.e2e
def test_tc_001_feature():
    ...
```

## Result Format

- **Output flags**: --json-report --json-report-file=results.json
- **Format type**: json-report
- **Execution command**: pytest --json-report --json-report-file=results.json tests/e2e/

## Import Patterns

Standard imports for e2e test files:

```python
import pytest
import subprocess
```

## Code Style

- **Test naming**: test_tc_nnn_description (e.g., test_tc_001_create_project)
- **Traceability**: Each test function name includes the TC ID
- **Fixtures**: Use pytest fixtures for setup/teardown

## Anti-patterns

- **Forbidden**: `time.sleep()` for synchronization -- use pytest-repeat or polling
- **Forbidden**: Bare `print()` for debugging -- use `capsys` fixture
- **Forbidden**: `unittest.TestCase` style -- use plain functions and fixtures
```

### JavaScript Vitest

```markdown
---
title: "Vitest Testing Convention"
domains: [testing, typescript, javascript, vitest]
---

<!-- auto-generated by forge:test-guide -->

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

```typescript
import { describe, it, expect } from 'vitest'

describe('Feature', { tags: ['@e2e'] }, () => {
  it('should work', () => {
    // ...
  })
})
```

## Result Format

- **Output flags**: --reporter=json
- **Format type**: json-report
- **Execution command**: npx vitest run --reporter=json tests/e2e/

## Import Patterns

Standard imports for e2e test files:

```typescript
import { describe, it, expect, beforeAll, afterAll } from 'vitest'
import { execSync } from 'child_process'
```

## Code Style

- **Test naming**: describe('Feature: ...', () => { it('should ...', ...) })
- **Traceability**: Describe block includes TC ID reference
- **Async**: Use async/await pattern for all async operations

## Anti-patterns

- **Forbidden**: `await new Promise(r => setTimeout(r, N))` -- use vi.waitFor or polling
- **Forbidden**: Synchronous exec in async tests -- use execAsync or execSync consistently
- **Forbidden**: `done()` callback style -- use async/await
```

## Validation Rules

When skills load Convention files, the following validation rules apply:

### Missing frontmatter

| Condition | Behavior | User-Facing Output |
|-----------|----------|--------------------|
| File has no YAML frontmatter | Skip file | "Convention file `<path>` has no frontmatter. Skipping." |
| `domains` field missing | Skip file | "Convention file `<path>` has no domains frontmatter. Skipping." |
| `domains` field empty | Skip file | "Convention file `<path>` has empty domains. Skipping." |

### Missing required sections

| Condition | Behavior | User-Facing Output |
|-----------|----------|--------------------|
| All 4 required sections missing | Proceed with LLM defaults | "Convention file `<path>` is missing all required sections. Using LLM defaults." |
| Some required sections missing | Proceed with LLM defaults for missing sections | "Convention file `<path>` is missing sections: Assertion, Tags. Using LLM defaults for those sections." |

### Invalid section content

| Condition | Behavior | User-Facing Output |
|-----------|----------|--------------------|
| Section heading present but content empty | Treat as missing section | "Convention file `<path>` has empty Framework section. Using LLM defaults." |
| Required field within section missing (e.g., Framework has no Name) | Treat field as absent, use LLM default for that field | "Convention file `<path>` Framework section is missing Name field." |

### File access errors

| Condition | Behavior | User-Facing Output |
|-----------|----------|--------------------|
| No Convention files found | Proceed with LLM defaults + Reconnaissance | "No test Convention files found in docs/conventions/. Generation will use LLM defaults. Run /forge:test-guide to create one." |
| File unreadable (permissions, encoding) | Skip file | "Cannot read Convention file `<path>`: `<error>`. Skipping." |

### Convention vs Reconnaissance conflict

| Condition | Behavior | User-Facing Output |
|-----------|----------|--------------------|
| Convention declares X but Reconnaissance detects Y | Convention wins | "Convention declares `<X>` but existing tests use `<Y>`. Using Convention value." |

## Merge Semantics

When multiple Convention files have overlapping `domains`, the LLM merges them at the **section level**.

### Section-level merge rules

1. **Last-loaded wins for conflicting sections**: If two files both declare a `Framework` section, the later file's `Framework` section overwrites the earlier one entirely.
2. **Unique sections are preserved**: If file A has `Helpers` and file B does not, file A's `Helpers` section is kept.
3. **Within a section, last-loaded wins at the field level**: If both files declare `Framework.Name`, the later file's value is used.
4. **The skill logs overlap notes**: When domain overlap is detected, the skill outputs a note listing which sections were overwritten.

### Merge example

Given two Convention files loaded in order:

**File 1: `testing-go.md` (domains: [testing, go])**

```yaml
Framework:
  Name: Go testing package + testify/assert
  File pattern: *_test.go

Assertion:
  Library: assert from testify

Tags:
  Build tag: //go:build e2e

Helpers:
  - runCLI()
```

**File 2: `testing-ginkgo.md` (domains: [testing, go, ginkgo])**

```yaml
Framework:
  Name: Ginkgo v2 + Gomega
  File pattern: *_test.go

Assertion:
  Library: Expect from Gomega

Tags:
  Build tag: //go:build e2e
```

**Merged result (for a Journey matching `[testing, go, ginkgo]`):**

```yaml
Framework:
  Name: Ginkgo v2 + Gomega        # Overwritten by File 2
  File pattern: *_test.go          # Overwritten by File 2

Assertion:
  Library: Expect from Gomega      # Overwritten by File 2

Tags:
  Build tag: //go:build e2e        # Overwritten by File 2

Helpers:
  - runCLI()                       # Preserved from File 1 (unique section)
```

**Skill output**: "Convention files testing-go.md and testing-ginkgo.md have overlapping domains [testing, go]. Sections Framework, Assertion, Tags from testing-ginkgo.md overwrite those in testing-go.md."

## Growth Path

Convention files support a progressive enrichment model. Start with the minimum, expand as needed.

### Level 1: Minimal (auto-generated by `/forge:test-guide`)

The minimal Convention file includes only the 4 required sections. This is what `/forge:test-guide` produces.

```
Framework + Assertion + Tags + Result Format
```

Sufficient for: new projects, first-time Convention users, standard framework configurations.

### Level 2: Extended (user adds optional sections)

After generating tests and observing patterns, users add optional sections to improve generation quality.

```
+ Import Patterns
+ Code Style
```

Sufficient for: projects with established test conventions, teams that want consistent test style across generated files.

### Level 3: Comprehensive (full knowledge capture)

Add Anti-patterns and Helpers to capture the complete testing knowledge for the project.

```
+ Anti-patterns
+ Helpers
```

Sufficient for: mature projects with custom test infrastructure, teams enforcing strict testing standards.

### When to expand

| Trigger | Section to add |
|---------|---------------|
| LLM generates wrong imports consistently | Import Patterns |
| Test naming differs from project convention | Code Style |
| LLM generates patterns that compile but violate team rules | Anti-patterns |
| LLM does not use existing helper functions | Helpers |

### Adding sections

Users can add optional sections by:

1. **Direct edit**: Open `docs/conventions/testing-<scope>.md` and add the section heading with content.
2. **Re-run `/forge:test-guide`**: After creating test files, re-running the skill detects patterns and proposes updated Convention content (presents diff for confirmation).
3. **Manual creation**: Write the optional section from scratch following the schema in this document.

## Relationship to Other Forge Components

| Component | Interaction |
|-----------|-------------|
| `/forge:test-guide` | Generates Level 1 Convention files (required sections only) |
| `/forge:gen-test-scripts` | Reads Convention files to generate test code |
| `/forge:run-e2e-tests` | Reads Result Format section to parse test output |
| `/forge:init-justfile` | Reads Framework section to generate justfile recipes |
| `/forge:consolidate-specs` | Treats Convention files as standard conventions in `docs/conventions/` for drift detection |
