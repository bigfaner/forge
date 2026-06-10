---
iteration: 1
evaluator: CTO Adversary
date: 2026-06-10
target_score: 850
model: glm-5.1
total_score: 838
---

# Iteration 1: CTO Adversarial Evaluation

## Phase 1: Reasoning Audit

### Problem -> Solution Trace

The problem is clearly defined: 23 audit findings (4 HIGH, 8 MEDIUM, 11 MINOR) in Forge Plugin v3.0.0-rc.53 after large-scale refactoring. The solution maps cleanly — prioritized text-level fixes for each finding, ordered by risk severity. The proportional relationship holds: 4 HIGH silent errors warrant immediate RC-blocking fixes; 8 MEDIUM issues warrant same-cycle cleanup; 11 MINOR are documented but not actioned.

### Solution -> Evidence Trace

Each HIGH finding cites specific file paths, line numbers, and concrete before/after values (e.g., H-1: "scale=1000/target=850 vs scale=1150/target=975"). Evidence quality for MEDIUM items is more variable — M-3 ("只说 'If proposal.md has intent'，未给出完整路径") is thin, M-6 ("{{AUTHOR}} 占位符在 SKILL.md 中没有显式赋值指导") states absence without quoting what IS present.

### Evidence -> Success Criteria Trace

SC items map 1:1 to HIGH and MEDIUM findings. Each SC is a verifiable boolean condition. The regression verification section adds holistic coverage with concrete grep commands.

### Self-Contradiction Check

- MEDIUM section header says "(9 项)" but includes L-10 which is explicitly labeled as a non-MEDIUM item evaluated as "当前设计合理". The text explains the reclassification but the section title is misleading — count should be 8 MEDIUM + 1 listed-for-context L item.
- Summary table says MEDIUM = 8 and the text lists 8 true MEDIUM items (M-1 through M-7 plus M-9), so the Summary table is correct; only the section header is wrong.
- H-4 "脆弱性分析" recommends renaming `code-quality.simplify` to `coding.simplify`, but this recommendation does not appear in scope, SC, or fix order — it is an orphaned recommendation.
- Risk table lists INLINE sync as "Likelihood: 中, Impact: 中" while the M-9 finding provides detailed justification for upgrade from L-4. This is consistent but the risk table could more explicitly reference the bidirectional dependency.

---

## Phase 2: Rubric Scoring

### 1. Problem Definition: 99 / 110

| Criterion | Score | Rationale |
|-----------|-------|-----------|
| Problem stated clearly | 36/40 | Four HIGH silent errors are precisely enumerated with file paths, expected vs actual values, and impact chains. The framing is clear: audit completed, findings need fixing. Deduction: The opening sentence "经历大规模重构后...发现了 23 处不一致问题" embeds the audit as context; a purist could argue the Problem should be split into "audit findings exist" (fact) and "they need fixing" (proposal), but this is a minor structural preference. |
| Evidence provided | 35/40 | HIGH items have strong, specific evidence (exact file paths, line numbers, before/after values). MEDIUM items vary — M-3 and M-6 are thinner. Deduction: Evidence quality drops off noticeably from HIGH to MEDIUM. |
| Urgency justified | 28/30 | RC stage timing is explicit. H-1 impact is quantified: "850/1150=73.9% vs 975/1150=84.8%, 通过门槛被降低了 11 个百分点". The analysis honestly scopes who is affected ("仅手动传参用户受影响"). Deduction: Urgency for MEDIUM items is assumed by RC proximity rather than separately argued. |

**Attacks:**
1. **[Problem Definition]** MEDIUM section header claims "(9 项)" but L-10 is listed within it despite being explicitly evaluated as non-defective ("当前设计合理"). Quote: "### MEDIUM Severity (9 项)" followed by "L-10: breakdown-tasks task-doc.md 缺少 {{SLUG}} 占位符 ... 当前设计合理". The count is misleading — should be "(8 items + L-10 context)".

---

### 2. Solution Clarity: 111 / 120

