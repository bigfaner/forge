---
date: "2026-05-08"
doc_dir: "docs/proposals/eval-test-cases/"
iteration: "3"
target_score: "90"
evaluator: Claude (automated, adversarial)
---

# Proposal Eval -- Iteration 3

**Score: 93/100** (target: 90)

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
| 2. Solution Clarity          |  18      |  20      | :warning:  |
|    Approach concrete         |  7/7     |          |            |
|    User-facing behavior      |  6/7     |          |            |
|    Differentiated            |  5/6     |          |            |
+------------------------------+----------+----------+------------+
| 3. Alternatives Analysis     |  14      |  15      | :warning:  |
|    Alternatives listed (>=2) |  5/5     |          |            |
|    Pros/cons honest          |  5/5     |          |            |
|    Rationale justified       |  4/5     |          |            |
+------------------------------+----------+----------+------------+
| 4. Scope Definition          |  15      |  15      | :white_check_mark:  |
|    In-scope concrete         |  5/5     |          |            |
|    Out-of-scope explicit     |  5/5     |          |            |
|    Scope bounded             |  5/5     |          |            |
+------------------------------+----------+----------+------------+
| 5. Risk Assessment           |  14      |  15      | :white_check_mark:  |
|    Risks identified (>=3)    |  5/5     |          |            |
|    Likelihood + impact rated |  5/5     |          |            |
|    Mitigations actionable    |  4/5     |          |            |
+------------------------------+----------+----------+------------+
| 6. Success Criteria          |  14      |  15      | :white_check_mark:  |
|    Measurable                |  4/5     |          |            |
|    Coverage complete         |  5/5     |          |            |
|    Testable                  |  5/5     |          |            |
+------------------------------+----------+----------+------------+
| TOTAL                        |  93      |  100     |            |
+------------------------------+----------+----------+------------+
```

---

## Deductions

| Location | Issue | Penalty |
|----------|-------|---------|
| Problem:evidence | Evidence is empirically grounded (forge test-e2e feature, 25% defect rate, 40k tokens) but relies on a single incident (n=1 feature, 12 TCs). No second data point or mention of whether this pattern recurs across other features. | -1 pt (evidence) |
| Problem:urgency | "每 4 个 feature 中就有 1 个需要额外修复循环" extrapolates from n=1. The 25% rate comes from one feature's 12 TCs -- not statistically meaningful enough to project across all features. | -1 pt (urgency) |
| Solution:user-facing | Invocation inputs are specified ("输入 test-cases.md + PRD 源文件") but the user interaction model is split between automatic (via T-test-1b task chain) and manual (independent skill invocation) without explicitly describing both modes. The architecture decision says "独立 skill" but the user-facing description does not address when/how a user invokes it outside the task chain. | -1 pt (user-facing) |
| Solution:differentiated | "复用现有 doc-scorer + doc-reviser adversarial loop 架构" means the core mechanism is identical to eval-prd/eval-design. The differentiator is the rubric content ("评估其'下游可执行性'"), which is stated but not contrasted explicitly against what eval-prd/eval-design evaluate. A one-sentence comparison would close this gap. | -1 pt (differentiated) |
| Alternatives:rationale | The "do nothing" alternative (last row) says "脚本生成后才发现问题，修复成本高" but does not quantify "高" using the 40k tokens data from the problem section or the full cost comparison from the risk section. The risk section quantifies no-eval cost at ~85k tokens, yet the alternatives table still uses vague "修复成本高". | -1 pt (rationale) |
| Risk:mitigations | "阈值可根据实战反馈调整；6 轮迭代给足修复空间" for the Step Actionability threshold risk: "可根据实战反馈调整" is vague (what feedback triggers adjustment? what is the adjustment process?), and "6 轮迭代给足修复空间" is an assertion that the problem won't occur, not a mitigation for when it does. | -1 pt (mitigations) |
| Success:measurable | "reviser 仅修改 test-cases.md" is a behavioral constraint, not a measurable criterion. How is this verified -- by checking that no other files were modified after an eval run? The criterion lacks a verification method for this constraint. | -1 pt (measurable) |

---

## Attack Points

### Attack 1: Alternatives Analysis -- "do nothing" rationale is vague despite available data

**Where**: Alternatives table, last row: "脚本生成后才发现问题，修复成本高"

**Why it's weak**: The problem section provides concrete repair cost data (40k tokens, 3 failed scripts, 25% defect rate). The risk section provides a full cost comparison (eval ~40k vs no-eval ~85k). Yet the alternatives table still uses the hand-wave "修复成本高" without citing any of these numbers. This is a cross-section quality inconsistency: the same document has rigorous quantification in two places but leaves the alternatives analysis vague. The "do nothing" alternative is the most important one to refute convincingly, and the data exists to do so -- it is just not used here.

**What must improve**: Replace "修复成本高" with a specific cost reference, e.g., "实测修复成本 ~40k tokens + 无效脚本生成 ~45k tokens = ~85k tokens/feature (见 Problem 节和 Risk 节量化数据)."

### Attack 2: Solution Clarity -- user interaction model is implicit

**Where**: Proposed Solution section describes inputs ("输入 test-cases.md + PRD 源文件") and task chain integration ("T-test-1 → T-test-1b → T-test-2"), but never explicitly states the two invocation modes a user experiences.

**Why it's weak**: The architecture decision says "形态: 独立 skill" with rationale "不增加 gen-test-cases 的复杂度." But the user-facing description focuses exclusively on the automatic task-chain path. A reader cannot determine: can I manually run `/eval-test-cases` on an existing test-cases.md outside the task chain? What does that experience look like? The dual invocation model (automatic via T-test-1b, manual as independent skill) is a key user-facing behavior that remains implicit.

**What must improve**: Add a brief user-facing behavior subsection or expand the task integration section to explicitly describe both modes: (1) automatic invocation via T-test-1b in the task chain, and (2) manual invocation as an independent skill (command, inputs, output location).

### Attack 3: Problem Definition -- urgency rests on a single-incident sample

**Where**: Urgency section: "基于上述实测数据（25% TC 质量缺陷率、40k tokens 修复成本），缺少质量门控意味着每 4 个 feature 中就有 1 个需要额外修复循环."

**Why it's weak**: The 25% defect rate is derived from a single feature (forge test-e2e, 12 TCs, 3 defective). Extrapolating "每 4 个 feature 中就有 1 个" from n=1 is not statistically defensible. The argument would be stronger by acknowledging the sample limitation and supplementing with a structural argument (e.g., "gen-test-cases has no validation step by design, so defect rate depends entirely on single-pass LLM output quality, which is inherently variable"). The current phrasing presents a single data point as if it were a stable rate.

**What must improve**: Either (a) add a second data point from another feature's experience, or (b) rephrase to acknowledge sample size and supplement with a structural argument about why the defect rate is likely non-trivial across features.

---

## Previous Issues Check

| Previous Attack | Addressed? | Evidence |
|----------------|------------|----------|
| Attack 1: Scope Definition -- scope is still unbounded in time and effort | Yes | Each in-scope item now has effort estimate (e.g., "预计 ~1 session", "预计 ~15 min"). Total estimate added: "总估算：~2 sessions（约 3-4 小时），其中 run-tasks dispatch 改动占约一半工作量." The riskiest item (run-tasks dispatch) is flagged as "涉及 run-tasks 核心分发路径...是本 scope 中复杂度最高的改动." |
| Attack 2: Success Criteria -- testability gap for main_session dispatch behavior | Yes | Success criterion 5 now includes detailed acceptance test: "验收方法：在一个测试 feature 上运行完整任务链（T-test-1 -> T-test-1b -> T-test-2），验证：(1) T-test-1b 执行时 task-executor 未被 spawn（检查 agent transcript 中无 task-executor 调用记录）；(2) eval 报告生成到 docs/features/<slug>/eval/iteration-1.md；(3) doc-scorer 和 doc-reviser 的 transcript 存在且包含评分/修订内容." This is a concrete, executable test scenario. |
| Attack 3: Risk Assessment -- ROI assertion without numbers | Yes | Risk table row 4 now includes full cost comparison: "估算单次 eval 消耗：doc-scorer ~5k tokens/轮 + doc-reviser ~8k tokens/轮 x 平均 3 轮 = ~40k tokens...综合比较：eval ~40k vs 无 eval ~85k（生成 + 修复），前置 eval 净节省 ~45k tokens." ROI claim is now numerically supported. |

---

## Verdict

- **Score**: 93/100
- **Target**: 90/100
- **Gap**: -3 points (above target)
- **Action**: Target reached. All three iteration-2 attacks have been addressed. Remaining deductions are minor (single-instance evidence, implicit dual invocation model, vague alternative rationale) and do not prevent implementation. Proceed to `/write-prd`.
