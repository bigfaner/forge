---
date: "2026-05-11"
doc_dir: "docs/proposals/typed-task-dispatch/"
iteration: "3"
target_score: "85"
evaluator: Claude (automated, adversarial)
---

# Proposal Eval — Iteration 3

**Score: 87/100** (target: 85)

```
┌─────────────────────────────────────────────────────────────────┐
│                    PROPOSAL QUALITY SCORECARD                     │
├──────────────────────────────┬──────────┬──────────┬────────────┤
│ Dimension                    │ Score    │ Max      │ Status     │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 1. Problem Definition        │  18      │  20      │ ✅          │
│    Problem clarity           │   7/7    │          │            │
│    Evidence provided         │   6/7    │          │            │
│    Urgency justified         │   5/6    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 2. Solution Clarity          │  18      │  20      │ ✅          │
│    Approach concrete         │   7/7    │          │            │
│    User-facing behavior      │   6/7    │          │            │
│    Differentiated            │   5/6    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 3. Alternatives Analysis     │  13      │  15      │ ✅          │
│    Alternatives listed (≥2)  │   5/5    │          │            │
│    Pros/cons honest          │   4/5    │          │            │
│    Rationale justified       │   4/5    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 4. Scope Definition          │  14      │  15      │ ✅          │
│    In-scope concrete         │   5/5    │          │            │
│    Out-of-scope explicit     │   5/5    │          │            │
│    Scope bounded             │   4/5    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 5. Risk Assessment           │  14      │  15      │ ✅          │
│    Risks identified (≥3)     │   5/5    │          │            │
│    Likelihood + impact rated │   5/5    │          │            │
│    Mitigations actionable    │   4/5    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 6. Success Criteria          │  12      │  15      │ ⚠️          │
│    Measurable                │   4/5    │          │            │
│    Coverage complete         │   4/5    │          │            │
│    Testable                  │   4/5    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ TOTAL                        │  87      │  100     │            │
└──────────────────────────────┴──────────┴──────────┴────────────┘
```

---

## Deductions

| Location | Issue | Penalty |
|----------|-------|---------|
| Success Criteria #6(c) | "任务产物（代码文件、文档文件、测试文件）与迁移前同类任务产物的文件路径模式一致" — "文件路径模式一致" 无量化定义，三次迭代未修复，不可自动验证 | -2 pts |

---

## Attack Points

### Attack 1: Solution Clarity — gate type behavior never shown, yet it contains the most complex routing decision

**Where**: "gate | — | 阶段出口验证，直接进入 quality gate。自修复还是追加 fix 任务的判断标准由 prompt 模板自描述"

**Why it's weak**: The proposal adds a `fix` example in this iteration, which covers a previously missing case. But `gate` is the type with the most consequential branching logic in the entire system: it decides whether to self-fix inline or append a new `fix` task to the index. This decision has downstream effects on task count, phase progression, and record integrity. The proposal defers the decision to "prompt 模板自描述" — meaning the reader cannot evaluate whether the logic is sound without reading a template that does not yet exist. The two observable behavior examples both show linear flows (TDD steps, fix steps). Neither demonstrates a conditional branch. A reader reviewing this proposal cannot verify that the two-layer separation actually handles the gate decision correctly.

**What must improve**: Add an observable behavior example for `gate` type showing the conditional branch: what the synthesized prompt looks like when the gate passes, and what it looks like (or what action it triggers) when it fails and must decide between self-fix and appending a fix task. This is the hardest case — show it.

---

### Attack 2: Success Criteria — criterion 6(c) is still untestable after three iterations

**Where**: "任务产物（代码文件、文档文件、测试文件）与迁移前同类任务产物的文件路径模式一致"

**Why it's weak**: This criterion has appeared in every iteration and has never been made concrete. "文件路径模式一致" is not a test condition — it is a description of intent. What constitutes a "path pattern"? Is it a directory prefix check (`docs/features/<slug>/`)? A glob pattern per type? A regex? Without a definition, two engineers could disagree on whether this criterion passes. The other criteria in this section are specific: exit codes, file existence at named paths, `task validate` pass/fail. This one is not.

**What must improve**: Replace "文件路径模式一致" with a concrete, type-specific check. For example: "implementation 类型任务产物路径符合 `src/<scope>/*.go` 模式；doc-generation.summary 产物路径为 `docs/features/<slug>/tasks/records/<phase>-summary.md`；test-pipeline.gen-cases 产物路径为 `docs/features/<slug>/test-cases.md`." One line per type is sufficient. The current wording is a placeholder dressed as a criterion.

---

### Attack 3: Alternatives Analysis — rejected approach verdicts are conclusions, not reasoning

**Where**: "Do nothing | ... | Rejected：不可持续" / "Agent 内部策略分发 | ... | Rejected：短期方案" / "多 agent | ... | Rejected：维护成本高"

**Why it's weak**: The chosen approach verdict is now well-reasoned — it names the cons, explains the mitigations, and justifies the selection. The three rejected approaches get one-line dismissals. "不可持续" and "短期方案" are labels, not arguments. For "多 agent", the cons column says "run-tasks 需要维护 dispatch 表；共享逻辑重复" but the verdict does not explain why this cost exceeds the chosen approach's implementation cost (CLI + 11 templates + all task templates). A reader who prefers the multi-agent approach has no counter-argument to engage with — only a label. The alternatives table is the decision record; thin verdicts on rejected options undermine its credibility.

**What must improve**: For each rejected approach, add one sentence of reasoning to the verdict: why the specific cons outweigh the pros in this context. For "多 agent": explain that shared execution constraints (ONE TASK, record-task mandatory) would need to be duplicated across all agents, making the maintenance cost scale with type count rather than being fixed. For "Agent 内部策略分发": explain why the testability gap is disqualifying given the 2h/incident debugging cost already documented in the evidence section.

---

## Previous Issues Check

| Previous Attack | Addressed? | Evidence |
|----------------|------------|----------|
| Attack 1: Alternatives cons understated (iter 2) | ✅ | Cons column now includes "调试面扩大——任务失败时需逐层排查 CLI 合成、模板内容、执行器约束三处" and "11 个 prompt 模板从零编写，初始质量无保障（与 Risk 表中"高/高"评级一致）"; verdict explains mitigations |
| Attack 2: Scope no effort signal (iter 2) | ✅ | T-shirt sizing added (L/M/S/S per phase), team size stated ("1 名开发者"), per-phase completion criteria added for all four phases |
| Attack 3: Observable behavior covers only 1 type (iter 2) | ✅ | `task prompt T-fix-1` example added with full 7-step fix flow and explicit note on two-layer separation |
| Attack 4: Risk likelihood inconsistency (iter 2) | ✅ | First risk row now reads "中（迁移期）/ 低（稳定后）" with explicit reconciliation against the 高/高 template quality risk |

---

## Verdict

- **Score**: 87/100
- **Target**: 85/100
- **Gap**: 0 points (target reached)
- **Action**: Target reached. Remaining weaknesses (gate example, criterion 6c specificity, rejected-approach reasoning) are refinements, not blockers. Proceed to `/write-prd`.
