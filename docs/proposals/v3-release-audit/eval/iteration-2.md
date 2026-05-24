---
iteration: 2
reviewer: CTO Adversary
date: 2026-05-24
status: completed
based_on: proposal.md (revision after iteration-1 feedback)
---

# Iteration 2: CTO Adversary Rubric Scoring

## Revision Delta Analysis

Iteration 1 scored **698/1000**. The proposal has been revised in response to that report. Key changes detected:

1. **Headline count fixed**: "27 个偏差项" corrected to "50 个偏差项" — resolved the self-contradiction.
2. **Industry Benchmarking expanded**: Added Kubernetes, Rust, Spring Boot references with concrete practices; expanded comparison table from 4 to 5 rows with named sources.
3. **Target State added**: README and ARCHITECTURE.md target structures now described.
4. **Innovation section added**: Cross-domain borrowing from financial audit and supply-chain traceability.
5. **Urgency strengthened**: Added beta user count (3), open issues count (5), concrete version mismatch (2.16.1 vs 3.0.0-rc.24).
6. **Audit methodology expanded**: 5-dimension method now described with measurement pipeline (doc self-check -> impl comparison -> contract validation -> compliance check -> dependency graph).
7. **P0 rollback plan added**: P0.4 has explicit rollback strategy (git stash, agent regression test, fallback to P1).
8. **Timeline revised**: Total from ~5h to ~6.5h; SKILL.md splitting from ~1.5h to ~3h.
9. **P0.5 decision resolved**: Option (c) now recommended with rationale; no longer an open decision.
10. **P1.12 bounded**: Explicit line budget (<=180 lines, <=20 lines per subsystem) with escalation clause.
11. **Risk table expanded**: 6->7 risks; cascade dependency risk and scope creep risk added.
12. **Success criteria expanded**: Now split into P0/P1 sections; P1 criteria added.
13. **Scope boundary clarified**: "不修改 Go 运行时代码" replaces vague "不修改运行时代码"; P0.4 runtime impact explicitly acknowledged.

## Phase 1: Reasoning Audit (Revision Quality)

The proposal argues: Forge v3.0.0 docs are systemically drifted from implementation -> 50 audit findings across 5 dimensions -> tiered remediation (P0/P1/P2) -> ~6.5h effort -> no Go runtime code changes.

**Chain trace:**

1. Problem claim: docs-impl drift is systematic and severe. Evidence: 50 items with severity breakdown and detailed methodology. **Chain holds** — methodology description enables reproducibility; severity thresholds are now inferable from the table.
2. Urgency: v3.0.0 major release, 3 beta users, 5 open issues, version 2.16.1 in README. **Chain holds** — concrete, project-specific, no longer resting on general principle.
3. Solution: tiered fix by severity. **Chain holds with minor gap** — the scope boundary is now more honest ("Go 运行时代码" not "运行时代码"), but the tension between "SKILL.md splitting changes agent runtime behavior" and the scope claim remains. The proposal acknowledges this in the P0.4 description ("此项改变 agent 运行时行为") which is commendable transparency, but the Solution section header still says "不修改 Go 运行时代码" — technically true but semantically misleading for anyone scanning headers.
4. Alternatives: 5 options with named sources (Kubernetes, Rust, Spring Boot). **Chain holds** — alternatives are now meaningfully differentiated. However, the CI automation alternative is listed with a con of "需编写验证脚本（~2h），无 CI 基础设施" — this underestimates the effort and conflates two separate concerns (script writing vs CI infrastructure availability).
5. Feasibility: ~6.5h, no external deps. **Chain improved but still optimistic** — eval SKILL.md splitting at 3h is more realistic than 1.5h, but P1.12 (9 subsystem docs, <=180 lines) has no time estimate in the timeline breakdown. If P1 is executed before release, this adds 1.5-2h.
6. Success criteria: split into P0 (7 criteria) and P1 (6 criteria). **Chain significantly improved** — P1 now has criteria. Gap: P2 items (20 items, ~40% of findings) have no criteria, but this is acceptable for "发布后迭代."

## Phase 2: Dimension Scoring

### 1. Problem Definition: 96/110

**Problem stated clearly (38/40):** "核心文档（README.md、ARCHITECTURE.md）与实际实现存在系统性偏差——过时的计数、不存在的组件、断裂的交叉引用和死代码" is unambiguous. Skills/ SKILL.md drift is explicitly mentioned. Two readers would converge on the same interpretation. Deduction: the structural cause of drift (v2->v3 migration fallout) is still not analyzed. Understanding *why* helps prevent recurrence and is relevant to the "漂移预防" NFR.

