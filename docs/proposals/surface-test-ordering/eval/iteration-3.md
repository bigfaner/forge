# Eval Iteration 3: Surface Test Ordering & Journey Unification

**Evaluator**: CTO Adversary
**Date**: 2026-05-25
**Proposal**: docs/proposals/surface-test-ordering/proposal.md
**Rubric**: plugins/forge/skills/eval/rubrics/proposal.md
**Previous Score**: 825/1000 (iteration-2)

---

## Attack Resolution Check (Iteration-2 → Iteration-3)

### Attack 1 (D5 Creativity): Zero cross-domain inspiration
**Status: RESOLVED.** Added "跨领域启发" paragraph citing Bazel's query filter/visibility (label-based disambiguation for same-type targets) and Kubernetes init containers (failure propagation in dependency chains). Both analogies are structurally isomorphic to the proposal's mechanisms (surface-key disambiguation, blocked state propagation). The critical difference is articulated: AI agent pipelines have inferred test types vs. explicit declarations in Bazel/K8s. Cross-domain inspiration now scores above zero.

### Attack 2 (D3 Benchmarking): Superficial industry references, execution-level rejection inconsistency
**Status: LARGELY RESOLVED.**
- Industry Solutions expanded from one sentence to two substantive paragraphs analyzing HOW GitHub Actions (full-parallel + explicit-serial model, no default priority) and GitLab CI (stage-based ordering, manual partitioning) handle ordering. Each analysis concludes with an explicit delta from the proposal.
- Execution-level ordering row in comparison table changed from "Rejected: 用户已否决" to "Rejected: 调度器语义不完整——见 Selected Approach 技术论证". This is now consistent with the body text justification.
- Remaining gap: No open-source project beyond CI tools, no published pattern references (e.g., Buildbot, Airflow DAG scheduling, Terraform dependency graphs).

### Attack 3 (D6 Feasibility): Already resolved in iteration-2. No new issues introduced.
**Status: MAINTAINED.**

### Attack 4 (D9 Success Criteria): Missing SC for gen-journeys SKILL.md adaptation
**Status: RESOLVED.** Added SC12: "gen-journeys SKILL.md 多 surface 适配：`surfaces: { frontend: web, backend: api }` 时，gen-journeys 输出的 Journey 文件中每个 Journey 明确标注覆盖的 surface type 集合（如 `[web, api]`），且无遗漏——所有配置的 surface type 至少被一个 Journey 覆盖". This is precise, testable, and covers the risk item.

### Attack 5 (D10 Logical Consistency): Naming inconsistency, backward compatibility vacuity
**Status: PARTIALLY RESOLVED.**
- "向后兼容" NFR now explicitly acknowledges: "多 surface 项目是新增功能，无先前状态需兼容". This resolves the vacuous claim.
- Coupling argument added as a standalone paragraph before the two solution items. Explicitly states shared root cause (per-surface task generation model semantic defects), shared code (autogen.go per-surface loop, surfaces config), and cost of partial implementation (inconsistent task topology granularity). This resolves the unexplained problem coupling.
- Naming coexistence still has no SC verifying user-facing task list coherence.

