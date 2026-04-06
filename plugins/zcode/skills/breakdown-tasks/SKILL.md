---
name: breakdown-tasks
description: Use when design.md is finalized to break down into executable tasks. Creates task files based on technical design.
---

# Breakdown Tasks

## Overview

从技术设计文档拆解成可执行的任务。

**核心原则**：任务粒度适中（1-4 小时），依赖关系明确，验收标准可测试。

## Position in Workflow

```
/write-prd → /design-tech → /breakdown-tasks
     ↓              ↓              ↓
   prd.md      design.md      tasks/*.md
```

## Directory Structure

```
docs/features/<feature-slug>/
├── prd.md              # PRD document
├── design.md           # Technical design (input)
├── index.json          # Task index for this feature
├── tasks/              # Task definitions
└── records/            # Execution records
```

## When to Use

**Trigger conditions:**
- Design document exists at `docs/features/<slug>/design.md`
- User asks to "break down" or "split" a design into tasks

**Skip when:**
- No design.md exists (use `/design-tech` first)
- Tasks already exist for the feature

## Workflow

```
1. Read Design → 2. Map interfaces → 3. Define order → 4. Create task files → 5. Create index.json → 6. Validate
```

## Step 1: Read Design

Read `docs/features/<slug>/design.md`:
- Understand architecture and component structure
- Map interfaces to implementation tasks
- Identify data models and their tasks

## Step 2: Map Interfaces to Tasks

| Design Element | Task Type |
|----------------|-----------|
| Interface definition | Interface task |
| Data model | Model task |
| Component | Implementation task |
| Error type | Error handling task |

## Step 3: Define Task Order

```
1.x Interfaces → 2.x Models → 3.x Implementation → 4.x Integration → 5.x Tests
```

## Step 4: Create Task Files

**Naming convention:**
```
<sequence>.<sub-sequence>-<slug>.md
```

## Step 5: Create index.json

Create `docs/features/<slug>/index.json` with task definitions.

## Step 6: Validate

```bash
task validate -file docs/features/<slug>/index.json
```

## Integration

Works well with:
- `/design-tech` - Creates the design.md input
- `/claim-task` - Starts working on tasks
- `/record-task` - Records task completion
