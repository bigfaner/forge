---
name: test-guide
description: Generate per-surface test Convention files. Reads Surface config from .forge/config.yaml, renders surface strategy templates, and produces a top-level testing quick-reference index.
allowed-tools: Bash Read Write Edit Grep Glob
disable-model-invocation: true
argument-hint: '[--force]'
---

# /forge:test-guide

MANUAL-ONLY. Do NOT auto-invoke -- only when user explicitly asks to `/forge:test-guide`.

Generate per-surface test Convention files (`docs/conventions/testing/{surface}/core.md`) by reading Surface configuration via `forge surfaces` CLI, rendering built-in surface strategy templates, and producing a top-level testing quick-reference index.

## Parameters

| Parameter  | Values  | Default | Description                                      |
| ---------- | ------- | ------- | ------------------------------------------------ |
| `--force`  | (flag)  | false   | Overwrite existing Convention files without confirmation |

## Process Flow

```
0. Check existing Convention -> 1. Read Surface config -> 2. Detect framework (auxiliary) -> 3. Render per-surface core.md -> 4. Generate index.md -> 5. User review -> 6. Write files
```

### Step 0: Check Existing Convention Files

Check whether Convention files already exist for any surface.

1. Glob `docs/conventions/testing/*/core.md` in the project root.
2. Read each file's YAML frontmatter `title` field.
3. For each existing file:
   - Read its full content.
   - Record it as `existing_conventions[surface]` for diff comparison in Step 5.
   - If `--force` is set: proceed to Step 1 (will overwrite).
   - If `--force` is NOT set: note existing files and proceed -- Step 5 will present diffs.

**If no existing Convention files found**: proceed to Step 1 (fresh generation).

### Step 1: Read Surface Configuration

Read the project's Surface configuration to determine which surfaces need Convention files.

#### 1a. Parse `forge surfaces` output

Run `forge surfaces` to obtain surface configuration. This ensures a consistent data source with other skills.

```bash
forge surfaces
```

**Parsing rule**: Use the unified `forge surfaces` text output parsing rule (see Forge Guide → Surface Output Parsing).

Parse each line to build `active_surfaces` array. Each entry has `key` (may be empty for scalar) and `type`.

#### 1b. Validate surface types

Each surface type must be one of the supported types:

| Surface Type | Strategy Template | Test Type |
|--------------|-------------------|-----------|
| `cli` | `templates/surfaces/cli.md` | CLI Functional Test |
| `api` | `templates/surfaces/api.md` | API Functional Test |
| `web` | `templates/surfaces/web.md` | Web E2E Test |
| `tui` | `templates/surfaces/tui.md` | Terminal Functional Test |
| `mobile` | `templates/surfaces/mobile.md` | Mobile E2E Test |

If an unsupported surface type is found: output error `"Unsupported surface: <type>. Supported: cli, api, web, tui, mobile"` and abort.

If `forge surfaces` returns empty output or fails: output error `"No surfaces configured. Run 'forge init' or manually add surfaces to .forge/config.yaml."` and abort.

### Step 2: Detect Framework (Auxiliary Step)

Detect the project's test framework(s) from file system signals. This step does NOT execute any code -- it reads file names, dependency lists, and file contents. Its results are used ONLY to fill the assertion preference table in core.md.

#### 2a. Detect language and framework signals

Follow the detection algorithm in `rules/signal-detection.md`:

1. Detect language from marker files.
2. For each detected language, probe for framework-specific signals.
3. Apply false-positive exclusion rules.
4. Record detected frameworks with confidence levels.

#### 2b. Map framework to assertion preference

For each detected framework, produce an assertion preference row:

| Detected Framework | Assertion Library | Mock Mechanism | Fixture Pattern |
|--------------------|-------------------|----------------|-----------------|
| Go testing + testify | testify/assert | interfaces + test doubles | TestMain setup |
| Ginkgo + Gomega | Gomega Expect | Ginkgo fake | BeforeEach/AfterEach |
| Vitest | expect (built-in) | vi.mock() | beforeEach/fixture |
| Jest | expect (built-in) | jest.mock() | beforeEach |
| pytest | Python assert + pytest.raises | unittest.mock / pytest-mock | @pytest.fixture |
| cargo test | assert! macro | mockall | #[cfg(test)] module |
| JUnit 5 | JUnit 5 Assertions | Mockito | @BeforeEach |

#### 2c. Cold start handling

If no test files exist (medium or low confidence detection): use the default assertion preference row from `rules/convention-structure.md` Cold Start Framework Candidates table for the detected language.

