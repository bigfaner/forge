---
created: "2026-05-28"
iteration: 2
role: adversary
reviewer: CTO-adversary
previous_report: iteration-1.md
---

# Adversarial Evaluation Report — Iteration 2

## Bias Detection Report

- Annotated regions (`<!-- pre-revised -->`): 8 attack points / 10 annotated paragraphs = density 0.80
- Unannotated regions: 8 attack points / 12 unannotated paragraphs = density 0.67
- Ratio (annotated/unannotated): 1.19

Conclusion: Marginal bias detected (annotated density 1.19x unannotated). Within acceptable range (< 1.5). Annotated regions received slightly more scrutiny, which is expected since revisions often introduce new issues. No corrective action needed.

## Iteration-1 Issue Tracking

| # | Iteration-1 Attack | Status | Assessment |
|---|-------------------|--------|------------|
| 1 | Scope algorithm vs Risk 4 multi-matching contradiction | **Resolved** | Proposal line 109 now says "归入所有匹配的测试文件" and Risk 4 (line 132) also says "归入所有匹配文件（宁可多包含）". Contradiction eliminated. |
| 2 | Cap bypass mechanism architecturally unspecified | **Resolved** | Line 111: "通过直接调用底层 task 创建 API（绕过 `addSingleFixTask` 的 cap 检查）来绕过 cap." Architectural decision stated. |
| 3 | File-lock mechanism claim unsubstantiated | **Partially Resolved** | Line 128: Changed to "多个 agent 可能并发编辑同一生产代码文件，需人工介入解决冲突." The old unsubstantiated claim is removed. However, the new mitigation is "accept this trade-off" which is honest but raises a different concern (see blindspot-4). |
| 4 | `extractFileLineMap` interface unspecified | **Resolved** | Line 103: `func extractFileLineMap(output string) map[string][]int`. Signature provided. |
| 5 | Comparison table "复用现有代码" stale | **Resolved** | Line 70: "新代码量可控" replaces the old "复用现有代码" claim. |
| 6 | Timeline underestimate | **Partially Resolved** | Revised to "6-8 小时." Closer to reality but still optimistic (see Phase 2). |
| 7 | No upper bound on regression fix tasks | **Resolved** | Line 130: "regression 路径最多创建 10 个 fix task，超出部分合并到按目录分组的 fallback." Soft cap of 10 added. |
| 8 | `sourceFileRe` regex not designed for test framework output | **Resolved** | Line 55: "使用与 `sourceFileRe` 不同的提取逻辑以适配测试框架输出格式（如 pytest 的 `FAILED file::class::method`、Go 的 `--- FAIL:` 块）." Acknowledges the problem and states new extraction logic. |
| 9 | Proposal addresses only layer 1 of lesson | **Resolved** | Line 118: "本提案仅实现第一层（按测试文件拆分），基线过滤作为后续迭代独立实现。lesson 文档指出两层可独立生效." Explicitly acknowledged as partial solution with deferral rationale. |
| 10 | Agent's view inside fix task unspecified | **Partially Resolved** | Line 105: "description 格式为：测试文件路径 + 筛选后的相关输出行（含上下文窗口）." Description format partially specified but still vague (see attack 8). |
| 11 | No rollback plan | **Resolved** | Line 135: Rollback plan explicitly documented. |
| 12 | Straw-man alternative | **Not Addressed** | "Go 专属 suite 解析（原提案 v1）" still present as a comparison row. Not a genuinely different alternative. |
| 13 | "Partially Overridden" qualifier | **Resolved** | Line 94: Now says "Overridden for regression path" with explicit acknowledgment that "历史 loop 恰好发生在 regression 上下文中." |
| 14 | Feasibility section cites irrelevant test coverage | **Not Addressed** | Line 86: "`extractSourceFiles` 已稳定运行" is still used as feasibility evidence, but the proposal builds `extractFileLineMap` (a new function). |
| 15 | False positive from stack trace helper files | **Resolved** | Line 43: "栈 trace 引用辅助测试文件" scenario added with explicit accept-and-document stance. |

