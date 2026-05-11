---
date: "2026-05-11"
doc_dir: "docs/features/typed-task-dispatch/prd/"
iteration: "2"
target_score: "90"
scoring_mode: "Mode B"
evaluator: Claude (automated, adversarial)
---

# PRD Eval — Iteration 2

**Score: 90/100** (target: 90, mode: Mode B)

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
│ 2. Flow Diagrams             │  19      │  20      │ ✅          │
│    Mermaid diagram exists    │  7/7     │          │            │
│    Main path complete        │  7/7     │          │            │
│    Decision + error branches │  5/6     │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 3b. Flow Completeness (B)    │  16      │  20      │ ⚠️          │
│    Flow steps complete       │  6/7     │          │            │
│    Data flow documented      │  6/7     │          │            │
│    Exception handling        │  4/6     │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 4. User Stories              │  28      │  30      │ ✅          │
│    Coverage per user type    │  7/7     │          │            │
│    Format correct            │  7/7     │          │            │
│    AC per story (G/W/T)      │  6/6     │          │            │
│    AC verifiability          │  8/10    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 5. Scope Clarity             │  13      │  15      │ ⚠️          │
│    In-scope concrete         │  5/5     │          │            │
│    Out-of-scope explicit     │  4/4     │          │            │
│    Consistent with specs     │  4/6     │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ TOTAL                        │  90      │  100     │            │
└──────────────────────────────┴──────────┴──────────┴────────────┘
```

> **Mode B** (prd-ui-functions.md absent): Dimension 3 evaluates Flow Completeness from prd-spec.md Flow Description.
> Sub-criteria: Flow steps describe complete business process /7, Data flow documented /7, Exception handling and edge cases /6.

---

## Deductions

| Location | Issue | Penalty |
|----------|-------|---------|
| prd-spec.md — Background/What | "What (Target)" describes the technical solution ("将策略逻辑从 agent markdown 下沉到 CLI"), not the business outcome — conflates solution with target | -1 pt (Dim 1 consistency) |
| prd-spec.md — Flow Diagram (task prompt internal) | Template-not-found and template-parse-failure error branches absent; only type-missing branch shown (`D{type 字段存在?} → E`) despite functional spec listing "模板不存在、模板解析失败等" as error cases | -1 pt (Dim 2 error branches) |
| prd-spec.md — Flow Description | `MarkBlocked` state has no recovery path described: what triggers unblocking? Is it recoverable? The text never answers this | -1 pt (Dim 3 flow steps) |
| prd-spec.md — Flow Description | Data flow table omits `task migrate` and `task validate` read/write paths; also omits who writes `index.json` and what fields | -1 pt (Dim 3 data flow) |
| prd-spec.md — Flow Description | System-level failures undocumented: `.forge/state.json` missing or malformed, task ID not found in index.json, task-executor agent execution failure (distinct from prompt generation failure) | -2 pts (Dim 3 exception handling) |
| prd-user-stories.md — Story 2 AC1 | "task prompt <id> 对该类型输出正确 prompt" — "正确" is unverifiable without specifying required fields/content in the synthesized prompt | -1 pt (Dim 4 verifiability) |
| prd-user-stories.md — Story 1 | No error AC: what happens if task-executor fails to execute the doc-generation or fix flow? Only happy-path ACs present | -1 pt (Dim 4 verifiability) |
| prd-spec.md — Scope | `execute-task.md 路由变更` is in-scope (Related Changes #6) but has no corresponding user story | -1 pt (Dim 5 consistency) |
| prd-spec.md — Scope | `error-fixer agent 废弃` is in-scope (Related Changes #10) but has no corresponding user story (forge 维护者 impact) | -1 pt (Dim 5 consistency) |

---

## Attack Points

### Attack 1: Flow Completeness — blocked state is a dead end with no recovery path

**Where**: prd-spec.md Business Flow Diagram — `MarkBlocked[标记任务 blocked]` → `End([Continue Loop])`. The flow description states "run-tasks 检测到后标记任务 blocked，不静默失败" but never describes what happens after.

**Why it's weak**: The PRD's own Background section cites "record 文件未写入的静默失败，排查额外耗时 1 小时" as a motivating incident. The fix is to mark tasks blocked instead of failing silently — but the PRD never specifies what a blocked task means operationally: Can run-tasks retry it? Does the operator need to manually intervene? Is there a `task unblock` command? Without this, "blocked" is just a different kind of silent failure — the task stops progressing with no documented recovery path. An implementer reading this PRD cannot determine what the system should do when a task enters the blocked state.

**What must improve**: Add a section (or AC) describing the blocked state lifecycle: what triggers it, how it is surfaced to the operator, and how it is resolved. At minimum, the flow description should state whether blocked tasks are retryable and what manual intervention looks like.

---

### Attack 2: User Stories — Story 2 AC1 "正确 prompt" is objectively unverifiable

**Where**: prd-user-stories.md Story 2 AC1 — "Then `task prompt <id>` 对该类型输出正确 prompt，Go 单元测试可覆盖该类型，无需修改 task-executor.md 或任务模板"

**Why it's weak**: "正确 prompt" has no definition. A tester cannot determine pass/fail without knowing what fields the synthesized prompt must contain for a new type. Story 3 AC1 was fixed in this iteration by enumerating required fields ("任务 ID、scope、任务 type、执行步骤、phase summary 路径"), but Story 2 AC1 still uses the same vague pattern. The rubric penalizes any "Then" that cannot be verified without subjective judgment — "正确" is entirely subjective here.

**What must improve**: Replace "正确 prompt" with an explicit enumeration of required fields, mirroring Story 3 AC1's structure. For example: "stdout 包含任务 ID、scope、新 type 的执行步骤（来自新增模板文件），且不包含其他 type 的执行步骤".

---

### Attack 3: Scope Clarity — two in-scope changes have no user stories or acceptance criteria

**Where**: prd-spec.md Scope In-Scope — "execute-task 命令同步更新：与 run-tasks 保持相同路由逻辑" and "error-fixer agent 废弃：编译/测试/lint 修复能力迁移到 `type: fix` prompt 模板". Related Changes table items #6 and #10.

**Why it's weak**: `execute-task` is a separate entry point from `run-tasks` — it has its own routing logic and its own callers. "与 run-tasks 保持相同路由逻辑" is a behavioral claim with no AC to verify it. `error-fixer` deprecation is a breaking change for any caller that dispatches to it directly; there is no story describing the migration path or the forge 维护者's experience of removing it. Both changes are non-trivial and could introduce regressions, yet neither has a single testable acceptance criterion.

**What must improve**: Add a user story (or extend an existing one) covering: (a) execute-task routing parity with run-tasks, with an AC verifying identical behavior for at least one task type; (b) error-fixer deprecation, with an AC confirming that fix-type tasks routed through the new path produce equivalent outcomes and that the old error-fixer dispatch path is removed.

---

## Previous Issues Check

| Previous Attack | Addressed? | Evidence |
|----------------|------------|----------|
| Attack 1: Flow Completeness — data flow between components undocumented | ✅ | Data flow table added ("数据流表（多组件数据契约）") covering state.json, index.json, task prompt stdout/stderr, and task-executor prompt parameter |
| Attack 2: Flow Diagrams — `--fix-record-missed` failure path and eval-cases error branch missing | ✅ | `PromptFailEval` decision node added to eval-cases path; `FixFail` decision node added to fix-record-missed path; both route to `MarkBlocked` |
| Attack 3: User Stories — Story 3 AC1 vague; Stories 2 and 5 missing error ACs | ✅ | Story 3 AC1 now enumerates required fields explicitly; Story 2 gained an error AC (template syntax error / unregistered type); Story 5 gained an error AC (ambiguous type → fallback to `implementation` with warning) |

---

## Verdict

- **Score**: 90/100
- **Target**: 90/100
- **Gap**: 0 points
- **Action**: Target reached. Remaining weaknesses (blocked state recovery, Story 2 AC1 verifiability, execute-task/error-fixer story coverage) are documented above for optional follow-up but do not block the target.
