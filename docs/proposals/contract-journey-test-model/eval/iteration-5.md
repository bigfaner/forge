# Eval Report: Proposal (Iteration 5)

**Score**: 926/1000
**Target**: 900
**Date**: 2026-05-17

## DIMENSIONS

### 1. Problem Definition (105/110)

| Criterion | Score | Notes |
|-----------|-------|-------|
| Problem stated clearly | 39/40 | Core problem is unambiguous: test organization by language x interfaces vs user workflows. The 5 sub-problems are each clearly stated with concrete descriptions. Minor residual tension: "Journey-Driven" is described as a "redesign" but backward compat via "单步骤 Journey 退化形式" creates slight ambiguity about whether this is a replacement or a migration path. Resolved sufficiently by the backward compat NFR and success criterion. |
| Evidence provided | 37/40 | Five concrete evidence points: (1) gen-test-cases output by interface type (cli-test-cases.md, api-test-cases.md), (2) 6 hardcoded language profiles, (3) CLI/API/TUI vs E2E conflation, (4) contract regression -- 11 fixes in task-pipeline-hardening, (5) risk flat-weighting. Citation path exists. Remaining gaps: no quantitative data on how many test cases are split across files, no user complaints or survey data. The "11 contract fixes" claim is specific but not independently verifiable from the document alone. |
| Urgency justified | 29/30 | Four urgency drivers with quantification. "每个新 language 约 2 人天" addresses the iteration-4 attack point on cost quantification. "11 个契约修复点，累计消耗约 3 人天" gives concrete sunk cost. "扩展成本线性增长" is now quantified. Minor gap: no estimate of how many regressions occur per month or projected growth rate of language profiles. |

### 2. Solution Clarity (117/120)

| Criterion | Score | Notes |
|-----------|-------|-------|
| Approach is concrete | 39/40 | The 4-step pipeline (gen-journeys -> gen-contracts -> gen-test-scripts -> run-tests) is clearly delineated with inputs/outputs per step. The Journey pseudocode with Step/Outcome structure is concrete. The multi-Outcome example with Preconditions mutual exclusion is well-illustrated. Minor gap: the batch splitting logic ("15 Contracts or 50k tokens") is described narratively but no pseudocode shows how the auto-split works in practice. |
| User-facing behavior described | 44/45 | CLI commands shown, expected outputs given, tag lifecycle clear. The `forge test verify` output example is now present: "BROKEN (2): ... Output dimension: expected 'claimed task <task_id>' -> actual stdout 'Task <task_id> claimed' ... Semantic match: partial (word order changed)". This directly addresses the iteration-4 attack point. Deduction: the verify example shows 2 broken contracts but the "OK" section says "unchanged contracts omitted for brevity" -- an acceptable shorthand but a minor omission of the full output format. |
| Technical direction clear | 34/35 | Technical hints are present: Fact Table for code reconnaissance, semantic descriptors deferred to gen-test-scripts, markdown-with-schema for Contract storage, config.yaml for framework declaration. The Contract file example provides a concrete format reference. The verify report example clarifies the output format. Gap: the "markdown with schema" format is shown via example but no formal schema definition (JSON Schema, YAML schema) is provided -- implementers must infer the schema from examples. |

### 3. Industry Benchmarking (111/120)