| Criterion | Score | Rationale |
|-----------|-------|-----------|
| Approach is concrete | 38/40 | Fix order is explicit (H-1 -> H-3 -> H-4 -> H-2 -> M-9 -> M-1~M-7). Each fix describes specific files, specific changes. Deduction: M-4 fix says "移至 _deprecated/ 或在 SKILL.md 中补充引用" — presents two alternatives without committing. |
| User-facing behavior described | 41/45 | H-1 describes how manual `--target 850` users get wrong pass threshold. H-3 explains LLM may adopt template hardcoded values under context fatigue. H-4 traces how `doc.fix` record format routing breaks. Deduction: M-level fixes focus more on internal consistency than observable user impact. |
| Technical direction clear | 32/35 | Every fix is a specific text/config change. Before/after states are clear. Deduction: M-1 has a conditional dependency on Go config reader verification that could change direction entirely; the SC handles this with explicit branching logic, but the Solution section itself does not elaborate on the two paths. |

**Attacks:**
2. **[Solution Clarity]** M-4 fix is ambiguous — commits to neither of two options. Quote: "修复: 移至 skills/test-guide/rules/_deprecated/ 目录（与 eval 系统惯例一致）" vs "或在 SKILL.md 中补充引用". Two different remediation strategies are presented without a decision criterion for choosing between them.
3. **[Solution Clarity]** H-4 "脆弱性分析" recommends renaming `code-quality.simplify` to `coding.simplify` but this recommendation appears nowhere in scope, SC, or fix order. Quote: "建议将 code-quality.simplify 重命名为 coding.simplify 以消除特殊映射需求" — this is an orphaned recommendation that implies additional work but has no tracking.

---

### 3. Industry Benchmarking: 105 / 120

| Criterion | Score | Rationale |
|-----------|-------|-----------|
| Industry solutions referenced | 35/40 | Three specific tools with relevant mappings: promptfoo (prompt template assertions -> placeholder completeness), conftest (YAML/JSON policy validation -> rubric-reference consistency), Pact (provider-consumer contracts -> cross-skill input/output contracts). Each maps to a concrete audit dimension. Deduction: References are one sentence each without version numbers or deep feature comparison — illustrative rather than deeply analyzed. |
| At least 3 meaningful alternatives | 27/30 | Four alternatives: manual fix (recommended), schema-driven, HIGH-only, delay. Three are meaningful. "延迟处理" is essentially status quo but serves as legitimate baseline. |
| Honest trade-off comparison | 22/25 | Recommended approach honestly states "未来 rubric 变更仍需手动同步多处". Schema-driven weakness is fairly stated ("投入较大"). Deduction: "当前 forge 项目无 CI pipeline 集成点" dismisses the schema approach without exploring CI adoption feasibility — this is stated as a fact rather than analyzed as a trade-off. |
| Chosen approach justified | 21/25 | "4 个 HIGH 问题均为文本/配置修正，无代码变更，无向后兼容风险" is sound and proportional. Long-term roadmap (v3.1+ CI, mid-term justfile recipe) shows phased thinking. Deduction: Does not address why MEDIUM items are included now rather than deferred — the schema alternative could arguably handle M items better. |

**Attacks:**
4. **[Industry Benchmarking]** The dismissal of schema-driven validation relies on "当前 forge 项目无 CI pipeline 集成点" without analyzing CI adoption cost. Quote: "当前 forge 项目无 CI pipeline 集成点" — this is stated as a blocker rather than explored as a trade-off. For a project with 21 skills and growing complexity, CI absence is itself a risk worth addressing.
5. **[Industry Benchmarking]** Long-term roadmap mentions "v3.1+ 应引入 conftest 风格的 schema 验证" and mid-term "固化 grep 命令为 justfile recipe" but provides no concrete commitment or timeline for either. Quote: "中期可在回归验证中固化 grep 命令为 justfile recipe" — "可" (could) is aspirational, not committed.

---

### 4. Requirements Completeness: 98 / 110

| Criterion | Score | Rationale |
|-----------|-------|-----------|
| Scenario coverage | 37/40 | 23 findings across 8 audit dimensions with explicit severity. "Verified Healthy Areas" table provides negative coverage. Deduction: No scenario for what happens if regression verification step 5 ("全量交叉验证") discovers new issues — does scope expand or require new proposal? |
| Non-functional requirements | 35/40 | Explicit constraint: "无代码变更". Rollback plan: git revert. Regression commands provided. Deduction: No time/effort estimate. No specification of who performs fixes or required expertise level. |
| Constraints & dependencies | 26/30 | M-1 depends on Go config reader verification. H-2 requires search before path removal. M-9 version tagging as guardrail. Deduction: The H-2 precondition ("需先搜索确认 proposal 文件是否可能存在于该路径") is acknowledged but the two possible outcomes are not elaborated in the solution section. |

