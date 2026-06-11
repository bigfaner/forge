# Evaluation Report: Autogen Test Task Path References

**Evaluator**: CTO Adversary (Iteration 0 Baseline)
**Date**: 2026-05-28
**Proposal**: `docs/proposals/autogen-test-task-paths/proposal.md`
**Rubric**: `plugins/forge/skills/eval/rubrics/proposal.md` (1000 pts)

---

## Phase 1: Reasoning Audit — Problem -> Solution -> Evidence -> SC Chain

### Chain Trace

1. **Problem**: Auto-generated test pipeline tasks lack feature-level path context. Subagent cannot locate testing artifacts from task .md files alone and must do extra exploration.
2. **Evidence**: Three-layer architecture table showing path information distribution across Task embed templates, Prompt templates, and Skills. Concrete example of prompt template output missing `FEATURE_SLUG`. Reference to run-tests skill Step 1.5 listing 3 slug acquisition sources as proof of the friction.
3. **Solution**: Three-layer coordination: embed templates get `## Feature Paths` discovery commands, prompt templates get `FEATURE_SLUG` rendering, skill unchanged.
4. **Success Criteria**: Checklist of template changes + build/test pass.

### Chain Verdict

The chain is **shallow but coherent**. The problem is real — verified against actual templates which confirm `FeatureSlug` is declared in context but not rendered to output in the six test-pipeline prompt templates. However, the solution's value proposition is overstated relative to its scope. The proposal admits urgency is low ("不影响正确性"), and the freeform review confirmed that embed templates `test-gen-journeys.md` and `test-gen-contracts.md` already contain extensive `{{.FeatureSlug}}` path references, making the proposed `## Feature Paths` section redundant for half the templates.

---

## Phase 2: Rubric Scoring

### Dimension 1: Problem Definition — 72/110

**Problem stated clearly (30/40)**: The core problem is unambiguous — `FeatureSlug` is declared but not rendered, forcing agents to reverse-engineer slug from path. However, the problem conflates two distinct issues: (a) prompt template missing slug rendering, and (b) embed template lacking discovery commands. These are separable problems with different severity levels. Quote: "生成内容过于简略，缺少 feature 级路径上下文" — this is vague. "过于简略" for which templates? As verified, `test-gen-journeys.md` and `test-gen-contracts.md` embed templates are rich with path context. The problem statement should distinguish between the thin templates (`test-run.md`, `test-gen-scripts.md`, and the two eval templates) that genuinely lack paths and the rich templates that already have them.

**Evidence provided (30/40)**: The three-layer architecture table is well-structured and the concrete prompt template output example is effective. The reference to skill Step 1.5's 3 slug sources is strong evidence. Deduction: the table's first row ("gen-journeys/gen-contracts 较完整；test-run/test-gen-scripts 几乎为空") acknowledges that some templates are not broken but buries this distinction rather than making it central to the problem framing. This leads to a one-size-fits-all solution for a problem that affects templates unevenly.

**Urgency justified (12/30)**: Quote: "不影响正确性（skill 能动态发现路径），但降低了 agent 执行效率". The proposal admits correctness is unaffected. No quantification of the efficiency loss. How many seconds/tokens are wasted? How often does this path occur? The urgency section reads as an afterthought — a single sentence acknowledging the problem is non-critical. If the problem is non-critical, the proposal should justify why it deserves implementation time now rather than later.

### Dimension 2: Solution Clarity — 78/120

**Approach is concrete (35/40)**: The three changes are clearly specified: embed templates get `## Feature Paths` section, prompt templates get `FEATURE_SLUG` line, skill unchanged. The template format is shown with code blocks. A reader can explain back what will be built.

**User-facing behavior described (22/45)**: This is the weakest part. The "user" here is the subagent, but the proposal never describes the observable behavior change from the agent's perspective. What does the agent do differently after this change? The "三层联动效果" table shows information flow but not behavioral change. Quote: "Subagent 收到 prompt，直接看到 `FEATURE_SLUG: my-feature`，无需从路径解析" — this is the closest to a user-facing behavior description, but it describes information availability, not action. Does the agent skip Step 1.5 entirely? Does it still run the skill's discovery logic? The proposal is silent on how the agent's workflow changes.

