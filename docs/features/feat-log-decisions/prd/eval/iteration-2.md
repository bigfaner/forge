---
date: "2026-04-22"
doc_dir: "docs/features/feat-log-decisions/prd/"
iteration: 2
target_score: "90"
evaluator: Claude (automated, adversarial)
---

# PRD Eval — Iteration 2

**Score: 94/100** (target: 90)

```
+─────────────────────────────────────────────────────────────────+
|                       PRD QUALITY SCORECARD                      |
+──────────────────────────────+──────────+──────────+────────────+
| Dimension                    | Score    | Max      | Status     |
+──────────────────────────────+──────────+──────────+────────────+
| 1. Background & Goals        |  17      |  20      | W          |
|    Background three elements |  7/7     |          |            |
|    Goals quantified          |  5/7     |          |            |
|    Logical consistency       |  5/6     |          |            |
+──────────────────────────────+──────────+──────────+────────────+
| 2. Flow Diagrams             |  20      |  20      | C          |
|    Mermaid diagram exists    |  7/7     |          |            |
|    Main path complete        |  7/7     |          |            |
|    Decision + error branches |  6/6     |          |            |
+──────────────────────────────+──────────+──────────+────────────+
| 3. Functional Specs          |  19      |  20      | C          |
|    Tables complete           |  7/7     |          |            |
|    Field descriptions clear  |  6/7     |          |            |
|    Validation rules explicit |  6/6     |          |            |
+──────────────────────────────+──────────+──────────+────────────+
| 4. User Stories              |  20      |  20      | C          |
|    Coverage per user type    |  7/7     |          |            |
|    Format correct            |  7/7     |          |            |
|    AC per story              |  6/6     |          |            |
+──────────────────────────────+──────────+──────────+────────────+
| 5. Scope Clarity             |  18      |  20      | C          |
|    In-scope concrete         |  7/7     |          |            |
|    Out-of-scope explicit     |  7/7     |          |            |
|    Consistent with specs     |  4/6     |          |            |
+──────────────────────────────+──────────+──────────+────────────+
| TOTAL                        |  94      |  100     |            |
+──────────────────────────────+──────────+──────────+────────────+
```

---

## Deductions

| Location | Issue | Penalty |
|----------|-------|---------|
| prd-spec.md:35 (Goals table row 1) | "命名方向与输出文件一致" is a qualitative description, not a quantified metric — no %, count, or time target | -2 pts (Goals quantified) |
| prd-spec.md:14-17 vs 19-28 | Background problem 1 (naming inconsistency) is trivial relative to the PRD's dominant scope (decision archiving system); the weight imbalance means the first goal feels bolted on | -1 pt (Logical consistency) |
| prd-spec.md:185 (Source field) | "来源文件和章节，格式 `<文件路径>#<章节>`" — the `#<章节>` portion assumes a markdown heading exists at the target, but no rule states what happens when the heading does not exist or the file has no matching section anchor | -1 pt (Field descriptions clear) |
| prd-spec.md:46 vs prd-user-stories.md | In-scope items "更新 hooks guide 和 exploration 示例中的引用" and "更新 `zcode/CLAUDE.md` 中的 skill 列表" are bundled into Story 4's compound ACs rather than having dedicated stories; "创建 `docs/decisions/manifest.md` 索引文件和 8 个类型模板文件" has no AC verifying all 8 template files — traceability is imprecise | -2 pts (Consistent with specs) |

---

## Attack Points

### Attack 1: Background & Goals — First and fourth goals are unquantified

**Where**: prd-spec.md line 35, goals table row 1: "命名方向与输出文件一致 | 命名方向与输出文件一致 | 调用 `/zcode:tech-design` 输出 `tech-design.md`" and line 38: "集中管理决策索引 | `docs/decisions/manifest.md` 实时反映所有决策"

**Why it's weak**: The "量化指标" column for row 1 literally repeats the goal text ("命名方向与输出文件一致") rather than providing a numeric or binary measurable target. Compare with row 2: "100% 的 design 阶段可产生归档" — that is a quantified metric. Row 4: "manifest.md 实时反映所有决策" also lacks a numeric target — "实时" is vague (does it mean instant? within 5 seconds? after each write?). Two of four goals use vague language in the metric column.