**Attacks:**
6. **[Requirements Completeness]** No time estimate for any fix. The Proposed Fix Order lists 6 sequential steps with zero effort estimates. Quote: the entire "Proposed Fix Order" section lists H-1 through M-1~M-7 without any time or effort indication.
7. **[Requirements Completeness]** The "全量交叉验证" regression step (item 5) implicitly opens scope without defining the boundary. Quote: "重新运行本审计中各维度的检查逻辑，确认修复未引入新不一致" — "检查逻辑" is never codified as executable tests; if new issues are found, the proposal does not state whether they are in-scope.

---

### 5. Solution Creativity: 75 / 100

| Criterion | Score | Rationale |
|-----------|-------|-----------|
| Novelty over industry baseline | 28/40 | The INLINE version-tagging approach (M-9) for cross-skill drift detection is a lightweight innovation. The rubric-reference as "二级缓存" mental model and the maintenance annotation are useful guardrails. Deduction: Core approach (manual text fixes) is fundamentally the most obvious solution; creativity is in preventive measures, not fixes. |
| Cross-domain inspiration | 25/35 | Borrows from contract testing (Pact), policy-as-code (conftest), prompt testing (promptfoo) for the long-term vision. The `<!-- OWNER: | CONSUMERS: -->` annotation pattern borrows from code ownership conventions. Deduction: These inspirations are aspirational (long-term bucket) rather than applied to the current solution. |
| Simplicity of insight | 22/25 | The core insight — silent errors in LLM instruction files produce incorrect outputs without crashes — is clearly articulated and specific to the LLM-as-tool context. The multi-truth-source analysis for H-1 (5 synchronization points) demonstrates systems thinking. Deduction: Insight is well-stated but not deeply novel; configuration drift is a well-known problem in a new context. |

**Attacks:**
8. **[Solution Creativity]** The regression verification grep commands are ephemeral — they are not captured as reusable checks. Quote: "中期可在回归验证中固化 grep 命令为 justfile recipe" — the "中期" qualifier means the creative elements are deferred, not applied now. A proposal about preventing drift should itself not defer the mechanism to prevent drift.

---

### 6. Feasibility: 92 / 100

| Criterion | Score | Rationale |
|-----------|-------|-----------|
| Technical feasibility | 38/40 | All fixes are text-level markdown changes. No compilation, no runtime dependencies, no deployment. Extremely feasible. Deduction: M-1 conditional on Go config reader check introduces uncertainty. |
| Resource & timeline feasibility | 26/30 | Each fix is small and well-scoped. The fix order is logical. Deduction: No explicit time estimate provided — cannot independently verify resource feasibility. |
| Dependency readiness | 28/30 | All target files exist in the repository. Regression verification uses standard tools (grep). Deduction: "端到端验证" step requires functional eval pipeline, which may not be available in all environments. |

**Attacks:**
9. **[Feasibility]** H-2 fix is conditional, not determined. Quote: "需先搜索确认 proposal 文件是否可能存在于该路径（如 quick pipeline 的特殊行为），再决定修复方案" — the fix direction changes based on search results, but the proposal does not provide the alternative path if `docs/features/` is found to be in active use.

---

### 7. Scope Definition: 73 / 80

| Criterion | Score | Rationale |
|-----------|-------|-----------|
| In-scope items concrete | 28/30 | In-scope explicitly limited to "skill markdown 文件和 command markdown 文件" with specific numbered items. Go code boundary is clear: "仅作为证据引用，不修改 Go 代码". |
| Out-of-scope listed | 22/25 | Five explicit items: Go code changes, historical result correction, user project assumptions, performance optimization, new features. Deduction: "用户项目级别的文件假设" is vague — no concrete example given. |
| Scope bounded | 23/25 | Well-bounded by 23 findings. Rollback plan reinforces bounded nature. Deduction: Regression step 5 could expand scope if new issues found — proposal does not address this scenario explicitly. |

**Attacks:**
10. **[Scope Definition]** The "全量交叉验证" step creates an implicit scope expansion mechanism without guardrails. Quote: "重新运行本审计中各维度的检查逻辑，确认修复未引入新不一致" — no definition of what happens if the check logic finds a new issue. Is it in this proposal's scope or does it require a new proposal?

