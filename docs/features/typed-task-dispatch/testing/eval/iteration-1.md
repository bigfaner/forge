---
date: "2026-05-11"
doc_dir: "docs/features/typed-task-dispatch/testing/"
iteration: "1"
target_score: "80"
evaluator: Claude (automated, adversarial)
---

# Test Cases Eval — Iteration 1

**Score: 77/100** (target: 80)

```
┌─────────────────────────────────────────────────────────────────┐
│                  TEST CASES QUALITY SCORECARD                     │
├──────────────────────────────────────────────────────────────────┤
│ Dimension                    │ Score    │ Max      │ Status     │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 1. PRD Traceability          │  20      │  25      │ ⚠️          │
│    TC-to-AC mapping          │   7/9    │          │            │
│    Traceability table        │   8/8    │          │            │
│    Reverse coverage          │   5/8    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 2. Step Actionability        │  16      │  25      │ ⚠️ BLOCKING │
│    Steps concrete            │   5/9    │          │            │
│    Expected results          │   6/9    │          │            │
│    Preconditions explicit    │   5/7    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 3. Route & Element Accuracy  │  16      │  20      │ ⚠️          │
│    Routes valid              │   6/7    │          │            │
│    Elements identifiable     │   5/7    │          │            │
│    Consistency               │   5/6    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 4. Completeness              │  15      │  20      │ ⚠️          │
│    Type coverage             │   7/7    │          │            │
│    Boundary cases            │   5/7    │          │            │
│    Integration scenarios     │   3/6    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 5. Structure & ID Integrity  │  10      │  10      │ ✅          │
│    IDs sequential/unique     │   4/4    │          │            │
│    Classification correct    │   3/3    │          │            │
│    Summary matches actual    │   3/3    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ TOTAL                        │  77      │  100     │            │
└──────────────────────────────┴──────────┴──────────┴────────────┘
```

---

## Deductions

| Location | Issue | Penalty |
|----------|-------|---------|
| TC-015, TC-016 Source field | References PRD section headings ("PRD Spec / Functional Specs — task validate command extension") instead of numbered AC identifiers — rubric requires "PRD AC-3.1" level specificity | -2 pts (Dim 1) |
| Reverse coverage | No TC for quick-tasks skill generating type fields (in-scope per prd-spec.md §Scope), no TC for eval-cases permanent exception (main session execution path), no TC for `task prompt --fix-record-missed` as a standalone command, no TC for `.forge/state.json` read failure mode | -3 pts (Dim 1) |
| TC-001 Step 2 | "Run `run-tasks` (or invoke the dispatch path for that task)" — the "or" branch is undefined; a test step must be a single unambiguous action | -1 pt (Dim 2) |
| TC-001, TC-007, TC-009 Steps | "Set up a feature", "Prepare an index.json", "Provide a complete tech-design document" — none specify a file path, a command, or a concrete mechanism | -2 pts (Dim 2) |
| TC-014 Step 1 | "Simulate a task completion where the record file was not written" — no concrete mechanism: delete the file? Never write it? Which command produces this state? | -1 pt (Dim 2) |
| TC-003 Expected | "Go unit tests can be written to cover the new type without modifying task-executor.md or any task template file" — this is a design property, not a verifiable test outcome; no step produces evidence of this | -1 pt (Dim 2) |
| TC-011 Expected | "the prompt content received by task-executor is identical to what run-tasks would produce for the same task" — no verification mechanism is described; no diff step, no capture step | -1 pt (Dim 2) |
| TC-005 Step 3 | "Measure elapsed time" — no measurement method specified (time command? wrapper script?); without this, the 500ms assertion in Expected is unverifiable | -1 pt (Dim 2) |
| TC-002, TC-013 Step 3 | "Inspect the prompt delivered to task-executor" — no mechanism for capturing what was delivered to a subagent; this is an internal parameter, not stdout | -1 pt (Dim 2) |
| TC-003 Preconditions | "A new task type has been added by placing a markdown template file... and registering it in the type enum" — this is identical to Steps 1-2, making the precondition circular | -1 pt (Dim 2) |
| TC-014 Preconditions | "run-tasks is configured to detect this condition" — no specifics on what configuration is required or how to verify it is in place | -1 pt (Dim 2) |
| All 16 TCs Element field | "Element: sitemap-missing" is a warning placeholder value, not a proper N/A declaration; the document's own WARNING block acknowledges this is a missing-sitemap artifact, not an intentional field value | -4 pts (Dim 3) |
| TC-001 Step 2 Route | "Route: N/A" is correct for CLI, but TC-001 Step 2 says "or invoke the dispatch path" without specifying the command pattern — CLI TCs should have explicit command patterns in lieu of routes | -1 pt (Dim 3) |
| Completeness — integration | No TC covers the eval-cases permanent exception (main session execution, not subagent dispatch) — this is the only type with a different execution path and has zero coverage | -2 pts (Dim 4) |
| Completeness — boundary | No TC for `.forge/state.json` missing or unreadable — prd-spec.md §Blocked State Lifecycle explicitly lists "state.json 读取失败" as a blocked trigger; no TC exercises it | -1 pt (Dim 4) |
| Completeness — boundary | No TC for task ID not found in index.json — a basic invalid-input case for `task prompt <id>` | -1 pt (Dim 4) |
| Completeness — integration | No TC covers the full chain: breakdown-tasks generates index.json with type fields → run-tasks dispatches → task-executor executes; the two halves are tested in isolation only | -1 pt (Dim 4) |

