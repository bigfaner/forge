---
feature: "spec-drift-detection"
sources:
  - docs/proposals/spec-drift-detection/proposal.md
generated: "2026-05-15"
---

# Test Cases: spec-drift-detection

## Summary

| Type | Count |
|------|-------|
| UI   | 0   |
| **Integration** | **0** |
| API  | 0  |
| CLI  | 19  |
| **Total** | **19** |

> **Note**: This is a CLI-only feature with no UI or API interface. The feature extends the forge CLI task type system, test pipeline, and skill documentation. All test cases verify CLI behavior or file content correctness.

---

## CLI Test Cases

### Task Type: doc-generation.drift

## TC-001: List types includes doc-generation.drift with correct description
- **Source**: Task 2 AC — TypeDocGenerationDrift constant and registry entry
- **Type**: CLI
- **Target**: cli/task-list-types
- **Test ID**: cli/task-list-types/drift-type-in-registry
- **Pre-conditions**: forge binary built and available in PATH
- **Steps**:
  1. Run `forge task list-types`
  2. Check exit code is 0
  3. Verify output contains `doc-generation.drift` type
  4. Verify the description for `doc-generation.drift` is "detect and fix spec drift against codebase"
  5. Verify total type count includes doc-generation.drift (14 types total)
- **Expected**: Exit code 0, output lists `doc-generation.drift` with description "detect and fix spec drift against codebase"
- **Priority**: P0

## TC-002: Prompt for doc-generation.drift type resolves to correct strategy template
- **Source**: Task 2 AC — prompt.go maps new type to strategy template path
- **Type**: CLI
- **Target**: cli/prompt-get
- **Test ID**: cli/prompt-get/drift-type-resolves-strategy-template
- **Pre-conditions**: A task exists in index.json with type `doc-generation.drift`; forge binary built
- **Steps**:
  1. Run `forge prompt get-by-task-id <drift-task-id>`
  2. Check exit code is 0
  3. Verify prompt output contains instructions to invoke `consolidate-specs` skill in drift-only mode
  4. Verify prompt output mentions Steps 9-11 (drift detection, auto-fix, commit)
  5. Verify prompt output does NOT contain Steps 1-8 (extraction, review, integration)
- **Expected**: Exit code 0, prompt content describes drift-only consolidation workflow (Steps 9-11 only)
- **Priority**: P0

## TC-003: Type doc-generation.drift is a valid type recognized by validate-index
- **Source**: Task 2 AC — Valid type map includes doc-generation.drift
- **Type**: CLI
- **Target**: cli/task-validate-index
- **Test ID**: cli/task-validate-index/drift-type-is-valid
- **Pre-conditions**: index.json contains a task with type `doc-generation.drift`
- **Steps**:
  1. Run `forge task validate-index`
  2. Check exit code is 0
  3. Verify no validation errors about unknown task type
  4. Verify output confirms all task types are valid
- **Expected**: Exit code 0, no unknown-type errors for `doc-generation.drift` tasks
- **Priority**: P0

### T-quick-6 Task Generation

## TC-004: Quick test tasks include T-quick-6 with doc-generation.drift type
- **Source**: Task 2 AC — T-quick-6 added after T-quick-5 in generateQuickTestTasks
- **Type**: CLI
- **Target**: cli/task-index
- **Test ID**: cli/task-index/quick-pipeline-includes-t-quick-6-drift
- **Pre-conditions**: Feature directory exists in quick mode (proposal.md present, no PRD/design); forge binary built
- **Steps**:
  1. Set feature context to a quick-mode feature
  2. Run `forge task index` to regenerate index.json
  3. Check exit code is 0
  4. Verify index.json contains a task with ID "T-quick-6"
  5. Verify T-quick-6 type is `doc-generation.drift`
  6. Verify T-quick-6 NoTest is true
  7. Verify T-quick-6 scope is "all"
  8. Verify T-quick-6 title contains "drift" or "Drift"
- **Expected**: Exit code 0, index.json includes T-quick-6 with type `doc-generation.drift`, NoTest: true, scope: all
- **Priority**: P0

## TC-005: T-quick-6 depends on T-quick-5 in quick pipeline
- **Source**: Task 2 AC — T-quick-6 depends on T-quick-5 in resolveQuickDeps
- **Type**: CLI
- **Target**: cli/task-index
- **Test ID**: cli/task-index/t-quick-6-depends-on-t-quick-5
- **Pre-conditions**: Feature directory exists in quick mode; forge binary built
- **Steps**:
  1. Set feature context to a quick-mode feature
  2. Run `forge task index` to regenerate index.json
  3. Check exit code is 0
  4. Run `cat index.json | jq '.tasks[] | select(.id=="T-quick-6") | .dependencies'` to get T-quick-6 dependencies
  5. Verify dependencies array contains "T-quick-5"
