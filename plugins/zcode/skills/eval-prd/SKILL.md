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
1. 定位 PRD → 2. 启动评估 Agent → 3. 汇报结果
```

## Step 1: Locate PRD

Check in order:
1. Path provided by user
2. `docs/features/<current-feature>/prd.md`
3. Ask user for path if not found

Determine `<feature-slug>` from the path (e.g. `docs/features/auth-flow/prd.md` → slug is `auth-flow`).

## Step 2: Launch Evaluation Agent

Use the **Agent tool** to spawn a subagent. Pass the full prompt below, substituting `{{PRD_PATH}}` and `{{FEATURE_SLUG}}`:

---

**Agent prompt template:**

```
You are a PRD quality evaluator. Your job: read the PRD, apply the rubric, write the report, return a summary.

## Inputs
- PRD path: {{PRD_PATH}}
- Feature slug: {{FEATURE_SLUG}}
- Report output: docs/features/{{FEATURE_SLUG}}/prd-eval.md
- Report template: plugins/zcode/skills/eval-prd/templates/report.md

## Steps
1. Read {{PRD_PATH}}
2. Read the report template
3. Apply the rubric below to every dimension
4. Fill in the template and write to docs/features/{{FEATURE_SLUG}}/prd-eval.md
5. Return: overall grade, top 2-3 issues, and whether it can proceed to /breakdown-tasks

## Structure Check

Required sections — mark missing as F:

| Section                 | Required | Notes                          |
|-------------------------|----------|--------------------------------|
| Background              | ✓        | Must have problem statement    |
| Goals                   | ✓        | Must have success metrics      |
| Scope                   | ✓        | Must have both in/out of scope |
| User Stories            | ✓        | At least one per target user   |
| Functional Requirements | ✓        | At least one FR                |
| Acceptance Criteria     | ✓        | Must be testable checkboxes    |
| Non-Functional Req.     | ○        | Optional                       |
| Dependencies            | ○        | Optional                       |
| Risks                   | ○        | Optional                       |

## Dimension 1: User Stories

Checks: coverage (one story per target user), format (As a/I want/So that), specificity (concrete action), AC per story (Given/When/Then).

- A: All stories present, correct format, specific, AC attached
- B: All present, minor format issues or 1 missing AC
- C: Stories vague, or AC missing on most
- F: No user stories, or only one user covered when multiple exist

## Dimension 2: Acceptance Criteria

Checks: testable (verifiable by human or test), specific (no vague terms: fast/easy/good), complete (covers all stories + FRs), format (- [ ]).

- A: All testable, specific, complete, checkbox format
- B: Most testable, 1-2 vague items
- C: Many vague, or missing coverage for key features
- F: No AC, or AC is just a copy of requirements

## Dimension 3: Goals & Metrics

Checks: measurable (numeric targets), non-goals (explicitly listed), motivation (why, not just what).

- A: Metrics quantified, non-goals explicit, motivation clear
- B: Metrics defined but not all quantified
- C: Goals stated but no metrics, or no non-goals
- F: No goals section, or goals are purely technical

## Dimension 4: Scope Clarity

Checks: in-scope (concrete deliverables), out-of-scope (deferred items listed), consistency (aligns with user stories).

- A: Both in/out defined, items concrete, consistent with stories
- B: Both defined, minor vagueness
- C: Only in-scope defined, or items vague
- F: No scope section

## Dimension 5: Requirements Quality

Checks: priority (P0/P1/P2 per FR), traceability (FRs trace to stories), completeness (no obvious gaps).

- A: All FRs prioritized, traceable, no gaps
- B: Most prioritized, minor gaps
- C: No priorities, or FRs don't map to stories
- F: No functional requirements

## Overall Grade

| Grade | Condition                                    |
|-------|----------------------------------------------|
| A     | All 5 dimensions A/B, at least 3 A's         |
| B     | No F, max 1 C                                |
| C     | 1 F or 2+ C's                                |
| D     | 2 F's                                        |
| F     | 3+ F's or User Stories missing entirely      |
```

---

## Step 3: Report to User

After the agent completes, relay its summary to the user: overall grade, top issues, and next step recommendation.

## Related

- `/write-prd` — Create or revise the PRD
- `/breakdown-tasks` — Next step after PRD passes evaluation
