---
iteration: 1
scorer: CTO-adversary
date: 2026-05-19
total: 692
rubric_total: 1000
target: 900
verdict: FAIL
---

# Proposal Evaluation: Iteration 1

**Score: 692 / 1000** (target: 900)

## Dimension Scores

### 1. Problem Definition: 89 / 110

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Problem stated clearly | 32/40 | Core problem is unambiguous: two coexisting models (Journey-Contract vs Profile) create conceptual contradictions. The three-point decomposition (conceptual contradiction, responsibility mixing, config overreach) is well-defined. Minor deduction: "conceptual contradiction" is stated at an abstract level; two readers might disagree on whether mixing concerns truly constitutes a contradiction vs. a manageable trade-off. |
| Evidence provided | 35/40 | The four-row detection failure table is concrete and verifiable. The Profile extension cost analysis (5 steps to add a framework) is specific. "19+ files deep consumption" is a measurable claim. Deduction: no empirical data on how often these failures occur in practice or their user impact. The claim "generate.md is hardcoded assumption" is logically argued but lacks real incident data (e.g., "we observed N instances of generated code failing compilation due to wrong assertion library"). |
| Urgency justified | 22/30 | Implicit urgency from "conceptual contradiction" and accumulating tech debt. Deduction: no answer to "what happens if we don't?" or "what's the cost of delay?" The proposal doesn't establish whether the Profile system is actively blocking features, causing user pain, or just aesthetically displeasing to the architects. |

### 2. Solution Clarity: 103 / 120

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Approach is concrete | 35/40 | Three-layer separation model is clear. Every consumer is mapped to delete/rewrite/retain. The tables for deleted components, rewritten components, and retained components are exhaustive. Deduction: "Convention files follow fixed structure" but then "Convention files can contain more sections" — the boundary between fixed and extensible is imprecise. |
| User-facing behavior described | 38/40 | test-guide interaction flow is detailed (Step 1/2a/2b/3). Cold start sequence is documented. Nine scenarios cover major user-facing paths. Deduction: the actual user experience of Convention files — editing them, debugging wrong Convention content, what happens when generation fails due to bad Convention — is not well described. The proposal focuses on internal mechanisms over user experience. |
| Technical direction clear | 30/35 | Three-layer separation, Convention loading mechanism, just abstraction, Code Reconnaissance extension all have technical precision. Deduction: "LLM reads markdown, not programmatic parsing" — this is a critical architectural decision (relying on LLM parsing vs programmatic) without robustness discussion. What happens if the LLM misinterprets Convention structure? No fallback described. |

### 3. Industry Benchmarking: 18 / 120

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Industry solutions referenced | 8/40 | No external solutions, open-source projects, or published patterns are cited. The proposal self-invents the "Convention file" approach without referencing any existing work on test framework abstraction, configuration-as-code, or code generation knowledge management. This is a significant gap for a well-documented problem domain. |
| At least 3 meaningful alternatives | 10/30 | No explicit alternatives analysis. The proposal implicitly considers "do nothing" (keep Profile) but does not structurally evaluate it. Missing alternatives: config-driven profiles (user selects framework in config), schema-based generation (structured templates), AST-based code generation (parsing existing tests), LLM-only approach (no Convention files at all). |
| Honest trade-off comparison | 0/25 | No trade-off comparison section exists. |
| Chosen approach justified against benchmarks | 0/25 | No benchmark justification exists. |

