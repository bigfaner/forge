# Proposal Evaluation Report — Iteration Baseline

**Proposal**: Unified Template Engine for Prompt and Task Templates
**Date**: 2026-05-27
**Scorer**: CTO Adversarial Review (Phase 1 + 2 + 3)

---

## Phase 1: Reasoning Audit (Pre-Score Anchors)

### Argument Chain Trace

1. **Problem -> Solution**: Problem is real and well-documented. Three separate `strings.ReplaceAll` + post-processing pipelines exist. Solution (unified `text/template`) eliminates all three post-processing functions. However, solution scope underestimates the actual migration surface — `cleanTemplateOutput` already contains `<!-- IF NOT_LOW -->` marker handling (4 conditional behaviors, not 3), and `renderTemplate()` has a `TASK_CATEGORY` injection (line 136-137) that is entirely absent from the proposal.

2. **Solution -> Evidence**: Template count is wrong (claims 36, actual is 35: 21+2+12). This is a factual error in the migration inventory, which is the foundation for exhaustiveness claims.

3. **Evidence -> Success Criteria**: SC claims to remove "all conditional deletion logic" from `cleanTemplateOutput` but only lists 2 of 4 conditional behaviors. The `IF NOT_LOW` marker system and `just` command trailing whitespace are inconsistently covered between Scope and SC.

4. **Self-contradiction**: Scope lists 3 conditional behaviors to remove from `cleanTemplateOutput`; SC lists only 2. `BodyContext.Scope` is `[]string` but the proposal treats it as a simple rename. Surface hard-failure scope is ambiguous between quality gate and manual `forge task add` paths.

### SC Consistency Deep-Dive

**Cluster: pkg/prompt rendering**
- SC: `renderTemplate()` uses `text/template.Execute()` — clear
- SC: `cleanTemplateOutput()` only whitespace collapse — contradicts In Scope which lists 3 conditional behaviors to remove, but SC only names 2
- SC: 36 templates have `{{.Placeholder}}` format — count is wrong (35)
- Conflict: SC says "移除所有条件删除逻辑（空标签行、空 backtick 条件句）" but omits `IF NOT_LOW` and `just` whitespace. Scope includes `just` but neither includes `IF NOT_LOW`.

**Cluster: autogen templates**
- SC: `removeLineContaining()` and `removeSection()` removed — clear
- SC: No `{{SCOPE}}` residual — clear
- SC: `BodyContext` no `Scope` field — conflicts with proposal's own "重命名或替换" language (rename ≠ remove)
- Ambiguous: `{{SCOPE}}` has two usage patterns (block vs inline) requiring different template constructs, but SC treats it as a single check

**Cluster: Surface hard-failure**
- SC: `addSingleFixTask()` returns error on surface inference failure — clear
- Gap: No SC for `forge task add` behavior when surface inference fails

---

## Phase 2: Rubric Scoring (Dimension-by-Dimension)

### 1. Problem Definition (110 pts)

**Problem stated clearly (35/40)**: The core problem is unambiguous — three independent rendering engines with `strings.ReplaceAll` + post-processing hacks. The table in Evidence section makes it concrete. Minor ambiguity: the problem conflates "no conditional logic in templates" with "stale concepts" and "frontmatter inconsistency" — these are three distinct problems bundled as one.

**Evidence provided (38/40)**: Excellent evidence. Code-level detail with function names, file paths, template counts, and specific fragile assumptions (e.g., `isLabelWithEmptyValue` space limitation, `removeSection` heading prefix matching). The evidence table is precise and verifiable.

**Urgency justified (28/30)**: Clear dependency on `task-pipeline-precision` proposal (already approved). "Third layer of post-processing hack" is a compelling argument. Minor gap: does not quantify the cost of delay in terms of accumulated technical debt interest.

**Score: 101/110**

---

### 2. Solution Clarity (120 pts)

