---
name: write-prd
description: Use when user provides requirements or feature requests that need to be formalized into a structured PRD document through collaborative dialogue.
---

# Write PRD

## Overview

从模糊需求产出清晰的 PRD（产品需求文档），通过协作对话逐步澄清需求。

**核心原则**：在编码前先澄清 "做什么" 和 "为什么做"，避免方向性错误。

## Prerequisites

无强制前置产物。若有 brainstorm 提案，作为可选输入：

```bash
ls docs/proposals/<slug>/proposal.md 2>/dev/null  # 可选，不阻塞
```

<HARD-GATE>
Do NOT write any code, scaffold any project, or take any implementation action until the PRD is finalized and approved. Present the PRD and get user approval first.
</HARD-GATE>

<HARD-RULE>
**禁止技术选型，允许技术约束**：

- **允许**：描述非功能性约束——性能要求（响应时间、并发量）、平台要求（浏览器、移动端）、兼容性、安全合规等。这些是业务级需求。
- **禁止**：提及具体技术栈——框架名称、编程语言、数据库、库、中间件、架构模式（如微服务、事件驱动）等。这些是技术选型，留给 `/tech-design` 阶段。

**判断标准**：如果描述的是"需要达到什么效果"→ 允许；如果描述的是"用什么工具实现"→ 禁止。
</HARD-RULE>

## When to Use

**Trigger conditions:**

- User describes a feature/requirement without clear specifications
- User says "I want to..." or "We need..." without details
- Starting a new phase or major feature

**Skip when:**

- Clear task definitions already exist
- Simple bug fix or small tweak

## Process Flow

```
Explore context → Check proposal → Assess scope → Ask questions → Propose approaches → Present PRD sections → Write PRD Spec + User Stories + UI Functions → Create Manifest → Commit
```

## Checklist

1. **Explore project context** — check files, docs, recent commits
2. **Check for existing proposal** — read `docs/proposals/<slug>/proposal.md` if it exists
3. **Assess scope** — determine if request needs decomposition
4. **Ask clarifying questions** — one at a time via AskUserQuestion tool
5. **Propose 2-3 approaches** — with trade-offs and your recommendation
6. **Present PRD sections** — get approval after each section
7. **Write PRD Spec** — save to `docs/features/<feature-slug>/prd/prd-spec.md`
8. **Write User Stories** — save to `docs/features/<feature-slug>/prd/prd-user-stories.md`
9. **Write UI Functions** (if applicable) — save to `docs/features/<feature-slug>/prd/prd-ui-functions.md`
10. **Create Manifest** — save to `docs/features/<feature-slug>/manifest.md`
11. **Commit** — commit all documents

## Output Documents

PRD 完成后输出以下文件：

| 文件 | 模板 | 说明 |
|------|------|------|
| `prd/prd-spec.md` | `templates/prd-spec.md` | 产品需求文档，包含背景、目标、Scope、流程、功能描述等 |
| `prd/prd-user-stories.md` | `templates/prd-user-stories.md` | 用户故事，从 PRD 背景中识别的用户角色推导而出 |
| `prd/prd-ui-functions.md` | `templates/prd-ui-functions.md` | UI 功能要点（需求层，仅适用于有 UI 表面的功能） |
| `manifest.md` | `templates/manifest.md` | Feature 索引和可追溯性映射 |

## Step 1: Explore Project Context

Before asking questions, understand the current state:

- Check `docs/proposals/<slug>/proposal.md` if a proposal exists — carry forward business context
- Check `docs/features/<slug>/tasks/index.json` for related tasks
- Review recent git commits for related work

**禁止**：不得读取 `ARCHITECTURE.md`、`DECISIONS.md` 等技术文档来引导需求讨论。技术约束不属于 PRD 范畴。

## Step 2: Assess Scope

Evaluate if the request is appropriately scoped:

- If request describes multiple independent subsystems → **Decompose first**
- If single focused feature → **Proceed with questions**

## Step 3: Ask Clarifying Questions

**CRITICAL**: Use `AskUserQuestion` tool for ALL questions.

### Question Guidelines

- **One question at a time** — never batch questions
- **Prefer multiple choice** — easier to answer than open-ended
- **Focus on understanding**: user roles, purpose, constraints, success criteria
- **Go back when needed** — if something doesn't make sense, clarify

See `examples/ask-questions.md` for concrete examples.

## Step 4: Propose Approaches

After understanding requirements, propose 2-3 **business approaches** (not technical implementations):

