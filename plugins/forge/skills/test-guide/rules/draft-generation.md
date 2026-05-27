# Convention Draft Generation Rules

This document defines the rules for generating Convention file drafts based on detected test framework signals.

## Overview

The draft generation process takes the output of signal detection and produces a Convention file draft that:
1. Conforms to the 4-section schema (framework, structure, assertions, tags + result format)
2. Can be validated against the Convention structure rules
3. Is presented to the user for review before being written to disk

## Schema Validation (4-Section Completeness)

Every Convention draft MUST include these required sections:

| Section | Required Fields | Purpose |
|---------|----------------|---------|
| **framework** | name, version, language, runner_command | Framework identity and execution |
| **discovery** | test_dir, file_pattern, exclude_pattern | How test files are found |
| **structure** | suite_pattern, case_pattern | Test organization patterns |
| **assertions** | style, library, key patterns | How assertions are written |
| **Tags** | Format, usage examples | Test categorization |
| **Result Format** | Output flags, format type, execution command | How results are reported |

**Validation rule**: If any required field is missing, the draft is INCOMPLETE and must not be presented to the user without noting the gap.

## Draft Generation Process

### Step 1: Map detection signals to Convention template

For each detected framework, map the signals to Convention fields:

```
Signal -> Convention Field Mapping:

Marker file        -> framework.language (via marker-to-language table)
Framework name     -> framework.name
Framework version  -> framework.version (from dependency version if available, else "latest")
Test file pattern  -> discovery.file_pattern
Test directory     -> discovery.test_dir
Import patterns    -> assertions.library, assertions.style
Tag format         -> Tags section
```

### Step 2: Fill from built-in Convention templates

If a built-in Convention file exists in `docs/conventions/testing/` for the detected framework, use it as the primary template:

| Detected Framework | Built-in Template |
|--------------------|-------------------|
| Go testing         | `docs/conventions/testing/go.md` |
| Ginkgo             | `docs/conventions/testing/ginkgo.md` |
| Vitest             | `docs/conventions/testing/vitest.md` |
| pytest             | `docs/conventions/testing/pytest.md` |
| JUnit              | `docs/conventions/testing/junit.md` |
| Rust               | `docs/conventions/testing/rust.md` |

**When built-in template exists**: Copy the template and customize fields based on extracted patterns (test_dir, file_pattern, naming conventions from actual test files).

**When no built-in template exists**: Generate from scratch using LLM knowledge of the framework + the extracted patterns from signal detection.

### Step 3: Customize with extracted patterns

Override template defaults with patterns extracted from actual test files (Step 2 of SKILL.md process):

| Extracted Pattern | Convention Field Override |
|-------------------|--------------------------|
| Actual test file paths | `discovery.test_dir` |
| Actual import statements | `assertions.library` |
| Actual test function names | `structure.case_pattern` |
| Actual tag formats | `Tags section` |
| Actual assertion functions used | `assertions.key patterns` |

### Step 4: Validate draft completeness

Before presenting to the user, validate the draft against these rules:

1. **Framework section**: name, language, runner_command are non-empty
2. **Discovery section**: file_pattern is non-empty
3. **Structure section**: suite_pattern is non-empty
4. **Assertions section**: style and library are non-empty
5. **Tags section**: tag format is defined
6. **Result Format section**: execution command is non-empty

If validation fails, fill gaps with LLM-inferred defaults for the detected framework and flag them as "inferred (not detected from project)".

## Draft Output Format

Present the draft to the user in this format:

```
Convention Draft for: <framework name>

Generated from: <detection signals summary>
Template source: <built-in template name or "LLM-generated">

--- Draft Content ---

<full Convention file content with YAML frontmatter and all sections>

--- End Draft ---

Validation: ALL required sections present / MISSING: <list>
Confidence: <high/medium/low>

Review options:
  (a)ccept - Write to docs/conventions/testing/<scope>.md
  (e)dit   - Tell me what to change
  (r)eject - Discard and start over
```

