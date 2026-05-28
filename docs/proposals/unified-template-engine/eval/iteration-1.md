# Proposal Evaluation Report — Iteration 1 (Re-Evaluation)

**Proposal**: Unified Template Engine for Prompt and Task Templates
**Date**: 2026-05-28
**Scorer**: CTO Adversarial Review (Independent Re-Evaluation)
**Previous iteration**: 0 (score: 764/1000)

---

## Phase 1: Reasoning Audit (Pre-Score Anchors)

### Argument Chain Trace

1. **Problem → Solution**: Three sub-problems stated: (a) three independent rendering engines with post-processing hacks, (b) stale `scope` concepts, (c) frontmatter inconsistency. Solution addresses all three via unified `text/template` + concept cleanup + frontmatter standardization. The chain is sound. However, the solution introduces significant scope beyond what the problem demands: directory restructuring, a metadata-frontmatter validation system with reflection-based cross-checking, and `forge init` integration. These are good engineering decisions but they expand the proposal's surface area beyond the stated problem boundaries.

2. **Solution → Evidence**: Template counts (21+2+12=35, 6 record) verified by expert review and correct throughout. The `promptTemplateData` struct now has all 11 fields (fixing iteration-0 gap). The `autogenTemplateData` struct now has all 10 fields. Post-processing function behaviors documented with specific fragility patterns. Evidence is strong.

3. **Evidence → Success Criteria**: 31 SC entries cover most Scope items. Gaps identified:
   - `forge init` surface configuration (Scope line 233) has no SC entry
   - PHASE_SUMMARY dual-block rendering (Scope line 156) has a single SC entry (line 306) that is ambiguous about whether it covers both blocks
   - No SC for the `CoverageStrategy` source-of-truth design decision (Go code vs template)

4. **Self-contradiction check**: The solution claims to eliminate all post-processing but retains `cleanTemplateOutput()` for "blank-line collapse" (line 154). `text/template` has built-in whitespace trimming (`{{-` and `-}}`), so the retained function may be unnecessary. The proposal does not discuss whether the retained function is a temporary bridge or a permanent design choice. This is a minor inconsistency: the stated goal is "移除三套后处理函数" (line 53) but the execution keeps one alive, albeit stripped down.

### SC Consistency Deep-Dive

**Cluster: pkg/prompt rendering**
- SC: `renderTemplate()` uses `text/template.Execute()` — satisfiable
- SC: `cleanTemplateOutput()` stripped to whitespace collapse, lists all 4 conditional behaviors — satisfiable
- SC: 35 templates `{{.Placeholder}}` format — satisfiable
- SC: `TASK_CATEGORY` migration — satisfiable
- SC: `<!-- IF NOT_LOW -->` markers removed — satisfiable given sequencing constraint
- Tension: SC "`{{if .PhaseSummary}}` 条件块正确渲染：有值时注入段落，无值时段落消失" (line 306) is ambiguous. The Scope (line 156) explicitly states two independent `{{if}}` blocks are needed. The SC says "段落" (paragraph) in singular. A strict reading might verify only one block works.

**Cluster: autogen templates**
- SC: `removeLineContaining()` and `removeSection()` removed — satisfiable
- SC: No `{{SCOPE}}` residual — satisfiable
- SC: `BodyContext.Scope` replaced by `ScopeDisplay string` — satisfiable
- Note: `{{SCOPE}}` migration classification at line 202 lists `test-run` as "段落级" but expert review found `test-run.md:4` has `- Scope: {{SCOPE}}` (inline with label), matching the `doc-consolidate` pattern, not a standalone block. This classification error could cause incorrect template migration (wrapping a section in `{{if}}` when only a line-level conditional is needed). The Scope text does say "段落级" includes templates where SCOPE is "独占一行或作为段落标题的一部分" — but `test-run` has neither pattern.

**Cluster: Surface behavior**
- SC: `addSingleFixTask()` returns error on surface failure — satisfiable
- SC: `forge task add` keeps soft behavior — satisfiable
- Gap: `forge init` surface configuration step (Scope line 233) has no SC entry. If this is a deliverable, the SC set is incomplete.

**Cluster: Frontmatter & validation**
- SC: All 41 templates have metadata frontmatter — satisfiable
- SC: Record template structure correct — satisfiable
- SC: Template loader strips metadata before Parse — satisfiable
- SC: `ValidateTemplates()` uses `variables` field cross-check — satisfiable
- Note: This cluster represents ~20% of the Scope but is not motivated by any specific problem statement. The frontmatter inconsistency problem (line 30) only mentions missing `SURFACE_KEY` headers and inconsistent presence of YAML frontmatter — it does not motivate a full metadata-frontmatter + reflection-validation system.

### Pre-Score Anchors

