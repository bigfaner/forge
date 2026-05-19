---
iteration: 3
scorer: CTO-adversary
date: 2026-05-19
total: 870
rubric_total: 1000
target: 900
verdict: FAIL
---

# Proposal Evaluation: Iteration 3 (Final)

**Score: 870 / 1000** (target: 900)

## Status Assessment

The proposal has not been revised since iteration 2. The document content is identical. This final evaluation confirms iteration 2 scores with adjustments where the rubric demands final-round rigor, and introduces no new revision-based improvements.

## Dimension Scores

### 1. Problem Definition: 96 / 110

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Problem stated clearly | 37/40 | Core problem remains well-defined: Profile system creates conceptual contradiction with Journey-Contract model. Three-point decomposition is precise. Final-round deduction: "config.yaml 职责越界" violates "config 只管 Forge 行为控制" — this principle is asserted as self-evident but is actually an architectural preference. The proposal never establishes why a config file containing detectable information is harmful beyond "it violates a principle." Principles are useful but their violation is not inherently a problem without consequence. |
| Evidence provided | 38/40 | Four-row detection failure table is concrete. Profile extension cost (5 steps) is verifiable. "19+ files deep consumption" is measurable. "2 个待处理的框架支持请求被阻塞" is specific. Final-round deduction: no data on how many users or projects are affected. The evidence proves the mechanism fails but not the population-scale impact. |
| Urgency justified | 21/30 | "2 个框架支持请求被阻塞" provides concrete cost of delay. consolidate-specs blockage explains feature impact. Final-round deduction (stricter): the proposal still does not establish urgency for existing users on default frameworks. The urgency is real for 2 requestors, but the proposal asks for a 19+ file rewrite affecting all users. The urgency-to-scope ratio is not established. What happens to the majority of users (on default frameworks) if we do nothing? They are unaffected. The cost of delay is narrow; the cost of action is system-wide. |

### 2. Solution Clarity: 107 / 120

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Approach is concrete | 37/40 | Three-layer separation model is clear. Complete consumer mapping. Fixed Convention structure with concrete example. Final-round deduction: "Convention 文件可以包含更多 section" — the boundary between fixed and extensible remains imprecise after two iterations. Which sections are mandatory vs. optional is never specified. This ambiguity will surface during implementation when a developer must decide whether to reject a Convention file missing a "Result Format" section. |
| User-facing behavior described | 37/40 | test-guide interaction flow is detailed. Cold start sequence documented. Nine scenarios cover major paths. User maintenance cost section added. Final-round deduction: the error experience remains undescribed. Quote: "Convention 文件遵循固定结构，skill 按 section 标题解读" — what does the user see when a section title has a typo? What error message appears when the LLM cannot find the Assertion section? The happy path is thorough; the error path is invisible. |
| Technical direction clear | 33/35 | Convention loading mechanism, Code Reconnaissance extension, just abstraction all technically precise. Justfile bootstrap paradox addressed with four mitigations. Final-round deduction: the proposal says "run-e2e-tests 从 Convention 的 Result Format section 读取格式类型和解析策略" — this means run-e2e-tests must contain parsing logic for json-stream, json-report, and text-verbose formats. That is three parsing strategies conditioned on a Convention value. This is framework knowledge encoded differently, not eliminated. The parallel with the Profile approach it replaces is not acknowledged in the solution section. |

