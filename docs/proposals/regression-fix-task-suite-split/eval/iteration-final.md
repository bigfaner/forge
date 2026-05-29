# CTO Adversarial Evaluation — iteration-final

**Proposal**: Regression Fix Task 按 Test File 拆分（语言无关）
**Date**: 2026-05-29
**Evaluator**: CTO Adversary

---

## Phase 1: Reasoning Audit

### Argument Chain Trace

1. **Problem → Solution**: The problem states that a single fix task covering 20+ failures across multiple files causes agent freeze. The proposed solution splits by test file. This is valid — scoping each task to one file reduces the cognitive load per task. However, the proposal identifies two distinct problems in Phase 0 and Phase 1 but conflates them into a single proposal without demonstrating that Phase 1 is still needed after Phase 0 lands.

2. **Solution → Evidence**: Evidence is a single incident (`docs/lessons/gotcha-fix-task-broad-scope.md`). This is one data point. The proposal does not present frequency data (how often does this occur?).

3. **Evidence → Success Criteria**: SCs test the mechanism (splitting, fallback, cap behavior) but do NOT test whether the agent actually succeeds at fixing the failures. The original problem was "agent 卡住" — none of the SCs measure agent task completion rate or time.

4. **Self-contradiction check**: The proposal states "两阶段无条件交付" and claims "Phase 0 和 Phase 1 解决不同层面的问题，不互斥". Yet if Phase 0 resolves the information visibility issue (agent can see all failures), the remaining problem is scope size — and there is no evidence presented that scope size alone (with full information visible) still causes agent freeze. The "unconditional" assertion is an assumption, not a conclusion.

### SC Consistency Deep-Dive

Cluster by affected area:

- **Cap/count mechanism**: SC "确认 countFixTasks 按 title prefix 匹配正常工作", SC "maxFixTasksPerStep 保留用于非 regression 步骤", SC "addRegressionFixTasks 受 regression 专用软上限（10 个）约束"
  - These are consistent with each other.
- **Splitting logic**: SC "N 个测试文件各有直接 ---FAIL: 条目时创建 N 个独立 fix task", SC "每个 fix task 的 description 包含该测试文件相关的输出行"
  - These are consistent.
- **Fallback**: SC "无法识别测试文件时 fallback 创建按目录分组的 fix task"
  - Consistent with scope definition.

**Bidirectional satisfiability issue**: In Scope states "新建 extractFileLineMap 函数" with specific algorithm steps (1-6). The SCs verify the *output* of this function (task per file, description content) but do NOT verify the algorithm correctness itself (e.g., context window overlap dedup). An SC for "overlapping context windows produce no duplicate lines" is missing.

---

## Phase 2: Rubric Scoring

### 1. Problem Definition (110 pts)

**Problem stated clearly (35/40)**: The problem is concrete — single fix task covering 20+ failures causes agent freeze. Quote: "quality-gate 创建的 fix-1 包含 4 个 test suite 共 20+ 失败，agent 执行后长时间无响应被用户手动中断". Deduction: the proposal title says "语言无关" but the problem only occurred with Go. The title creates a mismatch with the actual scope.

**Evidence provided (30/40)**: Single incident documented in `gotcha-fix-task-broad-scope.md`. No frequency data, no reproduction rate, no data on how often this pattern occurs. Quote: "实际发生：quality-gate 创建的 fix-1 包含 4 个 test suite 共 20+ 失败". One event is not a pattern.

**Urgency justified (20/30)**: Quote: "高。每次 regression 测试出现多文件失败都会触发此问题，agent 卡死浪费时间和 token。" This is stated but not quantified. How many times has this happened? What is the token cost? What is the time cost? Without data, urgency is an assertion.

**Subtotal: 85/110**

### 2. Solution Clarity (120 pts)

