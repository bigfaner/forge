---
date: "2026-04-30"
doc_dir: "docs/features/justfile-standard-vocabulary/prd/"
iteration: 1
target_score: 90
evaluator: Claude (automated, adversarial)
---

# PRD Eval — Iteration 1

**Score: 79/100** (target: 90)

```
┌─────────────────────────────────────────────────────────────────┐
│                       PRD QUALITY SCORECARD                      │
├──────────────────────────────┬──────────┬──────────┬────────────┤
│ Dimension                    │ Score    │ Max      │ Status     │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 1. Background & Goals        │  15      │  20      │ ⚠️         │
│    Background three elements │  6/7     │          │            │
│    Goals quantified          │  5/7     │          │            │
│    Logical consistency       │  4/6     │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 2. Flow Diagrams             │  15      │  20      │ ⚠️         │
│    Mermaid diagram exists    │  7/7     │          │            │
│    Main path complete        │  5/7     │          │            │
│    Decision + error branches │  3/6     │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 3. Functional Specs          │  13      │  20      │ ❌         │
│    Tables complete           │  5/7     │          │            │
│    Field descriptions clear  │  5/7     │          │            │
│    Validation rules explicit │  3/6     │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 4. User Stories              │  19      │  20      │ ✅         │
│    Coverage per user type    │  7/7     │          │            │
│    Format correct            │  7/7     │          │            │
│    AC per story              │  5/6     │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 5. Scope Clarity             │  17      │  20      │ ⚠️         │
│    In-scope concrete         │  6/7     │          │            │
│    Out-of-scope explicit     │  7/7     │          │            │
│    Consistent with specs     │  4/6     │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ TOTAL                        │  79      │  100     │            │
└──────────────────────────────┴──────────┴──────────┴────────────┘
```

---

## Deductions

| Location | Issue | Penalty |
|----------|-------|---------|
| prd-spec.md:44 | "支持混合项目" goal lacks numeric metric — "混合项目可通过 scope 参数选择性操作前端/后端" is a boolean capability, not quantified | -2 pts |
| prd-spec.md:45 | "自适应生成" goal vague — "init-justfile 根据项目结构自动生成匹配类型的 justfile" has no metric like "% of project types correctly detected" | -2 pts |
| prd-spec.md:42 | "当前至少 8 处散落的原始命令" but migration table (5.4) shows only 3 rows with actual command changes; other 8 rows say "保持不变" | -3 pts (inconsistency) |
| prd-spec.md:109-126 | Command table lists 15 commands but goals and scope repeatedly state "16 个标准命令" | -3 pts (inconsistency) |
| prd-spec.md:81-94 | No flow diagram for init-justfile generation process — a core feature | -2 pts |
| prd-spec.md:82-93 | No error branch for `just project-type` failure (missing recipe, unexpected output, non-zero exit) | -3 pts |
| prd-spec.md:109-126 | No validation rules for invalid scope arguments (e.g. `just build invalidscope`) | -3 pts |
| prd-user-stories.md:68 | Story 5 AC mandates specific warning format "[forge] scope=frontend but project-type=backend; falling back to just build" but prd-spec flow diagram only says "⚠ 警告: scope 不匹配项目类型" — message format mismatch | -2 pts (inconsistency) |
| prd-spec.md:147-153 | `scope` field in index.json lacks JSON schema: type (string enum?), optional/required, default value | -2 pts |

---

## Attack Points

### Attack 1: Functional Specs — 16 vs 15 command count inconsistency

**Where**: prd-spec.md:24 "将标准命令词汇从 6 个扩展至 16 个" vs prd-spec.md:109-126 command table with 15 rows
**Why it's weak**: The goals section, scope section, and background all assert "16 个标准命令" as a hard metric, but the authoritative command vocabulary table (section 5.1) contains only 15 entries. This means either the table is missing a command, or the count was wrong from the start. Either way, the document's core quantified target is internally contradicted. If the PRD cannot agree with itself on how many commands it defines, the migration list and acceptance criteria built on that number are unreliable.
**What must improve**: Either add the missing 16th command to the table (and update the migration checklist if needed), or correct all "16 个" references to "15 个". Then cross-check every downstream reference (goals, scope, AC) against the corrected count.

### Attack 2: Flow Diagrams — missing init-justfile generation flow and error handling

**Where**: prd-spec.md:81-94 — only one Mermaid diagram exists (scope resolution at runtime)
**Why it's weak**: The document describes two major functional flows: (a) the scope resolution flow at skill-execution time, and (b) the init-justfile adaptive generation flow. Only (a) has a diagram. The init-justfile generation — detecting project signals, selecting justfile template, handling pre-existing justfiles — is described only in prose (section 5.2) with a 3-row table. Furthermore, the existing diagram has no error branch for `just project-type` failure: what should a skill do when the justfile lacks a `project-type` recipe, or when the command returns an unexpected string? The agent-friendliness requirements (section "Agent 友好性需求") demand predictable error behavior, but the flow diagram ignores that path entirely.
**What must improve**: Add a second Mermaid flowchart for the init-justfile generation process (detect signals -> classify project -> generate template -> handle existing justfile). In the existing scope-resolution diagram, add an error branch for `just project-type` failure (non-zero exit, missing recipe, unexpected output) and define the fallback behavior.

### Attack 3: Functional Specs — missing validation rules for scope arguments and field schema

**Where**: prd-spec.md:127-133 "scope 参数规则" table and prd-spec.md:147-153 "scope 值" table
**Why it's weak**: The scope parameter rules table explains the happy path (which project types have scope, what each scope means) but never defines validation behavior for invalid inputs. What happens when a user or agent runs `just build frontend` on a pure-backend project justfile that has no scope parameter? What about `just build backned` (typo)? The document assumes only correct usage. Similarly, section 5.3 says breakdown-tasks adds a `scope` field to index.json but never specifies: is the field optional or required? What is the default if omitted? What is the JSON type — a string enum constrained to exactly three values? Without these details, two implementers could produce incompatible results, and the AC in Story 3 ("Then `index.json` 中每个任务包含 `scope` 字段") cannot be unambiguously verified.
**What must improve**: (1) Add a validation rules section under 5.1 specifying behavior for invalid scope values: reject with non-zero exit and stderr message, or silently treat as `all`. (2) Add a JSON schema snippet for the `scope` field in index.json — type: string, enum: [frontend, backend, all], required: true (or false with default "all"). (3) Define what happens when `just project-type` is called on a justfile that lacks the recipe.

---

## Previous Issues Check

<!-- Only for iteration > 1 — N/A for iteration 1 -->

---

## Verdict

- **Score**: 79/100
- **Target**: 90/100
- **Gap**: 11 points
- **Action**: Continue to iteration 2. Priority fixes: (1) resolve 16-vs-15 command count inconsistency, (2) add init-justfile generation flow diagram and error handling branches, (3) add validation rules for scope arguments and JSON schema for the scope field.