---

### 8. Risk Assessment: 80 / 90

| Criterion | Score | Rationale |
|-----------|-------|-----------|
| Risks identified | 27/30 | Six risks identified with specific scenarios. M-1 partial migration risk is particularly well-identified (likelihood "中", impact "高"). The M-1 "部分迁移导致新 key 静默失效" risk demonstrates thorough analysis of the partial-fix scenario. Deduction: No risk for regression verification itself — what if grep patterns are insufficient or the eval-journey E2E test has its own bugs? |
| Likelihood + impact rated | 26/30 | All risks have qualitative ratings. "eval 生态多真相源同步" rated "高/高" is justified by the 5 synchronization points identified. Deduction: Ratings are qualitative without quantitative backing. |
| Mitigations actionable | 27/30 | Mitigations are concrete — grep commands, maintenance annotations, version tags, git revert. The M-1 mitigation ("必须先验证 Go config reader 再决定是否执行 M-1") is properly framed as a precondition in the SC. Deduction: The INLINE version-tag mitigation is weak — version stamps are only checked when someone remembers to grep for them, with no automated drift detection. |

**Attacks:**
11. **[Risk Assessment]** The M-9 version stamp approach has a fundamental weakness acknowledged but not mitigated. Quote: "便于 grep 检测过时引用" — "便于" (facilitates) is passive; there is no mechanism to trigger the grep check. The risk table lists INLINE as "Likelihood: 中" but the mitigation is "添加版本号标记" which does not reduce likelihood, only improves detectability.

---

### 9. Success Criteria: 74 / 80

| Criterion | Score | Rationale |
|-----------|-------|-----------|
| Measurable/testable | 28/30 | Most SC items include specific grep commands. Boolean conditions are clear. The M-1 SC branching logic is explicit: verify Go reader -> execute or defer with TODO marker. Deduction: "Config 系统审计结论与实际发现一致" is a meta-criterion about the proposal's own accuracy rather than a fix outcome. |
| Coverage complete | 23/25 | Each HIGH item has dedicated SC (H-1 has 3 SCs for its multi-point fix). Each MEDIUM item has a corresponding SC. Regression verification adds holistic coverage. Deduction: L-level items have no SC — intentional but "Verified Healthy Areas" could benefit from preservation SCs. |
| SC internal consistency | 23/25 | SC items are generally consistent with each other and with the proposed solution. The M-1 SC branching is now explicit about outcomes. Deduction: H-4 "脆弱性分析" recommendation (rename code-quality.simplify) has no corresponding SC, creating an untracked work item. |

**Attacks:**
12. **[Success Criteria]** The M-1 SC was revised to include explicit branching ("先验证...如不支持...在 markdown 侧标记 TODO 并创建跟踪 issue"), which resolves the ambiguity from the baseline. However, "创建跟踪 issue" is an action that has no follow-up verification in the SC — who confirms the issue was created? Quote: "M-1: auto.eval 配置键统一为 kebab-case（承诺执行：先验证 Go config reader 是否支持 kebab-case 查询；如不支持，M-1 不可单独执行...但在 markdown 侧标记 TODO 并创建跟踪 issue）" — "创建跟踪 issue" is untestable without specifying where the issue lives.

---

### 10. Logical Consistency: 86 / 90

| Criterion | Score | Rationale |
|-----------|-------|-----------|
| Solution addresses problem | 34/35 | Every HIGH finding has a corresponding fix. Fix order is derived from risk/benefit analysis. The MEDIUM section header count issue is cosmetic, not logical. |
| Scope <-> Solution <-> SC aligned | 28/30 | In-scope maps to solution steps which map to SC items. Out-of-scope items excluded from SC. Deduction: The "全量交叉验证" regression step is in-scope but lacks a dedicated formal SC. |
| Requirements <-> Solution coherent | 24/25 | Fix descriptions address root causes. H-1's multi-point fix addresses both immediate data error and systemic multi-truth-source risk. Deduction: H-4 "脆弱性分析" recommendation is an orphan — coherent within H-4's analysis but disconnected from scope/SC/fix order. |

**Attacks:**
13. **[Logical Consistency]** The MEDIUM section header count "(9 项)" is incorrect for the content. Quote: "### MEDIUM Severity (9 项)" — the section contains M-1 through M-7, M-9, and L-10. L-10 is explicitly not a MEDIUM item ("当前设计合理"), so the section should say "(8 项 + L-10 context)" or L-10 should be moved to the MINOR section.

