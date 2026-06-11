---
feature: "task-stage-gates"
sources:
  - docs/proposals/task-stage-gates/proposal.md (quick mode: proposal serves as requirements source)
generated: "2026-05-14"
---

# Test Cases: task-stage-gates

> **Note**: Quick mode feature -- no formal PRD. Acceptance criteria extracted from proposal Success Criteria (section "Success Criteria") and Key Scenarios (section "Key Scenarios"). The proposal is the sole input source.

## Summary

| Type | Count |
|------|-------|
| UI   | 0   |
| **Integration** | **0** |
| API  | 0  |
| CLI  | 20  |
| **Total** | **20** |

---

## CLI Test Cases

### Happy Path & Phase Detection

## TC-001: Generates summary and gate files for phases with >=2 business tasks
- **Source**: Proposal Success Criteria #1, Key Scenario "Happy path"
- **Type**: CLI
- **Target**: cli/task-index
- **Test ID**: cli/task-index/generates-summary-and-gate-for-qualifying-phases
- **Pre-conditions**: Feature directory with 4 business tasks across 2 phases (e.g., tasks 1.1, 1.2, 2.1, 2.2). No pre-existing `.summary` or `.gate` files.
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Create a feature directory with tasks: `1.1-task-a.md`, `1.2-task-b.md`, `2.1-task-c.md`, `2.2-task-d.md` (all business task type)
  2. Run `forge task index --feature <slug>`
  3. Check that `1.summary.md`, `1.gate.md`, `2.summary.md`, `2.gate.md` are created
  4. Verify `1.summary.md` has `depends_on` containing `["1.1", "1.2"]`
  5. Verify `1.gate.md` has `depends_on` containing `["1.summary"]` and `breaking: true`
  6. Verify `2.summary.md` has `depends_on` containing `["2.1", "2.2"]`
  7. Verify `2.gate.md` has `depends_on` containing `["2.summary"]` and `breaking: true`
- **Expected**: Two `.summary` and two `.gate` files created with correct dependency wiring
- **Priority**: P0

## TC-002: Correct dependency wiring for gate tasks
- **Source**: Proposal Success Criteria #2
- **Type**: CLI
- **Target**: cli/task-index
- **Test ID**: cli/task-index/gate-depends-on-summary-summary-depends-on-business-tasks
- **Pre-conditions**: Feature directory with phase containing 3 business tasks (e.g., 1.1, 1.2, 1.3)
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Create tasks `1.1.md`, `1.2.md`, `1.3.md` in feature directory
  2. Run `forge task index --feature <slug>`
  3. Parse `1.summary.md` frontmatter and verify `depends_on` includes all of `["1.1", "1.2", "1.3"]`
  4. Parse `1.gate.md` frontmatter and verify `depends_on` is exactly `["1.summary"]`
  5. Verify `1.gate.md` has `breaking: true`
- **Expected**: Summary depends on all same-phase business tasks; gate depends on summary with breaking flag
- **Priority**: P0

## TC-003: Skips single-task phases
- **Source**: Proposal Success Criteria #1, Key Scenario "Single-task phase"
- **Type**: CLI
- **Target**: cli/task-index
- **Test ID**: cli/task-index/skips-single-task-phases
- **Pre-conditions**: Feature directory with phase 1 having 1 business task and phase 2 having 2 business tasks
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Create tasks: `1.1-task.md` (1 business task in phase 1), `2.1-task-a.md`, `2.2-task-b.md` (2 business tasks in phase 2)
  2. Run `forge task index --feature <slug>`
  3. Verify no `1.summary.md` or `1.gate.md` created
  4. Verify `2.summary.md` and `2.gate.md` are created
- **Expected**: Phase with <2 business tasks gets no gate/summary; multi-task phase does
- **Priority**: P0

