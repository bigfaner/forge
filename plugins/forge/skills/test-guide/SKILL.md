---
name: test-guide
description: Guide users through generating Convention files. Detects project file signals, scans test files, extracts patterns, and writes a minimal Convention file after user confirmation.
allowed-tools: Bash Read Write Edit Grep Glob
disable-model-invocation: true
argument-hint: '[--scope <scope>] [--force]'
---

# /forge:test-guide

MANUAL-ONLY. Do NOT auto-invoke -- only when user explicitly asks to `/forge:test-guide`.

Guide users through creating test Convention files (`docs/conventions/testing-<scope>.md`). The skill detects project file signals, scans existing test files, extracts patterns, and writes a minimal Convention file after user confirmation.

## Parameters

| Parameter  | Values                    | Default         | Description                                        |
| ---------- | ------------------------- | --------------- | -------------------------------------------------- |
| `--scope`  | language or framework name | (auto-detect)  | Override the Convention file scope (e.g., `go`, `javascript`, `python`) |
| `--force`  | (flag)                    | false           | Overwrite existing Convention file without confirmation |

## Workflow

```
0. Check existing Convention files -> 1. Scan file signals -> 2. Scan test files & extract patterns -> 3. Present findings & confirm -> 4. Write Convention file
```

### Step 0: Check Existing Convention Files

Check whether a Convention file already exists for the target scope.

1. Glob `docs/conventions/testing-*.md` in the project root.
2. Read each file's YAML frontmatter `domains` field.
3. If a Convention file with matching scope already exists:
   - Read its full content.
   - Record it as `existing_convention` for diff comparison in Step 3.
   - If `--force` is set: proceed to Step 1 (will overwrite).
   - If `--force` is NOT set: note the existing file and proceed -- Step 3 will present a diff.

**If no existing Convention files found**: proceed to Step 1 (fresh generation).

### Step 1: Scan File Signals

Detect the project's language and framework from file system signals. This step does NOT execute any code -- it reads file names and parses file contents.

#### 1a. Detect language signals

Scan the project root for marker files:

```bash
ls go.mod package.json Cargo.toml pom.xml build.gradle pyproject.toml setup.py build.sbt *.csproj 2>/dev/null
```

**Detection signal mapping:**

| Marker File       | Language       | Convention scope |
| ----------------- | -------------- | ---------------- |
| `go.mod`          | Go             | `go`             |
| `package.json`    | JavaScript/TS  | `javascript`     |
| `Cargo.toml`      | Rust           | `rust`           |
| `pom.xml`         | Java           | `java`           |
| `build.gradle`    | Java/Groovy    | `java`           |
| `pyproject.toml`  | Python         | `python`         |
| `setup.py`        | Python         | `python`         |
| `build.sbt`       | Scala          | `scala`          |
| `*.csproj`        | C# / .NET      | `dotnet`         |

**Classification algorithm:**

1. Check for each marker file's existence in the project root.
2. Collect all detected languages into `detected_languages`.
3. If `--scope` was provided: use that scope directly, skip language detection.
4. If `detected_languages` is empty: output error "No known project markers detected. Expected one of: go.mod, package.json, Cargo.toml, pom.xml, pyproject.toml, etc." and ask user to specify `--scope`.
5. If `detected_languages` has exactly one entry: use it as `target_scope`.
6. If `detected_languages` has multiple entries: list all detected languages and ask the user to select which one(s) to generate Conventions for (one Convention per language).

#### 1b. Detect framework details (warm start)

For each detected language, probe deeper for framework-specific signals:

