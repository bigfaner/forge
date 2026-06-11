---
date: "2026-05-14"
doc_dir: "docs/proposals/task-stage-gates/"
iteration: 2
target_score: 800
evaluator: Claude (automated, adversarial)
---

# Proposal Eval — Iteration 2

**Score: 802/1000** (target: 800)

```
┌──────────────────────────────────────────────────────────────────────────┐
│                     PROPOSAL QUALITY SCORECARD (1000 pts)                │
├─────────────────────────────────────────────┬──────────┬────────┬────────┤
│ Dimension                                   │ Score    │ Max    │ Status │
├─────────────────────────────────────────────┼──────────┼────────┼────────┤
│ 1. Problem Definition                       │   86     │  110   │ ⚠️     │
│    Problem clarity                          │   36/40  │        │        │
│    Evidence provided                        │   28/40  │        │        │
│    Urgency justified                        │   22/30  │        │        │
├─────────────────────────────────────────────┼──────────┼────────┼────────┤
│ 2. Solution Clarity                         │  100     │  120   │ ✅     │
│    Approach concrete                        │   37/40  │        │        │
│    User-facing behavior                     │   30/45  │        │        │
│    Technical direction                      │   33/35  │        │        │
├─────────────────────────────────────────────┼──────────┼────────┼────────┤
│ 3. Industry Benchmarking                    │   93     │  120   │ ⚠️     │
│    Industry solutions referenced            │   32/40  │        │        │
│    3+ meaningful alternatives               │   24/30  │        │        │
│    Honest trade-off comparison              │   19/25  │        │        │
│    Justified against benchmarks             │   18/25  │        │        │
├─────────────────────────────────────────────┼──────────┼────────┼────────┤
│ 4. Requirements Completeness                │   95     │  110   │ ✅     │
│    Scenario coverage                        │   35/40  │        │        │
│    Non-functional requirements              │   34/40  │        │        │
│    Constraints & dependencies               │   26/30  │        │        │
├─────────────────────────────────────────────┼──────────┼────────┼────────┤
│ 5. Solution Creativity                      │   55     │  100   │ ⚠️     │
│    Novelty over industry baseline           │   18/40  │        │        │
│    Cross-domain inspiration                 │   15/35  │        │        │
│    Simplicity of insight                    │   22/25  │        │        │
├─────────────────────────────────────────────┼──────────┼────────┼────────┤
│ 6. Feasibility                              │   92     │  100   │ ✅     │
│    Technical feasibility                    │   37/40  │        │        │
│    Resource & timeline feasibility          │   28/30  │        │        │
│    Dependency readiness                     │   27/30  │        │        │
├─────────────────────────────────────────────┼──────────┼────────┼────────┤
│ 7. Scope Definition                         │   74     │   80   │ ✅     │
│    In-scope concrete                        │   27/30  │        │        │
│    Out-of-scope explicit                    │   23/25  │        │        │
│    Scope bounded                            │   24/25  │        │        │
├─────────────────────────────────────────────┼──────────┼────────┼────────┤
│ 8. Risk Assessment                          │   70     │   90   │ ⚠️     │
│    Risks identified (≥3)                    │   26/30  │        │        │
│    Likelihood + impact rated                │   22/30  │        │        │
│    Mitigations actionable                   │   22/30  │        │        │
├─────────────────────────────────────────────┼──────────┼────────┼────────┤
│ 9. Success Criteria                         │   63     │   80   │ ⚠️     │
│    Measurable and testable                  │   42/55  │        │        │
│    Coverage complete                        │   21/25  │        │        │
├─────────────────────────────────────────────┼──────────┼────────┼────────┤
│ 10. Logical Consistency                     │   74     │   90   │ ⚠️     │
│     Solution ↔ Problem                      │   28/35  │        │        │
│     Scope ↔ Solution ↔ Criteria             │   24/30  │        │        │
│     Requirements ↔ Solution                 │   22/25  │        │        │
├─────────────────────────────────────────────┴──────────┴────────┴────────┤
│ TOTAL                                        │  802     │  1000  │        │
└──────────────────────────────────────────────────────────────────────────┘
```