**Evidence provided (38/40):** The 5-dimension audit table with severity breakdown is strong. The methodology is now described: "文档自检（逐行扫描→grep验证）→ 实现-比对（提取版本/计数/路径→diff源码）→ 契约验证（forge --help vs SKILL.md）→ 规范合规（wc -l + grep引用）→ 依赖图（skills→rules→templates有向图，入度=0为孤儿）." This is reproducible. The "Assumptions Challenged" table provides excellent falsification evidence. Deduction: severity threshold definitions remain implicit — what distinguishes Critical from Major? The reader must infer from the examples.

**Urgency justified (20/30):** Significantly improved. "3 个 beta 用户", "5 个 open issues 直接源于文档-实现偏差", "README 版本号仍停留在 2.16.1（实际已 3.0.0-rc.24）" are concrete. Deduction: (a) no release date deadline is stated — is there a target date? (b) the 5 open issues are referenced but not linked — which issues? This weakens verifiability. (c) The cost of delay is still not quantified: what happens to those 3 beta users if the release ships with current docs?

### 2. Solution Clarity: 100/120

**Approach is concrete (38/40):** P0/P1/P2 tiering with explicit execution ordering and detailed item descriptions. The execution order dependency (P0.4 -> P0.5 -> P0.2 -> P0.3 -> P0.1) is well-reasoned. Each P0 item specifies what changes, where, and why. Deduction: P1 items 10-11 and most P2 items are still single-line descriptions that could mean different things to different implementers.

**User-facing behavior described (35/45):** Improved with the "Target State" section. README target structure is outlined: "标题→安装→快速开始→命令速查→技能列表→任务类型表→架构→贡献→License." ARCHITECTURE.md target: "1 agent，无 PostToolUse，1000 分制." Deduction: the Target State describes structural targets but not content targets. What will the "命令速查" section contain? What will the "技能列表" look like? The target is a skeleton, not a mockup. For the most complex item (SKILL.md splitting), no target structure is shown — what does the post-split eval SKILL.md look like? What goes into the extracted rules/ files?

**Technical direction clear (27/35):** SKILL.md splitting direction (extract to rules/), CLI reference fix (rename commands), README rewrite (version/count/type corrections). The README rewrite direction is still somewhat vague — "全面重写" is used but the Target State suggests a structural reorganization, not just patching claims. The gap between "fix factual errors" and "restructure the document" is unaddressed.

### 3. Industry Benchmarking: 85/120

**Industry solutions referenced (32/40):** Significantly improved. Kubernetes (auto-generated API docs, CI verification), Rust (RFC process for doc updates), Spring Boot (docusaurus validated links + snippet tests), remark-lint (Markdown linting) are now cited with specific practices. The Forge-specific insight ("文档是 agent 运行时输入，偏差=功能 bug") is a genuine differentiator. Deduction: citations are at practice level, not source level — no links to K8s doc generation pipeline, Rust RFC process docs, or Spring Boot CI configs. The reader cannot follow up.

**At least 3 meaningful alternatives (22/30):** 5 alternatives including do-nothing. The CI automation alternative ("CI 文档验证自动化") is now a genuinely different strategy. However, "仅修复 README", "仅修复 Critical", and "分层全量修复" are still scope variations of the same manual approach. The "deferred" CI option partially addresses this — it acknowledges a different strategy exists but defers it. Missing: automated doc generation from code (e.g., generating CLI reference from Go cobra command definitions) is a genuinely different strategy not considered.

**Honest trade-off comparison (16/25):** Improved. "需编写验证脚本（~2h），无 CI 基础设施" for CI automation is specific. Deduction: (a) the ~2h estimate for CI validation scripts is likely an underestimate — writing robust doc-impl validation requires parsing Go code, understanding the skill structure, and handling edge cases. (b) The "do nothing" alternative is now better presented but still has only one pro ("零成本；多数项目带着过期文档发布并存活") — the insight that "most projects survive with stale docs" is actually valuable context, making this less of a straw man than in iteration 1.

**Chosen approach justified against benchmarks (15/25):** The justification is "问题间存在依赖，分批修复不如一次性对齐." The Forge-specific insight ("文档是 agent 运行时输入") is used to reject do-nothing and justify comprehensive coverage. This is better than iteration 1 but still self-referential — it argues for *scope* (do everything) not *method* (why manual fix vs. automated). The proposal does not explain why it rejects the CI automation approach other than "deferred" — a stronger argument would compare manual-fix-now + CI-later vs. CI-now.

### 4. Requirements Completeness: 92/110

**Scenario coverage (35/40):** Four key scenarios with current-state failures: user reading README, contributor reading ARCHITECTURE.md, agent calling CLI, agent loading orphaned rules. Deduction: (a) still missing: what happens when a user follows incorrect task type names — does the CLI reject them or produce wrong output? (b) no explicit regression scenario ("after fix, verify nothing broke") — though NFR addresses this partially. (c) no scenario for the most dangerous item: SKILL.md splitting causing agent context overflow or misrouting.