| Language       | Probe Command                                                                    | Signal                     |
| -------------- | -------------------------------------------------------------------------------- | -------------------------- |
| Go             | `grep -l 'ginkgo' go.mod 2>/dev/null`                                            | Ginkgo if present          |
| Go             | `grep -l 'testify' go.mod 2>/dev/null`                                           | testify (default)          |
| JavaScript/TS  | `node -e "const d=JSON.parse(require('fs').readFileSync('package.json'));console.log(Object.keys(d.devDependencies||{}).concat(Object.keys(d.dependencies||{})).filter(p=>/vitest|jest|mocha|cypress|playwright/.test(p)).join(' '))" 2>/dev/null` | Framework from deps |
| Rust           | Default: `cargo test`                                                            | Standard                   |
| Python         | `grep -l 'pytest\|unittest\|nose' pyproject.toml setup.py 2>/dev/null`           | pytest/unittest            |
| Java           | `grep -l 'junit\|testng\|spock' pom.xml build.gradle 2>/dev/null`                | JUnit/TestNG/Spock         |

Record detected frameworks in `detected_frameworks`. This is a warm-start signal -- it narrows the candidate list but does NOT override Step 2's test file analysis.

### Step 2: Scan Test Files & Extract Patterns

Scan existing test files to extract concrete patterns. This is the most important step -- real test code is the strongest signal.

#### 2a. Locate test files

Use Glob to find test files by language-specific patterns:

| Language       | Test file patterns                              |
| -------------- | ----------------------------------------------- |
| Go             | `**/*_test.go`                                  |
| JavaScript/TS  | `**/*.test.{ts,js,tsx,jsx}`, `**/*.spec.{ts,js,tsx,jsx}` |
| Rust           | `**/tests/*.rs` (integration), `**/*_test.rs`   |
| Python         | `**/test_*.py`, `**/*_test.py`                  |
| Java           | `**/*Test.java`, `**/*Tests.java`               |
| C#             | `**/*Tests.cs`, `**/*Test.cs`                   |

Focus on `tests/` and `tests/e2e/` directories first (forge convention), then project-wide.

#### 2b. Extract patterns from test files

For each test file found, read and extract the following:

**Imports** -- identify assertion library and test framework:
- Look for `import` / `require` / `#include` statements
- Match against known assertion libraries:
  - Go: `"github.com/stretchr/testify/assert"`, `"github.com/stretchr/testify/require"`, `"github.com/onsi/gomega"`, `. "github.com/onsi/ginkgo/v2"`
  - JS/TS: `from 'vitest'`, `from '@jest/globals'`, `from 'mocha'`, `from 'chai'`
  - Python: `import pytest`, `import unittest`, `from assertpy import assert_that`
  - Java: `import static org.junit.Assert.*`, `import static org.assertj.core.api.*`, `import org.testng.Assert`

**Tags / markers** -- identify test categorization convention:
- Go: `//go:build e2e`, `//go:build feature`, `// +build e2e`
- JS/TS: `describe('@feature', ...)`, `describe('@e2e', ...)`, `{ tags: ['@feature'] }`
- Python: `@pytest.mark.e2e`, `@pytest.mark.feature`
- Java: `@Tag("e2e")`, `@Tag("feature")`

**Test function naming** -- identify naming conventions:
- Go: `TestTC_NNN_Description`, `TestFeatureName`
- JS/TS: `describe('Feature: ...', () => { it('should ...', ...) })`
- Python: `test_tc_nnn_description`, `test_feature_name`
- Java: `testTcNnnDescription`, `shouldDoSomethingWhenCondition`

**Assertion style** -- identify which assertion functions are actually used:
- Extract function names from assert/expect/should calls
- Note frequency to identify the primary assertion library

#### 2c. Compile findings

Summarize extracted patterns into a structured finding:

```
Detected patterns:
  Framework: <framework name from imports>
  Assertion library: <library and style from imports + function calls>
  Test tags: <tag format from test files>
  Test naming: <naming pattern from function names>
  File pattern: <test file extension/pattern>
  Result format: <inferred from framework -- e.g., go test -json, vitest --reporter=json>
```

**If no test files found** (cold start): proceed to Step 3 with framework candidates only (no extracted patterns).

### Step 3: Present Findings & Confirm

Present the analysis results to the user for confirmation.

#### 3a. Warm start (test files found)

Present the extracted patterns and ask for confirmation:

```
Test Convention Analysis for: <scope>

Detected signals:
  Language: <language> (from <marker file>)
  Framework: <framework> (from test imports)

Extracted patterns from <N> test files:
  Assertion library: <library> (e.g., assert from testify)
  Key assertion functions: <function list>
  Test tags: <tag format> (e.g., //go:build e2e)
  Test naming: <pattern>
  File pattern: <pattern>

Proposed Convention file: docs/conventions/testing-<scope>.md

Sections:
  Framework: <framework name + file pattern + package>
  Assertion: <library + key functions>
  Tags: <tag format>
  Result Format: <output flags + format type>

Confirm? (y/n/edit)
```

- **y**: proceed to Step 4 to write the file.
- **n**: abort without writing.
- **edit**: ask user what to change, then re-present.

#### 3b. Cold start (no test files found)

Present framework candidates for the detected language:

```
No existing test files found for: <scope>

Language detected: <language> (from <marker file>)

Select a test framework:
  1. <primary framework> (most common for <language>)
  2. <alternative 1>
  3. <alternative 2>
  4. Custom (specify framework name)

Enter choice (1-4):
```

**Framework candidates by language:**

| Language      | Primary         | Alternatives                          |
| ------------- | --------------- | ------------------------------------- |
| Go            | go testing + testify | Ginkgo v2 + Gomega               |
| JavaScript/TS | Vitest          | Jest, Mocha + Chai                    |
| Rust          | cargo test      |                                      |
| Python        | pytest          | unittest, nose2                       |
| Java          | JUnit 5         | TestNG, Spock                         |
| C#            | xUnit           | NUnit, MSTest                         |