**Summary**: 15 attacks from iteration-1. 10 resolved, 3 partially resolved, 2 not addressed. The proposal made substantive improvements. However, resolving old issues introduced new ones (see below).

## Phase 1: Reasoning Audit

### Argument Chain Trace

**Problem -> Solution**: The revised proposal now explicitly acknowledges it implements only layer 1 (suite splitting) while layer 2 (baseline filtering) is deferred. The deferral rationale ("lesson 文档指出两层可独立生效，第一层已能解决当前 agent 卡死问题") is valid. The link from problem to solution is now cleaner — the solution addresses the stated symptom (scope too broad) even if it does not address the deeper cause (pre-existing failures).

**Solution -> Implementation**: The proposal now specifies the cap bypass mechanism ("直接调用底层 task 创建 API"). This resolves the architectural gap. However, it introduces a new concern: `addSingleFixTask` (line 708-802) does far more than cap checking — it performs surface inference (`inferSurface`), derives task type (`fixTypeFromStep`), populates template defaults, builds `AddTaskOpts`, calls `task.AddTask`, calls `task.CreateTaskMarkdown`, and calls `feature.EnsureForgeState`. The proposal's `addRegressionFixTasks` must replicate ALL of this logic (minus the cap check) or extract shared code into a lower-level function. Neither path is trivial, and neither is scoped.

**Soft cap -> SC consistency**: The soft cap of 10 (line 130) was added to address the "no upper bound" attack. But no Success Criterion verifies this soft cap. SC-1 tests "4 个独立 fix task" and SC-3 tests "不受 cap 限制." Neither tests "最多创建 10 个." The soft cap is in the Risk mitigation but not in the Success Criteria, making it unverifiable.

### Self-Contradiction Check

1. **Soft cap contradicts "绕过 cap" framing**: The Assumptions Challenged table (line 94) says "`addRegressionFixTasks` 绕过 cap" and the Success Criteria SC-3 (line 142) says "`addRegressionFixTasks` 不受 cap 限制." But Risk 3 (line 130) introduces "regression 路径最多创建 10 个 fix task" — this IS a cap, just a different one. The proposal says "bypasses cap" but then introduces a new cap. The framing is misleading — it should say "replaces the general cap with a regression-specific soft cap of 10."

2. **`extractFileLineMap` return type vs algorithm mismatch**: The function signature is `func extractFileLineMap(output string) map[string][]int` (file path to line numbers). But the algorithm in Scope (lines 106-109) describes extracting matched lines WITH context windows (前后各 2 行). Line numbers alone (`[]int`) are insufficient to represent context windows with deduplication. The function should return `map[string][]string` (file path to actual output lines) or a richer structure. As specified, the signature cannot support the described algorithm.

3. **Feasibility evidence still references old mechanism**: Line 86 says "`extractSourceFiles` 已稳定运行" as Dependency Readiness evidence. But the proposal builds `extractFileLineMap` as a NEW function with different logic ("与 `sourceFileRe` 不同的提取逻辑"). The stability of `extractSourceFiles` is irrelevant to the new function's feasibility.

### SC Consistency Deep-Dive

Cluster SC entries by affected area:

**Cluster A — `addRegressionFixTasks` function**: SC-1 (4 fix tasks), SC-2 (per-file output lines)
- SC-1 + SC-2: Satisfiable. The algorithm is now unambiguous (multi-matching lines go to all files).
- **Issue**: SC-2 says "包含该文件路径的行及上下文" but `extractFileLineMap` returns `map[string][]int` (line numbers, not lines). Who constructs the actual description content? This is the gap between the function signature and the algorithm.

**Cluster B — Cap policy**: SC-3 (bypass cap), SC-5 (other steps unaffected), Risk-3 (soft cap 10)
- SC-3 + SC-5: Satisfiable if `addRegressionFixTasks` calls the low-level API directly (as specified on line 111).
- **Issue**: SC-3 says "不受 cap 限制" but Risk-3 introduces a soft cap of 10. **Ambiguous — requires author clarification**: Is the soft cap of 10 a hard requirement (should be in SC) or an implementation detail (guideline)? If the former, SC-3 is misleading. If the latter, the soft cap has no enforcement mechanism.

