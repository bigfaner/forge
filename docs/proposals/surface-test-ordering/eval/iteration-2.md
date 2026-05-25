# Eval Iteration 2: Surface Test Ordering & Journey Unification

**Evaluator**: CTO Adversary
**Date**: 2026-05-25
**Proposal**: docs/proposals/surface-test-ordering/proposal.md
**Rubric**: plugins/forge/skills/eval/rubrics/proposal.md
**Previous Score**: 755/1000 (iteration-1)

---

## Attack Resolution Check

### Attack 1 (D5 Creativity): Zero novelty, zero cross-domain inspiration
**Status: PARTIALLY RESOLVED.** Innovation Assessment section now explicitly states "并非创新" and correctly frames the work as convention-over-configuration (Rails, 2005). This is honest self-assessment that stops overstating novelty. However, the section still claims "真正的设计价值在于两点工程决策" — these are engineering decisions, not creative contributions. The rubric asks for novelty OVER industry baseline and cross-domain inspiration. The honesty earns back some points but does not create novelty where none exists. Cross-domain inspiration remains at zero.

### Attack 2 (D3 Benchmarking): Superficial references, asymmetric trade-offs, justification by user approval
**Status: PARTIALLY RESOLVED.** The comparison table now has a "Selected Approach 技术论证" paragraph that replaces the previous "user confirmed" justification with actual technical analysis (scheduler visibility, blocked state propagation). The Cons 明细 is now a detailed list of 4 concrete code changes including function signature changes, InferType, migration, and naming coexistence. This is a significant improvement. However:
- Industry references remain a single sentence mentioning GitHub Actions `needs` and GitLab CI `dependencies` — still name-dropping.
- The execution-level ordering alternative (row 2) is still dismissed as "用户已否决" rather than with technical reasoning.
- No open-source projects, published patterns, or specification references beyond the two CI tool names.

### Attack 3 (D6 Feasibility): Missing function signature change, misplaced migration, surface-key conflict
**Status: LARGELY RESOLVED.** The Feasibility Assessment now explicitly lists:
- Function signature change from `capabilities []string` to surfaces map/surface-key list (line 124)
- InferType prefix matching with test case updates (line 125)
- Migration logic correctly placed in BuildIndex phase, not GetBreakdownTestTasks (line 126)
- Full affected code path inventory (line 127)
- Resource estimate updated from "3-5" to "5-7" coding tasks (line 131)

The surface-key conflict with justfile proposal is now addressed in Key Risks (line 177): both proposals unified to config load time normalized values (`/` to `-`), sharing the same normalization function.

### Attack 4 (D9 Success Criteria): Missing SCs for quick mode, InferType, migration, normalization, default priority
**Status: LARGELY RESOLVED.** The SC section now has 12 items (up from 6 in iteration-1). New additions:
- SC7: Quick mode ordering (line 195)
- SC8: InferType prefix matching (line 196)
- SC9: Migration correctness (line 197)
- SC10: Surface-key normalization/rejection (line 198)
- SC11: Default priority full ordering (line 199)
- SC5: Single-surface degradation (line 193)
- SC6: execution-order validation (line 194)

Remaining gap: no SC for gen-journeys SKILL.md adaptation (risk item M/M, in-scope implies SKILL.md changes, but no SC verifies the output quality).

### Attack 5 (D10 Logical Consistency): Naming inconsistency, minimal-change NFR contradiction, execution-level rejection
**Status: PARTIALLY RESOLVED.** The NFR section now accurately lists affected files beyond autogen.go + config.go: "涉及 autogen.go、config.go、infer.go 及 renderBody 模板" (line 82). The naming coexistence between gen-scripts type suffix and run-tests key suffix is now explicitly called out as "当前的设计取舍" (line 83). The execution-level ordering rejection now has technical justification in the Selected Approach 技术论证 paragraph (scheduler opacity, blocked state). However:
- The "命名一致性说明" in NFR (line 83) is a description, not a resolution — users will still see three different naming schemes in task lists.
- The "改动范围" NFR still understates — it does not mention the surface-key normalization function or the config load time validation logic as affected code paths.

