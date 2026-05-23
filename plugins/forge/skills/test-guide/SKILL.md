---
name: test-guide
description: Guide users through generating Convention files. Auto-detects test frameworks from project signals, generates Convention drafts, and writes after user review with retry feedback loop.
allowed-tools: Bash Read Write Edit Grep Glob
disable-model-invocation: true
argument-hint: '[--scope <scope>] [--force]'
---

# /forge:test-guide

MANUAL-ONLY. Do NOT auto-invoke -- only when user explicitly asks to `/forge:test-guide`.

Guide users through creating test Convention files (`docs/conventions/testing/<scope>.md`). The skill auto-detects project test frameworks from file signals, generates a Convention draft conforming to the 4-section schema, and writes the file after user review.

## Parameters

| Parameter  | Values                    | Default         | Description                                        |
| ---------- | ------------------------- | --------------- | -------------------------------------------------- |
| `--scope`  | language or framework name | (auto-detect)  | Override the Convention file scope (e.g., `go`, `javascript`, `python`) |
| `--force`  | (flag)                    | false           | Overwrite existing Convention file without confirmation |

## Process Flow

```
0. Check existing Convention -> 1. Auto-detect framework -> 2. Scan test files -> 3. Generate draft -> 4. User review -> 5. Write Convention
```

### Step 0: Check Existing Convention Files

Check whether a Convention file already exists for the target scope.

1. Glob `docs/conventions/testing/*.md` in the project root (exclude `index.md` and `*.draft.md`).
2. Read each file's YAML frontmatter `title` field.
3. If a Convention file with matching scope already exists:
   - Read its full content.
   - Record it as `existing_convention` for diff comparison in Step 4.
   - If `--force` is set: proceed to Step 1 (will overwrite).
   - If `--force` is NOT set: note the existing file and proceed -- Step 4 will present a diff.

**If no existing Convention files found**: proceed to Step 1 (fresh generation).

### Step 1: Auto-Detect Framework

Detect the project's test framework from file system signals. This step does NOT execute any code -- it reads file names, dependency lists, and file contents.

#### 1a. Detect language signals

Detect language from marker files and apply the classification algorithm per `rules/signal-detection.md`.

#### 1b. Detect framework details

For each detected language, probe for framework-specific signals per `rules/signal-detection.md`. Each framework requires three signal types to confirm detection:
1. Marker file (language marker present)
2. Dependency (framework listed in dependency file)
3. Test files (file pattern matches exist)

Record detected frameworks with confidence levels per `rules/signal-detection.md`.

#### 1c. Apply false-positive exclusion

Apply the exclusion rules from `rules/signal-detection.md` to eliminate misclassification:
- Framework dependency overrides generic dependency
- Test file patterns are language-specific (React dependency does not imply test framework)
- Ginkgo vs go testing precedence
- pytest vs unittest distinction
- Build tool ambiguity resolution

#### 1d. Handle detection results

- **High confidence** (all 3 signals confirmed): proceed to Step 2 with detected framework.
- **Medium confidence** (2 signals, usually missing test files): proceed to Step 2 with detected framework, note cold start.
- **Low confidence** (marker only): present cold start candidate list from `rules/convention-structure.md`, ask user to select framework.

### Step 2: Scan Test Files & Extract Patterns

Scan existing test files to extract concrete patterns. This step provides the strongest signal for customizing the draft.

#### 2a. Locate test files

Use Glob to find test files by language-specific patterns per `rules/pattern-extraction.md`. Focus on `tests/` and `tests/e2e/` directories first (forge convention), then project-wide.

#### 2b. Extract patterns from test files

For each test file found, extract imports, tags/markers, test function naming, and assertion style per `rules/pattern-extraction.md`.

#### 2c. Compile findings

Summarize extracted patterns into a structured finding per `rules/pattern-extraction.md`.

**If no test files found** (cold start): proceed to Step 3 with framework detection results only (no extracted patterns).

### Step 3: Generate Convention Draft

Generate a Convention file draft based on detected framework and extracted patterns per `rules/draft-generation.md`.

#### 3a. Select template source

Check if a built-in Convention template exists in `docs/conventions/testing/` for the detected framework:

| Detected Framework | Built-in Template |
|--------------------|-------------------|
| Go testing         | `docs/conventions/testing/go.md` |
| Ginkgo             | `docs/conventions/testing/ginkgo.md` |
| Vitest             | `docs/conventions/testing/vitest.md` |
| pytest             | `docs/conventions/testing/pytest.md` |
| JUnit              | `docs/conventions/testing/junit.md` |
| Rust               | `docs/conventions/testing/rust.md` |

If a built-in template exists, use it as the base and customize with extracted patterns.
If no template exists, generate from scratch using LLM knowledge + extracted patterns.

#### 3b. Customize with extracted patterns

Override template defaults with patterns extracted in Step 2:

| Extracted Pattern | Convention Field Override |
|-------------------|--------------------------|
| Actual test file paths | `discovery.test_dir` |
| Actual import statements | `assertions.library` |
| Actual test function names | `structure.case_pattern` |
| Actual tag formats | `Tags section` |
| Actual assertion functions used | `assertions.key patterns` |

#### 3c. Validate draft completeness

Validate the draft against the 4-section schema per `rules/draft-generation.md`:
- **framework**: name, language, runner_command are non-empty
- **discovery**: file_pattern is non-empty
- **structure**: suite_pattern is non-empty
- **assertions**: style and library are non-empty
- **Tags**: tag format is defined
- **Result Format**: execution command is non-empty

If validation fails, fill gaps with LLM-inferred defaults and flag as "inferred (not detected from project)".

### Step 4: Present Draft & User Review

Present the Convention draft to the user for review. This step implements a feedback loop with retry.

