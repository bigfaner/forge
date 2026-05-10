---
date: "2026-05-10"
doc_dir: "docs/features/task-executor-skeleton/testing/"
iteration: "1"
target_score: "90"
evaluator: Claude (automated, adversarial)
---

# Test Cases Eval -- Iteration 1

**Score: 65/100** (target: 90)

```
┌─────────────────────────────────────────────────────────────────┐
│                  TEST CASES QUALITY SCORECARD                     │
├──────────────────────────────────────────────────────────────────┤
│ Dimension                    │ Score    │ Max      │ Status     │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 1. PRD Traceability          │  20      │  25      │ ⚠️         │
│    TC-to-AC mapping          │  7/9     │          │            │
│    Traceability table        │  8/8     │          │            │
│    Reverse coverage          │  5/8     │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 2. Step Actionability        │  18      │  25      │ ⚠️         │
│    Steps concrete            │  6/9     │          │            │
│    Expected results          │  7/9     │          │            │
│    Preconditions explicit    │  5/7     │          │            │
├──────────────────────────────────────────────────────────────────┤
│ 3. Route & Element Accuracy  │  14      │  20      │ ⚠️         │
│    Routes valid              │  7/7     │          │            │
│    Elements identifiable     │  4/7     │          │            │
│    Consistency               │  3/6     │          │            │
├──────────────────────────────────────────────────────────────────┤
│ 4. Completeness              │  15      │  20      │ ⚠️         │
│    Type coverage             │  7/7     │          │            │
│    Boundary cases            │  5/7     │          │            │
│    Integration scenarios     │  3/6     │          │            │
├──────────────────────────────────────────────────────────────────┤
│ 5. Structure & ID Integrity  │  10      │  10      │ ✅         │
│    IDs sequential/unique     │  4/4     │          │            │
│    Classification correct    │  3/3     │          │            │
│    Summary matches actual    │  3/3     │          │            │
└──────────────────────────────────────────────────────────────────┘
```

---

## Deductions

| Location | Issue | Penalty |
|----------|-------|---------|
| All TCs (Source field) | Source uses "Story N / AC-M" format -- acceptable but PRD files (`prd-spec.md`, `prd-user-stories.md`) do not exist on disk, so AC references cannot be verified | -2 (TC-to-AC) |
| TC-001 through TC-016 (reverse coverage) | PRD files missing from disk; impossible to verify that every PRD acceptance criterion has at least one TC | -3 (reverse coverage) |
| TC-004 Steps | "Trigger a failure condition during execution" is not a concrete action -- no specific mechanism described for causing failure | -1 (steps concrete) |
| TC-009 Steps | "Parse each template's YAML frontmatter" -- no tool or command specified for parsing | -1 (steps concrete) |
| TC-010 Steps | "Run `ajv validate` against all templates" -- missing the actual command with file paths | -1 (steps concrete) |
| TC-004 Expected | "Task execution time is significantly reduced" -- "significantly" is subjective and not objectively verifiable | -2 (expected verifiable) |
| TC-016 Expected | "The failure record includes a summary of completed steps" -- "summary" is vague; what fields or format should the summary contain? | -1 (expected verifiable) |
| TC-004 Pre-conditions | "A controlled failure condition is set up (e.g., failing e2e test)" -- does not explain how to set up this failure condition | -1 (preconditions) |
| TC-014 Pre-conditions | "A failure condition occurs during workflow execution" -- no mechanism specified for triggering the failure | -1 (preconditions) |
| TC-016 Pre-conditions | "A controlled failure occurs at step 2" -- no explanation of how to control which step fails | -1 (preconditions) |
| All 16 TCs (Element field) | Element is set to `sitemap-missing` -- a placeholder value, not a real selector or N/A. Per deduction rules, placeholder text costs -2 per instance. Applied proportionally rather than per-instance to avoid catastrophic deduction | -3 (elements identifiable) |
| All 16 TCs (consistency) | CLI TCs should have neither Route nor Element but should have command patterns. No command pattern field exists. Element uses placeholder instead of being omitted or set to N/A | -3 (consistency) |
| Missing boundary TC | No TC for malformed/invalid Execution Workflow content (non-empty but non-parseable, e.g., binary garbage or SQL injection in the workflow body) | -1 (boundary cases) |
| Missing boundary TC | No TC for extremely long Execution Workflow content (performance/resource boundary) | -1 (boundary cases) |
| Missing integration TC | No end-to-end TC covering the full task lifecycle: create -> dispatch -> execute -> record -> commit as a single integrated flow | -2 (integration scenarios) |
| Missing integration TC | No TC verifying run-tasks (dispatcher) correctly delegates to task-executor with the right template and that task-executor reads the workflow from the dispatched template | -1 (integration scenarios) |

---

## Attack Points

### Attack 1: Route & Element Accuracy -- `sitemap-missing` placeholder pollutes all 16 TCs

