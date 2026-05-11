---
date: "2026-05-11"
doc_dir: "docs/proposals/typed-task-dispatch/"
iteration: "1"
target_score: "85"
evaluator: Claude (automated, adversarial)
---

# Proposal Eval — Iteration 1

**Score: 73/100** (target: 85)

```
┌─────────────────────────────────────────────────────────────────┐
│                    PROPOSAL QUALITY SCORECARD                     │
├──────────────────────────────┬──────────┬──────────┬────────────┤
│ Dimension                    │ Score    │ Max      │ Status     │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 1. Problem Definition        │  15      │  20      │ ⚠️          │
│    Problem clarity           │   6/7    │          │            │
│    Evidence provided         │   5/7    │          │            │
│    Urgency justified         │   4/6    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 2. Solution Clarity          │  14      │  20      │ ⚠️          │
│    Approach concrete         │   6/7    │          │            │
│    User-facing behavior      │   4/7    │          │            │
│    Differentiated            │   4/6    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 3. Alternatives Analysis     │  11      │  15      │ ⚠️          │
│    Alternatives listed (≥2)  │   5/5    │          │            │
│    Pros/cons honest          │   3/5    │          │            │
│    Rationale justified       │   3/5    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 4. Scope Definition          │  12      │  15      │ ⚠️          │
│    In-scope concrete         │   5/5    │          │            │
│    Out-of-scope explicit     │   4/5    │          │            │
│    Scope bounded             │   3/5    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 5. Risk Assessment           │  12      │  15      │ ⚠️          │
│    Risks identified (≥3)     │   5/5    │          │            │
│    Likelihood + impact rated │   3/5    │          │            │
│    Mitigations actionable    │   4/5    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 6. Success Criteria          │   9      │  15      │ ❌          │
│    Measurable                │   3/5    │          │            │
│    Coverage complete         │   3/5    │          │            │
│    Testable                  │   3/5    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ TOTAL                        │  73      │  100     │            │
└──────────────────────────────┴──────────┴──────────┴────────────┘
```

---

## Deductions

| Location | Issue | Penalty |
|----------|-------|---------|
| Success Criteria #6 | "行为与迁移前等价（运行验证）" — "等价" 无量化定义 | -2 pts |
| Success Criteria #4 | "输出结构与当前 T-test-1 执行结果一致" — "结构一致" 无量化定义 | -2 pts |
| Success Criteria #5 | "行为覆盖原 error-fixer 的编译错误、测试失败、lint 修复场景" — "覆盖" 无量化定义 | -2 pts |

---

## Attack Points

### Attack 1: Problem Definition — evidence is structural, not operational

**Where**: "task-executor.md 已有 259 行，且仍只覆盖一种任务类型" / "T-test-* 任务（7 个）全部通过 mainSession: true 绕过 task-executor"

**Why it's weak**: All evidence is code-structural (line counts, field names, workaround patterns). There is zero operational evidence: no failed task runs, no debugging incidents, no time-cost data, no "we spent N hours diagnosing a bug caused by this." The proposal argues the architecture is ugly, but never shows it has actually hurt anyone. A skeptic can reasonably say "the workarounds work fine, why pay the migration cost?"

**What must improve**: Add at least one concrete operational incident — a task that misbehaved because of the hardcoded TDD flow, a debugging session that was painful because strategy was in markdown, or a new task type that couldn't be cleanly added. Quantify the cost: "adding a new task type requires touching 3 files and takes ~N hours of manual verification."

---

### Attack 2: Solution Clarity — observable behavior never shown

**Where**: "task prompt <id> → Agent(task-executor, prompt=<synthesized>)" / "模板合成采用标准字符串替换（{{TASK_ID}}、{{SCOPE}} 等占位符）"

**Why it's weak**: The proposal describes the routing logic in detail but never shows what a synthesized prompt actually looks like. A reader cannot answer: what does `task prompt T-impl-1` output? What does the slimmed-down task-executor.md contain after the change? The "user-facing behavior" — meaning what a forge operator observes when running tasks — is never described. Will run-tasks logs look different? Will task execution be faster or slower? Will error messages change? The proposal is entirely about internal plumbing with no concrete before/after observable output.

**What must improve**: Include a representative example: show the output of `task prompt <id>` for at least one task type (e.g., `implementation`), and show the before/after diff of task-executor.md. This makes the solution tangible and reviewable.

---

### Attack 3: Success Criteria — "equivalent behavior" is untestable as written

**Where**: "每种 type 的 prompt 模板迁移后，实际运行对应任务，行为与迁移前等价（运行验证）" / "task-executor 执行 test-pipeline.gen-cases 任务时，输出结构与当前 T-test-1 执行结果一致"

**Why it's weak**: Three of the twelve criteria rely on "等价" or "一致" without defining what equivalence means. Does "equivalent behavior" mean: same files created? same commit message format? same number of subagent calls? same test pass rate? "Output structure consistent with T-test-1" — what is the structure? A checklist? A file path pattern? Without a concrete definition, these criteria cannot be objectively verified and cannot be turned into automated tests. They are acceptance criteria in name only.

**What must improve**: For each "equivalent behavior" criterion, define the observable outputs that constitute equivalence. Example: "task-executor executing test-pipeline.gen-cases produces a test-cases.md file at the expected path, with ≥N test cases, and task validate passes." The criterion must be checkable by someone who was not involved in writing the original behavior.

---

### Attack 4: Alternatives Analysis — chosen approach cons are understated

**Where**: "CLI 合成 prompt + 薄执行器 ... Cons: 实现成本高（需改 CLI + run-tasks + 所有模板）"

**Why it's weak**: The cons column for the selected approach lists only implementation cost. It omits: (1) debugging difficulty increases — when a task fails, the operator must now diagnose whether the failure is in the CLI synthesis, the template, or the executor; (2) the prompt synthesis layer is a new failure mode that didn't exist before; (3) 11 templates written from scratch with no quality baseline is a significant quality risk that appears in the risk section but not in the alternatives trade-off. The alternatives table is the decision record — it should be honest about the chosen approach's downsides, not just its implementation cost.

**What must improve**: Add at least 2 more cons to the CLI approach: increased debugging surface (CLI + template + executor vs. just executor), and the cold-start quality risk for 11 new templates. Then explain why these are acceptable given the benefits.

---

### Attack 5: Scope Definition — no effort estimate or timeframe

**Where**: "阶段一：CLI 基础能力 ... 阶段二：Schema 与模板 ... 阶段三：Agent 与命令更新 ... 阶段四：清理"

**Why it's weak**: The 4-phase implementation plan has no time estimates, no team size assumption, and no definition of "done" for each phase beyond the phase 3 trigger condition. "Can a team execute this in a defined timeframe?" — the answer is no, because no timeframe is defined. The scope section lists 10 in-scope items but gives no signal about relative effort. Is phase 1 a day or a week? Is writing 11 prompt templates the bulk of the work or a minor step?

**What must improve**: Add rough effort estimates per phase (even T-shirt sizing: S/M/L). State the assumed team size. This is not about committing to a deadline — it's about demonstrating the scope is bounded and the author has thought through the work.

---

## Previous Issues Check

<!-- Only for iteration > 1 — N/A for iteration 1 -->

---

## Verdict

- **Score**: 73/100
- **Target**: 85/100
- **Gap**: 12 points
- **Action**: Continue to iteration 2
