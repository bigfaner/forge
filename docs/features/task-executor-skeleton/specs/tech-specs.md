---
feature: "task-executor-skeleton"
generated: "2026-05-10"
status: draft
---

# Technical Specifications: task-executor-skeleton

## Step Architecture

### TECH-001: Step Renumbering (6 to 5)

**Requirement**: Merge the current 6-step execution into 5 steps. Old Steps 2 (TDD) and 3 (Quality Gate) merge into new Step 2 (Execute Workflow). Old Steps 4 (Record) and 5 (Commit) become new Steps 3 and 4. Step 0 (Claim) and Step 1 (Read task) remain unchanged.
**Scope**: [LOCAL]
**Source**: tech-design Architecture - Step Renumbering

### TECH-002: Skeleton Step 2 (Zero Hardcoded TDD)

**Requirement**: The new Step 2 in `task-executor.md` contains zero hardcoded TDD logic. It implements a 3-case detection procedure: (A) workflow heading + non-empty content = follow it, (B) no workflow heading = read default template at `plugins/forge/skills/breakdown-tasks/templates/task.md`, (C) workflow heading + empty content = log warning + Case B. Output format: `Step 2/4: [workflow description]... DONE` or `Step 2/4: [workflow description]... FAILED: [reason]`.
**Scope**: [LOCAL]
**Source**: tech-design Interface 1

### TECH-003: Workflow Content Model (W1-W5)

**Requirement**: Every `## Execution Workflow` section MUST satisfy 5 validation rules: W1 (steps are numbered and sequential), W2 (each step names a concrete command or action), W3 (each step states success criteria), W4 (failure handling is explicit per step), W5 (terminal step is a stop condition, no loops).
**Scope**: [CROSS]
**Source**: tech-design Interface 2 - Workflow Content Model

### TECH-004: Default Template Fallback Path

**Requirement**: When no `## Execution Workflow` is found in the task file, task-executor reads `plugins/forge/skills/breakdown-tasks/templates/task.md` and follows its workflow. This file MUST contain `## Execution Workflow` with the standard TDD + Quality Gate steps.
**Scope**: [LOCAL]
**Source**: tech-design Interface 1 Case B, Component Diagram

## Data Model Changes

### TECH-005: Remove NoTest from Task and TaskState Structs

**Requirement**: Delete `NoTest bool \`json:"noTest,omitempty"\`` from both `Task` struct (~line 32) and `TaskState` struct (~line 120) in `pkg/task/types.go`. No new fields are added.
**Scope**: [LOCAL]
**Source**: tech-design Data Models - Removed NoTest Field, Interface 3

### TECH-006: Claim Output Simplification

**Requirement**: In `internal/cmd/claim.go`, remove `PrintField("NO_TEST", ...)` from `printTaskDetails` and remove `NoTest: t.NoTest` from the state bootstrap in `executeClaim`. Downstream consumers (run-tasks.md, execute-task.md) stop parsing `NO_TEST` from claim output.
**Scope**: [LOCAL]
**Source**: tech-design Interface 3

### TECH-007: Record Behavior Unification

**Requirement**: In `internal/cmd/record.go`, delete 3 NoTest-related code blocks: (1) the coverage auto-set when `NoTest && no tests`, (2) the `!t.NoTest` condition from quality gate check (now runs for all tasks), (3) the `NoTest` parameter from `formatTestsExecuted`. In `internal/cmd/errors.go`, update `ErrNoTestEvidence` hint to remove the noTest option: "Either (1) run tests and report results, or (2) use --force to override".
**Scope**: [LOCAL]
**Source**: tech-design Interface 4

## Error Handling

### TECH-008: Structured Error Handoff (Agent to task-cli)

**Requirement**: The agent-to-task-cli error handoff works via structured CLI arguments, NOT freeform text parsing. When a workflow fails, the agent calls `task record --status failed --notes "reason"`. The step output string ("Step 2/4: ... FAILED: reason") is human-readable log only; structured data flows through CLI args.
**Scope**: [CROSS]
**Source**: tech-design Error Handling - Agent-to-task-cli Error Translation

### TECH-009: Error Propagation Channels

**Requirement**: Errors propagate through two channels: (1) task-cli Go code uses `*AIError` structs with Code/Message/Cause/Hint/Action fields, printed to stderr as `ERROR_CODE: ...`, (2) agent prompts use step output strings as human-readable log. No error codes cross the Go/agent boundary.
**Scope**: [CROSS]
**Source**: tech-design Error Handling - Error Propagation

### TECH-010: Task File Validation in Step 1

**Requirement**: Step 1 validates task file readability and frontmatter parseability. If either fails, status = `failed`, error is logged, Step 2 is skipped. This is the first gate before workflow execution begins.
**Scope**: [LOCAL]
**Source**: tech-design Error Handling - Error Cases (row 1)

## Template and Schema Changes

### TECH-011: Template Workflow Section Format

**Requirement**: Each `## Execution Workflow` section in task templates contains numbered sequential steps with explicit success criteria and failure handling per step, ending with a stop condition. The default template (`task.md`) contains the TDD + Quality Gate workflow. No template should have an open-ended workflow.
**Scope**: [LOCAL]
**Source**: tech-design Interface 2

### TECH-012: Schema Field Removal

**Requirement**: Remove `noTest` field definition from both `breakdown-tasks/templates/index.schema.json` and `quick-tasks/templates/index.schema.json`. No new schema fields are added.
**Scope**: [LOCAL]
**Source**: tech-design Appendix - Schemas

### TECH-013: File Change Inventory (29 files)

**Requirement**: The feature modifies ~29 files across 6 categories: (1) 3 agent prompts (task-executor.md, run-tasks.md, execute-task.md), (2) 10 breakdown task templates, (3) 6 quick task templates, (4) 2 schemas, (5) 4 Go source files + test files, (6) 3 skill docs.
**Scope**: [LOCAL]
**Source**: tech-design Appendix - File Change Inventory

## Testing

### TECH-014: Grep Verification for noTest Removal

**Requirement**: After all changes, `grep -ri noTest` across `.go`, `.md`, and `.json` files MUST return zero matches. This is a binary pass/fail gate.
**Scope**: [LOCAL]
**Source**: tech-design Testing Strategy - Per-Layer Test Plan

### TECH-015: Agent Prompt Test Scenarios

**Requirement**: 8 manual test scenarios covering: (1) record with no evidence, (2) record --force, (3) record with evidence, (4) task without workflow heading, (5) task with workflow heading, (6) T-test-3 execution timing, (7) grep verification, (8) empty workflow heading. Each has explicit pass/fail criteria.
**Scope**: [LOCAL]
**Source**: tech-design Testing Strategy - Key Test Scenarios
