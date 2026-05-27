# Freeform Expert Review

**Reviewer**: Template Engine Migration & Rendering Pipeline Architect
**Review Style**: migration-risk-first
**Date**: 2026-05-28

## Background Assessment

This proposal consolidates three independent `strings.ReplaceAll`-based template rendering mechanisms across `pkg/prompt`, `pkg/template`, and `pkg/task/autogen` into Go's `text/template`. Each package currently has its own post-processing pipeline: `cleanTemplateOutput()` in prompt.go (4 conditional deletion rules + blank-line collapse), `injectSurfaceFrontmatter()` in add.go (surface field injection into YAML frontmatter), and `removeLineContaining()`/`removeSection()` in autogen.go (line- and section-level conditional deletion). The proposal also cleans up stale `{{SCOPE}}` references in 4 autogen templates and introduces surface hard-failure for quality-gate fix-task creation.

**Template inventory verification** (confirmed against codebase):
- `pkg/prompt/data/`: 21 `.md` files -- proposal correctly states 21
- `pkg/template/data/`: 2 `.md` files (coding.fix.md, coding.cleanup.md) -- proposal correctly states 2
- `pkg/task/data/`: 18 `.md` files total = 12 non-record + 6 record -- proposal correctly states 12 + 6
- Total: 35 migratable + 6 record (metadata-only) = 41 -- matches proposal throughout

**Verification of `cleanTemplateOutput` behavior** (4 conditional rules confirmed in `prompt.go:256-309`):
1. Empty label lines: `isLabelWithEmptyValue()` removes lines like `SURFACE_KEY:` or `PHASE_SUMMARY:` when the value is empty. Checks `strings.Contains(before, " ")` to avoid matching multi-word labels.
2. Empty-backtick conditional sentences: removes lines containing `If `` is non-empty` -- triggered when `{{PHASE_SUMMARY}}` substitutes to empty.
3. Trailing whitespace on `just` command lines: strips trailing spaces on lines starting with `just ` -- triggered when `{{SURFACE_KEY}}` is empty in `just compile {{SURFACE_KEY}}`.
4. `<!-- IF NOT_LOW -->...<!-- END_IF -->` paragraph blocks: removes entire marked paragraphs when complexity is "low"; strips markers only for other complexity levels. Present in 4 coding templates.

**Verification of `renderBody` behavior** (6 placeholder categories in `autogen.go:369-431`):
- `{{FEATURE_SLUG}}` -- always substituted
- `{{MODE}}` -- omit line when empty via `removeLineContaining`
- `{{SCOPE}}` -- omit `## Scope` section when empty via `removeSection`; format as `- item` list when present
- `{{SURFACES}}` -- default to "See .forge/config.yaml" when empty
- `{{TEST_TYPE}}` -- omit line when empty via `removeLineContaining`
- `{{ACCEPTANCE_CRITERIA}}` -- default to "All acceptance criteria met" when empty
- `{{DOC_TASK_AC}}` -- serialized from `map[string]string` or empty string

**Verification of `renderTemplate` behavior** (8 substitution operations in `prompt.go:114-168`):
- `{{TASK_ID}}`, `{{TASK_FILE}}` -- direct substitution
- TASK_CATEGORY injection via `strings.Replace` after TASK_FILE value (line 136-137) -- NOT a placeholder, inserts a new line after the resolved value
- `{{SURFACE_KEY}}` -- direct substitution (empty string when cross-surface)
- `{{FEATURE_SLUG}}` -- direct substitution
- `{{PHASE_SUMMARY}}` -- substitutes to `PHASE_SUMMARY: <path>` or empty string; two downstream effects: (a) label line becomes empty, cleaned by rule 1; (b) conditional sentence becomes `If `` is non-empty`, cleaned by rule 2
- `{{TEST_TYPE_ARG}}` -- substitutes to ` --type <surfaceType>` or empty string (only in test-gen-scripts.md)
- `{{COVERAGE_STRATEGY}}` / `{{COVERAGE_TARGET}}` -- direct substitution; non-testable types get empty strings; cleanup/refactor get "maintain" strategy
- `{{COMPLEXITY}}` -- defaults to "medium" when empty

## Key Risks

### 1. `promptTemplateData` struct is incomplete -- missing 5 fields

> `type promptTemplateData struct { TaskID, TaskCategory, PhaseSummary, CoverageStrategy, SurfaceKey, SurfaceType, Complexity }` (proposal lines 165-174)

**Analysis**: The proposed struct has 7 fields. But `renderTemplate()` in `prompt.go:114-168` performs 8 distinct substitutions. The following fields used by the current rendering pipeline are absent from the struct definition:

| Missing field | Current placeholder | Used in |
|---|---|---|
| `TaskFile` (string) | `{{TASK_FILE}}` | All 21 prompt templates |
| `FeatureSlug` (string) | `{{FEATURE_SLUG}}` | All 21 prompt templates |
| `TestTypeArg` (string) | `{{TEST_TYPE_ARG}}` | test-gen-scripts.md |
| `CoverageTarget` (string) | `{{COVERAGE_TARGET}}` | 5 coding templates |

Wait -- `CoverageTarget` IS listed in the struct (field name `CoverageStrategy` and the struct actually has it). Let me re-check. Looking at the proposal struct definition at line 165-174:

```go
type promptTemplateData struct {
    TaskID           string
    TaskCategory     string
    PhaseSummary     string
    CoverageStrategy string
    SurfaceKey       string
    SurfaceType      string
    Complexity       string
}
```

Fields missing from the struct:
- `TaskFile` -- `{{TASK_FILE}}` is used in every prompt template (line 2 of most templates)
- `FeatureSlug` -- `{{FEATURE_SLUG}}` is substituted in all templates (line 139 of prompt.go)
- `TestTypeArg` -- `{{TEST_TYPE_ARG}}` is used in test-gen-scripts.md (line 147 of prompt.go)
- `CoverageTarget` -- `{{COVERAGE_TARGET}}` is used in 5 coding templates (line 157 of prompt.go)

If these fields are missing from the data struct, `text/template` with `missingkey=error` will fail at Execute time for every single template. The init-time validation gate (zero-value Execute) would catch this immediately, but it indicates the struct definition was designed considering only the conditional fields, not the full substitution set.

### 2. `autogenTemplateData` struct is critically incomplete -- missing 5 of 9 data dimensions

> `type autogenTemplateData struct { TaskID, TaskType, SurfaceKey, SurfaceType, ScopeDisplay }` (proposal lines 206-213)

**Analysis**: The `renderBody()` function in `autogen.go:369-431` consumes 7 distinct placeholder categories from `BodyContext` and `AutoGenTaskDef`. The proposed 5-field struct omits:

| Missing field | Current source | Template usage |
|---|---|---|
| `FeatureSlug` | `ctx.FeatureSlug` | `{{FEATURE_SLUG}}` in all 12 autogen templates |
| `Mode` | `ctx.Mode` | `{{MODE}}` in doc-review, test-gen-contracts, test-gen-journeys |
| `SurfaceTypes` (formatted) | `ctx.SurfaceTypes` | `{{SURFACES}}` -- defaults to "See .forge/config.yaml" when empty |
| `AcceptanceCriteria` (formatted) | `ctx.AcceptanceCriteria` | `{{ACCEPTANCE_CRITERIA}}` in validation-code, validation-ux |
| `DocTaskCriteria` (formatted) | `ctx.DocTaskCriteria` | `{{DOC_TASK_AC}}` in doc-review |

The proposal mentions only `ScopeDisplay` as a pre-formatted string field but does not address how the other `[]string` and `map[string]string` fields should be handled. `AcceptanceCriteria` and `SurfaceTypes` are slices that need pre-formatting (join with `\n- ` prefix), while `DocTaskCriteria` needs serialization into markdown sub-sections. The current `renderBody` handles all of these with imperative Go code before substitution -- the proposal needs to specify whether these pre-formatting operations move into the data-population code or into the template itself.

### 3. PHASE_SUMMARY requires two independent `{{if}}` blocks, not one

> "`{{if .PhaseSummary}}...{{end}}` 替换 `If {{PHASE_SUMMARY}} is non-empty` 模式" (proposal line 156)

**Analysis**: In all templates that use PHASE_SUMMARY, the placeholder appears in two separate locations with different semantics:

1. **Label line** (e.g., line 4 of coding-cleanup.md): `{{PHASE_SUMMARY}}` -- becomes `PHASE_SUMMARY: <path>` or empty string. When empty, the line becomes just whitespace or an empty label (`PHASE_SUMMARY:`), which `isLabelWithEmptyValue()` removes entirely.

2. **Conditional instruction** (e.g., line 24 of coding-cleanup.md): `If \`{{PHASE_SUMMARY}}\` is non-empty, read that file...` -- when PHASE_SUMMARY is empty, this becomes `If \`\` is non-empty`, and the entire line is removed by `cleanTemplateOutput` rule 2.

