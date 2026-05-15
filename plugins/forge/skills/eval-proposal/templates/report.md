---
date: "{{DATE}}"
doc_dir: "{{DOC_DIR}}"
iteration: "{{ITERATION}}"
target_score: "{{TARGET}}"
evaluator: Claude (automated, adversarial)
---

# Proposal Eval — Iteration {{ITERATION}}

**Score: {{SCORE}}/1000** (target: {{TARGET}})

```
┌──────────────────────────────────────────────────────────────────────────┐
│                     PROPOSAL QUALITY SCORECARD (1000 pts)                │
├─────────────────────────────────────┬──────────┬──────────┬─────────────┤
│ Dimension                           │ Score    │ Max      │ Status      │
├─────────────────────────────────────┼──────────┼──────────┼─────────────┤
│ 1. Problem Definition               │  ___     │  110     │ ✅/⚠️/❌     │
│    Problem clarity                  │  ___/40  │          │             │
│    Evidence provided                │  ___/40  │          │             │
│    Urgency justified                │  ___/30  │          │             │
├─────────────────────────────────────┼──────────┼──────────┼─────────────┤
│ 2. Solution Clarity                 │  ___     │  120     │ ✅/⚠️/❌     │
│    Approach concrete                │  ___/40  │          │             │
│    User-facing behavior             │  ___/45  │          │             │
│    Technical direction              │  ___/35  │          │             │
├─────────────────────────────────────┼──────────┼──────────┼─────────────┤
│ 3. Industry Benchmarking            │  ___     │  120     │ ✅/⚠️/❌     │
│    Industry solutions referenced    │  ___/40  │          │             │
│    3+ meaningful alternatives       │  ___/30  │          │             │
│    Honest trade-off comparison      │  ___/25  │          │             │
│    Justified against benchmarks     │  ___/25  │          │             │
├─────────────────────────────────────┼──────────┼──────────┼─────────────┤
│ 4. Requirements Completeness        │  ___     │  110     │ ✅/⚠️/❌     │
│    Scenario coverage                │  ___/40  │          │             │
│    Non-functional requirements      │  ___/40  │          │             │
│    Constraints & dependencies       │  ___/30  │          │             │
├─────────────────────────────────────┼──────────┼──────────┼─────────────┤
│ 5. Solution Creativity              │  ___     │  100     │ ✅/⚠️/❌     │
│    Novelty over industry baseline   │  ___/40  │          │             │
│    Cross-domain inspiration         │  ___/35  │          │             │
│    Simplicity of insight            │  ___/25  │          │             │
├─────────────────────────────────────┼──────────┼──────────┼─────────────┤
│ 6. Feasibility                      │  ___     │  100     │ ✅/⚠️/❌     │
│    Technical feasibility            │  ___/40  │          │             │
│    Resource & timeline feasibility  │  ___/30  │          │             │
│    Dependency readiness             │  ___/30  │          │             │
├─────────────────────────────────────┼──────────┼──────────┼─────────────┤
│ 7. Scope Definition                 │  ___     │  80      │ ✅/⚠️/❌     │
│    In-scope concrete                │  ___/30  │          │             │
│    Out-of-scope explicit            │  ___/25  │          │             │
│    Scope bounded                    │  ___/25  │          │             │
├─────────────────────────────────────┼──────────┼──────────┼─────────────┤
│ 8. Risk Assessment                  │  ___     │  90      │ ✅/⚠️/❌     │
│    Risks identified (≥3)            │  ___/30  │          │             │
│    Likelihood + impact rated        │  ___/30  │          │             │
│    Mitigations actionable           │  ___/30  │          │             │
├─────────────────────────────────────┼──────────┼──────────┼─────────────┤
│ 9. Success Criteria                 │  ___     │  80      │ ✅/⚠️/❌     │
│    Measurable and testable          │  ___/55  │          │             │
│    Coverage complete                │  ___/25  │          │             │
├─────────────────────────────────────┼──────────┼──────────┼─────────────┤
│ 10. Logical Consistency             │  ___     │  90      │ ✅/⚠️/❌     │
│     Solution ↔ Problem              │  ___/35  │          │             │
│     Scope ↔ Solution ↔ Criteria     │  ___/30  │          │             │
│     Requirements ↔ Solution         │  ___/25  │          │             │
├─────────────────────────────────────┼──────────┼──────────┼─────────────┤
│ TOTAL                               │  ___     │  1000    │             │
└─────────────────────────────────────┴──────────┴──────────┴─────────────┘
```

---

## Deductions

| Location | Issue | Penalty |
|----------|-------|---------|
| <!-- section:line --> | <!-- vague language / missing content / inconsistency --> | <!-- -N pts --> |

---

## Attack Points

### Attack 1: [dimension — specific weakness]

**Where**: <!-- quote from proposal -->
**Why it's weak**: <!-- concrete critique -->
**What must improve**: <!-- actionable fix -->

### Attack 2: [dimension — specific weakness]

**Where**: <!-- quote from proposal -->
**Why it's weak**: <!-- concrete critique -->
**What must improve**: <!-- actionable fix -->

### Attack 3: [dimension — specific weakness]

**Where**: <!-- quote from proposal -->
**Why it's weak**: <!-- concrete critique -->
**What must improve**: <!-- actionable fix -->

---

## Previous Issues Check

<!-- Only for iteration > 1 -->

| Previous Attack | Addressed? | Evidence |
|----------------|------------|----------|
| <!-- attack from iter N-1 --> | ✅/❌ | <!-- what changed / what didn't --> |

---

## Verdict

- **Score**: {{SCORE}}/1000
- **Target**: {{TARGET}}/1000
- **Gap**: {{GAP}} points
- **Action**: {{Continue to iteration N+1 / Target reached / Iterations exhausted}}