**Non-functional requirements (32/40):** Now includes 5 NFRs including "可维护性" (drift prevention assertions) and "可回归性" (P0 end-to-end regression). This is a significant improvement. Deduction: (a) performance NFR still missing — does adding Load directives for 15 orphaned rules increase agent context consumption measurably? (b) the drift prevention NFR ("README 计数与 ls skills/ | wc -l 断言") is good but speculative — no concrete mechanism is proposed (shell script? CI check? Makefile target?).

**Constraints & dependencies (25/30):** Three well-specified constraints referencing actual project conventions. The P0.4/P0.5 dependency ordering is now explicit. Deduction: the constraint "SKILL.md <= 350 行" is listed as a constraint but the proposal's P0.4 exists specifically because this constraint was violated. This should be framed as a compliance gap, not a constraint — the constraint was already there, it was violated.

### 5. Solution Creativity: 55/100

**Novelty over industry baseline (20/40):** The "5维度交叉审计" framework has minor novelty over standard doc audits. The explicit comparison table in Industry Benchmarking shows awareness of industry approaches but the solution is fundamentally manual remediation. The innovation highlights section claims cross-domain inspiration but the execution is standard fix-by-severity.

**Cross-domain inspiration (15/35):** Now explicitly claims borrowing from "财务审计重要性阈值" (financial audit materiality thresholds) and "供应链审计可追溯性" (supply chain audit traceability, using directed graphs for orphan detection). This is genuine cross-domain inspiration. Deduction: the claims are brief (one sentence) — no detail on how financial audit concepts map to doc drift (what is the materiality threshold? 17 Critical items out of 50 = 34% — is that above or below a materiality threshold?). The directed-graph orphan detection is a known graph algorithm, not a novel cross-domain application.

**Simplicity of insight (20/25):** The core insight remains appropriately simple: audit docs against code, fix what's wrong, prioritize by severity. The execution order dependency analysis (P0.4 before P0.1 because README counts depend on skill structure) is an elegant insight. The rollback plan for P0.4 (git stash -> agent regression test -> downgrade to P1) is clean. Deduction: the P0.5 decision to "remove harness type" is the simplest possible fix and demonstrates good simplicity-of-insight — but it also means the eval SKILL.md's type list shrinks, which is a feature removal disguised as documentation cleanup.

### 6. Feasibility: 82/100

**Technical feasibility (35/40):** All changes are docs and file cleanup. No Go runtime code changes. SKILL.md splitting is the most complex operation and is now described with specific risk factors. Deduction: the proposal does not address how the eval SKILL.md Mermaid flowchart will be updated after splitting. Mermaid diagrams are fragile to structural changes, and a flowchart that references non-existent file boundaries is actively misleading.

**Resource & timeline feasibility (25/30):** ~6.5h total is more realistic than the previous ~5h. SKILL.md splitting at ~3h is reasonable. Deduction: (a) P1.12 (9 subsystem docs, <=180 lines) has no time estimate. At 20 lines per subsystem, this is 180 lines of architecture documentation — realistic estimate is 1.5-2h for research + writing. If P1 is "发布前建议", the total becomes ~8.5h. (b) The breakdown does not include time for verification/testing after each P0 item. If P0.4 requires agent regression testing, add 30-60min.

**Dependency readiness (22/30):** Improved — P0.5 harness decision is now resolved (option c recommended). P0.4 is the prerequisite for all subsequent items. Deduction: the assumption that "所有信息已在审计中收集完毕" may not hold for P1.12 — writing architecture docs for 9 subsystems requires understanding each subsystem's design, which may not be fully captured in the audit data.

### 7. Scope Definition: 72/80

**In-scope items are concrete (27/30):** P0 items are highly concrete with file names, specific changes, and execution order. P1 items are adequately specified. P2 items remain one-liners but are explicitly "发布后迭代" which excuses the lower detail level. Deduction: P1.12's "总新增 <=180 行" is a good bound but "每子系统 1 段概述 + 1 行架构角色 + SKILL.md 交叉链接" is a format specification, not a content specification — what should each subsystem overview contain?

**Out-of-scope explicitly listed (22/25):** Six items listed (runtime code refactoring, new feature development, eval rubric content quality, performance optimization, i18n). Clear and relevant. Deduction: the distinction between "Go runtime code" and "agent runtime behavior" is now clearer but still not in the Out of Scope section itself. The Out of Scope says "运行时代码重构（Go CLI 代码变更）" which is precise about Go code but silent on agent behavior changes.