- **Strong**: Problem evidence, template data structs (complete), urgency argument, rollback strategy (per-package independent commits), sequencing constraint
- **Moderate concern**: Scope creep beyond stated problem, `test-run` SCOPE classification error, retained `cleanTemplateOutput()` not discussed as a design trade-off
- **Weakness**: `forge init` deliverable without SC, PHASE_SUMMARY dual-block ambiguity in SC, CoverageStrategy design trade-off not acknowledged, frontmatter validation system unmotivated by problem section

---

## Phase 2: Rubric Scoring (Dimension-by-Dimension)

### 1. Problem Definition (110 pts)

**Problem stated clearly (37/40)**: Core problem is unambiguous — three independent rendering engines with fragile post-processing hacks. The three sub-problems (conditional logic, stale concepts, frontmatter inconsistency) are clearly delineated. The evidence tables with function names, file paths, and specific fragility patterns are precise. Minor deduction: the "frontmatter 不一致" sub-problem is stated at line 30 but its connection to the proposed metadata-frontmatter validation system (Scope lines 236-249) is tenuous — the problem says templates have inconsistent frontmatter, but the solution creates an entirely new metadata layer rather than simply fixing the inconsistency.

**Evidence provided (38/40)**: Excellent evidence. Three evidence tables with code-level detail: rendering engine inventory, stale concept inventory, post-processing fragility inventory. Template counts are accurate (verified by expert review). The `cleanTemplateOutput()` four-rule breakdown is precise and verifiable. The `injectSurfaceFrontmatter()` description has been corrected to "替换 `surface-key: ""` 字面值" (line 20) from the earlier dual-behavior claim. Minor deduction: the "后处理函数脆弱假设" section (lines 36-39) lists three specific fragilities but the Scope addresses four conditional behaviors in `cleanTemplateOutput()` — the `<!-- IF NOT_LOW -->` block pattern is not listed as a fragility (it's a feature, not a fragility), creating a slight mismatch between problem framing and solution scope.

**Urgency justified (28/30)**: Clear dependency on approved `task-pipeline-precision` proposal. "Third layer of hack" argument is compelling and concrete — complexity branching would add a third conditional dimension on top of phase-summary and coverage. The sequencing constraint (line 95) adds weight: delay creates a dependency bottleneck. Minor deduction: urgency focuses on the conditional-logic problem but the proposal also includes significant scope (directory restructuring, frontmatter system) that is not urgent.

**Score: 103/110**

---

### 2. Solution Clarity (120 pts)

**Approach is concrete (38/40)**: The approach is highly concrete. Key strengths:
- All three template data structs have complete field definitions with comments explaining zero-value behavior (lines 166-179, 188-194, 209-220)
- All four `cleanTemplateOutput()` conditional behaviors are individually listed with replacement `{{if}}` patterns (lines 155-161)
- `{{SCOPE}}` two-pattern migration strategy is documented with specific templates named (lines 201-203)
- `BodyContext.Scope` → `ScopeDisplay string` replacement is explicit (line 205)
- `TASK_CATEGORY` injection migration has its own sub-section (lines 224-227)

Minor deduction: line 226 says "在 `promptTemplateData` 中增加 `TaskCategory string` 字段" but `promptTemplateData` already has `TaskCategory` listed at line 169. This is a drafting artifact from incomplete revision cleanup — it suggests the author added the field to the struct but forgot to remove the separate "add this field" instruction.

**User-facing behavior described (41/45)**: "No visible functional change" with one exception (surface hard-failure). The two-path distinction is explicit:
- Quality gate path: hard failure (line 231)
- CLI path: soft behavior (line 232)

The error message guidance is mentioned: "引导用户运行 `forge surfaces detect` 配置 surfaces" (line 68). Deductions:
- The actual error message text is not shown — "引导用户运行" describes intent but not the message
- The `forge init` integration (line 233) is a user-facing behavior change (new surface configuration step during project initialization) that is not described in the User-Facing Behavior section (lines 66-68)
- The record template structural change (metadata frontmatter → body) is not described in terms of developer experience — what does a developer see when they open a record template file?

**Technical direction clear (33/35)**: `text/template` from stdlib, already proven in `pkg/task/data/` typed-task-records. Init-time validation with `missingkey=error` + zero-value Execute. Pre-formatted string fields to avoid template-side logic. All clear. Minor deduction: the retained `cleanTemplateOutput()` function (blank-line collapse only) is not discussed as a design decision. `text/template` supports `{{-` / `-}}` for whitespace trimming — the proposal should explain why a Go function is preferred over template-side whitespace control.

**Score: 112/120**

---

### 3. Industry Benchmarking (120 pts)

**Industry solutions referenced (30/40)**: References Helm, Hugo, goreleaser, buffalo. States "无竞争力的第三方替代品" — accurate for Go's `text/template` ecosystem. The Innovation Highlights section is honest: "这不是技术创新" (line 62). However, the benchmarking is shallow:
- No URLs, version numbers, or specific feature comparisons with these tools
- No discussion of how these tools handle the specific challenges Forge faces (conditional sections, frontmatter management, validation)
- No discussion of alternative approaches at the architectural level (e.g., code generation, builder pattern, DSL)

**At least 3 meaningful alternatives (25/30)**: Four alternatives listed including "do nothing". "标记注释 + 后处理" is a real alternative from `task-pipeline-precision`. "仅迁移 pkg/template" is a straw-man — it has 2 templates and the proposal itself acknowledges it's rejected for non-unification, making it a token entry rather than a seriously considered option. The comparison table is concise but adequate.

**Honest trade-off comparison (22/25)**: Trade-offs are reasonable. The "Cons" for the selected approach says "35 个模板文件需更新占位符语法" — this understates the scope. The actual migration involves: (1) placeholder syntax updates, (2) conditional block insertion, (3) `{{SCOPE}}` two-pattern migration, (4) post-processing function removal, (5) frontmatter restructuring, (6) directory reorganization, (7) `//go:embed` path updates, (8) validation system implementation. The "Cons" column should at minimum mention the cross-proposal dependency and behavioral change.

**Chosen approach justified (22/25)**: Justification is sound — standard library, already in codebase, declarative over imperative. "Already in codebase" (typed-task-records) is the strongest argument. Minor deduction: the comparison table's "Source" column for the selected approach says "Go 标准库" — could have cited the specific existing codebase usage that validates the pattern.

**Score: 99/120**

---

### 4. Requirements Completeness (110 pts)

**Scenario coverage (33/40)**: Six key scenarios identified. Edge cases now well-covered:
- Coverage three-state semantics (line 172)
- `{{SCOPE}}` two-pattern migration (lines 201-203)
- `just compile` trailing whitespace (line 159)
- Surface inference two-path behavior (lines 231-232)

Gaps:
- No scenario for `embed.FS` + `template.Parse()` init-time failure — what happens to the application when a template fails to parse? The current lazy approach fails only when a specific template is used.
- No scenario for `ScopeDisplay` pre-formatting edge cases (empty slice, single-item slice, multi-item slice with special characters)
- No scenario for the metadata frontmatter stripping — what if a template's body content starts with `---`? The loader needs to distinguish between metadata frontmatter and body content.

**Non-functional requirements (34/40)**: Three NFRs stated:
- Backward compatibility: "现有 index.json 中无 complexity 字段的任务默认为 `medium`" — specific and testable
- Byte-equivalent migration: "允许空白行差异但要求内容等价" — clear criterion
- Zero performance regression: "性能与 `strings.ReplaceAll` 可比" — reasonable claim for template rendering

Deductions:
- No NFR for reliability/startup behavior. After migration, `template.Parse()` failure prevents application startup. This is a reliability trade-off (fail-fast vs fail-late) that is mentioned at line 86 as an intentional "权衡" but not as a formal NFR with acceptance criteria.
- No NFR for the metadata frontmatter validation system's performance impact — reflection-based cross-checking at startup has a cost that should be bounded.

**Constraints & dependencies (28/30)**: Good list. Key constraints:
- `//go:embed` + init-time parsing (line 90)
- Cross-proposal sequencing (line 95)
- `missingkey=error` + zero-value Execute (line 86)
- `pkg/template` dual consumer paths (line 91)

Deduction: the `forge init` surface configuration dependency (line 233) is mentioned in Scope but not in the Constraints section — if it's a prerequisite for the hard-failure behavior to be acceptable, it should be listed as a constraint.

**Score: 95/110**

---

### 5. Solution Creativity (100 pts)

**Novelty over industry baseline (15/40)**: The proposal explicitly states "这不是技术创新" — honest. This is standard `text/template` usage applied to a consolidation problem. No differentiation from the Helm/Hugo/goreleaser pattern. The "innovation" is unification of three engines, which is engineering value, not creative contribution.

**Cross-domain inspiration (8/35)**: No cross-domain inspiration. No borrowing from other domains (e.g., Rails template rendering, Jinja2 conditionals, TypeScript template literal types). The solution is the most obvious choice for a Go project.

**Simplicity of insight (22/25)**: "Replace 3 post-processing hacks with 1 declarative system" is clean and elegant. The `just compile{{if .SurfaceKey}} {{.SurfaceKey}}{{end}}` pattern (line 159) is a nice touch — it eliminates trailing whitespace declaratively without post-processing. The `ScopeDisplay` pre-formatted string approach avoids template-side `{{range}}` loops, keeping templates simple. Not a "why didn't I think of that" moment, but a solid, well-scoped engineering decision.

**Score: 45/100**

---

### 6. Feasibility (100 pts)

**Technical feasibility (36/40)**: Pattern already proven in `pkg/task/data/` (typed-task-records). Migration path is mechanical for placeholder replacement. Template data structs are now complete. The `ScopeDisplay` pre-formatting approach is straightforward. Deductions:
- The metadata frontmatter system (strip frontmatter, parse `variables`, reflect over struct fields, cross-validate) adds non-trivial complexity. The proposal treats it as a simple addition (~2 hours, line 128) but it involves YAML parsing, Go reflection, and a new validation framework. This may be underestimated.
- The `CoverageStrategy` three-state semantics (empty/maintain/percentage) relies on Go code to resolve the correct value. The template cannot distinguish states — all non-empty values render the same block. If future requirements need template-level state differentiation, the field must be redesigned.

**Resource & timeline feasibility (24/30)**: 15 coding tasks, ~15 hours (line 132). The estimate is broken down per sub-task, which is good. However:
- "Metadata frontmatter: 41 个模板文件添加 metadata frontmatter, record 模板输出 frontmatter 移入 body, 模板加载器解析 metadata：约 2 小时" — this includes writing a metadata parser, implementing reflection-based validation, and restructuring 6 record templates. 2 hours seems optimistic.
- The `forge init` surface configuration integration (Scope line 233) is not included in the timeline estimate.
- Total estimate may be 20-30% low given the scope breadth.

**Dependency readiness (28/30)**: No external dependencies. `text/template` is stdlib. Cross-proposal dependency is sequenced (line 95). The only risk is `task-pipeline-precision` delay, which is outside this proposal's control but acknowledged.

**Score: 88/100**

---

### 7. Scope Definition (80 pts)

**In-scope items are concrete (26/30)**: Most items are concrete deliverables. The template data structs with field definitions are excellent. The conditional-behavior replacement patterns are specific. Deductions:
- The `forge init` surface configuration step (line 233) is vaguely specified: "在 `forge init` 中增加 surface 配置步骤（`forge surfaces detect` 集成）" — what does "集成" mean? A new sub-command? A prompt during init? An automatic detection step?
- The `taskTemplateData` struct (line 188-194) has a `Scope` field with comment "任务作用域描述" — this is the only field that uses the deprecated "Scope" terminology while the rest of the proposal systematically replaces it. Is this intentional (the task-creation context still uses scope as a user-facing concept) or an oversight?

**Out-of-scope explicitly listed (24/25)**: Seven items explicitly listed. Good scope containment. Record templates correctly excluded (already use `text/template`). `mixed.just` beyond scope correctly bounded. Deduction: `ARCHITECTURE.md` updates are excluded (line 275) but the directory restructuring (pkg/template deletion, data/ split) will make `ARCHITECTURE.md` inaccurate — excluding documentation updates from a restructuring proposal creates a maintenance gap.

**Scope is bounded (22/25)**: 15 tasks, ~15 hours, quick mode. The sequencing constraint is clear. However:
- The scope includes directory restructuring, a new metadata-frontmatter system, behavioral changes (surface hard-failure), `forge init` integration, skill/command/agent concept alignment, and template engine migration — this is a lot for a "quick mode" proposal. The breadth increases coordination risk.
- The surface hard-failure behavioral change (line 231) is a behavior change, not just a refactoring. The proposal acknowledges this ("这是行为变更（当前行为为静默空值），而非纯重构") but the scope classification doesn't reflect the additional testing burden for a behavioral change.

**Score: 72/80**

---

### 8. Risk Assessment (90 pts)

**Risks identified (26/30)**: Eleven risks identified. Good coverage of migration risks, cross-proposal synchronization, and nil-pointer safety. Gaps:
- No risk for the metadata frontmatter system — YAML parsing, reflection-based validation, and frontmatter stripping are new code paths that could have edge cases (malformed frontmatter, missing `variables` field, body content starting with `---`).
- No risk for `text/template` whitespace behavior differences — `text/template` inserts newlines differently from `strings.ReplaceAll`, and the "允许空白行差异" criterion may not capture all whitespace differences (e.g., trailing newlines at end of file).
- No risk for the `ScopeDisplay` field replacing `BodyContext.Scope []string` — the type change from slice to pre-formatted string means all callers of `BodyContext.Scope` must be updated, not just `renderBody()`.

**Likelihood + impact rated (25/30)**: Ratings are generally reasonable. Specific concern: "BodyContext.Scope 重命名影响 `BuildIndex()` 调用链" is rated M/M (line 288), but the proposal also says "Scope 仅在 `autogen.go` 内部消费" — if true, the impact is M/L (BuildIndex is a core function), not M/M. The rating appears inconsistent with the scope of the change.

**Mitigations are actionable (25/30)**: Most mitigations are concrete and actionable:
- `missingkey=error` + zero-value Execute — technically correct
- Per-package independent commits with golden-file testing — excellent rollback strategy
- Two-step separation (engine migration first, directory restructuring second) — sound risk reduction
- `forge init` surface configuration — actionable but not scoped as deliverable

Deduction: "错误信息包含 `forge surfaces detect` 命令指引" (line 284) — the error message content is not specified. A mitigation should include the actual error text to verify it's helpful. Also, the `forge init` integration is mentioned as a mitigation for surface hard-failure but is not in Scope — if it's required for the mitigation to work, it must be a deliverable.

**Score: 76/90**

---

### 9. Success Criteria (80 pts)

**Criteria are measurable and testable (26/30)**: Most SC items are binary checkable. "Does this function use `text/template.Execute()`?" — yes/no. "Are `removeLineContaining()` and `removeSection()` removed?" — yes/no. Deductions:
- "`forge prompt get-by-task-id` 输出与迁移前功能等价（golden-file 对比，允许空白行差异）" (line 312) — "功能等价" is the criterion, "golden-file 对比" is the test method. But "功能等价" is ambiguous: does it mean identical markdown structure? Identical semantic content? Identical rendered output? The parenthetical hints at the method but not the threshold.
- "`{{if .PhaseSummary}}` 条件块正确渲染：有值时注入段落，无值时段落消失" (line 306) — "段落" is singular, but Scope (line 156) explicitly states two independent blocks are needed. A verifier could check that PhaseSummary is conditionally rendered but miss that both locations are covered.
- "`ValidateTemplates()` 使用 `variables` 字段与模板数据 struct 做交叉校验" (line 330) — "交叉校验" is vague. Does it verify that every declared variable exists in the struct? That every struct field is declared? That types match? The criterion should specify the validation direction and completeness.

**Coverage is complete (22/25)**: 31 SC entries cover most Scope items. Gaps:
- No SC for `forge init` surface configuration integration (Scope line 233)
- No SC for the `CoverageStrategy` source-of-truth design (Go code resolves text, template renders it) — the SC (line 308) verifies the three-state rendering but not that the coverage instruction text is correct for each state
- No SC for the `autogen body 模板的 metadata frontmatter 替代原有的 `<!-- body-only -->` 注释方案" (Scope line 248) — if `<!-- body-only -->` comments are being replaced, there should be an SC verifying they are removed
- No SC for the `taskTemplateData.Scope` field — the proposal replaces "scope" terminology everywhere else but this field retains it (line 193). Is this intentional? No SC verifies the field naming.

**SC internal consistency (22/25)**: The SC set is largely internally consistent. The `consistency_check_result: pass` (lines 337-342) is self-reported but the improved proposal makes it more credible. Remaining tensions:
- SC "35 个模板文件中无 `{{PLACEHOLDER}}` 格式" (line 303) — but 6 record templates also need migration (they already use `text/template` but are included in the 41-template metadata frontmatter count at line 324). The SC says "35" for syntax migration but "41" for frontmatter — this is correct because record templates don't need syntax migration, but it could confuse a verifier.
- SC "`BodyContext` struct 中无 `Scope` 字段（已替换为当前概念对齐的字段名）" (line 305) — the parenthetical "当前概念对齐的字段名" is vague. Scope line 205 specifies `ScopeDisplay string` — the SC should name the replacement field directly.

**Score: 70/80**

---

### 10. Logical Consistency (90 pts)

**Solution addresses the stated problem (33/35)**: Unified `text/template` directly addresses three-engine fragmentation. Conditional `{{if}}` replaces all post-processing hacks. Stale concept cleanup is explicit. The solution is well-targeted at the core problem. Deduction: the "frontmatter 不一致" sub-problem (line 30) identifies two specific issues: (a) 6 record templates have frontmatter, 12 don't; (b) 4 doc templates lack `SURFACE_KEY` headers. The solution addresses (b) with a Scope entry (line 162) and SC entry (line 315). For (a), the solution creates an entirely new metadata-frontmatter system (lines 236-249) rather than simply adding consistent frontmatter to the 12 missing templates. This is a leap from "inconsistency is a problem" to "build a validation framework" — the response is disproportionate to the problem.

**Scope ↔ Solution ↔ SC aligned (25/30)**: Significant improvement from iteration 0. Most Scope items have corresponding SC entries. Remaining misalignments:
- Scope line 233 (`forge init` surface configuration) has no SC entry
- Scope line 248 (metadata frontmatter replaces `<!-- body-only -->` comments) has no SC entry verifying `<!-- body-only -->` removal
- Scope line 156 (PHASE_SUMMARY dual-block pattern) is not precisely reflected in SC line 306 (singular "段落")
- Scope line 202 (`test-run` classified as "段落级") contradicts expert review finding that `test-run.md:4` has `- Scope: {{SCOPE}}` (inline pattern). If the migration applies paragraph-level wrapping to an inline pattern, it would produce incorrect output.

**Requirements ↔ Solution coherent (21/25)**: Requirements map cleanly to solution for the core template-engine migration. Deductions:
- The "目录结构分散" problem (line 32) motivates directory restructuring. The restructuring is coherent but adds significant scope without a dedicated requirements analysis — no scenarios or NFRs for the restructuring itself.
- The `taskTemplateData.Scope` field (line 193) retains "Scope" terminology while the rest of the proposal systematically replaces it. This is either intentional (the task-creation context uses scope as a user-facing concept distinct from SurfaceKey) or an oversight. The proposal does not explain the retention.
- The "Non-Functional Requirements" section mentions "字节等价迁移" (line 84) but the Scope includes behavioral changes (surface hard-failure, `forge init` integration) that are not "equivalent" to current behavior. The NFR should clarify which paths must be equivalent and which are allowed to change.

**Score: 79/90**

---

## Phase 3: Blindspot Hunt

### [blindspot-1] `test-run` SCOPE classification error persists in Scope text

The proposal at line 202 classifies `test-run` as "段落级" (paragraph-level), implying its `{{SCOPE}}` usage is a standalone block or section heading. Expert review found `test-run.md:4` has `- Scope: {{SCOPE}}` — an inline value with a label, identical to `doc-consolidate.md`'s pattern. This was identified in iteration 0 but not corrected. If the migration applies paragraph-level wrapping (`{{if .SurfaceKey}}## Scope\n...{{end}}`) to an inline pattern, the output structure would be wrong.

Quote: "段落级（test-gen-contracts, test-gen-journeys, test-run）：`{{SCOPE}}` 独占一行或作为段落标题的一部分 → 整段用 `{{if .SurfaceKey}}...{{end}}` 包裹"

What must improve: Reclassify `test-run` as "行内值" alongside `doc-consolidate`, or verify the actual template content and correct the classification.

### [blindspot-2] `taskTemplateData.Scope` field contradicts the proposal's own terminology migration

The proposal systematically replaces "scope" with "SurfaceKey"/"SurfaceType"/"ScopeDisplay" across all three template data structs and the `BodyContext` struct. Yet `taskTemplateData` at line 193 retains a `Scope string` field with comment "任务作用域描述". This is the only struct field in the entire proposal that uses the deprecated "scope" name. The proposal does not explain why this field is exempt from the terminology migration.

Quote: "Scope string // 任务作用域描述" (line 193)

What must improve: Either rename this field to align with the migration (e.g., `ScopeDescription`), or explicitly justify the retention as representing a different concept (user-provided task scope description vs. the deprecated SurfaceKey scope).

### [blindspot-3] `cleanTemplateOutput()` retention is an undiscussed design trade-off

The proposal states the goal is "移除三套后处理函数" (line 53) and lists removing all four conditional behaviors. But it retains the function for "空白行塌陷" (line 154). `text/template` has built-in whitespace trimming via `{{-` and `-}}` syntax. The proposal does not discuss whether this function could be fully eliminated, or why Go-code-based whitespace collapse is preferred over template-side whitespace control. This is a minor inconsistency: the framing says "remove post-processing" but the implementation keeps one form of post-processing alive.

Quote: "移除 `cleanTemplateOutput()` 中的全部四种条件删除逻辑...函数保留但仅做空白行塌陷"

What must improve: Either explain why template-side whitespace trimming (`{{-`) is insufficient, or commit to fully eliminating the function.

### [blindspot-4] `forge init` surface configuration is a phantom deliverable

Scope line 233 states: "在 `forge init` 中增加 surface 配置步骤（`forge surfaces detect` 集成）". This item:
- Is in Scope (In Scope section)
- Is referenced as a mitigation for the surface hard-failure risk (line 284)
- Has no corresponding SC entry
- Has no resource/timeline estimate (not in the Feasibility breakdown)
- Is not mentioned in User-Facing Behavior (lines 66-68)
- Is not listed in Constraints & Dependencies (lines 89-95)

A deliverable that appears in Scope and in Risk mitigation but nowhere else in the document is a phantom — it exists to make the risk mitigation look complete but has no verification path.

Quote: "在 `forge init` 中增加 surface 配置步骤（`forge surfaces detect` 集成），确保新项目初始化即具备 surfaces 配置，从源头避免硬性失败场景"

What must improve: Either add a SC entry and timeline estimate, or move to Out of Scope with an acknowledgment that the hard-failure mitigation is limited to error-message guidance.

### [blindspot-5] Metadata frontmatter validation system is an unmotivated expansion

The proposal's Problem section identifies frontmatter inconsistency (line 30): 6 record templates have frontmatter, 12 don't; 4 doc templates lack `SURFACE_KEY` headers. The proportional solution would be: add consistent frontmatter to the 12 missing templates, add `SURFACE_KEY` to the 4 doc templates. Instead, the proposal creates an entirely new system: metadata frontmatter with `variables` fields, a loader that strips metadata before `Parse()`, and reflection-based cross-validation between declared variables and struct fields (lines 236-249, SC lines 324-330).

This system is not motivated by any problem statement. The frontmatter inconsistency is a "fix the gaps" problem, not a "build a validation framework" problem. The proposal does not explain what future benefit the validation system provides beyond the current migration, or why simple frontmatter consistency is insufficient.

Quote: "模板校验：`ValidateTemplates()` 使用 `variables` 字段与模板数据 struct 的反射字段做交叉校验，确保声明与实现一致"

What must improve: Either motivate the metadata frontmatter system with a specific problem it solves (e.g., "future templates added by contributors need automated validation"), or reduce the solution to simple frontmatter consistency without the reflection-validation framework.

### [blindspot-6] No rollback plan for behavioral changes

The proposal has an excellent rollback strategy for the template engine migration: per-package independent commits, preserve old functions until golden-file tests pass (line 290). However, the behavioral changes (surface hard-failure in quality gate, `forge init` surface configuration) have no rollback plan. If the hard-failure behavior breaks production workflows for projects without surface configuration, reverting requires a separate code change — the per-package commit isolation doesn't help because the behavioral change is in `quality_gate.go`, not in the template engine.

Quote: "回滚策略：每个包独立迁移并独立提交...若需回滚，revert 对应包的提交即可，不影响其他包"

What must improve: Add a specific rollback plan for the surface hard-failure behavioral change (e.g., feature flag, configurable hard/soft behavior, or a revert commit strategy).

### [blindspot-7] `CoverageStrategy` design trade-off undermines declarative-template value proposition

The proposal's core value is replacing imperative post-processing with declarative templates. Yet for coverage, the design explicitly puts the coverage instruction text's source-of-truth in Go code, not in templates: "coverage 指令文本的 source-of-truth 在 Go 代码而非模板中" (line 172). This means coverage rendering logic is split: the Go code determines what text to emit, and the template just checks presence. If a new coverage state is added (e.g., "partial coverage"), both Go code and potentially the template must change. This partially reintroduces the imperative-logic-in-Go problem that the proposal claims to solve.

Quote: "coverage 指令文本的 source-of-truth 在 Go 代码而非模板中，这是有意的设计权衡：将分支逻辑集中在调用端，避免模板中增加二级条件判断"

What must improve: Acknowledge this trade-off explicitly and explain why template-level coverage differentiation is not worth the complexity. The current explanation ("避免模板中增加二级条件判断") is reasonable but framed as an unchallenged design decision rather than a trade-off.

---

## Score Summary

| Dimension | Score | Max |
|-----------|-------|-----|
| Problem Definition | 103 | 110 |
| Solution Clarity | 112 | 120 |
| Industry Benchmarking | 99 | 120 |
| Requirements Completeness | 95 | 110 |
| Solution Creativity | 45 | 100 |
| Feasibility | 88 | 100 |
| Scope Definition | 72 | 80 |
| Risk Assessment | 76 | 90 |
| Success Criteria | 70 | 80 |
| Logical Consistency | 79 | 90 |
| **Total** | **839** | **1000** |

---

## Attack Summary

### Critical Attacks

1. **[Scope/SC Alignment]** `forge init` surface configuration is scoped (line 233), used as risk mitigation (line 284), but has no SC entry, no timeline estimate, and no User-Facing Behavior description. It is a phantom deliverable. — Quote: "在 `forge init` 中增加 surface 配置步骤（`forge surfaces detect` 集成），确保新项目初始化即具备 surfaces 配置" — Either add SC + timeline, or move to Out of Scope and acknowledge incomplete mitigation.

2. **[Scope]** `test-run` SCOPE classification is incorrect — classified as "段落级" but has inline pattern `- Scope: {{SCOPE}}`. Incorrect classification will cause wrong migration (section wrap instead of line conditional). — Quote: "段落级（test-gen-contracts, test-gen-journeys, test-run）：`{{SCOPE}}` 独占一行或作为段落标题的一部分" — Verify `test-run.md` content and reclassify.

3. **[Logical Consistency]** Metadata frontmatter validation system (reflection-based `variables` cross-check) is not motivated by any problem statement. The problem section identifies frontmatter inconsistency, not a validation-framework gap. — Quote: "模板校验：`ValidateTemplates()` 使用 `variables` 字段与模板数据 struct 的反射字段做交叉校验" — Either add problem motivation or reduce scope to simple frontmatter consistency.

### Moderate Attacks

4. **[Solution Clarity]** `taskTemplateData.Scope` field (line 193) retains deprecated "scope" terminology while the entire proposal systematically replaces it. No justification provided for this exception. — Quote: "Scope string // 任务作用域描述" — Rename field or justify retention.

5. **[Solution Clarity]** `cleanTemplateOutput()` retention for whitespace collapse is an undiscussed design trade-off. `text/template` has `{{-`/`-}}` for whitespace trimming. — Quote: "函数保留但仅做空白行塌陷" — Explain why template-side trimming is insufficient or commit to full elimination.

6. **[Success Criteria]** PHASE_SUMMARY dual-block requirement (Scope line 156: "需分别用两个独立的 `{{if .PhaseSummary}}` 块包裹") is reflected in SC as singular "段落" (line 306). A verifier could pass the SC with only one block. — Quote: "`{{if .PhaseSummary}}` 条件块正确渲染：有值时注入段落，无值时段落消失" — SC should explicitly verify both the label-line and instruction-line conditionals.

7. **[Risk Assessment]** No rollback plan for behavioral changes (surface hard-failure). The per-package rollback strategy only covers template engine migration, not behavioral changes in `quality_gate.go`. — Quote: "回滚策略：每个包独立迁移并独立提交...revert 对应包的提交即可" — Add rollback plan for behavioral change (feature flag, configurable behavior, or revert strategy).

8. **[Requirements Completeness]** No scenario for `embed.FS` + `template.Parse()` init-time failure behavior. Currently templates fail lazily (per-template at runtime). After migration, ANY parse failure blocks application startup. This reliability regression is mentioned as an intentional trade-off (line 86) but not explored in Requirements or Risk Assessment. — Quote: "模板文件通过 `//go:embed` 嵌入二进制，`text/template.Parse()` 须在启动时完成——这是 fail-fast 权衡" — Add an NFR or Risk entry that specifies what happens on parse failure and why init-blocking is acceptable.

9. **[Solution Clarity]** Line 226 says "在 `promptTemplateData` 中增加 `TaskCategory string` 字段" but `TaskCategory` is already in the struct at line 169. This drafting artifact suggests incomplete revision cleanup. — Quote: "在 `promptTemplateData` 中增加 `TaskCategory string` 字段" — Remove redundant instruction.

10. **[Scope Definition]** `taskTemplateData.Scope` field uses deprecated terminology, contradicting the proposal's own systematic terminology migration from "scope" to SurfaceKey/SurfaceType/ScopeDisplay. — Quote: "Scope string // 任务作用域描述" — This field either needs renaming or explicit justification for being exempt.

### [blindspot] Attacks

11. **[blindspot-1]** `test-run` SCOPE classification error persists from iteration 0 — may cause incorrect template migration output.

12. **[blindspot-4]** `forge init` surface configuration is a phantom deliverable — scoped and used as mitigation but unverifiable.

13. **[blindspot-5]** Metadata frontmatter validation system is an unmotivated scope expansion — no problem statement drives the reflection-validation framework.

14. **[blindspot-6]** No rollback plan for surface hard-failure behavioral change — per-package rollback doesn't cover behavioral changes.

15. **[blindspot-7]** `CoverageStrategy` source-of-truth in Go code partially undermines the declarative-template value proposition — the trade-off is acknowledged but not analyzed against the proposal's own framing.

---

## Bias Detection Report

- Iteration-0 pre-revision annotated regions: 6 attack points were identified by pre-revision, and the revised proposal addressed all 6. This re-evaluation found the revisions adequate but identified new issues at the boundaries of the revised text.
- Scope creep pattern: The proposal exhibits a common pattern where the core solution (unified template engine) is sound, but ancillary scope items (metadata frontmatter system, `forge init` integration, reflection validation) expand the surface area without proportional problem motivation. These additions are individually reasonable but collectively increase the proposal's risk profile without explicit acknowledgment.
- Self-reported consistency check: The `consistency_check_result: pass` (lines 337-342) is self-reported. While the proposal's improvements make it more credible, an independent consistency check would be more trustworthy.

**Verdict**: The proposal is technically sound and well-evidenced for its core objective (unified template engine). Its weaknesses are in scope management: phantom deliverables, unmotivated scope expansions, and misclassified migration patterns. One more revision focusing on scope hygiene — either committing to or dropping the `forge init` deliverable, motivating the metadata validation system, and fixing the `test-run` classification — would bring it to approval quality.
