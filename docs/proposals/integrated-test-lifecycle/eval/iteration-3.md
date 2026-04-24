---
date: "2026-04-24"
doc_dir: "docs/proposals/integrated-test-lifecycle/"
iteration: "3"
target_score: "N/A"
evaluator: Claude (automated, adversarial)
---

# Proposal Eval — Iteration 3

**Score: 82/100** (target: N/A)

```
┌─────────────────────────────────────────────────────────────────┐
│                    PROPOSAL QUALITY SCORECARD                     │
├──────────────────────────────┬──────────┬──────────┬────────────┤
│ Dimension                    │ Score    │ Max      │ Status     │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 1. Problem Definition        │  17      │  20      │ ✅          │
│    Problem clarity           │   6/7    │          │            │
│    Evidence provided         │   6/7    │          │            │
│    Urgency justified         │   5/6    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 2. Solution Clarity          │  16      │  20      │ ⚠️          │
│    Approach concrete         │   6/7    │          │            │
│    User-facing behavior      │   5/7    │          │            │
│    Differentiated            │   5/6    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 3. Alternatives Analysis     │  13      │  15      │ ⚠️          │
│    Alternatives listed (≥2)  │   5/5    │          │            │
│    Pros/cons honest          │   4/5    │          │            │
│    Rationale justified       │   4/5    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 4. Scope Definition          │  15      │  15      │ ✅          │
│    In-scope concrete         │   5/5    │          │            │
│    Out-of-scope explicit     │   5/5    │          │            │
│    Scope bounded             │   5/5    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 5. Risk Assessment           │  14      │  15      │ ✅          │
│    Risks identified (≥3)     │   5/5    │          │            │
│    Likelihood + impact rated │   5/5    │          │            │
│    Mitigations actionable    │   4/5    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 6. Success Criteria          │  12      │  15      │ ⚠️          │
│    Measurable                │   5/5    │          │            │
│    Coverage complete         │   3/5    │          │            │
│    Testable                  │   4/5    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ Deductions                   │  -5      │          │            │
│    Scope/criteria gap (docs) │  -3      │          │            │
│    Alternatives table stale  │  -2      │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ TOTAL                        │  82      │  100     │            │
└──────────────────────────────┴──────────┴──────────┴────────────┘
```

---

## Deductions

| Location | Issue | Penalty |
|----------|-------|---------|
| Scope § In Scope vs. Success Criteria | In-scope lists "更新 OVERVIEW.md 和 test-capture-design.md" and "更新 Prerequisites 说明" — neither appears in any success criterion. Persistent from iteration 2, unaddressed. | -3 pts |
| Alternatives § 方案 A cons vs. Solution § 核心机制 1 | Alternatives table lists "/breakdown-tasks 需要感知是否需要 e2e 测试" as a con for 方案 A, but the solution resolves this with "始终追加" — the cons column was not updated to reflect the design decision | -2 pts |

---

## Attack Points

### Attack 1: Solution Clarity — graduation trigger "首次成功" has no implementation mechanism

**Where**: Data flow diagram: "成功（首次）→ 写 latest.md → 毕业：迁移脚本到 tests/e2e/<type>/<target>/ → exit 0" vs. "成功（非首次）→ 写 latest.md → exit 0"

**Why it's weak**: The data flow distinguishes two success paths — first-time success triggers graduation (file migration), subsequent successes do not. But the proposal never specifies how `task all-completed` determines which case it's in. There is no flag file, no state field in `index.json`, no check described. Does it look for the existence of `tests/e2e/<type>/<target>/`? Does it write a `.graduated` sentinel? If the directory already exists from a previous feature, does the current feature's first success still trigger graduation? The graduation model is the highest-complexity deliverable in Phase 2, yet its triggering condition is undefined. An implementer reading this proposal cannot write the conditional branch.

