---
date: "2026-05-08"
doc_dir: "docs/proposals/eval-test-cases/"
iteration: "2"
target_score: "90"
evaluator: Claude (automated, adversarial)
---

# Proposal Eval — Iteration 2

**Score: 88/100** (target: 90)

```
+-----------------------------------------------------------------+
|                    PROPOSAL QUALITY SCORECARD                     |
+------------------------------+----------+----------+------------+
| Dimension                    | Score    | Max      | Status     |
+------------------------------+----------+----------+------------+
| 1. Problem Definition        |  18      |  20      | :warning:  |
|    Problem clarity           |  7/7     |          |            |
|    Evidence provided         |  6/7     |          |            |
|    Urgency justified         |  5/6     |          |            |
+------------------------------+----------+----------+------------+
| 2. Solution Clarity          |  18      |  20      |  :white_check_mark:  |
|    Approach concrete         |  7/7     |          |            |
|    User-facing behavior      |  6/7     |          |            |
|    Differentiated            |  5/6     |          |            |
+------------------------------+----------+----------+------------+
| 3. Alternatives Analysis     |  14      |  15      |  :white_check_mark:  |
|    Alternatives listed (>=2) |  5/5     |          |            |
|    Pros/cons honest          |  5/5     |          |            |
|    Rationale justified       |  4/5     |          |            |
+------------------------------+----------+----------+------------+
| 4. Scope Definition          |  13      |  15      | :warning:  |
|    In-scope concrete         |  5/5     |          |            |
|    Out-of-scope explicit     |  5/5     |          |            |
|    Scope bounded             |  3/5     |          |            |
+------------------------------+----------+----------+------------+
| 5. Risk Assessment           |  13      |  15      |  :white_check_mark:  |
|    Risks identified (>=3)    |  5/5     |          |            |
|    Likelihood + impact rated |  4/5     |          |            |
|    Mitigations actionable    |  4/5     |          |            |
+------------------------------+----------+----------+------------+
| 6. Success Criteria          |  12      |  15      | :warning:  |
|    Measurable                |  4/5     |          |            |
|    Coverage complete         |  4/5     |          |            |
|    Testable                  |  4/5     |          |            |
+------------------------------+----------+----------+------------+
| TOTAL                        |  88      |  100     |            |
+------------------------------+----------+----------+------------+
```

---

## Deductions

| Location | Issue | Penalty |
|----------|-------|---------|
| Problem:evidence | Evidence is now empirically grounded (forge test-e2e feature, 25% defect rate, 40k tokens) but relies on a single incident. No second data point or mention of whether this pattern recurs across other features. | -1 pt (evidence) |
| Problem:urgency | "每 4 个 feature 中就有 1 个需要额外修复循环" extrapolates from n=1 feature. The 25% rate is from a single sample of 12 TCs in one feature, not a statistically meaningful dataset. | -1 pt (urgency) |
| Solution:user-facing | Invocation is described ("输入 test-cases.md + PRD 源文件") but no description of how the user triggers it in the workflow -- is it automatic via T-test-1b, or can a user manually invoke it? The interaction model is implicit, not explicit. | -1 pt (user-facing) |
| Solution:differentiated | "复用现有 doc-scorer + doc-reviser adversarial loop 架构" means the core mechanism is identical to eval-prd/eval-design. The differentiator is the rubric content, not the architecture. This is fine but not called out explicitly. | -1 pt (differentiated) |
| Alternatives:rationale | The "do nothing" alternative's con says "脚本生成后才发现问题，修复成本高" but does not quantify "高" using the 40k tokens data now available in the problem section -- a missed opportunity to strengthen the argument. | -1 pt (rationale) |
| Scope:bounded | Still no timeframe or effort estimate. In-scope items are concrete deliverables but there is no indication of expected complexity or implementation duration. | -2 pts (bounded) |
| Risk:likelihood+impact | "eval 增加每个 feature 的 agent token 消耗" rated High/Low is reasonable, but the mitigation "前置 eval 的 ROI 为正" is still an assertion without a comparative calculation. The problem section provides 40k tokens repair cost but no eval cost estimate to compute ROI. | -1 pt (likelihood+impact) |
| Risk:mitigations | "6 轮迭代给足修复空间" for the threshold-strictness risk is not a mitigation -- it is an assertion that the problem won't occur. A mitigation would specify what happens if 6 rounds are insufficient for Step Actionability to cross 20. | -1 pt (mitigations) |
| Success:measurable | "adversarial loop 默认 target 90, max iterations 6，reviser 仅修改 test-cases.md" -- "仅修改 test-cases.md" is a constraint, not a measurability criterion. How is "only modifies test-cases.md" verified? | -1 pt (measurable) |
| Success:coverage | No criterion for verifying the "更新 CLAUDE.md 中的工作流图" in-scope item. The criterion says "CLAUDE.md 工作流图反映 T-test-1 -> T-test-1b -> T-test-2 链路" but does not specify what "reflects" means -- is a diagram addition sufficient, or must the full task dependency graph be updated? | -1 pt (coverage) |
| Success:testable | "SKILLS.md 注册表包含 eval-test-cases 条目" is testable with a simple grep, but "run-tasks 遇到 main_session: true 标记的 task 时，在 main session 中直接调用 eval-test-cases skill" requires a specific test scenario. No test scenario is described. | -1 pt (testable) |

