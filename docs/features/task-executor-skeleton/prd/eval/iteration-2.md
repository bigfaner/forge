---
date: "2026-05-10"
doc_dir: "docs/features/task-executor-skeleton/prd/"
iteration: 2
target_score: 90
scoring_mode: "Mode B"
evaluator: Claude (automated, adversarial)
---

# PRD Eval — Iteration 2

**Score: 91/100** (target: 90, mode: Mode B)

```
┌─────────────────────────────────────────────────────────────────┐
│                       PRD QUALITY SCORECARD                      │
├──────────────────────────────┬──────────┬──────────┬────────────┤
│ Dimension                    │ Score    │ Max      │ Status     │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 1. Background & Goals        │  13      │  15      │ ⚠️         │
│    Three elements            │  5/5     │          │            │
│    Goals quantified          │  3/4     │          │            │
│    Logical consistency       │  5/6     │          │            │
├──────────────────────────────┼──────────┬──────────┬────────────┤
│ 2. Flow Diagrams             │  19      │  20      │ ✅         │
│    Mermaid diagram exists    │  7/7     │          │            │
│    Main path complete        │  7/7     │          │            │
│    Decision + error branches │  5/6     │          │            │
├──────────────────────────────┼──────────┬──────────┬────────────┤
│ 3b. Flow Completeness (B)    │  18      │  20      │ ✅         │
│    Flow steps complete       │  6/7     │          │            │
│    Data flow documented      │  7/7     │          │ (auto N/A) │
│    Exception + edge cases    │  5/6     │          │            │
├──────────────────────────────┼──────────┬──────────┬────────────┤
│ 4. User Stories              │  27      │  30      │ ⚠️         │
│    Coverage per user type    │  7/7     │          │            │
│    Format correct            │  7/7     │          │            │
│    AC per story (G/W/T)      │  6/6     │          │            │
│    AC verifiability          │  7/10    │          │            │
├──────────────────────────────┼──────────┬──────────┬────────────┤
│ 5. Scope Clarity             │  14      │  15      │ ✅         │
│    In-scope concrete         │  5/5     │          │            │
│    Out-of-scope explicit     │  4/4     │          │            │
│    Consistent with specs     │  5/6     │          │            │
├──────────────────────────────┼──────────┬──────────┬────────────┤
│ TOTAL                        │  91      │  100     │ ✅         │
└──────────────────────────────┴──────────┴──────────┴────────────┘
```

> **Mode B** (prd-ui-functions.md absent): Dimension 3 evaluates Flow Completeness from prd-spec.md Flow Description.
> Sub-criteria: Flow steps describe complete business process /7, Data flow documented /7, Exception handling and edge cases /6.

---

## Deductions

| Location | Issue | Penalty |
|----------|-------|---------|
| prd-spec.md:34 (Goals table) | Goal "向后兼容" has no quantified metric — "fallback 保障" is not measurable; no test criterion for verifying old templates still work with TDD fallback | -1 pt |
| prd-spec.md:76-95 (Flow diagram) | `Fallback` and `Warn` both merge into same `Execute` node, losing the distinction between TDD execution path and Workflow execution path | -1 pt |
| prd-spec.md:62 (Flow description) | Injection mechanism underspecified — "将其注入到 agent prompt 的 Step 2 指令区域" does not describe how (string replacement? paragraph extraction? template variable?) | -1 pt |
| prd-spec.md:66-70 (Flow description) | No handling for workflow containing invalid/unexecutable instructions (not empty, but nonsensical) — Quality Requirements offload this to template authors with no system-level validation | -1 pt |
| prd-user-stories.md:16 (Story 1 AC2) | "行为与当前一致" is subjective — carried over from iteration 1 without fix; should reference specific outputs or steps | -0.5 pt |
| prd-user-stories.md:44 (Story 3 AC4) | "无基于 noTest 的残留逻辑" requires human judgment; grep coverage is already in AC1, but "残留逻辑" has no objective definition | -0.5 pt |
| prd-user-stories.md:64 (Story 4 AC4) | "已完成步骤摘要和失败点描述" is not objectively verifiable — "摘要" and "描述" have no structural requirements | -0.5 pt |
| prd-user-stories.md (all stories) | No AC verifies backward compatibility: Goal 4 promises "无 Execution Workflow 的旧任务回退到 TDD，行为不变" but no story tests this claim | -0.5 pt |
| prd-spec.md:40 vs prd-user-stories.md | In-scope item #1 (task-executor.md Step 2-3 merge) has no story directly testing the merged step behavior — Story 3 AC7 checks NO_TEST removal only, not the new merged execution logic | -1 pt |

---

## Attack Points

### Attack 1: User Stories — backward compatibility is claimed but never tested