**Cluster C — Language coverage**: SC-4 (fallback behavior), SC-6 (5 languages)
- SC-4 + SC-6: Satisfiable. Fallback for unrecognized languages + 5 explicit naming conventions.

**Cross-cluster**: SC-1 (4 fix tasks) + Risk-3 (soft cap 10) + SC-3 (bypass cap)
- If 15 test files fail, the soft cap of 10 means some files are merged into directory-based fallback tasks. This means not every test file gets its own task — SC-1 would produce fewer than 15 tasks, but the example in SC-1 (4 files) stays under the cap. No logical contradiction, but the soft cap's interaction with the splitting algorithm is untested by any SC.

## Phase 2: Rubric Scoring with Verification Stance

### 1. Problem Definition: 78/110

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Problem stated clearly | 32/40 | Core problem is clear: single fix task with broad scope causes agent stall. Improved from iteration-1 with the addition of the "栈 trace 引用辅助测试文件" scenario showing nuance. Deduction: "agent 执行时卡住" remains imprecise — still does not distinguish between timeout, infinite loop, or output-quality-induced stall. The lesson document clarifies "长时间无响应被用户手动中断" but the proposal does not carry this precision forward. |
| Evidence provided | 28/40 | One concrete incident with lesson document. The proposal now explicitly acknowledges it is a partial solution (layer 1 only), which is honest. Deduction: still single data point, no frequency data, no severity classification beyond "高," no data on how many user sessions are affected. |
| Urgency justified | 18/30 | "每次 regression 测试出现多文件失败都会触发此问题" — but no data on frequency. The deferral of layer 2 (baseline filtering) partially undermines urgency: if baseline filtering is the higher-value improvement (as the lesson states "即使第一层不拆分，基线过滤也能将 scope 自然收窄"), then this proposal's urgency is reduced — layer 2 alone might be sufficient. |

### 2. Solution Clarity: 78/120

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Approach is concrete | 32/40 | Function names, signatures, and algorithm specified. Cap bypass mechanism now specified (line 111). Deduction: `extractFileLineMap` returns `map[string][]int` but the algorithm requires extracting actual lines with context windows. The return type cannot support the described algorithm — this is a specification-level bug. |
| User-facing behavior described | 34/45 | "创建 4 个独立 fix task" is observable. Description format now partially specified ("测试文件路径 + 筛选后的相关输出行（含上下文窗口）"). Deduction: "筛选后的相关输出行（含上下文窗口）" is still imprecise. What does the agent actually see? Raw output lines? Annotated with match markers? Deduplicated? The format of "上下文窗口" in the description is undefined. |
| Technical direction clear | 12/35 | The cap bypass mechanism is now specified (call low-level API directly). But this creates a significant unscoped implementation concern: `addSingleFixTask` (lines 708-802) contains 95 lines of logic beyond the cap check (surface inference, template defaults, task creation, markdown creation, state update). The proposal must either (a) replicate this logic (code duplication) or (b) extract shared code into a helper (refactoring). Neither option is scoped or estimated. |

### 3. Industry Benchmarking: 55/120

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Industry solutions referenced | 18/40 | Unchanged from iteration-1. Two parenthetical mentions: "GitHub Actions test grouping、JUnit XML testsuite 元素." No analysis of how these systems work, what patterns can be borrowed, or how they handle multi-language output. |
| At least 3 meaningful alternatives | 16/30 | The LLM-based grouping alternative is now better described with honest cons ("非确定性、增加 token 开销"). But "Go 专属 suite 解析（原提案 v1）" remains a straw man — it is the proposal's own previous iteration. Only 3 genuinely distinct alternatives exist (do nothing, LLM grouping, current behavior), not the required 3 meaningful alternatives PLUS the selected approach. |
| Honest trade-off comparison | 11/25 | Improved: "新代码量可控" replaces the stale "复用现有代码." The same-root-cause risk is now honestly documented in the comparison table. Deduction: the comparison table still lists no quantitative trade-offs. "最小改动" vs actual 6-8 hours + duplicated task creation logic is not minimal. |
| Chosen approach justified against benchmarks | 10/25 | "最小改动，最大通用性" remains a slogan. No quantitative comparison. The proposal could have cited JUnit XML's `<testsuite>` grouping as the industry standard and explained why naming-convention detection was chosen over structured XML parsing. |