These two locations are separated by 10-30 lines of other content in the template. They cannot be wrapped in a single `{{if .PhaseSummary}}...{{end}}` block. The migration requires two separate conditional blocks:
- `{{if .PhaseSummary}}{{.PhaseSummary}}{{end}}` for the label line
- `{{if .PhaseSummary}}If \`{{.PhaseSummary}}\` is non-empty, read that file...{{end}}` for the instruction line

The proposal's description implies a single block replacement, which would miss the label-line conditionalization. The 18 templates that contain PHASE_SUMMARY all have this two-location pattern.

### 4. `injectSurfaceFrontmatter` comment claims insertion behavior that code does not implement

> "injectSurfaceFrontmatter() 同时执行替换已有字段和插入缺失字段两种行为" (proposal line 181)

**Analysis**: The function comment in `add.go:269-271` states: "If the fields are absent, they are inserted before the closing `---`". However, the actual implementation (lines 272-280) only calls `strings.Replace` -- it replaces `surface-key: ""` with `surface-key: "value"` and `surface-type: ""` with `surface-type: "value"`. There is no insertion logic for missing fields.

The proposal's claim that this function has "两种行为" (two behaviors: replace and insert) is based on the function comment, not the implementation. This is a documentation-accuracy issue in the codebase, not a migration risk, but the proposal's migration plan is built on an incorrect understanding of the current behavior. In practice, both templates (coding.fix.md and coding.cleanup.md) already contain `surface-key: ""` and `surface-type: ""` lines, so the replace-only behavior has been sufficient. The proposal's instruction that "模板须包含 surface-key 和 surface-type 字段行" (templates must include these field lines) is correct and aligns with the actual code behavior.

**Verdict**: Low risk. The proposal's migration plan for `injectSurfaceFrontmatter` is correct despite the inaccurate claim about the function's behavior. The proposal should simply correct its description to say "replaces static empty values" rather than "simultaneously replaces and inserts."

### 5. `just compile {{SURFACE_KEY}}` trailing-whitespace handling needs explicit template design

> "`cleanTemplateOutput()` 仅保留空白行塌陷逻辑，移除全部四种条件删除逻辑" (proposal line 293)

**Analysis**: When `{{SURFACE_KEY}}` is empty, `just compile {{SURFACE_KEY}}` renders as `just compile ` with trailing whitespace. `cleanTemplateOutput` rule 3 strips this trailing whitespace. After migration to `text/template`, if the template contains `just compile {{.SurfaceKey}}` and SurfaceKey is empty, the output is `just compile ` -- `text/template` does not trim trailing whitespace.

This is used in 5 coding templates (coding-fix, coding-cleanup, coding-feature, coding-enhancement, coding-refactor) and the fix-record-missed template. Each has 3-4 `just <cmd> {{SURFACE_KEY}}` lines in a bash code block.

Solutions:
- Use `just compile{{if .SurfaceKey}} {{.SurfaceKey}}{{end}}` -- ensures no trailing space when empty
- Or use `just compile {{.SurfaceKey | trim}}` -- requires a custom trim function
- Or keep trailing-whitespace stripping in the minimal `cleanTemplateOutput` that the proposal retains

The proposal says `cleanTemplateOutput` will "仅保留空白行塌陷逻辑" (only keep blank-line collapse), which would drop the trailing-whitespace stripping. This means the templates must handle it declaratively. The `{{if}}` approach is cleanest.

### 6. Coverage three-state semantics correctly handled but underdocumented

> "CoverageStrategy ... 三态值由 Go 代码解析完整文本，模板侧仅判断 `{{if .CoverageStrategy}}`" (proposal line 169)

**Analysis**: Verified against `resolveCoverage()` in `prompt.go:337-368`:
- **Non-testable types** (doc, gate, test-gen, validation): `IsTestableType` returns false -> `resolveCoverage` is not called -> CoverageStrategy="" -> `{{if .CoverageStrategy}}` is false -> entire coverage block omitted
- **coding.cleanup/coding.refactor**: returns `("maintain", "Maintain existing coverage, no more than 2% decrease")` -> CoverageStrategy is non-empty -> block renders with "maintain existing coverage" text
- **Other coding types** (fix, feature, enhancement): returns percentage strategy -> CoverageStrategy is non-empty -> block renders with "Achieve N% test coverage"

The three-state design (empty / maintain / percentage) is correctly handled by using `{{if .CoverageStrategy}}` for presence and embedding the semantic difference in the CoverageTarget text. The proposal's approach is sound: the template just checks presence, and the Go code determines what text goes into the field. This is a good separation of concerns.

