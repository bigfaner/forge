---
date: "2026-05-11"
doc_dir: "docs/features/typed-task-dispatch/prd/"
iteration: "1"
target_score: "90"
scoring_mode: "Mode B"
evaluator: Claude (automated, adversarial)
---

# PRD Eval — Iteration 1

**Score: 85/100** (target: 90, mode: Mode B)

```
┌─────────────────────────────────────────────────────────────────┐
│                       PRD QUALITY SCORECARD                      │
├──────────────────────────────┬──────────┬──────────┬────────────┤
│ Dimension                    │ Score    │ Max      │ Status     │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 1. Background & Goals        │  14      │  15      │ ✅          │
│    Three elements            │  5/5     │          │            │
│    Goals quantified          │  4/4     │          │            │
│    Logical consistency       │  5/6     │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 2. Flow Diagrams             │  18      │  20      │ ✅          │
│    Mermaid diagram exists    │  7/7     │          │            │
│    Main path complete        │  7/7     │          │            │
│    Decision + error branches │  4/6     │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 3b. Flow Completeness (B)    │  13      │  20      │ ⚠️          │
│    Flow steps complete       │  5/7     │          │            │
│    Data flow documented      │  4/7     │          │            │
│    Exception handling        │  4/6     │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 4. User Stories              │  27      │  30      │ ✅          │
│    Coverage per user type    │  7/7     │          │            │
│    Format correct            │  7/7     │          │            │
│    AC per story (G/W/T)      │  6/6     │          │            │
│    AC verifiability          │  7/10    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 5. Scope Clarity             │  13      │  15      │ ⚠️          │
│    In-scope concrete         │  5/5     │          │            │
│    Out-of-scope explicit     │  4/4     │          │            │
│    Consistent with specs     │  4/6     │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ TOTAL                        │  85      │  100     │            │
└──────────────────────────────┴──────────┴──────────┴────────────┘
```

> **Mode B** (prd-ui-functions.md absent): Dimension 3 evaluates Flow Completeness from prd-spec.md Flow Description.
> Sub-criteria: Flow steps describe complete business process /7, Data flow documented /7, Exception handling and edge cases /6.

---

## Deductions

| Location | Issue | Penalty |
|----------|-------|---------|
| prd-spec.md — Background/What | "What (Target)" describes the technical solution, not the business outcome — conflates solution with target | -1 pt (Dim 1 consistency) |
| prd-spec.md — Flow Description | No data flow table for multi-component system (run-tasks ↔ task-cli ↔ task-executor ↔ index.json); only narrative mentions file reads | -3 pts (Dim 3 data flow) |
| prd-spec.md — Flow Description | Task state transitions after `MarkBlocked` not described: what triggers unblocking? Is it recoverable? | -2 pts (Dim 3 flow steps) |
| prd-spec.md — Flow Diagram | `--fix-record-missed` failure path absent from diagram; eval-cases path has no error branch | -2 pts (Dim 2 error branches) |
| prd-spec.md — Scope | `execute-task.md 路由变更` is in-scope but has no corresponding user story | -1 pt (Dim 5 consistency) |
| prd-spec.md — Scope | `error-fixer agent 废弃` is in-scope but has no corresponding user story (forge 维护者 impact) | -1 pt (Dim 5 consistency) |
| prd-user-stories.md — Story 3 AC1 | "包含任务 ID、scope、执行步骤**等上下文**" — "等上下文" is vague; exact required fields are unspecified | -2 pts (Dim 4 verifiability) |
| prd-user-stories.md — Story 2 | Only one AC; no error case (template syntax error, missing type enum registration) | -1 pt (Dim 4 verifiability) |

---

## Attack Points

### Attack 1: Flow Completeness — data flow between components undocumented

**Where**: prd-spec.md Flow Description — "`task prompt <id>` 内部逻辑：1. 从 `.forge/state.json` 和 `index.json` 读取 feature、scope、task type 等上下文"

**Why it's weak**: This is a multi-component feature touching run-tasks, task-cli, task-executor, index.json, and state.json. The flow description narrates what each component does but never specifies the data contracts: what exact fields does `task prompt` read from state.json? What does it write to stdout (schema)? What does task-executor receive as its prompt parameter? Without a data flow table, an implementer cannot verify that the components will interoperate correctly. The rubric awards 0-7 for "data flow documented (if multi-system)" and this section earns only 4/7.

**What must improve**: Add a data flow table showing: source component → data produced → destination component → field names. At minimum: state.json fields consumed by `task prompt`, stdout schema of `task prompt`, and the prompt parameter structure passed to task-executor.

---

### Attack 2: Flow Diagrams — `--fix-record-missed` failure path and eval-cases error branch are missing

**Where**: prd-spec.md Business Flow Diagram — `FixRecord["task prompt <id> --fix-record-missed\n→ Agent(forge:task-executor)"]` leads directly to `Record2` with no failure branch; `PromptEval["task prompt <id>"]` leads directly to `MainSession` with no failure branch.

**Why it's weak**: The diagram shows `PromptFail` as a decision node for the normal path but silently omits the same check for the eval-cases path and for the fix-record-missed recovery path. If `task prompt --fix-record-missed` fails, the diagram implies the task silently succeeds (reaches Record2). If `task prompt <id>` fails in the eval-cases branch, the diagram implies execution continues to `MainSession`. Both are silent failure modes — exactly the class of bug the PRD claims to fix ("record 文件未写入的静默失败").

**What must improve**: Add `PromptFail` decision nodes after `PromptEval` (eval-cases path) and after `FixRecord`. Both failure branches should route to `MarkBlocked` or an explicit error state.

---

### Attack 3: User Stories — Story 3 AC1 is partially unverifiable; Stories 2 and 5 lack error ACs

**Where**: prd-user-stories.md Story 3 AC1 — "stdout 输出该任务的完整合成 prompt，包含任务 ID、scope、**执行步骤等上下文**"; Story 2 — single AC with no error case; Story 5 — single AC with no error case.

**Why it's weak**: "执行步骤等上下文" cannot be objectively verified — a tester cannot determine pass/fail without knowing the exact required fields. This is the same vague-language pattern the PRD's own Background section criticizes ("agent 行为取决于多个隐式字段的组合"). Story 2 has no AC for the failure case (what if the template file has a syntax error, or the type enum registration is missing?). Story 5 has no AC for the case where breakdown-tasks cannot determine the type for an ambiguous task.

**What must improve**: Replace "等上下文" with an explicit enumeration of required fields in the synthesized prompt. Add error ACs to Story 2 (template parse failure → non-zero exit, stderr message) and Story 5 (ambiguous type → explicit error or fallback to `implementation` with warning).

---

## Previous Issues Check

<!-- Only for iteration > 1 — N/A for iteration 1 -->

---

## Verdict

- **Score**: 85/100
- **Target**: 90/100
- **Gap**: 5 points
- **Action**: Continue to iteration 2 — address Attack 1 (data flow table), Attack 2 (missing error branches in diagrams), and Attack 3 (AC verifiability gaps) to close the 5-point gap.