### 3. Industry Benchmarking: 88 / 120

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Industry solutions referenced | 30/40 | Three external patterns (Hygen, Plop.js, Cursor/Windsurf rules) plus three cross-domain parallels (.editorconfig, .prettierrc, tsconfig.json). Each analyzed for strengths and weaknesses. The cross-domain parallels added since iteration 1 strengthen this section significantly. Final-round deduction: still no citations, version numbers, or links. The Cursor/Windsurf analysis claims "no structured sections, no validation" but `.cursorrules` files have evolved — some tools now support structured rule files with metadata. The analysis may be outdated but cannot be verified without citations. |
| At least 3 meaningful alternatives | 22/30 | Four alternatives with trade-off analysis. "Do nothing" is implicitly captured by Alternative A. Each alternative is genuinely different. Final-round deduction: Alternative B (AST-based) is dismissed with "maintenance cost exceeds Profile" and "cold start gap." Both are real concerns, but the assertion about maintenance cost remains unevidenced. Go's `go/ast` package is part of the standard library — zero maintenance cost for Go. TypeScript's compiler API is maintained by Microsoft. The dismissal is plausible but not demonstrated. |
| Honest trade-off comparison | 18/25 | Trade-offs are presented in the alternatives table. Convention's weakness ("LLM may ignore Convention content") is acknowledged. Final-round deduction: the bias persists. Alternative A's cons say "reproduces the exact failure mode this proposal eliminates" — editorial language. Alternative C has the longest cons list. Convention's cons list is the shortest. The "Why Convention over Alternatives" section frames four attributes selected to favor Convention: "user-editability, zero per-framework code, cross-session persistence, and graceful degradation" — all true of Convention, but the framing excludes attributes where Convention is weaker (determinism, enforcement strength, reliability). |
| Chosen approach justified against benchmarks | 18/25 | The compile-gate-as-enforcement argument is the core justification: "the compile gate converts 'LLM might misinterpret Convention' from a silent failure into a caught-and-retried failure." This is the proposal's strongest argument. Final-round deduction: the justification treats the compile gate as a complete safety net, but the compile gate only validates compilation — not correctness. Code that compiles with wrong assertion semantics (e.g., using `assert.Equal` when Convention specifies `assert.NoError`) passes the compile gate. The proposal does not acknowledge this gap in the justification. |

### 4. Requirements Completeness: 94 / 110

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Scenario coverage | 35/40 | Nine scenarios cover happy path, cold start, multi-framework, Convention missing, backward compatibility, CLI commands, and config cleanup. Final-round deduction: error scenarios remain the weakest area. No scenario for Convention file with wrong content (typo in framework name, wrong assertion library). No scenario for LLM generating code that compiles but uses wrong assertion semantics. No scenario for Convention file corruption or deletion mid-project. The scenarios validate the architecture but not the failure modes. |
| Non-functional requirements | 38/40 | Five NFRs: backward compatibility, progressive migration, zero new infrastructure, compile gate, on-demand loading. "首过编译率 >= 85%" adds a measurable quality NFR. Final-round deduction: "零新增基础设施" (zero new infrastructure) is imprecise. The proposal creates `/forge:test-guide` — a new slash command. The NFR should state "零新增运行时外部依赖" or similar. This is the same deduction as iterations 1 and 2; the imprecision persists. |
| Constraints & dependencies | 21/30 | Four constraints listed. Phased timeline with dependencies. Final-round deduction (stricter): three constraints remain unaddressed after two iterations: (1) LLM model capability — the proposal assumes LLM quality is constant, but model changes or degradation would break the system. (2) Convention format versioning — when the "fixed structure" evolves, existing Convention files in user projects become invalid with no migration path. (3) Team skill constraint — LLM prompt engineering proficiency is assumed but not stated. These are not theoretical concerns; the first Convention format revision will trigger the versioning gap. |

