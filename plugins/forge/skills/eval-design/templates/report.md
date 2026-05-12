---
date: "{{DATE}}"
doc_dir: "{{DOC_DIR}}"
iteration: "{{ITERATION}}"
target_score: "{{TARGET}}"
evaluator: Claude (automated, adversarial)
---

# Design Eval — Iteration {{ITERATION}}

**Score: {{SCORE}}/1000** (target: {{TARGET}})

```
┌─────────────────────────────────────────────────────────────────┐
│                     DESIGN QUALITY SCORECARD                     │
├──────────────────────────────┬──────────┬──────────┬────────────┤
│ Dimension                    │ Score    │ Max      │ Status     │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 1. Architecture Clarity      │  ___     │  200     │ ✅/⚠️/❌    │
│    Layer placement explicit  │  ___/70  │          │            │
│    Component diagram present │  ___/70  │          │            │
│    Dependencies listed       │  ___/60  │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 2. Interface & Model Defs    │  ___     │  200     │ ✅/⚠️/❌    │
│    Interface signatures typed│  ___/70  │          │            │
│    Models concrete           │  ___/70  │          │            │
│    Directly implementable    │  ___/60  │          │            │
<!-- When HAS_DB_SCHEMA=true, replace lines 25-27 (the 3 sub-criteria above) with the 5 rows below. Keep the dimension header (line 24). Delete this comment afterward.
│    Interface signatures typed│  ___/50  │          │            │
│    Inline models concrete    │  ___/50  │          │            │
│    ER diagram complete       │  ___/30  │          │            │
│    SQL DDL directly usable   │  ___/40  │          │            │
│    Cross-layer consistency   │  ___/30  │          │            │
-->
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 3. Error Handling            │  ___     │  150     │ ✅/⚠️/❌    │
│    Error types defined       │  ___/50  │          │            │
│    Propagation strategy clear│  ___/50  │          │            │
│    HTTP status codes mapped  │  ___/50  │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 4. Testing Strategy          │  ___     │  150     │ ✅/⚠️/❌    │
│    Per-layer test plan       │  ___/50  │          │            │
│    Coverage target numeric   │  ___/50  │          │            │
│    Test tooling named        │  ___/50  │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 5. Breakdown-Readiness ★     │  ___     │  200     │ ✅/⚠️/❌    │
│    Components enumerable     │  ___/70  │          │            │
│    Tasks derivable           │  ___/70  │          │            │
│    PRD AC coverage           │  ___/60  │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 6. Security Considerations   │  ___     │  100     │ ✅/⚠️/N/A  │
│    Threat model present      │  ___/50  │          │            │
│    Mitigations concrete      │  ___/50  │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ TOTAL                        │  ___     │  1000    │            │
└──────────────────────────────┴──────────┴──────────┴────────────┘
```

★ Breakdown-Readiness < 180/200 blocks progression to `/breakdown-tasks`

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

- **Score**: {{SCORE}}/1000
- **Target**: {{TARGET}}/1000
- **Gap**: {{GAP}} points
- **Breakdown-Readiness**: {{BR_SCORE}}/200 — {{can/cannot proceed to /breakdown-tasks}}
- **Action**: {{Continue to iteration N+1 / Target reached / Iterations exhausted}}
