---
date: "2026-05-14"
doc_dir: "docs/proposals/task-stage-gates/"
iteration: 1
target_score: 800
evaluator: Claude (automated, adversarial)
---

# Proposal Eval — Iteration 1

**Score: 655/1000** (target: 800)

```
┌──────────────────────────────────────────────────────────────────────────┐
│                     PROPOSAL QUALITY SCORECARD (1000 pts)                │
├─────────────────────────────────────────────┬──────────┬────────┬────────┤
│ Dimension                                   │ Score    │ Max    │ Status │
├─────────────────────────────────────────────┼──────────┼────────┼────────┤
│ 1. Problem Definition                       │   80     │  110   │ ⚠️     │
│    Problem clarity                          │   32/40  │        │        │
│    Evidence provided                        │   28/40  │        │        │
│    Urgency justified                        │   20/30  │        │        │
├─────────────────────────────────────────────┼──────────┼────────┼────────┤
│ 2. Solution Clarity                         │   95     │  120   │ ⚠️     │
│    Approach concrete                        │   38/40  │        │        │
│    User-facing behavior                     │   30/45  │        │        │
│    Technical direction                      │   27/35  │        │        │
├─────────────────────────────────────────────┼──────────┼────────┼────────┤
│ 3. Industry Benchmarking                    │   40     │  120   │ ❌     │
│    Industry solutions referenced            │   12/40  │        │        │
│    3+ meaningful alternatives               │   10/30  │        │        │
│    Honest trade-off comparison              │   10/25  │        │        │
│    Justified against benchmarks             │    8/25  │        │        │
├─────────────────────────────────────────────┼──────────┼────────┼────────┤
│ 4. Requirements Completeness                │   70     │  110   │ ⚠️     │
│    Scenario coverage                        │   25/40  │        │        │
│    Non-functional requirements              │   20/40  │        │        │
│    Constraints & dependencies               │   25/30  │        │        │
├─────────────────────────────────────────────┼──────────┼────────┼────────┤
│ 5. Solution Creativity                      │   50     │  100   │ ⚠️     │
│    Novelty over industry baseline           │   15/40  │        │        │
│    Cross-domain inspiration                 │   15/35  │        │        │
│    Simplicity of insight                    │   20/25  │        │        │
├─────────────────────────────────────────────┼──────────┼────────┼────────┤
│ 6. Feasibility                              │   85     │  100   │ ✅     │
│    Technical feasibility                    │   35/40  │        │        │
│    Resource & timeline feasibility          │   25/30  │        │        │
│    Dependency readiness                     │   25/30  │        │        │
├─────────────────────────────────────────────┼──────────┼────────┼────────┤
│ 7. Scope Definition                         │   70     │   80   │ ✅     │
│    In-scope concrete                        │   25/30  │        │        │
│    Out-of-scope explicit                    │   22/25  │        │        │
│    Scope bounded                            │   23/25  │        │        │
├─────────────────────────────────────────────┼──────────┼────────┼────────┤
│ 8. Risk Assessment                          │   65     │   90   │ ⚠️     │
│    Risks identified (≥3)                    │   25/30  │        │        │
│    Likelihood + impact rated                │   20/30  │        │        │
│    Mitigations actionable                   │   20/30  │        │        │
├─────────────────────────────────────────────┼──────────┼────────┼────────┤
│ 9. Success Criteria                         │   60     │   80   │ ⚠️     │
│    Measurable and testable                  │   40/55  │        │        │
│    Coverage complete                        │   20/25  │        │        │
├─────────────────────────────────────────────┼──────────┼────────┼────────┤
│ 10. Logical Consistency                     │   40     │   90   │ ❌     │
│     Solution ↔ Problem                      │   15/35  │        │        │
│     Scope ↔ Solution ↔ Criteria             │   15/30  │        │        │
│     Requirements ↔ Solution                 │   10/25  │        │        │
├─────────────────────────────────────────────┴──────────┴────────┴────────┤
│ TOTAL                                        │  655     │  1000  │        │
└──────────────────────────────────────────────────────────────────────────┘
```

---

## Deductions

