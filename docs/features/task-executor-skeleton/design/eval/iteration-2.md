---
date: "2026-05-10"
doc_dir: "docs/features/task-executor-skeleton/design/"
iteration: 2
target_score: 90
evaluator: Claude (automated, adversarial)
---

# Design Eval — Iteration 2

**Score: 87/100** (target: 90)

```
┌─────────────────────────────────────────────────────────────────┐
│                     DESIGN QUALITY SCORECARD                     │
├──────────────────────────────┬──────────┬──────────┬────────────┤
│ Dimension                    │ Score    │ Max      │ Status     │
├──────────────────────────────┼──────────┼──────────┬────────────┤
│ 1. Architecture Clarity      │  20      │  20      │ OK         │
│    Layer placement explicit  │  7/7     │          │            │
│    Component diagram present │  7/7     │          │            │
│    Dependencies listed       │  6/6     │          │            │
├──────────────────────────────┼──────────┼──────────┬────────────┤
│ 2. Interface & Model Defs    │  17      │  20      │ WARN       │
│    Interface signatures typed│  6/7     │          │            │
│    Models concrete           │  6/7     │          │            │
│    Directly implementable    │  5/6     │          │            │
├──────────────────────────────┼──────────┼──────────┬────────────┤
│ 3. Error Handling            │  13      │  15      │ WARN       │
│    Error types defined       │  4/5     │          │            │
│    Propagation strategy clear│  4/5     │          │            │
│    HTTP status codes mapped  │  5/5     │          │            │ N/A
├──────────────────────────────┼──────────┼──────────┬────────────┤
│ 4. Testing Strategy          │  13      │  15      │ WARN       │
│    Per-layer test plan       │  4/5     │          │            │
│    Coverage target numeric   │  5/5     │          │            │
│    Test tooling named        │  4/5     │          │            │
├──────────────────────────────┼──────────┼──────────┬────────────┤
│ 5. Breakdown-Readiness ★     │  18      │  20      │ OK         │
│    Components enumerable     │  7/7     │          │            │
│    Tasks derivable           │  6/7     │          │            │
│    PRD AC coverage           │  5/6     │          │            │
├──────────────────────────────┼──────────┼──────────┬────────────┤
│ 6. Security Considerations   │  10      │  10      │ N/A        │
│    Threat model present      │  N/A     │          │            │
│    Mitigations concrete      │  N/A     │          │            │
├──────────────────────────────┼──────────┼──────────┬────────────┤
│ TOTAL                        │  87      │  100     │            │
└──────────────────────────────┴──────────┴──────────┴────────────┘
```

★ Breakdown-Readiness >= 18/20 — can proceed to /breakdown-tasks

---

## Deductions

| Location | Issue | Penalty |
|----------|-------|---------|
| Interface 2 (rules for workflow content) | Workflow content rules are prose enforcement criteria, not a typed model. No schema or formal definition of what constitutes a valid workflow step. | -1 pts (Dim 2) |
| Data Models (workflow content) | "The `## Execution Workflow` content is embedded in task template markdown — no schema field needed" — design decision is defensible but means the core new concept has no concrete model. | -1 pts (Dim 2) |
| Interface 2 (workflow validation) | "Rules for workflow content (enforced in template review)" — the enforcement mechanism is unspecified. Template review is a process, not a code interface. Developer must guess how validation works. | -1 pts (Dim 2) |
| Error Cases (no new codes) | Two new error situations (empty workflow warning, default template missing) introduced but no new error codes defined. "Warning only" for Case C is a log message, not a typed error. | -1 pts (Dim 3) |
| Error Propagation (agent→task-cli boundary) | Agent writes `FAILED: [reason]` as freeform text. How task-cli's `task record` determines status from agent output is unspecified. The two-channel model is described but the handoff is not. | -1 pts (Dim 3) |
| Per-Layer Test Plan (agent tests) | Interface 1 defines 3 cases (A, B, C) but agent prompt test plan has only 1 manual test scenario (T-test-3). Cases B and C have no explicit test coverage. | -1 pts (Dim 4) |
| Test tooling (agent prompt automation) | Agent prompt tests are all "Manual" with binary pass/fail. No approach to automated regression testing for prompt behavior. Previous attack called this out; still not addressed. | -1 pts (Dim 4) |
| PRD Coverage (fix task creation) | Error Cases row 6 says agent "creates fix task" and "sets status to blocked" but no interface specifies HOW — no function, no command, no structured output format for fix task creation. | -1 pts (Dim 5) |
| Tasks derivable (workflow validation) | "Rules for workflow content (enforced in template review)" is not a code-derivable task. It is a review checklist. No test or validation code can be written from this specification. | -1 pts (Dim 5) |