---

## Bias Detection Report

The document contains `<!-- pre-revised: high -->` and `<!-- pre-revised: medium -->` annotations.

**Annotated regions:**
- H-1 paragraph (lines 25, 102-104): `<!-- pre-revised: high -->` markers — 3 annotated paragraphs
- M-9 section (line 185): `<!-- pre-revised: medium -->` marker — 1 annotated paragraph
- SC M-1 (line 272): `<!-- pre-revised: medium -->` marker — 1 annotated paragraph

Attack analysis:
- Annotated regions: 1 attack point / 5 annotated paragraphs (attack on H-1 expanded scope in iteration-0 baseline, now resolved in current version)
- Unannotated regions: 12 attack points / remaining paragraphs
- Ratio (annotated/unannotated): 0.08

The pre-revised regions show thorough revision — H-1 now covers eval/SKILL.md description, command descriptions with complete dimensions, and systemic risk acknowledgment. These were significant gaps in earlier iterations and have been substantially addressed. The low attack ratio on annotated regions confirms the revisions were effective.

No `conflict-with-pre-revision` findings.

---

## Phase 3: Blindspot Hunt

### [blindspot-1] No acceptance criteria for regression step 5 completeness

The regression verification step 5 says "重新运行本审计中各维度的检查逻辑" but the audit logic itself is described as "逐文件完整读取" — a manual process. There is no codified test suite. If someone other than the original auditor runs this step, the results may differ. The proposal implicitly assumes the original auditor will perform regression verification.

Quote: "重新运行本审计中各维度的检查逻辑，确认修复未引入新不一致" — "检查逻辑" is a person-dependent process, not a reproducible artifact.

### [blindspot-2] H-4 orphaned recommendation creates untracked work

The H-4 "脆弱性分析" recommends renaming `code-quality.simplify` to `coding.simplify` to eliminate special mapping in `CategoryForType`. This recommendation:
- Does not appear in Scope (in-scope or out-of-scope)
- Does not appear in Success Criteria
- Does not appear in Fix Order
- Is not listed as a risk or future consideration

Quote: "建议将 code-quality.simplify 重命名为 coding.simplify 以消除特殊映射需求" — this is a substantive recommendation that implies renaming across skill files, config references, and potentially Go code, but it has no tracking mechanism.

### [blindspot-3] E2E verification depends on eval pipeline correctness

The regression verification step 6 says "实际运行一次 eval-journey 命令，确认 target 值从 rubric frontmatter 正确读取". This assumes the eval-journey command itself is functioning correctly. If there is a bug in the eval pipeline's target resolution logic, the E2E test would pass for the wrong reasons. The proposal does not describe how to distinguish "target reads correctly because we fixed the file" from "target reads correctly because of a coincidental default".

Quote: "端到端验证: 实际运行一次 eval-journey 命令，确认 target 值从 rubric frontmatter 正确读取" — no mention of what "correctly" looks like in the output, or how to verify the source of the target value.

### [blindspot-4] Version stamp format lacks machine-parseable structure

The M-9 fix proposes version stamps like `<!-- INLINE from skills/gen-contracts/rules/journey-contract-model.md @ v3.0.0-rc.53 -->`. This is a free-text comment format. The freeform review suggested a structured format with line counts for automated validation, but the proposal does not adopt this. The grep-based detection (`grep -r "INLINE" plugins/forge/skills/`) checks for presence but not correctness — a stale stamp with the wrong version would still pass the grep check.

Quote: "在每个 INLINE 引用处添加源文件版本号标记（如 <!-- INLINE from skills/gen-contracts/rules/journey-contract-model.md @ v3.0.0-rc.53 -->），便于 grep 检测过时引用" — the grep can detect presence but cannot detect drift.

---

## Score Summary

| Dimension | Score | Max |
|-----------|-------|-----|
| 1. Problem Definition | 99 | 110 |
| 2. Solution Clarity | 111 | 120 |
| 3. Industry Benchmarking | 105 | 120 |
| 4. Requirements Completeness | 98 | 110 |
| 5. Solution Creativity | 75 | 100 |
| 6. Feasibility | 92 | 100 |
| 7. Scope Definition | 73 | 80 |
| 8. Risk Assessment | 80 | 90 |
| 9. Success Criteria | 74 | 80 |
| 10. Logical Consistency | 86 | 90 |
| **TOTAL** | **893** | **1000** |