### 4. Requirements Completeness: 75/110

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Scenario coverage | 32/40 | Improved with "栈 trace 引用辅助测试文件" scenario (line 43). Happy path, single file, non-standard naming, all-passing, helper file false positive covered. Missing: (1) very large output (1000+ lines) performance, (2) concurrent quality-gate runs producing interleaved output, (3) test output containing file paths in stack traces that reference files not in `sourceExts` (e.g., `.mod` files). |
| Non-functional requirements | 22/40 | Performance: "时间可忽略" — still asserted without evidence. Compatibility: "不影响 compile/fmt/lint/unit-test 步骤" — improved but implementation path (direct API call) risks duplicating 95 lines of logic from `addSingleFixTask`. No mention of: correctness of line association, maximum output size handling, memory usage for large test outputs. |
| Constraints & dependencies | 21/30 | Improved: `extractFileLineMap` dependency on different extraction logic now stated (line 55). Missing: constraint on `extractFileLineMap`'s ability to handle all test output formats across the 5 languages despite using "与 `sourceFileRe` 不同的提取逻辑." The "different logic" is unspecified. |

### 5. Solution Creativity: 35/100

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Novelty over industry baseline | 16/40 | Acknowledged "无创新." The `extractFileLineMap` function preserving file-to-line mapping is engineering, not creativity. The soft cap of 10 for regression tasks is a minor practical addition. |
| Cross-domain inspiration | 5/35 | No cross-domain ideas. Could have drawn from: Sentry fingerprint-based grouping (same root cause deduplication), distributed tracing span-linking (cross-referencing related failures), IDE test runner failure tree views. |
| Simplicity of insight | 14/25 | The insight ("split by test file") is simple. The soft cap addition is a practical refinement. But the implementation complexity (duplicating or extracting 95 lines from `addSingleFixTask`) undermines the simplicity claim. |

### 6. Feasibility: 65/100

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Technical feasibility | 28/40 | Core mechanism is buildable. Cap bypass path specified. `extractFileLineMap` signature provided. Deduction: (1) the return type `map[string][]int` cannot support the described algorithm (context windows require actual line content, not just line numbers). (2) `addRegressionFixTasks` must replicate 95 lines of `addSingleFixTask` logic (surface inference, template defaults, task creation, markdown, state update) or extract shared code — neither option is scoped. |
| Resource & timeline | 20/30 | "6-8 小时" is closer to realistic than the original 2-3 hours. But still does not account for: (a) extracting or duplicating `addSingleFixTask`'s 95 lines of non-cap logic, (b) building test-framework-specific extraction logic for 5 languages, (c) implementing and testing the soft cap of 10. Realistic estimate: 1-2 days. |
| Dependency readiness | 17/30 | No external dependencies. `sourceFileRe` exists. Deduction: `extractFileLineMap` requires "与 `sourceFileRe` 不同的提取逻辑" — this new extraction logic is not specified and has no existing dependency to build on. Line 86 still cites "`extractSourceFiles` 已稳定运行" as evidence, which is irrelevant to the new function. |

### 7. Scope Definition: 64/80

