# Adversarial Evaluation: forge-skill-audit Proposal (Baseline)

**Evaluator**: Adversary Agent (Baseline)
**Date**: 2026-06-10
**Iteration**: 0 (Baseline)
**Proposal**: `docs/proposals/forge-skill-audit/proposal.md`

---

## Phase 1: Reasoning Audit

### Problem -> Solution Trace
The problem (23 audit findings in Forge Plugin v3.0.0-rc.53, 4 HIGH silent errors) maps cleanly to the solution (prioritized text-level fixes with regression verification). The solution is proportional to the problem scope.

### Solution -> Evidence Trace
Each HIGH item cites specific file paths, line numbers, and concrete values (e.g., "scale=1000/target=850 vs scale=1150/target=975"). MEDIUM items vary in specificity — M-3 and M-6 are less detailed.

### Evidence -> Success Criteria Trace
SC items are mapped 1:1 to each HIGH and MEDIUM finding. Each SC is phrased as a verifiable boolean condition.

### Self-Contradiction Check
- Summary table says 8 MEDIUM findings, but MEDIUM section lists 9 items (M-1 through M-7, M-9, plus L-10 which is correctly noted as downgraded). The text explains "L-4 upgraded to M-9, original M-8 downgraded to L-10", so the count is actually 8 MEDIUM if M-8 was removed and replaced by M-9. This is internally consistent after careful reading but the section title "MEDIUM Severity (9 items)" is a labeling error — it should say 8 items (M-1 through M-7 plus M-9, with L-10 listed separately as a non-MEDIUM item).
- M-1 scope says "only modify skill markdown" but SC for M-1 has a conditional that may defer to Go code changes. This is handled, not contradicted.

---

## Phase 2: Rubric Scoring

### 1. Problem Definition (110 pts)

**Problem stated clearly (38/40)**: The problem is precisely defined — 4 HIGH silent errors with specific file paths and impact descriptions. Deduction: The framing relies on an internal audit that is not independently verifiable from the proposal alone; the reader must trust the audit results. The statement "通过八维度深度审计...发现了 23 处不一致问题" describes the methodology but the audit itself is not reproducible from this document.

**Evidence provided (35/40)**: HIGH items have strong evidence (exact file paths, line numbers, actual vs expected values). MEDIUM items are less uniformly evidenced — M-6 ("{{AUTHOR}} 占位符在 SKILL.md 中没有显式赋值指导") states the absence but doesn't quote what IS present. Deduction: Evidence quality is inconsistent across findings.

**Urgency justified (28/30)**: RC stage timing is clear. The H-1 impact analysis (73.9% vs 84.8% pass threshold, 11 percentage point difference) is quantitative and compelling. Deduction: Urgency for MEDIUM items is assumed by RC proximity rather than explicitly argued.

**Dimension Total: 101/110**

---

### 2. Solution Clarity (120 pts)

**Approach is concrete (38/40)**: Fix order is explicit (H-1 -> H-3 -> H-4 -> H-2 -> M-9 -> M-1~M-7). Each fix is described at file-level granularity. Deduction: M-4 fix is ambiguous — "移至 _deprecated/ 或在 SKILL.md 中补充引用" presents two alternatives without committing to one.

**User-facing behavior described (42/45)**: The proposal describes how each fix changes LLM behavior (e.g., "LLM 可能直接采用模板值而非覆盖" for H-3). Deduction: User-facing behavior for M-level fixes is less explicitly described — the proposal focuses on internal consistency rather than observable user impact for M items.

**Technical direction clear (34/35)**: Every fix is a specific text/config change with before/after states. No ambiguity about what files to touch. Deduction: M-1 has a dependency on Go config reader verification that could change the fix direction entirely.

**Dimension Total: 114/120**

---

### 3. Industry Benchmarking (120 pts)

**Industry solutions referenced (32/40)**: Three tools are referenced — promptfoo (prompt template assertions), conftest (YAML/JSON policy validation), Pact (contract testing). These are relevant to the problem domain. Deduction: References are brief (one sentence each) and don't include version numbers or specific feature comparisons; they're illustrative rather than deeply analyzed.

**At least 3 meaningful alternatives (25/30)**: Four alternatives are presented in the comparison table. Three are meaningful (manual fix, schema-driven, fix HIGH only). "延迟处理" is a weak alternative — it's essentially "do nothing" which is a straw-man. However, it does serve as a legitimate status-quo baseline, so the deduction is moderate.

**Honest trade-off comparison (22/25)**: The recommended approach honestly states "未来 rubric 变更仍需手动同步多处" as a weakness. The schema-driven approach's weakness ("需要编写 schema 定义和验证规则，投入较大") is fair. Deduction: The "当前 forge 项目无 CI pipeline 集成点" weakness for schema-driven approach is stated without exploring whether CI could be added — this is dismissive rather than analytical.

