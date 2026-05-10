---
date: "2026-05-10"
doc_dir: "docs/features/task-executor-skeleton/testing/"
iteration: "2"
target_score: "90"
evaluator: Claude (automated, adversarial)
---

# Test Cases Eval — Iteration 2

**Score: 91/100** (target: 90)

```
┌─────────────────────────────────────────────────────────────────┐
│                  TEST CASES QUALITY SCORECARD                     │
├──────────────────────────────────────────────────────────────────┤
│ Dimension                    │ Score    │ Max      │ Status     │
├──────────────────────────────┼──────────┼──────────┤
│ 1. PRD Traceability          │  24      │  25      │ ⚠️         │
│    TC-to-AC mapping          │  8/9     │          │            │
│    Traceability table        │  8/8     │          │            │
│    Reverse coverage          │  8/8     │          │            │
├──────────────────────────────┼──────────┼──────────┤
│ 2. Step Actionability        │  22      │  25      │ ⚠️         │
│    Steps concrete            │  8/9     │          │            │
│    Expected results          │  8/9     │          │            │
│    Preconditions explicit    │  6/7     │          │            │
├──────────────────────────────┼──────────┼──────────┤
│ 3. Route & Element Accuracy  │  20      │  20      │ ✅         │
│    Routes valid              │  7/7     │          │            │
│    Elements identifiable     │  7/7     │          │            │
│    Consistency               │  6/6     │          │            │
├──────────────────────────────┼──────────┼──────────┤
│ 4. Completeness              │  15      │  20      │ ⚠️         │
│    Type coverage             │  7/7     │          │            │
│    Boundary cases            │  5/7     │          │            │
│    Integration scenarios     │  3/6     │          │            │
├──────────────────────────────┼──────────┤
│ 5. Structure & ID Integrity  │  10      │  10      │ ✅         │
│    IDs sequential/unique     │  4/4     │          │            │
│    Classification correct    │  3/3     │          │            │
│    Summary matches actual    │  3/3     │          │
├──────────────────────────────┼──────────┤
│ TOTAL                        │  91      │  100     │            │
└──────────────────────────────┴──────────┴──────────┘
```

---

## Deductions

| Location | Issue | Penalty |
|----------|-------|---------|
| TC-007 to TC-012 (TC-to-AC) | No TC explicitly targets the three skill docs (record-task, quick-tasks, consolidate-specs) listed in PRD Scope. TC-007's broad grep covers them implicitly, but no TC Source field claims coverage of skill doc cleanup | -1 (TC-to-AC) |
| TC-009 Step 4 | "For each template listed (if any), run `head -20 {file}`" — `{file}` is an unresolved variable, not a concrete action. Should say "for each file returned by the grep in step 3" | -1 (steps concrete) |
| TC-006 Expected | "The `grep -c` for 'compile', 'fmt', 'lint', 'test' in the record returns `0` or the matches are limited to the Execution Workflow body content" — the "or" clause creates an ambiguity. How does the tester distinguish a Quality Gate log match from a workflow body content match? No mechanism specified | -1 (expected verifiable) |
| TC-002 Expected | "Behavior is identical to the pre-feature baseline" — "identical to baseline" is not a self-contained assertion. The tester must possess external knowledge of what the baseline looks like. While steps 3-4 provide concrete checks, the expected result text itself is not independently verifiable | -1 (expected verifiable) |
| TC-017 Pre-conditions | "The working tree is clean" — no command given to verify or enforce this state. A pre-condition should either be verifiable by a step or paired with a setup action | -1 (preconditions) |
| TC-006 Expected | See above — the "or matches are limited to" escape hatch is under-specified | counted above |
| All boundary TCs | No TC for malformed/corrupt Execution Workflow content (binary garbage, extremely long strings). PRD spec says "正文非空" triggers injection — what happens with non-UTF8 or 10MB body content? Iteration 1 flagged this; partially addressed by TC-003 (empty body) but not for pathological content | -2 (boundary cases) |
| Integration coverage | No TC covers TDD-fallback end-to-end through the full run-tasks → task-executor → record → commit pipeline. TC-002 tests fallback detection in isolation; TC-017 tests workflow end-to-end. The backward-compatibility path (TDD fallback through the full pipeline) is untested at integration level | -2 (integration scenarios) |
| TC-016 | `completed_steps` and `failure_point` are referenced as JSON fields in steps and expected results, but no PRD document defines the record.json schema. The tester is assumed to know these field names | -1 (steps concrete) |

---

## Attack Points

### Attack 1: Completeness — missing TDD-fallback end-to-end integration TC