**Technical direction clear (21/35)**: The proposal correctly identifies `{{.FeatureSlug}}` as the template variable and names the rendering engines (`autogenTemplateData` for embed, `promptTemplateData` for prompt). However, it asserts "不依赖 `.forge/state.json` 或分支名" without evidence, and does not describe the rendering pipeline in enough detail to verify. Quote: "`{{.FeatureSlug}}` 由 CLI 在 `forge task index` 时从 `docs/features/<slug>/` 目录路径填充" — this is a claim about Go code behavior that the proposal does not substantiate with code references. The distinction between `autogenTemplateData` (embed) and `promptTemplateData` (prompt) rendering timing is not explained, which the freeform review correctly identified as a temporal coupling risk.

### Dimension 3: Industry Benchmarking — 32/120

**Industry solutions referenced (5/40)**: No industry solutions, patterns, open-source projects, or published approaches are referenced. The proposal is entirely self-contained. Quote: the Alternatives section lists four rows — all internal options with no external reference. For a proposal about context propagation in agent task systems, relevant industry patterns exist: environment variable propagation in CI/CD systems, context injection in agent frameworks (LangChain, AutoGen), or template variable rendering in static site generators. None are cited.

**At least 3 meaningful alternatives (18/30)**: Four alternatives are listed including "do nothing". However, two are straw men: "只改 embed 模板" is dismissed with "联动不完整" without explaining why incomplete联动 is worse than current state. "改 embed + prompt + 简化 skill" is dismissed with "skill 不能简化" which is an assertion without evidence. The chosen option is pre-ordained — the alternatives exist to be rejected.

**Honest trade-off comparison (5/25)**: The chosen approach's Cons column says "无" (none). Every solution has trade-offs. The freeform review identified at least three: temporal coupling between embed and prompt rendering, redundancy with existing rich templates, and new sync burden between embed paths and skill paths. The proposal acknowledges none of these.

**Chosen approach justified against benchmarks (4/25)**: No industry benchmarks exist to justify against. The justification is purely internal consistency ("三层联动，职责清晰"), which is circular reasoning — the proposal defines the three-layer model and then praises its own architecture.

### Dimension 4: Requirements Completeness — 60/110

**Scenario coverage (20/40)**: Three key scenarios are listed — all happy path. No edge cases identified. Critical missing scenarios: What happens when `FeatureSlug` is empty? What happens when the feature directory has been renamed since `forge task index`? What happens when a task is executed outside the pipeline (standalone invocation) — does the embed template still make sense? The proposal mentions standalone invocation for skills but does not address what happens when a user reads a task .md file directly without the pipeline.

**Non-functional requirements (18/40)**: No NFRs are stated. Relevant NFRs include: template rendering performance (adding `ls` commands to 12 templates), token budget impact (each `## Feature Paths` section adds ~100 tokens to agent context), and maintainability (12 templates now contain path patterns that must stay in sync with skills). The proposal mentions "不修改任何 Go 代码" as if this is inherently good, but the freeform review correctly noted that a Go-level validation of FeatureSlug non-empty would be a stronger guarantee.

**Constraints & dependencies (22/30)**: The constraints section is the strongest part of this dimension. It correctly identifies that `FeatureSlug` is already declared in identity/context groups and that Go structs already have the field. Quote: "Embed 模板：`{{.FeatureSlug}}` 已在 identity 组声明，不需新增变量" — this is precise and verifiable. Deduction: the section does not mention the dependency on `forge task index` being run before embed templates are materialized, nor does it mention the dependency on the dispatcher passing FeatureSlug at dispatch time.

### Dimension 5: Solution Creativity — 20/100

**Novelty over industry baseline (5/40)**: The proposal explicitly disclaims novelty. Quote: "无创新，纯信息补全。将已有的 `FeatureSlug` 变量从'声明但不渲染'变为'渲染到输出'". While honesty is appreciated, this means there is zero novelty. The solution is a straightforward bug-fix style change — rendering a variable that was already wired up but not emitted.

**Cross-domain inspiration (5/35)**: No cross-domain inspiration is evident or claimed. The solution does not borrow ideas from other domains.

