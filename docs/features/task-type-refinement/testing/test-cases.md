---
feature: "task-type-refinement"
sources:
  - docs/proposals/task-type-refinement/proposal.md
  - docs/features/task-type-refinement/tasks/1-type-constants.md
  - docs/features/task-type-refinement/tasks/2-pipeline-logic.md
  - docs/features/task-type-refinement/tasks/3-prompt-templates.md
  - docs/features/task-type-refinement/tasks/4-dynamic-type-and-record.md
  - docs/features/task-type-refinement/tasks/5-migration-and-intent.md
  - docs/features/task-type-refinement/tasks/6-skill-updates.md
generated: "2026-05-16"
---

# Test Cases: task-type-refinement

## Summary

| Type | Count |
|------|-------|
| UI   | 0     |
| TUI  | 0     |
| Mobile | 0   |
| API  | 0     |
| CLI  | 20    |
| **Total** | **20** |

---

## CLI Test Cases

### Task 1: Type Constants and Registry

## TC-001: forge list-types displays all four new business types
- **Source**: Task 1 AC-6, Proposal Success Criterion 1
- **Type**: CLI
- **Target**: cli/list-types
- **Test ID**: cli/list-types/displays-four-new-business-types
- **Pre-conditions**: forge binary built and available in PATH
- **Steps**:
  1. Run `forge list-types`
  2. Inspect output for type entries
- **Expected**: Output contains entries for `feature`, `enhancement`, `cleanup`, and `refactor`, each with category "Core business"
- **Priority**: P0

## TC-002: forge list-types still shows deprecated TypeImplementation
- **Source**: Task 1 AC-2, Task 1 Hard Rule (keep TypeImplementation in ValidTypes)
- **Type**: CLI
- **Target**: cli/list-types
- **Test ID**: cli/list-types/shows-deprecated-implementation
- **Pre-conditions**: forge binary built and available in PATH
- **Steps**:
  1. Run `forge list-types`
  2. Inspect output for `implementation` entry
- **Expected**: Output contains `implementation` type with a deprecation note indicating it is deprecated in favor of the four new types
- **Priority**: P1

## TC-003: forge validates index.json with new type values
- **Source**: Task 1 AC-3 (ValidTypes includes new types), Task 1 AC-7 (existing tests pass)
- **Type**: CLI
- **Target**: cli/validate
- **Test ID**: cli/validate/accepts-new-type-values-in-index
- **Pre-conditions**: A feature directory with index.json containing tasks of type `feature`, `enhancement`, `cleanup`, `refactor`
- **Steps**:
  1. Run `forge validate` against the feature directory
  2. Inspect exit code and output
- **Expected**: Command exits with code 0, no validation errors for the new type values
- **Priority**: P0

### Task 2: Pipeline Logic (needsTestPipeline / needsDocEval)

## TC-004: forge build-index generates test pipeline for feature-typed tasks
- **Source**: Task 2 AC-1, Proposal D2 (needsTestPipeline true for feature)
- **Type**: CLI
- **Target**: cli/build-index
- **Test ID**: cli/build-index/generates-pipeline-for-feature-tasks
- **Pre-conditions**: A feature directory with index.json containing a business task of type `feature`
- **Steps**:
  1. Run `forge build-index` for the feature
  2. Inspect generated index.json for auto-gen test pipeline tasks (T-quick-1 through T-quick-6 or equivalent)
- **Expected**: Auto-gen test pipeline tasks are present in the generated index.json
- **Priority**: P0

## TC-005: forge build-index generates test pipeline for enhancement-typed tasks
- **Source**: Task 2 AC-1, Proposal D2 (needsTestPipeline true for enhancement)
- **Type**: CLI
- **Target**: cli/build-index
- **Test ID**: cli/build-index/generates-pipeline-for-enhancement-tasks
- **Pre-conditions**: A feature directory with index.json containing a business task of type `enhancement`
- **Steps**:
  1. Run `forge build-index` for the feature
  2. Inspect generated index.json for auto-gen test pipeline tasks
