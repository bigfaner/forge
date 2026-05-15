---
date: "{{DATE}}"
doc_dir: "{{DOC_DIR}}"
iteration: "{{ITERATION}}"
target_score: "{{TARGET}}"
evaluator: Claude (automated, adversarial)
---

# UI Design Eval — Iteration {{ITERATION}}

**Score: {{SCORE}}/1000** (target: {{TARGET}})

<!-- Platform-adaptive: replace D2/D3 dimension names with the active rubric's names:
     web → "User Experience" / "Design Integrity"
     mobile → "Touch Experience" / "Adaptive Layout"
     tui → "Terminal Experience" / "Visual Specification" -->

```
┌─────────────────────────────────────────────────────────────────┐
│                    UI DESIGN QUALITY SCORECARD                   │
├──────────────────────────────┬──────────┬──────────┬────────────┤
│ Dimension / Perspective      │ Score    │ Max      │ Status     │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 1. Requirement Coverage (PM) │  ___     │  250     │ ✅/⚠️/❌    │
│    UI function coverage      │  ___/80  │          │            │
│    Navigation Arch coverage  │  ___/40  │          │            │
│    State requirement coverage│  ___/80  │          │            │
│    Edge case handling        │  ___/50  │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 2. {{D2_NAME}} (End User)    │  ___     │  250     │ ✅/⚠️/❌    │
<!-- D2 sub-criteria: web (default) -->
│    Information hierarchy     │  ___/80  │          │            │
│    Interaction intuitiveness │  ___/80  │          │            │
│    Accessibility             │  ___/90  │          │            │
<!-- D2 sub-criteria: mobile (Touch Experience)
│    Touch target sizing       │  ___/70  │          │            │
│    Thumb reachability        │  ___/60  │          │            │
│    Gesture intuitiveness     │  ___/60  │          │            │
│    Platform convention       │  ___/60  │          │            │
-->
<!-- D2 sub-criteria: tui (Terminal Experience)
│    Keyboard reachability     │  ___/70  │          │            │
│    Info scanning efficiency  │  ___/60  │          │            │
│    Panel switching           │  ___/60  │          │            │
│    Mode switching consistency│  ___/60  │          │            │
-->
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 3. {{D3_NAME}} (Designer)    │  ___     │  250     │ ✅/⚠️/❌    │
<!-- D3 sub-criteria: web (default) -->
│    Design system adherence   │  ___/80  │          │            │
│    Visual coherence          │  ___/90  │          │            │
│    State completeness        │  ___/80  │          │            │
<!-- D3 sub-criteria: mobile (Adaptive Layout)
│    Multi-screen adaptation   │  ___/70  │          │            │
│    Orientation handling      │  ___/60  │          │            │
│    Safe area handling        │  ___/60  │          │            │
│    Platform-native feel      │  ___/60  │          │            │
-->
<!-- D3 sub-criteria: tui (Visual Specification)
│    ASCII mockup completeness │  ___/70  │          │            │
│    Character palette         │  ___/60  │          │            │
│    Color mapping compliance  │  ___/60  │          │            │
│    Dimension specificity     │  ___/60  │          │            │
-->
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 4. Implementability (Dev)    │  ___     │  250     │ ✅/⚠️/❌    │
│    Layout specificity        │  ___/80  │          │            │
│    Data binding explicit     │  ___/80  │          │            │
│    Interaction unambiguity   │  ___/90  │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ TOTAL                        │  ___     │  1000    │            │
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

- **Score**: {{SCORE}}/1000
- **Target**: {{TARGET}}/1000
- **Gap**: {{GAP}} points
- **Action**: {{Continue to iteration N+1 / Target reached / Iterations exhausted}}