---

## Attack Points

### Attack 1: Interface & Model Definitions — workflow content has no typed model

**Where**: Data Models section states "The `## Execution Workflow` content is embedded in task template markdown — no schema field needed." Interface 2 lists 4 "Rules for workflow content" as prose.
**Why it's weak**: The core new concept in this design is the `## Execution Workflow` section. Yet it has no typed model, no schema, no parseable structure. The 4 rules (must include stop condition, must not use open-ended instructions, must specify commands, must specify failure handling) are review criteria for humans, not a specification a developer can validate in code. There is no grammar, no BNF, no example of an invalid workflow, no automated validation. The design says "enforced in template review" — which is a process answer to a technical question.
**What must improve**: Either (a) define a minimal schema for workflow content (e.g., each step must have a command field, a success criteria field, and a failure action field), or (b) explicitly acknowledge this is intentionally unstructured and describe the template review process as a gate. If (b), add a concrete example of an invalid workflow and the human review checklist.

### Attack 2: Error Handling — agent-to-task-cli error boundary is a black box

**Where**: Error Propagation section describes two channels — task-cli uses `*AIError` structs, agent prompts use step output strings like `Step 2/4: ... FAILED: [reason]`. Error Persistence says the `status` frontmatter field is set to `failed`.
**Why it's weak**: The design describes WHAT the agent writes (freeform text) and WHERE errors end up (task record, index.json), but never specifies HOW the status propagates from the agent's freeform text output to the `status` field in the record. Does the agent set `status: failed` in the frontmatter directly? Does task-cli parse the step output string? Does the agent call `task record --status failed`? This is the critical handoff point between the two error channels, and it is completely unspecified. The previous iteration had this same gap (scored 3/5 for propagation); the revision added `*AIError` struct details and error persistence locations but did not close the agent→task-cli boundary question.
**What must improve**: Add one sentence clarifying the mechanism: e.g., "The agent calls `task record --status failed` with the failure reason in the notes field" or "The agent writes `status: failed` into the task record frontmatter, which task-cli reads on next invocation." Specify who sets the status and how.

### Attack 3: Testing Strategy — agent prompt test coverage is incomplete

**Where**: Per-Layer Test Plan table, Agent prompt row: "Manual — Execute T-test-3 — Step 2 output = 'Execution Workflow', not 'TDD implementation'". Key Test Scenarios table has 7 scenarios, of which only scenarios 4-6 test agent prompt behavior.
**Why it's weak**: Interface 1 defines 3 distinct cases (A: workflow present with content, B: no workflow heading, C: empty workflow heading). The test plan covers Case A (scenarios 5, 6) and Case B (scenario 4). Case C (empty workflow → warning + fallback) has NO test scenario. This is a gap because Case C is the one that produces a warning and silently changes behavior — the exact kind of edge case that should be tested. Additionally, all agent tests remain "Manual" with no path toward automated regression.
**What must improve**: Add a test scenario for Case C: "Task file with `## Execution Workflow` heading but empty content → agent logs warning, reads default template, follows TDD + QG steps." Consider whether any agent prompt tests can be automated (e.g., grep for warning string in output).

---

## Previous Issues Check

| Previous Attack | Addressed? | Evidence |
|----------------|------------|----------|
| Attack 1 (untyped interfaces — Go structs/code diffs) | Partially | Interface 3 now has concrete Go struct diffs (types.go, claim.go). Interface 4 has before/after code for record.go. Interface 2 (workflow content model) remains untyped — design decision, not oversight. Improved from 4/7 to 6/7. |
| Attack 2 (no error types/codes) | Yes | TaskStatus enum shown from types.go. ErrorCode type shown from errors.go. Error cases table maps each case to a status and error code. Score improved from 2/5 to 4/5. |
| Attack 3 (missing coverage target) | Yes | Explicit "maintain >= 80% coverage on changed files" with `go test -race -cover ./...` verification. Score improved from 0/5 to 5/5. |

---

## Verdict

- **Score**: 87/100
- **Target**: 90/100
- **Gap**: 3 points
- **Breakdown-Readiness**: 18/20 — can proceed to /breakdown-tasks
- **Action**: Breakdown-Readiness passes gate (>=18). Gap is in Interface & Model Definitions (-3), Error Handling (-2), and Testing (-2). Consider one more iteration to close the agent→task-cli error boundary gap and add Case C test scenario, or proceed with breakdown accepting minor gaps.
