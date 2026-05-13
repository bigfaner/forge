---
name: simplify-skill
description: Refactor skill files by extracting templates/examples to separate files.
argument-hints:
  - name: skill-name
    description: Name of the skill to simplify (e.g. brainstorm, write-prd)
    required: true
allowed_tools: ["Read", "Write", "Edit", "AskUserQuestion"]
---

# /simplify-skill

Refactor skill files: **extract non-core content** to keep the main workflow clear.

## Core Principle

```
MAIN FILE = WORKFLOW ONLY

skill.md     → Process steps, decision points
templates/   → JSON templates, output format examples
examples/    → Full use cases, edge cases
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