**Chosen approach justified (23/25)**: The justification "4 个 HIGH 问题均为文本/配置修正，无代码变更，无向后兼容风险" is sound and proportional. The long-term roadmap (v3.1+ CI checks) shows awareness of the tactical limitation. Deduction: The justification doesn't address why MEDIUM items are included in this round rather than deferred — the schema alternative could handle M items as well.

**Dimension Total: 102/120**

---

### 4. Requirements Completeness (110 pts)

**Scenario coverage (36/40)**: 23 findings across 8 audit dimensions, with explicit severity classification and fix instructions. The "Verified Healthy Areas" table provides negative-scenario coverage. Deduction: No coverage of what happens if new inconsistencies are discovered during the fix process — is the scope open to expansion?

**Non-functional requirements (32/40)**: The proposal explicitly states "无代码变更" as a constraint, regression verification commands are provided, and rollback plan exists ("git revert"). Deduction: No time/effort estimate for the fixes. No mention of who will perform the fixes or required expertise level. No performance impact analysis (though arguably N/A for text changes).

**Constraints & dependencies (26/30)**: Dependencies are well-documented — M-1 depends on Go config reader verification, H-2 requires search confirmation before path removal, M-9 has version tagging as dependency. Deduction: The dependency between M-1 and Go code changes is noted but the decision tree (what if Go reader doesn't support kebab-case?) is only partially explored.

**Dimension Total: 94/110**

---

### 5. Solution Creativity (100 pts)

**Novelty over industry baseline (28/40)**: The INLINE version-tagging approach (M-9) for detecting stale cross-skill references is a lightweight innovation. The rubric-reference.md maintenance annotation is a simple but effective guard. Deduction: The core approach (manual text fixes) is fundamentally uncreative — it's the most obvious solution. The creativity is in the regression verification and preventive measures, not in the fixes themselves.

**Cross-domain inspiration (22/35)**: The proposal borrows from contract testing (Pact), policy-as-code (conftest), and prompt testing (promptfoo) for the long-term vision. Deduction: These inspirations are relegated to the "long-term" bucket and don't influence the actual proposed solution — they're aspirational rather than applied.

**Simplicity of insight (23/25)**: The key insight — that silent errors in LLM instruction files produce incorrect outputs without crashes — is clearly articulated. The rubric-reference as "二级缓存" framing is an elegant mental model. Deduction: The insight is well-stated but not deeply novel; it's a known configuration management problem in a new context.

**Dimension Total: 73/100**

---

### 6. Feasibility (100 pts)

**Technical feasibility (38/40)**: All fixes are text-level changes to markdown files. No compilation, no runtime dependencies, no deployment concerns. Extremely feasible. Deduction: M-1 has a conditional that could block execution (Go config reader verification), which introduces uncertainty.

**Resource & timeline feasibility (25/30)**: The fix order is logical and each fix is well-scoped. Deduction: No explicit time estimate is provided. The proposal doesn't state whether one person or multiple people will execute these fixes, or whether the fixes need to be coordinated.

**Dependency readiness (28/30)**: All files to be modified exist in the repository and are accessible. The regression verification commands use standard tools (grep). Deduction: The "端到端验证" step (running eval-journey) depends on a functional eval pipeline, which isn't guaranteed to be available in all environments.

**Dimension Total: 91/100**

---

### 7. Scope Definition (80 pts)

**In-scope items concrete (28/30)**: In-scope is explicitly limited to "skill markdown 文件和 command markdown 文件" with specific numbered items. The boundary with Go code is clearly drawn. Deduction: "M-1 的 config key 重命名仅修改 skill markdown 中的引用" — the partial migration nature of M-1 creates ambiguity about what constitutes "complete" execution.

**Out-of-scope listed (23/25)**: Five explicit out-of-scope items including Go code changes, historical result correction, user project assumptions, performance optimization, and new features. Deduction: "用户项目级别的文件假设" is vague — what specific assumptions? A concrete example would strengthen this.

**Scope bounded (22/25)**: The scope is well-bounded by the audit results (23 findings). The rollback plan reinforces the bounded nature. Deduction: The "全量交叉验证" regression step (item 5 in Regression Verification) could expand scope if new issues are found — the proposal doesn't address this scenario.

**Dimension Total: 73/80**

---

### 8. Risk Assessment (90 pts)

**Risks identified (28/30)**: Six risks are identified with specific scenarios. The M-1 partial migration risk is particularly well-identified. Deduction: No risk identified for the regression verification itself — what if the grep patterns are insufficient or the end-to-end eval-journey test has its own bugs?

**Likelihood + impact rated (26/30)**: All risks have likelihood/impact ratings. The M-1 partial migration risk correctly identifies "中" likelihood and "高" impact. Deduction: Ratings are qualitative (低/中/高) without quantitative backing. The "高" likelihood for "eval 生态多真相源同步" is asserted without justification — why is it "高" rather than "中"?

**Mitigations actionable (27/30)**: Mitigations are concrete — grep commands, maintenance annotations, version tags, rollback via git revert. Deduction: The mitigation for M-1 partial migration ("必须先验证 Go config reader 再决定是否执行 M-1") is a precondition, not a mitigation — it doesn't reduce risk, it defers the decision.

**Dimension Total: 81/90**

---

### 9. Success Criteria (80 pts)

**Measurable/testable (28/30)**: Most SC items include specific grep commands for verification. Boolean conditions are clear. Deduction: M-1 SC includes a conditional ("前提：验证 Go config reader 是否支持 kebab-case") which makes it partially untestable until the prerequisite is resolved.

**Coverage complete (23/25)**: Each HIGH item has dedicated SC. Each MEDIUM item has a corresponding SC. The regression verification section adds holistic coverage. Deduction: L-level items have no SC at all — this is intentional (they're not being fixed) but the "Verified Healthy Areas" could benefit from preservation SC (ensuring healthy areas remain healthy post-fix).

**SC internal consistency (22/25)**: SC items are generally consistent with each other and with the proposed solution. Deduction: The SC for M-1 creates a two-path outcome (execute now if Go supports kebab-case, defer if not) that could split the proposal into partial-completion scenarios. The SC doesn't define what happens to the M-1 SC if deferred — does it become "not applicable" or "transferred to next iteration"?

**Dimension Total: 73/80**

---

### 10. Logical Consistency (90 pts)

**Solution addresses problem (33/35)**: Every HIGH finding has a corresponding fix in the solution section. The fix order is logically derived from risk/benefit analysis. Deduction: MEDIUM section header says "9 items" but includes L-10 (which is explicitly not a MEDIUM item) — this is a presentation error that doesn't affect the solution logic but undermines document precision.

**Scope <-> Solution <-> SC aligned (28/30)**: In-scope items map to solution steps which map to SC items. Out-of-scope items are consistently excluded from SC. Deduction: The "全量交叉验证" regression step is in-scope but doesn't have a dedicated SC — it's mentioned in the Regression Verification section but not as a formal success criterion.

**Requirements <-> Solution coherent (23/25)**: The fix descriptions directly address the root causes identified in each finding. H-1's fix (data sync + maintenance annotation) addresses both the immediate error and the systemic multi-truth-source risk. Deduction: The H-4 "脆弱性分析" suggests renaming code-quality.simplify to coding.simplify but this recommendation is not included in scope or SC — it's an orphaned recommendation.

**Dimension Total: 84/90**

---

## Phase 3: Blindspot Hunt

### [blindspot-1] No acceptance criteria for regression verification completeness
The Regression Verification section lists 6 verification steps, but there's no definition of what "pass" looks like for steps 5 and 6. Step 5 says "重新运行本审计中各维度的检查逻辑，确认修复未引入新不一致" — but the audit logic itself is not codified or reproducible. It was a manual process. If a new inconsistency is found during regression, the proposal doesn't define whether this blocks release or triggers scope expansion.

**Quote**: "重新运行本审计中各维度的检查逻辑，确认修复未引入新不一致" — "检查逻辑" is never defined as executable tests or scripts.

### [blindspot-2] M-8 downgrade to L-10 lacks documented rationale
The proposal mentions "原 M-8 降级为 L-10（设计合理，非缺陷）" but L-10's entry in the MEDIUM section is simply "评估：breakdown-tasks 有更多输入文档...当前设计合理". The downgrade rationale is one sentence. For a proposal that claims rigorous 8-dimension audit, the downgrade of a finding to "by design" deserves more justification — was the original M-8 classification an error, or did the understanding evolve?

**Quote**: "原 M-8 降级为 L-10（设计合理，非缺陷）" — no supporting analysis for the downgrade.

### [blindspot-3] H-2 fix has a precondition that could expand scope
H-2 fix says "需先搜索确认 proposal 文件是否可能存在于该路径（如 quick pipeline 的特殊行为），再决定修复方案". This means the fix for H-2 is not yet determined — it's contingent on search results. If the search reveals that docs/features/ IS used by some pipeline, the fix direction changes entirely. The proposal presents H-2 as a straightforward "remove dead path" but it might not be.

**Quote**: "需先搜索确认 proposal 文件是否可能存在于该路径（如 quick pipeline 的特殊行为），再决定修复方案" — the fix is conditional, not determined.

### [blindspot-4] No sequencing dependency analysis
The fix order (H-1 -> H-3 -> H-4 -> H-2 -> M-9 -> M-1~M-7) is presented as a simple priority sequence. But there are dependencies between fixes that aren't analyzed: M-9 (INLINE version tags) could affect the grep patterns used in regression verification. M-1 (config key rename) affects the behavior tested in end-to-end verification. The proposal doesn't analyze whether fixes can be parallelized or must be strictly sequential.

**Quote**: The Proposed Fix Order section lists fixes sequentially but never states whether this is a strict dependency order or merely a priority recommendation.

### [blindspot-5] Success Criteria missing for "verified healthy areas"
The proposal devotes a table to "Verified Healthy Areas" (Surface system, Hook system, etc.) but has no SC ensuring these areas remain healthy after fixes. A fix to one skill could inadvertently break a previously-healthy area. The regression verification step 5 partially addresses this, but it's not a formal SC.

**Quote**: The "Verified Healthy Areas" table lists 7 healthy areas, but none appear in the Success Criteria section.

---

## Score Summary

| Dimension | Score | Max |
|-----------|-------|-----|
| 1. Problem Definition | 101 | 110 |
| 2. Solution Clarity | 114 | 120 |
| 3. Industry Benchmarking | 102 | 120 |
| 4. Requirements Completeness | 94 | 110 |
| 5. Solution Creativity | 73 | 100 |
| 6. Feasibility | 91 | 100 |
| 7. Scope Definition | 73 | 80 |
| 8. Risk Assessment | 81 | 90 |
| 9. Success Criteria | 73 | 80 |
| 10. Logical Consistency | 84 | 90 |
| **TOTAL** | **886** | **1000** |

---

## Attacks (Top Issues)

1. **Industry Benchmarking**: Alternatives analysis is shallow — the three industry tools (promptfoo, conftest, Pact) are mentioned with one-sentence descriptions and relegated to "long-term". The schema-driven alternative is the only serious competitor, and its dismissal relies on "无 CI pipeline 集成点" without exploring whether adding CI is feasible. **Quote**: "当前 forge 项目无 CI pipeline 集成点" — this dismisses the most promising alternative without analysis of CI adoption cost. **Must improve**: Provide a brief cost-benefit analysis for adding basic CI, or acknowledge the schema approach as the recommended v3.1 path with a concrete plan.

2. **Solution Creativity**: The proposed solution is entirely manual text editing — the least creative approach possible. While appropriate for the urgency, the creative elements (INLINE version tags, maintenance annotations) are minor. **Quote**: "按优先级执行文本级修复（无代码变更），每项修复后运行回归验证" — the entire solution is "find and replace with grep verification". **Must improve**: Propose at least one structural improvement beyond text fixes — e.g., a template validation script that could be added to prevent recurrence.

3. **Requirements Completeness**: No time/effort estimate, no resource allocation, no required expertise specification. **Quote**: The Proposed Fix Order lists 6 sequential steps but provides zero time estimates or resource requirements. **Must improve**: Add even a rough estimate (e.g., "estimated 2-3 hours for all HIGH fixes, 1-2 hours for MEDIUM") and specify who should perform the fixes.

4. **Success Criteria**: M-1 SC has a conditional that creates ambiguity about completion criteria. **Quote**: "M-1: auto.eval 配置键统一为 kebab-case（前提：验证 Go config reader 是否支持 kebab-case 查询；如不支持，将 M-1 与 Go alias 绑定为原子操作，推迟至 Go 代码变更窗口执行）" — this SC has three possible outcomes (execute, defer, bind-to-Go-task) which makes acceptance testing ambiguous. **Must improve**: Split M-1 into an investigation SC and an execution SC, or define explicit pass/fail/waive conditions.

5. **Logical Consistency**: The MEDIUM section header claims "9 items" but includes L-10 which is explicitly not a MEDIUM item. **Quote**: "### MEDIUM Severity (9 项)" followed by items including "L-10: breakdown-tasks task-doc.md 缺少 {{SLUG}} 占位符" which is labeled as L-10 (MINOR) and evaluated as "当前设计合理". **Must improve**: Correct the section header or move L-10 to the MINOR section to eliminate the classification inconsistency.

6. **Scope Definition**: The regression verification step 5 ("重新运行本审计中各维度的检查逻辑") implicitly opens scope to new findings but the proposal doesn't define the boundary. **Quote**: "全量交叉验证: 重新运行本审计中各维度的检查逻辑，确认修复未引入新不一致" — if new issues are found, are they in-scope or do they require a new proposal? **Must improve**: Add explicit scope handling: "Any new findings discovered during regression verification will be documented but require a separate proposal unless they directly block a current fix."
