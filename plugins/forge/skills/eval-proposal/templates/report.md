---
date: "{{DATE}}"
doc_dir: "{{DOC_DIR}}"
iteration: "{{ITERATION}}"
target_score: "{{TARGET}}"
evaluator: Claude (automated, adversarial)
---

# Proposal Eval — Iteration {{ITERATION}}

**Score: {{SCORE}}/100** (target: {{TARGET}})

```
┌─────────────────────────────────────────────────────────────────┐
│                    PROPOSAL QUALITY SCORECARD                     │
├──────────────────────────────┬──────────┬──────────┬────────────┤
│ Dimension                    │ Score    │ Max      │ Status     │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 1. Problem Definition        │  ___     │  20      │ ✅/⚠️/❌    │
│    Problem clarity           │  ___/7   │          │            │
│    Evidence provided         │  ___/7   │          │            │
│    Urgency justified         │  ___/6   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 2. Solution Clarity          │  ___     │  20      │ ✅/⚠️/❌    │
│    Approach concrete         │  ___/7   │          │            │
│    User-facing behavior      │  ___/7   │          │            │
│    Differentiated            │  ___/6   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 3. Alternatives Analysis     │  ___     │  15      │ ✅/⚠️/❌    │
│    Alternatives listed (≥2)  │  ___/5   │          │            │
│    Pros/cons honest          │  ___/5   │          │            │
│    Rationale justified       │  ___/5   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 4. Scope Definition          │  ___     │  15      │ ✅/⚠️/❌    │
│    In-scope concrete         │  ___/5   │          │            │
│    Out-of-scope explicit     │  ___/5   │          │            │
│    Scope bounded             │  ___/5   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 5. Risk Assessment           │  ___     │  15      │ ✅/⚠️/❌    │
│    Risks identified (≥3)     │  ___/5   │          │            │
│    Likelihood + impact rated │  ___/5   │          │            │
│    Mitigations actionable    │  ___/5   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 6. Success Criteria          │  ___     │  15      │ ✅/⚠️/❌    │
│    Measurable                │  ___/5   │          │            │
│    Coverage complete         │  ___/5   │          │            │
│    Testable                  │  ___/5   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ TOTAL                        │  ___     │  100     │            │
└──────────────────────────────┴──────────┴──────────┴────────────┘
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

- **Score**: {{SCORE}}/100
- **Target**: {{TARGET}}/100
- **Gap**: {{GAP}} points
- **Action**: {{Continue to iteration N+1 / Target reached / Iterations exhausted}}
