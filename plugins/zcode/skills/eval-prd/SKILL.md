---
name: eval-prd
description: Evaluate a PRD document against quality standards. Checks structure completeness, user story quality, flow diagrams, functional specs, and scope clarity. Outputs a scored report with actionable improvements.
---

# Eval PRD

评估 PRD 文档是否满足规范，输出评分报告和改进建议。

## Prerequisites

检查上一阶段产物，缺失则中止并提示用户：

```bash
ls docs/features/<slug>/prd/prd-spec.md docs/features/<slug>/prd/prd-user-stories.md
```

| 产物 | 缺失时提示 |
|------|-----------|
| `prd/prd-spec.md` | 先执行 `/write-prd` |
| `prd/prd-user-stories.md` | 先执行 `/write-prd` |

## When to Use

**Trigger:**
- User asks to "evaluate PRD" or "check PRD quality"
- User provides `/eval-prd` command
- Before handing off PRD to `/design-tech` or `/ui-design`

**Skip:**
- PRD doesn't exist yet (use `/write-prd` first)

## Workflow

```
1. 定位 PRD + User Stories → 2. 启动评估 Agent → 3. 汇报结果
```

## Step 1: Locate Documents

Check in order:
1. Path provided by user
2. Read `docs/features/<current-feature>/manifest.md` → locate PRD documents
3. Fall back to `docs/features/<current-feature>/prd/prd-spec.md` + `prd/prd-user-stories.md`
4. Ask user for path if not found

Determine `<feature-slug>` from the path (e.g. `docs/features/auth-flow/prd.md` → slug is `auth-flow`).

## Step 2: Launch Evaluation Agent

Use the **Agent tool** to spawn a subagent. Pass the full prompt below, substituting `{{PRD_PATH}}`, `{{USER_STORIES_PATH}}`, and `{{FEATURE_SLUG}}`:

---

**Agent prompt template:**

```
You are a PRD quality evaluator. Your job: read the PRD and User Stories, apply the rubric, write the report, return a summary.

## Inputs
- PRD path: {{PRD_PATH}} (default: prd/prd-spec.md)
- User Stories path: {{USER_STORIES_PATH}} (default: prd/prd-user-stories.md)
- UI Functions path: {{UI_FUNCTIONS_PATH}} (optional: prd/prd-ui-functions.md)
- Feature slug: {{FEATURE_SLUG}}
- Report output: docs/features/{{FEATURE_SLUG}}/prd-eval.md
- Report template: plugins/zcode/skills/eval-prd/templates/report.md

## Steps
1. Read {{PRD_PATH}}
2. Read {{USER_STORIES_PATH}} (if exists)
3. Read {{UI_FUNCTIONS_PATH}} (if exists)
4. Read the report template
5. Apply the rubric below to every dimension
6. Fill in the template and write to docs/features/{{FEATURE_SLUG}}/prd-eval.md
7. Return: overall grade, top 2-3 issues, and whether it can proceed to /design-tech (or /ui-design if prd-ui-functions.md exists)

## Structure Check

Required sections in prd-spec.md — mark missing as F:

| Section | Required | Notes |
|---------|----------|-------|
| 需求背景（原因/对象/人员） | ✓ | 必须包含三个维度 |
| 需求目标 | ✓ | 必须包含量化指标 |
| Scope（In/Out） | ✓ | 两者都必须有 |
| 流程说明 + 业务流程图 | ✓ | Mermaid 流程图必填 |
| 功能描述 | ✓ | 至少包含列表页/按钮/表单之一 |
| 其他说明 | ○ | 可选但建议有 |
| 质量检查 | ○ | 可选 |

User Stories file:

| Section | Required | Notes |
|---------|----------|-------|
| User Stories (独立文件) | ✓ | 至少每个目标用户一个故事 |
| Acceptance Criteria | ✓ | Given/When/Then 格式 |

## Dimension 1: 背景与目标

Checks: 背景三要素（原因、对象、人员），目标量化，背景与目标逻辑一致。

- A: 背景含三要素，目标量化，逻辑一致
- B: 背景缺一个要素，或目标部分量化
- C: 背景模糊，目标无量化的
- F: 无背景或无目标

## Dimension 2: 流程说明

Checks: 流程图存在（Mermaid），主流程完整，决策点明确，异常分支覆盖。

- A: 流程图完整，含主流程+决策点+异常分支
- B: 流程图存在，缺异常分支或部分决策点
- C: 仅文字描述，无流程图
- F: 无流程说明

## Dimension 3: 功能描述

Checks: 表格完整性（列表页7要素、按钮4要素、表单2要素），字段说明清晰，校验规则明确。

- A: 所有表格完整填写，字段和校验规则清晰
- B: 大部分完整，1-2处缺失
- C: 表格存在但内容不完整
- F: 无功能描述或仅有文字无表格

## Dimension 4: User Stories

Checks: coverage (one story per target user), format (As a/I want/So that), specificity (concrete action), AC per story (Given/When/Then).

- A: All stories present, correct format, specific, AC attached
- B: All present, minor format issues or 1 missing AC
- C: Stories vague, or AC missing on most
- F: No user stories, or only one user covered when multiple exist

## Dimension 5: Scope Clarity

Checks: in-scope (concrete deliverables), out-of-scope (deferred items listed), consistency (aligns with 功能描述 and user stories).

- A: Both in/out defined, items concrete, consistent with 功能描述
- B: Both defined, minor vagueness
- C: Only in-scope defined, or items vague
- F: No scope section

## Dimension 6: UI Functions (optional)

Only checked if `prd/prd-ui-functions.md` exists.

Checks: each UI function has description, interaction flow, data requirements, states, validation.

- A: All functions fully specified with all sub-sections
- B: Most specified, 1-2 missing sub-sections
- C: Functions listed but incomplete
- N/A: File doesn't exist (not an F)

## Overall Grade

| Grade | Condition |
|-------|-----------|
| A | All 5 dimensions A/B, at least 3 A's |
| B | No F, max 1 C |
| C | 1 F or 2+ C's |
| D | 2 F's |
| F | 3+ F's or User Stories missing entirely |
```

---

## Step 3: Report to User

After the agent completes, relay its summary to the user: overall grade, top issues, and next step recommendation.

## Related

- `/write-prd` — Create or revise the PRD
- `/design-tech` — Next step: produce technical design document (architecture, interfaces, data model)
- `/ui-design` — Next step (optional): produce UI design spec, if `prd-ui-functions.md` exists
- `/breakdown-tasks` — After design docs are finalized, break design into executable tasks