**Simplicity of insight (10/25)**: The insight is simple but not in an elegant way — it is more "obvious" than "why didn't I think of that". The observation that a declared variable should be rendered is a correctness fix, not a creative leap. The three-layer framing adds some organizational value but does not represent a new mental model.

### Dimension 6: Feasibility — 82/100

**Technical feasibility (36/40)**: High. The changes are template-only modifications to .md files. The variables are already wired. The rendering engines already support them. The main deduction: the proposal claims "不修改任何 Go 代码（仅改模板 .md 文件）" but the freeform review correctly identified that validating FeatureSlug non-empty at render time would require a small Go change. The proposal forecloses this option without justification.

**Resource & timeline feasibility (28/30)**: The scope is 12 template files with mechanical changes. A single developer can complete this in under an hour. No bandwidth concerns.

**Dependency readiness (18/30)**: The proposal assumes `FeatureSlug` is always populated. It states "CLI 在 `forge task index` 时从 `docs/features/<slug>/` 目录路径填充" but does not verify what happens when tasks are generated outside the `docs/features/` directory structure (e.g., project-level tasks that are not feature-scoped). If a task exists at a path that does not match `docs/features/<slug>/`, the slug extraction from path would fail or produce incorrect results. This is an unstated dependency on directory structure convention.

### Dimension 7: Scope Definition — 65/80

**In-scope items are concrete (25/30)**: The 12 template files are enumerated by name. The change to each is specified (add `## Feature Paths` section or add `FEATURE_SLUG` line). This is executable. Minor deduction: the proposal does not specify the exact location within each file where changes should be made (beyond "TASK_FILE 行之后" for prompt templates, which is imprecise for templates with conditional rendering blocks).

**Out-of-scope explicitly listed (20/25)**: Four items are listed: Skill, Go structs, verification/doc/cleanup templates, Go code. The "Go 代码" out-of-scope is asserted but may be suboptimal as noted above. Missing from out-of-scope: What about documentation updates? What about the freeform review's suggestion of a sync validation mechanism? These are not mentioned as deferred items.

**Scope is bounded (20/25)**: The scope is clearly bounded to 12 files with mechanical changes. The proposal does not address what happens if new test-pipeline task types are added in the future — should they automatically get the same treatment? This is an open-ended maintenance question that is not addressed.

### Dimension 8: Risk Assessment — 40/90

**Risks identified (12/30)**: Only two risks are listed. Both are low-likelihood. The freeform review identified at least five additional risks: (1) temporal coupling between embed and prompt rendering, (2) redundancy with existing rich template content, (3) embed empty-slug producing invalid commands, (4) skill vs. embed path drift over time, (5) agent executing discovery commands at the wrong time in the workflow. The proposal's risk table is dangerously incomplete.

**Likelihood + impact rated (12/30)**: Both risks are rated "L" likelihood. The FeatureSlug empty-slug risk is rated "L" impact, but as the freeform review noted, an empty slug in the embed template produces `ls docs/features//testing/` — an actively misleading command. This should be rated Medium impact at minimum. The assessment is not honest — it minimizes the downsides.

**Mitigations are actionable (16/30)**: Quote for "路径与实际目录不一致" mitigation: "路径基于 forge 硬编码目录约定". This is not a mitigation — it is an assertion that the problem won't happen because of a convention. Conventions can be violated. Quote for "FeatureSlug 渲染为空" mitigation: "prompt 模板中 FeatureSlug 已在 context 声明，dispatcher 始终传入". Again, an assertion, not a mitigation. A mitigation would be: "Add a build-time validation that FeatureSlug is non-empty before rendering" or "Add conditional rendering `{{if .FeatureSlug}}...{{end}}`".

### Dimension 9: Success Criteria — 55/80

**Criteria are measurable and testable (24/30)**: All four criteria are objectively verifiable: presence of `## Feature Paths` in 6 embed templates, presence of `FEATURE_SLUG` line in 6 prompt templates, `go build` passes, `go test` passes. This is good. Deduction: no criterion verifies that the rendered slug is correct (non-empty, matches the actual feature directory). A template could have `## Feature Paths` with an empty slug and still pass the stated criteria.