| Location | Issue | Penalty |
|----------|-------|---------|
| Problem:Urgency | "As more features use quick mode, the gap widens" — no data on how many features use quick mode, no quantified impact | -10 pts |
| Problem:Evidence | Evidence is all internal code references — no user feedback, bug reports, or incident data demonstrating actual quality failures caused by missing gates | -12 pts |
| Solution:User-facing | No description of what the AI agent or developer sees/experiences after running `forge task index` — what output confirms gates were generated? What if 0 gates are generated? | -15 pts |
| Benchmarking:Industry | "Stage-Gate(R) is a standard product development methodology (Cooper, 1990)" — a one-sentence name-drop with no analysis of how Stage-Gate principles apply to this specific CLI context | -28 pts |
| Benchmarking:Alternatives | "Do nothing" and "standalone command" are both straw-man alternatives — neither is a genuinely different architectural approach. Only 3 alternatives, and 2 are trivially dismissible | -20 pts (straw-man: -20 each for "do nothing" and "standalone command") |
| Benchmarking:Trade-offs | Comparison table has 1-2 word cons ("Extra step", "Slightly increases complexity") — not honest trade-offs based on project constraints | -15 pts |
| Benchmarking:Justification | "minimal API surface, maximum reliability" is a slogan, not a rationale grounded in the benchmarked alternatives | -17 pts |
| Requirements:Scenarios | Missing: What happens when a task ID has no phase prefix? What about unnumbered tasks? What if `.summary` exists but `.gate` does not (partial state)? What about concurrent runs? | -15 pts |
| Requirements:NFR | "Generation time: negligible" is vague. No security consideration (template injection?). No compatibility analysis (Go version, OS). No performance baseline measurement | -20 pts |
| Creativity:Novelty | "straightforward adoption of the existing test-task generation pattern" — self-admitted as a copy of an existing pattern, not novel | -25 pts |
| Creativity:Cross-domain | No cross-domain inspiration cited. The proposal explicitly says it copies an existing internal pattern | -20 pts |
| Risk:Impact ratings | All three risks are rated L/M or M/L — suspiciously uniform. No risk rated "High" impact despite the feature touching core indexing logic | -10 pts |
| Risk:Mitigations | "Use the same templates from `breakdown-tasks/templates/` via `go:embed`" — this is a design choice, not a mitigation for template divergence. A mitigation would specify a sync/review mechanism | -10 pts |
| Consistency:Solution↔Problem | Problem #2 is "No deterministic guarantee" but the solution still relies on the AI agent calling `forge task index` — if the agent never calls the command, the gap persists. The solution only makes generation deterministic *if* invoked | -20 pts |
| Consistency:Scope↔Criteria | In-scope includes "Version bump" but no success criterion checks for version number correctness | -10 pts |
| Consistency:Req↔Solution | Requirement "Phase with only test tasks" is listed but the solution's phase detection algorithm doesn't specify how to exclude test-only phases from the threshold count | -15 pts |

---

## Attack Points

### Attack 1: Industry Benchmarking — shallow and straw-manned

**Where**: "Stage-Gate(R) is a standard product development methodology (Cooper, 1990). Each phase ends with a gate that evaluates quality before proceeding. Most CI/CD systems implement similar patterns (pipeline stages with gate approvals)."

**Why it's weak**: This is a two-sentence gloss of an entire domain. No specific CI/CD system is named. No open-source project implementing stage-gate patterns is cited. No analysis of how Cooper's Stage-Gate maps to task dependency graphs. The three alternatives in the comparison table are "do nothing", "standalone command", and "integrate into index" — two of these are straw-men trivially dismissed. There is zero engagement with how any real tool (GitHub Actions, GitLab CI, Jenkins, Tekton) implements gate patterns.

**What must improve**: Cite at least 3 concrete industry tools or projects (e.g., GitHub Actions `environment` gates, Bazel's `query` for dependency verification). Replace "do nothing" and "standalone command" with genuinely different approaches (e.g., hook-based generation via `pre-commit`, template scaffolding via `cookiecutter`-like CLI). Provide a real trade-off analysis with quantified complexity/latency/maintenance costs.

### Attack 2: Logical Consistency — solution does not fully solve Problem #2

**Where**: Problem states "No deterministic guarantee: Gate/summary generation depends on the AI agent correctly following skill instructions" and Solution says "The AI agent still generates business task files, then calls `forge task index --feature <slug>`. One command handles everything."

**Why it's weak**: The solution still requires the AI agent to call `forge task index`. If the agent forgets this call (the exact failure mode described in Problem #2), no gates are generated. The determinism is conditional on AI compliance, which is the very problem stated. The solution shifts the failure point but does not eliminate it. A truly deterministic solution would hook into `forge task index` as a mandatory step in the task creation pipeline, or generate gates as a post-processing step in `forge task add`.

**What must improve**: Either (a) acknowledge this residual gap explicitly and explain why the reduced failure surface is acceptable, or (b) propose a mechanism where gate generation is truly automatic (e.g., triggered on every task file creation via filesystem watch, or embedded in `forge task add`).

### Attack 3: Requirements Completeness — missing edge cases and NFR gaps

**Where**: Requirements Analysis section lists 5 scenarios but misses critical edge cases: "Phase with only test tasks" is listed as a scenario but the algorithm in "How It Works" does not specify the exclusion logic for test-only phases.

**Why it's weak**: Multiple edge cases are unaddressed: What if a `.summary` file exists but `.gate` does not (partial generation)? What if task IDs use non-standard numbering (e.g., `1.10`, `1.2a`)? What if `forge task index` is run concurrently by two agent sessions? The NFR section lists "negligible" generation time (vague, unquantified) and ignores security (template injection if task IDs contain malicious content), compatibility (Go version requirements), and backward compatibility (what if existing projects have hand-crafted `.summary` files with different schemas?).

**What must improve**: Add edge cases for partial state (`.summary` exists without `.gate`), concurrent execution, malformed task IDs, and existing hand-crafted gate files. Replace "negligible" with a measured or estimated bound. Add NFRs for security (input validation on task IDs) and backward compatibility (schema versioning or migration).

---

## Previous Issues Check

<!-- Only for iteration > 1 — not applicable for iteration 1 -->

---

## Verdict

- **Score**: 655/1000
- **Target**: 800/1000
- **Gap**: 145 points
- **Action**: Continue to iteration 2 — primary gaps are Industry Benchmarking (40/120), Logical Consistency (40/90), Solution Creativity (50/100), and Requirements Completeness (70/110). Strengthening benchmarking with real tool citations and fixing the logical gap where the solution does not fully address Problem #2 would recover the most points.
