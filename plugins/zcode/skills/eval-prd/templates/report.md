---
date: "{{DATE}}"
prd_path: "{{PRD_PATH}}"
user_stories_path: "{{USER_STORIES_PATH}}"
iteration: "{{ITERATION}}"
target_score: "{{TARGET}}"
evaluator: Claude (automated, adversarial)
---

# PRD Eval — Iteration {{ITERATION}}

**Score: {{SCORE}}/100** (target: {{TARGET}})

```
┌─────────────────────────────────────────────────────────────────┐
│                       PRD QUALITY SCORECARD                      │
├──────────────────────────────┬──────────┬──────────┬────────────┤
│ Dimension                    │ Score    │ Max      │ Status     │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 1. Background & Goals        │  ___     │  20      │ ✅/⚠️/❌    │
│    Background three elements │  ___/7   │          │            │
│    Goals quantified          │  ___/7   │          │            │
│    Logical consistency       │  ___/6   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 2. Flow Diagrams             │  ___     │  20      │ ✅/⚠️/❌    │
│    Mermaid diagram exists    │  ___/7   │          │            │
│    Main path complete        │  ___/7   │          │            │
│    Decision + error branches │  ___/6   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 3. Functional Specs          │  ___     │  20      │ ✅/⚠️/❌    │
│    Tables complete           │  ___/7   │          │            │
│    Field descriptions clear  │  ___/7   │          │            │
│    Validation rules explicit │  ___/6   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 4. User Stories              │  ___     │  20      │ ✅/⚠️/❌    │
│    Coverage per user type    │  ___/7   │          │            │
│    Format correct            │  ___/7   │          │            │
│    AC per story              │  ___/6   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 5. Scope Clarity             │  ___     │  20      │ ✅/⚠️/❌    │
│    In-scope concrete         │  ___/7   │          │            │
│    Out-of-scope explicit     │  ___/7   │          │            │
│    Consistent with specs     │  ___/6   │          │            │
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

**Where**: <!-- quote from document -->
**Why it's weak**: <!-- concrete critique -->
**What must improve**: <!-- actionable fix -->

### Attack 2: [dimension — specific weakness]

**Where**: <!-- quote from document -->
**Why it's weak**: <!-- concrete critique -->
**What must improve**: <!-- actionable fix -->

### Attack 3: [dimension — specific weakness]

**Where**: <!-- quote from document -->
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
