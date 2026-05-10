---
date: "2026-05-10"
doc_dir: "docs/features/task-executor-skeleton/design/"
iteration: 3
target_score: 90
evaluator: Claude (automated, adversarial)
---

# Design Eval — Iteration 3

**Score: 93/100** (target: 90)

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
│ 2. Interface & Model Defs    │  19      │  20      │ OK         │
│    Interface signatures typed│  7/7     │          │            │
│    Models concrete           │  6/7     │          │            │
│    Directly implementable    │  6/6     │          │            │
├──────────────────────────────┼──────────┼──────────┬────────────┤
│ 3. Error Handling            │  15      │  15      │ OK         │
│    Error types defined       │  5/5     │          │            │
│    Propagation strategy clear│  5/5     │          │            │
│    HTTP status codes mapped  │  5/5     │          │ N/A        │
├──────────────────────────────┼──────────┼──────────┬────────────┤
│ 4. Testing Strategy          │  15      │  15      │ OK         │
│    Per-layer test plan       │  5/5     │          │            │
│    Coverage target numeric   │  5/5     │          │            │
│    Test tooling named        │  5/5     │          │            │
├──────────────────────────────┼──────────┼──────────┬────────────┤
│ 5. Breakdown-Readiness ★     │  19      │  20      │ OK         │
│    Components enumerable     │  7/7     │          │            │
│    Tasks derivable           │  7/7     │          │            │
│    PRD AC coverage           │  5/6     │          │            │
├──────────────────────────────┼──────────┼──────────┬────────────┤
│ 6. Security Considerations   │  10      │  10      │ N/A        │
│    Threat model present      │  N/A     │          │            │
│    Mitigations concrete      │  N/A     │          │            │
├──────────────────────────────┼──────────┼──────────┬────────────┤
│ TOTAL                        │  93      │  100     │            │
└──────────────────────────────┴──────────┴──────────┴────────────┘
```

★ Breakdown-Readiness >= 18/20 — can proceed to /breakdown-tasks

---

## Deductions

| Location | Issue | Penalty |
|----------|-------|---------|
| Interface 2 (Workflow Content Model) | W1-W5 rules are structured as a validation checklist with examples/counter-examples. This is a substantial improvement from iteration 2's prose-only rules, but it remains a human-readable checklist rather than a machine-parseable schema. A developer implementing automated validation must infer parsing rules from the examples. The design acknowledges "a future linter could parse numbered steps" but provides no grammar, regex, or AST for the workflow content structure. | -1 pts (Dim 2) |
| PRD Coverage Map + Error Cases row 6 | Cross-section inconsistency on fix-task status. PRD user story 4 AC says "agent 按指令执行后停止，任务状态设为 `failed`" but Error Cases row 6 says "Workflow declares failure instructions" → "Status Set: `blocked`". PRD says `failed`, design says `blocked`. These are semantically different terminal states in the TaskStatus enum. | -1 pts (Dim 5) |

---

## Attack Points

### Attack 1: Interface & Model Definitions — workflow content validation is process-bound, not code-bound

**Where**: Interface 2, "Workflow Content Model" table (W1-W5) and "Automated validation" paragraph: `grep -c "^## Execution Workflow" templates/*.md` confirms every template has the heading. Content compliance (W1-W5) is enforced in template review with the checklist above — a future linter could parse numbered steps and detect missing Stop / Proceed to Step markers.

**Why it's weak**: The `grep -c` check only validates heading presence, not content compliance. The W1-W5 rules are well-structured with examples and counter-examples (a significant improvement), but "enforced in template review" is still a process answer. The design defers the linter to "future" without specifying what the linter would parse. For a design doc, this is acceptable but not excellent -- the core innovation (Execution Workflow) has no formal content model that a developer can validate in CI.

**What must improve**: Either (a) add a minimal regex or structured format that W1-W5 can be checked against (e.g., "each step must start with `\d+\.` and contain at least one backtick-wrapped command"), or (b) accept this as a known limitation and document it in Open Questions instead of leaving it implicit.

### Attack 2: Breakdown-Readiness — PRD/design status mismatch for fix-task creation

**Where**: Error Cases table row 6: "Workflow declares failure instructions | Step 2 (workflow-specific) | `blocked`". PRD user story 4, second AC: "agent 按指令执行后停止，任务状态设为 `failed`".

**Why it's weak**: The PRD explicitly states the task status should be `failed` when the agent follows failure instructions (e.g., creating a fix task). The design sets it to `blocked`. These are distinct states in the TaskStatus enum: `failed` implies the task did not succeed, `blocked` implies the task is waiting on an external resolution. While `blocked` is arguably more precise (a fix task must resolve the issue), the design does not acknowledge or justify this deviation from the PRD. A developer implementing from the design would follow the Error Cases table (set `blocked`), while someone testing against PRD ACs would expect `failed`. This is a cross-section inconsistency per the rubric's -3 deduction rule (reduced to -1 since only one AC is affected and the design choice is defensible).

**What must improve**: Either (a) align the design with the PRD by changing Error Cases row 6 status from `blocked` to `failed`, or (b) update the PRD AC to say `blocked` and add a note in the design explaining why `blocked` is the correct status (the task is not failed -- it created a fix task and is waiting for resolution). The PRD Coverage Map should reflect whichever choice is made.

### Attack 3: Architecture Clarity — fallback path dependency on template file location is implicit

**Where**: Interface 1 Case B: "Read the default workflow template at: `plugins/forge/skills/breakdown-tasks/templates/task.md`". Dependencies section: "The fallback reads `plugins/forge/skills/breakdown-tasks/templates/task.md` — an existing file."

**Why it's weak**: This is a minor issue. The hardcoded path to the default template creates a tight coupling: if the template directory is reorganized, the agent prompt breaks silently (the agent would fall through to Case C behavior -- empty/not-found -- without a clear error). The design does not discuss this coupling risk. This is not a scoring deduction (the architecture dimension is already at max), but it is a design fragility worth noting for the implementer.

**What must improve**: Consider adding the default template path to the Dependencies section as an explicit runtime dependency with a note that the path must be kept in sync between the agent prompt and the actual file location. Alternatively, consider making the path configurable or discoverable at runtime rather than hardcoded in the prompt.

---

## Previous Issues Check

| Previous Attack | Addressed? | Evidence |
|----------------|------------|----------|
| Attack 1 (workflow content has no typed model) | Yes | Interface 2 now has a structured Workflow Content Model table (W1-W5) with rule, example, and counter-example columns. Automated `grep -c` validation specified. Significant improvement from prose-only rules. Remaining gap (-1) is the absence of a machine-parseable schema for content validation. |
| Attack 2 (agent-to-task-cli error boundary is a black box) | Yes | New "Agent-to-task-cli Error Translation" section (3 numbered steps) specifies exact mechanism: agent writes status in frontmatter, calls `task record --status failed`, task-cli persists via CLI arguments. Freeform text is explicitly described as human-readable only. Fully resolves the iteration-2 gap. |
| Attack 3 (agent prompt test coverage incomplete) | Yes | Test scenario 8 added: "Task file with `## Execution Workflow` heading but empty content" covers Case C. Agent prompt row in Per-Layer Test Plan explicitly references "test scenarios 4-8". All three Cases (A, B, C) now have explicit test coverage. |

---

## Verdict

- **Score**: 93/100
- **Target**: 90/100
- **Gap**: Target reached (+3 above target)
- **Breakdown-Readiness**: 19/20 — can proceed to /breakdown-tasks
- **Action**: Target reached. Proceed to /breakdown-tasks. Minor recommendations: (1) resolve the `failed` vs `blocked` status inconsistency with PRD before implementation, (2) consider adding a note about the hardcoded default template path dependency.
