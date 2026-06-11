---
feature: "quick-test-slim"
sources:
  - docs/proposals/quick-test-slim/proposal.md
  - docs/features/quick-test-slim/tasks/1-merge-gen-run.md
generated: "2026-05-16"
---

# Test Cases: quick-test-slim

## Summary

| Type | Count |
|------|-------|
| UI   | 0   |
| **Integration** | **0** |
| API  | 0  |
| CLI  | 16  |
| **Total** | **16** |

---

## CLI Test Cases

## TC-001: Quick mode single profile generates correct task count
- **Source**: Proposal Success Criteria — "Quick mode generates 4 test tasks (flat profile) or 5 (nested profile), no longer 6"
- **Type**: CLI
- **Target**: cli/testgen
- **Test ID**: cli/testgen/quick-mode-single-profile-task-count
- **Pre-conditions**: forge CLI built, no existing feature with slug `test-quick-slim-001`
- **Steps**:
  1. Run `forge quick` with a single `go-test` profile and a simple feature proposal
  2. Parse the generated manifest to count test tasks
  3. Verify the count is 5 (gen-cases, gen-and-run, graduate, verify-regression, drift-detection)
- **Expected**: Quick mode with single profile generates exactly 5 test pipeline tasks
- **Priority**: P0

## TC-002: Quick mode merged task has gen-and-run type
- **Source**: Task 1 AC — "GetQuickTestTasks() generates a merged gen-and-run task instead of separate gen-scripts + run tasks"
- **Type**: CLI
- **Target**: cli/testgen
- **Test ID**: cli/testgen/merged-task-has-gen-and-run-type
- **Pre-conditions**: forge CLI built
- **Steps**:
  1. Call `GetQuickTestTasks` with `profiles=["go-test"]` and `detectedTypes=[]`
  2. Find the task at index 1 (T-quick-2)
  3. Verify its Type field equals `"test-pipeline.gen-and-run"`
- **Expected**: The merged task has type `test-pipeline.gen-and-run`, not `test-pipeline.gen-scripts` or `test-pipeline.run`
- **Priority**: P0

## TC-003: Quick mode merged task generates correct prompt template
- **Source**: Task 1 AC — "prompt.go maps TypeTestPipelineGenAndRun to data/test-pipeline-gen-and-run.md"
- **Type**: CLI
- **Target**: cli/testgen
- **Test ID**: cli/testgen/merged-task-prompt-template-mapping
- **Pre-conditions**: forge CLI built, type-to-template mapping loaded
- **Steps**:
  1. Look up the type-to-template map for `TypeTestPipelineGenAndRun`
  2. Verify it maps to `data/test-pipeline-gen-and-run.md`
  3. Verify that template file exists on disk
- **Expected**: The merged type resolves to a valid prompt template file
- **Priority**: P0

## TC-004: Quick mode per-type creates independent gen-and-run tasks
- **Source**: Proposal Success Criteria — "per-type mode T-quick-2-tui and T-quick-2-api each independently generate and run"
- **Type**: CLI
- **Target**: cli/testgen
- **Test ID**: cli/testgen/per-type-creates-independent-gen-and-run
- **Pre-conditions**: forge CLI built
- **Steps**:
  1. Call `GetQuickTestTasks` with `profiles=["go-test"]` and `detectedTypes=["tui", "api"]`
  2. Find tasks with IDs matching `T-quick-2-tui` and `T-quick-2-api`
  3. Verify both have `Type = "test-pipeline.gen-and-run"`
  4. Verify both have correct `TestType` fields ("tui" and "api" respectively)
- **Expected**: Per-type mode creates separate gen-and-run tasks for each detected type
- **Priority**: P0

## TC-005: Quick mode dependency chain is correct after merge
- **Source**: Task 1 AC — "resolveQuickDeps() dependency chain updated: verify-regression depends on the merged task (or graduate if present), drift-detection depends on verify-regression"
- **Type**: CLI
- **Target**: cli/testgen
- **Test ID**: cli/testgen/dependency-chain-correct-after-merge
- **Pre-conditions**: forge CLI built
- **Steps**:
  1. Call `GetQuickTestTasks` with `profiles=["go-test"]` and `detectedTypes=[]`
  2. Verify T-quick-2 depends on T-quick-1
  3. Verify T-quick-3 (graduate) depends on T-quick-2
  4. Verify T-quick-4 (verify-regression) depends on T-quick-3
  5. Verify T-quick-5 (drift-detection) depends on T-quick-4
- **Expected**: The dependency chain is gen-cases -> gen-and-run -> graduate -> verify-regression -> drift-detection
- **Priority**: P0