- **Expected**: T-quick-6 dependencies include T-quick-5
- **Priority**: P0

## TC-006: T-quick-6 appears after T-quick-5 in task generation order
- **Source**: Task 2 AC — T-quick-6 added after T-quick-5
- **Type**: CLI
- **Target**: cli/task-index
- **Test ID**: cli/task-index/t-quick-6-order-after-t-quick-5
- **Pre-conditions**: Feature directory exists in quick mode; forge binary built
- **Steps**:
  1. Set feature context to a quick-mode feature
  2. Run `forge task index` to regenerate index.json
  3. Verify T-quick-6 task file appears after T-quick-5 task file in the tasks directory listing
  4. Verify index.json tasks map includes both T-quick-5 and T-quick-6
- **Expected**: T-quick-6 generated after T-quick-5 in both file listing and index structure
- **Priority**: P1

### Breakdown Mode: T-test-5 Drift Scope

## TC-007: T-test-5 description includes drift detection in breakdown mode
- **Source**: Task 3 AC — breakdown-tasks SKILL.md references drift detection
- **Type**: CLI
- **Target**: cli/task-index
- **Test ID**: cli/task-index/t-test-5-description-includes-drift
- **Pre-conditions**: Feature directory exists in breakdown mode (PRD/design present); forge binary built
- **Steps**:
  1. Run `forge task index` for a breakdown-mode feature
  2. Check exit code is 0
  3. Verify index.json contains task T-test-5 with type `doc-generation.consolidate`
  4. Verify T-test-5 title or description references consolidation and drift detection
- **Expected**: T-test-5 task present with type `doc-generation.consolidate`, title reflects full pipeline scope
- **Priority**: P1

### Skill Documentation Verification

## TC-008: consolidate-specs SKILL.md contains Steps 9-11
- **Source**: Task 1 AC — SKILL.md contains new Steps 9-11 after existing Step 8
- **Type**: CLI
- **Target**: cli/skill-verification
- **Test ID**: cli/skill-verification/consolidate-specs-has-steps-9-11
- **Pre-conditions**: forge repository checked out with feature implementation
- **Steps**:
  1. Read `plugins/forge/skills/consolidate-specs/SKILL.md`
  2. Verify "Step 9" section exists with title containing "Detect Drift"
  3. Verify "Step 10" section exists with title containing "Auto-Fix"
  4. Verify "Step 11" section exists with title containing "Commit"
  5. Verify Steps 9-11 appear after Step 8 in document order
- **Expected**: SKILL.md contains Step 9 (Detect Drift), Step 10 (Auto-Fix), Step 11 (Commit) after Step 8
- **Priority**: P0

## TC-009: consolidate-specs SKILL.md Step 9 validates rules against code
- **Source**: Task 1 AC — Step 9 reads all business-rules and conventions, validates each rule, classifies as current/drifted/orphaned
- **Type**: CLI
- **Target**: cli/skill-verification
- **Test ID**: cli/skill-verification/step-9-classifies-rules
- **Pre-conditions**: forge repository checked out with feature implementation
- **Steps**:
  1. Read `plugins/forge/skills/consolidate-specs/SKILL.md` Step 9 section
  2. Verify it instructs reading `docs/business-rules/*.md`
  3. Verify it instructs reading `docs/conventions/*.md`
  4. Verify it defines classification: `current`, `drifted`, `orphaned`
  5. Verify it instructs validating each rule against current code
- **Expected**: Step 9 defines the three-way classification (current/drifted/orphaned) and reads both business-rules and conventions directories
- **Priority**: P0

## TC-010: consolidate-specs SKILL.md Step 10 preserves project-global IDs
- **Source**: Task 1 AC — Step 10 updates drifted rules in-place preserving project-global IDs, removes orphaned rules with commit message
- **Type**: CLI
- **Target**: cli/skill-verification
- **Test ID**: cli/skill-verification/step-10-preserves-ids
- **Pre-conditions**: forge repository checked out with feature implementation
- **Steps**:
  1. Read `plugins/forge/skills/consolidate-specs/SKILL.md` Step 10 section
  2. Verify it instructs preserving project-global IDs (e.g., `BIZ-auth-001`) during updates
  3. Verify it instructs removing orphaned rules and recording rule ID and reason
  4. Verify it instructs detecting implicit new rules from code changes
