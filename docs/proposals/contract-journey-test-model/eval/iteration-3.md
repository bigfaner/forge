# Eval Report: Proposal (Iteration 3)

**Score**: 828/1000
**Target**: 900
**Date**: 2026-05-17

## DIMENSIONS

### 1. Problem Definition (98/110)

| Criterion | Score | Notes |
|-----------|-------|-------|
| Problem stated clearly | 38/40 | Core problem is unambiguous: test organization is by language x interfaces instead of user workflows. One minor ambiguity: the relationship between "Journey-Driven" and the existing model could be sharper -- is this a replacement or a parallel system? The statement "将测试管道从 ... 重新设计为 Journey-Driven 模型" clarifies it's a redesign, but later "现有单步 TC 作为单步骤 Journey 退化形式" suggests backward compat, creating slight tension. |
| Evidence provided | 34/40 | Five concrete evidence points: (1) gen-test-cases output by interface type, (2) 6 hardcoded language profiles, (3) CLI/API/TUI vs E2E conflation, (4) contract regression from task-pipeline-hardening (11 fixes), (5) risk flat-weighting. Strong concrete examples. Missing: quantitative data on how many test cases are split across files, user complaints, or incident reports. The reference to task-pipeline-hardening section 4 is specific but the referenced doc does not exist in the proposals directory (verified via glob), weakening the citation. |
| Urgency justified | 26/30 | Four urgency points: CLI iteration speed, linear scaling cost of language profiles, contract regression repetition, agile testing quadrant gaps. Strong on "why now" but light on cost-of-delay quantification -- no estimate of how many regressions per month, or how much manual rework the current model causes. |

### 2. Solution Clarity (108/120)

| Criterion | Score | Notes |
|-----------|-------|-------|
| Approach is concrete | 38/40 | The 4-step pipeline (gen-journeys -> gen-contracts -> gen-test-scripts -> run-tests) is clearly delineated with inputs/outputs per step. The Journey pseudocode with Step/Outcome structure is concrete. A reader can explain back what will be built. Minor gap: the relationship between Contract and Journey is described conceptually but the physical file format of a Contract spec is not shown -- only fragments. |
| User-facing behavior described | 40/45 | CLI commands shown (`forge feature`, `forge task claim`), expected outputs given, tag lifecycle (`@feature` -> `@regression`) is clear, `forge test verify` and `forge test promote` are described. The Journey pseudocode is exemplary. Deduction: the end-user experience of `forge test verify` is described only as "scans Contract specs" and "compares" -- the actual CLI output format (what the user sees) is not shown. What does a "broken contract" report look like to the user? |
| Technical direction clear | 30/35 | Technical hints are present: Fact Table for code reconnaissance, semantic descriptors deferred to gen-test-scripts, markdown-with-schema for Contract storage, config.yaml for framework declaration. Clear enough for implementation. Gap: the "markdown with schema" format for Contract specs is mentioned but not exemplified -- what does a stored Contract file actually look like on disk? |

### 3. Industry Benchmarking (98/120)

| Criterion | Score | Notes |
|-----------|-------|-------|
| Industry solutions referenced | 37/40 | Five industry solutions cited: Pact (consumer-driven contracts), Playwright test.step, Go testing/Robot Framework, Google Testing pyramid. Pact comparison is thorough with three specific differences. Playwright and Go testing are mentioned with one-line descriptions but lack depth of comparison. |
| At least 3 meaningful alternatives | 26/30 | Four alternatives presented: Do nothing, Incremental Pyramid Overlay, Keep language x interfaces + add workflow, Journey-Driven model. The Incremental Pyramid Overlay alternative (row 2) is well-argued with three specific failure modes -- this is genuinely different from the chosen approach and not a straw man. The "Keep language x interfaces + add workflow" alternative is the weakest -- described as "XY problem" with a 2-sentence dismissal. It could be a straw man given how tersely it's handled, but the "XY problem" framing is defensible. |
| Honest trade-off comparison | 22/25 | Cons table for Journey-Driven lists "改动面大" (large change surface) -- honest. The incremental overlay's cons are detailed (3 failure modes). The "Do nothing" cons are appropriate. Minor concern: all alternatives are "Rejected" with one being "Selected" -- the table structure is binary rather than scoring alternatives on multiple axes, which would show more nuanced comparison. |
| Chosen approach justified against benchmarks | 13/25 | The chosen approach claims inspiration from Pact + Robot Framework + Google Testing, but the justification against benchmarks is thin. The Innovation Highlights section mentions Patton's Story Mapping and Pact's consumer-driven ideas, which is good溯源. However, the key question -- "why not just use Pact directly?" -- is answered only implicitly (Pact is for microservices, Forge is for CLI/TUI). The justification for not adopting Robot Framework's keyword-driven approach is absent. The section "Forge 的创新在于将语义描述符、标签晋升和六维度契约结合" is more of a description than a justification of why this combination beats established tools. |