## TC-004: Excludes test-only phases
- **Source**: Proposal Key Scenario "Phase with only test tasks", Success Criteria #1
- **Type**: CLI
- **Target**: cli/task-index
- **Test ID**: cli/task-index/excludes-test-only-phases
- **Pre-conditions**: Feature directory with a phase containing only T-test and T-quick tasks
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Create tasks: `1.1.md` (type: developTask), `T-test-1.md` (type: testTask), `T-quick-1.md` (type: testTask)
  2. Phase 1 has 1 business task (1.1), phase with T-test/T-quick only has 0 business tasks
  3. Run `forge task index --feature <slug>`
  4. Verify no gate/summary generated for the test-only phase
  5. Verify phase 1 also gets no gate (only 1 business task)
- **Expected**: Phases with only T-test/T-quick tasks are excluded from gate generation
- **Priority**: P0

## TC-005: Filters T-test and T-quick tasks from business task count
- **Source**: Proposal Key Scenario "Phase with only test tasks", section "How It Works" step 1
- **Type**: CLI
- **Target**: cli/task-index
- **Test ID**: cli/task-index/filters-test-tasks-from-business-count
- **Pre-conditions**: Feature directory with phase containing 2 business tasks + 1 T-test task
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Create tasks: `1.1-task.md`, `1.2-task.md`, `T-test-1.md` (type: testTask, in same feature dir but ID doesn't match `<digit>.<digit>` so not in phase 1) -- OR create `1.3.md` with `type: testTask`
  2. Run `forge task index --feature <slug>`
  3. Verify business task count for phase 1 is 2 (excluding the test task)
  4. Verify `1.summary.md` is generated (count >=2)
  5. Verify `1.summary.md` `depends_on` does NOT include the test task
- **Expected**: Test tasks excluded from count and from dependency wiring
- **Priority**: P0

### Idempotency & Partial State

## TC-006: Idempotent re-run preserves existing files
- **Source**: Proposal Success Criteria #3, Key Scenario "Idempotent re-run"
- **Type**: CLI
- **Target**: cli/task-index
- **Test ID**: cli/task-index/idempotent-rerun-preserves-existing-files
- **Pre-conditions**: Feature directory where `forge task index` has already been run, producing `1.summary.md` and `1.gate.md`
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Run `forge task index --feature <slug>` to generate gate/summary files
  2. Record file contents and modification times of `1.summary.md` and `1.gate.md`
  3. Run `forge task index --feature <slug>` again
  4. Verify file contents are unchanged
  5. Verify modification times are unchanged
- **Expected**: Second run produces identical output, does not overwrite
- **Priority**: P0

## TC-007: Generates only missing gate when summary already exists
- **Source**: Proposal Success Criteria #4, Key Scenario "Partial state"
- **Type**: CLI
- **Target**: cli/task-index
- **Test ID**: cli/task-index/generates-only-missing-gate-when-summary-exists
- **Pre-conditions**: Feature directory with phase 1 having 2 business tasks. `1.summary.md` already exists but `1.gate.md` does not.
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Create `1.summary.md` manually in feature tasks directory with valid content
  2. Ensure `1.gate.md` does not exist
  3. Run `forge task index --feature <slug>`
  4. Verify `1.summary.md` content is unchanged (not overwritten)
  5. Verify `1.gate.md` is created with correct dependencies
- **Expected**: Only the missing `.gate.md` is generated; existing `.summary.md` preserved
- **Priority**: P0

## TC-008: Generates only missing summary when gate already exists
- **Source**: Proposal Key Scenario "Partial state" (reverse case)
- **Type**: CLI
- **Target**: cli/task-index
- **Test ID**: cli/task-index/generates-only-missing-summary-when-gate-exists
- **Pre-conditions**: Feature directory with `1.gate.md` already existing but `1.summary.md` missing
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Create `1.gate.md` manually in feature tasks directory with valid content
  2. Ensure `1.summary.md` does not exist
  3. Run `forge task index --feature <slug>`
  4. Verify `1.gate.md` content is unchanged
  5. Verify `1.summary.md` is created with correct dependencies
- **Expected**: Only the missing `.summary.md` is generated; existing `.gate.md` preserved
- **Priority**: P1

### Malformed & Edge Cases

## TC-009: Silently skips malformed task IDs
- **Source**: Proposal Success Criteria #5, Key Scenario "Malformed task IDs"
- **Type**: CLI
- **Target**: cli/task-index
- **Test ID**: cli/task-index/silently-skips-malformed-task-ids
- **Pre-conditions**: Feature directory with valid tasks plus malformed IDs
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Create tasks: `1.1-task.md`, `1.2-task.md`, `intro.md`, `1.2a-task.md`, `overview.md`
  2. Run `forge task index --feature <slug>`
  3. Verify command exits with code 0 (no crash)
  4. Verify `1.summary.md` and `1.gate.md` are generated (from valid IDs 1.1, 1.2)
  5. Verify no gate/summary generated for malformed IDs
- **Expected**: Malformed IDs silently skipped; valid phases processed normally
- **Priority**: P0

## TC-010: Pre-existing hand-crafted gate files are preserved
- **Source**: Proposal Key Scenario "Pre-existing hand-crafted gate files"
- **Type**: CLI
- **Target**: cli/task-index
- **Test ID**: cli/task-index/preserves-preexisting-hand-crafted-gate-files
- **Pre-conditions**: Feature directory with manually created `1.gate.md` containing custom content
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Create `1.gate.md` with custom hand-crafted content (different from template output)
  2. Create tasks `1.1.md`, `1.2.md`
  3. Run `forge task index --feature <slug>`
  4. Verify `1.gate.md` content is identical to the hand-crafted version (not overwritten)
  5. Verify `1.summary.md` is generated from template
- **Expected**: Existing gate file preserved; only missing files generated
- **Priority**: P1

### Index Output & CLI Behavior

## TC-011: Generated tasks appear in index.json with correct type
- **Source**: Proposal Success Criteria #6
- **Type**: CLI
- **Target**: cli/task-index
- **Test ID**: cli/task-index/generated-tasks-in-index-json-with-correct-type
- **Pre-conditions**: Feature directory with 2 phases each having >=2 business tasks
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Create tasks `1.1.md`, `1.2.md`, `2.1.md`, `2.2.md`
  2. Run `forge task index --feature <slug>`
  3. Parse `index.json`
  4. Verify `1.summary` entry exists with `type: "doc-generation.summary"`
  5. Verify `1.gate` entry exists with `type: "gate"`
  6. Verify `2.summary` entry exists with `type: "doc-generation.summary"`
  7. Verify `2.gate` entry exists with `type: "gate"`
- **Expected**: All generated gate/summary tasks appear in index.json with correct type field
- **Priority**: P0

## TC-012: CLI prints summary line per qualifying phase
- **Source**: Proposal section "CLI Output Behavior" - "Gates generated"
- **Type**: CLI
- **Target**: cli/task-index
- **Test ID**: cli/task-index/prints-summary-line-per-qualifying-phase
- **Pre-conditions**: Feature directory with 2 qualifying phases
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Create tasks for 2 qualifying phases (e.g., 1.1, 1.2, 2.1, 2.2)
  2. Run `forge task index --feature <slug>` and capture stderr
  3. Verify stderr contains line matching `Generated stage-gate: phase 1`
  4. Verify stderr contains line matching `Generated stage-gate: phase 2`
  5. Verify stderr contains "Indexed N tasks" line
- **Expected**: One summary line per qualifying phase printed to stderr
- **Priority**: P1

## TC-013: CLI prints no-qualification message when no phases qualify
- **Source**: Proposal section "CLI Output Behavior" - "Zero phases qualify"
- **Type**: CLI
- **Target**: cli/task-index
- **Test ID**: cli/task-index/prints-no-qualification-message
- **Pre-conditions**: Feature directory where all phases have <2 business tasks
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Create tasks: `1.1-task.md` (single task in phase 1)
  2. Run `forge task index --feature <slug>` and capture stderr
  3. Verify stderr contains "No phases qualified for stage-gate generation"
- **Expected**: Clear message when no phases meet the >=2 business task threshold
- **Priority**: P1

## TC-014: CLI exits with error on template rendering failure
- **Source**: Proposal section "CLI Output Behavior" - "Template rendering failure"
- **Type**: CLI
- **Target**: cli/task-index
- **Test ID**: cli/task-index/exits-with-error-on-template-render-failure
- **Pre-conditions**: Feature directory with qualifying phases; template embedded binary is corrupted (unit test scenario)
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. This is a unit-level test scenario: mock/modify embedded template to be malformed
  2. Run `forge task index --feature <slug>`
  3. Verify exit code is 1
  4. Verify error message indicates template rendering failure
- **Expected**: Non-zero exit with descriptive error when template rendering fails
- **Priority**: P1

### Quick Mode & Backward Compatibility

## TC-015: Quick mode generates stage-gates identically to full mode
- **Source**: Proposal Key Scenario "Quick mode", Success Criteria #7
- **Type**: CLI
- **Target**: cli/task-index
- **Test ID**: cli/task-index/quick-mode-generates-stage-gates-identically
- **Pre-conditions**: Feature directory with manifest `mode: quick` and qualifying phases
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Create a quick-mode feature directory with tasks across 2 phases
  2. Run `forge task index --feature <slug>`
  3. Verify gate/summary generation behavior is identical to full mode
  4. Verify same dependency wiring, same file naming, same index.json entries
- **Expected**: No difference in gate/summary generation between quick and full mode
- **Priority**: P1

## TC-016: Does not break existing forge task index behavior
- **Source**: Proposal section "Constraints & Dependencies"
- **Type**: CLI
- **Target**: cli/task-index
- **Test ID**: cli/task-index/does-not-break-existing-index-behavior
- **Pre-conditions**: Feature directory without qualifying phases (no stage-gate generation triggered)
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Use existing feature directory that previously indexed correctly
  2. Run `forge task index --feature <slug>`
  3. Verify `index.json` is generated correctly as before
  4. Verify all existing task entries are present
  5. Verify existing test task auto-generation still works
- **Expected**: No regression in existing index behavior
- **Priority**: P0

## TC-017: Respects no-test flag independently of stage-gates
- **Source**: Proposal section "Constraints & Dependencies" - "Must respect --no-test flag"
- **Type**: CLI
- **Target**: cli/task-index
- **Test ID**: cli/task-index/no-test-flag-does-not-affect-stage-gates
- **Pre-conditions**: Feature directory with qualifying phases for both test tasks and stage-gates
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Run `forge task index --feature <slug> --no-test`
  2. Verify test tasks are NOT generated (--no-test respected)
  3. Verify `.summary` and `.gate` files ARE generated (stage-gates unaffected by --no-test)
- **Expected**: --no-test suppresses test task generation but not stage-gate generation
- **Priority**: P0

### Concurrent & Deterministic Behavior

## TC-018: Concurrent execution produces identical output
- **Source**: Proposal Key Scenario "Concurrent execution"
- **Type**: CLI
- **Target**: cli/task-index
- **Test ID**: cli/task-index/concurrent-execution-identical-output
- **Pre-conditions**: Feature directory with qualifying phases
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Run two instances of `forge task index --feature <slug>` concurrently
  2. Verify both complete successfully (exit code 0)
  3. Verify generated files have identical content
- **Expected**: Deterministic templates produce identical output; no corruption from concurrent writes
- **Priority**: P2

### Security & Performance

## TC-019: Phase detection rejects path traversal in task IDs
- **Source**: Proposal section "Non-Functional Requirements" - Security
- **Type**: CLI
- **Target**: cli/task-index
- **Test ID**: cli/task-index/rejects-path-traversal-in-task-ids
- **Pre-conditions**: Feature directory with task IDs containing path traversal patterns
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Create task files with IDs like `../1.1` or `1.1/../2.1` (if possible via filename or frontmatter)
  2. Run `forge task index --feature <slug>`
  3. Verify command does not crash
  4. Verify no files are written outside the feature tasks directory
  5. Verify regex `^\d+\.\d+$` rejects these patterns
- **Expected**: Path traversal patterns silently rejected by phase detection regex
- **Priority**: P1

## TC-020: Generation completes within 5ms for 100 tasks and 20 phases
- **Source**: Proposal section "Non-Functional Requirements" - Performance
- **Type**: CLI
- **Target**: cli/task-index
- **Test ID**: cli/task-index/generation-completes-within-5ms
- **Pre-conditions**: Feature directory with 100 tasks across 20 phases, each phase having >=2 business tasks
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Create 100 task files across 20 phases (5 tasks per phase)
  2. Measure time for `forge task index --feature <slug>` gate/summary generation phase
  3. Verify total generation time < 5ms (measured by file-write syscall overhead)
- **Expected**: Fast deterministic generation well within performance budget
- **Priority**: P2

---

## Traceability

| TC ID | Source | Type | Target | Priority |
|-------|--------|------|--------|----------|
| TC-001 | Proposal Success Criteria #1, Key Scenario "Happy path" | CLI | cli/task-index | P0 |
| TC-002 | Proposal Success Criteria #2 | CLI | cli/task-index | P0 |
| TC-003 | Proposal Success Criteria #1, Key Scenario "Single-task phase" | CLI | cli/task-index | P0 |
| TC-004 | Proposal Key Scenario "Phase with only test tasks", Success Criteria #1 | CLI | cli/task-index | P0 |
| TC-005 | Proposal Key Scenario "Phase with only test tasks", How It Works step 1 | CLI | cli/task-index | P0 |
| TC-006 | Proposal Success Criteria #3, Key Scenario "Idempotent re-run" | CLI | cli/task-index | P0 |
| TC-007 | Proposal Success Criteria #4, Key Scenario "Partial state" | CLI | cli/task-index | P0 |
| TC-008 | Proposal Key Scenario "Partial state" (reverse case) | CLI | cli/task-index | P1 |
| TC-009 | Proposal Success Criteria #5, Key Scenario "Malformed task IDs" | CLI | cli/task-index | P0 |
| TC-010 | Proposal Key Scenario "Pre-existing hand-crafted gate files" | CLI | cli/task-index | P1 |
| TC-011 | Proposal Success Criteria #6 | CLI | cli/task-index | P0 |
| TC-012 | Proposal CLI Output Behavior "Gates generated" | CLI | cli/task-index | P1 |
| TC-013 | Proposal CLI Output Behavior "Zero phases qualify" | CLI | cli/task-index | P1 |
| TC-014 | Proposal CLI Output Behavior "Template rendering failure" | CLI | cli/task-index | P1 |
| TC-015 | Proposal Key Scenario "Quick mode", Success Criteria #7 | CLI | cli/task-index | P1 |
| TC-016 | Proposal Constraints & Dependencies | CLI | cli/task-index | P0 |
| TC-017 | Proposal Constraints & Dependencies "--no-test flag" | CLI | cli/task-index | P0 |
| TC-018 | Proposal Key Scenario "Concurrent execution" | CLI | cli/task-index | P2 |
| TC-019 | Proposal Non-Functional Requirements - Security | CLI | cli/task-index | P1 |
| TC-020 | Proposal Non-Functional Requirements - Performance | CLI | cli/task-index | P2 |

---

## Route Validation

_Omitted -- project is a CLI tool with no web routes. Route validation is not applicable._