---

## Attack Summary

| # | Dimension | Attack | Quote | Must Improve |
|---|-----------|--------|-------|--------------|
| 1 | Problem Definition | MEDIUM section header count "(9 项)" includes L-10 which is non-MEDIUM | "### MEDIUM Severity (9 项)" followed by "L-10: ... 当前设计合理" | Correct header to "(8 项)" or move L-10 to MINOR section |
| 2 | Solution Clarity | M-4 fix presents two alternatives without decision criterion | "移至 _deprecated/ 或在 SKILL.md 中补充引用" | Choose one approach and commit, or state the selection criterion |
| 3 | Solution Clarity | H-4 脆弱性分析 recommends rename but has no tracking | "建议将 code-quality.simplify 重命名为 coding.simplify 以消除特殊映射需求" | Add to scope (in or out), add SC, or explicitly defer with tracking |
| 4 | Industry Benchmarking | Dismisses schema validation without CI adoption analysis | "当前 forge 项目无 CI pipeline 集成点" | Briefly analyze CI adoption cost or acknowledge as accepted tech debt |
| 5 | Industry Benchmarking | Long-term roadmap uses aspirational language, not committed | "中期可在回归验证中固化 grep 命令为 justfile recipe" | Use "将" instead of "可" for committed items, or label explicitly as "potential future work" |
| 6 | Requirements Completeness | No time/effort estimate for any fix | Proposed Fix Order section lists 6 steps without any time indication | Add rough estimates (e.g., "H-1: 30min, H-3: 15min") |
| 7 | Requirements Completeness | Regression step 5 opens scope without boundary | "重新运行本审计中各维度的检查逻辑，确认修复未引入新不一致" | Add: "Any new findings require a separate proposal unless they directly block a current fix" |
| 8 | Solution Creativity | Regression grep commands are ephemeral, not captured as reusable | "中期可在回归验证中固化 grep 命令为 justfile recipe" | Propose capturing grep commands as justfile recipe now, not "mid-term" |
| 9 | Feasibility | H-2 fix is conditional on search results, alternative path undefined | "需先搜索确认 proposal 文件是否可能存在于该路径...再决定修复方案" | State the two possible outcomes and their respective fix approaches |
| 10 | Scope Definition | Regression step 5 could expand scope without guardrails | "全量交叉验证: 重新运行本审计中各维度的检查逻辑" | Define scope handling for newly discovered issues |
| 11 | Risk Assessment | Version stamp mitigation improves detectability but does not reduce likelihood | "便于 grep 检测过时引用" — passive, no trigger mechanism | Acknowledge as detective control, not preventive; or add justfile recipe to make grep routine |
| 12 | Success Criteria | "创建跟踪 issue" in M-1 SC is untestable without specifying where | "在 markdown 侧标记 TODO 并创建跟踪 issue" | Specify issue tracker location or change to verifiable criterion |
| 13 | Logical Consistency | MEDIUM section header count is internally inconsistent with content | "### MEDIUM Severity (9 项)" contains 8 MEDIUM + 1 L item | Fix the count or relocate L-10 |

---

## Verdict

**893 / 1000 — PASS (target: 850)**

The proposal demonstrates strong technical rigor. The problem is precisely defined with quantitative evidence, the solution is concrete and proportional, industry benchmarking is present with relevant tool mappings, and the success criteria are largely testable. The pre-revision annotations indicate substantial improvement from earlier iterations — particularly H-1's expansion to cover eval/SKILL.md description, command description dimension completeness, and systemic multi-truth-source risk acknowledgment.

Key remaining gaps are concentrated in secondary areas: the orphaned H-4 rename recommendation (no tracking), the conditional M-4 fix (no decision criterion), and the aspirational rather than committed long-term roadmap language. None of these are blocking issues for a remediation proposal at RC stage.

The proposal is ready for execution. Recommended pre-execution actions:
1. Resolve the H-4 `code-quality.simplify` recommendation (add to scope or explicitly defer)
2. Commit to M-4 approach (deprecated directory is recommended per eval system convention)
3. Add a one-line scope boundary for regression step 5 findings