### 4. Requirements Completeness: 89 / 110

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Scenario coverage | 32/40 | Nine scenarios cover happy path, cold start, multi-framework, Convention missing, and backward compatibility. Deduction: error scenarios are weak — what if Convention content is wrong (wrong framework name, misspelled assertion library)? What if LLM ignores Convention content? No scenario for test execution failures (not compilation). No CI/CD integration impact scenario. |
| Non-functional requirements | 32/40 | Five NFRs listed: backward compatibility, progressive migration, zero new infrastructure, compile gate, on-demand loading. These are relevant. Deduction: missing performance NFR (LLM-based Convention loading vs embedded files — what's the latency impact?). No reliability NFR (how often does LLM misinterpret Convention?). No maintainability NFR (who maintains Convention format as it evolves?). |
| Constraints & dependencies | 25/30 | Four constraints listed: existing conventions directory, just dependency, no external services, markdown format. Deduction: no team skill constraint (is the team proficient in LLM prompt engineering?). No timeline constraint. No backward migration constraint (how do existing users of forge-cli migrate?). |

### 5. Solution Creativity: 61 / 100

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Novelty over industry baseline | 28/40 | The Convention file approach is pragmatic rather than revolutionary — it essentially moves embedded knowledge to user-editable markdown files with domains-based loading. The novel part for test code generation knowledge management is the pragmatism (combining LLM's general knowledge with project conventions). It does not significantly transcend established patterns. |
| Cross-domain inspiration | 15/35 | No evidence of borrowing from other domains. The proposal appears entirely derived from the project's internal experience. The "Convention file" pattern resembles `.editorconfig`, `.prettierrc`, or any number of configuration-as-code approaches, but these parallels are neither acknowledged nor leveraged. |
| Simplicity of insight | 18/25 | The three-layer separation (methodology/framework-knowledge/project-actual) is elegant. The insight "LLM doesn't need embedded knowledge — it just needs to know project conventions" is valid. Deduction: the execution is not as simple as the insight — 19+ file rewrites, new slash command, inline result parsing in run-e2e-tests. |

### 6. Feasibility: 72 / 100

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Technical feasibility | 30/40 | The tech stack (Go, markdown, just, LLM) is standard. No showstopper dependencies. Deduction: the assumption that "LLM reads Convention and correctly applies it" is unvalidated. The transition from embedded templates to LLM-prompted generation is a major risk in non-trivial projects. No POC or spike discussed. |
| Resource & timeline feasibility | 20/30 | 19+ file rewrites, 3 skill rewrites, 1 new skill, CLI command changes, Go package rewrites. The scope is significant. The proposal mentions phased delivery but provides no time estimate or team allocation. The absence of timeline feasibility analysis is notable given the scope. |
| Dependency readiness | 22/30 | Existing infrastructure (conventions directory, domains loading, just) is ready. Deduction: `/forge:test-guide` is a new slash command requiring development. The `consolidate-specs` integration is a dependency whose readiness is not discussed. |

### 7. Scope Definition: 67 / 80

| Criterion | Score | Justification |
|-----------|-------|---------------|
| In-scope items are concrete | 25/30 | In-scope items are concrete — deletion by directory/file name, rewrites by component name, creation and integration items listed. Each is a deliverable. Deduction: some in-scope items are coarse-grained ("pkg/journey/ rewrite" — what exactly does this include?). |
| Out-of-scope explicitly listed | 22/25 | Out-of-scope items are explicitly named: gen-journeys, gen-contracts, verify/promote, existing test migration, unit tests, anti-pattern docs, Convention auto-sync. Deduction: "User manual update" is in-scope but has no explicit deliverable (which sections? which pages?). |
| Scope is bounded | 20/25 | Phased delivery approach (Phase 1/2/3) provides some bounding. Deduction: no time estimate, so scope cannot be assessed as realistic for a given team. Some success criteria are open-ended ("Convention file fixed structure, skill reliably interprets" — what constitutes "reliably"?). |

### 8. Risk Assessment: 69 / 90

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Risks identified | 25/30 | Seven risks listed — a respectable number. Includes technical risks (regression, LLM quality, tag generation, result parsing) and project risks (timeline overrun). Deduction: missing risks: Convention content poisoning (user puts wrong info), LLM hallucination ignoring Convention, maintenance burden shifting from system to user, Convention format drift over time. |
| Likelihood + impact rated | 22/30 | 6 of 7 risks rated M likelihood, 1 rated H. Impact is 2 H, 5 M. Ratings seem honest — not all low-likelihood/high-impact. Deduction: no quantification — H/M/L meanings are subjective. The highest risk ("full rewrite introduces regression") is H/H, which is correct, but the remaining risks lack differentiation in likelihood (overwhelmingly M). |
| Mitigations are actionable | 22/30 | Mitigations are specific: "126+ e2e tests as regression safety net," "compile gate," "test-guide to solidify patterns." Deduction: some mitigations are repetitive ("compile gate" appears in 4 of 7 risks) — this is relying on a single mitigation rather than diverse strategies. No rollback plan — if Convention approach proves insufficient after full rewrite, what happens? |

### 9. Success Criteria: 60 / 80

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Criteria are measurable and testable | 42/55 | Of 14 criteria, most are binary-verifiable (directory doesn't exist, no imports, commands work). `just e2e-compile` + test output diff is measurable. Deduction: "Convention file fixed structure, skill reliably interprets" — "reliably" is not measurable. "gen-test-scripts outputs prompt when Convention not found" — this criterion is weak (what if the prompt is unhelpful?). No criterion measures whether generated code quality is equivalent to or better than the Profile era. Missing: compilation rate, first-pass success rate of generated code. |
| Coverage is complete | 18/25 | Criteria cover Profile removal, config cleanup, skill rewrites, CLI simplification, backward compatibility, and Convention creation. Deduction: no criterion for run-e2e-tests result parsing accuracy. No criterion for test-guide user experience quality. No criterion for progressive migration (how does user know when to create a Convention?). |

### 10. Logical Consistency: 74 / 90

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Solution addresses stated problem | 30/35 | Yes — Profile system is fully removed, replaced by Convention + Reconnaissance + just. Three-layer separation resolves the "responsibility mixing" issue. Convention files resolve the "hardcoded assumption" issue. Deduction: the proposal argues Profile detection fails (evidence table), but the solution relies on LLM-based detection (Code Reconnaissance), which has the same fundamental weakness — it can also fail for non-standard setups, just silently. |
| Scope <-> Solution <-> Success Criteria aligned | 24/30 | Generally consistent. Scope covers Profile removal, skill rewrites, Convention creation. Success criteria cover these areas. Deduction: scope includes "user manual update" but no success criterion validates it. "consolidate-specs integration" is in both scope and criteria, but the criterion is just "incorporated into management," which is vague. |
| Requirements <-> Solution coherent | 20/25 | Requirements scenarios map well to solution features. Deduction: NFR "zero new infrastructure" contradicts creation of new slash command `/forge:test-guide`. NFR "Convention on-demand loading" claims to use existing mechanism, but the Convention file format and fixed structure are new conventions, not existing infrastructure. |

## Cross-Dimension Coherence Findings

1. **NFR vs Scope contradiction**: Dimension 4 (NFR "zero new infrastructure") claims no new infrastructure, but Dimension 7 (Scope) explicitly lists creating `/forge:test-guide` slash command. A slash command is new infrastructure. Scored in Dimension 10.

2. **Evidence gap spanning Problem and Success Criteria**: Dimension 1 provides strong evidence that Profile detection fails, but provides no evidence that Convention + LLM will succeed. Dimension 9's success criteria omit quality-equivalence measurement. Together, this means the proposal proves the old system is broken but does not prove the new system works. This gap manifests in Dimension 9 (missing measurable quality criteria).

3. **Detection skepticism self-contradiction**: The proposal's core evidence is "detection fails for non-default frameworks." The proposed replacement also relies on detection — just LLM-based instead of rule-based. This irony is not acknowledged anywhere. Manifests in Dimension 10.

## Blindspot Attacks

1. **[blindspot] No rollback plan for a full system rewrite**: The proposal requires removing 19+ files and rewriting core systems. Quote: "全面移除 Profile——不保留内部实现，所有消费者重写" (Section: Core Principles). If the Convention approach proves insufficient post-rewrite, there is no documented rollback strategy. For a CTO evaluating this proposal, the absence of a rollback plan for infrastructure-level changes is a critical oversight. The risk table mentions regression but no "if it fails, revert to..." plan.

2. **[blindspot] "LLM reliably applies Convention" is an unvalidated foundational assumption**: Quote: "Convention 加载是 LLM 行为指引，不是程序化过滤。Skill prompt 指示 LLM：1. 列出 docs/conventions/ 目录中所有文件 2. 读取每个文件的 domains frontmatter 3. 仅加载 domains 与当前任务匹配的文件内容" (Section: Convention Loading Mechanism). The entire proposal rests on LLMs correctly reading, interpreting, and applying Convention files during code generation. No evidence, POC, or fallback exists. If LLMs ignore Convention content (a well-documented behavior when prompts grow long), the system fails silently. This is outside all rubric dimensions — no dimension evaluates whether the proposal's core mechanism has been validated.

3. **[blindspot] Maintenance burden transfer from system to user is unacknowledged**: Quote: "Convention 文件是用户可编辑的项目级文档" and the growth path shows "用户手动编辑 → 加反模式引用、风格偏好" (Section: Convention Files). The proposal positions Convention files as user-editable, transferring maintenance burden from the system (embedded profiles, updated by developers) to users. When a project upgrades its testing framework (the exact scenario the proposal cites as a Profile failure), who updates the Convention? This cost is not acknowledged. No rubric dimension captures the total cost of ownership shift.

4. **[blindspot] Bootstrap paradox — who validates the justfile?**: Quote: "冷启动初始化顺序：init-justfile（生成 justfile，内含 e2e-compile recipe）→ gen-test-scripts（LLM 默认 + 编译门禁）→ test-guide（固化模式）" (Section: Cold Start Handling). The compile gate (`just e2e-compile`) is the safety net. But for a true cold start, the justfile itself is LLM-generated. If the justfile contains wrong recipes (e.g., `go test` for a Python project), the compile gate validates against a wrong recipe. Nothing validates the justfile's correctness. This circular dependency is outside all rubric dimensions.

5. **[blindspot] Reasoning audit flagged independently: the proposal solves a different problem than it diagnoses**: The diagnosis is "Profile detection fails for non-default frameworks." The solution is "remove Profile entirely and use user-editable Convention files." But the solution to "detection fails" could be "fix detection" or "let users configure." Instead, the proposal uses the detection failure as justification for a complete architectural overhaul. The reasoning chain jumps from "this component has bugs" to "remove the entire component category." This logical leap is not acknowledged. (Reasoning audit flagged this independently of dimension scoring.)
