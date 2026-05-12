---
feature: "typed-task-dispatch"
generated: "2026-05-12"
status: draft
---

# Business Rules: typed-task-dispatch

## Task Type System

### BIZ-001: Type field requirement

**Rule**: All tasks in index.json must include a `type` field with one of 11 valid enum values after migration
**Context**: Replaces the `noTest` and `mainSession` patch fields with an explicit type system
**Scope**: [LOCAL] - Specific to this feature's type system implementation
**Source**: prd-spec.md §Scope "task-cli 新增 task migrate 命令"

### BIZ-002: Task migrate precondition

**Rule**: The `task migrate` command refuses to run if any tasks have status `in_progress`
**Context**: Prevents data corruption during type field migration
**Scope**: [LOCAL] - Migration-specific validation rule
**Source**: prd-spec.md §Functional Specs "task migrate 命令规格"

### BIZ-003: Blocked state trigger

**Rule**: When `task prompt <id>` exits with non-zero code, the task is marked as blocked with blockedReason set to the stderr output
**Context**: Prevents silent failures; task execution failures are now visible
**Scope**: [LOCAL] - Specific to this feature's error handling pattern
**Source**: prd-spec.md §Flow Description "数据流表"

## Task Routing Rules

### BIZ-004: Eval-cases main session exception

**Rule**: Tasks with `type == test-pipeline.eval-cases` execute in the main Claude Code session instead of being dispatched to a subagent
**Context**: Platform limitation requires spawning subagents (doc-scorer, doc-reviser) from main session
**Scope**: [LOCAL] - Specific to eval-cases workflow
**Source**: prd-spec.md §Scope "type == test-pipeline.eval-cases（永久例外，平台限制）"

## Workflow Rules

### BIZ-005: Phase boundary detection

**Rule**: The first task of a new phase (phase number > max completed phase) receives a PHASE_SUMMARY placeholder in its prompt pointing to the previous phase's summary record
**Context**: Provides context from previous phases to improve decision-making
**Scope**: [LOCAL] - Specific to phased workflow execution
**Source**: prd-spec.md §Flow Description "task prompt 内部流程"