However, there is a subtlety in the cleanup/refactor templates: their coverage line contains both hardcoded text ("maintain existing coverage, no new tests required") AND the placeholder `{{COVERAGE_STRATEGY}}`. After migration, if the template hardcodes "Coverage strategy: maintain existing coverage" and also has `{{.CoverageStrategy}}`, the content would be redundant. The migration must decide: either the template contains only `{{.CoverageStrategy}} - {{.CoverageTarget}}` (like other coding templates), or the coverage field values encode the full text including the "no new tests" directive.

### 7. `{{SCOPE}}` two usage patterns correctly identified but migration examples missing

> "`{{SCOPE}}` 两种使用模式的迁移" (proposal lines 198-200)

**Analysis**: The proposal correctly identifies the two patterns:
- **Paragraph-level** (doc-consolidate, test-gen-contracts, test-gen-journeys, test-run): `{{SCOPE}}` as standalone block or `Scope: {{SCOPE}}` label -> entire section wrapped in `{{if .SurfaceKey}}...{{end}}`
- **Inline value** (other usage points): `{{SCOPE}}` as inline substitution -> replaced with `{{.SurfaceKey}}`

Verified against autogen templates:
- `test-gen-contracts.md:6` -- `{{SCOPE}}` standalone (paragraph-level)
- `test-gen-journeys.md:6` -- `{{SCOPE}}` standalone (paragraph-level)
- `doc-consolidate.md:4` -- `- Scope: {{SCOPE}}` (inline with label)
- `test-run.md:4` -- `- Scope: {{SCOPE}}` (inline with label)

The migration is well-defined. For paragraph-level, the `removeSection("Scope")` function removes everything from `## Scope` to the next `## ` heading. In `text/template`, this maps to `{{if .ScopeDisplay}}## Scope\n{{.ScopeDisplay}}\n{{end}}`. For inline, the `removeLineContaining("{{SCOPE}}")` call is replaced by `{{if .ScopeDisplay}}- Scope: {{.ScopeDisplay}}{{end}}`.

One edge case: `doc-consolidate.md` has `- Scope: {{SCOPE}}` (not under a `## Scope` heading), so it uses `removeLineContaining` semantics, not `removeSection`. The proposal's paragraph-level classification for doc-consolidate appears incorrect -- it should be inline.

### 8. Cross-proposal sequencing is correctly specified

> "`task-pipeline-precision` 先实施...本提案随后实施" (proposal line 95)

**Analysis**: The proposal explicitly states the sequencing: task-pipeline-precision introduces `<!-- IF NOT_LOW -->` markers first, then this proposal replaces them with `{{if ne .Complexity "low"}}...{{end}}`. The two proposals cannot be implemented in parallel because:
- If this proposal lands first, there are no `<!-- IF NOT_LOW -->` markers to replace, but task-pipeline-precision expects `cleanTemplateOutput` to process them
- If both land simultaneously, merge conflicts on the 4 coding templates

The proposal correctly identifies this dependency. The risk mitigation (each package migrated in independent commits, old functions preserved until golden-file tests pass) is sound.

### 9. Surface hard-failure scope is explicitly bounded

> "quality_gate.go 的 addSingleFixTask() 中硬性失败" vs "forge task add 命令路径保留软性行为" (proposal lines 222-224)

**Analysis**: Verified against `quality_gate.go:696-787`. The current `addSingleFixTask()` calls `inferSurface()` (line 715) which returns `("", "")` on failure (line 527-531, 557). The proposal changes this to return an error on inference failure. This is correctly scoped to the quality-gate path only.

The `forge task add` CLI path goes through `CreateTaskMarkdown()` in `add.go:240-267`, which calls `injectSurfaceFrontmatter()` only when surface values are non-empty (line 262). This is soft behavior -- no error on inference failure. The proposal preserves this by using empty string fallback in the template data struct for the CLI path.

Two paths, two behaviors, correctly documented. No issue.

## Improvement Suggestions

### Suggestion 1: Complete the `promptTemplateData` struct definition

**Addresses**: Key Risk #1 -- missing fields will cause `missingkey=error` failures

> The struct should include all 8 fields consumed by `renderTemplate`:
> ```go
> type promptTemplateData struct {
>     TaskID           string
>     TaskFile         string // absolute path to task markdown file
>     TaskCategory     string // task category for submit-task routing; empty = omit line
>     FeatureSlug      string // feature slug
>     PhaseSummary     string // "PHASE_SUMMARY: <path>" or empty
>     SurfaceKey       string // surface key; empty = omit label line
>     SurfaceType      string // surface type
>     TestTypeArg      string // " --type <type>" or empty; only used in test-gen-scripts
>     CoverageStrategy string // "percentage"/"maintain" or empty; empty = omit coverage block
>     CoverageTarget   string // human-readable coverage instruction text
>     Complexity       string // "low"/"medium"/"high"; defaults to "medium"
> }
> ```
> The proposal struct is missing `TaskFile`, `FeatureSlug`, `TestTypeArg`, and `CoverageTarget`. This will be caught immediately by the init-time validation gate (zero-value Execute with `missingkey=error`), but adding it to the proposal avoids implementor confusion.