### Attack 6 (Blindspot: Surface-key definition conflict with justfile proposal)
**Status: RESOLVED.** Key Risks row 2 (line 177) now explicitly addresses this: "统一为 config load time 归一化后的值（`/` 转 `-`），两提案共用同一归一化函数." This is concrete and actionable.

---

## Phase 1: Reasoning Audit

### Argument Chain Trace

**Problem**: (A) no cross-surface test ordering for fail-fast, (B) gen-journeys per-surface split contradicts Journey semantics.

**Solution**: (A) run-tests split into per-surface-key serial tasks, (B) gen-journeys merged to single task.

**Evidence**: Code references, semantic argument, token budget reasoning. HARD-RATE typo from iteration-1 remains unfixed ("SKILL.md HARD-RATE" at line 17).

### Self-Contradiction Check

1. Evidence says "gen-journeys 是纯叙事提取（读 PRD + 写 MD），不读代码" (line 18). But the solution requires merged gen-journeys to "内部遍历所有 surface type 加载对应规则" (line 31). If it loads surface-specific rules, it is not "pure narrative extraction." The revision did not address this contradiction — it added both claims in pre-revised text.

2. The comparison table row 2 says "Rejected: 用户已否决" while the Selected Approach 技术论证 paragraph gives actual technical reasons why execution-level ordering is inferior. These two justifications are inconsistent — if there is a technical case against it, the technical case should be in the comparison table, not "user vetoed."

3. The Constraints section says "所有配置校验...均在 config load time 执行（fail fast）" (line 92). But SC4 says "同类型冲突场景...在 config load time 报错" — this is consistent. However, the Key Scenarios section describes "提示用户配置 execution-order" (line 74) without specifying that this is a hard error vs. a warning. If it is a hard error at config load time, the user cannot proceed without fixing it — is this the intended UX?

---

## Phase 2: Dimension Scoring

### D1. Problem Definition (110 pts)

**Problem stated clearly**: 37/40
Two problems clearly stated. The coupling between them remains implicit — they share autogen.go and the surfaces config, but no explicit argument for why they must ship together. Pre-revision added "两种模式" note (line 31) but did not add a coupling argument.
> Deduction: -3 for bundled problems without explicit coupling justification.

**Evidence provided**: 33/40
Code references are verifiable. The "HARD-RATE" typo (should be "HARD-RULE") persists from iteration-1 — a factual error in evidence quality. The "parallel benefit ~0" claim remains an assertion without profiling data. The "两个独立 gen-journeys 任务...边界划分可能不一致" evidence is the strongest piece — concrete and specific.
> Deduction: -3 for persistent factual error; -4 for unsubstantiated performance claim.

**Urgency justified**: 23/30
"本提案应在 justfile 提案实现前落地，避免重复实现和依赖链冲突" (line 24). The cost of delay remains vague — no effort estimate for "duplicate implementation." The alternative of merging into justfile proposal is still not examined.
> Deduction: -4 for vague cost-of-delay; -3 for unexamined merge alternative.

**D1 Total: 93/110** (+1 from iteration-1)

---

### D2. Solution Clarity (120 pts)

**Approach is concrete**: 38/40
The 3-surface dependency chain example (lines 39-63) shows exact task topology for both breakdown and quick modes. Surface-key naming strategy explained with rationale (line 36). The task ID scheme is explicit.
> Deduction: -2 for task ID parsing ambiguity: `T-test-run-auth-service` — how does the parser know `auth-service` is a single surface-key and not task-type `auth` with suffix `-service`?

