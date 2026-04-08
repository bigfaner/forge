---
name: eval-prd
description: Evaluate a PRD document against quality standards. Checks structure completeness, user story quality, acceptance criteria testability, scope clarity, and goals measurability. Outputs a scored report with actionable improvements.
---

# Eval PRD

评估 PRD 文档是否满足规范，输出评分报告和改进建议。

## When to Use

**Trigger:**
- User asks to "evaluate PRD" or "check PRD quality"
- User provides `/eval-prd` command
- Before handing off PRD to `/breakdown-tasks`

**Skip:**
- PRD doesn't exist yet (use `/write-prd` first)

## Workflow

```
1. 定位 PRD → 2. 检查结构 → 3. 检查内容质量 → 4. 生成报告
```

## Step 1: Locate PRD

Check in order:
1. Path provided by user
2. `docs/features/<current-feature>/prd.md`
3. Ask user for path if not found

## Step 2: Check Structure Completeness

Required sections — mark missing ones as F immediately:

| Section              | Required | Notes                              |
| -------------------- | -------- | ---------------------------------- |
| Background           | ✓        | Must have problem statement        |
| Goals                | ✓        | Must have success metrics          |
| Scope                | ✓        | Must have both in/out of scope     |
| User Stories         | ✓        | At least one story per target user |
| Functional Requirements | ✓     | At least one FR                    |
| Acceptance Criteria  | ✓        | Must be testable checkboxes        |
| Non-Functional Requirements | ○  | Optional but recommended           |
| Dependencies         | ○        | Optional                           |
| Risks                | ○        | Optional                           |

## Step 3: Check Content Quality

### Dimension 1: User Stories (用户故事)

| Check | Criteria |
|-------|----------|
| Coverage | At least one story per target user identified in Background |
| Format | Each story follows "As a / I want / So that" |
| Specificity | Action is concrete, not vague ("view task list" not "manage tasks") |
| AC per story | Each story has Given/When/Then acceptance criteria |

**Grading:**
- A: All stories present, correct format, specific actions, AC attached
- B: All stories present, minor format issues or 1 missing AC
- C: Stories exist but vague, or missing AC on most stories
- F: No user stories, or only one user covered when multiple exist

### Dimension 2: Acceptance Criteria (验收标准)

| Check | Criteria |
|-------|----------|
| Testable | Each criterion can be verified by a human or automated test |
| Specific | No vague terms: "fast", "easy", "good", "support" without definition |
| Complete | Covers all user stories and key functional requirements |
| Format | Checkbox list `- [ ]` |

**Grading:**
- A: All criteria testable, specific, complete, checkbox format
- B: Most testable, 1-2 vague items
- C: Many vague items, or missing coverage for key features
- F: No AC, or AC is just a copy of requirements

### Dimension 3: Goals & Metrics (目标与指标)

| Check | Criteria |
|-------|----------|
| Measurable | Success metrics have numeric targets (not "improve X") |
| Non-goals | Explicitly lists what is NOT in scope |
| Motivation | Background explains *why* this matters, not just *what* |

**Grading:**
- A: Metrics are quantified, non-goals explicit, motivation clear
- B: Metrics defined but not all quantified
- C: Goals stated but no metrics, or no non-goals
- F: No goals section, or goals are purely technical with no user value

### Dimension 4: Scope Clarity (范围清晰度)

| Check | Criteria |
|-------|----------|
| In-scope | Each item is a concrete deliverable, not a vague capability |
| Out-of-scope | Explicitly lists deferred items (prevents scope creep) |
| Consistency | Scope items align with user stories and requirements |

**Grading:**
- A: Both in/out scope defined, items concrete, consistent with stories
- B: Both defined, minor vagueness
- C: Only in-scope defined, or items are vague
- F: No scope section

### Dimension 5: Requirements Quality (需求质量)

| Check | Criteria |
|-------|----------|
| Priority | Each FR has P0/P1/P2 priority |
| Traceability | FRs can be traced back to user stories |
| Completeness | No obvious gaps between stories and FRs |

**Grading:**
- A: All FRs prioritized, traceable to stories, no gaps
- B: Most prioritized, minor gaps
- C: No priorities, or FRs don't map to stories
- F: No functional requirements

## Step 4: Generate Report

### Grading Rules

**Overall:**

| Grade | Condition |
|-------|-----------|
| A | All 5 dimensions A/B, at least 3 A's |
| B | No F, max 1 C |
| C | 1 F or 2+ C's |
| D | 2 F's |
| F | 3+ F's or User Stories missing entirely |

Save report to `docs/features/<feature-slug>/prd-eval.md` using `templates/report.md`.

## Related

- `/write-prd` — Create or revise the PRD
- `/breakdown-tasks` — Next step after PRD passes evaluation
