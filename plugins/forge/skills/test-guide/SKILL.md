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

Detect language from marker files and apply the classification algorithm per `rules/signal-detection.md`.

#### 1b. Detect framework details (warm start)

For each detected language, probe for framework-specific signals per `rules/signal-detection.md`. Record detected frameworks in `detected_frameworks`. This is a warm-start signal -- it narrows the candidate list but does NOT override Step 2's test file analysis.

### Step 2: Scan Test Files & Extract Patterns

Scan existing test files to extract concrete patterns. This is the most important step -- real test code is the strongest signal.

#### 2a. Locate test files

Use Glob to find test files by language-specific patterns per `rules/pattern-extraction.md`. Focus on `tests/` and `tests/e2e/` directories first (forge convention), then project-wide.

#### 2b. Extract patterns from test files

For each test file found, extract imports, tags/markers, test function naming, and assertion style per `rules/pattern-extraction.md`.

#### 2c. Compile findings

Summarize extracted patterns into a structured finding per `rules/pattern-extraction.md`.

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

Present framework candidates for the detected language per `rules/convention-structure.md` cold start table. Ask user to select from listed options.

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

Write `docs/conventions/testing-<scope>.md` following the fixed section structure per `rules/convention-structure.md`.

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

See `rules/signal-detection.md` for the complete file signal detection reference.

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