**What must improve**: Specify the graduation detection mechanism explicitly. For example: "首次成功定义为：`tests/e2e/` 中不存在来自本 feature（`<slug>`）的任何脚本文件。`task all-completed` 在迁移前检查 `tests/e2e/**/*` 中是否有 `# source: <slug>` 注释头；若无，则执行迁移。" Or use a simpler sentinel: write `docs/features/<slug>/testing/.graduated` on first success and skip migration if it exists.

---

### Attack 2: Success Criteria — documentation deliverables have no criteria (persistent from iteration 2)

**Where**: Scope § In Scope: "更新 task-cli/docs/OVERVIEW.md 和 docs/todos/test-capture-design.md 反映新行为" and "更新 gen-test-cases 和 gen-test-scripts skill 的 Prerequisites 说明，明确它们也可以作为任务被 agent 调用"

**Why it's weak**: These are explicit in-scope deliverables — two documentation files and two skill Prerequisite sections. None of the eight success criteria mention them. This was flagged in iteration 2 as a -3 deduction and remains completely unaddressed. A reviewer cannot verify completion of these deliverables because there is no acceptance condition. "OVERVIEW.md 反映新行为" is not a criterion — it's a description of intent.

**What must improve**: Add criteria such as: "OVERVIEW.md 中包含对毕业模型目录结构（`tests/e2e/<type>/<target>/`）的说明" and "gen-test-cases skill 的 Prerequisites 中明确列出 '可作为 T-test-1 任务被 agent 通过 task claim 调用'". Each documentation deliverable needs at least one verifiable condition.

---

### Attack 3: Alternatives Analysis — method A's cons column is stale and contradicts the solution

**Where**: Alternatives table, 方案 A row, Cons column: "/breakdown-tasks 需要感知是否需要 e2e 测试"; Solution § 核心机制 1: "在所有业务任务之后，自动追加两个固定任务" and Risk 2 mitigation: "始终追加（与核心机制 1 一致）；agent 执行 T-test-1 时若 PRD 无 UI/API/CLI 需求，将任务标记为 skipped"

**Why it's weak**: The alternatives table was written before the "始终追加" design decision was made. The solution and risk sections have been updated to reflect the decision, but the alternatives table still lists the pre-decision uncertainty as a con. A reader comparing the table to the solution will find a contradiction: the table implies the approach requires sensing logic that the solution explicitly eliminates. The alternatives table is the primary artifact a skeptic reads to evaluate the tradeoff — having a stale con undermines the credibility of the entire analysis.

**What must improve**: Update the cons column for 方案 A to reflect the actual design: "始终追加两个固定任务；无 e2e 需求的 feature 会产生两个 skipped 任务（可接受的噪音）". This is honest about the real tradeoff (task list noise) rather than a resolved concern.

---

## Previous Issues Check

| Previous Attack | Addressed? | Evidence |
|----------------|------------|----------|
| Attack 1: Test ID is the linchpin of the graduation model but is never defined | ✅ | Now defines format explicitly: "`<type>/<target>/<slug>`，其中 slug 由测试用例标题规范化生成（小写、空格转连字符、去除标点）"，and specifies ID is auto-generated by `gen-test-cases`, stable if title unchanged |
| Attack 2: Risk 2 mitigation is an unresolved design fork contradicting the solution | ✅ | Risk 2 mitigation now reads "始终追加（与核心机制 1 一致）；agent 执行 T-test-1 时若 PRD 无 UI/API/CLI 需求，将任务标记为 skipped 并说明原因" — contradiction resolved, one behavior chosen |
| Attack 3: Urgency names the risk but not the consequence | ✅ | Now states concrete hypothetical: "若该 feature 引入了 hook 的回归缺陷，缺陷将在零 e2e 覆盖的情况下合并，此后所有 feature 的自动化闭环都将静默失效" |

---

## Verdict

- **Score**: 82/100
- **Target**: N/A
- **Gap**: N/A
- **Action**: All three iteration-2 attacks addressed. Three new attacks identified. Address graduation trigger mechanism (Attack 1), documentation success criteria (Attack 2, persistent), and stale alternatives table (Attack 3) to reach 88+.