**User-facing behavior described**: 40/45
Key scenarios now include: typical fullstack, multi-api conflict, single-surface degradation, failure propagation. Error behavior specified ("config load time 报错"). The single-surface scalar form degradation is clearly described.
> Remaining gaps: No error message content shown. No description of adding a new surface to existing config (migration UX). No description of how users diagnose WHY a task is blocked.
> Deduction: -5 for incomplete error/diagnostic UX.

**Technical direction clear**: 33/35
Pre-revision added: InferType prefix matching, verify-regression chain-tail dependency, surface-key regex, renderBody empty TestType, function signature change (line 124), migration logic placement (line 126), affected code path inventory (line 127). This is now comprehensive.
> Deduction: -2 for not explaining how `renderBody` handles empty TestType — does it omit the type line entirely? What does the rendered task body look like?

**D2 Total: 111/120** (+9 from iteration-1)

---

### D3. Industry Benchmarking (120 pts)

**Industry solutions referenced**: 25/40
Unchanged. One sentence: "CI 系统通常通过 job dependency graph 表达跨服务测试顺序（如 GitHub Actions 的 `needs` 字段、GitLab CI 的 `dependencies`）" (line 98). No analysis of HOW these systems handle ordering, default priority, or conflict resolution. No open-source projects, no published patterns.
> Deduction: -15 for superficial references.

**At least 3 meaningful alternatives**: 26/30
Four alternatives including "do nothing." The execution-level ordering alternative is genuinely different. Post-gen dependency injection remains weakly dismissed ("gen-scripts 排序无实际意义" — but ordering gen-scripts before run-tests for the first surface IS meaningful). The Selected Approach now has a proper technical justification paragraph.
> Deduction: -4 for one weak dismissal.

**Honest trade-off comparison**: 22/25
Significantly improved. Cons 明细 (lines 111-115) now lists 4 concrete code impacts: function signature change, InferType prefix match, index.json migration, naming coexistence. This is honest trade-off analysis.
> Deduction: -3 for comparison table itself not updated to reflect the fuller scope — the table's "Cons" column says "详见下方 Cons 明细" which is a redirect rather than inline comparison.

**Chosen approach justified against benchmarks**: 20/25
Major improvement. "Selected Approach 技术论证" (line 109) replaces the previous "user confirmed" with actual technical analysis: scheduler visibility, independent status per surface, native blocked state propagation. This is genuine justification.
> Remaining gap: The comparison table row 2 still says "Rejected: 用户已否决" — if there is a technical case, the table should reflect it. Two different justifications for the same rejection is inconsistent.
> Deduction: -5 for inconsistent rejection justification across table and text.

**D3 Total: 93/120** (+10 from iteration-1)

---

### D4. Requirements Completeness (110 pts)

**Scenario coverage**: 35/40
Four scenarios: typical fullstack, multi-api conflict, single-surface degradation, failure propagation. Pre-revision added single-surface scalar form (line 76).
> Missing: surface removal scenario (old per-surface-key tasks become orphans). Partial execution-order specification (2 of 3 surfaces listed). Quick mode is now covered by SC7 but not listed as a separate Key Scenario.
> Deduction: -5 for missing removal and partial-config scenarios.

**Non-functional requirements**: 37/40
Pre-revision improved significantly: surface-key regex constraint, validation timing (config load time), affected files list expanded. The "改动范围" NFR now lists 4 files (autogen.go, config.go, infer.go, renderBody). The naming coexistence is explicitly called out as a "设计取舍."
> Remaining gap: No performance NFR for serial execution — API tests block web tests even when independent. This is a latent performance regression with no measurement plan.
> Deduction: -3 for missing serial execution performance NFR.

**Constraints & dependencies**: 26/30
Surface-key constraint and validation timing well-specified (lines 91-92). Justfile proposal dependency acknowledged.
> Missing: InferType exact-match contract as a constraint (the current contract is being broken). Whether the justfile proposal has conflicting plans for run-test restructuring.
> Deduction: -4 for missing technical constraints.

**D4 Total: 98/110** (+5 from iteration-1)

---

### D5. Solution Creativity (100 pts)

