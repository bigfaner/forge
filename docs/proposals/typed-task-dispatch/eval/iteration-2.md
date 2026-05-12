---
date: "2026-05-11"
doc_dir: "docs/proposals/typed-task-dispatch/"
iteration: "2"
target_score: "85"
evaluator: Claude (automated, adversarial)
---

# Proposal Eval — Iteration 2

**Score: 75/100** (target: 85)

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
│    User-facing behavior      │   6/7    │          │            │
│    Differentiated            │   4/6    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 3. Alternatives Analysis     │  11      │  15      │ ⚠️          │
│    Alternatives listed (≥2)  │   5/5    │          │            │
│    Pros/cons honest          │   3/5    │          │            │
│    Rationale justified       │   3/5    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 4. Scope Definition          │  11      │  15      │ ⚠️          │
│    In-scope concrete         │   5/5    │          │            │
│    Out-of-scope explicit     │   4/5    │          │            │
│    Scope bounded             │   2/5    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 5. Risk Assessment           │  13      │  15      │ ✅          │
│    Risks identified (≥3)     │   5/5    │          │            │
│    Likelihood + impact rated │   4/5    │          │            │
│    Mitigations actionable    │   4/5    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 6. Success Criteria          │  12      │  15      │ ⚠️          │
│    Measurable                │   4/5    │          │            │
│    Coverage complete         │   4/5    │          │            │
│    Testable                  │   4/5    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ TOTAL                        │  75      │  100     │            │
└──────────────────────────────┴──────────┴──────────┴────────────┘
```

---

## Deductions

| Location | Issue | Penalty |
|----------|-------|---------|
| Success Criteria #6 | "任务产物（代码文件、文档文件、测试文件）与迁移前同类任务产物的文件路径模式一致" — "文件路径模式一致" 无量化定义，不可自动验证 | -2 pts |
| Alternatives vs Risk | Cross-section inconsistency: Risk section rates "11 个 prompt 模板从零编写，初始质量无保障" as 高/高，但 Alternatives 表中 CLI 方案的 Cons 列完全未提及此风险 | -3 pts |

---

## Attack Points

### Attack 1: Alternatives Analysis — chosen approach cons are still a single line

**Where**: "CLI 合成 prompt + 薄执行器 ... Cons: 实现成本高（需改 CLI + run-tasks + 所有模板）"

**Why it's weak**: This is identical to iteration 1. The cons column lists only implementation cost. Two material downsides are absent: (1) the debugging surface expands — when a task fails post-migration, the operator must now triage CLI synthesis, template content, and executor separately, whereas today there is only one place to look; (2) the cold-start quality risk for 11 templates written from scratch is rated 高/高 in the risk section but does not appear here at all. The alternatives table is the decision record. Omitting known downsides from it makes the trade-off analysis dishonest and undermines the rationale for the chosen approach.

**What must improve**: Add at least two more cons to the CLI approach row: increased debugging surface (CLI + template + executor vs. just executor), and cold-start quality risk for 11 new templates. Then explain in the Verdict column why these are acceptable given the benefits. The risk section already has the language — copy it here.

---

### Attack 2: Scope Definition — no effort signal after two iterations

**Where**: "阶段一：CLI 基础能力 ... 阶段二：Schema 与模板 ... 阶段三：Agent 与命令更新 ... 阶段四：清理"

**Why it's weak**: Identical to iteration 1. The 4-phase plan lists 12 numbered steps with zero time estimates, no T-shirt sizing, and no team size assumption. Phase 3 has a trigger condition ("所有 11 种类型的任务实际运行行为等价"), but phases 1, 2, and 4 have no completion criteria beyond "done." Writing 11 prompt templates from scratch is listed as a single bullet in phase 2 alongside schema changes — there is no signal that this is the bulk of the work. A reader cannot answer: is this a one-week effort or a one-month effort?

**What must improve**: Add rough effort estimates per phase (T-shirt sizing is sufficient: S/M/L/XL). State the assumed team size. Define a completion criterion for each phase, not just phase 3. This is not about committing to a deadline — it is about demonstrating the scope is bounded.

---

### Attack 3: Solution Clarity — observable behavior example covers 1 of 11 types

**Where**: "Observable Behavior Example — `task prompt T-impl-1` 的实际输出（implementation 类型）"

**Why it's weak**: The example is a genuine improvement over iteration 1, but it covers only `implementation`, the simplest and most TDD-like type. The types with meaningfully different flows — `fix` (diagnose → locate → fix → verify → commit), `gate` (quality gate with self-fix vs. append-fix decision), `test-pipeline.eval-cases` (the permanent mainSession exception) — have no example output. A reader reviewing the proposal cannot verify that the two-layer separation actually works for the edge cases. The example chosen is the easy case.

**What must improve**: Add one example for a non-implementation type — `fix` or `gate` would be most revealing, since these are the types that previously required workarounds. Show what the synthesized prompt looks like and how the executor constraints layer interacts with it.

---

### Attack 4: Risk Assessment — internal likelihood inconsistency

**Where**: "task prompt 命令失败（type 未知、模板错误）导致 dispatch 链路中断 | 低 | 高" vs. "11 个 prompt 模板从零编写，初始质量无保障 | 高 | 高"

**Why it's weak**: These two risks are causally linked — template quality failures are the primary mechanism by which `task prompt` fails. Yet the first is rated 低/高 and the second 高/高. If template quality is a 高 likelihood risk, then `task prompt` failure is also a 高 likelihood risk during the migration period. Rating them differently without explanation suggests the likelihood ratings were assigned independently rather than as a coherent model.

**What must improve**: Either reconcile the two ratings with an explanation (e.g., "task prompt failure is low likelihood post-stabilization, but high during initial rollout"), or merge them into a single risk entry that covers both the mechanism and the consequence. The current table implies the command is unlikely to fail even though the templates it depends on are likely to have quality issues.

---

## Previous Issues Check

| Previous Attack | Addressed? | Evidence |
|----------------|------------|----------|
| Attack 1: Problem Definition — evidence is structural, not operational | ✅ | Added 运营事故（2026-04）with specific time costs (2h verification + 1h debugging) and "新增任务类型的固定成本" (~2h/次) |
| Attack 2: Solution Clarity — observable behavior never shown | ✅ | Added Observable Behavior Example with actual `task prompt T-impl-1` output and before/after diff of task-executor.md |
| Attack 3: Success Criteria — "equivalent behavior" untestable | ✅ Partial | Criteria 4, 5, 6 now have concrete definitions (exit codes, file paths, ≥1 test case). Criterion 6(c) "文件路径模式一致" remains vague. |
| Attack 4: Alternatives Analysis — chosen approach cons understated | ❌ | Cons column unchanged: still only "实现成本高（需改 CLI + run-tasks + 所有模板）" |
| Attack 5: Scope Definition — no effort estimate or timeframe | ❌ | 4-phase plan unchanged, no T-shirt sizing, no team size, no per-phase completion criteria |

---

## Verdict

- **Score**: 75/100
- **Target**: 85/100
- **Gap**: 10 points
- **Action**: Continue to iteration 3