### Step 3: Render Per-Surface Convention Files

For each surface in `active_surfaces`, generate a Convention file from the corresponding template. Each entry has `type` (used for template and directory naming) and optionally `key` (non-empty for named surfaces).

#### 3a. Load surface strategy template

Read the template file for the current surface from `templates/surfaces/<type>.md` (where `<type>` is the surface type from the parsing step).

#### 3b. Fill assertion preference table

Replace the template placeholder row in the `## 断言偏好表` section:

Template row:
```
| {{ASSERTION_LIBRARY}} | {{MOCK_MECHANISM}} | {{FIXTURE_PATTERN}} |
```

Replace with the detected framework's assertion preference from Step 2b. If multiple frameworks are detected, add one row per framework.

If no framework was detected (completely cold start): keep the placeholder row and add a comment:
```
<!-- TODO: Run /forge:test-guide after adding test dependencies to fill this table -->
```

#### 3c. Generate index.md for the surface

Create `docs/conventions/testing/<type>/index.md` as a file index pointing to `core.md`:

```markdown
---
title: "<Surface Name> 测试约定"
domains: [testing, <surface>]
---

<!-- auto-generated by forge:test-guide -->

# <Surface Name> 测试约定

- [测试策略 (core.md)](core.md)
```

### Step 4: Generate Top-Level Index

Generate `docs/conventions/testing/index.md` as a quick-reference table covering all active surfaces.

#### 4a. Build quick-reference table

Create the index with a summary table:

```markdown
---
title: "测试约定速查表"
domains: [testing]
---

<!-- auto-generated by forge:test-guide -->

# 测试约定速查表

| Surface | 测试类型 | 文件位置 | 断言重点 | 详细策略 |
|---------|---------|---------|---------|---------|
| cli | CLI 功能测试 | tests/<surfaceKey>/<journey>/ (多 surface) 或 tests/<journey>/ (单 surface) | exit code + stdout + stderr | [cli/core.md](cli/core.md) |
| api | API 功能测试 | tests/<surfaceKey>/<journey>/ (多 surface) 或 tests/<journey>/ (单 surface) | status code + response body + headers | [api/core.md](api/core.md) |
| web | Web E2E 测试 | tests/<surfaceKey>/<journey>/ (多 surface) 或 tests/<journey>/ (单 surface) | DOM 可见性 + 用户操作 + URL | [web/core.md](web/core.md) |
| tui | 终端功能测试 | tests/<surfaceKey>/<journey>/ (多 surface) 或 tests/<journey>/ (单 surface) | 精确文本 + 正则匹配 + 快照 | [tui/core.md](tui/core.md) |
| mobile | Mobile E2E 测试 | tests/<surfaceKey>/<journey>/ (多 surface) 或 tests/<journey>/ (单 surface) | UI 可见性 + 操作响应 + 屏幕 ID | [mobile/core.md](mobile/core.md) |
```

Only include rows for surfaces present in `active_surfaces`.

### Step 5: Present Drafts & User Review

Present all generated Convention files to the user for review.

#### 5a. Present summary

```
Convention Drafts Generated
============================

Surfaces: <active_surfaces list>
Framework detected: <framework name> (confidence: <level>)

Files to write:
  docs/conventions/testing/index.md (quick-reference table)
  <for each surface>:
    docs/conventions/testing/<surface>/index.md
    docs/conventions/testing/<surface>/core.md

Review options:
  (a)ccept - Write all files
  (e)dit   - Tell me what to change
  (r)eject - Discard and regenerate with feedback
```

#### 5b. Handle user response

Initialize retry counter: `retry_count = 0`, `max_retries = 2`.

- **(a)ccept**: proceed to Step 6 to write all files.
- **(e)dit**:
  1. Ask: "What would you like to change?"
  2. Apply the requested changes to the affected draft(s).
  3. Re-present the updated summary with diff markers.
  4. Do NOT increment retry counter.
  5. Return to 5b for next response.
- **(r)eject with feedback**:
  1. Parse user feedback to identify approved and rejected files/sections.
  2. Regenerate ONLY rejected files/sections, preserving approved content verbatim.
  3. Increment `retry_count`.
  4. If `retry_count <= max_retries`: re-present summary, return to 5b.
  5. If `retry_count > max_retries`: proceed to 5c (retry exhausted).

#### 5c. Retry exhausted

After 2 retries are exhausted and the user still rejects:

1. Mark all files as `[DRAFT - needs manual review]`.
2. Write all drafts to `docs/conventions/testing/*.draft.md` (one per surface).
3. Output:

```
Retry limit reached. Drafts written to:
  docs/conventions/testing/<surface>/core.draft.md
  ...

Please manually edit the files, then rename from .draft.md to .md.
The pipeline will wait for confirmation before proceeding.
```

#### 5d. Existing Convention files (from Step 0)

If existing Convention files were found in Step 0:

1. Present the diff between existing content and new draft for each file.
2. Ask: "Existing Convention files found. (a)ccept update / (k)eep existing / (e)dit"
   - **a**: proceed to Step 6 to overwrite.
   - **k**: abort, keep existing files unchanged.
   - **e**: ask user what to change, then re-present (counts toward retry limit).

### Step 6: Write Convention Files

Write all confirmed Convention files.

#### 6a. Ensure directory structure exists

```bash
mkdir -p docs/conventions/testing
for surface_type in <active_surfaces types>; do
  mkdir -p "docs/conventions/testing/$surface_type"
done
```

#### 6b. Write files

For each surface in `active_surfaces` (using `type` for directory naming):
- Write `docs/conventions/testing/<type>/index.md` following the format from Step 3c.
- Write `docs/conventions/testing/<type>/core.md` following the Convention structure per `rules/convention-structure.md`.

Write `docs/conventions/testing/index.md` following the format from Step 4a.

Every generated file MUST include the auto-generated marker:

```markdown
<!-- auto-generated by forge:test-guide -->
```

#### 6c. Report result

After writing all files:

```
Created: docs/conventions/testing/index.md (quick-reference table)

<for each surface>:
  Created: docs/conventions/testing/<surface>/index.md
  Created: docs/conventions/testing/<surface>/core.md
    Sections: 文件位置, 隔离模型, 断言重点, 超时策略, 生命周期, Contract/Journey 比例, 反模式, 断言偏好表
    Framework: <detected framework>

Next steps:
  - Run `/forge:init-justfile` to generate test recipes using these Conventions
  - Run `/forge:gen-test-scripts` to generate tests using these Conventions
  - Run `/forge:run-tests` to execute tests per surface
```

## Surface Strategy Templates

Per-surface strategy templates are located at `templates/surfaces/<surface>.md`. Each template defines 7 mandatory sections plus an assertion preference table:

1. **文件位置**: Test directory, file naming, build tags
2. **隔离模型**: Isolation model specific to the surface
3. **断言重点**: Assertion dimensions and minimum requirements
4. **超时策略**: Timeout strategy at multiple levels
5. **生命周期**: Test lifecycle steps
6. **Contract/Journey 比例**: Balance between Contract and Journey tests
7. **反模式**: Common anti-patterns and alternatives
8. **断言偏好表**: Per-framework assertion, mock, and fixture preferences (filled from Step 2)

## Framework Detection Reference

See `rules/signal-detection.md` for the complete file signal detection reference. Framework detection is an auxiliary step -- its results populate only the assertion preference table in core.md, not the overall Convention structure.

## Test Type Model Reference

See `references/test-type-model.md` for the complete Surface -> Test Type mapping model, including classification criteria and "e2e" terminology constraints.

## Notes

- **No code execution**: This skill is entirely LLM-driven file analysis and generation. It reads files and writes files. It does NOT run `go test`, `npm test`, or any other test command.
- **Surface-first organization**: Convention files are organized by Surface (cli/api/web/tui/mobile), not by framework. Framework information is reduced to a single row in the assertion preference table within each surface's core.md.
- **Multi-surface projects**: If the project has multiple surfaces (e.g., cli + api), generate separate core.md for each surface. The top-level index.md provides a unified quick-reference. Test directory follows surface-key adaptive rules: multi-surface projects use `tests/<surfaceKey>/<journey>/`, single-surface projects (scalar or named) use `tests/<journey>/`. The `surfaceKey` is the key from `forge surfaces` output (e.g., `backend`, `frontend`), not the surface type (e.g., `api`, `web`).
- **Existing Convention files**: When Convention files already exist, the skill presents diffs and asks for confirmation. It never silently overwrites without `--force`.
- **Cold start**: When no test files exist, the skill uses dependency detection for framework identification and fills the assertion preference table from defaults.
- **Draft feedback loop**: User rejections trigger regeneration of only the rejected files/sections. After 2 retries, drafts are written as `.draft.md` for manual editing.

<EXTREMELY-IMPORTANT>
- When regenerating after rejection, preserve approved files/sections verbatim and regenerate ONLY rejected files/sections.
- After 2 retry rejections, write drafts as `.draft.md` for manual editing. Do NOT force-apply.
</EXTREMELY-IMPORTANT>