**Novelty over industry baseline**: 20/40
The Innovation Assessment (line 67) now honestly states "并非创新" and correctly identifies the core mechanism as convention-over-configuration (Rails, 2005). This honesty is valued — it stops overstating novelty. The two "engineering decisions" (gen-journeys merge, single-surface degradation rule) are practical but not creative contributions.
> Deduction: -20 for honest acknowledgment of zero novelty — the rubric scores novelty, and there is none.

**Cross-domain inspiration**: 5/35
Unchanged. No references to build systems (Make/Bazel), package managers (npm peerDeps), workflow engines (Airflow DAGs), or any domain beyond CI/CD.
> Deduction: -30 for zero cross-domain inspiration.

**Simplicity of insight**: 20/25
The gen-journeys merge remains elegant — journeys are cross-surface by definition, splitting by surface was always wrong. The single-surface degradation rule (scalar form, no suffix) is a clean touch. The run-test split introduces complexity that is proportional to benefit, not elegant.
> Deduction: -5 for run-test split complexity.

**D5 Total: 45/100** (+7 from iteration-1)

---

### D6. Feasibility (100 pts)

**Technical feasibility**: 36/40
Significantly improved. Pre-revision now explicitly lists:
- Function signature change: `capabilities []string` → surfaces map/surface-key list (line 124)
- InferType prefix matching with test updates (line 125)
- Migration logic correctly placed in BuildIndex (line 126)
- Full affected code path inventory (line 127)
> Remaining gap: The `renderBody` empty TestType mechanism is mentioned but not explained — what does the rendered output look like? Is the `{{TEST_TYPE}}` line omitted entirely?
> Deduction: -4 for incomplete renderBody explanation.

**Resource & timeline feasibility**: 26/30
Updated from "3-5" to "5-7 coding tasks" (line 131). The justification for the increase is provided (function signature change, InferType, migration, surface-key normalization). This is more realistic.
> Deduction: -4 for no breakdown of the 5-7 tasks — which are the 5-7 individual coding tasks?

**Dependency readiness**: 26/30
The surface-key definition conflict with justfile proposal is now explicitly addressed in Key Risks (line 177): unified to config load time normalized values, shared normalization function. This resolves the primary concern.
> Remaining gap: Whether the justfile proposal has agreed to this shared normalization function, or if this is a unilateral proposal.
> Deduction: -4 for unilateral dependency assumption.

**D6 Total: 88/100** (+11 from iteration-1)

---

### D7. Scope Definition (80 pts)

**In-scope items are concrete**: 29/30
Most items are now specific deliverables with file names and mechanism descriptions. InferType update, verify-regression dependency, migration step are all concrete.
> Remaining vagueness: "更新 resolveBreakdownDeps 和 resolveQuickDeps 依赖链" — the WHAT is now clear (serial chain) but the HOW (loop → serial chain construction from execution-order or default priority) is implied but not explicit.
> Deduction: -1 for one item with unspecified mechanism.

**Out-of-scope explicitly listed**: 23/25
Six out-of-scope items. The gen-journeys SKILL.md scope ambiguity from iteration-1 is partially resolved — it is now in the risks table with M/M and the mitigation references SC3. But it is still not explicitly in-scope or out-of-scope.
> "部分继续运行" (line 168) is out-of-scope, but the serial design makes this architecturally impossible without redesign. This should be noted as a foreclosed future option.
> Deduction: -2 for SKILL.md scope ambiguity and foreclosed option not noted.

**Scope is bounded**: 23/25
The expanded in-scope items provide better bounds. The affected code path inventory (line 127) helps bound the scope. The function signature change impact is now acknowledged.
> Deduction: -2 for scope boundary still not explicitly stating whether the surface-key normalization function is a new shared module or added to config.go.

**D7 Total: 75/80** (+3 from iteration-1)

---

### D8. Risk Assessment (90 pts)