**Scope is bounded (23/25):** P0/P1/P2 tiering with item counts and line budgets. P1.12 bounded to <=180 lines with escalation to P2. Deduction: the P0 rollback plan ("在 README/ARCHITECTURE.md 标注已知超标，将拆分降级为 P1") effectively means P0.4 failure expands P0 scope by adding annotations. This is a bounded expansion but still an expansion.

### 8. Risk Assessment: 75/90

**Risks identified (26/30):** 7 risks identified, up from 4 in iteration 1. The cascade dependency risk ("P0 串行依赖级联失败") and scope creep risk ("P1.12 范围蔓延") are now included. The SKILL.md splitting risk is well-detailed. Deduction: (a) still missing: validation regression risk — all success criteria are manually checked; a false-positive pass could ship with errors. (b) The "后续迭代漂移" risk (H likelihood, M impact) is the only risk with High likelihood — good calibration honesty. But the "SKILL.md 拆分破坏引用链" risk is rated H impact / H likelihood in the text but M likelihood in the table — wait, it is now H/H in the table. Correction: the table shows H likelihood and H impact for this risk. This is honest calibration. However, "P0 串行依赖级联失败" is M/H — the likelihood should arguably be higher given that the H/H item is the first in the dependency chain.

**Likelihood + impact rated (24/30):** Improved calibration. SKILL.md splitting is now H/H (honest). "后续迭代漂移" is H/M (honest — it *will* happen, but impact is deferred). Deduction: "README 重写引入新不准确" is M/M — for a document with 7 Critical + 6 Major issues being rewritten, the likelihood of introducing *some* new error seems higher than Medium. The P0.4 rollback plan mitigates the cascade risk but the cascade risk itself should reflect the pre-mitigation state.

**Mitigations are actionable (25/30):** Mitigations are specific. P0.4 rollback plan is concrete: "保留拆分前副本（git stash），若 agent 回归测试失败则回滚." P1.12 mitigation: "每子系统 <=20 行，总计 <=180 行." Deduction: (a) the "后续迭代漂移" mitigation is "NFR 检查点 + CI 文档验证（deferred）" — deferred mitigation is not actionable in this proposal's scope. (b) The "死代码误删" mitigation is "仅删已确认死代码" — tautological (only delete confirmed dead code = don't delete living code). What confirms code as dead? The audit's directed graph analysis? Manual review? The freeform review challenged the init-justfile template classification specifically.

### 9. Success Criteria: 68/80

**Criteria are measurable and testable (45/55):** Significantly improved. P0 criteria include grep patterns, line count checks, and agent regression tests. P1 criteria include "零跨技能路径违规", "零孤儿 rules（引用图入度 = 0）", "漂移预防：>=1 个自动化断言." Deduction: (a) "README 事实性声明 100% 一致" — "事实性声明" is still an undefined term. Who enumerates all factual claims? Without a checklist, "100%" is untestable. (b) "ARCHITECTURE.md 已有内容 100% 一致" — the "已有内容" scope is clear (fix existing, not add new), but the P1 criterion "ARCHITECTURE.md 含 9 个子系统概述" contradicts this by adding *new* content. (c) "P0.4 agent 回归：eval 无报错、gen-test-scripts 可生成脚本" — "无报错" is ambiguous. No errors during what operation? Loading the SKILL.md? Running the full eval pipeline? Generating test scripts? (d) The "漂移预防" criterion (">=1 个自动化断言") is the weakest — one assertion is a trivial bar. ">=1" could be a single `wc -l` check that catches nothing.

**Coverage is complete (23/25):** P0 and P1 now both have criteria. P2 items have no criteria but are explicitly post-release. Deduction: (a) P0.5 (harness rubric removal) has no explicit success criterion. What verifies that the removal was done correctly? (b) P1 item 10 (guide.md Pipeline update) and P1 item 11 (forge-distribution.md alignment) have no corresponding criteria. These are P1 items that can be released without verification.

### 10. Logical Consistency: 80/90

**Solution addresses the stated problem (32/35):** The tiered remediation directly addresses the 5-dimension audit findings. P0 maps to Critical. P1 maps to Major. The drift prevention NFR addresses recurrence. Deduction: the problem statement mentions "核心文档" and "SKILL.md" but the solution extends to rules/ restructuring, dead code cleanup, and subsystem documentation authoring. The revised problem statement is better ("核心文档... skills/ 下的 SKILL.md 亦存在超标... 这些文档是用户和 agent 运行时的第一入口") but does not fully reflect the scope of plugin-internal restructuring in the solution.

