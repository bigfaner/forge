---
feature: "typed-task-dispatch"
generated: "2026-05-12"
status: draft
---

# Technical Specifications: typed-task-dispatch

## CLI Commands

### TECH-001: task prompt command contract

**Requirement**: The `task prompt <id>` command outputs synthesized prompt to stdout (UTF-8, no prefix/suffix), errors to stderr, exit code 0 on success, 1 on failure
**Scope**: [LOCAL] - Specific to this feature's CLI implementation
**Source**: tech-design.md §Interface 3

### TECH-002: task prompt timing requirement

**Requirement**: The `task prompt <id>` command must complete in < 500ms (local file reads + string substitution, no network calls)
**Scope**: [LOCAL] - Performance requirement for this command
**Source**: tech-design.md §Interface 3

### TECH-003: task migrate command contract

**Requirement**: The `task migrate` command reads current feature's index.json, outputs summary line to stdout, errors to stderr, exits 0 on success, 1 if in_progress tasks exist or file I/O fails
**Scope**: [LOCAL] - Specific to this feature's migration command
**Source**: tech-design.md §Interface 4

## Architecture Patterns

### TECH-004: Two-layer model separation

**Requirement**: Task execution separates constraint layer (task-executor.md agent definition: ONE TASK, record-task mandatory, no background tasks, max 3 subagent calls) from strategy layer (injected via `task prompt <id>` stdout)
**Scope**: [LOCAL] - Architecture decision specific to this feature
**Source**: tech-design.md §Overview "两层模型"

### TECH-005: Go embed templates

**Requirement**: Prompt templates are stored as Go embed strings in `pkg/prompt/data/*.md`, not as separate markdown files
**Scope**: [LOCAL] - Implementation choice for template storage
**Source**: tech-design.md §New Package "pkg/prompt"

## Data Models

### TECH-006: Type field enum values

**Requirement**: Task `type` field accepts exactly 11 enum values: implementation, doc-generation.summary, doc-generation.consolidate, test-pipeline.gen-cases, test-pipeline.eval-cases, test-pipeline.gen-scripts, test-pipeline.run, test-pipeline.graduate, test-pipeline.verify-regression, fix, gate
**Scope**: [LOCAL] - Type system specific to this feature
**Source**: tech-design.md §Data Models "Type enum constants"

### TECH-007: BlockedReason field

**Requirement**: Task struct includes optional `blockedReason` field (string) written by run-tasks when task prompt fails
**Scope**: [LOCAL] - Error handling field for this feature
**Source**: tech-design.md §Data Models "Task struct additions"

## Type Inference Rules

### TECH-008: Task migrate type inference

**Requirement**: The `task migrate` command infers task type from ID patterns: .summary → doc-generation.summary, .gate → gate, T-test-* → test-pipeline.*, fix-/disc- prefix → fix, default → implementation
**Scope**: [LOCAL] - Migration logic specific to this feature
**Source**: prd-spec.md §Functional Specs "task migrate 命令规格"

## String Substitution

### TECH-009: Template placeholder variables

**Requirement**: Prompt templates use placeholders: {{TASK_ID}}, {{TASK_FILE}}, {{SCOPE}}, {{NO_TEST}}, {{PHASE_SUMMARY}}, {{FEATURE_SLUG}}
**Scope**: [LOCAL] - Template variable naming specific to this feature
**Source**: tech-design.md §New Package "pkg/prompt"

## Error Handling

### TECH-010: State.json fallback to git branch

**Requirement**: When `.forge/state.json` is missing or unreadable, the system falls back to git branch name to determine the current feature
**Scope**: [LOCAL] - Error handling pattern specific to task-cli
**Source**: pkg/feature/feature.go (implementation detail)

## Phase Boundary Detection

### TECH-011: Phase summary injection logic

**Requirement**: Phase boundary detection scans completed tasks to find max phase number; if current task phase > max completed and phase > 1, inject PHASE_SUMMARY placeholder
**Scope**: [LOCAL] - Workflow-specific logic
**Source**: prd-spec.md §Flow Description "task prompt 内部流程"