**Approach is concrete (35/40)**: Two phases with clear boundaries. Phase 0 modifies description generation. Phase 1 creates new `addRegressionFixTasks` with `extractFileLineMap`. Quote: "新建 extractFileLineMap 函数，签名：func extractFileLineMap(output string) map[string][]string". Deduction: The algorithm for line-to-file association has 6 numbered steps in Scope, but the description section itself is less precise — it says "含上下文窗口" without defining the window size there (it's ±2 lines, buried in Scope bullet 3). The core approach description is scattered.

**User-facing behavior described (35/45)**: Agent receives per-file fix tasks with filtered output. But "user-facing" here means agent-facing, and the behavior for the human operator (who monitors task creation) is not described. What does the task list look like? How does the operator know which task covers which file? The title format is not specified for split tasks — will it still be "fix test: just test failure in quality gate" for all of them? This creates ambiguity when 10 tasks have the same title.

**Technical direction clear (30/35)**: Function signatures, algorithms, and integration points with existing code are well specified. Quote: "runTestRegression 失败时调用新函数替代 addFixTask". Deduction: The `createFixTask` helper extraction is mentioned but its interface is not specified.

**Subtotal: 100/120**

### 3. Industry Benchmarking (120 pts)

**Industry solutions referenced (30/40)**: CI test grouping (GitHub Actions, JUnit XML), Sentry error fingerprinting. These are relevant analogies. Quote: "CI 系统普遍按 test module/suite 分组报告失败". Deduction: References are superficial — no specific product documentation or implementation details cited beyond names.

**At least 3 meaningful alternatives (25/30)**: Six alternatives in comparison table (do nothing, improved description, file grouping, LLM analysis, JUnit XML, Sentry fingerprinting, directory grouping). The LLM analysis alternative is described well enough to reject credibly.

**Honest trade-off comparison (20/25)**: The comparison table includes cons for each approach. However, the "Selected" approach's cons are understated: "同一根因 bug 创建多任务" is listed but its impact is not explored. If one bug in a utility function causes failures across 8 test files, the proposal creates 8 independent fix tasks that may all attempt to modify the same source file — this concurrent editing risk is acknowledged in risks but minimized in the trade-off table.

**Chosen approach justified (18/25)**: The justification is essentially "simplest thing that works" — which is valid for an MVP. But the proposal does not explain why file-level grouping is superior to suite-level grouping (the lesson document actually suggests `FAIL forge-tests/<suite>` extraction which maps more naturally to Go's test output). The pivot from suite-based (in the lesson) to file-based (in this proposal) is not justified.

**Subtotal: 93/120**

### 4. Requirements Completeness (110 pts)

**Scenario coverage (30/40)**: Five key scenarios listed. Missing scenarios:
- Mixed test files and non-test files in the same output
- Very large output (>10000 lines) where parsing performance matters (mentioned in NFR but no scenario)
- Nested test failures (sub-tests)
- Parallel test output interleaving (mentioned in risks as out-of-scope, but no scenario for fallback behavior)
- What happens when `extractFileLineMap` returns empty map but output contains failures

**Non-functional requirements (30/40)**: Performance bounds stated with specific numbers. Quote: "对 10000 行 output、F=15 的典型场景，解析时间 < 100ms". However, the 750MB worst-case memory estimate for F=15 is alarming and no mitigation is proposed beyond stating the number. For a CLI tool, 750MB is not acceptable.

**Constraints & dependencies (22/30)**: Dependencies on naming convention stated. MVP limited to Go. Quote: "依赖测试文件命名约定（MVP 仅 *_test.go）". Deduction: No discussion of what happens when `go test` output format changes across Go versions. The `sourceFileRe` regex behavior across Go versions is an unstated dependency.

**Subtotal: 82/110**

### 5. Solution Creativity (100 pts)

**Novelty over industry baseline (15/40)**: The proposal explicitly states "无创新". Quote: "无创新。原提案使用 Go 专属的 FAIL <package> 解析，改进为...结合测试文件命名约定识别。与 CI 系统中按 failing module 分组报告的常规做法一致。" This is a straightforward application of standard grouping.

**Cross-domain inspiration (10/35)**: Sentry fingerprinting is mentioned as inspiration but not actually used. The proposal uses simple file-path matching, not semantic fingerprinting.

**Simplicity of insight (20/25)**: The Phase 0 insight (improve description instead of restructuring) is genuinely simple and valuable. The recognition that the problem has two layers (information visibility + scope size) shows good decomposition.

**Subtotal: 45/100**

### 6. Feasibility (100 pts)

**Technical feasibility (35/40)**: The proposal demonstrates deep familiarity with the existing codebase, verified by my code review. The `countFixTasks` analysis is accurate. The `extractFileLineMap` design is implementable. Quote: "新建 extractFileLineMap 函数...输入输出类型简单（string → map[string][]string）". Deduction: The context window overlap deduplication algorithm is not trivial — merging overlapping ±2-line windows across multiple match points requires careful implementation.

**Resource & timeline feasibility (25/30)**: Phase 0 half a day, Phase 1 1-2 days. Reasonable for the described scope. However, the `createFixTask` helper extraction from `addSingleFixTask` is a refactor that could expand scope — the proposal acknowledges this risk but the timeline does not account for the refactor independently.

**Dependency readiness (25/30)**: No external dependencies. Quote: "无外部依赖。extractFileLineMap 以现有 sourceFileRe 正则为基线". The `sourceFileRe` regex is a dependency whose behavior must be verified for the new use case.

**Subtotal: 85/100**

### 7. Scope Definition (80 pts)

**In-scope items concrete (25/30)**: Very detailed scope items with specific function names and algorithmic steps. Quote: "新建 extractFileLineMap 函数，签名：func extractFileLineMap(output string) map[string][]string". Deduction: The scope numbering has a bug — there are two items numbered "5" in the algorithm steps (line 124 and line 126 in the proposal), suggesting insufficient proofreading.

**Out-of-scope listed (22/25)**: Eight items explicitly listed as out of scope. Good. Quote: "基线对比过滤预存在失败（lesson 第二层改进）：本提案仅实现第一层". Each item references the source of the deferred work.

**Scope bounded (20/25)**: The MVP is bounded to Go. But the proposal simultaneously describes multi-language extension paths in detail (Python/JS-TS/Java/Ruby mentioned 3+ times), creating ambiguity about whether the current scope is truly bounded or if partial multi-language work is expected.

**Subtotal: 67/80**

### 8. Risk Assessment (90 pts)

**Risks identified (25/30)**: Nine risks identified. Good coverage. Includes both technical risks (parsing accuracy) and process risks (claim priority, cross-feature pollution).

**Likelihood + impact rated (22/30)**: Each risk has Likelihood and Impact ratings. However, the ratings use H/M/L without defined scales. Quote: "测试文件命名不规范导致识别失败 | M | L". What does "M" mean? 50%? The risk "拆分后 fix task 数量增加，加剧 claim priority 问题" is rated H likelihood — yet this is listed as out of scope with no mitigation within the proposal. A High-likelihood risk with no in-scope mitigation is a gap.

**Mitigations actionable (25/30)**: Most mitigations are concrete. The fallback mechanism is well-designed: "fallback 到按目录分组，产生与改动前完全一致的 monolithic task——零回归、零改进". However, several mitigations defer to "短期" and "长期" without committing to either. Quote for claim priority risk: "短期：dispatcher 层面对同一 feature 下 fix task 添加优先 claim 逻辑" — this is actionable but out of scope, making the mitigation theoretical for this proposal.

**Subtotal: 72/90**

### 9. Success Criteria (80 pts)

**Measurable and testable (22/30)**: Most SCs are testable. Quote: "N（N ≤ 10）个测试文件各有直接 ---FAIL: 条目时，创建 N 个独立 fix task". Deduction: Several SCs use vague language:
- "addSingleFixTask 的 description 包含所有 --- FAIL 行条目" — what does "所有" mean when output has 10,000 lines?
- "现有 compile/fmt/lint/unit-test 步骤的 fix task 创建不受影响" — how is "不受影响" tested? Regression test suite? Manual verification?
- SC for "确认 countFixTasks 按 title prefix 匹配正常工作" is already checked — it's not a criterion for delivery.

**Coverage complete (18/25)**: Missing SCs:
- No SC for the 750MB memory concern raised in NFR
- No SC for Phase 0 specifically improving agent completion rate (the original problem)
- No SC for task title uniqueness (10 tasks with same title)
- No SC for the `createFixTask` helper extraction mentioned in scope
- No SC for context window deduplication correctness

**SC internal consistency (20/25)**: SCs are internally consistent with each other. However, SC "超过 10 个测试文件有失败时，第 11 至 N 个文件按 groupFilesByDir 归并为 1 个综合 overflow task" uses `groupFilesByDir` — but this function returns nil when files are in a single directory. If all overflow files are in the same directory, no overflow task is created, contradicting the SC.

**Subtotal: 60/80**

### 10. Logical Consistency (90 pts)

**Solution addresses stated problem (25/35)**: The solution addresses the mechanism (splitting tasks) but the connection to the original problem (agent freeze) is weak. The proposal does not explain why smaller-scope tasks will prevent agent freeze — it assumes this without evidence. The lesson document identifies three layers of root cause: (1) scope too large, (2) incomplete description, (3) pre-existing failures mixed in. The proposal addresses layers 1 and 2 but does layer 0 (description improvement) actually solve the freeze? If yes, Phase 1 may be unnecessary. If no, what evidence supports that file-level scoping alone prevents freeze?

**Scope ↔ Solution ↔ SC aligned (22/30)**: Generally aligned. Gap: Scope includes "共享 helper 须有独立单元测试覆盖" but no SC verifies this. Scope includes "输出结构化日志警告" and there is a matching SC — good alignment.

**Requirements ↔ Solution coherent (20/25)**: The solution design matches the requirements. However, the NFR "内存 < 750MB" for worst case is alarming and the solution does not address this. A 750MB memory spike in a CLI hook is a production risk that the solution should mitigate.

**Subtotal: 67/90**

---

## Phase 3: Cross-Dimension Coherence Check

1. **Feasibility vs NFR**: The 750MB worst-case memory estimate in NFR is inconsistent with the "完全可行" assessment in Feasibility. A CLI tool consuming 750MB is not "完全可行" without mitigation.

2. **Scope vs SC**: The scope item for `createFixTask` helper has no corresponding SC. If the refactor is in scope, its correctness should be verifiable.

3. **Problem vs SC**: The problem is "agent 卡住" but no SC measures agent behavior improvement. All SCs test the mechanism, not the outcome.

4. **Risks vs Out-of-Scope**: Two risks rated "H" likelihood (claim priority, cross-feature pollution) are explicitly out of scope. The proposal may make these worse (more tasks = more claim confusion) without addressing them.

---

## Phase 4: Blindspot Hunt

1. **[blindspot] Task title uniqueness**: The proposal does not address task title generation for split tasks. Quote: "每个 task 只包含该测试文件相关的输出行". Currently `addSingleFixTask` generates: `title := fmt.Sprintf("fix %s: %s failure in quality gate", step, testScript)`. With 10 split tasks, all would have the same title "fix test: just test failure in quality gate". This makes the task index un-navigable and breaks `countFixTasks`'s title prefix matching — the cap mechanism matches by prefix, and 10 tasks with the same prefix would all count toward the cap. The proposal does not specify how titles are differentiated for split tasks.

2. **[blindspot] groupFilesByDir returns nil for single-directory overflow**: The overflow strategy uses `groupFilesByDir` but this function returns nil when files are in one directory. Quote: "第 11 至 N 个文件按 groupFilesByDir 归并为一个综合 task". If overflow files share a directory, the function returns nil, producing zero overflow tasks. The SC "总 task 数 ≤ 11" would be violated (only 10 tasks created, overflow failures silently dropped).

3. **[blindspot] Phase 0 和 Phase 1 的无条件交付假设**: Quote: "Phase 0 独立交付且具有自身价值...Phase 1 无论 Phase 0 效果如何都会实施". If Phase 0 fully resolves the agent freeze (by giving the agent complete failure information), Phase 1's splitting becomes a nice-to-have optimization, not a necessity. The "unconditional" delivery of Phase 1 wastes engineering effort if Phase 0 is sufficient. The proposal should define a decision criterion for Phase 1 rather than declaring it unconditional.

4. **[blindspot] 750MB 内存开销无缓解措施**: Quote: "最坏情况（每个匹配行归入所有主文件，步 5 跨文件复制）扩展倍率 ≤ F，对 F=15 的场景总内存 < 750MB". A CLI hook consuming 750MB of memory is a serious concern, yet the proposal provides no mitigation. The NFR section identifies the risk but the solution design does not include streaming or chunking to bound memory.

5. **[blindspot] `extractFileLineMap` 的空输出场景未定义**: Quote: "仅为拥有直接 ---FAIL: 条目的文件生成映射条目". If the output contains failures but none match the Go `---FAIL:` pattern (e.g., panic output, build errors in test files), `extractFileLineMap` returns an empty map, and the code falls through to the existing `addFixTask` — but the proposal does not make this explicit. The fallback SC only covers "无法识别测试文件" (isTestFile returns zero), not "no FAIL lines parsed".

6. **[blindspot] Scope 算法步骤编号重复**: The In Scope section for the description algorithm has two items numbered "5": "一行匹配多个'主文件'时，该行及其上下文归入所有匹配的主文件" and "示例 description 内容". This suggests the document was not carefully reviewed before submission.

---

## Summary

| Dimension | Score | Max |
|-----------|-------|-----|
| Problem Definition | 85 | 110 |
| Solution Clarity | 100 | 120 |
| Industry Benchmarking | 93 | 120 |
| Requirements Completeness | 82 | 110 |
| Solution Creativity | 45 | 100 |
| Feasibility | 85 | 100 |
| Scope Definition | 67 | 80 |
| Risk Assessment | 72 | 90 |
| Success Criteria | 60 | 80 |
| Logical Consistency | 67 | 90 |
| **Total** | **756** | **1000** |
