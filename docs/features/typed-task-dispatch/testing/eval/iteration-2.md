---
date: "2026-05-11"
doc_dir: "docs/features/typed-task-dispatch/testing/"
iteration: "2"
target_score: "90"
evaluator: Claude (main session, adversarial)
---

# Test Cases Eval — Iteration 2

**Score: 84/100** (target: 90)

| Dimension | Score | Max |
|-----------|-------|-----|
| PRD Traceability | 22 | 25 |
| Step Actionability | 19 | 25 |
| Route & Element Accuracy | 16 | 20 |
| Completeness | 17 | 20 |
| Structure & ID Integrity | 10 | 10 |
| **TOTAL** | **84** | **100** |

⚠️ Step Actionability 19/25 — still below 20 blocking threshold

## Attack Points

### Attack 1: Step Actionability — 7 TCs still use vague setup verbs
TC-002 Step 1: "Set up a feature with an index.json containing a task with `type: fix`" — no file path, no JSON content.
TC-003 Steps 1-2: "Add a new markdown template file", "Register the new type in the type enum" — no specific paths or commands.
TC-005 Step 1: "Ensure `.forge/state.json` references the current feature" — no concrete command or file content.
TC-008 Step 1: "Prepare an index.json with one or more tasks having `status: in_progress`" — no file path, no JSON.
TC-010 Step 1: "Provide a tech-design document with an ambiguous or novel task description" — no file path, no content.
TC-012 Step 1: "Set up a task in index.json with a missing or invalid type" — no file path, no JSON.
TC-016 Step 1: "Set up index.json with completed tasks in phase 1 and a pending task as the first task of phase 2" — no file path, no JSON.
Fix: Replace every vague verb with `Create <path> with content: {...}` or a concrete shell command.

### Attack 2: Route & Element Accuracy — "sitemap-missing" placeholder in all 20 TCs
All TCs have `Element: sitemap-missing` — this is a warning artifact from gen-test-cases, not a valid field value. The document itself acknowledges this is a pure CLI feature. Replace with `Element: N/A` for all CLI TCs.

### Attack 3: Step Actionability — unverifiable expected results in TC-003 and TC-005
TC-003 Expected: "Go unit tests can be written to cover the new type without modifying task-executor.md or any task template file" — this is a design property, not a verifiable test outcome. No step produces evidence of this.
TC-005 Step 3: "Measure elapsed time" — no measurement method specified. Replace with `time task prompt <id>` and assert on the real output.
Fix: Replace TC-003 Expected design-property sentence with a verifiable assertion (e.g., "Running `go test ./pkg/prompt/...` exits 0 and covers the new type"). Replace TC-005 Step 3 with `time task prompt <id> > /dev/null` and assert elapsed < 0.5s.