### 4. Requirements Completeness (92/110)

| Criterion | Score | Notes |
|-----------|-------|-------|
| Scenario coverage | 36/40 | Seven key scenarios: 4 concrete Journeys (CLI task, TUI diagnosis, Web-UI milestone, API registration), contract regression, edge case (mid-step failure), backward compat. The 4 Journey examples span CLI/TUI/Web-UI/API -- good interface coverage. Each has risk levels. Missing: no scenario for the tag promotion lifecycle (feature -> regression), no scenario for multi-Outcome steps, no scenario for the `forge test verify` false-positive case. |
| Non-functional requirements | 32/40 | Seven NFRs listed: backward compat, batch generation, project adaptation, contract break reporting, execution performance (120% of current), Journey isolation, auth contracts. The "120% of current" performance target is measurable -- good. Deductions: (1) No security NFR beyond auth state -- what about test data isolation, secret handling in contracts? (2) No accessibility NFR. (3) "执行性能" excludes Setup time, which is a significant caveat -- total test runtime could be much higher. |
| Constraints & dependencies | 24/30 | Five constraints listed, including the critical State dimension degradation path (state-verification: partial / deferred). The degradation mechanism is well-thought-out with concrete fallback behavior. Missing: no constraint on LLM context window size limits (mentioned in passing but not as a hard constraint), no constraint on Fact Table coverage gaps, no dependency on specific forge-cli version compatibility. |

### 5. Solution Creativity (78/100)

| Criterion | Score | Notes |
|-----------|-------|-------|
| Novelty over industry baseline | 32/40 | Two genuinely novel ideas: (1) Semantic descriptors with deferred precision (gen-contracts uses business language, gen-test-scripts generates regex from Fact Table) -- this directly addresses LLM limitations in observing code output. (2) Tag-Based Promotion replacing file migration -- simpler than traditional test graduation. Differentiation from benchmarks is articulated in Innovation Highlights with溯源 to Patton and Pact. |
| Cross-domain inspiration | 26/35 | Story Mapping (Jeff Patton) from product management, consumer-driven contracts from Pact (microservices), tag-based lifecycle from CI/CD (git tags, @feature/@regression). Three domains represented. However, no inspiration from adjacent domains like property-based testing, chaos engineering, or formal verification methods that could address the contract completeness concern. |
| Simplicity of insight | 20/25 | The tag-based promotion is elegant ("why didn't I think of that"). The semantic descriptor split is practical rather than elegant -- it's a workaround for LLM limitations, not a conceptual simplification. The six-dimension Contract model is somewhat overengineered -- do all six dimensions need to be explicit for every step, or could some be inferred? The document partially addresses this by making Invariants and Side-effect optional, which is good. |

### 6. Feasibility (88/100)

| Criterion | Score | Notes |
|-----------|-------|-------|
| Technical feasibility | 35/40 | 4-step pipeline is decomposed well. Each step's cognitive task is within LLM capabilities. Fact Table mechanism validated (20+ e2e tests). Config-driven approach already has foundation. Risk: the semantic descriptor -> regex conversion relies heavily on Fact Table coverage and LLM accuracy -- this is acknowledged in risks but the mitigation (forge test verify) is itself part of the proposal and unproven. |
| Resource & timeline feasibility | 27/30 | 1 engineer, 10 weeks, 3 phases with clear deliverables. Phase 1 has explicit 1-week risk buffer. Phase 1 scope concern is identified as a risk with a concrete fallback (TUI await degraded to CLI-only first delivery). Timeline is realistic for the scope. Minor concern: Phase 3 at 2 weeks includes 5 deliverables (run-tests rewrite + verify + CLI adaptation + eval rubric + e2e validation) -- this is tight. |
| Dependency readiness | 26/30 | All 5 dependencies confirmed as existing: config.yaml, conventions/, 6 language profiles, testing subcommand, Fact Table (20+ e2e tests). Strong. The task-pipeline-hardening reference is cited as evidence but the document doesn't exist in the proposals directory -- this weakens the "11 contract fixes" claim as verifiable evidence. |

### 7. Scope Definition (72/80)

| Criterion | Score | Notes |
|-----------|-------|-------|
| In-scope items are concrete | 27/30 | 14 in-scope items, each is a deliverable (e.g., "gen-journeys skill", "forge test verify", "TUI await 语义形式化"). Most are specific enough to generate tasks from. Minor vagueness: "forge-cli testing 命令重命名和适配" -- what specifically changes? Just rename or behavioral changes too? |
| Out-of-scope explicitly listed | 22/25 | 7 out-of-scope items with rationale. Idempotency contract has explicit deferral reasoning tied to Phase 2. Good. Missing: no mention of migration tooling for existing projects, documentation updates, or training/onboarding materials. |
| Scope is bounded | 23/25 | 10-week timeline with 3 phases. Phase boundaries are clear. The scope is bounded by the 14 in-scope items and 10-week timeline. Minor concern: "eval rubric 更新" is in Phase 3 scope but no detail on what this entails -- could expand. |