**Risks identified**: 27/30
Now 8 risks listed, including 3 new ones from pre-revision:
- index.json orphan risk (line 181)
- InferType prefix match ambiguity (line 182)
- Function signature change impact (line 183)
- gen-scripts/run-tests naming coexistence (line 184)
> Missing: No rollback plan if serial execution causes unacceptable latency. No risk for serial execution being a performance regression on the happy path.
> Deduction: -3 for missing rollback/performance regression risk.

**Likelihood + impact rated**: 24/30
Ratings are now more diverse: M/L, M/M, L/L, M/M, M/M, M/M, M/L, L/M. The index.json orphan risk at M/M seems reasonable. The InferType ambiguity at M/M seems reasonable.
> However: "gen-journeys 单任务加载多 surface 规则增加 context 噪音" at M/L — but if gen-journeys is "pure narrative extraction" (Evidence section), noise should be negligible. The likelihood rating contradicts the evidence claim.
> The surface-key conflict with justfile proposal is now resolved (line 177) — but its risk row still shows M/M. If the mitigation is concrete and shared, the likelihood should be reduced.
> Deduction: -6 for ratings contradicting evidence claims and outdated likelihood.

**Mitigations are actionable**: 26/30
Significantly improved. Pre-revision added:
- Specific migration strategy for index.json orphans (copy status/blocked-reason to first surface-key, line 181)
- InferType fallback mechanism (prefix → lookup → fallback to exact match, line 182)
- Shared normalization function for justfile proposal alignment (line 177)
> Remaining gaps: "gen-journeys 的规则加载仅为参考信息...SKILL.md 中增加 Multi-Surface Rules Loading 段落" — this is a plan, not a mitigation. The risk is that noise causes incorrect generation; the mitigation should describe how to verify correctness.
> "文档中明确说明命名差异的原因" (line 184) — documentation is not a mitigation for user confusion. The confusion is real and will persist regardless of documentation.
> Deduction: -4 for two non-actionable mitigations.

**D8 Total: 77/90** (+9 from iteration-1)

---

### D9. Success Criteria (80 pts)

**Criteria are measurable and testable**: 27/30
12 SCs, most are specific and testable. Pre-revision added:
- SC8: `InferType("T-test-run-backend")` returns `api` with prefix matching — precise, testable
- SC9: Migration correctness with specific BuildIndex phase — testable
- SC10: Surface-key normalization/rejection with exact examples — excellent, testable
- SC11: Default priority with 4-surface example — testable
> Remaining issues:
> - SC3: "输出覆盖所有配置 surface 的 Journey 文件（非仅单个 surface）" — improved from iteration-1 (now specifies output coverage), but "所有配置 surface" is still vague. How many Journey files? What determines coverage completeness?
> - SC5: "任务 ID 和依赖列表与改动前一致" — "一致" needs enumeration: same task ID, same dependencies, same title, same body?
> Deduction: -3 for two criteria with residual vagueness.

**Coverage is complete**: 22/25
Major improvement. 12 SCs now cover: ordering, failure propagation, gen-journeys output, conflict detection, degradation, validation, quick mode, InferType, migration, normalization, default priority.
> Remaining gaps:
> - No SC for gen-journeys SKILL.md adaptation quality (risk item M/M, but no SC verifies the SKILL.md change works correctly).
> - No SC for the naming coexistence being comprehensible to users (risk item L/M, but no SC verifies user-facing behavior).
> Deduction: -3 for 2 risk items without corresponding SCs.

**SC internal consistency**: 22/25
Cluster analysis:

**gen-journeys cluster**: SC3 (single task, output covers all surfaces) is satisfiable. In-scope says TestType left empty, renderBody adapts. Consistent.

**run-test cluster**: SC1 (ordering), SC2 (failure propagation), SC4 (conflict detection), SC5 (degradation), SC6 (validation), SC7 (quick mode), SC8 (InferType), SC9 (migration), SC10 (normalization), SC11 (default priority). SC5 says "任务 ID 和依赖列表与改动前一致" — but InferType matching semantics change from exact to prefix (SC8). The BEHAVIOR may be identical but the code path is different. If SC5 means "same observable task IDs and dependencies," it is satisfiable alongside SC8.