## User Review and Feedback Loop

### Acceptance

When user accepts (`a`):
1. Write the draft to `docs/conventions/testing/<scope>.md`
2. Output confirmation with file path and section summary
3. Proceed to the report step in SKILL.md

### Edit Request

When user requests edits (`e`):
1. Ask: "What would you like to change?"
2. Apply the requested changes to the draft
3. Re-present the updated draft with diff markers for changed sections
4. Reset retry counter (edit requests do not count as rejections)

### Rejection with Feedback

When user rejects with feedback (`r`):
1. Parse the user's feedback to identify:
   - **Approved sections**: Sections the user did not mention or explicitly approved
   - **Rejected sections**: Sections the user identified as incorrect
2. Regenerate ONLY the rejected sections, preserving approved sections verbatim
3. Re-present the updated draft
4. Increment retry counter

**Retry limit**: Maximum 2 retries (3 total presentations including the initial one).

### Retry Exhausted

After 2 retries are exhausted and the user still rejects:
1. Output the current draft with all sections marked as `[DRAFT - needs manual review]`
2. Write the draft to `docs/conventions/testing/<scope>.draft.md` (note the `.draft.md` extension)
3. Output message:

```
Retry limit reached. Draft written to:
  docs/conventions/testing/<scope>.draft.md

Please manually edit the file, then rename to:
  docs/conventions/testing/<scope>.md

The pipeline will wait for confirmation before proceeding.
```

4. Ask the user: "Have you finished editing? (y)es / (n)o, I'll edit later"
   - `y`: Rename `.draft.md` to `.md`, continue pipeline
   - `n`: Abort, user will manually handle the draft file

## File Write Rules

### Directory

Convention drafts are written to `docs/conventions/testing/` directory.

```bash
mkdir -p docs/conventions/testing
```

### File naming

- Accepted drafts: `docs/conventions/testing/<scope>.md` (e.g., `go.md`, `pytest.md`, `vitest.md`)
- Pending drafts (retry exhausted): `docs/conventions/testing/<scope>.draft.md`

### Auto-generated marker

Every generated Convention file MUST include the auto-generated marker immediately after YAML frontmatter:

```markdown
---
title: "<Framework Name> Testing Convention"
---

<!-- auto-generated by forge:test-guide -->
```

### Existing file handling

If a Convention file already exists at the target path:
- Without `--force`: Present a diff and ask for confirmation (handled by SKILL.md Step 0)
- With `--force`: Overwrite directly

## Draft Quality Guidelines

### Section generation priority

When generating from scratch (no built-in template), follow this priority:

1. **framework**: Use exact dependency name and version from detection
2. **discovery**: Use actual file patterns found in the project
3. **structure**: Infer from test file naming conventions in the project
4. **assertions**: Use actual assertion library detected from imports
5. **Tags**: Use actual tag format found in test files, or infer from framework conventions
6. **Result Format**: Use framework's standard output format (JSON for CI integration preferred)

### Inferred defaults

When a field cannot be determined from project signals, use these defaults:

| Framework   | Default runner          | Default assertion        | Default tag format        |
|-------------|------------------------|--------------------------|---------------------------|
| Go testing  | `go test -v -json`     | testify/assert           | `//go:build <surface>-<type>`          |
| Ginkgo      | `ginkgo -v --json-report=report.json` | Gomega Expect | `Label("<surface>-<type>")` |
| Vitest      | `vitest run --reporter=verbose` | `expect` from vitest | `describe('@tag', ...)` |
| Jest        | `jest --verbose`       | `expect` from @jest/globals | `describe('@tag', ...)` |
| pytest      | `pytest -v`            | Python assert            | `@pytest.mark.<tag>`     |
| JUnit 5     | `mvn test`             | JUnit 5 Assertions       | `@Tag("name")`           |
| Rust        | `cargo test`           | `assert!` macro          | `#[cfg(feature = "<surface>-<type>")]` |