**Scope <-> Solution <-> Success Criteria aligned (25/30):** P0 alignment is strong — every P0 item has a criterion. P1 alignment is improved — most P1 items have criteria. Deduction: (a) P1.10 and P1.11 have no criteria (alignment gap). (b) P0.5 (harness rubric) has no criterion. (c) The success criterion "所有 rules/ 被父 SKILL.md 引用（入度 >= 1）" does not account for parameterized references (surface-<type>.md) — the scope section mentions these separately but the criterion does not distinguish.

**Requirements <-> Solution coherent (23/25):** NFRs map cleanly to the solution's approach. Constraints reference actual conventions. The drift prevention NFR is not fully addressed by the solution (the solution is purely remediation; prevention is deferred to CI automation). The "可回归性" NFR is addressed by the P0.4 agent regression criterion. Deduction: NFR "死代码清理不能影响分发后的路径解析" — the P1.9 solution does not describe verification of path resolution impact post-cleanup.

## Phase 3: Blindspot Hunt

### CTO Failure Pattern 1: Overstated Value

Appropriately restrained. The proposal does not overclaim. "可复用审计框架" in Innovation Highlights is mildly overstated — the 5-dimension framework is a structured checklist, not a reusable framework (no tooling, no templates, no generalized methodology document). But this is a minor overstatement in a non-critical section.

### CTO Failure Pattern 2: Hidden Costs

**Partially resolved.** Timeline revised to ~6.5h. P1.12 bounded to <=180 lines. However:

1. **P1.12 time estimate missing**: 180 lines of architecture documentation for 9 subsystems is still unestimated in the timeline. At the stated budget of ~20 lines per subsystem, this is ~2h of research and writing. If P1 is "发布前建议", the total is closer to ~8.5h.
2. **Verification time not budgeted**: P0 criteria require agent regression testing, grep verification, and manual fact-checking. The timeline covers execution but not verification.
3. **P0.4 rollback cost**: If P0.4 fails and is downgraded to P1, the rollback itself takes time (restoring from git stash, annotating README/ARCHITECTURE.md, re-running affected downstream items).

### CTO Failure Pattern 3: Solution Reintroducing the Problem

**Partially addressed.** The drift prevention NFR is now present. The success criterion includes "漂移预防：>=1 个自动化断言." However:

1. The bar of ">=1 个自动化断言" is trivially low. A single `wc -l` check is one assertion but catches only one type of drift.
2. CI doc validation is "deferred" — meaning there is no concrete commitment to when or how prevention will be implemented.
3. After this fix, the next feature addition could reintroduce drift unless the CI validation is actually built.

### CTO Failure Pattern 4: Unstated Assumptions

1. **Assumption: the audit is complete and correct.** The freeform review found several issues the audit missed (e.g., `/improve-harness` ghost command, `forge feature complete --if-done` undocumented hook, `tests/e2e/` directory confusion in ARCHITECTURE.md, `forge forge task claim` typo). The 50-item count may undercount.
2. **Assumption: one person executes the entire plan.** The serialized P0 chain (P0.4 -> P0.5 -> P0.2 -> P0.3 -> P0.1) means parallelization is impossible for P0. If the single executor is interrupted, the entire P0 chain stalls.
3. **Assumption: SKILL.md splitting does not affect downstream consumers.** If any external tooling or user scripts parse SKILL.md files directly (not through the Forge CLI), splitting changes their input format.

### CTO Failure Pattern 5: Missing Rollback Plan

**Resolved for P0.4.** The rollback plan is explicit: git stash, agent regression test, fallback to P1 with annotations. However:

1. No rollback plan for P0.1 (README rewrite). If the rewritten README introduces new errors, reverting to the old README restores known-bad state.
2. No rollback plan for P0.2/P0.3. If ARCHITECTURE.md fixes or CLI reference fixes introduce new errors, what is the revert strategy?
3. No branching strategy specified. Should work be on a feature branch? Should there be a tagged checkpoint before P0.4?

---

## Bias Detection Report

**Annotated regions (pre-revised):**

6 annotated regions at lines 15-16, 33-34, 122-123, 128-129, 130-131, 141-142.

Attack points targeting annotated regions:

1. [Problem Definition] Audit methodology improved but severity threshold definitions remain implicit — "Critical/Major/Minor/Advisory" criteria not stated
2. [Solution Clarity] Target State describes structure but not content — no mockup of corrected README or post-split SKILL.md
3. [Scope Definition] P0.4 runtime impact acknowledged in description but Solution header still says "不修改 Go 运行时代码" — semantically true but potentially misleading
4. [Scope Definition] P0.5 option (c) is feature removal disguised as documentation cleanup — removing a valid eval type reduces functionality
5. [Scope Definition] P1.12 bounded to <=180 lines but "每子系统 1 段概述" format is not a content spec
6. [Risk Assessment] Cascade risk is M/H but should be higher given that the triggering item (P0.4) is H/H

