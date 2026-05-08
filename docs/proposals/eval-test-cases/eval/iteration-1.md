---
date: "2026-05-08"
doc_dir: "docs/proposals/eval-test-cases/"
iteration: "1"
target_score: "90"
evaluator: Claude (automated, adversarial)
---

# Proposal Eval -- Iteration 1

**Score: 75/100** (target: 90)

```
+-----------------------------------------------------------------+
|                    PROPOSAL QUALITY SCORECARD                     |
+------------------------------+----------+----------+------------+
| Dimension                    | Score    | Max      | Status     |
+------------------------------+----------+----------+------------+
| 1. Problem Definition        |  13      |  20      | :warning:  |
|    Problem clarity           |  5/7     |          |            |
|    Evidence provided         |  4/7     |          |            |
|    Urgency justified         |  4/6     |          |            |
+------------------------------+----------+----------+------------+
| 2. Solution Clarity          |  16      |  20      | :warning:  |
|    Approach concrete         |  6/7     |          |            |
|    User-facing behavior      |  5/7     |          |            |
|    Differentiated            |  5/6     |          |            |
+------------------------------+----------+----------+------------+
| 3. Alternatives Analysis     |  12      |  15      | :warning:  |
|    Alternatives listed (>=2) |  5/5     |          |            |
|    Pros/cons honest          |  3/5     |          |            |
|    Rationale justified       |  4/5     |          |            |
+------------------------------+----------+----------+------------+
| 4. Scope Definition          |  12      |  15      | :warning:  |
|    In-scope concrete         |  4/5     |          |            |
|    Out-of-scope explicit     |  5/5     |          |            |
|    Scope bounded             |  3/5     |          |            |
+------------------------------+----------+----------+------------+
| 5. Risk Assessment           |  11      |  15      | :warning:  |
|    Risks identified (>=3)    |  4/5     |          |            |
|    Likelihood + impact rated |  4/5     |          |            |
|    Mitigations actionable    |  3/5     |          |            |
+------------------------------+----------+----------+------------+
| 6. Success Criteria          |  11      |  15      | :warning:  |
|    Measurable                |  4/5     |          |            |
|    Coverage complete         |  3/5     |          |            |
|    Testable                  |  4/5     |          |            |
+------------------------------+----------+----------+------------+
| TOTAL                        |  75      |  100     |            |
+------------------------------+----------+----------+------------+
```

---

## Deductions

| Location | Issue | Penalty |
|----------|-------|---------|
| Problem:line 11 | Cost claim "3-5 轮 agent 交互来修复（每轮 ~10k tokens）" is a rough estimate presented as evidence, not backed by logs or incident data | -2 pts (evidence) |
| Problem:line 15-17 | Evidence is purely architectural ("eval-prd has it, eval-design has it") with no empirical data -- no failure rates, no user complaints, no past run metrics | -3 pts (evidence) |
| Problem:line 21 | Urgency claims "每个 feature 都会走" but does not quantify how many features per sprint or the actual failure rate | -2 pts (urgency) |
| Solution:line 29 | "下游可执行性" is stated as the evaluation target but "直接驱动" is undefined -- what specific downstream contract must be satisfied? | -1 pt (concrete) |
| Solution:passim | No description of how the user invokes the skill (CLI syntax, parameters, interaction flow) -- only "输入 test-cases.md + PRD 源文件，输出 eval 报告" | -2 pts (user-facing) |
| Alternatives:line 77 | "gen-test-cases 过于复杂，违反单一职责" is architectural dogma without evidence of actual complexity | -2 pts (honest pros/cons) |
| Alternatives:line 80 | "修复成本高" for do-nothing alternative does not reuse the "3-5 轮 * ~10k tokens" data from the problem statement, making it a bare assertion | -2 pts (honest pros/cons) |
| Scope:line 88 | "修改 run-tasks dispatch 逻辑" modifies shared infrastructure but scope does not bound the change (how many files, what complexity) | -1 pt (in-scope concrete) |
| Scope:passim | No timeframe or effort estimate -- cannot assess whether scope is executable in a defined window | -2 pts (bounded) |
| Risk:line 105 | "Structure & ID Integrity 维度（10分）会捕获此类问题" is detection, not mitigation -- reviser could still break things between rounds | -2 pts (mitigations actionable) |
| Risk:line 106 | "用户可查看报告后决定是否接受或手动修复" pushes the decision to the user rather than providing a structural mitigation | -2 pts (mitigations actionable) |
| Risk:line 108 | "eval 的 ROI 为正" is asserted without numbers -- the problem section provides cost data but the risk section does not use it for comparison | -1 pt (likelihood+impact) |
| Success:passim | No success criterion for verifying the main_session:true dispatch logic works correctly (in-scope item with no corresponding criterion) | -2 pts (coverage) |
| Success:line 113 | "可独立调用" is subjective -- what constitutes successful invocation? Needs pass/fail conditions | -1 pt (measurable) |