1. **Present options conversationally** with your recommendation
2. **Lead with your recommended option** and explain why
3. **Include trade-offs** for each approach (business impact, user experience, scope)

**禁止**：方案中不得涉及具体技术选型。方案应描述不同的业务功能组合或用户流程，而非技术实现路径。但可以提及非功能性约束（如性能要求、平台要求、安全合规）。

See `examples/propose-approaches.md` for structure and tips.

## Step 5: Present PRD Sections

Present incrementally, getting approval after each section:

| Section | Content | Key Points |
|---------|---------|------------|
| 需求背景 | 原因、对象、人员 | 必须包含三个维度 |
| 需求目标 | 目标 + 量化指标 | 尽可能量化收益 |
| Scope | In Scope / Out of Scope | 明确边界 |
| 流程说明 | 业务流程 + Mermaid 流程图 | 流程图必填 |
| 功能描述 | 列表页 / 按钮 / 表单 / 关联改动 | 快速/详细模式按需选择，表格必填 |
| 其他说明 | 性能 / 数据 / 监控 / 安全 | 非功能性需求 |
| User Stories | As a / I want / So that + AC | 输出到独立文件 |

## Step 6: Write PRD Spec

使用 `templates/prd-spec.md` 模板填写。

**目录结构：**

```
docs/features/<feature-slug>/
├── manifest.md                # Feature index & traceability
├── prd/
│   ├── prd-spec.md            # PRD Spec
│   ├── prd-user-stories.md    # 用户故事
│   └── prd-ui-functions.md    # UI 功能要点（可选）
├── design/                    # (created by /tech-design)
├── ui/                        # (created by /ui-design)
└── tasks/                     # (created by /breakdown-tasks)
    └── records/
```

## Step 7: Write User Stories

从 PRD 背景中识别的用户角色推导用户故事，输出到 `prd/prd-user-stories.md`。

```
As a [user role from Background]
I want to [specific action]
So that [concrete benefit/goal]
```

**Coverage rules:**
- Every user type from 需求背景 must have at least one story
- Actions must be concrete — not "manage" or "handle" but "create X", "filter by Y"

**Acceptance Criteria** (Given/When/Then) 必须跟随每个故事，每条 AC 必须可客观验证。

## Step 8: Write UI Functions (if applicable)

For features with UI surfaces, create `prd/prd-ui-functions.md` using `templates/prd-ui-functions.md`.
Skip this step for backend/API/CLI features with no UI surface.

## Step 9: Create Manifest

Create `manifest.md` at the feature root using `templates/manifest.md`:
- Fill in PRD entries and summaries
- Set status to `prd`
- Include UI Functions row only if `prd/prd-ui-functions.md` was created

## Step 9.5: Self-Check

Before presenting to the user, verify the PRD passes these checks:

| Check | What to verify |
|-------|----------------|
| Background completeness | 原因 + 对象 + 人员 all present and specific |
| Goals quantified | At least one numeric target (% , count, time) |
| Flow diagram | Mermaid flowchart with decision points (diamond nodes) and at least one error/exception branch |
| Functional specs | All applicable tables filled — no placeholder rows |
| User stories | One story per user role, each with Given/When/Then AC |
| Scope consistency | In-scope items match what's described in 功能描述 and user stories |
| No vague language | No "better", "faster", "improved" without quantification |

## Step 10: Review & Commit

<HARD-RULE>
Do NOT commit documents automatically. Present all generated documents to the user for review and wait for explicit approval before committing.
</HARD-RULE>

1. Present the full PRD spec, user stories, and UI functions (if any) to the user
2. Wait for the user to review and approve (or request changes)
3. Only commit after explicit user approval:

```bash
git add docs/features/<feature-slug>/
git commit -m "docs: add PRD for <feature-slug>"
```

## Step 11: Adversarial Eval Prompt

After committing, use `AskUserQuestion` to ask:

> 是否运行 `/eval-prd` 对 PRD 进行对抗性评估？（默认 80 分 / 3 轮）

- **Yes** → invoke `/eval-prd` via `Skill` tool
- **Custom** → invoke `/eval-prd --target X --iterations Y` via `Skill` tool
- **No** → proceed to `/tech-design`

## Integration

Works well with skills:

- `/eval-prd` - Evaluate PRD quality before proceeding to design phase
- `/tech-design` - After PRD passes evaluation, produce technical design document
- `/ui-design` - After PRD passes evaluation, produce UI design spec (if prd-ui-functions.md exists)
- `docs/decisions/` - Record key decisions during PRD creation (category-based directory)