**What must improve**: Replace the non-numeric "metrics" with measurable targets. For row 1: "skill 名称与输出文件名完全一致（1/1 匹配）" or "重命名后 0 处残留旧名称引用". For row 4: "manifest.md 条目数与所有类型文件行数之和一致（100% 同步）" or give a time bound like "写入决策后 manifest.md 在 5 秒内同步完成".

### Attack 2: Functional Specs — Source field missing from record-decision input rounds

**Where**: prd-spec.md section 5.4 (lines 218-223) defines 4 interaction rounds vs section 5.2 (lines 179-185) which defines Source as a required field

**Why it's weak**: Section 5.2 defines 5 required fields per decision record: Date, Feature, Decision, Rationale, Source. Section 5.4's 4-round interaction collects only type, description, rationale, and feature slug. The Source field is never collected from the user in the record-decision path, and no auto-generation rule is specified. This means an implementer cannot know how to populate Source when a decision is recorded via `/zcode:record-decision`. Additionally, the Source format `<文件路径>#<章节>` has no edge-case handling: what if the heading does not exist?

**What must improve**: Either (a) add a 5th round to record-decision for Source input, or (b) explicitly state that Source is auto-generated (e.g., `record-decision (manual entry)` or `<feature>/tech-design.md`) when called from record-decision. Also specify what happens when a heading anchor in Source does not match any section in the target file.

### Attack 3: Scope Clarity — In-scope items have imprecise story traceability

**Where**: prd-spec.md lines 46-51 (In Scope) vs prd-user-stories.md Stories 1-6

**Why it's weak**: The In Scope list has 8 items. "创建 `docs/decisions/manifest.md` 索引文件和 8 个类型模板文件" bundles manifest + 8 templates into one checkbox, yet no story's AC explicitly verifies that all 8 type template files were created with correct headers. "更新 hooks guide 和 exploration 示例中的引用" and "更新 `zcode/CLAUDE.md` 中的 skill 列表" are bundled into Story 4's compound ACs. When a developer checks off scope items, they cannot cleanly map each to a verified story.

**What must improve**: Ensure every in-scope item has at least one story or AC that directly verifies it. Either expand Story 5's ACs to explicitly verify all 8 template files, or split the manifest/templates scope item into two separate items with dedicated verification. Make the traceability from scope to stories explicit.

---

## Previous Issues Check

| Previous Attack | Addressed? | Evidence |
|----------------|------------|----------|
| Attack 1: Validation rules absent in 5.1 and 5.4 | YES | New "输入校验规则" table (lines 165-171) with explicit rules per input field; section 5.4 per-round validation table (lines 218-223) with legal values, illegal handling, and re-prompt messages |
| Attack 2: Skill 开发者 zero story coverage | YES | Two new stories added: Story 4 (lines 50-68) and Story 5 (lines 71-85), both "As a Skill 开发者" covering rename/references and shared template creation |
| Attack 3: No error branches in flow diagrams | YES | Mermaid flowchart now includes: invalid input loops (ErrA1, ErrQ1-Q4), write failure retry paths (ErrA2, ErrWrite), and cancel exits (CancelEnd) for both Flow A and Flow B |

---

## Verdict

- **Score**: 94/100
- **Target**: 90/100
- **Gap**: 0 points (target exceeded)
- **Action**: Target reached

SCORE: 94/100
DIMENSIONS:
  Background & Goals: 17/20
  Flow Diagrams: 20/20
  Functional Specs: 19/20
  User Stories: 20/20
  Scope Clarity: 18/20
ATTACKS:
1. Background & Goals: First and fourth goals use qualitative text as "metrics" — prd-spec.md line 35 "命名方向与输出文件一致" in the metric column is not a number, and line 38 "manifest.md 实时反映所有决策" has no numeric/time target — Replace with measurable targets like "0 处残留旧名称引用" and "写入后 5 秒内 manifest 同步完成"
2. Functional Specs: Source field missing from record-decision input rounds — prd-spec.md section 5.4 (lines 218-223) collects 4 fields (type, description, rationale, feature) but section 5.2 defines Source as a required 5th field — Either add a 5th round for Source input or explicitly state how Source is auto-populated in the record-decision path
3. Scope Clarity: In-scope items have imprecise story traceability — "创建 docs/decisions/manifest.md 索引文件和 8 个类型模板文件" has no AC verifying all 8 template files; multiple scope items are bundled into Story 4's compound ACs — Add dedicated ACs or stories for each in-scope item that currently lacks explicit verification