**Validation timing**: SC4 says "config load time 报错" and SC6 says "config load time 报错" — both now consistently specify config load time. The iteration-1 timing inconsistency between SC4/SC6 is resolved.
> Deduction: -3 for SC5/singleton degradation requiring clarification of what "一致" means exactly.

**D9 Total: 71/80** (+7 from iteration-1)

---

### D10. Logical Consistency (90 pts)

**Solution addresses the stated problem**: 30/35
- Problem A (no cross-surface ordering): per-surface-key serial tasks with execution-order solves this directly.
- Problem B (gen-journeys semantic inconsistency): merging to single task directly addresses this.
- The execution-level ordering alternative is now rejected with technical justification (scheduler opacity, blocked state), not just "user vetoed." This is a stronger logical basis.
> Remaining: The two problems are still bundled without explicit coupling justification. The gen-journeys merge could ship independently as a separate proposal.
> Deduction: -5 for unexplained problem coupling.

**Scope <-> Solution <-> Success Criteria aligned**: 24/30
Significantly improved. Pre-revision added:
- SC7 (quick mode) covers resolveQuickDeps in-scope item
- SC8 (InferType) covers InferType update in-scope item
- SC9 (migration) covers migration step in-scope item
- SC10 (normalization) covers surface-key constraint
- SC11 (default priority) covers the priority convention

Remaining gaps:
- The naming coexistence (gen-scripts type suffix, run-tests key suffix, gen-journeys no suffix) is described in NFR (line 83) as a "设计取舍" but has no SC verifying user-facing task list coherence.
- verify-regression dependency on chain tail is in-scope but has no explicit SC — it is implicitly tested by SC1 and SC7 (ordering includes verify-regression position in the dependency chain diagrams).
> Deduction: -6 for naming coexistence without SC and implicit verify-regression coverage.

**Requirements <-> Solution coherent**: 20/25
The "改动范围" NFR now lists 4 files (line 82), a significant improvement over the previous "only autogen.go and config.go" claim. The naming coexistence is explicitly called out.
> Remaining gap: The "向后兼容" NFR says "单 surface 项目的任务结构和依赖链不变" — this is tested by SC5. But the NFR does not mention that multi-surface projects have NO backward state to be compatible with (they are new functionality). The backward compatibility claim is technically vacuous for the feature's primary use case.
> The Constraints section (line 92) says validation happens at config load time, but the solution's surface-key normalization (line 91) happens at config load time too. If normalization changes the key, the user's config file now contains keys that differ from what they wrote. This implicit mutation is not discussed.
> Deduction: -5 for vacuous backward compatibility claim and unmentioned config mutation.

**D10 Total: 74/90** (+8 from iteration-1)

---

## Phase 3: Blindspot Hunt

**[blindspot-1] No rollback plan.** If serial ordering causes unacceptable pipeline latency (API tests take 2h, web tests cannot start), there is no mechanism to revert to parallel execution without removing the execution-order configuration entirely. The proposal optimizes for fail-fast but may degrade the happy path. No measurement plan for latency impact.

**[blindspot-2] Config mutation is invisible to users.** Surface-key normalization (line 91) transforms `ADMIN PANEL` to `admin-panel` at config load time. The user's YAML file is not updated — the normalized value exists only in memory. If the user sees error messages referencing `admin-panel` but wrote `ADMIN PANEL`, the diagnosis path is unclear. No SC tests error message content for normalized keys.

**[blindspot-3] Serial execution is a performance regression for independent surfaces.** If backend API and frontend web are truly independent (no shared state), running them serially is pure waste. The proposal assumes dependencies exist between all surfaces by default. The fail-fast benefit only materializes when tests actually fail — the minority case. The proposal does not quantify the expected latency increase for the happy path.