---

## Deductions

| Location | Issue | Penalty |
|----------|-------|---------|
| Problem:Evidence | Evidence is entirely internal code/file references — no user feedback, bug reports, incident data, or quality metrics demonstrating actual failures caused by missing stage-gates | -12 pts |
| Problem:Urgency | "As more features use quick mode, the gap widens" — no data on quick-mode adoption rate, number of affected features, or quantified quality cost | -8 pts |
| Solution:User-facing | No description of CLI output — what does the developer see after `forge task index` runs? Is there a log line confirming "Generated 2 stage-gate files"? What if 0 gates qualify? Silent generation is an observable gap | -15 pts |
| Benchmarking:Industry | Each industry tool gets 1-2 sentences of analysis — GitHub Actions environment rules, Bazel test_suite, and GitLab CI rules are named but not deeply analyzed for how their gate patterns map to this context | -8 pts |
| Benchmarking:Alternatives | "Standalone `forge task stage-gates`" is still a borderline straw-man — described primarily by its failure mode ("does not reduce agent-dependent failure") rather than its merits | -6 pts |
| Benchmarking:Trade-offs | No quantified comparison — latency impact, maintenance burden over time, lines-of-code delta between alternatives are all absent | -6 pts |
| Benchmarking:Justification | "piggybacks on already-mandatory command" is a clear rationale, but no explicit mapping of "which industry pattern does each alternative follow?" | -7 pts |
| Requirements:Scenarios | No scenario for cross-feature phase number collision (e.g., feature A has phase 1, feature B has phase 1 — how are they distinguished?) | -5 pts |
| Requirements:NFR | Performance "< 5ms" is an estimate, not measured. No NFR for template rendering failure (what if `text/template.Execute` returns an error?) | -6 pts |
| Creativity:Novelty | "This is a straightforward adoption of the existing test-task generation pattern" — self-admitted as a copy. The "where to put it" insight is valuable but the implementation is mechanically identical to existing code | -22 pts |
| Creativity:Cross-domain | No explicit cross-domain inspiration cited. Industry benchmarking provides some cross-tool awareness but the proposal does not name a specific insight borrowed from another domain | -20 pts |
| Risk:Impact ratings | All three risks rated L/M or M/L — suspiciously uniform. "Template content diverges from skill-authored gate/summary" could have High impact if it causes incorrect gate dependencies in production | -8 pts |
| Risk:Mitigations | "Use the same templates from `breakdown-tasks/templates/` via `go:embed`" is a design decision, not a divergence mitigation. A real mitigation would specify a CI check or sync mechanism that detects drift | -8 pts |
| Risk:Missing risk | No risk for "existing projects with hand-crafted gate files using different schemas" — this scenario is listed in Requirements but absent from Risk Assessment | -4 pts |
| Criteria:Testable | "Unit tests cover: phase detection, test-task exclusion, template generation, idempotency, malformed task IDs" — "cover" is vague; no specification of what assertions or edge-case tests are required | -13 pts |
| Criteria:Coverage | "Version bump" appears in Scope (In Scope) but no success criterion validates the version number change | -4 pts |
| Consistency:Scope↔Criteria | "Version bump" in scope, no criterion — same inconsistency noted in iteration 1, still present | -6 pts |

---

## Attack Points

### Attack 1: Solution Creativity — self-admitted pattern copy with no cross-domain insight

**Where**: "This is a straightforward adoption of the existing test-task generation pattern (`generateTestTasks` in `build.go`) to a new task type."

**Why it's weak**: The proposal explicitly states it is copying an existing internal pattern. The only creative contribution is the insight that gate generation belongs in the deterministic CLI layer — a positioning decision, not a technical innovation. No cross-domain inspiration is cited (no borrowing from compiler design, database migration systems, or any domain outside the project). For a dimension worth 100 points, a "straightforward adoption" that scores 18/40 on novelty and 15/35 on cross-domain inspiration is the weakest area of the proposal.