### Suggestion 2: Complete the `autogenTemplateData` struct definition

**Addresses**: Key Risk #2 -- 5 of 9 data dimensions missing from struct

> The struct should mirror all fields consumed by `renderBody()` and the `BodyContext` struct:
> ```go
> type autogenTemplateData struct {
>     FeatureSlug        string // always substituted
>     Mode               string // empty = omit line containing {{.Mode}}
>     ScopeDisplay       string // pre-formatted "- item\n- item" string; empty = omit Scope section
>     SurfaceTypes       string // pre-formatted "- type" list; empty = "See .forge/config.yaml"
>     SurfaceKey         string // surface key for inline use
>     SurfaceType        string // surface type (also serves as TestType)
>     AcceptanceCriteria string // pre-formatted "- [ ] item" list; empty = default text
>     DocTaskCriteria    string // pre-formatted markdown sub-sections; empty = omit
> }
> ```
> Note: `TaskID` and `TaskType` from the proposal's struct are NOT consumed by `renderBody()` -- they are used in frontmatter generation (separate from template body rendering). The proposal should clarify whether the struct is for body rendering only or for the full `GenerateTestTaskMD()` pipeline.

### Suggestion 3: Document the two-block PHASE_SUMMARY pattern explicitly

**Addresses**: Key Risk #3 -- two separate `{{if}}` blocks needed

> For each template containing PHASE_SUMMARY, the migration needs two independent conditional blocks. Show a concrete example:
> ```markdown
> TASK_ID: {{.TaskID}}
> TASK_FILE: {{.TaskFile}}
> SURFACE_KEY: {{.SurfaceKey}}
> {{if .PhaseSummary}}{{.PhaseSummary}}{{end}}
> ...
> {{if .PhaseSummary}}If `{{.PhaseSummary}}` is non-empty, read that file for context from the previous phase.{{end}}
> ```
> The proposal should note that this pattern applies to all 18 templates that use PHASE_SUMMARY (all prompt templates except fix-record-missed). Alternatively, consider renaming the label-line field to avoid confusion -- e.g., `{{.PhaseSummaryLine}}` for the label and `{{.PhaseSummaryPath}}` for the path value used in the instruction.

### Suggestion 4: Add trailing-space prevention for `just` command lines in templates

**Addresses**: Key Risk #5 -- trailing whitespace on `just` lines after SurfaceKey substitution

> Use the conditional pattern in bash code blocks:
> ```markdown
> just compile{{if .SurfaceKey}} {{.SurfaceKey}}{{end}}
> just fmt{{if .SurfaceKey}} {{.SurfaceKey}}{{end}}
> just lint{{if .SurfaceKey}} {{.SurfaceKey}}{{end}}
> ```
> This eliminates trailing whitespace when SurfaceKey is empty without requiring post-processing. Apply to all 6 templates that have `just <cmd> {{SURFACE_KEY}}` lines: coding-fix, coding-cleanup, coding-feature, coding-enhancement, coding-refactor, and fix-record-missed.

### Suggestion 5: Correct the `injectSurfaceFrontmatter` description

**Addresses**: Key Risk #4 -- inaccurate description of current behavior

> The proposal states the function has "两种行为: 替换已有字段和插入缺失字段" but the actual implementation only performs `strings.Replace` on `surface-key: ""` and `surface-type: ""` literals. There is no insertion logic for absent fields. The proposal should correct this to: "替换 `surface-key: ""` 和 `surface-type: ""` 的字面值为推断值。" This correction does not change the migration plan -- both templates already contain the field lines.

### Suggestion 6: Clarify the `doc-consolidate.md` SCOPE migration pattern

**Addresses**: Key Risk #7 -- doc-consolidate is inline, not paragraph-level

> `doc-consolidate.md` line 4 has `- Scope: {{SCOPE}}` which is an inline value after a label, NOT a standalone paragraph block. The `renderBody()` function handles this via `strings.ReplaceAll` (not `removeSection`). The proposal classifies it as paragraph-level ("段落级"), but it should be classified as inline ("行内值") since there is no `## Scope` heading in this template. The correct migration is `{{if .ScopeDisplay}}- Scope: {{.ScopeDisplay}}{{end}}`, not wrapping a section in `{{if}}...{{end}}`.