**Where**: TC-017 covers the workflow path end-to-end, but no TC covers the TDD-fallback path end-to-end. TC-002 only checks Step 2/Step 3 output in isolation.
**Why it's weak**: The PRD's primary backward-compatibility promise is "无 Execution Workflow 的旧任务回退到 TDD，行为不变". This is tested at the unit level (TC-002 detects the fallback) but never verified through the full dispatch-to-commit pipeline. If run-tasks dispatches a task without an Execution Workflow template, and task-executor processes it, the entire TDD + Quality Gate + commit path must work. This is a distinct cross-component scenario from TC-017 (which tests the workflow path). The absence means the backward-compatibility guarantee is only partially verified.
**What must improve**: Add a TC-018 that mirrors TC-017 but uses a template WITHOUT `## Execution Workflow`. Dispatch through run-tasks, verify task-executor executes TDD + Quality Gate, verify record shows TDD steps, verify commit succeeds. This is the most important missing integration scenario.

### Attack 2: Step Actionability — TC-006 expected result has an untestable escape hatch

**Where**: TC-006 Expected (line 133): "The `grep -c` for 'compile', 'fmt', 'lint', 'test' in the record returns `0` or the matches are limited to the Execution Workflow body content, not Quality Gate execution logs."
**Why it's weak**: The "or" clause makes this assertion non-binary. If grep returns 3 matches for "test", the tester must now determine whether those matches came from the workflow body or from Quality Gate execution. No mechanism is provided for making this distinction. The expected result should be a single, unambiguous pass/fail criterion. A downstream test script cannot programmatically distinguish "workflow body content matches" from "Quality Gate log matches" without additional structural context (e.g., checking which step array index contains the match).
**What must improve**: Replace with: "grep -c returns 0 for all four keywords when run against steps[2].output specifically (the Step 3 output). If the Execution Workflow body itself contains these keywords, they appear only in steps[1].output, not steps[2].output." This scopes the check to a specific step index, making it binary.

### Attack 3: Completeness — no boundary TC for pathological Execution Workflow content

**Where**: TC-003 covers empty body. No TC covers non-empty but problematic body content.
**Why it's weak**: The PRD says "正文非空" triggers injection. What happens when the body is a 10MB string, or contains binary/null bytes, or includes control characters that break the agent prompt? The task-executor reads the body and injects it into an agent prompt — there is no documented size limit or sanitization step. TC-003 tests the empty boundary but not the large/malformed boundary. This is a real operational risk for an LLM-prompt-injection pipeline.
**What must improve**: Add a TC for a very large Execution Workflow body (e.g., >50KB) that verifies the task-executor either processes it correctly or fails gracefully with a specific error. This does not need to be P0 but should be P1 to prevent prompt-length regressions.

---

## Previous Issues Check

| Previous Attack (Iteration 1) | Addressed? | Evidence |
|-------------------------------|------------|----------|
| Attack 1: `sitemap-missing` placeholder in all 16 TC Element fields | ✅ Yes | All 17 TCs now show `**Element**: N/A`. The placeholder is completely eliminated. |
| Attack 2: Steps use agent-level instructions, not concrete CLI commands | ✅ Yes | Steps now use concrete commands: `task-cli execute-task --task-file`, `cat ... \| jq`, `grep -ri`, `ajv validate`, `find ... \| wc -l`. Failure mechanisms specify exact file and content (`assert(false)` in `tests/e2e/smoke.test.ts`, `exit 1` in multi-step workflow). |
| Attack 3: No end-to-end integration TC for full task lifecycle | ✅ Yes | TC-017 added: "Full dispatch-to-commit pipeline with Execution Workflow template" — exercises run-tasks → task-executor → workflow execution → record → commit with 6 verification steps. P0 priority. |
| Attack 4: PRD source files do not exist on disk | ✅ Yes | PRD files now exist at `docs/features/task-executor-skeleton/prd/prd-spec.md` and `prd/user-stories.md`. All Story/AC references verified against actual PRD content. |
| Attack 5: Expected results contain subjective language ("significantly reduced") | ✅ Yes | TC-004 Expected now reads: "record.json shows exactly one attempt at the workflow step -- no retry entries... The steps array does NOT contain any entries with 'RED', 'GREEN', 'REFACTOR', or 'TDD cycle'". Fully binary and verifiable. |

---

## Verdict

- **Score**: 91/100
- **Target**: 90/100
- **Gap**: 0 points (target met)
- **Step Actionability**: 22/25 (above 20 blocking threshold)
- **Action**: Target reached. Remaining improvements are optional: add TDD-fallback end-to-end TC (TC-018), tighten TC-006 expected result, add pathological-workflow-content boundary TC.