- **Expected**: Auto-gen test pipeline tasks are present in the generated index.json
- **Priority**: P0

## TC-006: forge build-index generates test pipeline for fix-typed tasks
- **Source**: Task 2 AC-1, Proposal D2 (needsTestPipeline true for fix)
- **Type**: CLI
- **Target**: cli/build-index
- **Test ID**: cli/build-index/generates-pipeline-for-fix-tasks
- **Pre-conditions**: A feature directory with index.json containing a business task of type `fix`
- **Steps**:
  1. Run `forge build-index` for the feature
  2. Inspect generated index.json for auto-gen test pipeline tasks
- **Expected**: Auto-gen test pipeline tasks are present in the generated index.json
- **Priority**: P0

## TC-007: forge build-index skips test pipeline for cleanup-only feature
- **Source**: Task 2 AC-1, Proposal D2 (needsTestPipeline false for cleanup), Proposal Success Criterion 2
- **Type**: CLI
- **Target**: cli/build-index
- **Test ID**: cli/build-index/skips-pipeline-for-cleanup-only
- **Pre-conditions**: A feature directory with index.json containing only business tasks of type `cleanup`
- **Steps**:
  1. Run `forge build-index` for the feature
  2. Inspect generated index.json for auto-gen test pipeline tasks
- **Expected**: No auto-gen test pipeline tasks are present in the generated index.json
- **Priority**: P0

## TC-008: forge build-index skips test pipeline for refactor-only feature
- **Source**: Task 2 AC-1, Proposal D2 (needsTestPipeline false for refactor), Proposal Success Criterion 2
- **Type**: CLI
- **Target**: cli/build-index
- **Test ID**: cli/build-index/skips-pipeline-for-refactor-only
- **Pre-conditions**: A feature directory with index.json containing only business tasks of type `refactor`
- **Steps**:
  1. Run `forge build-index` for the feature
  2. Inspect generated index.json for auto-gen test pipeline tasks
- **Expected**: No auto-gen test pipeline tasks are present in the generated index.json
- **Priority**: P0

## TC-009: forge build-index generates T-eval-doc for documentation-only feature
- **Source**: Task 2 AC-2, Proposal D2 (needsDocEval true for documentation), Proposal Success Criterion 3
- **Type**: CLI
- **Target**: cli/build-index
- **Test ID**: cli/build-index/generates-eval-doc-for-documentation-only
- **Pre-conditions**: A feature directory with index.json containing only business tasks of type `documentation`
- **Steps**:
  1. Run `forge build-index` for the feature
  2. Inspect generated index.json for T-eval-doc task
- **Expected**: A T-eval-doc auto-gen task is present; no test pipeline tasks are present
- **Priority**: P0

## TC-010: forge build-index generates neither pipeline nor eval-doc for mixed cleanup-refactor feature
- **Source**: Task 2 AC-5 (three-tier behavior), Proposal D2 table row "only cleanup/refactor"
- **Type**: CLI
- **Target**: cli/build-index
- **Test ID**: cli/build-index/no-pipeline-no-eval-for-cleanup-refactor
- **Pre-conditions**: A feature directory with index.json containing tasks of type `cleanup` and `refactor` only
- **Steps**:
  1. Run `forge build-index` for the feature
  2. Inspect generated index.json
- **Expected**: No test pipeline tasks and no T-eval-doc task are present in the generated index.json
- **Priority**: P1

## TC-011: forge quality gate skips for cleanup-only feature
- **Source**: Task 2 AC-5, Proposal D2 (quality_gate.go updated)
- **Type**: CLI
- **Target**: cli/quality-gate
- **Test ID**: cli/quality-gate/skips-for-cleanup-only
- **Pre-conditions**: A feature directory with only `cleanup`-typed tasks, all marked as completed
- **Steps**:
  1. Run `forge task submit` or quality gate command for the feature
  2. Observe quality gate behavior
- **Expected**: Quality gate is skipped (does not run compile/test/lint/fmt checks)
- **Priority**: P1