- **Expected**: Step 10 preserves project-global IDs during auto-fix, records orphaned rule deletions with ID and reason, detects implicit new rules
- **Priority**: P0

## TC-011: consolidate-specs SKILL.md HARD-GATE allows drift modification
- **Source**: Task 1 AC — HARD-GATE updated to allow modification when drift detected
- **Type**: CLI
- **Target**: cli/skill-verification
- **Test ID**: cli/skill-verification/hard-gate-allows-drift-modification
- **Pre-conditions**: forge repository checked out with feature implementation
- **Steps**:
  1. Read `plugins/forge/skills/consolidate-specs/SKILL.md` HARD-GATE section
  2. Verify the "Do NOT overwrite" gate includes an exception clause for "unless drift is detected in Step 9"
- **Expected**: HARD-GATE's second bullet includes the drift exception allowing modification when drift is detected
- **Priority**: P0

## TC-012: consolidate-specs SKILL.md supports drift-only mode
- **Source**: Task 1 AC — Skill supports drift-only mode when no PRD/design exists
- **Type**: CLI
- **Target**: cli/skill-verification
- **Test ID**: cli/skill-verification/drift-only-mode-supported
- **Pre-conditions**: forge repository checked out with feature implementation
- **Steps**:
  1. Read `plugins/forge/skills/consolidate-specs/SKILL.md`
  2. Verify it documents drift-only mode: when prd/prd-spec.md and design/tech-design.md do not exist
  3. Verify drift-only mode skips Steps 1-8 and runs only Steps 9-11
- **Expected**: SKILL.md documents drift-only mode entry condition and confirms Steps 1-8 are skipped
- **Priority**: P0

## TC-013: breakdown-tasks SKILL.md references drift detection for T-test-5
- **Source**: Task 3 AC — breakdown-tasks SKILL.md references T-test-5 drift detection
- **Type**: CLI
- **Target**: cli/skill-verification
- **Test ID**: cli/skill-verification/breakdown-tasks-t-test-5-drift
- **Pre-conditions**: forge repository checked out with feature implementation
- **Steps**:
  1. Read `plugins/forge/skills/breakdown-tasks/SKILL.md`
  2. Verify T-test-5 description mentions drift detection
  3. Verify description clarifies the full pipeline: extract, integrate, detect drift, auto-fix
- **Expected**: breakdown-tasks SKILL.md T-test-5 description includes drift detection and documents the full pipeline
- **Priority**: P1

## TC-014: guide.md reflects T-quick-6 and drift detection flow
- **Source**: Task 4 AC — guide.md updated with T-quick-6 and drift detection
- **Type**: CLI
- **Target**: cli/skill-verification
- **Test ID**: cli/skill-verification/guide-reflects-t-quick-6
- **Pre-conditions**: forge repository checked out with feature implementation
- **Steps**:
  1. Read `plugins/forge/hooks/guide.md`
  2. Verify Quick Mode section references T-quick-1 through T-quick-6 (not T-quick-1 through T-quick-5)
  3. Verify drift detection is mentioned as a test step
  4. Verify the mermaid diagram shows consolidate-specs with drift audit
- **Expected**: guide.md Quick Mode section lists T-quick-6, mentions drift detection, and the diagram reflects drift audit
- **Priority**: P1

## TC-015: guide.md specs rule mentions drift verification
- **Source**: Task 4 AC — specs/ rule in Directory Conventions updated to mention drift detection
- **Type**: CLI
- **Target**: cli/skill-verification
- **Test ID**: cli/skill-verification/guide-specs-drift-verification
- **Pre-conditions**: forge repository checked out with feature implementation
- **Steps**:
  1. Read `plugins/forge/hooks/guide.md`
  2. Verify the `docs/business-rules/` and `docs/conventions/` entries mention drift verification
- **Expected**: guide.md Directory Conventions section mentions drift verification for spec files
- **Priority**: P2

### Strategy Template Verification

## TC-016: doc-generation-drift.md strategy template exists
- **Source**: Task 2 AC — doc-generation-drift.md strategy template created
- **Type**: CLI
- **Target**: cli/skill-verification
- **Test ID**: cli/skill-verification/drift-strategy-template-exists
- **Pre-conditions**: forge repository checked out with feature implementation
- **Steps**:
  1. Verify file exists at `forge-cli/pkg/prompt/data/doc-generation-drift.md`
  2. Read the file content
  3. Verify it references `consolidate-specs` skill
  4. Verify it specifies drift-only mode (skip extraction, run Steps 9-11 only)