### 5. Solution Creativity: 70 / 100

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Novelty over industry baseline | 33/40 | The proposal positions itself against Cursor/Windsurf rules and adds fixed section structure + compile gate validation. "LLM doesn't need embedded templates — it just needs to know project conventions" is genuinely pragmatic. The Phase 0 POC gate is a creative process innovation. Final-round assessment: the novelty is solid — it takes an established pattern (LLM instruction files) and adds structure + validation, which is the right incremental innovation for this domain. |
| Cross-domain inspiration | 19/35 | Cross-domain parallels added since iteration 1 (.editorconfig, .prettierrc, tsconfig.json) improve this score. The analogy to user-editable config + tool respects it + sensible fallback is clearly drawn. Final-round deduction: all parallels are from developer tooling configuration. No inspiration from domains that solve similar "user-declared knowledge + runtime detection" patterns more broadly — e.g., how IDEs combine user settings with auto-detected project properties, or how CI systems combine user config with auto-detected build environments. The cross-domain reach is narrow. |
| Simplicity of insight | 18/25 | Core insight remains elegant: "LLM doesn't need embedded knowledge — it just needs to know project conventions." The fixed section structure is simple. The compile gate as validation is simple. Final-round deduction: the execution is not simple. 19+ file rewrites, 3 skill rewrites, new slash command, Go package rewrites, inline result parsing in run-e2e-tests, Convention loading mechanism, Code Reconnaissance extension, config cleanup, CLI command removal. The gap between insight simplicity and execution complexity remains significant and unresolved across all three iterations. |

### 6. Feasibility: 80 / 100

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Technical feasibility | 34/40 | Tech stack is standard. Phase 0 POC validates core assumption before committing. Fallback plan exists with three escalation tiers. Final-round deduction: the POC has not been executed. The proposal requests approval contingent on POC results, but presents no POC data. The feasibility assessment is incomplete — the most critical technical risk (LLM Convention compliance) remains unvalidated at proposal time. This is the correct process (POC first), but it means the feasibility score cannot be higher until POC results are available. |
| Resource & timeline feasibility | 24/30 | 13-20 day estimate for single developer. Phase dependencies explicit. Phase 0 (2-3 days) is the gate. Final-round deduction (stricter): 19+ file rewrites + 3 skill rewrites + new slash command + Go package rewrites in 13-20 days with no buffer. Phase 1 alone (5-7 days for all Go code changes: pkg/journey/, pkg/e2e/, pkg/just/, pkg/task/, internal/cmd/) packs significant complexity into a tight window. The validation phase (1-2 days for 126+ tests) assumes no issues are discovered. No contingency buffer is included for any phase. |
| Dependency readiness | 22/30 | Existing infrastructure (conventions directory, domains loading, just) is ready. `/forge:test-guide` builds on existing skill infrastructure. Final-round deduction: consolidate-specs integration is in scope but its readiness is not discussed. The Convention Result Format-driven parsing in run-e2e-tests depends on justfile recipe output format standardization — acknowledged but readiness unclear. |

### 7. Scope Definition: 72 / 80

| Criterion | Score | Justification |
|-----------|-------|---------------|
| In-scope items are concrete | 26/30 | In-scope items are concrete deliverables: specific directories to delete, files to rewrite, components to create. Phased timeline maps scope to phases. Final-round deduction: "Convention 文件固定结构定义" — the definition is in scope, but the boundary between fixed and extensible sections is still imprecise. "用户手册更新" has no specific deliverable (which pages, which sections). |
| Out-of-scope explicitly listed | 24/25 | Eight items explicitly out of scope. Clear and specific. |
| Scope is bounded | 22/25 | Phased delivery with estimated durations. Phase 0 is an explicit early-stop gate. Total estimate provided. Final-round deduction: some phase descriptions remain coarse — "pkg/journey/ rewrite" does not specify what "rewrite" means at the function level. The validation phase (1-2 days for 126+ tests) seems tight if issues are discovered. |