**Approach is concrete (30/40)**: The approach is mostly clear — migrate to `text/template`, use `{{if .Field}}...{{end}}` for conditionals, remove post-processing functions. However, key implementation details are missing:
- `TASK_CATEGORY` injection (line 136-137 of `prompt.go`) is not mentioned at all — this is a non-trivial migration because it injects after another field's substituted value
- The `IF NOT_LOW` marker system already in `cleanTemplateOutput` is not addressed
- Template data struct field designs are absent (only names given, no field definitions)

**User-facing behavior described (40/45)**: Good — explicitly states "no visible functional change" with one exception (surface hard-failure). The error message guidance is helpful. Deduction: the scope of the hard-failure behavior change is ambiguous — does it affect `forge task add` too?

**Technical direction clear (28/35)**: Direction is clear (`text/template` from stdlib, already used in codebase). But specific technical decisions are left open:
- `BodyContext.Scope` `[]string` rendering strategy (pre-formatted string vs `{{range}}`) — not decided
- `injectSurfaceFrontmatter` field-absent behavior — not addressed
- Init-time validation mechanism (`template.Parse` only validates syntax, not field references)

**Score: 98/120**

---

### 3. Industry Benchmarking (120 pts)

**Industry solutions referenced (30/40)**: References Helm, Hugo, goreleaser, buffalo. States "no competitive third-party alternatives" which is accurate for Go's `text/template` ecosystem. However, does not discuss alternative template approaches (e.g., Pongo2/templates, code generation, or functional options pattern) that could solve the same problem without template engines.

**At least 3 meaningful alternatives (25/30)**: Four alternatives listed including "do nothing". The "仅迁移 pkg/template" alternative is somewhat straw-man (2 templates only). The "标记注释 + 后处理" is a real alternative from task-pipeline-precision. One industry-validated solution present (text/template). Missing: functional approach (Go code with conditional builders instead of templates).

**Honest trade-off comparison (20/25)**: Trade-offs are reasonable but cherry-picked in favor of `text/template`. The "Cons" column for the selected approach says only "24 个模板文件需更新占位符语法" — this significantly understates the risk. The actual migration involves replacing 3 different post-processing systems with template-level conditionals, and the count is wrong (35 not 24, and not all are simple placeholder replacement).

**Chosen approach justified against benchmarks (22/25)**: Justification is sound — standard library, already used in codebase, declarative over imperative. The "why not do nothing" and "why not partial migration" arguments are convincing.

**Score: 97/120**

---

### 4. Requirements Completeness (110 pts)