## TC-006: Quick mode per-type dependency chain fans in correctly
- **Source**: Proposal Success Criteria — "verify-regression correctly depends on T-quick-2 (merged) rather than independent run task"
- **Type**: CLI
- **Target**: cli/testgen
- **Test ID**: cli/testgen/per-type-dependency-fan-in
- **Pre-conditions**: forge CLI built
- **Steps**:
  1. Call `GetQuickTestTasks` with `profiles=["go-test"]` and `detectedTypes=["tui", "api"]`
  2. Verify T-quick-2-tui and T-quick-2-api both depend on T-quick-1
  3. Verify T-quick-3 (graduate) depends on both T-quick-2-tui AND T-quick-2-api
  4. Verify T-quick-4 depends on T-quick-3
- **Expected**: Per-type tasks fan in to graduate, which fans in to verify-regression
- **Priority**: P0

## TC-007: Breakdown mode is unchanged by quick mode merge
- **Source**: Task 1 Hard Rules — "Do NOT modify breakdown mode task generation — this change is quick-mode only"
- **Type**: CLI
- **Target**: cli/testgen
- **Test ID**: cli/testgen/breakdown-mode-unchanged
- **Pre-conditions**: forge CLI built
- **Steps**:
  1. Call `GetBreakdownTestTasks` with `profiles=["go-test"]` and `detectedTypes=[]`
  2. Verify T-test-2 has type `test-pipeline.gen-scripts` (not gen-and-run)
  3. Verify T-test-3 has type `test-pipeline.run`
  4. Verify task count matches the original breakdown structure (7 tasks total)
- **Expected**: Breakdown mode still generates separate gen-scripts and run tasks, unaffected by the merge
- **Priority**: P0

## TC-008: Quick mode multi-profile with letter suffixes works
- **Source**: Task 1 AC — "Multi-profile: letter suffixes work correctly with merged task"
- **Type**: CLI
- **Target**: cli/testgen
- **Test ID**: cli/testgen/multi-profile-letter-suffixes
- **Pre-conditions**: forge CLI built
- **Steps**:
  1. Call `GetQuickTestTasks` with `profiles=["go-test", "web-playwright"]` and `detectedTypes=[]`
  2. Verify T-quick-1a and T-quick-1b exist with type `gen-cases`
  3. Verify T-quick-2a and T-quick-2b exist with type `gen-and-run`
  4. Verify T-quick-3a and T-quick-3b exist with type `graduate`
  5. Verify shared tasks T-quick-4 and T-quick-5 exist
- **Expected**: Multi-profile mode generates suffixed tasks with correct types, merged task replaces separate gen/run
- **Priority**: P1

## TC-009: Merged prompt template calls gen then run sequentially
- **Source**: Task 1 AC — "New prompt template test-pipeline-gen-and-run.md calls /gen-test-scripts then /run-e2e-tests sequentially, with in-session fix loop"
- **Type**: CLI
- **Target**: cli/prompt-template
- **Test ID**: cli/prompt-template/merged-template-calls-both-skills
- **Pre-conditions**: Prompt template file exists at `forge-cli/pkg/prompt/data/test-pipeline-gen-and-run.md`
- **Steps**:
  1. Read the prompt template file
  2. Verify it contains references to `gen-test-scripts` skill
  3. Verify it contains references to `run-e2e-tests` skill
  4. Verify gen phase appears before run phase in the document structure
  5. Verify it includes instructions for in-session fix loop (retry on failure)
- **Expected**: The merged prompt template orchestrates both skills sequentially with failure recovery
- **Priority**: P0

## TC-010: Type constant registered in types.go
- **Source**: Task 1 AC — "New type constant TypeTestPipelineGenAndRun = test-pipeline.gen-and-run registered in types.go (const block + registry + validTypes)"
- **Type**: CLI
- **Target**: cli/types
- **Test ID**: cli/types/gen-and-run-constant-registered
- **Pre-conditions**: forge CLI source code available
- **Steps**:
  1. Verify `TypeTestPipelineGenAndRun` constant exists in `types.go` with value `"test-pipeline.gen-and-run"`
  2. Verify it is registered in the `TaskTypeRegistry` slice
  3. Verify it is present in the `validTypes` map
- **Expected**: The type constant is fully registered in all required data structures
- **Priority**: P0

## TC-011: InferType maps merged IDs correctly
- **Source**: Task 1 AC — "infer.go maps merged IDs to TypeTestPipelineGenAndRun (handles profile suffix and type suffix patterns)"
- **Type**: CLI
- **Target**: cli/infer
- **Test ID**: cli/infer/merged-ids-mapped-correctly
- **Pre-conditions**: forge CLI built
- **Steps**:
  1. Call `InferType("T-quick-2")` — verify returns `TypeTestPipelineGenAndRun`
  2. Call `InferType("T-quick-2-api")` — verify returns `TypeTestPipelineGenAndRun`
  3. Call `InferType("T-quick-2-tui")` — verify returns `TypeTestPipelineGenAndRun`
  4. Call `InferType("T-quick-2a-api")` — verify returns `TypeTestPipelineGenAndRun`
  5. Call `InferType("T-quick-2b-tui")` — verify returns `TypeTestPipelineGenAndRun`
- **Expected**: All merged task ID patterns resolve to `TypeTestPipelineGenAndRun`
- **Priority**: P0