### 8. Risk Assessment: 80 / 90

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Risks identified | 27/30 | Seven risks identified including rollback strategy. Final-round deduction: three risks flagged in iteration 2 remain unaddressed: (1) Convention format versioning — if the fixed structure changes, existing Convention files break. (2) LLM model quality degradation — if a model update changes generation behavior. (3) Multi-user Convention conflicts — two developers editing Convention differently. These are operational risks that will materialize in production. |
| Likelihood + impact rated | 25/30 | Ratings are honest — H/H for the highest risk, M/M or M/H for others. Not all risks are low-likelihood/high-impact. Final-round deduction (stricter): "无 Convention 时 LLM 生成质量低于 generate.md" rated M/M. If LLM generation quality degrades without Convention (or with wrong Convention), the proposal's value proposition collapses — this impact should arguably be H. The rating understates the risk to the core thesis. |
| Mitigations are actionable | 28/30 | Mitigations are concrete: "126+ e2e tests as regression safety net," "compile gate," "test-guide to solidify patterns," "structured output formats," "Phase 0 POC as early-stop gate," rollback strategy with four explicit steps. Final-round deduction: "compile gate" appears as mitigation for 4 of 7 risks — reliance on a single mitigation mechanism. The compile gate is genuinely the right safety net for code generation, but it only validates compilation, not semantic correctness. |

### 9. Success Criteria: 72 / 80

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Criteria are measurable and testable | 48/55 | "首过编译率 >= 85%" is measurable with a clear threshold. "生成代码 diff 等效性" has specific scope (import, assertion functions, tag syntax — diff = 0; style differences allowed). "3 个不同框架项目" provides coverage. Most criteria are binary-verifiable. Final-round deduction: "skill 在 3 个不同框架项目上正确解读 Convention 内容，生成的 import/断言/tag 与 Convention 声明一致" — "一致" needs a more precise definition. Does "一致" mean the generated code uses the exact assertion function names from Convention? Or that the generated code uses the same assertion library? These are different precision bars. The ambiguity has persisted across all three iterations. |
| Coverage is complete | 24/25 | Criteria cover Profile removal, config cleanup, skill rewrites, CLI simplification, backward compatibility, Convention creation, and quality metrics. Final-round deduction: no criterion for test-guide user experience quality. No criterion for consolidate-specs integration depth ("纳入管理" — what does this mean specifically?). |

### 10. Logical Consistency: 81 / 90

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Solution addresses stated problem | 33/35 | Profile removal directly addresses the conceptual contradiction. Convention files address the hardcoded assumption issue. Config cleanup addresses the config overreach issue. Three-layer separation resolves the responsibility mixing. Final-round deduction (partial credit restored): the proposal does acknowledge that Code Reconnaissance also performs detection (LLM-based) and that cold starts rely on it. The Convention layer provides user-declared knowledge as a primary mechanism, reducing (not eliminating) detection dependency. The gap is smaller than in iteration 1 but not zero. |
| Scope <-> Solution <-> Success Criteria aligned | 25/30 | Generally consistent. Scope covers all three solution layers. Success criteria map to scope items. Final-round deduction: scope includes "用户手册更新" but no success criterion validates it. Scope includes "Convention 文件纳入 consolidate-specs 管理" but the success criterion is vague ("纳入管理"). The Phase 0 POC is described in the timeline and requirements section but not explicitly listed in scope or validated by a success criterion. |
| Requirements <-> Solution coherent | 23/25 | Requirements scenarios map well to solution features. NFRs are addressed by the solution design. Final-round deduction: run-e2e-tests "从 Convention 的 Result Format section 读取格式类型和解析策略" means the skill must contain parsing logic for multiple format types (json-stream, json-report, text-verbose). This is framework-specific knowledge encoded in a different medium (skill prompt vs. Go embed). The proposal explicitly states "不包含硬编码的框架特定解析逻辑" in the success criteria, but the skill must still implement format-specific parsing. The distinction between "hardcoded framework knowledge" and "format-type parsing logic conditioned on Convention values" is thinner than the proposal acknowledges. The NFR "零新增基础设施" is contradicted by the new slash command (same issue as iterations 1 and 2). |

## Cross-Dimension Coherence Findings