### Attack 6 (D1 Problem Definition): HARD-RATE typo, performance assertion
**Status: RESOLVED.** "HARD-RATE" corrected to "HARD-RULE". "纯叙事提取（读 PRD + 写 MD），不读代码" revised to "叙事提取型任务（读 PRD 为主输入，加载 surface 规则作为参考指导，但不读源代码）". This resolves the iteration-2 blindspot-4 contradiction (the claim now matches the solution's "内部遍历所有 surface type 加载对应规则").

### Attack 7 (D8 Risk Assessment): gen-journeys noise risk rating contradicts "pure narrative" evidence
**Status: RESOLVED.** Likelihood downgraded from M to L. Mitigation text updated to match revised evidence: "gen-journeys 以 PRD 为主要输入，surface 规则仅作参考指导...噪音影响可忽略". Consistent with evidence section.

### Attack 8 (D8 Risk Assessment): No rollback plan for serial execution latency
**Status: RESOLVED.** New risk row added: "串行执行导致 happy path 延迟回归 | M | M | 回滚策略...短期缓解：默认串行仅在失败时产生 fail-fast 收益，happy path 无额外开销——因为 per-surface 任务仍由调度器调度，串行仅影响启动时机". The mitigation includes a specific architectural argument for why happy-path overhead is minimal.

### Attack 9 (D9 Success Criteria): SC3 output coverage vague, SC5 "一致" ambiguous
**Status: NOT ADDRESSED.** SC3 still says "输出覆盖所有配置 surface 的 Journey 文件" — no change to this wording. SC5 still says "任务 ID 和依赖列表与改动前一致" — no enumeration of what "一致" means.

### Attack 10 (D10/D3): Comparison table inconsistency, missing performance NFR
**Status: PARTIALLY RESOLVED.** Comparison table row 2 now consistent with body text. Performance NFR for serial execution still absent from Requirements section, though the new risk row in Key Risks addresses the concern.

---

## Phase 1: Reasoning Audit

### Argument Chain Trace

**Problem**: (A) no cross-surface test ordering for fail-fast, (B) gen-journeys per-surface split contradicts Journey semantics.

**Solution**: (A) run-tests split into per-surface-key serial tasks, (B) gen-journeys merged to single task.

**Coupling**: Now explicitly argued — shared root cause (per-surface loop in autogen.go), shared data structure (surfaces config), partial implementation creates inconsistent granularity.

**Evidence**: Code references verifiable. HARD-RULE now correct. "叙事提取型任务" claim now consistent with solution's rule-loading mechanism.

### Self-Contradiction Check

1. ~~"纯叙事提取" vs "加载对应规则"~~ **RESOLVED.** Evidence now says "叙事提取型任务（读 PRD 为主输入，加载 surface 规则作为参考指导，但不读源代码）" — consistent with solution loading surface rules as reference.

2. ~~Comparison table "用户已否决" vs body text technical justification~~ **RESOLVED.** Both now reference "调度器语义不完整".

3. ~~gen-journeys noise risk M/L contradicts "pure narrative" evidence~~ **RESOLVED.** Risk downgraded to L/L, evidence revised.

4. **New potential issue**: The happy-path latency mitigation says "串行仅影响启动时机" — but if backend tests take 30 minutes, frontend tests still cannot START until backend finishes. This is a 30-minute delay on the happy path, not zero overhead. The claim "happy path 无额外开销" is misleading if "overhead" includes wall-clock time. The mitigation should say "happy path 延迟 = 上游 surface 测试耗时" rather than "无额外开销".

---

## Phase 2: Dimension Scoring

### D1. Problem Definition (110 pts)

**Problem stated clearly**: 39/40
Two problems clearly stated with code references. Coupling now explicitly argued: shared root cause, shared code paths, partial-implementation cost.
> Deduction: -1 for coupling argument still being dense — a reader may miss the coupling in the wall of text. A one-line summary ("both stem from autogen.go's per-surface loop") would tighten it.

**Evidence provided**: 38/40
HARD-RULE typo fixed. "纯叙事提取" contradiction resolved with accurate "叙事提取型任务" language. Code references remain verifiable. The "边界划分可能不一致" evidence remains the strongest piece.
> Deduction: -2 for "并行收益≈0" still being an assertion without profiling data or theoretical model (e.g., "subagent spawn overhead ~30s, gen-journeys wall-time ~20min, parallel speedup < 2.5%").

**Urgency justified**: 26/30
"本提案应在 justfile 提案实现前落地，避免重复实现和依赖链冲突" (line 24). The cost of delay is more credible now that the coupling argument exists — if justfile implements per-surface task generation independently, the gen-journeys merge and run-tests split would require rework.
> Remaining gap: The alternative of merging this INTO the justfile proposal (a combined change) is still not examined. If both proposals touch the same code paths, combining them might reduce total effort.
> Deduction: -4 for unexamined merge-into-justfile alternative.

**D1 Total: 103/110** (+8 from iteration-2)

---

### D2. Solution Clarity (120 pts)

**Approach is concrete**: 39/40
3-surface dependency chain diagrams show exact task topology for both modes. Surface-key naming strategy explained. Coupling argument added.
> Deduction: -1 for task ID parsing ambiguity remaining: `T-test-run-auth-service` — the parser cannot distinguish a multi-hyphen surface-key from a task-type suffix without a lookup table. The InferType fallback (prefix → surfaces map lookup → exact match) is described in Key Risks but not in the Solution section.

**User-facing behavior described**: 42/45
Key scenarios cover: fullstack, multi-api conflict, single-surface degradation, failure propagation. Error timing specified (config load time).
> Remaining gaps: No error message content. No description of how users diagnose why a task is blocked. No description of adding/removing a surface from existing config (migration UX beyond index.json fix-tasks).
> Deduction: -3 for incomplete error/diagnostic/migration UX.

**Technical direction clear**: 34/35
Comprehensive: InferType prefix matching, verify-regression chain-tail, surface-key regex, renderBody empty TestType, function signature change, migration logic placement, affected code path inventory.
> Deduction: -1 for renderBody empty TestType — mentioned but rendered output not shown.

**D2 Total: 115/120** (+4 from iteration-2)

---

### D3. Industry Benchmarking (120 pts)

**Industry solutions referenced**: 32/40
Significant improvement. Two substantive paragraphs analyzing GitHub Actions (full-parallel + explicit-serial, no default priority, no conflict detection) and GitLab CI (stage-based ordering, manual partitioning, no auto-inference). Each concludes with explicit delta from the proposal.
> Remaining gap: Still limited to CI/CD domain. No build system references (Make, Bazel, Pants), no workflow engine references (Airflow, Prefect), no package manager dependency resolution analogies.
> Deduction: -8 for domain-bounded references only.

**At least 3 meaningful alternatives**: 27/30
Four alternatives including "do nothing". Execution-level ordering now rejected with technical reasoning in comparison table (consistent with body text). Post-gen dependency injection dismissal still weak ("gen-scripts 排序无实际意义").
> Deduction: -3 for one weak dismissal.

**Honest trade-off comparison**: 23/25
Cons 明细 lists 4 concrete code impacts. Comparison table updated to be consistent with body text.
> Deduction: -2 for Cons column redirecting to "详见下方 Cons 明细" rather than inline.

**Chosen approach justified against benchmarks**: 23/25
"Selected Approach 技术论证" provides genuine technical analysis: scheduler visibility, independent status, native blocked state. The comparison with GitHub Actions/GitLab CI shows where the proposal innovates (default priority, auto-inference, conflict detection).
> Deduction: -2 for not explaining why the proposal doesn't adopt GitLab's stage concept (which is closer to execution-order than GitHub's needs-only model).

**D3 Total: 105/120** (+12 from iteration-2)

---

### D4. Requirements Completeness (110 pts)

**Scenario coverage**: 36/40
Four key scenarios: fullstack, multi-api conflict, single-surface degradation, failure propagation.
> Missing: surface removal scenario. Partial execution-order specification. Quick mode not listed as separate scenario (covered by SC7).
> Deduction: -4 for missing removal and partial-config scenarios.

**Non-functional requirements**: 38/40
Backward compatibility NFR now honestly scoped: "多 surface 项目是新增功能，无先前状态需兼容". Affected files list covers 4 code paths. Naming coexistence called out as design trade-off.
> Remaining gap: Serial execution performance — the happy-path wall-clock delay is not quantified as an NFR. The new risk row addresses it but risks are not requirements. No NFR says "serial execution MUST NOT increase total pipeline time by more than X%".
> Deduction: -2 for missing serial execution latency NFR.

**Constraints & dependencies**: 27/30
Surface-key regex, validation timing, justfile proposal dependency all specified.
> Deduction: -3 for not stating InferType exact-match contract as a breaking change constraint.

**D4 Total: 101/110** (+3 from iteration-2)

---

### D5. Solution Creativity (100 pts)

**Novelty over industry baseline**: 25/40
The Innovation Assessment honestly states "并非创新" for the core convention-over-configuration mechanism. The genuine differentiations from industry are: (1) default priority inferred from surface type (not user-specified like GitLab stages), (2) automatic same-type conflict detection at config load time. These are incremental innovations, not paradigm shifts.
> Deduction: -15 for limited novelty — the differentiations are incremental, not structural.

**Cross-domain inspiration**: 22/35
Major improvement. Bazel's query filter/visibility (label-based disambiguation for same-type targets) and Kubernetes init containers (failure propagation in dependency chains) are structurally isomorphic analogies. The critical difference is articulated (inferred vs. declared types).
> Remaining gap: No inspiration from build systems' incremental compilation models, package managers' dependency resolution algorithms, or database transaction isolation levels (serializability of independent transactions).
> Deduction: -13 for two analogies only, both from infrastructure/DevOps domain.

**Simplicity of insight**: 22/25
gen-journeys merge is elegant (journeys are cross-surface by definition). Single-surface degradation (scalar form, no suffix) is clean. The run-test split is proportional complexity.
> Deduction: -3 for the three naming schemes coexisting (gen-scripts type, run-tests key, gen-journeys none) — this is not simple.

**D5 Total: 69/100** (+24 from iteration-2)

---

### D6. Feasibility (100 pts)

**Technical feasibility**: 37/40
Function signature change, InferType prefix matching, migration logic placement, affected code path inventory — all explicit. Evidence contradiction resolved.
> Deduction: -3 for renderBody empty TestType mechanism still unexplained.

**Resource & timeline feasibility**: 27/30
"5-7 coding tasks" with justification for the increase.
> Deduction: -3 for no breakdown of the 5-7 tasks.

**Dependency readiness**: 27/30
Justfile proposal alignment addressed with shared normalization function.
> Deduction: -3 for unilateral proposal of shared function — no evidence of justfile proposal author agreement.

**D6 Total: 91/100** (+3 from iteration-2)

---

### D7. Scope Definition (80 pts)

**In-scope items are concrete**: 29/30
Most items are specific deliverables with file names and mechanism descriptions.
> Deduction: -1 for "更新 resolveBreakdownDeps 和 resolveQuickDeps 依赖链" — HOW (serial chain construction) implied but not explicit.

**Out-of-scope explicitly listed**: 24/25
Six out-of-scope items. "部分继续运行" out-of-scope noted.
> Deduction: -1 for "部分继续运行" being architecturally foreclosed by serial design without noting this.

**Scope is bounded**: 23/25
Affected code path inventory bounds the scope. Function signature change acknowledged.
> Deduction: -2 for surface-key normalization function module ownership unspecified.

**D7 Total: 76/80** (+1 from iteration-2)

---

### D8. Risk Assessment (90 pts)

**Risks identified**: 29/30
Now 9 risks including new "串行执行导致 happy path 延迟回归" row. This directly addresses iteration-2's blindspot-1 and blindspot-3.
> Deduction: -1 for no risk about config mutation visibility (blindspot-2 from iteration-2: normalized keys in error messages don't match user-written keys).

**Likelihood + impact rated**: 27/30
gen-journeys noise risk corrected to L/L (consistent with revised evidence). Serial execution risk at M/M (honest). Justfile conflict still M/M despite mitigation being concrete — but the mitigation is a plan ("共用同一归一化函数"), not yet agreed.
> Deduction: -3 for justfile conflict risk not downgraded given the concrete mitigation.

**Mitigations are actionable**: 27/30
Serial execution mitigation includes architectural argument ("串行仅影响启动时机") and rollback strategy. InferType fallback mechanism is concrete. Migration strategy specifies BuildIndex phase and status inheritance.
> Remaining gap: The serial execution mitigation claims "happy path 无额外开销" which is misleading — it means "no per-task overhead" not "no wall-clock delay". A backend test taking 30min still blocks frontend for 30min.
> Deduction: -3 for misleading "无额外开销" claim in serial execution mitigation.

**D8 Total: 83/90** (+6 from iteration-2)

---

### D9. Success Criteria (80 pts)

**Criteria are measurable and testable**: 28/30
13 SCs, most precise and testable. SC12 added for gen-journeys SKILL.md adaptation: "每个 Journey 明确标注覆盖的 surface type 集合（如 `[web, api]`），且无遗漏" — testable.
> Remaining issues:
> - SC3: "输出覆盖所有配置 surface 的 Journey 文件" — "覆盖" is still imprecise. How many files? What constitutes coverage?
> - SC5: "任务 ID 和依赖列表与改动前一致" — "一致" not enumerated.
> Deduction: -2 for two criteria with residual vagueness.

**Coverage is complete**: 24/25
13 SCs now cover: ordering, failure propagation, gen-journeys output, conflict detection, degradation, validation, quick mode, InferType, migration, normalization, default priority, SKILL.md adaptation.
> Remaining gap: No SC for user-facing task list coherence with mixed naming (gen-scripts type suffix, run-tests key suffix, gen-journeys no suffix).
> Deduction: -1 for naming coherence without SC.

**SC internal consistency**: 23/25
Clusters:
- **gen-journeys cluster**: SC3 (output coverage) + SC12 (SKILL.md adaptation) — satisfiable. SC12 adds a specific verification for the mechanism SC3 describes at the output level.
- **run-test cluster**: SC1-2, SC4-11 — no new contradictions.
- **Singleton degradation**: SC5 ("一致") remains ambiguous but not contradictory.
> Deduction: -2 for SC5 ambiguity still unresolved.

**D9 Total: 75/80** (+4 from iteration-2)

---

### D10. Logical Consistency (90 pts)

**Solution addresses the stated problem**: 33/35
Both problems directly addressed. Coupling now explicitly argued. Execution-level rejection consistent across table and text.
> Deduction: -2 for coupling argument dense but present.

**Scope <-> Solution <-> Success Criteria aligned**: 27/30
SC12 covers gen-journeys SKILL.md in-scope. SC7 covers quick mode. SC8 covers InferType. SC9 covers migration. SC10 covers normalization. SC11 covers default priority.
> Remaining gaps:
> - Naming coexistence (three schemes) still has no SC.
> - verify-regression chain-tail dependency has no explicit SC (implicit in dependency diagrams).
> - New risk (serial happy-path latency) has no corresponding SC measuring latency impact.
> Deduction: -3 for three items without SC coverage.

**Requirements <-> Solution coherent**: 22/25
Backward compatibility NFR now honestly scoped. Affected files comprehensive. Config mutation via normalization (iteration-2 blindspot-2) still not discussed — user writes `ADMIN PANEL`, system uses `admin-panel`, error messages may reference the normalized form.
> Deduction: -3 for unmentioned config mutation visibility gap.

**D10 Total: 82/90** (+8 from iteration-2)

---

## Phase 3: Blindspot Hunt

**[blindspot-1] "无额外开销" claim is misleading.** The serial execution risk mitigation states "happy path 无额外开销——因为 per-surface 任务仍由调度器调度，串行仅影响启动时机". "Startup timing" IS the overhead — if backend tests take 30 minutes, frontend waits 30 minutes before starting. The claim conflates "no per-task computational overhead" with "no wall-clock delay". The happy-path regression is real and should be quantified.

**[blindspot-2] Config mutation visibility gap persists.** Surface-key normalization transforms user input silently. Error messages, task IDs, and logs will reference `admin-panel` while the user wrote `ADMIN PANEL`. This is not discussed anywhere in the proposal. The SC10 tests that normalization works, but no SC tests that error messages are comprehensible to users who wrote the un-normalized form.

---

## Bias Detection Report

Annotated regions analysis:
- 10 pre-revised markers covering approximately 12 paragraphs (unchanged from iteration-2).
- Attacks on annotated regions: 3 new attack points found (renderBody gap, task ID parsing, scope mechanism).
- Annotated density: 3/12 = 0.25

Unannotated regions analysis:
- Approximately 24 paragraphs in unannotated sections (expanded by new additions).
- New unannotated additions: coupling argument, cross-domain inspiration, revised industry analysis, revised comparison table row, revised noise risk, serial execution risk, SC12.
- Attacks on unannotated regions: 4 attack points found ("无额外开销" misleading claim, config mutation gap, SC3/SC5 vagueness, naming coherence without SC).
- Unannotated density: 4/24 = 0.17

**Ratio (annotated/unannotated): 1.47**

The ratio is above 1.0, indicating the revised regions are being scrutinized more heavily. This is expected — the new additions (coupling argument, cross-domain inspiration, serial execution risk, SC12) are all unannotated and are well-constructed with few issues. The persistent weaknesses remain in the pre-revised annotated regions (SC3/SC5 vagueness, renderBody gap).

---

## Score Summary

| Dimension | Score | Max | Delta from iter-2 |
|-----------|-------|-----|-------------------|
| D1. Problem Definition | 103 | 110 | +8 |
| D2. Solution Clarity | 115 | 120 | +4 |
| D3. Industry Benchmarking | 105 | 120 | +12 |
| D4. Requirements Completeness | 101 | 110 | +3 |
| D5. Solution Creativity | 69 | 100 | +24 |
| D6. Feasibility | 91 | 100 | +3 |
| D7. Scope Definition | 76 | 80 | +1 |
| D8. Risk Assessment | 83 | 90 | +6 |
| D9. Success Criteria | 75 | 80 | +4 |
| D10. Logical Consistency | 82 | 90 | +8 |
| **Total** | **900** | **1000** | **+75** |

---

## Top Attacks (Priority Order)

1. **D5 Creativity**: Novelty remains incremental — default priority inference and same-type conflict detection are useful but not structurally novel. Cross-domain inspiration now has two solid analogies (Bazel, K8s) but both from infrastructure/DevOps domain. 31/100 lost. Must find analogies from more distant domains (e.g., database transaction ordering, compiler pipeline scheduling, game engine render pass ordering).

2. **D3 Industry Benchmarking**: Industry analysis improved substantially but still limited to two CI/CD tools. No build system, workflow engine, or dependency management references. 15/120 lost.

3. **D9 Success Criteria**: SC3 "输出覆盖" and SC5 "一致" remain vague after three iterations. Naming coexistence has no SC. 5/80 lost.

4. **D8 Risk Assessment**: "无额外开销" claim in serial execution mitigation is misleading (conflates per-task overhead with wall-clock delay). Config mutation visibility not identified as risk. 7/90 lost.

5. **D10 Logical Consistency**: Config mutation via normalization invisible to users (blindspot-2 persists through three iterations). verify-regression chain-tail and serial execution latency have no SCs. 8/90 lost.

6. **[blindspot-1]** "无额外开销" misleading claim — quote: "happy path 无额外开销——因为 per-surface 任务仍由调度器调度，串行仅影响启动时机" — "startup timing" IS the wall-clock overhead. Must quantify the expected happy-path delay or retract the "无额外开销" claim.