---

## Attack Points

### Attack 1: Success Criteria -- coverage gap for main_session dispatch logic

**Where**: Scope section lists "修改 run-tasks dispatch 逻辑（识别 main_session: true 标记，在 main session 中直接执行）" as in-scope, but Success Criteria section has no criterion verifying this dispatch behavior works correctly.

**Why it's weak**: The most novel and riskiest technical change in this proposal -- modifying run-tasks to support main_session:true dispatch -- has zero success criteria. The 7 criteria cover the skill's rubric, templates, task dependencies, and registry, but none test that run-tasks actually routes main_session tasks to the main session instead of task-executor. This is a cross-section inconsistency between Scope and Success Criteria.

**What must improve**: Add at least one success criterion like: "run-tasks dispatches a main_session:true task to the main session (not task-executor), and the task executes eval-test-cases successfully with doc-scorer and doc-reviser subagents spawned."

### Attack 2: Risk Assessment -- mitigations are detection or deflection, not prevention

**Where**: Risk table rows 1 and 2: "Rubric 中 Structure & ID Integrity 维度（10分）会捕获此类问题" and "用户可查看报告后决定是否接受或手动修复"

**Why it's weak**: The first mitigation confuses detection with prevention. The rubric can only catch damage after the reviser has already broken structure/IDs. The second mitigation defers to the user rather than providing a structural safety net. Of 5 risks, the two most actionable mitigations are "the scorer will notice" and "the user will decide." This leaves the adversarial loop without guardrails.

**What must improve**: Provide structural mitigations: e.g., "reviser prompt includes an invariant guard instruction to never alter TC IDs or traceability table structure" or "reviser output is diff-checked against original for ID integrity before acceptance."

### Attack 3: Problem Definition -- evidence is architectural analogy, not empirical data

**Where**: Evidence section (lines 15-17): "eval-prd 对 PRD、eval-design 对 tech-design 都有 adversarial loop 质量门控，唯独 gen-test-cases 的输出没有" and "gen-test-cases 的 Step Actionability 完全依赖 LLM 单次生成的质量，没有验证环节"

**Why it's weak**: All three evidence points are structural observations about what exists vs. what does not, not evidence of actual harm. No failure logs, no user reports, no metrics showing how often gen-test-scripts fails due to test-cases.md quality. The cost estimate ("3-5 轮 agent 交互来修复，每轮 ~10k tokens") is presented without attribution -- is this measured, estimated, or hypothetical? A problem without evidence of actual impact is a conjecture.

**What must improve**: Provide at least one concrete instance: a past feature where gen-test-scripts failed and the root cause was traceable to test-cases.md quality. Include measured or estimated token waste. Without empirical grounding, the problem reads as "we noticed a gap" rather than "we experienced pain."

---

## Previous Issues Check

<!-- Only for iteration > 1 -->

| Previous Attack | Addressed? | Evidence |
|----------------|------------|----------|

---

## Verdict

- **Score**: 75/100
- **Target**: 90/100
- **Gap**: 15 points
- **Action**: Continue to iteration 2