| Criterion | Score | Justification |
|-----------|-------|---------------|
| In-scope items are concrete | 26/30 | Function signatures, algorithm steps, and soft cap are concrete. Line 111 specifies the implementation path (direct API call). Deduction: the soft cap of 10 is in Risk mitigation but not in In Scope — should be explicitly listed as a deliverable. |
| Out-of-scope explicitly listed | 22/25 | Five items explicitly out of scope. Layer 2 deferral with rationale. Good. |
| Scope is bounded | 16/25 | "改动集中在 `quality_gate.go`" — bounded. But the direct API call path (bypassing `addSingleFixTask`) requires either code duplication or extraction of shared logic from `addSingleFixTask`. The proposal does not scope this refactoring/duplication work. The 95 lines of non-cap logic in `addSingleFixTask` represent hidden scope. |

### 8. Risk Assessment: 62/90

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Risks identified | 24/30 | 6 risks listed (up from 5 in iteration-1). Rollback plan added. Same-root-cause risk now honestly documented. Missing: (1) `addSingleFixTask` logic duplication risk (95 lines of code to replicate), (2) `extractFileLineMap` return type mismatch with algorithm (specification bug). |
| Likelihood + impact rated | 18/30 | Improved from iteration-1. Same-root-cause risk M/M — reasonable. Soft cap risk M/M — reasonable. Deduction: same-root-cause risk impact should arguably be H — two agents producing conflicting edits to the same production file, with no automated merge strategy, requiring manual intervention is a high-impact degradation. |
| Mitigations are actionable | 20/30 | Rollback plan is actionable (line 135). Soft cap of 10 is actionable. Same-root-cause mitigation ("需人工介入解决冲突") is honest. Deduction: the mitigation for same-root-cause risk is "accept this trade-off" which is a valid position but not a mitigation — it is an acceptance. A real mitigation would be: "add a `RELATED_TASKS` field to the fix task description so the agent knows about sibling tasks." |