**[blindspot-4] gen-journeys "pure narrative" claim is contradicted by implementation.** Evidence says "纯叙事提取（读 PRD + 写 MD），不读代码" (line 18). Solution says "内部遍历所有 surface type 加载对应规则" (line 31). These two claims are contradictory — if it loads rules, it is not pure narrative extraction.

---

## Bias Detection Report

Annotated regions analysis:
- 10 pre-revised markers covering approximately 12 paragraphs.
- Attacks on annotated regions: 5 attack points (task ID parsing in D2, renderBody gap in D2/D6, scope ambiguity in D7, SC5 vagueness in D9, naming coexistence in D10).
- Annotated density: 5/12 = 0.42

Unannotated regions analysis:
- Approximately 22 paragraphs in unannotated sections.
- Attacks on unannotated regions: 10 attack points (HARD-RATE typo in D1, performance assertion in D1, urgency alternatives in D1, benchmarking superficiality in D3, execution-level dismissal in D3, comparison table inconsistency in D3, missing performance NFR in D4, ratings vs evidence contradiction in D8, rollback missing in D8/blindspot-1, "pure narrative" contradiction in blindspot-4).
- Unannotated density: 10/22 = 0.45

**Ratio (annotated/unannotated): 0.93**

The ratio is near 1.0, indicating balanced scrutiny across revised and unrevised regions. Pre-revision regions introduced fewer new issues than in iteration-1, reflecting more careful revision. The unrevised regions continue to carry structural weaknesses (benchmarking, rollback, performance) that were not addressed in this iteration.

---

## Score Summary

| Dimension | Score | Max |
|-----------|-------|-----|
| D1. Problem Definition | 93 | 110 |
| D2. Solution Clarity | 111 | 120 |
| D3. Industry Benchmarking | 93 | 120 |
| D4. Requirements Completeness | 98 | 110 |
| D5. Solution Creativity | 45 | 100 |
| D6. Feasibility | 88 | 100 |
| D7. Scope Definition | 75 | 80 |
| D8. Risk Assessment | 77 | 90 |
| D9. Success Criteria | 71 | 80 |
| D10. Logical Consistency | 74 | 90 |
| **Total** | **825** | **1000** |

---

## Top Attacks (Priority Order)

1. **D5 Creativity**: Zero novelty beyond industry baseline and zero cross-domain inspiration. The Innovation Assessment honestly acknowledges this ("并非创新"), which is valued, but honesty does not create points where the rubric demands novelty. 55/100 lost. Must find genuine differentiation or cross-domain inspiration.

2. **D3 Industry Benchmarking**: Industry references remain a single sentence with two CI tool names — no analysis of HOW these systems handle the specific problems (ordering, default priority, conflict resolution). The execution-level ordering alternative is still dismissed as "用户已否决" in the comparison table while the body text provides technical justification — inconsistent. 27/120 lost.

3. **D10 Logical Consistency**: Two problems bundled without coupling argument; naming coexistence across gen-scripts/run-tests/gen-journeys has no SC; backward compatibility NFR is vacuous for the primary use case (multi-surface projects are new functionality); config mutation via normalization is invisible to users. 16/90 lost.

4. **D8 Risk Assessment**: No rollback plan for serial execution latency regression; gen-journeys noise risk rating (M/L) contradicts "pure narrative" evidence claim; "documentation" as mitigation for naming confusion is not actionable. 13/90 lost.

5. **D9 Success Criteria**: SC3 output coverage is vague ("所有配置 surface"); SC5 "一致" needs enumeration; no SC for gen-journeys SKILL.md adaptation quality; no SC for user-facing task list coherence with mixed naming. 9/80 lost.

6. **[blindspot-4] gen-journeys "pure narrative" contradiction**: Evidence claims "纯叙事提取...不读代码" but solution requires "内部遍历所有 surface type 加载对应规则." These two claims cannot both be true.
