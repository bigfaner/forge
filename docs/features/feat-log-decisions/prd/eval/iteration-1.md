---
date: "2026-04-22"
doc_dir: "docs/features/feat-log-decisions/prd/"
iteration: 1
target: "90"
evaluator: Claude (automated, adversarial)
---

# PRD Eval — Iteration 1

**Score: 81/100** (target: 90)

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
| 2. Flow Diagrams             |  17      |  20      | W          |
|    Mermaid diagram exists    |  7/7     |          |            |
|    Main path complete        |  7/7     |          |            |
|    Decision + error branches |  3/6     |          |            |
+──────────────────────────────+──────────+──────────+────────────+
| 3. Functional Specs          |  13      |  20      | X          |
|    Tables complete           |  6/7     |          |            |
|    Field descriptions clear  |  5/7     |          |            |
|    Validation rules explicit |  2/6     |          |            |
+──────────────────────────────+──────────+──────────+────────────+
| 4. User Stories              |  16      |  20      | W          |
|    Coverage per user type    |  3/7     |          |            |
|    Format correct            |  7/7     |          |            |
|    AC per story              |  6/6     |          |            |
+──────────────────────────────+──────────+──────────+────────────+
| 5. Scope Clarity             |  18      |  20      | C          |
|    In-scope concrete         |  7/7     |          |            |
|    Out-of-scope explicit     |  7/7     |          |            |
|    Consistent with specs     |  4/6     |          |            |
+──────────────────────────────+──────────+──────────+────────────+
| TOTAL                        |  81      |  100     |            |
+──────────────────────────────+──────────+──────────+────────────+
```

---

## Deductions

| Location | Issue | Penalty |
|----------|-------|---------|
| prd-spec.md:33-38 (Goals table) | "命名方向与输出文件一致" is a description, not a quantified metric; "manifest.md 实时反映所有决策" has no numeric target | -2 pts (Goals quantified) |
| prd-spec.md:14-17 vs 19-28 | Background problem 1 (naming) is trivial compared to the PRD's dominant focus on decision archiving; weight imbalance | -1 pt (Logical consistency) |
| prd-spec.md:81-103 | No error/exception branches in either flow diagram (e.g., invalid input, file write failure) | -3 pts (Decision + error branches) |
| prd-spec.md:135-139 | `edit:<N>` interaction lacks description of edit UI/flow in the table | -1 pt (Tables complete) |
| prd-spec.md:148-153 | "Source" field has no format specification; "Feature" slug has no validation criteria | -2 pts (Field descriptions clear) |
| prd-spec.md:135-139, 186-192 | Validation rules are shown by example (`1,3`, `all`, `none`) but never stated explicitly; no rules for invalid number, empty field, out-of-range type selection | -4 pts (Validation rules explicit) |
| prd-user-stories.md (all) | Zero stories for "Skill 开发者" user type defined in background | -4 pts (Coverage per user type) |
| prd-spec.md:46 vs 5.2 | Scope item "创建 templates/decision-entry.md 决策条目模板" has no corresponding description in functional specs; 5.2 describes table rows, not a template | -2 pts (Consistent with specs) |

---

## Attack Points

### Attack 1: Functional Specs — Validation rules are absent

**Where**: prd-spec.md sections 5.1 (user interaction table, lines 135-139) and 5.4 (record-decision rounds, lines 186-192)

**Why it's weak**: The user interaction table shows valid inputs by example (`1,3`, `all`, `none`, `edit:<N>`) but never states what happens when input is invalid. For record-decision's 4 rounds: what happens if the user enters empty text for "决策描述"? What if they type "0" or "9" for type selection (only 1-8 valid)? What if the feature slug does not match any existing feature directory? These are basic validation concerns that any implementer would need answered. The document received 2/6 on this criterion because examples are not rules.

**What must improve**: Add explicit validation rules per input field. For example: "Type selection must be integer 1-8; invalid input re-prompts with error message." "Decision description must be non-empty string; empty input re-prompts." "Feature slug must match an existing directory under `docs/features/`; mismatch warns but allows override."

### Attack 2: User Stories — Skill 开发者 has zero coverage

**Where**: prd-user-stories.md — all four stories use "As a Skill 使用者"

**Why it's weak**: The background section (prd-spec.md lines 27-29) explicitly defines two user types: "Skill 开发者" (maintains SKILL.md, references, templates) and "Skill 使用者" (calls skills). The PRD's scope includes significant developer-facing work: renaming directories, updating SKILL.md frontmatter, updating hooks guide, updating CLAUDE.md skill index, creating new template and reference files. None of these have a user story from the Skill 开发者 perspective. A developer tasked with implementing the rename or creating the new reference file has no acceptance criteria to verify against.

**What must improve**: Add at least 2 stories for "Skill 开发者". For example: (1) As a Skill 开发者, I want to rename the skill directory and update all references so that `/zcode:tech-design` resolves correctly. (2) As a Skill 开发者, I want the decision-logging reference to be shared between tech-design and record-decision so that extraction logic is not duplicated.

### Attack 3: Flow Diagrams — No error branches

**Where**: prd-spec.md lines 81-103 (Mermaid flowchart)

**Why it's weak**: Both flows (A and B) only show the happy path. Flow A has decision diamonds for user approval and selection, but no error branches: what if file write to `docs/decisions/<type>.md` fails? What if `manifest.md` is missing or corrupted? What if the user enters an invalid selection (e.g., "abc" when prompted for numbers)? Flow B's 4-round interaction has no path for user cancellation mid-way, no path for invalid type selection. A developer implementing this cannot know the expected error handling behavior.

**What must improve**: Add at least 2 error branches to the flow diagram. For example: (1) After "Archive" node, add a branch for "write failure" leading to "display error + retry". (2) In Flow B, add a path from Q1 for "invalid input" leading back to Q1 with error prompt. (3) Consider adding a "cancel" exit from any round in Flow B.

---

## Previous Issues Check

<!-- Only for iteration > 1 — not applicable for iteration 1 -->

---

## Verdict

- **Score**: 81/100
- **Target**: 90/100
- **Gap**: 9 points
- **Action**: Continue to iteration 2

SCORE: 81/100
DIMENSIONS:
  Background & Goals: 17/20
  Flow Diagrams: 17/20
  Functional Specs: 13/20
  User Stories: 16/20
  Scope Clarity: 18/20
ATTACKS:
1. Functional Specs: Validation rules are absent — prd-spec.md sections 5.1 and 5.4 show valid inputs by example (`1,3`, `all`, `none`) but never state what happens on invalid input — Add explicit validation rules per input field (type range, non-empty constraints, slug format)
2. User Stories: Skill 开发者 has zero coverage — All 4 stories use "As a Skill 使用者" while background defines two user types — Add at least 2 stories for Skill 开发者 covering rename, reference/template creation, and CLAUDE.md updates
3. Flow Diagrams: No error branches — Mermaid flowchart (lines 81-103) shows only happy paths with no error/exception handling paths — Add at least 2 error branches (write failure, invalid input, mid-flow cancellation)