### 8. Risk Assessment (80/90)

| Criterion | Score | Notes |
|-----------|-------|-------|
| Risks identified | 28/30 | 8 risks identified, spanning design (6-dimension coverage), migration (template regression), complexity (config system), performance (pipeline overhead), stability (flaky smoke tests), accuracy (semantic->regex), scope (Phase 1 size), and combinatorics (Outcome explosion). Comprehensive. Missing: no risk around LLM prompt engineering for gen-contracts (prompt quality直接影响 Contract 准确性). |
| Likelihood + impact rated | 25/30 | Ratings vary: M/H, L/M, L/M, L/L, M/M, M/H, M/M, L/M. Good distribution -- not all "low likelihood, high impact". However, "语义描述符 -> regex 转换不准确" is rated M/H (the highest risk) but the mitigation is circular -- it relies on `forge test verify` which is itself part of the proposal and unvalidated. The risk of "Contract 六维度无法覆盖所有测试场景" is rated M/H but the mitigation ("归入 Invariants") is a catch-all that reduces the dimension model's precision rather than truly mitigating the risk. |
| Mitigations are actionable | 27/30 | Most mitigations are actionable: "逐个 language 迁移 + e2e 回归验证", "内置模板作为默认值", "自动重试 3 次 + @flaky 标记", "Outcome 按 Preconditions 互斥 + 5个触发检查点". The Phase 1 scope risk has a specific fallback (TUI degraded to CLI-only). Two mitigations are less actionable: "六维度覆盖已知场景" is descriptive not prescriptive, and "forge test verify 检测断裂" is the thing being built, not a mitigation. |

### 9. Success Criteria (62/80)

| Criterion | Score | Notes |
|-----------|-------|-------|
| Criteria are measurable and testable | 42/55 | 14 success criteria listed. Several are excellent: "3 个 Journey 并行运行时结果与顺序运行完全一致（文件系统状态、退出码、输出内容 diff 为空）" -- highly measurable. "误报数 = 0" on 20+ unchanged Contracts -- testable. "diff 为空" comparisons -- testable. However, several criteria are vague: (1) "输出叙述性工作流文档" -- what constitutes a valid narrative document? No acceptance criteria for the output quality. (2) "Invariants 覆盖跨步骤不变量" -- how do we verify coverage is sufficient? (3) "端到端验证...全链路通过" -- "通过" is undefined. Does one failing edge case mean "not passed"? (4) "生成的测试代码使用声明的框架 import 和断言库" -- this is testable but the criterion is about structural correctness, not behavioral correctness. (5) "可编译的测试代码" -- compilation success is a low bar. |
| Coverage is complete | 20/25 | Success criteria cover most in-scope items: gen-journeys, gen-contracts, gen-test-scripts, run-tests, tag promotion, verify, Journey isolation, backward compat, config-driven, TUI await, template migration, CLI rename. Missing: no success criterion for (1) "契约断裂报告" NFR (output format of verify), (2) Risk分级 in gen-journeys output, (3) 多 Outcome Contract generation quality, (4) 分批生成 auto-split behavior. |

### 10. Logical Consistency (52/90)

