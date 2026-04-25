---
date: "2026-04-24"
doc_dir: "docs/proposals/integrated-test-lifecycle/"
iteration: "2"
target_score: "N/A"
evaluator: Claude (automated, adversarial)
---

# Proposal Eval — Iteration 2

**Score: 81/100** (target: N/A)

```
┌─────────────────────────────────────────────────────────────────┐
│                    PROPOSAL QUALITY SCORECARD                     │
├──────────────────────────────┬──────────┬──────────┬────────────┤
│ Dimension                    │ Score    │ Max      │ Status     │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 1. Problem Definition        │  16      │  20      │ ⚠️          │
│    Problem clarity           │   6/7    │          │            │
│    Evidence provided         │   6/7    │          │            │
│    Urgency justified         │   4/6    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 2. Solution Clarity          │  17      │  20      │ ✅          │
│    Approach concrete         │   6/7    │          │            │
│    User-facing behavior      │   6/7    │          │            │
│    Differentiated            │   5/6    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 3. Alternatives Analysis     │  13      │  15      │ ⚠️          │
│    Alternatives listed (≥2)  │   5/5    │          │            │
│    Pros/cons honest          │   4/5    │          │            │
│    Rationale justified       │   4/5    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 4. Scope Definition          │  14      │  15      │ ✅          │
│    In-scope concrete         │   5/5    │          │            │
│    Out-of-scope explicit     │   5/5    │          │            │
│    Scope bounded             │   4/5    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 5. Risk Assessment           │  13      │  15      │ ⚠️          │
│    Risks identified (≥3)     │   5/5    │          │            │
│    Likelihood + impact rated │   4/5    │          │            │
│    Mitigations actionable    │   4/5    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 6. Success Criteria          │  13      │  15      │ ⚠️          │
│    Measurable                │   5/5    │          │            │
│    Coverage complete         │   4/5    │          │            │
│    Testable                  │   4/5    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ Deductions                   │  -5      │          │            │
│    Scope/criteria gap (docs) │  -3      │          │            │
│    Solution/risk contradiction│ -2      │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ TOTAL                        │  81      │  100     │            │
└──────────────────────────────┴──────────┴──────────┴────────────┘
```

---

## Deductions

| Location | Issue | Penalty |
|----------|-------|---------|
| Scope § In Scope vs. Success Criteria | In-scope lists `OVERVIEW.md`, `test-capture-design.md`, and two skill Prerequisites updates — none appear in success criteria | -3 pts |
| Solution § 核心机制 1 vs. Risk § Risk 2 mitigation | Solution says "自动追加两个固定任务" (always append); Risk 2 mitigation says "追加时询问用户是否需要 e2e 测试，或检测 PRD 中是否有 e2e 相关需求" — two contradictory behaviors, neither chosen | -2 pts |

---

## Attack Points

### Attack 1: Solution Clarity — test ID is the linchpin of the graduation model but is never defined

**Where**: "相同 test ID 更新，新 ID 追加——测试自然聚合，不重复" and success criterion 4: "两个不同 feature 覆盖同一页面时，`tests/e2e/ui/<page>/` 中包含两者的测试，无重复"

**Why it's weak**: The entire deduplication guarantee rests on test IDs being stable and consistent across features. The proposal never defines what a test ID is, how it is generated, or who controls it. If `gen-test-cases` auto-generates IDs (e.g., from a hash or sequential counter), two features covering the same login flow will produce different IDs for semantically identical tests — the graduation model will append duplicates instead of deduplicating. If IDs are user-defined, the proposal needs a naming convention and enforcement mechanism. "无重复" in success criterion 4 is untestable without this definition: does it mean no duplicate IDs, no duplicate test logic, or no duplicate file paths?

**What must improve**: Define the test ID schema explicitly (e.g., `<type>/<target>/<slug>` where slug is derived from the test case title). Specify whether IDs are generated or authored. Add a success criterion that verifies ID stability across two features covering the same target.

---

### Attack 2: Risk Assessment — Risk 2 mitigation is an unresolved design fork that contradicts the solution

**Where**: Risk 2 mitigation: "追加时询问用户是否需要 e2e 测试，或检测 PRD 中是否有 e2e 相关需求"; Solution § 核心机制 1: "在所有业务任务之后，自动追加两个固定任务"

**Why it's weak**: The solution commits to unconditional appending ("固定任务" — fixed tasks). The risk mitigation then proposes two conditional alternatives without choosing either. "询问用户" (interactive prompt) and "检测 PRD" (automatic detection) are fundamentally different UX models — one requires human input mid-automation, the other requires the skill to parse PRD content. Neither is reconciled with the "always append" default in the solution. A reader cannot tell what the actual behavior will be for a feature with no e2e requirements.

**What must improve**: Make a decision: either always append (and accept that agents will `skip` the tasks for non-e2e features, which should be stated explicitly) or define the detection logic (e.g., "append only if PRD contains a `## UI` or `## API` section"). Remove the "or" and commit to one approach.

---

### Attack 3: Problem Definition — urgency names the risk but not the consequence

**Where**: "当前 7 个 feature 的跳过率已达 100%；每新增一个 feature 就多一次无测试门控的合并风险。"

**Why it's weak**: The urgency argument identifies the exposure (no test gate) but never states what actually happens when the gate is bypassed. The `feat-auto-test-after-run-all-tasks` incident is cited as evidence that the gate was bypassed — but the proposal doesn't say whether this caused any observable problem (a regression, a broken build, a silent failure). "合并风险" is a category of harm, not a harm. A skeptical reader can reasonably ask: "If 7 features merged without e2e tests and nothing broke, why is this urgent?"

**What must improve**: State the actual or hypothetical consequence of a missed test gate. Even a hypothetical is better than none: e.g., "If `feat-auto-test-after-run-all-tasks` had introduced a regression in the Stop hook, it would have shipped undetected — the very feature that was supposed to enforce testing bypassed its own gate." This makes the urgency concrete rather than categorical.

---

## Previous Issues Check

| Previous Attack | Addressed? | Evidence |
|----------------|------------|----------|
| Attack 1: Evidence sample too small (only 2 directories, skip rate unsupported) | ✅ | Now states "项目现有 7 个 feature 目录，无一包含 `testing/` 子目录，测试生成步骤跳过率为 100%（7/7）" with specific incident about `feat-auto-test-after-run-all-tasks` |
| Attack 2: 方案 B dismissed with straw-man ("收益不明显") | ✅ | Now explains overlap concretely ("两者都需要检查所有任务已完成、都需要运行 e2e 测试") and replaces vague dismissal with specific comparison ("方案 B 引入 1 个新命令 + 1 套调用约定，而方案 A 在现有原语内解决同一问题") |
| Attack 3: No effort bound, scope open-ended | ✅ | Added Phase 1 (~3 engineer-days) and Phase 2 (~3 engineer-days) with explicit phasing rationale and dependency between phases |

---

## Verdict

- **Score**: 81/100
- **Target**: N/A
- **Gap**: N/A
- **Action**: All three previous attacks addressed. Three new attacks identified. Address test ID definition (Attack 1), the solution/risk contradiction on conditional appending (Attack 2), and urgency concreteness (Attack 3) to reach 88+.
