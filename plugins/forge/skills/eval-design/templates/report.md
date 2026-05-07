---
date: "{{DATE}}"
doc_dir: "{{DOC_DIR}}"
iteration: "{{ITERATION}}"
target_score: "{{TARGET}}"
evaluator: Claude (automated, adversarial)
---

# Design Eval — Iteration {{ITERATION}}

**Score: {{SCORE}}/100** (target: {{TARGET}})

```
┌─────────────────────────────────────────────────────────────────┐
│                     DESIGN QUALITY SCORECARD                     │
├──────────────────────────────┬──────────┬──────────┬────────────┤
│ Dimension                    │ Score    │ Max      │ Status     │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 1. Architecture Clarity      │  ___     │  20      │ ✅/⚠️/❌    │
│    Layer placement explicit  │  ___/7   │          │            │
│    Component diagram present │  ___/7   │          │            │
│    Dependencies listed       │  ___/6   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 2. Interface & Model Defs    │  ___     │  20      │ ✅/⚠️/❌    │
│    Interface signatures typed│  ___/7   │          │            │
│    Models concrete           │  ___/7   │          │            │
│    Directly implementable    │  ___/6   │          │            │
<!-- When HAS_DB_SCHEMA=true, replace lines 25-27 (the 3 sub-criteria above) with the 5 rows below. Keep the dimension header (line 24). Delete this comment afterward.
│    Interface signatures typed│  ___/5   │          │            │
│    Inline models concrete    │  ___/5   │          │            │
│    ER diagram complete       │  ___/3   │          │            │
│    SQL DDL directly usable   │  ___/4   │          │            │
│    Cross-layer consistency   │  ___/3   │          │            │
-->
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 3. Error Handling            │  ___     │  15      │ ✅/⚠️/❌    │
│    Error types defined       │  ___/5   │          │            │
│    Propagation strategy clear│  ___/5   │          │            │
│    HTTP status codes mapped  │  ___/5   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 4. Testing Strategy          │  ___     │  15      │ ✅/⚠️/❌    │
│    Per-layer test plan       │  ___/5   │          │            │
│    Coverage target numeric   │  ___/5   │          │            │
│    Test tooling named        │  ___/5   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 5. Breakdown-Readiness ★     │  ___     │  20      │ ✅/⚠️/❌    │
│    Components enumerable     │  ___/7   │          │            │
│    Tasks derivable           │  ___/7   │          │            │
│    PRD AC coverage           │  ___/6   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 6. Security Considerations   │  ___     │  10      │ ✅/⚠️/N/A  │
│    Threat model present      │  ___/5   │          │            │
│    Mitigations concrete      │  ___/5   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ TOTAL                        │  ___     │  100     │            │
└──────────────────────────────┴──────────┴──────────┴────────────┘
```

★ Breakdown-Readiness < 12/20 blocks progression to `/breakdown-tasks`

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
- **Breakdown-Readiness**: {{BR_SCORE}}/20 — {{can/cannot proceed to /breakdown-tasks}}
- **Action**: {{Continue to iteration N+1 / Target reached / Iterations exhausted}}