**Annotated attack density: 6 / 6 regions = 1.00**

**Unannotated regions:**

~33 paragraphs without markers.

Attack points targeting unannotated regions:

1. [Problem Definition] Root cause of drift (v2->v3 migration) not analyzed — needed for prevention
2. [Problem Definition] Urgency: 5 open issues referenced but not linked — verifiability gap
3. [Solution Clarity] P2 items still one-liners — acceptable for post-release but weak for completeness
4. [Industry Benchmarking] Citations at practice level, not source level — no links to follow up
5. [Industry Benchmarking] CI automation script estimate (~2h) likely underestimated
6. [Industry Benchmarking] Automated doc generation from Go code not considered as alternative
7. [Requirements Completeness] Missing scenario: incorrect task type names — CLI behavior unknown
8. [Requirements Completeness] Drift prevention NFR is speculative — no concrete mechanism proposed
9. [Solution Creativity] "可复用审计框架" overclaimed — it is a structured checklist, not a framework
10. [Solution Creativity] Cross-domain inspiration claims are one-sentence — no depth
11. [Feasibility] P1.12 time not estimated — adds ~2h to total if executed pre-release
12. [Feasibility] Verification/testing time not budgeted
13. [Feasibility] Mermaid diagram update complexity not addressed
14. [Scope Definition] Out of Scope section still says "运行时代码重构" — should clarify "Go source code"
15. [Risk Assessment] Dead code mitigation is tautological — "仅删已确认死代码"
16. [Risk Assessment] Deferred mitigation (CI validation) is not actionable in proposal scope
17. [Success Criteria] "事实性声明" undefined — who enumerates the complete list?
18. [Success Criteria] P0.5 has no explicit success criterion
19. [Success Criteria] P1.10 and P1.11 have no corresponding criteria
20. [Success Criteria] "漂移预防：>=1 个自动化断言" is trivially low bar
21. [Success Criteria] Parameterized references (surface-<type>.md) not addressed in rules criterion
22. [Logical Consistency] Problem statement does not fully reflect plugin-internal restructuring scope
23. [Logical Consistency] P1 criterion "含 9 个子系统概述" contradicts "已有内容 100% 一致" P0 criterion scope
24. [Logical Consistency] NFR "死代码清理不能影响分发后的路径解析" — P1.9 solution does not describe verification
25. [Blindspot] P1.12 hidden cost: ~2h not in timeline
26. [Blindspot] Freeform review found additional audit misses (improve-harness, forge feature complete, tests/e2e confusion)
27. [Blindspot] No rollback plan for P0.1/P0.2/P0.3 — only P0.4 has rollback

**Unannotated attack density: 27 / 33 paragraphs = 0.82**

**Ratio (annotated/unannotated): 1.00 / 0.82 = 1.22**

The ratio of 1.22 indicates near-neutral bias between annotated and unannotated regions. The revised annotated regions received proportionate scrutiny to unannotated regions. This is a healthy ratio — no evidence of under-scrutinizing revisions.

---

## ATTACK POINTS (Tagged by Dimension)

### Annotated Region Attacks

1. **[Problem Definition]** Audit methodology improved but severity thresholds implicit — "Critical | Major | Minor | Advisory" — What distinguishes Critical from Major? A broken CLI reference that causes runtime failure is Critical, but a broken CLI reference in an orphaned rule file has zero runtime impact. The freeform review flagged this classification issue (`forge test run --tags` in orphaned files classified as P0). Without explicit severity definitions, the 17/13/15/5 split is an unverifiable assertion. — Improvement: Add a one-line severity definition for each tier (e.g., "Critical = causes runtime failure or blocks core user workflow").

2. **[Solution Clarity]** Target State is a skeleton, not a mockup — "标题→安装→快速开始→命令速查→技能列表→任务类型表→架构→贡献→License" — This tells the reader the section ordering but not the content. For the most complex P0 item (README rewrite, 7 Critical + 6 Major issues), the target structure should include: (a) what the command reference section contains (all 18 commands? subset?), (b) what the skill list format is (table? descriptions? links?), (c) what the task type table looks like (old names as aliases? migration note?). Without this, two implementers would produce different READMEs. — Improvement: Add 2-3 sentences of target content for each Target State bullet.

3. **[Scope Definition]** P0.4 runtime impact acknowledged but scope boundary still imprecise — "不修改 Go 运行时代码" — The revised scope boundary is technically precise (Go runtime code), but the P0.4 description says "此项改变 agent 运行时行为." The gap between the header claim and the item description creates a trust issue. A reader scanning the Solution section header will conclude "no runtime changes" and may miss the P0.4 admission. — Improvement: Add a sentence to the Solution header: "SKILL.md 结构重组可能改变 agent 上下文加载行为，详见 P0.4."