**Coverage is complete (16/25)**: The SC covers the template changes and build stability. Missing: No SC for the prompt template's FEATURE_SLUG being non-empty at runtime. No SC for the embed template's discovery commands producing valid output. No SC for the agent's behavioral change (does it actually skip path parsing?).

**SC internal consistency (15/25)**: The SC set is internally consistent — satisfying all four is possible without contradiction. However, the SC set does not verify the proposal's core claim of "三层联动" — the criteria only check that each layer has the information, not that the layers are coordinated or that the information is consistent across layers. A scenario where embed template says `my-feature` and prompt template says `my-renamed-feature` would pass all stated SC but violate the proposal's intent.

### Dimension 10: Logical Consistency — 62/90

**Solution addresses the stated problem (28/35)**: The solution directly addresses the "declared but not rendered" defect. The prompt template change resolves the core issue. The embed template change is less well-targeted — as verified, the rich templates (`test-gen-journeys.md`, `test-gen-contracts.md`) already contain extensive `{{.FeatureSlug}}` references, so adding `## Feature Paths` to them is redundant. The solution would be more logically consistent if it distinguished between templates that genuinely lack paths and those that already have them.

**Scope <-> Solution <-> SC aligned (20/30)**: Scope lists 12 template files. Solution describes changes to those files. SC verifies the changes. The alignment is present but shallow. The SC does not verify the solution's claim of "三层联动" — there is no criterion testing that the three layers produce consistent information. Additionally, the scope excludes the `doc-summary.md` template which already has `FEATURE_SLUG: {{.FeatureSlug}}` — the proposal does not explain why this template is not included as a precedent or reference, creating an inconsistency in the template fleet's treatment.

**Requirements <-> Solution coherent (14/25)**: The key scenarios map to the solution. However, Scenario 3 ("用户独立调用 `/run-tests`（不经过 pipeline），skill 的发现逻辑仍然工作") describes a scenario that the solution explicitly does not address — skill is unchanged. This is a non-requirement dressed as a scenario. Meanwhile, no scenario addresses what happens when FeatureSlug is empty (a real edge case), and no scenario addresses the behavioral change in the agent — the "so what" of the entire proposal.

---

## Phase 3: Blindspot Hunt — What the Rubric Missed

1. **Value quantification gap**: The rubric's Problem Definition dimension checks for "evidence provided" but does not penalize proposals that solve problems without quantifying the value of the solution. This proposal admits the problem "不影响正确性" but never quantifies the efficiency gain. How many agent turns are saved? How many tokens are avoided? Without this, the proposal is a correctness fix masquerading as an efficiency improvement.

2. **Redundancy audit gap**: The rubric's Solution Clarity dimension checks "approach is concrete" but does not check whether the approach is redundant with existing solutions. The freeform review identified that `test-gen-journeys.md` and `test-gen-contracts.md` already contain extensive FeatureSlug path references. Adding `## Feature Paths` to these templates is redundant. The rubric has no mechanism to catch "solving a problem that is already solved for some of the affected artifacts."

3. **Straw-man detection gap**: The rubric's deduction rules mention "Straw-man alternative: -20 pts per instance" but the Alternatives section's "只改 embed 模板" and "改 embed + prompt + 简化 skill" options are presented with dismissive single-sentence verdicts ("联动不完整", "skill 不能简化"). The rubric does not define what constitutes a straw man with enough specificity to catch dismissiveness without outright fabrication.

4. **Temporal coupling risk**: The proposal creates a new coupling between two independently rendered templates (embed at index time, prompt at dispatch time) but the rubric's Risk Assessment dimension has no criterion for "solution-introduced coupling risks." This is a blindspot — the rubric checks for risks the proposal identifies but not for risks the proposal creates.

5. **No rollback plan**: The rubric has no criterion for rollback or reversibility. If the embed template changes produce invalid paths (empty slug), the generated task .md files are broken until the next `forge task index`. There is no rollback mechanism described. For a proposal that modifies generated artifacts, this is a significant omission.

---

## Score Summary