**Where**: prd-spec.md line 34 — "向后兼容 | 无 Execution Workflow 的旧任务回退到 TDD，行为不变 | fallback 保障" and prd-user-stories.md — no story covers this.

**Why it's weak**: The PRD makes an explicit backward-compatibility promise as one of its four goals, yet no user story has an AC that verifies old templates (without `## Execution Workflow`) actually fall back to TDD correctly. Story 1 AC2 says "行为与当前一致" which is both subjective and only covers the positive case (workflow absence triggers TDD), not the full regression (old templates produce identical results to current behavior). This is a goal-without-verification gap.

**What must improve**: Add an AC to Story 1 (or create a new story): "Given a task template without `## Execution Workflow`, When task-executor runs the task, Then Step 2 output matches current TDD cycle output (RED/GREEN/REFACTOR steps present, Quality Gate triggered)". Make the expected behavior structurally specific, not "行为与当前一致".

### Attack 2: User Stories — subjective ACs carried over from iteration 1

**Where**: prd-user-stories.md line 16 — "行为与当前一致", line 44 — "无基于 noTest 的残留逻辑", line 64 — "已完成步骤摘要和失败点描述".

**Why it's weak**: Three ACs use language that cannot be objectively verified. (1) "行为与当前一致" — the tester must know and compare against undocumented current behavior. (2) "残留逻辑" — grep is already covered by AC1; this AC adds a manual code review requirement with no pass/fail criterion. (3) "摘要" and "描述" have no structural definition — one developer's "摘要" is another's "insufficient detail". These survived both iterations unaddressed.

**What must improve**: (1) Replace "行为与当前一致" with specific expected outputs: e.g., "Then Step 2 contains 'TDD implementation cycle' instructions and Step 3 triggers Quality Gate". (2) Replace "无残留逻辑" with a grep-able criterion: "Then `grep -r 'noTest\|NO_TEST\|no.test\|--no-test' --include='*.go' --include='*.md'` returns zero matches in conditional branches" — or remove it if already covered by AC1. (3) Replace "摘要和描述" with structural requirements: "Then the failure record contains a list of completed step names and a `failure_reason` field".

### Attack 3: Flow Completeness — injection mechanism is a black box

**Where**: prd-spec.md line 62 — "将其注入到 agent prompt 的 Step 2 指令区域，替换硬编码的 TDD 步骤".

**Why it's weak**: "注入" (inject) is the core mechanism of the entire feature — it is how workflow text gets from the task file into the agent's execution context. But the PRD treats it as a single verb with no detail. Is it string concatenation? Regex-based section replacement? Prompt template variable substitution? Does it strip the `## Execution Workflow` header? Does it preserve markdown formatting? Implementation depends entirely on answering these questions, but the PRD is silent. This is the one place where "describe WHAT not HOW" breaks down — the injection boundary is a WHAT question (what gets injected, in what format, with what boundaries), not a HOW question.

**What must improve**: Add a sentence or two clarifying the injection: e.g., "The raw text under `## Execution Workflow` (excluding the heading itself) replaces the entire Step 2 instruction block in the agent prompt. Markdown formatting is preserved. If the task template has no `## Execution Workflow` heading, the existing hardcoded TDD+Quality Gate Step 2 block is used unchanged."

---

## Previous Issues Check

| Previous Attack | Addressed? | Evidence |
|----------------|------------|----------|
| Attack 1 (Flow Diagrams — no execution failure path) | ✅ Yes | Mermaid diagram now has `ExecOK` decision node with `成功` and `失败/超时` branches, plus `FileOK` decision for file-read failure. Three diamond nodes and two error terminal states added. |
| Attack 2 (User Stories — missing ACs for execution failure and scope gap) | ✅ Yes | New Story 4 ("Task-executor Agent 处理执行失败") added with 4 ACs covering file failure, workflow failure with/without explicit directives, and partial completion. Story 3 expanded from 3 to 6 ACs now covering commands, skill docs, and task-executor.md explicitly. |
| Attack 3 (Flow Completeness — undocumented failure states) | ✅ Yes | Flow description now documents 4 failure paths (file unreadable, timeout, workflow execution failure, partial completion) and 2 terminal states (`completed`, `failed`). |

All three iteration-1 attacks were genuinely addressed. The remaining deductions are either carry-overs that were not fully fixed (subjective AC language) or new observations at a finer level of scrutiny.

---

## Verdict

- **Score**: 91/100
- **Target**: 90/100
- **Gap**: 0 points (target reached, +1 above)
- **Action**: Target reached. Remaining issues (subjective ACs, injection mechanism clarity) are minor and do not block implementation. Optional iteration 3 could focus on tightening AC verifiability for marginal improvement.
