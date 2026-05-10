---
date: "2026-05-10"
doc_dir: "docs/features/task-executor-skeleton/prd/"
iteration: 1
target_score: 90
scoring_mode: "Mode B"
evaluator: Claude (automated, adversarial)
---

# PRD Eval — Iteration 1

**Score: 88/100** (target: 90, mode: Mode B)

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
│ 2. Flow Diagrams             │  18      │  20      │ ⚠️         │
│    Mermaid diagram exists    │  7/7     │          │            │
│    Main path complete        │  7/7     │          │            │
│    Decision + error branches │  4/6     │          │            │
├──────────────────────────────┼──────────┬──────────┬────────────┤
│ 3b. Flow Completeness (B)    │  17      │  20      │ ⚠️         │
│    Flow steps complete       │  6/7     │          │            │
│    Data flow documented      │  7/7     │          │ (auto N/A) │
│    Exception + edge cases    │  4/6     │          │            │
├──────────────────────────────┼──────────┬──────────┬────────────┤
│ 4. User Stories              │  27      │  30      │ ⚠️         │
│    Coverage per user type    │  7/7     │          │            │
│    Format correct            │  7/7     │          │            │
│    AC per story (G/W/T)      │  6/6     │          │            │
│    AC verifiability          │  7/10    │          │            │
├──────────────────────────────┼──────────┬──────────┬────────────┤
│ 5. Scope Clarity             │  13      │  15      │ ⚠️         │
│    In-scope concrete         │  5/5     │          │            │
│    Out-of-scope explicit     │  4/4     │          │
│    Consistent with specs     │  4/6     │          │            │
├──────────────────────────────┼──────────┬──────────┬────────────┤
│ TOTAL                        │  88      │  100     │ ⚠️         │
└──────────────────────────────┴──────────┴──────────┴────────────┘
```

> **Mode B** (prd-ui-functions.md absent): Dimension 3 evaluates Flow Completeness from prd-spec.md Flow Description.
> Sub-criteria: Flow steps describe complete business process /7, Data flow documented /7, Exception handling and edge cases /6.

---

## Deductions

| Location | Issue | Penalty |
|----------|-------|---------|
| prd-spec.md:29 (Goals table) | Goal "向后兼容" has no quantified metric — the other 3 goals have numeric targets but this one is purely qualitative ("fallback 保障" is not measurable) | -1 pt |
| prd-spec.md:29 (Goals table) | Background says "所有非 MAIN_SESSION 任务" are affected, but only T-test-3 has a quantified improvement; the other 15 templates have no metric | -1 pt |
| prd-spec.md:67-78 (Flow diagram) | No execution failure branch — diagram assumes "Agent 执行 workflow 步骤" always succeeds and proceeds to Step 3 | -1 pt |
| prd-spec.md:67-78 (Flow diagram) | "Fallback" and "Warn" both merge into same "Execute" node, losing the distinction between TDD execution and Workflow execution | -1 pt |
| prd-spec.md:62 (Flow description) | No documented handling for: task file unreadable, workflow execution partial failure, agent timeout during execution | -1 pt |
| prd-spec.md:62 (Flow description) | Missing state transitions: what status does the task end in after each path? No mention of success/failure states | -1 pt |
| prd-user-stories.md:16 (Story 1 AC2) | "行为与当前一致" is subjective — requires knowing current behavior; should reference specific existing outputs or behaviors | -0.5 pt |
| prd-user-stories.md:44 (Story 3 AC2) | "审查条件分支无残留逻辑" requires human judgment rather than objective verification; what constitutes "residual logic"? | -0.5 pt |
| prd-user-stories.md (all stories) | No AC covers execution failure during workflow — what happens if the agent fails mid-execution? No error-path AC exists | -2 pt |
| prd-spec.md:44-48 (In Scope) vs prd-user-stories.md | Commands (run-tasks.md, execute-task.md) and 3 skill docs are in scope but have no corresponding user stories | -2 pt |

---

## Attack Points

### Attack 1: Flow Diagrams — no execution failure path

**Where**: prd-spec.md lines 67-78 — the Mermaid diagram shows `Execute[Agent 执行 workflow 步骤] --> S3[Step 3: 记录结果 + 提交]` with no failure branch.

**Why it's weak**: The diagram assumes workflow execution always succeeds. In reality, an agent following a custom Execution Workflow could fail, timeout, or produce invalid output. The flow diagram must show what happens when execution fails — does it retry? Does it record a failure? Does it create a fix task (as Story 2 AC1 suggests for T-test-3)? Without this branch, the diagram is incomplete.

**What must improve**: Add a failure branch after the Execute node: `Execute -->|Failure| FailHandler[记录失败 + 创建 fix task 或标记失败]`. Also consider splitting Execute into two nodes (TDD execution vs Workflow execution) to show they are genuinely different paths, not just the same node.

### Attack 2: User Stories — missing ACs for execution failure and scope gap

**Where**: prd-user-stories.md — no story covers what happens when workflow execution fails. Also, no story covers the command/skill doc changes listed in prd-spec.md In Scope items 5-6 and 9.

**Why it's weak**: Story 2 AC1 mentions "创建 fix task 并停止" but only in the specific context of T-test-3's workflow. There is no general-purpose AC for: "Given an Execution Workflow that fails during execution, Then the agent handles the failure gracefully." Additionally, In Scope lists changes to `commands/run-tasks.md`, `commands/execute-task.md`, and 3 skill docs, but no user story covers these components. The Forge Maintainer story only covers noTest removal, not the command/skill changes. This is a cross-section inconsistency between Scope and User Stories.

**What must improve**: (1) Add at least one AC covering general execution failure handling (agent encounters error during workflow execution). (2) Either add a user story for the command/skill doc changes, or expand Story 3 to cover "all references to NO_TEST are removed from commands and skill docs" with corresponding ACs.

### Attack 3: Flow Completeness — undocumented failure states and edge cases

**Where**: prd-spec.md line 62 — "执行完成后统一进入 Step 3（记录 + 提交）" assumes all executions complete successfully.

**Why it's weak**: The flow description text, like the diagram, has no failure handling. Three edge cases are missing: (1) What happens if the task file cannot be read? (2) What happens if the agent fails partway through executing a custom workflow? (3) What happens if the workflow contains instructions that produce an ambiguous or invalid result? The Quality Requirements mention "每个 Execution Workflow 模板须包含明确的结束条件" but this is a template-author obligation, not a system guarantee. The system itself needs documented failure handling.

**What must improve**: Add explicit handling for: unreadable task files (error + log + skip?), workflow execution failure (retry? fail-fast? create fix task?), and partial completion states. Document the terminal states (success / failure / partial) and what triggers each transition.

---

## Previous Issues Check

<!-- Only for iteration > 1 -->

| Previous Attack | Addressed? | Evidence |
|----------------|------------|----------|
| N/A (iteration 1) | — | — |

---

## Verdict

- **Score**: 88/100
- **Target**: 90/100
- **Gap**: 2 points
- **Action**: Continue to iteration 2 — gap is narrow; focus on adding execution failure handling to flow diagram + description, and closing the scope/story coverage gap for commands and skill docs.