### Task 3: Prompt Templates

## TC-012: forge prompt get-by-task-id returns feature template for feature-typed task
- **Source**: Task 3 AC-5 (typeToTemplate map), Proposal D3 table
- **Type**: CLI
- **Target**: cli/prompt
- **Test ID**: cli/prompt/returns-feature-template
- **Pre-conditions**: A feature with a task of type `feature` in index.json
- **Steps**:
  1. Run `forge prompt get-by-task-id <task-id>` where task type is `feature`
  2. Inspect output
- **Expected**: Output contains the feature-specific prompt template (with "implement functionality" workflow)
- **Priority**: P0

## TC-013: forge prompt get-by-task-id returns cleanup template for cleanup-typed task
- **Source**: Task 3 AC-3, Proposal D3 (cleanup.md: improve technical debt, no TDD)
- **Type**: CLI
- **Target**: cli/prompt
- **Test ID**: cli/prompt/returns-cleanup-template
- **Pre-conditions**: A feature with a task of type `cleanup` in index.json
- **Steps**:
  1. Run `forge prompt get-by-task-id <task-id>` where task type is `cleanup`
  2. Inspect output
- **Expected**: Output contains the cleanup-specific prompt template (with "improve technical debt" workflow, no TDD requirement)
- **Priority**: P0

## TC-014: forge prompt get-by-task-id returns refactor template for refactor-typed task
- **Source**: Task 3 AC-4, Proposal D3 (refactor.md: behavior preservation check)
- **Type**: CLI
- **Target**: cli/prompt
- **Test ID**: cli/prompt/returns-refactor-template
- **Pre-conditions**: A feature with a task of type `refactor` in index.json
- **Steps**:
  1. Run `forge prompt get-by-task-id <task-id>` where task type is `refactor`
  2. Inspect output
- **Expected**: Output contains the refactor-specific prompt template (with "restructure code, verify behavior unchanged" workflow)
- **Priority**: P0

### Task 4: Dynamic Fix Task Type and Type Reclassification

## TC-015: forge creates fix-typed dynamic task on compile failure
- **Source**: Task 4 AC-1, Proposal D4 (compile failure -> TypeFix)
- **Type**: CLI
- **Target**: cli/quality-gate
- **Test ID**: cli/quality-gate/creates-fix-type-on-compile-failure
- **Pre-conditions**: A feature with a task whose code change introduces a compile error
- **Steps**:
  1. Run `forge task submit` triggering a quality gate that fails at compile step
  2. Inspect the dynamically created fix task in index.json
- **Expected**: The new auto-gen fix task has `type: "fix"` in its frontmatter
- **Priority**: P0

## TC-016: forge creates cleanup-typed dynamic task on fmt failure
- **Source**: Task 4 AC-2, Proposal D4 (fmt failure -> TypeCleanup)
- **Type**: CLI
- **Target**: cli/quality-gate
- **Test ID**: cli/quality-gate/creates-cleanup-type-on-fmt-failure
- **Pre-conditions**: A feature with a task whose code change has formatting violations
- **Steps**:
  1. Run `forge task submit` triggering a quality gate that fails at fmt step
  2. Inspect the dynamically created fix task in index.json
- **Expected**: The new auto-gen fix task has `type: "cleanup"` in its frontmatter
- **Priority**: P0

## TC-017: forge creates cleanup-typed dynamic task on lint failure
- **Source**: Task 4 AC-2, Proposal D4 (lint failure -> TypeCleanup)
- **Type**: CLI
- **Target**: cli/quality-gate
- **Test ID**: cli/quality-gate/creates-cleanup-type-on-lint-failure
- **Pre-conditions**: A feature with a task whose code change has lint violations
- **Steps**:
  1. Run `forge task submit` triggering a quality gate that fails at lint step
  2. Inspect the dynamically created fix task in index.json
- **Expected**: The new auto-gen fix task has `type: "cleanup"` in its frontmatter
- **Priority**: P0