4. **[Scope Definition]** P0.5 option (c) is feature removal — "移除 harness 类型（纯文档删除，零运行时影响）" — Removing `harness` from the eval type list is not "纯文档删除" — it removes a feature that users could theoretically invoke (`/eval --type harness`). Even if the feature was broken (no rubric file), removing it is a feature deprecation, not a documentation fix. This should be acknowledged as such, even if the decision is correct. — Improvement: Reframe as "废弃无效 eval 类型" with explicit rationale (broken feature with no rubric).

5. **[Scope Definition]** P1.12 bounded by format but not by content — "每子系统 1 段概述 + 1 行架构角色 + SKILL.md 交叉链接" — This is a template, not a specification. What should the "概述" contain? Design rationale? Key implementation files? Integration points? Without content guidance, 9 subsystems could produce 9 inconsistent summaries. — Improvement: Add one example subsystem summary to anchor the template.

6. **[Risk Assessment]** Cascade dependency risk underrated — "P0 串行依赖级联失败 | M | H" — The first item in the chain (P0.4 SKILL.md splitting) is rated H/H. If P0.4 has High likelihood of failure, then the cascade risk should also be at least M/H, arguably H/H. The M likelihood rating implies P0.4 will probably succeed, which contradicts the H likelihood rating of P0.4 itself. — Improvement: Either lower P0.4 risk to M/H or raise cascade risk to H/H. Internal consistency requires one of these adjustments.

### Unannotated Region Attacks

7. **[Problem Definition]** Root cause of drift not analyzed — The proposal identifies *what* drifted but not *why*. Was it the v2->v3 migration? Feature additions without doc updates? Multi-contributor coordination gaps? Without root cause analysis, the drift prevention NFR ("漂移预防") is uninformed — you cannot prevent what you do not understand. — Improvement: Add 1-2 sentences analyzing drift causes and feed this into the prevention mechanism design.

8. **[Industry Benchmarking]** CI automation cost underestimated — "需编写验证脚本（~2h），无 CI 基础设施" — Writing a script that validates doc claims against codebase state requires: (a) parsing Go source to extract command definitions, (b) parsing SKILL.md to extract CLI references, (c) comparing the two, (d) handling edge cases (parameterized references, cross-skill paths). This is 4-8h minimum, not 2h. Additionally, "无 CI 基础设施" is stated as fact without evidence — does the project have GitHub Actions? Any CI at all? — Improvement: Either remove the time estimate or provide a realistic breakdown.

9. **[Industry Benchmarking]** Automated doc generation not considered — The proposal considers manual fix, CI validation, and partial scope fixes. Missing: generating CLI reference documentation from Go cobra command definitions (an industry-standard approach used by Kubernetes, Hugo, and many Go CLI tools). This would eliminate an entire category of drift (CLI reference accuracy) permanently. — Improvement: Add as an alternative or as a P2 follow-up with a note on feasibility.

10. **[Requirements Completeness]** Drift prevention mechanism unspecified — "建立漂移预防（如 README 计数与 ls skills/ | wc -l 断言）" — The "如" (e.g.) makes this speculative, not committed. The success criterion says ">=1 个自动化断言" but no mechanism is proposed. Is it a shell script? A Makefile target? A CI step? Without specifying the mechanism, the NFR is an aspiration, not a requirement. — Improvement: Replace "如" with a committed mechanism: "创建 scripts/validate-docs.sh，包含 README 计数断言和 CLI 引用 grep 检查."

11. **[Feasibility]** P1.12 time not estimated — "总新增 <=180 行" is a scope bound, not a time estimate. Writing accurate architecture summaries for 9 subsystems requires reading each subsystem's SKILL.md, understanding its design, and distilling it into 20 lines. At 10-15 min per subsystem, that is 1.5-2.25h. If P1 is executed pre-release, the total is ~8.5h, not ~6.5h. — Improvement: Add P1.12 estimate to timeline breakdown, or explicitly note that P1 timeline is separate from the ~6.5h P0 estimate.

12. **[Success Criteria]** "事实性声明" scope undefined — "README 事实性声明 100% 一致" — The proposal does not enumerate what constitutes a "事实性声明." Version numbers, counts, paths, and command names are listed as examples, but is the task type naming convention a "factual claim"? Is the architecture section a "factual claim"? Without an explicit checklist, "100%" is untestable — the implementer decides the scope of verification. — Improvement: Either enumerate all factual claims (a checklist of ~15-20 items based on the audit) or define the term precisely: "factual claims = version numbers, file counts, CLI command names, directory paths, agent/hook counts."