**Scenario coverage (30/40)**: Four key scenarios identified (prompt conditional rendering, task conditional rendering, complexity branching, data model unification). Edge cases partially covered:
- Missing: What happens when a template uses `{{}}` syntax for non-placeholder purposes (e.g., code examples in templates)? The current code has a WARNING comment about this.
- Missing: The `coding-cleanup` template's special coverage semantics ("maintain existing coverage" vs "no instruction") is not addressed
- Missing: `injectSurfaceFrontmatter` field-absent branch (inserts fields that don't exist in template)

**Non-functional requirements (32/40)**: Three NFRs listed:
- Backward compatibility: well-defined (default to "medium" complexity)
- Byte-equivalent migration: well-defined (allow whitespace differences)
- Zero performance regression: claimed but not measured; `text/template.Execute()` is generally slower than `strings.ReplaceAll` for simple cases, though the difference is negligible for this use case
- Missing: Error handling NFR — what happens when `template.Parse()` fails at init time? What's the fallback behavior?

**Constraints & dependencies (25/30)**: Good list — `//go:embed`, callers of each function, `text/template` delimiter compatibility. Missing:
- The `task-pipeline-precision` dependency ordering is mentioned but not specified
- `forge task add` CLI surface behavior change is a constraint not listed

**Score: 87/110**

---

### 5. Solution Creativity (100 pts)

**Novelty over industry baseline (15/40)**: The proposal explicitly states "这不是技术创新" — honest, but means novelty is near zero. The differentiation from industry baseline is "unification" which is an engineering value, not a creative contribution. The proposal copies the standard Go `text/template` pattern with no adaptation or innovation.

**Cross-domain inspiration (10/35)**: No cross-domain inspiration. The solution is the most obvious choice for a Go project — standard library template engine. No borrowing from other domains.

**Simplicity of insight (20/25)**: The insight "replace 3 post-processing hacks with 1 declarative system" is clean and elegant. It's the right solution. Not a "why didn't I think of that" moment, but a solid engineering decision.

**Score: 45/100**

---

### 6. Feasibility (100 pts)

**Technical feasibility (32/40)**: `text/template` is already used in `pkg/task/data/` for record templates — the pattern is proven in this codebase. The migration path is mechanical for simple placeholder replacement. Deductions:
- `TASK_CATEGORY` injection requires a design decision not yet made
- `BodyContext.Scope` `[]string` rendering needs a strategy choice
- `IF NOT_LOW` markers need removal plan alongside `{{if}}` addition
- Init-time validation as described (parse-only) is insufficient for catching field typos

**Resource & timeline feasibility (25/30)**: 10-12 coding tasks in ~10 hours is realistic for mechanical migration. However, the task count is wrong (35 templates, not 36) and some tasks are more complex than "mechanical replacement" (e.g., `TASK_CATEGORY` injection, `IF NOT_LOW` marker coexistence). The estimate may be 20-30% low.

**Dependency readiness (28/30)**: No external dependencies — `text/template` is stdlib. Already proven in codebase. The cross-proposal dependency with `task-pipeline-precision` is the only risk, and it's acknowledged.

**Score: 85/100**

---

### 7. Scope Definition (80 pts)

**In-scope items are concrete (22/30)**: Most items are concrete deliverables (specific files, specific functions). Deductions:
- `BodyContext.Scope` field "重命名或替换为与当前概念对齐的字段名" — vague, no specific field name chosen
- "移除 `cleanTemplateOutput()` 中的条件删除逻辑（空标签行、空 backtick 条件句、`just` 命令尾部空白）" — lists 3 of 4 conditional behaviors; `IF NOT_LOW` missing
- Template data struct definitions: only names given, no fields

**Out-of-scope explicitly listed (22/25)**: Seven items explicitly listed as out of scope. Good — prevents scope creep. "Record 模板" is correctly excluded (already uses `text/template`). Deduction: `task-pipeline-precision`'s `IF NOT_LOW` markers are not listed as either in-scope or out-of-scope.

**Scope is bounded (18/25)**: Mostly bounded — 10-12 tasks, ~10 hours. However, the cross-proposal dependency with `task-pipeline-precision` creates unbounded scope: "共同实施" could mean anything from "sequential with clear handoff" to "merged implementation". The ambiguity adds risk.

**Score: 62/80**

---

### 8. Risk Assessment (90 pts)

**Risks identified (22/30)**: Seven risks identified. Good coverage of technical risks. Missing risks:
- `TASK_CATEGORY` injection loss (submit-task routing failure)
- `IF NOT_LOW` / `{{if}}` double-conditional coexistence
- Template data struct field design errors causing silent output changes
- `forge task add` breaking change from surface hard-failure scope expansion

**Likelihood + impact rated (22/30)**: Ratings are reasonable but tilted toward "M/H" pattern — most are Medium likelihood, High impact. This is an honest assessment. Deduction: the "占位符迁移遗漏" risk is rated M/H but should be rated H/H given the count is already wrong (36 vs 35), suggesting the inventory is incomplete.

**Mitigations are actionable (22/30)**: Mitigations are mostly actionable:
- "ValidatePromptTemplates() 扩展" — actionable but technically insufficient (parse-only, not field-check)
- "Golden-file 测试对比" — excellent, actionable
- "错误信息包含 `forge surfaces detect` 命令指引" — good UX
- Deduction: "本提案与 `task-pipeline-precision` 共同实施" is not actionable — no sequencing, no interface contract

**Score: 66/90**

---

### 9. Success Criteria (80 pts)

**Criteria are measurable and testable (22/30)**: Most SC items are binary checkable ("uses X", "no Y residual", "correctly renders"). Deductions:
- "`cleanTemplateOutput()` 仅保留空白行塌陷逻辑，移除所有条件删除逻辑（空标签行、空 backtick 条件句）" — lists 2 of 4 conditional behaviors, so the criterion is incomplete. An implementer following this SC would leave `IF NOT_LOW` and `just` whitespace handling in place.
- "36 个模板文件中无 `{{PLACEHOLDER}}` 格式" — count is wrong (35)
- "`forge prompt get-by-task-id` 输出与迁移前功能等价" — "功能等价" is somewhat vague; golden-file comparison is mentioned but the criterion itself is not the test

**Coverage is complete (18/25)**: SC covers all three packages, template syntax migration, surface hard-failure, and frontmatter. Gaps:
- No SC for `TASK_CATEGORY` injection preservation
- No SC for `IF NOT_LOW` marker removal
- No SC for `forge task add` behavior with empty surface (only quality gate is covered)
- No SC for template data struct field completeness
- No SC for `BodyContext.Scope` rendering correctness

**SC internal consistency (17/25)**:
- SC-4 (`cleanTemplateOutput` only whitespace) partially contradicts Scope's list of 3 conditional behaviors to remove
- SC-6 (36 templates) is factually wrong, creating inconsistency with actual migration
- SC-8 (`BodyContext` no `Scope` field) conflicts with proposal text saying "重命名或替换" (rename or replace, not remove)
- The `consistency_check_result: pass` at the bottom is questionable given these contradictions

**Score: 57/80**

---

### 10. Logical Consistency (90 pts)

**Solution addresses the stated problem (30/35)**: The unified `text/template` approach directly addresses the three-engine fragmentation problem. Conditional logic via `{{if}}` directly replaces post-processing hacks. Stale concept cleanup (`{{SCOPE}}` -> `{{.SurfaceKey}}`) directly addresses the migration debt. Deduction: The `TASK_CATEGORY` injection gap means the solution does not fully address `renderTemplate()`'s post-processing — one post-processing step would remain.

**Scope <-> Solution <-> Success Criteria aligned (18/30)**: Significant alignment gaps:
- Scope lists 3 `cleanTemplateOutput` conditional behaviors to remove; SC lists 2; code has 4
- Scope says `BodyContext.Scope` "重命名或替换"; SC says "no `Scope` field" (remove)
- Solution claims to "移除三套后处理函数"; `TASK_CATEGORY` injection is a fourth post-processing step not mentioned
- Scope includes "task-pipeline-precision complexity conditional branches" but no SC verifies complexity-conditional rendering correctness

**Requirements <-> Solution coherent (18/25)**: Requirements map mostly cleanly to solution. Gaps:
- "Complexity 条件分支" requirement maps to solution's `{{if}}` blocks, but the existing `IF NOT_LOW` marker system is an additional requirement not captured
- "Surface 硬性失败" requirement only maps to quality gate path; the `forge task add` path is unaddressed
- Template data model requirement (3 structs) is stated but not designed — fields are unknown

**Score: 66/90**

---

## Phase 3: Blindspot Hunt

### [blindspot-1] Template content collision with `{{}}` delimiters

Current `prompt.go` has an explicit WARNING comment (lines 98-103): "Template content must not contain bare placeholder strings like {{TASK_ID}}... If literal {{...}} is ever needed in a template, an escaping mechanism must be implemented first." After migration to `text/template`, any template that contains `{{` for non-template purposes (e.g., code examples, Go template examples in documentation) will cause parse errors. The proposal does not address this risk at all.

### [blindspot-2] `text/template` whitespace control behavior differs from current

`text/template`'s `{{if}}...{{end}}` blocks include the whitespace/newlines around them. The current post-processing approach of "remove line if value is empty" has different whitespace semantics than `{{if .Field}}...{{end}}`. The proposal claims "允许空白行差异" but does not analyze how `text/template`'s whitespace handling (which requires `{{-` and `-}}` for trimming) will affect output formatting across 35 templates.

### [blindspot-3] `embed.FS` template loading requires all templates to parse at init

Currently, templates are loaded as raw strings and substitution is done at runtime. After migration, `template.Parse()` must succeed for ALL templates at init time. If ANY template has a syntax error (even in a rarely-used template like `fix-record-missed.md`), the entire application fails to start. This is a reliability regression — currently, a malformed template would only fail when that specific template is used. The proposal does not discuss this trade-off.

### [blindspot-4] The `resolveCoverage` function's three-state output is not representable with a single boolean field

`resolveCoverage` returns three semantic states: ("", "") = no instruction, ("percentage", "text") = percentage target, ("maintain", "text") = maintain existing. The proposal's `{{if .CoverageStrategy}}...{{end}}` collapses all non-empty states into "render coverage block", but the template text differs between "percentage" and "maintain" strategies. A single `CoverageStrategy` field cannot differentiate between these states without additional template logic or multiple fields.

### [blindspot-5] No rollback plan

The proposal describes a major refactoring (3 packages, 35 templates, 3 data structs, function removal) with no rollback strategy. If the migration introduces subtle rendering differences that are only discovered in production, there is no documented path to revert. Given that `strings.ReplaceAll` and `text/template` have fundamentally different whitespace and escaping semantics, subtle differences are likely.

### [blindspot-6] The "consistency_check_result: pass" is self-reported and unverified

The proposal ends with a consistency check result claiming 21 pairs checked, 0 conflicts found. But the analysis above found multiple contradictions (SC vs Scope, template count, `cleanTemplateOutput` coverage). This self-assessment appears to be unreliable and may give false confidence.

---

## Score Summary

| Dimension | Score | Max |
|-----------|-------|-----|
| Problem Definition | 101 | 110 |
| Solution Clarity | 98 | 120 |
| Industry Benchmarking | 97 | 120 |
| Requirements Completeness | 87 | 110 |
| Solution Creativity | 45 | 100 |
| Feasibility | 85 | 100 |
| Scope Definition | 62 | 80 |
| Risk Assessment | 66 | 90 |
| Success Criteria | 57 | 80 |
| Logical Consistency | 66 | 90 |
| **Total** | **764** | **1000** |

---

## Attack Summary

### Critical Attacks (require revision before approval)

1. **[Solution Clarity]** `TASK_CATEGORY` injection in `renderTemplate()` (line 136-137) is entirely absent from the proposal. This post-processing step injects a category line after `TASK_FILE` using `strings.Replace` on the already-substituted content. In `text/template`, this cannot be replicated with a simple `{{.TaskCategory}}` field because the injection point is relative to another field's value. If this is missed during migration, `submit-task` skill routing will break silently. — Quote: "为每个模板添加 `{{if}}` 条件块，替换当前靠后处理删除的段落" — Must add TASK_CATEGORY to the migration plan.

2. **[Scope/SC Consistency]** `cleanTemplateOutput()` has 4 conditional behaviors, not 3. The `<!-- IF NOT_LOW -->...<!-- END_IF -->` marker system (already active in 4 coding templates) is neither listed in Scope nor covered by SC. — Quote from Scope: "移除 `cleanTemplateOutput()` 中的条件删除逻辑（空标签行、空 backtick 条件句、`just` 命令尾部空白）" — missing `IF NOT_LOW`. Quote from SC: "`cleanTemplateOutput()` 仅保留空白行塌陷逻辑，移除所有条件删除逻辑（空标签行、空 backtick 条件句）" — missing `IF NOT_LOW` and `just` whitespace. Must list all 4 behaviors and add SC entries for each.

3. **[Success Criteria]** Template count is factually wrong. Claims "36 个模板" but actual count is 35 (21 prompt + 2 task creation + 12 autogen). — Quote: "36 个模板文件中无 `{{PLACEHOLDER}}` 格式（全部为 `{{.Placeholder}}` 格式）" — Must correct the count and provide a per-file checklist to ensure migration completeness.

4. **[Risk Assessment]** Init-time validation as described is insufficient. `template.Parse()` only validates syntax, not field references. A typo like `{{.TaskCatgory}}` would parse successfully but output empty string at execute time. — Quote: "启动时 `ValidatePromptTemplates()` 和 `ValidateAutogenTemplates()` 扩展为对所有模板执行 `template.Parse()` + 检查所有字段引用，编译期而非运行期捕获遗漏" — Must specify `missingkey=error` option or equivalent field-check mechanism.

5. **[Scope Definition]** Cross-proposal dependency with `task-pipeline-precision` is undefined. "共同实施" is ambiguous — implementation order affects template file content. — Quote: "本提案与 `task-pipeline-precision` 共同实施，complexity 条件直接用 `{{if}}` 实现" — Must define sequencing: which proposal lands first, and what the interface contract between them is.

6. **[Logical Consistency]** Surface hard-failure scope is ambiguous between quality gate path and `forge task add` path. — Quote from Scope: "`quality_gate.go` 的 `addSingleFixTask()` 中 `inferSurface()` 失败时返回错误（硬性失败）" — but `injectSurfaceFrontmatter` (which also handles surface injection) is called from `CreateTaskMarkdown` which is used by both paths. Must explicitly state whether `forge task add` also hard-fails on surface inference failure.

### Moderate Attacks (should improve before implementation)

7. **[Solution Clarity]** `BodyContext.Scope` is `[]string`, not a simple field to rename. The rendering strategy (pre-formatted string vs `{{range}}`) is undecided. — Quote: "`BodyContext.Scope` 字段重命名或替换为与当前概念对齐的字段名" — Must choose and specify the rendering strategy.

8. **[Solution Clarity]** Template data struct designs are absent. Only names are given (`promptTemplateData`, `taskTemplateData`, `autogenTemplateData`) with no field definitions. This is the core of the migration — without field designs, implementers cannot determine which template patterns are representable. — Must provide at least the key fields for each struct, especially for multi-state fields like coverage.

9. **[Requirements Completeness]** `{{SCOPE}}` has two semantically different usage patterns in autogen templates: block-level (handled by `removeSection`) and inline (handled by `strings.ReplaceAll`). These require different `text/template` constructs. — Quote: "`{{SCOPE}}` 替换为 `{{.Scope}}` 或等效的 surface 概念字段" — Must specify different migration strategies for the two patterns.

10. **[Scope Definition]** `injectSurfaceFrontmatter` has two behaviors: replacing existing empty values AND inserting fields when absent. The proposal only considers the replacement case. — Quote: "移除 `pkg/task/add.go` 的 `injectSurfaceFrontmatter()`——surface 值直接由模板渲染" — Must address what happens when templates lack surface frontmatter fields entirely.

### [blindspot] Attacks

11. **[blindspot-1]** No discussion of template content containing literal `{{` — the current code has an explicit WARNING about this. After migration, any template with `{{` for non-template purposes (code examples, documentation) will cause parse errors.

12. **[blindspot-2]** `text/template` whitespace control semantics differ from current line-removal approach. `{{if}}...{{end}}` blocks include surrounding whitespace, requiring `{{-` / `-}}` for trimming. This affects output formatting across all 35 templates but is not analyzed.

13. **[blindspot-3]** `embed.FS` + `template.Parse()` at init means ANY template syntax error prevents application startup. This is a reliability regression from the current lazy-loading behavior.

14. **[blindspot-4]** `resolveCoverage` has three output states that cannot be represented by a single boolean `{{if .CoverageStrategy}}`. The "maintain" and "percentage" strategies need different template text, requiring either multiple fields or conditional template logic.

15. **[blindspot-5]** No rollback plan for a major refactoring affecting 3 packages and 35 templates. Subtle rendering differences may only surface in production.

16. **[blindspot-6]** The self-reported `consistency_check_result: pass` is unreliable — multiple contradictions were found between SC, Scope, and actual code.