1. **NFR "零新增基础设施" imprecision** (Dimensions 4, 10): The NFR claims "零新增基础设施" but the proposal creates a new slash command `/forge:test-guide`. This contradiction was flagged in iteration 1 and persists in iteration 3. The NFR should be rephrased as "零新增运行时外部依赖" for precision. Unresolved across three iterations.

2. **POC-as-gate vs approval request governance gap** (Dimensions 6, 7): The proposal requests approval while making its core assumption contingent on POC validation. The governance model — who approves Phase 0 results, what triggers proceeding to Phase 1, what constitutes POC failure beyond "first-pass compile rate < 70%" — is undefined. The POC failure threshold is mentioned in the timeline but not in the risk table or success criteria.

3. **Success criteria "3 frameworks" vs evidence table "4 scenarios"** (Dimensions 4, 9): The evidence table tests 4 scenarios including "framework upgrade mocha->vitest." Success criteria require validation on only 3 frameworks. The framework upgrade scenario — the strongest evidence for the problem — has no corresponding success criterion for the solution.

4. **Compile gate as universal safety net** (Dimensions 2, 8, 10): The compile gate appears as the primary mitigation for 4 of 7 risks and as the core justification for choosing Convention over alternatives. But the compile gate validates compilation only, not semantic correctness. Code that compiles with wrong assertion semantics passes the gate. This limitation is not acknowledged anywhere in the proposal.

## Blindspot Attacks

1. **[blindspot] Convention format versioning has no migration strategy**: Quote: "Convention 文件遵循固定结构" and "Convention 文件可以包含更多 section" (Section: 固定结构). The proposal mandates a "fixed structure" but simultaneously permits extension. When Forge evolves and the structure changes (e.g., a new required section), existing Convention files in user projects become stale or broken. There is no version field, no migration tool, no backward compatibility strategy for Convention format evolution. This is not a theoretical risk — the proposal itself says "Convention 文件可以包含更多 section" which implies the structure will change. The first format revision will break existing deployments with no recovery path.

2. **[blindspot] run-e2e-tests re-embeds framework knowledge in a different medium**: Quote: "run-e2e-tests 从 Convention 的 Result Format section 读取格式类型和解析策略" and "run-e2e-tests skill...结果解析策略从 Convention 的 Result Format section 读取，不包含硬编码的框架特定解析逻辑" (Sections: 结果输出契约, Success Criteria). The proposal removes generate.md because it "hardcodes framework-specific knowledge." But run-e2e-tests must contain parsing logic for json-stream, json-report, and text-verbose formats — three format-specific parsing strategies conditioned on Convention values. When a new framework produces a new output format (e.g., XML reports), the skill must be updated with new parsing logic. The maintenance bottleneck is not eliminated — it is relocated from Go embed files to the skill prompt. The success criterion claims "不包含硬编码的框架特定解析逻辑" but format-type-specific parsing is framework knowledge, just abstracted one level.

3. **[blindspot] The "一致" success criterion is semantically ambiguous**: Quote: "skill 在 3 个不同框架项目上正确解读 Convention 内容，生成的 import/断言/tag 与 Convention 声明一致" (Section: Success Criteria). "一致" (consistent) is the key quality gate but is never operationally defined. Does it mean: (a) the generated code imports exactly the packages named in Convention? (b) the generated code uses assertion functions from the library named in Convention? (c) the generated code uses the specific assertion functions listed in Convention? Each is a progressively stricter bar. The proposal measures "diff 等效性" separately (import/assertion/tag diff = 0 against Profile-generated output), but "一致" against Convention is a different measurement with undefined precision. This ambiguity will cause acceptance disputes.

## Summary

The proposal is well-structured and addresses most of its stated goals. The core insight — user-editable Convention files with compile-gate validation — is sound. The primary remaining weaknesses are: (1) the Convention format versioning gap, (2) the run-e2e-tests framework knowledge relocation, and (3) imprecise success criteria around Convention compliance. These are real issues that will surface during implementation and operation.