---

## Attack Points

### Attack 1: Scope Definition — scope is still unbounded in time and effort

**Where**: Scope section lists 7 concrete in-scope deliverables and 5 out-of-scope items, but nowhere does the proposal estimate implementation effort or timeframe.

**Why it's weak**: The scope is itemized but not sized. "修改 run-tasks dispatch 逻辑" touches shared infrastructure with potential ripple effects -- is this a 2-hour change or a 2-day refactor? Without effort sizing, a team cannot commit to the scope. This was flagged in iteration 1 ("No timeframe or effort estimate") and remains unaddressed.

**What must improve**: Add a rough effort estimate per in-scope item or a total implementation estimate (e.g., "estimated 1-2 sessions" or "~4 hours"). At minimum, flag the riskiest item (run-tasks dispatch modification) with an expected complexity level.

### Attack 2: Success Criteria — testability gap for main_session dispatch behavior

**Where**: Success criterion: "run-tasks 遇到 main_session: true 标记的 task 时，在 main session 中直接调用 eval-test-cases skill（不 dispatch 给 task-executor），且 doc-scorer 和 doc-reviser subagent 正常 spawn 并完成评分与修订"

**Why it's weak**: This criterion describes the expected behavior but provides no test scenario. How does one verify this? A concrete test scenario would be: "Given a T-test-1b task with main_session:true after T-test-1 completes, run-tasks executes eval-test-cases in the main session and the eval report appears at docs/features/<slug>/eval/iteration-1.md." The current formulation is an assertion, not a test case.

**What must improve**: Either add a concrete acceptance test scenario alongside the criterion, or reframe the criterion as a pass/fail checklist item with a specific verification method (e.g., "Run the full task chain on a test feature and verify eval-test-cases runs in main session by checking that doc-scorer/doc-reviser subagent transcripts exist").

### Attack 3: Risk Assessment — ROI assertion without numbers

**Where**: Risk table row 4: "eval 增加每个 feature 的 agent token 消耗" with mitigation "相比 gen-test-scripts 失败后重试的消耗，前置 eval 的 ROI 为正"

**Why it's weak**: The problem section provides repair cost data (40k tokens for 3 failed scripts). But the risk mitigation claims positive ROI without estimating the eval's own cost. What is the expected token cost of one eval-test-cases run (doc-scorer + doc-reviser, 6 iterations max)? Without the denominator, "ROI is positive" is a claim, not a calculation. This was partially flagged in iteration 1 ("ROI 为正 is asserted without numbers") -- the problem section now has data, but the risk section still does not use it to compute ROI.

**What must improve**: Estimate eval-test-cases token cost (e.g., "each eval run: ~doc-scorer 5k + doc-reviser 8k per iteration x 3 avg iterations = ~40k tokens") and compare against the 40k repair cost from the problem section. Even a rough estimate makes the ROI claim defensible.

---

## Previous Issues Check

| Previous Attack | Addressed? | Evidence |
|----------------|------------|----------|
| Attack 1: Success Criteria coverage gap for main_session dispatch logic | Partially | Success criteria now includes: "run-tasks 遇到 main_session: true 标记的 task 时，在 main session 中直接调用 eval-test-cases skill（不 dispatch 给 task-executor），且 doc-scorer 和 doc-reviser subagent 正常 spawn 并完成评分与修订" -- coverage gap is closed, but testability remains weak (no concrete test scenario) |
| Attack 2: Risk mitigations are detection/deflection, not prevention | Yes | Risk table row 1 now specifies: "Reviser prompt 内置 invariant guard 指令：明确禁止修改 TC ID 编号格式和 traceability table 结构；每轮结束后 diff 检查 TC ID 集合是否与原始一致，不一致则拒绝本轮修订并回退" -- this is a structural prevention mechanism, not just detection |
| Attack 3: Problem evidence is architectural analogy, not empirical data | Yes | Evidence section now includes concrete incident: "forge 项目自身 test-e2e feature（2026-04）中，gen-test-cases 输出的 test-cases.md 含 2 条 Route 为空和 1 条 Element 描述模糊的 TC...25% TC 质量缺陷率、40k tokens 修复成本" -- empirically grounded with specific numbers |

---

## Verdict

- **Score**: 88/100
- **Target**: 90/100
- **Gap**: 2 points
- **Action**: Continue to iteration 3 — close the scope bounding gap (add effort estimates) and strengthen risk ROI with cost comparison
