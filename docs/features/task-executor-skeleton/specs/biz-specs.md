---
feature: "task-executor-skeleton"
generated: "2026-05-10"
status: draft
---

# Business Rules: task-executor-skeleton

## Task Execution Dispatch

### BIZ-001: Execution Workflow Detection

**Rule**: task-executor MUST detect `## Execution Workflow` heading in the task file. If the heading exists with non-empty content, inject that content into Step 2 instructions, replacing hardcoded TDD logic. If absent, fallback to TDD + Quality Gate via the default template (`task.md`).
**Context**: Replaces the `noTest` boolean flag which was semantically ambiguous ("don't do TDD" vs "no tests needed"). Template-driven workflow gives each task type explicit control over its execution steps.
**Scope**: [LOCAL]
**Source**: prd-spec.md Flow Description, tech-design Interface 1

### BIZ-002: Empty Workflow Configuration Error

**Rule**: If `## Execution Workflow` heading exists but the content is empty, task-executor MUST log a warning ("WARNING: ## Execution Workflow heading present but empty. Falling back to default template.") and fallback to the default template.
**Context**: Prevents silent misconfiguration. An empty workflow section likely indicates an incomplete template edit.
**Scope**: [LOCAL]
**Source**: prd-spec.md Flow Description, tech-design Interface 1 Case C

### BIZ-003: Terminal States

**Rule**: Task execution has exactly two terminal states: `completed` (success, proceeds to record + commit) and `failed` (any failure path). There is no `partial` state; partial completion is classified as `failed` with an attached summary of completed steps.
**Context**: Simplifies status tracking and prevents ambiguous intermediate states.
**Scope**: [CROSS]
**Source**: prd-spec.md Flow Description - Terminal States

### BIZ-004: No TDD Retry on Execution Failure

**Rule**: When a workflow step fails, the agent MUST NOT fall back to a TDD loop. If the workflow declares failure handling instructions (e.g., "create fix task"), the agent follows those instructions. If no explicit failure instructions exist, the agent records the failure reason and stops.
**Context**: Eliminates the 14-minute retry loop observed in execution tasks like T-test-3 where the TDD cycle was irrelevant and wasteful.
**Scope**: [LOCAL]
**Source**: prd-spec.md Flow Description - Failure Path

## Validation and Quality

### BIZ-005: Workflow Must Have Explicit Stop Condition

**Rule**: Every `## Execution Workflow` template MUST include a terminal stop condition (e.g., "6. Stop. Proceed to Step 3."). Open-ended instructions (e.g., "Repeat until all tests pass") are prohibited.
**Context**: Prevents agents from entering infinite loops or interpreting vague instructions inconsistently. Enforced by the Workflow Content Model checklist (W5).
**Scope**: [LOCAL]
**Source**: prd-spec.md Quality Requirements, tech-design Interface 2

### BIZ-006: noTest Complete Removal

**Rule**: The `noTest`/`NO_TEST` field must be completely removed from all code, prompts, templates, schemas, and documentation. Verification: `grep -ri noTest` returns zero matches across `.go`, `.md`, and `.json` files.
**Context**: The `noTest` flag was a workaround that conflated "no TDD cycle" with "no tests exist". With template-driven workflows, each template declares its own steps, making the boolean flag obsolete.
**Scope**: [LOCAL]
**Source**: prd-spec.md Goals, Scope

### BIZ-007: Backward Compatibility Fallback

**Rule**: Task files without `## Execution Workflow` MUST automatically fallback to TDD + Quality Gate behavior via the default template. The fallback is task-executor's built-in behavior and requires no template declaration.
**Context**: Ensures existing task files (created before this feature) continue to work without modification.
**Scope**: [CROSS]
**Source**: prd-spec.md Compatibility Requirements, tech-design Interface 1 Case B

## Failure Handling

### BIZ-008: Task File Unreadable

**Rule**: If the task file is missing or has corrupted frontmatter (cannot parse), task-executor sets status to `failed`, logs the error, and skips Step 2 entirely. A minimal error record is committed.
**Context**: Prevents agents from proceeding with incomplete task definitions.
**Scope**: [LOCAL]
**Source**: prd-spec.md Flow Description - Failure Path

### BIZ-009: Agent Timeout

**Rule**: If the agent exceeds the built-in timeout threshold during workflow execution, execution is forcibly terminated, status is set to `failed`, and timeout information is recorded.
**Context**: Prevents runaway agent sessions consuming resources.
**Scope**: [CROSS]
**Source**: prd-spec.md Flow Description - Failure Path

### BIZ-010: Quality Gate Safety Net

**Rule**: `task record` in task-cli runs the quality gate pre-check for ALL tasks uniformly. The `noTest` bypass is removed. Documentation-only tasks use `task record --force` or include a lightweight validation step (e.g., `just compile`) in their workflow.
**Context**: Maintains code quality enforcement even after removing the `noTest` escape hatch.
**Scope**: [CROSS]
**Source**: tech-design Interface 4
