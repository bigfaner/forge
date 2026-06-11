# Proposal Evaluation Report — Iteration 2 (Re-Evaluation)

**Proposal**: Unified Template Engine for Prompt and Task Templates
**Date**: 2026-05-28
**Scorer**: CTO Adversarial Review (Iteration 2)
**Previous iteration**: 1 (score: 839/1000)

---

## Revision Verification: Iteration-1 Key Issues

| # | Issue | Resolution Status | Evidence |
|---|-------|-------------------|----------|
| 1 | `forge init` surface configuration phantom deliverable | **Resolved** | Timeline (line 131: "约 1 小时"), SC (line 312), Risk mitigation (line 284: "见 Scope + SC"). Three cross-references now exist. |
| 2 | `test-run` SCOPE misclassified as paragraph-level | **Resolved** | Line 203-204 now classifies `test-run` as inline alongside `doc-consolidate`. Verified against actual `test-run.md` (line 4: `- Scope: {{SCOPE}}`). Paragraph-level list only contains `test-gen-contracts, test-gen-journeys`. |
| 3 | Metadata frontmatter validation lacks problem motivation | **Partially resolved** | Line 79 adds motivation: "35 个模板分散在 3 个包中，人工审查无法保证模板变量与 Go struct 字段的对应关系". However, the leap from "frontmatter inconsistency" to "reflection-validation framework" remains a scope expansion. |
| 4 | `taskTemplateData.Scope` uses deprecated terminology | **Resolved** | Line 194 renamed to `ScopeDescription string` with explicit justification: "非 deprecated scope 概念；这里是 task-level 语境描述，非 surface-level 的 SurfaceKey". |
| 5 | `cleanTemplateOutput()` retention undiscussed trade-off | **Resolved** | Line 155 now explains: "text/template 的 `{{-`/`-}}` 可消除单个 `{{if}}` 块周围的空行，但无法处理连续空行...保留 Go 级空白行塌陷作为最终后处理步骤是必要的". Concrete technical justification provided. |
| 6 | PHASE_SUMMARY SC singular | **Resolved** | Line 306 now reads: "有值时同时注入标签行和条件指令行（两个独立 `{{if}}` 块），无值时两处均消失——覆盖所有 18 个使用 PHASE_SUMMARY 的模板". Dual-block requirement is explicit. |

---

## Phase 1: Reasoning Audit (Pre-Score Anchors)

### Argument Chain Trace

1. **Problem -> Solution**: Three sub-problems stated. Solution addresses all three. The chain is sound. The scope has been tightened relative to iteration 1: `forge init` now has SC + timeline; `test-run` classification is corrected; `taskTemplateData.ScopeDescription` is justified. Remaining concern: the metadata frontmatter validation system (lines 236-249) remains disproportionate to the "frontmatter 不一致" problem statement. The problem says templates have inconsistent frontmatter; the solution builds a reflection-based validation framework.

2. **Solution -> Evidence**: Template counts (21+2+12=35, 6 record) consistent throughout. All three template data structs now have complete field definitions with behavioral comments. The `promptTemplateData` has 11 fields, `taskTemplateData` has 5 fields, `autogenTemplateData` has 10 fields. Post-processing function behaviors documented with specific replacement patterns. Evidence is strong.

3. **Evidence -> Success Criteria**: 32 SC entries (up from 31 in iteration 1, with the addition of `forge init` SC at line 312). Coverage is now near-complete. One remaining gap: no SC for the `<!-- body-only -->` comment replacement (Scope line 248: "autogen body 模板的 metadata frontmatter 替代原有的 `<!-- body-only -->` 注释方案").

4. **Self-contradiction check**: No significant contradictions found. The `cleanTemplateOutput()` retention is now explicitly justified (line 155). The `taskTemplateData.ScopeDescription` naming is explained (line 194). The PHASE_SUMMARY dual-block is explicit in both Scope (line 157) and SC (line 306).

### Pre-Score Anchors

- **Strong**: Problem evidence, template data structs, urgency argument, rollback strategy, sequencing constraint, `forge init` now fully cross-referenced, PHASE_SUMMARY dual-block explicit, `test-run` classification corrected
- **Moderate concern**: Metadata frontmatter validation system still unmotivated by proportional problem statement; scope breadth for "quick mode" classification
- **Weakness**: Minor SC gaps (`<!-- body-only -->` removal not verified), `CoverageStrategy` source-of-truth trade-off acknowledged but not analyzed against proposal's own declarative framing