```
SCORE: 566/1000
DIMENSIONS:
  Problem Definition: 72/110
  Solution Clarity: 78/120
  Industry Benchmarking: 32/120
  Requirements Completeness: 60/110
  Solution Creativity: 20/100
  Feasibility: 82/100
  Scope Definition: 65/80
  Risk Assessment: 40/90
  Success Criteria: 55/80
  Logical Consistency: 62/90
ATTACKS:
1. Industry Benchmarking: Zero external references — the entire Alternatives section is an internal-only comparison table with no citation of industry patterns, open-source projects, or published approaches to context propagation in agent/pipeline systems. Quote: "无创新，纯信息补全" — while this is honest, the rubric requires industry solutions referenced (0-40 pts) and at least one industry-validated alternative. Neither is present.

2. Risk Assessment: Only two risks listed, both rated Low/Low — the freeform review identified five additional material risks including temporal coupling between embed and prompt rendering, embed empty-slug producing invalid `ls docs/features//testing/` commands, and skill-vs-embed path drift. Quote: "路径基于 forge 硬编码目录约定" — this is a convention assertion, not a mitigation. The risk table is dangerously incomplete for a change affecting 12 templates in a rendering pipeline.

3. Solution Clarity: User-facing behavior is not described — the proposal shows what information will be available but never describes what the agent does differently. Quote: "三层各有清晰职责，不重复但互相补充" — this claims no duplication, but verified against actual templates, `test-gen-journeys.md` and `test-gen-contracts.md` already contain extensive `{{.FeatureSlug}}` path references. The proposed `## Feature Paths` section is redundant for these templates.

4. Problem Definition: The urgency justification is a single sentence admitting the problem does not affect correctness. Quote: "不影响正确性（skill 能动态发现路径），但降低了 agent 执行效率。三层之间缺乏联动，路径发现逻辑仅在 skill 中，task file 和 prompt 未提供有效上下文。" — No quantification of the efficiency loss. No analysis of what happens if this is never fixed.

5. Requirements Completeness: No edge cases identified — specifically the empty-FeatureSlug scenario, the directory-rename scenario, and the non-feature-scoped-task scenario. Quote: "`{{.FeatureSlug}}` 由 CLI 在 `forge task index` 时从 `docs/features/<slug>/` 目录路径填充" — this assumes all tasks live under `docs/features/<slug>/`, which may not hold for project-level tasks.

6. Logical Consistency: The proposal applies a uniform 6-template change to a non-uniform template fleet. The thin templates (`test-run.md`, `test-gen-scripts.md`) genuinely lack path context. The rich templates (`test-gen-journeys.md`, `test-gen-contracts.md`) already contain `{{.FeatureSlug}}` references in their body. Adding `## Feature Paths` to rich templates creates redundancy without acknowledging it. Quote: "6 个测试流水线 embed 模板统一添加两类发现命令" — "统一" (uniform) is the wrong approach for a heterogeneous template set.

7. Solution Creativity: The proposal explicitly disclaims novelty. Quote: "无创新，纯信息补全" — this is a bug-fix level change rendered as a proposal. The rubric allocates 100 pts to creativity; this proposal earns minimal points here because it does not attempt any creative solution to the underlying problem of context propagation across rendering layers.

8. Success Criteria: No criterion verifies that the rendered FeatureSlug is correct or non-empty. Quote: "6 个 prompt 模板输出 `FEATURE_SLUG: <slug>` 行" — this checks for the presence of the line but not for the value being a valid slug. A template emitting `FEATURE_SLUG:` (empty) would pass this criterion.

9. Scope Definition: The scope claims "不修改任何 Go 代码" but the freeform review identified that a Go-level validation of FeatureSlug non-empty at render time would be a stronger guarantee than template-only changes. This out-of-scope exclusion may be suboptimal. Quote: "不修改任何 Go 代码（仅改模板 .md 文件）" — the proposal forecloses a higher-quality solution without justification.

10. Feasibility: The dependency on directory structure convention is unstated. Quote: "`{{.FeatureSlug}}` 由 CLI 在 `forge task index` 时从 `docs/features/<slug>/` 目录路径填充，不依赖 `.forge/state.json` 或分支名，生成时即确定" — this asserts path-based slug extraction as reliable, but does not verify what happens when the task path does not match the expected `docs/features/<slug>/tasks/` pattern.
```