After user selects framework, present the proposed Convention sections (using LLM knowledge of that framework's defaults) and ask for confirmation.

#### 3c. Existing Convention file (from Step 0)

If an existing Convention file was found in Step 0:

1. Present the diff between the existing Convention content and the new proposed content.
2. Ask the user: "An existing Convention file was found at `docs/conventions/testing-<scope>.md`. The proposed changes are shown above. (a)ccept update / (k)eep existing / (e)dit"
   - **a**: proceed to Step 4 to overwrite.
   - **k**: abort, keep existing file unchanged.
   - **e**: ask user what to change, then re-present.

### Step 4: Write Convention File

Write the Convention file with the confirmed content.

#### 4a. Ensure directory exists

```bash
mkdir -p docs/conventions
```

#### 4b. Write Convention file

Write `docs/conventions/testing-<scope>.md` with the following structure.

<HARD-RULE>
The Convention file MUST follow the fixed section structure defined in the tech-design Data Models:
Framework, Assertion, Tags, Result Format (required sections).
Optional sections: Helpers, Import Patterns, Code Style, Anti-patterns.
</HARD-RULE>

The file must include:

1. **YAML frontmatter** with `title` and `domains`:
   ```yaml
   ---
   title: "<Framework Name> Testing Convention"
   domains: [testing, <scope>]
   ---
   ```

2. **Auto-generated marker**: The first comment after frontmatter MUST be:
   ```markdown
   <!-- auto-generated by forge:test-guide -->
   ```

3. **Required sections** (minimum set -- always present):

   **Framework**:
   ```markdown
   ## Framework

   - **Name**: <framework name + assertion library>
   - **File pattern**: <test file pattern>
   - **Package**: <test package name, e.g., e2e>
   - **Test runner**: <test runner command>
   - **Build tag**: <tag syntax if applicable>
   ```

   **Assertion**:
   ```markdown
   ## Assertion

   - **Library**: <library name and import>
   - **Key functions**:
     - <function signature> -- <purpose>
     - <function signature> -- <purpose>
   - **Rule**: <usage rule, e.g., "Always use assert, never require">
   ```

   **Tags**:
   ```markdown
   ## Tags

   - **Build tag**: <tag syntax with example>
   ```
   Include a code block showing the tag in context (e.g., with package declaration for Go).

   **Result Format**:
   ```markdown
   ## Result Format

   - **Output flags**: <output command flags>
   - **Format type**: <json-stream | json-report | text-verbose>
   - **Execution command**: <full command to run tests>
   ```

4. **Optional sections** -- include if extracted patterns provide enough detail, or the framework has well-known conventions:

   - **Import Patterns**: Standard import blocks for e2e tests.
   - **Code Style**: Test function naming, table-driven test patterns, traceability conventions.
   - **Anti-patterns**: Framework-specific forbidden patterns.
   - **Helpers**: Common helper functions (e.g., `runCLI`, `withRetry`).

   For cold start (no existing tests), include only the required sections. The user can add optional sections later by editing the Convention file directly or re-running `/forge:test-guide` after creating test files.

#### 4c. Report result

After writing the file:

```
Created: docs/conventions/testing-<scope>.md

Sections:
  Framework: <framework name>
  Assertion: <assertion library>
  Tags: <tag format>
  Result Format: <format type>

Next steps:
  - Run `/forge:init-justfile` to generate e2e recipes using this Convention
  - Run `/forge:gen-test-scripts` to generate tests using this Convention
  - Edit the Convention file to add optional sections (Code Style, Anti-patterns, Helpers)
```

If the file was overwritten (existing Convention updated):

```
Updated: docs/conventions/testing-<scope>.md

Changes:
  <summary of what changed>
```

## File Signal Reference

Complete reference for file signal detection used in Step 1:

| Signal File           | Language      | Framework candidates                            |
| --------------------- | ------------- | ----------------------------------------------- |
| `go.mod`              | Go            | go testing, Ginkgo                              |
| `package.json`        | JavaScript/TS | Vitest, Jest, Mocha, Cypress, Playwright        |
| `Cargo.toml`          | Rust          | cargo test (built-in)                           |
| `pom.xml`             | Java          | JUnit 4/5, TestNG, Spock                        |
| `build.gradle`        | Java/Groovy   | JUnit 4/5, TestNG, Spock, GroovyTestCase        |
| `pyproject.toml`      | Python        | pytest, unittest, nose2                         |
| `setup.py`            | Python        | pytest, unittest, nose2                         |
| `build.sbt`           | Scala         | ScalaTest, specs2, munit                        |
| `*.csproj`            | C# / .NET     | xUnit, NUnit, MSTest                            |
| `go.sum`              | Go            | (secondary -- use go.mod instead)               |
| `package-lock.json`   | JavaScript/TS | (secondary -- use package.json instead)         |
| `yarn.lock`           | JavaScript/TS | (secondary -- use package.json instead)         |

## Notes

- **No code execution**: This skill is entirely LLM-driven file analysis and generation. It reads files and writes files. It does NOT run `go test`, `npm test`, or any other test command.
- **Convention file structure**: Generated files follow the fixed section structure from tech-design.md Data Models. The four required sections (Framework, Assertion, Tags, Result Format) are always present. Optional sections are included based on available signal strength.
- **Multi-framework projects**: If the project uses multiple languages (e.g., Go backend + TypeScript frontend), generate separate Convention files for each. The user selects which languages to generate in Step 1.
- **Existing Convention files**: When a Convention file already exists, the skill presents a diff and asks for confirmation. It never silently overwrites without `--force`.
- **Cold start**: When no test files exist, the skill lists mainstream frameworks for the detected language and asks the user to select. The generated Convention file includes only required sections with LLM-inferred defaults.

<EXTREMELY-IMPORTANT>
- MANUAL-ONLY. Do NOT auto-invoke this skill from other skills or agents. Only invoke when user explicitly runs `/forge:test-guide`.
- The Convention file MUST include the `<!-- auto-generated by forge:test-guide -->` marker immediately after the YAML frontmatter.
- The Convention file MUST follow the fixed section structure: Framework, Assertion, Tags, Result Format (required) with optional Helpers, Import Patterns, Code Style, Anti-patterns sections.
- Do NOT execute any test commands. This skill is file analysis and generation only.
- If an existing Convention file is found and `--force` is NOT set, you MUST present the diff and ask for user confirmation before overwriting.
- For cold start (no test files), you MUST list framework candidates and ask the user to select. Do NOT silently default to a framework.
</EXTREMELY-IMPORTANT>