---

## Attack Points

### Attack 1: Step Actionability — vague setup verbs make TCs non-executable

**Where**: TC-001 Step 1: "Set up a feature with an index.json containing a task with `type: doc-generation.summary` and status `pending`"; TC-007 Step 1: "Prepare an index.json with tasks of various IDs"; TC-009 Step 1: "Provide a complete tech-design document"; TC-014 Step 1: "Simulate a task completion where the record file was not written"

**Why it's weak**: "Set up", "Prepare", "Provide", and "Simulate" are not executable actions. A test script generator reading these steps cannot produce runnable code. There is no file path, no command, no fixture format, and no mechanism. TC-014 is the worst case: "simulate a task completion where the record file was not written" could mean deleting the file after the fact, never writing it, or corrupting the write — three completely different setups that would exercise different code paths.

**What must improve**: Replace every setup verb with a concrete action. Example: "Create `docs/features/test-feature/tasks/index.json` with the following content: `{...}`" or "Run `task claim` and then manually delete `tasks/records/<id>.md` before running the next step." Every step must be a single, unambiguous, reproducible action.

---

### Attack 2: PRD Traceability — three in-scope behaviors have zero TC coverage

**Where**: prd-spec.md §Scope: "task-cli 新增 `task prompt <id> --fix-record-missed` 模式"; "type == test-pipeline.eval-cases（永久例外，平台限制）：调用 `task prompt <id>` 获取 prompt，在主会话中直接按 prompt 执行"; "breakdown-tasks 和 quick-tasks skill 生成任务时自动设置 type"

**Why it's weak**: The eval-cases permanent exception is the only task type that takes a fundamentally different execution path (main session instead of subagent dispatch). It has zero TCs. The `--fix-record-missed` flag is tested only indirectly through run-tasks behavior (TC-014) but never as a standalone `task prompt` command — the prompt output itself is never verified. quick-tasks is listed in scope alongside breakdown-tasks but TC-009/TC-010 cover only breakdown-tasks.

**What must improve**: Add TC-017 for eval-cases main session execution path (verify task-executor is NOT dispatched as subagent). Add TC-018 for `task prompt <id> --fix-record-missed` stdout content. Add TC-019 for quick-tasks generating type fields. These are not edge cases — they are explicitly in-scope behaviors.

---

### Attack 3: Completeness — integration scenarios stop at component boundaries

**Where**: TC-011 Expected: "the prompt content received by task-executor is identical to what run-tasks would produce for the same task"; TC-013 Step 3: "Inspect the prompt delivered to task-executor and the content of run-tasks.md"

**Why it's weak**: The integration TCs (TC-001, TC-002, TC-011, TC-013) all stop at "inspect the prompt delivered to task-executor" without specifying how to capture that internal parameter. More critically, there is no TC that exercises the full chain: breakdown-tasks generates a typed index.json → run-tasks dispatches via task prompt → task-executor executes and writes a record. The two halves of the system (generation and execution) are tested in isolation. The prd-spec.md §Blocked State Lifecycle also explicitly lists `.forge/state.json` read failure as a blocked trigger, but no TC exercises this path.

**What must improve**: Add a TC for the full generation-to-execution chain. Add a TC for `.forge/state.json` missing/unreadable (verify task is marked blocked with correct blocked_reason). For TCs that assert on "prompt delivered to task-executor", replace with a verifiable mechanism: capture stdout of `task prompt <id>` directly and assert on that, since that is the actual contract.

---

## Previous Issues Check

<!-- Only for iteration > 1 — N/A for iteration 1 -->

---

## Verdict

- **Score**: 77/100
- **Target**: 80/100
- **Gap**: 3 points
- **Step Actionability**: 16/25 ⚠️ BLOCKING — gen-test-scripts is blocked until this dimension reaches 20
- **Action**: Continue to iteration 2. Priority fixes: (1) replace vague setup verbs with concrete commands and file paths in Steps, (2) add TCs for eval-cases permanent exception and quick-tasks, (3) add TC for `.forge/state.json` failure mode.
