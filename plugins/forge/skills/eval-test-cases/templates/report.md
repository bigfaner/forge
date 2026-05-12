---
date: "{{DATE}}"
doc_dir: "{{DOC_DIR}}"
iteration: "{{ITERATION}}"
target_score: "{{TARGET}}"
evaluator: Claude (automated, adversarial)
---

# Test Cases Eval — Iteration {{ITERATION}}

**Score: {{SCORE}}/1000** (target: {{TARGET}})

```
┌─────────────────────────────────────────────────────────────────┐
│                  TEST CASES QUALITY SCORECARD                     │
├──────────────────────────────────────────────────────────────────┤
│ Dimension                    │ Score    │ Max      │ Status     │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 1. PRD Traceability          │  ___     │  250     │ ✅/⚠️/❌    │
│    TC-to-AC mapping          │  ___/90  │          │            │
│    Traceability table        │  ___/80  │          │            │
│    Reverse coverage          │  ___/80  │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 2. Step Actionability        │  ___     │  250     │ ✅/⚠️/❌    │
│    Steps concrete            │  ___/90  │          │            │
│    Expected results          │  ___/90  │          │            │
│    Preconditions explicit    │  ___/70  │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 3. Route & Element Accuracy  │  ___     │  200     │ ✅/⚠️/❌    │
│    Routes valid              │  ___/70  │          │            │
│    Elements identifiable     │  ___/70  │          │            │
│    Consistency               │  ___/60  │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 4. Completeness              │  ___     │  200     │ ✅/⚠️/❌    │
│    Type coverage             │  ___/70  │          │            │
│    Boundary cases            │  ___/70  │          │            │
│    Integration scenarios     │  ___/60  │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 5. Structure & ID Integrity  │  ___     │  100     │ ✅/⚠️/❌    │
│    IDs sequential/unique     │  ___/40  │          │            │
│    Classification correct    │  ___/30  │          │            │
│    Summary matches actual    │  ___/30  │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ TOTAL                        │  ___     │  1000    │            │
└──────────────────────────────┴──────────┴──────────┴────────────┘
```

---

## Deductions

| Location | Issue | Penalty |
|----------|-------|---------|
| <!-- TC ID or section --> | <!-- vague language / missing field / inconsistency --> | <!-- -N pts --> |

---

## Attack Points

### Attack 1: [dimension — specific weakness]

**Where**: <!-- quote from test-cases.md -->
**Why it's weak**: <!-- concrete critique -->
**What must improve**: <!-- actionable fix -->

### Attack 2: [dimension — specific weakness]

**Where**: <!-- quote from test-cases.md -->
**Why it's weak**: <!-- concrete critique -->
**What must improve**: <!-- actionable fix -->

### Attack 3: [dimension — specific weakness]

**Where**: <!-- quote from test-cases.md -->
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
- **Step Actionability**: {{SA_SCORE}}/250 {{if < 200: "⚠️ BLOCKING"}}
- **Action**: {{Continue to iteration N+1 / Target reached / Iterations exhausted}}