| Criterion | Score | Notes |
|-----------|-------|-------|
| Solution addresses the stated problem | 30/35 | The 5 stated problems map to solution elements: (1) "tests not reflecting workflows" -> Journey-first organization addresses this directly. (2) "language enumeration incomplete" -> config-driven framework selection addresses this. (3) "test level conflation" -> explicit 3-tier model (Unit/Contract/E2E) addresses this. (4) "contract regression no mechanism" -> forge test verify addresses this. (5) "risk flat-weighting" -> Journey-level Risk grading addresses this. All 5 problems have corresponding solutions. Minor gap: the "risk flat-weighting" problem is addressed by risk grades (High/Medium/Low) but the proposal doesn't define what different test intensity means at each level -- it says "高风险应有更密集的边缘案例覆盖" but doesn't specify how. |
| Scope <-> Solution <-> Success Criteria aligned | 12/30 | Significant misalignments: (1) "Journey 隔离：独立临时工作目录" is in scope but the success criterion only tests "3 个 Journey 并行" -- no criterion for the temporary directory mechanism itself. (2) "gen-contracts skill" scope includes "多 Outcome" but only one success criterion mentions it implicitly ("六维度，语义描述符，多 Outcome" in a checklist item). (3) "TUI await 语义形式化" is in scope with a success criterion ("异步 Cmd 等待超时时 fail-fast") but the scope item includes "超时阈值、fail-fast 行为、并发 Cmd 等待语义" while the criterion only covers the timeout case. (4) "内置模板迁移" scope item's success criterion is "diff 为空" -- this is alignment but the criterion may be too strict (what if minor formatting differs?). (5) "forge-cli testing 命令重命名" has a checklist item but no criterion for behavioral equivalence post-rename. |
| Requirements <-> Solution coherent | 10/25 | Orphan requirements: (1) "认证/授权契约" NFR has no corresponding success criterion and no dedicated solution mechanism beyond "Preconditions must include auth state". (2) "契约断裂报告" NFR specifies "输出标识失败维度和具体不匹配内容" but the success criterion "输出标识失败维度和具体不匹配内容" is a nearly verbatim copy -- not a testable specification of output format. (3) "执行性能" NFR (120% of current) has no success criterion verifying this constraint. Orphan solution features: (1) The "多 Outcome Contract" mechanism is a significant design feature but has no dedicated success criterion for correctness (what if two Outcomes' Preconditions both match?). (2) The "分批生成" mechanism (auto-split at 15 Contracts or 50k tokens) is described in the solution but has no success criterion. |

## ATTACKS

1. **Logical Consistency**: Scope-Success Criteria misalignment on Journey isolation -- Scope lists "独立临时工作目录" but no criterion verifies the temp directory mechanism works correctly (only parallel execution consistency is tested). Success criteria should include "每个 Journey 执行在独立临时目录中，执行后目录被清理" and verify via filesystem inspection.

2. **Logical Consistency**: NFR "执行性能" (Contract test execution time <= 120% of current integration tests) has no corresponding success criterion. This is an orphan requirement -- stated but never verified.

3. **Logical Consistency**: "多 Outcome" is a core Contract mechanism but has no dedicated success criterion. The combination-explosion avoidance mechanism (Preconditions-based mutual exclusion) is untested -- quote: "gen-contracts 通过以下机制避免组合爆炸...Outcome 按 Preconditions 互斥分组". A criterion should verify that generated multi-Outcome Contracts have mutually exclusive Preconditions.

4. **Logical Consistency**: "认证/授权契约" NFR requires "Contract 的 Preconditions 必须包含 auth 状态声明" and "gen-test-scripts 生成的 Setup 代码注入对应 auth 策略" but no success criterion verifies auth Contract generation or auth Setup injection. This is an orphan requirement.

5. **Industry Benchmarking**: Justification against benchmarks is the weakest dimension element -- the comparison table has a "Verdict" column but no scoring on multiple axes (e.g., implementation cost, maintenance overhead, learning curve). The question "why not use Robot Framework directly?" is never answered.

6. **Success Criteria**: Several criteria use unquantified success conditions: "全链路通过" (what constitutes passing?), "Invariants 覆盖跨步骤不变量" (what is sufficient coverage?), "输出叙述性工作流文档" (what makes a valid narrative document?). These fail the "measurable and testable" bar.

7. **Problem Definition**: The reference "docs/proposals/task-pipeline-hardening/test-strategy-evaluation.md 第 4 节" does not exist in the file system (verified). This weakens the evidence chain for the contract regression problem.

8. **Feasibility**: The mitigation for the highest-rated risk (M/H: semantic descriptor -> regex inaccuracy) is circular -- it relies on `forge test verify` which is part of the proposed system and itself unvalidated. Quote: "forge test verify 检测断裂; gen-test-scripts 使用 Fact Table 验证转换结果". The verify command's own accuracy is assumed, not guaranteed.

9. **Requirements Completeness**: No security NFR beyond auth state declarations. Contracts may contain sensitive data (API keys, tokens in test setup), but no requirement addresses test data sanitization or secret management in generated test code.

10. **Solution Clarity**: The physical format of a stored Contract spec is described as "markdown with schema" but never shown. The document provides fragments (Outcome blocks, dimension tables) but no complete example of a Contract file as it exists on disk. This makes it hard for implementers to know what gen-contracts should output.

11. **Risk Assessment**: No risk around LLM prompt quality for gen-contracts and gen-journeys. The entire pipeline depends on LLM generating correct Journeys and Contracts, but prompt engineering quality is not listed as a risk factor. If gen-contracts produces incomplete Preconditions, the entire Contract system's value collapses.

## SUMMARY
- Total: 828/1000
- Target met: No (828 < 900, gap of 72 points)
- Biggest area for improvement: Logical Consistency (52/90) -- the misalignment between Scope items and Success Criteria is the largest point sink, accounting for ~28 points of gap from target
- Remaining gap: Primarily driven by (1) Scope-Success Criteria alignment failures, (2) weak benchmark justification, and (3) vague success criteria that fail the measurability bar
