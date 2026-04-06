---
name: simplify-skill
description: Refactor skill files by extracting templates/examples to separate files.
argument-hints: skill name
allowed_tools: ["Read", "AskUserQuestion"]
---

# /simplify-skill

重构 skill 文件：**拆分非核心内容**，保持主流程清晰。

## Core Principle

```
MAIN FILE = WORKFLOW ONLY

skill.md     → 流程步骤、决策点
templates/   → JSON 模板、输出格式示例
examples/    → 完整用例、边界情况
```

## Workflow

```
1. Identify Skill → 2. Analyze Content → 3. Ask Approval → 4. Extract
```

## Phase 1: Identify Target

If no argument provided, ask user which skill to refactor.

Target locations:
- Skills: `.claude/skills/<name>/SKILL.md`
- Commands: `.claude/commands/<name>.md`

## Phase 2: Analyze Extractables

| Category | Indicators | Extract To |
|----------|-----------|------------|
| Templates | JSON blocks, output formats | `templates/` |
| Examples | Multi-line code samples | `examples/` |
| Reference tables | Field definitions | `reference.md` |
| Verbose context | Background explanations | `context.md` |

## Phase 3: Ask for Approval

Use `AskUserQuestion` with multiSelect for which content to extract.

## Phase 4: Execute Extraction

1. Create directory structure
2. Extract content to new files
3. Replace in skill.md with reference
4. Keep workflow steps intact

## Iron Law

<EXTREMELY-IMPORTANT>
NEVER extract without user approval
NEVER remove content, only relocate
ALWAYS add file references
KEEP workflow steps intact
</EXTREMELY-IMPORTANT>