## TC-012: Quick mode single profile produces 5 tasks total
- **Source**: Task 1 AC — "Single profile (no per-type): 5 test tasks total (gen-cases, gen-and-run, graduate, verify-regression, drift-detection)"
- **Type**: CLI
- **Target**: cli/testgen
- **Test ID**: cli/testgen/single-profile-five-tasks
- **Pre-conditions**: forge CLI built
- **Steps**:
  1. Call `GetQuickTestTasks` with `profiles=["go-test"]` and `detectedTypes=[]`
  2. Count the returned tasks
- **Expected**: Exactly 5 tasks returned for single profile without per-type
- **Priority**: P1

## TC-013: Existing gen-scripts and run prompt templates remain intact
- **Source**: Task 1 Implementation Notes — "Existing prompt templates (gen-scripts.md, run.md) remain for breakdown mode — only quick mode uses the merged template"
- **Type**: CLI
- **Target**: cli/prompt-template
- **Test ID**: cli/prompt-template/existing-templates-unchanged
- **Pre-conditions**: forge CLI source code available
- **Steps**:
  1. Verify `data/test-pipeline-gen-scripts.md` still exists and is unchanged
  2. Verify `data/test-pipeline-run.md` still exists and is unchanged
  3. Verify breakdown mode still references these templates (not the merged one)
- **Expected**: Original prompt templates for gen-scripts and run remain for breakdown mode
- **Priority**: P1

## TC-014: Merged task generates correct task .md file
- **Source**: Task 1 AC — prompt template uses standard placeholders
- **Type**: CLI
- **Target**: cli/testgen
- **Test ID**: cli/testgen/merged-task-generates-correct-md
- **Pre-conditions**: forge CLI built
- **Steps**:
  1. Create a `TestTaskDef` with type `TypeTestPipelineGenAndRun`
  2. Call `GenerateTestTaskMD` with the definition
  3. Verify the generated .md file contains correct frontmatter (id, title, type, profile)
  4. Verify it contains the profile-specific strategy content
- **Expected**: Merged task generates a valid .md file with all required fields
- **Priority**: P1

## TC-015: DetectTypesFromTestCases correctly parses test-cases.md summary table
- **Source**: Task 1 AC — test cases cover merged task count and dependency chain
- **Type**: CLI
- **Target**: cli/testgen
- **Test ID**: cli/testgen/detect-types-from-test-cases
- **Pre-conditions**: forge CLI built
- **Steps**:
  1. Create a test-cases.md content with summary table showing `CLI: 5, API: 3`
  2. Call `DetectTypesFromTestCases` with the content
  3. Verify it returns `["cli", "api"]`
  4. Create content with `Total: 10` only
  5. Verify it returns nil (no types with count > 0)
- **Expected**: Type detection correctly parses summary table and returns types with non-zero counts
- **Priority**: P1

## TC-016: Subagent completes gen-and-run in single session
- **Source**: Proposal Success Criteria — "merged task's subagent completes generation and running in a single session (including failure recovery)"
- **Type**: CLI
- **Target**: cli/prompt-template
- **Test ID**: cli/prompt-template/single-session-gen-and-run
- **Pre-conditions**: forge CLI built, feature with test cases ready
- **Steps**:
  1. Execute a merged gen-and-run task via the task executor
  2. Verify the task produces both generated test scripts and test execution results
  3. Verify both outputs come from the same task record
- **Expected**: A single task execution produces both generated scripts and test results
- **Priority**: P2

---

## Traceability

| TC ID | Source | Type | Target | Priority |
|-------|--------|------|--------|----------|
| TC-001 | Proposal Success Criteria | CLI | cli/testgen | P0 |
| TC-002 | Task 1 AC | CLI | cli/testgen | P0 |
| TC-003 | Task 1 AC | CLI | cli/testgen | P0 |
| TC-004 | Proposal Success Criteria | CLI | cli/testgen | P0 |
| TC-005 | Task 1 AC | CLI | cli/testgen | P0 |
| TC-006 | Proposal Success Criteria | CLI | cli/testgen | P0 |
| TC-007 | Task 1 Hard Rules | CLI | cli/testgen | P0 |
| TC-008 | Task 1 AC | CLI | cli/testgen | P1 |
| TC-009 | Task 1 AC | CLI | cli/prompt-template | P0 |
| TC-010 | Task 1 AC | CLI | cli/types | P0 |
| TC-011 | Task 1 AC | CLI | cli/infer | P0 |
| TC-012 | Task 1 AC | CLI | cli/testgen | P1 |
| TC-013 | Task 1 Implementation Notes | CLI | cli/prompt-template | P1 |
| TC-014 | Task 1 AC | CLI | cli/testgen | P1 |
| TC-015 | Task 1 AC | CLI | cli/testgen | P1 |
| TC-016 | Proposal Success Criteria | CLI | cli/prompt-template | P2 |

> No route validation section: this feature is a Go CLI tool with no web routes. The project interfaces are CLI commands and Go function calls, not HTTP endpoints.