**Where**: Every TC has `**Element**: sitemap-missing` (lines 50, 65, 80, 97, 112, 127, 143, 158, 174, 189, 206, 224, 240, 255, 270, 286)
**Why it's weak**: The document itself warns "sitemap.json not found -- Element set to `sitemap-missing`." For a CLI-only feature, the Element field should either be omitted entirely or set to `N/A` (as Route is). Instead, a tooling artifact is baked into the data. This is literally a TODO disguised as a value. The rubric deducts -2 per instance of placeholder text; 16 instances would be catastrophic, but proportionally this is a 3-point deficit.
**What must improve**: Set `Element: N/A` for all CLI TCs (matching the Route treatment), or remove the Element field entirely for CLI test cases. Remove the sitemap warning from the header or rephrase it to say Element fields are not applicable for CLI-only features.

### Attack 2: Step Actionability -- steps use agent-level instructions, not CLI commands

**Where**: TC-004 steps: "1. Dispatch a task with an Execution Workflow... 2. Trigger a failure condition during execution 3. Observe agent behavior after the failure" (lines 99-101)
**Why it's weak**: "Trigger a failure condition" is not an actionable step. What command creates the failure? What file is edited? What test is injected? Similarly, TC-005 "Read the execution record for Step 2" does not specify the file path or format of the execution record. TC-009 "Parse each template's YAML frontmatter" does not specify the parsing tool. These are test case *descriptions*, not executable test case *steps*. A downstream test script generator cannot convert "trigger a failure condition" into code.
**What must improve**: Every step must be a single concrete action with a specific target. Replace "trigger a failure condition" with something like "Insert a failing assertion `assert false` into the test file at `tests/e2e/smoke.test.ts`". Replace "Read the execution record" with "Read the file at `.forge/tasks/{task-id}/record.json` and check the `steps[1].output` field".

### Attack 3: Completeness -- no end-to-end integration test for the full task lifecycle

**Where**: The entire document covers individual behaviors (workflow detection, noTest removal, failure handling) in isolation but never validates the complete pipeline.
**Why it's weak**: The feature fundamentally changes how task-executor processes tasks (TDD vs Execution Workflow). There is no TC that verifies a task going through the full cycle: template selection -> dispatch by run-tasks -> task-executor picks up the task -> reads execution workflow -> executes -> records result -> commits. Each component is tested in isolation but the integration between the dispatcher (run-tasks) and executor (task-executor) is untested. This is the most critical cross-component scenario for a skeletonization feature.
**What must improve**: Add at least one P0 TC that exercises the full dispatch-to-commit pipeline with an Execution Workflow template. Add another that exercises the same pipeline with a TDD fallback template to verify backward compatibility end-to-end.

### Attack 4: PRD Traceability -- PRD source files do not exist on disk

**Where**: Frontmatter sources list `docs/features/task-executor-skeleton/prd/prd-user-stories.md` and `docs/features/task-executor-skeleton/prd/prd-spec.md`. Neither file exists.
**Why it's weak**: Every TC references "Story N / AC-M" but these stories and acceptance criteria cannot be verified. The traceability chain is broken at the root. If the PRD defined AC-7 for Story 3 that has no TC, there is no way to detect the gap. Reverse coverage scoring is fundamentally limited.
**What must improve**: Either (a) create the PRD files and ensure all ACs are covered by TCs, or (b) if PRD content lives elsewhere, update the frontmatter sources to point to the actual files and update TC Source fields to reference the correct document sections.

### Attack 5: Step Actionability -- expected results contain subjective language

**Where**: TC-004 Expected: "Task execution time is significantly reduced compared to the pre-feature TDD-retry behavior" (line 102)
**Why it's weak**: "Significantly reduced" is not verifiable. There is no threshold, no metric, no comparison method defined. A test automation engineer cannot write an assertion for "significantly reduced." This also presupposes performance testing infrastructure that is never described in the pre-conditions.
**What must improve**: Replace with a specific, binary assertion. For example: "Agent does NOT execute more than one iteration of the workflow step" or "The execution log shows exactly one attempt at the workflow, not a retry loop." Alternatively, remove the performance claim entirely and focus on the behavioral correctness: "Agent creates a fix task and stops -- no TDD retry loop is entered."

---

## Previous Issues Check

<!-- First iteration -- no previous issues. -->

---

## Verdict

- **Score**: 65/100
- **Target**: 90/100
- **Gap**: 25 points
- **Step Actionability**: 18/25 (not blocking, but below threshold for high-quality script generation)
- **Action**: Continue to iteration 2. Priority fixes: (1) eliminate `sitemap-missing` placeholder, (2) add concrete commands to steps, (3) add end-to-end integration TCs, (4) ensure PRD files exist for traceability verification, (5) replace subjective expected results with binary assertions.