---

## Phase 2: Rubric Scoring (Dimension-by-Dimension)

### 1. Problem Definition (110 pts)

**Problem stated clearly (38/40)**: Core problem is unambiguous — three independent rendering engines with fragile post-processing hacks. The three sub-problems are clearly delineated with evidence tables. The "目录结构分散" (line 32) and "Skill/Command 层概念残留" (line 34) are well-documented secondary problems. Minor deduction: the "frontmatter 不一致" sub-problem (line 30) lists specific gaps (6 record templates have frontmatter, 12 don't; 4 doc templates lack SURFACE_KEY headers) but the solution's metadata-validation framework goes well beyond fixing these gaps. The problem-to-solution ratio is 1:3 for this sub-problem.

**Evidence provided (39/40)**: Excellent evidence. Three evidence tables with code-level detail. Template counts accurate (35+6=41). Post-processing function behaviors documented with specific fragility patterns (lines 36-39). The `injectSurfaceFrontmatter()` description is now corrected to "替换 `surface-key: ""` 字面值" (line 20). The only gap: the evidence section lists three fragility patterns for post-processing functions but Scope addresses four conditional behaviors in `cleanTemplateOutput()` — the `<!-- IF NOT_LOW -->` pattern is a feature that `task-pipeline-precision` will introduce, not a current fragility, creating a slight evidence-solution boundary issue.

**Urgency justified (29/30)**: Clear dependency on approved `task-pipeline-precision`. "Third layer of hack" argument is compelling and concrete. The sequencing constraint (line 95) adds weight. Minor deduction: urgency focuses on conditional-logic problem but the proposal includes directory restructuring, metadata frontmatter system, and `forge init` integration — these are not urgent, and their inclusion in the same proposal could delay the urgent core.

**Score: 106/110**

---

### 2. Solution Clarity (120 pts)

**Approach is concrete (39/40)**: The approach is highly concrete. Key strengths:
- All three template data structs have complete field definitions with behavioral comments (lines 166-179, 188-194, 209-220)
- All four `cleanTemplateOutput()` conditional behaviors are individually listed with replacement `{{if}}` patterns (lines 155-161)
- `{{SCOPE}}` two-pattern migration strategy is correctly documented: paragraph-level (test-gen-contracts, test-gen-journeys) vs inline (doc-consolidate, test-run) (lines 203-204)
- `BodyContext.Scope` -> `ScopeDisplay string` replacement is explicit (line 205)
- `TASK_CATEGORY` injection migration has its own sub-section (lines 224-227)
- `cleanTemplateOutput()` retention is justified with a concrete technical reason (line 155: consecutive empty lines from multiple omitted conditionals)
- `taskTemplateData.ScopeDescription` rename is justified (line 194)

Minor deduction: line 226 says "迁移到 `promptTemplateData.TaskCategory` 字段" but the TASK_CATEGORY sub-section (lines 224-227) reads as a standalone instruction rather than being redundant with the struct definition at line 169. The struct already has `TaskCategory` — the sub-section should say "use the existing TaskCategory field" rather than implying it needs to be added.

**User-facing behavior described (43/45)**: "No visible functional change" with one exception (surface hard-failure). Two-path distinction is explicit and well-documented:
- Quality gate path: hard failure (line 231)
- CLI path: soft behavior (line 232)
- `forge init` integration: preventive measure (line 233)

Deductions:
- The `forge init` integration (line 233) is mentioned in Scope but not in the User-Facing Behavior section (lines 66-68). A user-facing behavior change during project initialization should be described in that section.
- The error message content for surface hard-failure is described as "引导用户运行 `forge surfaces detect`" (line 68) — the actual error text is not shown. A concrete error message would strengthen the spec.

**Technical direction clear (34/35)**: `text/template` from stdlib, already proven in `pkg/task/data/` typed-task-records. Init-time validation with `missingkey=error` + zero-value Execute. Pre-formatted string fields to avoid template-side logic. `cleanTemplateOutput()` retention justified with whitespace-collapse argument. All clear. Minor deduction: the `CoverageStrategy` three-state design (line 169-172) is described inline but the trade-off analysis is buried in a struct comment rather than surfaced as a design decision. A reader might miss the implication: coverage instruction text is resolved in Go code, not in templates, which partially contradicts the declarative-template value proposition.

**Score: 116/120**

---

### 3. Industry Benchmarking (120 pts)

**Industry solutions referenced (30/40)**: References Helm, Hugo, goreleaser, buffalo. States "无竞争力的第三方替代品" — accurate for Go's `text/template` ecosystem. The Innovation Highlights section is honest: "这不是技术创新" (line 62). However, benchmarking remains shallow:
- No URLs, version numbers, or specific feature comparisons
- No discussion of how these tools handle specific challenges Forge faces (conditional sections, frontmatter management, validation)
- No discussion of alternative architectural approaches (code generation, builder pattern, DSL)

**At least 3 meaningful alternatives (25/30)**: Four alternatives listed. "标记注释 + 后处理" is a real alternative from `task-pipeline-precision`. "仅迁移 pkg/template" is a straw-man (2 templates, no serious consideration). The comparison table is concise but adequate.

**Honest trade-off comparison (22/25)**: Trade-offs are reasonable. The "Cons" for the selected approach says "35 个模板文件需更新占位符语法" — this still understates the scope. The actual migration involves 8 distinct activities (syntax update, conditional block insertion, SCOPE two-pattern migration, post-processing removal, frontmatter restructuring, directory reorganization, embed path updates, validation system implementation). The "Cons" should at minimum mention the behavioral change and cross-proposal dependency.

**Chosen approach justified (23/25)**: Justification is sound — standard library, already in codebase, declarative over imperative. "Already in codebase" (typed-task-records) is the strongest argument. Minor improvement: could cite the specific existing codebase usage that validates the pattern.

**Score: 100/120**

---

### 4. Requirements Completeness (110 pts)

**Scenario coverage (35/40)**: Six key scenarios identified. Edge cases well-covered:
- Coverage three-state semantics (line 172)
- `{{SCOPE}}` two-pattern migration (lines 203-204) — now correctly classified
- `just compile` trailing whitespace (line 159)
- Surface inference two-path behavior (lines 231-232)
- Metadata frontmatter motivation (line 79)
- PHASE_SUMMARY dual-block pattern (line 157)

Gaps:
- No scenario for `embed.FS` + `template.Parse()` init-time failure — what happens to the application when a template fails to parse? The NFR section (line 86) discusses this as a trade-off but no scenario explores the failure mode.
- No scenario for `ScopeDisplay` pre-formatting edge cases (empty slice, single-item slice)

**Non-functional requirements (36/40)**: Four NFRs stated:
- Backward compatibility: specific and testable (default to "medium")
- Byte-equivalent migration: clear criterion (content equivalence, whitespace tolerance)
- Zero performance regression: reasonable claim
- Reliability trade-off acknowledged: init-time fail-fast vs runtime fail-late (line 86)

Deductions:
- The reliability NFR (line 86) is framed as a "已确认" trade-off but lacks a formal acceptance criterion. What is the maximum acceptable startup time increase? What is the fallback if init-time parsing adds unacceptable latency?
- No NFR for metadata frontmatter validation system's startup performance impact — reflection-based cross-checking of 41 templates at startup has a measurable cost.

**Constraints & dependencies (28/30)**: Good list. Key constraints:
- `//go:embed` + init-time parsing (line 90)
- Cross-proposal sequencing (line 95) — now explicit about order
- `missingkey=error` + zero-value Execute (line 86)
- `pkg/template` dual consumer paths (line 91)

Minor gap: the `forge init` surface configuration dependency is now in Scope and SC but not in Constraints section — if it's a prerequisite for the hard-failure behavior to be acceptable, it should be listed as a constraint.

**Score: 99/110**

---

### 5. Solution Creativity (100 pts)

**Novelty over industry baseline (15/40)**: The proposal explicitly states "这不是技术创新" — honest. This is standard `text/template` usage applied to a consolidation problem. No differentiation from the Helm/Hugo/goreleaser pattern.

**Cross-domain inspiration (8/35)**: No cross-domain inspiration. No borrowing from other domains. The solution is the most obvious choice for a Go project.

**Simplicity of insight (23/25)**: "Replace 3 post-processing hacks with 1 declarative system" is clean and elegant. The `just compile{{if .SurfaceKey}} {{.SurfaceKey}}{{end}}` pattern (line 159) is a nice declarative touch. The `ScopeDisplay` pre-formatted string approach avoids template-side `{{range}}` loops. The `cleanTemplateOutput()` retention justification (line 155) demonstrates thoughtful analysis of the whitespace-collapse problem. Solid engineering decisions, well-reasoned.

**Score: 46/100**

---

### 6. Feasibility (100 pts)

**Technical feasibility (37/40)**: Pattern already proven in `pkg/task/data/` (typed-task-records). Migration path is mechanical for placeholder replacement. Template data structs are complete. The `ScopeDisplay` pre-formatting approach is straightforward. The `cleanTemplateOutput()` retention with concrete whitespace-collapse justification removes ambiguity.

Deductions:
- The metadata frontmatter system (strip frontmatter, parse `variables`, reflect over struct fields, cross-validate) adds non-trivial complexity. The proposal estimates ~2 hours (line 128) but this includes: (1) writing a metadata parser, (2) implementing reflection-based validation, (3) restructuring 6 record templates, and (4) integrating validation into startup. 2 hours seems optimistic for this scope.
- The `CoverageStrategy` three-state semantics relies on Go code to resolve the correct value. The template cannot distinguish states. If future requirements need template-level state differentiation, the field must be redesigned.

**Resource & timeline feasibility (26/30)**: 15 coding tasks, ~16 hours total (line 131 adds `forge init` integration at 1 hour — this was missing in iteration 0). The breakdown is good. However:
- "Metadata frontmatter: ...模板加载器解析 metadata：约 2 小时" (line 128) — as noted above, 2 hours for 4 distinct sub-activities is optimistic.
- Total estimate may be 15-20% low given the scope breadth (metadata system + directory restructuring + behavioral change + template migration + `forge init` integration all in one proposal).

**Dependency readiness (28/30)**: No external dependencies. `text/template` is stdlib. Cross-proposal dependency is sequenced (line 95).

**Score: 91/100**

---

### 7. Scope Definition (80 pts)

**In-scope items are concrete (28/30)**: Most items are concrete deliverables. The template data structs with field definitions are excellent. The conditional-behavior replacement patterns are specific. The `{{SCOPE}}` two-pattern migration is now correctly classified and concrete. The `taskTemplateData.ScopeDescription` rename is justified. The `forge init` integration (line 233) still has a minor vagueness: "集成 surface 配置步骤" — what does integration mean? A prompt during init? An automatic detection step? Running `forge surfaces detect` as a sub-process? The term "集成" is less precise than the rest of the proposal.

Minor deduction: the Scope lists 15 tasks at ~16 hours for "quick mode" — this proposal includes template engine migration, behavioral changes, directory restructuring, metadata frontmatter system, skill/command concept alignment, and `forge init` integration. This is a substantial scope for a single quick-mode proposal.

**Out-of-scope explicitly listed (24/25)**: Seven items explicitly listed. Good scope containment. Minor deduction: `ARCHITECTURE.md` updates are excluded (line 275) but the directory restructuring will make `ARCHITECTURE.md` inaccurate — excluding documentation updates from a restructuring proposal creates a maintenance gap.

**Scope is bounded (22/25)**: 15 tasks, ~16 hours, quick mode. The sequencing constraint is clear. The `task-pipeline-precision` dependency is explicit. However:
- The scope breadth (engine migration + behavioral change + directory restructuring + metadata system + `forge init` + concept alignment) is large for a single proposal. The risk of scope interaction is not discussed — what if the metadata frontmatter system reveals issues that require template redesign?
- The surface hard-failure behavioral change (line 231) is acknowledged as a behavior change, not just refactoring, but the scope classification doesn't reflect the additional testing burden for a behavioral change.

**Score: 74/80**

---

### 8. Risk Assessment (90 pts)

**Risks identified (27/30)**: Eleven risks identified. Good coverage of migration risks, cross-proposal synchronization, and nil-pointer safety. The `cleanTemplateOutput()` retention justification (line 155) addresses the whitespace-difference risk proactively. Gaps:
- No risk for the metadata frontmatter system — YAML parsing, reflection-based validation, and frontmatter stripping are new code paths that could have edge cases (malformed frontmatter, missing `variables` field, body content starting with `---`).
- No risk for `ScopeDisplay` field replacing `BodyContext.Scope []string` — the type change from slice to pre-formatted string means all callers must be updated, not just `renderBody()`.

**Likelihood + impact rated (26/30)**: Ratings are generally reasonable. The "BodyContext.Scope 重命名影响 BuildIndex() 调用链" is rated M/M (line 288), and the mitigation correctly notes that "Scope 仅在 autogen.go 内部消费". The rating is consistent with the mitigation. One concern: "Surface 硬性失败阻塞无 surfaces 配置的项目" (line 284) is rated M/H — the mitigation now references "见 Scope + SC" which is an improvement from iteration 1, but the mitigation relies on `forge init` integration which has its own implementation risk.

**Mitigations are actionable (27/30)**: Most mitigations are concrete and actionable:
- `missingkey=error` + zero-value Execute — technically correct
- Per-package independent commits with golden-file testing — excellent rollback strategy
- Two-step separation (engine migration first, directory restructuring second) — sound risk reduction
- `forge init` surface configuration — now scoped with SC + timeline

Deductions:
- The rollback strategy (line 290) only covers template engine migration. The behavioral change (surface hard-failure in `quality_gate.go`) has no specific rollback plan. If the hard-failure behavior breaks production workflows for projects without surface configuration, reverting requires a targeted code change — the per-package commit isolation helps but the behavioral change is a separate concern from template engine correctness.
- The `forge init` integration is now a full deliverable, which is an improvement. But the mitigation for "Surface 硬性失败阻塞无 surfaces 配置的项目" (line 284) still partially relies on a new feature (`forge init` integration) as a risk mitigation — this creates a circular dependency: the feature reduces the risk, but the feature itself has implementation risk.

**Score: 80/90**

---

### 9. Success Criteria (80 pts)

**Criteria are measurable and testable (28/30)**: Most SC items are binary checkable. "Does this function use `text/template.Execute()`?" — yes/no. "Are `removeLineContaining()` and `removeSection()` removed?" — yes/no. The PHASE_SUMMARY dual-block SC (line 306) now explicitly states "两个独立 `{{if}}` 块" — unambiguous. The `cleanTemplateOutput()` SC (line 301) lists all four conditional behaviors individually.

Deductions:
- "`forge prompt get-by-task-id` 输出与迁移前功能等价（golden-file 对比，允许空白行差异）" (line 313) — "功能等价" is the criterion but the threshold is not defined. What percentage of templates must pass golden-file comparison? What constitutes a content difference vs. a whitespace difference? The method (golden-file) is clear, the pass criterion is vague.
- "`ValidateTemplates()` 使用 `variables` 字段与模板数据 struct 做交叉校验" (line 331) — "交叉校验" is vague. Does it verify every declared variable exists in the struct? That every struct field is declared? That types match? The direction and completeness of validation is unspecified.

**Coverage is complete (23/25)**: 32 SC entries cover most Scope items. Significant improvement from iteration 1. Gaps:
- No SC for the `<!-- body-only -->` comment replacement (Scope line 248: "autogen body 模板的 metadata frontmatter 替代原有的 `<!-- body-only -->` 注释方案"). If `<!-- body-only -->` comments are being replaced by metadata frontmatter, there should be an SC verifying they are removed.
- The `BodyContext` SC (line 305) says "已替换为当前概念对齐的字段名" — the parenthetical uses vague language ("当前概念对齐的字段名") rather than naming `ScopeDisplay string` directly. Scope line 205 specifies `ScopeDisplay` — the SC should name it.

**SC internal consistency (23/25)**: The SC set is largely internally consistent. The `consistency_check_result: pass` (lines 339-342) is more credible given the improvements. Remaining tension:
- SC "35 个模板文件中无 `{{PLACEHOLDER}}` 格式" (line 303) — but 6 record templates also need migration consideration. They already use `text/template` but are included in the 41-template metadata frontmatter count (line 325). The number "35" for syntax migration and "41" for frontmatter is correct because record templates don't need syntax migration, but a verifier might find this confusing without the explicit explanation in Out of Scope (line 271).

**Score: 74/80**

---

### 10. Logical Consistency (90 pts)

**Solution addresses the stated problem (34/35)**: Unified `text/template` directly addresses three-engine fragmentation. Conditional `{{if}}` replaces all post-processing hacks. Stale concept cleanup is explicit. Directory restructuring addresses the "目录结构分散" problem. The solution is well-targeted at the core problem. Minor deduction: the "frontmatter 不一致" sub-problem (line 30) identifies two specific issues (inconsistent frontmatter presence, missing SURFACE_KEY headers). The solution addresses both but also creates an entirely new metadata-validation framework that is not proportional to the stated problem. The motivation added at line 79 ("35 个模板分散在 3 个包中，人工审查无法保证...") is reasonable but is a new argument not present in the Problem section — it belongs in the Problem section if it's a primary motivator.

**Scope <-> Solution <-> SC aligned (27/30)**: Significant improvement. Most Scope items have corresponding SC entries. The `forge init` integration now has SC (line 312), timeline (line 131), and risk mitigation reference (line 284). Remaining misalignments:
- Scope line 248 (metadata frontmatter replaces `<!-- body-only -->` comments) has no SC entry verifying `<!-- body-only -->` removal
- Scope line 155 (`cleanTemplateOutput()` retention for whitespace collapse) is reflected in SC (line 301) but the SC frames it as "仅保留空白行塌陷逻辑" without verifying that the whitespace collapse behavior produces correct output

**Requirements <-> Solution coherent (22/25)**: Requirements map cleanly to solution for the core template-engine migration. Deductions:
- The "Non-Functional Requirements" section mentions "字节等价迁移" (line 84) but Scope includes behavioral changes (surface hard-failure, `forge init` integration) that are not "equivalent" to current behavior. The NFR should clarify which paths must be equivalent and which are allowed to change.
- The `taskTemplateData.ScopeDescription` field (line 194) is now well-justified, removing the terminology inconsistency from iteration 1.

**Score: 83/90**

---

## Phase 3: Blindspot Hunt

### [blindspot-1] Metadata frontmatter validation system: motivation gap between problem and solution scope

The Problem section (line 30) identifies frontmatter inconsistency: 6 record templates have frontmatter, 12 don't; 4 doc templates lack SURFACE_KEY headers. The proportional solution is: add consistent frontmatter to the 12 missing templates, add SURFACE_KEY to the 4 doc templates. Instead, the solution creates a full validation framework: metadata frontmatter with `variables` fields, a loader that strips metadata before `Parse()`, and reflection-based cross-validation (lines 236-249, SC lines 325-331). The motivation added at line 79 is reasonable ("35 个模板分散在 3 个包中，人工审查无法保证模板变量与 Go struct 字段的对应关系") but this is a future-maintenance argument, not a current-problem argument. The current problem is inconsistency, not validation-framework absence.

**Verdict**: Moderate. The validation system is good engineering but represents scope creep relative to the stated problem. The motivation at line 79 should be in the Problem section, not the Requirements section.

### [blindspot-2] No rollback plan for behavioral change

The rollback strategy (line 290) is excellent for template engine migration: per-package independent commits, preserve old functions until golden-file tests pass. However, the behavioral change (surface hard-failure in `quality_gate.go`, line 231) has no specific rollback plan. If the hard-failure behavior breaks production workflows for projects without surface configuration, the per-package commit isolation helps but the behavioral change is a separate concern. The proposal acknowledges this is a "行为变更" but does not discuss rollback specifically for this change.

**Verdict**: Moderate. A feature flag or configurable hard/soft behavior would reduce rollback risk for the behavioral change specifically.

### [blindspot-3] `CoverageStrategy` source-of-truth trade-off not fully analyzed

The proposal's core value is replacing imperative post-processing with declarative templates. Yet for coverage, the design explicitly puts coverage instruction text's source-of-truth in Go code (line 172: "coverage 指令文本的 source-of-truth 在 Go 代码而非模板中"). This means coverage rendering logic is split: Go code determines what text to emit, and the template checks presence. The three-state semantics (empty/maintain/percentage) are resolved in Go code, and the template cannot distinguish states. This is acknowledged as a "有意的设计权衡" but is not analyzed against the proposal's own declarative-template framing. If a fourth state is needed, Go code must change.

**Verdict**: Low. The trade-off is reasonable and acknowledged. The analysis is adequate but could be more explicit about what would trigger a redesign.

### [blindspot-4] `forge init` integration implementation detail

The `forge init` integration (line 233) says "集成 surface 配置步骤（`forge surfaces detect`）" — the term "集成" is vague. Does this mean: (a) running `forge surfaces detect` as a subprocess during init? (b) Adding a prompt asking the user if they want to run surface detection? (c) Integrating the detection logic directly into the init command? The implementation approach affects the timeline estimate (line 131: "约 1 小时") — subprocess integration is simpler than embedded logic. The SC (line 312) says "新项目初始化时具备 surfaces 配置" which is outcome-oriented but doesn't specify the mechanism.

**Verdict**: Low. The outcome is clear even if the mechanism is vague. An implementer can choose the simplest approach that satisfies the SC.

---

## Score Summary

| Dimension | Score | Max | Delta from Iter 1 |
|-----------|-------|-----|--------------------|
| Problem Definition | 106 | 110 | +3 |
| Solution Clarity | 116 | 120 | +4 |
| Industry Benchmarking | 100 | 120 | +1 |
| Requirements Completeness | 99 | 110 | +4 |
| Solution Creativity | 46 | 100 | +1 |
| Feasibility | 91 | 100 | +3 |
| Scope Definition | 74 | 80 | +2 |
| Risk Assessment | 80 | 90 | +4 |
| Success Criteria | 74 | 80 | +4 |
| Logical Consistency | 83 | 90 | +4 |
| **Total** | **869** | **1000** | **+30** |

---

## Attack Summary

### Critical Attacks (1)

1. **[Requirements/Logical Consistency]** Metadata frontmatter validation system (reflection-based `variables` cross-check, lines 236-249) is not motivated by the Problem section. The Problem (line 30) identifies frontmatter inconsistency; the solution builds a validation framework. The motivation at line 79 ("35 个模板分散在 3 个包中...") is a future-maintenance argument that belongs in the Problem section, not buried in Requirements. Without this motivation in Problem, the validation system appears as scope creep — it adds ~20% of the proposal's scope without proportional problem justification.
   - Quote: "模板校验：`ValidateTemplates()` 使用 `variables` 字段与模板数据 struct 的反射字段做交叉校验，确保声明与实现一致" (line 249)
   - What must improve: Either (a) promote line 79's motivation to the Problem section as a fourth sub-problem ("模板变量与 Go struct 字段的对应关系缺乏自动校验"), or (b) reduce the solution to simple frontmatter consistency (add frontmatter to missing templates, add SURFACE_KEY to doc templates) without the reflection-validation framework.

### Moderate Attacks (5)

2. **[Risk Assessment]** No rollback plan for the surface hard-failure behavioral change (line 231). The per-package rollback strategy (line 290) covers template engine migration but not behavioral changes. If hard-failure breaks production workflows, a targeted revert of `quality_gate.go` is needed — but the proposal does not discuss this as a separate rollback concern.
   - Quote: "回滚策略：每个包独立迁移并独立提交...若需回滚，revert 对应包的提交即可" (line 290)
   - What must improve: Add a specific rollback plan for the behavioral change (e.g., feature flag, configurable hard/soft behavior, or a targeted revert strategy for `quality_gate.go`).

3. **[Success Criteria]** No SC for `<!-- body-only -->` comment removal. Scope line 248 states: "autogen body 模板的 metadata frontmatter 替代原有的 `<!-- body-only -->` 注释方案" — this implies `<!-- body-only -->` comments should be removed. No SC verifies this removal.
   - Quote: "autogen body 模板的 metadata frontmatter 替代原有的 `<!-- body-only -->` 注释方案" (line 248)
   - What must improve: Add an SC entry: "`<!-- body-only -->` comments removed from all autogen body templates, replaced by metadata frontmatter".

4. **[Success Criteria]** `ValidateTemplates()` SC (line 331) is vague about validation direction and completeness. "交叉校验" does not specify: (a) every declared variable must exist in struct, (b) every struct field must be declared, or (c) both.
   - Quote: "`ValidateTemplates()` 使用 `variables` 字段与模板数据 struct 做交叉校验" (line 331)
   - What must improve: Specify the validation direction, e.g., "every variable declared in metadata frontmatter exists as a field in the corresponding template data struct (forward check only; struct may have computed fields not in variables)".

5. **[Success Criteria]** `BodyContext` SC (line 305) uses vague language. "已替换为当前概念对齐的字段名" does not name the replacement field directly. Scope (line 205) specifies `ScopeDisplay string` — the SC should match.
   - Quote: "`BodyContext` struct 中无 `Scope` 字段（已替换为当前概念对齐的字段名）" (line 305)
   - What must improve: Change parenthetical to "已替换为 `ScopeDisplay string`" to match Scope's explicit naming.

6. **[Requirements Completeness]** "字节等价迁移" NFR (line 84) conflicts with behavioral changes in Scope. Surface hard-failure (line 231) and `forge init` integration (line 233) are not "equivalent" to current behavior. The NFR should clarify scope: "模板渲染输出与当前行为功能等价（不包括 Scope 定义的明确行为变更）".
   - Quote: "迁移后 prompt 输出与当前 `strings.ReplaceAll` + `cleanTemplateOutput()` 的输出在功能上等价" (line 84)
   - What must improve: Scope the NFR to exclude paths where behavioral changes are explicitly defined in Scope.

### [blindspot] Attacks (2)

7. **[blindspot-1]** Metadata frontmatter validation system motivation gap — the argument is in Requirements (line 79) but should be in Problem section. See Critical Attack #1.

8. **[blindspot-2]** No rollback plan for behavioral change — per-package rollback covers template engine migration only. See Moderate Attack #2.

---

## Bias Detection Report

- **Iteration 1 issues**: All 6 key issues were addressed. The revisions are substantive, not cosmetic. The `test-run` classification was verified against actual template content and is now correct. The `taskTemplateData.ScopeDescription` rename includes explicit justification. The `cleanTemplateOutput()` retention has a concrete technical argument.
- **Scope creep pattern**: The metadata frontmatter validation system remains the only significant scope expansion without proportional problem motivation. This was flagged in iteration 1 and partially addressed with a motivation statement at line 79, but the motivation is in the wrong section (Requirements instead of Problem).
- **Improvement trajectory**: +75 points from baseline (764 -> 839), +30 from iteration 1 (839 -> 869). The rate of improvement is slowing, which is expected as the remaining issues are structural (scope creep, NFR scoping) rather than factual errors or missing information.
- **Self-reported consistency check**: The `consistency_check_result: pass` (lines 339-342) is more credible in iteration 2. The one remaining inconsistency (NFR "字节等价" vs behavioral changes in Scope) is a scoping issue, not a factual contradiction.

**Verdict**: The proposal is in strong shape for its core objective (unified template engine migration). The remaining issues are: (1) one moderate scope-justification gap (metadata validation framework), (2) one missing SC (`<!-- body-only -->` removal), (3) SC vagueness on validation direction and BodyContext field naming, and (4) NFR scoping conflict with behavioral changes. These are addressable in a focused revision without restructuring the proposal. **Approval threshold (850+) reached.**

---

SCORE: 869/1000
DIMENSIONS:
  Problem Definition: 106/110
  Solution Clarity: 116/120
  Industry Benchmarking: 100/120
  Requirements Completeness: 99/110
  Solution Creativity: 46/100
  Feasibility: 91/100
  Scope Definition: 74/80
  Risk Assessment: 80/90
  Success Criteria: 74/80
  Logical Consistency: 83/90
ATTACKS:
1. [Requirements/Logical Consistency]: Metadata frontmatter validation system not motivated by Problem section — quote: "模板校验：ValidateTemplates() 使用 variables 字段与模板数据 struct 的反射字段做交叉校验" (line 249) — promote line 79 motivation to Problem section as a fourth sub-problem, or reduce solution to simple frontmatter consistency.
2. [Risk Assessment]: No rollback plan for surface hard-failure behavioral change — quote: "回滚策略：每个包独立迁移并独立提交...revert 对应包的提交即可" (line 290) — add feature flag, configurable behavior, or targeted revert strategy for quality_gate.go.
3. [Success Criteria]: No SC for <!-- body-only --> comment removal — quote: "autogen body 模板的 metadata frontmatter 替代原有的 <!-- body-only --> 注释方案" (line 248) — add SC verifying <!-- body-only --> comments are removed from all autogen body templates.
4. [Success Criteria]: ValidateTemplates() SC vague on validation direction — quote: "ValidateTemplates() 使用 variables 字段与模板数据 struct 做交叉校验" (line 331) — specify: every declared variable must exist in struct (forward check only).
5. [Success Criteria]: BodyContext SC uses vague language — quote: "已替换为当前概念对齐的字段名" (line 305) — change to "已替换为 ScopeDisplay string".
6. [Requirements Completeness]: NFR "字节等价迁移" conflicts with behavioral changes — quote: "迁移后 prompt 输出...在功能上等价" (line 84) — scope NFR to exclude paths with explicitly defined behavioral changes.
7. [Logical Consistency]: Metadata validation framework motivation in wrong section — quote: line 79 motivation in Requirements, not Problem — move to Problem section.
8. [Risk Assessment]: Behavioral change rollback relies on feature that itself has implementation risk — quote: "forge init 中集成 surface 配置步骤" as mitigation (line 284) — circular dependency between feature and risk mitigation.
