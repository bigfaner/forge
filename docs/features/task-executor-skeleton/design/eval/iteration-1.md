---
date: "2026-05-10"
doc_dir: "docs/features/task-executor-skeleton/design/"
iteration: 1
target_score: 90
evaluator: Claude (automated, adversarial)
---

# Design Eval -- Iteration 1

**Score: 78/100** (target: 90)

```
+-----------------------------------------------------------------+
|                     DESIGN QUALITY SCORECARD                     |
+------------------------------+----------+----------+------------+
| Dimension                    | Score    | Max      | Status     |
+------------------------------+----------+----------+------------+
| 1. Architecture Clarity      |  19      |  20      | OK         |
|    Layer placement explicit  |  6/7     |          |            |
|    Component diagram present |  7/7     |          |            |
|    Dependencies listed       |  6/6     |          |            |
+------------------------------+----------+----------+------------+
| 2. Interface & Model Defs    |  13      |  20      | WARN       |
|    Interface signatures typed|  4/7     |          |            |
|    Models concrete           |  5/7     |          |            |
|    Directly implementable    |  4/6     |          |            |
+------------------------------+----------+----------+------------+
| 3. Error Handling            |  10      |  15      | WARN       |
|    Error types defined       |  2/5     |          |            |
|    Propagation strategy clear|  3/5     |          |            |
|    HTTP status codes mapped  |  5/5     |          | N/A        |
+------------------------------+----------+----------+------------+
| 4. Testing Strategy          |  8       |  15      | WARN       |
|    Per-layer test plan       |  4/5     |          |            |
|    Coverage target numeric   |  0/5     |          |            |
|    Test tooling named        |  4/5     |          |            |
+------------------------------+----------+----------+------------+
| 5. Breakdown-Readiness *     |  18      |  20      | OK         |
|    Components enumerable     |  7/7     |          |            |
|    Tasks derivable           |  6/7     |          |            |
|    PRD AC coverage           |  5/6     |          |            |
+------------------------------+----------+----------+------------+
| 6. Security Considerations   |  10      |  10      | N/A        |
|    Threat model present      |  N/A     |          |            |
|    Mitigations concrete      |  N/A     |          |            |
+------------------------------+----------+----------+------------+
| TOTAL                        |  78      |  100     |            |
+------------------------------+----------+----------+------------+
```

* Breakdown-Readiness >= 18/20 -- can proceed to /breakdown-tasks

---

## Deductions

| Location | Issue | Penalty |
|----------|-------|---------|
| Interfaces (whole section) | Interface signatures are prose descriptions, not typed signatures. Interface 3 shows claim output as space-delimited strings, not typed fields. | -3 pts (Dim 2) |
| Data Models | No typed model for the core new concept (Execution Workflow content). "embedded in task template markdown" is a design rationale, not a model. | -2 pts (Dim 2) |
| Interface 4: record behavior | "require test evidence OR explicit --force" is a behavioral description, not code-level specification. Developer must infer implementation. | -2 pts (Dim 2) |
| Error Handling | No custom error types, error codes, or error enums defined. Table describes situations, not error types. | -3 pts (Dim 3) |
| Error Handling | No explicit propagation strategy statement. PRD Flow Description has more error detail than the design doc. | -2 pts (Dim 3) |
| Testing Strategy | No numeric coverage target anywhere in the document. | -5 pts (Dim 4) |
| Testing Strategy | Agent prompt tests are "Manual" only. No automated verification approach. | -1 pts (Dim 4) |
| Breakdown-Readiness | PRD Flow Description mentions "workflow declares failure instructions" sub-case (e.g., T-test-3 "create fix task") -- Error Handling table does not explicitly cover this. | -1 pts (Dim 5) |
| Interface 4 | Record behavior change lacks task-derivable specificity -- no concrete code diff or function signature. | -1 pts (Dim 5) |

---

## Attack Points

### Attack 1: Interface & Model Definitions -- untyped interfaces

**Where**: Interface 3 shows `BEFORE: KEY, TASK_ID, FILE, BREAKING, MAIN_SESSION, NO_TEST, SCOPE, FEATURE` and `AFTER: KEY, TASK_ID, FILE, BREAKING, MAIN_SESSION, SCOPE, FEATURE`
**Why it's weak**: These are space-delimited string formats, not typed interfaces. There is no Go struct definition, no function signature, no typed parameter list. A developer must open the existing code to understand field types, ordering, and parsing logic. Interface 1 (agent prompt) is well-specified as copy-paste markdown, but Interfaces 3 and 4 are prose behavioral changes.
**What must improve**: Add concrete Go struct/type definitions for the claim output and record behavior. Show actual function signature changes (e.g., `func (r *RecordCmd) Execute(task Task) error` with the updated Task struct). Provide before/after code snippets, not just prose.

### Attack 2: Error Handling -- no error types or codes

**Where**: Error Handling table lists 7 error cases but none have named error types. Quote: `Status = 'failed', log error, skip Step 2` -- this is a behavior description, not an error definition.
**Why it's weak**: The rubric requires "custom error types or error codes explicitly defined." The design describes *situations* where errors occur and *behaviors* to take, but never defines what `Status = 'failed'` is as a type, what error codes exist, or how errors are represented in the record system. There is no `TaskStatus` enum, no `WorkflowError` type, no error codes like `ERR_WORKFLOW_EMPTY` or `ERR_TASK_UNREADABLE`.
**What must improve**: Define a `TaskStatus` enum/type with values `pending`, `in_progress`, `completed`, `failed`. Define error codes for each error case in the table. Specify how errors are persisted (in the task record? in a log file?).

### Attack 3: Testing Strategy -- missing coverage target

**Where**: Testing Strategy section, Per-Layer Test Plan table. No mention of coverage percentage anywhere.
**Why it's weak**: The rubric requires a "numeric coverage target (e.g., 80%)." The document lists test types and scenarios but never quantifies what "enough" means. For the Go code removal (which touches types.go, claim.go, record.go, errors.go), there is no statement like "all existing tests must pass, new tests cover record behavior change." The phrase "remaining tests pass" implies regression but not coverage.
**What must improve**: Add explicit coverage target for task-cli Go code (e.g., "maintain >= 80% coverage on changed files" or "all existing test cases pass, plus 2 new test cases for record behavior change"). Define pass criteria for manual agent tests.

---

## Previous Issues Check

N/A -- Iteration 1.

---

## Verdict

- **Score**: 78/100
- **Target**: 90/100
- **Gap**: 12 points
- **Breakdown-Readiness**: 18/20 -- can proceed to /breakdown-tasks
- **Action**: Continue to iteration 2. Primary gaps are Testing Strategy (-7), Interface & Model Definitions (-7), and Error Handling (-5). These are addressable without architectural changes.