## TC-018: forge record contains Type Reclassification block when type shifts
- **Source**: Task 4 AC-3, AC-4, Proposal D5
- **Type**: CLI
- **Target**: cli/submit
- **Test ID**: cli/submit/record-has-reclassification-when-type-shifts
- **Pre-conditions**: A task with original type `fix` that executor reclassifies to `cleanup`
- **Steps**:
  1. Run `forge task submit` for a task where type reclassification occurred during execution
  2. Read the generated record file
- **Expected**: Record file contains a "## Type Reclassification" section with original type, actual type, and reason fields
- **Priority**: P1

## TC-019: forge record omits Type Reclassification block when no type shift
- **Source**: Task 4 AC-4, Proposal Success Criterion 7
- **Type**: CLI
- **Target**: cli/submit
- **Test ID**: cli/submit/record-omits-reclassification-when-no-shift
- **Pre-conditions**: A task executed with no type change from its declared type
- **Steps**:
  1. Run `forge task submit` for a normally executed task (no type reclassification)
  2. Read the generated record file
- **Expected**: Record file does NOT contain a "## Type Reclassification" section
- **Priority**: P1

### Task 5: Migration and Intent

## TC-020: forge task migrate maps implementation to feature
- **Source**: Task 5 AC-1, AC-2, Proposal Success Criterion 8
- **Type**: CLI
- **Target**: cli/migrate
- **Test ID**: cli/migrate/maps-implementation-to-feature
- **Pre-conditions**: An index.json file containing tasks with `type: "implementation"`
- **Steps**:
  1. Run `forge task migrate` on the feature
  2. Read the migrated index.json
- **Expected**: All tasks with `type: "implementation"` now have `type: "feature"`. Other task types are unchanged.
- **Priority**: P0

---

## Traceability

| TC ID | Source | Type | Target | Priority |
|-------|--------|------|--------|----------|
| TC-001 | Task 1 AC-6, Proposal SC-1 | CLI | cli/list-types | P0 |
| TC-002 | Task 1 AC-2 | CLI | cli/list-types | P1 |
| TC-003 | Task 1 AC-3, AC-7 | CLI | cli/validate | P0 |
| TC-004 | Task 2 AC-1, Proposal D2 | CLI | cli/build-index | P0 |
| TC-005 | Task 2 AC-1, Proposal D2 | CLI | cli/build-index | P0 |
| TC-006 | Task 2 AC-1, Proposal D2 | CLI | cli/build-index | P0 |
| TC-007 | Task 2 AC-1, Proposal SC-2 | CLI | cli/build-index | P0 |
| TC-008 | Task 2 AC-1, Proposal SC-2 | CLI | cli/build-index | P0 |
| TC-009 | Task 2 AC-2, Proposal SC-3 | CLI | cli/build-index | P0 |
| TC-010 | Task 2 AC-5, Proposal D2 | CLI | cli/build-index | P1 |
| TC-011 | Task 2 AC-5 | CLI | cli/quality-gate | P1 |
| TC-012 | Task 3 AC-5, Proposal D3 | CLI | cli/prompt | P0 |
| TC-013 | Task 3 AC-3, Proposal D3 | CLI | cli/prompt | P0 |
| TC-014 | Task 3 AC-4, Proposal D3 | CLI | cli/prompt | P0 |
| TC-015 | Task 4 AC-1, Proposal D4 | CLI | cli/quality-gate | P0 |
| TC-016 | Task 4 AC-2, Proposal D4 | CLI | cli/quality-gate | P0 |
| TC-017 | Task 4 AC-2, Proposal D4 | CLI | cli/quality-gate | P0 |
| TC-018 | Task 4 AC-3, AC-4, Proposal D5 | CLI | cli/submit | P1 |
| TC-019 | Task 4 AC-4, Proposal SC-7 | CLI | cli/submit | P1 |
| TC-020 | Task 5 AC-1, AC-2, Proposal SC-8 | CLI | cli/migrate | P0 |