13. **[Success Criteria]** P0.5 and P1.10/P1.11 have no criteria — The scope section lists these items as deliverables but the success criteria section does not mention them. P0.5 (harness rubric removal) is supposedly P0 (release-blocking) but has no release-blocking criterion. P1.10 (guide.md update) and P1.11 (forge-distribution.md alignment) are P1 but have no P1 criterion. — Improvement: Add criteria for each: P0.5 = "harness 类型已从 eval SKILL.md 类型列表中移除"; P1.10/P1.11 = "guide.md 和 forge-distribution.md Pipeline 描述与实际流程一致."

14. **[Logical Consistency]** P0 and P1 criteria scope contradiction — P0 criterion: "ARCHITECTURE.md 已有内容 100% 一致." P1 criterion: "ARCHITECTURE.md 含 9 个子系统概述." The P0 criterion scopes verification to "已有内容" (existing content only), but P1.12 adds 180 lines of new content. The P1 criterion for ARCHITECTURE.md should distinguish between "existing content accuracy" and "new content completeness" to avoid conflating the two verification tasks. — Improvement: Split ARCHITECTURE.md P1 criterion into: (a) "P1.12 新增 9 个子系统概述符合格式规范（每子系统 <=20 行，含架构角色 + SKILL.md 链接）" and keep the existing "已有内容 100%" as P0-only.

15. **[Blindspot]** Freeform review findings not incorporated — The independent freeform review found several issues the proposal's audit missed: `/improve-harness` ghost command, `forge feature complete --if-done` undocumented in ARCHITECTURE.md, `tests/e2e/` directory confusion, `forge forge task claim` typo. The proposal's 50-item count may be an undercount. The success criterion "100% 一致" is unachievable if the audit itself is incomplete. — Improvement: Either (a) acknowledge the freeform review findings as supplement to the audit, or (b) explain why they were excluded from scope. Incorporating them as additional P1/P2 items would be the strongest response.

16. **[Blindspot]** No rollback for P0.1/P0.2/P0.3 — P0.4 has an explicit rollback plan, but P0.1 (README rewrite), P0.2 (ARCHITECTURE.md fix), and P0.3 (CLI reference fix) have none. If the README rewrite introduces new errors, reverting to the old README restores known-bad state. The proposal should define a go/no-go checkpoint after each P0 item, not just after P0.4. — Improvement: Add a general P0 rollback principle: "each P0 item is committed on a feature branch; any item that fails its criterion is reverted and documented before proceeding."

---

## Scoring Summary

| Dimension | Score | Max |
|-----------|-------|-----|
| Problem Definition | 96 | 110 |
| Solution Clarity | 100 | 120 |
| Industry Benchmarking | 85 | 120 |
| Requirements Completeness | 92 | 110 |
| Solution Creativity | 55 | 100 |
| Feasibility | 82 | 100 |
| Scope Definition | 72 | 80 |
| Risk Assessment | 75 | 90 |
| Success Criteria | 68 | 80 |
| Logical Consistency | 80 | 90 |
| **Total** | **805** | **1000** |

### Iteration-over-Iteration Improvement

| Dimension | Iter 1 | Iter 2 | Delta |
|-----------|--------|--------|-------|
| Problem Definition | 88 | 96 | +8 |
| Solution Clarity | 95 | 100 | +5 |
| Industry Benchmarking | 55 | 85 | +30 |
| Requirements Completeness | 85 | 92 | +7 |
| Solution Creativity | 40 | 55 | +15 |
| Feasibility | 75 | 82 | +7 |
| Scope Definition | 68 | 72 | +4 |
| Risk Assessment | 65 | 75 | +10 |
| Success Criteria | 55 | 68 | +13 |
| Logical Consistency | 72 | 80 | +8 |
| **Total** | **698** | **805** | **+107** |

### Remaining Gaps to Target (900)

The proposal needs +95 points to reach the 900 target. The most efficient improvements by remaining gap size:

1. **Industry Benchmarking** (35 remaining): Add source-level citations with links; include automated doc generation as alternative; justify why manual fix is preferred over automated approaches.
2. **Solution Creativity** (45 remaining): This is the hardest gap to close — the proposal is fundamentally a tech-debt cleanup. Deepen the cross-domain claims with actual methodology mapping; consider proposing a novel doc-impl contract testing pattern.
3. **Success Criteria** (12 remaining): Define "事实性声明" explicitly; add criteria for P0.5 and P1.10/P1.11; strengthen drift prevention from ">=1 assertion" to specific checks.
4. **Problem Definition** (14 remaining): Add root cause analysis; link the 5 open issues.
5. **Feasibility** (18 remaining): Add P1.12 time estimate; budget verification time; address Mermaid update complexity.