### 9. Success Criteria: 58/80

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Criteria are measurable and testable | 24/30 | SC-1 through SC-6 are testable. The contradiction from iteration-1 (multi-matching lines) is resolved. Deduction: SC-2 ("包含该文件路径的行及上下文") is testable in principle but the format of "上下文" is unspecified (raw lines? annotated? how many lines?). |
| Coverage is complete | 16/25 | Missing SC for: (1) soft cap of 10 regression fix tasks (Risk-3 specifies this but no SC verifies it), (2) correctness of line association (SC-2 tests presence of lines but not that overlapping context windows are correctly deduplicated), (3) description format quality (is the agent's description actually useful for fixing?). |
| SC internal consistency | 18/25 | SC-3 vs SC-5 contradiction resolved. SC-1 + SC-3 + SC-5: satisfiable. **Issue**: SC-3 says "不受 cap 限制" but Risk-3 introduces a soft cap of 10. The SC does not reflect the soft cap, creating a discrepancy between SC and risk mitigation. **Ambiguous — requires author clarification**: Is the soft cap an enforced requirement or a guideline? |

### 10. Logical Consistency: 55/90

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Solution addresses the stated problem | 25/35 | Improved. The proposal now honestly positions itself as layer 1 only. The deferral rationale is valid. Deduction: the same-root-cause scenario (acknowledged in Risk-2) means the solution can create worse outcomes for this case — two agents editing the same production file simultaneously. The proposal acknowledges this trade-off, which is honest, but the logical consistency suffers because the stated goal is "解决 agent 卡死" yet the solution introduces a new failure mode (conflicting edits). |
| Scope <-> Solution <-> SC aligned | 15/30 | Improved from iteration-1. The cap bypass mechanism is now specified. Remaining misalignment: (1) soft cap of 10 is in Risk mitigation but not in Scope or SC, (2) `extractFileLineMap` return type (`map[string][]int`) cannot support the algorithm (which requires actual line content for context windows), (3) direct API call path bypasses `addSingleFixTask` but the 95 lines of shared logic are not scoped. |
| Requirements <-> Solution coherent | 15/25 | Improved. The naming convention constraint is coherent. The partial-solution acknowledgment is coherent. Deduction: the NFR "性能：时间可忽略" is asserted without evidence. The NFR "兼容性：不影响 compile/fmt/lint/unit-test 步骤" depends on the direct API call path not affecting `addSingleFixTask`'s behavior — achievable but requires careful separation of the 95 lines of shared logic. |

## Phase 3: Blindspot Hunt

### [blindspot-1] `extractFileLineMap` return type incompatible with described algorithm

The function signature on line 103 is `func extractFileLineMap(output string) map[string][]int` — returning file paths to line NUMBER slices. But the algorithm on lines 105-109 requires constructing descriptions containing actual output lines (matched lines + context windows of 2 lines before/after + deduplication of overlapping windows). Line numbers alone cannot represent this — you need the actual line content, and you need to handle deduplication of overlapping windows. The return type should be `map[string][]string` (file path to actual output lines) or a richer structure. This is a specification-level bug that would surface immediately during implementation.

### [blindspot-2] `addSingleFixTask` contains 95 lines of essential non-cap logic that must be replicated

Lines 708-802 of `quality_gate.go` show that `addSingleFixTask` performs: surface inference (line 727), task type derivation (line 746), template defaults loading (lines 753-756), `AddTaskOpts` construction (lines 758-774), template validation (lines 778-783), `task.AddTask` call (line 785), `task.CreateTaskMarkdown` call (line 792), and `feature.EnsureForgeState` call (line 796). The proposal says `addRegressionFixTasks` "直接调用底层 task 创建 API（绕过 `addSingleFixTask` 的 cap 检查）" — but "底层 task 创建 API" means `task.AddTask`, `task.CreateTaskMarkdown`, and `feature.EnsureForgeState`. The proposal must either: (a) duplicate all the preparation logic (surface inference, opts construction, template validation) in `addRegressionFixTasks`, or (b) extract a shared `createFixTaskInternal` helper from `addSingleFixTask`. Option (a) creates code duplication. Option (b) is a refactoring of `addSingleFixTask` that affects all its callers. Neither option is scoped or estimated.

### [blindspot-3] Soft cap of 10 is in Risk mitigation but not in Success Criteria

Line 130 introduces "regression 路径最多创建 10 个 fix task，超出部分合并到按目录分组的 fallback." This is a critical behavioral constraint, but it appears only in a Risk mitigation cell. No Success Criterion verifies it. If the soft cap is a requirement, it should have an SC entry: "regression 路径创建的 fix task 总数不超过 10，超出部分 fallback 到按目录分组." Without this SC, the soft cap is untestable and may be forgotten during implementation.

### [blindspot-4] Same-root-cause mitigation is acceptance, not mitigation

Risk 2 (line 128) acknowledges that "同一根因 bug 导致多个测试文件失败时创建冲突修复任务" with "mitigation" being "接受此 trade-off：...需人工介入解决冲突." This is not a mitigation — it is acceptance of the risk. A real mitigation would actively reduce the risk: e.g., adding a `RELATED_TASKS` field to each fix task's description so the agent knows about sibling tasks targeting the same production files, or detecting shared production file references and co-locating those tasks. The proposal's SC-5 ("现有 compile/fmt/lint/unit-test 步骤的 fix task 创建不受影响") would be better supplemented with a new SC: "每个 fix task 的 description 包含相关联的其他 fix task 信息."

### [blindspot-5] `sourceFileRe` was NOT designed for test framework output — new extraction logic is unspecified

Line 55 acknowledges the need for "与 `sourceFileRe` 不同的提取逻辑以适配测试框架输出格式" and cites examples (pytest `FAILED file::class::method`, Go `--- FAIL:` 块). But the actual extraction logic for `extractFileLineMap` is completely unspecified. What regex or parsing strategy handles `FAILED tests/test_foo.py::TestClass::test_method`? The file path is `tests/test_foo.py` but the `::class::method` suffix is not a `:line` pattern. Go's `--- FAIL: TestName` header doesn't contain a file path at all — the file path appears on subsequent lines. The proposal has acknowledged the problem but provided no solution for it. The `extractFileLineMap` function is the core technical deliverable and its extraction strategy is undefined.

### [blindspot-6] Feasibility section still cites irrelevant evidence

Line 86: "无外部依赖。`extractSourceFiles` 已稳定运行。" But the proposal builds `extractFileLineMap` with different extraction logic (line 55). The stability of `extractSourceFiles` says nothing about the feasibility of the new extraction logic, which must handle fundamentally different output formats across 5 languages. This was flagged in iteration-1 attack #14 and remains unaddressed.

### [blindspot-7] `extractFileLineMap` needs per-language extraction strategies but none are specified

The 5 supported languages (Go, Python, JS/TS, Java, Ruby) have fundamentally different test output formats:
- Go: `--- FAIL: TestName (0.00s)` header, then `file_test.go:42: Error message` on indented lines
- Python (pytest): `FAILED tests/test_foo.py::TestClass::test_method` on one line, then multi-line traceback
- JS/TS (Jest): `FAIL tests/foo.test.ts` header, then `● TestSuite > test name` blocks
- Java (JUnit/Maven): `Tests run: N, Failures: M` summary, then stack traces with `at com.example.TestClass.testMethod(TestClass.java:42)`
- Ruby (Minitest/RSpec): Different formats depending on runner

The proposal provides no extraction strategy for any of these. The `extractFileLineMap` function's "与 `sourceFileRe` 不同的提取逻辑" is a placeholder for what is actually the hardest technical problem in the entire proposal.

### [blindspot-8] "Go 专属 suite 解析（原提案 v1）" remains a straw-man alternative

The comparison table (line 67) includes "Go 专属 suite 解析（原提案 v1）" as a separate alternative. This is the proposal's own previous iteration, presented to be rejected. It is not a genuinely different approach. A real alternative would be: JUnit XML structured output parsing, stack-trace fingerprint-based grouping (like Sentry), or AST-based test-to-source mapping. This was flagged in iteration-1 attack #12 and remains unaddressed.

## Score Summary

| Dimension | Score | Max |
|-----------|-------|-----|
| Problem Definition | 78 | 110 |
| Solution Clarity | 78 | 120 |
| Industry Benchmarking | 55 | 120 |
| Requirements Completeness | 75 | 110 |
| Solution Creativity | 35 | 100 |
| Feasibility | 65 | 100 |
| Scope Definition | 64 | 80 |
| Risk Assessment | 62 | 90 |
| Success Criteria | 58 | 80 |
| Logical Consistency | 55 | 90 |
| **Total** | **625** | **1000** |

## Attack Points

1. **Solution Clarity**: `extractFileLineMap` return type `map[string][]int` cannot support the described algorithm — Line 103: `func extractFileLineMap(output string) map[string][]int` returns line numbers, but the algorithm (lines 105-109) requires constructing descriptions with actual output lines + context windows + deduplication. Line numbers alone cannot represent context windows or deduplicated content. — Must change return type to `map[string][]string` or define a richer structure that can carry actual line content.

2. **Feasibility**: `addSingleFixTask` contains 95 lines of essential non-cap logic (surface inference, template defaults, opts construction, task creation, markdown, state update) that `addRegressionFixTasks` must replicate — Line 111: "直接调用底层 task 创建 API（绕过 `addSingleFixTask` 的 cap 检查）" — but bypassing `addSingleFixTask` means bypassing ALL its logic, not just the cap check. The proposal must either scope code duplication or a shared-helper extraction. — Must add to In Scope: either "extract shared `createFixTaskCore` helper from `addSingleFixTask`" or "replicate task creation logic in `addRegressionFixTasks`" with explicit acknowledgment of the maintenance trade-off.

3. **Success Criteria**: Soft cap of 10 is only in Risk mitigation, not in Success Criteria — Line 130: "regression 路径最多创建 10 个 fix task" — this is a behavioral constraint with no SC verifying it. If the soft cap is a requirement, it must be testable. — Must add SC: "regression 路径创建的 fix task 总数不超过 10，超出部分 fallback 到按目录分组."

4. **Risk Assessment**: Same-root-cause "mitigation" is acceptance, not mitigation — Line 128: "接受此 trade-off：...需人工介入解决冲突" — acceptance is honest but does not reduce the risk. — Must add an active mitigation: e.g., "每个 fix task 的 description 包含 `RELATED_FILES` 字段标注所有相关 fix task 引用的生产代码文件，供 agent 检查并发冲突."

5. **Feasibility**: `extractFileLineMap` extraction strategy for 5 languages is completely unspecified — Line 55: "使用与 `sourceFileRe` 不同的提取逻辑以适配测试框架输出格式（如 pytest 的 `FAILED file::class::method`、Go 的 `--- FAIL:` 块）" — acknowledges the need but provides no extraction strategy. How does the function parse `FAILED tests/test_foo.py::TestClass::test_method` to extract `tests/test_foo.py`? How does it handle Go's `--- FAIL:` headers that don't contain file paths? — Must specify the extraction strategy: either (a) per-language regex patterns, (b) a generalized parser, or (c) an extension of `sourceFileRe` with additional patterns. The current specification is a placeholder for the hardest technical problem.

6. **Feasibility**: Dependency Readiness section cites irrelevant evidence — Line 86: "`extractSourceFiles` 已稳定运行" — but the proposal builds `extractFileLineMap` with different extraction logic. The stability of `extractSourceFiles` is irrelevant to the new function. — Must replace with: "`extractFileLineMap` 的提取逻辑需要针对 5 种语言的测试输出格式分别验证，无现有覆盖."

7. **Industry Benchmarking**: Straw-man alternative persists — Line 67: "Go 专属 suite 解析（原提案 v1）" is the proposal's own previous iteration. — Must replace with a genuinely different approach: e.g., "JUnit XML 结构化解析（解析标准测试报告 XML）" or "Sentry 风格指纹去重（按栈 trace 相似度分组）."

8. **Solution Clarity**: Description format inside fix task still underspecified — Line 105: "description 格式为：测试文件路径 + 筛选后的相关输出行（含上下文窗口）" — "筛选后的相关输出行（含上下文窗口）" is vague. Are the lines raw? Annotated with match markers? How is the context window delineated from the matched line? What if the context window overlaps with another test file's context? — Must provide a concrete example of the description format, e.g., a template string showing exact output structure.

9. **Logical Consistency**: Soft cap contradicts "bypasses cap" framing — SC-3 (line 142) says "`addRegressionFixTasks` 不受 cap 限制" but Risk-3 (line 130) introduces a soft cap of 10. The proposal simultaneously claims "no cap" and "cap of 10." — Must reconcile: either update SC-3 to "regression 路径使用独立的 soft cap（上限 10）" or remove the soft cap and document that regression task count is unbounded.

10. **Feasibility**: Timeline still optimistic given unspec'd extraction logic — Line 82: "6-8 小时实现 + 测试" — but building per-language test output extraction for 5 frameworks + context window algorithm + soft cap + direct API call path + deduplication logic + unit tests is realistically 1.5-2 days. — Must revise to "1-2 天" or provide a work breakdown showing how each deliverable fits within 6-8 hours.

11. **Requirements Completeness**: Missing edge case — very large output (1000+ lines) — `extractFileLineMap` must parse the entire output string, build line-number mappings, construct context windows, and deduplicate. For a 1000+ line output with 20+ matching files, the memory and CPU cost is unstated. — Must add NFR: "`extractFileLineMap` 处理 1000 行输出耗时不超过 [N]ms，内存占用不超过 [N]MB."

12. **Logical Consistency**: `extractFileLineMap` return type vs algorithm is a specification bug — The function returns line numbers (`[]int`) but the algorithm on lines 106-109 constructs descriptions with actual line content and context windows. The function cannot both return line numbers AND produce the described description format. Something must bridge the gap: either the function returns actual lines, or a second function constructs descriptions from line numbers. — Must either change the return type to `map[string][]string` or add a second function `buildFileDescription(output string, lineMap map[string][]int) map[string]string` to In Scope.
