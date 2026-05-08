---
date: "{{DATE}}"
doc_dir: "{{DOC_DIR}}"
iteration: "{{ITERATION}}"
target_score: "{{TARGET}}"
evaluator: Claude (automated, adversarial)
---

# UI Design Eval — Iteration {{ITERATION}}

**Score: {{SCORE}}/100** (target: {{TARGET}})

```
┌─────────────────────────────────────────────────────────────────┐
│                    UI DESIGN QUALITY SCORECARD                   │
├──────────────────────────────┬──────────┬──────────┬────────────┤
│ Dimension / Perspective      │ Score    │ Max      │ Status     │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 1. Requirement Coverage (PM) │  ___     │  25      │ ✅/⚠️/❌    │
│    UI function coverage      │  ___/8   │          │            │
│    Navigation Arch coverage  │  ___/4   │          │            │
│    State requirement coverage│  ___/8   │          │            │
│    Edge case handling        │  ___/5   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 2. User Experience (User)    │  ___     │  25      │ ✅/⚠️/❌    │
│    Information hierarchy     │  ___/8   │          │            │
│    Interaction intuitiveness │  ___/8   │          │            │
│    Accessibility             │  ___/9   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 3. Design Integrity (Design) │  ___     │  25      │ ✅/⚠️/❌    │
│    Design system adherence   │  ___/8   │          │            │
│    Visual coherence          │  ___/9   │          │            │
│    State completeness        │  ___/8   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 4. Implementability (Dev)    │  ___     │  25      │ ✅/⚠️/❌    │
│    Layout specificity        │  ___/8   │          │            │
│    Data binding explicit     │  ___/8   │          │            │
│    Interaction unambiguity   │  ___/9   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ TOTAL                        │  ___     │  100     │            │
└──────────────────────────────┴──────────┴──────────┴────────────┘
```

---

## Deductions

| Location | Issue | Penalty |
|----------|-------|---------|
| <!-- section:line --> | <!-- missing content / happy-path only / orphan element --> | <!-- -N pts --> |

---

## Attack Points

### Attack 1: [perspective — specific weakness]

**Where**: <!-- quote from document -->
**Why it's weak**: <!-- concrete critique from this perspective -->
**What must improve**: <!-- actionable fix -->

### Attack 2: [perspective — specific weakness]

**Where**: <!-- quote from document -->
**Why it's weak**: <!-- concrete critique from this perspective -->
**What must improve**: <!-- actionable fix -->

### Attack 3: [perspective — specific weakness]

**Where**: <!-- quote from document -->
**Why it's weak**: <!-- concrete critique from this perspective -->
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
