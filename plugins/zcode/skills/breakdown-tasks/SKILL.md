---
name: breakdown-tasks
description: Use when design.md is finalized to break down into executable tasks. Creates task files based on technical design.
---

# Breakdown Tasks

## Overview

从技术设计文档拆解成可执行的任务。

**核心原则**：任务粒度适中（1-4 小时），依赖关系明确，验收标准可测试。

## Prerequisites

检查上一阶段产物，缺失则中止并提示用户：

```bash
ls docs/features/<slug>/design/tech-design.md
```

| 产物 | 缺失时提示 |
|------|-----------|
| `design/tech-design.md` | 先执行 `/design-tech`，再执行 `/eval-design` |

## Directory Structure

```
docs/features/<feature-slug>/
├── manifest.md                    # Feature index & traceability
├── prd/
│   ├── prd-spec.md
│   ├── prd-user-stories.md
│   └── prd-ui-functions.md
├── design/
│   ├── tech-design.md             # Technical design (input)
│   └── api-handbook.md
├── ui/
│   └── ui-design.md               # (if applicable)
├── tasks/
│   ├── index.json                 # Task index
│   ├── 1.1-<title>.md            # Task detail files
│   ├── process/
│   └── records/
```

## When to Use

**Trigger conditions:**

- Manifest exists at `docs/features/<slug>/manifest.md` with status `design` or `tasks`
- Design document exists at `design/tech-design.md`
- User asks to "break down" or "split" a design into tasks

**Skip when:**

- No design.md exists (use `/design-tech` first)
- Tasks already exist for the feature

## Workflow

```
1. Read Design → 2. Map interfaces → 3. Define order → 4. Create task files → 5. Create index.json → 6. Validate
```

## Step 1: Read Manifest → All Documents

1. Read `manifest.md` to locate all documents
2. Read `prd/prd-spec.md` — understand WHAT
3. Read `design/tech-design.md` — understand HOW
4. Read `design/api-handbook.md` — understand interfaces (if exists)
5. Read `ui/ui-design.md` — understand UI components (if exists)
6. Read `prd/prd-user-stories.md` — understand user scenarios (if exists)

## Step 2: Map Interfaces to Tasks

| Design Element       | Task Type           |
| -------------------- | ------------------- |
| Interface definition | Interface task      |
| Data model           | Model task          |
| Component            | Implementation task |
| Error type           | Error handling task |

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

Create `docs/features/<slug>/tasks/index.json` with task definitions.

**Schema reference:** [templates/index.schema.json](templates/index.schema.json)

### Required Fields

| Field         | Type   | Description                      |
| ------------- | ------ | -------------------------------- |
| `version`     | string | Semver format (e.g., `1.0.0`)    |
| `lastUpdated` | date   | ISO date (e.g., `2026-04-06`)    |
| `tasks`       | object | Map of task ID → task definition |

### Task Fields

| Field           | Required | Type   | Description                                                     |
| --------------- | -------- | ------ | --------------------------------------------------------------- |
| `id`            | ✓        | string | Task identifier (e.g., `1.1`)                                   |
| `phase`         | ✓        | int    | Phase number (≥1)                                               |
| `title`         | ✓        | string | Task title                                                      |
| `priority`      | ✓        | enum   | `P0` / `P1` / `P2`                                              |
| `status`        | ✓        | enum   | `pending` / `in_progress` / `completed` / `blocked` / `skipped` |
| `file`          | ✓        | string | Task file path                                                  |
| `dependencies`  |          | array  | Task IDs this depends on                                        |
| `estimatedTime` |          | string | Time estimate                                                   |
| `record`        |          | string | Record file path                                                |

## Step 6: Validate

```bash
task validate -file docs/features/\<slug\>/tasks/index.json
```

## Step 7: Update Manifest

Update `manifest.md`:
- Fill Tasks column in Traceability table with task IDs linked to design sections
- Advance status to `tasks`

## Integration

Works well with skills:

- `/design-tech` - Creates the design.md input
- `/eval-design` - Evaluate design.md before breakdown (recommended gate)
- `/claim-task` - Starts working on tasks
- `/record-task` - Records task completion
