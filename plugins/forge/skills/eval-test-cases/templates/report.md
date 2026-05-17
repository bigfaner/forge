---
date: "{{DATE}}"
doc_dir: "{{DOC_DIR}}"
iteration: "{{ITERATION}}"
target_score: "{{TARGET}}"
evaluator: Claude (automated, adversarial)
---

# Test Cases Eval — Iteration {{ITERATION}}

**Score: {{SCORE}}/100** (target: {{TARGET}})

```
┌─────────────────────────────────────────────────────────────────┐
│                  TEST CASES QUALITY SCORECARD                     │
├──────────────────────────────────────────────────────────────────┤
│ Dimension                    │ Score    │ Max      │ Status     │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 1. PRD Traceability          │  ___     │  25      │ ✅/⚠️/❌    │
│    TC-to-AC mapping          │  ___/9   │          │            │
│    Traceability table        │  ___/8   │          │            │
│    Reverse coverage          │  ___/8   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 2. Step Actionability        │  ___     │  25      │ ✅/⚠️/❌    │
│    Steps concrete            │  ___/9   │          │            │
│    Expected results          │  ___/9   │          │            │
│    Preconditions explicit    │  ___/7   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 3. Route & Element Accuracy  │  ___     │  20      │ ✅/⚠️/❌    │
│    Routes valid              │  ___/7   │          │            │
│    Elements identifiable     │  ___/7   │          │            │
│    Consistency               │  ___/6   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 4. Completeness              │  ___     │  20      │ ✅/⚠️/❌    │
│    Type coverage             │  ___/7   │          │            │
│    Boundary cases            │  ___/7   │          │            │
│    Integration scenarios     │  ___/6   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 5. Structure & ID Integrity  │  ___     │  10      │ ✅/⚠️/❌    │
│    IDs sequential/unique     │  ___/4   │          │            │
│    Classification correct    │  ___/3   │          │            │
│    Summary matches actual    │  ___/3   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ TOTAL                        │  ___     │  100     │            │
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

- **Score**: {{SCORE}}/100
- **Target**: {{TARGET}}/100
- **Gap**: {{GAP}} points
- **Step Actionability**: {{SA_SCORE}}/25 {{if < 20: "⚠️ BLOCKING"}}
- **Action**: {{Continue to iteration N+1 / Target reached / Iterations exhausted}}