#### 4a. Present the draft

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
  (r)eject - Discard and regenerate with feedback
```

#### 4b. Handle user response

Initialize retry counter: `retry_count = 0`, `max_retries = 2`.

- **(a)ccept**: proceed to Step 5 to write the file.
- **(e)dit**:
  1. Ask: "What would you like to change?"
  2. Apply the requested changes to the draft (preserving unchanged sections).
  3. Re-present the updated draft with diff markers for changed sections.
  4. Do NOT increment retry counter.
  5. Return to 4b for next response.
- **(r)eject with feedback**:
  1. Parse user feedback to identify approved sections and rejected sections.
  2. Regenerate ONLY rejected sections, preserving approved sections verbatim.
  3. Increment `retry_count`.
  4. If `retry_count <= max_retries`: re-present updated draft, return to 4b.
  5. If `retry_count > max_retries`: proceed to 4c (retry exhausted).

#### 4c. Retry exhausted

After 2 retries are exhausted and the user still rejects:

1. Mark all sections as `[DRAFT - needs manual review]`.
2. Write the draft to `docs/conventions/testing/<scope>.draft.md`.
3. Output:

```
Retry limit reached. Draft written to:
  docs/conventions/testing/<scope>.draft.md

Please manually edit the file, then rename to:
  docs/conventions/testing/<scope>.md

The pipeline will wait for confirmation before proceeding.
```

4. Ask the user: "Have you finished editing? (y)es / (n)o, I'll edit later"
   - `y`: Rename `.draft.md` to `.md`, proceed to Step 5 for confirmation.
   - `n`: Abort. User will handle the draft file manually.

#### 4d. Existing Convention file (from Step 0)

If an existing Convention file was found in Step 0:

1. Present the diff between the existing Convention content and the new draft content.
2. Ask the user: "An existing Convention file was found at `docs/conventions/testing/<scope>.md`. The proposed changes are shown above. (a)ccept update / (k)eep existing / (e)dit"
   - **a**: proceed to Step 5 to overwrite.
   - **k**: abort, keep existing file unchanged.
   - **e**: ask user what to change, then re-present (counts toward retry limit).

### Step 5: Write Convention File

Write the Convention file with the confirmed content.

#### 5a. Ensure directory exists

```bash
mkdir -p docs/conventions/testing
```

#### 5b. Write Convention file

Write `docs/conventions/testing/<scope>.md` following the Convention structure per `rules/convention-structure.md`.

Every generated file MUST include the auto-generated marker:

```markdown
---
title: "<Framework Name> Testing Convention"
---

<!-- auto-generated by forge:test-guide -->
```

#### 5c. Report result

After writing the file:

```
Created: docs/conventions/testing/<scope>.md

Sections:
  Framework: <framework name>
  Discovery: <test dir and file pattern>
  Structure: <suite and case patterns>
  Assertion: <assertion library>
  Tags: <tag format>
  Result Format: <format type>

Next steps:
  - Run `/forge:init-justfile` to generate e2e recipes using this Convention
  - Run `/forge:gen-test-scripts` to generate tests using this Convention
  - Edit the Convention file to add optional sections (Import Patterns, Code Style, Anti-patterns, Helpers)
```

If the file was overwritten (existing Convention updated):

```
Updated: docs/conventions/testing/<scope>.md

Changes:
  <summary of what changed>
```

## File Signal Reference

See `rules/signal-detection.md` for the complete file signal detection reference.

## Draft Generation Reference

See `rules/draft-generation.md` for the complete draft generation rules including schema validation, retry logic, and file write rules.

## Notes

- **No code execution**: This skill is entirely LLM-driven file analysis and generation. It reads files and writes files. It does NOT run `go test`, `npm test`, or any other test command.
- **Convention file structure**: Generated files follow the section structure from built-in Convention templates. The required sections (framework, discovery, structure, assertions, Tags, Result Format) are always present. Optional sections (Import Patterns, Code Style, Anti-patterns, Helpers) are included when signal strength supports them.
- **Multi-framework projects**: If the project uses multiple languages (e.g., Go backend + TypeScript frontend), generate separate Convention files for each. The user selects which languages to generate in Step 1.
- **Existing Convention files**: When a Convention file already exists, the skill presents a diff and asks for confirmation. It never silently overwrites without `--force`.
- **Cold start**: When no test files exist, the skill uses dependency detection for framework identification and presents candidates for user selection.
- **Draft feedback loop**: User rejections trigger regeneration of only the rejected sections. After 2 retries, the draft is written as `.draft.md` for manual editing.

<EXTREMELY-IMPORTANT>
- MANUAL-ONLY. Do NOT auto-invoke this skill from other skills or agents. Only invoke when user explicitly runs `/forge:test-guide`.
- The Convention file MUST include the `<!-- auto-generated by forge:test-guide -->` marker immediately after the YAML frontmatter.
- The Convention file MUST follow the fixed section structure: framework, discovery, structure, assertions, Tags, Result Format (required) with optional Helpers, Import Patterns, Code Style, Anti-patterns sections.
- Do NOT execute any test commands. This skill is file analysis and generation only.
- If an existing Convention file is found and `--force` is NOT set, you MUST present the diff and ask for user confirmation before overwriting.
- For cold start (no test files), you MUST list framework candidates and ask the user to select. Do NOT silently default to a framework.
- Drafts MUST be reviewed by the user before being applied. Do NOT auto-apply drafts.
- Generated drafts MUST pass the 4-section schema validation defined in `rules/draft-generation.md`.
- When regenerating after rejection, preserve approved sections verbatim and regenerate ONLY rejected sections.
- After 2 retry rejections, write the draft as `.draft.md` for manual editing. Do NOT force-apply.
</EXTREMELY-IMPORTANT>