| Criterion | Score | Notes |
|-----------|-------|-------|
| Industry solutions referenced | 37/40 | Five industry solutions cited: Pact (consumer-driven contracts with thorough comparison), Playwright test.step, Go testing/Robot Framework, Google Testing pyramid. The Pact comparison is thorough with three specific differences. The Robot Framework justification has three concrete reasons for not adopting it. Playwright and Go testing still have only brief one-line descriptions. |
| At least 3 meaningful alternatives | 28/30 | Four alternatives: Do nothing, Incremental Pyramid Overlay, Keep language x interfaces + add workflow, Journey-Driven model. The Incremental Pyramid Overlay has specific failure modes. The "Keep language x interfaces + add workflow" alternative is dismissed as "XY problem" which borders on straw-man but is defensible since the proposal argues the fundamental model is wrong. |
| Honest trade-off comparison | 23/25 | Multi-axis comparison table with quantitative backing: "2 人天/language" for Do Nothing maintenance, "约 3 周" for Incremental, "约 1 周" for Keep, "10 周，后续新增 language 零改动" for Journey-Driven. This addresses the iteration-4 attack point on single-word ratings. Deduction: "Workflow Alignment" and "Contract Verification" columns still use qualitative assertions -- "工作流被拆散到 5 个接口文件" vs "一个用户工作流 = 一个文件" are descriptive, not measured. A quantitative metric (e.g., "avg files per workflow trace: 5 vs 1") would strengthen the comparison. |
| Chosen approach justified against benchmarks | 23/25 | The Robot Framework justification is explicit with three specific reasons. Pact's consumer-driven insight is clearly distinguished from Forge's needs. The Innovation Highlights section traces inspiration (Patton's Story Mapping, Pact). The comparison table shows concrete cost figures. Deduction: "Workflow Alignment" and "Contract Verification" for the Journey-Driven approach are self-assessed as "一个用户工作流 = 一个文件" and "六维度 Contract + forge test verify" without external validation of these claims -- the comparison is directional but not evidenced. |

### 4. Requirements Completeness (100/110)

| Criterion | Score | Notes |
|-----------|-------|-------|
| Scenario coverage | 38/40 | Nine key scenarios: 4 concrete Journeys (CLI task, TUI diagnosis, Web-UI milestone, API registration), contract regression, edge case (mid-step failure), backward compat, tag promotion lifecycle ("gen-test-scripts 生成测试注入 @feature 标签 -> CI 运行 --tags feature -> forge test promote 通过后标签更新为 @regression -> CI 运行 --tags regression 作为回归门禁"), batch-split auto-trigger ("单个 Journey 的 Contract 数超过 15 或估计 token 数超过 50k 时...自动拆分为多次调用"). The tag promotion lifecycle and batch-split auto-trigger scenarios directly address iteration-4 attack points. Minor gap: no scenario for `forge test verify` false-positive resolution workflow. |
| Non-functional requirements | 35/40 | Eight NFRs listed: backward compat, batch generation, project adaptation, contract break reporting, execution performance (120% of current), Journey isolation, auth contracts, test data security. The "测试数据安全" NFR is specific and actionable. Deductions: (1) No accessibility NFR for Web-UI Journey scenarios. (2) "执行性能" excludes Setup time with caveat "不含 Setup 时间" -- total runtime could be much higher than 120%. (3) No NFR for maintainability or extensibility of the Contract model itself (e.g., adding a 7th dimension in the future). |
| Constraints & dependencies | 27/30 | Five constraints listed, including the State dimension degradation path with specific behavior for each degradation level. The Fact Table staleness is addressed as a formal risk with mitigation. Deduction: LLM context window size limits (the 50k token threshold mentioned in passing) are not listed as a formal constraint. No dependency versioning or compatibility constraint specified. |

### 5. Solution Creativity (80/100)

| Criterion | Score | Notes |
|-----------|-------|-------|
| Novelty over industry baseline | 33/40 | Two genuinely novel ideas: (1) Semantic descriptors with deferred precision -- gen-contracts uses business language, gen-test-scripts generates regex from Fact Table. (2) Tag-Based Promotion replacing file migration. The six-dimension Contract model is the most novel structural contribution. Differentiation from benchmarks is articulated. Gap: the semantic descriptor mechanism is a workaround for LLM limitations rather than a conceptual breakthrough -- it's practical rather than paradigm-shifting. |
| Cross-domain inspiration | 26/35 | Story Mapping (Jeff Patton) from product management, consumer-driven contracts from Pact (microservices), tag-based lifecycle from CI/CD. Three domains represented. No inspiration from adjacent domains like property-based testing (QuickCheck/Hypothesis), chaos engineering, or formal verification methods that could address the contract completeness concern. The cross-domain borrowing is solid but not surprising. |
| Simplicity of insight | 21/25 | The tag-based promotion is elegant ("why didn't I think of that"). The semantic descriptor split is practical. The six-dimension Contract model is somewhat overengineered -- the proposal itself acknowledges that Invariants and Side-effect are optional, and the catch-all "归入 Invariants" fallback for unclassifiable scenarios reduces the model's discriminative power. |

### 6. Feasibility (92/100)

| Criterion | Score | Notes |
|-----------|-------|-------|
| Technical feasibility | 37/40 | 4-step pipeline is well-decomposed. Each step's cognitive task is within LLM capabilities. Fact Table mechanism validated (20+ e2e tests). Config-driven approach has foundation. The bootstrap strategy for verify is concrete. The Fact Table staleness mitigation (auto-refresh on each verify run) resolves the circularity concern. Remaining risk: semantic descriptor -> regex conversion still relies on LLM accuracy, but the verify mechanism provides a feedback loop. |
| Resource & timeline feasibility | 27/30 | 1 engineer, 10 weeks, 3 phases with clear deliverables. Phase 1 has explicit 1-week risk buffer. Phase 3 at 2 weeks includes 5 deliverables (run-tests rewrite + verify + CLI adaptation + eval rubric + e2e validation) -- this is tight and no risk buffer is allocated. No contingency for the single-engineer dependency (illness, context-switching). |
| Dependency readiness | 28/30 | All 5 dependencies confirmed as existing: config.yaml, conventions/, 6 language profiles, testing subcommand, Fact Table (20+ e2e tests). Minor gap: no dependency versioning or compatibility constraint specified -- "forge-cli 的 testing 子命令架构已稳定" is a claim without version pinning. |

### 7. Scope Definition (75/80)

| Criterion | Score | Notes |
|-----------|-------|-------|
| In-scope items are concrete | 28/30 | 14 in-scope items, each is a deliverable. Most are specific enough to generate tasks from. Minor vagueness: "forge-cli testing 命令重命名和适配" -- what specific adaptations beyond renaming? The success criterion clarifies "所有子命令行为不变（仅命令前缀变更）" but the scope item itself is slightly ambiguous. |
| Out-of-scope explicitly listed | 22/25 | 7 out-of-scope items with rationale. Idempotency contract has explicit deferral reasoning tied to Phase 2 with a specific expansion path. Good. Missing: no mention of documentation or onboarding materials for the new model. |
| Scope is bounded | 25/25 | 10-week timeline with 3 phases. Phase boundaries are clear. 14 in-scope items are enumerable. The scope is bounded. |

### 8. Risk Assessment (86/90)

| Criterion | Score | Notes |
|-----------|-------|-------|
| Risks identified | 29/30 | 10 risks identified (up from 9 in iteration 4). The new risk "Fact Table 快照与代码库演进不同步" directly addresses the iteration-4 attack point. Risks span design, migration, complexity, performance, stability, accuracy, scope, combinatorics, LLM quality, and data staleness. Minor gap: no risk around single-engineer knowledge concentration (the timeline depends on 1 engineer). |
| Likelihood + impact rated | 28/30 | Ratings vary: M/H, L/M, L/M, L/L, M/M, M/H, M/M, L/M, M/H, M/M. Good distribution -- not all low-likelihood-high-impact. The Fact Table staleness risk is honestly rated M/M. Deduction: "内置模板迁移引入回归" is rated L/M but the success criterion sets a very high bar ("零配置时 gen-test-scripts 输出与现有 profile 输出 diff 为空") that could make the likelihood higher than estimated. |
| Mitigations are actionable | 29/30 | Most mitigations are actionable with specific steps. The new Fact Table staleness mitigation is concrete: "forge test verify 每次运行时自动重新采集 Fact Table（基于当前代码库），不依赖历史快照；bootstrap 快照仅用于 verify 自身的首次正确性验证". The LLM quality mitigation is actionable. One mitigation remains less precise: "六维度无法覆盖所有测试场景" relies on the Invariants catch-all which reduces model precision rather than addressing the root cause. |

### 9. Success Criteria (78/80)

| Criterion | Score | Notes |
|-----------|-------|-------|
| Criteria are measurable and testable | 53/55 | 22 success criteria (up from 19 in iteration 4). Three new criteria address all three iteration-4 orphan requirements: (1) Auth Contract criterion: "Web-UI Journey 的 Contract Preconditions 包含 auth 状态声明（authenticated/unauthenticated/expired token 之一），gen-test-scripts 生成的 Setup 代码包含对应的 auth 策略代码（如 cookie 注入或 token header 设置）" -- measurable and specific. (2) Batch-split criterion: "当单个 Journey 的 Contract 数 > 15 时，gen-contracts 自动拆分为多批调用，合并后的 Contract 文档与单批生成的 Contract 在六维度内容上完全一致（diff 为空）" -- objectively verifiable. (3) Journey smoke test criterion: "烟测试端到端运行 Journey 的 happy path，每步输出与 Contract 中 'success' Outcome 声明的 Output/State 完全匹配" -- behavioral correctness specified. Remaining deductions: (1) "通过六维度完整性校验（必选维度非空，语义描述符不含 regex 语法）" -- the "不含 regex 语法" check is not formally defined (what about escaped characters? Unicode patterns? character classes in semantic descriptions?); (2) "输出叙述性工作流文档" with "叙述性" quality remains somewhat subjective despite structural format requirements. |
| Coverage is complete | 25/25 | Success criteria now cover all in-scope items: gen-journeys (2 criteria), gen-contracts (3 criteria), gen-test-scripts (2 criteria), run-tests (1 criterion), tag promotion (1 criterion), verify (1 criterion), Journey isolation (1 criterion), backward compat (1 criterion), config-driven (1 criterion), TUI await (1 criterion), template migration (1 criterion), CLI rename (1 criterion), multi-Outcome (1 criterion), execution performance (1 criterion), Risk grading (1 criterion), auth contracts (1 criterion), batch-split (1 criterion), smoke test correctness (1 criterion). No remaining coverage gaps. |

### 10. Logical Consistency (82/90)

| Criterion | Score | Notes |
|-----------|-------|-------|
| Solution addresses the stated problem | 33/35 | All 5 stated problems map to solution elements: (1) tests not reflecting workflows -> Journey-first organization; (2) language enumeration incomplete -> config-driven framework selection; (3) test level conflation -> 3-tier model; (4) contract regression no mechanism -> forge test verify; (5) risk flat-weighting -> Journey-level Risk grading. All 5 problems have corresponding solutions with measurable success criteria. Minor gap: the "risk flat-weighting" problem is addressed with a success criterion for edge case density ("高风险 Journey 的边缘场景数量 >= happy path 步骤数量") but no definition of what different test intensity means at each risk level beyond edge case count -- no mention of assertion density, coverage requirements, or test depth variation. |
| Scope <-> Solution <-> Success Criteria aligned | 27/30 | The three iteration-4 misalignments are now resolved: (1) Auth contracts NFR has a dedicated success criterion verifying Contract Preconditions contain auth state declarations and gen-test-scripts generates corresponding auth strategy code. (2) Batch-split NFR has a dedicated criterion verifying auto-split at > 15 Contracts with diff-based equivalence check. (3) Journey smoke test has a behavioral correctness criterion. Remaining minor misalignments: (1) TUI await scope item includes "超时阈值、fail-fast 行为、并发 Cmd 等待语义" but the success criterion only covers timeout ("fail-fast 并报告超时 Cmd 名称") and batch await ("tea.Batch(cmd1,cmd2) 等待全部完成后再进入下一步") -- no criterion for "超时阈值" being configurable or having a sensible default value. (2) "契约断裂报告" NFR specifies "输出标识失败维度和具体不匹配内容" and the success criterion echoes this, but there is no criterion for the report format being machine-parseable for CI integration. |
| Requirements <-> Solution coherent | 22/25 | The iteration-4 orphan requirements are now resolved: auth contracts NFR maps to auth Contract success criterion; batch-split NFR maps to batch-split criterion; smoke test maps to behavioral correctness criterion. The "测试数据安全" NFR maps to Contract placeholder behavior. Remaining minor gaps: (1) "执行性能" NFR excludes Setup time ("不含 Setup 时间") which creates a coherence gap between the 120% performance claim and actual total test runtime experienced by users -- if Setup dominates, the NFR may be misleadingly narrow; (2) the Risk grading mechanism ("高风险 Journey 的边缘场景数量 >= happy path 步骤数量") is a quantitative criterion but the solution description does not specify how gen-journeys determines risk level -- is it from PRD severity tags, heuristics, or manual annotation? The success criterion assumes the output exists but the mechanism is underspecified. |

## ATTACKS

1. **Logical Consistency**: TUI await scope item claims "超时阈值、fail-fast 行为、并发 Cmd 等待语义" but the success criterion only tests "fail-fast 并报告超时 Cmd 名称" and "tea.Batch(cmd1,cmd2) 等待全部完成后再进入下一步" -- the configurable timeout threshold ("超时阈值") has no success criterion verifying it is configurable or has a sensible default. Quote from scope: "TUI await 语义形式化：定义为 Contract 维度的 await 规范（超时阈值、fail-fast 行为、并发 Cmd 等待语义）" vs criterion: "TUI await 语义：异步 Cmd 等待超时时 fail-fast 并报告超时 Cmd 名称；tea.Batch(cmd1,cmd2) 等待全部完成后再进入下一步". Add a criterion for timeout configurability or default value verification.

2. **Solution Clarity**: The markdown-with-schema Contract format is shown via example but no formal schema definition exists. Implementers must infer the schema from the `step-2-task-claim.md` example. Quote: the Contract file example shows "## Outcome 'success'" with nested "- Preconditions:", "- Input:", "- Output:", "- State:", "- Side-effect:" but no JSON Schema, YAML schema, or formal grammar defines required fields, allowed values, or validation rules. A formal schema would eliminate ambiguity in gen-contracts output validation.

3. **Industry Benchmarking**: "Workflow Alignment" and "Contract Verification" columns in comparison table use qualitative assertions rather than quantitative metrics. Quote: "工作流被拆散到 5 个接口文件" vs "一个用户工作流 = 一个文件" -- these are descriptive but not measured. A quantitative comparison (e.g., "avg files per workflow trace: 5 (current) vs 1 (Journey-Driven)") would strengthen the trade-off analysis and make the self-assessment more credible.

4. **Requirements Completeness**: "执行性能" NFR excludes Setup time with "不含 Setup 时间" which creates a gap between the stated performance guarantee and user-perceived performance. Quote: "Contract 测试的执行时间不超过当前集成测试的 120%（不含 Setup 时间）" -- if Setup time dominates total runtime, the 120% claim may be misleading. The success criterion should clarify what fraction of total runtime Setup typically represents, or the NFR should specify total runtime including Setup.

5. **Success Criteria**: "通过六维度完整性校验（必选维度非空，语义描述符不含 regex 语法）" -- the "不含 regex 语法" check is not formally defined. Quote: "语义描述符不含 regex 语法". What counts as regex syntax? Escaped characters? Character classes? Quantifiers in natural language? This criterion needs a precise definition (e.g., "no characters from the set `[ ] { } ( ) * + ? | \ ^ $ .` used in a pattern-matching context") to be truly testable.

6. **Solution Creativity**: No inspiration drawn from adjacent testing domains like property-based testing, chaos engineering, or formal verification that could address the contract completeness concern. Quote from Innovation Highlights: the section traces inspiration only from Story Mapping, Pact, and CI/CD tags. Property-based testing (e.g., QuickCheck/Hypothesis) could generate adversarial Contract inputs automatically, addressing the "六维度无法覆盖所有测试场景" risk more systematically than the Invariants catch-all fallback.

## SUMMARY
- Total: 926/1000
- Target met: Yes (926 >= 900)
- Biggest improvements from iteration 4: Logical Consistency improved from 64 to 82 (+18) via aligned NFR-criterion pairs for auth, batch-split, and smoke test; Success Criteria improved from 76 to 78 (+2) via 3 new criteria closing orphan requirements; Risk Assessment improved from 82 to 86 (+4) via Fact Table staleness risk with actionable mitigation; Industry Benchmarking improved from 104 to 111 (+7) via quantitative backing in comparison table
- Key strength: All 8 iteration-4 attack points have been addressed with specific, measurable content
- Remaining gaps are minor: TUI timeout configurability criterion, formal Contract schema, quantitative workflow alignment metric, Setup time exclusion clarity, regex syntax definition in completeness check
