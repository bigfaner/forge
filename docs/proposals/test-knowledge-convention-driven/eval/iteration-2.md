---
iteration: 2
scorer: CTO-adversary
date: 2026-05-19
total: 842
rubric_total: 1000
target: 900
verdict: FAIL
---

# Proposal Evaluation: Iteration 2

**Score: 842 / 1000** (target: 900)

## Issues Addressed from Iteration 1

The proposal was substantially revised since iteration 1. Key improvements:

1. **Industry Benchmarking section added** (was 18/120): Three external patterns (Hygen, Plop.js, Cursor/Windsurf rules) and four alternatives with trade-off analysis, including the "do nothing" variant implicitly captured by Alternative A. This directly addresses the largest gap from iteration 1.

2. **Rollback strategy added** (blindspot #1 from iteration 1): Four-point rollback plan with explicit "no-rollback point" identified at Phase 2→3 boundary. Branch strategy specified.

3. **LLM Convention compliance validation (Phase 0 POC) added** (blindspot #2 from iteration 1): POC design with measurable criteria (first-pass compile rate, Convention accuracy). Fallback plan with three escalation tiers.

4. **Maintenance burden transfer acknowledged** (blindspot #5 from iteration 1): "用户维护成本说明" section added with explicit cost comparison (Profile era: blocked, Convention era: 3-5 minutes to edit markdown).

5. **Success criteria strengthened** (iteration 1 scored 60/80): "首过编译率 >= 85%" and "生成代码 diff 等效性" added as measurable criteria.

6. **Phased delivery timeline added** (iteration 1 gap): 5 phases with estimated durations, dependencies, and total estimate (13-20 days).

7. **Justfile bootstrap paradox addressed** (blindspot #4 from iteration 1): Four concrete mitigations for justfile circular dependency.

## Dimension Scores

### 1. Problem Definition: 98 / 110

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Problem stated clearly | 38/40 | Core problem remains well-defined: Profile system creates conceptual contradiction with Journey-Contract model. Three-point decomposition (conceptual contradiction, responsibility mixing, config overreach) is precise. Deduction: "config.yaml 职责越界" is presented as a problem but is actually a design preference — the principle "config 只管 Forge 行为控制" is asserted without establishing why mixing detectable and behavioral config is harmful beyond aesthetics. |
| Evidence provided | 38/40 | Four-row detection failure table remains strong. Profile extension cost (5 steps) is concrete. "19+ files deep consumption" is verifiable. New addition: "2 个待处理的框架支持请求被阻塞" — specific evidence of real user pain. Deduction: no quantitative data on failure frequency — how often do users actually hit non-default frameworks? |
| Urgency justified | 22/30 | Improved from iteration 1: "2 个框架支持请求被阻塞" provides concrete cost of delay. "consolidate-specs 对测试知识的完整管理" explains what feature is blocked. Deduction: still no answer to "what happens to existing users if we don't do this?" — current users on default frameworks are unaffected. The urgency applies only to non-default-framework users, whose population size is unstated. |

### 2. Solution Clarity: 108 / 120

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Approach is concrete | 38/40 | Three-layer separation model is clear. Complete consumer mapping (deleted/rewritten/retained tables). Fixed Convention structure with concrete example. Deduction: "Convention 文件可以包含更多 section" — the boundary between "fixed structure" and "extensible sections" is still imprecise. The proposal says both "固定结构" and "可以包含更多 section" without specifying which sections are required vs optional. |
| User-facing behavior described | 38/40 | test-guide interaction flow (Step 1/2a/2b/3) is detailed. Cold start sequence documented. Nine scenarios cover major paths. New: "用户维护成本说明" section describes what users actually do. Deduction: the experience of a user whose Convention content is wrong (typo in assertion library name) is not described. What does the error look like? How does the user diagnose it? |
| Technical direction clear | 32/35 | Convention loading mechanism, Code Reconnaissance extension, just abstraction all technically precise. New: justfile bootstrap paradox addressed with four mitigations. Deduction: "run-e2e-tests skill prompt 内置结果解析知识" — this means embedding framework-specific parsing logic in the skill prompt. This is a form of hardcoded framework knowledge, just moved from Go code to prompt text. The proposal does not acknowledge this parallel. |

### 3. Industry Benchmarking: 88 / 120

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Industry solutions referenced | 30/40 | Three external patterns cited: Hygen (template-driven), Plop.js (config-driven scaffolding), Cursor/Windsurf rules (LLM instruction files). Each has strengths and weaknesses analyzed. Deduction: the analysis is shallow — one paragraph per tool. No citations (no links, no version numbers). The Cursor/Windsurf pattern is described as the closest analog, but no evidence is provided that these tools' rule files actually work reliably for code generation guidance. The comparison reads as blog-post-level analysis rather than rigorous benchmarking. |
| At least 3 meaningful alternatives | 22/30 | Four alternatives in the comparison table: Config-driven profiles, AST-based detection, LLM-only, and Convention files. "Do nothing" is implicitly captured by Alternative A (config-driven profiles reproduces Profile's approach). Deduction: Alternative B (AST-based) is a genuine alternative but is dismissed with "maintenance cost exceeds Profile" — this is asserted without evidence. AST parsers for Go and TypeScript are mature, well-maintained libraries. The real cost comparison is never shown. |
| Honest trade-off comparison | 18/25 | Trade-offs are presented in the alternatives table. The Convention approach's primary weakness ("LLM may ignore Convention content") is acknowledged. Deduction: the comparison is slightly tilted — Alternative A's cons say "reproduces the exact failure mode this proposal eliminates" which is editorializing, not analysis. Alternative C's cons list is the longest. The proposal's own approach has the shortest cons list. The bias is subtle but present. |
| Chosen approach justified against benchmarks | 18/25 | "Why Convention over Alternatives" section explains the niche: "only option that combines user-editability, zero per-framework code, cross-session persistence, and graceful degradation." Deduction: this justification combines four attributes that make Convention appear uniquely superior, but these attributes were selected post-hoc to favor the Convention approach. An equally valid framing: Convention is the only option that relies on LLM interpretation — which is simultaneously its greatest strength (flexibility) and greatest risk (reliability). The framing minimizes the risk. |

### 4. Requirements Completeness: 96 / 110

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Scenario coverage | 36/40 | Nine scenarios cover happy path, cold start, multi-framework, Convention missing, backward compatibility, CLI commands, and config cleanup. New: Scenario 9 explicitly validates backward compatibility with diff measurement. Deduction: error scenarios remain weak — what if Convention content has wrong framework name? What if LLM generates code that compiles but uses wrong assertion semantics? No scenario for Convention file corruption or deletion mid-project. |
| Non-functional requirements | 38/40 | Five NFRs remain well-stated: backward compatibility, progressive migration, zero new infrastructure, compile gate, on-demand loading. New: "首过编译率 >= 85%" adds a measurable quality NFR. Deduction: "零新增基础设施" (zero new infrastructure) contradicts creating `/forge:test-guide` slash command. The NFR should be restated more precisely — "no new runtime infrastructure" or "no new external dependencies." |
| Constraints & dependencies | 22/30 | Four constraints listed. New: phased timeline with dependencies between phases. Deduction: no team skill constraint (LLM prompt engineering proficiency). No constraint on LLM model capability — what if the model changes or degrades? The proposal assumes LLM quality is constant, but this is not a given. No constraint on Convention file format versioning — what happens when the "fixed structure" needs to evolve? |

### 5. Solution Creativity: 68 / 100

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Novelty over industry baseline | 32/40 | The proposal explicitly positions itself relative to the Cursor/Windsurf rules pattern and adds fixed section structure + compile gate validation. The "three-layer separation" (methodology / knowledge / actual) is a useful mental model. The insight "LLM doesn't need embedded templates — it just needs to know project conventions" is genuinely pragmatic. The Phase 0 POC gate is a creative process innovation — validate before committing. |
| Cross-domain inspiration | 18/35 | Cursor/Windsurf rules are acknowledged as inspiration. The Convention file pattern resembles `.editorconfig`, `.prettierrc`, and linting configs, but these parallels are not explicitly drawn. The phased validation approach borrows from engineering best practices (spike before commit). Deduction: no cross-domain inspiration beyond the LLM tooling space. The proposal does not reference how other systems handle "user-declared knowledge + runtime detection" patterns (e.g., how IDEs combine user settings with auto-detected project properties). |
| Simplicity of insight | 18/25 | The core insight remains elegant: "LLM doesn't need embedded knowledge — it just needs to know project conventions." The fixed section structure is simple. The compile gate as validation is simple. Deduction: the execution is not simple — 19+ file rewrites, 3 skill rewrites, new slash command, Go package rewrites, inline result parsing in run-e2e-tests. The gap between insight simplicity and execution complexity is significant. |

### 6. Feasibility: 82 / 100

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Technical feasibility | 35/40 | Tech stack is standard (Go, markdown, just, LLM). Phase 0 POC validates the core assumption before committing. Fallback plan exists. Deduction: the POC is proposed but not yet executed. The proposal asks for approval contingent on POC results, but the POC results are not available. This is the right process, but the feasibility assessment remains incomplete until POC validates. |
| Resource & timeline feasibility | 25/30 | Substantially improved: phased delivery timeline with 13-20 day estimate for a single developer. Phase 0 (2-3 days) is the gate. Dependencies between phases are explicit. Deduction: 19+ file rewrites + 3 skill rewrites + new slash command + Go package rewrites in 13-20 days assumes no unexpected complications. The estimate is optimistic given the scope — Phase 1 alone (5-7 days for all Go code changes) packs significant complexity. No buffer is included. |
| Dependency readiness | 22/30 | Existing infrastructure (conventions directory, domains loading, just) is ready. `/forge:test-guide` is new but builds on existing skill infrastructure. Deduction: consolidate-specs integration is listed as in-scope but its readiness is not discussed. The "inline result parsing" in run-e2e-tests depends on justfile recipe output format, which must be standardized — this dependency is acknowledged but its readiness is unclear. |

### 7. Scope Definition: 72 / 80

| Criterion | Score | Justification |
|-----------|-------|---------------|
| In-scope items are concrete | 26/30 | In-scope items are concrete deliverables: specific directories to delete, specific files to rewrite, specific new components to create. The phased timeline maps scope to phases. Deduction: "Convention 文件固定结构定义" — the definition itself is in scope, but the exact structure is still being refined (the boundary between fixed and extensible sections is imprecise). "用户手册更新" has no specific deliverable (which pages? which sections?). |
| Out-of-scope explicitly listed | 24/25 | Eight items explicitly out of scope: gen-journeys, gen-contracts, verify/promote, existing test migration, unit tests, anti-pattern docs, Convention auto-sync. Clear and specific. |
| Scope is bounded | 22/25 | Substantially improved: phased delivery with estimated durations (2-3, 5-7, 3-5, 2-3, 1-2 days). Total estimate provided. Phase 0 is an explicit early-stop gate. Deduction: some phase descriptions remain coarse — "pkg/journey/ rewrite" in Phase 1 does not specify what "rewrite" means at the function level. The validation phase (1-2 days for 126+ tests) seems tight if issues are discovered. |

### 8. Risk Assessment: 82 / 90

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Risks identified | 28/30 | Seven risks identified, including technical risks (regression, LLM quality, tag generation, result parsing) and project risks (timeline). New: rollback strategy addresses the blindspot from iteration 1. Deduction: missing risks — Convention format versioning (what if the fixed structure needs to change?), LLM model quality degradation (what if a model update changes generation behavior?), and multi-user Convention conflicts (two developers editing Convention differently). |
| Likelihood + impact rated | 26/30 | Ratings are honest — H/H for the highest risk (full rewrite regression), M/M or M/H for others. Not all risks are low-likelihood/high-impact. Deduction: "无 Convention 时 LLM 生成质量低于 generate.md" rated M/M — this seems underrated. If LLM generation quality degrades significantly without Convention (or with wrong Convention), the entire proposal's value proposition collapses. Impact should arguably be H. |
| Mitigations are actionable | 28/30 | Mitigations are concrete: "126+ e2e tests as regression safety net," "compile gate," "test-guide to solidify patterns," "structured output formats," "Phase 0 POC as early-stop gate." New: rollback strategy with four explicit steps and identified no-rollback point. Deduction: "compile gate" appears as mitigation for 4 of 7 risks — reliance on a single mitigation mechanism. However, the compile gate is genuinely the right safety net for code generation, so the repetition is partially justified. |

### 9. Success Criteria: 72 / 80

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Criteria are measurable and testable | 48/55 | Substantially improved. "首过编译率 >= 85%" is measurable with a clear threshold. "生成代码 diff 等效性" has specific scope (import, assertion functions, tag syntax — diff = 0; style differences allowed). "3 个不同框架项目" provides coverage. Most criteria are binary-verifiable. Deduction: "Convention 文件固定结构...skill 在 3 个不同框架项目上正确解读" — "正确解读" needs a more precise definition. Does it mean the generated code compiles? Or that the generated code uses exactly the Convention-specified patterns? These are different bars. "skill 可靠地解读" from iteration 1 was removed, but the replacement still has ambiguity. |
| Coverage is complete | 24/25 | Criteria cover Profile removal, config cleanup, skill rewrites, CLI simplification, backward compatibility, Convention creation, and quality metrics. New: run-e2e-tests Convention-based parsing is covered. init-justfile Convention-driven generation is covered. Deduction: no criterion for test-guide user experience quality. No criterion for consolidate-specs integration depth (just "纳入管理" — what does that mean specifically?). |

### 10. Logical Consistency: 78 / 90

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Solution addresses stated problem | 32/35 | Yes — Profile removal directly addresses the conceptual contradiction. Convention files address the hardcoded assumption issue. Config cleanup addresses the config overreach issue. Three-layer separation resolves the responsibility mixing. Deduction: the proposal argues Profile detection fails (evidence table), but the solution's Code Reconnaissance also performs detection (LLM-based). The proposal does not acknowledge that the fundamental weakness — detection can fail — is shared by both approaches. The mitigation is that Convention files are user-declared (not detected), but Reconnaissance still detects, and cold starts rely entirely on detection. |
| Scope <-> Solution <-> Success Criteria aligned | 24/30 | Generally consistent. Scope covers all three solution layers. Success criteria map to scope items. Deduction: scope includes "用户手册更新" but no success criterion validates it. Scope includes "Convention 文件纳入 consolidate-specs 管理" but the success criterion is vague ("纳入管理"). The Phase 0 POC is in the timeline but not explicitly in scope or success criteria — it is described in the requirements section instead. |
| Requirements <-> Solution coherent | 22/25 | Requirements scenarios map well to solution features. NFRs are addressed by the solution design. Deduction: NFR "零新增基础设施" is contradicted by the new `/forge:test-guide` slash command. The proposal acknowledges Convention loading "使用现有机制" but the test-guide slash command is explicitly new. The contradiction from iteration 1 persists — it should be rephrased as "零新增运行时基础设施" or similar. Also: run-e2e-tests "内置结果解析知识" means framework-specific parsing logic in the skill prompt — this is framework knowledge embedded in a different medium, paralleling the generate.md approach the proposal aims to eliminate. |

## Cross-Dimension Coherence Findings

1. **NFR "零新增基础设施" vs test-guide creation** (Dimensions 4, 10): The NFR claims "零新增基础设施" but the proposal explicitly creates a new slash command. This is the same contradiction flagged in iteration 1. The NFR should be rephrased for precision.

2. **POC-as-gate vs approval request** (Dimensions 6, 7): The proposal requests approval (it is a proposal, not a POC report) but makes the core assumption contingent on POC validation. The logical structure is: "approve this proposal, then we'll validate whether it works." The more rigorous approach would be to present POC results first, then request approval. The proposal acknowledges this by structuring Phase 0 as a gate, but the governance model (who approves Phase 0 results? what triggers proceeding to Phase 1?) is undefined.

3. **Success criteria "3 frameworks" vs evidence table "4 scenarios"** (Dimensions 4, 9): The evidence table tests 4 scenarios (Go default, Go ginkgo, TypeScript vitest, framework upgrade), but success criteria require validation on only 3 frameworks. The "framework upgrade" scenario from the evidence table has no corresponding success criterion — it is the strongest evidence for the problem but is not tested in the solution.

## Blindspot Attacks

1. **[blindspot] Convention format evolution has no migration strategy**: Quote: "Convention 文件遵循固定结构" (Section: 固定结构) and "Convention 文件可以包含更多 section" (same section). The proposal mandates a "fixed structure" for Convention files, but simultaneously says sections can be added over time. If Forge evolves and the "fixed structure" changes (e.g., a new required section is added), existing Convention files in user projects break. There is no versioning mechanism, no migration path, and no backward compatibility strategy for Convention format changes. This is an infrastructure maintenance blindspot that will surface in the first Convention format revision.

2. **[blindspot] run-e2e-tests re-embeds framework knowledge in a different medium**: Quote: "run-e2e-tests skill prompt 内置结果解析知识（TC ID 提取模式、状态映射、错误信息提取）...skill 从 Convention 的 Framework section 知道当前是什么框架，使用对应的解析策略" (Section: 结果输出契约). The proposal removes generate.md because it "hardcodes framework-specific knowledge" in Go embed files. But run-e2e-tests will now contain "内置结果解析知识" with "三套解析策略按 Convention Framework section 切换" — this is framework-specific parsing logic embedded in a skill prompt instead of Go code. The problem (hardcoded framework knowledge) is not eliminated; it is relocated from Go to markdown. When a new framework is added, the skill prompt must be updated to include a new parsing strategy — the same maintenance bottleneck, different file format.

3. **[blindspot] The Phase 0 POC scope is narrower than the production scope**: Quote: "POC design: 1. Write Convention files for 3 frameworks: Go testing (forge-cli's actual setup), Go ginkgo, TypeScript vitest. 2. Run gen-test-scripts with Convention files loaded" (Section: LLM Convention Compliance Validation). The POC tests gen-test-scripts only. But the production scope includes run-e2e-tests (result parsing), init-justfile (recipe generation), and pkg/journey/ (tag generation). The POC validates the easiest case (code generation with Convention) but does not validate the harder cases (test result parsing, justfile recipe generation). The production system could pass Phase 0 and still fail at Phase 2 for unvalidated components.