**What must improve**: Identify at least one cross-domain analogy (e.g., how database migration tools like Flyway use "baseline" detection to skip already-applied migrations — structurally identical to the skip-if-exists logic). Articulate what is novel beyond "copy existing pattern to new task type" — even if the novelty is framing it as a general principle for deterministic task generation in AI-assisted workflows.

### Attack 2: Solution Clarity — user-facing behavior remains undescribed

**Where**: "The AI agent still generates business task files, then calls `forge task index --feature <slug>`. One command handles everything: stage-gate generation + index building."

**Why it's weak**: This describes the agent's workflow, not the developer's observable experience. After `forge task index` runs, what does the developer see? Is there a log line like "Generated stage-gates for phases 1, 2"? If zero gates qualify (e.g., all phases have <2 business tasks), is there any output? Silent success is indistinguishable from silent failure. The iteration-1 report specifically called this out: "No description of what the AI agent or developer sees/experiences after running `forge task index`." The proposal added the "Workflow Guarantee" section addressing the logical gap, but the observable-behavior gap persists.

**What must improve**: Add a concrete description of CLI output behavior: what is printed to stdout/stderr when gates are generated, when zero gates qualify, and when an error occurs during template rendering.

### Attack 3: Risk Assessment — mitigations are design choices, not risk-reduction actions

**Where**: Risk table row 2 — "Template content diverges from skill-authored gate/summary" with mitigation "Use the same templates from `breakdown-tasks/templates/` via `go:embed`".

**Why it's weak**: Using the same template files via `go:embed` is how the feature works — it is the design, not a mitigation for the risk of those templates diverging from what the skill expects over time. If someone updates the skill template without updating the embedded template (or vice versa), the divergence still occurs. A real mitigation would be: "CI check that diffs `breakdown-tasks/templates/` against `pkg/task/templates/` on every PR" or "Single source of truth: skill reads templates from the Go binary at runtime." The same issue applies to risk 3's mitigation: "`forge task validate-index` already checks for cycles" — this tool already exists and is not a new mitigation action.

**What must improve**: Replace design-description mitigations with concrete risk-reduction actions. For template divergence: specify a sync mechanism (CI check, single source of truth, or documented review step). For cycle detection: specify that the existing tool will be added to CI or documented as a required validation step.

---

## Previous Issues Check

| Previous Attack | Addressed? | Evidence |
|----------------|------------|----------|
| Attack 1 (iter 1): Industry Benchmarking — shallow and straw-manned | Partially | Industry tools now named (GitHub Actions, Bazel, GitLab CI) with specific features. "Do nothing" removed. "Standalone command" remains as borderline straw-man. Each tool gets only 1-2 sentences. |
| Attack 2 (iter 1): Logical Consistency — solution does not fully solve Problem #2 | Yes | New "Workflow Guarantee" section explicitly acknowledges the residual gap and provides three concrete reasons for acceptance. Honest and well-argued. |
| Attack 3 (iter 1): Requirements Completeness — missing edge cases and NFR gaps | Yes | Added scenarios for partial state, concurrent execution, malformed task IDs, and pre-existing hand-crafted files. NFRs now quantified ("< 5ms", security analysis, Go version, backward compatibility). |
| Deduction: Risk impact ratings uniform | No | Still L/M, M/L, L/M across all three risks. No risk rated High impact. |
| Deduction: Version bump scope/criteria mismatch | No | "Version bump" still in scope, still no success criterion for it. |

---

## Verdict

- **Score**: 802/1000
- **Target**: 800/1000
- **Gap**: -2 points (target reached)
- **Action**: Target reached. Primary remaining weaknesses are Solution Creativity (55/100), Risk Assessment (70/90), and Success Criteria (63/80). These are structural — creativity is inherently limited by the "straightforward adoption" nature of the feature, and further iterations are unlikely to significantly improve scores without fundamentally changing the proposal.
