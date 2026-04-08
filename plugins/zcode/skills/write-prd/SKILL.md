---
name: write-prd
description: Use when user provides requirements or feature requests that need to be formalized into a structured PRD document through collaborative dialogue.
---

# Write PRD

## Overview

从模糊需求产出清晰的 PRD（产品需求文档），通过协作对话逐步澄清需求。

**核心原则**：在编码前先澄清 "做什么" 和 "为什么做"，避免方向性错误。

<HARD-GATE>
Do NOT write any code, scaffold any project, or take any implementation action until the PRD is finalized and approved. Present the PRD and get user approval first.
</HARD-GATE>

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
Explore context → Assess scope → Ask questions → Propose approaches → Present PRD sections → Write PRD + User Stories → Commit
```

## Checklist

1. **Explore project context** — check files, docs, recent commits
2. **Assess scope** — determine if request needs decomposition
3. **Ask clarifying questions** — one at a time via AskUserQuestion tool
4. **Propose 2-3 approaches** — with trade-offs and your recommendation
5. **Present PRD sections** — get approval after each section
6. **Write PRD document** — save to `docs/features/<feature-slug>/prd.md`
7. **Write User Stories** — save to `docs/features/<feature-slug>/user-stories.md`
8. **Commit** — commit both documents

## Output Documents

PRD 完成后输出两个独立文件：

| 文件 | 模板 | 说明 |
|------|------|------|
| `prd.md` | `templates/prd.md` | 产品需求文档，包含背景、目标、Scope、流程、功能描述等 |
| `user-stories.md` | `templates/user-stories.md` | 用户故事，从 PRD 背景中识别的用户角色推导而出 |

## Step 1: Explore Project Context

Before asking questions, understand the current state:

- Read `docs/ARCHITECTURE.md` for architecture constraints
- Read `docs/DECISIONS.md` for existing technical decisions
- Check `docs/features/<slug>/tasks/index.json` for related tasks
- Review recent git commits

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

After understanding requirements, propose 2-3 implementation approaches:

1. **Present options conversationally** with your recommendation
2. **Lead with your recommended option** and explain why
3. **Include trade-offs** for each approach

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

## Step 6: Write PRD Document

使用 `templates/prd.md` 模板填写。

**目录结构：**

```
docs/features/<feature-slug>/
├── prd.md              # PRD 文档
├── user-stories.md     # 用户故事（独立文件）
├── tasks/              # Task definitions (created by /breakdown-tasks)
└── records/            # Execution records (created by /record-task)
```

## Step 7: Write User Stories

从 PRD 背景中识别的用户角色推导用户故事，输出到 `user-stories.md`。

```
As a [user role from Background]
I want to [specific action]
So that [concrete benefit/goal]
```

**Acceptance Criteria** (Given/When/Then) 必须跟随每个故事。

See `examples/user-stories.md` for a concrete example.

## Integration

Works well with skills:

- `/eval-prd` - Evaluate PRD quality before handing off to breakdown-tasks
- `/breakdown-tasks` - After PRD passes evaluation, break into tasks
- `docs/DECISIONS.md` - Record key decisions during PRD creation