- **Expected**: Strategy template file exists, invokes consolidate-specs in drift-only mode
- **Priority**: P0

### Existing Tests Pass

## TC-017: All existing tests pass after feature changes
- **Source**: Task 2 AC — All existing tests pass
- **Type**: CLI
- **Target**: cli/test-suite
- **Test ID**: cli/test-suite/all-existing-tests-pass
- **Pre-conditions**: forge repository checked out with feature implementation; Go toolchain installed
- **Steps**:
  1. Run `go test -race -cover ./...` in forge-cli directory
  2. Verify exit code is 0
  3. Verify no test failures
- **Expected**: Exit code 0, all existing tests pass without failure
- **Priority**: P0

### Type Inference

## TC-018: Task ID T-quick-6 infers type as doc-generation.drift
- **Source**: Task 2 AC — Type system correctly infers doc-generation.drift for T-quick-6 pattern
- **Type**: CLI
- **Target**: cli/task-index
- **Test ID**: cli/task-index/t-quick-6-infers-drift-type
- **Pre-conditions**: forge binary built; feature directory exists in quick mode
- **Steps**:
  1. Run `forge task index` for a quick-mode feature
  2. Verify task T-quick-6 has type `doc-generation.drift` (not `doc-generation.consolidate` or any other type)
  3. Verify T-quick-6a variant also infers `doc-generation.drift`
  4. Verify T-quick-6b variant also infers `doc-generation.drift`
- **Expected**: T-quick-6, T-quick-6a, T-quick-6b all infer type as `doc-generation.drift`
- **Priority**: P0

### Workflow Diagram

## TC-019: consolidate-specs workflow diagram includes Steps 9-11
- **Source**: Task 1 AC — Workflow diagram updated to include Steps 9-11
- **Type**: CLI
- **Target**: cli/skill-verification
- **Test ID**: cli/skill-verification/workflow-diagram-includes-drift-steps
- **Pre-conditions**: forge repository checked out with feature implementation
- **Steps**:
  1. Read `plugins/forge/skills/consolidate-specs/SKILL.md`
  2. Locate the Workflow section with the step diagram
  3. Verify the diagram includes Step 9 (Detect Drift), Step 10 (Auto-Fix), Step 11 (Commit)
  4. Verify these steps appear after Step 8 in the diagram flow
- **Expected**: Workflow diagram in SKILL.md includes Steps 9-11 in the correct position after Step 8
- **Priority**: P1

---

## Traceability

| TC ID | Source | Type | Target | Priority |
|-------|--------|------|--------|----------|
| TC-001 | Task 2 AC (type constant, registry) | CLI | cli/task-list-types | P0 |
| TC-002 | Task 2 AC (prompt.go mapping) | CLI | cli/prompt-get | P0 |
| TC-003 | Task 2 AC (valid type map) | CLI | cli/task-validate-index | P0 |
| TC-004 | Task 2 AC (T-quick-6 generation) | CLI | cli/task-index | P0 |
| TC-005 | Task 2 AC (T-quick-6 deps) | CLI | cli/task-index | P0 |
| TC-006 | Task 2 AC (task order) | CLI | cli/task-index | P1 |
| TC-007 | Task 3 AC (T-test-5 drift description) | CLI | cli/task-index | P1 |
| TC-008 | Task 1 AC (Steps 9-11 exist) | CLI | cli/skill-verification | P0 |
| TC-009 | Task 1 AC (Step 9 classification) | CLI | cli/skill-verification | P0 |
| TC-010 | Task 1 AC (Step 10 ID preservation) | CLI | cli/skill-verification | P0 |
| TC-011 | Task 1 AC (HARD-GATE drift exception) | CLI | cli/skill-verification | P0 |
| TC-012 | Task 1 AC (drift-only mode) | CLI | cli/skill-verification | P0 |
| TC-013 | Task 3 AC (breakdown-tasks drift) | CLI | cli/skill-verification | P1 |
| TC-014 | Task 4 AC (guide T-quick-6) | CLI | cli/skill-verification | P1 |
| TC-015 | Task 4 AC (guide specs drift) | CLI | cli/skill-verification | P2 |
| TC-016 | Task 2 AC (strategy template) | CLI | cli/skill-verification | P0 |
| TC-017 | Task 2 AC (existing tests pass) | CLI | cli/test-suite | P0 |
| TC-018 | Task 2 AC (type inference) | CLI | cli/task-index | P0 |
| TC-019 | Task 1 AC (workflow diagram) | CLI | cli/skill-verification | P1 |
